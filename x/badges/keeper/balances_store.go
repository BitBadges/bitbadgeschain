package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k Keeper) CreateBadgeBalanceInStore(ctx sdk.Context, balance_id string, badgeBalanceInfo types.BadgeBalanceInfo) error {
	if k.StoreHasBadgeBalance(ctx, balance_id) {
		return sdkerrors.Wrap(ErrBadgeBalanceExists, balance_id)
	}

	marshaled_badge_balance_info, err := k.cdc.Marshal(&badgeBalanceInfo)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.BadgeBalanceInfo failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(badgeBalanceStoreKey(balance_id), marshaled_badge_balance_info)
	return nil
}

func (k Keeper) UpdateBadgeBalanceInStore(ctx sdk.Context, balance_id string, badgeBalanceInfo types.BadgeBalanceInfo) error {
	if !k.StoreHasBadgeBalance(ctx, balance_id) {
		return sdkerrors.Wrap(ErrBadgeBalanceNotExists, balance_id)
	}
	marshaled_badge_balance_info, err := k.cdc.Marshal(&badgeBalanceInfo)
	if err != nil {
		return sdkerrors.Wrap(err, "Marshal types.BadgeBalanceInfo failed")
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(badgeBalanceStoreKey(balance_id), marshaled_badge_balance_info)
	return nil
}

func (k Keeper) GetBadgeBalanceFromStore(ctx sdk.Context, balance_id string) (types.BadgeBalanceInfo, bool) {
	store := ctx.KVStore(k.storeKey)
	marshaled_badge_balance_info := store.Get(badgeBalanceStoreKey(balance_id))

	var badgeBalanceInfo types.BadgeBalanceInfo
	if len(marshaled_badge_balance_info) == 0 {
		return badgeBalanceInfo, false
	}
	k.cdc.MustUnmarshal(marshaled_badge_balance_info, &badgeBalanceInfo)
	return badgeBalanceInfo, true
}

func (k Keeper) GetBadgeBalancesFromStore(ctx sdk.Context) (addresses []*types.BadgeBalanceInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, BadgeBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var badgeBalanceInfo types.BadgeBalanceInfo
		k.cdc.MustUnmarshal(iterator.Value(), &badgeBalanceInfo)
		addresses = append(addresses, &badgeBalanceInfo)
	}
	return
}

func (k Keeper) GetBadgeBalanceIdsFromStore(ctx sdk.Context) (ids []string) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, BadgeBalanceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		ids = append(ids, string(iterator.Value()))
	}
	return
}

func (k Keeper) StoreHasBadgeBalance(ctx sdk.Context, balance_id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(badgeBalanceStoreKey(balance_id))
}
