package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

// The UserApprovalsToCheck struct is used to keep track of which incoming / outgoing approvals for which addresses we need to check.
type UserApprovalsToCheck struct {
	Address       string
	Balances      []*types.Balance
	Outgoing      bool
	UserRoyalties *types.UserRoyalties
}

// All in one approval deduction function. We also return the user approvals to check (only used when collection approvals)
func (k Keeper) DeductAndGetUserApprovals(
	ctx sdk.Context,
	collection *types.TokenCollection,
	originalTransferBalances []*types.Balance,
	transfer *types.Transfer,
	_approvals []*types.CollectionApproval,
	transferMetadata TransferMetadata,
	eventTracking *EventTracking,
	royalties *types.UserRoyalties,
	approvalLevel string,
) ([]*UserApprovalsToCheck, error) {
	fromAddress := transferMetadata.From
	toAddress := transferMetadata.To
	initiatedBy := transferMetadata.InitiatedBy
	// Use the approvalLevel parameter passed in (not from transferMetadata)
	approverAddress := transferMetadata.ApproverAddress

	originalTransferBalances = types.DeepCopyBalances(originalTransferBalances)
	remainingBalances := types.DeepCopyBalances(transfer.Balances) //Keep a running tally of all the tokens we still have to handle
	approvals, err := FilterApprovalsWithPrioritizedHandling(_approvals, transfer, approvalLevel, approverAddress)
	if err != nil {
		return []*UserApprovalsToCheck{}, err
	}

	potentialErrors := []string{}
	errorsWithIdx := []ErrorWithIdx{}

	addToPotentialErrors := func(errorMsg string, idx int) {
		exists := false
		for _, errorWithIdx := range errorsWithIdx {
			if errorWithIdx.Idx == idx {
				errorWithIdx.ErrorMsgs = append(errorWithIdx.ErrorMsgs, errorMsg)
				exists = true
				break
			}
		}
		if !exists {
			errorsWithIdx = append(errorsWithIdx, ErrorWithIdx{ErrorMsgs: []string{errorMsg}, Idx: idx})
		}
	}

	// Helper function to add potential errors for prioritized approvals
	addPotentialError := func(isExplicitlyPrioritized bool, idx int, errorMsg string) {
		addToPotentialErrors(errorMsg, idx)

		if isExplicitlyPrioritized {
			potentialErrors = append(potentialErrors, errorMsg)
		}
	}

	approvalIdxsChecked := []int{} //We keep track of the indexes of the approvals that we have checked

	//For each approved transfer, we check if the transfer is allowed
	// 1: If transfer meets all criteria, we deduct, get user approvals to check, and continue (if there are any remaining balances)
	// 2. If transfer does not meet all criteria, we continue and do not mark anything as handled
	// 3. At the end, if there are any unhandled transfers, we throw (not enough approvals = transfer disallowed)
	userApprovalsToCheck := []*UserApprovalsToCheck{}
	for idx, approval := range approvals {
		remainingBalances = types.FilterZeroBalances(remainingBalances)
		if len(remainingBalances) == 0 {
			break
		}

		userRoyalties := &types.UserRoyalties{
			Percentage:    sdkmath.NewUint(0),
			PayoutAddress: "",
		}
		if approval.ApprovalCriteria != nil && approval.ApprovalCriteria.UserRoyalties != nil {
			userRoyalties = approval.ApprovalCriteria.UserRoyalties
		}

		isExplicitlyPrioritized := false
		for _, prioritizedApproval := range transfer.PrioritizedApprovals {
			if prioritizedApproval.ApprovalId == approval.ApprovalId && prioritizedApproval.ApprovalLevel == approvalLevel && prioritizedApproval.ApproverAddress == approverAddress {
				isExplicitlyPrioritized = true
				break
			}
		}

		//Initial checks: Make sure (from, to, initiatedBy) match the approval's collection list IDs
		doAddressesMatch := k.CheckIfAddressesMatchCollectionListIds(ctx, approval, fromAddress, toAddress, initiatedBy)
		if !doAddressesMatch {
			addPotentialError(isExplicitlyPrioritized, idx, fmt.Sprintf("addresses do not match (from, to, initiatedBy): %s, %s, %s", fromAddress, toAddress, initiatedBy))
			continue
		}

		// Check valid transfer times
		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currTimeFound, err := types.SearchUintRangesForUint(currTime, approval.TransferTimes)
		if !currTimeFound || err != nil {
			addPotentialError(isExplicitlyPrioritized, idx, "transfer time not in range")
			continue
		}

		// Only after we sanity check the addresses and times
		approvalIdxsChecked = append(approvalIdxsChecked, idx)

		transferStr := "attempting to transfer ID " + approval.TokenIds[0].Start.String()

		markApprovalAsUsed := func(approval *types.CollectionApproval) {
			*eventTracking.ApprovalsUsed = append(*eventTracking.ApprovalsUsed, ApprovalsUsed{
				ApprovalId:      approval.ApprovalId,
				ApprovalLevel:   approvalLevel,
				ApproverAddress: approverAddress,
				Version:         approval.Version.String(),
			})
		}

		addToUserApprovalsToCheck := func(address string, balances []*types.Balance, outgoing bool, userRoyalties *types.UserRoyalties) {
			userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
				Address:       address,
				Balances:      balances,
				Outgoing:      outgoing,
				UserRoyalties: userRoyalties,
			})
		}

		//If there are no restrictions or criteria, it is a full match and we can deduct all the overlapping (tokenIds, ownershipTimes) from the remaining balances
		approvalCriteria := approval.ApprovalCriteria
		if approvalCriteria == nil {
			allBalancesForIdsAndTimes, err := types.GetBalancesForIds(ctx, approval.TokenIds, approval.OwnershipTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for transfer: %s", transferStr)
			}

			remainingBalances, err = types.SubtractBalances(ctx, allBalancesForIdsAndTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: underflow error subtracting balances for transfer: %s", transferStr)
			}

			//If we do not override the approved outgoing / incoming transfers, we need to check the user approvals
			addToUserApprovalsToCheck(fromAddress, allBalancesForIdsAndTimes, true, userRoyalties)
			addToUserApprovalsToCheck(toAddress, allBalancesForIdsAndTimes, false, userRoyalties)
			markApprovalAsUsed(approval)
		} else {
			// Create a cached context - we'll perform all operations on this and only commit if everything succeeds
			cachedCtx, writeCache := ctx.CacheContext()

			// Run all applicable checkers dynamically (includes basic validation and approval criteria)
			checkers := k.GetApprovalCriteriaCheckers(approval)
			checkerFailed := false
			for _, checker := range checkers {
				detErrMsg, err := checker.Check(cachedCtx, approval, collection, toAddress, fromAddress, initiatedBy, approvalLevel, approverAddress, transfer.MerkleProofs, transfer.EthSignatureProofs, transfer.Memo, isExplicitlyPrioritized)
				if err != nil {
					checkerName := checker.Name()
					if detErrMsg != "" {
						addPotentialError(isExplicitlyPrioritized, idx, fmt.Sprintf("%s: %s", checkerName, detErrMsg))
					} else {
						addPotentialError(isExplicitlyPrioritized, idx, fmt.Sprintf("%s: approval criteria check failed", checkerName))
					}
					checkerFailed = true
					break
				}
			}
			if checkerFailed {
				continue
			}

			//Get max balances allowed for this approvalCriteria element
			transferBalancesToCheck, err := types.GetBalancesForIds(cachedCtx, approval.TokenIds, approval.OwnershipTimes, remainingBalances)
			if err != nil {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "transfer disallowed: err fetching balances for transfer: %s", transferStr)
			}

			transferBalancesToCheck = types.FilterZeroBalances(transferBalancesToCheck)
			if len(transferBalancesToCheck) == 0 {
				addPotentialError(isExplicitlyPrioritized, idx, "no balances to check")
				continue
			}

			// Handle coin transfers on cached context
			detErrMsg, err := k.ExecuteCoinTransfers(cachedCtx, approvalCriteria.CoinTransfers, transferMetadata, eventTracking.CoinTransfers, collection, royalties)
			if err != nil {
				if detErrMsg != "" {
					addPotentialError(isExplicitlyPrioritized, idx, detErrMsg)
				} else {
					addPotentialError(isExplicitlyPrioritized, idx, "coin transfer error")
				}
				continue
			}

			// Handle Merkle challenges on cached context
			detErrMsg, challengeNumIncrements, err := k.HandleMerkleChallenges(
				cachedCtx,
				collection.CollectionId,
				transfer,
				approval,
				transferMetadata,
				false,
			)
			if err != nil {
				if detErrMsg != "" {
					addPotentialError(isExplicitlyPrioritized, idx, detErrMsg)
				} else {
					addPotentialError(isExplicitlyPrioritized, idx, "merkle challenge error")
				}
				continue
			}

			// Handle ETH signature challenges on cached context
			detErrMsg, err = k.HandleETHSignatureChallenges(
				cachedCtx,
				collection.CollectionId,
				transfer,
				approval,
				transferMetadata,
			)
			if err != nil {
				if detErrMsg != "" {
					addPotentialError(isExplicitlyPrioritized, idx, detErrMsg)
				} else {
					addPotentialError(isExplicitlyPrioritized, idx, "eth signature challenge error")
				}
				continue
			}

			trackerTypes := []string{"overall", "to", "from", "initiatedBy"}
			// Initialize default values if ApprovalAmounts or MaxNumTransfers are nil
			var approvedAmounts []sdkmath.Uint
			if approvalCriteria.ApprovalAmounts != nil {
				approvedAmounts = []sdkmath.Uint{
					approvalCriteria.ApprovalAmounts.OverallApprovalAmount,
					approvalCriteria.ApprovalAmounts.PerToAddressApprovalAmount,
					approvalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount,
					approvalCriteria.ApprovalAmounts.PerInitiatedByAddressApprovalAmount,
				}
			} else {
				approvedAmounts = []sdkmath.Uint{
					sdkmath.NewUint(0),
					sdkmath.NewUint(0),
					sdkmath.NewUint(0),
					sdkmath.NewUint(0),
				}
			}
			var maxNumTransfers []sdkmath.Uint
			if approvalCriteria.MaxNumTransfers != nil {
				maxNumTransfers = []sdkmath.Uint{
					approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers,
					approvalCriteria.MaxNumTransfers.PerToAddressMaxNumTransfers,
					approvalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers,
					approvalCriteria.MaxNumTransfers.PerInitiatedByAddressMaxNumTransfers,
				}
			} else {
				maxNumTransfers = []sdkmath.Uint{
					sdkmath.NewUint(0),
					sdkmath.NewUint(0),
					sdkmath.NewUint(0),
					sdkmath.NewUint(0),
				}
			}
			approvedAddresses := []string{"", toAddress, fromAddress, initiatedBy}

			//Get max balances allowed for this approvalCriteria element
			// Use original context for GetMaxPossible since it reads current approval tracker state
			// (which may have been modified by previous transfers in the same transaction)
			failed := false
			for i, trackerType := range trackerTypes {
				//Get max allowed by criteria
				detErrMsg, maxPossible, err := k.GetMaxPossible(
					cachedCtx,
					collection,
					approval,
					transfer,
					originalTransferBalances,
					approvedAmounts[i],
					challengeNumIncrements,
					transferMetadata,
					trackerType,
					approvedAddresses[i],
				)
				if err != nil {
					if detErrMsg != "" {
						addPotentialError(isExplicitlyPrioritized, idx, detErrMsg)
					} else {
						addPotentialError(isExplicitlyPrioritized, idx, "get max possible error")
					}
					failed = true
					break
				}

				//Get max allowed by remaining balances to check
				transferBalancesToCheck, err = types.GetOverlappingBalances(cachedCtx, maxPossible, transferBalancesToCheck)
				if err != nil {
					addPotentialError(isExplicitlyPrioritized, idx, "get overlapping balances error")
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

			// Increment approvals and assert within threshold on cached context
			for i, trackerType := range trackerTypes {
				var precalculationOptions *types.PrecalculationOptions
				if transfer.PrecalculateBalancesFromApproval != nil {
					precalculationOptions = transfer.PrecalculateBalancesFromApproval.PrecalculationOptions
				}
				detErrMsg, err := k.IncrementApprovalsAndAssertWithinThreshold(cachedCtx, collection, approval, originalTransferBalances, approvedAmounts[i], maxNumTransfers[i], transferBalancesToCheck, challengeNumIncrements, transferMetadata, trackerType, approvedAddresses[i], precalculationOptions)
				if err != nil {
					if detErrMsg != "" {
						addPotentialError(isExplicitlyPrioritized, idx, detErrMsg)
					} else {
						addPotentialError(isExplicitlyPrioritized, idx, "increment approval and assert within threshold error")
					}
					failed = true
					break
				}
			}
			if failed {
				continue
			}

			// Validate that we can subtract balances before committing state changes
			// This is just a local calculation, no state change
			remainingBalances, err = types.SubtractBalances(cachedCtx, transferBalancesToCheck, remainingBalances)
			if err != nil {
				addPotentialError(isExplicitlyPrioritized, idx, "error subtracting balances")
				continue
			}

			// If we get here, all operations succeeded on the cached context
			// Write the cache back to the main context atomically
			writeCache()

			if !approvalCriteria.OverridesFromOutgoingApprovals {
				addToUserApprovalsToCheck(fromAddress, transferBalancesToCheck, true, userRoyalties)
			}

			if !approvalCriteria.OverridesToIncomingApprovals {
				addToUserApprovalsToCheck(toAddress, transferBalancesToCheck, false, userRoyalties)
			}

			markApprovalAsUsed(approval)
		}
	}

	//If we didn't find a successful approval, we throw
	if len(remainingBalances) > 0 {
		// If we used approvals and had partial success for some balances, we need to add an error for that
		errorsWithIdx = addPartialSuccessErrors(errorsWithIdx, *eventTracking.ApprovalsUsed, approvals)
		transferStr := buildTransferString(remainingBalances, fromAddress, toAddress, initiatedBy)
		potentialErrorsStr := buildPotentialErrorsString(potentialErrors, approvalIdxsChecked, errorsWithIdx)
		return []*UserApprovalsToCheck{}, buildApprovalFailureError(ctx, approvalLevel, transferStr, potentialErrorsStr)
	}

	// cannot have two different user royalty percentages
	if len(userApprovalsToCheck) > 0 {
		userRoyalties := userApprovalsToCheck[0].UserRoyalties
		for _, userApproval := range userApprovalsToCheck {
			if userApproval.UserRoyalties == nil || userApproval.UserRoyalties.Percentage.IsNil() {
				continue
			}

			if !userApproval.UserRoyalties.Percentage.Equal(userRoyalties.Percentage) {
				return []*UserApprovalsToCheck{}, buildMultipleRoyaltiesError(ctx)
			}
		}
	}

	return userApprovalsToCheck, nil
}

