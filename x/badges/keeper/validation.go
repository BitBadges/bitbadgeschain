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

	isArchived := GetIsArchived(ctx, collection)
	if isArchived {
		return types.BadgeCollection{}, ErrCollectionIsArchived
	}

	// Assert all permissions
	if params.MustBeManager {
		currManager := GetCurrentManager(ctx, collection)
		if currManager != params.Creator {
			return types.BadgeCollection{}, ErrSenderIsNotManager
		}
	}

	return collection, nil
}
