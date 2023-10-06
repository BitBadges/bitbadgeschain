package types

import (
	"fmt"
	"math"
	"regexp"
	"strings"

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
		return sdkerrors.Wrapf(ErrInvalidURI, "invalid uri: %s", uri)
	}

	return nil
}

func ValidateAddress(address string, allowAliases bool) error {
	if allowAliases && (address == "Mint") {
		return nil
	}

	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}

func DoRangesOverlap(ids []*UintRange) bool {
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
func ValidateRangesAreValid(badgeUintRanges []*UintRange, allowAllUints bool, errorOnEmpty bool) error {
	if len(badgeUintRanges) == 0 {
		if errorOnEmpty {
			return sdkerrors.Wrapf(ErrInvalidUintRangeSpecified, "these id ranges can not be empty (length == 0)")
		}
	}

	for _, badgeUintRange := range badgeUintRanges {
		if badgeUintRange == nil {
			return ErrRangesIsNil
		}

		if badgeUintRange.Start.IsNil() || badgeUintRange.End.IsNil() {
			return sdkerrors.Wrapf(ErrUintUnititialized, "id range start and/or end is nil")
		}

		if badgeUintRange.Start.GT(badgeUintRange.End) {
			return ErrStartGreaterThanEnd
		}
		

		if !allowAllUints {
			if badgeUintRange.Start.IsZero() || badgeUintRange.End.IsZero() {
				return sdkerrors.Wrapf(ErrUintUnititialized, "id range start and/or end is zero")
			}

			if badgeUintRange.Start.GT(sdkmath.NewUint(math.MaxUint64)) || badgeUintRange.End.GT(sdkmath.NewUint(math.MaxUint64)) {
				return ErrUintGreaterThanMax
			}
		}
	}

	overlap := DoRangesOverlap(badgeUintRanges)
	if overlap {
		return ErrRangesOverlap
	}

	return nil
}

// Validates no element is X
func ValidateNoElementIsX(amounts []sdkmath.Uint, x sdkmath.Uint) error {
	for _, amount := range amounts {
		if amount.Equal(x) {
			return sdkerrors.Wrapf(ErrElementCantEqualThis, "amount can not equal %s", x.String())
		}
	}
	return nil
}

// Validates no element is X
func ValidateNoStringElementIsX(addresses []string, x string) error {
	for _, amount := range addresses {
		if amount == x {
			return sdkerrors.Wrapf(ErrElementCantEqualThis, "address can not equal %s", x)
		}
	}
	return nil
}

func ValidateAddressMapping(addressMapping *AddressMapping) error {
	if addressMapping.MappingId == "" ||
		addressMapping.MappingId == "Mint" ||
		addressMapping.MappingId == "Manager" ||
		addressMapping.MappingId == "AllWithoutMint" ||
		addressMapping.MappingId == "None" {
		return sdkerrors.Wrapf(ErrInvalidAddress, "mapping id is uninitialized")
	}

	if err := ValidateAddress(addressMapping.MappingId, false); err == nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "mapping id can not be a valid address")
	}

	if strings.Contains(addressMapping.MappingId, ":") || strings.Contains(addressMapping.MappingId, "!") {
		return sdkerrors.Wrapf(ErrInvalidAddress, "mapping id can not contain : or !")
	}

	if addressMapping.Uri != "" {
		if err := ValidateURI(addressMapping.Uri); err != nil {
			return err
		}
	}

	for _, address := range addressMapping.Addresses {
		if err := ValidateAddress(address, false); err != nil {
			return err
		}
	}

	return nil
}

func ValidateUserOutgoingApprovals(userOutgoingApprovals []*UserOutgoingApproval, fromAddress string) error {
	castedTransfers := CastOutgoingTransfersToCollectionTransfers(userOutgoingApprovals, fromAddress)
	err := ValidateCollectionApprovals(castedTransfers)
	return err
}

func ValidateUserIncomingApprovals(userIncomingApprovals []*UserIncomingApproval, toAddress string) error {
	castedTransfers := CastIncomingTransfersToCollectionTransfers(userIncomingApprovals, toAddress)
	err := ValidateCollectionApprovals(castedTransfers)
	return err
}

