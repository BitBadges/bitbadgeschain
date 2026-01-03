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

const (
	MaxMerkleProofLength = 10
)

// HandleMerkleChallenges processes merkle challenges for approval validation
// Returns (deterministicErrorMsg, numIncrements, error) where deterministicErrorMsg is a deterministic error string
func (k Keeper) HandleMerkleChallenges(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	transfer *types.Transfer,
	approval *types.CollectionApproval,
	transferMetadata TransferMetadata,
	simulation bool,
) (string, sdkmath.Uint, error) {
	creatorAddress := transferMetadata.InitiatedBy
	approverAddress := transferMetadata.ApproverAddress
	approvalLevel := transferMetadata.ApprovalLevel
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
			detErrMsg := fmt.Sprintf("expected to calculate balances from challenge but no matching challenge for challenge tracker id %s", approval.ApprovalCriteria.PredeterminedBalances.OrderCalculationMethod.ChallengeTrackerId)
			return detErrMsg, numIncrements, sdkerrors.Wrap(ErrNoMatchingChallengeForChallengeTrackerId, detErrMsg)
		}
	}

	for _, challenge := range challenges {
		if challenge == nil || challenge.Root == "" {
			detErrMsg := "challenge is nil or has empty root"
			return detErrMsg, numIncrements, sdkerrors.Wrap(types.ErrChallengeTrackerIdIsNil, detErrMsg)
		}

		// Early validation of proof length to prevent DoS attacks
		// This check must happen before processing any proofs to prevent gas consumption attacks
		if challenge.ExpectedProofLength.GT(sdkmath.NewUint(MaxMerkleProofLength)) {
			detErrMsg := fmt.Sprintf("expected proof length %s exceeds maximum allowed %d", challenge.ExpectedProofLength.String(), MaxMerkleProofLength)
			return detErrMsg, numIncrements, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
		}

		// Additional validation: check actual proof lengths before processing to prevent DoS
		for _, proof := range merkleProofs {
			if len(proof.Aunts) > MaxMerkleProofLength {
				detErrMsg := fmt.Sprintf("proof length %d exceeds maximum allowed %d", len(proof.Aunts), MaxMerkleProofLength)
				return detErrMsg, numIncrements, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
			}
		}

		challengeId := challenge.ChallengeTrackerId
		root := challenge.Root
		hasValidSolution := false

		baseErrorStr := ""
		if challenge.UseCreatorAddressAsLeaf {
			baseErrorStr = "does not satisfy whitelist"
		} else {
			baseErrorStr = "invalid code / password"
		}

		// We check that 1 of N proofs is valid
		detailedErrorStr := ""
		for _, proof := range merkleProofs {
			detailedErrorStr = ""
			if root != "" {
				// Must be proper length to avoid preimage attacks
				if len(proof.Aunts) != int(challenge.ExpectedProofLength.Uint64()) {
					detailedErrorStr = "invalid proof length"
					continue
				}

				// Use local variable to avoid mutating the original transfer object
				leafValue := proof.Leaf
				if challenge.UseCreatorAddressAsLeaf {
					leafValue = creatorAddress
				}

				if leafValue == "" {
					detailedErrorStr = "empty leaf"
					continue
				}

				if challenge.LeafSigner != "" {
					leafSignature := proof.LeafSignature

					leafSignerEthAddress := challenge.LeafSigner
					if leafSignerEthAddress == "" {
						detailedErrorStr = "empty leaf signer"
						continue
					}

					ethAddress := ethcommon.HexToAddress(leafSignerEthAddress)

					leafSignatureString := leafValue + "-" + creatorAddress
					isValid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
						ethAddress,
						[]byte(leafSignatureString),
						leafSignature,
					)

					if !isValid || err != nil {
						detailedErrorStr = "invalid leaf signature"
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
					// Prevent underflow: ensure leafIndex >= leftmostLeafIndex
					// If leafIndex < leftmostLeafIndex, the proof is invalid for predetermined balance calculation
					if leafIndex.LT(leftmostLeafIndex) {
						detErrMsg := fmt.Sprintf("leaf index %s is less than leftmost leaf index %s, invalid for predetermined balance calculation", leafIndex.String(), leftmostLeafIndex.String())
						return detErrMsg, numIncrements, sdkerrors.Wrap(types.ErrInvalidRequest, detErrMsg)
					}
					numIncrements = leafIndex.Sub(leftmostLeafIndex)
				}

				err := CheckMerklePath(leafValue, root, proof.Aunts)
				if err != nil {
					detailedErrorStr = ""
					continue
				}

				//If there is a max uses per leaf, we need to check it has not exceeded the treshold uses
				if !challenge.MaxUsesPerLeaf.IsNil() && challenge.MaxUsesPerLeaf.GT(sdkmath.NewUint(0)) {
					numUsed, err := k.GetChallengeTrackerFromStore(ctx, collectionId, approverAddress, approvalLevel, approval.ApprovalId, challengeId, leafIndex.Sub(leftmostLeafIndex))
					if err != nil {
						detailedErrorStr = "error getting num processed"
						continue
					}
					numUsed = numUsed.Add(sdkmath.NewUint(1))

					maxUses := challenge.MaxUsesPerLeaf
					if numUsed.GT(maxUses) {
						detailedErrorStr = "exceeded max number of uses"
						continue
					}

					//Increment the number of uses in store if we are doing it for real
					if !simulation {
						newNumUsed, err := k.IncrementChallengeTrackerInStore(ctx, collectionId, approverAddress, approvalLevel, approval.ApprovalId, challengeId, leafIndex.Sub(leftmostLeafIndex))
						if err != nil {
							detErrMsg := "failed to increment challenge tracker"
							return detErrMsg, numIncrements, sdkerrors.Wrap(err, detErrMsg)
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
			detErrMsg := fmt.Sprintf("%s - %s", baseErrorStr, detailedErrorStr)
			return detErrMsg, numIncrements, sdkerrors.Wrap(ErrNoValidSolutionForChallenge, detErrMsg)
		}
	}

	return "", numIncrements, nil
}

// HandleETHSignatureChallenges processes ETH signature challenges for approval validation
// Returns (deterministicErrorMsg, error) where deterministicErrorMsg is a deterministic error string
func (k Keeper) HandleETHSignatureChallenges(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	transfer *types.Transfer,
	approval *types.CollectionApproval,
	transferMetadata TransferMetadata,
) (string, error) {
	initiatorAddress := transferMetadata.InitiatedBy
	approverAddress := transferMetadata.ApproverAddress
	approvalLevel := transferMetadata.ApprovalLevel
	challenges := approval.ApprovalCriteria.EthSignatureChallenges
	ethSignatureProofs := transfer.EthSignatureProofs

	for _, challenge := range challenges {
		if challenge == nil || challenge.Signer == "" {
			detErrMsg := "challenge is nil or has empty signer"
			return detErrMsg, sdkerrors.Wrap(types.ErrChallengeTrackerIdIsNil, detErrMsg)
		}

		challengeId := challenge.ChallengeTrackerId
		signerAddress := challenge.Signer
		hasValidSolution := false

		// We check that 1 of N proofs is valid
		for _, proof := range ethSignatureProofs {
			if proof.Nonce == "" || proof.Signature == "" {
				continue
			}

			// Verify the signature
			// Signature scheme: ETHSign(nonce + "-" + initiatorAddress + "-" + collectionId + "-" + approverAddress + "-" + approvalLevel + "-" + approvalId + "-" + challengeId)
			ethAddress := ethcommon.HexToAddress(signerAddress)
			signatureString := proof.Nonce + "-" + initiatorAddress + "-" + collectionId.String() + "-" + approverAddress + "-" + approvalLevel + "-" + approval.ApprovalId + "-" + challengeId

			isValid, err := sigverify.VerifyEllipticCurveHexSignatureEx(
				ethAddress,
				[]byte(signatureString),
				proof.Signature,
			)

			if !isValid || err != nil {
				continue
			}

			// Check if this signature has already been used
			signatureKey := ConstructETHSignatureTrackerKey(collectionId, approverAddress, approvalLevel, approval.ApprovalId, challengeId, proof.Signature)
			numUsed, exists := k.GetETHSignatureTrackerFromStore(ctx, signatureKey)
			if !exists {
				numUsed = sdkmath.NewUint(0)
			}

			// Each signature can only be used once
			if numUsed.GT(sdkmath.NewUint(0)) {
				continue
			}

			// Increment the usage count
			newNumUsed, err := k.IncrementETHSignatureTrackerInStore(ctx, signatureKey)
			if err != nil {
				continue
			}

			// Currently added for indexer, but note that it is planned to be deprecated
			ctx.EventManager().EmitEvent(
				sdk.NewEvent("ethSignatureChallenge"+fmt.Sprint(approval.ApprovalId)+fmt.Sprint(challengeId)+fmt.Sprint(proof.Signature)+fmt.Sprint(approverAddress)+fmt.Sprint(approvalLevel)+fmt.Sprint(newNumUsed),
					sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
					sdk.NewAttribute("creator", initiatorAddress),
					sdk.NewAttribute("initiator", initiatorAddress),
					sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
					sdk.NewAttribute("challengeTrackerId", fmt.Sprint(challengeId)),
					sdk.NewAttribute("approvalId", fmt.Sprint(approval.ApprovalId)),
					sdk.NewAttribute("signature", fmt.Sprint(proof.Signature)),
					sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
					sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
					sdk.NewAttribute("numUsed", fmt.Sprint(newNumUsed)),
				),
			)

			hasValidSolution = true
			break
		}

		if !hasValidSolution {
			detErrMsg := "invalid ETH signature - signature not provided or already used"
			return detErrMsg, sdkerrors.Wrap(ErrNoValidSolutionForChallenge, detErrMsg)
		}
	}

	return "", nil
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
