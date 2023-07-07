package types

import (
	sdkerrors "cosmossdk.io/errors"
)

func ValidatePermittedTimes(permittedTimes []*UintRange, forbiddenTimes []*UintRange) error {
	//Check if any overlap between permittedTimes and forbiddenTimes
	err := ValidateRangesAreValid(permittedTimes, false)
	if err != nil {
		return sdkerrors.Wrap(err, "permittedTimes is invalid")
	}

	err = ValidateRangesAreValid(forbiddenTimes, false)
	if err != nil {
		return sdkerrors.Wrap(err, "forbiddenTimes is invalid")
	}

	err = AssertRangesDoNotOverlapAtAll(permittedTimes, forbiddenTimes)
	if err != nil {
		return sdkerrors.Wrap(err, "permittedTimes and forbiddenTimes overlap")
	}

	return nil
}

func ValidateCollectionApprovedTransferPermissions(permissions []*CollectionApprovedTransferPermission) error {
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

		permission.DefaultValues.TimelineTimes = SortAndMergeOverlapping(permission.DefaultValues.TimelineTimes)
		permission.DefaultValues.TransferTimes = SortAndMergeOverlapping(permission.DefaultValues.TransferTimes)
		permission.DefaultValues.BadgeIds = SortAndMergeOverlapping(permission.DefaultValues.BadgeIds)

		//Assert no two combinations are the same
		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.BadgeIdsOptions == combination2.BadgeIdsOptions &&
					combination.TransferTimesOptions == combination2.TransferTimesOptions &&
					combination.ToMappingOptions == combination2.ToMappingOptions &&
					combination.FromMappingOptions == combination2.FromMappingOptions &&
					combination.InitiatedByMappingOptions == combination2.InitiatedByMappingOptions &&
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

			permission.DefaultValues.PermittedTimes = SortAndMergeOverlapping(permission.DefaultValues.PermittedTimes)
			permission.DefaultValues.ForbiddenTimes = SortAndMergeOverlapping(permission.DefaultValues.ForbiddenTimes)
		}
	}

	return nil
}

func ValidateUserApprovedOutgoingTransferPermissions(permissions []*UserApprovedOutgoingTransferPermission) error {
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

		permission.DefaultValues.TimelineTimes = SortAndMergeOverlapping(permission.DefaultValues.TimelineTimes)
		permission.DefaultValues.TransferTimes = SortAndMergeOverlapping(permission.DefaultValues.TransferTimes)
		permission.DefaultValues.BadgeIds = SortAndMergeOverlapping(permission.DefaultValues.BadgeIds)

		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.BadgeIdsOptions == combination2.BadgeIdsOptions &&
					combination.TransferTimesOptions == combination2.TransferTimesOptions &&
					combination.ToMappingOptions == combination2.ToMappingOptions &&
					combination.InitiatedByMappingOptions == combination2.InitiatedByMappingOptions &&
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

			permission.DefaultValues.PermittedTimes = SortAndMergeOverlapping(permission.DefaultValues.PermittedTimes)
			permission.DefaultValues.ForbiddenTimes = SortAndMergeOverlapping(permission.DefaultValues.ForbiddenTimes)
		}
	}

	return nil
}

func ValidateUserApprovedIncomingTransferPermissions(permissions []*UserApprovedIncomingTransferPermission) error {
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

		permission.DefaultValues.TimelineTimes = SortAndMergeOverlapping(permission.DefaultValues.TimelineTimes)
		permission.DefaultValues.TransferTimes = SortAndMergeOverlapping(permission.DefaultValues.TransferTimes)
		permission.DefaultValues.BadgeIds = SortAndMergeOverlapping(permission.DefaultValues.BadgeIds)

		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.BadgeIdsOptions == combination2.BadgeIdsOptions &&
					combination.TransferTimesOptions == combination2.TransferTimesOptions &&
					combination.FromMappingOptions == combination2.FromMappingOptions &&
					combination.InitiatedByMappingOptions == combination2.InitiatedByMappingOptions &&
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

			permission.DefaultValues.PermittedTimes = SortAndMergeOverlapping(permission.DefaultValues.PermittedTimes)
			permission.DefaultValues.ForbiddenTimes = SortAndMergeOverlapping(permission.DefaultValues.ForbiddenTimes)
		}
	}

	return nil
}

