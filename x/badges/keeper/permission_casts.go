package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//HACK: We cast the permissions to their UniversalPermission equivalents, so we can reuse the UniversalPermission functions

func (k Keeper) CastUserIncomingApprovalPermissionToUniversalPermission(ctx sdk.Context, permissions []*types.UserIncomingApprovalPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		approvalTrackerMapping := &types.AddressMapping{}
		if permission.AmountTrackerId == "All" {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses:        []string{},
				IncludeAddresses: false,
			}
		} else {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses:        []string{permission.AmountTrackerId},
				IncludeAddresses: true,
			}
		}

		challengeTrackerMapping := &types.AddressMapping{}
		if permission.ChallengeTrackerId == "All" {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses:        []string{},
				IncludeAddresses: false,
			}
		} else {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses:        []string{permission.ChallengeTrackerId},
				IncludeAddresses: true,
			}
		}

		fromMapping, err := k.GetAddressMappingById(ctx, permission.FromMappingId)
		if err != nil {
			return nil, err
		}

		initiatedByMapping, err := k.GetAddressMappingById(ctx, permission.InitiatedByMappingId)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:                  permission.BadgeIds,
			TransferTimes:             permission.TransferTimes,
			OwnershipTimes:            permission.OwnershipTimes,
			FromMapping:               fromMapping,
			InitiatedByMapping:        initiatedByMapping,
			AmountTrackerIdMapping:    approvalTrackerMapping,
			ChallengeTrackerIdMapping: challengeTrackerMapping,

			UsesBadgeIds:           true,
			UsesTransferTimes:      true,
			UsesOwnershipTimes:     true,
			UsesFromMapping:        true,
			UsesInitiatedByMapping: true,
			UsesAmountTrackerId:    true,
			UsesChallengeTrackerId: true,
			PermittedTimes:         permission.PermittedTimes,
			ForbiddenTimes:         permission.ForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastUserOutgoingApprovalPermissionToUniversalPermission(ctx sdk.Context, permissions []*types.UserOutgoingApprovalPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		approvalTrackerMapping := &types.AddressMapping{}
		if permission.AmountTrackerId == "All" {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses:        []string{},
				IncludeAddresses: false,
			}
		} else {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses:        []string{permission.AmountTrackerId},
				IncludeAddresses: true,
			}
		}

		challengeTrackerMapping := &types.AddressMapping{}
		if permission.ChallengeTrackerId == "All" {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses:        []string{},
				IncludeAddresses: false,
			}
		} else {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses:        []string{permission.ChallengeTrackerId},
				IncludeAddresses: true,
			}
		}
		initiatedByMapping, err := k.GetAddressMappingById(ctx, permission.InitiatedByMappingId)
		if err != nil {
			return nil, err
		}

		toMapping, err := k.GetAddressMappingById(ctx, permission.ToMappingId)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:                  permission.BadgeIds,
			TransferTimes:             permission.TransferTimes,
			OwnershipTimes:            permission.OwnershipTimes,
			ToMapping:                 toMapping,
			InitiatedByMapping:        initiatedByMapping,
			AmountTrackerIdMapping:    approvalTrackerMapping,
			ChallengeTrackerIdMapping: challengeTrackerMapping,
			UsesAmountTrackerId:       true,
			UsesChallengeTrackerId:    true,
			UsesBadgeIds:              true,
			UsesTransferTimes:         true,
			UsesOwnershipTimes:        true,
			UsesToMapping:             true,
			UsesInitiatedByMapping:    true,
			PermittedTimes:            permission.PermittedTimes,
			ForbiddenTimes:            permission.ForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastActionPermissionToUniversalPermission(actionPermission []*types.ActionPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, actionPermission := range actionPermission {

		castedPermissions = append(castedPermissions, &types.UniversalPermission{

			PermittedTimes: actionPermission.PermittedTimes,
			ForbiddenTimes: actionPermission.ForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastCollectionApprovalPermissionToUniversalPermission(ctx sdk.Context, collectionUpdatePermission []*types.CollectionApprovalPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, collectionUpdatePermission := range collectionUpdatePermission {
		approvalTrackerMapping := &types.AddressMapping{}
		if collectionUpdatePermission.AmountTrackerId == "All" {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses:        []string{},
				IncludeAddresses: false,
			}
		} else {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses:        []string{collectionUpdatePermission.AmountTrackerId},
				IncludeAddresses: true,
			}
		}

		challengeTrackerMapping := &types.AddressMapping{}
		if collectionUpdatePermission.ChallengeTrackerId == "All" {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses:        []string{},
				IncludeAddresses: false,
			}
		} else {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses:        []string{collectionUpdatePermission.ChallengeTrackerId},
				IncludeAddresses: true,
			}
		}

		fromMapping, err := k.GetAddressMappingById(ctx, collectionUpdatePermission.FromMappingId)
		if err != nil {
			return nil, err
		}

		initiatedByMapping, err := k.GetAddressMappingById(ctx, collectionUpdatePermission.InitiatedByMappingId)
		if err != nil {
			return nil, err
		}

		toMapping, err := k.GetAddressMappingById(ctx, collectionUpdatePermission.ToMappingId)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{

			TransferTimes:             collectionUpdatePermission.TransferTimes,
			OwnershipTimes:            collectionUpdatePermission.OwnershipTimes,
			ToMapping:                 toMapping,
			FromMapping:               fromMapping,
			InitiatedByMapping:        initiatedByMapping,
			BadgeIds:                  collectionUpdatePermission.BadgeIds,
			AmountTrackerIdMapping:    approvalTrackerMapping,
			ChallengeTrackerIdMapping: challengeTrackerMapping,

			UsesAmountTrackerId:    true,
			UsesChallengeTrackerId: true,
			UsesBadgeIds:           true,
			UsesTransferTimes:      true,
			UsesOwnershipTimes:     true,
			UsesToMapping:          true,
			UsesFromMapping:        true,
			UsesInitiatedByMapping: true,
			PermittedTimes:         collectionUpdatePermission.PermittedTimes,
			ForbiddenTimes:         collectionUpdatePermission.ForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(timedUpdateWithBadgeIdsPermission []*types.TimedUpdateWithBadgeIdsPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, timedUpdateWithBadgeIdsPermission := range timedUpdateWithBadgeIdsPermission {

		castedPermissions = append(castedPermissions, &types.UniversalPermission{

			TimelineTimes:     timedUpdateWithBadgeIdsPermission.TimelineTimes,
			BadgeIds:          timedUpdateWithBadgeIdsPermission.BadgeIds,
			UsesTimelineTimes: true,
			UsesBadgeIds:      true,
			PermittedTimes:    timedUpdateWithBadgeIdsPermission.PermittedTimes,
			ForbiddenTimes:    timedUpdateWithBadgeIdsPermission.ForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastTimedUpdatePermissionToUniversalPermission(timedUpdatePermission []*types.TimedUpdatePermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, timedUpdatePermission := range timedUpdatePermission {

		castedPermissions = append(castedPermissions, &types.UniversalPermission{

			TimelineTimes:     timedUpdatePermission.TimelineTimes,
			UsesTimelineTimes: true,
			PermittedTimes:    timedUpdatePermission.PermittedTimes,
			ForbiddenTimes:    timedUpdatePermission.ForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastBalancesActionPermissionToUniversalPermission(BalancesActionPermission []*types.BalancesActionPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, BalancesActionPermission := range BalancesActionPermission {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{

			BadgeIds:           BalancesActionPermission.BadgeIds,
			OwnershipTimes:     BalancesActionPermission.OwnershipTimes,
			UsesBadgeIds:       true,
			UsesOwnershipTimes: true,
			PermittedTimes:     BalancesActionPermission.PermittedTimes,
			ForbiddenTimes:     BalancesActionPermission.ForbiddenTimes,
		})
	}
	return castedPermissions, nil
}
