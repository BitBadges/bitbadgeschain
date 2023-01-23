package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateUris(goCtx context.Context, msg *types.MsgUpdateUris) (*types.MsgUpdateUrisResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		CollectionId:       msg.CollectionId,
		MustBeManager: true,
		CanUpdateUris: true,
	})
	if err != nil {
		return nil, err
	}

	//Already validated in ValidateBasic
	badge.BadgeUri = msg.BadgeUri
	badge.CollectionUri = msg.CollectionUri


	if err := k.SetCollectionInStore(ctx, badge); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "UpdateURIs"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.CollectionId)),
		),
	)

	return &types.MsgUpdateUrisResponse{}, nil
}
