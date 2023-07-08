package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

// The UserApprovalsToCheck struct is used to keep track of which incoming / outgoing approvals for which addresses we need to check.
type UserApprovalsToCheck struct {
	Address  string
	BadgeIds []*types.UintRange
	Outgoing bool
}

// DeductUserOutgoingApprovals will check if the current transfer is approved from the from's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserOutgoingApprovals(ctx sdk.Context, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.ChallengeSolution) error {
	currApprovedTransfers := types.GetCurrentUserApprovedOutgoingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForOutgoing(currApprovedTransfers, from)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastOutgoingTransfersToCollectionTransfers(currApprovedTransfers, from)
	_, err := k.DeductAndGetUserApprovals(castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "outgoing")
	return err
}

// DeductUserIncomingApprovals will check if the current transfer is approved from the to's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserIncomingApprovals(ctx sdk.Context, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.ChallengeSolution) error {
	currApprovedTransfers := types.GetCurrentUserApprovedIncomingTransfers(ctx, userBalance)
	currApprovedTransfers = AppendDefaultForIncoming(currApprovedTransfers, to)

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedTransfers := types.CastIncomingTransfersToCollectionTransfers(currApprovedTransfers, to)
	_, err := k.DeductAndGetUserApprovals(castedTransfers, ctx, collection, badgeIds, times, from, to, requester, amount, solutions, "incoming")
	return err
}

// DeductCollectionApprovalsAndGetUserApprovalsToCheck will check if the current transfer is allowed via the collection's approved transfers and handle any tallying accordingly
func (k Keeper) DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx sdk.Context, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.ChallengeSolution) ([]*UserApprovalsToCheck, error) {
	approvedTransfers := types.GetCurrentCollectionApprovedTransfers(ctx, collection)
	return k.DeductAndGetUserApprovals(approvedTransfers, ctx, collection, badgeIds, times, fromAddress, toAddress, initiatedBy, amount, solutions, "overall")
}

func (k Keeper) DeductAndGetUserApprovals(approvedTransfers []*types.CollectionApprovedTransfer, ctx sdk.Context, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromAddress string, toAddress string, initiatedBy string, amount sdkmath.Uint, solutions []*types.ChallengeSolution, timelineType string) ([]*UserApprovalsToCheck, error) {
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
	unhandledBadgeIds := make([]*types.UintRange, len(badgeIds))
	copy(unhandledBadgeIds, badgeIds)

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

		remaining, overlaps := types.RemoveUintRangeFromUintRange([]*types.UintRange{match.BadgeId}, unhandledBadgeIds)
		unhandledBadgeIds = remaining
		if len(overlaps) > 0 {
			//For the overlapping badges, we have a match because mapping IDs, time, and badge IDs match.
			//We can now proceed to check any restrictions.
			//If any restriction fails in this if statement, we MUST throw because the transfer is invalid for the badge IDs (since we use first match only)

			transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + overlaps[0].Start.String() + ", time: " + currTime.String() + ")" //for error msgs

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
			useLeafIndexForNumIncrements, numIncrements, err := k.AssertValidSolutionForEveryChallenge(ctx, collection.CollectionId, transferVal.Challenges, solutions, initiatedBy, "overall", transferVal.ApprovalId)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed because of invalid challenges / solutions: %s", transferStr)
			}

			transferBalance := &types.Balance{Amount: amount, OwnershipTimes: times, BadgeIds: overlaps}

			if transferVal.ApprovalId == "" && (transferVal.OverallApprovals != nil || transferVal.PerAddressApprovals != nil) {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because approvalId cannot be blank: %s", transferStr)
			}

			//Handle the incrementing of overall approvals
			if transferVal.OverallApprovals != nil {
				err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.OverallApprovals.Amounts, transferVal.OverallApprovals.NumTransfers, transferBalance, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "overall", "")
				if err != nil {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing overall approvals: %s", transferStr)
				}
			}

			//Handle the per-address approvals
			if transferVal.PerAddressApprovals != nil {
				if transferVal.PerAddressApprovals.ApprovalsPerToAddress != nil {
					err := k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.PerAddressApprovals.ApprovalsPerToAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers, transferBalance, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "per-to", toAddress)
					if err != nil {
						return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing per-to approvals: %s", transferStr)
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerFromAddress != nil {
					err := k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerFromAddress.NumTransfers, transferBalance, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "per-from", fromAddress)
					if err != nil {
						return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing per-from approvals: %s", transferStr)
					}
				}

				if transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress != nil {
					err := k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress.Amounts, transferVal.PerAddressApprovals.ApprovalsPerInitiatedByAddress.NumTransfers, transferBalance, transferVal.ApprovalId, transferVal.IncrementBadgeIdsBy, transferVal.IncrementOwnershipTimesBy, useLeafIndexForNumIncrements, numIncrements, timelineType, "per-initiated-by", initiatedBy)
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
					BadgeIds: overlaps,
					Outgoing: true,
				})
			}

			if !transferVal.OverridesToApprovedIncomingTransfers {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address:  toAddress,
					BadgeIds: overlaps,
					Outgoing: false,
				})
			}
		}
	}

	//If all are not explicitly allowed, we return that it is disallowed by default
	if len(unhandledBadgeIds) > 0 {
		return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrInadequateApprovals, "transfer disallowed because no approved transfer was found for badge ids: %v", unhandledBadgeIds)
	}

	return userApprovalsToCheck, nil
}

