package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AssertBadgeAndSubBadgeExists(ctx sdk.Context, badge_id uint64, subbadge_id uint64) error {
	badge, found := k.GetBadgeFromStore(ctx, badge_id)
	if !found {
		return ErrBadgeNotExists
	}

	if subbadge_id >= badge.NextSubassetId {
		return ErrSubBadgeNotExists
	}
	return nil
}
