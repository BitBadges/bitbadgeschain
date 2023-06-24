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

	// newCollectionMetadata, newBadgeMetadata, newOffChainBalancesMetadata, needToValidateUpdateCollectionMetadata, needToValidateUpdateBadgeMetadata, needToValidateUpdateBalanceUri := GetUrisToStoreAndPermissionsToCheck(collection, msg.CollectionMetadata, msg.BadgeMetadata, msg.OffChainBalancesMetadata)

	// _, err = k.UniversalValidate(ctx, UniversalValidationParams{
	// 	Creator:                     msg.Creator,
	// 	CollectionId:                msg.CollectionId,
	// 	MustBeManager:               true,
	// 	CanUpdateOffChainBalancesMetadata:   needToValidateUpdateBalanceUri,
	// 	CanUpdateBadgeMetadata:      needToValidateUpdateBadgeMetadata,
	// 	CanUpdateCollectionMetadata: needToValidateUpdateCollectionMetadata,
	// 	CanUpdateContractAddress:    msg.ContractAddress != "" && collection.ContractAddress != msg.ContractAddress,
	// 	CanUpdateCustomData:         msg.CustomData != "" && collection.CustomData != msg.CustomData,
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// //Check badge metadata for isFrozen logic
	// err = AssertIsFrozenLogicIsMaintained(collection.BadgeMetadata, newBadgeMetadata)
	// if err != nil {
	// 	return nil, err
	// }

	// collection.BadgeMetadata = newBadgeMetadata
	// collection.CollectionMetadata = newCollectionMetadata
	// collection.OffChainBalancesMetadata = newOffChainBalancesMetadata

	// if msg.ContractAddress != "" {
	// 	collection.ContractAddress = msg.ContractAddress
	// }

	// if msg.CustomData != "" {
	// 	collection.CustomData = msg.CustomData
	// }

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
