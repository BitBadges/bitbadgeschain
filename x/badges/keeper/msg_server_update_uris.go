package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdateUris(goCtx context.Context, msg *types.MsgUpdateUris) (*types.MsgUpdateUrisResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		BadgeId:       msg.BadgeId,
		MustBeManager: true,
		CanUpdateUris: true,
	})
	if err != nil {
		return nil, err
	}

	badge.Uri = msg.Uri

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "UpdateURIs"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
		),
	)

	return &types.MsgUpdateUrisResponse{}, nil
}
