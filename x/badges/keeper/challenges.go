package keeper

import (
	"fmt"
	"crypto/sha256"
	"encoding/hex"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AssertValidSolutionForEveryChallenge(ctx sdk.Context, collectionId sdkmath.Uint, challenges []*types.MerkleChallenge, merkleProofs []*types.MerkleProof, creatorAddress string, simulation bool, approverAddress string, challengeLevel string) (sdkmath.Uint, error) {
	numIncrements := sdkmath.NewUint(0)

	for _, challenge := range challenges {
		root := challenge.Root
		hasValidSolution := false

		for _, proof := range merkleProofs {
			if root != "" {
				// Must be proper length to avoid preimage attacks
				if len(proof.Aunts) != int(challenge.ExpectedProofLength.Uint64()) {
					continue
				}

				if challenge.UseCreatorAddressAsLeaf {
					proof.Leaf = creatorAddress //overwrites it
				}

				if proof.Leaf == "" {
					continue
				}

				leafIndex := GetLeafIndex(proof.Aunts)
				if challenge.UseLeafIndexForTransferOrder {
					//Get leftmost leaf index for layer === challenge.ExpectedProofLength
					leftmostLeafIndex := sdkmath.NewUint(1)
					for i := sdkmath.NewUint(0); i.LT(challenge.ExpectedProofLength); i = i.Add(sdkmath.NewUint(1)) {
						leftmostLeafIndex = leftmostLeafIndex.Mul(sdkmath.NewUint(2))
					}

					numIncrements = leafIndex.Sub(leftmostLeafIndex)
				}

				err := CheckMerklePath(proof.Leaf, root, proof.Aunts)
				if err != nil {
					continue
				}

				
				if challenge.MaxOneUsePerLeaf {
					challengeId := challenge.ChallengeId
					numUsed, err := k.GetNumUsedForMerkleChallengeFromStore(ctx, collectionId, approverAddress, challengeLevel, challengeId, leafIndex)
					if err != nil {
						continue
					}
					numUsed = numUsed.Add(sdkmath.NewUint(1))

					maxUses := sdkmath.NewUint(1)
					if numUsed.GT(maxUses) {
						continue
					}

					if !simulation {
						newNumUsed, err := k.IncrementNumUsedForMerkleChallengeInStore(ctx, collectionId, approverAddress, challengeLevel, challengeId, leafIndex)
						if err != nil {
							continue
						}

						//Currently added for indexer, but note that it is planned to be deprecated
						ctx.EventManager().EmitEvent(
							sdk.NewEvent("challenge",
								sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
								sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
								sdk.NewAttribute("challengeId", fmt.Sprint(challengeId)),
								sdk.NewAttribute("leafIndex", fmt.Sprint(leafIndex)),
								sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
								sdk.NewAttribute("challengeLevel", fmt.Sprint(challengeLevel)),
								sdk.NewAttribute("numUsed", fmt.Sprint(newNumUsed)),
							),
						)
					}


				}

				hasValidSolution = true
				break
			}
		}

		if !hasValidSolution {
			return numIncrements, ErrNoValidSolutionForChallenge
		}
	}

	return numIncrements, nil
}

func CheckMerklePath(leaf string, expectedRoot string, aunts []*types.MerklePathItem) error {
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

func GetLeafIndex(aunts []*types.MerklePathItem) sdkmath.Uint {
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
