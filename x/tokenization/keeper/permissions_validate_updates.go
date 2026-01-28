package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// This file checks that we are able to update some permission from A to B according to the rules:
//-First match only (subsequent matches are ignored)
//-If we have previously defined an explicit permitted time or forbidden time, we cannot remove or switch it. If it was previously unhandled, we can make it explciit and permanent
//-Forbidden and permitted times cannot overlap (i.e. a time cannot be both permitted and forbidden)
//-Other permission-specific validation checks
//-Duplicate combinations are disallowed because they are redundant

// HACK: For reusable code, we cast the permissions to their UniversalPermission equivalents, so we can reuse the UniversalPermission functions
func (k Keeper) ValidateTokenIdsActionPermissionUpdate(ctx sdk.Context, oldPermissions []*types.TokenIdsActionPermission, newPermissions []*types.TokenIdsActionPermission) error {
	if err := types.ValidateTokenIdsActionPermission(oldPermissions, true); err != nil {
		return err
	}

	if err := types.ValidateTokenIdsActionPermission(newPermissions, true); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastTokenIdsActionPermissionToUniversalPermission(oldPermissions)
	if err != nil {
		return err
	}

	castedNewPermissions, err := k.CastTokenIdsActionPermissionToUniversalPermission(newPermissions)
	if err != nil {
		return err
	}

	err = types.ValidateUniversalPermissionUpdate(ctx, types.GetFirstMatchOnly(ctx, castedOldPermissions), types.GetFirstMatchOnly(ctx, castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateCollectionApprovalPermissionsUpdate(ctx sdk.Context, oldPermissions []*types.CollectionApprovalPermission, newPermissions []*types.CollectionApprovalPermission) error {
	if err := types.ValidateCollectionApprovalPermissions(oldPermissions, true); err != nil {
		return err
	}

	if err := types.ValidateCollectionApprovalPermissions(newPermissions, true); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastCollectionApprovalPermissionToUniversalPermission(ctx, oldPermissions)
	if err != nil {
		return err
	}

	castedNewPermissions, err := k.CastCollectionApprovalPermissionToUniversalPermission(ctx, newPermissions)
	if err != nil {
		return err
	}

	err = types.ValidateUniversalPermissionUpdate(ctx, types.GetFirstMatchOnly(ctx, castedOldPermissions), types.GetFirstMatchOnly(ctx, castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateActionPermissionUpdate(ctx sdk.Context, oldPermissions []*types.ActionPermission, newPermissions []*types.ActionPermission) error {
	if err := types.ValidateActionPermission(oldPermissions, true); err != nil {
		return err
	}

	if err := types.ValidateActionPermission(newPermissions, true); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastActionPermissionToUniversalPermission(oldPermissions)
	if err != nil {
		return err
	}

	castedNewPermissions, err := k.CastActionPermissionToUniversalPermission(newPermissions)
	if err != nil {
		return err
	}

	err = types.ValidateUniversalPermissionUpdate(ctx, types.GetFirstMatchOnly(ctx, castedOldPermissions), types.GetFirstMatchOnly(ctx, castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserIncomingApprovalPermissionsUpdate(ctx sdk.Context, oldPermissions []*types.UserIncomingApprovalPermission, newPermissions []*types.UserIncomingApprovalPermission) error {
	if err := types.ValidateUserIncomingApprovalPermissions(oldPermissions, true); err != nil {
		return err
	}

	if err := types.ValidateUserIncomingApprovalPermissions(newPermissions, true); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastUserIncomingApprovalPermissionToUniversalPermission(ctx, oldPermissions)
	if err != nil {
		return err
	}

	castedNewPermissions, err := k.CastUserIncomingApprovalPermissionToUniversalPermission(ctx, newPermissions)
	if err != nil {
		return err
	}

	err = types.ValidateUniversalPermissionUpdate(ctx, types.GetFirstMatchOnly(ctx, castedOldPermissions), types.GetFirstMatchOnly(ctx, castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserOutgoingApprovalPermissionsUpdate(ctx sdk.Context, oldPermissions []*types.UserOutgoingApprovalPermission, newPermissions []*types.UserOutgoingApprovalPermission) error {
	if err := types.ValidateUserOutgoingApprovalPermissions(oldPermissions, true); err != nil {
		return err
	}

	if err := types.ValidateUserOutgoingApprovalPermissions(newPermissions, true); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastUserOutgoingApprovalPermissionToUniversalPermission(ctx, oldPermissions)
	if err != nil {
		return err
	}

	castedNewPermissions, err := k.CastUserOutgoingApprovalPermissionToUniversalPermission(ctx, newPermissions)
	if err != nil {
		return err
	}

	err = types.ValidateUniversalPermissionUpdate(ctx, types.GetFirstMatchOnly(ctx, castedOldPermissions), types.GetFirstMatchOnly(ctx, castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserPermissionsUpdate(ctx sdk.Context, oldPermissions *types.UserPermissions, newPermissions *types.UserPermissions) error {
	if err := types.ValidateUserPermissions(oldPermissions, true); err != nil {
		return err
	}

	if err := types.ValidateUserPermissions(newPermissions, true); err != nil {
		return err
	}

	if err := k.ValidateUserIncomingApprovalPermissionsUpdate(ctx, oldPermissions.CanUpdateIncomingApprovals, newPermissions.CanUpdateIncomingApprovals); err != nil {
		return err
	}

	if err := k.ValidateUserOutgoingApprovalPermissionsUpdate(ctx, oldPermissions.CanUpdateOutgoingApprovals, newPermissions.CanUpdateOutgoingApprovals); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanUpdateAutoApproveSelfInitiatedIncomingTransfers, newPermissions.CanUpdateAutoApproveSelfInitiatedIncomingTransfers); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers, newPermissions.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanUpdateAutoApproveAllIncomingTransfers, newPermissions.CanUpdateAutoApproveAllIncomingTransfers); err != nil {
		return err
	}

	return nil
}

// Validate that the new permissions are valid and is not changing anything that they can't.
func (k Keeper) ValidatePermissionsUpdate(ctx sdk.Context, oldPermissions *types.CollectionPermissions, newPermissions *types.CollectionPermissions) error {
	if err := types.ValidatePermissions(newPermissions, true); err != nil {
		return err
	}

	if err := types.ValidatePermissions(oldPermissions, true); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanDeleteCollection, newPermissions.CanDeleteCollection); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanUpdateManager, newPermissions.CanUpdateManager); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanUpdateCustomData, newPermissions.CanUpdateCustomData); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanUpdateStandards, newPermissions.CanUpdateStandards); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanArchiveCollection, newPermissions.CanArchiveCollection); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanUpdateCollectionMetadata, newPermissions.CanUpdateCollectionMetadata); err != nil {
		return err
	}

	if err := k.ValidateTokenIdsActionPermissionUpdate(ctx, oldPermissions.CanUpdateValidTokenIds, newPermissions.CanUpdateValidTokenIds); err != nil {
		return err
	}

	if err := k.ValidateTokenIdsActionPermissionUpdate(ctx, oldPermissions.CanUpdateTokenMetadata, newPermissions.CanUpdateTokenMetadata); err != nil {
		return err
	}

	if err := k.ValidateCollectionApprovalPermissionsUpdate(ctx, oldPermissions.CanUpdateCollectionApprovals, newPermissions.CanUpdateCollectionApprovals); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanAddMoreAliasPaths, newPermissions.CanAddMoreAliasPaths); err != nil {
		return err
	}

	if err := k.ValidateActionPermissionUpdate(ctx, oldPermissions.CanAddMoreCosmosCoinWrapperPaths, newPermissions.CanAddMoreCosmosCoinWrapperPaths); err != nil {
		return err
	}

	return nil
}
