package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

const (
	// ProtocolFeeDenominator represents the denominator for protocol fee calculation (0.1% = 1/1000)
	ProtocolFeeDenominator = 1000
)

// CalculateAndDistributeProtocolFees calculates protocol fees from coin transfers and distributes them
// to the community pool
func (k Keeper) CalculateAndDistributeProtocolFees(
	ctx sdk.Context,
	coinTransfers []CoinTransfers,
	initiatedBy string,
) ([]CoinTransfers, error) {
	// Calculate protocol fees for all denoms (0.1% of each denom transferred)
	protocolFees := sdk.NewCoins()
	denomAmounts := make(map[string]sdkmath.Uint)

	for _, coinTransfer := range coinTransfers {
		amount := sdkmath.NewUintFromString(coinTransfer.Amount)
		//initialize it if it doesn't exist
		if _, ok := denomAmounts[coinTransfer.Denom]; !ok {
			denomAmounts[coinTransfer.Denom] = sdkmath.NewUint(0)
		}

		denomAmounts[coinTransfer.Denom] = denomAmounts[coinTransfer.Denom].Add(amount)
	}

	for denom, totalAmount := range denomAmounts {
		// 0.1% of the total amount for this denom
		// Safety check to prevent division by zero
		if ProtocolFeeDenominator == 0 {
			return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "protocol fee denominator cannot be zero")
		}
		protocolFee := totalAmount.Quo(sdkmath.NewUint(ProtocolFeeDenominator))

		// For other denoms, just use 0.1%
		if !protocolFee.IsZero() {
			protocolFees = protocolFees.Add(sdk.NewCoin(denom, sdkmath.NewIntFromUint64(protocolFee.Uint64())))
		}
	}

	fromAddressAcc, err := sdk.AccAddressFromBech32(initiatedBy)
	if err != nil {
		return nil, err
	}

	var protocolFeeTransfers []CoinTransfers

	if !protocolFees.IsZero() {
		// Send all fees to community pool using FundCommunityPoolWithAliasRouting to support wrapped badge denoms
		err = k.sendManagerKeeper.FundCommunityPoolWithAliasRouting(ctx, fromAddressAcc, protocolFees)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "error funding community pool with protocol fees: %s", protocolFees)
		}

		// Add all protocol fees to coinTransfers for community pool
		for _, protocolFee := range protocolFees {
			protocolFeeTransfers = append(protocolFeeTransfers, CoinTransfers{
				From:          initiatedBy,
				To:            authtypes.NewModuleAddress(distrtypes.ModuleName).String(),
				Amount:        protocolFee.Amount.String(),
				Denom:         protocolFee.Denom,
				IsProtocolFee: true,
			})
		}
	}

	return protocolFeeTransfers, nil
}

