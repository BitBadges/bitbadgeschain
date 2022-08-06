package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k Keeper) AssertBadgeAndSubBadgeExistsAndReturnBadge(ctx sdk.Context, badge_id uint64, subbadgeRanges []*types.NumberRange) (types.BitBadge, error) {
	badge, found := k.GetBadgeFromStore(ctx, badge_id)
	if !found {
		return types.BitBadge{}, ErrBadgeNotExists
	}

	for _, subbadgeRange := range subbadgeRanges {
		if subbadgeRange.End >= badge.NextSubassetId {
			return types.BitBadge{}, ErrSubBadgeNotExists
		}
	}

	return badge, nil
}
