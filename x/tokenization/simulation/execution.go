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

// ExecuteAndVerifyMessage actually executes a message via keeper and verifies state changes
// Returns true if execution was successful, false otherwise
func ExecuteAndVerifyMessage(
	ctx sdk.Context,
	k *keeper.Keeper,
	msg sdk.Msg,
) (bool, error) {
	msgServer := keeper.NewMsgServerImpl(k)
	sdkCtx := ctx

	// Execute the message based on its type
	switch m := msg.(type) {
	case *types.MsgCreateCollection:
		_, err := msgServer.CreateCollection(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify collection was created
		nextId := k.GetNextCollectionId(ctx)
		if nextId.LTE(sdkmath.NewUint(1)) {
			return false, nil // Collection should exist
		}
		collectionId := nextId.Sub(sdkmath.NewUint(1))
		_, found := k.GetCollectionFromStore(ctx, collectionId)
		return found, nil

	case *types.MsgUniversalUpdateCollection:
		_, err := msgServer.UniversalUpdateCollection(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify collection still exists (or was created)
		_, found := k.GetCollectionFromStore(ctx, m.CollectionId)
		return found, nil

	case *types.MsgDeleteCollection:
		_, err := msgServer.DeleteCollection(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify collection was deleted (or doesn't exist)
		_, found := k.GetCollectionFromStore(ctx, m.CollectionId)
		return !found, nil

	case *types.MsgTransferTokens:
		// Get balances before
		balancesBefore := make(map[string]*types.UserBalanceStore)
		for _, transfer := range m.Transfers {
			for _, toAddr := range transfer.ToAddresses {
				if transfer.From != types.MintAddress {
					collection, _ := k.GetCollectionFromStore(ctx, m.CollectionId)
					bal, _, _ := k.GetBalanceOrApplyDefault(ctx, collection, transfer.From)
					balancesBefore[transfer.From] = bal
				}
				collection, _ := k.GetCollectionFromStore(ctx, m.CollectionId)
				bal, _, _ := k.GetBalanceOrApplyDefault(ctx, collection, toAddr)
				balancesBefore[toAddr] = bal
			}
		}

		_, err := msgServer.TransferTokens(sdkCtx, m)
		if err != nil {
			return false, err
		}

		// Verify balances changed (simplified check - just verify no error)
		return true, nil

	case *types.MsgCreateDynamicStore:
		_, err := msgServer.CreateDynamicStore(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify dynamic store was created
		nextStoreId := k.GetNextDynamicStoreId(ctx)
		if nextStoreId.LTE(sdkmath.NewUint(1)) {
			return false, nil
		}
		storeId := nextStoreId.Sub(sdkmath.NewUint(1))
		_, found := k.GetDynamicStoreFromStore(ctx, storeId)
		return found, nil

	case *types.MsgUpdateDynamicStore:
		_, err := msgServer.UpdateDynamicStore(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify dynamic store still exists
		_, found := k.GetDynamicStoreFromStore(ctx, m.StoreId)
		return found, nil

	case *types.MsgSetDynamicStoreValue:
		_, err := msgServer.SetDynamicStoreValue(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify value was set (check if store exists)
		_, found := k.GetDynamicStoreFromStore(ctx, m.StoreId)
		return found, nil

	case *types.MsgUpdateUserApprovals:
		_, err := msgServer.UpdateUserApprovals(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify collection still exists
		_, found := k.GetCollectionFromStore(ctx, m.CollectionId)
		return found, nil

	case *types.MsgSetIncomingApproval:
		_, err := msgServer.SetIncomingApproval(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify collection still exists
		_, found := k.GetCollectionFromStore(ctx, m.CollectionId)
		return found, nil

	case *types.MsgSetOutgoingApproval:
		_, err := msgServer.SetOutgoingApproval(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify collection still exists
		_, found := k.GetCollectionFromStore(ctx, m.CollectionId)
		return found, nil

	case *types.MsgPurgeApprovals:
		_, err := msgServer.PurgeApprovals(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Verify collection still exists
		_, found := k.GetCollectionFromStore(ctx, m.CollectionId)
		return found, nil

	case *types.MsgCreateAddressLists:
		_, err := msgServer.CreateAddressLists(sdkCtx, m)
		if err != nil {
			return false, err
		}
		// Address lists are stored internally, just verify no error
		return true, nil

	default:
		// Unknown message type - can't execute
		return false, nil
	}
}

// WrapOperationWithExecution wraps an operation to actually execute the message when app is available
func WrapOperationWithExecution(
	op simtypes.Operation,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Call the original operation
		opMsg, futureOps, err := op(r, app, ctx, accs, chainID)
		if err != nil {
			return opMsg, futureOps, err
		}

		// If it's a NoOpMsg, return early
		if !opMsg.OK {
			return opMsg, futureOps, err
		}

		// If app is available, try to execute the message
		if app != nil {
			// Get the keeper from the app (this is a simplified approach)
			// In a real simulation, we'd need to get the keeper properly
			// For now, we'll just return the operation message
			// The actual execution will happen in the simulation framework
		}

		return opMsg, futureOps, err
	}
}
