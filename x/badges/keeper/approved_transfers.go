package keeper

import (
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
func (k Keeper) DeductUserOutgoingApprovals(ctx sdk.Context, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.ChallengeSolution, transferAsMuchAsPossible bool) error {
	currApprovedTransfers := types.GetCurrentUserApprovedOutgoingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForOutgoing(currApprovedTransfers, from)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastOutgoingTransfersToCollectionTransfers(currApprovedTransfers, from)
	_, err := k.DeductAndGetUserApprovals(castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "outgoing", transferAsMuchAsPossible)
	return err
}

// DeductUserIncomingApprovals will check if the current transfer is approved from the to's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserIncomingApprovals(ctx sdk.Context, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.ChallengeSolution, transferAsMuchAsPossible bool) error {
	currApprovedTransfers := types.GetCurrentUserApprovedIncomingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForIncoming(currApprovedTransfers, to)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastIncomingTransfersToCollectionTransfers(currApprovedTransfers, to)
	_, err := k.DeductAndGetUserApprovals(castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "incoming", transferAsMuchAsPossible)
	return err
}

// DeductCollectionApprovalsAndGetUserApprovalsToCheck will check if the current transfer is allowed via the collection's approved transfers and handle any tallying accordingly
func (k Keeper) DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx sdk.Context, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.ChallengeSolution, transferAsMuchAsPossible bool) ([]*UserApprovalsToCheck, error) {
	approvedTransfers := types.GetCurrentCollectionApprovedTransfers(ctx, collection)
	return k.DeductAndGetUserApprovals(approvedTransfers, ctx, collection, badgeIds, times, fromAddress, toAddress, initiatedBy, amount, solutions, "overall", transferAsMuchAsPossible)
}

func (k Keeper) DeductAndGetUserApprovals(approvedTransfers []*types.CollectionApprovedTransfer, ctx sdk.Context, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.ChallengeSolution, timelineType string, transferAsMuchAsPossible bool) ([]*UserApprovalsToCheck, error) {
	//HACK: We first expand all transfers to have just a len == 1 AllowedCombination[] so that we can easily check IsAllowed later
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
			continue
		}

		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currTimeFound := types.SearchUintRangesForUint(currTime, []*types.UintRange{match.TransferTime})
		if !currTimeFound {
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

		// remaining, overlaps := types.RemoveUintRangeFromUintRange([]*types.UintRange{match.BadgeId}, unhandledBadgeIds)
		// timesRemaining, overlapsTimes := types.RemoveUintRangeFromUintRange([]*types.UintRange{match.OwnershipTime}, unhandledBadgeIds)
		
		unhandled = remaining

		for _, overlap := range overlaps {
			//For the overlapping badges and ownership times, we have a match because mapping IDs, time, and badge IDs match.
			//We can now proceed to check any restrictions.
			//If any restriction fails in this if statement, we MUST throw because the transfer is invalid for the badge IDs (since we use first match only)

			transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + overlap.BadgeId.Start.String() + ", time: " + currTime.String() + ", ownershipTime: " + overlap.OwnershipTime.Start.String() + ")"

			allowed := transferVal.AllowedCombinations[0].IsAllowed //HACK: can do this because we expanded the allowed combinations above
			if !allowed {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed explicitly: %s", transferStr)
			}

			if transferVal.RequireFromDoesNotEqualInitiatedBy && fromAddress == initiatedBy {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because from == initiatedBy: %s", transferStr)
			}

			if transferVal.RequireFromEqualsInitiatedBy && fromAddress != initiatedBy {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because from != initiatedBy: %s", transferStr)
			}

			if transferVal.RequireToDoesNotEqualInitiatedBy && toAddress == initiatedBy {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because to == initiatedBy: %s", transferStr)
			}

			if transferVal.RequireToEqualsInitiatedBy && toAddress != initiatedBy {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because to != initiatedBy: %s", transferStr)
			}

			//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
			//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
			useLeafIndexForNumIncrements, numIncrements, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, transferVal.Challenges, solutions, initiatedBy, "overall", transferVal.ApprovalId, false)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed because of invalid challenges / solutions: %s", transferStr)
			}

			transferBalances := []*types.Balance{{Amount: amount, OwnershipTimes: []*types.UintRange{overlap.OwnershipTime}, BadgeIds: []*types.UintRange{overlap.BadgeId}}}

			if transferVal.ApprovalId == "" && (transferVal.OverallApprovals != nil || transferVal.PerAddressApprovals != nil) {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because approvalId cannot be blank: %s", transferStr)
			}

			//Handle the incrementing of overall approvals
			if transferVal.OverallApprovals != nil {
				_, err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.OverallApprovals.Amounts, transferVal.OverallApprovals.NumTransfers, transferBalances, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "overall", "", false)
				if err != nil {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing overall approvals: %s", transferStr)
				}
			}

			//Handle the per-address approvals
			if transferVal.PerAddressApprovals != nil {
				if transferVal.PerAddressApprovals.ApprovalsPerToAddress != nil {
					_, err := k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.PerAddressApprovals.ApprovalsPerToAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers, transferBalances, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "per-to", toAddress, false)
					if err != nil {
						return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing per-to approvals: %s", transferStr)
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerFromAddress != nil {
					_, err := k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers, transferBalances, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "per-from", fromAddress, false)
					if err != nil {
						return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing per-from approvals: %s", transferStr)
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil {
					_, err := k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress.NumTransfers, transferBalances, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "per-initiated-by", initiatedBy, false)
					if err != nil {
						return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing per-initiated-by approvals: %s", transferStr)
					}
				}
			}

			//If we are overriding the approved outgoing / incoming transfers, we don't need to check the user approvals
			//Else, we do
			if !transferVal.OverridesFromApprovedOutgoingTransfers {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address:  fromAddress,
					Balances: transferBalances,
					Outgoing: true,
				})
			}

			if !transferVal.OverridesToApprovedIncomingTransfers {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address:  toAddress,
					Balances: transferBalances,
					Outgoing: false,
				})
			}
		}
	}

	//If all are not explicitly allowed, we return that it is disallowed by default
	if len(unhandled) > 0 {
		return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrInadequateApprovals, "transfer disallowed because no approved transfer was found for badge ids: %v", unhandled)
	}

	return userApprovalsToCheck, nil
}

