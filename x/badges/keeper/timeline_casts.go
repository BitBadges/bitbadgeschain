package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//For TimedUpdates that involve checking additonal information (i.e. TimedUpdateWithBadgeIds, CollectionApprovals, UserIncomingApprovals, etc.)
//We cast the values to a UniversalPermission struct, which is compatible with the permissions.go file in types
//This allows us to easily check overlaps and get the correct permissions

//HACK: We use the ArbitraryValue field to store the original value, so we can cast it back later
//HACK: We cast to a UniversalPermission for reusable code.

func (k Keeper) CastInheritedBalancesToUniversalPermission(inheritedBalances []*types.InheritedBalance) []*types.UniversalPermission {
	castedPermissions := []*types.UniversalPermission{}
	for _, inheritedBalance := range inheritedBalances {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:       inheritedBalance.BadgeIds,
				UsesBadgeIds:   true,
				ArbitraryValue: inheritedBalance,
		})
	}

	return castedPermissions
}

func (k Keeper) CastBadgeMetadataToUniversalPermission(badgeMetadata []*types.BadgeMetadata) []*types.UniversalPermission {
	castedPermissions := []*types.UniversalPermission{}
	for _, badgeMetadata := range badgeMetadata {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:       badgeMetadata.BadgeIds,
			UsesBadgeIds:   true,
			ArbitraryValue: badgeMetadata.Uri + "<><><>" + badgeMetadata.CustomData,
		})
	}
	return castedPermissions
}

func (k Keeper) CastCollectionApprovalToUniversalPermission(ctx sdk.Context, approvals []*types.CollectionApproval) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, approval := range approvals {
		fromMapping, err := k.GetAddressMappingById(ctx, approval.FromMappingId)
		if err != nil {
			return nil, err
		}

		initiatedByMapping, err := k.GetAddressMappingById(ctx, approval.InitiatedByMappingId)
		if err != nil {
			return nil, err
		}

		toMapping, err := k.GetAddressMappingById(ctx, approval.ToMappingId)
		if err != nil {
			return nil, err
		}

		approvalTrackerMapping := &types.AddressMapping{}
		if approval.AmountTrackerId == "All" {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{approval.AmountTrackerId},
				IncludeAddresses: true,
			}
		}

		challengeTrackerMapping := &types.AddressMapping{}
		if approval.ChallengeTrackerId == "All" {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{approval.ChallengeTrackerId},
				IncludeAddresses: true,
			}
		}



		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:               approval.BadgeIds,
				TransferTimes:          approval.TransferTimes,
				OwnershipTimes: 			 	approval.OwnershipTimes,
				FromMapping:            fromMapping,
				ToMapping:              toMapping,
				InitiatedByMapping:     initiatedByMapping,
				AmountTrackerIdMapping: approvalTrackerMapping,
				ChallengeTrackerIdMapping: challengeTrackerMapping,
				UsesBadgeIds:           true,
				UsesTransferTimes:      true,
				UsesToMapping:          true,
				UsesFromMapping:        true,
				UsesInitiatedByMapping: true,
				UsesOwnershipTimes: 	 	true,
				UsesAmountTrackerId: 	true,
				UsesChallengeTrackerId: true,
				ArbitraryValue:         approval,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastUserOutgoingApprovalToUniversalPermission(ctx sdk.Context, approvals []*types.UserOutgoingApproval) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, approval := range approvals {
		initiatedByMapping, err := k.GetAddressMappingById(ctx, approval.InitiatedByMappingId)
		if err != nil {
			return nil, err
		}

		toMapping, err := k.GetAddressMappingById(ctx, approval.ToMappingId)
		if err != nil {
			return nil, err
		}

		approvalTrackerMapping := &types.AddressMapping{}
		if approval.AmountTrackerId == "All" {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{approval.AmountTrackerId},
				IncludeAddresses: true,
			}
		}

		challengeTrackerMapping := &types.AddressMapping{}
		if approval.ChallengeTrackerId == "All" {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{approval.ChallengeTrackerId},
				IncludeAddresses: true,
			}
		}


		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:               approval.BadgeIds,
				TransferTimes:          approval.TransferTimes,
				OwnershipTimes: 			 	approval.OwnershipTimes,
				ToMapping:              toMapping,
				InitiatedByMapping:     initiatedByMapping,
				AmountTrackerIdMapping: approvalTrackerMapping,
				ChallengeTrackerIdMapping: challengeTrackerMapping,
				UsesAmountTrackerId: 	true,
				UsesChallengeTrackerId: true,
				UsesBadgeIds:           true,
				UsesTransferTimes:      true,
				UsesOwnershipTimes: 	 true,
				UsesToMapping:          true,
				UsesInitiatedByMapping: true,
				ArbitraryValue:         approval,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastUserIncomingApprovalToUniversalPermission(ctx sdk.Context, approvals []*types.UserIncomingApproval) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, approval := range approvals {

		fromMapping, err := k.GetAddressMappingById(ctx, approval.FromMappingId)
		if err != nil {
			return nil, err
		}

		initiatedByMapping, err := k.GetAddressMappingById(ctx, approval.InitiatedByMappingId)
		if err != nil {
			return nil, err
		}

		approvalTrackerMapping := &types.AddressMapping{}
		if approval.AmountTrackerId == "All" {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			approvalTrackerMapping = &types.AddressMapping{
				Addresses: []string{approval.AmountTrackerId},
				IncludeAddresses: true,
			}
		}

		challengeTrackerMapping := &types.AddressMapping{}
		if approval.ChallengeTrackerId == "All" {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{},
				IncludeAddresses: false,
			}
		} else {
			challengeTrackerMapping = &types.AddressMapping{
				Addresses: []string{approval.ChallengeTrackerId},
				IncludeAddresses: true,
			}
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:               approval.BadgeIds,
				TransferTimes:          approval.TransferTimes,
				OwnershipTimes: 			 			approval.OwnershipTimes,
				FromMapping:            fromMapping,
				InitiatedByMapping:     initiatedByMapping,
				AmountTrackerIdMapping: approvalTrackerMapping,
				ChallengeTrackerIdMapping: challengeTrackerMapping,
				UsesAmountTrackerId: 	true,
				UsesChallengeTrackerId: true,
				UsesBadgeIds:           true,
				UsesTransferTimes:      true,
				UsesOwnershipTimes: 	 			true,
				UsesFromMapping:        true,
				UsesInitiatedByMapping: true,
				ArbitraryValue:         approval,
		})
	}
	return castedPermissions, nil
}
