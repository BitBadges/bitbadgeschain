package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) SelfDestructBadge(goCtx context.Context, msg *types.MsgSelfDestructBadge) (*types.MsgSelfDestructBadgeResponse, error) {
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

	permissions := types.GetPermissions(badge.PermissionFlags)
	//If manager has permissions to revoke, he can theoretically self destruct the badge by revoking all supply of everyone
	if !permissions.CanRevoke() {
		return nil, ErrBadgeCanNotBeSelfDestructed
	}

	k.DeleteBadgeFromStore(ctx, msg.BadgeId)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "SelfDestructBadge"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
		),
	)

	return &types.MsgSelfDestructBadgeResponse{}, nil
}
