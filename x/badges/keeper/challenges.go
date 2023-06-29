package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmath "cosmossdk.io/math"
)

func (k Keeper) AssertValidSolutionForEveryChallenge(ctx sdk.Context, collectionId sdkmath.Uint,  challenges []*types.Challenge, solutions []*types.ChallengeSolution, creatorAddress string, level string) (bool, sdkmath.Uint, error) {
	numIncrements := sdk.NewUint(0)

	for _, challenge := range challenges {
		root := challenge.Root
		useCreatorAddressAsLeaf := challenge.UseCreatorAddressAsLeaf
		expectedProofLength := challenge.ExpectedProofLength
		hasValidSolution := false
		challengeId := challenge.ChallengeId
		for _, solution := range solutions {
			

			if root != "" {
				if len(solution.Proof.Aunts) != int(expectedProofLength.Uint64()) {
					continue
				}

				if useCreatorAddressAsLeaf {
					solution.Proof.Leaf = creatorAddress //overwrites it
				}

				leafIndex := GetLeafIndex(solution.Proof.Aunts)
				if challenge.UseLeafIndexForDistributionOrder {
					//Get leftmost leaf index for layer === expectedProofLength
					leftmostLeafIndex := sdk.NewUint(1)
					for i := sdk.NewUint(0); i.LT(expectedProofLength); i = i.Add(sdk.NewUint(1)) {
						leftmostLeafIndex = leftmostLeafIndex.Mul(sdk.NewUint(2))
					}

					numIncrements = leafIndex.Sub(leftmostLeafIndex)
				}

				if challenge.MaxOneUsePerLeaf {
					numUsed, err := k.IncrementNumUsedForChallengeInStore(ctx, collectionId, challengeId, leafIndex, level)
					if err != nil {
						continue
					}

					maxUses := sdk.NewUint(1)
					if numUsed.GT(maxUses) {
						continue
					}
				}

				//Check if claim is valid
				if solution.Proof.Leaf == "" {
					continue
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

	return true, numIncrements, nil
}