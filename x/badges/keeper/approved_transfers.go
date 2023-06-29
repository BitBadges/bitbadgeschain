package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)

//The UserApprovalsToCheck struct is used to keep track of which incoming / outgoing approvals for which addresses we need to check.
type UserApprovalsToCheck struct {
	Address string
	BadgeIds []*types.IdRange
	Outgoing bool
}

//DeductUserOutgoingApprovals will check if the current transfer is approved from the from's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserOutgoingApprovals(ctx sdk.Context, collection types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.IdRange, times []*types.IdRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.ChallengeSolution) (error) {
	currApprovedTransfers := GetCurrentUserApprovedOutgoingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForOutgoing(currApprovedTransfers, from)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastOutgoingTransfersToCollectionTransfers(currApprovedTransfers, from)
	_, err := k.DeductAndGetUserApprovals(castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "outgoing")
	return err
}

//DeductUserIncomingApprovals will check if the current transfer is approved from the to's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserIncomingApprovals(ctx sdk.Context, collection types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.IdRange, times []*types.IdRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.ChallengeSolution) (error) {
	currApprovedTransfers := GetCurrentUserApprovedIncomingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForIncoming(currApprovedTransfers, to)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastIncomingTransfersToCollectionTransfers(currApprovedTransfers, to)
	_, err := k.DeductAndGetUserApprovals(castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "incoming")
	return err
}

//DeductCollectionApprovalsAndGetUserApprovalsToCheck will check if the current transfer is allowed via the collection's approved transfers and handle any tallying accordingly
func (k Keeper) DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx sdk.Context, collection types.BadgeCollection, badgeIds []*types.IdRange, times []*types.IdRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.ChallengeSolution) ([]*UserApprovalsToCheck, error) {
	approvedTransfers := GetCurrentCollectionApprovedTransfers(ctx, collection)
	return k.DeductAndGetUserApprovals(approvedTransfers, ctx, collection, badgeIds, times, fromAddress, toAddress, initiatedBy, amount, solutions, "overall")
}	

func (k Keeper) DeductAndGetUserApprovals(approvedTransfers []*types.CollectionApprovedTransfer, ctx sdk.Context, collection types.BadgeCollection, badgeIds []*types.IdRange, times []*types.IdRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.ChallengeSolution, timelineType string) ([]*UserApprovalsToCheck, error) {
	//HACK: We first expand all transfers to have just a len == 1 AllowedCombination[] so that we can easily check IsAllowed later
	//		  This is because GetFirstMatchOnly will break down the transfers into smaller parts and without expansion, fetching if a certain transfer is allowed is impossible. 
	expandedApprovedTransfers := ExpandCollectionApprovedTransfers(approvedTransfers) 
	castedApprovedTransfers := CastCollectionApprovedTransferToUniversalPermission(expandedApprovedTransfers)
	firstMatches := types.GetFirstMatchOnly(castedApprovedTransfers) //Note: This filters only on the mapping ID level. We need to ensure we use first match only for each (to, from, initiatedBy) here. 

	manager := GetCurrentManager(ctx, collection)

	//Keep a running tally of all the badges we still have to handle
	unhandledBadgeIds := make([]*types.IdRange, len(badgeIds))
	copy(unhandledBadgeIds, badgeIds)

	//We will need to return a list of all the user incoming/outgoing approvals we need to check
	userApprovalsToCheck := []*UserApprovalsToCheck{}
	for _, match := range firstMatches {
		transferVal := match.ArbitraryValue.(*types.CollectionApprovedTransfer)

		doAddressesMatch := k.CheckIfAddressesMatchCollectionMappingIds(ctx, transferVal, fromAddress, toAddress, initiatedBy, manager)
		if !doAddressesMatch {
			continue
		}

		currTimeFound := types.SearchIdRangesForId(sdk.NewUint(uint64(ctx.BlockTime().UnixMilli())), []*types.IdRange{match.TransferTime})
		if !currTimeFound {
			continue
		}

		remaining, overlaps := types.RemoveIdRangeFromIdRange([]*types.IdRange{match.BadgeId}, unhandledBadgeIds)
		unhandledBadgeIds = remaining
		if len(overlaps) > 0 {
			//For the overlapping badges, we have a match.
			//We can now proceed to check restrictions. 
			//If any restriction fails, we MUST throw because the transfer is invalid for these badges

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

			//If the approval has challenges, we need to check that the solutions are valid
			//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
			useLeafIndexForNumIncrements, numIncrements, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, transferVal.Challenges, solutions, initiatedBy, "overall")
			if err != nil {
				return  []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			transferBalance := &types.Balance{ Amount: amount,	Times: times,	BadgeIds: overlaps }


			if transferVal.TrackerId == "" && (transferVal.OverallApprovals != nil || transferVal.PerAddressApprovals != nil) {
				return []*UserApprovalsToCheck{}, ErrDisallowedTransfer
			}

			if transferVal.OverallApprovals != nil {
				err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, &collection, transferVal.OverallApprovals.Amounts, transferVal.OverallApprovals.NumTransfers, transferBalance, transferVal.TrackerId, transferVal.IncrementIdsBy, transferVal.IncrementTimesBy,	useLeafIndexForNumIncrements, numIncrements, timelineType, "overall", "")
				if err != nil {
					return []*UserApprovalsToCheck{}, err
				}
			}

			if transferVal.PerAddressApprovals != nil {
				if transferVal.PerAddressApprovals.ApprovalsPerToAddress != nil {
					err := k.IncrementApprovalsAndAssertWithinThreshold(ctx, &collection, transferVal.PerAddressApprovals.ApprovalsPerToAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers, transferBalance, transferVal.TrackerId, transferVal.IncrementIdsBy, transferVal.IncrementTimesBy,	useLeafIndexForNumIncrements, numIncrements, timelineType, "per-to", toAddress)
					if err != nil {
						return []*UserApprovalsToCheck{}, err
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerFromAddress != nil {
					err := k.IncrementApprovalsAndAssertWithinThreshold(ctx, &collection, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers, transferBalance, transferVal.TrackerId, transferVal.IncrementIdsBy, transferVal.IncrementTimesBy,	useLeafIndexForNumIncrements, numIncrements, timelineType, "per-from", fromAddress)
					if err != nil {
						return []*UserApprovalsToCheck{}, err
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil {
					err := k.IncrementApprovalsAndAssertWithinThreshold(ctx, &collection, transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress.NumTransfers, transferBalance, transferVal.TrackerId, transferVal.IncrementIdsBy, transferVal.IncrementTimesBy,	useLeafIndexForNumIncrements, numIncrements, timelineType, "per-initiated-by", initiatedBy)
					if err != nil {
						return []*UserApprovalsToCheck{}, err
					}
				}
			}

			if !transferVal.OverridesFromApprovedOutgoingTransfers {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address: fromAddress,
					BadgeIds: overlaps,
					Outgoing: true,
				})
			}

			if !transferVal.OverridesToApprovedIncomingTransfers {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address: toAddress,
					BadgeIds: overlaps,
					Outgoing: false,
				})
			}
		}
	}

	if len(unhandledBadgeIds) > 0 {
		return []*UserApprovalsToCheck{}, ErrInadequateApprovals
	}

	//If not explicitly allowed, we return that it is disallowed
	return userApprovalsToCheck, nil
}

