package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Deduct approvals from requester if requester != from
func (k Keeper) DeductOutgoingApprovalsIfNeeded(ctx sdk.Context, collection types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.IdRange, times []*types.IdRange, from string, to string, requester string, amount sdk.Uint, solutions []*types.ChallengeSolution) (error) {
	//check all approvals, default disallow at the end unless from === initiator	
	currApprovedTransfers := GetCurrentUserApprovedOutgoingTransfers(ctx, userBalance)
	expandedApprovedTransfers := ExpandUserApprovedOutgoingTransfers(currApprovedTransfers)
	castedApprovedTransfers := CastUserApprovedOutgoingTransferToUniversalPermission(expandedApprovedTransfers)
	firstMatches := types.GetFirstMatchOnly(castedApprovedTransfers) //but could be duplicate mapping IDs so we need to be careful here

	//filter out those that don't match current time
	newFirstMatches := []*types.UniversalPermissionDetails{}
	for _, match := range firstMatches {
		timeFound := types.SearchIdRangesForId(sdk.NewUint(uint64(ctx.BlockTime().UnixMilli())), []*types.IdRange{match.TransferTime})
		if timeFound {
			newFirstMatches = append(newFirstMatches, match)
		}
	}

	unhandledBadgeIds := badgeIds
	overlaps := []*types.IdRange{}
	for _, match := range newFirstMatches {
		
		approvedTransfer := match.ArbitraryValue.(*types.UserApprovedOutgoingTransfer)

		doAddressesMatch := k.CheckIfAddressesMatchUserOutgoingMappingIds(ctx, approvedTransfer, from, to, requester, GetCurrentManager(ctx, collection))
		if !doAddressesMatch {
			continue
		}

		unhandledBadgeIds, overlaps = types.RemoveIdRangeFromIdRange([]*types.IdRange{match.BadgeId}, unhandledBadgeIds)
		if len(overlaps) > 0 {
			//We have a valid match for at least some badges, procees to check restrictions
			approvedTransfer := match.ArbitraryValue.(*types.UserApprovedOutgoingTransfer)
			isAllowed := approvedTransfer.AllowedCombinations[0].IsAllowed
			if !isAllowed {
				return ErrInadequateApprovals
			}
			
			requireToEqualsInitiatedBy := approvedTransfer.RequireToEqualsInitiatedBy
			requireToDoesNotEqualInitiatedBy := approvedTransfer.RequireToDoesNotEqualInitiatedBy

			if requireToEqualsInitiatedBy && to != requester {
				return ErrInadequateApprovals
			}

			if requireToDoesNotEqualInitiatedBy && to == requester {
				return ErrInadequateApprovals
			}

			//Get the current approvals for this transfer
			//If nil, we are approved for any amount / times in the badge ID range
			//Note we filter any excess badge IDs later and apply num increments as well
			approvals := approvedTransfer.Approvals
			if approvedTransfer.Approvals == nil {
				approvals = &types.ApprovalAmounts{
					StartAmounts: []*types.Balance{
						{
							Amount: amount,
							Times: times,
							BadgeIds: overlaps,
						},
					},
					IncrementIdsBy: sdk.NewUint(0),
					IncrementTimesBy: sdk.NewUint(0),
				}
			}

			//TODO: optimize this to avoid fetching unnecessarily (i guess only case would be when we use leaf index for order and have no perAddress or numTransfers restrictions)
			approvalTrackerDetails, found := k.GetTransferTrackerFromStore(ctx, collection.CollectionId, approvedTransfer.TrackerId, false, true, false, "")
			if !found {
				return ErrDisallowedTransfer
			}

			challengeHasNumIncrements, newNumIncrements, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, approvedTransfer.Challenges, solutions, requester, false, false, true)
			if err != nil {
				return ErrDisallowedTransfer
			}

			numIncrements := approvalTrackerDetails.NumTransfers
			if challengeHasNumIncrements {
				numIncrements = newNumIncrements
			}

			//allApprovals is the total amount approved by the user (i.e. the initial total amounts for all addresses)
			//we will check this transfer + tallies doesn't exceed this amount later
			allApprovals := []*types.Balance{}
			for _, startAmount := range approvals.StartAmounts {
				incrementedTimes := []*types.IdRange{}
				for _, time := range startAmount.Times {
					incrementedTimes = append(incrementedTimes, &types.IdRange{
						Start: time.Start.Add(numIncrements.Mul(approvals.IncrementTimesBy)),
						End: time.End.Add(numIncrements.Mul(approvals.IncrementTimesBy)),
					})
				}

				incrementedBadgeIds := []*types.IdRange{}
				for _, badgeId := range startAmount.BadgeIds {
					incrementedBadgeIds = append(incrementedBadgeIds, &types.IdRange{
						Start: badgeId.Start.Add(numIncrements.Mul(approvals.IncrementIdsBy)),
						End: badgeId.End.Add(numIncrements.Mul(approvals.IncrementIdsBy)),
					})
				}

				allApprovals = append(allApprovals, &types.Balance{
					Amount: startAmount.Amount,
					Times: incrementedTimes,
					BadgeIds: incrementedBadgeIds,
				})
			}
			
			//These are the amounts we need to check for the current transfer
			unhandledCurrTransferApprovals := []*types.Balance{{
					Amount: amount,
					Times: times,
					BadgeIds: overlaps,
			}}

			//here, we do two things at once:
			//1. remove all non-overlapping badge IDs (could sneak in with the num increments calculation)
			//2. check that the entire current transfer is handled via all approvals (if not, we can just fail below if leftovers)
			//	 note this just checks the overall initial amounts, not the tallies
			for _, approval := range allApprovals {
				//we only want to handle the overlaps 
				_, approval.BadgeIds = types.RemoveIdRangeFromIdRange(overlaps, approval.BadgeIds)
				if len(approval.BadgeIds) > 0 {
					unhandledCurrTransferApprovals, err = SubtractBalancesForIdRanges(unhandledCurrTransferApprovals, approval.BadgeIds, approval.Times, approval.Amount)
					if err != nil {
						return ErrDisallowedTransfer
					}
				}
			}

			if len(unhandledCurrTransferApprovals) > 0 {
				return ErrDisallowedTransfer
			}

			toTrackerDetails := types.ApprovalsTracker{}
			needToFetchTo := (approvedTransfer.PerAddressApprovals != nil && approvedTransfer.PerAddressApprovals.ApprovalsPerToAddress != nil) ||
				(approvedTransfer.PerAddressMaxNumTransfers != nil && !approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerToAddress.IsNil()) 
			if needToFetchTo {
				toTrackerDetails, found = k.GetTransferTrackerFromStore(ctx, collection.CollectionId, approvedTransfer.TrackerId, false, true, false, to)
				if !found {
					return ErrDisallowedTransfer
				}
			}

			initiatedByTrackerDetails := types.ApprovalsTracker{}
			needToFetchInitiatedBy := (approvedTransfer.PerAddressApprovals != nil && approvedTransfer.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil) ||
				(approvedTransfer.PerAddressMaxNumTransfers != nil && !approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerInitiatedByAddress.IsNil()) 
			if needToFetchInitiatedBy {
				initiatedByTrackerDetails, found = k.GetTransferTrackerFromStore(ctx, collection.CollectionId, approvedTransfer.TrackerId, false, true, false, requester)
				if !found {
					return ErrDisallowedTransfer
				}
			}

			currTransferBalanceToAdd := &types.Balance{
				Amount: amount,
				Times: times,
				BadgeIds: overlaps,
			}
			
			//if === nil, no restrictions so we don't have to check
			//if increments > 0, then it doesn't make sense to check approvals bc they are not automatically calculated
			if approvedTransfer.Approvals != nil && (approvedTransfer.Approvals.IncrementIdsBy == sdk.NewUint(0) && approvedTransfer.Approvals.IncrementTimesBy == sdk.NewUint(0)) {
				approvalTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(approvalTrackerDetails.Amounts, currTransferBalanceToAdd, allApprovals)
				if err != nil {
					return err
				}
			}

			if approvedTransfer.PerAddressApprovals != nil {
				if approvedTransfer.PerAddressApprovals.ApprovalsPerToAddress != nil {
					toTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(toTrackerDetails.Amounts, currTransferBalanceToAdd, approvedTransfer.PerAddressApprovals.ApprovalsPerToAddress)
					if err != nil {
						return err
					}
				}

				if approvedTransfer.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil {
					initiatedByTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(initiatedByTrackerDetails.Amounts, currTransferBalanceToAdd, approvedTransfer.PerAddressApprovals.ApprovalsPerInitiatedByAddress)
					if err != nil {
						return err
					}	
				}
			}

			if !approvedTransfer.MaxNumTransfers.IsNil() {
				newTally := approvalTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
				if newTally.GT(approvedTransfer.MaxNumTransfers) {
					return ErrDisallowedTransfer
				}

				approvalTrackerDetails.NumTransfers = newTally
			}

			if approvedTransfer.PerAddressMaxNumTransfers != nil {
				if !approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerToAddress.IsNil() {
					newTally := toTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
					if newTally.GT(approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerToAddress) {
						return ErrDisallowedTransfer
					}

					toTrackerDetails.NumTransfers = newTally
				}

				if !approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerInitiatedByAddress.IsNil() {
					newTally := initiatedByTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
					if newTally.GT(approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerInitiatedByAddress) {
						return ErrDisallowedTransfer
					}

					initiatedByTrackerDetails.NumTransfers = newTally
				}
			}


			err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, approvalTrackerDetails, false, true, false, "")
			if err != nil {
				return ErrDisallowedTransfer
			}

			if needToFetchTo {
				err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, toTrackerDetails, false, true, false, to)
				if err != nil {
					return ErrDisallowedTransfer
				}
			}

			if needToFetchInitiatedBy {
				err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, initiatedByTrackerDetails, false, true, false, requester)
				if err != nil {
					return ErrDisallowedTransfer
				}
			}
		}
	}

	if len(unhandledBadgeIds) > 0 {
		return ErrInadequateApprovals
	}

	return nil
}

