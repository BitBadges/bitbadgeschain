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
					UserPermissions: &types.UserPermissions{
						CanUpdateApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransferPermission{},
						CanUpdateApprovedIncomingTransfers: []*types.UserApprovedIncomingTransferPermission{},
					},
				}
			}

			for _, balance := range transfer.Balances {
				fromUserBalance, toUserBalance, err = k.HandleTransfer(ctx, collection, balance.BadgeIds, balance.OwnershipTimes, fromUserBalance, toUserBalance, balance.Amount, transfer.From, to, initiatedBy, transfer.Solutions, transfer.TransferAsMuchAsPossible)
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
func (k Keeper) HandleTransfer(ctx sdk.Context, collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, fromUserBalance *types.UserBalanceStore, toUserBalance *types.UserBalanceStore, amount sdkmath.Uint, from string, to string, initiatedBy string, solutions []*types.ChallengeSolution, transferAsMuchAsPossible bool) (*types.UserBalanceStore, *types.UserBalanceStore, error) {
	err := *new(error)
	
	transferBalances :=  []*types.Balance{ { Amount: amount, BadgeIds: badgeIds, OwnershipTimes: times } }


	if transferAsMuchAsPossible {
		userApprovals := []*UserApprovalsToCheck{}
		userApprovals, transferBalances, err = k.GetUnapprovedBalancesForCollection(ctx, collection, transferBalances, from, to, initiatedBy, solutions, true)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
		}

		//here, we may double remove some balances
		if len(userApprovals) > 0 {
			for _, userApproval := range userApprovals {
				if userApproval.Outgoing {
					_, transferBalances, err = k.GetUnapprovedBalancesForOutgoing(ctx, collection, fromUserBalance, transferBalances, userApproval.Balances, from, to, initiatedBy, solutions, true)
					if err != nil {
						return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
					}
				} else {
					_, transferBalances, err = k.GetUnapprovedBalancesForIncoming(ctx, collection, toUserBalance, transferBalances, userApproval.Balances, from, to, initiatedBy, solutions, true)
					if err != nil {
						return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
					}
				}
			}
		}
	} 

	//TODO: make this balances[] instead of passing in times, etc

	for _, balance := range transferBalances {
		badgeIds := balance.BadgeIds
		times := balance.OwnershipTimes
		amount := balance.Amount
		userApprovals, err := k.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx, collection, badgeIds, times, from, to, initiatedBy, amount, solutions, false)
		if err != nil {
			return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
		}

		if len(userApprovals) > 0 {
			for _, userApproval := range userApprovals {
				for _, balance := range userApproval.Balances {
					if userApproval.Outgoing {
						err = k.DeductUserOutgoingApprovals(ctx, collection, fromUserBalance, balance.BadgeIds, balance.OwnershipTimes, from, to, initiatedBy, amount, solutions, false)
						if err != nil {
							return &types.UserBalanceStore{}, &types.UserBalanceStore{}, err
						}
					} else {
						err = k.DeductUserIncomingApprovals(ctx, collection, toUserBalance, balance.BadgeIds, balance.OwnershipTimes, from, to, initiatedBy, amount, solutions, false)
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
