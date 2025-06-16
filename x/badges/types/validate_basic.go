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

func CollectionApprovalHasNoSideEffects(approvalCriteria *ApprovalCriteria) bool {
	if approvalCriteria == nil {
		return true
	}

	if approvalCriteria.CoinTransfers != nil && len(approvalCriteria.CoinTransfers) > 0 {
		return false
	}

	if approvalCriteria.PredeterminedBalances != nil && !PredeterminedBalancesIsBasicallyNil(approvalCriteria.PredeterminedBalances) {
		return false
	}

	if approvalCriteria.MerkleChallenges != nil && len(approvalCriteria.MerkleChallenges) > 0 {
		return false
	}

	if approvalCriteria.MaxNumTransfers != nil && !MaxNumTransfersIsBasicallyNil(approvalCriteria.MaxNumTransfers) {
		return false
	}

	if approvalCriteria.ApprovalAmounts != nil && !ApprovalAmountsIsBasicallyNil(approvalCriteria.ApprovalAmounts) {
		return false
	}

	return true
}

func ValidateCollectionApprovals(ctx sdk.Context, collectionApprovals []*CollectionApproval, canChangeValues bool) error {
	for i := 0; i < len(collectionApprovals); i++ {
		if collectionApprovals[i].ApprovalId == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id is uninitialized")
		}

		if collectionApprovals[i].ApprovalId == "All" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval id can not be All")
		}

		reservedApprovalIds := []string{"default-outgoing", "default-incoming", "self-initiated-outgoing", "self-initiated-incoming", "all-incoming-transfers"}

		for _, reservedApprovalId := range reservedApprovalIds {
			if collectionApprovals[i].ApprovalId == reservedApprovalId {
				return sdkerrors.Wrapf(ErrInvalidRequest, "approval id can not be %s", reservedApprovalId)
			}
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

					if coinToTransfer.Denom != "ubadge" {
						return sdkerrors.Wrapf(ErrInvalidRequest, "coin denom must be badge")
					}

					if coinToTransfer.Amount.GT(sdkmath.NewInt(100000000000)) {
						return sdkerrors.Wrapf(ErrInvalidRequest, "coin amount is too large - the max amount is 100000000000ubadge")
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

					if numTrue > 1 {
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
								return sdkerrors.Wrapf(ErrInvalidRequest, "grace period length cannot be longer than or equal tothe interval length")
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

						badgeIdOverrideCount := 0
						if !sequentialTransfer.IncrementBadgeIdsBy.IsNil() && !sequentialTransfer.IncrementBadgeIdsBy.IsZero() {
							badgeIdOverrideCount++
						}

						if sequentialTransfer.AllowOverrideWithAnyValidBadge {
							badgeIdOverrideCount++
						}

						if badgeIdOverrideCount > 1 {
							return sdkerrors.Wrapf(ErrInvalidRequest, "only one of increment badge ids by, or allow override with any valid badge can be set")
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

func IsManualBalancesBasicallyNil(manualBalances []*ManualBalances) bool {
	return manualBalances == nil || len(manualBalances) == 0
}

func IsOrderCalculationMethodBasicallyNil(orderCalculationMethod *PredeterminedOrderCalculationMethod) bool {
	return orderCalculationMethod == nil || (orderCalculationMethod.UseMerkleChallengeLeafIndex == false &&
		orderCalculationMethod.UseOverallNumTransfers == false &&
		orderCalculationMethod.UsePerToAddressNumTransfers == false &&
		orderCalculationMethod.UsePerFromAddressNumTransfers == false &&
		orderCalculationMethod.UsePerInitiatedByAddressNumTransfers == false)
}

func IsSequentialTransferBasicallyNil(incrementedBalances *IncrementedBalances) bool {
	return incrementedBalances == nil || ((incrementedBalances.StartBalances == nil || len(incrementedBalances.StartBalances) == 0) &&
		(incrementedBalances.AllowOverrideWithAnyValidBadge == false) &&
		(incrementedBalances.AllowOverrideTimestamp == false) &&
		(incrementedBalances.IncrementBadgeIdsBy.IsNil() ||
			incrementedBalances.IncrementBadgeIdsBy.IsZero()) &&
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

	if len(transfer.PrioritizedApprovals) > 0 {
		for _, prioritizedApproval := range transfer.PrioritizedApprovals {
			if prioritizedApproval.Version.IsNil() {
				return sdkerrors.Wrapf(ErrUintUnititialized, "version is uninitialized")
			}
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

			if transfer.PrecalculateBalancesFromApproval.Version.IsNil() {
				return sdkerrors.Wrapf(ErrUintUnititialized, "version is uninitialized")
			}
		}
	}

	if canChangeValues {
		if transfer.PrecalculationOptions == nil {
			transfer.PrecalculationOptions = &PrecalculationOptions{
				OverrideTimestamp: sdkmath.NewUint(0),
				BadgeIdsOverride:  nil,
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
