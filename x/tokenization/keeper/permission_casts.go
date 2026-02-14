package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HACK: We cast the permissions to their UniversalPermission equivalents, so we can reuse the UniversalPermission functions
func (k Keeper) CastUserIncomingApprovalPermissionToUniversalPermission(ctx sdk.Context, permissions []*types.UserIncomingApprovalPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		approvalTrackerList, err := k.GetTrackerListById(ctx, permission.ApprovalId)
		if err != nil {
			return nil, err
		}

		fromList, err := k.GetAddressListById(ctx, permission.FromListId)
		if err != nil {
			return nil, err
		}

		initiatedByList, err := k.GetAddressListById(ctx, permission.InitiatedByListId)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			TokenIds:        permission.TokenIds,
			TransferTimes:   permission.TransferTimes,
			OwnershipTimes:  permission.OwnershipTimes,
			FromList:        fromList,
			InitiatedByList: initiatedByList,
			ApprovalIdList:  approvalTrackerList,

			UsesTokenIds:              true,
			UsesTransferTimes:         true,
			UsesOwnershipTimes:        true,
			UsesFromList:              true,
			UsesInitiatedByList:       true,
			UsesApprovalId:            true,
			PermanentlyPermittedTimes: permission.PermanentlyPermittedTimes,
			PermanentlyForbiddenTimes: permission.PermanentlyForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastUserOutgoingApprovalPermissionToUniversalPermission(ctx sdk.Context, permissions []*types.UserOutgoingApprovalPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, permission := range permissions {
		approvalTrackerList, err := k.GetTrackerListById(ctx, permission.ApprovalId)
		if err != nil {
			return nil, err
		}

		initiatedByList, err := k.GetAddressListById(ctx, permission.InitiatedByListId)
		if err != nil {
			return nil, err
		}

		toList, err := k.GetAddressListById(ctx, permission.ToListId)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			TokenIds:                  permission.TokenIds,
			TransferTimes:             permission.TransferTimes,
			OwnershipTimes:            permission.OwnershipTimes,
			ToList:                    toList,
			InitiatedByList:           initiatedByList,
			ApprovalIdList:            approvalTrackerList,
			UsesApprovalId:            true,
			UsesTokenIds:              true,
			UsesTransferTimes:         true,
			UsesOwnershipTimes:        true,
			UsesToList:                true,
			UsesInitiatedByList:       true,
			PermanentlyPermittedTimes: permission.PermanentlyPermittedTimes,
			PermanentlyForbiddenTimes: permission.PermanentlyForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastActionPermissionToUniversalPermission(actionPermission []*types.ActionPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for i, actionPermission := range actionPermission {
		if actionPermission == nil {
			return nil, sdkerrors.Wrapf(types.ErrInvalidRequest, "action permission at index %d is nil", i)
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			PermanentlyPermittedTimes: actionPermission.PermanentlyPermittedTimes,
			PermanentlyForbiddenTimes: actionPermission.PermanentlyForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastCollectionApprovalPermissionToUniversalPermission(ctx sdk.Context, collectionUpdatePermission []*types.CollectionApprovalPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, collectionUpdatePermission := range collectionUpdatePermission {
		approvalTrackerList, err := k.GetTrackerListById(ctx, collectionUpdatePermission.ApprovalId)
		if err != nil {
			return nil, err
		}

		fromList, err := k.GetAddressListById(ctx, collectionUpdatePermission.FromListId)
		if err != nil {
			return nil, err
		}

		initiatedByList, err := k.GetAddressListById(ctx, collectionUpdatePermission.InitiatedByListId)
		if err != nil {
			return nil, err
		}

		toList, err := k.GetAddressListById(ctx, collectionUpdatePermission.ToListId)
		if err != nil {
			return nil, err
		}

		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			TransferTimes:       collectionUpdatePermission.TransferTimes,
			OwnershipTimes:      collectionUpdatePermission.OwnershipTimes,
			ToList:              toList,
			FromList:            fromList,
			InitiatedByList:     initiatedByList,
			TokenIds:            collectionUpdatePermission.TokenIds,
			ApprovalIdList:      approvalTrackerList,
			UsesApprovalId:      true,
			UsesTokenIds:        true,
			UsesTransferTimes:   true,
			UsesOwnershipTimes:  true,
			UsesToList:          true,
			UsesFromList:        true,
			UsesInitiatedByList: true,

			PermanentlyPermittedTimes: collectionUpdatePermission.PermanentlyPermittedTimes,
			PermanentlyForbiddenTimes: collectionUpdatePermission.PermanentlyForbiddenTimes,
		})
	}
	return castedPermissions, nil
}

func (k Keeper) CastTokenIdsActionPermissionToUniversalPermission(TokenIdsActionPermission []*types.TokenIdsActionPermission) ([]*types.UniversalPermission, error) {
	castedPermissions := []*types.UniversalPermission{}
	for _, TokenIdsActionPermission := range TokenIdsActionPermission {
		castedPermissions = append(castedPermissions, &types.UniversalPermission{
			TokenIds:     TokenIdsActionPermission.TokenIds,
			UsesTokenIds: true,

			PermanentlyPermittedTimes: TokenIdsActionPermission.PermanentlyPermittedTimes,
			PermanentlyForbiddenTimes: TokenIdsActionPermission.PermanentlyForbiddenTimes,
		})
	}
	return castedPermissions, nil
}
