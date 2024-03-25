package keeper

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

// The UserApprovalsToCheck struct is used to keep track of which incoming / outgoing approvals for which addresses we need to check.
type UserApprovalsToCheck struct {
	Address  string
	Balances []*types.Balance
	Outgoing bool
}

// DeductUserOutgoingApprovals will check if the current transfer is approved from the from's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserOutgoingApprovals(ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.MerkleProof, challengeIdsIncremented *[]string, trackerIdsIncremented *[]string, prioritizedApprovals []*types.ApprovalIdentifierDetails, onlyCheckPrioritized bool) error {
	currApprovals := userBalance.OutgoingApprovals
	if userBalance.AutoApproveSelfInitiatedOutgoingTransfers {
		currApprovals = AppendSelfInitiatedOutgoingApproval(currApprovals, from)
	}

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastOutgoingTransfersToCollectionTransfers(currApprovals, from)
	_, err := k.DeductAndGetUserApprovals(overallTransferBalances, castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "outgoing", from, challengeIdsIncremented, trackerIdsIncremented, prioritizedApprovals, onlyCheckPrioritized)
	return err
}

// DeductUserIncomingApprovals will check if the current transfer is approved from the to's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserIncomingApprovals(ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.MerkleProof, challengeIdsIncremented *[]string, trackerIdsIncremented *[]string, prioritizedApprovals []*types.ApprovalIdentifierDetails, onlyCheckPrioritized bool) error {
	currApprovals := userBalance.IncomingApprovals
	if userBalance.AutoApproveSelfInitiatedIncomingTransfers {
		currApprovals = AppendSelfInitiatedIncomingApproval(currApprovals, to)
	}

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastIncomingTransfersToCollectionTransfers(currApprovals, to)
	_, err := k.DeductAndGetUserApprovals(overallTransferBalances, castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "incoming", to, challengeIdsIncremented, trackerIdsIncremented, prioritizedApprovals, onlyCheckPrioritized)
	return err
}

// DeductCollectionApprovalsAndGetUserApprovalsToCheck will check if the current transfer is allowed via the collection's approved transfers and handle any tallying accordingly
func (k Keeper) DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.MerkleProof, challengeIdsIncremented *[]string, trackerIdsIncremented *[]string, prioritizedApprovals []*types.ApprovalIdentifierDetails, onlyCheckPrioritized bool) ([]*UserApprovalsToCheck, error) {
	approvals := collection.CollectionApprovals
	return k.DeductAndGetUserApprovals(overallTransferBalances, approvals, ctx, collection, badgeIds, times, fromAddress, toAddress, initiatedBy, amount, solutions, "collection", "", challengeIdsIncremented, trackerIdsIncremented, prioritizedApprovals, onlyCheckPrioritized)
}

