package types

import (
	"encoding/json"
	"sort"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func FilterZeroBalances(balances []*Balance) []*Balance {
	newBalances := []*Balance{}
	for _, balance := range balances {
		if balance.Amount.GT(sdkmath.NewUint(0)) {
			newBalances = append(newBalances, balance)
		}
	}

	return newBalances
}

func DoBalancesExceedThreshold(ctx sdk.Context, balances []*Balance, thresholdBalances []*Balance) error {
	thresholdCopy := DeepCopyBalances(thresholdBalances)
	_, err := SubtractBalances(ctx, balances, thresholdCopy) // underflow if exceeds
	return err
}

func AddBalancesAndAssertDoesntExceedThreshold(ctx sdk.Context, currTally []*Balance, toAdd []*Balance, threshold []*Balance) ([]*Balance, error) {
	var err error
	//If we transferAsMuchAsPossible, we need to increment the currTally by all that we can
	//We then need to return the updated toAdd

	currTally, err = AddBalances(ctx, toAdd, currTally)
	if err != nil {
		return []*Balance{}, err
	}

	//Check if we exceed the threshold; will underflow if we do exceed it
	err = DoBalancesExceedThreshold(ctx, currTally, threshold)
	if err != nil {
		return []*Balance{}, sdkerrors.Wrapf(err, "curr tally plus added amounts exceeds threshold")
	}

	return currTally, nil
}

func AreBalancesEqual(ctx sdk.Context, expected []*Balance, actual []*Balance, checkZeroBalances bool) bool {
	expected = DeepCopyBalances(expected)
	actual = DeepCopyBalances(actual)

	if !checkZeroBalances {
		expected = FilterZeroBalances(expected)
		actual = FilterZeroBalances(actual)
	}

	actual, err := SubtractBalances(ctx, expected, actual)
	if len(actual) != 0 || err != nil {
		return false
	}

	return true
}

func DeepCopyBalances(balances []*Balance) []*Balance {
	newBalances := []*Balance{}
	for _, approval := range balances {
		balanceToAdd := &Balance{
			Amount: approval.Amount,
		}
		for _, badgeId := range approval.BadgeIds {
			balanceToAdd.BadgeIds = append(balanceToAdd.BadgeIds, &UintRange{
				Start: badgeId.Start,
				End:   badgeId.End,
			})
		}

		for _, time := range approval.OwnershipTimes {
			balanceToAdd.OwnershipTimes = append(balanceToAdd.OwnershipTimes, &UintRange{
				Start: time.Start,
				End:   time.End,
			})
		}

		newBalances = append(newBalances, balanceToAdd)
	}

	return newBalances
}

// We handle the following cases:
// 1) {amount: 1, badgeIds: [1 to 10, 5 to 20]} -> {amount: 1, badgeIds: [1 to 4, 11 to 20]}, {amount: 2: badgeIds: [5 to 10]}
// 2) {amount: 1, badgeIds: [5 to 10]}, {amount: 2: badgeIds: [5 to 10]} -> {amount: 3: badgeIds: [5 to 10]}
func HandleDuplicateBadgeIds(ctx sdk.Context, balances []*Balance, canChangeValues bool) ([]*Balance, error) {
	if !canChangeValues {
		return balances, nil
	}

	newBalances := []*Balance{}
	var err error
	newBalances, err = AddBalances(ctx, balances, newBalances)
	if err != nil {
		return []*Balance{}, err
	}

	return newBalances, nil
}

// Updates the balance for a specific ids from what it currently is to newAmount. No add/subtract logic. Just set it.
func UpdateBalance(ctx sdk.Context, newBalance *Balance, balances []*Balance) ([]*Balance, error) {
	//Can maybe optimize this in the future by doing this all in one loop instead of deleting then setting
	// ranges = SortUintRangesAndMerge(ranges)
	var err error
	balances, err = DeleteBalances(ctx, newBalance.BadgeIds, newBalance.OwnershipTimes, balances)
	if err != nil {
		return balances, err
	}

	bals := []*Balance{newBalance}
	balances, err = SetBalances(bals, balances)
	if err != nil {
		return balances, err
	}

	return balances, nil
}

// Gets the balances for specified ID ranges. Returns a new []*Balance where only the specified ID ranges and their balances are included. Appends balance == 0 objects so all IDs are accounted for, even if not found.
func GetBalancesForIds(ctx sdk.Context, idRanges []*UintRange, times []*UintRange, balances []*Balance) (newBalances []*Balance, err error) {
	fetchedBalances := []*Balance{}

	currPermissionDetails := []*UniversalPermissionDetails{}
	for _, balanceObj := range balances {
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.OwnershipTimes {
				currPermissionDetails = append(currPermissionDetails, &UniversalPermissionDetails{
					BadgeId:        currRange,
					OwnershipTime:  currTime,
					ArbitraryValue: balanceObj.Amount,
				})
			}
		}
	}

	toFetchPermissionDetails := []*UniversalPermissionDetails{}
	for _, rangeToFetch := range idRanges {
		for _, timeToFetch := range times {
			toFetchPermissionDetails = append(toFetchPermissionDetails, &UniversalPermissionDetails{
				BadgeId:       rangeToFetch,
				OwnershipTime: timeToFetch,
			},
			)
		}
	}

	overlaps, _, inNewButNotOld := GetOverlapsAndNonOverlaps(ctx, currPermissionDetails, toFetchPermissionDetails)
	//For all overlaps, we simply return the amount
	for _, overlapObject := range overlaps {
		overlap := overlapObject.Overlap
		amount := overlapObject.FirstDetails.ArbitraryValue.(sdkmath.Uint)

		fetchedBalances = append(fetchedBalances, &Balance{
			Amount:         amount,
			BadgeIds:       []*UintRange{overlap.BadgeId},
			OwnershipTimes: []*UintRange{overlap.OwnershipTime},
		})
	}

	//For those that were in toFetch but not currBalances, we return amount == 0
	for _, detail := range inNewButNotOld {
		fetchedBalances = append(fetchedBalances, &Balance{
			Amount:         sdkmath.NewUint(0),
			BadgeIds:       []*UintRange{detail.BadgeId},
			OwnershipTimes: []*UintRange{detail.OwnershipTime},
		})
	}

	return fetchedBalances, nil
}

