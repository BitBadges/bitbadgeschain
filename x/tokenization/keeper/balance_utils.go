package keeper

import (
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	"github.com/cosmos/gogoproto/proto"

	sdkerrors "cosmossdk.io/errors"
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

// GetBalanceOrApplyDefaultReadOnly retrieves a user balance or constructs the collection default
// balance without writing approval version records to state. Callers that only inspect balances or
// delegate persistence to a subsequent write helper should use this variant.
func (k Keeper) GetBalanceOrApplyDefaultReadOnly(ctx sdk.Context, collection *types.TokenCollection, userAddress string) (*types.UserBalanceStore, bool, error) {
	return k.getBalanceOrApplyDefault(ctx, collection, userAddress, false)
}

// GetBalanceOrApplyDefault retrieves user balance or applies default balance store.
// When initializing defaults for a state-writing flow, approval versions are recorded in store.
func (k Keeper) GetBalanceOrApplyDefault(ctx sdk.Context, collection *types.TokenCollection, userAddress string) (*types.UserBalanceStore, bool, error) {
	return k.getBalanceOrApplyDefault(ctx, collection, userAddress, true)
}

func (k Keeper) getBalanceOrApplyDefault(ctx sdk.Context, collection *types.TokenCollection, userAddress string, persistDefaultVersions bool) (*types.UserBalanceStore, bool, error) {
	// Mint has unlimited balances
	if types.IsMintAddress(userAddress) {
		return &types.UserBalanceStore{}, false, nil
	}

	// Total address tracks minted supply — read from store or return empty.
	// Never inherits DefaultBalances since it's a protocol accumulator, not a user.
	if types.IsTotalAddress(userAddress) {
		balanceKey := ConstructBalanceKey(userAddress, collection.CollectionId)
		balance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
		if !found {
			return &types.UserBalanceStore{}, false, nil
		}
		return balance, false, nil
	}

	// Special backed addresses also have unlimited balances
	if k.IsSpecialBackedAddress(ctx, collection, userAddress) {
		return &types.UserBalanceStore{
			Balances:                                  []*types.Balance{},
			AutoApproveAllIncomingTransfers:           true,
			AutoApproveSelfInitiatedIncomingTransfers: true,
			AutoApproveSelfInitiatedOutgoingTransfers: true,
		}, false, nil
	}

	// We get current balances or fallback to default balances
	balanceKey := ConstructBalanceKey(userAddress, collection.CollectionId)
	balance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
	appliedDefault := false
	if !found {
		balance = getDefaultBalanceStoreForCollection(collection)
		appliedDefault = true

		if persistDefaultVersions {
			for _, approval := range balance.IncomingApprovals {
				version, err := k.IncrementApprovalVersion(ctx, collection.CollectionId, "incoming", userAddress, approval.ApprovalId)
				if err != nil {
					return nil, false, sdkerrors.Wrap(err, "failed to increment approval version")
				}
				approval.Version = version
			}
			for _, approval := range balance.OutgoingApprovals {
				version, err := k.IncrementApprovalVersion(ctx, collection.CollectionId, "outgoing", userAddress, approval.ApprovalId)
				if err != nil {
					return nil, false, sdkerrors.Wrap(err, "failed to increment approval version")
				}
				approval.Version = version
			}
		} else {
			for _, approval := range balance.IncomingApprovals {
				vKey := ConstructApprovalVersionKey(collection.CollectionId, "incoming", userAddress, approval.ApprovalId)
				version, _ := k.GetApprovalTrackerVersionFromStore(ctx, vKey)
				approval.Version = version
			}
			for _, approval := range balance.OutgoingApprovals {
				vKey := ConstructApprovalVersionKey(collection.CollectionId, "outgoing", userAddress, approval.ApprovalId)
				version, _ := k.GetApprovalTrackerVersionFromStore(ctx, vKey)
				approval.Version = version
			}
		}
	}

	if balance.UserPermissions == nil {
		balance.UserPermissions = &types.UserPermissions{}
	}

	return balance, appliedDefault, nil
}

// SetBalanceForAddress stores a user balance for a specific address.
// Updates holder count for the collection when balance crosses zero.
func (k Keeper) SetBalanceForAddress(ctx sdk.Context, collection *types.TokenCollection, userAddress string, balance *types.UserBalanceStore) error {
	balanceKey := ConstructBalanceKey(userAddress, collection.CollectionId)

	// Get old balance for holder count tracking
	oldBalance, _ := k.GetUserBalanceFromStore(ctx, balanceKey)

	// Update holder count BEFORE setting new balance
	if err := k.UpdateHolderCount(ctx, collection, userAddress, oldBalance, balance); err != nil {
		return err
	}

	return k.SetUserBalanceInStore(ctx, balanceKey, balance, false)
}
