package keeper

import (
	"fmt"
	"slices"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DynamicCoinTransferOptions holds the new optional parameters for
// enhanced coin transfer execution (dynamic fees, metered balances,
// time-based refunds). All fields are optional — nil means the feature
// is not used and the transfer falls back to existing fixed-amount behavior.
type DynamicCoinTransferOptions struct {
	// DynamicFeeSchedule computes percentage-based fees scaled to token transfer amount.
	DynamicFeeSchedule *types.DynamicFeeSchedule
	// CoinPerTokenMultiplier scales coin amounts by the token transfer amount.
	CoinPerTokenMultiplier *types.CoinPerTokenMultiplier
	// TimeBasedRefundFormula computes pro-rated refund amounts.
	TimeBasedRefundFormula *types.TimeBasedRefundFormula
	// TokenTransferBalances is the actual token balances being transferred.
	// Required when DynamicFeeSchedule or CoinPerTokenMultiplier is set.
	TokenTransferBalances []*types.Balance
}

// ExecuteCoinTransfersWithDynamicOptions extends the existing ExecuteCoinTransfers
// with support for dynamic fees, metered balances, and time-based refunds.
//
// When dynamicOpts is nil, this behaves identically to ExecuteCoinTransfers.
func (k Keeper) ExecuteCoinTransfersWithDynamicOptions(
	ctx sdk.Context,
	coinTransfers []*types.CoinTransfer,
	transferMetadata TransferMetadata,
	coinTransfersUsed *[]CoinTransfers,
	collection *types.TokenCollection,
	royalties *types.UserRoyalties,
	dynamicOpts *DynamicCoinTransferOptions,
) (string, error) {
	// If no dynamic options, delegate to existing logic
	if dynamicOpts == nil {
		return k.ExecuteCoinTransfers(ctx, coinTransfers, transferMetadata, coinTransfersUsed, collection, royalties)
	}

	initiatedBy := transferMetadata.InitiatedBy
	approverAddress := transferMetadata.ApproverAddress
	approvalLevel := transferMetadata.ApprovalLevel

	if royalties == nil {
		detErrMsg := "royalties is nil"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}
	if collection == nil {
		detErrMsg := "collection is nil"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	royaltyPercentage := royalties.Percentage
	royaltyPayoutAddress := royalties.PayoutAddress

	if royaltyPercentage.GT(sdkmath.NewUint(RoyaltyDivisor)) {
		detErrMsg := fmt.Sprintf("royalty percentage cannot exceed %d (100%%), got %s", RoyaltyDivisor, royaltyPercentage.String())
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

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

	allowedDenoms := k.GetAllowedDenoms(ctx)
	if len(allowedDenoms) == 0 {
		detErrMsg := "allowed denoms is empty"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	// --- Feature B: Dynamic Fee Schedule ---
	// Execute dynamic fee transfers before fixed coin transfers.
	if dynamicOpts.DynamicFeeSchedule != nil {
		detErrMsg, err := k.executeDynamicFee(ctx, dynamicOpts, transferMetadata, coinTransfersUsed, collection, allowedDenoms)
		if err != nil {
			return detErrMsg, err
		}
	}

	// --- Feature C: Time-Based Refund ---
	// Execute time-based refund if configured.
	if dynamicOpts.TimeBasedRefundFormula != nil {
		detErrMsg, err := k.executeTimeBasedRefund(ctx, dynamicOpts, transferMetadata, coinTransfersUsed, collection, allowedDenoms)
		if err != nil {
			return detErrMsg, err
		}
	}

	// --- Feature A + standard coin transfers ---
	// Process each coin transfer, possibly with metered amounts.
	for _, coinTransfer := range coinTransfers {
		if coinTransfer == nil {
			detErrMsg := "coin transfer is nil"
			return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Determine coins to transfer: metered or fixed
		var coinsToTransfer []*sdk.Coin
		if dynamicOpts.CoinPerTokenMultiplier != nil && len(coinTransfer.Coins) == 0 {
			// Feature A: Metered balances — compute coin amounts from token transfer amount
			tokenAmount := types.GetTotalTokenTransferAmount(dynamicOpts.TokenTransferBalances)
			if tokenAmount.IsZero() {
				detErrMsg := "token transfer amount is zero but CoinPerTokenMultiplier is set"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
			meteredCoins := dynamicOpts.CoinPerTokenMultiplier.ComputeMeteredAmount(tokenAmount)
			for i := range meteredCoins {
				coinsToTransfer = append(coinsToTransfer, &meteredCoins[i])
			}
		} else {
			if len(coinTransfer.Coins) == 0 {
				detErrMsg := "coin transfer cannot have empty coins slice"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
			coinsToTransfer = coinTransfer.Coins
		}

		// Validate denoms
		for _, coin := range coinsToTransfer {
			if !slices.Contains(allowedDenoms, coin.Denom) {
				detErrMsg := fmt.Sprintf("denom %s is not allowed", coin.Denom)
				return detErrMsg, sdkerrors.Wrap(ErrInvalidDenom, detErrMsg)
			}
		}

		to := coinTransfer.To
		if coinTransfer.OverrideToWithInitiator {
			to = initiatedBy
		}

		toAddressAcc, err := sdk.AccAddressFromBech32(to)
		if err != nil {
			detErrMsg := fmt.Sprintf("invalid to address: %s", to)
			return detErrMsg, sdkerrors.Wrapf(err, "invalid to address: %s", to)
		}
		fromAddressAcc, err := sdk.AccAddressFromBech32(initiatedBy)
		if err != nil {
			detErrMsg := fmt.Sprintf("invalid initiatedBy address: %s", initiatedBy)
			return detErrMsg, sdkerrors.Wrapf(err, "invalid initiatedBy address: %s", initiatedBy)
		}
		if coinTransfer.OverrideFromWithApproverAddress {
			if approverAddress == "" && approvalLevel == "collection" {
				approverAddress = collection.MintEscrowAddress
			}
			if approverAddress == "" {
				detErrMsg := "approver address is required when overrideFromWithApproverAddress is true"
				return detErrMsg, sdkerrors.Wrap(types.ErrInvalidAddress, detErrMsg)
			}
			fromAddressAcc, err = sdk.AccAddressFromBech32(approverAddress)
			if err != nil {
				detErrMsg := fmt.Sprintf("invalid approver address: %s", approverAddress)
				return detErrMsg, sdkerrors.Wrapf(err, "invalid approver address: %s", approverAddress)
			}
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
				detErrMsg := fmt.Sprintf("insufficient %s balance to complete transfer", formatDenomForDisplay(coin.Denom))
				return detErrMsg, sdkerrors.Wrap(types.ErrUnderflow, err.Error())
			}
		}
	}

	return "", nil
}

// executeDynamicFee computes and sends a percentage-based fee scaled to
// the token transfer amount.
func (k Keeper) executeDynamicFee(
	ctx sdk.Context,
	dynamicOpts *DynamicCoinTransferOptions,
	transferMetadata TransferMetadata,
	coinTransfersUsed *[]CoinTransfers,
	collection *types.TokenCollection,
	allowedDenoms []string,
) (string, error) {
	schedule := dynamicOpts.DynamicFeeSchedule

	// Validate fee denom
	if !slices.Contains(allowedDenoms, schedule.FeeDenom) {
		detErrMsg := fmt.Sprintf("dynamic fee denom %s is not allowed", schedule.FeeDenom)
		return detErrMsg, sdkerrors.Wrap(ErrInvalidDenom, detErrMsg)
	}

	// Validate fee recipient
	feeRecipientAcc, err := sdk.AccAddressFromBech32(schedule.FeeRecipient)
	if err != nil {
		detErrMsg := fmt.Sprintf("invalid dynamic fee recipient: %s", schedule.FeeRecipient)
		return detErrMsg, sdkerrors.Wrapf(err, "invalid dynamic fee recipient")
	}

	// Compute total token transfer amount
	tokenAmount := types.GetTotalTokenTransferAmount(dynamicOpts.TokenTransferBalances)
	if tokenAmount.IsZero() {
		// No tokens being transferred, no fee
		return "", nil
	}

	// Compute fee from schedule
	feeAmount, found := schedule.ComputeDynamicFee(tokenAmount)
	if !found {
		detErrMsg := fmt.Sprintf("no matching fee tier for token transfer amount %s", tokenAmount.String())
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	if feeAmount.IsZero() {
		return "", nil
	}

	// Send fee from initiator to fee recipient
	fromAddressAcc, err := sdk.AccAddressFromBech32(transferMetadata.InitiatedBy)
	if err != nil {
		detErrMsg := fmt.Sprintf("invalid initiatedBy address: %s", transferMetadata.InitiatedBy)
		return detErrMsg, sdkerrors.Wrapf(err, "invalid initiatedBy address")
	}

	feeCoin := sdk.NewCoin(schedule.FeeDenom, feeAmount)
	err = k.sendManagerKeeper.SendCoinWithAliasRouting(ctx, fromAddressAcc, feeRecipientAcc, &feeCoin)
	if err != nil {
		detErrMsg := fmt.Sprintf("insufficient %s balance for dynamic fee", formatDenomForDisplay(schedule.FeeDenom))
		return detErrMsg, sdkerrors.Wrap(types.ErrUnderflow, err.Error())
	}

	*coinTransfersUsed = append(*coinTransfersUsed, CoinTransfers{
		From:   fromAddressAcc.String(),
		To:     feeRecipientAcc.String(),
		Amount: feeAmount.String(),
		Denom:  schedule.FeeDenom,
	})

	return "", nil
}

// executeTimeBasedRefund computes and sends a pro-rated refund based on
// remaining subscription time.
func (k Keeper) executeTimeBasedRefund(
	ctx sdk.Context,
	dynamicOpts *DynamicCoinTransferOptions,
	transferMetadata TransferMetadata,
	coinTransfersUsed *[]CoinTransfers,
	collection *types.TokenCollection,
	allowedDenoms []string,
) (string, error) {
	formula := dynamicOpts.TimeBasedRefundFormula

	if formula.TotalDuration.IsZero() {
		detErrMsg := "time-based refund formula has zero total duration"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
	}

	// Validate all refund coin denoms
	for _, coin := range formula.BaseRefundAmount {
		if !slices.Contains(allowedDenoms, coin.Denom) {
			detErrMsg := fmt.Sprintf("refund denom %s is not allowed", coin.Denom)
			return detErrMsg, sdkerrors.Wrap(ErrInvalidDenom, detErrMsg)
		}
	}

	// Compute refund amounts based on current time
	currentTimeMs := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	refundCoins := formula.ComputeRefundAmount(currentTimeMs)

	// The refund is sent from the approver (escrow) to the initiator.
	// This requires OverrideFromWithApproverAddress semantics.
	refundRecipientAcc, err := sdk.AccAddressFromBech32(transferMetadata.InitiatedBy)
	if err != nil {
		detErrMsg := fmt.Sprintf("invalid initiatedBy address for refund: %s", transferMetadata.InitiatedBy)
		return detErrMsg, sdkerrors.Wrapf(err, "invalid initiatedBy address for refund")
	}

	// Determine refund source (approver/escrow address)
	refundSourceAddr := transferMetadata.ApproverAddress
	if refundSourceAddr == "" && transferMetadata.ApprovalLevel == "collection" {
		refundSourceAddr = collection.MintEscrowAddress
	}
	if refundSourceAddr == "" {
		detErrMsg := "approver address is required for time-based refund (refund source)"
		return detErrMsg, sdkerrors.Wrap(types.ErrInvalidAddress, detErrMsg)
	}

	refundSourceAcc, err := sdk.AccAddressFromBech32(refundSourceAddr)
	if err != nil {
		detErrMsg := fmt.Sprintf("invalid refund source address: %s", refundSourceAddr)
		return detErrMsg, sdkerrors.Wrapf(err, "invalid refund source address")
	}

	for _, coin := range refundCoins {
		if coin.Amount.IsZero() {
			continue
		}

		err = k.sendManagerKeeper.SendCoinWithAliasRouting(ctx, refundSourceAcc, refundRecipientAcc, &coin)
		if err != nil {
			detErrMsg := fmt.Sprintf("insufficient %s balance for time-based refund", formatDenomForDisplay(coin.Denom))
			return detErrMsg, sdkerrors.Wrap(types.ErrUnderflow, err.Error())
		}

		*coinTransfersUsed = append(*coinTransfersUsed, CoinTransfers{
			From:   refundSourceAcc.String(),
			To:     refundRecipientAcc.String(),
			Amount: coin.Amount.String(),
			Denom:  coin.Denom,
		})
	}

	return "", nil
}
