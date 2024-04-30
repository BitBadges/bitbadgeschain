package keeper

import (
	"encoding/json"
	"fmt"
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
func (k Keeper) DeductUserOutgoingApprovals(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	originalTransferBalances []*types.Balance,
	transfer *types.Transfer,
	from string,
	to string,
	requester string,
	userBalance *types.UserBalanceStore,
) error {
	currApprovals := userBalance.OutgoingApprovals
	if userBalance.AutoApproveSelfInitiatedOutgoingTransfers {
		currApprovals = AppendSelfInitiatedOutgoingApproval(currApprovals, from)
	}

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedApprovals := types.CastOutgoingTransfersToCollectionTransfers(currApprovals, from)
	_, err := k.DeductAndGetUserApprovals(
		ctx,
		collection,
		originalTransferBalances,
		transfer,
		castedApprovals,
		to,
		requester,
		"outgoing",
		from,
	)
	return err
}

// DeductUserIncomingApprovals will check if the current transfer is approved from the to's outgoing approvals and handle the approval tallying accordingly
func (k Keeper) DeductUserIncomingApprovals(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	originalTransferBalances []*types.Balance,
	transfer *types.Transfer,
	to string,
	initiatedBy string,
	userBalance *types.UserBalanceStore,
) error {
	currApprovals := userBalance.IncomingApprovals
	if userBalance.AutoApproveSelfInitiatedIncomingTransfers {
		currApprovals = AppendSelfInitiatedIncomingApproval(currApprovals, to)
	}

	//Little hack to reuse the same function for all transfer objects (we cast everything to a collection transfer)
	castedApprovals := types.CastIncomingTransfersToCollectionTransfers(currApprovals, to)
	_, err := k.DeductAndGetUserApprovals(
		ctx,
		collection,
		originalTransferBalances,
		transfer,
		castedApprovals,
		to,
		initiatedBy,
		"incoming",
		to,
	)
	return err
}

// DeductCollectionApprovalsAndGetUserApprovalsToCheck will check if the current transfer is allowed via the collection's approved transfers and handle any tallying accordingly
func (k Keeper) DeductCollectionApprovalsAndGetUserApprovalsToCheck(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	transfer *types.Transfer,
	toAddress string,
	initiatedBy string,
) ([]*UserApprovalsToCheck, error) {
	return k.DeductAndGetUserApprovals(
		ctx,
		collection,
		transfer.Balances,
		transfer,
		collection.CollectionApprovals,
		toAddress,
		initiatedBy,
		"collection",
		"",
	)
}

func (k Keeper) DeductAndGetUserApprovals(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	originalTransferBalances []*types.Balance,
	transfer *types.Transfer,
	_approvals []*types.CollectionApproval,
	toAddress string,
	initiatedBy string,
	approvalLevel string,
	approverAddress string,
) ([]*UserApprovalsToCheck, error) {
	prioritizedApprovals := transfer.PrioritizedApprovals
	fromAddress := transfer.From
	onlyCheckPrioritized := false
	if approvalLevel == "collection" && transfer.OnlyCheckPrioritizedCollectionApprovals {
		onlyCheckPrioritized = true
	} else if approvalLevel == "outgoing" && transfer.OnlyCheckPrioritizedOutgoingApprovals {
		onlyCheckPrioritized = true
	} else if approvalLevel == "incoming" && transfer.OnlyCheckPrioritizedIncomingApprovals {
		onlyCheckPrioritized = true
	}

	zkProofSolutions := transfer.ZkProofSolutions
	originalTransferBalances = types.DeepCopyBalances(originalTransferBalances)

	//Reorder approvals based on prioritized approvals
	//We want to check prioritized approvals first
	//If onlyCheckPrioritized is true, we only check prioritized approvals and ignore the rest
	approvals := []*types.CollectionApproval{}
	prioritizedTransfers := []*types.CollectionApproval{}
	for _, approval := range _approvals {
		prioritized := false

		for _, prioritizedApproval := range prioritizedApprovals {
			if approval.ApprovalId == prioritizedApproval.ApprovalId && prioritizedApproval.ApprovalLevel == approvalLevel && approverAddress == prioritizedApproval.ApproverAddress {
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

	remainingBalances := types.DeepCopyBalances(transfer.Balances) //Keep a running tally of all the badges we still have to handle

	//For each approved transfer, we check if the transfer is allowed
	//1: If transfer meets all criteria, we deduct, get user approvals to check, and continue (if there are any remaining balances)
	//2. If transfer does not meet all criteria, we continue and do not mark anything as handled
	//3. At the end, if there are any unhandled transfers, we throw (not enough approvals = transfer disallowed)
	userApprovalsToCheck := []*UserApprovalsToCheck{}
	for _, approval := range approvals {
		approvalId := approval.ApprovalId
		remainingBalances = types.FilterZeroBalances(remainingBalances)
		if len(remainingBalances) == 0 {
			break
		}

		//Initial checks: Make sure (from, to, initiatedBy) match the approval's collection list IDs
		//								Make sure the current time is within the approval's transfer times
		doAddressesMatch := k.CheckIfAddressesMatchCollectionListIds(ctx, approval, fromAddress, toAddress, initiatedBy)
		if !doAddressesMatch {
			continue
		}

		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currTimeFound, err := types.SearchUintRangesForUint(currTime, approval.TransferTimes)
		if !currTimeFound || err != nil {
			continue
		}

		transferStr := "(from: " + fromAddress + ", to: " + toAddress + ", initiatedBy: " + initiatedBy + ", badgeId: " + approval.BadgeIds[0].Start.String() + ", time: " + approval.TransferTimes[0].Start.String() + ", ownershipTime: " + approval.OwnershipTimes[0].Start.String() + ")"

		//If there are no restrictions or criteria, it is a full match and we can deduct all the overlapping (badgeIds, ownershipTimes) from the remaining balances
		if approval.ApprovalCriteria == nil {
			allBalancesForIdsAndTimes, err := types.GetBalancesForIds(ctx, approval.BadgeIds, approval.OwnershipTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for transfer: %s", transferStr)
			}

			remainingBalances, err = types.SubtractBalances(ctx, allBalancesForIdsAndTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: underflow error subtracting balances for transfer: %s", transferStr)
			}
		}

		//Else, we have a match and we can proceed to check the restrictions
		//This is split into a two part process:
		//	1. Simulate the transfer to see if it is allowed
		//	2. If all simulations pass, we deduct as much as possible from the approval
		if approval.ApprovalCriteria != nil {
			/**** SECTION 1: NO STORAGE WRITES (just simulate everything and continue if it doesn't pass) ****/

			approvalCriteria := approval.ApprovalCriteria
			if approvalCriteria.RequireFromDoesNotEqualInitiatedBy && fromAddress == initiatedBy {
				continue
			}

			if approvalCriteria.RequireFromEqualsInitiatedBy && fromAddress != initiatedBy {
				continue
			}

			if approvalCriteria.RequireToDoesNotEqualInitiatedBy && toAddress == initiatedBy {
				continue
			}

			if approvalCriteria.RequireToEqualsInitiatedBy && toAddress != initiatedBy {
				continue
			}

			err := k.HandleCoinTransfers(ctx, approvalCriteria.CoinTransfers, initiatedBy, true) //simulate = true
			if err != nil {
				continue
			}

			err = k.CheckMustOwnBadges(ctx, approvalCriteria.MustOwnBadges, initiatedBy)
			if err != nil {
				continue
			}

			validZKPSolutionIdxs, err := k.SimulateZKPs(ctx, collection, approvalCriteria.ZkProofs, zkProofSolutions, approverAddress, approvalLevel, approvalId)
			if err != nil {
				continue
			}

			//Get max balances allowed for this approvalCriteria element
			//Get the max balances allowed for this approvalCriteria element WITHOUT incrementing
			transferBalancesToCheck, err := types.GetBalancesForIds(ctx, approval.BadgeIds, approval.OwnershipTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for transfer: %s", transferStr)
			}

			transferBalancesToCheck = types.FilterZeroBalances(transferBalancesToCheck)
			if len(transferBalancesToCheck) == 0 {
				continue
			}

			challengeNumIncrements, err := k.HandleMerkleChallenges(
				ctx,
				collection.CollectionId,
				transfer,
				approval,
				initiatedBy,
				approverAddress,
				approvalLevel,
				true, //simulation = true
			)
			if err != nil {
				continue
			}

			trackerTypes := []string{"overall", "to", "from", "initiatedBy"}
			approvedAmounts := []sdkmath.Uint{
				approvalCriteria.ApprovalAmounts.OverallApprovalAmount,
				approvalCriteria.ApprovalAmounts.PerToAddressApprovalAmount,
				approvalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount,
				approvalCriteria.ApprovalAmounts.PerInitiatedByAddressApprovalAmount,
			}
			maxNumTransfers := []sdkmath.Uint{
				approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers,
				approvalCriteria.MaxNumTransfers.PerToAddressMaxNumTransfers,
				approvalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers,
				approvalCriteria.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers,
			}
			approvedAddresses := []string{"", toAddress, fromAddress, initiatedBy}

			//Get max balances allowed for this approvalCriteria element
			failed := false
			for i, trackerType := range trackerTypes {
				//Get max allowed by criteria
				maxPossible, err := k.GetMaxPossible(
					ctx,
					collection,
					approval,
					transfer,
					originalTransferBalances,
					approvedAmounts[i],
					challengeNumIncrements,
					approverAddress,
					approvalLevel,
					trackerType,
					approvedAddresses[i],
					true,
				)
				if err != nil {
					failed = true
					break
				}

				//Get max allowed by remaining balances to check
				transferBalancesToCheck, err = types.GetOverlappingBalances(ctx, maxPossible, transferBalancesToCheck)
				if err != nil {
					failed = true
					break
				}
			}
			if failed {
				continue
			}

			transferBalancesToCheck = types.FilterZeroBalances(transferBalancesToCheck)
			if len(transferBalancesToCheck) == 0 {
				continue
			}

			//here, we assert that the transfer can be incremented and is within the threshold for all trackers (this is a simulation)
			for i, trackerType := range trackerTypes {
				err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, approval, originalTransferBalances, approvedAmounts[i], maxNumTransfers[i], transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, trackerType, approvedAddresses[i], true)
				if err != nil {
					failed = true
					break
				}
			}
			if failed {
				continue
			}

			/**** SECTION 2: ONCE HERE, EVERYTHING BELOW SHOULD BE SUCCESSFUL BC IT WAS SIMULATED ****/
			remainingBalances, err = types.SubtractBalances(ctx, transferBalancesToCheck, remainingBalances)
			if err != nil {
				continue
			}

			err = k.HandleCoinTransfers(ctx, approvalCriteria.CoinTransfers, initiatedBy, false) //simulate = false
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error handling coin transfers")
			}

			//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
			//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
			//    If so, useLeafIndexForNumIncrements will be true
			challengeNumIncrements, err = k.HandleMerkleChallenges(
				ctx,
				collection.CollectionId,
				transfer,
				approval,
				initiatedBy,
				approverAddress,
				approvalLevel,
				false, //simulation = false
			)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "%s", transferStr)
			}

			err = k.HandleZKPs(ctx, collection, validZKPSolutionIdxs, approvalCriteria.ZkProofs, zkProofSolutions, approverAddress, approvalLevel, approvalId)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error handling zk proofs")
			}

			for i, trackerType := range trackerTypes {
				err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, approval, originalTransferBalances, approvedAmounts[i], maxNumTransfers[i], transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, trackerType, approvedAddresses[i], false)
				if err != nil {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing approvals")
				}
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
		//convert ownership time unix milliseconds number to string
		timeToConvert := remainingBalances[0].OwnershipTimes[0].Start //unix milliseconds
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
	collection *types.BadgeCollection,
	approval *types.CollectionApproval,
	transfer *types.Transfer,
	originalTransferBalances []*types.Balance,
	approvedAmount sdkmath.Uint,
	challengeNumIncrements sdkmath.Uint,
	approverAddress string,
	approvalLevel string,
	trackerType string,
	address string,
	simulate bool,
) ([]*types.Balance, error) {
	transferBalances := types.DeepCopyBalances(transfer.Balances)
	approvalCriteria := approval.ApprovalCriteria
	amountsTrackerId := approvalCriteria.ApprovalAmounts.AmountTrackerId

	predeterminedBalances := approvalCriteria.PredeterminedBalances
	allApprovals := []*types.Balance{{
		Amount:         approvedAmount,
		OwnershipTimes: approval.OwnershipTimes,
		BadgeIds:       approval.BadgeIds,
	}}

	//Get the current approvals for this transfer
	//If nil, no restrictions and we are approved for the entire transfer
	//Note we filter any excess badge IDs later and apply num increments as well
	if approvedAmount.IsNil() {
		approvedAmount = sdkmath.NewUint(0)
	}

	needToFetchApprovalTrackerDetails := approvedAmount.GT(sdkmath.NewUint(0)) ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers && trackerType == "overall") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy")

	amountsTrackerDetails := types.ApprovalTracker{
		Amounts:      []*types.Balance{},
		NumTransfers: sdkmath.NewUint(0),
	}

	if needToFetchApprovalTrackerDetails {
		fetchedDetails, found := k.GetApprovalTrackerFromStore(ctx, collection.CollectionId, approverAddress, approval.ApprovalId, amountsTrackerId, approvalLevel, trackerType, address)
		if found {
			amountsTrackerDetails = fetchedDetails
		}
	}

	//Get the max amount that we can add from transferBalances without exceeding the amount threshold
	if approvedAmount.GT(sdkmath.NewUint(0)) {
		currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(ctx, approval.BadgeIds, approval.OwnershipTimes, amountsTrackerDetails.Amounts)
		if err != nil {
			return nil, err
		}

		maxBalancesWeCanAdd, err := types.SubtractBalances(ctx, currTallyForCurrentIdsAndTimes, allApprovals)
		if err != nil {
			return nil, err
		}

		return maxBalancesWeCanAdd, nil
	} else {
		return transferBalances, nil
	}
}

