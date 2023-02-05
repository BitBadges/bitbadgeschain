package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) MintBadge(goCtx context.Context, msg *types.MsgMintBadge) (*types.MsgMintBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	checkIfCanCreateMoreBadges := msg.BadgeSupplys != nil && len(msg.BadgeSupplys) > 0


	_, collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:             msg.Creator,
		CollectionId:        msg.CollectionId,
		MustBeManager:       true,
		CanCreateMoreBadges: checkIfCanCreateMoreBadges,
	})
	if err != nil {
		return nil, err
	}

	originalSubassetId := collection.NextBadgeId

	collection, err = k.CreateBadges(ctx, collection, msg.BadgeSupplys, msg.Transfers, msg.Claims, msg.Creator)
	if err != nil {
		return nil, err
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "CreatedBadges"),
			sdk.NewAttribute("BadgeId", fmt.Sprint(collection.CollectionId)),
			sdk.NewAttribute("FirstId", fmt.Sprint(originalSubassetId)),
			sdk.NewAttribute("LastId", fmt.Sprint(collection.NextBadgeId-1)),
		),
	)

	return &types.MsgMintBadgeResponse{
		NextBadgeId: collection.NextBadgeId,
	}, nil
}
