package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

func GetDefaultBalanceStoreForCollection(collection *types.BadgeCollection) *types.UserBalanceStore {
	return &types.UserBalanceStore{
		Balances:          collection.DefaultBalances.Balances,
		OutgoingApprovals: collection.DefaultBalances.OutgoingApprovals,
		IncomingApprovals: collection.DefaultBalances.IncomingApprovals,
		AutoApproveSelfInitiatedOutgoingTransfers: collection.DefaultBalances.AutoApproveSelfInitiatedOutgoingTransfers,
		AutoApproveSelfInitiatedIncomingTransfers: collection.DefaultBalances.AutoApproveSelfInitiatedIncomingTransfers,
		AutoApproveAllIncomingTransfers:           collection.DefaultBalances.AutoApproveAllIncomingTransfers,
		UserPermissions:                           collection.DefaultBalances.UserPermissions,
	}
}

func (k Keeper) GetBalanceOrApplyDefault(ctx sdk.Context, collection *types.BadgeCollection, userAddress string) (*types.UserBalanceStore, bool) {
	//Mint has unlimited balances
	if userAddress == "Total" || userAddress == "Mint" {
		return &types.UserBalanceStore{}, false
	}

	//We get current balances or fallback to default balances
	balanceKey := ConstructBalanceKey(userAddress, collection.CollectionId)
	balance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
	appliedDefault := false
	if !found {
		balance = GetDefaultBalanceStoreForCollection(collection)
		appliedDefault = true
		// We need to set the version to "0" for all incoming and outgoing approvals
		for _, approval := range balance.IncomingApprovals {
			approval.Version = k.IncrementApprovalVersion(ctx, collection.CollectionId, "incoming", userAddress, approval.ApprovalId)
		}
		for _, approval := range balance.OutgoingApprovals {
			approval.Version = k.IncrementApprovalVersion(ctx, collection.CollectionId, "outgoing", userAddress, approval.ApprovalId)
		}
	}

	return balance, appliedDefault
}

func (k Keeper) SetBalanceForAddress(ctx sdk.Context, collection *types.BadgeCollection, userAddress string, balance *types.UserBalanceStore) error {
	balanceKey := ConstructBalanceKey(userAddress, collection.CollectionId)
	return k.SetUserBalanceInStore(ctx, balanceKey, balance)
}

type ApprovalsUsed struct {
	ApprovalId      string
	ApprovalLevel   string
	ApproverAddress string
	Version         string
}

type CoinTransfers struct {
	From          string
	To            string
	Amount        string
	Denom         string
	IsProtocolFee bool
}

