package keeper

import (
	"encoding/json"
	"fmt"
	"strings"

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
	collection *types.BadgeCollection,
	originalTransferBalances []*types.Balance,
	transfer *types.Transfer,
	_approvals []*types.CollectionApproval,
	toAddress string,
	initiatedBy string,
	approvalLevel string,
	approverAddress string,
	approvalsUsed *[]ApprovalsUsed,
	coinTransfersUsed *[]CoinTransfers,
	royalties *types.UserRoyalties,
) ([]*UserApprovalsToCheck, error) {
	fromAddress := transfer.From
	originalTransferBalances = types.DeepCopyBalances(originalTransferBalances)
	remainingBalances := types.DeepCopyBalances(transfer.Balances) //Keep a running tally of all the badges we still have to handle
	approvals, err := SortViaPrioritizedApprovals(_approvals, transfer, approvalLevel, approverAddress)
	if err != nil {
		return []*UserApprovalsToCheck{}, err
	}

	potentialErrors := []string{}

	//For each approved transfer, we check if the transfer is allowed
	// 1: If transfer meets all criteria, we deduct, get user approvals to check, and continue (if there are any remaining balances)
	// 2. If transfer does not meet all criteria, we continue and do not mark anything as handled
	// 3. At the end, if there are any unhandled transfers, we throw (not enough approvals = transfer disallowed)
	userApprovalsToCheck := []*UserApprovalsToCheck{}
	for _, approval := range approvals {
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

		isPrioritizedApproval := false
		for _, prioritizedApproval := range transfer.PrioritizedApprovals {
			if prioritizedApproval.ApprovalId == approval.ApprovalId && prioritizedApproval.ApprovalLevel == approvalLevel && prioritizedApproval.ApproverAddress == approverAddress {
				isPrioritizedApproval = true
				break
			}
		}

		//Initial checks: Make sure (from, to, initiatedBy) match the approval's collection list IDs
		doAddressesMatch := k.CheckIfAddressesMatchCollectionListIds(ctx, approval, fromAddress, toAddress, initiatedBy)
		if !doAddressesMatch {
			if isPrioritizedApproval {
				potentialErrors = append(potentialErrors, fmt.Sprintf("addresses do not match (from, to, initiatedBy): %s, %s, %s", fromAddress, toAddress, initiatedBy))
			}
			continue
		}

		// Check valid transfer times
		currTime := sdkmath.NewUint(uint64(ctx.BlockTime().UnixMilli()))
		currTimeFound, err := types.SearchUintRangesForUint(currTime, approval.TransferTimes)
		if !currTimeFound || err != nil {
			if isPrioritizedApproval {
				potentialErrors = append(potentialErrors, fmt.Sprintf("transfer time not in range"))
			}
			continue
		}

		transferStr := "attempting to transfer badge ID " + approval.BadgeIds[0].Start.String()

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

			//If we do not override the approved outgoing / incoming transfers, we need to check the user approvals
			userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
				Address:       fromAddress,
				Balances:      allBalancesForIdsAndTimes,
				Outgoing:      true,
				UserRoyalties: userRoyalties,
			})

			userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
				Address:       toAddress,
				Balances:      allBalancesForIdsAndTimes,
				Outgoing:      false,
				UserRoyalties: userRoyalties,
			})

			*approvalsUsed = append(*approvalsUsed, ApprovalsUsed{
				ApprovalId:      approval.ApprovalId,
				ApprovalLevel:   approvalLevel,
				ApproverAddress: approverAddress,
				Version:         approval.Version.String(),
			})
		} else {
			//Else, we have a match and we can proceed to check the restrictions
			//This is split into a two part process:
			//	1. Simulate the transfer to see if it is allowed
			//	2. If all simulations pass, we deduct as much as possible from the approval

			approvalCriteria := approval.ApprovalCriteria
			if approvalCriteria.RequireFromDoesNotEqualInitiatedBy && fromAddress == initiatedBy {
				if isPrioritizedApproval {
					potentialErrors = append(potentialErrors, "from address equals initiated by")
				}
				continue
			}

			if approvalCriteria.RequireFromEqualsInitiatedBy && fromAddress != initiatedBy {
				if isPrioritizedApproval {
					potentialErrors = append(potentialErrors, "from address does not equal initiated by")
				}
				continue
			}

			if approvalCriteria.RequireToDoesNotEqualInitiatedBy && toAddress == initiatedBy {
				if isPrioritizedApproval {
					potentialErrors = append(potentialErrors, "to address equals initiated by")
				}
				continue
			}

			if approvalCriteria.RequireToEqualsInitiatedBy && toAddress != initiatedBy {
				if isPrioritizedApproval {
					potentialErrors = append(potentialErrors, "to address does not equal initiated by")
				}
				continue
			}

			// Must own badges check
			err = k.CheckMustOwnBadges(ctx, approvalCriteria.MustOwnBadges, initiatedBy)
			if err != nil {
				continue
			}

			// Dynamic store challenges check - all challenges must pass
			for _, challenge := range approvalCriteria.DynamicStoreChallenges {
				storeId := challenge.StoreId
				dynamicStoreValue, found := k.GetDynamicStoreValueFromStore(ctx, storeId, initiatedBy)
				val := false
				if found {
					val = dynamicStoreValue.Value
				} else {
					dynamicStore, foundStore := k.GetDynamicStoreFromStore(ctx, storeId)
					if !foundStore {
						if isPrioritizedApproval {
							potentialErrors = append(potentialErrors, fmt.Sprintf("dynamic store not found for storeId %s", storeId.String()))
						}
						goto skipApproval
					}
					val = dynamicStore.DefaultValue
				}
				if !val {
					if isPrioritizedApproval {
						potentialErrors = append(potentialErrors, fmt.Sprintf("initiator did not pass dynamic store challenge for storeId %s", storeId.String()))
					}
					goto skipApproval
				}
			}

		skipApproval:

			/**** SECTION 1: NO STORAGE WRITES (just simulate everything and continue if it doesn't pass) ****/
			err := k.HandleCoinTransfers(ctx, approvalCriteria.CoinTransfers, initiatedBy, approverAddress, approvalLevel, true, coinTransfersUsed, collection, royalties) //simulate = true
			if err != nil {
				if isPrioritizedApproval {
					potentialErrors = append(potentialErrors, fmt.Sprintf("coin transfer error: %s", err))
				}
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
				if isPrioritizedApproval {
					potentialErrors = append(potentialErrors, "no balances to check")
				}
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
				if isPrioritizedApproval {
					potentialErrors = append(potentialErrors, fmt.Sprintf("merkle challenge error: %s", err))
				}
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
					if isPrioritizedApproval {
						potentialErrors = append(potentialErrors, fmt.Sprintf("get max possible error: %s", err))
					}
					failed = true
					break
				}

				//Get max allowed by remaining balances to check
				transferBalancesToCheck, err = types.GetOverlappingBalances(ctx, maxPossible, transferBalancesToCheck)
				if err != nil {
					if isPrioritizedApproval {
						potentialErrors = append(potentialErrors, fmt.Sprintf("get overlapping balances error: %s", err))
					}
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
				err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, approval, originalTransferBalances, approvedAmounts[i], maxNumTransfers[i], transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, trackerType, approvedAddresses[i], true, transfer.PrecalculationOptions)
				if err != nil {
					failed = true
					break
				}
			}
			if failed {
				if isPrioritizedApproval {
					potentialErrors = append(potentialErrors, fmt.Sprintf("increment approval and assert within threshold error: %s", err))
				}
				continue
			}

			/**** SECTION 2: ONCE HERE, EVERYTHING BELOW SHOULD BE SUCCESSFUL BC IT WAS SIMULATED ****/
			remainingBalances, err = types.SubtractBalances(ctx, transferBalancesToCheck, remainingBalances)
			if err != nil {
				continue
			}

			err = k.HandleCoinTransfers(ctx, approvalCriteria.CoinTransfers, initiatedBy, approverAddress, approvalLevel, false, coinTransfersUsed, collection, royalties) //simulate = false
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

			for i, trackerType := range trackerTypes {
				err = k.IncrementApprovalsAndAssertWithinThreshold(ctx, collection, approval, originalTransferBalances, approvedAmounts[i], maxNumTransfers[i], transferBalancesToCheck, challengeNumIncrements, approverAddress, approvalLevel, trackerType, approvedAddresses[i], false, transfer.PrecalculationOptions)
				if err != nil {
					return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(err, "error incrementing approvals")
				}
			}

			//If we do not override the approved outgoing / incoming transfers, we need to check the user approvals
			if !approvalCriteria.OverridesFromOutgoingApprovals {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address:       fromAddress,
					Balances:      transferBalancesToCheck,
					Outgoing:      true,
					UserRoyalties: userRoyalties,
				})
			}

			if !approvalCriteria.OverridesToIncomingApprovals {
				userApprovalsToCheck = append(userApprovalsToCheck, &UserApprovalsToCheck{
					Address:       toAddress,
					Balances:      transferBalancesToCheck,
					Outgoing:      false,
					UserRoyalties: userRoyalties,
				})
			}

			*approvalsUsed = append(*approvalsUsed, ApprovalsUsed{
				ApprovalId:      approval.ApprovalId,
				ApprovalLevel:   approvalLevel,
				ApproverAddress: approverAddress,
				Version:         approval.Version.String(),
			})

		}
	}

	//If we didn't find a successful approval, we throw
	if len(remainingBalances) > 0 {
		transferStr := "attempting to transfer badge ID " + remainingBalances[0].BadgeIds[0].Start.String()
		potentialErrorsStr := ""
		if len(potentialErrors) > 0 {
			potentialErrorsStr = " - errors w/ prioritized approvals: " + strings.Join(potentialErrors, ", ")
		} else {
			potentialErrorsStr = " - auto-scan failed (no prioritized approvals found): "
		}
		return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrInadequateApprovals, "no approval satisfied for transfer: %s%s", transferStr, potentialErrorsStr)
	}

	// cannot have two different user royalty percentages
	if len(userApprovalsToCheck) > 0 {
		userRoyalties := userApprovalsToCheck[0].UserRoyalties
		for _, userApproval := range userApprovalsToCheck {
			if userApproval.UserRoyalties == nil || userApproval.UserRoyalties.Percentage.IsNil() {
				continue
			}

			if !userApproval.UserRoyalties.Percentage.Equal(userRoyalties.Percentage) {
				return []*UserApprovalsToCheck{}, sdkerrors.Wrapf(ErrInadequateApprovals, "multiple user-level royalties found - please split your transfer up to use one collection approval w/ royalty per transfer")
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
	// Initialize with transfer balances
	transferBalances := types.DeepCopyBalances(transfer.Balances)

	// If no amount restrictions, return full transfer balances
	if approvedAmount.IsNil() || approvedAmount.IsZero() {
		return transferBalances, nil
	}

	// Get approval tracker details
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
		return nil, err
	}

	// Calculate current tally for specific badge IDs and times
	currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(
		ctx,
		approval.BadgeIds,
		approval.OwnershipTimes,
		amountsTrackerDetails.Amounts,
	)
	if err != nil {
		return nil, err
	}

	// Calculate maximum balances that can be added without exceeding threshold
	allApprovals := []*types.Balance{{
		Amount:         approvedAmount,
		OwnershipTimes: approval.OwnershipTimes,
		BadgeIds:       approval.BadgeIds,
	}}

	maxBalancesWeCanAdd, err := types.SubtractBalances(ctx, currTallyForCurrentIdsAndTimes, allApprovals)
	if err != nil {
		return nil, err
	}

	return maxBalancesWeCanAdd, nil
}

// handlePredeterminedBalances checks if the transfer matches predetermined balance requirements
func (k Keeper) handlePredeterminedBalances(
	ctx sdk.Context,
	predeterminedBalances *types.PredeterminedBalances,
	originalTransferBalances []*types.Balance,
	trackerType string,
	trackerNumTransfers sdkmath.Uint,
	challengeNumIncrements sdkmath.Uint,
	precalculationOptions *types.PrecalculationOptions,
	collection *types.BadgeCollection,
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
			i.StartBalances,
			numIncrements,
			i.IncrementOwnershipTimesBy,
			i.IncrementBadgeIdsBy,
			i.DurationFromTimestamp,
			i.RecurringOwnershipTimes,
			precalculationOptions.OverrideTimestamp,
			i.AllowOverrideTimestamp,
			precalculationOptions.BadgeIdsOverride,
			i.AllowOverrideWithAnyValidBadge,
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

// emitApprovalEvent emits an event for approval tracking
func emitApprovalEvent(
	ctx sdk.Context,
	collectionId sdkmath.Uint,
	approverAddress string,
	approvalId string,
	amountsTrackerId string,
	approvalLevel string,
	trackerType string,
	address string,
	amountsStr string,
	numTransfersStr string,
	lastUpdatedAt sdkmath.Uint,
) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"approval"+fmt.Sprint(collectionId)+fmt.Sprint(approverAddress)+fmt.Sprint(approvalId)+fmt.Sprint(amountsTrackerId)+fmt.Sprint(approvalLevel)+fmt.Sprint(trackerType)+fmt.Sprint(address),
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
			sdk.NewAttribute("approvalId", fmt.Sprint(approvalId)),
			sdk.NewAttribute("approverAddress", fmt.Sprint(approverAddress)),
			sdk.NewAttribute("amountTrackerId", fmt.Sprint(amountsTrackerId)),
			sdk.NewAttribute("approvalLevel", fmt.Sprint(approvalLevel)),
			sdk.NewAttribute("trackerType", fmt.Sprint(trackerType)),
			sdk.NewAttribute("approvedAddress", fmt.Sprint(address)),
			sdk.NewAttribute("amounts", amountsStr),
			sdk.NewAttribute("numTransfers", numTransfersStr),
			sdk.NewAttribute("lastUpdatedAt", fmt.Sprint(lastUpdatedAt)),
		),
	)
}

