package types

import (
	"math/big"
	"strings"

	errorsmod "cosmossdk.io/errors"
)

func ParseChainID(chainID string) (*big.Int, error) {
	chainID = strings.TrimSpace(chainID)

	if chainID == "bitbadges-1" {
		chainIDInt, ok := new(big.Int).SetString("1", 10)
		if !ok {
			return nil, errorsmod.Wrapf(ErrInvalidChainID, "")
		}

		return chainIDInt, nil
	} else if chainID == "bitbadges-2" {
		chainIDInt, ok := new(big.Int).SetString("2", 10)
		if !ok {
			return nil, errorsmod.Wrapf(ErrInvalidChainID, "")
		}

		return chainIDInt, nil
	}

	return nil, errorsmod.Wrapf(ErrInvalidChainID, "")
}