func isCustomChallengeOrderCalculation(predeterminedBalances *types.PredeterminedBalances, trackerType string) bool {
	return (predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers && trackerType == "overall") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from") ||
		(predeterminedBalances != nil && predeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy")
}

func (k Keeper) ResetApprovalTrackerIfNeeded(ctx sdk.Context, approvalTracker *types.ApprovalTracker, resetTimeIntervals *types.ResetTimeIntervals, isNumTransfers bool) types.ApprovalTracker {
	now := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
	if resetTimeIntervals != nil && resetTimeIntervals.StartTime.GT(sdkmath.NewUint(0)) && resetTimeIntervals.IntervalLength.GT(sdkmath.NewUint(0)) {
		startTime := resetTimeIntervals.StartTime
		intervalLength := resetTimeIntervals.IntervalLength
		lastResetAt := approvalTracker.LastUpdatedAt

		//If the first reset time is in the future, we don't need to reset
		if startTime.GT(now) {
			return *approvalTracker
		}

		//1. Calculate what interval we are in
		currInterval := now.Sub(startTime).Quo(intervalLength)
		currIntervalStart := startTime.Add(currInterval.Mul(intervalLength))

		if currIntervalStart.GT(lastResetAt) {
			if !isNumTransfers {
				approvalTracker.Amounts = []*types.Balance{}
			} else {
				approvalTracker.NumTransfers = sdkmath.NewUint(0)
			}
		}
	}

	// We can set the last updated no matter what
	// If it is N/A, it doesn't matter
	// If we need to update it, we update it
	approvalTracker.LastUpdatedAt = now

	return *approvalTracker
}

func (k Keeper) GetApprovalTrackerFromStoreAndResetIfNeeded(ctx sdk.Context, collectionId sdkmath.Uint, addressForApproval string, approvalId string, amountTrackerId string, level string, trackerType string, address string, resetTimeIntervals *types.ResetTimeIntervals, isNumTransfers bool) (types.ApprovalTracker, error) {
	approvalTracker, found := k.GetApprovalTrackerFromStore(ctx, collectionId, addressForApproval, approvalId, amountTrackerId, level, trackerType, address)
	if !found {
		return types.ApprovalTracker{
			Amounts:       []*types.Balance{},
			NumTransfers:  sdkmath.NewUint(0),
			LastUpdatedAt: sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli())),
		}, nil
	}

	return k.ResetApprovalTrackerIfNeeded(ctx, &approvalTracker, resetTimeIntervals, isNumTransfers), nil
}

