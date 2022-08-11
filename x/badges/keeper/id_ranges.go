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
		median := int(uint(low + high) >> 1)
		
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
func GetIdxToInsertForNewId(id uint64, targetIds []*types.IdRange) (int) {
	low := 0
	high := len(targetIds) - 2
	median := 0
	for low <= high {
		median = int(uint(low + high) >> 1)
		if targetIds[median].Start < id && targetIds[median + 1].Start > id {
			break;
		} else if targetIds[median].Start > id {
			high = median - 1
		} else {
			low = median + 1
		}
	}
	
	insertIdx := median + 1
	if (targetIds[median].Start <= id) {
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
		prevStartIdx = ids[insertedAtIdx - 1].Start
		prevEndIdx := ids[insertedAtIdx - 1].End
		if prevEndIdx == 0 {
			prevEndIdx = ids[insertedAtIdx - 1].Start
		}

		if prevEndIdx + 1 == id {
			needToMergeWithPrev = true
		}
	}

	if insertedAtIdx < len(ids) - 2 {
		nextStartIdx := ids[insertedAtIdx + 1].Start
		nextEndIdx = ids[insertedAtIdx + 1].End
		if nextEndIdx == 0 {
			nextEndIdx = ids[insertedAtIdx + 1].Start
		}

		if nextStartIdx - 1 == id {
			needToMergeWithNext = true
		}
	}


	mergedIds := []*types.IdRange{}
	// 4 Cases: Need to merge with both, just next, just prev, or none
	if needToMergeWithPrev && needToMergeWithNext {
		mergedIds = append(mergedIds, ids[:insertedAtIdx - 1]...)
		mergedIds = append(mergedIds, GetIdRangeToInsert(prevStartIdx, nextEndIdx))
		mergedIds = append(mergedIds, ids[insertedAtIdx + 2:]...)
	} else if needToMergeWithPrev {
		mergedIds = append(mergedIds, ids[:insertedAtIdx - 1]...)
		mergedIds = append(mergedIds, GetIdRangeToInsert(prevStartIdx, id))
		mergedIds = append(mergedIds, ids[insertedAtIdx + 1:]...)
	} else if needToMergeWithNext {
		mergedIds = append(mergedIds, ids[:insertedAtIdx]...)
		mergedIds = append(mergedIds, GetIdRangeToInsert(id, nextEndIdx))
		mergedIds = append(mergedIds, ids[insertedAtIdx + 2:]...)
	} else {
		mergedIds = ids
	}

	return mergedIds
}

//Inserts an id to the id ranges. Handles merging if necessary
func InsertIdRange(id uint64, ids []*types.IdRange) ([]*types.IdRange) {
	newIds := []*types.IdRange{}
	insertIdAtIdx := 0
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

// Removes an id from a single id range. Removing can make this range be split into 0, 1, or 2 new ranges. 
func RemoveIdFromIdRange(id uint64, rangeObject types.IdRange) []*types.IdRange {
	newRanges := []*types.IdRange{}
	if id >= 1 && id - 1 >= rangeObject.Start {
		newRanges = append(newRanges, GetIdRangeToInsert(rangeObject.Start, id - 1))
	}

	if id <= math.MaxUint64-1 && id + 1 <= rangeObject.End {
		newRanges = append(newRanges, GetIdRangeToInsert(id + 1, rangeObject.End)) //Note rangeObject.End could == 0 but by removing the id, the range would just be removed
	}
	
	return newRanges
}