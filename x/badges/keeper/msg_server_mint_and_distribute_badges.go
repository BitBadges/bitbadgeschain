package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) MintAndDistributeBadges(goCtx context.Context, msg *types.MsgMintAndDistributeBadges) (*types.MsgMintAndDistributeBadgesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:             msg.Creator,
		CollectionId:        msg.CollectionId,
		MustBeManager:       true,
	})
	if err != nil {
		return nil, err
	}

	for _, addressMapping := range msg.AddressMappings {
		if err := k.CreateAddressMapping(ctx, addressMapping); err != nil {
			return nil, err
		}
	}

	collection, err = k.CreateBadges(ctx, collection, msg.BadgesToCreate, msg.Transfers, msg.Creator)
	if err != nil {
		return nil, err
	}


	if err := k.ValidateCollectionApprovedTransfersUpdate(ctx, collection, collection.CollectionApprovedTransfersTimeline, msg.CollectionApprovedTransfersTimeline, collection.Permissions.CanUpdateCollectionApprovedTransfers, msg.Creator); err != nil {
		return nil, err
	}
	collection.CollectionApprovedTransfersTimeline = msg.CollectionApprovedTransfersTimeline

	if err := k.ValidateCollectionMetadataUpdate(ctx, collection.CollectionMetadataTimeline, msg.CollectionMetadataTimeline, collection.Permissions.CanUpdateCollectionMetadata); err != nil {
		return nil, err
	}
	collection.CollectionMetadataTimeline = msg.CollectionMetadataTimeline


	if err := k.ValidateOffChainBalancesMetadataUpdate(ctx, collection, collection.OffChainBalancesMetadataTimeline, msg.OffChainBalancesMetadataTimeline, collection.Permissions.CanUpdateOffChainBalancesMetadata); err != nil {
		return nil, err
	}
	collection.OffChainBalancesMetadataTimeline = msg.OffChainBalancesMetadataTimeline


	if err := k.ValidateInheritedBalancesUpdate(ctx, collection, collection.InheritedBalancesTimeline, msg.InheritedBalancesTimeline, collection.Permissions.CanUpdateInheritedBalances); err != nil {
		return nil, err
	}
	collection.InheritedBalancesTimeline = msg.InheritedBalancesTimeline


	if err := k.ValidateBadgeMetadataUpdate(ctx, collection.BadgeMetadataTimeline, msg.BadgeMetadataTimeline, collection.Permissions.CanUpdateBadgeMetadata); err != nil {
		return nil, err
	}
	collection.BadgeMetadataTimeline = msg.BadgeMetadataTimeline
	

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgMintAndDistributeBadgesResponse{}, nil
}