func GetOverlappingBalances(ctx sdk.Context, transferBalancesToCheck []*Balance, maxPossible []*Balance) ([]*Balance, error) {
	newBalances := []*Balance{}
	for _, balance := range transferBalancesToCheck {
		prevAmount := balance.Amount

		fetchedBalances, err := GetBalancesForIds(ctx, balance.BadgeIds, balance.OwnershipTimes, maxPossible)
		if err != nil {
			return nil, err
		}

		for _, fetchedBalance := range fetchedBalances {
			balanceToAdd := &Balance{
				Amount:         fetchedBalance.Amount,
				BadgeIds:       fetchedBalance.BadgeIds,
				OwnershipTimes: fetchedBalance.OwnershipTimes,
			}

			//Take min amount
			if balanceToAdd.Amount.GT(prevAmount) {
				balanceToAdd.Amount = prevAmount
			}

			newBalances = append(newBalances, balanceToAdd)
		}
	}
	return newBalances, nil
}

func AddBalances(ctx sdk.Context, balancesToAdd []*Balance, balances []*Balance) ([]*Balance, error) {
	var err error
	for _, balance := range balancesToAdd {
		balances, err = AddBalance(ctx, balances, balance)
		if err != nil {
			return balances, err
		}
	}

	return balances, nil
}