// GetMaxPossible calculates the maximum possible transfer amounts based on approval criteria
// Returns (deterministicErrorMsg, balances, error) where deterministicErrorMsg is a deterministic error string
func (k Keeper) GetMaxPossible(
	ctx sdk.Context,
	collection *types.TokenCollection,
	approval *types.CollectionApproval,
	transfer *types.Transfer,
	originalTransferBalances []*types.Balance,
	approvedAmount sdkmath.Uint,
	challengeNumIncrements sdkmath.Uint,
	transferMetadata TransferMetadata,
	trackerType string,
	address string,
) (string, []*types.Balance, error) {
	approverAddress := transferMetadata.ApproverAddress
	approvalLevel := transferMetadata.ApprovalLevel
	// Initialize with transfer balances
	transferBalances := types.DeepCopyBalances(transfer.Balances)

	// If no amount restrictions, return full transfer balances
	if approvedAmount.IsNil() || approvedAmount.IsZero() {
		return "", transferBalances, nil
	}

	// Get approval tracker details
	// Check if ApprovalAmounts is nil before accessing
	if approval.ApprovalCriteria.ApprovalAmounts == nil {
		panic("ApprovalAmounts is nil")
	}
	amountsTrackerId := approval.ApprovalCriteria.ApprovalAmounts.AmountTrackerId

	// Fetch current approval tracker state
	amountsTrackerDetails, err := k.GetApprovalTrackerFromStoreAndResetIfNeeded(
		ctx,
		collection.CollectionId,
		approverAddress,
		approval.ApprovalId,
		amountsTrackerId,
		approvalLevel,
		trackerType,
		address,
		approval.ApprovalCriteria.ApprovalAmounts.ResetTimeIntervals,
		false,
	)
	if err != nil {
		return "", nil, err
	}

	// Calculate current tally for specific IDs and times
	currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(
		ctx,
		approval.TokenIds,
		approval.OwnershipTimes,
		amountsTrackerDetails.Amounts,
	)
	if err != nil {
		return "", nil, err
	}

	// Calculate maximum balances that can be added without exceeding threshold
	allApprovals := []*types.Balance{{
		Amount:         approvedAmount,
		OwnershipTimes: approval.OwnershipTimes,
		TokenIds:       approval.TokenIds,
	}}

	maxBalancesWeCanAdd, err := types.SubtractBalances(ctx, currTallyForCurrentIdsAndTimes, allApprovals)
	if err != nil {
		return "", nil, err
	}

	return "", maxBalancesWeCanAdd, nil
}