func AssertBalancesDoNotExceedThreshold(balancesToCheck []*types.Balance, threshold []*types.Balance) error {
	err := *new(error)

	
	//Check if we exceed the threshold; will underflow if we do exceed it
	thresholdCopy := types.DeepCopyBalances(threshold)
	for _, balance := range balancesToCheck {
		thresholdCopy, err = types.SubtractBalance(thresholdCopy, balance)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddTallyAndAssertDoesntExceedThreshold(currTally []*types.Balance, toAdd []*types.Balance, threshold []*types.Balance, transferAsMuchAsPossible bool) ([]*types.Balance, []*types.Balance, error) {
	err := *new(error)
	//If we transferAsMuchAsPossible, we need to increment the currTally by all that we can
	//We then need to return the updated toAdd


	for _, balance := range toAdd {
		//Add the new tally to existing
		currTally, err = types.AddBalance(currTally, balance)
		if err != nil {
			return []*types.Balance{}, []*types.Balance{}, err
		}
	}

	
	castedCurrTally := []*types.UniversalPermissionDetails{}
	for _, balance := range currTally {
		for _, badgeId := range balance.BadgeIds {
			for _, time := range balance.OwnershipTimes {
				castedCurrTally = append(castedCurrTally, &types.UniversalPermissionDetails{
					BadgeId:            badgeId,
					OwnershipTime:       time,
					TransferTime:       &types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					TimelineTime: 		 	&types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					ToMapping:          &types.AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					FromMapping:        &types.AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					InitiatedByMapping: &types.AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					ArbitraryValue:     balance.Amount,
				})
			}
		}
	}

	castedThreshold := []*types.UniversalPermissionDetails{}
	for _, balance := range threshold {
		for _, badgeId := range balance.BadgeIds {
			for _, time := range balance.OwnershipTimes {
				castedThreshold = append(castedThreshold, &types.UniversalPermissionDetails{
					BadgeId:            badgeId,
					OwnershipTime:       time,
					TransferTime:       &types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					TimelineTime: 		 &types.UintRange{Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64)}, //dummy range
					ToMapping:          &types.AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					FromMapping:        &types.AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					InitiatedByMapping: &types.AddressMapping{Addresses: []string{}, IncludeAddresses: false},
					ArbitraryValue:     balance.Amount,
				})
			}
		}
	}

	overlaps, inTallyButNotThreshold, _ := types.GetOverlapsAndNonOverlaps(castedCurrTally, castedThreshold)
	if transferAsMuchAsPossible {
		for _, overlapObject := range overlaps {
			tallyAmount := overlapObject.FirstDetails.ArbitraryValue.(sdkmath.Uint)
			thresholdAmount := overlapObject.SecondDetails.ArbitraryValue.(sdkmath.Uint)

			if thresholdAmount.LT(tallyAmount) {
				//we overflowed so we need to reduce the toAdd by the difference
				toAdd, err = types.SubtractBalance(toAdd, &types.Balance{
					Amount: 			 tallyAmount.Sub(thresholdAmount),
					BadgeIds:       []*types.UintRange{overlapObject.Overlap.BadgeId},
					OwnershipTimes: []*types.UintRange{overlapObject.Overlap.OwnershipTime},
				})
				if err != nil {
					return []*types.Balance{}, []*types.Balance{}, err
				}
			}
		}

		for _, details := range inTallyButNotThreshold {
			toAdd, err = types.DeleteBalances([]*types.UintRange{details.BadgeId}, []*types.UintRange{details.OwnershipTime}, toAdd)
			if err != nil {
				return []*types.Balance{}, []*types.Balance{}, err
			}
		}
		
		return currTally, toAdd, nil
	}


	//HACK: Because threshold can change with increments, we assume that the existing tally is already valid because it was fetched from store
	//		  We then remove any badges that are not in the toAdd balances, so we are only checking if the current transfer exceeds the threshold
	// 			Ex: for an increment-based approval, the curr tally might be 1-1 (first transfer) and the toAdd might be 2-2 (second transfer) with increments
	//				  the threshold would then change from 1-1 to 2-2 because the number of increments changed. however, if we do not remove the 1-1 from the tally,
	//					 we will incorrectly throw an error because the tally would be 1-1, 2-2, which exceeds the threshold of 2-2 
	currTallyWithNonAddedBadgesRemoved := []*types.Balance{}
	for _, balance := range toAdd {
		balances, err := types.GetBalancesForIds(balance.BadgeIds, balance.OwnershipTimes, currTally)
		if err != nil {
			return []*types.Balance{}, []*types.Balance{}, err
		}

		currTallyWithNonAddedBadgesRemoved = append(currTallyWithNonAddedBadgesRemoved, balances...)
	}


	//Check if we exceed the threshold; will underflow if we do exceed it
	err = AssertBalancesDoNotExceedThreshold(currTallyWithNonAddedBadgesRemoved, threshold)
	return currTally, toAdd, err
}

func (k Keeper) IncrementApprovalsAndAssertWithinThreshold(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	approvals []*types.Balance,
	maxNumTransfers sdkmath.Uint,
	transferBalances []*types.Balance,
	approvalId string,
	incrementBadgeIdsBy sdkmath.Uint,
	incrementOwnershipTimesBy sdkmath.Uint,
	precalculatedNumIncrements bool,
	numIncrements sdkmath.Uint,
	timelineType string,
	depth string,
	address string,
	transferAsMuchAsPossible bool,
) ([]*types.Balance, error) {
	//Get the current approvals for this transfer
	//If nil, no restrictions and we are approved for the entire transfer
	//Note we filter any excess badge IDs later and apply num increments as well
	err := *new(error)
	if approvals == nil {
		return nil, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because no approval amounts were found")
	}

	approvalTrackerDetails, found := k.GetApprovalsTrackerFromStore(ctx, collection.CollectionId, approvalId, timelineType, depth, address)
	if !found {
		approvalTrackerDetails = types.ApprovalsTracker{
			Amounts:      []*types.Balance{},
			NumTransfers: sdkmath.NewUint(0),
		}
	}

	if !precalculatedNumIncrements {
		numIncrements = approvalTrackerDetails.NumTransfers
	}

	//allApprovals is the total amount approved (i.e. the initial total amounts plus all increments)
	allApprovals := types.DeepCopyBalances(approvals)

	for _, startAmount := range allApprovals {
		for _, time := range startAmount.OwnershipTimes {
			time.Start = time.Start.Add(numIncrements.Mul(incrementOwnershipTimesBy))
			time.End = time.End.Add(numIncrements.Mul(incrementOwnershipTimesBy))
		}

		for _, badgeId := range startAmount.BadgeIds {
			badgeId.Start = badgeId.Start.Add(numIncrements.Mul(incrementBadgeIdsBy))
			badgeId.End = badgeId.End.Add(numIncrements.Mul(incrementBadgeIdsBy))
		}
	}

	//Increment amounts and numTransfers and add back to store
	approvalTrackerDetails.Amounts, transferBalances, err = AddTallyAndAssertDoesntExceedThreshold(approvalTrackerDetails.Amounts, transferBalances, allApprovals, transferAsMuchAsPossible)
	if err != nil {
		return transferBalances, err
	}

	approvalTrackerDetails.NumTransfers = approvalTrackerDetails.NumTransfers.Add(sdkmath.NewUint(1))
	if approvalTrackerDetails.NumTransfers.GT(maxNumTransfers) {
		if transferAsMuchAsPossible {
			return []*types.Balance{}, nil //can't transfer anything 
		} else {
			return transferBalances, sdkerrors.Wrapf(ErrDisallowedTransfer, "exceeded max num transfers - %s", maxNumTransfers.String())
		}
	}

	//Simulation so we do not want to actually increment
	if transferAsMuchAsPossible {
		return transferBalances, nil
	}

	err = k.SetApprovalsTrackerInStore(ctx, collection.CollectionId, approvalId, approvalTrackerDetails, timelineType, depth, address)
	if err != nil {
		return transferBalances, err
	}

	return transferBalances, nil
}
















func (k Keeper) GetUnapprovedBalancesForOutgoing(ctx sdk.Context, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, overallTransferBalances []*types.Balance, balancesToCheck []*types.Balance, from string, to string, requester string, solutions []*types.ChallengeSolution, transferAsMuchAsPossible bool) ([]*UserApprovalsToCheck, []*types.Balance, error) {
	currApprovedTransfers := types.GetCurrentUserApprovedOutgoingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForOutgoing(currApprovedTransfers, from)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastOutgoingTransfersToCollectionTransfers(currApprovedTransfers, from)
	return k.GetUnapprovedTransferableBalances(castedTransfers, ctx, collection, overallTransferBalances, userBalance.Balances, from, to, requester, solutions, "outgoing", transferAsMuchAsPossible)
}

func (k Keeper) GetUnapprovedBalancesForIncoming(ctx sdk.Context, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, overallTransferBalances []*types.Balance, balancesToCheck []*types.Balance, from string, to string, requester string, solutions []*types.ChallengeSolution, transferAsMuchAsPossible bool) ([]*UserApprovalsToCheck, []*types.Balance, error) {
	currApprovedTransfers := types.GetCurrentUserApprovedIncomingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForIncoming(currApprovedTransfers, to)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastIncomingTransfersToCollectionTransfers(currApprovedTransfers, to)
	return k.GetUnapprovedTransferableBalances(castedTransfers, ctx, collection, overallTransferBalances, userBalance.Balances, from, to, requester, solutions, "incoming", transferAsMuchAsPossible)
}

func (k Keeper) GetUnapprovedBalancesForCollection(ctx sdk.Context, collection *types.BadgeCollection, balances []*types.Balance, fromAddress string, toAddress string, initiatedBy string, solutions []*types.ChallengeSolution, transferAsMuchAsPossible bool) ([]*UserApprovalsToCheck, []*types.Balance, error) {
	approvedTransfers := types.GetCurrentCollectionApprovedTransfers(ctx, collection)
	balancesCopy := types.DeepCopyBalances(balances)
	return k.GetUnapprovedTransferableBalances(approvedTransfers, ctx, collection, balances, balancesCopy, fromAddress, toAddress, initiatedBy, solutions, "overall", transferAsMuchAsPossible)
}

func (k Keeper) GetUnapprovedTransferableBalances(approvedTransfers []*types.CollectionApprovedTransfer, ctx sdk.Context, collection *types.BadgeCollection, overallTransferBalances []*types.Balance, balancesToCheck []*types.Balance, fromAddress string, toAddress string, initiatedBy string, solutions []*types.ChallengeSolution, timelineType string, transferAsMuchAsPossible bool) ([]*UserApprovalsToCheck, []*types.Balance, error) {
	//HACK: We first expand all transfers to have just a len == 1 AllowedCombination[] so that we can easily check IsAllowed later
	//		  This is because GetFirstMatchOnly will break down the transfers into smaller parts and without expansion, fetching if a certain transfer is allowed is impossible.
	expandedApprovedTransfers := ExpandCollectionApprovedTransfers(approvedTransfers)
	manager := types.GetCurrentManager(ctx, collection)
	castedApprovedTransfers, err := k.CastCollectionApprovedTransferToUniversalPermission(ctx, expandedApprovedTransfers, manager)
	if err != nil {
		return []*UserApprovalsToCheck{}, []*types.Balance{}, err
	}

	firstMatches := types.GetFirstMatchOnly(castedApprovedTransfers)

	//Keep a running tally of all the badges we still have to handle
	unhandled := []*types.UniversalPermissionDetails{}
	unapproved := []*types.Balance{}

	for _, balance := range balancesToCheck {
		for _, badgeId := range balance.BadgeIds {
			for _, time := range balance.OwnershipTimes {
				unhandled = append(unhandled, &types.UniversalPermissionDetails{
					BadgeId:       badgeId,
					OwnershipTime: time,
					TimelineTime: &types.UintRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
					TransferTime: &types.UintRange{ Start: sdkmath.NewUint(math.MaxUint64), End: sdkmath.NewUint(math.MaxUint64) }, //dummy range
					ToMapping: &types.AddressMapping{},
					FromMapping: &types.AddressMapping{},
					InitiatedByMapping: &types.AddressMapping{},

					ArbitraryValue: balance.Amount,
				})
			}
		}
	}

	//Keep a list of all the user incoming/outgoing approvals we need to check
	userApprovalsToCheck := []*UserApprovalsToCheck{}
	for _, match := range firstMatches {
		transferVal := match.ArbitraryValue.(*types.CollectionApprovedTransfer)

		doAddressesMatch := k.CheckIfAddressesMatchCollectionMappingIds(ctx, transferVal, fromAddress, toAddress, initiatedBy, manager)
		if !doAddressesMatch {
			continue
		}

		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currTimeFound := types.SearchUintRangesForUint(currTime, []*types.UintRange{match.TransferTime})
		if !currTimeFound {
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

			ArbitraryValue: match.ArbitraryValue,
		}, unhandled)

		unhandled = remaining

		for _, overlap := range overlaps {
			totalBalances := []*types.Balance{{Amount: overlap.ArbitraryValue.(sdkmath.Uint), OwnershipTimes: []*types.UintRange{overlap.OwnershipTime}, BadgeIds: []*types.UintRange{overlap.BadgeId}}}

			allowed := transferVal.AllowedCombinations[0].IsAllowed //HACK: can do this because we expanded the allowed combinations above
			if !allowed {
				unapproved = append(unapproved, totalBalances...)
				continue
			}

			if transferVal.RequireFromDoesNotEqualInitiatedBy && fromAddress == initiatedBy {
				unapproved = append(unapproved, totalBalances...)
				continue
			}

			if transferVal.RequireFromEqualsInitiatedBy && fromAddress != initiatedBy {
				unapproved = append(unapproved, totalBalances...)
				continue
			}

			if transferVal.RequireToDoesNotEqualInitiatedBy && toAddress == initiatedBy {
				unapproved = append(unapproved, totalBalances...)
				continue
			}

			if transferVal.RequireToEqualsInitiatedBy && toAddress != initiatedBy {
				unapproved = append(unapproved, totalBalances...)
				continue
			}

			//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
			//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
			useLeafIndexForNumIncrements, numIncrements, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, transferVal.Challenges, solutions, initiatedBy, "overall", transferVal.ApprovalId, true)
			if err != nil {
				unapproved = append(unapproved, totalBalances...)
				continue
			}

			transferBalances := []*types.Balance{{Amount: overlap.ArbitraryValue.(sdkmath.Uint), OwnershipTimes: []*types.UintRange{overlap.OwnershipTime}, BadgeIds: []*types.UintRange{overlap.BadgeId}}}
			unapprovedBalances := types.DeepCopyBalances(transferBalances)


			if transferVal.ApprovalId == "" && (transferVal.OverallApprovals != nil || transferVal.PerAddressApprovals != nil) {
				unapproved = append(unapproved, totalBalances...)
				continue
			}


			//Handle the incrementing of overall approvals
			if transferVal.OverallApprovals != nil {
				transferBalances, err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.OverallApprovals.Amounts, transferVal.OverallApprovals.NumTransfers, transferBalances, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "overall", "", true)
				if err != nil {
					unapproved = append(unapproved, totalBalances...)
					continue
				}
			}

			//Handle the per-address approvals
			if transferVal.PerAddressApprovals != nil {
				if transferVal.PerAddressApprovals.ApprovalsPerToAddress != nil {
					transferBalances, err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.PerAddressApprovals.ApprovalsPerToAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers, transferBalances, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "per-to", toAddress, true)
					if err != nil {
						unapproved = append(unapproved, totalBalances...)
						continue
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerFromAddress != nil {
					transferBalances, err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers, transferBalances, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "per-from", fromAddress, true)
					if err != nil {
						unapproved = append(unapproved, totalBalances...)
						continue
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil {
					transferBalances, err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress.NumTransfers, transferBalances, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "per-initiated-by", initiatedBy, true)
					if err != nil {
						unapproved = append(unapproved, totalBalances...)
						continue
					}
				}
			}

			for _, transferBalance := range transferBalances {
				unapprovedBalances, err = types.SubtractBalance(unapprovedBalances, transferBalance)
				if err != nil {
					return []*UserApprovalsToCheck{}, []*types.Balance{}, err
				}


				//If we are overriding the approved outgoing / incoming transfers, we don't need to check the user approvals
				//Else, we do
				if !transferVal.OverridesFromApprovedOutgoingTransfers {
					userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
						Address:  fromAddress,
						Balances: []*types.Balance{transferBalance},
						Outgoing: true,
					})
				}

				if !transferVal.OverridesToApprovedIncomingTransfers {
					userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
						Address:  toAddress,
						Balances: []*types.Balance{transferBalance},
						Outgoing: false,
					})
				}
			}

			unapproved = append(unapproved, unapprovedBalances...)
		}
	}

	for _, unhandledDetails := range unhandled {
		balanceToRemove := &types.Balance{Amount: unhandledDetails.ArbitraryValue.(sdkmath.Uint), OwnershipTimes: []*types.UintRange{unhandledDetails.OwnershipTime}, BadgeIds: []*types.UintRange{unhandledDetails.BadgeId}}
		overallTransferBalances, err = types.SubtractBalance(overallTransferBalances, balanceToRemove)
		if err != nil {
			return []*UserApprovalsToCheck{}, []*types.Balance{}, err
		}
	}

	for _, unapprovedBalance := range unapproved {
		overallTransferBalances, err = types.SubtractBalance(overallTransferBalances, unapprovedBalance)
		if err != nil {
			return []*UserApprovalsToCheck{}, []*types.Balance{}, err
		}
	}

	return userApprovalsToCheck, overallTransferBalances, nil
}