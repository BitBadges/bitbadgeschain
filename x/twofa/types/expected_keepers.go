package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	badgestypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// BadgesKeeper defines the expected interface for the badges keeper
// This allows the twofa module to interact with badges without tight coupling
type BadgesKeeper interface {
	GetCollectionFromStore(ctx sdk.Context, collectionId sdkmath.Uint) (*badgestypes.TokenCollection, bool)
	GetBalanceOrApplyDefault(ctx sdk.Context, collection *badgestypes.TokenCollection, address string) (*badgestypes.UserBalanceStore, bool, error)
	GetDynamicStoreFromStore(ctx sdk.Context, storeId sdkmath.Uint) (badgestypes.DynamicStore, bool)
	GetDynamicStoreValueFromStore(ctx sdk.Context, storeId sdkmath.Uint, address string) (*badgestypes.DynamicStoreValue, bool)
}