func (k Keeper) DeductAndGetUserApprovals(overallTransferBalances []*types.Balance, _approvals []*types.CollectionApproval, ctx sdk.Context, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.MerkleProof, approvalLevel string, approverAddress string, challengeIdsIncremented *[]string, trackerIdsIncremented *[]string, prioritizedApprovals []*types.ApprovalIdentifierDetails, onlyCheckPrioritized bool) ([]*UserApprovalsToCheck, error) {
	//Reorder approvals based on prioritized approvals
	approvals := []*types.CollectionApproval{}
	prioritizedTransfers := []*types.CollectionApproval{}
	for _, approval := range _approvals {
		prioritized := false

		for _, prioritizedApproval := range prioritizedApprovals {
			if approval.ApprovalId == prioritizedApproval.ApprovalId && prioritizedApproval.ApprovalLevel == approvalLevel {
				if (prioritizedApproval.ApprovalLevel == "incoming" && toAddress != prioritizedApproval.ApproverAddress) ||
					(prioritizedApproval.ApprovalLevel == "outgoing" && fromAddress != prioritizedApproval.ApproverAddress) {
					continue
				}
				prioritized = true
				break
			}
		}

		if prioritized {
			prioritizedTransfers = append(prioritizedTransfers, approval)
		} else {
			approvals = append(approvals, approval)
		}
	}

	if onlyCheckPrioritized {
		approvals = prioritizedTransfers
	} else {
		approvals = append(prioritizedTransfers, approvals...)
	}

	//HACK: We first expand all transfers to have just a len == 1 AllowedCombination[] so that we can easily check IsApproved later
	//		  This is because GetFirstMatchOnly will break down the transfers into smaller parts and without expansion, fetching if a certain transfer is allowed is impossible.
	expandedApprovals := approvals
	unhandled := []*types.UniversalPermissionDetails{}
	for _, badgeId := range badgeIds {
		for _, time := range times {
			unhandled = append(unhandled, &types.UniversalPermissionDetails{
				BadgeId:         badgeId,
				OwnershipTime:   time,
				TimelineTime:    &types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
				TransferTime:    &types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
				ToList:          &types.AddressList{},
				FromList:        &types.AddressList{},
				InitiatedByList: &types.AddressList{},
				ApprovalIdList:  &types.AddressList{},
			})
		}
	}
	//Keep a running tally of all the badges we still have to handle
	remainingBalances := []*types.Balance{
		{
			Amount:         amount,
			BadgeIds:       badgeIds,
			OwnershipTimes: times,
		},
	}

	//Step 2: For each approved transfer, we check if the transfer is allowed
	//1: If transfer meets all criteria FULLY (we don't support overflows), we mark it as handled and continue
	//2. If transfer does not meet all criteria, we continue and do not mark as handled
	//3. At the end, if there are any unhandled transfers, we throw
	userApprovalsToCheck := []*UserApprovalsToCheck{}
	for _, transferVal := range expandedApprovals {
		remainingBalances = types.FilterZeroBalances(remainingBalances)
		if len(remainingBalances) == 0 {
			break
		}

		doAddressesMatch := k.CheckIfAddressesMatchCollectionListIds(ctx, transferVal, fromAddress, toAddress, initiatedBy)
		if !doAddressesMatch {
			continue
		}

		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currTimeFound, err := types.SearchUintRangesForUint(currTime, transferVal.TransferTimes)
		if !currTimeFound || err != nil {
			continue
		}

		//From here on, we must have a full match or else continue

		//For the overlapping badges and ownership times, we have a match because list IDs, time, and badge IDs match.
		//We can now proceed to check any restrictions.
		transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + transferVal.BadgeIds[0].Start.String() + ", time: " + transferVal.TransferTimes[0].Start.String() + ", ownershipTime: " + transferVal.OwnershipTimes[0].Start.String() + ")"

		if transferVal.ApprovalCriteria == nil {
			//If there are no restrictions, it is a full match
			//Setting remainingBalances to remaining will set everything to handled since remaining is empty
			allBalancesForIdsAndTimes, err := types.GetBalancesForIds(ctx, transferVal.BadgeIds, transferVal.OwnershipTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for transfer: %s", transferStr)
			}

			remainingBalances, err = types.SubtractBalances(ctx, allBalancesForIdsAndTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: underflow error subtracting balances for transfer: %s", transferStr)
			}
		}

		//If we are here, we have a match and we can proceed to check the restrictions
		//We have to satisfy at least one of the approval details in full to be allowed
		//We scan linearly through them
		//Note that we do not overflow into the next. Each must match in full
		if transferVal.ApprovalCriteria != nil {
			approvalCriteria := transferVal.ApprovalCriteria
			if approvalCriteria.RequireFromDoesNotEqualInitiatedBy && fromAddress == initiatedBy {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because from == initiatedBy: %s", transferStr)
			}

			if approvalCriteria.RequireFromEqualsInitiatedBy && fromAddress != initiatedBy {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because from != initiatedBy: %s", transferStr)
			}

			if approvalCriteria.RequireToDoesNotEqualInitiatedBy && toAddress == initiatedBy {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because to == initiatedBy: %s", transferStr)
			}

			if approvalCriteria.RequireToEqualsInitiatedBy && toAddress != initiatedBy {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because to != initiatedBy: %s", transferStr)
			}

			//Assert that initiatedBy owns the required badges
			failedMustOwnBadges := false
			for _, mustOwnBadge := range approvalCriteria.MustOwnBadges {
				balances := []*types.Balance{}
				balancesFound := false

				alreadyChecked := []sdkmath.Uint{} //Do not want to circularly check the same collection
				currCollectionId := mustOwnBadge.CollectionId
				for !balancesFound {
					//Check if we have already searched this collection
					for _, alreadyCheckedId := range alreadyChecked {
						if alreadyCheckedId.Equal(currCollectionId) {
							return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrCircularInheritance, "circular inheritance detected for collection %s", currCollectionId)
						}
					}

					//Check if the collection has inherited balances
					collection, found := k.GetCollectionFromStore(ctx, currCollectionId)
					if !found {
						//Just continue with blank balances
						balancesFound = true
					} else {
						isStandardBalances := collection.BalancesType == "Standard"
						// isInheritedBalances := collection.BalancesType == "Inherited"
						if isStandardBalances {
							balancesFound = true
							initiatedByBalanceKey := ConstructBalanceKey(initiatedBy, currCollectionId)
							initiatedByBalance, found := k.GetUserBalanceFromStore(ctx, initiatedByBalanceKey)
							if found {
								balances = initiatedByBalance.Balances
							}
						} else {
							return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrWrongBalancesType, "must own badges must have Standard balances type %s", collection.CollectionId)
						}
					}
				}

				if mustOwnBadge.OverrideWithCurrentTime {
					mustOwnBadge.OwnershipTimes = []*types.UintRange{{Start: currTime, End: currTime}}
				}

				fetchedBalances, err := types.GetBalancesForIds(ctx, mustOwnBadge.BadgeIds, mustOwnBadge.OwnershipTimes, balances)
				if err != nil {
					failedMustOwnBadges = true
					break
					// continue
					// return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for mustOwnBadges: %s", transferStr)
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
				continue
			}

			//Get max balances allowed for this approvalCriteria element
			//Get the max balances allowed for this approvalCriteria element WITHOUT incrementing
			transferBalancesToCheck, err := types.GetBalancesForIds(ctx, transferVal.BadgeIds, transferVal.OwnershipTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for transfer: %s", transferStr)
			}

			transferBalancesToCheck = types.FilterZeroBalances(transferBalancesToCheck)
			if len(transferBalancesToCheck) == 0 {
				continue
			}
			//The section below are simply simulations seeing if the eventual increments will be allowed
			//This is because we do not want to increment something then continue and find out that it is not allowed
			//We prefer to check if it is allowed first, then increment if it is
			//This is why all the simulate params are set to true

			//Simulate to get challengeNumIncrements
			challengeNumIncrements, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId,
				transferVal.ChallengeTrackerId,
				[]*types.MerkleChallenge{
					approvalCriteria.MerkleChallenge,
				}, solutions, initiatedBy, true, approverAddress, approvalLevel, challengeIdsIncremented, transferVal)
			if err != nil {
				continue
				// return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "%s", transferStr)
			}

			//here, we assert the transfer is good for each level of approvals and increment if necessary
			maxPossible, err := k.GetMaxPossible(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.OverallApprovalAmount, approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "overall", "", true, trackerIdsIncremented)
			if err != nil {
				continue
			}
			transferBalancesToCheck, err = types.GetOverlappingBalances(ctx, maxPossible, transferBalancesToCheck)
			if err != nil {
				continue
			}

			maxPossible, err = k.GetMaxPossible(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.PerToAddressApprovalAmount, approvalCriteria.MaxNumTransfers.PerToAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "to", toAddress, true, trackerIdsIncremented)
			if err != nil {
				continue
			}
			transferBalancesToCheck, err = types.GetOverlappingBalances(ctx, maxPossible, transferBalancesToCheck)
			if err != nil {
				continue
			}

			maxPossible, err = k.GetMaxPossible(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount, approvalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "from", fromAddress, true, trackerIdsIncremented)
			if err != nil {
				continue
			}
			transferBalancesToCheck, err = types.GetOverlappingBalances(ctx, maxPossible, transferBalancesToCheck)
			if err != nil {
				continue
			}

			maxPossible, err = k.GetMaxPossible(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.PerInitiatedByAddressApprovalAmount, approvalCriteria.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "initiatedBy", initiatedBy, true, trackerIdsIncremented)
			if err != nil {
				continue
			}
			transferBalancesToCheck, err = types.GetOverlappingBalances(ctx, maxPossible, transferBalancesToCheck)
			if err != nil {
				continue
			}

			transferBalancesToCheck = types.FilterZeroBalances(transferBalancesToCheck)
			if len(transferBalancesToCheck) == 0 {
				continue
			}

			//here, we assert the transfer is good for each level of approvals and increment if necessary
			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.OverallApprovalAmount, approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "overall", "", true, trackerIdsIncremented)
			if err != nil {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded overall approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.PerToAddressApprovalAmount, approvalCriteria.MaxNumTransfers.PerToAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "to", toAddress, true, trackerIdsIncremented)
			if err != nil {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded to approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount, approvalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "from", fromAddress, true, trackerIdsIncremented)
			if err != nil {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded from approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.PerInitiatedByAddressApprovalAmount, approvalCriteria.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "initiatedBy", initiatedBy, true, trackerIdsIncremented)
			if err != nil {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded initiatedBy approvals: %s", transferStr)
			}

			//Finally, increment everything in store
			remainingBalances, err = types.SubtractBalances(ctx, transferBalancesToCheck, remainingBalances)
			if err != nil {
				continue
			}

			//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
			//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
			//    If so, useLeafIndexForNumIncrements will be true
			challengeNumIncrements, err = k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId,
				transferVal.ChallengeTrackerId,
				[]*types.MerkleChallenge{approvalCriteria.MerkleChallenge}, solutions, initiatedBy, false, approverAddress, approvalLevel, challengeIdsIncremented, transferVal)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "%s", transferStr)
			}

			//here, we assert the transfer is good for each level of approvals and increment if necessary
			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.OverallApprovalAmount, approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "overall", "", false, trackerIdsIncremented)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded overall approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.PerToAddressApprovalAmount, approvalCriteria.MaxNumTransfers.PerToAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "to", toAddress, false, trackerIdsIncremented)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded to approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount, approvalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "from", fromAddress, false, trackerIdsIncremented)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded from approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalCriteria, overallTransferBalances, collection, approvalCriteria.ApprovalAmounts.PerInitiatedByAddressApprovalAmount, approvalCriteria.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "initiatedBy", initiatedBy, false, trackerIdsIncremented)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded initiatedBy approvals: %s", transferStr)
			}

			//If we do not override the approved outgoing / incoming transfers, we need to check the user approvals
			if !approvalCriteria.OverridesFromOutgoingApprovals {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address:  fromAddress,
					Balances: transferBalancesToCheck,
					Outgoing: true,
				})
			}

			if !approvalCriteria.OverridesToIncomingApprovals {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address:  toAddress,
					Balances: transferBalancesToCheck,
					Outgoing: false,
				})
			}
		}

	}

	//If we didn't find a successful approval, we throw
	if len(remainingBalances) > 0 {
		// return []*UserApprovalsToCheck{}, ErrInadequateApprovals
		//convert ownership time unix milliseconds number to string
		timeToConvert := remainingBalances[0].OwnershipTimes[0].Start //unix milliseconds
		//convert timeToConvert to human readable date
		dateStr := time.Unix(0, int64(timeToConvert.Uint64())).Format(time.UnixDate)

		transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + remainingBalances[0].BadgeIds[0].Start.String() + ", ownership time: " + dateStr + ")"

		return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrInadequateApprovals, "no approval found for transfer: %s", transferStr)
	}

	return userApprovalsToCheck, nil
}

