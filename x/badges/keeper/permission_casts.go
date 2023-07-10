package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//HACK: We cast the permissions to their UniversalPermission equivalents, so we can reuse the UniversalPermission functions

func (k Keeper) CastUserApprovedIncomingTransferPermissionToUniversalPermission(ctx sdk.Context, managerAddress string, permissions []*types.UserApprovedIncomingTransferPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		castedCombinations := []*types.UniversalCombination{}
		for _, combination := range permission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				BadgeIdsOptions:           combination.BadgeIdsOptions,
				PermittedTimesOptions:     combination.PermittedTimesOptions,
				ForbiddenTimesOptions:     combination.ForbiddenTimesOptions,
				TimelineTimesOptions:      combination.TimelineTimesOptions,
				TransferTimesOptions:      combination.TransferTimesOptions,
				OwnedTimesOptions: 	 combination.OwnedTimesOptions,
				FromMappingOptions:        combination.FromMappingOptions,
				InitiatedByMappingOptions: combination.InitiatedByMappingOptions,
				
			})
		}

		fromMapping, err := k.GetAddressMappingById(ctx, permission.DefaultValues.FromMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		initiatedByMapping, err := k.GetAddressMappingById(ctx, permission.DefaultValues.InitiatedByMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds:               permission.DefaultValues.BadgeIds,
				TimelineTimes:          permission.DefaultValues.TimelineTimes,
				TransferTimes:          permission.DefaultValues.TransferTimes,
				OwnedTimes: 				permission.DefaultValues.OwnedTimes,
				FromMapping:            fromMapping,
				InitiatedByMapping:     initiatedByMapping,
				UsesBadgeIds:           true,
				UsesTimelineTimes:      true,
				UsesTransferTimes:      true,
				UsesOwnedTimes: 		true,
				UsesFromMapping:        true,
				UsesInitiatedByMapping: true,
				PermittedTimes:         permission.DefaultValues.PermittedTimes,
				ForbiddenTimes:         permission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastUserApprovedOutgoingTransferPermissionToUniversalPermission(ctx sdk.Context, managerAddress string, permissions []*types.UserApprovedOutgoingTransferPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		castedCombinations := []*types.UniversalCombination{}
		for _, combination := range permission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				BadgeIdsOptions:           combination.BadgeIdsOptions,
				PermittedTimesOptions:     combination.PermittedTimesOptions,
				ForbiddenTimesOptions:     combination.ForbiddenTimesOptions,
				TimelineTimesOptions:      combination.TimelineTimesOptions,
				TransferTimesOptions:      combination.TransferTimesOptions,
				OwnedTimesOptions: 	 combination.OwnedTimesOptions,
				ToMappingOptions:          combination.ToMappingOptions,
				InitiatedByMappingOptions: combination.InitiatedByMappingOptions,
			})
		}

		initiatedByMapping, err := k.GetAddressMappingById(ctx, permission.DefaultValues.InitiatedByMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		toMapping, err := k.GetAddressMappingById(ctx, permission.DefaultValues.ToMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds:               permission.DefaultValues.BadgeIds,
				TimelineTimes:          permission.DefaultValues.TimelineTimes,
				TransferTimes:          permission.DefaultValues.TransferTimes,
				OwnedTimes: 				permission.DefaultValues.OwnedTimes,
				ToMapping:              toMapping,
				InitiatedByMapping:     initiatedByMapping,
				UsesBadgeIds:           true,
				UsesTimelineTimes:      true,
				UsesTransferTimes:      true,
				UsesOwnedTimes: 		true,
				UsesToMapping:          true,
				UsesInitiatedByMapping: true,
				PermittedTimes:         permission.DefaultValues.PermittedTimes,
				ForbiddenTimes:         permission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastActionPermissionToUniversalPermission(actionPermission []*types.ActionPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, actionPermission := range actionPermission {
		castedCombinations := []*types.UniversalCombination{}
		for _, actionCombination := range actionPermission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				PermittedTimesOptions: actionCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: actionCombination.ForbiddenTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				PermittedTimes: actionPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes: actionPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastCollectionApprovedTransferPermissionToUniversalPermission(ctx sdk.Context, managerAddress string, collectionUpdatePermission []*types.CollectionApprovedTransferPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, collectionUpdatePermission := range collectionUpdatePermission {
		castedCombinations := []*types.UniversalCombination{}
		for _, collectionUpdateCombination := range collectionUpdatePermission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				PermittedTimesOptions:     collectionUpdateCombination.PermittedTimesOptions,
				ForbiddenTimesOptions:     collectionUpdateCombination.ForbiddenTimesOptions,
				TimelineTimesOptions:      collectionUpdateCombination.TimelineTimesOptions,
				TransferTimesOptions:      collectionUpdateCombination.TransferTimesOptions,
				OwnedTimesOptions: 	 collectionUpdateCombination.OwnedTimesOptions,
				ToMappingOptions:          collectionUpdateCombination.ToMappingOptions,
				FromMappingOptions:        collectionUpdateCombination.FromMappingOptions,
				InitiatedByMappingOptions: collectionUpdateCombination.InitiatedByMappingOptions,
				BadgeIdsOptions:           collectionUpdateCombination.BadgeIdsOptions,
			})
		}

		fromMapping, err := k.GetAddressMappingById(ctx, collectionUpdatePermission.DefaultValues.FromMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		initiatedByMapping, err := k.GetAddressMappingById(ctx, collectionUpdatePermission.DefaultValues.InitiatedByMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		toMapping, err := k.GetAddressMappingById(ctx, collectionUpdatePermission.DefaultValues.ToMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				TimelineTimes:          collectionUpdatePermission.DefaultValues.TimelineTimes,
				TransferTimes:          collectionUpdatePermission.DefaultValues.TransferTimes,
				OwnedTimes: 				collectionUpdatePermission.DefaultValues.OwnedTimes,
				ToMapping:              toMapping,
				FromMapping:            fromMapping,
				InitiatedByMapping:     initiatedByMapping,
				BadgeIds:               collectionUpdatePermission.DefaultValues.BadgeIds,
				UsesBadgeIds:           true,
				UsesTimelineTimes:      true,
				UsesTransferTimes:      true,
				UsesOwnedTimes: 		true,
				UsesToMapping:          true,
				UsesFromMapping:        true,
				UsesInitiatedByMapping: true,
				PermittedTimes:         collectionUpdatePermission.DefaultValues.PermittedTimes,
				ForbiddenTimes:         collectionUpdatePermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(timedUpdateWithBadgeIdsPermission []*types.TimedUpdateWithBadgeIdsPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, timedUpdateWithBadgeIdsPermission := range timedUpdateWithBadgeIdsPermission {
		castedCombinations := []*types.UniversalCombination{}
		for _, timedUpdateWithBadgeIdsCombination := range timedUpdateWithBadgeIdsPermission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				BadgeIdsOptions:       timedUpdateWithBadgeIdsCombination.BadgeIdsOptions,
				PermittedTimesOptions: timedUpdateWithBadgeIdsCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: timedUpdateWithBadgeIdsCombination.ForbiddenTimesOptions,
				TimelineTimesOptions:  timedUpdateWithBadgeIdsCombination.TimelineTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				TimelineTimes:     timedUpdateWithBadgeIdsPermission.DefaultValues.TimelineTimes,
				BadgeIds:          timedUpdateWithBadgeIdsPermission.DefaultValues.BadgeIds,
				UsesTimelineTimes: true,
				UsesBadgeIds:      true,
				PermittedTimes:    timedUpdateWithBadgeIdsPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes:    timedUpdateWithBadgeIdsPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastTimedUpdatePermissionToUniversalPermission(timedUpdatePermission []*types.TimedUpdatePermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, timedUpdatePermission := range timedUpdatePermission {
		castedCombinations := []*types.UniversalCombination{}
		for _, timedUpdateCombination := range timedUpdatePermission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				PermittedTimesOptions: timedUpdateCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: timedUpdateCombination.ForbiddenTimesOptions,
				TimelineTimesOptions:  timedUpdateCombination.TimelineTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				TimelineTimes:     timedUpdatePermission.DefaultValues.TimelineTimes,
				UsesTimelineTimes: true,
				PermittedTimes:    timedUpdatePermission.DefaultValues.PermittedTimes,
				ForbiddenTimes:    timedUpdatePermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastBalancesActionPermissionToUniversalPermission(BalancesActionPermission []*types.BalancesActionPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, BalancesActionPermission := range BalancesActionPermission {
		castedCombinations := []*types.UniversalCombination{}
		for _, BalancesActionCombination := range BalancesActionPermission.Combinations {
			castedCombinations = append(castedCombinations, &types.UniversalCombination{
				BadgeIdsOptions:       BalancesActionCombination.BadgeIdsOptions,
				OwnedTimesOptions:  BalancesActionCombination.OwnedTimesOptions, 
				PermittedTimesOptions: BalancesActionCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: BalancesActionCombination.ForbiddenTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds:          BalancesActionPermission.DefaultValues.BadgeIds,
				OwnedTimes:     BalancesActionPermission.DefaultValues.OwnedTimes,
				UsesBadgeIds:      true,
				UsesOwnedTimes: true,
				PermittedTimes:    BalancesActionPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes:    BalancesActionPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions, nil
}
