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
func (k Keeper) DeductUserOutgoingApprovals(ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.MerkleProof) error {
	currApprovedTransfers := types.GetCurrentUserApprovedOutgoingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForOutgoing(currApprovedTransfers, from)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastOutgoingTransfersToCollectionTransfers(currApprovedTransfers, from)
	_, err := k.DeductAndGetUserApprovals(overallTransferBalances, castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "outgoing", from)
	return err
}

// DeductUserIncomingApprovals will check if the current transfer is approved from the to's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserIncomingApprovals(ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.MerkleProof) error {
	currApprovedTransfers := types.GetCurrentUserApprovedIncomingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForIncoming(currApprovedTransfers, to)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastIncomingTransfersToCollectionTransfers(currApprovedTransfers, to)
	_, err := k.DeductAndGetUserApprovals(overallTransferBalances, castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "incoming", to)
	return err
}

// DeductCollectionApprovalsAndGetUserApprovalsToCheck will check if the current transfer is allowed via the collection's approved transfers and handle any tallying accordingly
func (k Keeper) DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.MerkleProof) ([]*UserApprovalsToCheck, error) {
	approvedTransfers := types.GetCurrentCollectionApprovedTransfers(ctx, collection)
	return k.DeductAndGetUserApprovals(overallTransferBalances, approvedTransfers, ctx, collection, badgeIds, times, fromAddress, toAddress, initiatedBy, amount, solutions, "collection", "")
}

