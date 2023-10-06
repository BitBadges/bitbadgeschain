package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateUserApprovals(goCtx context.Context, msg *types.MsgUpdateUserApprovals) (*types.MsgUpdateUserApprovalsResponse, error) {
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
			OutgoingApprovals: collection.DefaultUserOutgoingApprovals,
			IncomingApprovals: collection.DefaultUserIncomingApprovals,
			UserPermissions:                   collection.DefaultUserPermissions,
		}
	}

	if userBalance.UserPermissions == nil {
		userBalance.UserPermissions = &types.UserPermissions{}
	}

	manager := types.GetCurrentManager(ctx, collection)

	if msg.UpdateOutgoingApprovals {
		if err := k.ValidateUserOutgoingApprovalsUpdate(ctx, collection, userBalance.OutgoingApprovals, msg.OutgoingApprovals, userBalance.UserPermissions.CanUpdateOutgoingApprovals, manager, msg.Creator); err != nil {
			return nil, err
		}
		userBalance.OutgoingApprovals = msg.OutgoingApprovals
	}

	if msg.UpdateIncomingApprovals {
		if err := k.ValidateUserIncomingApprovalsUpdate(ctx, collection,  userBalance.IncomingApprovals, msg.IncomingApprovals, userBalance.UserPermissions.CanUpdateIncomingApprovals, manager, msg.Creator); err != nil {
			return nil, err
		}
		userBalance.IncomingApprovals = msg.IncomingApprovals
	}

	if msg.UpdateUserPermissions {
		err := k.ValidateUserPermissionsUpdate(ctx, userBalance.UserPermissions, msg.UserPermissions, manager)
		if err != nil {
			return nil, err
		}

		//iterate through the non-nil values
		if msg.UserPermissions.CanUpdateIncomingApprovals != nil {
			userBalance.UserPermissions.CanUpdateIncomingApprovals = msg.UserPermissions.CanUpdateIncomingApprovals
		}

		if msg.UserPermissions.CanUpdateOutgoingApprovals != nil {
			userBalance.UserPermissions.CanUpdateOutgoingApprovals = msg.UserPermissions.CanUpdateOutgoingApprovals
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

	return &types.MsgUpdateUserApprovalsResponse{}, nil
}
