package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"

	"math"
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

func DoBalancesExceedThreshold(balances []*Balance, thresholdBalances []*Balance) bool {
	//Check if we exceed the threshold; will underflow if we do exceed it
	thresholdCopy := DeepCopyBalances(thresholdBalances)
	_, err := SubtractBalances(balances, thresholdCopy)
	if err != nil {
		return true
	}

	return false
}

func AddBalancesAndAssertDoesntExceedThreshold(currTally []*Balance, toAdd []*Balance, threshold []*Balance) ([]*Balance, error) {
	err := *new(error)
	//If we transferAsMuchAsPossible, we need to increment the currTally by all that we can
	//We then need to return the updated toAdd

	currTally, err = AddBalances(toAdd, currTally)
	if err != nil {
		return []*Balance{}, err
	}

	//Check if we exceed the threshold; will underflow if we do exceed it
	doExceed := DoBalancesExceedThreshold(currTally, threshold)
	if doExceed {
		return []*Balance{}, sdkerrors.Wrapf(ErrExceedsThreshold, "curr tally plus added amounts exceeds threshold")
	}

	return currTally, nil	
}

func AreBalancesEqual(expected []*Balance, actual []*Balance, checkZeroBalances bool) bool {
	expected = DeepCopyBalances(expected)
	actual = DeepCopyBalances(actual)

	if !checkZeroBalances {
		expected = FilterZeroBalances(expected)
		actual = FilterZeroBalances(actual)
	}

	actual, err := SubtractBalances(expected, actual)
	if len(actual) != 0 || err != nil {
		return false
	}

	return true
}

