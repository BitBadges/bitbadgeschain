package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"math"
)

//Assumes that these parameters are obtained from GetPotentialUpdatesForTimelineValues
func GetUpdateCombinationsToCheck(
	ctx sdk.Context,
	firstMatchesForOld []*types.UniversalPermissionDetails, 
	firstMatchesForNew []*types.UniversalPermissionDetails, 
	emptyValue interface{},
	managerAddress string,
	compareAndGetUpdateCombosToCheck func(ctx sdk.Context, oldValue interface{}, newValue interface{}, managerAddress string) ([]*types.UniversalPermissionDetails, error),
) ([]*types.UniversalPermissionDetails, error) {
	detailsToCheck := []*types.UniversalPermissionDetails{}

	overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)
	for _, detail := range inOldButNotNew {
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, detail.ArbitraryValue, emptyValue, managerAddress)
		if err != nil {
			return nil, err 
		}
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime: detail.TimelineTime,
				BadgeId: detailToAdd.BadgeId,
				TransferTime: detailToAdd.TransferTime,
				ToMapping: detailToAdd.ToMapping,
				FromMapping: detailToAdd.FromMapping,
				InitiatedByMapping: detailToAdd.InitiatedByMapping,
			})
		}
	}

	for _, detail := range inNewButNotOld {
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, detail.ArbitraryValue, emptyValue, managerAddress)
		if err != nil {
			return nil, err 
		}
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime: detail.TimelineTime,
				BadgeId: detailToAdd.BadgeId,
				TransferTime: detailToAdd.TransferTime,
				ToMapping: detailToAdd.ToMapping,
				FromMapping: detailToAdd.FromMapping,
				InitiatedByMapping: detailToAdd.InitiatedByMapping,
			})
		}
	}

	for _, overlapObj := range overlapObjects {
		overlap := overlapObj.Overlap
		oldDetails := overlapObj.FirstDetails
		newDetails := overlapObj.SecondDetails
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, oldDetails.ArbitraryValue, newDetails.ArbitraryValue, managerAddress)
		if err != nil {
			return nil, err
		}
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime: overlap.TimelineTime,
				BadgeId: detailToAdd.BadgeId,
				TransferTime: detailToAdd.TransferTime,
				ToMapping: detailToAdd.ToMapping,
				FromMapping: detailToAdd.FromMapping,
				InitiatedByMapping: detailToAdd.InitiatedByMapping,
			})
		}
	}

	return detailsToCheck, nil
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

func (k Keeper) CheckActionPermission(ctx sdk.Context, permissions []*types.ActionPermission) error {
	castedPermissions, err := k.CastActionPermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}
	
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

func (k Keeper) CheckTimedUpdatePermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdatePermission) error {
	castedPermissions, err := k.CastTimedUpdatePermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}
	
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func (k Keeper) CheckActionWithBadgeIdsAndTimesPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.ActionWithBadgeIdsAndTimesPermission) error {
	castedPermissions, err := k.CastActionWithBadgeIdsAndTimesPermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}
	
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func (k Keeper) CheckTimedUpdateWithBadgeIdsPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdateWithBadgeIdsPermission) error {
	castedPermissions, err := k.CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}
	
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func (k Keeper) CheckCollectionApprovedTransferPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.CollectionApprovedTransferPermission, managerAddress string) error {
	castedPermissions, err := k.CastCollectionApprovedTransferPermissionToUniversalPermission(ctx, managerAddress, permissions)
	if err != nil {
		return err
	}
	
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}


func (k Keeper) CheckUserApprovedOutgoingTransferPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.UserApprovedOutgoingTransferPermission, managerAddress string) error {
	castedPermissions, err := k.CastUserApprovedOutgoingTransferPermissionToUniversalPermission(ctx, managerAddress, permissions)
	if err != nil {
		return err
	}
	
	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func (k Keeper) CheckUserApprovedIncomingTransferPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.UserApprovedIncomingTransferPermission, managerAddress string) error {
	castedPermissions, err := k.CastUserApprovedIncomingTransferPermissionToUniversalPermission(ctx, managerAddress, permissions)
	if err != nil {
		return err
	}
	
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

		if detailToCheck.ToMapping == nil {
			detailToCheck.ToMapping = &types.AddressMapping{Addresses: []string{}, IncludeOnlySpecified: false}
		}

		if detailToCheck.FromMapping == nil {
			detailToCheck.FromMapping = &types.AddressMapping{Addresses: []string{}, IncludeOnlySpecified: false}
		}

		if detailToCheck.InitiatedByMapping == nil {
			detailToCheck.InitiatedByMapping = &types.AddressMapping{Addresses: []string{}, IncludeOnlySpecified: false}
		}
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