package cli

import (
	"strings"

	"github.com/spf13/cast"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func GetUriObject(uri string, subassetUri string) (*types.UriObject, error) {
	//TODO: get uri object from uri and subasset uri
	return &types.UriObject{
		Uri:         []byte(uri),
	}, nil
}

func GetIdRange(start uint64, end uint64) *types.IdRange {
	return &types.IdRange{
		Start: start,
		End:   end,
	}
}

func GetIdArrFromString(str string) ([]uint64, error) {
	argStartingNonces := strings.Split(str, listSeparator)

	argStartingNoncesUint64 := []uint64{}
	for _, nonce := range argStartingNonces {
		nonceAsUint64, err := cast.ToUint64E(nonce)
		if err != nil {
			return nil, err
		}

		argStartingNoncesUint64 = append(argStartingNoncesUint64, nonceAsUint64)
	}

	return argStartingNoncesUint64, nil
}
 
func GetIdRanges(startStrs string, endStrs string) ([]*types.IdRange, error) {
	argStartingNoncesUint64, err := GetIdArrFromString(startStrs)
	if err != nil {
		return nil, err
	}

	argEndingNoncesUint64, err := GetIdArrFromString(endStrs)
	if err != nil {
		return nil, err
	}

	if len(argStartingNoncesUint64) != len(argEndingNoncesUint64) {
		return nil, types.ErrInvalidArgumentLengths
	}

	nonceRanges := []*types.IdRange{}
	for i := 0; i < len(argStartingNoncesUint64); i++ {
		nonceRanges = append(nonceRanges, &types.IdRange{
			Start: argStartingNoncesUint64[i],
			End:   argEndingNoncesUint64[i],
		})
	}

	return nonceRanges, nil
}