func ValidateCollectionApprovals(collectionApprovals []*CollectionApproval) error {
	for i := 0; i < len(collectionApprovals); i++ {
		if collectionApprovals[i].IsApproved && collectionApprovals[i].ApprovalId == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id is uninitialized")
		}

		if collectionApprovals[i].ApprovalId == "default-outgoing" || collectionApprovals[i].ApprovalId == "default-incoming" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id can not be default-outgoing or default-incoming")
		}

		for j := i + 1; j < len(collectionApprovals); j++ {
			if !collectionApprovals[i].IsApproved || !collectionApprovals[j].IsApproved {
				continue
			}
			
			if collectionApprovals[i].ApprovalId == collectionApprovals[j].ApprovalId {
				return sdkerrors.Wrapf(ErrInvalidRequest, "duplicate approval ids")
			}
		}
	}

	for _, collectionApproval := range collectionApprovals {
		if collectionApproval == nil {
			return sdkerrors.Wrapf(ErrInvalidRequest, "collection approved transfer is nil")
		}

		if collectionApproval.FromMappingId == "" {
			return sdkerrors.Wrapf(ErrInvalidAddress, "from mapping id is uninitialized")
		}

		if collectionApproval.ToMappingId == "" {
			return sdkerrors.Wrapf(ErrInvalidAddress, "to mapping id is uninitialized")
		}

		if collectionApproval.InitiatedByMappingId == "" {
			return sdkerrors.Wrapf(ErrInvalidAddress, "initiated by mapping id is uninitialized")
		}

		if err := ValidateRangesAreValid(collectionApproval.BadgeIds, false, false); err != nil {
			return sdkerrors.Wrapf(err, "invalid badge IDs")
		}

		if err := ValidateRangesAreValid(collectionApproval.TransferTimes, false, false); err != nil {
			return sdkerrors.Wrapf(err, "invalid transfer times")
		}

		if err := ValidateRangesAreValid(collectionApproval.OwnershipTimes, false, false); err != nil {
			return sdkerrors.Wrapf(err, "invalid transfer times")
		}

		if collectionApproval.Uri != "" {
			if err := ValidateURI(collectionApproval.Uri); err != nil {
				return err
			}
		}

		
		approvalCriteria := collectionApproval.ApprovalCriteria
		if approvalCriteria != nil {
			usingLeafIndexForTransferOrder := false
			if (approvalCriteria.PredeterminedBalances != nil && approvalCriteria.PredeterminedBalances.OrderCalculationMethod != nil && approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex) {
				usingLeafIndexForTransferOrder = true
			}
			if err := ValidateMerkleChallenge(approvalCriteria.MerkleChallenge, collectionApproval.ChallengeTrackerId, usingLeafIndexForTransferOrder); err != nil {
				return sdkerrors.Wrapf(err, "invalid challenges")
			}

			if collectionApproval.ApprovalTrackerId == "All" {
				return sdkerrors.Wrapf(ErrInvalidRequest, "approval tracker id can not be All")
			}

			if strings.Contains(collectionApproval.ApprovalTrackerId, ":") || strings.Contains(collectionApproval.ApprovalTrackerId, "!") {
				return sdkerrors.Wrapf(ErrIdsContainsInvalidChars, "approval tracker id can not contain : or !")
			}
			
			for _, mustOwnBadgeBalance := range approvalCriteria.MustOwnBadges {
				if mustOwnBadgeBalance == nil {
					return sdkerrors.Wrapf(ErrInvalidRequest, "mustOwnBadges balance is nil")
				}

				if err := ValidateRangesAreValid(mustOwnBadgeBalance.BadgeIds, false, false); err != nil {
					return sdkerrors.Wrapf(err, "invalid badge IDs")
				}

				if err := ValidateRangesAreValid(mustOwnBadgeBalance.OwnershipTimes, false, false); err != nil {
					return sdkerrors.Wrapf(err, "invalid owned times")
				}

				if err := ValidateRangesAreValid([]*UintRange{mustOwnBadgeBalance.AmountRange}, true, true); err != nil {
					return sdkerrors.Wrapf(err, "invalid transfer times")
				}

				if mustOwnBadgeBalance.CollectionId.IsNil() || mustOwnBadgeBalance.CollectionId.IsZero() {
					return sdkerrors.Wrapf(ErrUintUnititialized, "collection id is uninitialized")
				}
			}
			
			
	
			if approvalCriteria.PredeterminedBalances != nil {
				orderCalculationMethodIsBasicallyNil := !approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex &&
					!approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers &&
					!approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers &&
					!approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers &&
					!approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers

				sequentialTransferIsBasicallyNil :=
				approvalCriteria.PredeterminedBalances.IncrementedBalances == nil || (approvalCriteria.PredeterminedBalances.IncrementedBalances.StartBalances == nil &&
					approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementBadgeIdsBy.IsZero() &&
					approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementOwnershipTimesBy.IsZero())


					manualBalancesIsBasicallyNil := approvalCriteria.PredeterminedBalances.ManualBalances == nil

				isBasicallyNil := orderCalculationMethodIsBasicallyNil && sequentialTransferIsBasicallyNil && manualBalancesIsBasicallyNil

				if (!isBasicallyNil) {
					orderType := approvalCriteria.PredeterminedBalances.OrderCalculationMethod
					if orderType == nil {
						return sdkerrors.Wrapf(ErrInvalidRequest, "order type is nil")
					}
		
					numTrue := 0
					if orderType.UseMerkleChallengeLeafIndex {
						numTrue++
					}
		
					if orderType.UseOverallNumTransfers {
						numTrue++
					}
		
					if orderType.UsePerToAddressNumTransfers {
						numTrue++
					}
		
					if orderType.UsePerFromAddressNumTransfers {
						numTrue++
					}
		
					if orderType.UsePerInitiatedByAddressNumTransfers {
						numTrue++
					}
		
		
					if numTrue != 1 {
						return sdkerrors.Wrapf(ErrInvalidRequest, "only one of use challenge leaf index, use overall num transfers, use per to address num transfers, use per from address num transfers, use per initiated by address num transfers can be true")
					}
					
		
					
		
					err := *new(error)
					if (manualBalancesIsBasicallyNil && !sequentialTransferIsBasicallyNil) {
						sequentialTransfer := approvalCriteria.PredeterminedBalances.IncrementedBalances 
						sequentialTransfer.StartBalances, err = ValidateBalances(sequentialTransfer.StartBalances)
						if err != nil {
							return err
						}
		
						if sequentialTransfer.IncrementBadgeIdsBy.IsNil() {
							return sdkerrors.Wrapf(ErrUintUnititialized, "increment ids by is uninitialized")
						}
					
						if sequentialTransfer.IncrementOwnershipTimesBy.IsNil() {
							return sdkerrors.Wrapf(ErrUintUnititialized, "max num transfers is uninitialized")
						}
					} else if (!manualBalancesIsBasicallyNil && sequentialTransferIsBasicallyNil) {
						for _, manualTransfer := range approvalCriteria.PredeterminedBalances.ManualBalances {
							manualTransfer.Balances, err = ValidateBalances(manualTransfer.Balances)
							if err != nil {
								return err
							}
						}
					} else {
						return sdkerrors.Wrapf(ErrInvalidRequest, "manual transfers and sequential transfers cannot be both nil or both defined")
					}
				}

				if approvalCriteria.ApprovalAmounts == nil {
					return sdkerrors.Wrapf(ErrInvalidRequest, "approval amounts is uninitialized")
				}

				if approvalCriteria.MaxNumTransfers == nil {
					return sdkerrors.Wrapf(ErrInvalidRequest, "max num transfers must not be nil")
				}
			} else {
				approvalCriteria.PredeterminedBalances = nil
			}
		}
	}

	return nil
}

