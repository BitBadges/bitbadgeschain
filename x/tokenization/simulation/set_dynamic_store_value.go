package simulation

import (
	"math/rand"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

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
		// Ensure we have valid accounts
		if len(accs) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetDynamicStoreValue, "no accounts available"), nil, nil
		}

		simAccount := EnsureAccountExists(r, accs)

		// Try to get a known-good dynamic store ID first
		storeId, found := GetKnownGoodDynamicStoreId(ctx, k)
		if !found {
			// Fallback: try to get a random existing store ID
			nextStoreId := k.GetNextDynamicStoreId(ctx)
			if nextStoreId.LTE(sdkmath.NewUint(1)) {
				// No dynamic stores exist - try to create one first
				createMsg := &types.MsgCreateDynamicStore{
					Creator:      simAccount.Address.String(),
					DefaultValue: r.Intn(2) == 0,
				}
				msgServer := keeper.NewMsgServerImpl(k)
				createResp, err := msgServer.CreateDynamicStore(ctx, createMsg)
				if err != nil {
					return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetDynamicStoreValue, "no dynamic stores exist and failed to create one"), nil, nil
				}
				storeId = createResp.StoreId
			} else {
				// Get a random existing store ID (stores exist from 1 to (nextStoreId - 1))
				maxId := nextStoreId.Sub(sdkmath.NewUint(1))
				if maxId.IsZero() {
					return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetDynamicStoreValue, "no dynamic stores exist"), nil, nil
				}
				// Random ID between 1 and maxId
				storeId = sdkmath.NewUint(uint64(r.Int63n(int64(maxId.Uint64()))) + 1)
			}
		}

		// Verify the store exists and get it to check creator
		store, found := k.GetDynamicStoreFromStore(ctx, storeId)
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgSetDynamicStoreValue, "dynamic store not found"), nil, nil
		}

		// Only the creator can set values, so use the creator account
		creatorAccount, found := FindAccount(accs, store.CreatedBy)
		if !found {
			// Creator not in simulation accounts - this will fail, but try anyway for simulation coverage
			creatorAccount = simAccount
		}

		// Random address (could be the creator or another account)
		targetAccount := EnsureAccountExists(r, accs)

		// Random boolean value
		value := r.Intn(2) == 0

		// Use creator as the message creator (required for setting values)
		msg := &types.MsgSetDynamicStoreValue{
			Creator: creatorAccount.Address.String(),
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
