package keeper

import (
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Search ID ranges for a specific ID. Return (idx, true) if found. And (-1, false) if not.
func SearchIdRangesForId(id uint64, idRanges []*types.IdRange) (int, bool) {
	//Binary search because ID ranges will be sorted
	low := 0
	high := len(idRanges) - 1
	for low <= high {
		median := int(uint(low+high) >> 1)

		currRange := NormalizeIdRange(idRanges[median])

		if currRange.Start <= id && currRange.End >= id {
			return median, true
		} else if currRange.Start > id {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	return -1, false
}

// Search a set of ranges to find what indexes a specific ID range overlaps. Return overlapping idxs as a IdRange, true if found. And empty IdRang, false if not
func GetIdxSpanForRange(targetRange *types.IdRange, targetIdRanges []*types.IdRange) (*types.IdRange, bool) {
	//its search for start, if found set to that
	//if not found, set to insertIdx + 0 (because we already incremented by 1)
	//if end is found, set to that
	//else set to insertIdx - 1 (because we already incremented by 1)
	targetRange = NormalizeIdRange(targetRange)
	idRanges := targetIdRanges

	startIdx, startFound := SearchIdRangesForId(targetRange.Start, idRanges)
	if !startFound {
		startIdx = GetIdxToInsertForNewId(targetRange.Start, idRanges)
	}

	endIdx, endFound := SearchIdRangesForId(targetRange.End, idRanges)
	if !endFound {
		endIdx = GetIdxToInsertForNewId(targetRange.End, idRanges) - 1
	}

	if startIdx <= endIdx {
		return &types.IdRange{
			Start: uint64(startIdx),
			End:   uint64(endIdx),
		}, true
	} else {
		return &types.IdRange{}, false
	}
}

// Handle the case where it omits an empty IdRange because Start && End == 0. This is in the case where we have a non-empty balance and an empty idRanges.
func GetIdRangesWithOmitEmptyCaseHandled(ids []*types.IdRange) []*types.IdRange {
	if len(ids) == 0 {
		ids = append(ids, &types.IdRange{})
	}
	return ids
}

// Gets the number range to insert with the additional convention of storing end = 0 when end == start
func GetIdRangeToInsert(start uint64, end uint64) *types.IdRange {
	if end == start {
		end = 0
	}

	return &types.IdRange{
		Start: start,
		End:   end,
	}
}

// Normalizes an existing ID range with the additional convention of storing end == 0 when end == start
func NormalizeIdRange(rangeToNormalize *types.IdRange) *types.IdRange {
	if rangeToNormalize.End == 0 {
		rangeToNormalize.End = rangeToNormalize.Start
	}

	return &types.IdRange{
		Start: rangeToNormalize.Start,
		End:   rangeToNormalize.End,
	}
}

// Assumes id is not already in a range. Gets the index to insert at. Ex. [10, 20, 30] and inserting 25 would return index 2
func GetIdxToInsertForNewId(id uint64, targetIds []*types.IdRange) int {
	//Since we assume the id is not already in there, we can just compare start positions of the existing idRanges and see where it falls between
	ids := targetIds
	if len(ids) == 0 {
		return 0
	}

	if ids[0].Start > id { //assumes not in already so we don't have to handle that case
		return 0
	} else if ids[len(ids)-1].End < id {
		return len(ids)
	}

	low := 0
	high := len(ids) - 2
	median := 0
	for low <= high {
		median = int(uint(low+high) >> 1)
		currRange := NormalizeIdRange(ids[median])
		nextRange := NormalizeIdRange(ids[median+1])

		if currRange.Start < id && nextRange.Start > id {
			break
		} else if currRange.Start > id {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	currRange := NormalizeIdRange(ids[median])
	insertIdx := median + 1
	if currRange.Start <= id {
		insertIdx = median
	}

	//We return insertIdx + 1 because this function uses (curr, next) pairs so if we find that we have to insert at a certain (curr, next) pair where insertIdx == currIdx, we actually insert at the "next" position
	return insertIdx + 1
}

// We inserted a new id at insertedAtIdx, this can cause the prev or next to have to merge if id + 1 or id - 1 overlaps with prev or next range. Handle this here.
func MergePrevOrNextIfPossible(targetIds []*types.IdRange, insertedAtIdx int) []*types.IdRange {
	//Handle cases where we need to merge with the previous or next range
	needToMergeWithPrev := false
	needToMergeWithNext := false
	prevStartIdx := uint64(0)
	nextEndIdx := uint64(0)
	ids := targetIds
	
	id := NormalizeIdRange(ids[insertedAtIdx])
	idStart := id.Start
	idEnd := id.End

	if insertedAtIdx > 0 {
		prev := NormalizeIdRange(ids[insertedAtIdx-1])
		prevStartIdx = prev.Start
		prevEndIdx := prev.End

		if prevEndIdx+1 == idStart {
			needToMergeWithPrev = true
		}
	}

	if insertedAtIdx < len(ids)-1 {
		next := NormalizeIdRange(ids[insertedAtIdx+1])
		nextStartIdx := next.Start
		nextEndIdx = next.End

		if nextStartIdx-1 == idEnd {
			needToMergeWithNext = true
		}
	}

	mergedIds := []*types.IdRange{}
	// 4 Cases: Need to merge with both, just next, just prev, or neither
	if needToMergeWithPrev && needToMergeWithNext {
		mergedIds = append(mergedIds, ids[:insertedAtIdx-1]...)
		mergedIds = append(mergedIds, GetIdRangeToInsert(prevStartIdx, nextEndIdx))
		mergedIds = append(mergedIds, ids[insertedAtIdx+2:]...)
	} else if needToMergeWithPrev {
		mergedIds = append(mergedIds, ids[:insertedAtIdx-1]...)
		mergedIds = append(mergedIds, GetIdRangeToInsert(prevStartIdx, idEnd))
		mergedIds = append(mergedIds, ids[insertedAtIdx+1:]...)
	} else if needToMergeWithNext {
		mergedIds = append(mergedIds, ids[:insertedAtIdx]...)
		mergedIds = append(mergedIds, GetIdRangeToInsert(idStart, nextEndIdx))
		mergedIds = append(mergedIds, ids[insertedAtIdx+2:]...)
	} else {
		mergedIds = ids
	}

	return mergedIds
}

// Inserts a range into its correct position. Assumes range is already deleted and not present at all, so we only search for where start fits in.
func InsertRangeToIdRanges(rangeToAdd *types.IdRange, targetIds []*types.IdRange) []*types.IdRange {
	ids := targetIds
	newIds := []*types.IdRange{}
	insertIdAtIdx := 0
	rangeToAdd = NormalizeIdRange(rangeToAdd)
	lastRange := NormalizeIdRange(ids[len(ids)-1])

	//Three cases: Goes at beginning, end, or somewhere in the middle
	if ids[0].Start > rangeToAdd.End {
		newIds = append(newIds, GetIdRangeToInsert(rangeToAdd.Start, rangeToAdd.End))
		newIds = append(newIds, ids...)
	} else if lastRange.End < rangeToAdd.Start {
		insertIdAtIdx = len(ids)
		newIds = append(newIds, ids...)
		newIds = append(newIds, GetIdRangeToInsert(rangeToAdd.Start, rangeToAdd.End))
	} else {
		insertIdAtIdx = GetIdxToInsertForNewId(rangeToAdd.Start, ids) //Only lookup start since we assume the whole range isn't included already
		newIds = append(newIds, ids[:insertIdAtIdx]...)
		newIds = append(newIds, GetIdRangeToInsert(rangeToAdd.Start, rangeToAdd.End))
		newIds = append(newIds, ids[insertIdAtIdx:]...)
	}

	newIds = MergePrevOrNextIfPossible(newIds, insertIdAtIdx)

	return newIds
}

// Removes all ids within an id range from an id range. Removing can make this range be split into 0, 1, or 2 new ranges.
func RemoveIdsFromIdRange(rangeToRemove *types.IdRange, rangeObject *types.IdRange) []*types.IdRange {
	newRanges := []*types.IdRange{}
	rangeToRemove = NormalizeIdRange(rangeToRemove)
	rangeObject = NormalizeIdRange(rangeObject)


	if rangeToRemove.Start > rangeObject.Start && rangeToRemove.End < rangeObject.End {
		// Completely in the middle
		newRanges = append(newRanges, GetIdRangeToInsert(rangeObject.Start, rangeToRemove.Start-1))
		newRanges = append(newRanges, GetIdRangeToInsert(rangeToRemove.End+1, rangeObject.End))
	} else if rangeToRemove.Start <= rangeObject.Start && rangeToRemove.End >= rangeObject.End {
		// Overlaps both; remove whole thing
		// Do nothing
	} else if rangeToRemove.Start <= rangeObject.Start && rangeToRemove.End < rangeObject.End && rangeToRemove.End >= rangeObject.Start {
		// Still have some at the end
		newRanges = append(newRanges, GetIdRangeToInsert(rangeToRemove.End+1, rangeObject.End))
	} else if rangeToRemove.Start > rangeObject.Start && rangeToRemove.End >= rangeObject.End && rangeToRemove.Start <= rangeObject.End {
		// Still have some at the start
		newRanges = append(newRanges, GetIdRangeToInsert(rangeObject.Start, rangeToRemove.Start-1))	
	} else {
		// Doesn't overlap at all
		newRanges = append(newRanges, GetIdRangeToInsert(rangeObject.Start, rangeObject.End))
	}

	return newRanges
}

//Will sort the ID ranges in order and merge overlapping IDs if we can
func SortIdRangesAndMergeIfNecessary(ids []*types.IdRange) []*types.IdRange {
	//Insertion sort in order of range.Start. If two have same range.Start, sort by range.End.
	var n = len(ids)
    for i := 1; i < n; i++ {
        j := i
        for j > 0 {
            if ids[j-1].Start > ids[j].Start {
                ids[j-1], ids[j] = ids[j], ids[j-1]
            } else if ids[j-1].Start == ids[j].Start && ids[j-1].End > ids[j].End {
				ids[j-1], ids[j] = ids[j], ids[j-1]
			}
            j = j - 1
        }
    }
	
	//Merge overlapping ranges
	if n > 0 {
		newIdRanges := []*types.IdRange{}
		newIdRanges = append(newIdRanges, GetIdRangeToInsert(ids[0].Start, ids[0].End))

		for i := 1; i < n; i++ {
			prevInsertedRange := NormalizeIdRange(newIdRanges[len(newIdRanges)-1])
			currRange := NormalizeIdRange(ids[i])

			if currRange.Start == prevInsertedRange.Start {
				//Both have same start, so we set to currRange.End because currRange.End is greater due to our sorting
				newIdRanges[len(newIdRanges)-1].End = currRange.End
			} else if currRange.End > prevInsertedRange.End {
				//We have different starts and curr end is greater than prev end
				if currRange.Start > prevInsertedRange.End + 1 {
					//We have a gap between the prev range end and curr range start, so we just append currRange
					newIdRanges = append(newIdRanges, GetIdRangeToInsert(currRange.Start, currRange.End))
				} else {
					//they overlap and we can merge them
					newIdRanges[len(newIdRanges)-1].End = currRange.End
				}
			} 
			// else if currRange.End <= prevInsertedRange.End {
				//Start must be >= because it is sorted, so we can just skip this range since currRange is already completely enclosed by prevRange
			// }
		}
		return newIdRanges
	} else {
		return ids
	}
}

func GetIdRangesToInsertToStorage(idRanges []*types.IdRange) []*types.IdRange {
	newIdRanges := []*types.IdRange{}
	for _, idRange := range idRanges {
		newIdRanges = append(newIdRanges, GetIdRangeToInsert(idRange.Start, idRange.End))
	}
	return newIdRanges
}