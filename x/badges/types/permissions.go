package types

const (
	NumPermissions = 8
)

/*
	Flag bits are in the following order from left to right; leading zeroes are applied and any future additions will be appended to the right

	CanUpdateBytes: can the manager update the bytes of the badge; if false, locked forever
	CanManagerTransfer: can the manager transfer the managerial ownership of the badge to another account
	CanUpdateUris: can the manager update the uris of the class and subassets; if false, locked forever
	ForcefulTransfers: if true, one can send a badge to an account without pending approval; these badges should not by default be displayed on public profiles (can also use collections)
	CanCreate: when true, manager can create more subassets of the class; once set to false, it is locked
	CanRevoke: when true, manager can revoke subassets of the class (including null address); once set to false, it is locked
	CanFreeze: when true, manager can freeze addresseses from transferring; once set to false, it is locked
	FrozenByDefault: when true, all addresses are considered frozen and must be unfrozen to transfer; when false, all addresses are considered unfrozen and must be frozen to freeze
*/

type Permissions struct {
	CanUpdateBytes     bool
	CanManagerTransfer bool
	CanUpdateUris      bool
	ForcefulTransfers   bool
	CanCreate           bool
	CanRevoke           bool
	CanFreeze           bool
	FrozenByDefault    bool
}

const (
	CanUpdateBytesDigit = 8
	CanManagerTransferDigit = 7
	CanUpdateUrisDigit = 6
	ForcefulTransfersDigit = 5
	CanCreateDigit = 4
	CanRevokeDigit = 3
	CanFreezeDigit = 2
	FrozenByDefaultDigit = 1
)

//Validate permissions are validly formed., Disallows leading zeroes.
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

	if !oldFlags.CanUpdateUris && newFlags.CanUpdateUris {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanUpdateBytes && newFlags.CanUpdateBytes {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanCreate && newFlags.CanCreate {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanFreeze && newFlags.CanFreeze {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanRevoke && newFlags.CanRevoke {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.CanManagerTransfer && newFlags.CanManagerTransfer {
		return ErrInvalidPermissionsUpdateLocked
	}

	if oldFlags.ForcefulTransfers != newFlags.ForcefulTransfers {
		return ErrInvalidPermissionsUpdatePermanent
	}

	if oldFlags.FrozenByDefault != newFlags.FrozenByDefault {
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

//Sets the permission flags for a digit.
func SetPermissionsFlags(permission bool, digit int, flags *Permissions) {
	if digit == CanUpdateBytesDigit {
		flags.CanUpdateBytes = permission
	} else if digit == CanManagerTransferDigit {
		flags.CanManagerTransfer = permission
	} else if digit == CanUpdateUrisDigit {
		flags.CanUpdateUris = permission
	} else if digit == ForcefulTransfersDigit {
		flags.ForcefulTransfers = permission
	} else if digit == CanCreateDigit {
		flags.CanCreate = permission
	} else if digit == CanRevokeDigit {
		flags.CanRevoke = permission
	} else if digit == CanFreezeDigit {
		flags.CanFreeze = permission
	} else if digit == FrozenByDefaultDigit {
		flags.FrozenByDefault = permission
	}
}
