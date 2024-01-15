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
	compareAndGetUpdateCombosToCheck func(ctx sdk.Context, oldValue interface{}, newValue interface{}) ([]*types.UniversalPermissionDetails, error),
) ([]*types.UniversalPermissionDetails, error) {
	detailsToCheck := []*types.UniversalPermissionDetails{}

	overlapObjects, inOldButNotNew, inNewButNotOld := types.GetOverlapsAndNonOverlaps(firstMatchesForOld, firstMatchesForNew)

	//Handle all old combinations that are not in the new (by comparing to empty value)
	for _, detail := range inOldButNotNew {
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, detail.ArbitraryValue, emptyValue)
		if err != nil {
			return nil, err
		}
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime:              detail.TimelineTime,
				BadgeId:                   detailToAdd.BadgeId,
				TransferTime:              detailToAdd.TransferTime,
				OwnershipTime:             detailToAdd.OwnershipTime,
				ToList:                 detailToAdd.ToList,
				FromList:               detailToAdd.FromList,
				InitiatedByList:        detailToAdd.InitiatedByList,
				AmountTrackerIdList:    detailToAdd.AmountTrackerIdList,
				ChallengeTrackerIdList: detailToAdd.ChallengeTrackerIdList,
			})
		}
	}

	//Handle all new combinations that are not in the old (by comparing to empty value)
	for _, detail := range inNewButNotOld {
		detailsToAdd, err := compareAndGetUpdateCombosToCheck(ctx, detail.ArbitraryValue, emptyValue)
		if err != nil {
			return nil, err
		}
		for _, detailToAdd := range detailsToAdd {
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime:              detail.TimelineTime,
				BadgeId:                   detailToAdd.BadgeId,
				TransferTime:              detailToAdd.TransferTime,
				OwnershipTime:             detailToAdd.OwnershipTime,
				ToList:                 detailToAdd.ToList,
				FromList:               detailToAdd.FromList,
				InitiatedByList:        detailToAdd.InitiatedByList,
				AmountTrackerIdList:    detailToAdd.AmountTrackerIdList,
				ChallengeTrackerIdList: detailToAdd.ChallengeTrackerIdList,
			})
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
			detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
				TimelineTime:              overlap.TimelineTime,
				BadgeId:                   detailToAdd.BadgeId,
				TransferTime:              detailToAdd.TransferTime,
				OwnershipTime:             detailToAdd.OwnershipTime,
				ToList:                 detailToAdd.ToList,
				FromList:               detailToAdd.FromList,
				InitiatedByList:        detailToAdd.InitiatedByList,
				AmountTrackerIdList:    detailToAdd.AmountTrackerIdList,
				ChallengeTrackerIdList: detailToAdd.ChallengeTrackerIdList,
			})
		}
	}

	return detailsToCheck, nil
}

func CheckNotForbidden(ctx sdk.Context, permission *types.UniversalPermissionDetails, permissionStr string) error {
	blockTime := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	found, err := types.SearchUintRangesForUint(blockTime, permission.PermanentlyForbiddenTimes)
	if found || err != nil {
		return sdkerrors.Wrapf(ErrForbiddenTime, "current time %s is forbidden for permission %s", blockTime.String(), permissionStr)
	}

	return nil
}

func (k Keeper) CheckIfActionPermissionPermits(ctx sdk.Context, permissions []*types.ActionPermission, permissionStr string) error {
	castedPermissions, err := k.CastActionPermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	//In this case we only care about the first match since we have no extra criteria
	for _, permissionDetail := range permissionDetails {
		err := CheckNotForbidden(ctx, permissionDetail, permissionStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) CheckIfTimedUpdatePermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdatePermission, permissionStr string) error {
	castedPermissions, err := k.CastTimedUpdatePermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfBalancesActionPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.BalancesActionPermission, permissionStr string) error {
	castedPermissions, err := k.CastBalancesActionPermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfTimedUpdateWithBadgeIdsPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdateWithBadgeIdsPermission, permissionStr string) error {
	castedPermissions, err := k.CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(permissions)
	if err != nil {
		return err
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfCollectionApprovalPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.CollectionApprovalPermission, permissionStr string) error {
	castedPermissions, err := k.CastCollectionApprovalPermissionToUniversalPermission(ctx, permissions)
	if err != nil {
		return err
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfUserOutgoingApprovalPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.UserOutgoingApprovalPermission, permissionStr string) error {
	castedPermissions, err := k.CastUserOutgoingApprovalPermissionToUniversalPermission(ctx, permissions)
	if err != nil {
		return err
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck, permissionStr)
}

func (k Keeper) CheckIfUserIncomingApprovalPermissionPermits(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.UserIncomingApprovalPermission, permissionStr string) error {
	castedPermissions, err := k.CastUserIncomingApprovalPermissionToUniversalPermission(ctx, permissions)
	if err != nil {
		return err
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck, permissionStr)
}

func CheckNotForbiddenForAllOverlaps(ctx sdk.Context, permissionDetails []*types.UniversalPermissionDetails, detailsToCheck []*types.UniversalPermissionDetails, permissionStr string) error {
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

		if detailToCheck.OwnershipTime == nil {
			detailToCheck.OwnershipTime = &types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)} //dummy range
		}

		if detailToCheck.AmountTrackerIdList == nil {
			detailToCheck.AmountTrackerIdList = &types.AddressList{Addresses: []string{}, Whitelist: false}
		}

		if detailToCheck.ChallengeTrackerIdList == nil {
			detailToCheck.ChallengeTrackerIdList = &types.AddressList{Addresses: []string{}, Whitelist: false}
		}

		if detailToCheck.ToList == nil {
			detailToCheck.ToList = &types.AddressList{Addresses: []string{}, Whitelist: false}
		}

		if detailToCheck.FromList == nil {
			detailToCheck.FromList = &types.AddressList{Addresses: []string{}, Whitelist: false}
		}

		if detailToCheck.InitiatedByList == nil {
			detailToCheck.InitiatedByList = &types.AddressList{Addresses: []string{}, Whitelist: false}
		}
	}

	//Validate that for each updated timeline time, the current time is permitted
	//We iterate through all explicitly defined permissions (permissionDetails)
	//If we find a match for some timeline time, we check that the current time is not forbidden
	for _, permissionDetail := range permissionDetails {
		for _, detailToCheck := range detailsToCheck {
			_, overlap := types.UniversalRemoveOverlaps(permissionDetail, detailToCheck)
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