func AddTallyAndAssertDoesntExceedThreshold(currTally []*types.Balance, toAdd *types.Balance, threshold []*types.Balance) ([]*types.Balance, error) {
	//Add the new tally
	err := *new(error)
	currTally, err = types.AddBalancesForIdRanges(currTally, toAdd.BadgeIds, toAdd.Times, toAdd.Amount)
	if err != nil {
		return []*types.Balance{}, err
	}

	//Check if we exceed the threshold; will underflow if we do exceed
	thresholdCopy := make([]*types.Balance, len(threshold))
	copy(thresholdCopy, threshold)
	for _, newTalliedAmount := range currTally {
		thresholdCopy, err = types.SubtractBalancesForIdRanges(thresholdCopy, newTalliedAmount.BadgeIds, newTalliedAmount.Times, newTalliedAmount.Amount)
		if err != nil {
			return []*types.Balance{}, err
		}
	}

	return currTally, nil
}


func (k Keeper) IncrementApprovalsAndAssertWithinThreshold(
	ctx sdk.Context, 
	collection *types.BadgeCollection, 
	approvals []*types.Balance, 
	maxNumTransfers sdkmath.Uint,
	transferAmounts *types.Balance, 
	trackerId string, 
	incrementIdsBy sdkmath.Uint,
	incrementTimesBy sdkmath.Uint,
	precalculatedNumIncrements bool, 
	numIncrements sdkmath.Uint,
	timelineType string,
	depth string,
	address string,
) (error) {
	//Get the current approvals for this transfer
	//If nil, no restrictions and we are approved for the entire transfer
	//Note we filter any excess badge IDs later and apply num increments as well
	err := *new(error)
	if approvals == nil {
		approvals = []*types.Balance{{
			Amount: transferAmounts.Amount,
			Times: transferAmounts.Times,
			BadgeIds: transferAmounts.BadgeIds,
		}}
	}

	approvalTrackerDetails, found := k.GetTransferTrackerFromStore(ctx, collection.CollectionId, trackerId, timelineType, depth, address)
	if !found {
		approvalTrackerDetails = types.ApprovalsTracker{
			Amounts: []*types.Balance{},
			NumTransfers: sdk.NewUint(0),
		}
	}

	if !precalculatedNumIncrements {
		numIncrements = approvalTrackerDetails.NumTransfers
	}

	//allApprovals is the total amount approved (i.e. the initial total amounts plus all increments)
	allApprovals := make([]*types.Balance, len(approvals))
	copy(allApprovals, approvals)
	
	for _, startAmount := range allApprovals {
		for _, time := range startAmount.Times {
			time.Start = time.Start.Add(numIncrements.Mul(incrementTimesBy))
			time.End = time.End.Add(numIncrements.Mul(incrementTimesBy))
		}

		for _, badgeId := range startAmount.BadgeIds {
			badgeId.Start = badgeId.Start.Add(numIncrements.Mul(incrementIdsBy))
			badgeId.End = badgeId.End.Add(numIncrements.Mul(incrementIdsBy))
		}
	}
	
	approvalTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(approvalTrackerDetails.Amounts, transferAmounts, allApprovals)
	if err != nil {
		return err
	}

	approvalTrackerDetails.NumTransfers = approvalTrackerDetails.NumTransfers.Add(sdk.NewUint(1))
	if approvalTrackerDetails.NumTransfers.GT(maxNumTransfers) {
		return ErrDisallowedTransfer
	}

	err = k.SetTransferTrackerInStore(ctx, collection.CollectionId, trackerId, approvalTrackerDetails, timelineType, depth, address)
	if err != nil {
		return ErrDisallowedTransfer
	}

	return nil
}