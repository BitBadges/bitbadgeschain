package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//For TimedUpdates that involve checking additonal information (i.e. TimedUpdateWithBadgeIds, CollectionApprovedTransfers, UserApprovedIncomingTransfers, etc.)
//We cast the values to a UniversalPermission struct, which is compatible with the permissions.go file in types
//This allows us to easily check overlaps and get the correct permissions

func (k Keeper) CastInheritedBalancesToUniversalPermission(inheritedBalances []*types.InheritedBalance) []*types.UniversalPermission {
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

func (k Keeper) CastBadgeMetadataToUniversalPermission(badgeMetadata []*types.BadgeMetadata) []*types.UniversalPermission {
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



func (k Keeper) CastCollectionApprovedTransferToUniversalPermission(ctx sdk.Context, approvedTransfers []*types.CollectionApprovedTransfer, managerAddress string) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, approvedTransfer := range approvedTransfers {
		fromMapping, err := k.GetAddressMappingFromStore(ctx, approvedTransfer.FromMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		initiatedByMapping, err := k.GetAddressMappingFromStore(ctx, approvedTransfer.InitiatedByMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		toMapping, err := k.GetAddressMappingFromStore(ctx, approvedTransfer.ToMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: approvedTransfer.BadgeIds,
				TransferTimes: approvedTransfer.TransferTimes,
				FromMapping: fromMapping,
				ToMapping: toMapping,
				InitiatedByMapping: initiatedByMapping,
				UsesBadgeIds: true,
				UsesTransferTimes: true,
				UsesToMapping: true,
				UsesFromMapping: true,
				UsesInitiatedByMapping: true,
				ArbitraryValue: approvedTransfer,
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastUserApprovedOutgoingTransferToUniversalPermission(ctx sdk.Context, approvedTransfers []*types.UserApprovedOutgoingTransfer, managerAddress string) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, approvedTransfer := range approvedTransfers {
		initiatedByMapping, err := k.GetAddressMappingFromStore(ctx, approvedTransfer.InitiatedByMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		toMapping, err := k.GetAddressMappingFromStore(ctx, approvedTransfer.ToMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: approvedTransfer.BadgeIds,
				TransferTimes: approvedTransfer.TransferTimes,
				ToMapping: toMapping,
				InitiatedByMapping: initiatedByMapping,
				UsesBadgeIds: true,
				UsesTransferTimes: true,
				UsesToMapping: true,
				UsesInitiatedByMapping: true,
				ArbitraryValue: approvedTransfer,
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}
	return castedPermissions, nil
}



func (k Keeper) CastUserApprovedIncomingTransferToUniversalPermission(ctx sdk.Context, approvedTransfers []*types.UserApprovedIncomingTransfer, managerAddress string) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, approvedTransfer := range approvedTransfers {
		
		fromMapping, err := k.GetAddressMappingFromStore(ctx, approvedTransfer.FromMappingId, managerAddress)
		if err != nil {
			return nil, err
		}

		initiatedByMapping, err := k.GetAddressMappingFromStore(ctx, approvedTransfer.InitiatedByMappingId, managerAddress)
		if err != nil {
			return nil, err
		}
		
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			DefaultValues: &types.UniversalDefaultValues{
				BadgeIds: approvedTransfer.BadgeIds,
				TransferTimes: approvedTransfer.TransferTimes,
				FromMapping: fromMapping,
				InitiatedByMapping: initiatedByMapping,
				UsesBadgeIds: true,
				UsesTransferTimes: true,
				UsesToMapping: true,
				UsesInitiatedByMapping: true,
				ArbitraryValue: approvedTransfer,
			},
			Combinations: []*types.UniversalCombination{{}},
		})
	}
	return castedPermissions, nil
}