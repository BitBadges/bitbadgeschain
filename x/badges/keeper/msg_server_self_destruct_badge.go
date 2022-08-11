package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) SelfDestructBadge(goCtx context.Context, msg *types.MsgSelfDestructBadge) (*types.MsgSelfDestructBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, _, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		BadgeId:       msg.BadgeId,
		MustBeManager: true,
		CanRevoke:     true,
	})
	if err != nil {
		return nil, err
	}

	k.DeleteBadgeFromStore(ctx, msg.BadgeId)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "SelfDestructBadge"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
		),
	)

	return &types.MsgSelfDestructBadgeResponse{}, nil
}