func IncrementBalances(startBalances []*types.Balance, numIncrements sdkmath.Uint, incrementOwnershipTimesBy sdkmath.Uint, incrementBadgeIdsBy sdkmath.Uint) ([]*types.Balance, error) {
	balances := types.DeepCopyBalances(startBalances)

	for _, startBalance := range balances {
		for _, time := range startBalance.OwnershipTimes {
			time.Start = time.Start.Add(numIncrements.Mul(incrementOwnershipTimesBy))
			time.End = time.End.Add(numIncrements.Mul(incrementOwnershipTimesBy))
		}

		for _, badgeId := range startBalance.BadgeIds {
			badgeId.Start = badgeId.Start.Add(numIncrements.Mul(incrementBadgeIdsBy))
			badgeId.End = badgeId.End.Add(numIncrements.Mul(incrementBadgeIdsBy))
		}
	}

	return balances, nil
}

func (k Keeper) GetMaxPossible(
	ctx sdk.Context,
	transferVal *types.CollectionApproval,
	approvalCriteria *types.ApprovalCriteria,
	overallTransferBalances []*types.Balance,
	collection *types.BadgeCollection,
	approvedAmount sdkmath.Uint,
	maxNumTransfers sdkmath.Uint,
	transferBalances []*types.Balance,
	challengeNumIncrements sdkmath.Uint,
	approverAddress string,
	approvalLevel string,
	trackerType string,
	address string,
	simulate bool,
	trackerIdsIncremented *[]string,
) ([]*types.Balance, error) {
	amountTrackerId := transferVal.AmountTrackerId
	predeterminedBalances := approvalCriteria.PredeterminedBalances
	allApprovals := []*types.Balance{{
		Amount:         approvedAmount,
		OwnershipTimes: transferVal.OwnershipTimes,
		BadgeIds:       transferVal.BadgeIds,
	}}

	//Get the current approvals for this transfer
	//If nil, no restrictions and we are approved for the entire transfer
	//Note we filter any excess badge IDs later and apply num increments as well
	if approvedAmount.IsNil() {
		approvedAmount = sdkmath.NewUint(0)
	}

	if maxNumTransfers.IsNil() {
		maxNumTransfers = sdkmath.NewUint(0)
	}

	needToFetchApprovalTrackerDetails := maxNumTransfers.GT(sdkmath.NewUint(0)) || approvedAmount.GT(sdkmath.NewUint(0)) ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers && trackerType == "overall") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy")

	approvalTrackerDetails := types.ApprovalTracker{
		Amounts:      []*types.Balance{},
		NumTransfers: sdkmath.NewUint(0),
	}

	if needToFetchApprovalTrackerDetails {
		fetchedDetails, found := k.GetApprovalTrackerFromStore(ctx, collection.CollectionId, approverAddress, amountTrackerId, approvalLevel, trackerType, address)
		if found {
			approvalTrackerDetails = fetchedDetails
		}
	}

	//Increment amounts and numTransfers and add back to store
	if approvedAmount.GT(sdkmath.NewUint(0)) {
		//Assume that if approvalTrackerDetails.Amounts is already not nil, it is correct and has been incremented properly
		//Here, we ONLY check if the NEW transferBalances makes it exceed the threshold
		currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(ctx, transferVal.BadgeIds, transferVal.OwnershipTimes, approvalTrackerDetails.Amounts)
		if err != nil {
			return nil, err
		}

		maxBalancesWeCanAdd, err := types.SubtractBalances(ctx, currTallyForCurrentIdsAndTimes, allApprovals)
		if err != nil {
			return nil, err
		}

		return maxBalancesWeCanAdd, nil
	}

	return transferBalances, nil
}

