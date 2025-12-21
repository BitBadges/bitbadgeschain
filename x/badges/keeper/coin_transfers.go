package keeper

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	RoyaltyDivisor = 10000
)

// formatDenomForDisplay formats a denom for display in error messages
// Shows "BADGE" for "ubadge" and prints others as-is
func formatDenomForDisplay(denom string) string {
	if denom == "ubadge" {
		return "BADGE"
	}
	return denom
}

// SimulateCoinTransfers simulates coin transfers for approval validation
// Returns (deterministicErrorMsg, error) where deterministicErrorMsg is a deterministic error string
func (k Keeper) SimulateCoinTransfers(
	ctx sdk.Context,
	coinTransfers []*types.CoinTransfer,
	transferMetadata TransferMetadata,
	collection *types.TokenCollection,
	royalties *types.UserRoyalties,
) (string, error) {
	initiatedBy := transferMetadata.InitiatedBy
	approverAddress := transferMetadata.ApproverAddress
	approvalLevel := transferMetadata.ApprovalLevel
	if len(coinTransfers) == 0 {
		return "", nil
	}

	if royalties == nil {
		detErrMsg := "royalties is nil"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	if collection == nil {
		detErrMsg := "collection is nil"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	allowedDenoms := k.GetAllowedDenoms(ctx)
	if len(allowedDenoms) == 0 {
		detErrMsg := "allowed denoms is empty"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	royaltyPercentage := royalties.Percentage
	royaltyPayoutAddress := royalties.PayoutAddress
	if royaltyPercentage.GT(sdkmath.NewUint(0)) {
		if royaltyPayoutAddress == "" {
			detErrMsg := "payout address is required when royalty percentage is greater than 0"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidAddress, detErrMsg)
		}

		_, err := sdk.AccAddressFromBech32(royaltyPayoutAddress)
		if err != nil {
			return "", sdkerrors.Wrapf(err, "invalid payout address %s", royaltyPayoutAddress)
		}
	}

	for _, coinTransfer := range coinTransfers {
		if coinTransfer == nil {
			detErrMsg := "coin transfer is nil"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
		if len(coinTransfer.Coins) == 0 {
			detErrMsg := "coin transfer cannot have empty coins slice"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}
		if !slices.Contains(allowedDenoms, coinTransfer.Coins[0].Denom) {
			detErrMsg := fmt.Sprintf("denom %s is not allowed", coinTransfer.Coins[0].Denom)
			return detErrMsg, sdkerrors.Wrap(ErrInvalidDenom, detErrMsg)
		}
	}

	spendableCoinsMap := make(map[string]sdk.Coins)

	for _, coinTransfer := range coinTransfers {
		toTransfer := coinTransfer.Coins
		fromAddress := initiatedBy
		if coinTransfer.OverrideFromWithApproverAddress {
			// collection-level
			if approverAddress == "" && approvalLevel == "collection" {
				approverAddress = collection.MintEscrowAddress
			}

			fromAddress = approverAddress
			if fromAddress == "" {
				detErrMsg := "approver address is required when overrideFromWithApproverAddress is true"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidAddress, detErrMsg)
			}
		}

		for _, coin := range toTransfer {
			// Only badgeslp: denoms use the wrapped approach (check badge balances)
			// All other denoms (badges: and non-badges) use standard bank logic
			parts := strings.Split(coin.Denom, ":")
			if len(parts) >= 2 && parts[0] == "badgeslp" {
				// badgeslp: denom - check badge balances (wrapped approach)
				fromAddressAcc := sdk.MustAccAddressFromBech32(fromAddress)
				balanceCoin, err := k.sendManagerKeeper.GetBalanceWithAliasRouting(ctx, fromAddressAcc, coin.Denom)
				if err != nil {
					return "", sdkerrors.Wrapf(err, "error checking wrapped denom balance for %s", coin.Denom)
				}

				requiredAmount := coin.Amount
				if balanceCoin.Amount.LT(requiredAmount) {
					detErrMsg := fmt.Sprintf("insufficient %s balance to complete transfer", formatDenomForDisplay(coin.Denom))
					return detErrMsg, sdkerrors.Wrap(types.ErrUnderflow, detErrMsg)
				}
			} else {
				// For all other denoms (badges: and non-badges), use standard bank coin logic
				// Lazy-load spendable coins only when needed
				if _, exists := spendableCoinsMap[fromAddress]; !exists {
					spendableCoinsMap[fromAddress] = k.bankKeeper.SpendableCoins(ctx, sdk.MustAccAddressFromBech32(fromAddress))
				}

				newCoins, underflow := spendableCoinsMap[fromAddress].SafeSub(*coin)
				if underflow {
					detErrMsg := fmt.Sprintf("insufficient %s balance to complete transfer", formatDenomForDisplay(coin.Denom))
					return detErrMsg, sdkerrors.Wrap(types.ErrUnderflow, detErrMsg)
				}
				spendableCoinsMap[fromAddress] = newCoins
			}
		}
	}

	return "", nil
}

// sendCoinWithRoyalty handles sending a coin with royalty deduction
// It sends the royalty to the payout address and the remaining amount to the recipient
func (k Keeper) sendCoinWithRoyalty(
	ctx sdk.Context,
	coin *sdk.Coin,
	royaltyAmountInt sdkmath.Int,
	remainingAmount sdkmath.Int,
	fromAddressAcc sdk.AccAddress,
	toAddressAcc sdk.AccAddress,
	royaltyPayoutAddress string,
	coinTransfersUsed *[]CoinTransfers,
) error {
	// Send royalty to payout address
	if !royaltyAmountInt.IsZero() {
		payoutAddressAcc := sdk.MustAccAddressFromBech32(royaltyPayoutAddress)
		royaltyCoin := sdk.NewCoin(coin.Denom, royaltyAmountInt)

		err := k.sendManagerKeeper.SendCoinWithAliasRouting(ctx, fromAddressAcc, payoutAddressAcc, &royaltyCoin)
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

		err := k.sendManagerKeeper.SendCoinWithAliasRouting(ctx, fromAddressAcc, toAddressAcc, &remainingCoin)
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

	return nil
}

func (k Keeper) ExecuteCoinTransfers(
	ctx sdk.Context,
	coinTransfers []*types.CoinTransfer,
	transferMetadata TransferMetadata,
	coinTransfersUsed *[]CoinTransfers,
	collection *types.TokenCollection,
	royalties *types.UserRoyalties,
) error {
	initiatedBy := transferMetadata.InitiatedBy
	approverAddress := transferMetadata.ApproverAddress
	approvalLevel := transferMetadata.ApprovalLevel
	if len(coinTransfers) == 0 {
		return nil
	}

	if royalties == nil {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "royalties is nil")
	}

	if collection == nil {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "collection is nil")
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

	allowedDenoms := k.GetAllowedDenoms(ctx)
	if len(allowedDenoms) == 0 {
		return sdkerrors.Wrapf(types.ErrInvalidRequest, "allowed denoms is empty")
	}

	for _, coinTransfer := range coinTransfers {
		if coinTransfer == nil {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin transfer is nil")
		}
		if len(coinTransfer.Coins) == 0 {
			return sdkerrors.Wrapf(types.ErrInvalidRequest, "coin transfer cannot have empty coins slice")
		}
		if !slices.Contains(allowedDenoms, coinTransfer.Coins[0].Denom) {
			return sdkerrors.Wrapf(ErrInvalidDenom, "denom %s is not allowed", coinTransfer.Coins[0].Denom)
		}
	}

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

			if approverAddress == "" {
				return sdkerrors.Wrapf(types.ErrInvalidAddress, "approver address is required when overrideFromWithApproverAddress is true")
			}

			fromAddressAcc = sdk.MustAccAddressFromBech32(approverAddress)
		}

		for _, coin := range coinsToTransfer {
			coinAmountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())
			royaltyAmountUint := coinAmountUint.Mul(royaltyPercentage).Quo(sdkmath.NewUint(RoyaltyDivisor))
			royaltyAmountInt := sdkmath.NewIntFromBigInt(royaltyAmountUint.BigInt())
			remainingAmount := coin.Amount.Sub(royaltyAmountInt)

			err := k.sendCoinWithRoyalty(
				ctx,
				coin,
				royaltyAmountInt,
				remainingAmount,
				fromAddressAcc,
				toAddressAcc,
				royaltyPayoutAddress,
				coinTransfersUsed,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
