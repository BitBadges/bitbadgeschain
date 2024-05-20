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
	if uri == "" {
		return nil
	}

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

func ValidateAddressList(addressList *AddressList) error {
	if addressList.ListId == "" ||
		addressList.ListId == "Mint" ||
		addressList.ListId == "Manager" ||
		addressList.ListId == "AllWithoutMint" ||
		addressList.ListId == "None" {
		return sdkerrors.Wrapf(ErrInvalidAddress, "list id is uninitialized")
	}

	if err := ValidateAddress(addressList.ListId, false); err == nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "list id can not be a valid address")
	}

	if strings.Contains(addressList.ListId, ":") || strings.Contains(addressList.ListId, "!") {
		return sdkerrors.Wrapf(ErrInvalidAddress, "list id can not contain : or !")
	}

	if addressList.Uri != "" {
		if err := ValidateURI(addressList.Uri); err != nil {
			return err
		}
	}

	for _, address := range addressList.Addresses {
		if err := ValidateAddress(address, false); err != nil {
			return err
		}
	}

	//check duplicate addresses
	if duplicateInStringArray(addressList.Addresses) {
		return ErrDuplicateAddresses
	}

	return nil
}

func ValidateUserOutgoingApprovals(ctx sdk.Context, userOutgoingApprovals []*UserOutgoingApproval, fromAddress string, canChangeValues bool) error {
	castedTransfers := CastOutgoingTransfersToCollectionTransfers(userOutgoingApprovals, fromAddress)
	err := ValidateCollectionApprovals(ctx, castedTransfers, canChangeValues)
	return err
}

func ValidateUserIncomingApprovals(ctx sdk.Context, userIncomingApprovals []*UserIncomingApproval, toAddress string, canChangeValues bool) error {
	castedTransfers := CastIncomingTransfersToCollectionTransfers(userIncomingApprovals, toAddress)
	err := ValidateCollectionApprovals(ctx, castedTransfers, canChangeValues)
	return err
}

