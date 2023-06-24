package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"math"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type UniversalTimelineDetails struct {
	TimelineTime *types.IdRange
	ArbitraryValue interface{}
}

//Assumes that all times are present (0 to math.MaxUint64)
func GetDetailsToCheck(oldTimeline []*UniversalTimelineDetails, newTimeline []*UniversalTimelineDetails, compareFunc func(oldValue interface{}, newValue interface{}) []*types.UniversalPermissionDetails) []*types.UniversalPermissionDetails {
	detailsToCheck := []*types.UniversalPermissionDetails{}
	
	for len(oldTimeline) > 0 {
		timeToCheck := oldTimeline[0].TimelineTime.Start
		
		oldValue := oldTimeline[0].ArbitraryValue
		
		newValue := interface{}(nil)
		newFound := false
		newTimelineTime := &types.IdRange{}
		for _, newTimelineDetail := range newTimeline {
			_, found := types.SearchIdRangesForId(timeToCheck, []*types.IdRange{newTimelineDetail.TimelineTime})
			if found {
				newValue = newTimelineDetail.ArbitraryValue
				newTimelineTime = newTimelineDetail.TimelineTime
				break
			}
		}

		if !newFound {
			panic(sdkerrors.Wrap(types.ErrInvalidPermissions, "old timeline has times that new timeline does not"))
		}

		remainingTimes, removedTimes := types.RemoveIdsFromIdRange(newTimelineTime, oldTimeline[0].TimelineTime)
		if len(removedTimes) > 0 {
			detailsToAdd := compareFunc(oldValue, newValue)
			for _, detailToAdd := range detailsToAdd {
				for _, removedTime := range removedTimes {
					detailsToCheck = append(detailsToCheck, &types.UniversalPermissionDetails{
						TimelineTime: removedTime,
						BadgeId: detailToAdd.BadgeId,
						TransferTime: detailToAdd.TransferTime,
						ToMappingId: detailToAdd.ToMappingId,
						FromMappingId: detailToAdd.FromMappingId,
						InitiatedByMappingId: detailToAdd.InitiatedByMappingId,
					})
				}
			}
		}
		
		newOldTimelineDetails := []*UniversalTimelineDetails{}
		for _, remaining := range remainingTimes {
			newOldTimelineDetails = append(newOldTimelineDetails, &UniversalTimelineDetails{
				ArbitraryValue: oldValue,
				TimelineTime: remaining,
			})
		}
		newOldTimelineDetails = append(newOldTimelineDetails, oldTimeline[1:]...)
		oldTimeline = newOldTimelineDetails
	}

	return detailsToCheck
}

func GetFirstMatchOnlyForTimeline(times [][]*types.IdRange, values []interface{}) []*UniversalTimelineDetails {
	castedPermissions := []*types.UniversalPermission{}
	for idx, time := range times {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				TimelineTimes: time, 
				ArbitraryValue: values[idx],
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}

	firstMatches := types.GetFirstMatchOnly(castedPermissions)
	
	details := []*UniversalTimelineDetails{}
	for _, firstMatch := range firstMatches {
		details = append(details, &UniversalTimelineDetails{
			ArbitraryValue: firstMatch.ArbitraryValue,
			TimelineTime: firstMatch.TimelineTime,
		})
	}

	return details
}

func CheckNotForbidden(ctx sdk.Context, permission *types.UniversalPermissionDetails) error {
	//Throw if we are in a forbidden time
	blockTime := sdk.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	_, found := types.SearchIdRangesForId(blockTime, permission.ForbiddenTimes)
	if found {
		return ErrInvalidPermissions
	}

	return nil
}

func CheckActionPermission(ctx sdk.Context, permissions []*types.ActionPermission) error {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		castedCombinations := []*types.UniversalCombination{}
		for _, combination := range permission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				PermittedTimesOptions: combination.PermittedTimesOptions,
				ForbiddenTimesOptions: combination.ForbiddenTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				PermittedTimes: permission.DefaultValues.PermittedTimes,
				ForbiddenTimes: permission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
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

func CheckTimedUpdatePermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdatePermission) error {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		castedCombinations := []*types.UniversalCombination{}
		for _, combination := range permission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				PermittedTimesOptions: combination.PermittedTimesOptions,
				ForbiddenTimesOptions: combination.ForbiddenTimesOptions,
				TimelineTimesOptions:  combination.TimelineTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				PermittedTimes: permission.DefaultValues.PermittedTimes,
				ForbiddenTimes: permission.DefaultValues.ForbiddenTimes,
				TimelineTimes:  permission.DefaultValues.TimelineTimes,
				UsesTimelineTimes: true,
			},
			Combinations: castedCombinations,
		})
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func CheckActionWithBadgeIdsPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.ActionWithBadgeIdsPermission) error {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		castedCombinations := []*types.UniversalCombination{}
		for _, combination := range permission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				PermittedTimesOptions: combination.PermittedTimesOptions,
				ForbiddenTimesOptions: combination.ForbiddenTimesOptions,
				BadgeIdsOptions:       combination.BadgeIdsOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				PermittedTimes: permission.DefaultValues.PermittedTimes,
				ForbiddenTimes: permission.DefaultValues.ForbiddenTimes,
				BadgeIds:       permission.DefaultValues.BadgeIds,
				UsesBadgeIds: true,
			},
			Combinations: castedCombinations,
		})
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