func (k Keeper) HandleTransfers(ctx sdk.Context, collection *types.BadgeCollection, transfers []*types.Transfer, initiatedBy string) error {
	err := *new(error)

	// Validate transfers with invariants
	for _, transfer := range transfers {
		if err := types.ValidateTransferWithInvariants(ctx, transfer, true, collection); err != nil {
			return err
		}
	}

	for _, transfer := range transfers {
		numAttempts := sdkmath.NewUint(1)
		if !transfer.NumAttempts.IsNil() {
			numAttempts = transfer.NumAttempts
		}

		for i := sdkmath.NewUint(0); i.LT(numAttempts); i = i.Add(sdkmath.NewUint(1)) {
			fromUserBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, transfer.From)
			totalMinted := []*types.Balance{}

			for _, to := range transfer.ToAddresses {
				approvalsUsed := []ApprovalsUsed{}
				coinTransfers := []CoinTransfers{}

				toUserBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, to)

				if transfer.PrecalculateBalancesFromApproval != nil && transfer.PrecalculateBalancesFromApproval.ApprovalId != "" {
					//Here, we precalculate balances from a specified approval
					approvals := collection.CollectionApprovals
					if transfer.PrecalculateBalancesFromApproval.ApprovalLevel == "collection" {
						if transfer.PrecalculateBalancesFromApproval.ApproverAddress != "" {
							return sdkerrors.Wrapf(ErrNotImplemented, "approver address must be blank for collection level approvals")
						}
					} else {
						if transfer.PrecalculateBalancesFromApproval.ApproverAddress != to && transfer.PrecalculateBalancesFromApproval.ApproverAddress != transfer.From {
							return sdkerrors.Wrapf(ErrNotImplemented, "approver address %s must match to or from address for user level precalculations", transfer.PrecalculateBalancesFromApproval.ApproverAddress)
						}

						handled := false
						if transfer.PrecalculateBalancesFromApproval.ApproverAddress == to && transfer.PrecalculateBalancesFromApproval.ApprovalLevel == "incoming" {
							userApprovals := toUserBalance.IncomingApprovals
							approvals = types.CastIncomingTransfersToCollectionTransfers(userApprovals, to)
							handled = true
						}

						if transfer.PrecalculateBalancesFromApproval.ApprovalLevel == "outgoing" && !handled && transfer.PrecalculateBalancesFromApproval.ApproverAddress == transfer.From {
							userApprovals := fromUserBalance.OutgoingApprovals
							approvals = types.CastOutgoingTransfersToCollectionTransfers(userApprovals, transfer.From)
							handled = true
						}

						if !handled {
							return sdkerrors.Wrapf(ErrNotImplemented, "could not determine approval to precalculate from %s", transfer.PrecalculateBalancesFromApproval.ApproverAddress)
						}
					}

					//Precaluclate the balances that will be transferred
					transfer.Balances, err = k.GetPredeterminedBalancesForPrecalculationId(
						ctx,
						collection,
						approvals,
						transfer,
						transfer.PrecalculateBalancesFromApproval,
						to,
						initiatedBy,
						transfer.PrecalculationOptions,
					)
					if err != nil {
						return err
					}

					//TODO: Deprecate this in favor of actually calculating the balances in indexer
					amountsJsonData, err := json.Marshal(transfer)
					if err != nil {
						return err
					}
					amountsStr := string(amountsJsonData)

					ctx.EventManager().EmitEvent(
						sdk.NewEvent(sdk.EventTypeMessage,
							sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
							sdk.NewAttribute("creator", initiatedBy),
							sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
							sdk.NewAttribute("transfer", amountsStr),
						),
					)

					ctx.EventManager().EmitEvent(
						sdk.NewEvent("indexer",
							sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
							sdk.NewAttribute("creator", initiatedBy),
							sdk.NewAttribute("collectionId", fmt.Sprint(collection.CollectionId)),
							sdk.NewAttribute("transfer", amountsStr),
						),
					)
				}

				if transfer.From == "Mint" {
					copiedBalances := types.DeepCopyBalances(transfer.Balances)
					totalMinted, err = types.AddBalances(ctx, totalMinted, copiedBalances)
					if err != nil {
						return err
					}
				}

				fromUserBalance, toUserBalance, err = k.HandleTransfer(
					ctx,
					collection,
					transfer,
					fromUserBalance,
					toUserBalance,
					transfer.From,
					to,
					initiatedBy,
					&approvalsUsed,
					&coinTransfers,
				)
				if err != nil {
					return err
				}

				if err := k.SetBalanceForAddress(ctx, collection, to, toUserBalance); err != nil {
					return err
				}

				// Calculate protocol fees for all denoms (0.5% of each denom transferred)
				protocolFees := sdk.NewCoins()
				denomAmounts := make(map[string]sdkmath.Uint)

				for _, coinTransfer := range coinTransfers {
					amount := sdkmath.NewUintFromString(coinTransfer.Amount)
					//initialize it if it doesn't exist
					if _, ok := denomAmounts[coinTransfer.Denom]; !ok {
						denomAmounts[coinTransfer.Denom] = sdkmath.NewUint(0)
					}

					denomAmounts[coinTransfer.Denom] = denomAmounts[coinTransfer.Denom].Add(amount)
				}

				for denom, totalAmount := range denomAmounts {
					// 0.5% of the total amount for this denom
					protocolFee := totalAmount.Mul(sdkmath.NewUint(5)).Quo(sdkmath.NewUint(1000))

					// For other denoms, just use 0.5%
					if !protocolFee.IsZero() {
						protocolFees = protocolFees.Add(sdk.NewCoin(denom, sdkmath.NewIntFromUint64(protocolFee.Uint64())))
					}
				}

				affiliatePercentage := k.GetParams(ctx).AffiliatePercentage // 0 to 10000

				fromAddressAcc, err := sdk.AccAddressFromBech32(initiatedBy)
				if err != nil {
					return err
				}

				if !protocolFees.IsZero() {
					// If no affiliate address is specified, send all fees to community pool
					if transfer.AffiliateAddress == "" {
						err = k.distributionKeeper.FundCommunityPool(ctx, protocolFees, fromAddressAcc)
						if err != nil {
							return sdkerrors.Wrapf(err, "error funding community pool with protocol fees: %s", protocolFees)
						}

						// Add all protocol fees to coinTransfers for community pool
						for _, protocolFee := range protocolFees {
							coinTransfers = append(coinTransfers, CoinTransfers{
								From:          initiatedBy,
								To:            authtypes.NewModuleAddress(distrtypes.ModuleName).String(),
								Amount:        protocolFee.Amount.String(),
								Denom:         protocolFee.Denom,
								IsProtocolFee: true,
							})
						}
					} else {
						// Split protocol fees between community pool and affiliate
						affiliateFees := sdk.NewCoins()
						communityPoolFees := sdk.NewCoins()

						for _, protocolFee := range protocolFees {
							// Calculate affiliate portion (affiliatePercentage is 0-10000, so divide by 10000)
							affiliateAmount := protocolFee.Amount.Mul(sdkmath.NewIntFromUint64(affiliatePercentage.Uint64())).Quo(sdkmath.NewInt(10000))
							communityPoolAmount := protocolFee.Amount.Sub(affiliateAmount)

							if affiliateAmount.GT(sdkmath.ZeroInt()) {
								affiliateFees = affiliateFees.Add(sdk.NewCoin(protocolFee.Denom, affiliateAmount))
							}
							if communityPoolAmount.GT(sdkmath.ZeroInt()) {
								communityPoolFees = communityPoolFees.Add(sdk.NewCoin(protocolFee.Denom, communityPoolAmount))
							}
						}

						// Send affiliate fees to affiliate address
						if !affiliateFees.IsZero() {
							affiliateAddressAcc, err := sdk.AccAddressFromBech32(transfer.AffiliateAddress)
							if err != nil {
								return err
							}

							err = k.bankKeeper.SendCoins(ctx, fromAddressAcc, affiliateAddressAcc, affiliateFees)
							if err != nil {
								return sdkerrors.Wrapf(err, "error sending affiliate fees: %s", affiliateFees)
							}
						}

						// Send remaining fees to community pool
						if !communityPoolFees.IsZero() {
							err = k.distributionKeeper.FundCommunityPool(ctx, communityPoolFees, fromAddressAcc)
							if err != nil {
								return sdkerrors.Wrapf(err, "error funding community pool with protocol fees: %s", communityPoolFees)
							}
						}

						// Add protocol fee transfers to coinTransfers for each denom
						// Add affiliate fee transfers
						for _, affiliateFee := range affiliateFees {
							coinTransfers = append(coinTransfers, CoinTransfers{
								From:          initiatedBy,
								To:            transfer.AffiliateAddress,
								Amount:        affiliateFee.Amount.String(),
								Denom:         affiliateFee.Denom,
								IsProtocolFee: true,
							})
						}

						// Add community pool fee transfers
						for _, communityPoolFee := range communityPoolFees {
							coinTransfers = append(coinTransfers, CoinTransfers{
								From:          initiatedBy,
								To:            authtypes.NewModuleAddress(distrtypes.ModuleName).String(),
								Amount:        communityPoolFee.Amount.String(),
								Denom:         communityPoolFee.Denom,
								IsProtocolFee: true,
							})
						}
					}
				}

				err = emitUsedApprovalDetailsEvent(ctx, collection.CollectionId, transfer.From, to, initiatedBy, coinTransfers, approvalsUsed, transfer.Balances)
				if err != nil {
					return err
				}
			}

			if transfer.From != "Mint" {
				if err := k.SetBalanceForAddress(ctx, collection, transfer.From, fromUserBalance); err != nil {
					return err
				}
			} else {
				// Get current Total
				totalBalances, _ := k.GetBalanceOrApplyDefault(ctx, collection, "Total")
				totalBalances.Balances, err = types.AddBalances(ctx, totalBalances.Balances, totalMinted)
				if err != nil {
					return err
				}

				if err := k.SetBalanceForAddress(ctx, collection, "Total", totalBalances); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func emitUsedApprovalDetailsEvent(ctx sdk.Context, collectionId sdkmath.Uint, from string, to string, initiatedBy string, coinTransfers []CoinTransfers, approvalsUsed []ApprovalsUsed, balances []*types.Balance) (err error) {
	marshalToString := func(v interface{}) (string, error) {
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}

	coinTransfersStr, err := marshalToString(coinTransfers)
	if err != nil {
		return err
	}

	approvalsUsedStr, err := marshalToString(approvalsUsed)
	if err != nil {
		return err
	}

	balancesStr, err := marshalToString(balances)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("usedApprovalDetails",
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute("collectionId", fmt.Sprint(collectionId)),
			sdk.NewAttribute("from", from),
			sdk.NewAttribute("to", to),
			sdk.NewAttribute("initiatedBy", initiatedBy),
			sdk.NewAttribute("coinTransfers", coinTransfersStr),
			sdk.NewAttribute("approvalsUsed", approvalsUsedStr),
			sdk.NewAttribute("balances", balancesStr),
		),
	)

	return nil
}

// Step 1: Check if transfer is allowed on collection level (deducting collection approvals if needed). Will return what userApprovals we need to check.
// Step 2: Check necessary approvals on user level (deducting corresponding approvals if needed)
// Step 3: If all good, we can transfer the balances
func (k Keeper) HandleTransfer(
	ctx sdk.Context,
	collection *types.BadgeCollection,
	transfer *types.Transfer,
	fromUserBalance *types.UserBalanceStore,
	toUserBalance *types.UserBalanceStore,
	from string,
	to string,
	initiatedBy string,
	approvalsUsed *[]ApprovalsUsed,
	coinTransfers *[]CoinTransfers,
) (*types.UserBalanceStore, *types.UserBalanceStore, error) {
	err := *new(error)

	transferBalances := types.DeepCopyBalances(transfer.Balances)
	userApprovals, err := k.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx, collection, transfer, to, initiatedBy, approvalsUsed, coinTransfers)
	if err != nil {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "collection approvals not satisfied")
	}

	if len(userApprovals) > 0 {
		for _, userApproval := range userApprovals {
			newTransfer := &types.Transfer{
				From:                                    from,
				ToAddresses:                             []string{to},
				Balances:                                userApproval.Balances,
				MerkleProofs:                            transfer.MerkleProofs,
				PrioritizedApprovals:                    transfer.PrioritizedApprovals,
				OnlyCheckPrioritizedCollectionApprovals: transfer.OnlyCheckPrioritizedCollectionApprovals,
				OnlyCheckPrioritizedIncomingApprovals:   transfer.OnlyCheckPrioritizedIncomingApprovals,
				OnlyCheckPrioritizedOutgoingApprovals:   transfer.OnlyCheckPrioritizedOutgoingApprovals,
				PrecalculationOptions:                   transfer.PrecalculationOptions,
				AffiliateAddress:                        transfer.AffiliateAddress,
				NumAttempts:                             transfer.NumAttempts,
			}

			if userApproval.Outgoing {
				err = k.DeductUserOutgoingApprovals(ctx, collection, transferBalances, newTransfer, from, to, initiatedBy, fromUserBalance, approvalsUsed, coinTransfers, userApproval.UserRoyalties)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "outgoing approvals for %s not satisfied", from)
				}
			} else {
				err = k.DeductUserIncomingApprovals(ctx, collection, transferBalances, newTransfer, to, initiatedBy, toUserBalance, approvalsUsed, coinTransfers, userApproval.UserRoyalties)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "incoming approvals for %s not satisfied", to)
				}
			}
		}
	}

	for _, balance := range transferBalances {
		//Mint has unlimited balances
		if from != "Mint" {
			fromUserBalance.Balances, err = types.SubtractBalance(ctx, fromUserBalance.Balances, balance, false)
			if err != nil {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "inadequate balances for transfer from %s", from)
			}
		}

		toUserBalance.Balances, err = types.AddBalance(ctx, toUserBalance.Balances, balance)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
		}
	}

	// Get denomination information
	denomInfo := &types.CosmosCoinWrapperPath{}
	isSendingToSpecialAddress := false
	isSendingFromSpecialAddress := false
	for _, path := range collection.CosmosCoinWrapperPaths {
		if path.Address == to {
			isSendingToSpecialAddress = true
			denomInfo = path
		}
		if path.Address == from {
			isSendingFromSpecialAddress = true
			denomInfo = path
		}
	}

	if to == from {
		return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(ErrNotImplemented, "cannot send to self")
	}

	if isSendingFromSpecialAddress || isSendingToSpecialAddress {
		if denomInfo.Denom == "" {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(ErrNotImplemented, "no denom info found for %s", denomInfo.Address)
		}

		conversionBalances := types.DeepCopyBalances(denomInfo.Balances)

		// Little hacky but we find the amount for a specific time and ID
		// Then we will check if it is evenly divisible by the number of transfer balances

		firstBadgeId := transferBalances[0].BadgeIds[0].Start
		firstOwnershipTime := transferBalances[0].OwnershipTimes[0].Start
		firstAmount := transferBalances[0].Amount

		multiplier := sdkmath.NewUint(0)
		for _, balance := range conversionBalances {
			foundBadgeId, err := types.SearchUintRangesForUint(firstBadgeId, balance.BadgeIds)
			if err != nil {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
			}
			foundOwnershipTime, err := types.SearchUintRangesForUint(firstOwnershipTime, balance.OwnershipTimes)
			if err != nil {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
			}
			if foundBadgeId && foundOwnershipTime {
				multiplier = firstAmount.Quo(balance.Amount)
				break
			}
		}

		if multiplier.IsZero() {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(ErrInvalidConversion, "conversion is not evenly divisible")
		}

		conversionBalancesMultiplied := types.DeepCopyBalances(conversionBalances)
		for _, balance := range conversionBalancesMultiplied {
			balance.Amount = balance.Amount.Mul(multiplier)
		}

		transferBalancesCopy := types.DeepCopyBalances(transferBalances)
		remainingBalances, err := types.SubtractBalances(ctx, transferBalancesCopy, conversionBalancesMultiplied)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "conversion is not evenly divisible")
		}

		if len(remainingBalances) > 0 {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(ErrInvalidConversion, "conversion is not evenly divisible")
		}

		ibcDenom := "badges:" + collection.CollectionId.String() + ":" + denomInfo.Denom
		bankKeeper := k.bankKeeper
		amountInt := multiplier.BigInt()
		if isSendingToSpecialAddress {
			if from == "Mint" {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(ErrNotImplemented, "the Mint address cannot perform wrap / unwrap actions")
			}

			userAddressAcc := sdk.MustAccAddressFromBech32(from)

			err = bankKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
			if err != nil {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
			}

			err = bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, userAddressAcc, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
			if err != nil {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
			}
		}

		if isSendingFromSpecialAddress {
			userAddressAcc := sdk.MustAccAddressFromBech32(to)

			err = bankKeeper.SendCoinsFromAccountToModule(ctx, userAddressAcc, types.ModuleName, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
			if err != nil {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
			}

			err = bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.Coins{sdk.NewCoin(ibcDenom, sdkmath.NewIntFromBigInt(amountInt))})
			if err != nil {
				return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
			}
		}
	}

	IsDeleteAfterOneUse := func(autoDeletionOptions *types.AutoDeletionOptions) bool {
		if autoDeletionOptions == nil {
			return false
		}

		if autoDeletionOptions.AfterOneUse {
			return true
		}

		return false
	}

	IsDeleteAfterOverallMaxNumTransfersForCollection := func(autoDeletionOptions *types.AutoDeletionOptions, approvalCriteria *types.ApprovalCriteria, approvalUsed ApprovalsUsed) bool {
		if autoDeletionOptions == nil || !autoDeletionOptions.AfterOverallMaxNumTransfers {
			return false
		}

		// Check if overall max number of transfers threshold is set
		if approvalCriteria == nil || approvalCriteria.MaxNumTransfers == nil || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsNil() || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsZero() {
			return false
		}

		// Get the tracker to check current number of transfers
		maxNumTransfersTrackerId := approvalCriteria.MaxNumTransfers.AmountTrackerId
		if maxNumTransfersTrackerId == "" {
			return false
		}

		// Get the current tracker details
		trackerDetails, err := k.GetApprovalTrackerFromStoreAndResetIfNeeded(
			ctx,
			collection.CollectionId,
			approvalUsed.ApproverAddress,
			approvalUsed.ApprovalId,
			maxNumTransfersTrackerId,
			approvalUsed.ApprovalLevel,
			"overall",
			"",
			approvalCriteria.MaxNumTransfers.ResetTimeIntervals,
			true,
		)
		if err != nil {
			return false
		}

		// Check if the current number of transfers has reached or exceeded the threshold
		return trackerDetails.NumTransfers.GTE(approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers)
	}

	IsDeleteAfterOverallMaxNumTransfersForOutgoing := func(autoDeletionOptions *types.AutoDeletionOptions, approvalCriteria *types.OutgoingApprovalCriteria, approvalUsed ApprovalsUsed) bool {
		if autoDeletionOptions == nil || !autoDeletionOptions.AfterOverallMaxNumTransfers {
			return false
		}

		// Check if overall max number of transfers threshold is set
		if approvalCriteria == nil || approvalCriteria.MaxNumTransfers == nil || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsNil() || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsZero() {
			return false
		}

		// Get the tracker to check current number of transfers
		maxNumTransfersTrackerId := approvalCriteria.MaxNumTransfers.AmountTrackerId
		if maxNumTransfersTrackerId == "" {
			return false
		}

		// Get the current tracker details
		trackerDetails, err := k.GetApprovalTrackerFromStoreAndResetIfNeeded(
			ctx,
			collection.CollectionId,
			approvalUsed.ApproverAddress,
			approvalUsed.ApprovalId,
			maxNumTransfersTrackerId,
			approvalUsed.ApprovalLevel,
			"overall",
			"",
			approvalCriteria.MaxNumTransfers.ResetTimeIntervals,
			true,
		)
		if err != nil {
			return false
		}

		// Check if the current number of transfers has reached or exceeded the threshold
		return trackerDetails.NumTransfers.GTE(approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers)
	}

	IsDeleteAfterOverallMaxNumTransfersForIncoming := func(autoDeletionOptions *types.AutoDeletionOptions, approvalCriteria *types.IncomingApprovalCriteria, approvalUsed ApprovalsUsed) bool {
		if autoDeletionOptions == nil || !autoDeletionOptions.AfterOverallMaxNumTransfers {
			return false
		}

		// Check if overall max number of transfers threshold is set
		if approvalCriteria == nil || approvalCriteria.MaxNumTransfers == nil || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsNil() || approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers.IsZero() {
			return false
		}

		// Get the tracker to check current number of transfers
		maxNumTransfersTrackerId := approvalCriteria.MaxNumTransfers.AmountTrackerId
		if maxNumTransfersTrackerId == "" {
			return false
		}

		// Get the current tracker details
		trackerDetails, err := k.GetApprovalTrackerFromStoreAndResetIfNeeded(
			ctx,
			collection.CollectionId,
			approvalUsed.ApproverAddress,
			approvalUsed.ApprovalId,
			maxNumTransfersTrackerId,
			approvalUsed.ApprovalLevel,
			"overall",
			"",
			approvalCriteria.MaxNumTransfers.ResetTimeIntervals,
			true,
		)
		if err != nil {
			return false
		}

		// Check if the current number of transfers has reached or exceeded the threshold
		return trackerDetails.NumTransfers.GTE(approvalCriteria.MaxNumTransfers.OverallMaxNumTransfers)
	}

	// Per-transfer, we handle auto-deletions if applicable
	for _, approvalUsed := range *approvalsUsed {
		if approvalUsed.ApprovalLevel == "incoming" {
			newIncomingApprovals := []*types.UserIncomingApproval{}
			for _, incomingApproval := range toUserBalance.IncomingApprovals {
				if incomingApproval.ApprovalId != approvalUsed.ApprovalId {
					newIncomingApprovals = append(newIncomingApprovals, incomingApproval)
				} else {
					shouldDelete := false

					// Check if should delete after one use (doesn't depend on ApprovalCriteria)
					if incomingApproval.ApprovalCriteria != nil && incomingApproval.ApprovalCriteria.AutoDeletionOptions != nil {
						shouldDelete = IsDeleteAfterOneUse(incomingApproval.ApprovalCriteria.AutoDeletionOptions)
					}

					// Check if should delete after overall max transfers (depends on ApprovalCriteria)
					if !shouldDelete && incomingApproval.ApprovalCriteria != nil {
						shouldDelete = IsDeleteAfterOverallMaxNumTransfersForIncoming(incomingApproval.ApprovalCriteria.AutoDeletionOptions, incomingApproval.ApprovalCriteria, approvalUsed)
					}

					if !shouldDelete {
						newIncomingApprovals = append(newIncomingApprovals, incomingApproval)
					} else {
						// Delete the approval
					}
				}
			}
			toUserBalance.IncomingApprovals = newIncomingApprovals
		} else if approvalUsed.ApprovalLevel == "outgoing" {
			newOutgoingApprovals := []*types.UserOutgoingApproval{}
			for _, outgoingApproval := range fromUserBalance.OutgoingApprovals {
				if outgoingApproval.ApprovalId != approvalUsed.ApprovalId {
					newOutgoingApprovals = append(newOutgoingApprovals, outgoingApproval)
				} else {
					shouldDelete := false

					// Check if should delete after one use (doesn't depend on ApprovalCriteria)
					if outgoingApproval.ApprovalCriteria != nil && outgoingApproval.ApprovalCriteria.AutoDeletionOptions != nil {
						shouldDelete = IsDeleteAfterOneUse(outgoingApproval.ApprovalCriteria.AutoDeletionOptions)
					}

					// Check if should delete after overall max transfers (depends on ApprovalCriteria)
					if !shouldDelete && outgoingApproval.ApprovalCriteria != nil {
						shouldDelete = IsDeleteAfterOverallMaxNumTransfersForOutgoing(outgoingApproval.ApprovalCriteria.AutoDeletionOptions, outgoingApproval.ApprovalCriteria, approvalUsed)
					}

					if !shouldDelete {
						newOutgoingApprovals = append(newOutgoingApprovals, outgoingApproval)
					} else {
						// Delete the approval
					}
				}
			}
			fromUserBalance.OutgoingApprovals = newOutgoingApprovals
		} else if approvalUsed.ApprovalLevel == "collection" {
			newCollectionApprovals := []*types.CollectionApproval{}
			edited := false
			for _, collectionApproval := range collection.CollectionApprovals {
				if collectionApproval.ApprovalId != approvalUsed.ApprovalId {
					newCollectionApprovals = append(newCollectionApprovals, collectionApproval)
				} else {
					shouldDelete := false

					// Check if should delete after one use (doesn't depend on ApprovalCriteria)
					if collectionApproval.ApprovalCriteria != nil && collectionApproval.ApprovalCriteria.AutoDeletionOptions != nil {
						shouldDelete = IsDeleteAfterOneUse(collectionApproval.ApprovalCriteria.AutoDeletionOptions)
					}

					// Check if should delete after overall max transfers (depends on ApprovalCriteria)
					if !shouldDelete && collectionApproval.ApprovalCriteria != nil {
						shouldDelete = IsDeleteAfterOverallMaxNumTransfersForCollection(collectionApproval.ApprovalCriteria.AutoDeletionOptions, collectionApproval.ApprovalCriteria, approvalUsed)
					}

					if !shouldDelete {
						newCollectionApprovals = append(newCollectionApprovals, collectionApproval)
					} else {
						// Delete the approval
						edited = true
					}
				}
			}

			collection.CollectionApprovals = newCollectionApprovals
			if edited {
				err = k.SetCollectionInStore(ctx, collection)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
				}
			}
		}
	}

	return fromUserBalance, toUserBalance, nil
}
