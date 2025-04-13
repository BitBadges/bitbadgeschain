package types

import (
	fmt "fmt"
	"math"
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

// Helper functions to reduce repetition
func createDummyUintRange() *UintRange {
	return &UintRange{
		Start: sdkmath.NewUint(math.MaxUint64),
		End:   sdkmath.NewUint(math.MaxUint64),
	}
}

func createFullAddressList() *AddressList {
	return &AddressList{
		Addresses: []string{},
		Whitelist: false,
	}
}

func IsAddressListEmpty(list *AddressList) bool {
	return len(list.Addresses) == 0 && list.Whitelist
}

func AddDefaultsIfNil(permission *UniversalPermissionDetails) *UniversalPermissionDetails {
	if permission.BadgeId == nil {
		permission.BadgeId = createDummyUintRange()
	}
	if permission.TimelineTime == nil {
		permission.TimelineTime = createDummyUintRange()
	}
	if permission.TransferTime == nil {
		permission.TransferTime = createDummyUintRange()
	}
	if permission.OwnershipTime == nil {
		permission.OwnershipTime = createDummyUintRange()
	}

	// Initialize empty address lists if nil
	if permission.ToList == nil {
		permission.ToList = createFullAddressList()
	}
	if permission.FromList == nil {
		permission.FromList = createFullAddressList()
	}
	if permission.InitiatedByList == nil {
		permission.InitiatedByList = createFullAddressList()
	}
	if permission.ApprovalIdList == nil {
		permission.ApprovalIdList = createFullAddressList()
	}

	return permission
}

func GetUintRangesWithOptions(ranges []*UintRange, uses bool) []*UintRange {
	if !uses {
		ranges = []*UintRange{{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}} //dummy range
		return ranges
	} else {
		return ranges
	}
}

func GetListIdWithOptions(listId string, uses bool) string {
	if !uses {
		listId = "All" //dummy all-inclusive list
		return listId
	} else {
		return listId
	}
}

func GetListWithOptions(list *AddressList, uses bool) *AddressList {
	if !uses {
		list = &AddressList{Addresses: []string{}, Whitelist: false} //All addresses
	}

	return list
}

func GetPermissionString(permission *UniversalPermissionDetails) string {
	var parts []string

	// Helper to add field if it's at max value
	addIfMax := func(name, value string) {
		if permission.BadgeId.Start.Equal(sdkmath.NewUint(math.MaxUint64)) ||
			permission.BadgeId.End.Equal(sdkmath.NewUint(math.MaxUint64)) {
			parts = append(parts, fmt.Sprintf("%s: %s", name, value))
		}
	}

	// Add each field
	addIfMax("badgeId", permission.BadgeId.Start.String())
	addIfMax("timelineTime", permission.TimelineTime.Start.String())
	addIfMax("transferTime", permission.TransferTime.Start.String())
	addIfMax("ownershipTime", permission.OwnershipTime.Start.String())

	// Add address lists
	if permission.ToList != nil {
		parts = append(parts, formatAddressList("toList", permission.ToList))
	}
	if permission.FromList != nil {
		parts = append(parts, formatAddressList("fromList", permission.FromList))
	}
	if permission.InitiatedByList != nil {
		parts = append(parts, formatAddressList("initiatedByList", permission.InitiatedByList))
	}

	return "(" + strings.Join(parts, " ") + ") "
}

// Helper for formatting address lists
func formatAddressList(name string, list *AddressList) string {
	result := fmt.Sprintf("%s: ", name)
	if !list.Whitelist {
		result += fmt.Sprintf("%d addresses", len(list.Addresses))
	} else {
		result += fmt.Sprintf("all except %d addresses", len(list.Addresses))
	}

	if len(list.Addresses) > 0 && len(list.Addresses) <= 5 {
		result += " (" + strings.Join(list.Addresses, " ") + ")"
	}

	return result
}

// ValidateNoMissingPermissions ensures all old permissions exist in new permissions
func ValidateNoMissingPermissions(missingPermissions []*UniversalPermissionDetails) error {
	if len(missingPermissions) == 0 {
		return nil
	}

	errMsg := fmt.Sprintf(
		"permission %sfound in old permissions but not in new permissions",
		GetPermissionString(missingPermissions[0]),
	)

	if len(missingPermissions) > 1 {
		errMsg += fmt.Sprintf(" (along with %s more)",
			sdkmath.NewUint(uint64(len(missingPermissions)-1)).String(),
		)
	}

	return sdkerrors.Wrapf(ErrInvalidPermissions, errMsg)
}

// ValidateOverlappingPermissions checks that overlapping permissions maintain consistent rules
func ValidateOverlappingPermissions(overlaps []*Overlap) error {
	for _, overlap := range overlaps {
		if err := validatePermissionTimes(overlap.FirstDetails, overlap.SecondDetails); err != nil {
			return err
		}
	}
	return nil
}

// validatePermissionTimes ensures new permissions don't invalidate old permission times
func validatePermissionTimes(oldPerm, newPerm *UniversalPermissionDetails) error {
	leftoverPermitted, _ := RemoveUintRangesFromUintRanges(
		newPerm.PermanentlyPermittedTimes,
		oldPerm.PermanentlyPermittedTimes,
	)
	leftoverForbidden, _ := RemoveUintRangesFromUintRanges(
		newPerm.PermanentlyForbiddenTimes,
		oldPerm.PermanentlyForbiddenTimes,
	)

	if len(leftoverPermitted) == 0 && len(leftoverForbidden) == 0 {
		return nil
	}

	return buildPermissionTimeError(oldPerm, leftoverPermitted, leftoverForbidden)
}

// buildPermissionTimeError constructs detailed error message for invalid permission times
func buildPermissionTimeError(perm *UniversalPermissionDetails, permittedChanges, forbiddenChanges []*UintRange) error {
	var errParts []string
	errParts = append(errParts, "permission "+GetPermissionString(perm)+"found in both new and old permissions but")

	if len(permittedChanges) > 0 {
		errParts = append(errParts, formatTimeRangeChanges(
			"previously explicitly allowed",
			"now set to disallowed",
			permittedChanges,
		))
	}

	if len(forbiddenChanges) > 0 {
		if len(permittedChanges) > 0 {
			errParts = append(errParts, "and")
		}
		errParts = append(errParts, formatTimeRangeChanges(
			"previously explicitly disallowed",
			"now set to allowed",
			forbiddenChanges,
		))
	}

	return sdkerrors.Wrapf(ErrInvalidPermissions, strings.Join(errParts, " "))
}

// formatTimeRangeChanges creates a formatted string describing time range changes
func formatTimeRangeChanges(beforeMsg, afterMsg string, times []*UintRange) string {
	ranges := make([]string, len(times))
	for i, t := range times {
		ranges[i] = t.Start.String() + "-" + t.End.String()
	}

	return fmt.Sprintf(
		"%s the times ( %s ) which are %s",
		beforeMsg,
		strings.Join(ranges, " "),
		afterMsg,
	)
}
