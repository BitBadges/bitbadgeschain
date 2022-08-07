package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdatePermissions(goCtx context.Context, msg *types.MsgUpdatePermissions) (*types.MsgUpdatePermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)
	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)

	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")
	if !found {
		return nil, ErrBadgeNotExists
	}

	if badge.Manager != CreatorAccountNum {
		return nil, ErrSenderIsNotManager
	}

	ctx.GasMeter().ConsumeGas(BadgeUpdate, "badge update")

	err := types.ValidatePermissionsUpdate(badge.PermissionFlags, msg.Permissions)
	if err != nil {
		return nil, err
	}

	badge.PermissionFlags = msg.Permissions

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "UpdatePermissions"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("NewPermissions", fmt.Sprint(msg.Permissions)),
		),
	)

	return &types.MsgUpdatePermissionsResponse{}, nil
}