// Adds a balance to all ids specified in []ranges
func AddBalance(ctx sdk.Context, existingBalances []*Balance, balanceToAdd *Balance) ([]*Balance, error) {
	currBalances, err := GetBalancesForIds(ctx, balanceToAdd.BadgeIds, balanceToAdd.OwnershipTimes, existingBalances)
	if err != nil {
		return existingBalances, err
	}

	existingBalances, err = DeleteBalances(ctx, balanceToAdd.BadgeIds, balanceToAdd.OwnershipTimes, existingBalances)
	if err != nil {
		return existingBalances, err
	}

	for _, balance := range currBalances {
		balance.Amount, err = SafeAdd(balance.Amount, balanceToAdd.Amount)
		if err != nil {
			return existingBalances, err
		}
	}

	existingBalances, err = SetBalances(currBalances, existingBalances)
	if err != nil {
		return existingBalances, err
	}

	return existingBalances, nil
}

func SubtractBalances(ctx sdk.Context, balancesToSubtract []*Balance, balances []*Balance) ([]*Balance, error) {
	var err error

	for _, balance := range balancesToSubtract {
		balances, err = SubtractBalance(ctx, balances, balance, false)
		if err != nil {
			return balances, err
		}
	}

	return balances, nil
}

func SubtractBalancesWithZeroForUnderflows(ctx sdk.Context, balancesToSubtract []*Balance, balances []*Balance) ([]*Balance, error) {
	var err error

	for _, balance := range balancesToSubtract {
		balances, err = SubtractBalance(ctx, balances, balance, true)
		if err != nil {
			return balances, err
		}
	}

	return balances, nil
}

// Subtracts a balance to all ids specified in []ranges
func SubtractBalance(ctx sdk.Context, balances []*Balance, balanceToRemove *Balance, setToZeroOnUnderflow bool) ([]*Balance, error) {
	currBalances, err := GetBalancesForIds(ctx, balanceToRemove.BadgeIds, balanceToRemove.OwnershipTimes, balances)
	if err != nil {
		return balances, err
	}

	balances, err = DeleteBalances(ctx, balanceToRemove.BadgeIds, balanceToRemove.OwnershipTimes, balances)
	if err != nil {
		return balances, err
	}

	for _, currBalanceObj := range currBalances {
		currBalanceObj.Amount, err = SafeSubtract(currBalanceObj.Amount, balanceToRemove.Amount)
		if err != nil {
			if setToZeroOnUnderflow {
				currBalanceObj.Amount = sdkmath.NewUint(0)
			} else {
				return balances, err
			}
		}
	}

	balances, err = SetBalances(currBalances, balances)
	if err != nil {
		return balances, err
	}

	return balances, nil
}

// Deletes the balance for a specific id.
func DeleteBalances(ctx sdk.Context, rangesToDelete []*UintRange, timesToDelete []*UintRange, balances []*Balance) ([]*Balance, error) {
	newBalances := []*Balance{}

	for _, balanceObj := range balances {
		currPermissionDetails := []*UniversalPermissionDetails{}
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.OwnershipTimes {
				currPermissionDetails = append(currPermissionDetails, &UniversalPermissionDetails{
					BadgeId:       currRange,
					OwnershipTime: currTime,
				})
			}
		}

		toDeletePermissionDetails := []*UniversalPermissionDetails{}
		for _, rangeToDelete := range rangesToDelete {
			for _, timeToDelete := range timesToDelete {
				toDeletePermissionDetails = append(toDeletePermissionDetails, &UniversalPermissionDetails{
					BadgeId:       rangeToDelete,
					OwnershipTime: timeToDelete,
				})
			}
		}

		_, inOldButNotNew, _ := GetOverlapsAndNonOverlaps(ctx, currPermissionDetails, toDeletePermissionDetails)
		for _, remainingBalance := range inOldButNotNew {
			newBalances = append(newBalances, &Balance{
				Amount:         balanceObj.Amount,
				BadgeIds:       []*UintRange{remainingBalance.BadgeId},
				OwnershipTimes: []*UintRange{remainingBalance.OwnershipTime},
			})
		}
	}

	return newBalances, nil
}

