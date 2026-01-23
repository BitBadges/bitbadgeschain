package keeper

import (
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// For TimedUpdates that involve checking additional information (i.e. TimedUpdateWithTokenIds, CollectionApprovals, UserIncomingApprovals, etc.)
// We cast the values to a UniversalPermission struct, which is compatible with the permissions.go file in types
// This allows us to easily check overlaps and get the correct permissions

// HACK: We use the ArbitraryValue field to store the original value, so we can cast it back later
// HACK: We cast to a UniversalPermission for reusable code.
func (k Keeper) CastTokenMetadataToUniversalPermission(tokenMetadata []*types.TokenMetadata) []*types.UniversalPermission {
	castedPermissions := []*types.UniversalPermission{}
	for _, tokenMetadata := range tokenMetadata {
		marshaledMetadata, _ := json.Marshal(tokenMetadata)
		stringifiedMetadata := string(marshaledMetadata)

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			TokenIds:       tokenMetadata.TokenIds,
			UsesTokenIds:   true,
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
			TokenIds:            approval.TokenIds,
			TransferTimes:       approval.TransferTimes,
			OwnershipTimes:      approval.OwnershipTimes,
			FromList:            fromList,
			ToList:              toList,
			InitiatedByList:     initiatedByList,
			ApprovalIdList:      approvalTrackerList,
			UsesTokenIds:        true,
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
