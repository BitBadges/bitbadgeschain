package types

import (
	sdkmath "cosmossdk.io/math"
)

// Helper function to make code more readable
func CreateUintRange(start sdkmath.Uint, end sdkmath.Uint) *UintRange {
	return &UintRange{
		Start: start,
		End:   end,
	}
}

// Search ID ranges for a specific ID. Return (true) if found. And (false) if not.
func SearchUintRangesForUint(id sdkmath.Uint, uintRanges []*UintRange) bool {
	ranges := make([]*UintRange, len(uintRanges))
	copy(ranges, uintRanges)

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

func InvertUintRanges(uintRanges []*UintRange, maxId sdkmath.Uint) []*UintRange {
	ranges := []*UintRange{}
	ranges = append(ranges, CreateUintRange(sdkmath.NewUint(0), maxId))

	for _, uintRange := range uintRanges {
		newRanges := []*UintRange{}
		for _, rangeObject := range ranges {
			rangesAfterRemoval, _ := RemoveUintsFromUintRange(uintRange, rangeObject)
			newRanges = append(newRanges, rangesAfterRemoval...)
		}
		ranges = newRanges
	}

	return ranges
}

// Removes all ids within an id range from an id range.
// Removing can make this range be split into 0, 1, or 2 new ranges.
// Returns if anything was removed or not
func RemoveUintsFromUintRange(idxsToRemove *UintRange, rangeObject *UintRange) ([]*UintRange, []*UintRange) {
	if idxsToRemove.End.LT(rangeObject.Start) || idxsToRemove.Start.GT(rangeObject.End) {
		// idxsToRemove doesn't overlap with rangeObject, so nothing is removed
		return []*UintRange{rangeObject}, []*UintRange{}
	}

	var newRanges []*UintRange
	var removedRanges []*UintRange
	if idxsToRemove.Start.LTE(rangeObject.Start) && idxsToRemove.End.GTE(rangeObject.End) {
		// idxsToRemove fully contains rangeObject, so nothing is left
		return newRanges, []*UintRange{rangeObject}
	}

	if idxsToRemove.Start.GT(rangeObject.Start) {
		// There's a range before idxsToRemove
		// Underflow is not possible because idxsToRemove.Start.GT(rangeObject.Start
		newRanges = append(newRanges, &UintRange{
			Start: rangeObject.Start,
			End:   idxsToRemove.Start.SubUint64(1),
		})

		//get min of idxsToRemove.End and rangeObject.End
		minEnd := idxsToRemove.End
		if idxsToRemove.End.GT(rangeObject.End) {
			minEnd = rangeObject.End
		}

		removedRanges = append(removedRanges, &UintRange{
			Start: idxsToRemove.Start,
			End:   minEnd,
		})
	}

	if idxsToRemove.End.LT(rangeObject.End) {
		// There's a range after idxsToRemove
		// Overflow is not possible because idxsToRemove.End.LT(rangeObject.End
		newRanges = append(newRanges, &UintRange{
			Start: idxsToRemove.End.AddUint64(1),
			End:   rangeObject.End,
		})

		maxStart := idxsToRemove.Start
		if idxsToRemove.Start.LT(rangeObject.Start) {
			maxStart = rangeObject.Start
		}

		removedRanges = append(removedRanges, &UintRange{
			Start: maxStart,
			End:   idxsToRemove.End,
		})
	}

	return newRanges, removedRanges
}

func RemoveUintRangeFromUintRange(idsToRemove []*UintRange, rangeToRemoveFrom []*UintRange) ([]*UintRange, []*UintRange) {
	if len(idsToRemove) == 0 {
		return rangeToRemoveFrom, []*UintRange{}
	}

	removedRanges := []*UintRange{}
	for _, handledValue := range idsToRemove {
		newRanges := []*UintRange{}
		for _, oldPermittedTime := range rangeToRemoveFrom {
			rangesAfterRemoval, removed := RemoveUintsFromUintRange(handledValue, oldPermittedTime)
			newRanges = append(newRanges, rangesAfterRemoval...)
			removedRanges = append(removedRanges, removed...)
		}
		rangeToRemoveFrom = newRanges
	}

	return rangeToRemoveFrom, removedRanges
}


func AssertRangesDoNotOverlapAtAll(rangeToCheck []*UintRange, overlappingRange []*UintRange) error {
	//Check that for old times, there is 100% overlap with new times and 0% overlap with the opposite
	for _, oldAllowedTime := range rangeToCheck {
		for _, newAllowedTime := range overlappingRange {
			//Check that the new time completely overlaps with the old time
			_, removed := RemoveUintsFromUintRange(newAllowedTime, oldAllowedTime)
			if len(removed) > 0 {
				return ErrRangesOverlap
			}
		}
	}

	return nil
}


// Will sort the ID ranges in order and merge overlapping IDs if we can
func SortAndMergeOverlapping(ids []*UintRange) []*UintRange {
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
		newUintRanges := []*UintRange{CreateUintRange(ids[0].Start, ids[0].End)}
		//Iterate through and compare with previously inserted range
		for i := 1; i < n; i++ {
			prevInsertedRange := newUintRanges[len(newUintRanges)-1]
			currRange := ids[i]

			if currRange.Start.Equal(prevInsertedRange.Start) {
				//Both have same start, so we set to currRange.End because currRange.End is greater due to our sorting
				//Example: prevRange = [1, 5], currRange = [1, 10] -> newRange = [1, 10]
				newUintRanges[len(newUintRanges)-1].End = currRange.End
			} else if currRange.End.GT(prevInsertedRange.End) {
				//We have different starts and curr end is greater than prev end

				if currRange.Start.GT(prevInsertedRange.End.AddUint64(1)) {
					//We have a gap between the prev range end and curr range start, so we just append currRange
					//Example: prevRange = [1, 5], currRange = [7, 10] -> newRange = [1, 5], [7, 10]
					newUintRanges = append(newUintRanges, CreateUintRange(currRange.Start, currRange.End))
				} else {
					//They overlap and we can merge them
					//Example: prevRange = [1, 5], currRange = [2, 10] -> newRange = [1, 10]
					newUintRanges[len(newUintRanges)-1].End = currRange.End
				}
			} 
			// else {
				//Note: If currRange.End <= prevInsertedRange.End, it is already fully contained within the previous. We can just continue.
			// }
		}
		return newUintRanges
	} else {
		return ids
	}
}