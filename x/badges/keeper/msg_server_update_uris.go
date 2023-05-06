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

	newCollectionUri, newBadgeUris, newBalancesUri, needToValidateUpdateMetadataUris, needToValidateUpdateBalanceUri := GetUrisToStoreAndPermissionsToCheck(collection, msg.CollectionUri, msg.BadgeUris, msg.BalancesUri)
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

	collection.BadgeUris = newBadgeUris
	collection.CollectionUri = newCollectionUri
	collection.BalancesUri = newBalancesUri

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
