package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Only handles from => to (pending and forceful) (not other way around)
func (k msgServer) TransferBadge(goCtx context.Context, msg *types.MsgTransferBadge) (*types.MsgTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	rangesToValidate := []*types.IdRange{}
	for _, transfer := range msg.Transfers {
		for _, balance := range transfer.Balances {
			for _, badgeIdRange := range balance.BadgeIds {
				rangesToValidate = append(rangesToValidate, badgeIdRange)
			}
		}
	}

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		BadgeIdRangesToValidate:     rangesToValidate,
	})
	if err != nil {
		return nil, err
	}

	fromBalanceKey := ConstructBalanceKey(msg.From, msg.CollectionId)
	fromUserBalance, found := k.Keeper.GetUserBalanceFromStore(ctx, fromBalanceKey)
	if !found {
		return nil, ErrUserBalanceNotExists
	}

	for _, transfer := range msg.Transfers {
		for _, to := range transfer.ToAddresses {
			toBalanceKey := ConstructBalanceKey(to, msg.CollectionId)
			toUserBalance, found := k.Keeper.GetUserBalanceFromStore(ctx, toBalanceKey)
			if !found {
				toUserBalance = types.UserBalanceStore{}
			}

			for _, balance := range transfer.Balances {
				amount := balance.Amount

				for _, badgeIdRange := range balance.BadgeIds {
					fromUserBalance, toUserBalance, err = HandleTransfer(collection, badgeIdRange, fromUserBalance, toUserBalance, amount, msg.From, to, msg.Creator)
					if err != nil {
						return nil, err
					}
				}

			}

			if err := k.SetUserBalanceInStore(ctx, toBalanceKey, toUserBalance); err != nil {
				return nil, err
			}
		}
	}

	if err := k.SetUserBalanceInStore(ctx, fromBalanceKey, fromUserBalance); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute("collection_id", fmt.Sprint(msg.CollectionId)),
		),
	)

	return &types.MsgTransferBadgeResponse{}, nil
}
