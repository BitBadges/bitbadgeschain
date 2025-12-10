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

// GetSpendableCoinAmountWithWrapping gets the spendable amount for a specific denom, handling wrapped badgeslp denoms
// For badgeslp: denoms, calculates from badge balances (wrapped approach)
// For all other denoms (including badges:), uses bank module directly
func (k Keeper) GetSpendableCoinAmountWithWrapping(ctx sdk.Context, address sdk.AccAddress, denom string) (sdkmath.Int, error) {
	// Only badgeslp: denoms use the wrapped approach
	parts := strings.Split(denom, ":")
	if len(parts) >= 2 && parts[0] == "badgeslp" {
		// badgeslp: denom - calculate from badge balances
		// Parse collection from denom
		collection, err := k.ParseCollectionFromDenom(ctx, denom)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}

		// Get the corresponding wrapper path
		path, err := GetCorrespondingPath(collection, denom)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}

		// Get user's badge balance
		userBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, address.String())

		// Use the same calculation as GetWrappableBalances
		maxWrappableAmount, err := k.calculateMaxWrappableAmount(ctx, userBalances.Balances, path.Balances)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}

		return sdkmath.NewIntFromBigInt(maxWrappableAmount.BigInt()), nil
	}

	// For all other denoms (badges: and non-badges), use bank keeper
	spendableCoins := k.bankKeeper.SpendableCoins(ctx, address)
	coin := spendableCoins.AmountOf(denom)
	return coin, nil
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
				availableAmount, err := k.GetSpendableCoinAmountWithWrapping(ctx, fromAddressAcc, coin.Denom)
				if err != nil {
					return "", sdkerrors.Wrapf(err, "error checking wrapped denom balance for %s", coin.Denom)
				}

				requiredAmount := coin.Amount
				if availableAmount.LT(requiredAmount) {
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
	useWrappedApproach bool,
) error {
	// Send royalty to payout address
	if !royaltyAmountInt.IsZero() {
		payoutAddressAcc := sdk.MustAccAddressFromBech32(royaltyPayoutAddress)
		var err error

		if useWrappedApproach {
			// Use badge transfer functions for badgeslp: denoms
			royaltyAmountUint := sdkmath.NewUintFromBigInt(royaltyAmountInt.BigInt())
			err = k.SendNativeTokensToAddress(ctx, fromAddressAcc.String(), payoutAddressAcc.String(), coin.Denom, royaltyAmountUint)
			if err != nil {
				return sdkerrors.Wrapf(err, "error sending royalty to payout address for badgeslp denom")
			}
		} else {
			// Use bank keeper for all other denoms
			royaltyCoin := sdk.NewCoin(coin.Denom, royaltyAmountInt)
			err = k.bankKeeper.SendCoins(ctx, fromAddressAcc, payoutAddressAcc, sdk.NewCoins(royaltyCoin))
			if err != nil {
				return sdkerrors.Wrapf(err, "error sending royalty to payout address")
			}
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
		var err error

		if useWrappedApproach {
			// Use badge transfer functions for badgeslp: denoms
			remainingAmountUint := sdkmath.NewUintFromBigInt(remainingAmount.BigInt())
			err = k.SendNativeTokensToAddress(ctx, fromAddressAcc.String(), toAddressAcc.String(), coin.Denom, remainingAmountUint)
			if err != nil {
				return sdkerrors.Wrapf(err, "error sending remaining amount to recipient for badgeslp denom")
			}
		} else {
			// Use bank keeper for all other denoms
			remainingCoin := sdk.NewCoin(coin.Denom, remainingAmount)
			err = k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(remainingCoin))
			if err != nil {
				return sdkerrors.Wrapf(err, "error sending remaining amount to recipient")
			}
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

			fromAddressAcc = sdk.MustAccAddressFromBech32(approverAddress)
		}

		for _, coin := range coinsToTransfer {
			// Calculate royalty amount
			coinAmountUint := sdkmath.NewUintFromBigInt(coin.Amount.BigInt())

			// Validate royalty divisor to prevent division by zero
			if RoyaltyDivisor == 0 {
				return sdkerrors.Wrapf(types.ErrInvalidRequest, "royalty divisor cannot be zero")
			}

			royaltyAmountUint := coinAmountUint.Mul(royaltyPercentage).Quo(sdkmath.NewUint(RoyaltyDivisor))

			// Convert royalty amount to Int for coin operations
			royaltyAmountInt := sdkmath.NewIntFromBigInt(royaltyAmountUint.BigInt())

			// Calculate remaining amount after royalty
			remainingAmount := coin.Amount.Sub(royaltyAmountInt)

			// Only badgeslp: denoms use the wrapped approach (badge transfer functions)
			// All other denoms (badges: and non-badges) use bank.SendCoins
			parts := strings.Split(coin.Denom, ":")
			useWrappedApproach := len(parts) >= 2 && parts[0] == "badgeslp"

			err := k.sendCoinWithRoyalty(
				ctx,
				coin,
				royaltyAmountInt,
				remainingAmount,
				fromAddressAcc,
				toAddressAcc,
				royaltyPayoutAddress,
				coinTransfersUsed,
				useWrappedApproach,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
