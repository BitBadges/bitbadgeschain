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

func ValidateCollectionApprovalPermissions(permissions []*CollectionApprovalPermission, canChangeValues bool) error {
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

		if permission.AmountTrackerId == "" {
			return ErrAmountTrackerIdIsNil
		}

		if permission.ChallengeTrackerId == "" {
			return ErrChallengeTrackerIdIsNil
		}

		if permission.ToMappingId == "" {
			return sdkerrors.Wrap(ErrInvalidRequest, "toMappingId is nil")
		}

		if permission.FromMappingId == "" {
			return sdkerrors.Wrap(ErrInvalidRequest, "fromMappingId is nil")
		}

		if permission.InitiatedByMappingId == "" {
			return sdkerrors.Wrap(ErrInvalidRequest, "initiatedByMappingId is nil")
		}

		err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.TransferTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TransferTimes)
			permission.OwnershipTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.OwnershipTimes)
			permission.BadgeIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.BadgeIds)

			permission.PermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermittedTimes)
			permission.ForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.ForbiddenTimes)
		}
	}

	return nil
}

func ValidateUserOutgoingApprovalPermissions(permissions []*UserOutgoingApprovalPermission, canChangeValues bool) error {
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

		err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.OwnershipTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.OwnershipTimes)
			permission.TransferTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TransferTimes)
			permission.BadgeIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.BadgeIds)
			permission.PermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermittedTimes)
			permission.ForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.ForbiddenTimes)
		}
	}

	return nil
}

func ValidateUserIncomingApprovalPermissions(permissions []*UserIncomingApprovalPermission, canChangeValues bool) error {
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

		err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.OwnershipTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.OwnershipTimes)
			permission.TransferTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TransferTimes)
			permission.BadgeIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.BadgeIds)

			permission.PermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermittedTimes)
			permission.ForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.ForbiddenTimes)
		}

	}

	return nil
}

func ValidateTimedUpdateWithBadgeIdsPermission(permissions []*TimedUpdateWithBadgeIdsPermission, canChangeValues bool) error {
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

		err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.TimelineTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TimelineTimes)
			permission.BadgeIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.BadgeIds)

			permission.PermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermittedTimes)
			permission.ForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.ForbiddenTimes)
		}

		//Note we can check overlap here with other badgeIds but
		//that would take away from the flexibility of the BadgeIdsOptions.
		//Because if we have > 1 badgeIds[], then BadgeIdsOptions on the second
		//will always overlap with the first.
	}

	return nil
}

func ValidateBalancesActionPermission(permissions []*BalancesActionPermission, canChangeValues bool) error {
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

		err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.OwnershipTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.OwnershipTimes)
			permission.BadgeIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.BadgeIds)
			permission.PermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermittedTimes)
			permission.ForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.ForbiddenTimes)
		}

	}

	return nil
}

func ValidateTimedUpdatePermission(permissions []*TimedUpdatePermission, canChangeValues bool) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.TimelineTimes, false, false)
		if err != nil {
			return err
		}

		err = ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.TimelineTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TimelineTimes)

			permission.PermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermittedTimes)
			permission.ForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.ForbiddenTimes)
		}

	}

	return nil
}

func ValidateActionPermission(permissions []*ActionPermission, canChangeValues bool) error {
	if len(permissions) > 1 {
		return sdkerrors.Wrap(ErrInvalidCombinations, "only one action permission allowed")
	}

	for _, permission := range permissions {
		err := ValidatePermittedTimes(permission.PermittedTimes, permission.ForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.PermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermittedTimes)
			permission.ForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.ForbiddenTimes)
		}
	}

	return nil
}

func ValidateUserPermissions(permissions *UserPermissions, canChangeValues bool) error {
	if permissions.CanUpdateIncomingApprovals != nil {
		if err := ValidateUserIncomingApprovalPermissions(permissions.CanUpdateIncomingApprovals, canChangeValues); err != nil {
			return err
		}
	}

	if permissions.CanUpdateOutgoingApprovals != nil {
		if err := ValidateUserOutgoingApprovalPermissions(permissions.CanUpdateOutgoingApprovals, canChangeValues); err != nil {
			return err
		}
	}

	return nil
}

// Validate permissions are validly formed. Disallows leading zeroes.
func ValidatePermissions(permissions *CollectionPermissions, canChangeValues bool) error {
	if permissions == nil {
		return ErrPermissionsIsNil
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateCustomData, canChangeValues); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateStandards, canChangeValues); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateManager, canChangeValues); err != nil {
		return err
	}

	if err := ValidateTimedUpdateWithBadgeIdsPermission(permissions.CanUpdateBadgeMetadata, canChangeValues); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateCollectionMetadata, canChangeValues); err != nil {
		return err
	}

	if err := ValidateBalancesActionPermission(permissions.CanCreateMoreBadges, canChangeValues); err != nil {
		return err
	}

	if err := ValidateCollectionApprovalPermissions(permissions.CanUpdateCollectionApprovals, canChangeValues); err != nil {
		return err
	}

	if err := ValidateActionPermission(permissions.CanDeleteCollection, canChangeValues); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateOffChainBalancesMetadata, canChangeValues); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanArchiveCollection, canChangeValues); err != nil {
		return err
	}

	return nil
}
