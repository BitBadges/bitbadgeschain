package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetDefaultBalanceStoreForCollection creates a default balance store from collection defaults
func GetDefaultBalanceStoreForCollection(collection *types.TokenCollection) *types.UserBalanceStore {
	return &types.UserBalanceStore{
		Balances:          collection.DefaultBalances.Balances,
		OutgoingApprovals: collection.DefaultBalances.OutgoingApprovals,
		IncomingApprovals: collection.DefaultBalances.IncomingApprovals,
		AutoApproveSelfInitiatedOutgoingTransfers: collection.DefaultBalances.AutoApproveSelfInitiatedOutgoingTransfers,
		AutoApproveSelfInitiatedIncomingTransfers: collection.DefaultBalances.AutoApproveSelfInitiatedIncomingTransfers,
		AutoApproveAllIncomingTransfers:           collection.DefaultBalances.AutoApproveAllIncomingTransfers,
		UserPermissions:                           collection.DefaultBalances.UserPermissions,
	}
}

// GetBalanceOrApplyDefault retrieves user balance or applies default balance store
func (k Keeper) GetBalanceOrApplyDefault(ctx sdk.Context, collection *types.TokenCollection, userAddress string) (*types.UserBalanceStore, bool) {
	//Mint has unlimited balances
	if types.IsSpecialAddress(userAddress) {
		return &types.UserBalanceStore{}, false
	}

	//We get current balances or fallback to default balances
	balanceKey := ConstructBalanceKey(userAddress, collection.CollectionId)
	balance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
	appliedDefault := false
	if !found {
		balance = GetDefaultBalanceStoreForCollection(collection)
		appliedDefault = true

		// We need to set the version to "0" for all incoming and outgoing approvals
		for _, approval := range balance.IncomingApprovals {
			approval.Version = k.IncrementApprovalVersion(ctx, collection.CollectionId, "incoming", userAddress, approval.ApprovalId)
		}
		for _, approval := range balance.OutgoingApprovals {
			approval.Version = k.IncrementApprovalVersion(ctx, collection.CollectionId, "outgoing", userAddress, approval.ApprovalId)
		}
	}

	return balance, appliedDefault
}

// SetBalanceForAddress stores a user balance for a specific address
func (k Keeper) SetBalanceForAddress(ctx sdk.Context, collection *types.TokenCollection, userAddress string, balance *types.UserBalanceStore) error {
	balanceKey := ConstructBalanceKey(userAddress, collection.CollectionId)
	return k.SetUserBalanceInStore(ctx, balanceKey, balance)
}
