package types

const (
	NumPermissions = 7
)

/*
	Flag bits are in the following order from left to right; leading zeroes are applied and any future additions will be appended to the right
*/

type Permissions struct {
	CanUpdateBalancesUri    bool //can the manager update the balances uri of the collection; if false, locked forever
	CanDelete               bool //when true, manager can delete the collection; once set to false, it is locked
	CanUpdateBytes          bool //can the manager update the bytes of the badge; if false, locked forever
	CanManagerBeTransferred bool //can the manager transfer the managerial ownership of the badge to another account; if false, locked forever
	CanUpdateMetadataUris   bool //can the manager update the uris (metadata) of the collection and badges; if false, locked forever
	CanCreateMoreBadges     bool //when true, manager can create more badges of the collection; once set to false, the number of badges in the collection is locked
	CanUpdateAllowed     		bool //when true, manager can freeze and unfreeze addresseses from transferring; once set to false, it is locked
}

const (
	CanUpdateBalancesUriDigit    = 7
	CanDeleteDigit               = 6
	CanUpdateBytesDigit          = 5
	CanManagerBeTransferredDigit = 4
	CanUpdateMetadataUrisDigit   = 3
	CanCreateMoreBadgesDigit     = 2
	CanUpdateAllowedDigit        = 1
)

// Validate permissions are validly formed. Disallows leading zeroes.
func ValidatePermissions(permissions uint64) error {
	tempPermissions := permissions >> NumPermissions

	if tempPermissions != 0 {
		return ErrInvalidPermissionsLeadingZeroes
	}

	return nil
}

// Validate that the new permissions are valid and is not changing anything that they can't.
func ValidatePermissionsUpdate(oldPermissions uint64, newPermissions uint64) error {
	if err := ValidatePermissions(newPermissions); err != nil {
		return err
	}

	if err := ValidatePermissions(oldPermissions); err != nil {
		return err
	}

	oldFlags := GetPermissions(oldPermissions)
	newFlags := GetPermissions(newPermissions)

	if !oldFlags.CanUpdateMetadataUris && newFlags.CanUpdateMetadataUris {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanUpdateBytes && newFlags.CanUpdateBytes {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanCreateMoreBadges && newFlags.CanCreateMoreBadges {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanUpdateAllowed && newFlags.CanUpdateAllowed {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanManagerBeTransferred && newFlags.CanManagerBeTransferred {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanDelete && newFlags.CanDelete {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanUpdateBalancesUri && newFlags.CanUpdateBalancesUri {
		return ErrInvalidPermissionsUpdateLocked
	}

	return nil
}

// Get permissions from permissions number
func GetPermissions(permissions uint64) Permissions {
	var flags Permissions = Permissions{}
	for i := 0; i <= NumPermissions; i++ {
		mask := uint64(1) << i
		masked_n := permissions & mask
		bit := masked_n >> i
		bit_as_bool := bit == 1

		SetPermissionsFlags(bit_as_bool, i+1, &flags)
	}

	return flags
}

// Sets the permission flags for a digit.
func SetPermissionsFlags(permission bool, digit int, flags *Permissions) {
	if digit == CanUpdateBytesDigit {
		flags.CanUpdateBytes = permission
	} else if digit == CanManagerBeTransferredDigit {
		flags.CanManagerBeTransferred = permission
	} else if digit == CanUpdateMetadataUrisDigit {
		flags.CanUpdateMetadataUris = permission
	} else if digit == CanCreateMoreBadgesDigit {
		flags.CanCreateMoreBadges = permission
	} else if digit == CanUpdateAllowedDigit {
		flags.CanUpdateAllowed = permission
	} else if digit == CanDeleteDigit {
		flags.CanDelete = permission
	} else if digit == CanUpdateBalancesUriDigit {
		flags.CanUpdateBalancesUri = permission
	}
}
