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

func SimulateMsgUpdateDynamicStore(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		// Get next dynamic store ID to determine if any stores exist
		nextStoreId := k.GetNextDynamicStoreId(ctx)
		if nextStoreId.LTE(sdkmath.NewUint(1)) {
			// No dynamic stores exist yet
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateDynamicStore, "no dynamic stores exist"), nil, nil
		}

		// Get a random existing store ID (stores exist from 1 to (nextStoreId - 1))
		maxId := nextStoreId.Sub(sdkmath.NewUint(1))
		if maxId.IsZero() {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateDynamicStore, "no dynamic stores exist"), nil, nil
		}

		// Random ID between 1 and maxId
		storeId := sdkmath.NewUint(uint64(r.Int63n(int64(maxId.Uint64()))) + 1)

		// Verify the store exists and get it to check creator
		store, found := k.GetDynamicStoreFromStore(ctx, storeId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgUpdateDynamicStore, "dynamic store not found"), nil, nil
		}

		// Only the creator can update, so find the creator account
		creatorAccount, found := FindAccount(accs, store.CreatedBy)
		if !found {
			// Creator not in simulation accounts, use random account (will fail validation but that's ok for simulation)
			creatorAccount = simAccount
		}

		// Random boolean for defaultValue
		defaultValue := r.Intn(2) == 0

		// Random boolean for globalEnabled (global kill switch)
		globalEnabled := r.Intn(2) == 0

		msg := types.NewMsgUpdateDynamicStoreWithGlobalEnabled(
			creatorAccount.Address.String(),
			storeId,
			defaultValue,
			globalEnabled,
		)

		// Validate message
		if err := msg.ValidateBasic(); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), err.Error()), nil, nil
		}

		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
