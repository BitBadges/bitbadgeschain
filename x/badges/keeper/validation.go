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
	CanArchive                           bool
	CanUpdateContractAddress             bool
	CanCreateMoreBadges                  bool
	CanUpdateCollectionApprovedTransfers bool
	CanUpdateManager              bool
	CanUpdateBadgeMetadata               bool
	CanUpdateCollectionMetadata          bool
	CanUpdateCustomData                  bool
	CanDeleteCollection                  bool
	CanUpdateOffChainBalancesMetadata    bool
}

func CheckPermission(ctx sdk.Context, permission *types.Permission) bool {
	if !permission.IsFrozen {
		return true
	}

	for _, interval := range permission.TimeIntervals {
		time := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		if time.GT(interval.Start) && time.LT(interval.End) {
			return true
		}
	}

	return false
}

// Validates everything about the Msg is valid and returns (creatorNum, badge, permissions, error).
func (k Keeper) UniversalValidate(ctx sdk.Context, params UniversalValidationParams) (types.BadgeCollection, error) {

	if len(params.AccountsThatCantEqualCreator) > 0 {
		for _, account := range params.AccountsThatCantEqualCreator {
			if account == params.Creator {
				return types.BadgeCollection{}, ErrAccountCanNotEqualCreator
			}
		}
	}

	// Assert collection and badgeId ranges exist and are well-formed
	badge, err := k.GetCollectionAndAssertBadgeIdsAreValid(ctx, params.CollectionId, params.BadgeIdRangesToValidate)
	if err != nil {
		return types.BadgeCollection{}, err
	}

	if badge.IsArchived {
		return types.BadgeCollection{}, ErrCollectionIsArchived
	}

	// Assert all permissions
	if params.MustBeManager && badge.Manager != params.Creator {
		return types.BadgeCollection{}, ErrSenderIsNotManager
	}

	permissions := badge.Permissions
	if params.CanUpdateCollectionApprovedTransfers && !CheckPermission(ctx, permissions.CanUpdateCollectionApprovedTransfers) {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanCreateMoreBadges && !CheckPermission(ctx, permissions.CanCreateMoreBadges) {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanUpdateManager && !CheckPermission(ctx, permissions.CanUpdateManager) {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanUpdateBadgeMetadata && !CheckPermission(ctx, permissions.CanUpdateBadgeMetadata) {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanUpdateCollectionMetadata && !CheckPermission(ctx, permissions.CanUpdateCollectionMetadata) {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanUpdateCustomData && !CheckPermission(ctx, permissions.CanUpdateCustomData) {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanDeleteCollection && !CheckPermission(ctx, permissions.CanDeleteCollection) {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanUpdateOffChainBalancesMetadata && !CheckPermission(ctx, permissions.CanUpdateOffChainBalancesMetadata) {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	if params.CanUpdateContractAddress && !CheckPermission(ctx, permissions.CanUpdateContractAddress) {
		return types.BadgeCollection{}, ErrInvalidPermissions
	}

	return badge, nil
}
