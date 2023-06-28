package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidatePermittedTimes(permittedTimes []*IdRange, forbiddenTimes []*IdRange) error {
	if permittedTimes == nil || forbiddenTimes == nil {
		return ErrInvalidCombinations
	}

	//Check if any overlap between permittedTimes and forbiddenTimes
	err := ValidateRangesAreValid(permittedTimes, false)
	if err != nil {
		return sdkerrors.Wrap(err, "permittedTimes is invalid")
	}

	err = ValidateRangesAreValid(forbiddenTimes, false)
	if err != nil {
		return sdkerrors.Wrap(err, "forbiddenTimes is invalid")
	}

	err = AssertRangeDoesNotOverlapAtAll(permittedTimes, forbiddenTimes)
	if err != nil {
		return sdkerrors.Wrap(err, "permittedTimes and forbiddenTimes overlap")
	}

	return nil
}

func ValidateCollectionApprovedTransferPermissions(permissions []*CollectionApprovedTransferPermission) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}
	
	for _, permission := range permissions {
		if permission.DefaultValues == nil {
			return ErrPermissionsValueIsNil
		}

		if permission.Combinations == nil || len(permission.Combinations) == 0 {
			return ErrCombinationsIsNil
		}

		err := ValidateRangesAreValid(permission.DefaultValues.BadgeIds, false)
		if err != nil {
			return err
		}
		
		err = ValidateRangesAreValid(permission.DefaultValues.TransferTimes, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.DefaultValues.TimelineTimes, false)
		if err != nil {
			return err
		}
		
		//Assert no two combinations are the same
		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.BadgeIdsOptions == combination2.BadgeIdsOptions && 
				combination.TransferTimesOptions == combination2.TransferTimesOptions && 
				combination.ToMappingIdOptions == combination2.ToMappingIdOptions &&
				combination.FromMappingIdOptions == combination2.FromMappingIdOptions &&
				combination.InitiatedByMappingIdOptions == combination2.InitiatedByMappingIdOptions &&
				combination.TimelineTimesOptions == combination2.TimelineTimesOptions &&
				combination.PermittedTimesOptions == combination2.PermittedTimesOptions &&
				combination.ForbiddenTimesOptions == combination2.ForbiddenTimesOptions {
					return ErrInvalidCombinations
				}
			}

			err := ValidatePermittedTimes(permission.DefaultValues.PermittedTimes, permission.DefaultValues.ForbiddenTimes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}


func ValidateUserApprovedTransferPermissions(permissions []*UserApprovedTransferPermission) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}
	
	for _, permission := range permissions {
		if permission.DefaultValues == nil {
			return ErrPermissionsValueIsNil
		}

		if permission.Combinations == nil || len(permission.Combinations) == 0 {
			return ErrCombinationsIsNil
		}

		err := ValidateRangesAreValid(permission.DefaultValues.BadgeIds, false)
		if err != nil {
			return err
		}
		
		err = ValidateRangesAreValid(permission.DefaultValues.TransferTimes, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.DefaultValues.TimelineTimes, false)
		if err != nil {
			return err
		}
		
		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.BadgeIdsOptions == combination2.BadgeIdsOptions && 
				combination.TransferTimesOptions == combination2.TransferTimesOptions && 
				combination.ToMappingIdOptions == combination2.ToMappingIdOptions &&
				combination.InitiatedByMappingIdOptions == combination2.InitiatedByMappingIdOptions &&
				combination.TimelineTimesOptions == combination2.TimelineTimesOptions &&
				combination.PermittedTimesOptions == combination2.PermittedTimesOptions &&
				combination.ForbiddenTimesOptions == combination2.ForbiddenTimesOptions {
					return ErrInvalidCombinations
				}
			}
			err := ValidatePermittedTimes(permission.DefaultValues.PermittedTimes, permission.DefaultValues.ForbiddenTimes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}



