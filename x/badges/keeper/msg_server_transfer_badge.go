package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Only handles from => to (pending and forceful) (not other way around)
func (k msgServer) TransferBadge(goCtx context.Context, msg *types.MsgTransferBadge) (*types.MsgTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	accsToCheck := []uint64{msg.From}
	for _, transfer := range msg.Transfers {
		accsToCheck = append(accsToCheck, transfer.ToAddresses...)
	}

	rangesToValidate := []*types.IdRange{}
	for _, transfer := range msg.Transfers {
		for _, balance := range transfer.Balances {
			for _, badgeIdRange := range balance.BadgeIds {
				rangesToValidate = append(rangesToValidate, badgeIdRange)
			}
		}
	}

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		BadgeIdRangesToValidate:     rangesToValidate,
		AccountsToCheckRegistration: accsToCheck,
	})
	if err != nil {
		return nil, err
	}

	fromBalanceKey := ConstructBalanceKey(msg.From, msg.CollectionId)
	fromUserBalance, found := k.Keeper.GetUserBalanceFromStore(ctx, fromBalanceKey)
	if !found {
		return nil, ErrUserBalanceNotExists
	}

	newUserBalances := []types.UserBalance{}
	newUserBalanceAccounts := []uint64{}

	for _, transfer := range msg.Transfers {
		for _, to := range transfer.ToAddresses {
			toBalanceKey := ConstructBalanceKey(to, msg.CollectionId)
			toUserBalance, found := k.Keeper.GetUserBalanceFromStore(ctx, toBalanceKey)
			if !found {
				toUserBalance = types.UserBalance{}
			}

			for _, balance := range transfer.Balances {
				amount := balance.Balance

				for _, badgeIdRange := range balance.BadgeIds {
					fromUserBalance, toUserBalance, err = HandleTransfer(badge, badgeIdRange, fromUserBalance, toUserBalance, amount, msg.From, to, CreatorAccountNum)
					if err != nil {
						return nil, err
					}
				}

			}

			if err := k.SetUserBalanceInStore(ctx, toBalanceKey, toUserBalance); err != nil {
				return nil, err
			}
			newUserBalances = append(newUserBalances, toUserBalance)
			newUserBalanceAccounts = append(newUserBalanceAccounts, to)
		}
	}



	if err := k.SetUserBalanceInStore(ctx, fromBalanceKey, fromUserBalance); err != nil {
		return nil, err
	}
	newUserBalances = append(newUserBalances, fromUserBalance)
	newUserBalanceAccounts = append(newUserBalanceAccounts, msg.From)

	newUserBalancesJson, err := json.Marshal(newUserBalances)
	if err != nil {
		return nil, err
	}

	newUserBalanceAccountsJson, err := json.Marshal(newUserBalanceAccounts)
	if err != nil {
		return nil, err
	}
	

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute("collection_id", fmt.Sprint(msg.CollectionId)),
			sdk.NewAttribute("new_balances", string(newUserBalancesJson)),
			sdk.NewAttribute("new_balance_accounts", string(newUserBalanceAccountsJson)),
		),
	)

	return &types.MsgTransferBadgeResponse{}, nil
}
