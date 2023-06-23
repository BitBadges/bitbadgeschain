package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NewCollection(goCtx context.Context, msg *types.MsgNewCollection) (*types.MsgNewCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	NextCollectionId := k.GetNextCollectionId(ctx)
	k.IncrementNextCollectionId(ctx)

	collection := types.BadgeCollection{
		CollectionId:       NextCollectionId,
		CollectionMetadata: msg.CollectionMetadata,
		OffChainBalancesMetadata:   msg.OffChainBalancesMetadata,
		BadgeMetadata:      msg.BadgeMetadata,
		Manager:            msg.Creator,
		Permissions:        msg.Permissions,
		ApprovedTransfers:  msg.ApprovedTransfers,
		CustomData:         msg.CustomData,
		ContractAddress:    msg.ContractAddress,
		Standard:           msg.Standard,
		NextBadgeId:        sdk.NewUint(1),
		NextClaimId:        sdk.NewUint(1),
		ParentCollectionId: sdk.NewUint(0),
		IsOffChainBalances: msg.OffChainBalancesMetadata.Uri != "" || msg.OffChainBalancesMetadata.CustomData != "",
		IsArchived:         false,
		UnmintedSupplys:    []*types.Balance{},
		MaxSupplys:         []*types.Balance{},
	}

	//Check badge metadata for isFrozen logic
	err := AssertIsFrozenLogicIsMaintained([]*types.BadgeMetadata{}, collection.BadgeMetadata)
	if err != nil {
		return nil, err
	}

	err = AssertIsFrozenLogicForApprovedTransfers([]*types.CollectionApprovedTransfer{}, collection.ApprovedTransfers)
	if err != nil {
		return nil, err
	}

	if len(msg.BadgesToCreate) != 0 {
		err := *new(error)
		collection, err = k.CreateBadges(ctx, collection, msg.BadgesToCreate, msg.Transfers, msg.Claims, msg.Creator)
		if err != nil {
			return nil, err
		}
	}

	//If we set the permissions to be frozen permanently, we can safely set all the individual isFrozen flags to true
	if collection.Permissions.CanUpdateCollectionApprovedTransfers.IsFrozen && len(collection.Permissions.CanUpdateCollectionApprovedTransfers.TimeIntervals) == 0 {
		for _, allowedTransfer := range collection.ApprovedTransfers {
			allowedTransfer.IsFrozen = true
		}
	}

	if collection.Permissions.CanUpdateBadgeMetadata.IsFrozen && len(collection.Permissions.CanUpdateBadgeMetadata.TimeIntervals) == 0 {
		for _, badgeMetadata := range collection.BadgeMetadata {
			badgeMetadata.IsFrozen = true
		}
	}

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
			sdk.NewAttribute("collectionId", fmt.Sprint(NextCollectionId)),
		),
	)

	return &types.MsgNewCollectionResponse{
		CollectionId: NextCollectionId,
	}, nil
}
