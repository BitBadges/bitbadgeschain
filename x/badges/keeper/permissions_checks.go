package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"math"
)

//Assumes that these parameters are obtained from GetPotentialUpdatesForTimelineValues
func GetUpdateCombinationsToCheck(
	firstMatchesForOld []*types.UniversalPermissionDetails, 
	firstMatchesForNew []*types.UniversalPermissionDetails, 
	emptyValue interface{},
	compareAndGetUpdateCombosToCheck func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails,
) []*types.UniversalPermissionDetails {
	detailsToCheck := []*types.UniversalPermissionDetails{}

	overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
	for _, detail := range inOldButNotNew {
		detailsToAdd := compareAndGetUpdateCombosToCheck(detail.ArbitraryValue, emptyValue)
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime: detail.TimelineTime,
				BadgeId: detailToAdd.BadgeId,
				TransferTime: detailToAdd.TransferTime,
				ToMappingId: detailToAdd.ToMappingId,
				FromMappingId: detailToAdd.FromMappingId,
				InitiatedByMappingId: detailToAdd.InitiatedByMappingId,
			})
		}
	}

	for _, detail := range inNewButNotOld {
		detailsToAdd := compareAndGetUpdateCombosToCheck(detail.ArbitraryValue, emptyValue)
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime: detail.TimelineTime,
				BadgeId: detailToAdd.BadgeId,
				TransferTime: detailToAdd.TransferTime,
				ToMappingId: detailToAdd.ToMappingId,
				FromMappingId: detailToAdd.FromMappingId,
				InitiatedByMappingId: detailToAdd.InitiatedByMappingId,
			})
		}
	}

	for _, overlapObj := range overlapObjects {
		overlap := overlapObj.Overlap
		oldDetails := overlapObj.FirstDetails
		newDetails := overlapObj.SecondDetails
		detailsToAdd := compareAndGetUpdateCombosToCheck(oldDetails.ArbitraryValue, newDetails.ArbitraryValue)
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime: overlap.TimelineTime,
				BadgeId: detailToAdd.BadgeId,
				TransferTime: detailToAdd.TransferTime,
				ToMappingId: detailToAdd.ToMappingId,
				FromMappingId: detailToAdd.FromMappingId,
				InitiatedByMappingId: detailToAdd.InitiatedByMappingId,
			})
		}
	}

	return detailsToCheck
}

//Returns all combinations of timeline times and values
func GetPotentialUpdatesForTimelineValues(times [][]*types.IdRange, values []interface{}) []*types.UniversalPermissionDetails {
	castedPermissions := []*types.UniversalPermission{}
	for idx, time := range times {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				TimelineTimes: time, 
				ArbitraryValue: values[idx],
				UsesTimelineTimes: true,
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}

	firstMatches := types.GetFirstMatchOnly(castedPermissions)
	
	return firstMatches
}

func CheckNotForbidden(ctx sdk.Context, permission *types.UniversalPermissionDetails) error {
	//Throw if we are in a forbidden time
	blockTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	found := types.SearchIdRangesForId(blockTime, permission.ForbiddenTimes)
	if found {
		return ErrInvalidPermissions
	}

	return nil
}

func CheckActionPermission(ctx sdk.Context, permissions []*types.ActionPermission) error {
	castedPermissions := types.CastActionPermissionToUniversalPermission(permissions)
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)
	
	//In this case we only care about the first match since we have no extra criteria
	for _, permissionDetail := range permissionDetails {
		err := CheckNotForbidden(ctx, permissionDetail)
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckTimedUpdatePermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdatePermission) error {
	castedPermissions := types.CastTimedUpdatePermissionToUniversalPermission(permissions)
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func CheckActionWithBadgeIdsAndTimesPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.ActionWithBadgeIdsAndTimesPermission) error {
	castedPermissions := types.CastActionWithBadgeIdsAndTimesPermissionToUniversalPermission(permissions)
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func CheckTimedUpdateWithBadgeIdsPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdateWithBadgeIdsPermission) error {
	castedPermissions := types.CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(permissions)
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func CheckCollectionApprovedTransferPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.CollectionApprovedTransferPermission) error {
	castedPermissions := types.CastCollectionApprovedTransferPermissionToUniversalPermission(permissions)
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}


func CheckUserApprovedTransferPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.UserApprovedTransferPermission) error {
	castedPermissions := types.CastUserApprovedTransferPermissionToUniversalPermission(permissions)
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func CheckNotForbiddenForAllOverlaps(ctx sdk.Context, permissionDetails []*types.UniversalPermissionDetails, detailsToCheck []*types.UniversalPermissionDetails) error {
	//Apply dummy ranges to all detailsToCheck
	for _, detailToCheck := range detailsToCheck {
		if detailToCheck.BadgeId == nil {
			detailToCheck.BadgeId = &types.IdRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) } //dummy range
		}

		if detailToCheck.TimelineTime == nil {
			detailToCheck.TimelineTime = &types.IdRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) } //dummy range
		}

		if detailToCheck.TransferTime == nil {
			detailToCheck.TransferTime = &types.IdRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) } //dummy range
		}

		//Note we are okay with the mapping IDs being "" because they will equal each other
	}
	
	
	//Validate that for each updated timeline time, the current time is permitted
	//We iterate through all explicitly defined permissions (permissionDetails)
	//If we find a match for some timeline time, we check that the current time is not forbidden
	for _, permissionDetail := range permissionDetails {
		for _, detailToCheck := range detailsToCheck {
			_, overlap := types.UniversalRemoveOverlaps(permissionDetail, detailToCheck)
			if len(overlap) > 0 {
				err := CheckNotForbidden(ctx, permissionDetail)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}