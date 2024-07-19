package cli

import (
	"encoding/json"

	"bitbadgeschain/x/badges/types"
)

func GetUintRange(start types.Uint, end types.Uint) *types.UintRange {
	return &types.UintRange{
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

func GetIdArrFromString(str string) ([]types.Uint, error) {
	vals, err := parseJsonArr(str)
	if err != nil {
		return nil, err
	}

	//convert vals to []uint64
	argStartValuesUint64 := []types.Uint{}
	for _, val := range vals {
		valAsUint64 := types.NewUintFromString(val.(string))
		if err != nil {
			return nil, err
		}

		argStartValuesUint64 = append(argStartValuesUint64, valAsUint64)
	}

	return argStartValuesUint64, nil
}

// Start and end strings should be comma separated list of ids
func GetUintRanges(uintRangesStr string) ([]*types.UintRange, error) {
	vals, err := parseJsonArr(uintRangesStr)
	if err != nil {
		return nil, err
	}

	ranges := []*types.UintRange{}
	for _, val := range vals {
		valAsMap, ok := val.(types.UintRange)
		if !ok {
			return nil, types.ErrInvalidUintRangeSpecified
		}

		ranges = append(ranges, &valAsMap)
	}

	return ranges, nil
}