func ValidateMerkleChallenge(challenge *MerkleChallenge, challengeId string, usingLeafIndexForTransferOrder bool) error {
	if challenge == nil || challenge.Root == "" {
		return nil
	}

	if challenge.ExpectedProofLength.IsNil() {
		return sdkerrors.Wrapf(ErrUintUnititialized, "expected proof length is uninitialized")
	}

	if challengeId == "All" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "approval tracker id can not be All")
	}

	if strings.Contains(challengeId, ":") || strings.Contains(challengeId, "!") {
		return sdkerrors.Wrapf(ErrIdsContainsInvalidChars, "approval tracker id can not contain : or !")
	}

	if !challenge.MaxOneUsePerLeaf && usingLeafIndexForTransferOrder {
		return ErrPrimaryChallengeMustBeOneUsePerLeaf
	}

	if !challenge.MaxOneUsePerLeaf && !challenge.UseCreatorAddressAsLeaf {
		return ErrCanOnlyUseMaxOneUsePerLeafWithWhitelistTree
	}
	

	return nil
}

func ValidateBalances(balances []*Balance) ([]*Balance, error) {
	err := *new(error)
	for _, balance := range balances {
		if balance == nil {
			return balances, sdkerrors.Wrapf(ErrInvalidLengthBalances, "balances is nil")
		}

		if balance.Amount.IsNil() || balance.Amount.IsZero() {
			return balances, sdkerrors.Wrapf(ErrAmountEqualsZero, "amount is uninitialized")
		}

		err = ValidateRangesAreValid(balance.BadgeIds, false, true)
		if err != nil {
			return balances, sdkerrors.Wrapf(err, "invalid balance badge ids")
		}

		err = ValidateRangesAreValid(balance.OwnershipTimes, false, true)
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
		err = ValidateAddress(address, false)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddress, "invalid to address (%s)", err)
		}
	}

	if transfer.PrecalculateBalancesFromApproval != nil {
		if transfer.PrecalculateBalancesFromApproval.ApprovalLevel == "" && transfer.PrecalculateBalancesFromApproval.ApproverAddress == "" && transfer.PrecalculateBalancesFromApproval.ApprovalId == "" {
			//basically nil
		} else {
			if transfer.PrecalculateBalancesFromApproval.ApprovalLevel != "collection" && transfer.PrecalculateBalancesFromApproval.ApprovalLevel != "incoming" && transfer.PrecalculateBalancesFromApproval.ApprovalLevel != "outgoing" {
				return sdkerrors.Wrapf(ErrInvalidRequest, "approval level must be collection, incoming, or outgoing")
			}

			if transfer.PrecalculateBalancesFromApproval.ApproverAddress != "" {
				if err := ValidateAddress(transfer.PrecalculateBalancesFromApproval.ApproverAddress, false); err != nil {
					return sdkerrors.Wrapf(ErrInvalidAddress, "invalid approval id address (%s)", err)
				}
			}
		}
	}

	return nil
}

