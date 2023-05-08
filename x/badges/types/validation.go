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

func duplicateInArray(arr []sdk.Uint) bool {
	visited := make(map[sdk.Uint]bool, 0)
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


func DoRangesOverlap(ids []*IdRange) bool {
	//Insertion sort in order of range.Start. If two have same range.Start, sort by range.End.
	var n = len(ids)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			if ids[j-1].Start.GT(ids[j].Start) {
				ids[j-1], ids[j] = ids[j], ids[j-1]
			} else if ids[j-1].Start.Equal(ids[j].Start) && ids[j-1].End.GT(ids[j].End) {
				ids[j-1], ids[j] = ids[j], ids[j-1]
			}
			j = j - 1
		}
	}

	//Check if any overlap
	for i := 1; i < n; i++ {
		prevInsertedRange := ids[i-1]
		currRange := ids[i]

		if currRange.Start.LTE(prevInsertedRange.End) {
			return true
		}
	}

	return false
}

// Validates ranges are valid. If end.IsZero(), we assume end == start.
func ValidateRangesAreValid(badgeIdRanges []*IdRange) error {
	for _, badgeIdRange := range badgeIdRanges {
		if badgeIdRange.Start.IsNil() || badgeIdRange.End.IsNil() {
			return ErrUintUnititialized
		}

		if badgeIdRange == nil {
			return ErrRangesIsNil
		}

		if badgeIdRange.Start.GT(badgeIdRange.End) {
			return ErrStartGreaterThanEnd
		}
	}

	overlap := DoRangesOverlap(badgeIdRanges)
	if overlap {
		return ErrRangesOverlap
	}

	return nil
}

// Validates no element is X
func ValidateNoElementIsX(amounts []sdk.Uint, x sdk.Uint) error {
	for _, amount := range amounts {
		if amount.Equal(x) {
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

func ValidateAddressesMapping(addressesMapping AddressesMapping, allowMintAddress bool) error {
	if addressesMapping.ManagerOptions.IsNil() {
		return ErrUintUnititialized
	}

	for _, address := range addressesMapping.Addresses {



		if allowMintAddress && address == MintAddress {
			continue
		}

		if err := ValidateAddress(address); err != nil {
			return err
		}
	}

	return nil
}


func ValidateTransferMapping(transferMapping TransferMapping) error {
	if err := ValidateAddressesMapping(*transferMapping.To, true); err != nil {
		return err
	}

	if err := ValidateAddressesMapping(*transferMapping.From, false); err != nil {
		return err
	}

	return nil
}	

func ValidateClaim(claim *Claim) error {
	err := *new(error)

	if claim.NumClaimsPerAddress.IsNil() || claim.IncrementIdsBy.IsNil() {
		return ErrUintUnititialized
	}

	if claim.Uri != "" {
		err = ValidateURI(claim.Uri)
		if err != nil {
			return err
		}
	}

	for _, challenge := range claim.Challenges {
		if challenge.ExpectedProofLength.IsNil() {
			return ErrUintUnititialized
		}
	}

	err = ValidateBalances(claim.UndistributedBalances)
	if err != nil {
		return err
	}

	err = ValidateBalances(claim.CurrentClaimAmounts)
	if err != nil {
		return err
	}


	if claim.TimeRange == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
	}

	err = ValidateRangesAreValid([]*IdRange{claim.TimeRange})
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
	}


	err = ValidateRangesAreValid([]*IdRange{claim.TimeRange})
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid time range")
	}



	return nil
}

func ValidateBalances(balances []*Balance) error {
	for _, balance := range balances {
		if balance == nil {
			return ErrInvalidLengthBalances
		}

		if balance.Amount.IsZero() || balance.Amount.IsNil() {
			return ErrAmountEqualsZero
		}

		err := ValidateRangesAreValid(balance.BadgeIds)
		if err != nil {
			return err
		}
	}

	return nil
}


func ValidateTransfer(transfer *Transfer) error {
	err := *new(error)
	err = ValidateBalances(transfer.Balances)
	if err != nil {
		return err
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

func ValidateBadgeUris(badgeUris []*BadgeUri) error {
	err := *new(error)
	if badgeUris != nil && len(badgeUris) > 0 {
		for _, badgeUri := range badgeUris {
			//Validate well-formedness of the message entries
			if err := ValidateURI(badgeUri.Uri); err != nil {
				return err
			}

			err = ValidateRangesAreValid(badgeUri.BadgeIds)
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid badgeIds")
			}
		}
	}
	return nil
}

//IMPORTANT: Note this was copied from the keeper id_range.go file. If you change this, change that as well and vice versa.

// Will sort the ID ranges in order and merge overlapping IDs if we can
func SortAndMergeOverlapping(ids []*IdRange) []*IdRange {
	//Insertion sort in order of range.Start. If two have same range.Start, sort by range.End.
	var n = len(ids)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			if ids[j-1].Start.GT(ids[j].Start) {
				ids[j-1], ids[j] = ids[j], ids[j-1]
			} else if ids[j-1].Start.Equal(ids[j].Start) && ids[j-1].End.GT(ids[j].End) {
				ids[j-1], ids[j] = ids[j], ids[j-1]
			}
			j = j - 1
		}
	}

	//Merge overlapping ranges
	if n > 0 {
		newIdRanges := []*IdRange{ids[0]}
		//Iterate through and compare with previously inserted range
		for i := 1; i < n; i++ {
			prevInsertedRange := newIdRanges[len(newIdRanges)-1]
			currRange := ids[i]

			if currRange.Start.Equal(prevInsertedRange.Start) {
				//Both have same start, so we set to currRange.End because currRange.End is greater due to our sorting
				//Example: prevRange = [1, 5], currRange = [1, 10] -> newRange = [1, 10]
				newIdRanges[len(newIdRanges)-1].End = currRange.End
			} else if currRange.End.GT(prevInsertedRange.End) {
				//We have different starts and curr end is greater than prev end
				
				
				if currRange.Start.GT(prevInsertedRange.End.AddUint64(1)) {
					//We have a gap between the prev range end and curr range start, so we just append currRange
					//Example: prevRange = [1, 5], currRange = [7, 10] -> newRange = [1, 5], [7, 10]
					newIdRanges = append(newIdRanges, currRange)
				} else {
					//They overlap and we can merge them
					//Example: prevRange = [1, 5], currRange = [2, 10] -> newRange = [1, 10]
					newIdRanges[len(newIdRanges)-1].End = currRange.End
				}
			} else {
				//Note: If currRange.End <= prevInsertedRange.End, it is already fully contained within the previous. We can just continue.
			}
		}
		return newIdRanges
	} else {
		return ids
	}
}