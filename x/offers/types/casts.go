package types

import (
	badgetypes "bitbadgeschain/x/badges/types"
)

func CastUintRanges(ranges []*UintRange) []*badgetypes.UintRange {
	castedRanges := make([]*badgetypes.UintRange, len(ranges))
	for i, rangeVal := range ranges {
		castedRanges[i] = &badgetypes.UintRange{
			Start: rangeVal.Start,
			End:   rangeVal.End,
		}
	}
	return castedRanges
}
