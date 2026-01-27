package keeper

import (
	sdkmath "cosmossdk.io/math"
	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// twofaCollectionService adapts the twofa keeper's badgesKeeper to the CollectionService interface
type twofaCollectionService struct {
	keeper Keeper
}

// GetCollection implements CollectionService interface
func (s *twofaCollectionService) GetCollection(ctx sdk.Context, collectionId sdkmath.Uint) (*badgestypes.TokenCollection, bool) {
	return s.keeper.badgesKeeper.GetCollectionFromStore(ctx, collectionId)
}

// GetBalanceOrApplyDefault implements CollectionService interface
func (s *twofaCollectionService) GetBalanceOrApplyDefault(ctx sdk.Context, collection *badgestypes.TokenCollection, userAddress string) (*badgestypes.UserBalanceStore, bool, error) {
	return s.keeper.badgesKeeper.GetBalanceOrApplyDefault(ctx, collection, userAddress)
}

// twofaDynamicStoreService adapts the twofa keeper's badgesKeeper to the DynamicStoreService interface
type twofaDynamicStoreService struct {
	keeper Keeper
}

// GetDynamicStore implements DynamicStoreService interface
func (s *twofaDynamicStoreService) GetDynamicStore(ctx sdk.Context, storeId sdkmath.Uint) (*badgestypes.DynamicStore, bool) {
	store, found := s.keeper.badgesKeeper.GetDynamicStoreFromStore(ctx, storeId)
	if !found {
		return nil, false
	}
	return &store, true
}

// GetDynamicStoreValue implements DynamicStoreService interface
func (s *twofaDynamicStoreService) GetDynamicStoreValue(ctx sdk.Context, storeId sdkmath.Uint, address string) (*badgestypes.DynamicStoreValue, bool) {
	value, found := s.keeper.badgesKeeper.GetDynamicStoreValueFromStore(ctx, storeId, address)
	if !found {
		return nil, false
	}
	return &value, true
}

