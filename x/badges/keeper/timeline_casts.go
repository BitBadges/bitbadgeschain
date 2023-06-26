package keeper 

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

//For TimedUpdates that involve checking additonal information (i.e. TimedUpdateWithBadgeIds, CollectionApprovedTransfers, UserApprovedTransfers, etc.)
//We cast the values to a UniversalPermission struct, which is compatible with the permissions.go file in types
//This allows us to easily check overlaps and get the correct permissions

func CastInheritedBalancesToUniversalPermission(inheritedBalances []*types.InheritedBalance) []*types.UniversalPermission {
	castedPermissions := []*types.UniversalPermission{}
	for _, inheritedBalance := range inheritedBalances {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: inheritedBalance.BadgeIds,
				UsesBadgeIds: true,
				ArbitraryValue: inheritedBalance,
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}

	return castedPermissions
}

func CastBadgeMetadataToUniversalPermission(badgeMetadata []*types.BadgeMetadata) []*types.UniversalPermission {
	castedPermissions := []*types.UniversalPermission{}
	for _, badgeMetadata := range badgeMetadata {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: badgeMetadata.BadgeIds,
				UsesBadgeIds: true,
				ArbitraryValue: badgeMetadata.Uri + "<><><>" + badgeMetadata.CustomData,
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}
	return castedPermissions
}



func CastCollectionApprovedTransferToUniversalPermission(approvedTransfers []*types.CollectionApprovedTransfer) []*types.UniversalPermission {
	castedPermissions := []*types.UniversalPermission{}
	for _, approvedTransfer := range approvedTransfers {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: approvedTransfer.BadgeIds,
				TransferTimes: approvedTransfer.TransferTimes,
				ToMappingId: approvedTransfer.ToMappingId,
				FromMappingId: approvedTransfer.FromMappingId,
				InitiatedByMappingId: approvedTransfer.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesFromMappingId: true,
				UsesInitiatedByMappingId: true,
				ArbitraryValue: approvedTransfer,
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}
	return castedPermissions
}

func CastUserApprovedTransferToUniversalPermission(approvedTransfers []*types.UserApprovedTransfer) []*types.UniversalPermission {
	castedPermissions := []*types.UniversalPermission{}
	for _, approvedTransfer := range approvedTransfers {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: approvedTransfer.BadgeIds,
				TransferTimes: approvedTransfer.TransferTimes,
				ToMappingId: approvedTransfer.ToMappingId,
				InitiatedByMappingId: approvedTransfer.InitiatedByMappingId,
				UsesBadgeIds: true,
				UsesTransferTimes: true,
				UsesToMappingId: true,
				UsesInitiatedByMappingId: true,
				ArbitraryValue: approvedTransfer,
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}
	return castedPermissions
}