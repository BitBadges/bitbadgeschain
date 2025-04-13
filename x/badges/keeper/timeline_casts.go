package keeper

import (
	"encoding/json"

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
		marshaledMetadata, _ := json.Marshal(badgeMetadata)
		stringifiedMetadata := string(marshaledMetadata)

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:       badgeMetadata.BadgeIds,
			UsesBadgeIds:   true,
			ArbitraryValue: stringifiedMetadata,
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

		approvalTrackerList := &types.AddressList{
			Addresses: []string{approval.ApprovalId},
			Whitelist: true,
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			BadgeIds:            approval.BadgeIds,
			TransferTimes:       approval.TransferTimes,
			OwnershipTimes:      approval.OwnershipTimes,
			FromList:            fromList,
			ToList:              toList,
			InitiatedByList:     initiatedByList,
			ApprovalIdList:      approvalTrackerList,
			UsesBadgeIds:        true,
			UsesTransferTimes:   true,
			UsesToList:          true,
			UsesFromList:        true,
			UsesInitiatedByList: true,
			UsesOwnershipTimes:  true,
			UsesApprovalId:      true,
			ArbitraryValue:      approval,
		})
	}
	return castedPermissions, nil
}
