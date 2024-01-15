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
		fromList, err := k.GetAddressListById(ctx, approval.FromListId)
		if err != nil {
			return nil, err
		}

		initiatedByList, err := k.GetAddressListById(ctx, approval.InitiatedByListId)
		if err != nil {
			return nil, err
		}

		toList, err := k.GetAddressListById(ctx, approval.ToListId)
		if err != nil {
			return nil, err
		}

		approvalTrackerList := &types.AddressList{}
		if approval.AmountTrackerId == "All" {
			approvalTrackerList = &types.AddressList{
				Addresses:        []string{},
				Whitelist: false,
			}
		} else {
			approvalTrackerList = &types.AddressList{
				Addresses:        []string{approval.AmountTrackerId},
				Whitelist: true,
			}
		}

		challengeTrackerList := &types.AddressList{}
		if approval.ChallengeTrackerId == "All" {
			challengeTrackerList = &types.AddressList{
				Addresses:        []string{},
				Whitelist: false,
			}
		} else {
			challengeTrackerList = &types.AddressList{
				Addresses:        []string{approval.ChallengeTrackerId},
				Whitelist: true,
			}
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:                  approval.BadgeIds,
			TransferTimes:             approval.TransferTimes,
			OwnershipTimes:            approval.OwnershipTimes,
			FromList:               fromList,
			ToList:                 toList,
			InitiatedByList:        initiatedByList,
			AmountTrackerIdList:    approvalTrackerList,
			ChallengeTrackerIdList: challengeTrackerList,
			UsesBadgeIds:              true,
			UsesTransferTimes:         true,
			UsesToList:             true,
			UsesFromList:           true,
			UsesInitiatedByList:    true,
			UsesOwnershipTimes:        true,
			UsesAmountTrackerId:       true,
			UsesChallengeTrackerId:    true,
			ArbitraryValue:            approval,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastUserOutgoingApprovalToUniversalPermission(ctx sdk.Context, approvals []*types.UserOutgoingApproval) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, approval := range approvals {
		initiatedByList, err := k.GetAddressListById(ctx, approval.InitiatedByListId)
		if err != nil {
			return nil, err
		}

		toList, err := k.GetAddressListById(ctx, approval.ToListId)
		if err != nil {
			return nil, err
		}

		approvalTrackerList := &types.AddressList{}
		if approval.AmountTrackerId == "All" {
			approvalTrackerList = &types.AddressList{
				Addresses:        []string{},
				Whitelist: false,
			}
		} else {
			approvalTrackerList = &types.AddressList{
				Addresses:        []string{approval.AmountTrackerId},
				Whitelist: true,
			}
		}

		challengeTrackerList := &types.AddressList{}
		if approval.ChallengeTrackerId == "All" {
			challengeTrackerList = &types.AddressList{
				Addresses:        []string{},
				Whitelist: false,
			}
		} else {
			challengeTrackerList = &types.AddressList{
				Addresses:        []string{approval.ChallengeTrackerId},
				Whitelist: true,
			}
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:                  approval.BadgeIds,
			TransferTimes:             approval.TransferTimes,
			OwnershipTimes:            approval.OwnershipTimes,
			ToList:                 toList,
			InitiatedByList:        initiatedByList,
			AmountTrackerIdList:    approvalTrackerList,
			ChallengeTrackerIdList: challengeTrackerList,
			UsesAmountTrackerId:       true,
			UsesChallengeTrackerId:    true,
			UsesBadgeIds:              true,
			UsesTransferTimes:         true,
			UsesOwnershipTimes:        true,
			UsesToList:             true,
			UsesInitiatedByList:    true,
			ArbitraryValue:            approval,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastUserIncomingApprovalToUniversalPermission(ctx sdk.Context, approvals []*types.UserIncomingApproval) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, approval := range approvals {

		fromList, err := k.GetAddressListById(ctx, approval.FromListId)
		if err != nil {
			return nil, err
		}

		initiatedByList, err := k.GetAddressListById(ctx, approval.InitiatedByListId)
		if err != nil {
			return nil, err
		}

		approvalTrackerList := &types.AddressList{}
		if approval.AmountTrackerId == "All" {
			approvalTrackerList = &types.AddressList{
				Addresses:        []string{},
				Whitelist: false,
			}
		} else {
			approvalTrackerList = &types.AddressList{
				Addresses:        []string{approval.AmountTrackerId},
				Whitelist: true,
			}
		}

		challengeTrackerList := &types.AddressList{}
		if approval.ChallengeTrackerId == "All" {
			challengeTrackerList = &types.AddressList{
				Addresses:        []string{},
				Whitelist: false,
			}
		} else {
			challengeTrackerList = &types.AddressList{
				Addresses:        []string{approval.ChallengeTrackerId},
				Whitelist: true,
			}
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:                  approval.BadgeIds,
			TransferTimes:             approval.TransferTimes,
			OwnershipTimes:            approval.OwnershipTimes,
			FromList:               fromList,
			InitiatedByList:        initiatedByList,
			AmountTrackerIdList:    approvalTrackerList,
			ChallengeTrackerIdList: challengeTrackerList,
			UsesAmountTrackerId:       true,
			UsesChallengeTrackerId:    true,
			UsesBadgeIds:              true,
			UsesTransferTimes:         true,
			UsesOwnershipTimes:        true,
			UsesFromList:           true,
			UsesInitiatedByList:    true,
			ArbitraryValue:            approval,
		})
	}
	return castedPermissions, nil
}