func (k Keeper) IncrementApprovalsAndAssertWithinThreshold(
	ctx sdk.Context,
	transferVal *types.CollectionApproval,
	approvalCriteria *types.ApprovalCriteria,
	overallTransferBalances []*types.Balance,
	collection *types.BadgeCollection,
	approvedAmount sdkmath.Uint,
	maxNumTransfers sdkmath.Uint,
	transferBalances []*types.Balance,
	challengeNumIncrements sdkmath.Uint,
	approverAddress string,
	approvalLevel string,
	trackerType string,
	address string,
	simulate bool,
	trackerIdsIncremented *[]string,
) error {
	amountTrackerId := transferVal.AmountTrackerId
	predeterminedBalances := approvalCriteria.PredeterminedBalances
	allApprovals := []*types.Balance{{
		Amount:         approvedAmount,
		OwnershipTimes: transferVal.OwnershipTimes,
		BadgeIds:       transferVal.BadgeIds,
	}}

	//Get the current approvals for this transfer
	//If nil, no restrictions and we are approved for the entire transfer
	//Note we filter any excess badge IDs later and apply num increments as well
	err := *new(error)
	if approvedAmount.IsNil() {
		approvedAmount = sdkmath.NewUint(0)
	}

	if maxNumTransfers.IsNil() {
		maxNumTransfers = sdkmath.NewUint(0)
	}

	needToFetchApprovalTrackerDetails := maxNumTransfers.GT(sdkmath.NewUint(0)) || approvedAmount.GT(sdkmath.NewUint(0)) ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers && trackerType == "overall") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy")

	approvalTrackerDetails := types.ApprovalTracker{
		Amounts:      []*types.Balance{},
		NumTransfers: sdkmath.NewUint(0),
	}

	trackerIncrementId := fmt.Sprintf("%s-%s-%s-%s-%s-%s", amountTrackerId, approverAddress, approvalLevel, trackerType, address, collection.CollectionId)
	if needToFetchApprovalTrackerDetails {
		fetchedDetails, found := k.GetApprovalTrackerFromStore(ctx, collection.CollectionId, approverAddress, amountTrackerId, approvalLevel, trackerType, address)
		if found {
			approvalTrackerDetails = fetchedDetails
		}
	}

	alreadyIncremented := false
	for _, trackerId := range *trackerIdsIncremented {
		if trackerId == trackerIncrementId {
			alreadyIncremented = true
			break
		}
	}

	//Here, we handle if predeterminedBalances is set
	//If it is, we need to check that the calculated predetermined transfer (based on numIncrements) matches the current transfer
	//We say it matches if for the transferVal's badge IDs and times (note not the overlaps), the amounts of the overallTransferBalances
	// 	 match exactly, including not specifying any badge IDs or times that are not within the predetermined transfer
	//
	//Note that the specific overlap is not used here, but it is inherently taken into account because the overlap is a subset of the transferVal's badge IDs and times
	//We only do this on overall because we only need to do it once (and num increments for "Sequential" order types is based on the "Overall" timeline)

	if predeterminedBalances != nil {
		numIncrements := sdkmath.NewUint(0)
		toBeCalculated := true
		trackerNumTransfers := approvalTrackerDetails.NumTransfers
		if alreadyIncremented {
			trackerNumTransfers = trackerNumTransfers.Sub(sdkmath.NewUint(1))
		}

		if predeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {
			numIncrements = challengeNumIncrements
		} else if predeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers && trackerType == "overall" {
			numIncrements = trackerNumTransfers
		} else if predeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to" {
			numIncrements = trackerNumTransfers
		} else if predeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from" {
			numIncrements = trackerNumTransfers
		} else if predeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy" {
			numIncrements = trackerNumTransfers
		} else {
			toBeCalculated = false
		}

		if toBeCalculated {
			calculatedBalances := []*types.Balance{}

			if predeterminedBalances.ManualBalances != nil {
				if numIncrements.LT(sdkmath.NewUint(uint64(len(predeterminedBalances.ManualBalances)))) {
					calculatedBalances = types.DeepCopyBalances(predeterminedBalances.ManualBalances[numIncrements.Uint64()].Balances)
				}
			} else if predeterminedBalances.IncrementedBalances != nil {
				calculatedBalances, err = IncrementBalances(predeterminedBalances.IncrementedBalances.StartBalances, numIncrements, predeterminedBalances.IncrementedBalances.IncrementOwnershipTimesBy, predeterminedBalances.IncrementedBalances.IncrementBadgeIdsBy)
				if err != nil {
					return err
				}
			}

			// //From the original transfer balances (the overall ones, not overlap), fetch the balances for the badge IDs and times of the current transfer
			// //Filter out any balances that have an amount of 0
			// fetchedBalances := []*types.Balance{}
			// for _, calculatedBalance := range calculatedBalances {
			// 	fetchedBalancesOfOverall, err := types.GetBalancesForIds( ctx, calculatedBalance.BadgeIds, calculatedBalance.OwnershipTimes, overallTransferBalances)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	fetchedBalances, err  = types.AddBalances(fetchedBalances, fetchedBalancesOfOverall)
			// 	if err != nil {
			// 		return err
			// 	}
			// }

			//Assert that we have exactly the amount specified in the original transfers
			equal := types.AreBalancesEqual(ctx, overallTransferBalances, calculatedBalances, false)
			if !equal {
				return sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because predetermined balances do not match: %s", amountTrackerId)
			}
		}
	}

	//Increment amounts and numTransfers and add back to store
	if approvedAmount.GT(sdkmath.NewUint(0)) {
		//Assume that if approvalTrackerDetails.Amounts is already not nil, it is correct and has been incremented properly
		//Here, we ONLY check if the NEW transferBalances makes it exceed the threshold
		currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(ctx, transferVal.BadgeIds, transferVal.OwnershipTimes, approvalTrackerDetails.Amounts)
		if err != nil {
			return err
		}

		//If this passes, the new transferBalances are okay
		_, err = types.AddBalancesAndAssertDoesntExceedThreshold(ctx, currTallyForCurrentIdsAndTimes, transferBalances, allApprovals)
		if err != nil {
			return err
		}

		//We then add them to the current tally of ALL ids and times
		approvalTrackerDetails.Amounts, err = types.AddBalances(ctx, approvalTrackerDetails.Amounts, transferBalances)
		if err != nil {
			return err
		}
	}

	if maxNumTransfers.GT(sdkmath.NewUint(0)) ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers && trackerType == "overall") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy") {
		if !alreadyIncremented {
			approvalTrackerDetails.NumTransfers = approvalTrackerDetails.NumTransfers.Add(sdkmath.NewUint(1))
			//only check exceeds if maxNumTransfers is not 0 (because 0 means no limit)
			if maxNumTransfers.GT(sdkmath.NewUint(0)) {
				if approvalTrackerDetails.NumTransfers.GT(maxNumTransfers) {
					return sdkerrors.Wrapf(ErrDisallowedTransfer, "exceeded max transfers allowed - %s", maxNumTransfers.String())
				}
			}
		}
	}

	if needToFetchApprovalTrackerDetails && !simulate {
		if !alreadyIncremented {
			*trackerIdsIncremented = append(*trackerIdsIncremented, trackerIncrementId)
		}

		//Currently added for indexer, but note that it is planned to be deprecated
		amountsJsonData, err := json.Marshal(approvalTrackerDetails.Amounts)
		if err != nil {
			return err
		}
		amountsStr := string(amountsJsonData)

		numTransfersJsonData, err := json.Marshal(approvalTrackerDetails.NumTransfers)
		if err != nil {
			return err
		}
		numTransfersStr := string(numTransfersJsonData)

		//TODO: Deprecate this eventually in favor of doing exclusively in indexer
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("approval"+fmt.Sprint(collection.CollectionId)+fmt.Sprint(approverAddress)+fmt.Sprint(amountTrackerId)+fmt.Sprint(approvalLevel)+fmt.Sprint(trackerType)+fmt.Sprint(address),
				sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
				sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
				sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
				sdk.NewAttribute("amountTrackerId", fmt.Sprint(amountTrackerId)),
				sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
				sdk.NewAttribute("trackerType", fmt.Sprint(trackerType)),
				sdk.NewAttribute("approvedAddress", fmt.Sprint(address)),
				sdk.NewAttribute("amounts", amountsStr),
				sdk.NewAttribute("numTransfers", numTransfersStr),
			),
		)

		err = k.SetApprovalTrackerInStore(ctx, collection.CollectionId, approverAddress, amountTrackerId, approvalTrackerDetails, approvalLevel, trackerType, address)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) GetPredeterminedBalancesForPrecalculationId(ctx sdk.Context, approvals []*types.CollectionApproval, collection *types.BadgeCollection, approverAddress string, precalculationId string, approvalLevel string, address string, solutions []*types.MerkleProof, initiatedBy string) ([]*types.Balance, error) {
	approvalId := ""
	for _, transfer := range approvals {
		approvalCriteria := transfer.ApprovalCriteria
		approvalId = transfer.ApprovalId
		amountTrackerId := transfer.AmountTrackerId
		if approvalCriteria == nil {
			continue
		}

		if approvalId != precalculationId {
			continue
		}

		if approvalId == "" || amountTrackerId == "" {
			continue
		}

		if approvalCriteria.PredeterminedBalances != nil {
			numIncrements := sdkmath.NewUint(0)
			if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {

				//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
				//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
				numIncrementsFetched, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId,
					transfer.ChallengeTrackerId,
					[]*types.MerkleChallenge{
						approvalCriteria.MerkleChallenge,
					}, solutions, initiatedBy, true, address, approvalLevel, &[]string{}, transfer)
				if err != nil {
					return []*types.Balance{}, sdkerrors.Wrapf(err, "invalid challenges / solutions")
				}

				numIncrements = numIncrementsFetched
			} else {
				trackerType := "overall"
				if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers {
					trackerType = "from"
				} else if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers {
					trackerType = "to"
				} else if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers {
					trackerType = "initiatedBy"
				}

				approvalTrackerDetails, found := k.GetApprovalTrackerFromStore(ctx, collection.CollectionId, approverAddress, amountTrackerId, approvalLevel, trackerType, address)
				if !found {
					approvalTrackerDetails = types.ApprovalTracker{
						Amounts:      []*types.Balance{},
						NumTransfers: sdkmath.NewUint(0),
					}
				}

				numIncrements = approvalTrackerDetails.NumTransfers
			}

			predeterminedBalances := []*types.Balance{}
			if approvalCriteria.PredeterminedBalances.ManualBalances != nil {
				if numIncrements.LT(sdkmath.NewUint(uint64(len(approvalCriteria.PredeterminedBalances.ManualBalances)))) {
					predeterminedBalances = types.DeepCopyBalances(approvalCriteria.PredeterminedBalances.ManualBalances[numIncrements.Uint64()].Balances)
				}
			} else if approvalCriteria.PredeterminedBalances.IncrementedBalances != nil {
				err := *new(error)
				predeterminedBalances, err = IncrementBalances(approvalCriteria.PredeterminedBalances.IncrementedBalances.StartBalances, numIncrements, approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementOwnershipTimesBy, approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementBadgeIdsBy)
				if err != nil {
					return []*types.Balance{}, err
				}
			}

			return predeterminedBalances, nil
		} else {
			return []*types.Balance{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "no predetermined transfers found for approval id: %s", precalculationId)
		}

	}

	return []*types.Balance{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "no predetermined transfers found for approval id: %s", precalculationId)
}
