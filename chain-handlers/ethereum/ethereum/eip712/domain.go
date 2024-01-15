package eip712

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// createEIP712Domain creates the typed data domain for the given chainID.
func CreateEIP712Domain(chainID uint64) apitypes.TypedDataDomain {
	domain := apitypes.TypedDataDomain{
		Name:              "BitBadges",
		Version:           "1.0.0",
		ChainId:           math.NewHexOrDecimal256(int64(chainID)),
		VerifyingContract: "0xa607FcD07cfe8d84cA839e4D6EdEE4B1A6789603",
		Salt:              "0x5d1e2c0e9b8a5c395979525d5f6d5f0c595d5a5c5e5e5b5d5ecd5a5e5d2e5412",
	}

	return domain
}
