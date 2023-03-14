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

	needToValidateUpdateUris := false
	newCollectionUri := collection.CollectionUri
	newBadgeUris := collection.BadgeUris

	if msg.CollectionUri != "" && msg.CollectionUri != collection.CollectionUri {
		needToValidateUpdateUris = true
		newCollectionUri = msg.CollectionUri
	}
	if len(msg.BadgeUris) > 0 {
		newBadgeUris = msg.BadgeUris

		for idx, badgeUri := range collection.BadgeUris {
			if msg.BadgeUris[idx].Uri != badgeUri.Uri {
				needToValidateUpdateUris = true
				break
			}
			if len(msg.BadgeUris[idx].BadgeIds) != len(badgeUri.BadgeIds) {
				needToValidateUpdateUris = true
				break
			}

			for j, badgeIdRange := range badgeUri.BadgeIds {
				if badgeIdRange.Start != msg.BadgeUris[idx].BadgeIds[j].Start || badgeIdRange.End != msg.BadgeUris[idx].BadgeIds[j].End {
					needToValidateUpdateUris = true
					break
				}
			}
		}
	}


	_, _, err = k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		CollectionId:  msg.CollectionId,
		MustBeManager: true,
		CanUpdateUris: needToValidateUpdateUris,
	})
	if err != nil {
		return nil, err
	}

	//Already validated in ValidateBasic
	collection.BadgeUris = newBadgeUris
	collection.CollectionUri = newCollectionUri

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
