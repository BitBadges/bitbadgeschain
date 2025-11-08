package types

import (
	sdkerrors "cosmossdk.io/errors"
)

func ValidatePermanentlyPermittedTimes(permanentlyPermittedTimes []*UintRange, permanentlyForbiddenTimes []*UintRange) error {
	//Check if any overlap between permanentlyPermittedTimes and permanentlyForbiddenTimes
	err := ValidateRangesAreValid(permanentlyPermittedTimes, false, false)
	if err != nil {
		return sdkerrors.Wrap(err, "permanentlyPermittedTimes is invalid")
	}

	err = ValidateRangesAreValid(permanentlyForbiddenTimes, false, false)
	if err != nil {
		return sdkerrors.Wrap(err, "permanentlyForbiddenTimes is invalid")
	}

	err = AssertRangesDoNotOverlapAtAll(permanentlyPermittedTimes, permanentlyForbiddenTimes)
	if err != nil {
		return sdkerrors.Wrap(err, "permanentlyPermittedTimes and permanentlyForbiddenTimes overlap")
	}

	return nil
}

func ValidateCollectionApprovalPermissions(permissions []*CollectionApprovalPermission, canChangeValues bool) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.TokenIds, false, false)
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

		if permission.ApprovalId == "" {
			return ErrAmountTrackerIdIsNil
		}

		if permission.ToListId == "" {
			return sdkerrors.Wrap(ErrInvalidRequest, "toListId is nil")
		}

		if permission.FromListId == "" {
			return sdkerrors.Wrap(ErrInvalidRequest, "fromListId is nil")
		}

		if permission.InitiatedByListId == "" {
			return sdkerrors.Wrap(ErrInvalidRequest, "initiatedByListId is nil")
		}

		err = ValidatePermanentlyPermittedTimes(permission.PermanentlyPermittedTimes, permission.PermanentlyForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.TransferTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TransferTimes)
			permission.OwnershipTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.OwnershipTimes)
			permission.TokenIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TokenIds)

			permission.PermanentlyPermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyPermittedTimes)
			permission.PermanentlyForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyForbiddenTimes)
		}
	}

	return nil
}

func ValidateUserOutgoingApprovalPermissions(permissions []*UserOutgoingApprovalPermission, canChangeValues bool) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.TokenIds, false, false)
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

		err = ValidatePermanentlyPermittedTimes(permission.PermanentlyPermittedTimes, permission.PermanentlyForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.OwnershipTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.OwnershipTimes)
			permission.TransferTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TransferTimes)
			permission.TokenIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TokenIds)
			permission.PermanentlyPermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyPermittedTimes)
			permission.PermanentlyForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyForbiddenTimes)
		}
	}

	return nil
}

func ValidateUserIncomingApprovalPermissions(permissions []*UserIncomingApprovalPermission, canChangeValues bool) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.TokenIds, false, false)
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

		err = ValidatePermanentlyPermittedTimes(permission.PermanentlyPermittedTimes, permission.PermanentlyForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.OwnershipTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.OwnershipTimes)
			permission.TransferTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TransferTimes)
			permission.TokenIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TokenIds)

			permission.PermanentlyPermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyPermittedTimes)
			permission.PermanentlyForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyForbiddenTimes)
		}

	}

	return nil
}

func ValidateTimedUpdateWithTokenIdsPermission(permissions []*TimedUpdateWithTokenIdsPermission, canChangeValues bool) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.TokenIds, false, false)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(permission.TimelineTimes, false, false)
		if err != nil {
			return err
		}

		err = ValidatePermanentlyPermittedTimes(permission.PermanentlyPermittedTimes, permission.PermanentlyForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.TimelineTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TimelineTimes)
			permission.TokenIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TokenIds)

			permission.PermanentlyPermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyPermittedTimes)
			permission.PermanentlyForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyForbiddenTimes)
		}

		//Note we can check overlap here with other tokenIds but
		//that would take away from the flexibility of the TokenIdsOptions.
		//Because if we have > 1 tokenIds[], then TokenIdsOptions on the second
		//will always overlap with the first.
	}

	return nil
}

func ValidateTokenIdsActionPermission(permissions []*TokenIdsActionPermission, canChangeValues bool) error {
	for _, permission := range permissions {
		if permission == nil {
			return ErrPermissionsValueIsNil
		}

		err := ValidateRangesAreValid(permission.TokenIds, false, false)
		if err != nil {
			return err
		}

		err = ValidatePermanentlyPermittedTimes(permission.PermanentlyPermittedTimes, permission.PermanentlyForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.TokenIds = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TokenIds)
			permission.PermanentlyPermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyPermittedTimes)
			permission.PermanentlyForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyForbiddenTimes)
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

		err = ValidatePermanentlyPermittedTimes(permission.PermanentlyPermittedTimes, permission.PermanentlyForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.TimelineTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.TimelineTimes)

			permission.PermanentlyPermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyPermittedTimes)
			permission.PermanentlyForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyForbiddenTimes)
		}

	}

	return nil
}

func ValidateActionPermission(permissions []*ActionPermission, canChangeValues bool) error {
	if len(permissions) > 1 {
		return sdkerrors.Wrap(ErrInvalidCombinations, "only one action permission allowed")
	}

	for _, permission := range permissions {
		err := ValidatePermanentlyPermittedTimes(permission.PermanentlyPermittedTimes, permission.PermanentlyForbiddenTimes)
		if err != nil {
			return err
		}

		if canChangeValues {
			permission.PermanentlyPermittedTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyPermittedTimes)
			permission.PermanentlyForbiddenTimes = SortUintRangesAndMergeAdjacentAndIntersecting(permission.PermanentlyForbiddenTimes)
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

	if err := ValidateTimedUpdateWithTokenIdsPermission(permissions.CanUpdateTokenMetadata, canChangeValues); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanUpdateCollectionMetadata, canChangeValues); err != nil {
		return err
	}

	if err := ValidateTokenIdsActionPermission(permissions.CanUpdateValidTokenIds, canChangeValues); err != nil {
		return err
	}

	if err := ValidateCollectionApprovalPermissions(permissions.CanUpdateCollectionApprovals, canChangeValues); err != nil {
		return err
	}

	if err := ValidateActionPermission(permissions.CanDeleteCollection, canChangeValues); err != nil {
		return err
	}

	if err := ValidateTimedUpdatePermission(permissions.CanArchiveCollection, canChangeValues); err != nil {
		return err
	}

	return nil
}
