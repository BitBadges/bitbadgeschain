package keeper

import (
	"slices"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleCoinTransfers(ctx sdk.Context, coinTransfers []*types.CoinTransfer, initiatedBy string, simulate bool) error {
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

	//simulate the sdk.Coin transfers
	initiatedByAcc := sdk.MustAccAddressFromBech32(initiatedBy)

	if simulate {
		spendableCoins := k.bankKeeper.SpendableCoins(ctx, initiatedByAcc)
		for _, coinTransfer := range coinTransfers {
			toTransfer := coinTransfer.Coins
			for _, coin := range toTransfer {
				newCoins, underflow := spendableCoins.SafeSub(*coin)
				if underflow {
					return sdkerrors.Wrapf(types.ErrUnderflow, "insufficient $BADGE balance to complete transfer")
				}
				spendableCoins = newCoins
			}
		}
	} else {
		for _, coinTransfer := range coinTransfers {
			coinsToTransfer := coinTransfer.Coins
			toAddressAcc := sdk.MustAccAddressFromBech32(coinTransfer.To)
			fromAddressAcc := initiatedByAcc
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