func ValidateCollectionApprovals(ctx sdk.Context, collectionApprovals []*CollectionApproval, canChangeValues bool) error {
	for i := 0; i < len(collectionApprovals); i++ {
		if collectionApprovals[i].ApprovalId == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id is uninitialized")
		}

		if collectionApprovals[i].ApprovalId == "All" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id can not be All")
		}

		if collectionApprovals[i].ApprovalId == "default-outgoing" || collectionApprovals[i].ApprovalId == "default-incoming" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id can not be default-outgoing or default-incoming")
		}

		if collectionApprovals[i].ApprovalId == "self-initiated-outgoing" || collectionApprovals[i].ApprovalId == "self-initiated-incoming" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id can not be default-outgoing or default-incoming")
		}

		for j := i + 1; j < len(collectionApprovals); j++ {
			if collectionApprovals[i].ApprovalId == collectionApprovals[j].ApprovalId {
				return sdkerrors.Wrapf(ErrInvalidRequest, "duplicate approval ids")
			}
		}
	}

	for _, collectionApproval := range collectionApprovals {
		if collectionApproval == nil {
			return sdkerrors.Wrapf(ErrInvalidRequest, "collection approved transfer is nil")
		}

		if collectionApproval.FromListId == "" {
			return sdkerrors.Wrapf(ErrInvalidAddress, "from list id is uninitialized")
		}

		if collectionApproval.ToListId == "" {
			return sdkerrors.Wrapf(ErrInvalidAddress, "to list id is uninitialized")
		}

		if collectionApproval.InitiatedByListId == "" {
			return sdkerrors.Wrapf(ErrInvalidAddress, "initiated by list id is uninitialized")
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
			for _, coinTransfer := range approvalCriteria.CoinTransfers {
				if coinTransfer == nil {
					return sdkerrors.Wrapf(ErrInvalidRequest, "coin transfer is nil")
				}

				if ValidateAddress(coinTransfer.To, false) != nil {
					return sdkerrors.Wrapf(ErrInvalidRequest, "invalid from address")
				}

				for _, coinToTransfer := range coinTransfer.Coins {
					if coinToTransfer == nil {
						return sdkerrors.Wrapf(ErrInvalidRequest, "coin to transfer is nil")
					}

					if coinToTransfer.Amount.IsNil() {
						return sdkerrors.Wrapf(ErrInvalidRequest, "coin amount is uninitialized")
					}

					if coinToTransfer.Amount.IsZero() {
						return sdkerrors.Wrapf(ErrInvalidRequest, "coin amount is zero")
					}

					if coinToTransfer.Denom == "" {
						return sdkerrors.Wrapf(ErrInvalidRequest, "coin denom is uninitialized")
					}

					if coinToTransfer.Denom != "badge" {
						return sdkerrors.Wrapf(ErrInvalidRequest, "coin denom must be badge")
					}
				}
			}

			usingLeafIndexForTransferOrder := false
			challengeTrackerIdForTransferOrder := ""
			if approvalCriteria.PredeterminedBalances != nil && approvalCriteria.PredeterminedBalances.OrderCalculationMethod != nil && approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {
				usingLeafIndexForTransferOrder = true
				challengeTrackerIdForTransferOrder = approvalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId
			}

			if err := ValidateMerkleChallenges(approvalCriteria.MerkleChallenges, usingLeafIndexForTransferOrder, challengeTrackerIdForTransferOrder); err != nil {
				return sdkerrors.Wrapf(err, "invalid challenges")
			}

			if approvalCriteria.MustOwnBadges == nil {
				approvalCriteria.MustOwnBadges = []*MustOwnBadges{}
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

			if canChangeValues {
				if approvalCriteria.ApprovalAmounts == nil {
					approvalCriteria.ApprovalAmounts = &ApprovalAmounts{}
				}

				if approvalCriteria.ApprovalAmounts.OverallApprovalAmount.IsNil() {
					approvalCriteria.ApprovalAmounts.OverallApprovalAmount = sdkmath.NewUint(0)
				}

				if approvalCriteria.ApprovalAmounts.PerToAddressApprovalAmount.IsNil() {
					approvalCriteria.ApprovalAmounts.PerToAddressApprovalAmount = sdkmath.NewUint(0)
				}

				if approvalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount.IsNil() {
					approvalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(0)
				}

				if approvalCriteria.ApprovalAmounts.PerInitiatedByAddressApprovalAmount.IsNil() {
					approvalCriteria.ApprovalAmounts.PerInitiatedByAddressApprovalAmount = sdkmath.NewUint(0)
				}
			}

			if canChangeValues {
				if approvalCriteria.MaxNumTransfers == nil {
					approvalCriteria.MaxNumTransfers = &MaxNumTransfers{}
				}

				if approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsNil() {
					approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers = sdkmath.NewUint(0)
				}

				if approvalCriteria.MaxNumTransfers.PerToAddressMaxNumTransfers.IsNil() {
					approvalCriteria.MaxNumTransfers.PerToAddressMaxNumTransfers = sdkmath.NewUint(0)
				}

				if approvalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers.IsNil() {
					approvalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(0)
				}

				if approvalCriteria.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers.IsNil() {
					approvalCriteria.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers = sdkmath.NewUint(0)
				}
			}

			if approvalCriteria.PredeterminedBalances != nil {
				orderCalculationMethodIsBasicallyNil := !approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex &&
					!approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers &&
					!approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers &&
					!approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers &&
					!approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers

				sequentialTransferIsBasicallyNil := approvalCriteria.PredeterminedBalances.IncrementedBalances == nil || ((approvalCriteria.PredeterminedBalances.IncrementedBalances.StartBalances == nil || len(approvalCriteria.PredeterminedBalances.IncrementedBalances.StartBalances) == 0) &&
					(approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementBadgeIdsBy.IsNil() ||
						approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementBadgeIdsBy.IsZero()) &&
					(approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementOwnershipTimesBy.IsNil() ||
						approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementOwnershipTimesBy.IsZero()))

				manualBalancesIsBasicallyNil := approvalCriteria.PredeterminedBalances.ManualBalances == nil || len(approvalCriteria.PredeterminedBalances.ManualBalances) == 0

				isBasicallyNil := orderCalculationMethodIsBasicallyNil && sequentialTransferIsBasicallyNil && manualBalancesIsBasicallyNil

				if !isBasicallyNil {
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
					if manualBalancesIsBasicallyNil && !sequentialTransferIsBasicallyNil {
						sequentialTransfer := approvalCriteria.PredeterminedBalances.IncrementedBalances
						sequentialTransfer.StartBalances, err = ValidateBalances(ctx, sequentialTransfer.StartBalances, canChangeValues)
						if err != nil {
							return err
						}

						if sequentialTransfer.IncrementBadgeIdsBy.IsNil() {
							return sdkerrors.Wrapf(ErrUintUnititialized, "increment ids by is uninitialized")
						}

						if sequentialTransfer.IncrementOwnershipTimesBy.IsNil() {
							return sdkerrors.Wrapf(ErrUintUnititialized, "max num transfers is uninitialized")
						}
					} else if !manualBalancesIsBasicallyNil && sequentialTransferIsBasicallyNil {
						for _, manualTransfer := range approvalCriteria.PredeterminedBalances.ManualBalances {
							manualTransfer.Balances, err = ValidateBalances(ctx, manualTransfer.Balances, canChangeValues)
							if err != nil {
								return err
							}
						}
					} else {
						return sdkerrors.Wrapf(ErrInvalidRequest, "manual transfers and sequential transfers cannot be both nil or both defined")
					}
				}

			} else {
				approvalCriteria.PredeterminedBalances = nil
			}
		} else {
			approvalCriteria = &ApprovalCriteria{}
		}
	}

	return nil
}

