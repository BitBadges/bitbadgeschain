package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteCollection(goCtx context.Context, msg *types.MsgDeleteCollection) (*types.MsgDeleteCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := k.GetCreator(ctx, msg.Creator, msg.CreatorOverride)
	if err != nil {
		return nil, err
	}
	msg.Creator = creator

	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	err = k.UniversalValidate(ctx, collection, UniversalValidationParams{
		Creator:       msg.Creator,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	//Check deleted permission is valid for current time
	err = k.CheckIfActionPermissionPermits(ctx, collection.CollectionPermissions.CanDeleteCollection, "can delete collection")
	if err != nil {
		return nil, err
	}

	k.DeleteCollectionFromStore(ctx, collection.CollectionId)
	if err != nil {
		return nil, err
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	//TODO: should we prune all balances and challenge stores here too?
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("msg_type", "delete_collection"),
			sdk.NewAttribute("msg", string(msgBytes)),
		),
	)

	return &types.MsgDeleteCollectionResponse{}, nil
}
