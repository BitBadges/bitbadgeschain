package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) MintAndDistributeBadges(goCtx context.Context, msg *types.MsgMintAndDistributeBadges) (*types.MsgMintAndDistributeBadgesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	checkIfCanCreateMoreBadges := msg.BadgesToCreate != nil && len(msg.BadgesToCreate) > 0
	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:             msg.Creator,
		CollectionId:        msg.CollectionId,
		MustBeManager:       true,
		CanCreateMoreBadges: checkIfCanCreateMoreBadges,
	})
	if err != nil {
		return nil, err
	}

	if collection.IsOffChainBalances {
		return nil, ErrOffChainBalances
	}

	newCollectionMetadata, newBadgeMetadata, newOffChainBalancesMetadata, needToValidateUpdateCollectionMetadata, needToValidateUpdateBadgeMetadata, needToValidateUpdateBalanceUri := GetUrisToStoreAndPermissionsToCheck(collection, msg.CollectionMetadata, msg.BadgeMetadata, msg.OffChainBalancesMetadata)
	newApprovedTransfers, needToValidateUpdateCollectionApprovedTransfers := GetApprovedTransfersToStore(collection, msg.ApprovedTransfers)

	_, err = k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:                              msg.Creator,
		CollectionId:                         msg.CollectionId,
		MustBeManager:                        true,
		CanUpdateOffChainBalancesMetadata:            needToValidateUpdateBalanceUri,
		CanUpdateBadgeMetadata:               needToValidateUpdateBadgeMetadata,
		CanUpdateCollectionMetadata:          needToValidateUpdateCollectionMetadata,
		CanUpdateCollectionApprovedTransfers: needToValidateUpdateCollectionApprovedTransfers,
	})
	if err != nil {
		return nil, err
	}

	//Check badge metadata for isFrozen logic
	err = AssertIsFrozenLogicIsMaintained(collection.BadgeMetadata, newBadgeMetadata)
	if err != nil {
		return nil, err
	}

	err = AssertIsFrozenLogicForApprovedTransfers(collection.ApprovedTransfers, newApprovedTransfers)
	if err != nil {
		return nil, err
	}

	collection.BadgeMetadata = newBadgeMetadata
	collection.CollectionMetadata = newCollectionMetadata
	collection.OffChainBalancesMetadata = newOffChainBalancesMetadata
	collection.ApprovedTransfers = newApprovedTransfers

	collection, err = k.CreateBadges(ctx, collection, msg.BadgesToCreate, msg.Transfers, msg.Claims, msg.Creator)
	if err != nil {
		return nil, err
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

	return &types.MsgMintAndDistributeBadgesResponse{
		NextBadgeId: collection.NextBadgeId,
	}, nil
}