// HandleAutoDeletions processes auto-deletion logic for approvals after transfers
func (k Keeper) HandleAutoDeletions(
	ctx sdk.Context,
	collection *types.TokenCollection,
	fromUserBalance *types.UserBalanceStore,
	toUserBalance *types.UserBalanceStore,
	approvalsUsed []ApprovalsUsed,
) (*types.UserBalanceStore, *types.UserBalanceStore, error) {
	var err error

	// Helper functions for auto-deletion checks
	isDeleteAfterOneUse := func(autoDeletionOptions *types.AutoDeletionOptions) bool {
		if autoDeletionOptions == nil {
			return false
		}
		return autoDeletionOptions.AfterOneUse
	}

	isDeleteAfterOverallMaxNumTransfersForCollection := func(autoDeletionOptions *types.AutoDeletionOptions, approvalCriteria *types.ApprovalCriteria, approvalUsed ApprovalsUsed) bool {
		if autoDeletionOptions == nil || !autoDeletionOptions.AfterOverallMaxNumTransfers {
			return false
		}

		// Check if overall max number of transfers threshold is set
		if approvalCriteria == nil || approvalCriteria.MaxNumTransfers == nil || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsNil() || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsZero() {
			return false
		}

		// Get the tracker to check current number of transfers
		maxNumTransfersTrackerId := approvalCriteria.MaxNumTransfers.AmountTrackerId
		if maxNumTransfersTrackerId == "" {
			return false
		}

		// Get the current tracker details
		trackerDetails, err := k.GetApprovalTrackerFromStoreAndResetIfNeeded(
			ctx,
			collection.CollectionId,
			approvalUsed.ApproverAddress,
			approvalUsed.ApprovalId,
			maxNumTransfersTrackerId,
			approvalUsed.ApprovalLevel,
			"overall",
			"",
			approvalCriteria.MaxNumTransfers.ResetTimeIntervals,
			true,
		)
		if err != nil {
			return false
		}

		// Check if the current number of transfers has reached or exceeded the threshold
		return trackerDetails.NumTransfers.GTE(approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers)
	}

	isDeleteAfterOverallMaxNumTransfersForOutgoing := func(autoDeletionOptions *types.AutoDeletionOptions, approvalCriteria *types.OutgoingApprovalCriteria, approvalUsed ApprovalsUsed) bool {
		if autoDeletionOptions == nil || !autoDeletionOptions.AfterOverallMaxNumTransfers {
			return false
		}

		// Check if overall max number of transfers threshold is set
		if approvalCriteria == nil || approvalCriteria.MaxNumTransfers == nil || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsNil() || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsZero() {
			return false
		}

		// Get the tracker to check current number of transfers
		maxNumTransfersTrackerId := approvalCriteria.MaxNumTransfers.AmountTrackerId
		if maxNumTransfersTrackerId == "" {
			return false
		}

		// Get the current tracker details
		trackerDetails, err := k.GetApprovalTrackerFromStoreAndResetIfNeeded(
			ctx,
			collection.CollectionId,
			approvalUsed.ApproverAddress,
			approvalUsed.ApprovalId,
			maxNumTransfersTrackerId,
			approvalUsed.ApprovalLevel,
			"overall",
			"",
			approvalCriteria.MaxNumTransfers.ResetTimeIntervals,
			true,
		)
		if err != nil {
			return false
		}

		// Check if the current number of transfers has reached or exceeded the threshold
		return trackerDetails.NumTransfers.GTE(approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers)
	}

	isDeleteAfterOverallMaxNumTransfersForIncoming := func(autoDeletionOptions *types.AutoDeletionOptions, approvalCriteria *types.IncomingApprovalCriteria, approvalUsed ApprovalsUsed) bool {
		if autoDeletionOptions == nil || !autoDeletionOptions.AfterOverallMaxNumTransfers {
			return false
		}

		// Check if overall max number of transfers threshold is set
		if approvalCriteria == nil || approvalCriteria.MaxNumTransfers == nil || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsNil() || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsZero() {
			return false
		}

		// Get the tracker to check current number of transfers
		maxNumTransfersTrackerId := approvalCriteria.MaxNumTransfers.AmountTrackerId
		if maxNumTransfersTrackerId == "" {
			return false
		}

		// Get the current tracker details
		trackerDetails, err := k.GetApprovalTrackerFromStoreAndResetIfNeeded(
			ctx,
			collection.CollectionId,
			approvalUsed.ApproverAddress,
			approvalUsed.ApprovalId,
			maxNumTransfersTrackerId,
			approvalUsed.ApprovalLevel,
			"overall",
			"",
			approvalCriteria.MaxNumTransfers.ResetTimeIntervals,
			true,
		)
		if err != nil {
			return false
		}

		// Check if the current number of transfers has reached or exceeded the threshold
		return trackerDetails.NumTransfers.GTE(approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers)
	}

	// Per-transfer, we handle auto-deletions if applicable
	for _, approvalUsed := range approvalsUsed {
		if approvalUsed.ApprovalLevel == "incoming" {
			newIncomingApprovals := []*types.UserIncomingApproval{}
			for _, incomingApproval := range toUserBalance.IncomingApprovals {
				if incomingApproval.ApprovalId != approvalUsed.ApprovalId {
					newIncomingApprovals = append(newIncomingApprovals, incomingApproval)
				} else {
					shouldDelete := false

					// Check if should delete after one use (doesn't depend on ApprovalCriteria)
					if incomingApproval.ApprovalCriteria != nil && incomingApproval.ApprovalCriteria.AutoDeletionOptions != nil {
						shouldDelete = isDeleteAfterOneUse(incomingApproval.ApprovalCriteria.AutoDeletionOptions)
					}

					// Check if should delete after overall max transfers (depends on ApprovalCriteria)
					if !shouldDelete && incomingApproval.ApprovalCriteria != nil {
						shouldDelete = isDeleteAfterOverallMaxNumTransfersForIncoming(incomingApproval.ApprovalCriteria.AutoDeletionOptions, incomingApproval.ApprovalCriteria, approvalUsed)
					}

					if !shouldDelete {
						newIncomingApprovals = append(newIncomingApprovals, incomingApproval)
					}
					// If shouldDelete is true, we simply don't add it to the new slice (effectively deleting it)
				}
			}
			toUserBalance.IncomingApprovals = newIncomingApprovals
		} else if approvalUsed.ApprovalLevel == "outgoing" {
			newOutgoingApprovals := []*types.UserOutgoingApproval{}
			for _, outgoingApproval := range fromUserBalance.OutgoingApprovals {
				if outgoingApproval.ApprovalId != approvalUsed.ApprovalId {
					newOutgoingApprovals = append(newOutgoingApprovals, outgoingApproval)
				} else {
					shouldDelete := false

					// Check if should delete after one use (doesn't depend on ApprovalCriteria)
					if outgoingApproval.ApprovalCriteria != nil && outgoingApproval.ApprovalCriteria.AutoDeletionOptions != nil {
						shouldDelete = isDeleteAfterOneUse(outgoingApproval.ApprovalCriteria.AutoDeletionOptions)
					}

					// Check if should delete after overall max transfers (depends on ApprovalCriteria)
					if !shouldDelete && outgoingApproval.ApprovalCriteria != nil {
						shouldDelete = isDeleteAfterOverallMaxNumTransfersForOutgoing(outgoingApproval.ApprovalCriteria.AutoDeletionOptions, outgoingApproval.ApprovalCriteria, approvalUsed)
					}

					if !shouldDelete {
						newOutgoingApprovals = append(newOutgoingApprovals, outgoingApproval)
					}
					// If shouldDelete is true, we simply don't add it to the new slice (effectively deleting it)
				}
			}
			fromUserBalance.OutgoingApprovals = newOutgoingApprovals
		} else if approvalUsed.ApprovalLevel == "collection" {
			newCollectionApprovals := []*types.CollectionApproval{}
			edited := false
			for _, collectionApproval := range collection.CollectionApprovals {
				if collectionApproval.ApprovalId != approvalUsed.ApprovalId {
					newCollectionApprovals = append(newCollectionApprovals, collectionApproval)
				} else {
					shouldDelete := false

					// Check if should delete after one use (doesn't depend on ApprovalCriteria)
					if collectionApproval.ApprovalCriteria != nil && collectionApproval.ApprovalCriteria.AutoDeletionOptions != nil {
						shouldDelete = isDeleteAfterOneUse(collectionApproval.ApprovalCriteria.AutoDeletionOptions)
					}

					// Check if should delete after overall max transfers (depends on ApprovalCriteria)
					if !shouldDelete && collectionApproval.ApprovalCriteria != nil {
						shouldDelete = isDeleteAfterOverallMaxNumTransfersForCollection(collectionApproval.ApprovalCriteria.AutoDeletionOptions, collectionApproval.ApprovalCriteria, approvalUsed)
					}

					if !shouldDelete {
						newCollectionApprovals = append(newCollectionApprovals, collectionApproval)
					} else {
						// Delete the approval
						edited = true
					}
				}
			}

			collection.CollectionApprovals = newCollectionApprovals
			if edited {
				err = k.SetCollectionInStore(ctx, collection, false)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
				}
			}
		}
	}

	return fromUserBalance, toUserBalance, nil
}

