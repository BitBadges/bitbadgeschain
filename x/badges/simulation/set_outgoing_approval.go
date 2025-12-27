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

func SimulateMsgSetOutgoingApproval(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Get a random existing collection
		collectionId, found := GetRandomCollectionId(r, ctx, k)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetOutgoingApproval, "no collections exist"), nil, nil
		}
		
		// Check if collection exists
		_, found = k.GetCollectionFromStore(ctx, collectionId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetOutgoingApproval, "collection not found"), nil, nil
		}
		
		simAccount, _ := simtypes.RandomAcc(r, accs)
		
		// Generate outgoing approval
		approvalId := simtypes.RandStringOfLength(r, 10)
		toListId := "All"
		if r.Intn(3) == 0 {
			toListId = GetRandomAddresses(r, 1, accs)[0]
		}
		
		approval := &types.UserOutgoingApproval{
			ApprovalId:        approvalId,
			ToListId:          toListId,
			InitiatedByListId: "All",
			TransferTimes:     GetTimelineTimes(r, 1),
			TokenIds:          GetTimelineTimes(r, 1),
			OwnershipTimes:    GetTimelineTimes(r, 1),
			ApprovalCriteria:  &types.OutgoingApprovalCriteria{},
			Version:           sdkmath.NewUint(0),
		}
		
		msg := &types.MsgSetOutgoingApproval{
			Creator:      simAccount.Address.String(),
			CollectionId: collectionId,
			Approval:     approval,
		}
		
		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}
		
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

