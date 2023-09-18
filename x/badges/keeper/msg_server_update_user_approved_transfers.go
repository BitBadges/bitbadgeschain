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
			ApprovedOutgoingTransfers: collection.DefaultUserApprovedOutgoingTransfers,
			ApprovedIncomingTransfers: collection.DefaultUserApprovedIncomingTransfers,
			UserPermissions:                   collection.DefaultUserPermissions,
		}
	}

	if userBalance.UserPermissions == nil {
		userBalance.UserPermissions = &types.UserPermissions{}
	}

	manager := types.GetCurrentManager(ctx, collection)

	if msg.UpdateApprovedOutgoingTransfers {
		if err := k.ValidateUserApprovedOutgoingTransfersUpdate(ctx, userBalance.ApprovedOutgoingTransfers, msg.ApprovedOutgoingTransfers, userBalance.UserPermissions.CanUpdateApprovedOutgoingTransfers, manager); err != nil {
			return nil, err
		}
		userBalance.ApprovedOutgoingTransfers = msg.ApprovedOutgoingTransfers
	}

	if msg.UpdateApprovedIncomingTransfers {
		if err := k.ValidateUserApprovedIncomingTransfersUpdate(ctx, userBalance.ApprovedIncomingTransfers, msg.ApprovedIncomingTransfers, userBalance.UserPermissions.CanUpdateApprovedIncomingTransfers, manager); err != nil {
			return nil, err
		}
		userBalance.ApprovedIncomingTransfers = msg.ApprovedIncomingTransfers
	}

	if msg.UpdateUserPermissions {
		err := k.ValidateUserPermissionsUpdate(ctx, userBalance.UserPermissions, msg.UserPermissions, manager)
		if err != nil {
			return nil, err
		}

		//iterate through the non-nil values
		if msg.UserPermissions.CanUpdateApprovedIncomingTransfers != nil {
			userBalance.UserPermissions.CanUpdateApprovedIncomingTransfers = msg.UserPermissions.CanUpdateApprovedIncomingTransfers
		}

		if msg.UserPermissions.CanUpdateApprovedOutgoingTransfers != nil {
			userBalance.UserPermissions.CanUpdateApprovedOutgoingTransfers = msg.UserPermissions.CanUpdateApprovedOutgoingTransfers
		}

		userBalance.UserPermissions = msg.UserPermissions
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
