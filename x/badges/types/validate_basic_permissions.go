package types

import (
	sdkerrors "cosmossdk.io/errors"
)

func ValidatePermittedTimes(permittedTimes []*UintRange, forbiddenTimes []*UintRange) error {
	//Check if any overlap between permittedTimes and forbiddenTimes
	err := ValidateRangesAreValid(permittedTimes, false, false)
	if err != nil {
		return sdkerrors.Wrap(err, "permittedTimes is invalid")
	}

	err = ValidateRangesAreValid(forbiddenTimes, false, false)
	if err != nil {
		return sdkerrors.Wrap(err, "forbiddenTimes is invalid")
	}

	err = AssertRangesDoNotOverlapAtAll(permittedTimes, forbiddenTimes)
	if err != nil {
		return sdkerrors.Wrap(err, "permittedTimes and forbiddenTimes overlap")
	}

	return nil
}

func ValidateCollectionApprovalPermissions(permissions []*CollectionApprovalPermission) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.BadgeIds, false, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.TransferTimes, false, false)
		if err != nil {
			return err
		}
		
		err = ValidateRangesAreValid(permission.OwnershipTimes, false, false)
		if err != nil {
			return err
		}

		// if permission.ApprovalTrackerId == "" {
		// 	return ErrApprovalTrackerIdIsNil
		// }

		// if permission.ChallengeTrackerId == "" {
		// 	return ErrChallengeTrackerIdIsNil
		// }
		
		permission.TransferTimes = SortAndMergeOverlapping(permission.TransferTimes)
		permission.OwnershipTimes = SortAndMergeOverlapping(permission.OwnershipTimes)
		permission.BadgeIds = SortAndMergeOverlapping(permission.BadgeIds)

		err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		permission.PermittedTimes = SortAndMergeOverlapping(permission.PermittedTimes)
		permission.ForbiddenTimes = SortAndMergeOverlapping(permission.ForbiddenTimes)
	}

	return nil
}

func ValidateUserOutgoingApprovalPermissions(permissions []*UserOutgoingApprovalPermission) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.BadgeIds, false, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.TransferTimes, false, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.OwnershipTimes, false, false)
		if err != nil {
			return err
		}

		permission.OwnershipTimes = SortAndMergeOverlapping(permission.OwnershipTimes)
		permission.TransferTimes = SortAndMergeOverlapping(permission.TransferTimes)
		permission.BadgeIds = SortAndMergeOverlapping(permission.BadgeIds)

		err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		permission.PermittedTimes = SortAndMergeOverlapping(permission.PermittedTimes)
		permission.ForbiddenTimes = SortAndMergeOverlapping(permission.ForbiddenTimes)
		
	}

	return nil
}

func ValidateUserIncomingApprovalPermissions(permissions []*UserIncomingApprovalPermission) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.BadgeIds, false, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.TransferTimes, false, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.OwnershipTimes, false, false)
		if err != nil {
			return err
		}

		permission.OwnershipTimes = SortAndMergeOverlapping(permission.OwnershipTimes)
		permission.TransferTimes = SortAndMergeOverlapping(permission.TransferTimes)
		permission.BadgeIds = SortAndMergeOverlapping(permission.BadgeIds)

		err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		permission.PermittedTimes = SortAndMergeOverlapping(permission.PermittedTimes)
		permission.ForbiddenTimes = SortAndMergeOverlapping(permission.ForbiddenTimes)

	}

	return nil
}

func ValidateTimedUpdateWithBadgeIdsPermission(permissions []*TimedUpdateWithBadgeIdsPermission) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.BadgeIds, false, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.TimelineTimes, false, false)
		if err != nil {
			return err
		}

		permission.TimelineTimes = SortAndMergeOverlapping(permission.TimelineTimes)
		permission.BadgeIds = SortAndMergeOverlapping(permission.BadgeIds)


			err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
			if err != nil {
				return err
			}

			permission.PermittedTimes = SortAndMergeOverlapping(permission.PermittedTimes)
			permission.ForbiddenTimes = SortAndMergeOverlapping(permission.ForbiddenTimes)
		

		//Note we can check overlap here with other badgeIds but
		//that would take away from the flexibility of the BadgeIdsOptions.
		//Because if we have > 1 badgeIds[], then BadgeIdsOptions on the second
		//will always overlap with the first.
	}

	return nil
}

func ValidateBalancesActionPermission(permissions []*BalancesActionPermission) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.BadgeIds, false, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.OwnershipTimes, false, false)
		if err != nil {
			return err
		}

		permission.OwnershipTimes = SortAndMergeOverlapping(permission.OwnershipTimes)
		permission.BadgeIds = SortAndMergeOverlapping(permission.BadgeIds)


			err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
			if err != nil {
				return err
			}

			permission.PermittedTimes = SortAndMergeOverlapping(permission.PermittedTimes)
			permission.ForbiddenTimes = SortAndMergeOverlapping(permission.ForbiddenTimes)
		
	}

	return nil
}

func ValidateTimedUpdatePermission(permissions []*TimedUpdatePermission) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.TimelineTimes, false, false)
		if err != nil {
			return err
		}

		permission.TimelineTimes = SortAndMergeOverlapping(permission.TimelineTimes)

		

			err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
			if err != nil {
				return err
			}

			permission.PermittedTimes = SortAndMergeOverlapping(permission.PermittedTimes)
			permission.ForbiddenTimes = SortAndMergeOverlapping(permission.ForbiddenTimes)
		
	}

	return nil
}

func ValidateActionPermission(permissions []*ActionPermission) error {
	if len(permissions) > 1 {
		return sdkerrors.Wrap(ErrInvalidCombinations, "only one action permission allowed")
	}

	for _, permission := range permissions {
		err := ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		permission.PermittedTimes = SortAndMergeOverlapping(permission.PermittedTimes)
		permission.ForbiddenTimes = SortAndMergeOverlapping(permission.ForbiddenTimes)
	}

	return nil
}

func ValidateUserPermissions(permissions *UserPermissions) error {
	if permissions.CanUpdateIncomingApprovals != nil {
		if err := ValidateUserIncomingApprovalPermissions(permissions.CanUpdateIncomingApprovals); err != nil {
			return err
		}
	}

	if permissions.CanUpdateOutgoingApprovals != nil {
		if err := ValidateUserOutgoingApprovalPermissions(permissions.CanUpdateOutgoingApprovals); err != nil {
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

	if err := ValidateCollectionApprovalPermissions(permissions.CanUpdateCollectionApprovals); err != nil {
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

	if err := ValidateTimedUpdatePermission(permissions.CanArchiveCollection); err != nil {
		return err
	}

	return nil
}
