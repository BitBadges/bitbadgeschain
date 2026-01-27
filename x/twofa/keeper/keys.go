package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/twofa/types"
)

const StoreKey = types.ModuleName

var (
	User2FARequirementsKey = []byte{0x01}
)

// user2FARequirementsStoreKey returns the store key for a user's 2FA requirements
func user2FARequirementsStoreKey(address string) []byte {
	key := make([]byte, len(User2FARequirementsKey)+len(address))
	copy(key, User2FARequirementsKey)
	copy(key[len(User2FARequirementsKey):], []byte(address))
	return key
}

// getStore returns a prefix store for the twofa module
func (k Keeper) getStore(ctx sdk.Context) storetypes.KVStore {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return prefix.NewStore(storeAdapter, []byte{})
}

