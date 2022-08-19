package types

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	// URI must be a valid URI. Method <= 10 characters long. Path <= 90 characters long.
	reUriString = `\w{0,10}:(\/?\/?)[^\s]{0,90}`
	reUri       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reUriString))
)

// Validate uri and subasset uri returns whether both the uri and subasset uri is valid. Max 100 characters each.
func ValidateURI(uriObject UriObject) error {
	uri, err := GetUriFromUriObject(uriObject)
	if err != nil || !reUri.MatchString(uri) {
		return sdkerrors.Wrapf(ErrInvalidBadgeURI, "invalid uri: %s", uri)
	}

	subassetUri, err := GetSubassetUriFromUriObject(uriObject)
	if err != nil || !reUri.MatchString(subassetUri) {
		return sdkerrors.Wrapf(ErrInvalidBadgeURI, "invalid uri: %s", uri)
	}

	return nil
}

func ValidateAddress(address string) error {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}

// Validate bytes we store are valid. We don't allow users to store anything > 256 bytes in a badge.
func ValidateBytes(bytesToCheck []byte) error {
	if len(bytesToCheck) > 256 {
		return ErrBytesGreaterThan256
	}
	return nil
}

//Validates ranges are valid. If end == 0, we assume end == start.
func ValidateRangesAreValid(subbadgeRanges []*IdRange) error {

	for _, subbadgeRange := range subbadgeRanges {
		if subbadgeRange == nil {
			return ErrRangesIsNil
		}

		if subbadgeRange.End == 0 {
			subbadgeRange.End = subbadgeRange.Start
		}

		if subbadgeRange.Start > subbadgeRange.End {
			return ErrStartGreaterThanEnd
		}
	}
	return nil
}

//Validates no element is X
func ValidateNoElementIsX(amounts []uint64, x uint64) error {
	for _, amount := range amounts {
		if amount == x {
			return ErrElementCantEqualThis
		}
	}
	return nil
}
