package keeper

import (
	"context"
	"encoding/binary"
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

// Keeper holds references to the abstract validator set and tokenization keepers
// plus its own KV store (for params, compliance-jailed set, and saved power).
type Keeper struct {
	cdc    codec.BinaryCodec
	logger log.Logger

	// storeKey is for persistent state (params + compliance-jailed set + saved power).
	storeKey storetypes.StoreKey

	// authority is the address capable of executing MsgUpdateParams (typically x/gov).
	authority string

	// External keeper interfaces — never concrete types.
	validatorSet       types.ValidatorSetKeeper
	tokenizationKeeper types.TokenizationKeeper
}

// NewKeeper creates a new x/pot keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	logger log.Logger,
	authority string,
	tokenizationKeeper types.TokenizationKeeper,
	validatorSet types.ValidatorSetKeeper,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		logger:             logger,
		authority:          authority,
		validatorSet:       validatorSet,
		tokenizationKeeper: tokenizationKeeper,
	}
}

// GetAuthority returns the module's governance authority address.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// --------------------------------------------------------------------------
// Compliance-jailed validator set (persistent KV store)
// --------------------------------------------------------------------------

// SetComplianceJailed marks a validator (by consensus address) as compliance-jailed
// in x/pot's own state.
func (k Keeper) SetComplianceJailed(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ComplianceJailedKey(consAddr), []byte{1})
}

// RemoveComplianceJailed removes a validator from the compliance-jailed set.
func (k Keeper) RemoveComplianceJailed(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ComplianceJailedKey(consAddr))
}

// IsComplianceJailed returns true if the validator is in x/pot's compliance-jailed set.
func (k Keeper) IsComplianceJailed(ctx sdk.Context, consAddr sdk.ConsAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ComplianceJailedKey(consAddr))
}

// GetAllComplianceJailed returns all consensus addresses in the compliance-jailed set.
func (k Keeper) GetAllComplianceJailed(ctx sdk.Context) [][]byte {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(
		types.ComplianceJailedPrefix,
		storetypes.PrefixEndBytes(types.ComplianceJailedPrefix),
	)
	defer iter.Close()

	var addrs [][]byte
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		addr := make([]byte, len(key)-len(types.ComplianceJailedPrefix))
		copy(addr, key[len(types.ComplianceJailedPrefix):])
		addrs = append(addrs, addr)
	}
	return addrs
}

// --------------------------------------------------------------------------
// Saved power store (for PoA adapter — implements types.PowerStore)
// --------------------------------------------------------------------------

// GetSavedPower returns the saved voting power for a validator, plus a bool
// indicating whether an entry existed.
func (k Keeper) GetSavedPower(ctx context.Context, consAddr sdk.ConsAddress) (int64, bool) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	bz := store.Get(types.SavedPowerKey(consAddr))
	if bz == nil {
		return 0, false
	}
	return int64(binary.BigEndian.Uint64(bz)), true
}

// SetSavedPower persists the saved voting power for a validator.
func (k Keeper) SetSavedPower(ctx context.Context, consAddr sdk.ConsAddress, power int64) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(power))
	store.Set(types.SavedPowerKey(consAddr), bz)
}

// RemoveSavedPower deletes the saved voting power entry for a validator.
func (k Keeper) RemoveSavedPower(ctx context.Context, consAddr sdk.ConsAddress) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := sdkCtx.KVStore(k.storeKey)
	store.Delete(types.SavedPowerKey(consAddr))
}

// Compile-time check that Keeper implements PowerStore.
var _ types.PowerStore = Keeper{}