// Sets the balance for a specific id.
// Important precondition: assumes balance does not exist.
func SetBalances(newBalancesToSet []*Balance, balances []*Balance) ([]*Balance, error) {
	newBalancesWithoutZeroes := []*Balance{}
	for _, balance := range newBalancesToSet {
		if balance.Amount.GT(sdkmath.NewUint(0)) {
			newBalancesWithoutZeroes = append(newBalancesWithoutZeroes, balance)
		}
	}
	if len(newBalancesWithoutZeroes) == 0 {
		return balances, nil
	}

	balances = append(balances, newBalancesWithoutZeroes...)

	var err error

	//Little clean up to start. We sort and if we have adjacent (note not intersecting), we  merge them
	for _, balance := range balances {
		balance.OwnershipTimes, err = SortUintRangesAndMerge(balance.OwnershipTimes, false)
		if err != nil {
			return balances, nil
		}

		balance.BadgeIds, err = SortUintRangesAndMerge(balance.BadgeIds, false)
		if err != nil {
			return balances, nil
		}
	}

	//Sort by amount in increasing order
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].Amount.LT(balances[j].Amount)
	})

	//See if we can merge on a cross-balance level
	newBalances := []*Balance{}
	for i := 0; i < len(balances); i++ {
		currBalance := balances[i]

		merged := false
		for _, existingBalance := range newBalances {
			if currBalance.Amount.Equal(existingBalance.Amount) {
				if compareSlices(currBalance.BadgeIds, existingBalance.BadgeIds) {
					existingBalance.OwnershipTimes = append(existingBalance.OwnershipTimes, currBalance.OwnershipTimes...)
					existingBalance.OwnershipTimes = SortUintRangesAndMergeAdjacentAndIntersecting(existingBalance.OwnershipTimes)
					merged = true
					break
				} else if compareSlices(currBalance.OwnershipTimes, existingBalance.OwnershipTimes) {
					existingBalance.BadgeIds = append(existingBalance.BadgeIds, currBalance.BadgeIds...)
					existingBalance.BadgeIds = SortUintRangesAndMergeAdjacentAndIntersecting(existingBalance.BadgeIds)
					merged = true
					break
				}
			} else if currBalance.Amount.GT(existingBalance.Amount) {
				//We can't merge if the current balance has a greater amount (arr is sorted by amount)
				break
			}
		}

		if !merged {
			newBalances = append(newBalances, currBalance)
		}
	}

	return newBalances, err
}

func compareSlices(slice1, slice2 []*UintRange) bool {
	// Compare two slices for equality
	return jsonEqual(slice1, slice2)
}

func jsonEqual(a, b interface{}) bool {
	// Compare JSON representations of two values
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}

	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}

	return string(aJSON) == string(bJSON)
}