func ValidateBadgeMetadata(badgeMetadata []*BadgeMetadata) error {
	err := *new(error)

	handledBadgeIds := []*UintRange{}
	if len(badgeMetadata) > 0 {
		for _, badgeMetadata := range badgeMetadata {
			//Validate well-formedness of the message entries
			if err := ValidateURI(badgeMetadata.Uri); err != nil {
				return err
			}

			err = ValidateRangesAreValid(badgeMetadata.BadgeIds, false, false)
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

	handledBadgeIds := []*UintRange{}
	if len(inheritedBalances) > 0 {
		for _, inheritedBalance := range inheritedBalances {
			err = ValidateRangesAreValid(inheritedBalance.BadgeIds, false, false)
			if err != nil {
				return sdkerrors.Wrapf(ErrInvalidRequest, "invalid badgeIds")
			}

			err = ValidateRangesAreValid(inheritedBalance.ParentBadgeIds, false, false)
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
			for _, badgeUintRange := range inheritedBalance.BadgeIds {
				totalBadgeLength = totalBadgeLength.Add(badgeUintRange.End.Sub(badgeUintRange.Start).AddUint64(1))
			}

			totalParentBadgeLength := sdk.NewUint(0)
			for _, badgeUintRange := range inheritedBalance.ParentBadgeIds {
				totalParentBadgeLength = totalParentBadgeLength.Add(badgeUintRange.End.Sub(badgeUintRange.Start).AddUint64(1))
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
