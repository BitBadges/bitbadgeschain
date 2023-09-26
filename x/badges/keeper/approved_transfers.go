package keeper

import (
	"encoding/json"
	"fmt"
	"math"

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
	currApprovedTransfers := userBalance.ApprovedOutgoingTransfers
	currApprovedTransfers = AppendDefaultForOutgoing(currApprovedTransfers, from)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastOutgoingTransfersToCollectionTransfers(currApprovedTransfers, from)
	_, err := k.DeductAndGetUserApprovals(overallTransferBalances, castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "outgoing", from, challengeIdsIncremented, trackerIdsIncremented, prioritizedApprovals, onlyCheckPrioritized)
	return err
}

// DeductUserIncomingApprovals will check if the current transfer is approved from the to's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserIncomingApprovals(ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.MerkleProof, challengeIdsIncremented *[]string, trackerIdsIncremented *[]string, prioritizedApprovals []*types.ApprovalIdentifierDetails, onlyCheckPrioritized bool) error {
	currApprovedTransfers := userBalance.ApprovedIncomingTransfers
	currApprovedTransfers = AppendDefaultForIncoming(currApprovedTransfers, to)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastIncomingTransfersToCollectionTransfers(currApprovedTransfers, to)
	_, err := k.DeductAndGetUserApprovals(overallTransferBalances, castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "incoming", to, challengeIdsIncremented, trackerIdsIncremented, prioritizedApprovals, onlyCheckPrioritized)
	return err
}

// DeductCollectionApprovalsAndGetUserApprovalsToCheck will check if the current transfer is allowed via the collection's approved transfers and handle any tallying accordingly
func (k Keeper) DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.MerkleProof, challengeIdsIncremented *[]string, trackerIdsIncremented *[]string, prioritizedApprovals []*types.ApprovalIdentifierDetails, onlyCheckPrioritized bool) ([]*UserApprovalsToCheck, error) {
	approvedTransfers := collection.CollectionApprovedTransfers
	return k.DeductAndGetUserApprovals(overallTransferBalances, approvedTransfers, ctx, collection, badgeIds, times, fromAddress, toAddress, initiatedBy, amount, solutions, "collection", "", challengeIdsIncremented, trackerIdsIncremented, prioritizedApprovals, onlyCheckPrioritized)
}

