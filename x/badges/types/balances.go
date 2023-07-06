package types

import (
	sdkmath "cosmossdk.io/math"

	"math"
)

func HandleDuplicateBadgeIds(balances []*Balance) ([]*Balance, error) {
	newBalances := []*Balance{}
	err := *new(error)
	for _, balance := range balances {
		for _, badgeId := range balance.BadgeIds {
			for _, time := range balance.Times {
				newBalances, err = AddBalance(newBalances, &Balance{
					Amount:   balance.Amount,
					BadgeIds: []*IdRange{badgeId},
					Times: []*IdRange{time},
				})
				if err != nil {
					return []*Balance{}, err
				}
			}
		}
	}

	return newBalances, nil
}

// Gets the balances for a specific ID. Assumes balances are sorted
func GetBalancesForId(id sdkmath.Uint, balances []*Balance) []*Balance {
	matchingBalances := []*Balance{}
	for _, balance := range balances {
		found := SearchIdRangesForId(id, balance.BadgeIds)
		if found {
			matchingBalances = append(matchingBalances, &Balance{
				Amount:   balance.Amount,
				BadgeIds: []*IdRange{{Start: id, End: id}},
				Times: balance.Times,
			})
		}
	}
	return matchingBalances
}

// Updates the balance for a specific ids from what it currently is to newAmount.
func UpdateBalance(newBalance *Balance, balances []*Balance) ([]*Balance, error) {
	//Can maybe optimize this in the future by doing this all in one loop instead of deleting then setting
	// ranges = SortAndMergeOverlapping(ranges)
	err := *new(error)
	balances, err = DeleteBalances(newBalance.BadgeIds, newBalance.Times, balances)
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
func GetBalancesForIds(idRanges []*IdRange, times []*IdRange, balances []*Balance) (newBalances []*Balance, err error) {
	fetchedBalances := []*Balance{}

	currPermissionDetails := []*UniversalPermissionDetails{}
	for _, balanceObj := range balances {
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.Times {
				currPermissionDetails = append(currPermissionDetails, &UniversalPermissionDetails{
					BadgeId: currRange,
					TimelineTime: currTime,
					TransferTime: &IdRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
					ToMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
					FromMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
					InitiatedByMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
					ArbitraryValue: balanceObj.Amount,
				})
			}
		}
	}

	toFetchPermissionDetails := []*UniversalPermissionDetails{}
	for _, rangeToFetch := range idRanges {
		for _, timeToFetch := range times {
			toFetchPermissionDetails = append(toFetchPermissionDetails, &UniversalPermissionDetails{
					BadgeId: rangeToFetch,
					TimelineTime: timeToFetch,
					TransferTime: &IdRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
					ToMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
					FromMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
					InitiatedByMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
				},
			)
		}
	}


	overlaps, _, inNewButNotOld := GetOverlapsAndNonOverlaps(currPermissionDetails, toFetchPermissionDetails)
	for _, overlapObject := range overlaps {
		overlap := overlapObject.Overlap
		amount := overlapObject.FirstDetails.ArbitraryValue.(sdkmath.Uint)

		fetchedBalances = append(fetchedBalances, &Balance{
			Amount:  amount,
			BadgeIds: []*IdRange{overlap.BadgeId},
			Times: []*IdRange{overlap.TimelineTime},
		})
	}

	for _, detail := range inNewButNotOld {
		fetchedBalances = append(fetchedBalances, &Balance{
			Amount:   sdkmath.NewUint(0),
			BadgeIds: []*IdRange{detail.BadgeId},
			Times: []*IdRange{detail.TimelineTime},
		})
	}
	

	return fetchedBalances, nil
}

// Adds a balance to all ids specified in []ranges
func AddBalance(existingBalances []*Balance, balanceToAdd *Balance) ([]*Balance, error) {
	currBalances, err := GetBalancesForIds(balanceToAdd.BadgeIds, balanceToAdd.Times, existingBalances)
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

// Subtracts a balance to all ids specified in []ranges
func SubtractBalance(balances []*Balance, balanceToRemove *Balance) ([]*Balance, error) {

	currBalances, err := GetBalancesForIds(balanceToRemove.BadgeIds, balanceToRemove.Times, balances)
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

// Deletes the balance for a specific id.
func DeleteBalances(rangesToDelete []*IdRange, timesToDelete []*IdRange, balances []*Balance) ([]*Balance, error) {
	newBalances := []*Balance{}

	for _, balanceObj := range balances {
		currPermissionDetails := []*UniversalPermissionDetails{}
		for _, currRange := range balanceObj.BadgeIds {
			for _, currTime := range balanceObj.Times {
				currPermissionDetails = append(currPermissionDetails, &UniversalPermissionDetails{
					BadgeId: currRange,
					TimelineTime: currTime,
					TransferTime: &IdRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
					ToMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
					FromMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
					InitiatedByMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
				})
			}
		}

		toDeletePermissionDetails := []*UniversalPermissionDetails{}
		for _, rangeToDelete := range rangesToDelete {
			for _, timeToDelete := range timesToDelete {
				toDeletePermissionDetails = append(toDeletePermissionDetails, &UniversalPermissionDetails{
						BadgeId: rangeToDelete,
						TimelineTime: timeToDelete,
						TransferTime: &IdRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
						ToMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
						FromMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
						InitiatedByMapping: &AddressMapping{ Addresses: []string{}, IncludeOnlySpecified: false },
					},
				)
			}
		}

		_, inOldButNotNew, _ := GetOverlapsAndNonOverlaps(currPermissionDetails, toDeletePermissionDetails)
		for _, remainingBalance := range inOldButNotNew {
			newBalances = append(newBalances, &Balance{
				Amount:   balanceObj.Amount,
				BadgeIds: []*IdRange{remainingBalance.BadgeId},
				Times: []*IdRange{remainingBalance.TimelineTime},
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
