package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// SaveBadge defines a method for creating a new badge class
func (k Keeper) SetBadgeInStore(ctx sdk.Context, badge types.BitBadge) error {
	if k.StoreHasBadgeID(ctx, badge.Id) {
		return sdkerrors.Wrap(ErrBadgeExists, strconv.FormatUint(badge.Id, 10))
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
func (k Keeper) UpdateBadgeInStore(ctx sdk.Context, badge types.BitBadge) error {
	if !k.StoreHasBadgeID(ctx, badge.Id) {
		return sdkerrors.Wrap(ErrBadgeNotExists, strconv.FormatUint(badge.Id, 10))
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
func (k Keeper) GetBadgeFromStore(ctx sdk.Context, badgeID uint64) (types.BitBadge, bool) {
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
func (k Keeper) GetBadgesFromStore(ctx sdk.Context) (badges []*types.BitBadge) {
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
func (k Keeper) StoreHasBadgeID(ctx sdk.Context, badgeID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(badgeStoreKey(badgeID))
}

func (k Keeper) HasAddressRequestedManagerTransfer(ctx sdk.Context, badgeId uint64, address uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := GetManagerRequestKey(badgeId, address)
	return store.Has(managerTransferRequestKey(key))
}

func (k Keeper) CreateTransferManagerRequest(ctx sdk.Context, badgeId uint64, address uint64) error {
	request := []byte{}
	store := ctx.KVStore(k.storeKey)
	key := GetManagerRequestKey(badgeId, address)
	store.Set(managerTransferRequestKey(key), request)
	return nil
}

func (k Keeper) RemoveTransferManagerRequest(ctx sdk.Context, badgeId uint64, address uint64) error {
	key := GetManagerRequestKey(badgeId, address)
	store := ctx.KVStore(k.storeKey)
	// store.Set(managerTransferRequestKey(key), request)

	if store.Has(managerTransferRequestKey(key)) {
		store.Delete(managerTransferRequestKey(key))
	}
	return nil
}
