package keeper

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"bitbadgeschain/third_party/go-rapidsnark"
	"bitbadgeschain/x/badges/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	rapidsnarktypes "github.com/iden3/go-rapidsnark/types"
)

func (k Keeper) HandleCoinTransfers(ctx sdk.Context, coinTransfers []*types.CoinTransfer, initiatedBy string, simulate bool) error {
	//simulate the sdk.Coin transfers
	initiatedByAcc := sdk.MustAccAddressFromBech32(initiatedBy)

	if simulate {
		spendableCoins := k.bankKeeper.SpendableCoins(ctx, initiatedByAcc)
		for _, coinTransfer := range coinTransfers {
			toTransfer := coinTransfer.Coins
			for _, coin := range toTransfer {
				newCoins, underflow := spendableCoins.SafeSub(*coin)
				if underflow {
					return sdkerrors.Wrapf(types.ErrUnderflow, "insufficient $BADGE balance to complete transfer")
				}
				spendableCoins = newCoins
			}
		}
	} else {
		for _, coinTransfer := range coinTransfers {
			coinsToTransfer := coinTransfer.Coins
			toAddressAcc := sdk.MustAccAddressFromBech32(coinTransfer.To)
			fromAddressAcc := initiatedByAcc
			for _, coin := range coinsToTransfer {
				err := k.bankKeeper.SendCoins(ctx, fromAddressAcc, toAddressAcc, sdk.NewCoins(*coin))
				if err != nil {
					return sdkerrors.Wrapf(err, "error sending $BADGE, passed simulation but not actual transfers")
				}
			}
		}
	}

	return nil
}

func (k Keeper) SimulateZKPs(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	zkps []*types.ZkProof, zkpSolutions []*types.ZkProofSolution, approverAddress string, approvalLevel string, approvalId string) ([]int, error) {
	validProofs := make([]int, len(zkps))
	for i := range validProofs {
		validProofs[i] = -1
	}

	//Assert valid solution for every ZKP
	for i, zkProof := range zkps {
		for j, zkProofSolution := range zkpSolutions {
			if validProofs[i] >= 0 {
				continue //already found a valid solution
			}

			verificationKey := []byte(zkProof.VerificationKey)

			proof := rapidsnarktypes.ProofData{}
			err := json.Unmarshal([]byte(zkProofSolution.Proof), &proof)
			if err != nil {
				continue
			}

			pubSignals := []string{}
			err = json.Unmarshal([]byte(zkProofSolution.PublicInputs), &pubSignals)
			if err != nil {
				continue
			}

			proofData := rapidsnarktypes.ZKProof{
				Proof:      &proof,
				PubSignals: pubSignals,
			}

			// verify the proof with the given verificationKey & publicSignals
			err = rapidsnark.VerifyGroth16(proofData, verificationKey)
			if err != nil {
				fmt.Println("Error verifying proof:", err)
				continue
			}

			proofHash := sha256.Sum256([]byte(zkProofSolution.Proof))
			proofHashStr := fmt.Sprintf("%x", proofHash)

			found, err := k.GetZKPFromStore(ctx, collection.CollectionId, approverAddress, approvalLevel, approvalId, zkProof.ZkpTrackerId, proofHashStr)
			if !found && err == nil {
				validProofs[i] = j
			}
		}
	}

	someProofIsInvalid := false
	for _, valid := range validProofs {
		if valid < 0 {
			someProofIsInvalid = true
		}
	}

	if someProofIsInvalid {
		return []int{}, sdkerrors.Wrapf(ErrInadequateApprovals, "zkp proofs are invalid")
	}

	return validProofs, nil
}

func (k Keeper) HandleZKPs(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	validZkpSolutionIdxs []int,
	zkps []*types.ZkProof, zkpSolutions []*types.ZkProofSolution, approverAddress string, approvalLevel string, approvalId string) error {

	//Assert valid solution for every ZKP
	for i, zkProof := range zkps {
		if validZkpSolutionIdxs[i] >= 0 {
			solutionIdx := validZkpSolutionIdxs[i]

			hashStr := sha256.Sum256([]byte(zkpSolutions[solutionIdx].Proof))
			proofHashStr := fmt.Sprintf("%x", hashStr)
			err := k.SetZKPAsUsedInStore(ctx, collection.CollectionId, approverAddress, approvalLevel, approvalId, zkProof.ZkpTrackerId, proofHashStr)
			if err != nil {
				return sdkerrors.Wrapf(err, "error setting zk proof as used in store")
			}
		} else {
			return sdkerrors.Wrapf(ErrInadequateApprovals, "zkp proofs are invalid. passed simulation but not actual")
		}
	}

	return nil
}

func (k Keeper) CheckMustOwnBadges(
	ctx sdk.Context,
	mustOwnBadges []*types.MustOwnBadges,
	initiatedBy string,
) error {
	//Assert that initiatedBy owns the required badges
	failedMustOwnBadges := false
	for _, mustOwnBadge := range mustOwnBadges {
		collection, found := k.GetCollectionFromStore(ctx, mustOwnBadge.CollectionId)
		if !found {
			failedMustOwnBadges = true
			break
		}

		initiatorBalances := k.GetBalanceOrApplyDefault(ctx, collection, initiatedBy)
		balances := initiatorBalances.Balances

		if mustOwnBadge.OverrideWithCurrentTime {
			currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
			mustOwnBadge.OwnershipTimes = []*types.UintRange{{Start: currTime, End: currTime}}
		}

		fetchedBalances, err := types.GetBalancesForIds(ctx, mustOwnBadge.BadgeIds, mustOwnBadge.OwnershipTimes, balances)
		if err != nil {
			failedMustOwnBadges = true
			break
		}

		satisfiesRequirementsForOne := false
		for _, fetchedBalance := range fetchedBalances {
			//check if amount is within range
			minAmount := mustOwnBadge.AmountRange.Start
			maxAmount := mustOwnBadge.AmountRange.End

			if fetchedBalance.Amount.LT(minAmount) || fetchedBalance.Amount.GT(maxAmount) {
				failedMustOwnBadges = true
			} else {
				satisfiesRequirementsForOne = true
			}
		}

		if mustOwnBadge.MustSatisfyForAllAssets && failedMustOwnBadges {
			break
		} else if !mustOwnBadge.MustSatisfyForAllAssets && satisfiesRequirementsForOne {
			failedMustOwnBadges = false
			break
		}
	}

	if failedMustOwnBadges {
		return sdkerrors.Wrapf(ErrInadequateApprovals, "failed must own badges")
	}

	return nil
}
