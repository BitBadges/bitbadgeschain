package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdateUris(goCtx context.Context, msg *types.MsgUpdateUris) (*types.MsgUpdateUrisResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validationParams := UniversalValidationParams{
		Creator: msg.Creator,
		BadgeId: msg.BadgeId,
		MustBeManager: true,
		CanUpdateUris: true,
	}

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, validationParams)
	if err != nil {
		return nil, err
	}

	badge.Uri = msg.Uri
	badge.SubassetUriFormat = msg.SubassetUri

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
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