func ValidateMerkleChallenges(challenges []*MerkleChallenge, usingLeafIndexForTransferOrder bool, challengeTrackerIdForTransferOrder string) error {

	for i, challenge := range challenges {
		if challenge == nil || challenge.Root == "" {
			challenge = &MerkleChallenge{}
			return nil
		}

		if challenge.ExpectedProofLength.IsNil() {
			return sdkerrors.Wrapf(ErrUintUnititialized, "expected proof length is uninitialized")
		}

		if challenge.MaxUsesPerLeaf.IsNil() {
			return sdkerrors.Wrapf(ErrUintUnititialized, "max uses per leaf is uninitialized")
		}

		maxOneUsePerLeaf := challenge.MaxUsesPerLeaf.Equal(sdkmath.NewUint(1))

		if !maxOneUsePerLeaf && usingLeafIndexForTransferOrder && challenge.ChallengeTrackerId == challengeTrackerIdForTransferOrder {
			return ErrPrimaryChallengeMustBeOneUsePerLeaf
		}

		//For non-whitelist trees, we can only use max one use per leaf (bc as soon as we use a leaf, the merkle path is public so anyone can use it)
		if !maxOneUsePerLeaf && !challenge.UseCreatorAddressAsLeaf {
			return ErrCanOnlyUseMaxOneUsePerLeafWithWhitelistTree
		}

		for j := i + 1; j < len(challenges); j++ {
			if challenge.ChallengeTrackerId == challenges[j].ChallengeTrackerId {
				return sdkerrors.Wrapf(ErrInvalidRequest, "duplicate challenge ids")
			}
		}
	}

	return nil
}

func ValidateBalances(ctx sdk.Context, balances []*Balance, canChangeValues bool) ([]*Balance, error) {
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

	balances, err = HandleDuplicateBadgeIds(ctx, balances, canChangeValues)
	if err != nil {
		return balances, err
	}

	return balances, nil
}

func ValidateTransfer(ctx sdk.Context, transfer *Transfer, canChangeValues bool) error {
	err := *new(error)

	transfer.Balances, err = ValidateBalances(ctx, transfer.Balances, canChangeValues)
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

func ValidateBadgeMetadata(badgeMetadata []*BadgeMetadata, canChangeValues bool) error {
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

			if err := AssertRangesDoNotOverlapAtAll(handledBadgeIds, badgeMetadata.BadgeIds); err != nil {
				return sdkerrors.Wrapf(err, "badge metadata has duplicate badge ids")
			}

			handledBadgeIds = append(handledBadgeIds, SortUintRangesAndMergeAdjacentAndIntersecting(badgeMetadata.BadgeIds)...)

			if canChangeValues {
				badgeMetadata.BadgeIds = SortUintRangesAndMergeAdjacentAndIntersecting(badgeMetadata.BadgeIds)
			}
		}
	}

	return nil
}