type TimeWithBadgeId struct {
	BadgeIds []*types.IdRange
	TimelineTimes []*types.IdRange
}

func CheckTimedUpdateWithBadgeIdsPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.TimedUpdateWithBadgeIdsPermission) error {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		castedCombinations := []*types.UniversalCombination{}
		for _, combination := range permission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				PermittedTimesOptions: combination.PermittedTimesOptions,
				ForbiddenTimesOptions: combination.ForbiddenTimesOptions,
				BadgeIdsOptions:       combination.BadgeIdsOptions,
				TimelineTimesOptions:  combination.TimelineTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: 		  permission.DefaultValues.BadgeIds,
				PermittedTimes: permission.DefaultValues.PermittedTimes,
				ForbiddenTimes: permission.DefaultValues.ForbiddenTimes,
				TimelineTimes:  permission.DefaultValues.TimelineTimes,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
			},
			Combinations: castedCombinations,
		})
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func CheckCollectionApprovedTransferPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.CollectionApprovedTransferPermission) error {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		castedCombinations := []*types.UniversalCombination{}
		for _, combination := range permission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				PermittedTimesOptions: combination.PermittedTimesOptions,
				ForbiddenTimesOptions: combination.ForbiddenTimesOptions,
				BadgeIdsOptions:       combination.BadgeIdsOptions,
				TimelineTimesOptions:  combination.TimelineTimesOptions,
				TransferTimesOptions:  combination.TransferTimesOptions,
				ToMappingIdOptions:    combination.ToMappingIdOptions,
				FromMappingIdOptions:  combination.FromMappingIdOptions,
				InitiatedByMappingIdOptions: combination.InitiatedByMappingIdOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: 		  permission.DefaultValues.BadgeIds,
				PermittedTimes: permission.DefaultValues.PermittedTimes,
				ForbiddenTimes: permission.DefaultValues.ForbiddenTimes,
				TimelineTimes:  permission.DefaultValues.TimelineTimes,
				TransferTimes:  permission.DefaultValues.TransferTimes,
				ToMappingId:    permission.DefaultValues.ToMappingId,
				FromMappingId:  permission.DefaultValues.FromMappingId,
				InitiatedByMappingId: permission.DefaultValues.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesFromMappingId: true,
				UsesInitiatedByMappingId: true,
			},
			Combinations: castedCombinations,
		})
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}


func CheckUserApprovedTransferPermission(ctx sdk.Context, detailsToCheck []*types.UniversalPermissionDetails, permissions []*types.UserApprovedTransferPermission) error {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		castedCombinations := []*types.UniversalCombination{}
		for _, combination := range permission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				PermittedTimesOptions: combination.PermittedTimesOptions,
				ForbiddenTimesOptions: combination.ForbiddenTimesOptions,
				BadgeIdsOptions:       combination.BadgeIdsOptions,
				TimelineTimesOptions:  combination.TimelineTimesOptions,
				TransferTimesOptions:  combination.TransferTimesOptions,
				ToMappingIdOptions:    combination.ToMappingIdOptions,
				InitiatedByMappingIdOptions: combination.InitiatedByMappingIdOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: 		  permission.DefaultValues.BadgeIds,
				PermittedTimes: permission.DefaultValues.PermittedTimes,
				ForbiddenTimes: permission.DefaultValues.ForbiddenTimes,
				TimelineTimes:  permission.DefaultValues.TimelineTimes,
				TransferTimes:  permission.DefaultValues.TransferTimes,
				ToMappingId:    permission.DefaultValues.ToMappingId,
				InitiatedByMappingId: permission.DefaultValues.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTimelineTimes: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesFromMappingId: true,
				UsesInitiatedByMappingId: true,
			},
			Combinations: castedCombinations,
		})
	}

	permissionDetails := types.GetFirstMatchOnly(castedPermissions)

	return CheckNotForbiddenForAllOverlaps(ctx, permissionDetails, detailsToCheck)
}

func CheckNotForbiddenForAllOverlaps(ctx sdk.Context, permissionDetails []*types.UniversalPermissionDetails, detailsToCheck []*types.UniversalPermissionDetails) error {
	for _, detailToCheck := range detailsToCheck {
		if detailToCheck.BadgeId == nil {
			detailToCheck.BadgeId = &types.IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) } //dummy range
		}

		if detailToCheck.TimelineTime == nil {
			detailToCheck.TimelineTime = &types.IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) } //dummy range
		}

		if detailToCheck.TransferTime == nil {
			detailToCheck.TransferTime = &types.IdRange{ Start: sdk.NewUint(math.MaxUint64), End: sdk.NewUint(math.MaxUint64) } //dummy range
		}

		//Note we are okay with the mapping IDs being ""
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