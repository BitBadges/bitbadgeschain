package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdateUris(goCtx context.Context, msg *types.MsgUpdateUris) (*types.MsgUpdateUrisResponse, error) {
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

	permissions := types.GetPermissions(badge.PermissionFlags)
	if !permissions.CanUpdateUris() {
		return nil, ErrInvalidPermissions
	}

	badge.Uri = msg.Uri
	badge.SubassetUriFormat = msg.SubassetUri

	if err := k.UpdateBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "UpdatePermissions"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("NewUri", fmt.Sprint(msg.Uri)),
			sdk.NewAttribute("NewSubbadgeUri", fmt.Sprint(msg.SubassetUri)),
		),
	)

	return &types.MsgUpdateUrisResponse{}, nil
}
