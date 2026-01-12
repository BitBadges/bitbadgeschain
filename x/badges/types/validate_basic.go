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

const (
	// MaxUint64Value represents the maximum value for uint64
	MaxUint64Value = math.MaxUint64
)

var (
	// URI must be a valid URI. Method <= 10 characters long. Path <= 90 characters long.
	reUriString = `\w+:(\/?\/?)[^\s]+`
	reUri       = regexp.MustCompile(fmt.Sprintf(`^%s$`, reUriString))

	// Cosmos wrapper path denom/symbol validation: only a-zA-Z, _, {, }, and -
	reCosmosWrapperPathString = `^[a-zA-Z_{}-]+$`
	reCosmosWrapperPath       = regexp.MustCompile(reCosmosWrapperPathString)
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
		return sdkerrors.Wrapf(ErrInvalidURI, "invalid URI: %s", uri)
	}

	return nil
}

// ValidateCosmosWrapperPathDenom validates that a cosmos wrapper path denom contains only allowed characters: a-zA-Z, _, {, }, and -
func ValidateCosmosWrapperPathDenom(denom string) error {
	if denom == "" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "denom cannot be empty")
	}

	regexMatch := reCosmosWrapperPath.MatchString(denom)
	if !regexMatch {
		return sdkerrors.Wrapf(ErrInvalidRequest, "denom contains invalid characters - only a-zA-Z, _, {, }, and - are allowed: %s", denom)
	}

	return nil
}

// ValidateCosmosWrapperPathSymbol validates that a cosmos wrapper path symbol contains only allowed characters: a-zA-Z, _, {, }, and -
func ValidateCosmosWrapperPathSymbol(symbol string) error {
	if symbol == "" {
		return nil // Symbol can be empty
	}

	regexMatch := reCosmosWrapperPath.MatchString(symbol)
	if !regexMatch {
		return sdkerrors.Wrapf(ErrInvalidRequest, "symbol contains invalid characters - only a-zA-Z, _, {, }, and - are allowed: %s", symbol)
	}

	return nil
}

func ValidateAddress(address string, alowMint bool) error {
	if alowMint && (address == "Mint") {
		return nil
	}

	// Validate address using global SDK config (should be "bb" prefix)
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid address: %s", err)
	}
	return nil
}

func DoRangesOverlap(ids []*UintRange) bool {
	// Create a copy to avoid modifying the input slice
	idsCopy := make([]*UintRange, len(ids))
	copy(idsCopy, ids)

	// Insertion sort in order of range.Start. If two have same range.Start, sort by range.End.
	n := len(idsCopy)
	for i := 1; i < n; i++ {
		j := i
		for j > 0 {
			if idsCopy[j-1].Start.GT(idsCopy[j].Start) {
				idsCopy[j-1], idsCopy[j] = idsCopy[j], idsCopy[j-1]
			} else if idsCopy[j-1].Start.Equal(idsCopy[j].Start) && idsCopy[j-1].End.GT(idsCopy[j].End) {
				idsCopy[j-1], idsCopy[j] = idsCopy[j], idsCopy[j-1]
			}
			j = j - 1
		}
	}

	// Check if any overlap
	for i := 1; i < n; i++ {
		prevInsertedRange := idsCopy[i-1]
		currRange := idsCopy[i]

		if currRange.Start.LTE(prevInsertedRange.End) {
			return true
		}
	}

	return false
}

// Validates ranges are valid. If end.IsZero(), we assume end == start.
func ValidateRangesAreValid(tokenUintRanges []*UintRange, allowAllUints bool, errorOnEmpty bool) error {
	if len(tokenUintRanges) == 0 {
		if errorOnEmpty {
			return sdkerrors.Wrapf(ErrInvalidUintRangeSpecified, "ID ranges cannot be empty (length == 0)")
		}
	}

	for _, tokenUintRange := range tokenUintRanges {
		if tokenUintRange == nil {
			return ErrRangesIsNil
		}

		if tokenUintRange.Start.IsNil() || tokenUintRange.End.IsNil() {
			return sdkerrors.Wrapf(ErrUintUnititialized, "ID range start and/or end is nil")
		}

		if tokenUintRange.Start.GT(tokenUintRange.End) {
			return ErrStartGreaterThanEnd
		}

		if !allowAllUints {
			if tokenUintRange.Start.IsZero() || tokenUintRange.End.IsZero() {
				return sdkerrors.Wrapf(ErrUintUnititialized, "ID range start and/or end is zero")
			}

			if tokenUintRange.Start.GT(sdkmath.NewUint(MaxUint64Value)) || tokenUintRange.End.GT(sdkmath.NewUint(MaxUint64Value)) {
				return ErrUintGreaterThanMax
			}
		}
	}

	overlap := DoRangesOverlap(tokenUintRanges)
	if overlap {
		return ErrRangesOverlap
	}

	return nil
}

// Validates no element is X
func ValidateNoElementIsX(amounts []sdkmath.Uint, x sdkmath.Uint) error {
	for _, amount := range amounts {
		if amount.Equal(x) {
			return sdkerrors.Wrapf(ErrElementCantEqualThis, "amount cannot equal %s", x.String())
		}
	}
	return nil
}

