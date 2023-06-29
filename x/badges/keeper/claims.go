package keeper

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
)

func CheckMerklePath(leaf string, expectedRoot string, aunts []*types.ClaimProofItem) error {
	hashedMsgLeaf := sha256.Sum256([]byte(leaf))
	currHash := hashedMsgLeaf[:]

	for _, aunt := range aunts {
		decodedAunt, err := hex.DecodeString(aunt.Aunt)
		if err != nil {
			return ErrDecodingHexString
		}

		if aunt.OnRight {
			parentHash := sha256.Sum256(append(currHash, decodedAunt...))
			currHash = parentHash[:]
		} else {
			parentHash := sha256.Sum256(append(decodedAunt, currHash...))
			currHash = parentHash[:]
		}
	}

	hexCurrHash := hex.EncodeToString(currHash)
	if hexCurrHash != expectedRoot {
		return ErrRootHashInvalid
	}

	return nil
}

func GetLeafIndex(aunts []*types.ClaimProofItem) sdkmath.Uint {
	leafIndex := sdkmath.NewUint(1)
	//iterate through msg.WhitelistProof.Aunts backwards
	for i := len(aunts) - 1; i >= 0; i-- {
		aunt := aunts[i]
		onRight := aunt.OnRight

		if onRight {
			leafIndex = leafIndex.Mul(sdkmath.NewUint(2))
		} else {
			leafIndex = leafIndex.Mul(sdkmath.NewUint(2)).Add(sdkmath.NewUint(1))
		}
	}
	return leafIndex
}