func IncrementBalances(
	ctx sdk.Context,
	startBalances []*Balance,
	numIncrements sdkmath.Uint,
	incrementOwnershipTimesBy sdkmath.Uint,
	incrementBadgeIdsBy sdkmath.Uint,
	durationFromTimestamp sdkmath.Uint,
	recurringOwnershipTimes *RecurringOwnershipTimes,
	overrideTimestamp sdkmath.Uint,
	allowOverrideTimestamp bool,
	overrideBadgeIds []*UintRange,
	allowOverrideBadgeIdsWithAnyValidBadgeId bool,
	collection *BadgeCollection,
) ([]*Balance, error) {
	balances := DeepCopyBalances(startBalances)
	now := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))

	for _, startBalance := range balances {

		if recurringOwnershipTimes != nil && !(recurringOwnershipTimes.IntervalLength.IsNil() || recurringOwnershipTimes.IntervalLength.IsZero()) {
			startTime := recurringOwnershipTimes.StartTime
			intervalLength := recurringOwnershipTimes.IntervalLength
			chargePeriodLength := recurringOwnershipTimes.ChargePeriodLength

			// Pre: Assert outside charge periods
			// Edge case: Handle first charge period (throws with negative interval number if we don't handle this)
			if now.LT(startTime) {
				if chargePeriodLength.GT(startTime) {
					return balances, sdkerrors.Wrapf(ErrOutsideChargePeriod, "outside charge period")
				}

				// Check within first charge period
				firstChargeTime := startTime.Sub(chargePeriodLength)
				if now.GTE(firstChargeTime) && now.LT(startTime) {
					startBalance.OwnershipTimes = []*UintRange{
						{
							Start: startTime,
							End:   startTime.Add(intervalLength).Sub(sdkmath.OneUint()),
						},
					}

					return balances, nil
				}

				return balances, sdkerrors.Wrapf(ErrOutsideChargePeriod, "outside charge period")
			}

			//1. Calculate what interval we are in
			interval := now.Sub(startTime).Quo(intervalLength)
			nextInterval := interval.Add(sdkmath.OneUint())

			//2. Calculate the new intervals
			newStartTime := startTime.Add(intervalLength.Mul(nextInterval))
			newEndTime := newStartTime.Add(intervalLength).Sub(sdkmath.OneUint())
			chargeAfterTime := newStartTime.Sub(chargePeriodLength)

			//3. Assert that we are in the charge period
			if now.GTE(chargeAfterTime) && now.LT(newStartTime) {
				startBalance.OwnershipTimes = []*UintRange{
					{
						Start: newStartTime,
						End:   newEndTime,
					},
				}
			} else {
				return balances, sdkerrors.Wrapf(ErrOutsideChargePeriod, "outside charge period")
			}

		} else if durationFromTimestamp.IsZero() || durationFromTimestamp.IsNil() {
			for _, time := range startBalance.OwnershipTimes {
				time.Start = time.Start.Add(numIncrements.Mul(incrementOwnershipTimesBy))
				time.End = time.End.Add(numIncrements.Mul(incrementOwnershipTimesBy))
			}
		} else {
			if allowOverrideTimestamp && !overrideTimestamp.IsNil() && overrideTimestamp.GT(sdkmath.ZeroUint()) {
				now = overrideTimestamp
			}

			startBalance.OwnershipTimes = []*UintRange{
				{
					Start: now,
					End:   now.Add(durationFromTimestamp).Sub(sdkmath.OneUint()),
				},
			}
		}

		//Handle token IDs override
		if allowOverrideBadgeIdsWithAnyValidBadgeId {
			//Verify that the token IDs are valid

			//1. Check size == 1
			if len(overrideBadgeIds) != 1 {
				return balances, sdkerrors.Wrapf(ErrInvalidBadgeIds, "invalid token IDs override (length != 1)")
			}

			//2. Check that the token IDs are the same
			if !overrideBadgeIds[0].Start.Equal(overrideBadgeIds[0].End) {
				return balances, sdkerrors.Wrapf(ErrInvalidBadgeIds, "invalid token IDs override (start != end)")
			}

			//3. Check that the token IDs are specified as valid in the collection
			isValid, err := SearchUintRangesForUint(overrideBadgeIds[0].Start, collection.ValidBadgeIds)
			if err != nil || !isValid {
				return balances, sdkerrors.Wrapf(ErrInvalidBadgeIds, "invalid token IDs override (not valid ID in collection)")
			}

			startBalance.BadgeIds = overrideBadgeIds
		} else {
			for _, badgeId := range startBalance.BadgeIds {
				badgeId.Start = badgeId.Start.Add(numIncrements.Mul(incrementBadgeIdsBy))
				badgeId.End = badgeId.End.Add(numIncrements.Mul(incrementBadgeIdsBy))
			}
		}
	}

	return balances, nil
}