func ValidateTimedUpdateWithBadgeIdsPermission(permissions []*TimedUpdateWithBadgeIdsPermission) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}
	
	for _, permission := range permissions {
		if permission.DefaultValues == nil {
			return ErrPermissionsValueIsNil
		}

		if permission.Combinations == nil || len(permission.Combinations) == 0 {
			return ErrCombinationsIsNil
		}

		err := ValidateRangesAreValid(permission.DefaultValues.BadgeIds, true)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.DefaultValues.TimelineTimes, true)
		if err != nil {
			return err
		}

		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.BadgeIdsOptions == combination2.BadgeIdsOptions &&
				combination.TimelineTimesOptions == combination2.TimelineTimesOptions &&
				combination.PermittedTimesOptions == combination2.PermittedTimesOptions &&
				combination.ForbiddenTimesOptions == combination2.ForbiddenTimesOptions {
					return ErrInvalidCombinations
				}
			}

			err := ValidatePermittedTimes(permission.DefaultValues.PermittedTimes, permission.DefaultValues.ForbiddenTimes)
			if err != nil {
				return err
			}
		}

		//Note we can check overlap here with other badgeIds but
		//that would take away from the flexibility of the BadgeIdsOptions.
		//Because if we have > 1 badgeIds[], then BadgeIdsOptions on the second
		//will always overlap with the first.
	}

	return nil
}

func ValidateActionWithBadgeIdsPermission(permissions []*ActionWithBadgeIdsPermission) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}

	for _, permission := range permissions {
		if permission.DefaultValues == nil {
			return ErrPermissionsValueIsNil
		}

		if permission.Combinations == nil || len(permission.Combinations) == 0 {
			return ErrCombinationsIsNil
		}

		err := ValidateRangesAreValid(permission.DefaultValues.BadgeIds, false)
		if err != nil {
			return err
		}

		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.BadgeIdsOptions == combination2.BadgeIdsOptions &&
				combination.PermittedTimesOptions == combination2.PermittedTimesOptions &&
				combination.ForbiddenTimesOptions == combination2.ForbiddenTimesOptions {
					return ErrInvalidCombinations
				}
			}

			err := ValidatePermittedTimes(permission.DefaultValues.PermittedTimes, permission.DefaultValues.ForbiddenTimes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidateTimedUpdatePermission(permissions []*TimedUpdatePermission) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}

	for _, permission := range permissions {
		if permission.DefaultValues == nil {
			return ErrPermissionsValueIsNil
		}

		if permission.Combinations == nil || len(permission.Combinations) == 0 {
			return ErrCombinationsIsNil
		}

		err := ValidateRangesAreValid(permission.DefaultValues.TimelineTimes, false)
		if err != nil {
			return err
		}

		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.TimelineTimesOptions == combination2.TimelineTimesOptions &&
				combination.PermittedTimesOptions == combination2.PermittedTimesOptions &&
				combination.ForbiddenTimesOptions == combination2.ForbiddenTimesOptions {
					return ErrInvalidCombinations
				}
			}

			err := ValidatePermittedTimes(permission.DefaultValues.PermittedTimes, permission.DefaultValues.ForbiddenTimes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ValidateActionPermission(permissions []*ActionPermission) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}

	if len(permissions) > 1 {
		return sdkerrors.Wrap(ErrInvalidCombinations, "only one action permission allowed")
	}

	for _, permission := range permissions {
		err := ValidatePermittedTimes(permission.DefaultValues.PermittedTimes, permission.DefaultValues.ForbiddenTimes)
		if err != nil {
			return err
		}

		if permission.Combinations == nil || len(permission.Combinations) == 0 {
			return ErrCombinationsIsNil
		}

		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.PermittedTimesOptions == combination2.PermittedTimesOptions &&
				combination.ForbiddenTimesOptions == combination2.ForbiddenTimesOptions {
					return ErrInvalidCombinations
				}
			}
		}
	}

	return nil
}