func ValidateTimedUpdateWithBadgeIdsPermission(permissions []*TimedUpdateWithBadgeIdsPermission) error {
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

		permission.DefaultValues.TimelineTimes = SortAndMergeOverlapping(permission.DefaultValues.TimelineTimes)
		permission.DefaultValues.BadgeIds = SortAndMergeOverlapping(permission.DefaultValues.BadgeIds)

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

			permission.DefaultValues.PermittedTimes = SortAndMergeOverlapping(permission.DefaultValues.PermittedTimes)
			permission.DefaultValues.ForbiddenTimes = SortAndMergeOverlapping(permission.DefaultValues.ForbiddenTimes)
		}

		//Note we can check overlap here with other badgeIds but
		//that would take away from the flexibility of the BadgeIdsOptions.
		//Because if we have > 1 badgeIds[], then BadgeIdsOptions on the second
		//will always overlap with the first.
	}

	return nil
}

func ValidateBalancesActionPermission(permissions []*BalancesActionPermission) error {
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

		err = ValidateRangesAreValid(permission.DefaultValues.OwnershipTimes, false)
		if err != nil {
			return err
		}

		permission.DefaultValues.OwnershipTimes = SortAndMergeOverlapping(permission.DefaultValues.OwnershipTimes)
		permission.DefaultValues.BadgeIds = SortAndMergeOverlapping(permission.DefaultValues.BadgeIds)

		for idx, combination := range permission.Combinations {
			for _, combination2 := range permission.Combinations[idx+1:] {
				if combination.BadgeIdsOptions == combination2.BadgeIdsOptions &&
					combination.PermittedTimesOptions == combination2.PermittedTimesOptions &&
					combination.ForbiddenTimesOptions == combination2.ForbiddenTimesOptions &&
					combination.OwnershipTimesOptions == combination2.OwnershipTimesOptions {
					return ErrInvalidCombinations
				}
			}

			err := ValidatePermittedTimes(permission.DefaultValues.PermittedTimes, permission.DefaultValues.ForbiddenTimes)
			if err != nil {
				return err
			}

			permission.DefaultValues.PermittedTimes = SortAndMergeOverlapping(permission.DefaultValues.PermittedTimes)
			permission.DefaultValues.ForbiddenTimes = SortAndMergeOverlapping(permission.DefaultValues.ForbiddenTimes)
		}
	}

	return nil
}

func ValidateTimedUpdatePermission(permissions []*TimedUpdatePermission) error {
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

		permission.DefaultValues.TimelineTimes = SortAndMergeOverlapping(permission.DefaultValues.TimelineTimes)

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

			permission.DefaultValues.PermittedTimes = SortAndMergeOverlapping(permission.DefaultValues.PermittedTimes)
			permission.DefaultValues.ForbiddenTimes = SortAndMergeOverlapping(permission.DefaultValues.ForbiddenTimes)
		}
	}

	return nil
}

func ValidateActionPermission(permissions []*ActionPermission) error {
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

		permission.DefaultValues.PermittedTimes = SortAndMergeOverlapping(permission.DefaultValues.PermittedTimes)
		permission.DefaultValues.ForbiddenTimes = SortAndMergeOverlapping(permission.DefaultValues.ForbiddenTimes)

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

func ValidateUserPermissions(permissions *UserPermissions) error {
	if permissions.CanUpdateApprovedIncomingTransfers != nil {
		if err := ValidateUserApprovedIncomingTransferPermissions(permissions.CanUpdateApprovedIncomingTransfers); err != nil {
			return err
		}
	}

	if permissions.CanUpdateApprovedOutgoingTransfers != nil {
		if err := ValidateUserApprovedOutgoingTransferPermissions(permissions.CanUpdateApprovedOutgoingTransfers); err != nil {
			return err
		}
	}

	return nil
}

// Validate permissions are validly formed. Disallows leading zeroes.
func ValidatePermissions(permissions *CollectionPermissions) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateCustomData); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateStandards); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateManager); err != nil {
		return err
	}

	if err := ValidateTimedUpdateWithBadgeIdsPermission(permissions.CanUpdateBadgeMetadata); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateCollectionMetadata); err != nil {
		return err
	}

	if err := ValidateBalancesActionPermission(permissions.CanCreateMoreBadges); err != nil {
		return err
	}

	if err := ValidateCollectionApprovedTransferPermissions(permissions.CanUpdateCollectionApprovedTransfers); err != nil {
		return err
	}

	if err := ValidateActionPermission(permissions.CanDeleteCollection); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateOffChainBalancesMetadata); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateContractAddress); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanArchive); err != nil {
		return err
	}

	if err := ValidateTimedUpdateWithBadgeIdsPermission(permissions.CanUpdateInheritedBalances); err != nil {
		return err
	}

	return nil
}
