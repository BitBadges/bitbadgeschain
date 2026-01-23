package cli

import (
	"encoding/json"
	"os"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

// ReadJSONFromFileOrString reads JSON from either a file path or inline string.
// If the input is a valid file path, it reads the file contents.
// Otherwise, it treats the input as inline JSON string.
func ReadJSONFromFileOrString(input string) (string, error) {
	// Check if input is a file path
	if fileInfo, err := os.Stat(input); err == nil && !fileInfo.IsDir() {
		// It's a file, read it
		contents, err := os.ReadFile(input)
		if err != nil {
			return "", err
		}
		return string(contents), nil
	}
	// Not a file, treat as inline JSON
	return input, nil
}

// ReadJSONBytesFromFileOrString reads JSON from either a file path or inline string and returns bytes.
// If the input is a valid file path, it reads the file contents.
// Otherwise, it treats the input as inline JSON string.
func ReadJSONBytesFromFileOrString(input string) ([]byte, error) {
	jsonStr, err := ReadJSONFromFileOrString(input)
	if err != nil {
		return nil, err
	}
	return []byte(jsonStr), nil
}

func GetUintRange(start types.Uint, end types.Uint) *types.UintRange {
	return &types.UintRange{
		Start: start,
		End:   end,
	}
}

func parseJsonArr(jsonStr string) ([]interface{}, error) {
	// Support file or inline JSON
	jsonBytes, err := ReadJSONBytesFromFileOrString(jsonStr)
	if err != nil {
		return nil, err
	}

	var result []interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetIdArrFromString parses JSON array of ID strings from either a file path or inline string.
// If the input is a valid file path, it reads the JSON from that file.
// Otherwise, it treats the input as inline JSON string.
func GetIdArrFromString(str string) ([]types.Uint, error) {
	vals, err := parseJsonArr(str)
	if err != nil {
		return nil, err
	}

	// convert vals to []uint64
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

// GetUintRanges parses JSON array of UintRange objects from either a file path or inline string.
// If the input is a valid file path, it reads the JSON from that file.
// Otherwise, it treats the input as inline JSON string.
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
