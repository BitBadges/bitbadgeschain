package types

import (
	sdkmath "cosmossdk.io/math"
)

// Helper function to make code more readable
func CreateIdRange(start sdkmath.Uint, end sdkmath.Uint) *IdRange {
	return &IdRange{
		Start: start,
		End:   end,
	}
}

// Search ID ranges for a specific ID. Return (true) if found. And (false) if not.
func SearchIdRangesForId(id sdkmath.Uint, idRanges []*IdRange) bool {
	ranges := make([]*IdRange, len(idRanges))
	copy(ranges, idRanges)

	ranges = SortAndMergeOverlapping(ranges) 

	//Binary search because ID ranges will be sorted
	low := 0
	high := len(ranges) - 1
	for low <= high {
		median := int(uint(low+high) >> 1)

		currRange := ranges[median]

		if currRange.Start.LTE(id) && currRange.End.GTE(id) {
			return true
		} else if currRange.Start.GT(id) {
			high = median - 1
		} else {
			low = median + 1
		}
	}

	return false
}

func InvertIdRanges(idRanges []*IdRange, maxId sdkmath.Uint) []*IdRange {
	ranges := []*IdRange{}
	ranges = append(ranges, CreateIdRange(sdkmath.NewUint(0), maxId))

	for _, idRange := range idRanges {
		newRanges := []*IdRange{}
		for _, rangeObject := range ranges {
			rangesAfterRemoval, _ := RemoveIdsFromIdRange(idRange, rangeObject)
			newRanges = append(newRanges, rangesAfterRemoval...)
		}
		ranges = newRanges
	}

	return ranges
}

// Removes all ids within an id range from an id range.
// Removing can make this range be split into 0, 1, or 2 new ranges.
// Returns if anything was removed or not
func RemoveIdsFromIdRange(idxsToRemove *IdRange, rangeObject *IdRange) ([]*IdRange, []*IdRange) {
	if idxsToRemove.End.LT(rangeObject.Start) || idxsToRemove.Start.GT(rangeObject.End) {
		// idxsToRemove doesn't overlap with rangeObject, so nothing is removed
		return []*IdRange{rangeObject}, []*IdRange{}
	}

	var newRanges []*IdRange
	var removedRanges []*IdRange
	if idxsToRemove.Start.LTE(rangeObject.Start) && idxsToRemove.End.GTE(rangeObject.End) {
		// idxsToRemove fully contains rangeObject, so nothing is left
		return newRanges, []*IdRange{rangeObject}
	}

	if idxsToRemove.Start.GT(rangeObject.Start) {
		// There's a range before idxsToRemove
		// Underflow is not possible because idxsToRemove.Start.GT(rangeObject.Start
		newRanges = append(newRanges, &IdRange{
			Start: rangeObject.Start,
			End:   idxsToRemove.Start.SubUint64(1),
		})

		//get min of idxsToRemove.End and rangeObject.End
		minEnd := idxsToRemove.End
		if idxsToRemove.End.GT(rangeObject.End) {
			minEnd = rangeObject.End
		}

		removedRanges = append(removedRanges, &IdRange{
			Start: idxsToRemove.Start,
			End:   minEnd,
		})
	}

	if idxsToRemove.End.LT(rangeObject.End) {
		// There's a range after idxsToRemove
		// Overflow is not possible because idxsToRemove.End.LT(rangeObject.End
		newRanges = append(newRanges, &IdRange{
			Start: idxsToRemove.End.AddUint64(1),
			End:   rangeObject.End,
		})

		maxStart := idxsToRemove.Start
		if idxsToRemove.Start.LT(rangeObject.Start) {
			maxStart = rangeObject.Start
		}

		removedRanges = append(removedRanges, &IdRange{
			Start: maxStart,
			End:   idxsToRemove.End,
		})
	}

	return newRanges, removedRanges
}

func RemoveIdRangeFromIdRange(idsToRemove []*IdRange, rangeToRemoveFrom []*IdRange) ([]*IdRange, []*IdRange) {
	if len(idsToRemove) == 0 {
		return rangeToRemoveFrom, []*IdRange{}
	}

	removedRanges := []*IdRange{}
	for _, handledValue := range idsToRemove {
		newRanges := []*IdRange{}
		for _, oldPermittedTime := range rangeToRemoveFrom {
			rangesAfterRemoval, removed := RemoveIdsFromIdRange(handledValue, oldPermittedTime)
			newRanges = append(newRanges, rangesAfterRemoval...)
			removedRanges = append(removedRanges, removed...)
		}
		rangeToRemoveFrom = newRanges
	}

	return rangeToRemoveFrom, removedRanges
}


func AssertRangesDoNotOverlapAtAll(rangeToCheck []*IdRange, overlappingRange []*IdRange) error {
	//Check that for old times, there is 100% overlap with new times and 0% overlap with the opposite
	for _, oldAllowedTime := range rangeToCheck {
		for _, newAllowedTime := range overlappingRange {
			//Check that the new time completely overlaps with the old time
			_, removed := RemoveIdsFromIdRange(newAllowedTime, oldAllowedTime)
			if len(removed) > 0 {
				return ErrInvalidPermissionsUpdateLocked
			}
		}
	}

	return nil
}


// Will sort the ID ranges in order and merge overlapping IDs if we can
func SortAndMergeOverlapping(ids []*IdRange) []*IdRange {
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
		newIdRanges := []*IdRange{CreateIdRange(ids[0].Start, ids[0].End)}
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
			} 
			// else {
				//Note: If currRange.End <= prevInsertedRange.End, it is already fully contained within the previous. We can just continue.
			// }
		}
		return newIdRanges
	} else {
		return ids
	}
}