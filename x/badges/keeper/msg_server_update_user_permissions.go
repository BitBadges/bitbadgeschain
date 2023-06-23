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

	balanceKey := ConstructBalanceKey(msg.Creator, msg.CollectionId)
	userBalance, found := k.Keeper.GetUserBalanceFromStore(ctx, balanceKey)
	if !found {
		userBalance = types.UserBalanceStore{
			Permissions: types.UserPermissions{
				CanUpdateApprovedTransfers: {
					DefaultValues: types.DefaultValues{
						PermittedTimes: []types.IdRange{},
						ForbiddenTimes: []types.IdRange{},
					},
					Combinations: {},
				},
			},
		}
	}

	err = types.ValidateUserPermissionsUpdate(userBalance.Permissions, msg.Permissions, true)
	if err != nil {
		return nil, err
	}

	//iterate through the non-nil values
	if msg.Permissions.CanUpdateApprovedTransfers != nil {
		collection.Permissions.CanUpdateApprovedTransfers = msg.Permissions.CanUpdateApprovedTransfers
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
