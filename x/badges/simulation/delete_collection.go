package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgDeleteCollection(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Get a random existing collection
		collectionId, found := GetRandomCollectionId(r, ctx, k)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDeleteCollection, "no collections exist"), nil, nil
		}
		
		// Check if collection exists
		collection, found := k.GetCollectionFromStore(ctx, collectionId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgDeleteCollection, "collection not found"), nil, nil
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
		
		msg := &types.MsgDeleteCollection{
			Creator:      simAccount.Address.String(),
			CollectionId: collectionId,
		}
		
		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}
		
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