func (k Keeper) IncrementApprovalsAndAssertWithinThreshold(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	approval *types.CollectionApproval,
	originalTransferBalances []*types.Balance,
	approvedAmount sdkmath.Uint,
	maxNumTransfers sdkmath.Uint,
	transferBalances []*types.Balance,
	challengeNumIncrements sdkmath.Uint,
	approverAddress string,
	approvalLevel string,
	trackerType string,
	address string,
	simulate bool,
) error {
	approvalCriteria := approval.ApprovalCriteria
	amountsTrackerId := approvalCriteria.ApprovalAmounts.AmountTrackerId
	maxNumTransfersTrackerId := approvalCriteria.MaxNumTransfers.AmountTrackerId
	predeterminedBalances := approvalCriteria.PredeterminedBalances
	thresholdAmounts := []*types.Balance{{
		Amount:         approvedAmount,
		OwnershipTimes: approval.OwnershipTimes,
		BadgeIds:       approval.BadgeIds,
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

	amountsTrackerDetails := types.ApprovalTracker{
		Amounts:      []*types.Balance{},
		NumTransfers: sdkmath.NewUint(0),
	}

	maxNumTransfersTrackerDetails := types.ApprovalTracker{
		Amounts:      []*types.Balance{},
		NumTransfers: sdkmath.NewUint(0),
	}

	if needToFetchApprovalTrackerDetails {
		fetchedDetails, found := k.GetApprovalTrackerFromStore(ctx, collection.CollectionId, approverAddress, approval.ApprovalId, amountsTrackerId, approvalLevel, trackerType, address)
		if found {
			amountsTrackerDetails = fetchedDetails
		}

		fetchedDetails, found = k.GetApprovalTrackerFromStore(ctx, collection.CollectionId, approverAddress, approval.ApprovalId, maxNumTransfersTrackerId, approvalLevel, trackerType, address)
		if found {
			maxNumTransfersTrackerDetails = fetchedDetails
		}
	}

	//Here, we check that if predeterminedBalances is set and non-nil, the ORIGINAL transfer balances must be exactly as predetermined
	//Typically, the original transfer balances are the balances that are being transferred, but in some cases, the approval may only be applicable
	//to a subset of the transfer balances. In these cases, we check the original and only approve the subset (the rest must be approved by a different approval)
	if predeterminedBalances != nil {
		numIncrements := sdkmath.NewUint(0)
		toBeCalculated := true
		trackerNumTransfers := maxNumTransfersTrackerDetails.NumTransfers

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

			//Assert that we have exactly the amount specified in the original transfers
			equal := types.AreBalancesEqual(ctx, originalTransferBalances, calculatedBalances, false)
			if !equal {
				return sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because predetermined balances do not match: %s", approval.String())
			}
		}
	}

	//Increment amounts and add back to store
	if approvedAmount.GT(sdkmath.NewUint(0)) {
		currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(ctx, approval.BadgeIds, approval.OwnershipTimes, amountsTrackerDetails.Amounts)
		if err != nil {
			return err
		}

		//If this passes, the new transferBalances are okay
		_, err = types.AddBalancesAndAssertDoesntExceedThreshold(ctx, currTallyForCurrentIdsAndTimes, transferBalances, thresholdAmounts)
		if err != nil {
			return err
		}

		//We then add them to the current tally of ALL ids and times
		amountsTrackerDetails.Amounts, err = types.AddBalances(ctx, amountsTrackerDetails.Amounts, transferBalances)
		if err != nil {
			return err
		}
	}

	//We need to increment if we are using it to assign predetermined balances or there is an explicit num transfers limit
	if maxNumTransfers.GT(sdkmath.NewUint(0)) ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers && trackerType == "overall") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy") {

		maxNumTransfersTrackerDetails.NumTransfers = maxNumTransfersTrackerDetails.NumTransfers.Add(sdkmath.NewUint(1))
		//only check exceeds if maxNumTransfers is not 0 (because 0 means no limit)
		if maxNumTransfers.GT(sdkmath.NewUint(0)) {
			if maxNumTransfersTrackerDetails.NumTransfers.GT(maxNumTransfers) {
				return sdkerrors.Wrapf(ErrDisallowedTransfer, "exceeded max transfers allowed - %s", maxNumTransfers.String())
			}
		}
	}

	//
	if needToFetchApprovalTrackerDetails && !simulate {
		//Currently added for indexer, but note that it is planned to be deprecated
		amountsJsonData, err := json.Marshal(amountsTrackerDetails.Amounts)
		if err != nil {
			return err
		}
		amountsStr := string(amountsJsonData)

		numTransfersJsonData, err := json.Marshal(maxNumTransfersTrackerDetails.NumTransfers)
		if err != nil {
			return err
		}
		numTransfersStr := string(numTransfersJsonData)

		amountsNumTransfersJsonData, err := json.Marshal(amountsTrackerDetails.NumTransfers)
		if err != nil {
			return err
		}
		amountsNumTransfersStr := string(amountsNumTransfersJsonData)

		maxNumTransfersAmountsJsonData, err := json.Marshal(maxNumTransfersTrackerDetails.Amounts)
		if err != nil {
			return err
		}
		maxNumTransfersAmountsStr := string(maxNumTransfersAmountsJsonData)

		//TODO: Deprecate this eventually in favor of doing exclusively in indexer

		//If same ID, we only need to set the tracker once (but make sure to include both the amounts and the num transfers bc we handled them separately above)
		//If not, we need to set both trackers individually
		isSameId := amountsTrackerId == maxNumTransfersTrackerId
		if isSameId {
			ctx.EventManager().EmitEvent(
				sdk.NewEvent("approval"+fmt.Sprint(collection.CollectionId)+fmt.Sprint(approverAddress)+fmt.Sprint(approval.ApprovalId)+fmt.Sprint(amountsTrackerId)+fmt.Sprint(approvalLevel)+fmt.Sprint(trackerType)+fmt.Sprint(address),
					sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
					sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
					sdk.NewAttribute("approvalId", fmt.Sprint(approval.ApprovalId)),
					sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
					sdk.NewAttribute("amountTrackerId", fmt.Sprint(amountsTrackerId)),
					sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
					sdk.NewAttribute("trackerType", fmt.Sprint(trackerType)),
					sdk.NewAttribute("approvedAddress", fmt.Sprint(address)),
					sdk.NewAttribute("amounts", amountsStr),
					sdk.NewAttribute("numTransfers", numTransfersStr),
				),
			)

			amountsTrackerDetails.NumTransfers = maxNumTransfersTrackerDetails.NumTransfers
			err = k.SetApprovalTrackerInStore(ctx, collection.CollectionId, approverAddress, approval.ApprovalId, amountsTrackerId, amountsTrackerDetails, approvalLevel, trackerType, address)
			if err != nil {
				return err
			}
		} else {

			ctx.EventManager().EmitEvent(
				sdk.NewEvent("approval"+fmt.Sprint(collection.CollectionId)+fmt.Sprint(approverAddress)+fmt.Sprint(approval.ApprovalId)+fmt.Sprint(amountsTrackerId)+fmt.Sprint(approvalLevel)+fmt.Sprint(trackerType)+fmt.Sprint(address),
					sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
					sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
					sdk.NewAttribute("approvalId", fmt.Sprint(approval.ApprovalId)),
					sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
					sdk.NewAttribute("amountTrackerId", fmt.Sprint(amountsTrackerId)),
					sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
					sdk.NewAttribute("trackerType", fmt.Sprint(trackerType)),
					sdk.NewAttribute("approvedAddress", fmt.Sprint(address)),
					sdk.NewAttribute("amounts", amountsStr),
					sdk.NewAttribute("numTransfers", amountsNumTransfersStr),
				),
			)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent("approval"+fmt.Sprint(collection.CollectionId)+fmt.Sprint(approverAddress)+fmt.Sprint(approval.ApprovalId)+fmt.Sprint(maxNumTransfersTrackerId)+fmt.Sprint(approvalLevel)+fmt.Sprint(trackerType)+fmt.Sprint(address),
					sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
					sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
					sdk.NewAttribute("approvalId", fmt.Sprint(approval.ApprovalId)),
					sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
					sdk.NewAttribute("amountTrackerId", fmt.Sprint(maxNumTransfersTrackerId)),
					sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
					sdk.NewAttribute("trackerType", fmt.Sprint(trackerType)),
					sdk.NewAttribute("approvedAddress", fmt.Sprint(address)),
					sdk.NewAttribute("amounts", maxNumTransfersAmountsStr),
					sdk.NewAttribute("numTransfers", numTransfersStr),
				),
			)

			err = k.SetApprovalTrackerInStore(ctx, collection.CollectionId, approverAddress, approval.ApprovalId, amountsTrackerId, amountsTrackerDetails, approvalLevel, trackerType, address)
			if err != nil {
				return err
			}

			err = k.SetApprovalTrackerInStore(ctx, collection.CollectionId, approverAddress, approval.ApprovalId, maxNumTransfersTrackerId, maxNumTransfersTrackerDetails, approvalLevel, trackerType, address)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (k Keeper) GetPredeterminedBalancesForPrecalculationId(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	approvals []*types.CollectionApproval,
	transfer *types.Transfer,
	precalcDetails *types.ApprovalIdentifierDetails,
	to string,
	initiatedBy string,
) ([]*types.Balance, error) {
	approvalId := ""
	approverAddress := precalcDetails.ApproverAddress
	approvalLevel := precalcDetails.ApprovalLevel
	precalculationId := precalcDetails.ApprovalId

	for _, approval := range approvals {
		maxNumTransfersTrackerId := approval.ApprovalCriteria.MaxNumTransfers.AmountTrackerId
		approvalCriteria := approval.ApprovalCriteria
		approvalId = approval.ApprovalId
		if approvalCriteria == nil {
			continue
		}

		if approvalId != precalculationId {
			continue
		}

		if approvalId == "" {
			continue
		}

		if approvalCriteria.PredeterminedBalances != nil {
			numIncrements := sdkmath.NewUint(0)
			if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {

				//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
				//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
				numIncrementsFetched, err := k.HandleMerkleChallenges(
					ctx,
					collection.CollectionId,
					transfer,
					approval,
					initiatedBy,
					approverAddress,
					approvalLevel,
					true, //simulation = true
				)
				if err != nil {
					return []*types.Balance{}, sdkerrors.Wrapf(err, "invalid challenges / solutions")
				}

				numIncrements = numIncrementsFetched
			} else {
				trackerType := "overall"
				approvedAddress := ""
				if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers {
					trackerType = "from"
					approvedAddress = transfer.From
				} else if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers {
					trackerType = "to"
					approvedAddress = to
				} else if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers {
					trackerType = "initiatedBy"
					approvedAddress = initiatedBy
				}

				numTransfersTracker, found := k.GetApprovalTrackerFromStore(
					ctx,
					collection.CollectionId,
					approverAddress,
					approval.ApprovalId,
					maxNumTransfersTrackerId,
					approvalLevel,
					trackerType,
					approvedAddress,
				)
				if !found {
					numTransfersTracker = types.ApprovalTracker{
						Amounts:      []*types.Balance{},
						NumTransfers: sdkmath.NewUint(0),
					}
				}

				numIncrements = numTransfersTracker.NumTransfers
			}

			//calculate the current approved balances from the numIncrements and predeterminedBalances
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
