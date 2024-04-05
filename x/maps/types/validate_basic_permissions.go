package types

import (
	badgetypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// Validate permissions are validly formed. Disallows leading zeroes.
func ValidatePermissions(permissions *MapPermissions, canChangeValues bool) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}

	if err := badgetypes.ValidateTimedUpdatePermission(CastTimedUpdatePermissions(permissions.CanUpdateManager), canChangeValues); err != nil {
		return err
	}

	if err := badgetypes.ValidateCollectionApprovalPermissions(CastIsEditablePermissions(permissions.CanUpdateIsEditable), canChangeValues); err != nil {
		return err
	}

	if err := badgetypes.ValidateTimedUpdatePermission(CastTimedUpdatePermissions(permissions.CanUpdateMetadata), canChangeValues); err != nil {
		return err
	}

	if err := badgetypes.ValidateActionPermission(CastActionPermissions(permissions.CanDeleteMap), canChangeValues); err != nil {
		return err
	}

	if err := badgetypes.ValidateCollectionApprovalPermissions(CastIsEditablePermissions(permissions.CanUpdateIsForceEditable), canChangeValues); err != nil {
		return err
	}

	return nil
}