// IncrementApprovalsAndAssertWithinThreshold handles approval tracking and threshold checks
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
	precalculationOptions *types.PrecalculationOptions,
) error {
	approvalCriteria := approval.ApprovalCriteria
	amountsTrackerId := approvalCriteria.ApprovalAmounts.AmountTrackerId
	maxNumTransfersTrackerId := approvalCriteria.MaxNumTransfers.AmountTrackerId

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

	err := *new(error)
	if needToFetchApprovalTrackerDetails {
		amountsTrackerDetails, err = k.GetApprovalTrackerFromStoreAndResetIfNeeded(
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
			return err
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
			approval.ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals,
			true,
		)
		if err != nil {
			return err
		}
	}

	// Handle predetermined balances check
	_, err = k.handlePredeterminedBalances(
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
		return err
	}

	// Handle amount approvals
	if approvedAmount.GT(sdkmath.NewUint(0)) {
		currTallyForCurrentIdsAndTimes, err := types.GetBalancesForIds(
			ctx,
			approval.BadgeIds,
			approval.OwnershipTimes,
			amountsTrackerDetails.Amounts,
		)
		if err != nil {
			return err
		}

		thresholdAmounts := []*types.Balance{{
			Amount:         approvedAmount,
			OwnershipTimes: approval.OwnershipTimes,
			BadgeIds:       approval.BadgeIds,
		}}

		_, err = types.AddBalancesAndAssertDoesntExceedThreshold(
			ctx,
			currTallyForCurrentIdsAndTimes,
			transferBalances,
			thresholdAmounts,
		)
		if err != nil {
			return err
		}

		amountsTrackerDetails.Amounts, err = types.AddBalances(
			ctx,
			amountsTrackerDetails.Amounts,
			transferBalances,
		)
		if err != nil {
			return err
		}
	}

	// Handle max transfers tracking
	if maxNumTransfers.GT(sdkmath.NewUint(0)) || isCustomChallengeOrderCalculation(approvalCriteria.PredeterminedBalances, trackerType) {
		maxNumTransfersTrackerDetails.NumTransfers = maxNumTransfersTrackerDetails.NumTransfers.Add(sdkmath.NewUint(1))
		if maxNumTransfers.GT(sdkmath.NewUint(0)) && maxNumTransfersTrackerDetails.NumTransfers.GT(maxNumTransfers) {
			return sdkerrors.Wrapf(ErrDisallowedTransfer, "exceeded max transfers allowed - %s", maxNumTransfers.String())
		}
	}

	// Handle event emission and store updates
	if needToFetchApprovalTrackerDetails && !simulate {
		marshalToString := func(v interface{}) (string, error) {
			data, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(data), nil
		}

		amountsStr, err := marshalToString(amountsTrackerDetails.Amounts)
		if err != nil {
			return err
		}

		numTransfersStr, err := marshalToString(maxNumTransfersTrackerDetails.NumTransfers)
		if err != nil {
			return err
		}

		amountsNumTransfersStr, err := marshalToString(amountsTrackerDetails.NumTransfers)
		if err != nil {
			return err
		}

		maxNumTransfersAmountsStr, err := marshalToString(maxNumTransfersTrackerDetails.Amounts)
		if err != nil {
			return err
		}

		isSameId := amountsTrackerId == maxNumTransfersTrackerId
		if isSameId {
			emitApprovalEvent(
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
				return err
			}
		} else {
			emitApprovalEvent(
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

			emitApprovalEvent(
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
				return err
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
	precalculationOptions *types.PrecalculationOptions,
) ([]*types.Balance, error) {
	approvalId := ""
	approverAddress := precalcDetails.ApproverAddress
	approvalLevel := precalcDetails.ApprovalLevel
	precalculationId := precalcDetails.ApprovalId

	for _, approval := range approvals {
		maxNumTransfersTrackerId := approval.ApprovalCriteria.MaxNumTransfers.AmountTrackerId
		approvalCriteria := approval.ApprovalCriteria
		approvalId = approval.ApprovalId
		if approvalCriteria == nil || approvalId != precalculationId || approvalId == "" {
			continue
		}

		if !approval.Version.Equal(transfer.PrecalculateBalancesFromApproval.Version) {
			return []*types.Balance{}, sdkerrors.Wrapf(types.ErrMismatchedVersions, "versions are mismatched for a prioritized approval")
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

			//calculate the current approved balances from the numIncrements and predeterminedBalances
			predeterminedBalances := []*types.Balance{}
			if approvalCriteria.PredeterminedBalances.ManualBalances != nil {
				if numIncrements.LT(sdkmath.NewUint(uint64(len(approvalCriteria.PredeterminedBalances.ManualBalances)))) {
					predeterminedBalances = types.DeepCopyBalances(approvalCriteria.PredeterminedBalances.ManualBalances[numIncrements.Uint64()].Balances)
				}
			} else if approvalCriteria.PredeterminedBalances.IncrementedBalances != nil {
				err := *new(error)
				predeterminedBalances, err = types.IncrementBalances(
					ctx,
					approvalCriteria.PredeterminedBalances.IncrementedBalances.StartBalances,
					numIncrements,
					approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementOwnershipTimesBy,
					approvalCriteria.PredeterminedBalances.IncrementedBalances.IncrementBadgeIdsBy,
					approvalCriteria.PredeterminedBalances.IncrementedBalances.DurationFromTimestamp,
					approvalCriteria.PredeterminedBalances.IncrementedBalances.RecurringOwnershipTimes,
					precalculationOptions.OverrideTimestamp,
					approvalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideTimestamp,
					precalculationOptions.BadgeIdsOverride,
					approvalCriteria.PredeterminedBalances.IncrementedBalances.AllowOverrideWithAnyValidBadge,
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
