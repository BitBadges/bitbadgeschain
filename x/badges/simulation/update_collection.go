package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgUniversalUpdateCollection(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Get a random existing collection
		collectionId, found := GetRandomCollectionId(r, ctx, k)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUniversalUpdateCollection, "no collections exist"), nil, nil
		}
		
		// Check if collection exists
		collection, found := k.GetCollectionFromStore(ctx, collectionId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUniversalUpdateCollection, "collection not found"), nil, nil
		}
		
		// Use the collection's manager as the creator (or random account with 50% chance)
		var simAccount simtypes.Account
		if r.Intn(2) == 0 && collection.Manager != "" {
			// Try to use the collection's manager
			managerAcc, found := FindAccount(accs, collection.Manager)
			if found {
				simAccount = managerAcc
			} else {
				simAccount, _ = simtypes.RandomAcc(r, accs)
			}
		} else {
			simAccount, _ = simtypes.RandomAcc(r, accs)
		}
		
		// Randomly decide which fields to update
		updateCollectionPermissions := r.Intn(2) == 0
		updateIsArchived := r.Intn(2) == 0
		updateManager := r.Intn(2) == 0
		updateCollectionMetadata := r.Intn(2) == 0
		updateTokenMetadata := r.Intn(2) == 0
		updateCustomData := r.Intn(2) == 0
		updateCollectionApprovals := r.Intn(2) == 0
		updateStandards := r.Intn(2) == 0
		updateValidTokenIds := r.Intn(2) == 0
		
		msg := &types.MsgUniversalUpdateCollection{
			Creator:                     simAccount.Address.String(),
			CollectionId:                collectionId,
			UpdateCollectionPermissions: updateCollectionPermissions,
			UpdateIsArchived:            updateIsArchived,
			UpdateManager:               updateManager,
			UpdateCollectionMetadata:    updateCollectionMetadata,
			UpdateTokenMetadata:         updateTokenMetadata,
			UpdateCustomData:            updateCustomData,
			UpdateCollectionApprovals:   updateCollectionApprovals,
			UpdateStandards:             updateStandards,
			UpdateValidTokenIds:         updateValidTokenIds,
		}
		
		// Set values only if updating
		if updateIsArchived {
			msg.IsArchived = r.Intn(2) == 0
		}
		if updateValidTokenIds {
			msg.ValidTokenIds = GetRandomValidTokenIds(r, collection, r.Intn(3)+1)
		}
		if updateCollectionApprovals {
			approvals := []*types.CollectionApproval{}
			if r.Intn(2) == 0 {
				approvals = append(approvals, GetRandomCollectionApproval(r, accs))
			}
			msg.CollectionApprovals = approvals
		}
		if updateManager {
			newManager, _ := simtypes.RandomAcc(r, accs)
			msg.Manager = newManager.Address.String()
		}
		if updateCollectionPermissions {
			msg.CollectionPermissions = GetRandomCollectionPermissions(r, accs)
		}
		if updateCollectionMetadata {
			msg.CollectionMetadata = &types.CollectionMetadata{
				Uri:        "https://example.com/metadata/" + simtypes.RandStringOfLength(r, 10),
				CustomData: simtypes.RandStringOfLength(r, 20),
			}
		}
		if updateTokenMetadata {
			msg.TokenMetadata = []*types.TokenMetadata{}
		}
		if updateCustomData {
			msg.CustomData = simtypes.RandStringOfLength(r, 20)
		}
		if updateStandards {
			msg.Standards = []string{}
		}
		
		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}
		
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
