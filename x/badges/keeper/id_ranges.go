package keeper

import (
	"math"

	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

// Binary search ID ranges for a specific ID. Return idx, true if found. And -1, false if not
func SearchIdRangesForId(id uint64, idRanges []*types.IdRange) (int, bool) {
	low := 0
	high := len(idRanges) - 1
	for low <= high {
		median := int(uint(low+high) >> 1)

		currRange := idRanges[median]
		if currRange.End == 0 {
			currRange.End = currRange.Start // If end == 0, set it to start by convention (done in order to save space)
		}

		if idRanges[median].Start <= id && idRanges[median].End >= id {
			return median, true
		} else if idRanges[median].Start > id {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	return -1, false
}


// Binary search ID ranges for a specific ID range idxs. Return idxs, true if found. And -1, false if not
func GetIdxSpanForRange(targetRange types.IdRange, idRanges []*types.IdRange) (types.IdRange, bool) {
	low := 0
	high := len(idRanges) - 1
	targetId := targetRange.Start
	for low <= high {
		median := int(uint(low+high) >> 1)

		currRange := idRanges[median]
		if currRange.End == 0 {
			currRange.End = currRange.Start // If end == 0, set it to start by convention (done in order to save space)
		}

		if idRanges[median].Start <= targetId && idRanges[median].End >= targetId {
			//check both sides of this median
			start := median
			end := median

			targetEnd := targetRange.End
			if targetEnd == 0 {
				targetEnd = targetRange.Start
			}

			median += 1

			for median < len(idRanges) {
				currRange = idRanges[median]
				if currRange.End == 0 {
					currRange.End = currRange.Start // If end == 0, set it to start by convention (done in order to save space)
				}

				if (currRange.Start <= targetEnd && currRange.End >= targetEnd) {
					end = median
					median += 1
				} else {
					break
				}
			}

			return *GetIdRangeToInsert(uint64(start), uint64(end)), true
		} else if idRanges[median].Start > targetId {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	return types.IdRange{}, false
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

// Assumes id is not already in a range. Gets the index to insert at. Ex. [10, 20, 30] and inserting 25 would return index 2
func GetIdxToInsertForNewId(id uint64, targetIds []*types.IdRange) int {
	low := 0
	high := len(targetIds) - 2
	median := 0
	for low <= high {
		median = int(uint(low+high) >> 1)
		if targetIds[median].Start < id && targetIds[median+1].Start > id {
			break
		} else if targetIds[median].Start > id {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	insertIdx := median + 1
	if targetIds[median].Start <= id {
		insertIdx = median
	}

	//We return insertIdx + 1 because this function uses (curr, next) pairs so if we find that we have to insert at a certain (curr, next) pair where insertIdx == currIdx, we actually insert in between
	return insertIdx + 1
}

// Assumes id is not already in a range. Gets the index to insert at. Ex. [10, 20, 30] and inserting 25 would return index 2
func GetIdxToInsertForNewIdRange(rangeToAdd types.IdRange, targetIds []*types.IdRange) int {
	low := 0
	high := len(targetIds) - 2
	median := 0
	for low <= high {
		median = int(uint(low+high) >> 1)
		currRange := targetIds[median]
		if currRange.End == 0 {
			currRange.End = currRange.Start // If end == 0, set it to start by convention (done in order to save space)
		}

		if currRange.End < rangeToAdd.Start && targetIds[median+1].Start > rangeToAdd.End {
			break
		} else if currRange.End > rangeToAdd.Start {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	insertIdx := median + 1
	if targetIds[median].End <= rangeToAdd.Start {
		insertIdx = median
	}

	return insertIdx
}

// We inserted a new id at insertedAtIdx, this can cause the prev or next to have to merge if id + 1 or id - 1 overlaps with prev or next range. Handle this here.
func MergePrevOrNextIfNecessary(ids []*types.IdRange, insertedAtIdx int) []*types.IdRange {
	//Handle cases where we need to merge with the previous or next range
	needToMergeWithPrev := false
	needToMergeWithNext := false
	prevStartIdx := uint64(0)
	nextEndIdx := uint64(0)
	id := ids[insertedAtIdx].Start

	if insertedAtIdx > 0 {
		prevStartIdx = ids[insertedAtIdx-1].Start
		prevEndIdx := ids[insertedAtIdx-1].End
		if prevEndIdx == 0 {
			prevEndIdx = ids[insertedAtIdx-1].Start
		}

		if prevEndIdx+1 == id {
			needToMergeWithPrev = true
		}
	}

	if insertedAtIdx < len(ids)-1 {
		nextStartIdx := ids[insertedAtIdx+1].Start
		nextEndIdx = ids[insertedAtIdx+1].End
		if nextEndIdx == 0 {
			nextEndIdx = ids[insertedAtIdx+1].Start
		}

		if nextStartIdx-1 == id {
			needToMergeWithNext = true
		}
	}

	mergedIds := []*types.IdRange{}
	// 4 Cases: Need to merge with both, just next, just prev, or none
	if needToMergeWithPrev && needToMergeWithNext {
		mergedIds = append(mergedIds, ids[:insertedAtIdx-1]...)
		mergedIds = append(mergedIds, GetIdRangeToInsert(prevStartIdx, nextEndIdx))
		mergedIds = append(mergedIds, ids[insertedAtIdx+2:]...)
	} else if needToMergeWithPrev {
		mergedIds = append(mergedIds, ids[:insertedAtIdx-1]...)
		mergedIds = append(mergedIds, GetIdRangeToInsert(prevStartIdx, id))
		mergedIds = append(mergedIds, ids[insertedAtIdx+1:]...)
	} else if needToMergeWithNext {
		mergedIds = append(mergedIds, ids[:insertedAtIdx]...)
		mergedIds = append(mergedIds, GetIdRangeToInsert(id, nextEndIdx))
		mergedIds = append(mergedIds, ids[insertedAtIdx+2:]...)
	} else {
		mergedIds = ids
	}

	return mergedIds
}

//Inserts an id to the id ranges. Handles merging if necessary
func InsertIdRange(id uint64, ids []*types.IdRange) []*types.IdRange {
	newIds := []*types.IdRange{}
	insertIdAtIdx := 0
	if (ids[len(ids)-1].End == 0) {
		ids[len(ids)-1].End = ids[len(ids)-1].Start
	}

	if ids[0].Start > id {
		newIds = append(newIds, GetIdRangeToInsert(id, id))
		newIds = append(newIds, ids...)
	} else if ids[len(ids)-1].End < id {
		insertIdAtIdx = len(ids)
		newIds = append(newIds, ids...)
		newIds = append(newIds, GetIdRangeToInsert(id, id))
	} else {
		insertIdAtIdx = GetIdxToInsertForNewId(id, ids)
		newIds = append(newIds, ids[:insertIdAtIdx]...)
		newIds = append(newIds, GetIdRangeToInsert(id, id))
		newIds = append(newIds, ids[insertIdAtIdx:]...)
	}

	newIds = MergePrevOrNextIfNecessary(newIds, insertIdAtIdx)

	return newIds
}

func InsertRangeToIdRanges(rangeToAdd types.IdRange, ids[]*types.IdRange) []*types.IdRange {
	newIds := []*types.IdRange{}
	insertIdAtIdx := 0
	if rangeToAdd.End == 0 {
		rangeToAdd.End = rangeToAdd.Start
	}

	if (ids[len(ids)-1].End == 0) {
		ids[len(ids)-1].End = ids[len(ids)-1].Start
	}

	if ids[0].Start > rangeToAdd.End {
		newIds = append(newIds, GetIdRangeToInsert(rangeToAdd.Start, rangeToAdd.End))
		newIds = append(newIds, ids...)
	} else if ids[len(ids)-1].End < rangeToAdd.Start {
		insertIdAtIdx = len(ids)
		newIds = append(newIds, ids...)
		newIds = append(newIds, GetIdRangeToInsert(rangeToAdd.Start, rangeToAdd.End))
	} else {
		insertIdAtIdx = GetIdxToInsertForNewIdRange(rangeToAdd, ids)
		newIds = append(newIds, ids[:insertIdAtIdx]...)
		newIds = append(newIds, GetIdRangeToInsert(rangeToAdd.Start, rangeToAdd.End))
		newIds = append(newIds, ids[insertIdAtIdx:]...)
	}

	newIds = MergePrevOrNextIfNecessary(newIds, insertIdAtIdx)

	return newIds
}

// Removes an id from a single id range. Removing can make this range be split into 0, 1, or 2 new ranges.
func RemoveIdFromIdRange(id uint64, rangeObject types.IdRange) []*types.IdRange {
	newRanges := []*types.IdRange{}
	if id >= 1 && id-1 >= rangeObject.Start {
		newRanges = append(newRanges, GetIdRangeToInsert(rangeObject.Start, id-1))
	}

	if id <= math.MaxUint64-1 && id+1 <= rangeObject.End {
		newRanges = append(newRanges, GetIdRangeToInsert(id+1, rangeObject.End)) //Note rangeObject.End could == 0 but by removing the id, the range would just be removed
	}

	return newRanges
}

// Removes an id range from a single id range. Removing can make this range be split into 0, 1, or 2 new ranges.
func RemoveIdsFromIdRange(rangeToRemove types.IdRange, rangeObject types.IdRange) []*types.IdRange {
	newRanges := []*types.IdRange{}
	if rangeToRemove.End == 0 {
		rangeToRemove.End = rangeToRemove.Start
	}

	if rangeObject.End == 0 {
		rangeObject.End = rangeObject.Start
	}

	if rangeToRemove.Start > rangeObject.Start && rangeToRemove.End < rangeObject.End {
		// Completely in the middle
		newRanges = append(newRanges, GetIdRangeToInsert(rangeObject.Start, rangeToRemove.Start-1))
		newRanges = append(newRanges, GetIdRangeToInsert(rangeToRemove.End+1, rangeObject.End))
	} else if rangeToRemove.Start <= rangeObject.Start && rangeToRemove.End >= rangeObject.End {
		// Overlaps both; remove whole thing
		// Do nothing
	} else if rangeToRemove.Start <= rangeObject.Start && rangeToRemove.End < rangeObject.End {
		// Still have some at the end
		newRanges = append(newRanges, GetIdRangeToInsert(rangeToRemove.End+1, rangeObject.End))
	} else if rangeToRemove.Start > rangeObject.Start && rangeToRemove.End >= rangeObject.End {
		// Still have some at the start
		newRanges = append(newRanges, GetIdRangeToInsert(rangeObject.Start, rangeToRemove.Start-1))
	} else {
		// Doesn't overlap at all
		newRanges = append(newRanges, GetIdRangeToInsert(rangeObject.Start, rangeObject.End))
	}

	return newRanges
}


func SortIdRangesAndMergeIfNecessary(ids []*types.IdRange) []*types.IdRange {
	origIds := ids

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
	
	newIdRanges := []*types.IdRange{}
	if n > 0 {
		newIdRanges = append(newIdRanges, GetIdRangeToInsert(ids[0].Start, ids[0].End))

		for i := 1; i < n; i++ {
			prevRange := newIdRanges[i-1]
			if prevRange.End == 0 {
				prevRange.End = prevRange.Start
			}

			currRange := ids[i]
			if currRange.End == 0 {
				currRange.End = currRange.Start
			}


			if ids[i].Start == prevRange.Start {
				newIdRanges[i - 1].End = currRange.End
			} else if ids[i].End > prevRange.End {
				if ids[i].Start > prevRange.End + 1 {
					newIdRanges = append(newIdRanges, GetIdRangeToInsert(ids[i].Start, ids[i].End))
				} else {
					newIdRanges[i - 1].End = currRange.End
				}
			} else {
				newIdRanges = append(newIdRanges, GetIdRangeToInsert(ids[i].Start, ids[i].End))
			}
		}
	}
	
	for i := 0; i < len(newIdRanges); i++ {
		if newIdRanges[i].End != origIds[i].End || newIdRanges[i].Start != origIds[i].Start {
			x := 2 +2
			_ = x
		}
	}
	
	return newIdRanges
}
