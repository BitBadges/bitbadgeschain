package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
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

	for _, transfer := range transfers {
		fromUserBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, transfer.From)

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

			totalUbadgeTransferred := sdkmath.NewUint(0)
			for _, coinTransfer := range coinTransfers {
				if coinTransfer.Denom == "ubadge" {
					amount := sdkmath.NewUintFromString(coinTransfer.Amount)
					totalUbadgeTransferred = totalUbadgeTransferred.Add(amount)
				}
			}

			// We take max(0.5% or k.FixedCostPerTransfer) as protocol fee
			fixedCost, err := sdk.ParseCoinNormalized(k.FixedCostPerTransfer)
			if err != nil {
				return err
			}

			cost := fixedCost
			//0.5% of the total ubadge transferred
			protocolFee := totalUbadgeTransferred.Mul(sdkmath.NewUint(5)).Quo(sdkmath.NewUint(1000))
			if protocolFee.GTE(sdkmath.Uint(fixedCost.Amount)) {
				cost = sdk.NewCoin("ubadge", sdkmath.NewIntFromUint64(protocolFee.Uint64()))
			}

			payoutAddress := k.PayoutAddress
			if transfer.AffiliateAddress != "" {
				payoutAddress = transfer.AffiliateAddress
			}

			payoutAddressAcc, err := sdk.AccAddressFromBech32(payoutAddress)
			if err != nil {
				return err
			}

			fromAddressAcc, err := sdk.AccAddressFromBech32(initiatedBy)
			if err != nil {
				return err
			}

			err = k.bankKeeper.SendCoins(ctx, fromAddressAcc, payoutAddressAcc, sdk.NewCoins(cost))
			if err != nil {
				return sdkerrors.Wrapf(err, "error completing required payout. each transfer costs %s", cost)
			}

			coinTransfers = append(coinTransfers, CoinTransfers{
				From:          initiatedBy,
				To:            payoutAddress,
				Amount:        cost.Amount.String(),
				Denom:         cost.Denom,
				IsProtocolFee: true,
			})

			err = emitUsedApprovalDetailsEvent(ctx, collection.CollectionId, transfer.From, to, initiatedBy, coinTransfers, approvalsUsed, transfer.Balances)
			if err != nil {
				return err
			}
		}

		if transfer.From != "Mint" {
			if err := k.SetBalanceForAddress(ctx, collection, transfer.From, fromUserBalance); err != nil {
				return err
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
			}

			if userApproval.Outgoing {
				err = k.DeductUserOutgoingApprovals(ctx, collection, transferBalances, newTransfer, from, to, initiatedBy, fromUserBalance, approvalsUsed, coinTransfers)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "outgoing approvals for %s not satisfied", from)
				}
			} else {
				err = k.DeductUserIncomingApprovals(ctx, collection, transferBalances, newTransfer, to, initiatedBy, toUserBalance, approvalsUsed, coinTransfers)
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

	IsDeleteAfterOneUse := func(autoDeletionOptions *types.AutoDeletionOptions) bool {
		if autoDeletionOptions == nil {
			return false
		}

		if autoDeletionOptions.AfterOneUse {
			return true
		}

		return false
	}

	// Per-transfer, we handle auto-deletions if applicable
	for _, approvalUsed := range *approvalsUsed {
		if approvalUsed.ApprovalLevel == "incoming" {
			newIncomingApprovals := []*types.UserIncomingApproval{}
			for _, incomingApproval := range fromUserBalance.IncomingApprovals {
				if incomingApproval.ApprovalId != approvalUsed.ApprovalId {
					newIncomingApprovals = append(newIncomingApprovals, incomingApproval)
				} else {
					if incomingApproval.ApprovalCriteria == nil || !IsDeleteAfterOneUse(incomingApproval.ApprovalCriteria.AutoDeletionOptions) {
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
					if outgoingApproval.ApprovalCriteria == nil || !IsDeleteAfterOneUse(outgoingApproval.ApprovalCriteria.AutoDeletionOptions) {
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
					if collectionApproval.ApprovalCriteria == nil || !IsDeleteAfterOneUse(collectionApproval.ApprovalCriteria.AutoDeletionOptions) {
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