// handlePredeterminedBalances checks if the transfer matches predetermined balance requirements
func handlePredeterminedBalances(
	ctx sdk.Context,
	predeterminedBalances *types.PredeterminedBalances,
	originalTransferBalances []*types.Balance,
	trackerType string,
	trackerNumTransfers sdkmath.Uint,
	challengeNumIncrements sdkmath.Uint,
	precalculationOptions *types.PrecalculationOptions,
	collection *types.TokenCollection,
) ([]*types.Balance, error) {
	if predeterminedBalances == nil {
		return nil, nil
	}

	numIncrements := sdkmath.NewUint(0)
	toBeCalculated := true
	orderCalculationMethod := predeterminedBalances.OrderCalculationMethod

	// Determine how to calculate the number of increments
	switch {
	case orderCalculationMethod.UseMerkleChallengeLeafIndex:
		numIncrements = challengeNumIncrements
	case orderCalculationMethod.UseOverallNumTransfers && trackerType == "overall":
		numIncrements = trackerNumTransfers
	case orderCalculationMethod.UsePerToAddressNumTransfers && trackerType == "to":
		numIncrements = trackerNumTransfers
	case orderCalculationMethod.UsePerFromAddressNumTransfers && trackerType == "from":
		numIncrements = trackerNumTransfers
	case orderCalculationMethod.UsePerInitiatedByAddressNumTransfers && trackerType == "initiatedBy":
		numIncrements = trackerNumTransfers
	default:
		toBeCalculated = false
	}

	if !toBeCalculated {
		return nil, nil
	}

	var calculatedBalances []*types.Balance
	var err error

	if predeterminedBalances.ManualBalances != nil {
		if numIncrements.LT(sdkmath.NewUint(uint64(len(predeterminedBalances.ManualBalances)))) {
			calculatedBalances = types.DeepCopyBalances(predeterminedBalances.ManualBalances[numIncrements.Uint64()].Balances)
		}
	} else if predeterminedBalances.IncrementedBalances != nil {
		i := predeterminedBalances.IncrementedBalances
		calculatedBalances, err = types.IncrementBalances(
			ctx,
			i,
			numIncrements,
			precalculationOptions,
			collection,
		)
		if err != nil {
			return nil, err
		}
	}

	// Assert that we have exactly the amount specified in the original transfers
	if !types.AreBalancesEqual(ctx, originalTransferBalances, calculatedBalances, false) {
		return nil, sdkerrors.Wrapf(ErrDisallowedTransfer, "transfer disallowed because predetermined balances do not match")
	}

	return calculatedBalances, nil
}