func DeepCopyBalances(balances []*Balance) []*Balance {
	newBalances := []*Balance{}
	for _, approval := range balances {
		balanceToAdd := &Balance{
			Amount: 			 approval.Amount,
		}
		for _, badgeId := range approval.BadgeIds {
			balanceToAdd.BadgeIds = append(balanceToAdd.BadgeIds, &UintRange{
				Start: badgeId.Start,
				End:   badgeId.End,
			})
		}

		for _, time := range approval.OwnedTimes {
			balanceToAdd.OwnedTimes = append(balanceToAdd.OwnedTimes, &UintRange{
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
func HandleDuplicateBadgeIds(balances []*Balance) ([]*Balance, error) {
	newBalances := []*Balance{}
	err := *new(error)
	for _, balance := range balances {
		for _, badgeId := range balance.BadgeIds {
			for _, time := range balance.OwnedTimes {
				newBalances, err = AddBalance(newBalances, &Balance{
					Amount:         balance.Amount,
					BadgeIds:       []*UintRange{badgeId},
					OwnedTimes: []*UintRange{time},
				})
				if err != nil {
					return []*Balance{}, err
				}
			}
		}
	}

	return newBalances, nil
}

// // Gets the balances for a specific ID
// func GetBalancesForId(id sdkmath.Uint, balances []*Balance) []*Balance {
// 	matchingBalances := []*Balance{}
// 	for _, balance := range balances {
// 		found := SearchUintRangesForUint(id, balance.BadgeIds)
// 		if found {
// 			matchingBalances = append(matchingBalances, &Balance{
// 				Amount:   balance.Amount,
// 				BadgeIds: []*UintRange{{Start: id, End: id}},
// 				OwnedTimes: balance.OwnedTimes,
// 			})
// 		}
// 	}
// 	return matchingBalances
// }

// Updates the balance for a specific ids from what it currently is to newAmount. No add/subtract logic. Just set it.
func UpdateBalance(newBalance *Balance, balances []*Balance) ([]*Balance, error) {
	//Can maybe optimize this in the future by doing this all in one loop instead of deleting then setting
	// ranges = SortAndMergeOverlapping(ranges)
	err := *new(error)
	balances, err = DeleteBalances(newBalance.BadgeIds, newBalance.OwnedTimes, balances)
	if err != nil {
		return balances, err
	}

	balances, err = SetBalance(newBalance, balances)
	if err != nil {
		return balances, err
	}

	return balances, nil
}

// Gets the balances for specified ID ranges. Returns a new []*Balance where only the specified ID ranges and their balances are included. Appends balance == 0 objects so all IDs are accounted for, even if not found.
func GetBalancesForIds(idRanges []*UintRange, times []*UintRange, balances []*Balance) (newBalances []*Balance, err error) {
	fetchedBalances := []*Balance{}

	currPermissionDetails := []*UniversalPermissionDetails{}
	for _, balanceObj := range balances {
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.OwnedTimes {
				currPermissionDetails = append(currPermissionDetails, &UniversalPermissionDetails{
					BadgeId:            currRange,
					OwnershipTime:       currTime,
					TransferTime:       &UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					TimelineTime: 		 &UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					ToMapping:          &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					FromMapping:        &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					InitiatedByMapping: &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					ArbitraryValue:     balanceObj.Amount,
				})
			}
		}
	}

	toFetchPermissionDetails := []*UniversalPermissionDetails{}
	for _, rangeToFetch := range idRanges {
		for _, timeToFetch := range times {
			toFetchPermissionDetails = append(toFetchPermissionDetails, &UniversalPermissionDetails{
				BadgeId:            rangeToFetch,
				OwnershipTime:       timeToFetch,
				TransferTime:       &UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
				TimelineTime: 			&UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
				ToMapping:          &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
				FromMapping:        &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
				InitiatedByMapping: &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
			},
			)
		}
	}

	overlaps, _, inNewButNotOld := GetOverlapsAndNonOverlaps(currPermissionDetails, toFetchPermissionDetails)
	//For all overlaps, we simply return the amount
	for _, overlapObject := range overlaps {
		overlap := overlapObject.Overlap
		amount := overlapObject.FirstDetails.ArbitraryValue.(sdkmath.Uint)

		fetchedBalances = append(fetchedBalances, &Balance{
			Amount:         amount,
			BadgeIds:       []*UintRange{overlap.BadgeId},
			OwnedTimes: []*UintRange{overlap.OwnershipTime},
		})
	}

	//For those that were in toFetch but not currBalances, we return amount == 0
	for _, detail := range inNewButNotOld {
		fetchedBalances = append(fetchedBalances, &Balance{
			Amount:         sdkmath.NewUint(0),
			BadgeIds:       []*UintRange{detail.BadgeId},
			OwnedTimes: []*UintRange{detail.OwnershipTime},
		})
	}

	return fetchedBalances, nil
}

func AddBalances(balancesToAdd []*Balance, balances []*Balance) ([]*Balance, error) {
	err := *new(error)
	for _, balance := range balancesToAdd {
		balances, err = AddBalance(balances, balance)
		if err != nil {
			return balances, err
		}
	}

	return balances, nil
}

// Adds a balance to all ids specified in []ranges
func AddBalance(existingBalances []*Balance, balanceToAdd *Balance) ([]*Balance, error) {
	currBalances, err := GetBalancesForIds(balanceToAdd.BadgeIds, balanceToAdd.OwnedTimes, existingBalances)
	if err != nil {
		return existingBalances, err
	}

	for _, balance := range currBalances {
		balance.Amount, err = SafeAdd(balance.Amount, balanceToAdd.Amount)
		if err != nil {
			return existingBalances, err
		}

		existingBalances, err = UpdateBalance(balance, existingBalances)
		if err != nil {
			return existingBalances, err
		}
	}

	return existingBalances, nil
}

func SubtractBalances(balancesToSubtract []*Balance, balances []*Balance) ([]*Balance, error) {
	err := *new(error)

	for _, balance := range balancesToSubtract {
		balances, err = SubtractBalance(balances, balance)
		if err != nil {
			return balances, err
		}
	}

	return balances, nil
}

// Subtracts a balance to all ids specified in []ranges
func SubtractBalance(balances []*Balance, balanceToRemove *Balance) ([]*Balance, error) {

	currBalances, err := GetBalancesForIds(balanceToRemove.BadgeIds, balanceToRemove.OwnedTimes, balances)
	if err != nil {
		return balances, err
	}

	for _, currBalanceObj := range currBalances {
		currBalanceObj.Amount, err = SafeSubtract(currBalanceObj.Amount, balanceToRemove.Amount)
		if err != nil {
			return balances, err
		}

		balances, err = UpdateBalance(currBalanceObj, balances)
		if err != nil {
			return balances, err
		}
	}

	return balances, nil
}

// Subtracts a balance to all ids specified in []ranges
func GetMaxSubtractableBalances(balances []*Balance, balanceToRemove *Balance) ([]*Balance, error) {
	currBalances, err := GetBalancesForIds(balanceToRemove.BadgeIds, balanceToRemove.OwnedTimes, balances)
	if err != nil {
		return balances, err
	}

	subtractableBalances := []*Balance{}
	for _, currBalanceObj := range currBalances {
		subtractableBalances = append(subtractableBalances, &Balance{
			Amount:         balanceToRemove.Amount,
			BadgeIds: 			currBalanceObj.BadgeIds,
			OwnedTimes: currBalanceObj.OwnedTimes,
		})

		if currBalanceObj.Amount.LT(balanceToRemove.Amount) {
			subtractableBalances[len(subtractableBalances)-1].Amount = currBalanceObj.Amount
		}
	}

	return balances, nil
}

// Deletes the balance for a specific id.
func DeleteBalances(rangesToDelete []*UintRange, timesToDelete []*UintRange, balances []*Balance) ([]*Balance, error) {
	newBalances := []*Balance{}

	for _, balanceObj := range balances {
		currPermissionDetails := []*UniversalPermissionDetails{}
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.OwnedTimes {
				currPermissionDetails = append(currPermissionDetails, &UniversalPermissionDetails{
					BadgeId:            currRange,
					OwnershipTime:       currTime,
					TransferTime:       &UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					TimelineTime: 		 &UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					ToMapping:          &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					FromMapping:        &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					InitiatedByMapping: &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
				})
			}
		}

		toDeletePermissionDetails := []*UniversalPermissionDetails{}
		for _, rangeToDelete := range rangesToDelete {
			for _, timeToDelete := range timesToDelete {
				toDeletePermissionDetails = append(toDeletePermissionDetails, &UniversalPermissionDetails{
					BadgeId:            rangeToDelete,
					OwnershipTime:       timeToDelete,
					TransferTime:       &UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					TimelineTime: 		 &UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					ToMapping:          &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					FromMapping:        &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					InitiatedByMapping: &AddressMapping{Addresses: []string{}, IncludeAddresses: false},
				},
				)
			}
		}

		_, inOldButNotNew, _ := GetOverlapsAndNonOverlaps(currPermissionDetails, toDeletePermissionDetails)
		for _, remainingBalance := range inOldButNotNew {
			newBalances = append(newBalances, &Balance{
				Amount:         balanceObj.Amount,
				BadgeIds:       []*UintRange{remainingBalance.BadgeId},
				OwnedTimes: []*UintRange{remainingBalance.OwnershipTime},
			})
		}
	}

	return newBalances, nil
}

// Sets the balance for a specific id. Precondition: assumes balance does not exist.
func SetBalance(newBalance *Balance, balances []*Balance) ([]*Balance, error) {
	if newBalance.Amount.IsZero() {
		return balances, nil
	}

	balances = append(balances, newBalance)

	return balances, nil
}
