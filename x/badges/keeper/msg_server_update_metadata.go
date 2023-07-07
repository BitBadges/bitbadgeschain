package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateMetadata(goCtx context.Context, msg *types.MsgUpdateMetadata) (*types.MsgUpdateMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		CollectionId:  msg.CollectionId,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	if err := k.ValidateBadgeMetadataUpdate(ctx, collection.BadgeMetadataTimeline, msg.BadgeMetadataTimeline, collection.Permissions.CanUpdateBadgeMetadata); err != nil {
		return nil, err
	}
	collection.BadgeMetadataTimeline = msg.BadgeMetadataTimeline
	

	if err := k.ValidateCollectionMetadataUpdate(ctx, collection.CollectionMetadataTimeline, msg.CollectionMetadataTimeline, collection.Permissions.CanUpdateCollectionMetadata); err != nil {
		return nil, err
	}
	collection.CollectionMetadataTimeline = msg.CollectionMetadataTimeline

	if err := k.ValidateOffChainBalancesMetadataUpdate(ctx, collection, collection.OffChainBalancesMetadataTimeline, msg.OffChainBalancesMetadataTimeline, collection.Permissions.CanUpdateOffChainBalancesMetadata); err != nil {
		return nil, err
	}
	collection.OffChainBalancesMetadataTimeline = msg.OffChainBalancesMetadataTimeline

	if err := k.ValidateContractAddressUpdate(ctx, collection.ContractAddressTimeline, msg.ContractAddressTimeline, collection.Permissions.CanUpdateContractAddress); err != nil {
		return nil, err
	}
	collection.ContractAddressTimeline = msg.ContractAddressTimeline

	if err := k.ValidateStandardsUpdate(ctx, collection.StandardsTimeline, msg.StandardsTimeline, collection.Permissions.CanUpdateStandards); err != nil {
		return nil, err
	}
	collection.StandardsTimeline = msg.StandardsTimeline

	if err := k.ValidateCustomDataUpdate(ctx, collection.CustomDataTimeline, msg.CustomDataTimeline, collection.Permissions.CanUpdateCustomData); err != nil {
		return nil, err
	}
	collection.CustomDataTimeline = msg.CustomDataTimeline

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)
	return &types.MsgUpdateMetadataResponse{}, nil
}