// Validates no element is X
func ValidateNoStringElementIsX(addresses []string, x string) error {
	for _, amount := range addresses {
		if amount == x {
			return sdkerrors.Wrapf(ErrElementCantEqualThis, "address cannot equal %s", x)
		}
	}
	return nil
}

// ValidateAltTimeChecks validates alt time checks for offline hours and days
func ValidateAltTimeChecks(altTimeChecks *AltTimeChecks) error {
	if altTimeChecks == nil {
		return nil
	}

	// Validate offline hours (0-23)
	if err := validateTimeRanges(altTimeChecks.OfflineHours, 0, 23, "offline hours"); err != nil {
		return err
	}

	// Validate offline days (0-6)
	if err := validateTimeRanges(altTimeChecks.OfflineDays, 0, 6, "offline days"); err != nil {
		return err
	}

	return nil
}

// validateTimeRanges validates time ranges for hours or days
// It ensures:
// 1. Ranges are valid numbers within min-max (inclusive, allowing zero)
// 2. No duplicate values across multiple range objects
// 3. No wrapping allowed (if you need to wrap, create two separate ranges)
func validateTimeRanges(ranges []*UintRange, min, max uint64, fieldName string) error {
	if len(ranges) == 0 {
		return nil
	}

	// Track all values to check for duplicates
	seenValues := make(map[uint64]bool)

	for i, r := range ranges {
		if r == nil {
			return sdkerrors.Wrapf(ErrInvalidRequest, "%s range at index %d is nil", fieldName, i)
		}

		if r.Start.IsNil() || r.End.IsNil() {
			return sdkerrors.Wrapf(ErrUintUnititialized, "%s range at index %d has nil start or end", fieldName, i)
		}

		start := r.Start.Uint64()
		end := r.End.Uint64()

		// Check that start and end are within valid range (0 to max, inclusive)
		if start > max {
			return sdkerrors.Wrapf(ErrInvalidRequest, "%s range at index %d has start value %d which exceeds maximum %d", fieldName, i, start, max)
		}

		if end > max {
			return sdkerrors.Wrapf(ErrInvalidRequest, "%s range at index %d has end value %d which exceeds maximum %d", fieldName, i, end, max)
		}

		// Check that start <= end (no wrapping allowed)
		if start > end {
			return sdkerrors.Wrapf(ErrInvalidRequest, "%s range at index %d has start %d greater than end %d (wrapping not allowed, use two separate ranges)", fieldName, i, start, end)
		}

		// Check for duplicates and collect all values in the range
		for val := start; val <= end; val++ {
			if seenValues[val] {
				return sdkerrors.Wrapf(ErrInvalidRequest, "%s range at index %d contains duplicate value %d that was already defined in another range", fieldName, i, val)
			}
			seenValues[val] = true
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
		return sdkerrors.Wrapf(ErrInvalidAddress, "list ID is uninitialized")
	}

	if err := ValidateAddress(addressList.ListId, false); err == nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "list ID cannot be a valid address")
	}

	if strings.Contains(addressList.ListId, ":") || strings.Contains(addressList.ListId, "!") {
		return sdkerrors.Wrapf(ErrInvalidAddress, "list ID cannot contain : or !")
	}

	if addressList.Uri != "" {
		if err := ValidateURI(addressList.Uri); err != nil {
			return err
		}
	}

	for _, address := range addressList.Addresses {
		// Check for empty addresses
		if address == "" {
			return sdkerrors.Wrapf(ErrInvalidAddress, "address list cannot contain empty addresses")
		}
		if err := ValidateAddress(address, false); err != nil {
			return err
		}
	}

	// check duplicate addresses
	if duplicateInStringArray(addressList.Addresses) {
		return ErrDuplicateAddresses
	}

	return nil
}

