package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
)

func (k Keeper) HandleTransfers(ctx sdk.Context, collection *types.BadgeCollection, transfers []*types.Transfer, initiatedBy string) error {
	err := *new(error)

	for _, transfer := range transfers {
		fromBalanceKey := ConstructBalanceKey(transfer.From, collection.CollectionId)
		fromUserBalance, found := k.GetUserBalanceFromStore(ctx, fromBalanceKey)
		if !found {
			if transfer.From == "Mint" {
				return sdkerrors.Wrapf(ErrUserBalanceNotExists, "sender user balance (Mint) for %s is empty or does not exist", transfer.From)
			} else {
				fromUserBalance = &types.UserBalanceStore{
					Balances:         collection.DefaultBalances.Balances,
					OutgoingApprovals: collection.DefaultBalances.OutgoingApprovals,
					IncomingApprovals: collection.DefaultBalances.IncomingApprovals,
					AutoApproveSelfInitiatedOutgoingTransfers: collection.DefaultBalances.AutoApproveSelfInitiatedOutgoingTransfers,
					AutoApproveSelfInitiatedIncomingTransfers: collection.DefaultBalances.AutoApproveSelfInitiatedIncomingTransfers,
					UserPermissions: collection.DefaultBalances.UserPermissions,
				}
			}
		}

		for _, to := range transfer.ToAddresses {
			toBalanceKey := ConstructBalanceKey(to, collection.CollectionId)
			toUserBalance, found := k.GetUserBalanceFromStore(ctx, toBalanceKey)
			if !found {
				toUserBalance = &types.UserBalanceStore{
					Balances:         collection.DefaultBalances.Balances,
					OutgoingApprovals: collection.DefaultBalances.OutgoingApprovals,
					IncomingApprovals: collection.DefaultBalances.IncomingApprovals,
					AutoApproveSelfInitiatedOutgoingTransfers: collection.DefaultBalances.AutoApproveSelfInitiatedOutgoingTransfers,
					AutoApproveSelfInitiatedIncomingTransfers: collection.DefaultBalances.AutoApproveSelfInitiatedIncomingTransfers,
					UserPermissions: collection.DefaultBalances.UserPermissions,
				}
			}

			if transfer.PrecalculateBalancesFromApproval != nil && transfer.PrecalculateBalancesFromApproval.ApprovalId != "" {
				approvals := collection.CollectionApprovals
				if transfer.PrecalculateBalancesFromApproval.ApproverAddress != "" {
					if transfer.PrecalculateBalancesFromApproval.ApproverAddress != to && transfer.PrecalculateBalancesFromApproval.ApproverAddress != transfer.From {
						return sdkerrors.Wrapf(ErrNotImplemented, "approval id address %s does not match to or from address", transfer.PrecalculateBalancesFromApproval.ApproverAddress)
					}

					if transfer.PrecalculateBalancesFromApproval.ApproverAddress == to {
						if transfer.PrecalculateBalancesFromApproval.ApprovalLevel == "incoming" {
							userApprovals := toUserBalance.IncomingApprovals
							approvals = types.CastIncomingTransfersToCollectionTransfers(userApprovals, to)
						} else {
							userApprovals := toUserBalance.OutgoingApprovals
							approvals = types.CastOutgoingTransfersToCollectionTransfers(userApprovals, to)
						}
					} else {
						if transfer.PrecalculateBalancesFromApproval.ApprovalLevel == "outgoing" {
							userApprovals := fromUserBalance.OutgoingApprovals
							approvals = types.CastOutgoingTransfersToCollectionTransfers(userApprovals, transfer.From)
						} else {
							userApprovals := fromUserBalance.IncomingApprovals
							approvals = types.CastIncomingTransfersToCollectionTransfers(userApprovals, transfer.From)
						}
					}
				}

				//Precaluclate the balances that will be transferred
				transfer.Balances, err = k.GetPredeterminedBalancesForPrecalculationId(ctx, approvals, collection, "", transfer.PrecalculateBalancesFromApproval.ApprovalId, transfer.PrecalculateBalancesFromApproval.ApprovalLevel, transfer.PrecalculateBalancesFromApproval.ApproverAddress, transfer.MerkleProofs, initiatedBy)
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
			}

			challengeIdsIncremented := &[]string{}
			trackerIdsIncremented := &[]string{}
			fromUserBalance, toUserBalance, err = k.HandleTransfer(ctx, collection, transfer.Balances, fromUserBalance, toUserBalance, transfer.From, to, initiatedBy, transfer.MerkleProofs, challengeIdsIncremented, trackerIdsIncremented, transfer.PrioritizedApprovals, transfer.OnlyCheckPrioritizedApprovals)
			if err != nil {
				return err
			}

			if err := k.SetUserBalanceInStore(ctx, toBalanceKey, toUserBalance); err != nil {
				return err
			}
		}

		if err := k.SetUserBalanceInStore(ctx, fromBalanceKey, fromUserBalance); err != nil {
			return err
		}
	}

	return nil
}

// Step 1: Check if transfer is allowed on collection level (deducting approvals if needed)
// Step 2: If not overriden by collection, check necessary approvals on user level (deducting approvals if needed)
// Step 3: If all good, we can transfer the balances
func (k Keeper) HandleTransfer(ctx sdk.Context, collection *types.BadgeCollection, transferBalances []*types.Balance, fromUserBalance *types.UserBalanceStore, toUserBalance *types.UserBalanceStore, from string, to string, initiatedBy string, solutions []*types.MerkleProof, challengeIdsIncremented *[]string, trackerIdsIncremented *[]string, prioritizedApprovals []*types.ApprovalIdentifierDetails, onlyCheckPrioritized bool) (*types.UserBalanceStore, *types.UserBalanceStore, error) {
	err := *new(error)

	for _, balance := range transferBalances {
		badgeIds := balance.BadgeIds
		times := balance.OwnershipTimes
		amount := balance.Amount
		userApprovals, err := k.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx, transferBalances, collection, badgeIds, times, from, to, initiatedBy, amount, solutions, challengeIdsIncremented, trackerIdsIncremented, prioritizedApprovals, onlyCheckPrioritized)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "collection approvals not satisfied")
		}

		if len(userApprovals) > 0 {
			for _, userApproval := range userApprovals {
				for _, balance := range userApproval.Balances {
					if userApproval.Outgoing {
						err = k.DeductUserOutgoingApprovals(ctx, transferBalances, collection, fromUserBalance, balance.BadgeIds, balance.OwnershipTimes, from, to, initiatedBy, amount, solutions, challengeIdsIncremented, trackerIdsIncremented, prioritizedApprovals, onlyCheckPrioritized)
						if err != nil {
							return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "outgoing approvals for %s not satisfied", from)
						}
					} else {
						err = k.DeductUserIncomingApprovals(ctx, transferBalances, collection, toUserBalance, balance.BadgeIds, balance.OwnershipTimes, from, to, initiatedBy, amount, solutions, challengeIdsIncremented, trackerIdsIncremented, prioritizedApprovals, onlyCheckPrioritized)
						if err != nil {
							return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "incoming approvals for %s not satisfied", to)
						}
					}
				}
			}
		}
	}

	for _, balance := range transferBalances {
		fromUserBalance.Balances, err = types.SubtractBalance(ctx, fromUserBalance.Balances, balance, false)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, sdkerrors.Wrapf(err, "inadequate balances for transfer from %s", from)
		}

		toUserBalance.Balances, err = types.AddBalance(ctx, toUserBalance.Balances, balance)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
		}
	}

	return fromUserBalance, toUserBalance, nil
}
