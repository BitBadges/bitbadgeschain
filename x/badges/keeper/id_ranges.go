package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Helper function to make code more readable
func CreateIdRange(start sdk.Uint, end sdk.Uint) *types.IdRange {
	return &types.IdRange{
		Start: start,
		End:   end,
	}
}

// Search ID ranges for a specific ID. Return (idx, true) if found. And (-1, false) if not.
// Assumes ID ranges are sorted.
func SearchIdRangesForId(id sdk.Uint, idRanges []*types.IdRange) (int, bool) {
	idRanges = SortAndMergeOverlapping(idRanges) // Just in case

	//Binary search because ID ranges will be sorted
	low := 0
	high := len(idRanges) - 1
	for low <= high {
		median := int(uint(low+high) >> 1)

		currRange := idRanges[median]

		if currRange.Start.LTE(id) && currRange.End.GTE(id) {
			return median, true
		} else if currRange.Start.GT(id) {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	return -1, false
}

// Search a set of ranges to find what indexes a specific ID range overlaps.
// Return overlapping idxs as a IdRange, true if found. And empty IdRange, false if not
// Inclusive (aka start and end idx will both have overlaps somewhere)
func GetIdxSpanForRange(rangeToCheck *types.IdRange, currRanges []*types.IdRange) (*types.IdRange, bool) {
	//Note GetIdxToInsertForNewId returns the index to insert at (i.e. the following idx)
	//For start, this is what we want because we want the first non-overlapping range
	//For end, we want the idx before (i.e. idx - 1) because that is the last overlapping range
	idRanges := currRanges

	startIdx, startFound := SearchIdRangesForId(rangeToCheck.Start, idRanges)
	if !startFound {
		startIdx, _ = GetIdxToInsertForNewId(rangeToCheck.Start, idRanges) //ignore error because we know it's not found
	}

	endIdx, endFound := SearchIdRangesForId(rangeToCheck.End, idRanges)
	if !endFound {
		endIdx, _ = GetIdxToInsertForNewId(rangeToCheck.End, idRanges) //ignore error because we know it's not found
		endIdx--                         	
	}

	if startIdx <= endIdx {
		return &types.IdRange{
			Start: sdk.NewUint(uint64(startIdx)),
			End:   sdk.NewUint(uint64(endIdx)),
		}, true
	} else {
		return &types.IdRange{}, false
	}
}

// Assumes given ID is not already in a range. We recommend calling SearchIdRangesForId first.
// Gets the index to insert at. Ex. [{0-10}, {10-20}, {30-40}] and inserting 25 would return index 2
func GetIdxToInsertForNewId(id sdk.Uint, targetIds []*types.IdRange) (int, error) {
	targetIds = SortAndMergeOverlapping(targetIds) // Just in case

	_, found := SearchIdRangesForId(id, targetIds)
	if found {
		return -1, ErrIdInRange
	}

	//Since we assume the id is not already in there, we can just compare start positions of the existing idRanges
	ids := targetIds
	if len(ids) == 0 {
		return 0, nil
	}

	// Check if id is before the first range or after the last range
	if ids[0].Start.GT(id) {
		return 0, nil
	} else if ids[len(ids)-1].End.LT(id) {
		return len(ids), nil
	}

	//If length == 1, then it should never reach here because we already checked if it's before or after
	//and assume it's not in the range
	if len(ids) == 1 {
		return -1, ErrIdInRange //Should never reach here but just in case
	}


	//Binary search by looking at two ranges at a time [..., {curr}, {next}, ...]
	low := 0
	high := len(ids) - 2
	median := 0
	for low <= high {
		median = int(uint(low+high) >> 1)
		currRange := ids[median]
		nextRange := ids[median+1]

		//If id is in between curr and next, then we found the index to insert at
		//Note that we assume id is not already in the range and sorted so can just check starts
		if currRange.Start.LT(id) && nextRange.Start.GT(id) {
			break
		} else if currRange.Start.GT(id) {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	//We return median + 1 because median == curr and we want to insert in between {curr} and {next}
	return median + 1, nil
}

// Inserts a range into its correct position.
// Assumes whole range is not present at all. Thus, we only search for where start fits in.
func InsertRangeToIdRanges(rangeToAdd *types.IdRange, targetIds []*types.IdRange) ([]*types.IdRange, error) {
	//Validation check; make sure rangeToAdd is not already in targetIds
	for _, idRange := range targetIds {
		_, removed := RemoveIdsFromIdRange(idRange, rangeToAdd)
		if removed {
			return nil, ErrIdAlreadyInRanges
		}
	}

	ids := targetIds
	newIds := []*types.IdRange{}
	insertIdAtIdx := 0
	lastRange := ids[len(ids)-1]

	err := *new(error)
	//Three cases: Goes at beginning, end, or somewhere in the middle
	if ids[0].Start.GT(rangeToAdd.End) {
		newIds = append(newIds, CreateIdRange(rangeToAdd.Start, rangeToAdd.End))
		newIds = append(newIds, ids...)
	} else if lastRange.End.LT(rangeToAdd.Start) {
		newIds = append(newIds, ids...)
		newIds = append(newIds, CreateIdRange(rangeToAdd.Start, rangeToAdd.End))
	} else {
		insertIdAtIdx, err = GetIdxToInsertForNewId(rangeToAdd.Start, ids) //Only lookup start since we assume the whole range isn't included already
		if err != nil {
			return nil, err
		}
		newIds = append(newIds, ids[:insertIdAtIdx]...)
		newIds = append(newIds, CreateIdRange(rangeToAdd.Start, rangeToAdd.End))
		newIds = append(newIds, ids[insertIdAtIdx:]...)
	}

	newIds = SortAndMergeOverlapping(newIds)

	return newIds, nil
}

// Removes all ids within an id range from an id range.
// Removing can make this range be split into 0, 1, or 2 new ranges.
// Returns if anything was removed or not
func RemoveIdsFromIdRange(idxsToRemove *types.IdRange, rangeObject *types.IdRange) ([]*types.IdRange, bool) {
	if idxsToRemove.End.LT(rangeObject.Start) || idxsToRemove.Start.GT(rangeObject.End) {
		// idxsToRemove doesn't overlap with rangeObject, so nothing is removed
		return []*types.IdRange{rangeObject}, false
	}

	var newRanges []*types.IdRange
	if idxsToRemove.Start.LTE(rangeObject.Start) && idxsToRemove.End.GTE(rangeObject.End) {
		// idxsToRemove fully contains rangeObject, so nothing is left
		return newRanges, true
	}

	if idxsToRemove.Start.GT(rangeObject.Start) {
		// There's a range before idxsToRemove
		// Underflow is not possible because idxsToRemove.Start.GT(rangeObject.Start
		newRanges = append(newRanges, &types.IdRange{
			Start: rangeObject.Start,
			End:   idxsToRemove.Start.SubUint64(1),
		})
	}

	if idxsToRemove.End.LT(rangeObject.End) {
		// There's a range after idxsToRemove
		// Overflow is not possible because idxsToRemove.End.LT(rangeObject.End
		newRanges = append(newRanges, &types.IdRange{
			Start: idxsToRemove.End.AddUint64(1),
			End:   rangeObject.End,
		})
	}

	return newRanges, true
}


//IMPORTANT: Note this function was copied to the types validation.go file. If you change this, change that as well and vice versa.

// Will sort the ID ranges in order and merge overlapping IDs if we can
func SortAndMergeOverlapping(ids []*types.IdRange) []*types.IdRange {
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
		newIdRanges := []*types.IdRange{CreateIdRange(ids[0].Start, ids[0].End)}
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
					newIdRanges = append(newIdRanges, CreateIdRange(currRange.Start, currRange.End))
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