func (k Keeper) DeductAndGetUserApprovals(overallTransferBalances []*types.Balance, _approvedTransfers []*types.CollectionApprovedTransfer, ctx sdk.Context, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.MerkleProof, approvalLevel string, approverAddress string, challengeIdsIncremented *[]string, trackerIdsIncremented *[]string, prioritizedApprovals []*types.ApprovalIdentifierDetails, onlyCheckPrioritized bool) ([]*UserApprovalsToCheck, error) {
	//Reorder approvedTransfers based on prioritized approvals
	approvedTransfers := []*types.CollectionApprovedTransfer{}
	prioritizedTransfers := []*types.CollectionApprovedTransfer{}
	for _, approvedTransfer := range _approvedTransfers {
    prioritized := false

    for _, prioritizedApproval := range prioritizedApprovals {
        if approvedTransfer.ApprovalId == prioritizedApproval.ApprovalId && prioritizedApproval.ApprovalLevel == approvalLevel {
            if (prioritizedApproval.ApprovalLevel == "incoming" && toAddress != prioritizedApproval.ApproverAddress) ||
                (prioritizedApproval.ApprovalLevel == "outgoing" && fromAddress != prioritizedApproval.ApproverAddress) {
                continue
            }
            prioritized = true
            break
        }
    }

    if prioritized {
        prioritizedTransfers = append(prioritizedTransfers, approvedTransfer)
    } else {
        approvedTransfers = append(approvedTransfers, approvedTransfer)
    }
	}

	if onlyCheckPrioritized {
		approvedTransfers = prioritizedTransfers
	} else {
		approvedTransfers = append(prioritizedTransfers, approvedTransfers...)
	}

	//HACK: We first expand all transfers to have just a len == 1 AllowedCombination[] so that we can easily check IsApproved later
	//		  This is because GetFirstMatchOnly will break down the transfers into smaller parts and without expansion, fetching if a certain transfer is allowed is impossible.
	expandedApprovedTransfers := ExpandCollectionApprovedTransfers(approvedTransfers)
	manager := types.GetCurrentManager(ctx, collection)

	unhandled := []*types.UniversalPermissionDetails{}
	for _, badgeId := range badgeIds {
		for _, time := range times {
			unhandled = append(unhandled, &types.UniversalPermissionDetails{
				BadgeId:       badgeId,
				OwnershipTime: time,
				TimelineTime: &types.UintRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
				TransferTime: &types.UintRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
				ToMapping: &types.AddressMapping{},
				FromMapping: &types.AddressMapping{},
				InitiatedByMapping: &types.AddressMapping{},
				ApprovalTrackerIdMapping: &types.AddressMapping{},
				ChallengeTrackerIdMapping: &types.AddressMapping{},
			})
		}
	}

	//Step 1: we pre-check to make sure that there are no explicit disapprovals for the balances being transferred
	//If there are, we throw
	for _, transferVal := range expandedApprovedTransfers {
		transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + transferVal.BadgeIds[0].Start.String() + ", time: " + transferVal.TransferTimes[0].Start.String() + ", ownershipTime: " + transferVal.OwnershipTimes[0].Start.String() + ")"

		allowed := transferVal.AllowedCombinations[0].IsApproved //HACK: can do this because we expanded the allowed combinations above
		if allowed {
			continue
		}

		doAddressesMatch := k.CheckIfAddressesMatchCollectionMappingIds(ctx, transferVal, fromAddress, toAddress, initiatedBy, manager)
		if !doAddressesMatch {
			continue
		}

		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currTimeFound := types.SearchUintRangesForUint(currTime, transferVal.TransferTimes)
		if !currTimeFound {
			continue
		}

		for _, badgeId := range transferVal.BadgeIds {
			for _, time := range transferVal.OwnershipTimes {
				_, overlaps := types.UniversalRemoveOverlapFromValues(&types.UniversalPermissionDetails{
					BadgeId:       badgeId,
					OwnershipTime: time,
					//Dummy values
					TimelineTime: &types.UintRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
					TransferTime: &types.UintRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
					ToMapping: &types.AddressMapping{},
					FromMapping: &types.AddressMapping{},
					InitiatedByMapping: &types.AddressMapping{},
					ApprovalTrackerIdMapping: &types.AddressMapping{},
					ChallengeTrackerIdMapping: &types.AddressMapping{},
					}, unhandled)
		
				//If any of the requested badge IDs or ownership times are disallowed, we throw and disallow the entire transfer
				if len(overlaps) > 0 {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed: %s", transferStr)
				}
			}
		}
	}

	//Keep a running tally of all the badges we still have to handle
	remainingBalances := []*types.Balance{
		{
			Amount: amount,
			BadgeIds: badgeIds,
			OwnershipTimes: times,
		},
	}

	//Step 2: For each approved transfer, we check if the transfer is allowed
	//1: If transfer meets all criteria FULLY (we don't support overflows), we mark it as handled and continue
	//2. If transfer does not meet all criteria, we continue and do not mark as handled
	//3. At the end, if there are any unhandled transfers, we throw
	userApprovalsToCheck := []*UserApprovalsToCheck{}
	for _, transferVal := range expandedApprovedTransfers {
		remainingBalances = types.FilterZeroBalances(remainingBalances)
		if len(remainingBalances) == 0 {
			break
		}

		doAddressesMatch := k.CheckIfAddressesMatchCollectionMappingIds(ctx, transferVal, fromAddress, toAddress, initiatedBy, manager)
		if !doAddressesMatch {
			continue
		}

		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currTimeFound := types.SearchUintRangesForUint(currTime, transferVal.TransferTimes)
		if !currTimeFound {
			continue
		}

		//From here on, we must have a full match or else continue

		//For the overlapping badges and ownership times, we have a match because mapping IDs, time, and badge IDs match.
		//We can now proceed to check any restrictions.
		transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + transferVal.BadgeIds[0].Start.String() + ", time: " + transferVal.TransferTimes[0].Start.String() + ", ownershipTime: " + transferVal.OwnershipTimes[0].Start.String() + ")"
		
		//Technically, this is probably not necessary because we already checked if any transfer is disallowed in the first step
		allowed := transferVal.AllowedCombinations[0].IsApproved //HACK: can do this because we expanded the allowed combinations above
		if !allowed {
			continue
		}


		if transferVal.ApprovalDetails == nil {
			//If there are no restrictions, it is a full match
			//Setting remainingBalances to remaining will set everything to handled since remaining is empty
			allBalancesForIdsAndTimes, err := types.GetBalancesForIds(transferVal.BadgeIds, transferVal.OwnershipTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for transfer: %s", transferStr)
			}
			
			remainingBalances, err = types.SubtractBalances(allBalancesForIdsAndTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err subtracting balances for transfer: %s", transferStr)
			}
		}

		//If we are here, we have a match and we can proceed to check the restrictions
		//We have to satisfy at least one of the approval details in full to be allowed
		//We scan linearly through them
		//Note that we do not overflow into the next. Each must match in full
		if transferVal.ApprovalDetails != nil {
			approvalDetails := transferVal.ApprovalDetails
			if approvalDetails.RequireFromDoesNotEqualInitiatedBy && fromAddress == initiatedBy {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because from == initiatedBy: %s", transferStr)
			}

			if approvalDetails.RequireFromEqualsInitiatedBy && fromAddress != initiatedBy {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because from != initiatedBy: %s", transferStr)
			}

			if approvalDetails.RequireToDoesNotEqualInitiatedBy && toAddress == initiatedBy {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because to == initiatedBy: %s", transferStr)
			}

			if approvalDetails.RequireToEqualsInitiatedBy && toAddress != initiatedBy {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because to != initiatedBy: %s", transferStr)
			}

			//Assert that initiatedBy owns the required badges
			failedMustOwnBadges := false
			for _, mustOwnBadge := range approvalDetails.MustOwnBadges {
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

				fetchedBalances, err := types.GetBalancesForIds(mustOwnBadge.BadgeIds, mustOwnBadge.OwnershipTimes, balances)
				if err != nil {
					failedMustOwnBadges = true
					break
					// continue
					// return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for mustOwnBadges: %s", transferStr)
				}

				ownedOne := false
				for _, fetchedBalance := range fetchedBalances {
					//check if amount is within range
					minAmount := mustOwnBadge.AmountRange.Start
					maxAmount := mustOwnBadge.AmountRange.End

					if fetchedBalance.Amount.LT(minAmount) || fetchedBalance.Amount.GT(maxAmount) {
						failedMustOwnBadges = true
						
					} else if !mustOwnBadge.MustOwnAll {
						ownedOne = true
					}
					
				}

				if !mustOwnBadge.MustOwnAll && ownedOne {
					failedMustOwnBadges = false
					break
				}
			}

			if failedMustOwnBadges {
				continue
			}

			//Get max balances allowed for this approvalDetails element
			//Get the max balances allowed for this approvalDetails element WITHOUT incrementing
			transferBalancesToCheck, err := types.GetBalancesForIds(transferVal.BadgeIds, transferVal.OwnershipTimes, remainingBalances)
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
					approvalDetails.MerkleChallenge,
				}, solutions, initiatedBy, true, approverAddress, approvalLevel, challengeIdsIncremented)
			if err != nil {
				continue
				// return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "%s", transferStr)
			}

			//here, we assert the transfer is good for each level of approvals and increment if necessary
			maxPossible, err :=  k.GetMaxPossible(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.OverallApprovalAmount, approvalDetails.MaxNumTransfers.OverallMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "overall", "", true, trackerIdsIncremented)
			if err != nil {
				continue
			}
			transferBalancesToCheck, err = types.GetOverlappingBalances(maxPossible, transferBalancesToCheck)
			if err != nil {
				continue
			}
			
			maxPossible, err = k.GetMaxPossible(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerToAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerToAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "to", toAddress, true, trackerIdsIncremented)
			if err != nil {
				continue
			}
			transferBalancesToCheck, err = types.GetOverlappingBalances(maxPossible, transferBalancesToCheck)
			if err != nil {
				continue
			}
			

			maxPossible, err = k.GetMaxPossible(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerFromAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerFromAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "from", fromAddress, true, trackerIdsIncremented)
			if err != nil {
				continue
			}
			transferBalancesToCheck, err = types.GetOverlappingBalances(maxPossible, transferBalancesToCheck)
			if err != nil {
				continue
			}
			

			maxPossible, err = k.GetMaxPossible(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerInitiatedByAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "initiatedBy", initiatedBy, true, trackerIdsIncremented)
			if err != nil {
				continue
			}
			transferBalancesToCheck, err = types.GetOverlappingBalances(maxPossible, transferBalancesToCheck)
			if err != nil {
				continue
			}

			transferBalancesToCheck = types.FilterZeroBalances(transferBalancesToCheck)
			if len(transferBalancesToCheck) == 0 {
				continue
			}


			//here, we assert the transfer is good for each level of approvals and increment if necessary
			err =  k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.OverallApprovalAmount, approvalDetails.MaxNumTransfers.OverallMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "overall", "", true, trackerIdsIncremented)
			if err != nil {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded overall approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerToAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerToAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "to", toAddress, true, trackerIdsIncremented)
			if err != nil {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded to approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerFromAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerFromAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "from", fromAddress, true, trackerIdsIncremented)
			if err != nil {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded from approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerInitiatedByAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "initiatedBy", initiatedBy, true, trackerIdsIncremented)
			if err != nil {
				continue
				//return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded initiatedBy approvals: %s", transferStr)
			}


			//Finally, increment everything in store 
			remainingBalances, err = types.SubtractBalances(transferBalancesToCheck, remainingBalances)
			if err != nil {
				continue
			}


			//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
			//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
			//    If so, useLeafIndexForNumIncrements will be true 
			challengeNumIncrements, err = k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, 
				transferVal.ChallengeTrackerId,
				[]*types.MerkleChallenge{approvalDetails.MerkleChallenge}, solutions, initiatedBy, false, approverAddress, approvalLevel, challengeIdsIncremented)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "%s", transferStr)
			}

			//here, we assert the transfer is good for each level of approvals and increment if necessary
			err =  k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.OverallApprovalAmount, approvalDetails.MaxNumTransfers.OverallMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "overall", "", false, trackerIdsIncremented)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded overall approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerToAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerToAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "to", toAddress, false, trackerIdsIncremented)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded to approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerFromAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerFromAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "from", fromAddress, false, trackerIdsIncremented)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded from approvals: %s", transferStr)
			}

			err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerInitiatedByAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "initiatedBy", initiatedBy, false, trackerIdsIncremented)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "exceeded initiatedBy approvals: %s", transferStr)
			}

			//If we do not override the approved outgoing / incoming transfers, we need to check the user approvals
			if !approvalDetails.OverridesFromApprovedOutgoingTransfers {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address:  fromAddress,
					Balances: transferBalancesToCheck,
					Outgoing: true,
				})
			}

			if !approvalDetails.OverridesToApprovedIncomingTransfers {
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
		transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + remainingBalances[0].BadgeIds[0].Start.String() + ", ownershipTime (unix milliseconds): " + remainingBalances[0].OwnershipTimes[0].Start.String() + ")"

		return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrInadequateApprovals, "no approval found: %s", transferStr)
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
	transferVal *types.CollectionApprovedTransfer,
	approvalDetails *types.ApprovalDetails,
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
	approvalTrackerId := transferVal.ApprovalTrackerId
	predeterminedBalances := approvalDetails.PredeterminedBalances
	allApprovals := []*types.Balance{{
		Amount: approvedAmount,
		OwnershipTimes: transferVal.OwnershipTimes,
		BadgeIds: transferVal.BadgeIds,
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

	approvalTrackerDetails := types.ApprovalsTracker{
		Amounts:      []*types.Balance{},
		NumTransfers: sdkmath.NewUint(0),
	}

	if needToFetchApprovalTrackerDetails {
		fetchedDetails, found := k.GetApprovalsTrackerFromStore(ctx, collection.CollectionId, approverAddress, approvalTrackerId, approvalLevel, trackerType, address)
		if found {
			approvalTrackerDetails = fetchedDetails
		}
	}

	//Increment amounts and numTransfers and add back to store
	if approvedAmount.GT(sdkmath.NewUint(0)) {
		//Assume that if approvalTrackerDetails.Amounts is already not nil, it is correct and has been incremented properly
		//Here, we ONLY check if the NEW transferBalances makes it exceed the threshold
		currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(transferVal.BadgeIds, transferVal.OwnershipTimes, approvalTrackerDetails.Amounts)
		if err != nil {
			return nil, err
		}

		maxBalancesWeCanAdd, err := types.SubtractBalances(currTallyForCurrentIdsAndTimes, allApprovals)
		if err != nil {
			return nil, err
		}

		return maxBalancesWeCanAdd, nil
	}

	return transferBalances, nil
}

func (k Keeper) IncrementApprovalsAndAssertWithinThreshold(
	ctx sdk.Context,
	transferVal *types.CollectionApprovedTransfer,
	approvalDetails *types.ApprovalDetails,
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
) (error) {
	approvalTrackerId := transferVal.ApprovalTrackerId
	predeterminedBalances := approvalDetails.PredeterminedBalances
	allApprovals := []*types.Balance{{
		Amount: approvedAmount,
		OwnershipTimes: transferVal.OwnershipTimes,
		BadgeIds: transferVal.BadgeIds,
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

	approvalTrackerDetails := types.ApprovalsTracker{
		Amounts:      []*types.Balance{},
		NumTransfers: sdkmath.NewUint(0),
	}

	trackerIncrementId := fmt.Sprintf("%s-%s-%s-%s-%s-%s", approvalTrackerId, approverAddress, approvalLevel, trackerType, address, collection.CollectionId)
	if needToFetchApprovalTrackerDetails {
		fetchedDetails, found := k.GetApprovalsTrackerFromStore(ctx, collection.CollectionId, approverAddress, approvalTrackerId, approvalLevel, trackerType, address)
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
			// 	fetchedBalancesOfOverall, err := types.GetBalancesForIds(calculatedBalance.BadgeIds, calculatedBalance.OwnershipTimes, overallTransferBalances)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	fetchedBalances, err  = types.AddBalances(fetchedBalances, fetchedBalancesOfOverall)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
			
			//Assert that we have exactly the amount specified in the original transfers
			equal := types.AreBalancesEqual(overallTransferBalances, calculatedBalances, false)
			if !equal {
				return sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because predetermined balances do not match: %s", approvalTrackerId)
			}
		}
	}

	//Increment amounts and numTransfers and add back to store
	if approvedAmount.GT(sdkmath.NewUint(0)) {
		//Assume that if approvalTrackerDetails.Amounts is already not nil, it is correct and has been incremented properly
		//Here, we ONLY check if the NEW transferBalances makes it exceed the threshold
		currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(transferVal.BadgeIds, transferVal.OwnershipTimes, approvalTrackerDetails.Amounts)
		if err != nil {
			return err
		}

		//If this passes, the new transferBalances are okay
		_, err = types.AddBalancesAndAssertDoesntExceedThreshold(currTallyForCurrentIdsAndTimes, transferBalances, allApprovals)
		if err != nil {
			return err
		}

		//We then add them to the current tally of ALL ids and times
		approvalTrackerDetails.Amounts, err = types.AddBalances(approvalTrackerDetails.Amounts, transferBalances)
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
			sdk.NewEvent("approval" + fmt.Sprint(collection.CollectionId) + fmt.Sprint(approverAddress) + fmt.Sprint(approvalTrackerId) + fmt.Sprint(approvalLevel) + fmt.Sprint(trackerType) + fmt.Sprint(address),
				sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
				sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
				sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
				sdk.NewAttribute("approvalTrackerId", fmt.Sprint(approvalTrackerId)),
				sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
				sdk.NewAttribute("trackerType", fmt.Sprint(trackerType)),
				sdk.NewAttribute("approvedAddress", fmt.Sprint(address)),
				sdk.NewAttribute("amounts", amountsStr),
				sdk.NewAttribute("numTransfers", numTransfersStr),
			),
		)

		err = k.SetApprovalsTrackerInStore(ctx, collection.CollectionId, approverAddress, approvalTrackerId, approvalTrackerDetails, approvalLevel, trackerType, address)
		if err != nil {
			return err
		}
	}

	return nil
}



func (k Keeper) GetPredeterminedBalancesForPrecalculationId(ctx sdk.Context, approvedTransfers []*types.CollectionApprovedTransfer, collection *types.BadgeCollection, approverAddress string, precalculationId string, approvalLevel string, address string, solutions []*types.MerkleProof, initiatedBy string) ([]*types.Balance, error) {
	approvalId := ""
	for _, transfer := range approvedTransfers {
		approvalDetails := transfer.ApprovalDetails
		approvalId = transfer.ApprovalId
		approvalTrackerId := transfer.ApprovalTrackerId
		if approvalDetails == nil {
			continue
		}

		if approvalId != precalculationId {
			continue
		}

		if approvalId == "" || approvalTrackerId == "" {
			continue
		}

		if approvalDetails.PredeterminedBalances != nil {
				numIncrements := sdkmath.NewUint(0)
				if approvalDetails.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {

					//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
					//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
					numIncrementsFetched, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId,
						transfer.ChallengeTrackerId,
						[]*types.MerkleChallenge{
						approvalDetails.MerkleChallenge,
					}, solutions, initiatedBy, true,  address, approvalLevel,  &[]string{})
					if err != nil {
						return []*types.Balance{}, sdkerrors.Wrapf(err, "invalid challenges / solutions")
					}

					numIncrements = numIncrementsFetched
				} else {
					trackerType := "overall"
					if approvalDetails.PredeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers {
						trackerType = "from"
					} else if approvalDetails.PredeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers {
						trackerType = "to"
					} else if approvalDetails.PredeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers {
						trackerType = "initiatedBy"
					}

					approvalTrackerDetails, found := k.GetApprovalsTrackerFromStore(ctx, collection.CollectionId, approverAddress, approvalTrackerId, approvalLevel, trackerType, address)
					if !found {
						approvalTrackerDetails = types.ApprovalsTracker{
							Amounts:      []*types.Balance{},
							NumTransfers: sdkmath.NewUint(0),
						}
					}

					numIncrements = approvalTrackerDetails.NumTransfers
				}

				predeterminedBalances := []*types.Balance{}
				if approvalDetails.PredeterminedBalances.ManualBalances != nil {
					if numIncrements.LT(sdkmath.NewUint(uint64(len(approvalDetails.PredeterminedBalances.ManualBalances)))) {
						predeterminedBalances = types.DeepCopyBalances(approvalDetails.PredeterminedBalances.ManualBalances[numIncrements.Uint64()].Balances)
					}
				} else if approvalDetails.PredeterminedBalances.IncrementedBalances != nil {
					err := *new(error)
					predeterminedBalances, err = IncrementBalances(approvalDetails.PredeterminedBalances.IncrementedBalances.StartBalances, numIncrements, approvalDetails.PredeterminedBalances.IncrementedBalances.IncrementOwnershipTimesBy, approvalDetails.PredeterminedBalances.IncrementedBalances.IncrementBadgeIdsBy)
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