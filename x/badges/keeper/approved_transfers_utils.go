package keeper

import (
	"bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DeductUserOutgoingApprovals will check if the current transfer is approved from the from's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserOutgoingApprovals(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	originalTransferBalances []*types.Balance,
	transfer *types.Transfer,
	from string,
	to string,
	requester string,
	userBalance *types.UserBalanceStore,
) error {
	currApprovals := userBalance.OutgoingApprovals
	if userBalance.AutoApproveSelfInitiatedOutgoingTransfers {
		currApprovals = AppendSelfInitiatedOutgoingApproval(currApprovals, from)
	}

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedApprovals := types.CastOutgoingTransfersToCollectionTransfers(currApprovals, from)
	_, err := k.DeductAndGetUserApprovals(
		ctx,
		collection,
		originalTransferBalances,
		transfer,
		castedApprovals,
		to,
		requester,
		"outgoing",
		from,
	)
	return err
}

// DeductUserIncomingApprovals will check if the current transfer is approved from the to's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserIncomingApprovals(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	originalTransferBalances []*types.Balance,
	transfer *types.Transfer,
	to string,
	initiatedBy string,
	userBalance *types.UserBalanceStore,
) error {
	currApprovals := userBalance.IncomingApprovals
	if userBalance.AutoApproveSelfInitiatedIncomingTransfers {
		currApprovals = AppendSelfInitiatedIncomingApproval(currApprovals, to)
	}

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedApprovals := types.CastIncomingTransfersToCollectionTransfers(currApprovals, to)
	_, err := k.DeductAndGetUserApprovals(
		ctx,
		collection,
		originalTransferBalances,
		transfer,
		castedApprovals,
		to,
		initiatedBy,
		"incoming",
		to,
	)
	return err
}

// DeductCollectionApprovalsAndGetUserApprovalsToCheck will check if the current transfer is allowed via the collection's approved transfers and handle any tallying accordingly
func (k Keeper) DeductCollectionApprovalsAndGetUserApprovalsToCheck(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	transfer *types.Transfer,
	toAddress string,
	initiatedBy string,
) ([]*UserApprovalsToCheck, error) {
	return k.DeductAndGetUserApprovals(
		ctx,
		collection,
		transfer.Balances,
		transfer,
		collection.CollectionApprovals,
		toAddress,
		initiatedBy,
		"collection",
		"",
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

func SortViaPrioritizedApprovals(_approvals []*types.CollectionApproval, transfer *types.Transfer, approvalLevel string, approverAddress string) []*types.CollectionApproval {
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

	return approvals
}
