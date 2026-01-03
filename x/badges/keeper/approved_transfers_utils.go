package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DeductUserOutgoingApprovals will check if the current transfer is approved from the from's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserOutgoingApprovals(
	ctx sdk.Context,
	collection *types.TokenCollection,
	originalTransferBalances []*types.Balance,
	transfer *types.Transfer,
	transferMetadata TransferMetadata,
	userBalance *types.UserBalanceStore,
	eventTracking *EventTracking,
	royalties *types.UserRoyalties,
) error {
	from := transferMetadata.From
	currApprovals := userBalance.OutgoingApprovals
	if userBalance.AutoApproveSelfInitiatedOutgoingTransfers {
		currApprovals = AppendSelfInitiatedOutgoingApproval(currApprovals, from)
	}

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedApprovals := types.CastOutgoingTransfersToCollectionTransfers(currApprovals, from)

	// We do not care about the return value here because it is user-level
	_, err := k.DeductAndGetUserApprovals(
		ctx,
		collection,
		originalTransferBalances,
		transfer,
		castedApprovals,
		transferMetadata,
		eventTracking,
		royalties,
		"outgoing",
	)
	return err
}

func (k Keeper) DeductUserIncomingApprovals(
	ctx sdk.Context,
	collection *types.TokenCollection,
	originalTransferBalances []*types.Balance,
	transfer *types.Transfer,
	transferMetadata TransferMetadata,
	userBalance *types.UserBalanceStore,
	eventTracking *EventTracking,
	royalties *types.UserRoyalties,
) error {
	if userBalance.AutoApproveAllIncomingTransfers {
		return nil
	}

	to := transferMetadata.To
	currApprovals := userBalance.IncomingApprovals
	if userBalance.AutoApproveSelfInitiatedIncomingTransfers {
		currApprovals = AppendSelfInitiatedIncomingApproval(currApprovals, to)
	}

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedApprovals := types.CastIncomingTransfersToCollectionTransfers(currApprovals, to)

	// We do not care about the return value here because it is user-level
	_, err := k.DeductAndGetUserApprovals(
		ctx,
		collection,
		originalTransferBalances,
		transfer,
		castedApprovals,
		transferMetadata,
		eventTracking,
		royalties,
		"incoming",
	)
	return err
}

// DeductCollectionApprovalsAndGetUserApprovalsToCheck will check if the current transfer is allowed via the collection's approved transfers and handle any tallying accordingly
func (k Keeper) DeductCollectionApprovalsAndGetUserApprovalsToCheck(
	ctx sdk.Context,
	collection *types.TokenCollection,
	transfer *types.Transfer,
	transferMetadata TransferMetadata,
	eventTracking *EventTracking,
	approvalLevel string,
) ([]*UserApprovalsToCheck, error) {
	blankRoyalties := &types.UserRoyalties{
		Percentage:    sdkmath.NewUint(0),
		PayoutAddress: "",
	}

	return k.DeductAndGetUserApprovals(
		ctx,
		collection,
		transfer.Balances,
		transfer,
		collection.CollectionApprovals,
		transferMetadata,
		eventTracking,
		blankRoyalties,
		approvalLevel,
	)
}

func onlyCheckPrioritizedApprovals(transfer *types.Transfer, approvalLevel string) bool {
	//prioritized approvals are checked first and if the "only" flag is set, we only check prioritized approvals
	onlyCheckPrioritized := false
	if approvalLevel == "collection" && transfer.OnlyCheckPrioritizedCollectionApprovals {
		onlyCheckPrioritized = true
	} else if approvalLevel == "outgoing" && transfer.OnlyCheckPrioritizedOutgoingApprovals {
		onlyCheckPrioritized = true
	} else if approvalLevel == "incoming" && transfer.OnlyCheckPrioritizedIncomingApprovals {
		onlyCheckPrioritized = true
	}

	return onlyCheckPrioritized
}

func FilterApprovalsWithPrioritizedHandling(
	_approvals []*types.CollectionApproval,
	transfer *types.Transfer,
	approvalLevel string,
	approverAddress string,
) ([]*types.CollectionApproval, error) {
	prioritizedApprovals := transfer.PrioritizedApprovals
	onlyCheckPrioritized := onlyCheckPrioritizedApprovals(transfer, approvalLevel)

	//Reorder approvals based on prioritized approvals
	//We want to check prioritized approvals first
	//If onlyCheckPrioritized is true, we only check prioritized approvals and ignore the rest
	approvals := []*types.CollectionApproval{}
	prioritizedTransfers := []*types.CollectionApproval{}
	for _, approval := range _approvals {
		prioritized := false

		for _, prioritizedApproval := range prioritizedApprovals {
			if approval.ApprovalId == prioritizedApproval.ApprovalId && prioritizedApproval.ApprovalLevel == approvalLevel && approverAddress == prioritizedApproval.ApproverAddress {
				prioritized = true
				break
			}
		}

		if prioritized {
			prioritizedTransfers = append(prioritizedTransfers, approval)
		} else {
			approvals = append(approvals, approval)
		}
	}

	if onlyCheckPrioritized {
		approvals = prioritizedTransfers
	} else {
		approvals = append(prioritizedTransfers, approvals...)
	}

	//Filter approvals where approvalCriteria != nil and not in prioritizedApprovals
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range approvals {
		if types.CollectionApprovalIsAutoScannable(approval.ApprovalCriteria) {
			filteredApprovals = append(filteredApprovals, approval)
			continue
		}

		// Check if this approval is explicitly prioritized and version is correct
		prioritizedAndVersionCorrect := false
		for _, prioritizedApproval := range prioritizedApprovals {
			if approval.ApprovalId == prioritizedApproval.ApprovalId && prioritizedApproval.ApprovalLevel == approvalLevel && approverAddress == prioritizedApproval.ApproverAddress {
				if prioritizedApproval.Version.IsNil() {
					return nil, sdkerrors.Wrapf(types.ErrUintUnititialized, "version is uninitialized")
				}

				if !prioritizedApproval.Version.Equal(approval.Version) {
					return nil, sdkerrors.Wrapf(types.ErrMismatchedVersions, "versions are mismatched for a prioritized approval %s %s %s", prioritizedApproval.ApprovalId, prioritizedApproval.ApproverAddress, prioritizedApproval.ApprovalLevel)
				}

				prioritizedAndVersionCorrect = true
				break
			}
		}

		// Check if this approval has mustPrioritize set
		// If mustPrioritize is true, only include if it's explicitly prioritized
		if approval.ApprovalCriteria.MustPrioritize && !prioritizedAndVersionCorrect {
			continue
		}

		// Include the approval if it's prioritized (maintain original behavior)
		if prioritizedAndVersionCorrect {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}

	return filteredApprovals, nil
}
