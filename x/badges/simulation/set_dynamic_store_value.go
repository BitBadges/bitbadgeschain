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

func SimulateMsgSetDynamicStoreValue(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		// Get a random existing collection (dynamic stores are associated with collections)
		collectionId, found := GetRandomCollectionId(r, ctx, k)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetDynamicStoreValue, "no collections exist"), nil, nil
		}

		// Check if collection exists
		_, found = k.GetCollectionFromStore(ctx, collectionId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetDynamicStoreValue, "collection not found"), nil, nil
		}

		// For now, use a random store ID (we'd need to track dynamic stores to get existing ones)
		// Using a small random ID that might exist
		storeId := sdkmath.NewUint(uint64(r.Int63n(10)) + 1)

		// Random address (could be the creator or another account)
		targetAccount, _ := simtypes.RandomAcc(r, accs)

		// Random boolean value
		value := r.Intn(2) == 0

		msg := &types.MsgSetDynamicStoreValue{
			Creator: simAccount.Address.String(),
			StoreId: storeId,
			Address: targetAccount.Address.String(),
			Value:   value,
		}

		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
