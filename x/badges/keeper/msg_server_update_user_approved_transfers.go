package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateUserApprovedTransfers(goCtx context.Context, msg *types.MsgUpdateUserApprovedTransfers) (*types.MsgUpdateUserApprovedTransfersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                              msg.Creator,
		CollectionId:                         msg.CollectionId,
		MustBeManager:                        true,
	})
	if err != nil {
		return nil, err
	}

	if collection.BalancesType != sdk.NewUint(0) {
		return nil, ErrOffChainBalances
	}

	balanceKey := ConstructBalanceKey(msg.Creator, collection.CollectionId)
	userBalance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
	if !found {
		userBalance = types.UserBalanceStore{
			Balances : []*types.Balance{},
			ApprovedTransfersTimeline: []*types.UserApprovedTransferTimeline{},
			NextTransferTrackerId: sdk.NewUint(1),
			Permissions: &types.UserPermissions{
				CanUpdateApprovedTransfers: []*types.UserApprovedTransferPermission{},
			},
		}
	}

	if err := ValidateUserApprovedTransfersUpdate(ctx, userBalance.ApprovedTransfersTimeline, msg.ApprovedTransfersTimeline, userBalance.Permissions.CanUpdateApprovedTransfers); err != nil {
		return nil, err
	}
	userBalance.ApprovedTransfersTimeline = msg.ApprovedTransfersTimeline

	err = k.SetUserBalanceInStore(ctx, balanceKey, userBalance)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgUpdateUserApprovedTransfersResponse{}, nil
}
