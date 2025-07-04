package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/storyicon/sigverify"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func (k Keeper) HandleMerkleChallenges(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	transfer *types.Transfer,
	approval *types.CollectionApproval,
	creatorAddress string,
	approverAddress string,
	approvalLevel string,
	simulation bool,
) (sdkmath.Uint, error) {
	numIncrements := sdkmath.NewUint(0)
	challenges := approval.ApprovalCriteria.MerkleChallenges
	merkleProofs := transfer.MerkleProofs

	// Sanity check to make sure the challenge tracker id is valid
	if approval.ApprovalCriteria != nil && approval.ApprovalCriteria.PredeterminedBalances != nil && approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId != "" && approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {
		hasMatchingChallenge := false
		for _, challenge := range challenges {
			if challenge.ChallengeTrackerId == approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId {
				hasMatchingChallenge = true
				break
			}
		}

		if !hasMatchingChallenge {
			return numIncrements, sdkerrors.Wrapf(ErrNoMatchingChallengeForChallengeTrackerId, "expected to calculate balances from challenge but no matching challenge for challenge tracker id %s", approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId)
		}
	}

	for _, challenge := range challenges {
		if challenge == nil || challenge.Root == "" {
			return numIncrements, sdkerrors.Wrapf(types.ErrChallengeTrackerIdIsNil, "challenge is nil or has empty root")
		}

		challengeId := challenge.ChallengeTrackerId
		root := challenge.Root
		hasValidSolution := false

		errStr := ""
		if challenge.UseCreatorAddressAsLeaf {
			errStr = "does not satisfy whitelist"
		} else {
			errStr = "invalid code / password"
		}

		// We check that 1 of N proofs is valid
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
					proof.Leaf = creatorAddress //overwrites it with creator address
				}

				if proof.Leaf == "" {
					additionalDetailsErrorStr = "empty leaf"
					continue
				}

				if challenge.LeafSigner != "" {
					leafSignature := proof.LeafSignature

					leafSignerEthAddress := challenge.LeafSigner
					if leafSignerEthAddress == "" {
						additionalDetailsErrorStr = "empty leaf signer"
						continue
					}

					ethAddress := ethcommon.HexToAddress(leafSignerEthAddress)

					leafSignatureString := proof.Leaf + "-" + creatorAddress
					isValid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
						ethAddress,
						[]byte(leafSignatureString),
						leafSignature,
					)

					if !isValid || err != nil {
						additionalDetailsErrorStr = "invalid leaf signature"
						continue
					}
				}

				//Get leftmost leaf index for layer === challenge.ExpectedProofLength
				leafIndex := GetLeafIndex(proof.Aunts)
				leftmostLeafIndex := sdkmath.NewUint(1)
				for i := sdkmath.NewUint(0); i.LT(challenge.ExpectedProofLength); i = i.Add(sdkmath.NewUint(1)) {
					leftmostLeafIndex = leftmostLeafIndex.Mul(sdkmath.NewUint(2))
				}

				//Predefined balances challenge tracker = current challenge tracker
				useLeafIndexForTransferOrder := false
				if approval.ApprovalCriteria != nil && approval.ApprovalCriteria.PredeterminedBalances != nil && approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod != nil && approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {
					if approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId == challengeId {
						useLeafIndexForTransferOrder = true
					}
				}

				if useLeafIndexForTransferOrder {
					numIncrements = leafIndex.Sub(leftmostLeafIndex)
				}

				err := CheckMerklePath(proof.Leaf, root, proof.Aunts)
				if err != nil {
					additionalDetailsErrorStr = ""
					continue
				}

				//If there is a max uses per leaf, we need to check it has not exceeded the treshold uses
				if !challenge.MaxUsesPerLeaf.IsNil() && challenge.MaxUsesPerLeaf.GT(sdkmath.NewUint(0)) {
					numUsed, err := k.GetChallengeTrackerFromStore(ctx, collectionId, approverAddress, approvalLevel, approval.ApprovalId, challengeId, leafIndex.Sub(leftmostLeafIndex))
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

					//Increment the number of uses in store if we are doing it for real
					if !simulation {
						newNumUsed, err := k.IncrementChallengeTrackerInStore(ctx, collectionId, approverAddress, approvalLevel, approval.ApprovalId, challengeId, leafIndex.Sub(leftmostLeafIndex))
						if err != nil {
							continue
						}

						//Currently added for indexer, but note that it is planned to be deprecated
						ctx.EventManager().EmitEvent(
							sdk.NewEvent("challenge"+fmt.Sprint(approval.ApprovalId)+fmt.Sprint(challengeId)+fmt.Sprint(leafIndex)+fmt.Sprint(approverAddress)+fmt.Sprint(approvalLevel)+fmt.Sprint(newNumUsed),
								sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
								sdk.NewAttribute("creator", creatorAddress),
								sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
								sdk.NewAttribute("challengeTrackerId", fmt.Sprint(challengeId)),
								sdk.NewAttribute("approvalId", fmt.Sprint(approval.ApprovalId)),
								sdk.NewAttribute("leafIndex", fmt.Sprint(leafIndex.Sub(leftmostLeafIndex))),
								sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
								sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
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
