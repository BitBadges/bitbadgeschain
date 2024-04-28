package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//This file checks that we are able to update some permission from A to B according to the rules:
//-First match only (subsequent matches are ignored)
//-If we have previously defined an explicit permitted time or forbidden time, we cannot remove or switch it. If it was previously unhandled, we can make it explciit and permanent
//-Forbidden and permitted times cannot overlap (i.e. a time cannot be both permitted and forbidden)
//-Other permission-specific validation checks
//-Duplicate combinations are disallowed because they are redundant

//HACK: For resuable code, we cast the permissions to their UniversalPermission equivalents, so we can reuse the UniversalPermission functions

func (k Keeper) ValidateBalancesActionPermissionUpdate(ctx sdk.Context, oldPermissions []*types.BalancesActionPermission, newPermissions []*types.BalancesActionPermission) error {
	if err := types.ValidateBalancesActionPermission(oldPermissions, true); err != nil {
		return err
	}

	if err := types.ValidateBalancesActionPermission(newPermissions, true); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastBalancesActionPermissionToUniversalPermission(oldPermissions)
	if err != nil {
		return err
	}

	castedNewPermissions, err := k.CastBalancesActionPermissionToUniversalPermission(newPermissions)
	if err != nil {
		return err
	}

	err = types.ValidateUniversalPermissionUpdate(ctx, types.GetFirstMatchOnly(ctx, castedOldPermissions), types.GetFirstMatchOnly(ctx, castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateTimedUpdatePermissionUpdate(ctx sdk.Context, oldPermissions []*types.TimedUpdatePermission, newPermissions []*types.TimedUpdatePermission) error {
	if err := types.ValidateTimedUpdatePermission(oldPermissions, true); err != nil {
		return err
	}

	if err := types.ValidateTimedUpdatePermission(newPermissions, true); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastTimedUpdatePermissionToUniversalPermission(oldPermissions)
	if err != nil {
		return err
	}

	castedNewPermissions, err := k.CastTimedUpdatePermissionToUniversalPermission(newPermissions)
	if err != nil {
		return err
	}

	err = types.ValidateUniversalPermissionUpdate(ctx, types.GetFirstMatchOnly(ctx, castedOldPermissions), types.GetFirstMatchOnly(ctx, castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateTimedUpdateWithBadgeIdsPermissionUpdate(ctx sdk.Context, oldPermissions []*types.TimedUpdateWithBadgeIdsPermission, newPermissions []*types.TimedUpdateWithBadgeIdsPermission) error {
	if err := types.ValidateTimedUpdateWithBadgeIdsPermission(oldPermissions, true); err != nil {
		return err
	}

	if err := types.ValidateTimedUpdateWithBadgeIdsPermission(newPermissions, true); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(oldPermissions)
	if err != nil {
		return err
	}

	castedNewPermissions, err := k.CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(newPermissions)
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

	if err := k.ValidateTimedUpdatePermissionUpdate(ctx, oldPermissions.CanUpdateManager, newPermissions.CanUpdateManager); err != nil {
		return err
	}

	if err := k.ValidateTimedUpdatePermissionUpdate(ctx, oldPermissions.CanUpdateCustomData, newPermissions.CanUpdateCustomData); err != nil {
		return err
	}

	if err := k.ValidateTimedUpdatePermissionUpdate(ctx, oldPermissions.CanUpdateStandards, newPermissions.CanUpdateStandards); err != nil {
		return err
	}

	if err := k.ValidateTimedUpdatePermissionUpdate(ctx, oldPermissions.CanArchiveCollection, newPermissions.CanArchiveCollection); err != nil {
		return err
	}

	if err := k.ValidateTimedUpdatePermissionUpdate(ctx, oldPermissions.CanUpdateOffChainBalancesMetadata, newPermissions.CanUpdateOffChainBalancesMetadata); err != nil {
		return err
	}

	if err := k.ValidateTimedUpdatePermissionUpdate(ctx, oldPermissions.CanUpdateCollectionMetadata, newPermissions.CanUpdateCollectionMetadata); err != nil {
		return err
	}

	if err := k.ValidateBalancesActionPermissionUpdate(ctx, oldPermissions.CanCreateMoreBadges, newPermissions.CanCreateMoreBadges); err != nil {
		return err
	}

	if err := k.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(ctx, oldPermissions.CanUpdateBadgeMetadata, newPermissions.CanUpdateBadgeMetadata); err != nil {
		return err
	}

	if err := k.ValidateCollectionApprovalPermissionsUpdate(ctx, oldPermissions.CanUpdateCollectionApprovals, newPermissions.CanUpdateCollectionApprovals); err != nil {
		return err
	}

	return nil
}
