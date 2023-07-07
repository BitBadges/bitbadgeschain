package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateUserApprovedTransfers(goCtx context.Context, msg *types.MsgUpdateUserApprovedTransfers) (*types.MsgUpdateUserApprovedTransfersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	if !IsStandardBalances(collection) {
		return nil, ErrWrongBalancesType
	}

	balanceKey := ConstructBalanceKey(msg.Creator, collection.CollectionId)
	userBalance, found := k.GetUserBalanceFromStore(ctx, balanceKey)
	if !found {
		userBalance = &types.UserBalanceStore{
			Balances:                          []*types.Balance{},
			ApprovedOutgoingTransfersTimeline: collection.DefaultUserApprovedOutgoingTransfersTimeline,
			ApprovedIncomingTransfersTimeline: collection.DefaultUserApprovedIncomingTransfersTimeline,
			Permissions: &types.UserPermissions{
				CanUpdateApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransferPermission{},
				CanUpdateApprovedIncomingTransfers: []*types.UserApprovedIncomingTransferPermission{},
			},
		}
	}

	manager := types.GetCurrentManager(ctx, collection)

	if msg.UpdateApprovedOutgoingTransfersTimeline {
		if err := k.ValidateUserApprovedOutgoingTransfersUpdate(ctx, userBalance.ApprovedOutgoingTransfersTimeline, msg.ApprovedOutgoingTransfersTimeline, userBalance.Permissions.CanUpdateApprovedOutgoingTransfers, manager); err != nil {
			return nil, err
		}
		userBalance.ApprovedOutgoingTransfersTimeline = msg.ApprovedOutgoingTransfersTimeline
	}

	if msg.UpdateApprovedIncomingTransfersTimeline {
		if err := k.ValidateUserApprovedIncomingTransfersUpdate(ctx, userBalance.ApprovedIncomingTransfersTimeline, msg.ApprovedIncomingTransfersTimeline, userBalance.Permissions.CanUpdateApprovedIncomingTransfers, manager); err != nil {
			return nil, err
		}
		userBalance.ApprovedIncomingTransfersTimeline = msg.ApprovedIncomingTransfersTimeline
	}

	if msg.UpdateApprovedTransfersUserPermissions {
		err := k.ValidateUserPermissionsUpdate(ctx, userBalance.Permissions, msg.Permissions, true, manager)
		if err != nil {
			return nil, err
		}

		//iterate through the non-nil values
		if msg.Permissions.CanUpdateApprovedIncomingTransfers != nil {
			userBalance.Permissions.CanUpdateApprovedIncomingTransfers = msg.Permissions.CanUpdateApprovedIncomingTransfers
		}

		if msg.Permissions.CanUpdateApprovedOutgoingTransfers != nil {
			userBalance.Permissions.CanUpdateApprovedOutgoingTransfers = msg.Permissions.CanUpdateApprovedOutgoingTransfers
		}

		userBalance.Permissions = msg.Permissions
	}



	err := k.SetUserBalanceInStore(ctx, balanceKey, userBalance)
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
