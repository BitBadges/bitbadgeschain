// Copyright 2023 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/bitbadges/bitbadgeschain/x/ethermint/blob/main/LICENSE
package eip712

import (
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// createEIP712Domain creates the typed data domain for the given chainID.
func createEIP712Domain(chainID uint64) apitypes.TypedDataDomain {
	domain := apitypes.TypedDataDomain{
		Name:              "BitBadges",
		Version:           "1.0.0",
		ChainId:           math.NewHexOrDecimal256(int64(chainID)),
		VerifyingContract: "0x1a16c87927570239fecd343ad2654fd81682725e",
		Salt:              "0x5d1e2c0e9b8a5c395979525d5f6d5f0c595d5a5c5e5e5b5d5ecd5a5e5d2e5412",
	}

	return domain
}