// IncrementApprovalsAndAssertWithinThreshold handles approval tracking and threshold checks
// Returns (deterministicErrorMsg, error) where deterministicErrorMsg is a deterministic error string
func (k Keeper) IncrementApprovalsAndAssertWithinThreshold(
	ctx sdk.Context,
	collection *types.TokenCollection,
	approval *types.CollectionApproval,
	originalTransferBalances []*types.Balance,
	approvedAmount sdkmath.Uint,
	maxNumTransfers sdkmath.Uint,
	transferBalances []*types.Balance,
	challengeNumIncrements sdkmath.Uint,
	transferMetadata TransferMetadata,
	trackerType string,
	address string,
	precalculationOptions *types.PrecalculationOptions,
) (string, error) {
	approverAddress := transferMetadata.ApproverAddress
	approvalLevel := transferMetadata.ApprovalLevel
	approvalCriteria := approval.ApprovalCriteria
	// Initialize default values if ApprovalAmounts or MaxNumTransfers are nil
	var amountsTrackerId string
	if approvalCriteria.ApprovalAmounts != nil {
		amountsTrackerId = approvalCriteria.ApprovalAmounts.AmountTrackerId
	}
	var maxNumTransfersTrackerId string
	if approvalCriteria.MaxNumTransfers != nil {
		maxNumTransfersTrackerId = approvalCriteria.MaxNumTransfers.AmountTrackerId
	}

	// Initialize default values
	if approvedAmount.IsNil() {
		approvedAmount = sdkmath.NewUint(0)
	}
	if maxNumTransfers.IsNil() {
		maxNumTransfers = sdkmath.NewUint(0)
	}

	// Fetch approval tracker details
	amountsTrackerDetails := types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli())),
	}

	maxNumTransfersTrackerDetails := types.ApprovalTracker{
		Amounts:       []*types.Balance{},
		NumTransfers:  sdkmath.NewUint(0),
		LastUpdatedAt: sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli())),
	}

	needToFetchApprovalTrackerDetails := maxNumTransfers.GT(sdkmath.NewUint(0)) ||
		approvedAmount.GT(sdkmath.NewUint(0)) ||
		isCustomChallengeOrderCalculation(approvalCriteria.PredeterminedBalances, trackerType)

	var err error
	if needToFetchApprovalTrackerDetails {
		var resetTimeIntervals *types.ResetTimeIntervals
		if approvalCriteria.ApprovalAmounts != nil {
			resetTimeIntervals = approvalCriteria.ApprovalAmounts.ResetTimeIntervals
		}
		amountsTrackerDetails, err = k.GetApprovalTrackerFromStoreAndResetIfNeeded(
			ctx,
			collection.CollectionId,
			approverAddress,
			approval.ApprovalId,
			amountsTrackerId,
			approvalLevel,
			trackerType,
			address,
			resetTimeIntervals,
			false,
		)
		if err != nil {
			return "", err
		}

		var maxNumTransfersResetTimeIntervals *types.ResetTimeIntervals
		if approvalCriteria.MaxNumTransfers != nil {
			maxNumTransfersResetTimeIntervals = approvalCriteria.MaxNumTransfers.ResetTimeIntervals
		}
		maxNumTransfersTrackerDetails, err = k.GetApprovalTrackerFromStoreAndResetIfNeeded(
			ctx,
			collection.CollectionId,
			approverAddress,
			approval.ApprovalId,
			maxNumTransfersTrackerId,
			approvalLevel,
			trackerType,
			address,
			maxNumTransfersResetTimeIntervals,
			true,
		)
		if err != nil {
			return "", err
		}
	}

	// Handle predetermined balances check
	_, err = handlePredeterminedBalances(
		ctx,
		approvalCriteria.PredeterminedBalances,
		originalTransferBalances,
		trackerType,
		maxNumTransfersTrackerDetails.NumTransfers,
		challengeNumIncrements,
		precalculationOptions,
		collection,
	)
	if err != nil {
		return "", err
	}

	// Handle amount approvals
	if approvedAmount.GT(sdkmath.NewUint(0)) {
		currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(
			ctx,
			approval.TokenIds,
			approval.OwnershipTimes,
			amountsTrackerDetails.Amounts,
		)
		if err != nil {
			return "", err
		}

		thresholdAmounts := []*types.Balance{{
			Amount:         approvedAmount,
			OwnershipTimes: approval.OwnershipTimes,
			TokenIds:       approval.TokenIds,
		}}

		_, err = types.AddBalancesAndAssertDoesntExceedThreshold(
			ctx,
			currTallyForCurrentIdsAndTimes,
			transferBalances,
			thresholdAmounts,
		)
		if err != nil {
			return "", err
		}

		amountsTrackerDetails.Amounts, err = types.AddBalances(
			ctx,
			amountsTrackerDetails.Amounts,
			transferBalances,
		)
		if err != nil {
			return "", err
		}
	}

	// Handle max transfers tracking
	if maxNumTransfers.GT(sdkmath.NewUint(0)) || isCustomChallengeOrderCalculation(approvalCriteria.PredeterminedBalances, trackerType) {
		maxNumTransfersTrackerDetails.NumTransfers = maxNumTransfersTrackerDetails.NumTransfers.Add(sdkmath.NewUint(1))
		if maxNumTransfers.GT(sdkmath.NewUint(0)) && maxNumTransfersTrackerDetails.NumTransfers.GT(maxNumTransfers) {
			detErrMsg := fmt.Sprintf("exceeded max transfers allowed - %s", maxNumTransfers.String())
			return detErrMsg, sdkerrors.Wrap(ErrDisallowedTransfer, detErrMsg)
		}
	}

	// Handle event emission and store updates
	if needToFetchApprovalTrackerDetails {
		marshalToString := func(v interface{}) (string, error) {
			data, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(data), nil
		}

		amountsStr, err := marshalToString(amountsTrackerDetails.Amounts)
		if err != nil {
			return "", err
		}

		numTransfersStr, err := marshalToString(maxNumTransfersTrackerDetails.NumTransfers)
		if err != nil {
			return "", err
		}

		amountsNumTransfersStr, err := marshalToString(amountsTrackerDetails.NumTransfers)
		if err != nil {
			return "", err
		}

		maxNumTransfersAmountsStr, err := marshalToString(maxNumTransfersTrackerDetails.Amounts)
		if err != nil {
			return "", err
		}

		isSameId := amountsTrackerId == maxNumTransfersTrackerId
		if isSameId {
			EmitApprovalEvent(
				ctx,
				collection.CollectionId,
				approverAddress,
				approval.ApprovalId,
				amountsTrackerId,
				approvalLevel,
				trackerType,
				address,
				amountsStr,
				numTransfersStr,
				maxNumTransfersTrackerDetails.LastUpdatedAt,
			)

			amountsTrackerDetails.NumTransfers = maxNumTransfersTrackerDetails.NumTransfers
			err = k.SetApprovalTrackerInStore(
				ctx,
				collection.CollectionId,
				approverAddress,
				approval.ApprovalId,
				amountsTrackerId,
				amountsTrackerDetails,
				approvalLevel,
				trackerType,
				address,
			)
			if err != nil {
				return "", err
			}
		} else {
			EmitApprovalEvent(
				ctx,
				collection.CollectionId,
				approverAddress,
				approval.ApprovalId,
				amountsTrackerId,
				approvalLevel,
				trackerType,
				address,
				amountsStr,
				amountsNumTransfersStr,
				maxNumTransfersTrackerDetails.LastUpdatedAt,
			)

			EmitApprovalEvent(
				ctx,
				collection.CollectionId,
				approverAddress,
				approval.ApprovalId,
				maxNumTransfersTrackerId,
				approvalLevel,
				trackerType,
				address,
				maxNumTransfersAmountsStr,
				numTransfersStr,
				maxNumTransfersTrackerDetails.LastUpdatedAt,
			)

			err = k.SetApprovalTrackerInStore(
				ctx,
				collection.CollectionId,
				approverAddress,
				approval.ApprovalId,
				amountsTrackerId,
				amountsTrackerDetails,
				approvalLevel,
				trackerType,
				address,
			)
			if err != nil {
				return "", err
			}

			err = k.SetApprovalTrackerInStore(
				ctx,
				collection.CollectionId,
				approverAddress,
				approval.ApprovalId,
				maxNumTransfersTrackerId,
				maxNumTransfersTrackerDetails,
				approvalLevel,
				trackerType,
				address,
			)
			if err != nil {
				return "", err
			}
		}
	}

	return "", nil
}

