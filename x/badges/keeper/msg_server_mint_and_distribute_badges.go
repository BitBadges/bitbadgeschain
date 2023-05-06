package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) MintAndDistributeBadges(goCtx context.Context, msg *types.MsgMintAndDistributeBadges) (*types.MsgMintAndDistributeBadgesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	checkIfCanCreateMoreBadges := msg.BadgeSupplys != nil && len(msg.BadgeSupplys) > 0
	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
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


	//Already validated in ValidateBasic
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

	return &types.MsgMintAndDistributeBadgesResponse{
		NextBadgeId: collection.NextBadgeId,
	}, nil
}
