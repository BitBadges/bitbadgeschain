package keeper

import (
	"context"
	"encoding/json"

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

	collection, err = k.CreateBadges(ctx, collection, msg.BadgeSupplys, msg.Transfers, msg.Claims, msg.Creator)
	if err != nil {
		return nil, err
	}

	if (msg.CollectionUri != "" && msg.CollectionUri != collection.CollectionUri) || (len(msg.BadgeUris) > 0) {
		_, _, err = k.UniversalValidate(ctx, UniversalValidationParams{
			Creator:       msg.Creator,
			CollectionId:  msg.CollectionId,
			MustBeManager: true,
			CanUpdateUris: true,
		})
		if err != nil {
			return nil, err
		}

		//Already validated in ValidateBasic
		collection.BadgeUris = msg.BadgeUris
		collection.CollectionUri = msg.CollectionUri
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	collectionJson, err := json.Marshal(collection)
	if err != nil {
		return nil, err
	}

	transfersJson, err := json.Marshal(msg.Transfers)
	if err != nil {
		return nil, err
	}
	
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("collection", string(collectionJson)),
			sdk.NewAttribute("transfers", string(transfersJson)),
		),
	)

	return &types.MsgMintBadgeResponse{
		NextBadgeId: collection.NextBadgeId,
	}, nil
}
