package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/cosmos/gogoproto/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetDefaultBalanceStoreForCollection creates a default balance store from collection defaults
// This performs a deep copy using proto Marshal/Unmarshal to ensure modifications don't affect the original collection defaults
func getDefaultBalanceStoreForCollection(collection *types.TokenCollection) *types.UserBalanceStore {
	if collection.DefaultBalances == nil {
		return &types.UserBalanceStore{}
	}

	// Use proto Marshal/Unmarshal for deep copy (standard proto deep copy pattern)
	// If marshal/unmarshal fails, it indicates a programming error and we panic
	data, err := proto.Marshal(collection.DefaultBalances)
	if err != nil {
		panic(fmt.Sprintf("failed to marshal UserBalanceStore for deep copy: %v", err))
	}

	var copied types.UserBalanceStore
	if err := proto.Unmarshal(data, &copied); err != nil {
		panic(fmt.Sprintf("failed to unmarshal UserBalanceStore for deep copy: %v", err))
	}

	return &copied
}

// GetBalanceOrApplyDefault retrieves user balance or applies default balance store
func (k Keeper) GetBalanceOrApplyDefault(ctx sdk.Context, collection *types.TokenCollection, userAddress string) (*types.UserBalanceStore, bool) {
	//Mint has unlimited balances
	if types.IsMintOrTotalAddress(userAddress) {
		return &types.UserBalanceStore{}, false
	}

	//Special backed addresses also have unlimited balances
	if k.IsSpecialBackedAddress(ctx, collection, userAddress) {
		return &types.UserBalanceStore{
			Balances:                        []*types.Balance{},
			AutoApproveAllIncomingTransfers: true,
			AutoApproveSelfInitiatedIncomingTransfers: true,
			AutoApproveSelfInitiatedOutgoingTransfers: true,
		}, false
	}

	//We get current balances or fallback to default balances
	balanceKey := ConstructBalanceKey(userAddress, collection.CollectionId)
	balance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
	appliedDefault := false
	if !found {
		balance = getDefaultBalanceStoreForCollection(collection)
		appliedDefault = true

		// Initialize approval versions for default approvals when balance is first accessed.
		// This is necessary to ensure approval versioning works correctly for default approvals.
		// Note: This has a side effect - first access to a balance increments approval versions,
		// which is intentional to prevent replay attacks using old default approval versions.
		// The version is incremented (not set to 0) to ensure uniqueness and prevent conflicts.
		for _, approval := range balance.IncomingApprovals {
			approval.Version = k.IncrementApprovalVersion(ctx, collection.CollectionId, "incoming", userAddress, approval.ApprovalId)
		}
		for _, approval := range balance.OutgoingApprovals {
			approval.Version = k.IncrementApprovalVersion(ctx, collection.CollectionId, "outgoing", userAddress, approval.ApprovalId)
		}
	}

	if balance.UserPermissions == nil {
		balance.UserPermissions = &types.UserPermissions{}
	}

	return balance, appliedDefault
}

// SetBalanceForAddress stores a user balance for a specific address
func (k Keeper) SetBalanceForAddress(ctx sdk.Context, collection *types.TokenCollection, userAddress string, balance *types.UserBalanceStore) error {
	balanceKey := ConstructBalanceKey(userAddress, collection.CollectionId)
	return k.SetUserBalanceInStore(ctx, balanceKey, balance, false)
}
