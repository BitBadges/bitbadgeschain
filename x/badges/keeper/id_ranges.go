package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func CreateIdRange(start uint64, end uint64) *types.IdRange {
	return &types.IdRange{
		Start: start,
		End:   end,
	}
}

// Search ID ranges for a specific ID. Return (idx, true) if found. And (-1, false) if not. 
// Assumes ID ranges are sorted.
func SearchIdRangesForId(id uint64, idRanges []*types.IdRange) (int, bool) {
	//Binary search because ID ranges will be sorted
	low := 0
	high := len(idRanges) - 1
	for low <= high {
		median := int(uint(low+high) >> 1)

		currRange := idRanges[median]

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

// Search a set of ranges to find what indexes a specific ID range overlaps. 
// Return overlapping idxs as a IdRange, true if found. And empty IdRange, false if not
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
		endIdx-- //We want the idx before because we want the last overlapping range
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

// Assumes given ID is not already in a range. 
// Gets the index to insert at. Ex. [{0-10}, {10-20}, {30-40}] and inserting 25 would return index 2
func GetIdxToInsertForNewId(id uint64, targetIds []*types.IdRange) (int, error) {
	_, found := SearchIdRangesForId(id, targetIds) 
	if found {
		return -1, ErrIdInRange
	}
	
	//Since we assume the id is not already in there, we can just compare start positions of the existing idRanges
	ids := targetIds
	if len(ids) == 0 {
		return 0, nil
	}

	if ids[0].Start > id { //assumes not in already so we don't have to handle that case
		return 0, nil
	} else if ids[len(ids)-1].End < id {
		return len(ids), nil
	}

	//Binary search by looking at two ranges at a time [..., {curr}, {next}, ...]
	low := 0
	high := len(ids) - 2
	median := 0
	for low <= high {
		median = int(uint(low+high) >> 1)
		currRange := ids[median]
		nextRange := ids[median+1]

		if currRange.Start < id && nextRange.Start > id {
			break
		} else if currRange.Start > id {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	//We return median + 1 to insert in between {curr} and {next}
	return median + 1, nil
}

// Inserts a range into its correct position. 
// Assumes whole range is not present at all. Thus, we only search for where start fits in.
func InsertRangeToIdRanges(rangeToAdd *types.IdRange, targetIds []*types.IdRange) ([]*types.IdRange, error) {
	//Validation check; make sure rangeToAdd is not already in targetIds
	for _, id := range targetIds {
		_, removed := RemoveIdsFromIdRange(id, rangeToAdd)
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
	if ids[0].Start > rangeToAdd.End {
		newIds = append(newIds, CreateIdRange(rangeToAdd.Start, rangeToAdd.End))
		newIds = append(newIds, ids...)
	} else if lastRange.End < rangeToAdd.Start {
		insertIdAtIdx = len(ids)
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
	newRanges := []*types.IdRange{}
	removed := true

	//5 cases: Completely before, completely after, completely in the middle, overlaps both, overlaps one side
	if idxsToRemove.Start > rangeObject.Start && idxsToRemove.End < rangeObject.End {
		// Completely in the middle; Split into two ranges
		newRanges = append(newRanges, CreateIdRange(rangeObject.Start, idxsToRemove.Start-1))
		newRanges = append(newRanges, CreateIdRange(idxsToRemove.End+1, rangeObject.End))
	} else if idxsToRemove.Start <= rangeObject.Start && idxsToRemove.End >= rangeObject.End {
		// Overlaps both; remove whole thing
		// Do nothing
	} else if idxsToRemove.Start <= rangeObject.Start && idxsToRemove.End < rangeObject.End && idxsToRemove.End >= rangeObject.Start {
		// Still have some left at the end
		newRanges = append(newRanges, CreateIdRange(idxsToRemove.End+1, rangeObject.End))
	} else if idxsToRemove.Start > rangeObject.Start && idxsToRemove.End >= rangeObject.End && idxsToRemove.Start <= rangeObject.End {
		// Still have some left at the start
		newRanges = append(newRanges, CreateIdRange(rangeObject.Start, idxsToRemove.Start-1))
	} else {
		// Doesn't overlap at all; keep everything
		newRanges = append(newRanges, CreateIdRange(rangeObject.Start, rangeObject.End))
		removed = false //Didn't remove anything
	}

	return newRanges, removed
}

//Will sort the ID ranges in order and merge overlapping IDs if we can
func SortAndMergeOverlapping(ids []*types.IdRange) []*types.IdRange {
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
		newIdRanges = append(newIdRanges, CreateIdRange(ids[0].Start, ids[0].End))
		//Iterate through and compare with previously inserted range
		for i := 1; i < n; i++ {
			prevInsertedRange := newIdRanges[len(newIdRanges)-1]
			currRange := ids[i]

			if currRange.Start == prevInsertedRange.Start {
				//Both have same start, so we set to currRange.End because currRange.End is greater due to our sorting
				newIdRanges[len(newIdRanges)-1].End = currRange.End
			} else if currRange.End > prevInsertedRange.End {
				//We have different starts and curr end is greater than prev end
				if currRange.Start > prevInsertedRange.End+1 {
					//We have a gap between the prev range end and curr range start, so we just append currRange
					newIdRanges = append(newIdRanges, CreateIdRange(currRange.Start, currRange.End))
				} else {
					//they overlap and we can merge them
					newIdRanges[len(newIdRanges)-1].End = currRange.End
				}
			}
			///If currRange.End <= prevInsertedRange.End, it is already fully contained within the previous. We can just continue.
		}
		return newIdRanges
	} else {
		return ids
	}
}


func AddManagerAddressToRanges(collection types.BadgeCollection, ranges []*types.IdRange, options types.AddressOptions) []*types.IdRange {
	idx, found := SearchIdRangesForId(collection.Manager, ranges)
	//Add or remove the manager to the ranges if specified according to the options
	if options == types.AddressOptions_IncludeManager {
		if !found {
			ranges = append(ranges, CreateIdRange(collection.Manager, collection.Manager))
			ranges = SortAndMergeOverlapping(ranges)
			return ranges
		}
	} else if options == types.AddressOptions_ExcludeManager {
		if found {
			newRanges := []*types.IdRange{}
			newRanges = append(newRanges, ranges[:idx]...)
			removedRanges, _ := RemoveIdsFromIdRange(CreateIdRange(collection.Manager, collection.Manager), ranges[idx])
			newRanges = append(newRanges, removedRanges...)
			newRanges = append(newRanges, ranges[idx+1:]...)
			ranges = newRanges
		}
	}

	return ranges
}