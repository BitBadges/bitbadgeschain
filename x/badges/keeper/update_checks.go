package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"math"
)

//See top of update_checks_helpers.go for documentation

// Precondition: Assumes that all passed-in matches have .ArbitraryValue set
//
//	We do this by getting from GetPotentialUpdatesForTimelineValues
//
// This is a generic function that is used to get the "updated" field combinations
// Ex: If we go from [badgeIDs 1 to 10 -> www.example.com] to [badgeIDs 1 to 2 -> www.example2.com, badgeIDs 3 to 10 -> www.example.com]
//
//	This will return a UniversalPermissionDetails with badgeIDs 1 to 2 because they are the only ones that changed
//
// Note that updates are field-specific, so the comparison logic is handled via a custom passed-in function - compareAndGetUpdateCombosToCheck
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

	//Handle all old combinations that are not in the new (by comparing to empty value)
	for _, detail := range inOldButNotNew {
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, detail.ArbitraryValue, emptyValue, managerAddress)
		if err != nil {
			return nil, err
		}
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime:       detail.TimelineTime,
				BadgeId:            detailToAdd.BadgeId,
				TransferTime:       detailToAdd.TransferTime,
				ToMapping:          detailToAdd.ToMapping,
				FromMapping:        detailToAdd.FromMapping,
				InitiatedByMapping: detailToAdd.InitiatedByMapping,
			})
		}
	}

	//Handle all new combinations that are not in the old (by comparing to empty value)
	for _, detail := range inNewButNotOld {
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, detail.ArbitraryValue, emptyValue, managerAddress)
		if err != nil {
			return nil, err
		}
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime:       detail.TimelineTime,
				BadgeId:            detailToAdd.BadgeId,
				TransferTime:       detailToAdd.TransferTime,
				ToMapping:          detailToAdd.ToMapping,
				FromMapping:        detailToAdd.FromMapping,
				InitiatedByMapping: detailToAdd.InitiatedByMapping,
			})
		}
	}

	//Handle all overlaps (by comparing old and new values directly)
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
				TimelineTime:       overlap.TimelineTime,
				BadgeId:            detailToAdd.BadgeId,
				TransferTime:       detailToAdd.TransferTime,
				ToMapping:          detailToAdd.ToMapping,
				FromMapping:        detailToAdd.FromMapping,
				InitiatedByMapping: detailToAdd.InitiatedByMapping,
			})
		}
	}

	return detailsToCheck, nil
}

func CheckNotForbidden(ctx sdk.Context, permission *types.UniversalPermissionDetails) error {
	//Throw if current block time is a forbidden time
	blockTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	found := types.SearchUintRangesForUint(blockTime, permission.ForbiddenTimes)
	if found {
		return sdkerrors.Wrapf(ErrForbiddenTime, "current time %s is forbidden", blockTime.String())
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

func (k Keeper) CheckBalancesActionPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.BalancesActionPermission) error {
	castedPermissions, err := k.CastBalancesActionPermissionToUniversalPermission(permissions)
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
			detailToCheck.BadgeId = &types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)} //dummy range
		}

		if detailToCheck.TimelineTime == nil {
			detailToCheck.TimelineTime = &types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)} //dummy range
		}

		if detailToCheck.TransferTime == nil {
			detailToCheck.TransferTime = &types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)} //dummy range
		}

		if detailToCheck.ToMapping == nil {
			detailToCheck.ToMapping = &types.AddressMapping{Addresses: []string{}, IncludeAddresses: false}
		}

		if detailToCheck.FromMapping == nil {
			detailToCheck.FromMapping = &types.AddressMapping{Addresses: []string{}, IncludeAddresses: false}
		}

		if detailToCheck.InitiatedByMapping == nil {
			detailToCheck.InitiatedByMapping = &types.AddressMapping{Addresses: []string{}, IncludeAddresses: false}
		}
	}

	//Validate that for each updated timeline time, the current time is permitted
	//We iterate through all explicitly defined permissions (permissionDetails)
	//If we find a match for some timeline time, we check that the current time is not forbidden
	for _, permissionDetail := range permissionDetails {
		for _, detailToCheck := range detailsToCheck {
			_, overlap := types.UniversalRemoveOverlaps(permissionDetail, detailToCheck)
			if len(overlap) > 0 {
				err := CheckNotForbidden(ctx, permissionDetail) //forbiddenTimes and permittedTimes are stored in here
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
