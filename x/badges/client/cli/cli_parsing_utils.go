package cli

import (
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/spf13/cast"
)

func GetIdRange(start uint64, end uint64) *types.IdRange {
	return &types.IdRange{
		Start: start,
		End:   end,
	}
}

func parseJson(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseJsonArr(jsonStr string) ([]interface{}, error) {
	var result []interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetIdArrFromString(str string) ([]uint64, error) {
	vals, err := parseJsonArr(str)
	if err != nil {
		return nil, err
	}

	//convert vals to []uint64
	argStartValuesUint64 := []uint64{}
	for _, val := range vals {
		valAsUint64, err := cast.ToUint64E(val)
		if err != nil {
			return nil, err
		}

		argStartValuesUint64 = append(argStartValuesUint64, valAsUint64)
	}

	return argStartValuesUint64, nil
}

// Start and end strings should be comma separated list of ids
func GetIdRanges(idRangesStr string) ([]*types.IdRange, error) {
	vals, err := parseJsonArr(idRangesStr)
	if err != nil {
		return nil, err
	}

	ranges := []*types.IdRange{}
	for _, val := range vals {
		valAsMap, ok := val.(types.IdRange)
		if !ok {
			return nil, types.ErrInvalidIdRangeSpecified
		}

		ranges = append(ranges, &valAsMap)
	}

	return ranges, nil
}