func AssertBalancesDoNotExceedThreshold(balancesToCheck []*types.Balance, threshold []*types.Balance) error {
	err := *new(error)

	//Check if we exceed the threshold; will underflow if we do exceed it
	thresholdCopy := make([]*types.Balance, len(threshold))
	copy(thresholdCopy, threshold)
	for _, balance := range balancesToCheck {
		thresholdCopy, err = types.SubtractBalance(thresholdCopy, balance)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddTallyAndAssertDoesntExceedThreshold(currTally []*types.Balance, toAdd *types.Balance, threshold []*types.Balance) ([]*types.Balance, error) {
	//Add the new tally to existing
	err := *new(error)
	currTally, err = types.AddBalance(currTally, toAdd)
	if err != nil {
		return []*types.Balance{}, err
	}

	//Check if we exceed the threshold; will underflow if we do exceed it
	err = AssertBalancesDoNotExceedThreshold(currTally, threshold)
	return currTally, err
}

func (k Keeper) IncrementApprovalsAndAssertWithinThreshold(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	approvals []*types.Balance,
	maxNumTransfers sdkmath.Uint,
	transferAmounts *types.Balance,
	approvalId string,
	incrementBadgeIdsBy sdkmath.Uint,
	incrementOwnershipTimesBy sdkmath.Uint,
	precalculatedNumIncrements bool,
	numIncrements sdkmath.Uint,
	timelineType string,
	depth string,
	address string,
) error {
	//Get the current approvals for this transfer
	//If nil, no restrictions and we are approved for the entire transfer
	//Note we filter any excess badge IDs later and apply num increments as well
	err := *new(error)
	if approvals == nil {
		return sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because no approval amounts were found")
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
	allApprovals := make([]*types.Balance, len(approvals))
	copy(allApprovals, approvals)

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
	approvalTrackerDetails.Amounts, err = AddTallyAndAssertDoesntExceedThreshold(approvalTrackerDetails.Amounts, transferAmounts, allApprovals)
	if err != nil {
		return err
	}

	approvalTrackerDetails.NumTransfers = approvalTrackerDetails.NumTransfers.Add(sdkmath.NewUint(1))
	if approvalTrackerDetails.NumTransfers.GT(maxNumTransfers) {
		return sdkerrors.Wrapf(ErrDisallowedTransfer, "exceeded max num transfers - %s", maxNumTransfers.String())
	}

	err = k.SetApprovalsTrackerInStore(ctx, collection.CollectionId, approvalId, approvalTrackerDetails, timelineType, depth, address)
	if err != nil {
		return err
	}

	return nil
}
