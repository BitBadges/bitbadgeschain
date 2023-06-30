package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) ArchiveCollection(goCtx context.Context, msg *types.MsgArchiveCollection) (*types.MsgArchiveCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                 msg.Creator,
		CollectionId:            msg.CollectionId,
		MustBeManager:           true,
		OverrideArchive:  			 true,
	})
	if err != nil {
		return nil, err
	}

	if err := ValidateIsArchivedUpdate(ctx, collection.IsArchivedTimeline, msg.IsArchivedTimeline, collection.Permissions.CanArchive); err != nil {
		return nil, err
	}
	collection.IsArchivedTimeline = msg.IsArchivedTimeline


	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgArchiveCollectionResponse{}, nil
}
