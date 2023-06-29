package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) UpdateUserPermissions(goCtx context.Context,  msg *types.MsgUpdateUserPermissions) (*types.MsgUpdateUserPermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

  collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		CollectionId:  msg.CollectionId,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	for _, addressMapping := range msg.AddressMappings {
		if err := k.CreateAddressMapping(ctx, addressMapping); err != nil {
			return nil, err
		}
	}

	balanceKey := ConstructBalanceKey(msg.Creator, msg.CollectionId)
	userBalance, found := k.Keeper.GetUserBalanceFromStore(ctx, balanceKey)
	if !found {
		userBalance = &types.UserBalanceStore{
			Balances : []*types.Balance{},
			ApprovedOutgoingTransfersTimeline: collection.DefaultUserApprovedOutgoingTransfersTimeline,
			ApprovedIncomingTransfersTimeline: collection.DefaultUserApprovedIncomingTransfersTimeline,
			Permissions: &types.UserPermissions{
				CanUpdateApprovedIncomingTransfers: []*types.UserApprovedTransferPermission{},
				CanUpdateApprovedOutgoingTransfers: []*types.UserApprovedTransferPermission{},
			},
		}
	}

	err = types.ValidateUserPermissionsUpdate(userBalance.Permissions, msg.Permissions, true)
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

	if err := k.SetUserBalanceInStore(ctx, balanceKey, userBalance); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgUpdateUserPermissionsResponse{}, nil
}
