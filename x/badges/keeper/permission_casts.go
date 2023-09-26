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
				TransferTimesOptions:      combination.TransferTimesOptions,
				OwnershipTimesOptions: 	 combination.OwnershipTimesOptions,
				FromMappingOptions:        combination.FromMappingOptions,
				InitiatedByMappingOptions: combination.InitiatedByMappingOptions,
				
				ApprovalTrackerIdOptions: combination.ApprovalTrackerIdOptions,
				ChallengeTrackerIdOptions: combination.ChallengeTrackerIdOptions,
			})
		}

		approvalTrackerMapping := &types.AddressMapping{}
		if permission.DefaultValues.ApprovalTrackerId == "All" {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{permission.DefaultValues.ApprovalTrackerId},
				IncludeAddresses: true,
			}
		}

		challengeTrackerMapping := &types.AddressMapping{}
		if permission.DefaultValues.ChallengeTrackerId == "All" {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{permission.DefaultValues.ChallengeTrackerId},
				IncludeAddresses: true,
			}
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
				TransferTimes:          permission.DefaultValues.TransferTimes,
				OwnershipTimes: 				permission.DefaultValues.OwnershipTimes,
				FromMapping:            fromMapping,
				InitiatedByMapping:     initiatedByMapping,
				ApprovalTrackerIdMapping: approvalTrackerMapping,
				ChallengeTrackerIdMapping: challengeTrackerMapping,
				
				UsesBadgeIds:           true,
				UsesTransferTimes:      true,
				UsesOwnershipTimes: 		true,
				UsesFromMapping:        true,
				UsesInitiatedByMapping: true,
				UsesApprovalTrackerId: 	true,
				UsesChallengeTrackerId: true,
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
				TransferTimesOptions:      combination.TransferTimesOptions,
				OwnershipTimesOptions: 	 combination.OwnershipTimesOptions,
				ToMappingOptions:          combination.ToMappingOptions,
				InitiatedByMappingOptions: combination.InitiatedByMappingOptions,
				ApprovalTrackerIdOptions: combination.ApprovalTrackerIdOptions,
				ChallengeTrackerIdOptions: combination.ChallengeTrackerIdOptions,
			})
		}
		approvalTrackerMapping := &types.AddressMapping{}
		if permission.DefaultValues.ApprovalTrackerId == "All" {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{permission.DefaultValues.ApprovalTrackerId},
				IncludeAddresses: true,
			}
		}

		challengeTrackerMapping := &types.AddressMapping{}
		if permission.DefaultValues.ChallengeTrackerId == "All" {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{permission.DefaultValues.ChallengeTrackerId},
				IncludeAddresses: true,
			}
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
				TransferTimes:          permission.DefaultValues.TransferTimes,
				OwnershipTimes: 				permission.DefaultValues.OwnershipTimes,
				ToMapping:              toMapping,
				InitiatedByMapping:     initiatedByMapping,
				ApprovalTrackerIdMapping: approvalTrackerMapping,
				ChallengeTrackerIdMapping: challengeTrackerMapping,
				UsesApprovalTrackerId: 	true,
				UsesChallengeTrackerId: true,
				UsesBadgeIds:           true,
				UsesTransferTimes:      true,
				UsesOwnershipTimes: 		true,
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
				TransferTimesOptions:      collectionUpdateCombination.TransferTimesOptions,
				OwnershipTimesOptions: 	 collectionUpdateCombination.OwnershipTimesOptions,
				ToMappingOptions:          collectionUpdateCombination.ToMappingOptions,
				FromMappingOptions:        collectionUpdateCombination.FromMappingOptions,
				InitiatedByMappingOptions: collectionUpdateCombination.InitiatedByMappingOptions,
				BadgeIdsOptions:           collectionUpdateCombination.BadgeIdsOptions,
				ApprovalTrackerIdOptions: collectionUpdateCombination.ApprovalTrackerIdOptions,
				ChallengeTrackerIdOptions: collectionUpdateCombination.ChallengeTrackerIdOptions,
			})
		}

		approvalTrackerMapping := &types.AddressMapping{}
		if collectionUpdatePermission.DefaultValues.ApprovalTrackerId == "All" {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{collectionUpdatePermission.DefaultValues.ApprovalTrackerId},
				IncludeAddresses: true,
			}
		}

		challengeTrackerMapping := &types.AddressMapping{}
		if collectionUpdatePermission.DefaultValues.ChallengeTrackerId == "All" {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{collectionUpdatePermission.DefaultValues.ChallengeTrackerId},
				IncludeAddresses: true,
			}
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
				TransferTimes:          collectionUpdatePermission.DefaultValues.TransferTimes,
				OwnershipTimes: 				collectionUpdatePermission.DefaultValues.OwnershipTimes,
				ToMapping:              toMapping,
				FromMapping:            fromMapping,
				InitiatedByMapping:     initiatedByMapping,
				BadgeIds:               collectionUpdatePermission.DefaultValues.BadgeIds,
				ApprovalTrackerIdMapping: approvalTrackerMapping,
				ChallengeTrackerIdMapping: challengeTrackerMapping,

				UsesApprovalTrackerId: 	true,
				UsesChallengeTrackerId: true,
				UsesBadgeIds:           true,
				UsesTransferTimes:      true,
				UsesOwnershipTimes: 		true,
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
				OwnershipTimesOptions:  BalancesActionCombination.OwnershipTimesOptions, 
				PermittedTimesOptions: BalancesActionCombination.PermittedTimesOptions,
				ForbiddenTimesOptions: BalancesActionCombination.ForbiddenTimesOptions,
			})
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds:          BalancesActionPermission.DefaultValues.BadgeIds,
				OwnershipTimes:     BalancesActionPermission.DefaultValues.OwnershipTimes,
				UsesBadgeIds:      true,
				UsesOwnershipTimes: true,
				PermittedTimes:    BalancesActionPermission.DefaultValues.PermittedTimes,
				ForbiddenTimes:    BalancesActionPermission.DefaultValues.ForbiddenTimes,
			},
			Combinations: castedCombinations,
		})
	}
	return castedPermissions, nil
}
