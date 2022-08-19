package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdateBytes(goCtx context.Context, msg *types.MsgUpdateBytes) (*types.MsgUpdateBytesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:        msg.Creator,
		BadgeId:        msg.BadgeId,
		MustBeManager:  true,
		CanUpdateBytes: true,
	})
	if err != nil {
		return nil, err
	}

	err = types.ValidateBytes(msg.NewBytes)
	if err != nil {
		return nil, err
	}

	badge.ArbitraryBytes = msg.NewBytes

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "UpdateArbitraryBytes"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
		),
	)

	return &types.MsgUpdateBytesResponse{}, nil
}