func ValidateUserPermissions(permissions *UserPermissions, canBeNil bool) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}
	
	if !canBeNil && (permissions.CanUpdateApprovedIncomingTransfers != nil || permissions.CanUpdateApprovedOutgoingTransfers != nil) {
		return ErrPermissionsIsNil
	}

	if permissions.CanUpdateApprovedIncomingTransfers != nil {
		if err := ValidateUserApprovedTransferPermissions(permissions.CanUpdateApprovedIncomingTransfers); err != nil {
			return err
		}
	}

	if permissions.CanUpdateApprovedOutgoingTransfers != nil {
		if err := ValidateUserApprovedTransferPermissions(permissions.CanUpdateApprovedOutgoingTransfers); err != nil {
			return err
		}
	}

	return nil
}


// Validate permissions are validly formed. Disallows leading zeroes.
func ValidatePermissions(permissions *CollectionPermissions, canBeNil bool) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}

	if !canBeNil && (permissions.CanUpdateBadgeMetadata != nil || permissions.CanUpdateManager != nil || permissions.CanUpdateStandard != nil || permissions.CanUpdateCustomData != nil || permissions.CanUpdateCollectionMetadata != nil || permissions.CanCreateMoreBadges != nil || permissions.CanUpdateApprovedTransfers != nil || permissions.CanDeleteCollection != nil || permissions.CanUpdateOffChainBalancesMetadata != nil || permissions.CanUpdateContractAddress != nil || permissions.CanArchive != nil || permissions.CanUpdateInheritedBalances != nil) {
		return ErrPermissionsIsNil
	}
	
	if permissions.CanUpdateCustomData != nil {
		if err := ValidateTimedUpdatePermission(permissions.CanUpdateCustomData); err != nil {
			return err
		}
	}

	if permissions.CanUpdateStandard != nil {
		if err := ValidateTimedUpdatePermission(permissions.CanUpdateStandard); err != nil {
			return err
		}
	}

	if permissions.CanUpdateManager != nil {
		if err := ValidateTimedUpdatePermission(permissions.CanUpdateManager); err != nil {
			return err
		}
	}

	if permissions.CanUpdateBadgeMetadata != nil {
		if err := ValidateTimedUpdateWithBadgeIdsPermission(permissions.CanUpdateBadgeMetadata); err != nil {
			return err
		}
	}

	if permissions.CanUpdateCollectionMetadata != nil {
		if err := ValidateTimedUpdatePermission(permissions.CanUpdateCollectionMetadata); err != nil {
			return err
		}
	}

	if permissions.CanCreateMoreBadges != nil {
		if err := ValidateActionWithBadgeIdsPermission(permissions.CanCreateMoreBadges); err != nil {
			return err
		}
	}

	if permissions.CanUpdateApprovedTransfers != nil {
		if err := ValidateCollectionApprovedTransferPermissions(permissions.CanUpdateApprovedTransfers); err != nil {
			return err
		}
	}

	if permissions.CanDeleteCollection != nil {
		if err := ValidateActionPermission(permissions.CanDeleteCollection); err != nil {
			return err
		}
	}

	if permissions.CanUpdateOffChainBalancesMetadata != nil {
		if err := ValidateTimedUpdatePermission(permissions.CanUpdateOffChainBalancesMetadata); err != nil {
			return err
		}
	}

	if permissions.CanUpdateContractAddress != nil {
		if err := ValidateTimedUpdatePermission(permissions.CanUpdateContractAddress); err != nil {
			return err
		}
	}

	if permissions.CanArchive != nil {
		if err := ValidateTimedUpdatePermission(permissions.CanArchive); err != nil {
			return err
		}
	}

	if permissions.CanUpdateInheritedBalances != nil {
		if err := ValidateTimedUpdateWithBadgeIdsPermission(permissions.CanUpdateInheritedBalances); err != nil {
			return err
		}
	}

	return nil
}