func (k Keeper) GetPredeterminedBalancesForPrecalculationId(
	ctx sdk.Context,
	collection *types.TokenCollection,
	approvals []*types.CollectionApproval,
	transfer *types.Transfer,
	transferMetadata TransferMetadata,
) ([]*types.Balance, error) {
	to := transferMetadata.To
	initiatedBy := transferMetadata.InitiatedBy
	approvalId := ""
	precalcDetails := transfer.PrecalculateBalancesFromApproval
	var precalculationOptions *types.PrecalculationOptions
	var approverAddress string
	var approvalLevel string
	var precalculationId string
	if precalcDetails != nil {
		precalculationOptions = precalcDetails.PrecalculationOptions
		approverAddress = precalcDetails.ApproverAddress
		approvalLevel = precalcDetails.ApprovalLevel
		precalculationId = precalcDetails.ApprovalId
	}

	for _, approval := range approvals {
		approvalCriteria := approval.ApprovalCriteria
		approvalId = approval.ApprovalId
		if approvalCriteria == nil || approvalId != precalculationId || approvalId == "" {
			continue
		}

		// Check if MaxNumTransfers is nil before accessing
		if approvalCriteria.MaxNumTransfers == nil {
			panic("MaxNumTransfers is nil")
		}
		maxNumTransfersTrackerId := approvalCriteria.MaxNumTransfers.AmountTrackerId

		if transfer.PrecalculateBalancesFromApproval != nil {
			if !approval.Version.Equal(transfer.PrecalculateBalancesFromApproval.Version) {
				return []*types.Balance{}, sdkerrors.Wrapf(types.ErrMismatchedVersions, "versions are mismatched for a prioritized approval")
			}
		}

		if approvalCriteria.PredeterminedBalances != nil {
			numIncrements := sdkmath.NewUint(0)
			hasOrderCalculationMethod := false
			if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseMerkleChallengeLeafIndex {
				hasOrderCalculationMethod = true

				//If the approval has challenges, we need to check that a valid solutions is provided for every challenge
				//If the challenge specifies to use the leaf index for the number of increments, we use this value for the number of increments later
				_, numIncrementsFetched, err := k.HandleMerkleChallenges(
					ctx,
					collection.CollectionId,
					transfer,
					approval,
					transferMetadata,
					true,
				)
				if err != nil {
					return []*types.Balance{}, sdkerrors.Wrapf(err, "invalid challenges / solutions")
				}

				numIncrements = numIncrementsFetched
			} else {
				trackerType := ""
				approvedAddress := ""

				if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UseOverallNumTransfers {
					trackerType = "overall"
				} else if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerFromAddressNumTransfers {
					trackerType = "from"
					approvedAddress = transfer.From
				} else if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerToAddressNumTransfers {
					trackerType = "to"
					approvedAddress = to
				} else if approvalCriteria.PredeterminedBalances.OrderCalculationMethod.UsePerInitiatedByAddressNumTransfers {
					trackerType = "initiatedBy"
					approvedAddress = initiatedBy
				}

				if trackerType != "" {
					hasOrderCalculationMethod = true
				}

				numTransfersTracker, err := k.GetApprovalTrackerFromStoreAndResetIfNeeded(
					ctx,
					collection.CollectionId,
					approverAddress,
					approval.ApprovalId,
					maxNumTransfersTrackerId,
					approvalLevel,
					trackerType,
					approvedAddress,
					approval.ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals,
					true,
				)
				if err != nil {
					return nil, err
				}

				numIncrements = numTransfersTracker.NumTransfers
			}

			if !hasOrderCalculationMethod {
				return []*types.Balance{}, sdkerrors.Wrapf(ErrDisallowedTransfer, "no order calculation method found for approval id: %s", precalculationId)
			}

			//calculate the current approved balances from the numIncrements and predeterminedBalances
			predeterminedBalances := []*types.Balance{}
			if approvalCriteria.PredeterminedBalances.ManualBalances != nil {
				if numIncrements.LT(sdkmath.NewUint(uint64(len(approvalCriteria.PredeterminedBalances.ManualBalances)))) {
					predeterminedBalances = types.DeepCopyBalances(approvalCriteria.PredeterminedBalances.ManualBalances[numIncrements.Uint64()].Balances)
				}
			} else if approvalCriteria.PredeterminedBalances.IncrementedBalances != nil {
				var err error
				predeterminedBalances, err = types.IncrementBalances(
					ctx,
					approvalCriteria.PredeterminedBalances.IncrementedBalances,
					numIncrements,
					precalculationOptions,
					collection,
				)
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
