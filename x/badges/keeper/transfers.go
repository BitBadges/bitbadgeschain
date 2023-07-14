package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
)

func (k Keeper) HandleTransfers(ctx sdk.Context, collection *types.BadgeCollection, transfers []*types.Transfer, initiatedBy string) error {
	err := *new(error)

	for _, transfer := range transfers {
		fromBalanceKey := ConstructBalanceKey(transfer.From, collection.CollectionId)
		fromUserBalance, found := k.GetUserBalanceFromStore(ctx, fromBalanceKey)
		if !found {
			return sdkerrors.Wrapf(ErrUserBalanceNotExists, "from user balance for %s does not exist", transfer.From)
		}

		for _, to := range transfer.ToAddresses {
			toBalanceKey := ConstructBalanceKey(to, collection.CollectionId)
			toUserBalance, found := k.GetUserBalanceFromStore(ctx, toBalanceKey)
			if !found {
				toUserBalance = &types.UserBalanceStore{
					Balances:                          []*types.Balance{},
					ApprovedOutgoingTransfersTimeline: collection.DefaultUserApprovedOutgoingTransfersTimeline,
					ApprovedIncomingTransfersTimeline: collection.DefaultUserApprovedIncomingTransfersTimeline,
					UserPermissions: 								   collection.DefaultUserPermissions,
				}
			}

			if transfer.PrecalculateFromApproval != nil {
				approvedTransfers := types.GetCurrentCollectionApprovedTransfers(ctx, collection)
				if transfer.PrecalculateFromApproval.Address != "" {
					if transfer.PrecalculateFromApproval.Address != to && transfer.PrecalculateFromApproval.Address != transfer.From {
						return sdkerrors.Wrapf(ErrNotImplemented, "approval id address %s does not match to or from address", transfer.PrecalculateFromApproval.Address)
					}

					if transfer.PrecalculateFromApproval.Address == to {
						if transfer.PrecalculateFromApproval.ApprovalLevel == "incoming" {
							userApprovedTransfers := types.GetCurrentUserApprovedIncomingTransfers(ctx, toUserBalance)
							approvedTransfers = types.CastIncomingTransfersToCollectionTransfers(userApprovedTransfers, to)
						} else {
							userApprovedTransfers := types.GetCurrentUserApprovedOutgoingTransfers(ctx, toUserBalance)
							approvedTransfers = types.CastOutgoingTransfersToCollectionTransfers(userApprovedTransfers, to)
						}
					} else {
						if transfer.PrecalculateFromApproval.ApprovalLevel == "outgoing" {
							userApprovedTransfers := types.GetCurrentUserApprovedOutgoingTransfers(ctx, fromUserBalance)
							approvedTransfers = types.CastOutgoingTransfersToCollectionTransfers(userApprovedTransfers, transfer.From)
						} else {
							userApprovedTransfers := types.GetCurrentUserApprovedIncomingTransfers(ctx, fromUserBalance)
							approvedTransfers = types.CastIncomingTransfersToCollectionTransfers(userApprovedTransfers, transfer.From)
						}
					}
				}

				//Precaluclate the balances that will be transferred
				transfer.Balances, err = k.GetPredeterminedBalancesForApprovalId(ctx, approvedTransfers, collection, transfer.PrecalculateFromApproval.ApprovalId, transfer.PrecalculateFromApproval.ApprovalLevel, transfer.PrecalculateFromApproval.Address, transfer.MerkleProofs, initiatedBy)
				if err != nil {
					return err
				}
			}

			for _, balance := range transfer.Balances {
				fromUserBalance, toUserBalance, err = k.HandleTransfer(ctx, collection, balance.BadgeIds, balance.OwnedTimes, fromUserBalance, toUserBalance, balance.Amount, transfer.From, to, initiatedBy, transfer.MerkleProofs)
				if err != nil {
					return err
				}
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
func (k Keeper) HandleTransfer(ctx sdk.Context, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromUserBalance *types.UserBalanceStore, toUserBalance *types.UserBalanceStore, amount sdkmath.Uint, from string, to string, initiatedBy string, solutions []*types.MerkleProof) (*types.UserBalanceStore, *types.UserBalanceStore, error) {
	err := *new(error)
	transferBalances :=  []*types.Balance{ { Amount: amount, BadgeIds: badgeIds, OwnedTimes: times } }
	

	for _, balance := range transferBalances {
		badgeIds := balance.BadgeIds
		times := balance.OwnedTimes
		amount := balance.Amount
		userApprovals, err := k.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx, transferBalances, collection, badgeIds, times, from, to, initiatedBy, amount, solutions)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
		}

		if len(userApprovals) > 0 {
			for _, userApproval := range userApprovals {
				for _, balance := range userApproval.Balances {
					if userApproval.Outgoing {
						err = k.DeductUserOutgoingApprovals(ctx, transferBalances, collection, fromUserBalance, balance.BadgeIds, balance.OwnedTimes, from, to, initiatedBy, amount, solutions)
						if err != nil {
							return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
						}
					} else {
						err = k.DeductUserIncomingApprovals(ctx, transferBalances, collection, toUserBalance, balance.BadgeIds, balance.OwnedTimes, from, to, initiatedBy, amount, solutions)
						if err != nil {
							return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
						}
					}
				}
			}
		}
	}

	for _, balance := range transferBalances {
		fromUserBalance.Balances, err = types.SubtractBalance(fromUserBalance.Balances, balance)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
		}

		toUserBalance.Balances, err = types.AddBalance(toUserBalance.Balances, balance)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
		}
	}


	return fromUserBalance, toUserBalance, nil
}