func (k Keeper) DeductAndGetUserApprovals(overallTransferBalances []*types.Balance, approvedTransfers []*types.CollectionApprovedTransfer, ctx sdk.Context, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.MerkleProof, approvalLevel string, approverAddress string) ([]*UserApprovalsToCheck, error) {
	//HACK: We first expand all transfers to have just a len == 1 AllowedCombination[] so that we can easily check IsApproved later
	//		  This is because GetFirstMatchOnly will break down the transfers into smaller parts and without expansion, fetching if a certain transfer is allowed is impossible.
	expandedApprovedTransfers := ExpandCollectionApprovedTransfers(approvedTransfers)
	manager := types.GetCurrentManager(ctx, collection)
	castedApprovedTransfers, err := k.CastCollectionApprovedTransferToUniversalPermission(ctx, expandedApprovedTransfers, manager)
	if err != nil {
		return []*UserApprovalsToCheck{}, err
	}

	firstMatches := types.GetFirstMatchOnly(castedApprovedTransfers)

	//Keep a running tally of all the badges we still have to handle
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
			})
		}
	}

	//Keep a list of all the user incoming/outgoing approvals we need to check
	userApprovalsToCheck := []*UserApprovalsToCheck{}
	for _, match := range firstMatches {
		transferVal := match.ArbitraryValue.(*types.CollectionApprovedTransfer)

		doAddressesMatch := k.CheckIfAddressesMatchCollectionMappingIds(ctx, transferVal, fromAddress, toAddress, initiatedBy, manager)
		if !doAddressesMatch {
			// return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "addresses do not match for transfer: (from: %s, to: %s, initiatedBy: %s)", fromAddress, toAddress, initiatedBy)
			continue
		}

		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currTimeFound := types.SearchUintRangesForUint(currTime, []*types.UintRange{match.TransferTime})
		if !currTimeFound {
			// return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "current time not found in transfer time: %s", currTime.String())
			continue
		}
		
		remaining, overlaps := types.UniversalRemoveOverlapFromValues(&types.UniversalPermissionDetails{
			BadgeId:       match.BadgeId,
			OwnershipTime: match.OwnershipTime,
			TimelineTime: &types.UintRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
			TransferTime: &types.UintRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
			ToMapping: &types.AddressMapping{},
			FromMapping: &types.AddressMapping{},
			InitiatedByMapping: &types.AddressMapping{},
		}, unhandled)

		
		unhandled = remaining
		for _, overlap := range overlaps {
			//For the overlapping badges and ownership times, we have a match because mapping IDs, time, and badge IDs match.
			//We can now proceed to check any restrictions.
			//If any restriction fails in this if statement, we MUST throw because the transfer is invalid for the badge IDs (since we use first match only)

			transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + overlap.BadgeId.Start.String() + ", time: " + currTime.String() + ", ownershipTime: " + overlap.OwnershipTime.Start.String() + ")"
			// return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "overlap transfer disallowed because no approved transfer was found for badge id: %s, %s", overlap.BadgeId.Start, transferStr)


			allowed := transferVal.AllowedCombinations[0].IsApproved //HACK: can do this because we expanded the allowed combinations above
			if !allowed {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed explicitly: %s", transferStr)
			}

			for _, approvalDetails := range transferVal.ApprovalDetails {
				if approvalDetails.RequireFromDoesNotEqualInitiatedBy && fromAddress == initiatedBy {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because from == initiatedBy: %s", transferStr)
				}

				if approvalDetails.RequireFromEqualsInitiatedBy && fromAddress != initiatedBy {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because from != initiatedBy: %s", transferStr)
				}

				if approvalDetails.RequireToDoesNotEqualInitiatedBy && toAddress == initiatedBy {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because to == initiatedBy: %s", transferStr)
				}

				if approvalDetails.RequireToEqualsInitiatedBy && toAddress != initiatedBy {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because to != initiatedBy: %s", transferStr)
				}

				//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
				//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
				//    If so, useLeafIndexForNumIncrements will be true 
				challengeNumIncrements, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, approvalDetails.MerkleChallenges, solutions, initiatedBy, false, approverAddress, approvalLevel)
				if err != nil {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed because of invalid challenges / solutions: %s", transferStr)
				}

				//TODO: Support inherited balances
				//Assert that initiatedBy owns the required badges
				for _, mustOwnBadge := range approvalDetails.MustOwnBadges {
					initiatedByBalanceKey := ConstructBalanceKey(initiatedBy, mustOwnBadge.CollectionId)
					initiatedByBalance, found := k.GetUserBalanceFromStore(ctx, initiatedByBalanceKey)
					balances := []*types.Balance{}
					if found {
						balances = initiatedByBalance.Balances
					}

					if mustOwnBadge.OverrideWithCurrentTime {
						mustOwnBadge.OwnershipTimes = []*types.UintRange{{Start: currTime, End: currTime}}
					}

					fetchedBalances, err := types.GetBalancesForIds(mustOwnBadge.BadgeIds, mustOwnBadge.OwnershipTimes, balances)
					if err != nil {
						return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for mustOwnBadges: %s", transferStr)
					}

					for _, fetchedBalance := range fetchedBalances {
						//check if amount is within range
						minAmount := mustOwnBadge.AmountRange.Start
						maxAmount := mustOwnBadge.AmountRange.End

						if fetchedBalance.Amount.LT(minAmount) || fetchedBalance.Amount.GT(maxAmount) {
							return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because user does not own the required badges in mustOwnBadges: %s", transferStr)
						}
					}
				}

				//transferBalances is the current balances we are checking if we can transfer
				transferBalancesToCheck := []*types.Balance{{Amount: amount, OwnershipTimes: []*types.UintRange{overlap.OwnershipTime}, BadgeIds: []*types.UintRange{overlap.BadgeId}}}

				//here, we assert the transfer is good for each level of approvals and increment if necessary
				err =  k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.OverallApprovalAmount, approvalDetails.MaxNumTransfers.OverallMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "overall", "")
				if err != nil {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing overall approvals: %s", transferStr)
				}

				err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerToAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerToAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "to", toAddress)
				if err != nil {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing to approvals: %s", transferStr)
				}

				err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerFromAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerFromAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "from", fromAddress)
				if err != nil {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing from approvals: %s", transferStr)
				}

				err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, transferVal, approvalDetails, overallTransferBalances, collection, approvalDetails.ApprovalAmounts.PerInitiatedByAddressApprovalAmount, approvalDetails.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers, transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, "initiatedBy", initiatedBy)
				if err != nil {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing initiatedBy approvals: %s", transferStr)
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
	}

	//If all are not explicitly allowed, we return that it is disallowed by default
	if len(unhandled) > 0 {
		transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + unhandled[0].BadgeId.Start.String() + ", ownershipTime: " + unhandled[0].OwnershipTime.Start.String() + ")"

		return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrInadequateApprovals, "transfer disallowed because no approved transfer was found for: %s", transferStr)
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
) (error) {
	approvalId := approvalDetails.ApprovalId
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
	if needToFetchApprovalTrackerDetails {
		fetchedDetails, found := k.GetApprovalsTrackerFromStore(ctx, collection.CollectionId, approverAddress, approvalId, approvalLevel, trackerType, address)
		if found {
			approvalTrackerDetails = fetchedDetails
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
		if predeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {
			numIncrements = challengeNumIncrements
		} else if predeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers && trackerType == "overall" {
			numIncrements = approvalTrackerDetails.NumTransfers
		} else if predeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to" {
			numIncrements = approvalTrackerDetails.NumTransfers
		} else if predeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from" {
			numIncrements = approvalTrackerDetails.NumTransfers
		} else if predeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy" {
			numIncrements = approvalTrackerDetails.NumTransfers
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
		
			//From the original transfer balances (the overall ones, not overlap), fetch the balances for the badge IDs and times of the current transfer
			//Filter out any balances that have an amount of 0
			fetchedBalances, err := types.GetBalancesForIds(allApprovals[0].BadgeIds, allApprovals[0].OwnershipTimes, overallTransferBalances)
			if err != nil {
				return err
			}

			//Assert that we have exactly the amount specified in the original transfers
			//This also asserts that calculated balances does not specify any out of bounds badge IDs or times
			equal := types.AreBalancesEqual(fetchedBalances, calculatedBalances, false)
			if !equal {
				return sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because predetermined balances do not match: %s", approvalId)
			}
		}
	}

	
	//Increment amounts and numTransfers and add back to store
	if approvedAmount.GT(sdkmath.NewUint(0)) {
		approvalTrackerDetails.Amounts, err = types.AddBalancesAndAssertDoesntExceedThreshold(approvalTrackerDetails.Amounts, transferBalances, allApprovals)
		if err != nil {
			return err
		}
	}

	if maxNumTransfers.GT(sdkmath.NewUint(0)) ||
	(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers && trackerType == "overall") ||
	(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to") ||
	(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from") ||
	(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy") {
		approvalTrackerDetails.NumTransfers = approvalTrackerDetails.NumTransfers.Add(sdkmath.NewUint(1))
		//only check exceeds if maxNumTransfers is not 0 (because 0 means no limit)
		if maxNumTransfers.GT(sdkmath.NewUint(0)) {
			if approvalTrackerDetails.NumTransfers.GT(maxNumTransfers) {
				return sdkerrors.Wrapf(ErrDisallowedTransfer, "exceeded max num transfers - %s", maxNumTransfers.String())	
			}
		}
	}

	if needToFetchApprovalTrackerDetails {
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

		ctx.EventManager().EmitEvent(
			sdk.NewEvent("approval" + fmt.Sprint(collection.CollectionId) + fmt.Sprint(approverAddress) + fmt.Sprint(approvalId) + fmt.Sprint(approvalLevel) + fmt.Sprint(trackerType) + fmt.Sprint(address),
				sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
				sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
				sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
				sdk.NewAttribute("approvalId", fmt.Sprint(approvalId)),
				sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
				sdk.NewAttribute("trackerType", fmt.Sprint(trackerType)),
				sdk.NewAttribute("approvedAddress", fmt.Sprint(address)),
				sdk.NewAttribute("amounts", amountsStr),
				sdk.NewAttribute("numTransfers", numTransfersStr),
			),
		)

		err = k.SetApprovalsTrackerInStore(ctx, collection.CollectionId, approverAddress, approvalId, approvalTrackerDetails, approvalLevel, trackerType, address)
		if err != nil {
			return err
		}
	}

	return nil
}



func (k Keeper) GetPredeterminedBalancesForApprovalId(ctx sdk.Context, approvedTransfers []*types.CollectionApprovedTransfer, collection *types.BadgeCollection, approverAddress string, approvalId string, approvalLevel string, address string, solutions []*types.MerkleProof, initiatedBy string) ([]*types.Balance, error) {
	for _, transfer := range approvedTransfers {
		for _, approvalDetails := range transfer.ApprovalDetails {
			if approvalDetails.ApprovalId == approvalId {
				if approvalDetails.PredeterminedBalances != nil {
					numIncrements := sdkmath.NewUint(0)
					if approvalDetails.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {
						//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
						//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
						numIncrementsFetched, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, approvalDetails.MerkleChallenges, solutions, initiatedBy, true,  address, approvalLevel)
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

						approvalTrackerDetails, found := k.GetApprovalsTrackerFromStore(ctx, collection.CollectionId, approverAddress, approvalId, approvalLevel, trackerType, address)
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
					return []*types.Balance{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "no predetermined transfers found for approval id: %s", approvalId)
				}
			}
		}
	}

	return []*types.Balance{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "no predetermined transfers found for approval id: %s", approvalId)
}