// Deduct approvals from requester if requester != from
func (k Keeper) DeductIncomingApprovalsIfNeeded(ctx sdk.Context, collection types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.IdRange, times []*types.IdRange, from string, to string, requester string, amount sdk.Uint, solutions []*types.ChallengeSolution) (error) {
	//check all approvals, default disallow at the end unless from === initiator
	currApprovedTransfers := GetCurrentUserApprovedIncomingTransfers(ctx, userBalance)
	expandedApprovedTransfers := ExpandUserApprovedIncomingTransfers(currApprovedTransfers)
	castedApprovedTransfers := CastUserApprovedIncomingTransferToUniversalPermission(expandedApprovedTransfers)
	firstMatches := types.GetFirstMatchOnly(castedApprovedTransfers) //but could be duplicate mapping IDs so we need to be careful here

	//filter out those that don't match current time
	newFirstMatches := []*types.UniversalPermissionDetails{}
	for _, match := range firstMatches {
		timeFound := types.SearchIdRangesForId(sdk.NewUint(uint64(ctx.BlockTime().UnixMilli())), []*types.IdRange{match.TransferTime})
		if timeFound {
			newFirstMatches = append(newFirstMatches, match)
		}
	}

	unhandledBadgeIds := badgeIds
	overlaps := []*types.IdRange{}
	for _, match := range newFirstMatches {
		
		approvedTransfer := match.ArbitraryValue.(*types.UserApprovedIncomingTransfer)

		doAddressesMatch := k.CheckIfAddressesMatchUserIncomingMappingIds(ctx, approvedTransfer, from, to, requester, GetCurrentManager(ctx, collection))
		if !doAddressesMatch {
			continue
		}

		unhandledBadgeIds, overlaps = types.RemoveIdRangeFromIdRange([]*types.IdRange{match.BadgeId}, unhandledBadgeIds)
		if len(overlaps) > 0 {
			//We have a valid match for at least some badges, procees to check restrictions
			approvedTransfer := match.ArbitraryValue.(*types.UserApprovedIncomingTransfer)
			isAllowed := approvedTransfer.AllowedCombinations[0].IsAllowed
			if !isAllowed {
				return ErrInadequateApprovals
			}
			
			requireFromEqualsInitiatedBy := approvedTransfer.RequireFromEqualsInitiatedBy
			requireFromDoesNotEqualInitiatedBy := approvedTransfer.RequireFromDoesNotEqualInitiatedBy

			if requireFromEqualsInitiatedBy && from != requester {
				return ErrInadequateApprovals
			}

			if requireFromDoesNotEqualInitiatedBy && from == requester {
				return ErrInadequateApprovals
			}

			//Get the current approvals for this transfer
			//If nil, we are approved for any amount / times in the badge ID range
			//Note we filter any excess badge IDs later and apply num increments as well
			approvals := approvedTransfer.Approvals
			if approvedTransfer.Approvals == nil {
				approvals = &types.ApprovalAmounts{
					StartAmounts: []*types.Balance{
						{
							Amount: amount,
							Times: times,
							BadgeIds: overlaps,
						},
					},
					IncrementIdsBy: sdk.NewUint(0),
					IncrementTimesBy: sdk.NewUint(0),
				}
			}

			

			//TODO: optimize this to avoid fetching unnecessarily (i guess only case would be when we use leaf index for order and have no perAddress or numTransfers restrictions)
			approvalTrackerDetails, found := k.GetTransferTrackerFromStore(ctx, collection.CollectionId, approvedTransfer.TrackerId, false, true, false, "")
			if !found {
				return ErrDisallowedTransfer
			}
			
		  challengeHasNumIncrements, newNumIncrements, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, approvedTransfer.Challenges, solutions, requester, true, false, false)
			if err != nil {
				return ErrDisallowedTransfer
			}

			numIncrements := approvalTrackerDetails.NumTransfers
			if challengeHasNumIncrements {
				numIncrements = newNumIncrements
			}

			//allApprovals is the total amount approved by the user (i.e. the initial total amounts for all addresses)
			//we will check this transfer + tallies doesn't exceed this amount later
			allApprovals := []*types.Balance{}
			for _, startAmount := range approvals.StartAmounts {
				incrementedTimes := []*types.IdRange{}
				for _, time := range startAmount.Times {
					incrementedTimes = append(incrementedTimes, &types.IdRange{
						Start: time.Start.Add(numIncrements.Mul(approvals.IncrementTimesBy)),
						End: time.End.Add(numIncrements.Mul(approvals.IncrementTimesBy)),
					})
				}

				incrementedBadgeIds := []*types.IdRange{}
				for _, badgeId := range startAmount.BadgeIds {
					incrementedBadgeIds = append(incrementedBadgeIds, &types.IdRange{
						Start: badgeId.Start.Add(numIncrements.Mul(approvals.IncrementIdsBy)),
						End: badgeId.End.Add(numIncrements.Mul(approvals.IncrementIdsBy)),
					})
				}

				allApprovals = append(allApprovals, &types.Balance{
					Amount: startAmount.Amount,
					Times: incrementedTimes,
					BadgeIds: incrementedBadgeIds,
				})
			}
			
			//These are the amounts we need to check for the current transfer
			unhandledCurrTransferApprovals := []*types.Balance{{
					Amount: amount,
					Times: times,
					BadgeIds: overlaps,
			}}

			//here, we do two things at once:
			//1. remove all non-overlapping badge IDs (could sneak in with the num increments calculation)
			//2. check that the entire current transfer is handled via all approvals (if not, we can just fail below if leftovers)
			//	 note this just checks the overall initial amounts, not the tallies
			for _, approval := range allApprovals {
				//we only want to handle the overlaps 
				_, approval.BadgeIds = types.RemoveIdRangeFromIdRange(overlaps, approval.BadgeIds)
				if len(approval.BadgeIds) > 0 {
					unhandledCurrTransferApprovals, err = SubtractBalancesForIdRanges(unhandledCurrTransferApprovals, approval.BadgeIds, approval.Times, approval.Amount)
					if err != nil {
						return ErrDisallowedTransfer
					}
				}
			}

			if len(unhandledCurrTransferApprovals) > 0 {
				return ErrDisallowedTransfer
			}

			fromTrackerDetails := types.ApprovalsTracker{}
			needToFetchFrom := (approvedTransfer.PerAddressApprovals != nil && approvedTransfer.PerAddressApprovals.ApprovalsPerFromAddress != nil) ||
				(approvedTransfer.PerAddressMaxNumTransfers != nil && !approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerFromAddress.IsNil()) 
			if needToFetchFrom {
				fromTrackerDetails, found = k.GetTransferTrackerFromStore(ctx, collection.CollectionId, approvedTransfer.TrackerId, false, true, false, from)
				if !found {
					return ErrDisallowedTransfer
				}
			}

			initiatedByTrackerDetails := types.ApprovalsTracker{}
			needToFetchInitiatedBy := (approvedTransfer.PerAddressApprovals != nil && approvedTransfer.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil) ||
				(approvedTransfer.PerAddressMaxNumTransfers != nil && !approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerInitiatedByAddress.IsNil()) 
			if needToFetchInitiatedBy {
				initiatedByTrackerDetails, found = k.GetTransferTrackerFromStore(ctx, collection.CollectionId, approvedTransfer.TrackerId, false, true, false, requester)
				if !found {
					return ErrDisallowedTransfer
				}
			}

			currTransferBalanceToAdd := &types.Balance{
				Amount: amount,
				Times: times,
				BadgeIds: overlaps,
			}
			
			//if === nil, no restrictions so we don't have to check
			//if increments > 0, then it doesn't make sense to check approvals bc they are not automatically calculated
			if approvedTransfer.Approvals != nil && (approvedTransfer.Approvals.IncrementIdsBy == sdk.NewUint(0) && approvedTransfer.Approvals.IncrementTimesBy == sdk.NewUint(0)) {
				approvalTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(approvalTrackerDetails.Amounts, currTransferBalanceToAdd, allApprovals)
				if err != nil {
					return err
				}
			}

			if approvedTransfer.PerAddressApprovals != nil {
				if approvedTransfer.PerAddressApprovals.ApprovalsPerFromAddress != nil {
					fromTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(fromTrackerDetails.Amounts, currTransferBalanceToAdd, approvedTransfer.PerAddressApprovals.ApprovalsPerFromAddress)
					if err != nil {
						return err
					}
				}

				if approvedTransfer.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil {
					initiatedByTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(initiatedByTrackerDetails.Amounts, currTransferBalanceToAdd, approvedTransfer.PerAddressApprovals.ApprovalsPerInitiatedByAddress)
					if err != nil {
						return err
					}	
				}
			}

			if !approvedTransfer.MaxNumTransfers.IsNil() {
				newTally := approvalTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
				if newTally.GT(approvedTransfer.MaxNumTransfers) {
					return ErrDisallowedTransfer
				}

				approvalTrackerDetails.NumTransfers = newTally
			}

			if approvedTransfer.PerAddressMaxNumTransfers != nil {
				if !approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerFromAddress.IsNil() {
					newTally := fromTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
					if newTally.GT(approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerFromAddress) {
						return ErrDisallowedTransfer
					}

					fromTrackerDetails.NumTransfers = newTally
				}

				if !approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerInitiatedByAddress.IsNil() {
					newTally := initiatedByTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
					if newTally.GT(approvedTransfer.PerAddressMaxNumTransfers.MaxNumTransfersPerInitiatedByAddress) {
						return ErrDisallowedTransfer
					}

					initiatedByTrackerDetails.NumTransfers = newTally
				}
			}


			err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, approvalTrackerDetails, false, true, false, "")
			if err != nil {
				return ErrDisallowedTransfer
			}

			if needToFetchFrom {
				err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, fromTrackerDetails, false, true, false, to)
				if err != nil {
					return ErrDisallowedTransfer
				}
			}

			if needToFetchInitiatedBy {
				err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, initiatedByTrackerDetails, false, true, false, requester)
				if err != nil {
					return ErrDisallowedTransfer
				}
			}
		}
	}

	if len(unhandledBadgeIds) > 0 {
		return ErrInadequateApprovals
	}

	return nil
}