// EmitUsedApprovalDetailsEvent emits an event with details about approvals used and coin transfers
func EmitUsedApprovalDetailsEvent(ctx sdk.Context, collectionId sdkmath.Uint, from string, to string, initiatedBy string, coinTransfers []CoinTransfers, approvalsUsed []ApprovalsUsed, balances []*types.Balance) (err error) {
	marshalToString := func(v interface{}) (string, error) {
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	coinTransfersStr, err := marshalToString(coinTransfers)
	if err != nil {
		return err
	}

	approvalsUsedStr, err := marshalToString(approvalsUsed)
	if err != nil {
		return err
	}

	balancesStr, err := marshalToString(balances)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("usedApprovalDetails",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
			sdk.NewAttribute("from", from),
			sdk.NewAttribute("to", to),
			sdk.NewAttribute("initiatedBy", initiatedBy),
			sdk.NewAttribute("coinTransfers", coinTransfersStr),
			sdk.NewAttribute("approvalsUsed", approvalsUsedStr),
			sdk.NewAttribute("balances", balancesStr),
		),
	)

	return nil
}

// EmitApprovalEvent emits an event for approval tracking
func EmitApprovalEvent(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	approverAddress string,
	approvalId string,
	amountsTrackerId string,
	approvalLevel string,
	trackerType string,
	address string,
	amountsStr string,
	numTransfersStr string,
	lastUpdatedAt sdkmath.Uint,
) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"approval"+fmt.Sprint(collectionId)+fmt.Sprint(approverAddress)+fmt.Sprint(approvalId)+fmt.Sprint(amountsTrackerId)+fmt.Sprint(approvalLevel)+fmt.Sprint(trackerType)+fmt.Sprint(address),
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
			sdk.NewAttribute("approvalId", fmt.Sprint(approvalId)),
			sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
			sdk.NewAttribute("amountTrackerId", fmt.Sprint(amountsTrackerId)),
			sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
			sdk.NewAttribute("trackerType", fmt.Sprint(trackerType)),
			sdk.NewAttribute("approvedAddress", fmt.Sprint(address)),
			sdk.NewAttribute("amounts", amountsStr),
			sdk.NewAttribute("numTransfers", numTransfersStr),
			sdk.NewAttribute("lastUpdatedAt", fmt.Sprint(lastUpdatedAt)),
		),
	)
}
