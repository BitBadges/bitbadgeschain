package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// SaveBadge defines a method for creating a new badge class
func (k Keeper) SetBadge(ctx sdk.Context, badge types.BitBadge) error {
	if k.HasBadge(ctx, badge.Id) {
		return sdkerrors.Wrap(ErrBadgeExists, badge.Id)
	}
	marshaled_badge, err := k.cdc.Marshal(&badge)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.BitBadge failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(badgeStoreKey(badge.Id), marshaled_badge)
	return nil
}

// UpdateBadge defines a method for updating an existing badge
func (k Keeper) UpdateBadge(ctx sdk.Context, badge types.BitBadge) error {
	if !k.HasBadge(ctx, badge.Id) {
		return sdkerrors.Wrap(ErrBadgeNotExists, badge.Id)
	}
	marshaled_badge, err := k.cdc.Marshal(&badge)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.BitBadge failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(badgeStoreKey(badge.Id), marshaled_badge)
	return nil
}

// GetBadge defines a method for returning the badge information of the specified id
func (k Keeper) GetBadge(ctx sdk.Context, badgeID string) (types.BitBadge, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_badge := store.Get(badgeStoreKey(badgeID))

	var badge types.BitBadge
	if len(marshaled_badge) == 0 {
		return badge, false
	}
	k.cdc.MustUnmarshal(marshaled_badge, &badge)
	return badge, true
}

// GetBadges defines a method for returning all badges information
func (k Keeper) GetBadges(ctx sdk.Context) (badges []*types.BitBadge) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, BadgeKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var badge types.BitBadge
		k.cdc.MustUnmarshal(iterator.Value(), &badge)
		badges = append(badges, &badge)
	}
	return
}

// HasBadge determines whether the specified badgeID exists
func (k Keeper) HasBadge(ctx sdk.Context, badgeID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(badgeStoreKey(badgeID))
}
