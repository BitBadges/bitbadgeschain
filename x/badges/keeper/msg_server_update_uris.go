package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateUris(goCtx context.Context, msg *types.MsgUpdateUris) (*types.MsgUpdateUrisResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		CollectionId:  msg.CollectionId,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	needToValidateUpdateMetadataUris := false
	needToValidateUpdateBalanceUri := false
	newCollectionUri := collection.CollectionUri
	newBadgeUris := collection.BadgeUris
	newBalanceUri := collection.BalancesUri

	if msg.CollectionUri != "" && msg.CollectionUri != collection.CollectionUri {
		needToValidateUpdateMetadataUris = true
		newCollectionUri = msg.CollectionUri
	}

	if msg.BalancesUri != "" && msg.BalancesUri != collection.BalancesUri {
		needToValidateUpdateBalanceUri = true
		newBalanceUri = msg.BalancesUri
	}

	if msg.BadgeUris != nil && len(msg.BadgeUris) > 0 {
		newBadgeUris = msg.BadgeUris

		for idx, badgeUri := range collection.BadgeUris {
			if msg.BadgeUris[idx].Uri != badgeUri.Uri {
				needToValidateUpdateMetadataUris = true
				break
			}

			if len(msg.BadgeUris[idx].BadgeIds) != len(badgeUri.BadgeIds) {
				needToValidateUpdateMetadataUris = true
				break
			}

			for j, badgeIdRange := range badgeUri.BadgeIds {
				if badgeIdRange.Start != msg.BadgeUris[idx].BadgeIds[j].Start || badgeIdRange.End != msg.BadgeUris[idx].BadgeIds[j].End {
					needToValidateUpdateMetadataUris = true
					break
				}
			}
		}
	}

	_, err = k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		CollectionId:  msg.CollectionId,
		MustBeManager: true,
		CanUpdateMetadataUris: needToValidateUpdateMetadataUris,
		CanUpdateBalancesUri: needToValidateUpdateBalanceUri,
	})
	if err != nil {
		return nil, err
	}

	

	//Already validated in ValidateBasic
	collection.BadgeUris = newBadgeUris
	collection.CollectionUri = newCollectionUri
	collection.BalancesUri = newBalanceUri

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)
	return &types.MsgUpdateUrisResponse{}, nil
}
