package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Determines what to validate for each Msg
type UniversalValidationParams struct {
	Creator                              string
	CollectionId                         sdk.Uint
	AccountsThatCantEqualCreator         []string
	BadgeIdRangesToValidate              []*types.IdRange
	MustBeManager                        bool
}

// Validates everything about the Msg is valid and returns (creatorNum, collection, permissions, error).
func (k Keeper) UniversalValidate(ctx sdk.Context, params UniversalValidationParams) (types.BadgeCollection, error) {
	blockTime := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	if len(params.AccountsThatCantEqualCreator) > 0 {
		for _, account := range params.AccountsThatCantEqualCreator {
			if account == params.Creator {
				return types.BadgeCollection{}, ErrAccountCanNotEqualCreator
			}
		}
	}

	// Assert collection and badgeId ranges exist and are well-formed
	collection, err := k.GetCollectionAndAssertBadgeIdsAreValid(ctx, params.CollectionId, params.BadgeIdRangesToValidate)
	if err != nil {
		return types.BadgeCollection{}, err
	}


	for _, timelineVal := range collection.IsArchivedTimeline {
		idx, found := types.SearchIdRangesForId(blockTime, timelineVal.Times)
		if found {
			if timelineVal.IsArchived { 
				return types.BadgeCollection{}, ErrCollectionIsArchived
			} 
			break
		} 
	}

	// Assert all permissions
	if params.MustBeManager {
		for _, timelineVal := range collection.ManagerTimeline {
			idx, found := types.SearchIdRangesForId(blockTime, timelineVal.Times)
			if found {
				if timelineVal.Manager != params.Creator {
					return types.BadgeCollection{}, ErrSenderIsNotManager
				}
				break
			}
		}
	}

	return collection, nil
}