// ValidateAddressListInput validates an AddressListInput (same validation as AddressList but without createdBy check)
func ValidateAddressListInput(addressListInput *AddressListInput) error {
	if addressListInput.ListId == "" ||
		addressListInput.ListId == "Mint" ||
		addressListInput.ListId == "Manager" ||
		addressListInput.ListId == "AllWithoutMint" ||
		addressListInput.ListId == "None" {
		return sdkerrors.Wrapf(ErrInvalidAddress, "list ID is uninitialized")
	}

	if err := ValidateAddress(addressListInput.ListId, false); err == nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "list ID cannot be a valid address")
	}

	if strings.Contains(addressListInput.ListId, ":") || strings.Contains(addressListInput.ListId, "!") {
		return sdkerrors.Wrapf(ErrInvalidAddress, "list ID cannot contain : or !")
	}

	if addressListInput.Uri != "" {
		if err := ValidateURI(addressListInput.Uri); err != nil {
			return err
		}
	}

	for _, address := range addressListInput.Addresses {
		// Check for empty addresses
		if address == "" {
			return sdkerrors.Wrapf(ErrInvalidAddress, "address list cannot contain empty addresses")
		}
		if err := ValidateAddress(address, false); err != nil {
			return err
		}
	}

	// check duplicate addresses
	if duplicateInStringArray(addressListInput.Addresses) {
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

func MaxNumTransfersIsBasicallyNil(maxNumTransfers *MaxNumTransfers) bool {
	return maxNumTransfers == nil || ((maxNumTransfers.OverallMaxNumTransfers.IsNil() || maxNumTransfers.OverallMaxNumTransfers.IsZero()) &&
		(maxNumTransfers.PerToAddressMaxNumTransfers.IsNil() || maxNumTransfers.PerToAddressMaxNumTransfers.IsZero()) &&
		(maxNumTransfers.PerFromAddressMaxNumTransfers.IsNil() || maxNumTransfers.PerFromAddressMaxNumTransfers.IsZero()) &&
		(maxNumTransfers.PerInitiatedByAddressMaxNumTransfers.IsNil() || maxNumTransfers.PerInitiatedByAddressMaxNumTransfers.IsZero()) &&
		(maxNumTransfers.ResetTimeIntervals == nil || IsResetTimeIntervalBasicallyNil(maxNumTransfers.ResetTimeIntervals)))
}

func IsResetTimeIntervalBasicallyNil(resetTimeInterval *ResetTimeIntervals) bool {
	return resetTimeInterval == nil || (resetTimeInterval.StartTime.IsNil() || resetTimeInterval.StartTime.IsZero()) &&
		(resetTimeInterval.IntervalLength.IsNil() || resetTimeInterval.IntervalLength.IsZero())
}

func ApprovalAmountsIsBasicallyNil(approvalAmounts *ApprovalAmounts) bool {
	return approvalAmounts == nil || ((approvalAmounts.OverallApprovalAmount.IsNil() || approvalAmounts.OverallApprovalAmount.IsZero()) &&
		(approvalAmounts.PerToAddressApprovalAmount.IsNil() || approvalAmounts.PerToAddressApprovalAmount.IsZero()) &&
		(approvalAmounts.PerFromAddressApprovalAmount.IsNil() || approvalAmounts.PerFromAddressApprovalAmount.IsZero()) &&
		(approvalAmounts.PerInitiatedByAddressApprovalAmount.IsNil() || approvalAmounts.PerInitiatedByAddressApprovalAmount.IsZero()) &&
		(approvalAmounts.ResetTimeIntervals == nil || IsResetTimeIntervalBasicallyNil(approvalAmounts.ResetTimeIntervals)))
}

func CollectionApprovalIsAutoScannable(approvalCriteria *ApprovalCriteria) bool {
	if approvalCriteria == nil {
		return true
	}

	if approvalCriteria.MustPrioritize {
		return false
	}

	if approvalCriteria.CoinTransfers != nil && len(approvalCriteria.CoinTransfers) > 0 {
		return false
	}

	if approvalCriteria.PredeterminedBalances != nil && !PredeterminedBalancesIsBasicallyNil(approvalCriteria.PredeterminedBalances) {
		return false
	}

	// Theoretically, we might be able to remove this but two things:
	// 1. It could potentially change which IDs are received (but that only makes sense if predetermined balances is true)
	// 2. We need to pass stuff to MsgTransferTokens so this doesn't really make sense for auto-scanning
	if approvalCriteria.MerkleChallenges != nil && len(approvalCriteria.MerkleChallenges) > 0 {
		return false
	}

	// I guess ETH signatures also fall under same category
	if approvalCriteria.EthSignatureChallenges != nil && len(approvalCriteria.EthSignatureChallenges) > 0 {
		return false
	}

	// if approvalCriteria.MaxNumTransfers != nil && !MaxNumTransfersIsBasicallyNil(approvalCriteria.MaxNumTransfers) {
	// 	return false
	// }

	// if approvalCriteria.DynamicStoreChallenges != nil && len(approvalCriteria.DynamicStoreChallenges) > 0 {
	// 	return false
	// }

	// if approvalCriteria.ApprovalAmounts != nil && !ApprovalAmountsIsBasicallyNil(approvalCriteria.ApprovalAmounts) {
	// 	return false
	// }

	// Note: mustOwnTokens, etc are fine since they are read-only during MsgTransferTokens

	return true
}

func ValidateCollectionApprovals(ctx sdk.Context, collectionApprovals []*CollectionApproval, canChangeValues bool) error {
	for i := 0; i < len(collectionApprovals); i++ {
		if collectionApprovals[i].ApprovalId == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id is uninitialized at index %d", i)
		}

		if collectionApprovals[i].ApprovalId == "All" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id can not be All at index %d", i)
		}

		reservedApprovalIds := []string{"default-outgoing", "default-incoming", "self-initiated-outgoing", "self-initiated-incoming", "all-incoming-transfers"}

		for _, reservedApprovalId := range reservedApprovalIds {
			if collectionApprovals[i].ApprovalId == reservedApprovalId {
				return sdkerrors.Wrapf(ErrInvalidRequest, "approval id can not be %s at index %d", reservedApprovalId, i)
			}
		}

		for j := i + 1; j < len(collectionApprovals); j++ {
			if collectionApprovals[i].ApprovalId == collectionApprovals[j].ApprovalId {
				return sdkerrors.Wrapf(ErrInvalidRequest, "duplicate approval ids at indices %d and %d: %s", i, j, collectionApprovals[i].ApprovalId)
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

		if err := ValidateRangesAreValid(collectionApproval.TokenIds, false, false); err != nil {
			return sdkerrors.Wrapf(err, "invalid token IDs")
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

				if !coinTransfer.OverrideToWithInitiator && ValidateAddress(coinTransfer.To, false) != nil {
					return sdkerrors.Wrapf(ErrInvalidRequest, "invalid to address")
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
				}
			}

			// This is a sanity check to preventt accidental unlimited approvals from the current address
			// If they really do want a very, very large number, they can just set max transfers to a large number
			// Validate that if maxNumTransfers is unlimited, coinTransfers cannot have overrideFromWithApproverAddress
			if MaxNumTransfersIsBasicallyNil(approvalCriteria.MaxNumTransfers) {
				for _, coinTransfer := range approvalCriteria.CoinTransfers {
					if coinTransfer != nil && coinTransfer.OverrideFromWithApproverAddress {
						return sdkerrors.Wrapf(ErrInvalidRequest, "overrideFromWithApproverAddress cannot be used when maxNumTransfers is unlimited (nothing set)")
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

			if canChangeValues {
				if approvalCriteria.MustOwnTokens == nil {
					approvalCriteria.MustOwnTokens = []*MustOwnTokens{}
				}

				if approvalCriteria.DynamicStoreChallenges == nil {
					approvalCriteria.DynamicStoreChallenges = []*DynamicStoreChallenge{}
				}

				for _, mustOwnTokenBalance := range approvalCriteria.MustOwnTokens {
					if mustOwnTokenBalance == nil {
						return sdkerrors.Wrapf(ErrInvalidRequest, "mustOwnTokens balance is nil")
					}

					if err := ValidateRangesAreValid(mustOwnTokenBalance.TokenIds, false, false); err != nil {
						return sdkerrors.Wrapf(err, "invalid token IDs")
					}

					if err := ValidateRangesAreValid(mustOwnTokenBalance.OwnershipTimes, false, false); err != nil {
						return sdkerrors.Wrapf(err, "invalid owned times")
					}

					if err := ValidateRangesAreValid([]*UintRange{mustOwnTokenBalance.AmountRange}, true, true); err != nil {
						return sdkerrors.Wrapf(err, "invalid transfer times")
					}

					if mustOwnTokenBalance.CollectionId.IsNil() || mustOwnTokenBalance.CollectionId.IsZero() {
						return sdkerrors.Wrapf(ErrUintUnititialized, "collection id is uninitialized")
					}
				}

				// Validate dynamic store challenges
				if approvalCriteria.DynamicStoreChallenges == nil {
					approvalCriteria.DynamicStoreChallenges = []*DynamicStoreChallenge{}
				}

				// Check for duplicate store IDs
				storeIds := make(map[string]bool)
				for _, challenge := range approvalCriteria.DynamicStoreChallenges {
					if challenge == nil {
						return sdkerrors.Wrapf(ErrInvalidRequest, "dynamic store challenge is nil")
					}

					if challenge.StoreId.IsNil() {
						return sdkerrors.Wrapf(ErrUintUnititialized, "dynamic store challenge storeId is uninitialized")
					}

					if challenge.StoreId.IsZero() {
						return sdkerrors.Wrapf(ErrUintUnititialized, "dynamic store challenge storeId is zero")
					}

					storeIdStr := challenge.StoreId.String()
					if storeIds[storeIdStr] {
						return sdkerrors.Wrapf(ErrInvalidRequest, "duplicate dynamic store challenge storeId: %s", storeIdStr)
					}
					storeIds[storeIdStr] = true
				}

				// Validate voting challenges
				if approvalCriteria.VotingChallenges == nil {
					approvalCriteria.VotingChallenges = []*VotingChallenge{}
				}

				// Check for duplicate proposal IDs
				proposalIds := make(map[string]bool)
				for _, challenge := range approvalCriteria.VotingChallenges {
					if challenge == nil {
						return sdkerrors.Wrapf(ErrInvalidRequest, "voting challenge is nil")
					}

					if err := challenge.ValidateBasic(); err != nil {
						return sdkerrors.Wrapf(err, "invalid voting challenge")
					}

					if proposalIds[challenge.ProposalId] {
						return sdkerrors.Wrapf(ErrInvalidRequest, "duplicate voting challenge proposalId: %s", challenge.ProposalId)
					}
					proposalIds[challenge.ProposalId] = true
				}

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

				if approvalCriteria.ApprovalAmounts.ResetTimeIntervals == nil {
					approvalCriteria.ApprovalAmounts.ResetTimeIntervals = &ResetTimeIntervals{}
				}

				if approvalCriteria.ApprovalAmounts.ResetTimeIntervals.StartTime.IsNil() {
					approvalCriteria.ApprovalAmounts.ResetTimeIntervals.StartTime = sdkmath.NewUint(0)
				}

				if approvalCriteria.ApprovalAmounts.ResetTimeIntervals.IntervalLength.IsNil() {
					approvalCriteria.ApprovalAmounts.ResetTimeIntervals.IntervalLength = sdkmath.NewUint(0)
				}
			}

			// Validate ApprovalAmounts: if amountTrackerId is empty/nil but other fields are set, throw error
			if approvalCriteria.ApprovalAmounts != nil {
				hasNonNilFields := (!approvalCriteria.ApprovalAmounts.OverallApprovalAmount.IsNil() && !approvalCriteria.ApprovalAmounts.OverallApprovalAmount.IsZero()) ||
					(!approvalCriteria.ApprovalAmounts.PerToAddressApprovalAmount.IsNil() && !approvalCriteria.ApprovalAmounts.PerToAddressApprovalAmount.IsZero()) ||
					(!approvalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount.IsNil() && !approvalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount.IsZero()) ||
					(!approvalCriteria.ApprovalAmounts.PerInitiatedByAddressApprovalAmount.IsNil() && !approvalCriteria.ApprovalAmounts.PerInitiatedByAddressApprovalAmount.IsZero()) ||
					(approvalCriteria.ApprovalAmounts.ResetTimeIntervals != nil && !IsResetTimeIntervalBasicallyNil(approvalCriteria.ApprovalAmounts.ResetTimeIntervals))

				if hasNonNilFields && approvalCriteria.ApprovalAmounts.AmountTrackerId == "" {
					return sdkerrors.Wrapf(ErrInvalidRequest, "approvalAmounts has non-nil fields set but amountTrackerId is empty or nil")
				}
			}

			if canChangeValues {
				if approvalCriteria.AutoDeletionOptions == nil {
					approvalCriteria.AutoDeletionOptions = &AutoDeletionOptions{}
				}

				if approvalCriteria.UserRoyalties == nil {
					approvalCriteria.UserRoyalties = &UserRoyalties{}
				}

				if approvalCriteria.UserRoyalties.Percentage.IsNil() {
					approvalCriteria.UserRoyalties.Percentage = sdkmath.NewUint(0)
				}

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

				if approvalCriteria.MaxNumTransfers.ResetTimeIntervals == nil {
					approvalCriteria.MaxNumTransfers.ResetTimeIntervals = &ResetTimeIntervals{}
				}

				if approvalCriteria.MaxNumTransfers.ResetTimeIntervals.StartTime.IsNil() {
					approvalCriteria.MaxNumTransfers.ResetTimeIntervals.StartTime = sdkmath.NewUint(0)
				}

				if approvalCriteria.MaxNumTransfers.ResetTimeIntervals.IntervalLength.IsNil() {
					approvalCriteria.MaxNumTransfers.ResetTimeIntervals.IntervalLength = sdkmath.NewUint(0)
				}
			}

			// Validate MaxNumTransfers: if amountTrackerId is empty/nil but other fields are set, throw error
			if approvalCriteria.MaxNumTransfers != nil {
				hasNonNilFields := (!approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsNil() && !approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsZero()) ||
					(!approvalCriteria.MaxNumTransfers.PerToAddressMaxNumTransfers.IsNil() && !approvalCriteria.MaxNumTransfers.PerToAddressMaxNumTransfers.IsZero()) ||
					(!approvalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers.IsNil() && !approvalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers.IsZero()) ||
					(!approvalCriteria.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers.IsNil() && !approvalCriteria.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers.IsZero()) ||
					(approvalCriteria.MaxNumTransfers.ResetTimeIntervals != nil && !IsResetTimeIntervalBasicallyNil(approvalCriteria.MaxNumTransfers.ResetTimeIntervals))

				if hasNonNilFields && approvalCriteria.MaxNumTransfers.AmountTrackerId == "" {
					return sdkerrors.Wrapf(ErrInvalidRequest, "maxNumTransfers has non-nil fields set but amountTrackerId is empty or nil")
				}
			}

			if approvalCriteria.PredeterminedBalances != nil {
				isBasicallyNil := PredeterminedBalancesIsBasicallyNil(approvalCriteria.PredeterminedBalances)
				manualBalancesIsBasicallyNil := IsManualBalancesBasicallyNil(approvalCriteria.PredeterminedBalances.ManualBalances)
				sequentialTransferIsBasicallyNil := IsSequentialTransferBasicallyNil(approvalCriteria.PredeterminedBalances.IncrementedBalances)
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

					var err error
					if manualBalancesIsBasicallyNil && !sequentialTransferIsBasicallyNil {
						sequentialTransfer := approvalCriteria.PredeterminedBalances.IncrementedBalances
						sequentialTransfer.StartBalances, err = ValidateBalances(ctx, sequentialTransfer.StartBalances, canChangeValues)
						if err != nil {
							return err
						}

						if sequentialTransfer.IncrementTokenIdsBy.IsNil() {
							return sdkerrors.Wrapf(ErrUintUnititialized, "increment ids by is uninitialized")
						}

						if sequentialTransfer.IncrementOwnershipTimesBy.IsNil() {
							return sdkerrors.Wrapf(ErrUintUnititialized, "increment ownership times by is uninitialized")
						}

						if sequentialTransfer.DurationFromTimestamp.IsNil() {
							return sdkerrors.Wrapf(ErrUintUnititialized, "approval duration from now is uninitialized")
						}

						if sequentialTransfer.RecurringOwnershipTimes != nil {
							startTimeIsPositive := !sequentialTransfer.RecurringOwnershipTimes.StartTime.IsNil() && !sequentialTransfer.RecurringOwnershipTimes.StartTime.IsZero()

							if sequentialTransfer.RecurringOwnershipTimes.IntervalLength.IsNil() {
								return sdkerrors.Wrapf(ErrUintUnititialized, "interval length is uninitialized")
							}

							if sequentialTransfer.RecurringOwnershipTimes.StartTime.IsNil() {
								return sdkerrors.Wrapf(ErrUintUnititialized, "start time is uninitialized")
							}

							if sequentialTransfer.RecurringOwnershipTimes.ChargePeriodLength.IsNil() {
								return sdkerrors.Wrapf(ErrUintUnititialized, "grace period length is uninitialized")
							}

							if startTimeIsPositive {
								if sequentialTransfer.RecurringOwnershipTimes.IntervalLength.IsZero() {
									return sdkerrors.Wrapf(ErrInvalidRequest, "interval length cannot be zero if start time is positive")
								}

								if sequentialTransfer.RecurringOwnershipTimes.ChargePeriodLength.IsZero() {
									return sdkerrors.Wrapf(ErrInvalidRequest, "grace period length cannot be zero if start time is positive")
								}
							}

							// grace period cannot be longer than the interval length
							if sequentialTransfer.RecurringOwnershipTimes.ChargePeriodLength.GT(sequentialTransfer.RecurringOwnershipTimes.IntervalLength) {
								return sdkerrors.Wrapf(ErrInvalidRequest, "grace period length cannot be longer than or equal to the interval length")
							}
						}

						// Cant use both increment ownership times by and approval duration from now
						isApprovalDurationZero := sequentialTransfer.DurationFromTimestamp.IsZero() || sequentialTransfer.DurationFromTimestamp.IsNil()
						isIncrementOwnershipTimesByZero := sequentialTransfer.IncrementOwnershipTimesBy.IsZero() || sequentialTransfer.IncrementOwnershipTimesBy.IsNil()
						isRecurringOwnershipTimesZero := sequentialTransfer.RecurringOwnershipTimes == nil || ((sequentialTransfer.RecurringOwnershipTimes.IntervalLength.IsZero() || sequentialTransfer.RecurringOwnershipTimes.IntervalLength.IsNil()) ||
							(sequentialTransfer.RecurringOwnershipTimes.StartTime.IsZero() || sequentialTransfer.RecurringOwnershipTimes.StartTime.IsNil()) ||
							(sequentialTransfer.RecurringOwnershipTimes.ChargePeriodLength.IsZero() || sequentialTransfer.RecurringOwnershipTimes.ChargePeriodLength.IsNil()))

						count := 0
						if !isApprovalDurationZero {
							count++
						}
						if !isIncrementOwnershipTimesByZero {
							count++
						}
						if !isRecurringOwnershipTimesZero {
							count++
						}

						if count > 1 {
							return sdkerrors.Wrapf(ErrInvalidRequest, "only one of increment ownership times by, approval duration from now, or recurring ownership times can be set")
						}

						tokenIdOverrideCount := 0
						if !sequentialTransfer.IncrementTokenIdsBy.IsNil() && !sequentialTransfer.IncrementTokenIdsBy.IsZero() {
							tokenIdOverrideCount++
						}

						if sequentialTransfer.AllowOverrideWithAnyValidToken {
							tokenIdOverrideCount++
						}

						if tokenIdOverrideCount > 1 {
							return sdkerrors.Wrapf(ErrInvalidRequest, "only one of increment token ids by, or allow override with any valid ID can be set")
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

			// Validate AltTimeChecks
			if approvalCriteria.AltTimeChecks != nil {
				if err := ValidateAltTimeChecks(approvalCriteria.AltTimeChecks); err != nil {
					return sdkerrors.Wrapf(err, "invalid alt time checks")
				}
			}
		} else {
			approvalCriteria = &ApprovalCriteria{}
		}
	}

	return nil
}

// ValidateCollectionApprovalsWithInvariants validates collection approvals and checks invariants
func ValidateCollectionApprovalsWithInvariants(ctx sdk.Context, collectionApprovals []*CollectionApproval, canChangeValues bool, collection *TokenCollection) error {
	// First validate the basic collection approvals
	if err := ValidateCollectionApprovals(ctx, collectionApprovals, canChangeValues); err != nil {
		return err
	}

	// Check invariants if collection is provided
	if collection != nil && collection.Invariants != nil {
		if collection.Invariants.NoCustomOwnershipTimes {
			for _, collectionApproval := range collectionApprovals {
				if err := ValidateNoCustomOwnershipTimesInvariant(collectionApproval.OwnershipTimes, true); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func IsManualBalancesBasicallyNil(manualBalances []*ManualBalances) bool {
	return manualBalances == nil || len(manualBalances) == 0
}

func IsOrderCalculationMethodBasicallyNil(orderCalculationMethod *PredeterminedOrderCalculationMethod) bool {
	return orderCalculationMethod == nil || (!orderCalculationMethod.UseMerkleChallengeLeafIndex &&
		!orderCalculationMethod.UseOverallNumTransfers &&
		!orderCalculationMethod.UsePerToAddressNumTransfers &&
		!orderCalculationMethod.UsePerFromAddressNumTransfers &&
		!orderCalculationMethod.UsePerInitiatedByAddressNumTransfers)
}

func IsSequentialTransferBasicallyNil(incrementedBalances *IncrementedBalances) bool {
	return incrementedBalances == nil || ((incrementedBalances.StartBalances == nil || len(incrementedBalances.StartBalances) == 0) &&
		(!incrementedBalances.AllowOverrideWithAnyValidToken) &&
		(!incrementedBalances.AllowOverrideTimestamp) &&
		(incrementedBalances.IncrementTokenIdsBy.IsNil() ||
			incrementedBalances.IncrementTokenIdsBy.IsZero()) &&
		(incrementedBalances.IncrementOwnershipTimesBy.IsNil() ||
			incrementedBalances.IncrementOwnershipTimesBy.IsZero()) &&
		(incrementedBalances.DurationFromTimestamp.IsNil() ||
			incrementedBalances.DurationFromTimestamp.IsZero()) && (incrementedBalances.RecurringOwnershipTimes == nil ||
		(incrementedBalances.RecurringOwnershipTimes.StartTime.IsNil() ||
			incrementedBalances.RecurringOwnershipTimes.StartTime.IsZero()) &&
			(incrementedBalances.RecurringOwnershipTimes.IntervalLength.IsNil() ||
				incrementedBalances.RecurringOwnershipTimes.IntervalLength.IsZero())))
}

func PredeterminedBalancesIsBasicallyNil(predeterminedBalances *PredeterminedBalances) bool {
	orderCalculationMethodIsBasicallyNil := IsOrderCalculationMethodBasicallyNil(predeterminedBalances.OrderCalculationMethod)
	sequentialTransferIsBasicallyNil := IsSequentialTransferBasicallyNil(predeterminedBalances.IncrementedBalances)
	manualBalancesIsBasicallyNil := IsManualBalancesBasicallyNil(predeterminedBalances.ManualBalances)

	isBasicallyNil := orderCalculationMethodIsBasicallyNil && sequentialTransferIsBasicallyNil && manualBalancesIsBasicallyNil

	return isBasicallyNil
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

		// For non-whitelist trees, we can only use max one use per leaf (bc as soon as we use a leaf, the merkle path is public so anyone can use it)
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
	var err error
	for _, balance := range balances {
		if balance == nil {
			return balances, sdkerrors.Wrapf(ErrInvalidLengthBalances, "balances is nil")
		}

		if balance.Amount.IsNil() || balance.Amount.IsZero() {
			return balances, sdkerrors.Wrapf(ErrAmountEqualsZero, "amount is uninitialized")
		}

		err = ValidateRangesAreValid(balance.TokenIds, false, true)
		if err != nil {
			return balances, sdkerrors.Wrapf(err, "invalid balance token ids")
		}

		err = ValidateRangesAreValid(balance.OwnershipTimes, false, true)
		if err != nil {
			return balances, sdkerrors.Wrapf(err, "invalid balance times")
		}
	}

	balances, err = HandleDuplicateTokenIds(ctx, balances, canChangeValues)
	if err != nil {
		return balances, err
	}

	return balances, nil
}

// ValidateTransferWithInvariants validates a transfer and checks invariants
func ValidateTransferWithInvariants(ctx sdk.Context, transfer *Transfer, canChangeValues bool, collection *TokenCollection) error {
	// First validate the basic transfer
	if err := ValidateTransfer(ctx, transfer, canChangeValues); err != nil {
		return err
	}

	// Check invariants if collection is provided
	if collection != nil && collection.Invariants != nil {
		if collection.Invariants.NoCustomOwnershipTimes {
			for _, balance := range transfer.Balances {
				if err := ValidateNoCustomOwnershipTimesInvariant(balance.OwnershipTimes, true); err != nil {
					return err
				}
			}
		}

		// If cosmosCoinBackedPath is set, transfers from Mint address are not allowed
		if collection.Invariants.CosmosCoinBackedPath != nil {
			if IsMintAddress(transfer.From) {
				return sdkerrors.Wrapf(ErrInvalidRequest, "transfers from Mint address are not allowed when cosmosCoinBackedPath is set")
			}
		}
	}

	return nil
}

func ValidateTransfer(ctx sdk.Context, transfer *Transfer, canChangeValues bool) error {
	var err error

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

	if len(transfer.PrioritizedApprovals) > 0 {
		for _, prioritizedApproval := range transfer.PrioritizedApprovals {
			if prioritizedApproval.Version.IsNil() {
				return sdkerrors.Wrapf(ErrUintUnititialized, "version is uninitialized")
			}
		}
	}

	if transfer.PrecalculateBalancesFromApproval != nil {
		precalcDetails := transfer.PrecalculateBalancesFromApproval
		if precalcDetails.ApprovalLevel == "" && precalcDetails.ApproverAddress == "" && precalcDetails.ApprovalId == "" {
			// basically nil
		} else {
			if precalcDetails.ApprovalLevel != "collection" && precalcDetails.ApprovalLevel != "incoming" && precalcDetails.ApprovalLevel != "outgoing" {
				return sdkerrors.Wrapf(ErrInvalidRequest, "approval level must be collection, incoming, or outgoing")
			}

			if precalcDetails.ApproverAddress != "" {
				if err := ValidateAddress(precalcDetails.ApproverAddress, false); err != nil {
					return sdkerrors.Wrapf(ErrInvalidAddress, "invalid approval id address (%s)", err)
				}
			}

			if precalcDetails.Version.IsNil() {
				return sdkerrors.Wrapf(ErrUintUnititialized, "version is uninitialized")
			}
		}
	}

	if canChangeValues {
		if transfer.PrecalculateBalancesFromApproval != nil && transfer.PrecalculateBalancesFromApproval.PrecalculationOptions == nil {
			transfer.PrecalculateBalancesFromApproval.PrecalculationOptions = &PrecalculationOptions{
				OverrideTimestamp: sdkmath.NewUint(0),
				TokenIdsOverride:  nil,
			}
		}
	}

	return nil
}

func ValidateTokenMetadata(tokenMetadata []*TokenMetadata, canChangeValues bool) error {
	var err error

	handledTokenIds := []*UintRange{}
	if len(tokenMetadata) > 0 {
		for _, tokenMetadata := range tokenMetadata {
			// Validate well-formedness of the message entries
			if err := ValidateURI(tokenMetadata.Uri); err != nil {
				return err
			}

			err = ValidateRangesAreValid(tokenMetadata.TokenIds, false, false)
			if err != nil {
				return sdkerrors.Wrapf(ErrInvalidRequest, "invalid IDIds")
			}

			if err := AssertRangesDoNotOverlapAtAll(handledTokenIds, tokenMetadata.TokenIds); err != nil {
				return sdkerrors.Wrapf(err, "token metadata has duplicate token ids")
			}

			handledTokenIds = append(handledTokenIds, SortUintRangesAndMergeAdjacentAndIntersecting(tokenMetadata.TokenIds)...)

			if canChangeValues {
				tokenMetadata.TokenIds = SortUintRangesAndMergeAdjacentAndIntersecting(tokenMetadata.TokenIds)
			}
		}
	}

	return nil
}

// IsFullOwnershipTimesRange checks if the ownership times represent a full range from 1 to MaxUint64
func IsFullOwnershipTimesRange(ownershipTimes []*UintRange) bool {
	if len(ownershipTimes) != 1 {
		return false
	}

	range_ := ownershipTimes[0]
	return range_.Start.Equal(sdkmath.NewUint(1)) && range_.End.Equal(sdkmath.NewUint(MaxUint64Value))
}

// ValidateNoCustomOwnershipTimesInvariant validates that all ownership times are full ranges when the invariant is enabled
func ValidateNoCustomOwnershipTimesInvariant(ownershipTimes []*UintRange, invariantEnabled bool) error {
	if !invariantEnabled {
		return nil
	}

	if !IsFullOwnershipTimesRange(ownershipTimes) {
		return sdkerrors.Wrapf(ErrInvalidRequest, "noCustomOwnershipTimes invariant is enabled: ownership times must be full range [{ start: 1, end: 18446744073709551615 }]")
	}

	return nil
}

// ValidateBasic validates a VotingChallenge
func (vc *VotingChallenge) ValidateBasic() error {
	if vc.ProposalId == "" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "proposalId cannot be empty")
	}

	if vc.QuorumThreshold.GT(sdkmath.NewUint(100)) {
		return sdkerrors.Wrapf(ErrInvalidRequest, "quorumThreshold must be between 0 and 100, got %s", vc.QuorumThreshold.String())
	}

	if len(vc.Voters) == 0 {
		return sdkerrors.Wrapf(ErrInvalidRequest, "voters list cannot be empty")
	}

	// Check for duplicate voters
	voterAddresses := make(map[string]bool)
	for _, voter := range vc.Voters {
		if err := voter.ValidateBasic(); err != nil {
			return err
		}
		if voterAddresses[voter.Address] {
			return sdkerrors.Wrapf(ErrInvalidRequest, "duplicate voter address: %s", voter.Address)
		}
		voterAddresses[voter.Address] = true
	}

	if vc.Uri != "" {
		if err := ValidateURI(vc.Uri); err != nil {
			return err
		}
	}

	return nil
}

// ValidateBasic validates a Voter
func (v *Voter) ValidateBasic() error {
	if v.Address == "" {
		return sdkerrors.Wrapf(ErrInvalidAddress, "voter address cannot be empty")
	}

	if err := ValidateAddress(v.Address, false); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid voter address (%s)", err)
	}

	if v.Weight.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "voter weight cannot be zero")
	}

	return nil
}

// ValidateBasic validates a VoteProof
func (vp *VoteProof) ValidateBasic() error {
	if vp.ProposalId == "" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "proposalId cannot be empty")
	}

	if vp.Voter == "" {
		return sdkerrors.Wrapf(ErrInvalidAddress, "voter address cannot be empty")
	}

	if err := ValidateAddress(vp.Voter, false); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid voter address (%s)", err)
	}

	if vp.YesWeight.GT(sdkmath.NewUint(100)) {
		return sdkerrors.Wrapf(ErrInvalidRequest, "yesWeight must be between 0 and 100, got %s", vp.YesWeight.String())
	}

	return nil
}
