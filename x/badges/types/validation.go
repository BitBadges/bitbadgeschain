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

func duplicateInArray(arr []uint64) bool {
	visited := make(map[uint64]bool, 0)
	for i := 0; i < len(arr); i++ {

		if visited[arr[i]] == true {
			return true
		} else {
			visited[arr[i]] = true
		}
	}
	return false
}

func duplicateInStringArray(arr []string) bool {
	visited := make(map[string]bool, 0)
	for i := 0; i < len(arr); i++ {
		if visited[arr[i]] {
			return true
		} else {
			visited[arr[i]] = true
		}
	}
	return false
}

// Validate uri and subasset uri returns whether both the uri and subasset uri is valid. Max 100 characters each.
func ValidateURI(uri string) error {
	regexMatch := reUri.MatchString(uri)
	if !regexMatch {
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
func ValidateBytes(bytesToCheck string) error {
	if len(bytesToCheck) > 256 {
		return ErrBytesGreaterThan256
	}
	return nil
}

// Validates ranges are valid. If end == 0, we assume end == start.
func ValidateRangesAreValid(badgeIdRanges []*IdRange) error {
	for _, badgeIdRange := range badgeIdRanges {
		if badgeIdRange == nil {
			return ErrRangesIsNil
		}

		if badgeIdRange.Start > badgeIdRange.End {
			return ErrStartGreaterThanEnd
		}
	}
	return nil
}

func ValidateActionsAreValid(actions []uint64, rangesLength int) error {
	for _, action := range actions {
		if action > 2 {
			return ErrActionOutOfRange
		}
	}

	if len(actions) == 0 {
		return ErrActionsEmpty
	}

	if len(actions) == 1 {
		return nil
	} else if len(actions) != rangesLength {
		return ErrActionsLengthNotEqualToRangesLength
	}

	return nil
}

// Validates no element is X
func ValidateNoElementIsX(amounts []uint64, x uint64) error {
	for _, amount := range amounts {
		if amount == x {
			return ErrElementCantEqualThis
		}
	}
	return nil
}


// Validates no element is X
func ValidateNoStringElementIsX(addresses []string, x string) error {
	for _, amount := range addresses {
		if amount == x {
			return ErrElementCantEqualThis
		}
	}
	return nil
}

func ValidateAddressesMapping(addressesMapping AddressesMapping) error {
	for _, address := range addressesMapping.Addresses {
		if err := ValidateAddress(address); err != nil {
			return err
		}
	}

	return nil
}


func ValidateTransferMapping(transferMapping TransferMapping) error {
	if err := ValidateAddressesMapping(*transferMapping.To); err != nil {
		return err
	}

	if err := ValidateAddressesMapping(*transferMapping.From); err != nil {
		return err
	}

	return nil
}	

func ValidateClaim(claim *Claim) error {
	err := *new(error)

	if claim.Uri != "" {
		err = ValidateURI(claim.Uri)
		if err != nil {
			return err
		}
	}

	if claim.TimeRange == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
	}

	err = ValidateRangesAreValid([]*IdRange{claim.TimeRange})
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
	}

	for _, balance := range claim.Balances {
		if balance == nil {
			return ErrInvalidLengthBalances
		}

		if balance.Amount == 0 {
			return ErrAmountEqualsZero
		}

		err = ValidateRangesAreValid(balance.BadgeIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateTransfer(transfer *Transfers) error {
	err := *new(error)
	for _, balance := range transfer.Balances {
		if balance == nil {
			return ErrInvalidLengthBalances
		}

		if balance.Amount == 0 {
			return ErrAmountEqualsZero
		}

		err = ValidateRangesAreValid(balance.BadgeIds)
		if err != nil {
			return err
		}
	}

	if duplicateInStringArray(transfer.ToAddresses) {
		return ErrDuplicateAddresses
	}

	for _, address := range transfer.ToAddresses {
		_, err = sdk.AccAddressFromBech32(address)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid to address (%s)", err)
		}
	}

	return nil
}