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

	if msg.ApprovedTransfersTimeline != nil && len(msg.ApprovedTransfersTimeline) > 0 {
		if collection.BalancesType != sdk.NewUint(0) {
			return nil, ErrOffChainBalances
		}
		
		if err := ValidateCollectionApprovedTransfersUpdate(ctx, collection.ApprovedTransfersTimeline, msg.ApprovedTransfersTimeline, collection.Permissions.CanUpdateApprovedTransfers); err != nil {
			return nil, err
		}

		collection.ApprovedTransfersTimeline = msg.ApprovedTransfersTimeline
	}
	
	if msg.CollectionMetadataTimeline != nil && len(msg.CollectionMetadataTimeline) > 0 {
		if err := ValidateCollectionMetadataUpdate(ctx, collection.CollectionMetadataTimeline, msg.CollectionMetadataTimeline, collection.Permissions.CanUpdateCollectionMetadata); err != nil {
			return nil, err
		}

		collection.CollectionMetadataTimeline = msg.CollectionMetadataTimeline
	}

	if msg.OffChainBalancesMetadataTimeline != nil && len(msg.OffChainBalancesMetadataTimeline) > 0 {
		if collection.BalancesType != sdk.NewUint(1) {
			return nil, ErrOffChainBalances
		}
		
		if err := ValidateOffChainBalancesMetadataUpdate(ctx, collection.OffChainBalancesMetadataTimeline, msg.OffChainBalancesMetadataTimeline, collection.Permissions.CanUpdateOffChainBalancesMetadata); err != nil {
			return nil, err
		}

		collection.OffChainBalancesMetadataTimeline = msg.OffChainBalancesMetadataTimeline
	}

	if msg.BadgesToCreate != nil && len(msg.BadgesToCreate) > 0 {
		if collection.BalancesType == sdk.NewUint(1) {
			return nil, ErrOffChainBalances
		}

		collection, err = k.CreateBadges(ctx, collection, msg.BadgesToCreate, msg.Transfers, msg.Creator)
		if err != nil {
			return nil, err
		}
	}

	if msg.InheritedBalancesTimeline != nil && len(msg.InheritedBalancesTimeline) > 0 {
		if collection.BalancesType != sdk.NewUint(2) {
			return nil, ErrOffChainBalances
		}
		
		if err := ValidateInheritedBalancesUpdate(ctx, collection.InheritedBalancesTimeline, msg.InheritedBalancesTimeline, collection.Permissions.CanUpdateInheritedBalances); err != nil {
			return nil, err
		}

		maxBadgeId := collection.NextBadgeId.Sub(sdk.NewUint(1))

		for _, timelineVal := range msg.InheritedBalancesTimeline {
			for _, inheritedBalance := range timelineVal.InheritedBalances {
				for _, badgeId := range inheritedBalance.BadgeIds {
					if badgeId.End.GT(maxBadgeId) {
						return nil, ErrBadgeIdTooHigh
					}
				}
			}
		}
	

		collection.InheritedBalancesTimeline = msg.InheritedBalancesTimeline
	}

	

	if msg.BadgeMetadataTimeline != nil && len(msg.BadgeMetadataTimeline) > 0 {
		if err := ValidateBadgeMetadataUpdate(ctx, collection.BadgeMetadataTimeline, msg.BadgeMetadataTimeline, collection.Permissions.CanUpdateBadgeMetadata); err != nil {
			return nil, err
		}
		
		collection.BadgeMetadataTimeline = msg.BadgeMetadataTimeline
	}

	
	//TODO: address mappings as well

	
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
