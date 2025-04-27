package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
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

func (k Keeper) HandleTransfers(ctx sdk.Context, collection *types.BadgeCollection, transfers []*types.Transfer, initiatedBy string) error {
	err := *new(error)

	for _, transfer := range transfers {
		fromUserBalance, _ := k.GetBalanceOrApplyDefault(ctx, collection, transfer.From)

		for _, to := range transfer.ToAddresses {
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

			fromUserBalance, toUserBalance, err = k.HandleTransfer(ctx, collection, transfer, fromUserBalance, toUserBalance, transfer.From, to, initiatedBy)
			if err != nil {
				return err
			}

			if err := k.SetBalanceForAddress(ctx, collection, to, toUserBalance); err != nil {
				return err
			}

			if k.PayoutAddress != "" && k.FixedCostPerTransfer != "" {
				cost, err := sdk.ParseCoinNormalized(k.FixedCostPerTransfer)
				if err != nil {
					return err
				}

				payoutAddressAcc, err := sdk.AccAddressFromBech32(k.PayoutAddress)
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
) (*types.UserBalanceStore, *types.UserBalanceStore, error) {
	err := *new(error)

	transferBalances := types.DeepCopyBalances(transfer.Balances)
	userApprovals, err := k.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx, collection, transfer, to, initiatedBy)
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
			}

			if userApproval.Outgoing {
				err = k.DeductUserOutgoingApprovals(ctx, collection, transferBalances, newTransfer, from, to, initiatedBy, fromUserBalance)
				if err != nil {
					return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "outgoing approvals for %s not satisfied", from)
				}
			} else {

				err = k.DeductUserIncomingApprovals(ctx, collection, transferBalances, newTransfer, to, initiatedBy, toUserBalance)
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

	return fromUserBalance, toUserBalance, nil
}
