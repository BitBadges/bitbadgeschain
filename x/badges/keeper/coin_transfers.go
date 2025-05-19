package keeper

import (
	"slices"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleCoinTransfers(ctx sdk.Context, coinTransfers []*types.CoinTransfer, initiatedBy string, approverAddress string, simulate bool) error {
	if len(coinTransfers) == 0 {
		return nil
	}

	if !k.EnableCoinTransfers {
		return sdkerrors.Wrap(ErrCoinTransfersDisabled, "coin transfers are disabled")
	}

	if len(k.AllowedDenoms) > 0 {
		for _, coinTransfer := range coinTransfers {
			if !slices.Contains(k.AllowedDenoms, coinTransfer.Coins[0].Denom) {
				return sdkerrors.Wrapf(ErrInvalidDenom, "denom %s is not allowed", coinTransfer.Coins[0].Denom)
			}
		}
	}

	if simulate {
		spendableCoinsMap := make(map[string]sdk.Coins)
		for _, coinTransfer := range coinTransfers {
			fromAddress := initiatedBy
			if coinTransfer.OverrideFromWithApproverAddress {
				fromAddress = approverAddress
				if fromAddress == "" {
					return sdkerrors.Wrap(types.ErrInvalidAddress, "approver address is required when overrideFromWithApproverAddress is true")
				}
			}

			spendableCoinsMap[fromAddress] = k.bankKeeper.SpendableCoins(ctx, sdk.MustAccAddressFromBech32(fromAddress))
		}

		for _, coinTransfer := range coinTransfers {
			toTransfer := coinTransfer.Coins
			fromAddress := initiatedBy
			if coinTransfer.OverrideFromWithApproverAddress {
				fromAddress = approverAddress
			}

			for _, coin := range toTransfer {
				newCoins, underflow := spendableCoinsMap[fromAddress].SafeSub(*coin)
				if underflow {
					return sdkerrors.Wrapf(types.ErrUnderflow, "insufficient $BADGE balance to complete transfer")
				}
				spendableCoinsMap[fromAddress] = newCoins
			}
		}
	} else {
		for _, coinTransfer := range coinTransfers {
			coinsToTransfer := coinTransfer.Coins

			to := coinTransfer.To
			if coinTransfer.OverrideToWithInitiator {
				to = initiatedBy
			}

			toAddressAcc := sdk.MustAccAddressFromBech32(to)
			fromAddressAcc := sdk.MustAccAddressFromBech32(initiatedBy)
			if coinTransfer.OverrideFromWithApproverAddress {
				fromAddressAcc = sdk.MustAccAddressFromBech32(approverAddress)
			}

			for _, coin := range coinsToTransfer {
				err := k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(*coin))
				if err != nil {
					return sdkerrors.Wrapf(err, "error sending $BADGE, passed simulation but not actual transfers")
				}
			}
		}
	}

	return nil
}
