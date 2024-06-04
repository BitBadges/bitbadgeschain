package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Precondition: Assumes that all passed-in matches have .ArbitraryValue set
//
// This is a generic function that is used to get the "updated" field combinations
// Ex: If we go from [badgeIDs 1 to 10 -> www.example.com] to [badgeIDs 1 to 2 -> www.example2.com, badgeIDs 3 to 10 -> www.example.com]
//
// # This will return a UniversalPermissionDetails with badgeIDs 1 to 2 because they changed and are the ones we need to check
//
// Note that updates are field-specific, so the comparison logic is handled via a custom passed-in function - compareAndGetUpdateCombosToCheck
func GetUpdateCombinationsToCheck(
	ctx sdk.Context,
	firstMatchesForOld []*types.UniversalPermissionDetails,
	firstMatchesForNew []*types.UniversalPermissionDetails,
	emptyValue interface{},
	compareAndGetUpdateCombosToCheck func(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error),
) ([]*types.UniversalPermissionDetails, error) {

	overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(ctx, firstMatchesForOld, firstMatchesForNew)

	detailsToCheck := []*types.UniversalPermissionDetails{}
	//Handle all old combinations that are not in the new (by comparing to empty value)
	for _, detail := range inOldButNotNew {
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, detail.ArbitraryValue, emptyValue)
		if err != nil {
			return nil, err
		}

		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, addTimelineTimeToDetails(detailToAdd, detail.TimelineTime))
		}
	}

	//Handle all new combinations that are not in the old (by comparing to empty value)
	for _, detail := range inNewButNotOld {
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, detail.ArbitraryValue, emptyValue)
		if err != nil {
			return nil, err
		}

		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, addTimelineTimeToDetails(detailToAdd, detail.TimelineTime))
		}
	}

	//Handle all overlaps (by comparing old and new values directly)
	for _, overlapObj := range overlapObjects {
		overlap := overlapObj.Overlap
		oldDetails := overlapObj.FirstDetails
		newDetails := overlapObj.SecondDetails
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, oldDetails.ArbitraryValue, newDetails.ArbitraryValue)
		if err != nil {
			return nil, err
		}

		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, addTimelineTimeToDetails(detailToAdd, overlap.TimelineTime))
		}
	}

	return detailsToCheck, nil
}

func addTimelineTimeToDetails(details *types.UniversalPermissionDetails, timelineTime *types.UintRange) *types.UniversalPermissionDetails {
	return &types.UniversalPermissionDetails{
		TimelineTime:    timelineTime,
		BadgeId:         details.BadgeId,
		TransferTime:    details.TransferTime,
		OwnershipTime:   details.OwnershipTime,
		ToList:          details.ToList,
		FromList:        details.FromList,
		InitiatedByList: details.InitiatedByList,
		ApprovalIdList:  details.ApprovalIdList,
	}
}

func (k Keeper) CheckIfActionPermissionPermits(ctx sdk.Context, permissions []*types.ActionPermission, permissionStr string) error {
	castedPermissions, err := k.CastActionPermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}

	return CheckNotForbiddenForAllOverlaps(ctx, castedPermissions, []*types.UniversalPermissionDetails{
		{}, //Little hacky but we just need to check if the current time is permitted for any value (so we check an all dummy value since ActionPermissions have no extra criteria)
	}, permissionStr)
}

func (k Keeper) CheckIfTimedUpdatePermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdatePermission, permissionStr string) error {
	castedPermissions, err := k.CastTimedUpdatePermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}

	return CheckNotForbiddenForAllOverlaps(ctx, castedPermissions, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfBalancesActionPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.BalancesActionPermission, permissionStr string) error {
	castedPermissions, err := k.CastBalancesActionPermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}

	return CheckNotForbiddenForAllOverlaps(ctx, castedPermissions, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfTimedUpdateWithBadgeIdsPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdateWithBadgeIdsPermission, permissionStr string) error {
	castedPermissions, err := k.CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}

	return CheckNotForbiddenForAllOverlaps(ctx, castedPermissions, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfCollectionApprovalPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.CollectionApprovalPermission, permissionStr string) error {
	castedPermissions, err := k.CastCollectionApprovalPermissionToUniversalPermission(ctx, permissions)
	if err != nil {
		return err
	}

	return CheckNotForbiddenForAllOverlaps(ctx, castedPermissions, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfUserOutgoingApprovalPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.UserOutgoingApprovalPermission, permissionStr string) error {
	castedPermissions, err := k.CastUserOutgoingApprovalPermissionToUniversalPermission(ctx, permissions)
	if err != nil {
		return err
	}

	return CheckNotForbiddenForAllOverlaps(ctx, castedPermissions, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfUserIncomingApprovalPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.UserIncomingApprovalPermission, permissionStr string) error {
	castedPermissions, err := k.CastUserIncomingApprovalPermissionToUniversalPermission(ctx, permissions)
	if err != nil {
		return err
	}

	return CheckNotForbiddenForAllOverlaps(ctx, castedPermissions, detailsToCheck, permissionStr)
}

func CheckNotForbiddenForAllOverlaps(ctx sdk.Context, castedPermissions []*types.UniversalPermission, detailsToCheck []*types.UniversalPermissionDetails, permissionStr string) error {
	//Get the permissions first match only and apply dummy values to detailsToCheck if nil or missing
	permissionDetails := types.GetFirstMatchOnly(ctx, castedPermissions)
	for i, detailToCheck := range detailsToCheck {
		detailsToCheck[i] = types.AddDefaultsIfNil(detailToCheck)
	}

	//Validate that for each updated combination, the current time is permitted (unhandled or explicitly permitted)
	//If we find a match for some timeline time, we check that the current time is not explicitly forbidden
	for _, permissionDetail := range permissionDetails {
		for _, detailToCheck := range detailsToCheck {
			_, overlap := types.UniversalRemoveOverlaps(ctx, permissionDetail, detailToCheck)
			if len(overlap) > 0 {
				err := CheckNotForbidden(ctx, permissionDetail, permissionStr) //permanentlyForbiddenTimes and permanentlyPermittedTimes are stored in here
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func CheckNotForbidden(ctx sdk.Context, permission *types.UniversalPermissionDetails, permissionStr string) error {
	blockTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	found, err := types.SearchUintRangesForUint(blockTime, permission.PermanentlyForbiddenTimes)
	if found || err != nil {
		return sdkerrors.Wrapf(ErrForbiddenTime, "current time %s is forbidden for permission %s", blockTime.String(), permissionStr)
	}

	return nil
}
