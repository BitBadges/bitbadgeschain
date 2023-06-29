package types



func ValidateActionWithBadgeIdsAndTimesPermissionUpdate(oldPermissions []*ActionWithBadgeIdsAndTimesPermission, newPermissions []*ActionWithBadgeIdsAndTimesPermission) error {
	if err := ValidateActionWithBadgeIdsAndTimesPermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateActionWithBadgeIdsAndTimesPermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastActionWithBadgeIdsAndTimesPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastActionWithBadgeIdsAndTimesPermissionToUniversalPermission(newPermissions)
	
	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}


func ValidateTimedUpdatePermissionUpdate(oldPermissions []*TimedUpdatePermission, newPermissions []*TimedUpdatePermission) error {
	if err := ValidateTimedUpdatePermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastTimedUpdatePermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastTimedUpdatePermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}



func ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions []*TimedUpdateWithBadgeIdsPermission, newPermissions []*TimedUpdateWithBadgeIdsPermission) error {
	if err := ValidateTimedUpdateWithBadgeIdsPermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateTimedUpdateWithBadgeIdsPermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastTimedUpdateWithBadgeIdsPermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}


func ValidateCollectionApprovedTransferPermissionsUpdate(oldPermissions []*CollectionApprovedTransferPermission, newPermissions []*CollectionApprovedTransferPermission) error {
	if err := ValidateCollectionApprovedTransferPermissions(oldPermissions); err != nil {
		return err
	}

	if err := ValidateCollectionApprovedTransferPermissions(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastCollectionApprovedTransferPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastCollectionApprovedTransferPermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func ValidateActionPermissionUpdate(oldPermissions []*ActionPermission, newPermissions []*ActionPermission) error {
	if err := ValidateActionPermission(oldPermissions); err != nil {
		return err
	}

	if err := ValidateActionPermission(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastActionPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastActionPermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}



func ValidateUserApprovedTransferPermissionsUpdate(oldPermissions []*UserApprovedTransferPermission, newPermissions []*UserApprovedTransferPermission) error {
	if err := ValidateUserApprovedTransferPermissions(oldPermissions); err != nil {
		return err
	}

	if err := ValidateUserApprovedTransferPermissions(newPermissions); err != nil {
		return err
	}

	castedOldPermissions := CastUserApprovedTransferPermissionToUniversalPermission(oldPermissions)
	castedNewPermissions := CastUserApprovedTransferPermissionToUniversalPermission(newPermissions)

	err := ValidateUniversalPermissionUpdate(GetFirstMatchOnly(castedOldPermissions), GetFirstMatchOnly(castedNewPermissions))
	if err != nil {
		return err
	}

	return nil
}

func ValidateUserPermissionsUpdate(oldPermissions *UserPermissions, newPermissions *UserPermissions, canBeNil bool) error {
	if err := ValidateUserPermissions(oldPermissions, canBeNil); err != nil {
		return err
	}

	if err := ValidateUserPermissions(newPermissions, canBeNil); err != nil {
		return err
	}

	if oldPermissions.CanUpdateApprovedIncomingTransfers != nil && newPermissions.CanUpdateApprovedIncomingTransfers != nil {
		if err := ValidateUserApprovedTransferPermissionsUpdate(oldPermissions.CanUpdateApprovedIncomingTransfers, newPermissions.CanUpdateApprovedIncomingTransfers); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateApprovedOutgoingTransfers != nil && newPermissions.CanUpdateApprovedOutgoingTransfers != nil {
		if err := ValidateUserApprovedTransferPermissionsUpdate(oldPermissions.CanUpdateApprovedOutgoingTransfers, newPermissions.CanUpdateApprovedOutgoingTransfers); err != nil {
			return err
		}
	}

	return nil
}


// Validate that the new permissions are valid and is not changing anything that they can't.
func ValidatePermissionsUpdate(oldPermissions *CollectionPermissions, newPermissions *CollectionPermissions, canBeNil bool) error {
	if err := ValidatePermissions(newPermissions, canBeNil); err != nil {
		return err
	}

	if err := ValidatePermissions(oldPermissions, canBeNil); err != nil {
		return err
	}

	if oldPermissions.CanDeleteCollection != nil && newPermissions.CanDeleteCollection != nil {
		if err := ValidateActionPermissionUpdate(oldPermissions.CanDeleteCollection, newPermissions.CanDeleteCollection); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateManager != nil && newPermissions.CanUpdateManager != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateManager, newPermissions.CanUpdateManager); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateCustomData != nil && newPermissions.CanUpdateCustomData != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateCustomData, newPermissions.CanUpdateCustomData); err != nil {
			return err
		}
	}
	
	if oldPermissions.CanUpdateStandards != nil && newPermissions.CanUpdateStandards != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateStandards, newPermissions.CanUpdateStandards); err != nil {
			return err
		}
	}

	if oldPermissions.CanArchive != nil && newPermissions.CanArchive != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanArchive, newPermissions.CanArchive); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateOffChainBalancesMetadata != nil && newPermissions.CanUpdateOffChainBalancesMetadata != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateOffChainBalancesMetadata, newPermissions.CanUpdateOffChainBalancesMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateCollectionMetadata != nil && newPermissions.CanUpdateCollectionMetadata != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateCollectionMetadata, newPermissions.CanUpdateCollectionMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateContractAddress != nil && newPermissions.CanUpdateContractAddress != nil {
		if err := ValidateTimedUpdatePermissionUpdate(oldPermissions.CanUpdateContractAddress, newPermissions.CanUpdateContractAddress); err != nil {
			return err
		}
	}

	if oldPermissions.CanCreateMoreBadges != nil && newPermissions.CanCreateMoreBadges != nil {
		if err := ValidateActionWithBadgeIdsAndTimesPermissionUpdate(oldPermissions.CanCreateMoreBadges, newPermissions.CanCreateMoreBadges); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateBadgeMetadata != nil && newPermissions.CanUpdateBadgeMetadata != nil {
		if err := ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions.CanUpdateBadgeMetadata, newPermissions.CanUpdateBadgeMetadata); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateInheritedBalances != nil && newPermissions.CanUpdateInheritedBalances != nil {
		if err := ValidateTimedUpdateWithBadgeIdsPermissionUpdate(oldPermissions.CanUpdateInheritedBalances, newPermissions.CanUpdateInheritedBalances); err != nil {
			return err
		}
	}

	if oldPermissions.CanUpdateCollectionApprovedTransfers != nil && newPermissions.CanUpdateCollectionApprovedTransfers != nil {
		if err := ValidateCollectionApprovedTransferPermissionsUpdate(oldPermissions.CanUpdateCollectionApprovedTransfers, newPermissions.CanUpdateCollectionApprovedTransfers); err != nil {
			return err
		}
	}

	return nil
}