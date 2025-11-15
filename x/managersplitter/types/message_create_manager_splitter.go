package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgCreateManagerSplitter) ValidateBasic() error {
	// Validate admin address format
	adminAddr, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAdmin, "invalid admin address (%s)", err)
	}
	if adminAddr.Empty() {
		return sdkerrors.Wrap(ErrInvalidAdmin, "admin address cannot be empty")
	}

	// Validate permissions are not nil
	if msg.Permissions == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "permissions cannot be nil")
	}

	// Validate that all permission criteria are properly initialized
	// (they can be empty, but the struct should exist)
	if err := validatePermissions(msg.Permissions); err != nil {
		return err
	}

	return nil
}

// validatePermissions ensures all permission criteria are properly initialized
func validatePermissions(perms *ManagerSplitterPermissions) error {
	if perms == nil {
		return sdkerrors.Wrap(ErrInvalidRequest, "permissions cannot be nil")
	}

	// Check each permission criteria - they should be initialized (can be empty)
	permissionFields := []struct {
		name     string
		criteria *PermissionCriteria
	}{
		{"canDeleteCollection", perms.CanDeleteCollection},
		{"canArchiveCollection", perms.CanArchiveCollection},
		{"canUpdateStandards", perms.CanUpdateStandards},
		{"canUpdateCustomData", perms.CanUpdateCustomData},
		{"canUpdateManager", perms.CanUpdateManager},
		{"canUpdateCollectionMetadata", perms.CanUpdateCollectionMetadata},
		{"canUpdateValidTokenIds", perms.CanUpdateValidTokenIds},
		{"canUpdateTokenMetadata", perms.CanUpdateTokenMetadata},
		{"canUpdateCollectionApprovals", perms.CanUpdateCollectionApprovals},
	}

	for _, field := range permissionFields {
		if field.criteria == nil {
			return sdkerrors.Wrapf(ErrInvalidRequest, "permission criteria for %s cannot be nil", field.name)
		}

		// Validate approved addresses format if they exist
		for i, addr := range field.criteria.ApprovedAddresses {
			if addr == "" {
				return sdkerrors.Wrapf(ErrInvalidAddress, "approved address at index %d for %s cannot be empty", i, field.name)
			}
			if _, err := sdk.AccAddressFromBech32(addr); err != nil {
				return sdkerrors.Wrapf(ErrInvalidAddress, "invalid approved address at index %d for %s (%s)", i, field.name, err)
			}
		}
	}

	return nil
}

