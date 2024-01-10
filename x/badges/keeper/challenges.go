package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AssertValidSolutionForEveryChallenge(ctx sdk.Context, collectionId sdkmath.Uint, challengeId string, challenges []*types.MerkleChallenge, merkleProofs []*types.MerkleProof, creatorAddress string, simulation bool, approverAddress string, challengeLevel string, challengeIdsIncremented *[]string, approval *types.CollectionApproval) (sdkmath.Uint, error) {
	numIncrements := sdkmath.NewUint(0)

	for _, challenge := range challenges {
		if challenge == nil || challenge.Root == "" {
			//No challenge specified
			continue
		}

		root := challenge.Root
		hasValidSolution := false
		errStr := ""
		if challenge.UseCreatorAddressAsLeaf {
			errStr = "does not satisfy allowlist"
		} else {
			errStr = "invalid code / password"
		}

		additionalDetailsErrorStr := ""

		for _, proof := range merkleProofs {
			additionalDetailsErrorStr = ""
			if root != "" {
				// Must be proper length to avoid preimage attacks
				if len(proof.Aunts) != int(challenge.ExpectedProofLength.Uint64()) {
					additionalDetailsErrorStr = "invalid proof length"
					continue
				}

				if challenge.UseCreatorAddressAsLeaf {
					proof.Leaf = creatorAddress //overwrites it
				}

				if proof.Leaf == "" {
					additionalDetailsErrorStr = "empty leaf"
					continue
				}

				leafIndex := GetLeafIndex(proof.Aunts)

				//Get leftmost leaf index for layer === challenge.ExpectedProofLength
				leftmostLeafIndex := sdkmath.NewUint(1)
				for i := sdkmath.NewUint(0); i.LT(challenge.ExpectedProofLength); i = i.Add(sdkmath.NewUint(1)) {
					leftmostLeafIndex = leftmostLeafIndex.Mul(sdkmath.NewUint(2))
				}

				useLeafIndexForTransferOrder := false
				if approval.ApprovalCriteria != nil && approval.ApprovalCriteria.PredeterminedBalances != nil && approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod != nil && approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {
					useLeafIndexForTransferOrder = true
				}

				if useLeafIndexForTransferOrder {
					numIncrements = leafIndex.Sub(leftmostLeafIndex)
				}

				err := CheckMerklePath(proof.Leaf, root, proof.Aunts)
				if err != nil {
					additionalDetailsErrorStr = ""
					continue
				}

				newNumUsed := sdk.NewUint(0)
				if !challenge.MaxUsesPerLeaf.IsNil() && challenge.MaxUsesPerLeaf.GT(sdkmath.NewUint(0)) {

					numUsed, err := k.GetChallengeTrackerFromStore(ctx, collectionId, approverAddress, challengeLevel, challengeId, leafIndex)
					if err != nil {
						additionalDetailsErrorStr = "error getting num processed"
						continue
					}
					numUsed = numUsed.Add(sdkmath.NewUint(1))

					maxUses := challenge.MaxUsesPerLeaf
					if numUsed.GT(maxUses) {
						additionalDetailsErrorStr = "exceeded max number of uses"
						continue
					}

					if !simulation {
						incrementId := fmt.Sprint(collectionId) + "-" + fmt.Sprint(approverAddress) + "-" + fmt.Sprint(challengeLevel) + "-" + fmt.Sprint(challengeId) + "-" + fmt.Sprint(leafIndex)
						alreadyIncremented := false
						for _, id := range *challengeIdsIncremented {
							if id == incrementId {
								alreadyIncremented = true
								break
							}
						}

						if !alreadyIncremented {
							*challengeIdsIncremented = append(*challengeIdsIncremented, incrementId)
							newNumUsed, err = k.IncrementChallengeTrackerInStore(ctx, collectionId, approverAddress, challengeLevel, challengeId, leafIndex.Sub(leftmostLeafIndex))
							if err != nil {
								continue
							}

							//Currently added for indexer, but note that it is planned to be deprecated
							ctx.EventManager().EmitEvent(
								sdk.NewEvent("challenge"+fmt.Sprint(challengeId)+fmt.Sprint(challengeId)+fmt.Sprint(leafIndex)+fmt.Sprint(approverAddress)+fmt.Sprint(challengeLevel)+fmt.Sprint(newNumUsed),
									sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
									sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
									sdk.NewAttribute("challengeId", fmt.Sprint(challengeId)),
									sdk.NewAttribute("leafIndex", fmt.Sprint(leafIndex.Sub(leftmostLeafIndex))),
									sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
									sdk.NewAttribute("challengeLevel", fmt.Sprint(challengeLevel)),
									sdk.NewAttribute("numUsed", fmt.Sprint(newNumUsed)),
								),
							)
						}
					}
				}

				hasValidSolution = true
				break
			}
		}

		if !hasValidSolution {
			return numIncrements, sdkerrors.Wrapf(ErrNoValidSolutionForChallenge, "%s - %s", errStr, additionalDetailsErrorStr)
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
	//iterate through msg.AllowlistProof.Aunts backwards
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
