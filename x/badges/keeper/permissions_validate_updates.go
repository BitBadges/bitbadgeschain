package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) ValidateBalancesActionPermissionUpdate(oldPermissions []*types.BalancesActionPermission, newPermissions []*types.BalancesActionPermission) error {
	if err := types.ValidateBalancesActionPermission(oldPermissions); err != nil {
		return err
	}

	if err := types.ValidateBalancesActionPermission(newPermissions); err != nil {
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

	err = types.ValidateUniversalPermissionUpdate(types.GetFirstMatchOnly(castedOldPermissions), types.GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}


func (k Keeper) ValidateTimedUpdatePermissionUpdate(oldPermissions []*types.TimedUpdatePermission, newPermissions []*types.TimedUpdatePermission) error {
	if err := types.ValidateTimedUpdatePermission(oldPermissions); err != nil {
		return err
	}

	if err := types.ValidateTimedUpdatePermission(newPermissions); err != nil {
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


	err = types.ValidateUniversalPermissionUpdate(types.GetFirstMatchOnly(castedOldPermissions), types.GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}



func (k Keeper) ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions []*types.TimedUpdateWithBadgeIdsPermission, newPermissions []*types.TimedUpdateWithBadgeIdsPermission) error {
	if err := types.ValidateTimedUpdateWithBadgeIdsPermission(oldPermissions); err != nil {
		return err
	}

	if err := types.ValidateTimedUpdateWithBadgeIdsPermission(newPermissions); err != nil {
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


	err = types.ValidateUniversalPermissionUpdate(types.GetFirstMatchOnly(castedOldPermissions), types.GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}


func (k Keeper) ValidateCollectionApprovedTransferPermissionsUpdate(ctx sdk.Context, oldPermissions []*types.CollectionApprovedTransferPermission, newPermissions []*types.CollectionApprovedTransferPermission, managerAddress string) error {
	if err := types.ValidateCollectionApprovedTransferPermissions(oldPermissions); err != nil {
		return err
	}

	if err := types.ValidateCollectionApprovedTransferPermissions(newPermissions); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastCollectionApprovedTransferPermissionToUniversalPermission(ctx, managerAddress, oldPermissions)
	if err != nil {
		return err
	}
	
	castedNewPermissions, err := k.CastCollectionApprovedTransferPermissionToUniversalPermission(ctx, managerAddress, newPermissions)
	if err != nil {
		return err 
	}


	err = types.ValidateUniversalPermissionUpdate(types.GetFirstMatchOnly(castedOldPermissions), types.GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateActionPermissionUpdate(oldPermissions []*types.ActionPermission, newPermissions []*types.ActionPermission) error {
	if err := types.ValidateActionPermission(oldPermissions); err != nil {
		return err
	}

	if err := types.ValidateActionPermission(newPermissions); err != nil {
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


	err = types.ValidateUniversalPermissionUpdate(types.GetFirstMatchOnly(castedOldPermissions), types.GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserApprovedIncomingTransferPermissionsUpdate(ctx sdk.Context, oldPermissions []*types.UserApprovedIncomingTransferPermission, newPermissions []*types.UserApprovedIncomingTransferPermission, managerAddress string) error {
	if err := types.ValidateUserApprovedIncomingTransferPermissions(oldPermissions); err != nil {
		return err
	}

	if err := types.ValidateUserApprovedIncomingTransferPermissions(newPermissions); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastUserApprovedIncomingTransferPermissionToUniversalPermission(ctx, managerAddress, oldPermissions)
	if err != nil {
		return err
	}
	
	castedNewPermissions, err := k.CastUserApprovedIncomingTransferPermissionToUniversalPermission(ctx, managerAddress, newPermissions)
if err != nil {
	return err 
}


	err = types.ValidateUniversalPermissionUpdate(types.GetFirstMatchOnly(castedOldPermissions), types.GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}


func (k Keeper) ValidateUserApprovedOutgoingTransferPermissionsUpdate(ctx sdk.Context, oldPermissions []*types.UserApprovedOutgoingTransferPermission, newPermissions []*types.UserApprovedOutgoingTransferPermission, managerAddress string) error {
	if err := types.ValidateUserApprovedOutgoingTransferPermissions(oldPermissions); err != nil {
		return err
	}

	if err := types.ValidateUserApprovedOutgoingTransferPermissions(newPermissions); err != nil {
		return err
	}

	castedOldPermissions, err := k.CastUserApprovedOutgoingTransferPermissionToUniversalPermission(ctx, managerAddress, oldPermissions)
	if err != nil {
		return err
	}
	
	castedNewPermissions, err := k.CastUserApprovedOutgoingTransferPermissionToUniversalPermission(ctx, managerAddress, newPermissions)
if err != nil {
	return err 
}


	err = types.ValidateUniversalPermissionUpdate(types.GetFirstMatchOnly(castedOldPermissions), types.GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) ValidateUserPermissionsUpdate(ctx sdk.Context, oldPermissions *types.UserPermissions, newPermissions *types.UserPermissions, canBeNil bool, managerAddress string) error {
	if err := types.ValidateUserPermissions(oldPermissions, canBeNil); err != nil {
		return err
	}

	if err := types.ValidateUserPermissions(newPermissions, canBeNil); err != nil {
		return err
	}

	if oldPermissions.CanUpdateApprovedIncomingTransfers != nil && newPermissions.CanUpdateApprovedIncomingTransfers != nil {
		if err := k.ValidateUserApprovedIncomingTransferPermissionsUpdate(ctx, oldPermissions.CanUpdateApprovedIncomingTransfers, newPermissions.CanUpdateApprovedIncomingTransfers, managerAddress); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateApprovedOutgoingTransfers != nil && newPermissions.CanUpdateApprovedOutgoingTransfers != nil {
		if err := k.ValidateUserApprovedOutgoingTransferPermissionsUpdate(ctx, oldPermissions.CanUpdateApprovedOutgoingTransfers, newPermissions.CanUpdateApprovedOutgoingTransfers, managerAddress); err != nil {
			return err
		}
	}

	return nil
}


// Validate that the new permissions are valid and is not changing anything that they can't.
func (k Keeper) ValidatePermissionsUpdate(ctx sdk.Context, oldPermissions *types.CollectionPermissions, newPermissions *types.CollectionPermissions, canBeNil bool, managerAddress string) error {
	if err := types.ValidatePermissions(newPermissions); err != nil {
		return err
	}

	if err := types.ValidatePermissions(oldPermissions); err != nil {
		return err
	}

	if oldPermissions.CanDeleteCollection != nil && newPermissions.CanDeleteCollection != nil {
		if err := k.ValidateActionPermissionUpdate(oldPermissions.CanDeleteCollection, newPermissions.CanDeleteCollection); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateManager != nil && newPermissions.CanUpdateManager != nil {
		if err := k.ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateManager, newPermissions.CanUpdateManager); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateCustomData != nil && newPermissions.CanUpdateCustomData != nil {
		if err := k.ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateCustomData, newPermissions.CanUpdateCustomData); err != nil {
			return err
		}
	}
	
	if oldPermissions.CanUpdateStandards != nil && newPermissions.CanUpdateStandards != nil {
		if err := k.ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateStandards, newPermissions.CanUpdateStandards); err != nil {
			return err
		}
	}

	if oldPermissions.CanArchive != nil && newPermissions.CanArchive != nil {
		if err := k.ValidateTimedUpdatePermissionUpdate(oldPermissions.CanArchive, newPermissions.CanArchive); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateOffChainBalancesMetadata != nil && newPermissions.CanUpdateOffChainBalancesMetadata != nil {
		if err := k.ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateOffChainBalancesMetadata, newPermissions.CanUpdateOffChainBalancesMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateCollectionMetadata != nil && newPermissions.CanUpdateCollectionMetadata != nil {
		if err := k.ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateCollectionMetadata, newPermissions.CanUpdateCollectionMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateContractAddress != nil && newPermissions.CanUpdateContractAddress != nil {
		if err := k.ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateContractAddress, newPermissions.CanUpdateContractAddress); err != nil {
			return err
		}
	}

	if oldPermissions.CanCreateMoreBadges != nil && newPermissions.CanCreateMoreBadges != nil {
		if err := k.ValidateBalancesActionPermissionUpdate(oldPermissions.CanCreateMoreBadges, newPermissions.CanCreateMoreBadges); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateBadgeMetadata != nil && newPermissions.CanUpdateBadgeMetadata != nil {
		if err := k.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions.CanUpdateBadgeMetadata, newPermissions.CanUpdateBadgeMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateInheritedBalances != nil && newPermissions.CanUpdateInheritedBalances != nil {
		if err := k.ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions.CanUpdateInheritedBalances, newPermissions.CanUpdateInheritedBalances); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateCollectionApprovedTransfers != nil && newPermissions.CanUpdateCollectionApprovedTransfers != nil {
		if err := k.ValidateCollectionApprovedTransferPermissionsUpdate(ctx, oldPermissions.CanUpdateCollectionApprovedTransfers, newPermissions.CanUpdateCollectionApprovedTransfers, managerAddress); err != nil {
			return err
		}
	}

	return nil
}