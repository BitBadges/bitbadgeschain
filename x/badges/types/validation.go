package types

import (
	"fmt"
	"regexp"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	// URI must be a valid URI. Method <= 10 characters long. Path <= 90 characters long.
	reUriString = `\w+:(\/?\/?)[^\s]+`
	reUri       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reUriString))
)

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

func ValidateAddress(address string, allowAliases bool) error {
	if allowAliases && (address == "Manager" || address == "Mint") {
		return nil
	}

	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid address (%s)", err)
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
func ValidateRangesAreValid(badgeIdRanges []*IdRange, errorOnEmpty bool) error {
	if len(badgeIdRanges) == 0 {
		if errorOnEmpty {
			return sdkerrors.Wrapf(ErrInvalidIdRangeSpecified, "these id ranges can not be empty (length == 0)")
		} 
	}


	for _, badgeIdRange := range badgeIdRanges {
		if badgeIdRange == nil {
			return ErrRangesIsNil
		}

		if badgeIdRange.Start.IsNil() || badgeIdRange.End.IsNil() {
			return sdkerrors.Wrapf(ErrUintUnititialized, "id range start and/or end is uninitialized")
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
func ValidateNoElementIsX(amounts []sdkmath.Uint, x sdkmath.Uint) error {
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

func ValidateAddressMapping(addressMapping *AddressMapping) error {
	//TODO:

	// for _, address := range addressMapping.Addresses {
	// 	if err := ValidateAddress(address); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

//TODO: tracker ID must be defined if approvals are not nil
func ValidateUserApprovedOutgoingTransfer(userApprovedOutgoingTransfer *UserApprovedOutgoingTransfer) error {
	//TODO:
	return nil
}

//TODO:
func ValidateUserApprovedIncomingTransfer(userApprovedIncomingTransfer *UserApprovedIncomingTransfer) error {
	return nil
}

func ValidateCollectionApprovedTransfer(collectionApprovedTransfer *CollectionApprovedTransfer) error {
	
// 	approvalAmountsArr := collectionApprovedTransfer.AmountRestrictions
// 	if approvalAmountsArr == nil {
// 		return ErrAmountRestrictionsIsNil
// 	}

// 	for idx, allowedCombination := range collectionApprovedTransfer.AllowedCombinations {
// 		for _, compCombination := range collectionApprovedTransfer.AllowedCombinations[idx+1:] {
// 			if allowedCombination.InvertBadgeIds == compCombination.InvertBadgeIds &&
// 				allowedCombination.InvertTransferTimes == compCombination.InvertTransferTimes &&
// 				allowedCombination.InvertTo == compCombination.InvertTo &&
// 				allowedCombination.InvertFrom == compCombination.InvertFrom &&
// 				allowedCombination.InvertInitiatedBy == compCombination.InvertInitiatedBy {
// 				return ErrInvalidCombinations
// 			}
// 		}
// 	}

// 	for _, approvalAmounts := range approvalAmountsArr {
		
		
// 		if approvalAmounts != nil {
// 			if err := ValidateRangesAreValid(approvalAmounts.BalancesTimes, true); err != nil {
// 				return sdkerrors.Wrapf(err, "invalid balances times")
// 			}
	
			
// 			if approvalAmounts.Amount.IsNil() {
// 				return sdkerrors.Wrapf(ErrUintUnititialized, "amount is uninitialized")
// 			}

// 			if approvalAmounts.MaxNumTransfers.IsNil() {
// 				return sdkerrors.Wrapf(ErrUintUnititialized, "max num transfers is uninitialized")
// 			}

// 			if approvalAmounts.FromRestrictions != nil {
// 				if approvalAmounts.FromRestrictions.AmountPerAddress.IsNil() {
// 					return sdkerrors.Wrapf(ErrUintUnititialized, "amount per address is uninitialized")
// 				}

// 				if approvalAmounts.FromRestrictions.TransfersPerAddress.IsNil() {
// 					return sdkerrors.Wrapf(ErrUintUnititialized, "transfers per address is uninitialized")
// 				}
// 			}

// 			if approvalAmounts.ToRestrictions != nil {
// 				if approvalAmounts.ToRestrictions.AmountPerAddress.IsNil() {
// 					return sdkerrors.Wrapf(ErrUintUnititialized, "amount per address is uninitialized")
// 				}

// 				if approvalAmounts.ToRestrictions.TransfersPerAddress.IsNil() {
// 					return sdkerrors.Wrapf(ErrUintUnititialized, "transfers per address is uninitialized")
// 				}
// 			}

// 			if approvalAmounts.InitiatedByRestrictions != nil {
// 				if approvalAmounts.InitiatedByRestrictions.AmountPerAddress.IsNil() {
// 					return sdkerrors.Wrapf(ErrUintUnititialized, "amount per address is uninitialized")
// 				}

// 				if approvalAmounts.InitiatedByRestrictions.TransfersPerAddress.IsNil() {
// 					return sdkerrors.Wrapf(ErrUintUnititialized, "transfers per address is uninitialized")
// 				}
// 			}
// 		}
// 	}

// 	if err := ValidateRangesAreValid(collectionApprovedTransfer.TransferTimes, false); err != nil {
// 		return sdkerrors.Wrapf(err, "invalid transfer times")
// 	}

// 	if err := ValidateRangesAreValid(collectionApprovedTransfer.BadgeIds, false); err != nil {
// 		return sdkerrors.Wrapf(err, "invalid badge ids")
// 	}

// 	if err := ValidateClaim(collectionApprovedTransfer.Claim); err != nil {
// 		return sdkerrors.Wrapf(err, "invalid claim")
// 	}

// 	if collectionApprovedTransfer.TransferTrackerId.IsNil() {
// 		return sdkerrors.Wrapf(ErrUintUnititialized, "transfer tracker id is uninitialized")
// 	}

// 	return nil
// }

// func ValidateClaim(claim *Claim) error {
// 	err := *new(error)

// 	if claim.Uri != "" {
// 		err = ValidateURI(claim.Uri)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	if claim.IncrementIdsBy.IsNil() {
// 		return sdkerrors.Wrapf(ErrUintUnititialized, "increment ids by is uninitialized")
// 	}

// 	hasOrderMatters := false
// 	for _, challenge := range claim.Challenges {
// 		if challenge.ExpectedProofLength.IsNil() {
// 			return sdkerrors.Wrapf(ErrUintUnititialized, "expected proof length is uninitialized")
// 		}

// 		if challenge.UseLeafIndexForDistributionOrder {
// 			if hasOrderMatters {
// 				return ErrCanOnlyUseLeafIndexForBadgeIdsOnce
// 			}

// 			hasOrderMatters = true
// 		}

// 		if !challenge.MaxOneUsePerLeaf && challenge.UseLeafIndexForDistributionOrder {
// 			return ErrPrimaryChallengeMustBeOneUsePerLeaf
// 		}

// 		if !challenge.MaxOneUsePerLeaf && !challenge.UseCreatorAddressAsLeaf {
// 			return ErrCanOnlyUseMaxOneUsePerLeafWithWhitelistTree
// 		}
// 	}

// 	err = ValidateBalances(claim.StartAmounts)
// 	if err != nil {
// 		return sdkerrors.Wrapf(err, "invalid start amounts")
// 	}

// 	return nil
	return nil
}

func ValidateBalances(balances []*Balance) ([]*Balance, error) {
	err := *new(error)
	for _, balance := range balances {
		if balance == nil {
			return balances, sdkerrors.Wrapf(ErrInvalidLengthBalances, "balances is nil")
		}

		if balance.Amount.IsNil() {
			return balances, sdkerrors.Wrapf(ErrAmountEqualsZero, "amount is uninitialized")
		}

		err = ValidateRangesAreValid(balance.BadgeIds, true)
		if err != nil {
			return balances, sdkerrors.Wrapf(err, "invalid balance badge ids")
		}

		err = ValidateRangesAreValid(balance.Times, true)
		if err != nil {
			return balances, sdkerrors.Wrapf(err, "invalid balance times")
		}
	}

	balances, err = HandleDuplicateBadgeIds(balances)
	if err != nil {
		return balances, err
	}

	return balances, nil
}

func ValidateTransfer(transfer *Transfer) error {
	err := *new(error)

	transfer.Balances, err = ValidateBalances(transfer.Balances)
	if err != nil {
		return err
	}

	err = ValidateNoStringElementIsX(transfer.ToAddresses, transfer.From)
	if err != nil {
		return ErrSenderAndReceiverSame
	}

	err = ValidateAddress(transfer.From, true)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid from address (%s)", err)
	}

	if duplicateInStringArray(transfer.ToAddresses) {
		return ErrDuplicateAddresses
	}

	for _, address := range transfer.ToAddresses {
		err = ValidateAddress(address, true)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddress, "invalid to address (%s)", err)
		}
	}

	return nil
}

func ValidateBadgeMetadata(badgeMetadata []*BadgeMetadata) error {
	err := *new(error)


	handledBadgeIds := []*IdRange{}
	if len(badgeMetadata) > 0 {
		for _, badgeMetadata := range badgeMetadata {
			//Validate well-formedness of the message entries
			if err := ValidateURI(badgeMetadata.Uri); err != nil {
				return err
			}

			err = ValidateRangesAreValid(badgeMetadata.BadgeIds, true)
			if err != nil {
				return sdkerrors.Wrapf(ErrInvalidRequest, "invalid badgeIds")
			}

			badgeMetadata.BadgeIds = SortAndMergeOverlapping(badgeMetadata.BadgeIds)

			if err := AssertRangesDoNotOverlapAtAll(handledBadgeIds, badgeMetadata.BadgeIds); err != nil {
				return sdkerrors.Wrapf(err, "badge metadata has duplicate badge ids")
			}

			handledBadgeIds = append(handledBadgeIds, badgeMetadata.BadgeIds...)
		}
	}
	return nil
}


func ValidateInheritedBalances(inheritedBalances []*InheritedBalance) error {
	err := *new(error)


	handledBadgeIds := []*IdRange{}
	if len(inheritedBalances) > 0 {
		for _, inheritedBalance := range inheritedBalances {
			err = ValidateRangesAreValid(inheritedBalance.BadgeIds, true)
			if err != nil {
				return sdkerrors.Wrapf(ErrInvalidRequest, "invalid badgeIds")
			}

			err = ValidateRangesAreValid(inheritedBalance.ParentBadgeIds, true)
			if err != nil {
				return sdkerrors.Wrapf(ErrInvalidRequest, "invalid badgeIds")
			}

			inheritedBalance.BadgeIds = SortAndMergeOverlapping(inheritedBalance.BadgeIds)
			inheritedBalance.ParentBadgeIds = SortAndMergeOverlapping(inheritedBalance.ParentBadgeIds)

			if err := AssertRangesDoNotOverlapAtAll(handledBadgeIds, inheritedBalance.BadgeIds); err != nil {
				return sdkerrors.Wrapf(err, "inherited balances has duplicate badge ids")
			}

			handledBadgeIds = append(handledBadgeIds, inheritedBalance.BadgeIds...)


			if inheritedBalance.ParentCollectionId.IsNil() || inheritedBalance.ParentCollectionId.IsZero() {
				return sdkerrors.Wrapf(ErrUintUnititialized, "parent collection id is uninitialized")
			}

			totalBadgeLength := sdk.NewUint(0)
			for _, badgeIdRange := range inheritedBalance.BadgeIds {
				totalBadgeLength = totalBadgeLength.Add(badgeIdRange.End.Sub(badgeIdRange.Start).AddUint64(1))
			}

			totalParentBadgeLength := sdk.NewUint(0)
			for _, badgeIdRange := range inheritedBalance.ParentBadgeIds {
				totalParentBadgeLength = totalParentBadgeLength.Add(badgeIdRange.End.Sub(badgeIdRange.Start).AddUint64(1))
			}

			if !totalParentBadgeLength.Equal(sdk.NewUint(1)) {
				if !totalParentBadgeLength.Equal(totalBadgeLength) {
					return ErrInvalidInheritedBadgeLength
				}
			}
		}
	}
	return nil
}