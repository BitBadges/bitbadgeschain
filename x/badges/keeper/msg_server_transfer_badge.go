package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TransferBadge(goCtx context.Context, msg *types.MsgTransferBadge) (*types.MsgTransferBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	rangesToValidate := []*types.IdRange{}
	for _, transfer := range msg.Transfers {
		for _, balance := range transfer.Balances {
			rangesToValidate = append(rangesToValidate, balance.BadgeIds...)
		}
	}

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                 msg.Creator,
		CollectionId:            msg.CollectionId,
		BadgeIdRangesToValidate: rangesToValidate,
	})
	if err != nil {
		return nil, err
	}

	if collection.BalancesType.LTE(sdk.NewUint(0)) {
		return nil, ErrOffChainBalances
	}

	

	for _, transfer := range msg.Transfers {
		fromBalanceKey := ConstructBalanceKey(transfer.From, msg.CollectionId)
		fromUserBalance, found := k.Keeper.GetUserBalanceFromStore(ctx, fromBalanceKey)
		if !found {
			return nil, ErrUserBalanceNotExists
		}
		
		for _, to := range transfer.ToAddresses {
			toBalanceKey := ConstructBalanceKey(to, msg.CollectionId)
			toUserBalance, found := k.Keeper.GetUserBalanceFromStore(ctx, toBalanceKey)
			if !found {
				toUserBalance = types.UserBalanceStore{
					Balances : []*types.Balance{},
					ApprovedTransfers: []*types.UserApprovedTransfer{},
					NextTransferTrackerId: sdk.NewUint(1),
					Permissions: &types.UserPermissions{
						CanUpdateApprovedTransfers: []*types.UserApprovedTransferPermission{
							{
								DefaultValues: &types.UserApprovedTransferDefaultValues{
									PermittedTimes: []*types.IdRange{},
									ForbiddenTimes: []*types.IdRange{},
								},
								Combinations: []*types.UserApprovedTransferCombination{{}},
							},
						},
					},
				}
			}

			for _, balance := range transfer.Balances {
				amount := balance.Amount
				fromUserBalance, toUserBalance, err = HandleTransfer(ctx, collection, balance.BadgeIds, fromUserBalance, toUserBalance, amount, transfer.From, to, msg.Creator)
				if err != nil {
					return nil, err
				}
			}
			
			//TODO: solutions

			if err := k.SetUserBalanceInStore(ctx, toBalanceKey, toUserBalance); err != nil {
				return nil, err
			}
		}

		if err := k.SetUserBalanceInStore(ctx, fromBalanceKey, fromUserBalance); err != nil {
			return nil, err
		}
	}

	

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute("collection_id", fmt.Sprint(msg.CollectionId)),
		),
	)

	return &types.MsgTransferBadgeResponse{}, nil
}
