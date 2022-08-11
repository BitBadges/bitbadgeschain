package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdatePermissions(goCtx context.Context, msg *types.MsgUpdatePermissions) (*types.MsgUpdatePermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	err = types.ValidatePermissionsUpdate(badge.Permissions, msg.Permissions)
	if err != nil {
		return nil, err
	}

	badge.Permissions = msg.Permissions

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "UpdatePermissions"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
		),
	)

	return &types.MsgUpdatePermissionsResponse{}, nil
}