type UserApprovalsToCheck struct {
	Address string
	BadgeIds []*types.IdRange
	Outgoing bool
}


func AddTallyAndAssertDoesntExceedThreshold(curr []*types.Balance, toAdd *types.Balance, threshold []*types.Balance) ([]*types.Balance, error) {
	threhsoldCopy := make([]*types.Balance, len(threshold))
	copy(threhsoldCopy, threshold)

	
	newTallied, err := AddBalancesForIdRanges(curr, toAdd.BadgeIds, toAdd.Times, toAdd.Amount)
	if err != nil {
		return []*types.Balance{}, ErrDisallowedTransfer
	}

	for _, newTalliedAmount := range newTallied {
		threshold, err = SubtractBalancesForIdRanges(threshold, newTalliedAmount.BadgeIds, newTalliedAmount.Times, newTalliedAmount.Amount)
		if err != nil {
			return []*types.Balance{}, ErrDisallowedTransfer
		}
	}

	curr = newTallied

	return curr, nil
}


// Checks if account is frozen or not.
func (k Keeper) CheckIfApprovedOnCollectionLevelAndGetUserApprovalsToCheck(ctx sdk.Context, collection types.BadgeCollection, badgeIds []*types.IdRange, times []*types.IdRange, fromAddress string, toAddress string, initiatedBy string, amount sdk.Uint, solutions []*types.ChallengeSolution) ([]*UserApprovalsToCheck, error) {
	approvedTransfers := GetCurrentCollectionApprovedTransfers(ctx, collection)
	expandedApprovedTransfers := ExpandCollectionApprovedTransfers(approvedTransfers)
	castedApprovedTransfers := CastCollectionApprovedTransferToUniversalPermission(expandedApprovedTransfers)
	firstMatches := types.GetFirstMatchOnly(castedApprovedTransfers) //but could be duplicate mapping IDs so we need to be careful here

	//filter out those that don't match current time
	newMatches := []*types.UniversalPermissionDetails{}
	for _, match := range firstMatches {
		timeFound := types.SearchIdRangesForId(sdk.NewUint(uint64(ctx.BlockTime().UnixMilli())), []*types.IdRange{match.TransferTime})
		if timeFound {
			newMatches = append(newMatches, match)
		}
	}

	//keep a running tally of all the badges we still have to handle
	unhandledBadgeIds := make([]*types.IdRange, len(badgeIds))
	copy(unhandledBadgeIds, badgeIds)


	overlaps := []*types.IdRange{}
	userApprovalsToCheck := []*UserApprovalsToCheck{}
	manager := GetCurrentManager(ctx, collection)
	for _, match := range newMatches {
		//Check if addresses match
		transferVal := match.ArbitraryValue.(*types.CollectionApprovedTransfer)

		doAddressesMatch := k.CheckIfAddressesMatchCollectionMappingIds(ctx, transferVal, fromAddress, toAddress, initiatedBy, manager)
		if !doAddressesMatch {
			continue
		}

		//Note we already filtered by transfer time above so all we have left is badge IDs
		unhandledBadgeIds, overlaps = types.RemoveIdRangeFromIdRange([]*types.IdRange{match.BadgeId}, unhandledBadgeIds)
		if len(overlaps) > 0 {
			//We have a valid match for at least some badges, procees to check restrictions
			//If something is invalid within this, we must throw an error because the transfer for those bages are invalid and thus the whole transfer is invalid

			allowed := transferVal.AllowedCombinations[0].IsAllowed //HACK: can do this because we expanded the allowed combinations above
			if !allowed {
				return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			if transferVal.RequireFromDoesNotEqualInitiatedBy && fromAddress == initiatedBy {
				return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			if transferVal.RequireFromEqualsInitiatedBy && fromAddress != initiatedBy {
				return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			if transferVal.RequireToDoesNotEqualInitiatedBy && toAddress == initiatedBy {
				return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			if transferVal.RequireToEqualsInitiatedBy && toAddress != initiatedBy {
				return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			//Get the current approvals for this transfer
			//If nil, we are approved for any amount / times in the badge ID range
			//Note we filter any excess badge IDs later and apply num increments as well
			approvals := transferVal.Approvals
			if transferVal.Approvals == nil {
				approvals = &types.ApprovalAmounts{
					StartAmounts: []*types.Balance{
						{
							Amount: amount,
							Times: times,
							BadgeIds: overlaps,
						},
					},
					IncrementIdsBy: sdk.NewUint(0),
					IncrementTimesBy: sdk.NewUint(0),
				}
			}

			
			//TODO: optimize this to avoid fetching unnecessarily (i guess only case would be when we use leaf index for order and have no perAddress or numTransfers restrictions)
			approvalTrackerDetails, found := k.GetTransferTrackerFromStore(ctx, collection.CollectionId, transferVal.TrackerId, true, false, false, "")
			if !found {
				return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			challengeHasNumIncrements, newNumIncrements, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, transferVal.Challenges, solutions, initiatedBy, false, true, false)
			if err != nil {
				return  []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			numIncrements := approvalTrackerDetails.NumTransfers
			if challengeHasNumIncrements {
				numIncrements = newNumIncrements
			}

			//allApprovals is the total amount approved by the user (i.e. the initial total amounts for all addresses)
			//we will check this transfer + tallies doesn't exceed this amount later
			allApprovals := []*types.Balance{}
			for _, startAmount := range approvals.StartAmounts {
				incrementedTimes := []*types.IdRange{}
				for _, time := range startAmount.Times {
					incrementedTimes = append(incrementedTimes, &types.IdRange{
						Start: time.Start.Add(numIncrements.Mul(approvals.IncrementTimesBy)),
						End: time.End.Add(numIncrements.Mul(approvals.IncrementTimesBy)),
					})
				}

				incrementedBadgeIds := []*types.IdRange{}
				for _, badgeId := range startAmount.BadgeIds {
					incrementedBadgeIds = append(incrementedBadgeIds, &types.IdRange{
						Start: badgeId.Start.Add(numIncrements.Mul(approvals.IncrementIdsBy)),
						End: badgeId.End.Add(numIncrements.Mul(approvals.IncrementIdsBy)),
					})
				}

				allApprovals = append(allApprovals, &types.Balance{
					Amount: startAmount.Amount,
					Times: incrementedTimes,
					BadgeIds: incrementedBadgeIds,
				})
			}
			
			//These are the amounts we need to check for the current transfer
			unhandledCurrTransferApprovals := []*types.Balance{{
					Amount: amount,
					Times: times,
					BadgeIds: overlaps,
			}}

			//here, we do two things at once:
			//1. remove all non-overlapping badge IDs (could sneak in with the num increments calculation)
			//2. check that the entire current transfer is handled via all approvals (if not, we can just fail below if leftovers)
			//	 note this just checks the overall initial amounts, not the tallies
			for _, approval := range allApprovals {
				//we only want to handle the overlaps 
				_, approval.BadgeIds = types.RemoveIdRangeFromIdRange(overlaps, approval.BadgeIds)
				if len(approval.BadgeIds) > 0 {
					unhandledCurrTransferApprovals, err = SubtractBalancesForIdRanges(unhandledCurrTransferApprovals, approval.BadgeIds, approval.Times, approval.Amount)
					if err != nil {
						return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
					}
				}
			}

			if len(unhandledCurrTransferApprovals) > 0 {
				return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			//Only fetch if needed
			fromTrackerDetails := types.ApprovalsTracker{}
			needToFetchFrom := (transferVal.PerAddressApprovals != nil && transferVal.PerAddressApprovals.ApprovalsPerFromAddress != nil) ||
				(transferVal.PerAddressMaxNumTransfers != nil && !transferVal.PerAddressMaxNumTransfers.MaxNumTransfersPerFromAddress.IsNil())
			if needToFetchFrom {
				fromTrackerDetails, found = k.GetTransferTrackerFromStore(ctx, collection.CollectionId, transferVal.TrackerId, true, false, false, fromAddress)
				if !found {
					return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
				}
			}

			toTrackerDetails := types.ApprovalsTracker{}
			needToFetchTo := (transferVal.PerAddressApprovals != nil && transferVal.PerAddressApprovals.ApprovalsPerToAddress != nil) ||
				(transferVal.PerAddressMaxNumTransfers != nil && !transferVal.PerAddressMaxNumTransfers.MaxNumTransfersPerToAddress.IsNil()) 
			if needToFetchTo {
				toTrackerDetails, found = k.GetTransferTrackerFromStore(ctx, collection.CollectionId, transferVal.TrackerId, true, false, false, toAddress)
				if !found {
					return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
				}
			}

			initiatedByTrackerDetails := types.ApprovalsTracker{}
			needToFetchInitiatedBy := (transferVal.PerAddressApprovals != nil && transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil) ||
				(transferVal.PerAddressMaxNumTransfers != nil && !transferVal.PerAddressMaxNumTransfers.MaxNumTransfersPerInitiatedByAddress.IsNil()) 
			if needToFetchInitiatedBy {
				initiatedByTrackerDetails, found = k.GetTransferTrackerFromStore(ctx, collection.CollectionId, transferVal.TrackerId, true, false, false, initiatedBy)
				if !found {
					return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
				}
			}

			currTransferBalanceToAdd := &types.Balance{
				Amount: amount,
				Times: times,
				BadgeIds: overlaps,
			}
			
			//if === nil, no restrictions so we don't have to check
			//if increments > 0, then it doesn't make sense to check approvals bc they are not automatically calculated
			if transferVal.Approvals != nil && (transferVal.Approvals.IncrementIdsBy == sdk.NewUint(0) && transferVal.Approvals.IncrementTimesBy == sdk.NewUint(0)) {
				approvalTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(approvalTrackerDetails.Amounts, currTransferBalanceToAdd, allApprovals)
				if err != nil {
					return []*UserApprovalsToCheck{}, err
				}
			}

			if transferVal.PerAddressApprovals != nil {
				if transferVal.PerAddressApprovals.ApprovalsPerFromAddress != nil {
					fromTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(fromTrackerDetails.Amounts, currTransferBalanceToAdd, transferVal.PerAddressApprovals.ApprovalsPerFromAddress)
					if err != nil {
						return []*UserApprovalsToCheck{}, err
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerToAddress != nil {
					toTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(toTrackerDetails.Amounts, currTransferBalanceToAdd, transferVal.PerAddressApprovals.ApprovalsPerToAddress)
					if err != nil {
						return []*UserApprovalsToCheck{}, err
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil {
					initiatedByTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(initiatedByTrackerDetails.Amounts, currTransferBalanceToAdd, transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress)
					if err != nil {
						return []*UserApprovalsToCheck{}, err
					}	
				}
			}

			if !transferVal.MaxNumTransfers.IsNil() {
				newTally := approvalTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
				if newTally.GT(transferVal.MaxNumTransfers) {
					return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
				}

				approvalTrackerDetails.NumTransfers = newTally
			}

			if transferVal.PerAddressMaxNumTransfers != nil {
				if !transferVal.PerAddressMaxNumTransfers.MaxNumTransfersPerFromAddress.IsNil() {
					newTally := fromTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
					if newTally.GT(transferVal.PerAddressMaxNumTransfers.MaxNumTransfersPerFromAddress) {
						return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
					}

					fromTrackerDetails.NumTransfers = newTally
				}

				if !transferVal.PerAddressMaxNumTransfers.MaxNumTransfersPerToAddress.IsNil() {
					newTally := toTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
					if newTally.GT(transferVal.PerAddressMaxNumTransfers.MaxNumTransfersPerToAddress) {
						return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
					}

					toTrackerDetails.NumTransfers = newTally
				}

				if !transferVal.PerAddressMaxNumTransfers.MaxNumTransfersPerInitiatedByAddress.IsNil() {
					newTally := initiatedByTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
					if newTally.GT(transferVal.PerAddressMaxNumTransfers.MaxNumTransfersPerInitiatedByAddress) {
						return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
					}

					initiatedByTrackerDetails.NumTransfers = newTally
				}
			}


			err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, approvalTrackerDetails, true, false, false, "")
			if err != nil {
				return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			if needToFetchFrom {
				err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, fromTrackerDetails, true, false, false, fromAddress)
				if err != nil {
					return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
				}
			}

			if needToFetchTo {
				err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, toTrackerDetails, true, false, false, toAddress)
				if err != nil {
					return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
				}
			}

			if needToFetchInitiatedBy {
				err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, initiatedByTrackerDetails, true, false, false, initiatedBy)
				if err != nil {
					return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
				}
			}
			
			if !transferVal.OverridesToApprovedIncomingTransfers {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address: toAddress,
					BadgeIds: overlaps,
				})
			}

			if !transferVal.OverridesFromApprovedOutgoingTransfers {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address: fromAddress,
					BadgeIds: overlaps,
					Outgoing: true,
				})
			}
		}
	}

	//If not explicitly allowed, we return that it is disallowed
	return userApprovalsToCheck, nil
}
