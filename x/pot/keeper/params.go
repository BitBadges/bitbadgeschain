package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/pot/types"
)

// GetParams returns the current x/pot params from the persistent store.
// It tries protobuf first, falls back to JSON (old format), and self-heals
// by re-writing as protobuf on JSON success.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return types.DefaultParams()
	}

	var params types.Params
	// Try protobuf first (new format).
	if err := k.cdc.Unmarshal(bz, &params); err == nil {
		return params
	}
	// Fall back to JSON (old format) — self-healing migration.
	if err := json.Unmarshal(bz, &params); err == nil {
		// Re-write as protobuf for future reads.
		newBz, err := k.cdc.Marshal(&params)
		if err == nil {
			store.Set(types.ParamsKey, newBz)
		}
		return params
	}
	// Both failed — return defaults.
	return types.DefaultParams()
}

// SetParams persists the x/pot params to the store using protobuf encoding.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
	return nil
}
