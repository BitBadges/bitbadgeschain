package keeper

import (
	"crypto/sha256"
	"encoding/hex"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AssertValidSolutionForEveryChallenge(ctx sdk.Context, collectionId sdkmath.Uint, challenges []*types.Challenge, solutions []*types.ChallengeSolution, creatorAddress string, level string) (bool, sdkmath.Uint, error) {
	numIncrements := sdkmath.NewUint(0)
	useLeafIndexForDistributionOrder := false

	for _, challenge := range challenges {
		root := challenge.Root
		hasValidSolution := false

		for _, solution := range solutions {
			if root != "" {
				// Must be proper length to avoid preimage attacks
				if len(solution.Proof.Aunts) != int(challenge.ExpectedProofLength.Uint64()) {
					continue
				}

				if challenge.UseCreatorAddressAsLeaf {
					solution.Proof.Leaf = creatorAddress //overwrites it
				}

				if solution.Proof.Leaf == "" {
					continue
				}

				leafIndex := GetLeafIndex(solution.Proof.Aunts)
				if challenge.UseLeafIndexForDistributionOrder {
					useLeafIndexForDistributionOrder = true

					//Get leftmost leaf index for layer === challenge.ExpectedProofLength
					leftmostLeafIndex := sdkmath.NewUint(1)
					for i := sdkmath.NewUint(0); i.LT(challenge.ExpectedProofLength); i = i.Add(sdkmath.NewUint(1)) {
						leftmostLeafIndex = leftmostLeafIndex.Mul(sdkmath.NewUint(2))
					}

					numIncrements = leafIndex.Sub(leftmostLeafIndex)
				}

				if challenge.MaxOneUsePerLeaf {
					numUsed, err := k.IncrementNumUsedForChallengeInStore(ctx, collectionId, challenge.ChallengeId, leafIndex, level)
					if err != nil {
						continue
					}

					maxUses := sdkmath.NewUint(1)
					if numUsed.GT(maxUses) {
						continue
					}
				}

				err := CheckMerklePath(solution.Proof.Leaf, root, solution.Proof.Aunts)
				if err != nil {
					continue
				}

				hasValidSolution = true
			}
		}

		if !hasValidSolution {
			return false, numIncrements, ErrNoValidSolutionForChallenge
		}
	}

	return useLeafIndexForDistributionOrder, numIncrements, nil
}

func CheckMerklePath(leaf string, expectedRoot string, aunts []*types.ClaimProofItem) error {
	hashedMsgLeaf := sha256.Sum256([]byte(leaf))
	currHash := hashedMsgLeaf[:]

	for _, aunt := range aunts {
		decodedAunt, err := hex.DecodeString(aunt.Aunt)
		if err != nil {
			return sdkerrors.Wrapf(ErrDecodingHexString, "error decoding aunt %s", aunt.Aunt)
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
		return sdkerrors.Wrapf(ErrRootHashInvalid, "expected root %s, got %s", expectedRoot, hexCurrHash)
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
