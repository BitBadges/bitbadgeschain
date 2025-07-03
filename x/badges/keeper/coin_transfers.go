package keeper

import (
	"slices"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// formatDenomForDisplay formats a denom for display in error messages
// Shows "BADGE" for "ubadge" and prints others as-is
func formatDenomForDisplay(denom string) string {
	if denom == "ubadge" {
		return "BADGE"
	}
	return denom
}

func (k Keeper) HandleCoinTransfers(
	ctx sdk.Context,
	coinTransfers []*types.CoinTransfer,
	initiatedBy string,
	approverAddress string,
	approvalLevel string,
	simulate bool,
	coinTransfersUsed *[]CoinTransfers,
	collection *types.BadgeCollection,
	royalties *types.UserRoyalties,
) error {
	if len(coinTransfers) == 0 {
		return nil
	}

	royaltyPercentage := royalties.Percentage
	royaltyPayoutAddress := royalties.PayoutAddress
	if royaltyPercentage.GT(sdkmath.NewUint(0)) {
		if royaltyPayoutAddress == "" {
			return sdkerrors.Wrap(types.ErrInvalidAddress, "payout address is required when royalty percentage is greater than 0")
		}

		_, err := sdk.AccAddressFromBech32(royaltyPayoutAddress)
		if err != nil {
			return sdkerrors.Wrapf(err, "invalid payout address %s", royaltyPayoutAddress)
		}
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
				// collection-level
				if approverAddress == "" && approvalLevel == "collection" {
					approverAddress = collection.MintEscrowAddress
				}

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
				// collection-level
				if approverAddress == "" && approvalLevel == "collection" {
					approverAddress = collection.MintEscrowAddress
				}

				fromAddress = approverAddress
			}

			for _, coin := range toTransfer {
				newCoins, underflow := spendableCoinsMap[fromAddress].SafeSub(*coin)
				if underflow {
					return sdkerrors.Wrapf(types.ErrUnderflow, "insufficient %s balance to complete transfer", formatDenomForDisplay(coin.Denom))
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
				// collection-level
				if approverAddress == "" && approvalLevel == "collection" {
					approverAddress = collection.MintEscrowAddress
				}

				fromAddressAcc = sdk.MustAccAddressFromBech32(approverAddress)
			}

			for _, coin := range coinsToTransfer {
				// Calculate royalty amount
				coinAmountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
				royaltyAmountUint := coinAmountUint.Mul(royaltyPercentage).Quo(sdkmath.NewUint(10000)) // Assuming royalty percentage is in basis points (e.g., 250 = 2.5%)

				// Convert royalty amount to Int for coin operations
				royaltyAmountInt := sdkmath.NewIntFromBigInt(royaltyAmountUint.BigInt())

				// Calculate remaining amount after royalty
				remainingAmount := coin.Amount.Sub(royaltyAmountInt)

				// Send royalty to mint escrow address
				if !royaltyAmountInt.IsZero() {
					royaltyCoin := sdk.NewCoin(coin.Denom, royaltyAmountInt)
					payoutAddressAcc := sdk.MustAccAddressFromBech32(royaltyPayoutAddress)
					err := k.bankKeeper.SendCoins(ctx, fromAddressAcc, payoutAddressAcc, sdk.NewCoins(royaltyCoin))
					if err != nil {
						return sdkerrors.Wrapf(err, "error sending royalty to payout address")
					}

					*coinTransfersUsed = append(*coinTransfersUsed, CoinTransfers{
						From:   fromAddressAcc.String(),
						To:     payoutAddressAcc.String(),
						Amount: royaltyAmountInt.String(),
						Denom:  coin.Denom,
					})
				}

				// Send remaining amount to recipient
				if !remainingAmount.IsZero() {
					remainingCoin := sdk.NewCoin(coin.Denom, remainingAmount)
					err := k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(remainingCoin))
					if err != nil {
						return sdkerrors.Wrapf(err, "error sending remaining amount to recipient")
					}

					*coinTransfersUsed = append(*coinTransfersUsed, CoinTransfers{
						From:   fromAddressAcc.String(),
						To:     toAddressAcc.String(),
						Amount: remainingAmount.String(),
						Denom:  coin.Denom,
					})
				}
			}
		}
	}

	return nil
}
