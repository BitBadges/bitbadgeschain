package cli

import (
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/spf13/cast"
)

func GetIdRange(start uint64, end uint64) *types.IdRange {
	return &types.IdRange{
		Start: start,
		End:   end,
	}
}

func GetIdArrFromString(str string) ([]uint64, error) {
	argStartValues := strings.Split(str, listSeparator)

	argStartValuesUint64 := []uint64{}
	for _, val := range argStartValues {
		valAsUint64, err := cast.ToUint64E(val)
		if err != nil {
			return nil, err
		}

		argStartValuesUint64 = append(argStartValuesUint64, valAsUint64)
	}

	return argStartValuesUint64, nil
}

// Start and end strings should be comma separated list of ids
func GetIdRanges(startStr string, endStr string) ([]*types.IdRange, error) {
	argStartValuesUint64, err := GetIdArrFromString(startStr)
	if err != nil {
		return nil, err
	}

	argEndingValuesUint64, err := GetIdArrFromString(endStr)
	if err != nil {
		return nil, err
	}

	if len(argStartValuesUint64) != len(argEndingValuesUint64) {
		return nil, types.ErrInvalidArgumentLengths
	}

	ranges := []*types.IdRange{}
	for i := 0; i < len(argStartValuesUint64); i++ {
		ranges = append(ranges, &types.IdRange{
			Start: argStartValuesUint64[i],
			End:   argEndingValuesUint64[i],
		})
	}

	return ranges, nil
}
