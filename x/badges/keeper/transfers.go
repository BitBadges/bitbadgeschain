package keeper

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func HandleManagerAlias(address string, managerAddress string) string {
	if address == "Manager" {
		return managerAddress
	}

	return address
}

func (k Keeper) HandleTransfers(ctx sdk.Context, collection types.BadgeCollection, transfers []*types.Transfer, initiatedBy string, onlyDeductApprovals bool) (error) {
	//If "Mint" or "Manager" 
	unmintedBalances := types.UserBalanceStore{
		Balances: collection.UnmintedSupplys,
	}
	err := *new(error)
	manager := GetCurrentManager(ctx, collection)

	initiatedBy = HandleManagerAlias(initiatedBy, manager)

	for _, transfer := range transfers {
		fromBalanceKey := ConstructBalanceKey(transfer.From, collection.CollectionId)
		fromUserBalance := unmintedBalances
		found := true
		transfer.From = HandleManagerAlias(transfer.From, manager)
		
		if transfer.From != "Mint" {
			fromUserBalance, found = k.GetUserBalanceFromStore(ctx, fromBalanceKey)
			if !found {
				return ErrUserBalanceNotExists
			}
		}
		
		for _, to := range transfer.ToAddresses {
			to = HandleManagerAlias(to, manager)

			toBalanceKey := ConstructBalanceKey(to, collection.CollectionId)
			toUserBalance, found := k.GetUserBalanceFromStore(ctx, toBalanceKey)
			if !found {
				toUserBalance = types.UserBalanceStore{
					Balances : []*types.Balance{},
					ApprovedOutgoingTransfersTimeline: collection.DefaultUserApprovedOutgoingTransfersTimeline,
					ApprovedIncomingTransfersTimeline: collection.DefaultUserApprovedIncomingTransfersTimeline,
					Permissions: &types.UserPermissions{
						CanUpdateApprovedOutgoingTransfers: []*types.UserApprovedTransferPermission{},
						CanUpdateApprovedIncomingTransfers: []*types.UserApprovedTransferPermission{},
					},
				}
			}

			for _, balance := range transfer.Balances {
				fromUserBalance, toUserBalance, err = k.HandleTransfer(ctx, collection, balance.BadgeIds, balance.Times, fromUserBalance, toUserBalance, balance.Amount, transfer.From, to, initiatedBy, transfer.Solutions, onlyDeductApprovals)
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

// Forceful transfers will transfer the balances and deduct from approvals directly without adding it to pending.
func (k Keeper) HandleTransfer(ctx sdk.Context, collection types.BadgeCollection, badgeIds []*types.IdRange, times []*types.IdRange, fromUserBalance types.UserBalanceStore, toUserBalance types.UserBalanceStore, amount sdk.Uint, from string, to string, initiatedBy string, solutions []*types.ChallengeSolution, onlyDeductApprovals bool) (types.UserBalanceStore, types.UserBalanceStore, error) {
	err := *new(error)

	userApprovals, err := k.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx,  collection, badgeIds, times,  from, to, initiatedBy, amount, solutions)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	if len(userApprovals) > 0 {
		for _, userApproval := range userApprovals {
			if userApproval.Outgoing {
				err = k.DeductUserOutgoingApprovals(ctx, collection, &fromUserBalance, userApproval.BadgeIds, times, from, to, initiatedBy, amount, solutions)
				if err != nil {
					return types.UserBalanceStore{}, types.UserBalanceStore{}, err
				}
			} else {
				err = k.DeductUserIncomingApprovals(ctx, collection, &toUserBalance, userApproval.BadgeIds, times, from, to, initiatedBy, amount, solutions)
				if err != nil {
					return types.UserBalanceStore{}, types.UserBalanceStore{}, err
				}
			}
		}
	}

	if onlyDeductApprovals {
		return fromUserBalance, toUserBalance, nil
	}

	fromUserBalance.Balances, err = types.SubtractBalancesForIdRanges(fromUserBalance.Balances, badgeIds, times, amount)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	toUserBalance.Balances, err = types.AddBalancesForIdRanges(toUserBalance.Balances, badgeIds, times, amount)
	if err != nil {
		return types.UserBalanceStore{}, types.UserBalanceStore{}, err
	}

	return fromUserBalance, toUserBalance, nil
}
