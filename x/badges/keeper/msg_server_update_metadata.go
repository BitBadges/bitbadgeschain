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

	if msg.CollectionMetadataTimeline != nil && len(msg.CollectionMetadataTimeline) > 0 {
		if err := ValidateCollectionMetadataUpdate(ctx, collection.CollectionMetadataTimeline, msg.CollectionMetadataTimeline, collection.Permissions.CanUpdateCollectionMetadata); err != nil {
			return nil, err
		}
		collection.CollectionMetadataTimeline = msg.CollectionMetadataTimeline
	}

	if msg.OffChainBalancesMetadataTimeline != nil && len(msg.OffChainBalancesMetadataTimeline) > 0 {
		if err := ValidateOffChainBalancesMetadataUpdate(ctx, collection, collection.OffChainBalancesMetadataTimeline, msg.OffChainBalancesMetadataTimeline, collection.Permissions.CanUpdateOffChainBalancesMetadata); err != nil {
			return nil, err
		}
		collection.OffChainBalancesMetadataTimeline = msg.OffChainBalancesMetadataTimeline
	}

	if msg.ContractAddressTimeline != nil && len(msg.ContractAddressTimeline) > 0 {
		if err := ValidateContractAddressUpdate(ctx, collection.ContractAddressTimeline, msg.ContractAddressTimeline, collection.Permissions.CanUpdateContractAddress); err != nil {
			return nil, err
		}
		collection.ContractAddressTimeline = msg.ContractAddressTimeline
	}

	if msg.StandardsTimeline != nil && len(msg.StandardsTimeline) > 0 {
		if err := ValidateStandardsUpdate(ctx, collection.StandardsTimeline, msg.StandardsTimeline, collection.Permissions.CanUpdateStandards); err != nil {
			return nil, err
		}
		collection.StandardsTimeline = msg.StandardsTimeline
	}

	if msg.CustomDataTimeline != nil && len(msg.CustomDataTimeline) > 0 {
		if err := ValidateCustomDataUpdate(ctx, collection.CustomDataTimeline, msg.CustomDataTimeline, collection.Permissions.CanUpdateCustomData); err != nil {
			return nil, err
		}
		collection.CustomDataTimeline = msg.CustomDataTimeline
	}

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
