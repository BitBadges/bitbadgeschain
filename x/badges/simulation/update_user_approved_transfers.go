package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgUpdateUserApprovals(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Get a random existing collection
		collectionId, found := GetRandomCollectionId(r, ctx, k)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateUserApprovals, "no collections exist"), nil, nil
		}
		
		// Check if collection exists
		_, found = k.GetCollectionFromStore(ctx, collectionId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateUserApprovals, "collection not found"), nil, nil
		}
		
		simAccount, _ := simtypes.RandomAcc(r, accs)
		
		// Randomly decide which fields to update
		updateOutgoingApprovals := r.Intn(2) == 0
		updateIncomingApprovals := r.Intn(2) == 0
		updateAutoApproveAllIncomingTransfers := r.Intn(2) == 0
		updateAutoApproveSelfInitiatedOutgoingTransfers := r.Intn(2) == 0
		updateAutoApproveSelfInitiatedIncomingTransfers := r.Intn(2) == 0
		
		msg := &types.MsgUpdateUserApprovals{
			Creator:      simAccount.Address.String(),
			CollectionId: collectionId,
			UpdateOutgoingApprovals:                        updateOutgoingApprovals,
			UpdateIncomingApprovals:                        updateIncomingApprovals,
			UpdateAutoApproveAllIncomingTransfers:          updateAutoApproveAllIncomingTransfers,
			UpdateAutoApproveSelfInitiatedOutgoingTransfers: updateAutoApproveSelfInitiatedOutgoingTransfers,
			UpdateAutoApproveSelfInitiatedIncomingTransfers: updateAutoApproveSelfInitiatedIncomingTransfers,
		}
		
		// Set values only if updating
		if updateOutgoingApprovals {
			outgoingApprovals := []*types.UserOutgoingApproval{}
			if r.Intn(2) == 0 {
				approvalId := simtypes.RandStringOfLength(r, 10)
				toListId := "All"
				if r.Intn(3) == 0 {
					toListId = GetRandomAddresses(r, 1, accs)[0]
				}
				outgoingApprovals = append(outgoingApprovals, &types.UserOutgoingApproval{
					ApprovalId:        approvalId,
					ToListId:          toListId,
					InitiatedByListId: "All",
					TransferTimes:     GetTimelineTimes(r, 1),
					TokenIds:          GetTimelineTimes(r, 1),
					OwnershipTimes:    GetTimelineTimes(r, 1),
					ApprovalCriteria:  &types.OutgoingApprovalCriteria{},
					Version:           sdkmath.NewUint(0),
				})
			}
			msg.OutgoingApprovals = outgoingApprovals
		}
		
		if updateIncomingApprovals {
			incomingApprovals := []*types.UserIncomingApproval{}
			if r.Intn(2) == 0 {
				approvalId := simtypes.RandStringOfLength(r, 10)
				fromListId := "All"
				if r.Intn(3) == 0 {
					fromListId = GetRandomAddresses(r, 1, accs)[0]
				}
				incomingApprovals = append(incomingApprovals, &types.UserIncomingApproval{
					ApprovalId:        approvalId,
					FromListId:        fromListId,
					InitiatedByListId: "All",
					TransferTimes:     GetTimelineTimes(r, 1),
					TokenIds:          GetTimelineTimes(r, 1),
					OwnershipTimes:    GetTimelineTimes(r, 1),
					ApprovalCriteria:  &types.IncomingApprovalCriteria{},
					Version:           sdkmath.NewUint(0),
				})
			}
			msg.IncomingApprovals = incomingApprovals
		}
		
		if updateAutoApproveAllIncomingTransfers {
			msg.AutoApproveAllIncomingTransfers = r.Intn(2) == 0
		}
		if updateAutoApproveSelfInitiatedOutgoingTransfers {
			msg.AutoApproveSelfInitiatedOutgoingTransfers = r.Intn(2) == 0
		}
		if updateAutoApproveSelfInitiatedIncomingTransfers {
			msg.AutoApproveSelfInitiatedIncomingTransfers = r.Intn(2) == 0
		}
		
		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}
		
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
