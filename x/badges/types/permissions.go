package types

const (
	NumPermissions = 7
)

/*
	Flag bits are in the following order from left to right; leading zeroes are applied and any future additions will be appended to the right

	can_manager_transfer: can the manager transfer the managerial ownership of the badge to another account
	can_update_uris: can the manager update the uris of the class and subassets; if false, locked forever
	forceful_transfers: if true, one can send a badge to an account without pending approval; these badges should not by default be displayed on public profiles (can also use collections)
	can_create: when true, manager can create more subassets of the class; once set to false, it is locked
	can_revoke: when true, manager can revoke subassets of the class (including null address); once set to false, it is locked
	can_freeze: when true, manager can freeze addresseses from transferring; once set to false, it is locked
	frozen_by_default: when true, all addresses are considered frozen and must be unfrozen to transfer; when false, all addresses are considered unfrozen and must be frozen to freeze
*/
//TODO: Add test cases for these permissions
type Permissions struct {
	can_manager_transfer bool
	can_update_uris      bool
	forceful_transfers   bool
	can_create           bool
	can_revoke           bool
	can_freeze           bool
	frozen_by_default    bool
}

//Check leading zeroes
func ValidatePermissions(permissions uint64) error {
	tempPermissions := permissions
	tempPermissions = tempPermissions >> NumPermissions

	if tempPermissions != 0 {
		return ErrInvalidPermissionsLeadingZeroes
	}

	return nil
}

func ValidatePermissionsUpdate(oldPermissions uint64, newPermissions uint64) error {
	if err := ValidatePermissions(newPermissions); err != nil {
		return err
	}

	if err := ValidatePermissions(oldPermissions); err != nil {
		return err
	}

	oldFlags := GetPermissions(oldPermissions)
	newFlags := GetPermissions(newPermissions)

	if !oldFlags.can_update_uris && newFlags.can_update_uris {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.can_create && newFlags.can_create {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.can_freeze && newFlags.can_freeze {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.can_revoke && newFlags.can_revoke {
		return ErrInvalidPermissionsUpdateLocked
	}

	if !oldFlags.can_manager_transfer && newFlags.can_manager_transfer {
		return ErrInvalidPermissionsUpdateLocked
	}

	if oldFlags.forceful_transfers != newFlags.forceful_transfers {
		return ErrInvalidPermissionsUpdatePermanent
	}
	if oldFlags.frozen_by_default != newFlags.frozen_by_default {
		return ErrInvalidPermissionsUpdateLocked
	}

	return nil
}

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

func (p *Permissions) CanManagerTransfer() bool {
	return p.can_manager_transfer
}

func (p *Permissions) CanUpdateUris() bool {
	return p.can_update_uris
}

func (p *Permissions) ForcefulTransfers() bool {
	return p.forceful_transfers
}

func (p *Permissions) CanCreateSubbadges() bool {
	return p.can_create
}

func (p *Permissions) CanRevoke() bool {
	return p.can_revoke
}

func (p *Permissions) CanFreeze() bool {
	return p.can_freeze
}
func (p *Permissions) FrozenByDefault() bool {
	return p.frozen_by_default
}

func SetPermissionsFlags(permission bool, digit_index int, flags *Permissions) {
	if digit_index == 7 {
		flags.can_manager_transfer = permission
	} else if digit_index == 6 {
		flags.can_update_uris = permission
	} else if digit_index == 5 {
		flags.forceful_transfers = permission
	} else if digit_index == 4 {
		flags.can_create = permission
	} else if digit_index == 3 {
		flags.can_revoke = permission
	} else if digit_index == 2 {
		flags.can_freeze = permission
	} else if digit_index == 1 {
		flags.frozen_by_default = permission
	}
}
