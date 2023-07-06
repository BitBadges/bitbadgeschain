package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) UpdateCollectionPermissions(goCtx context.Context, msg *types.MsgUpdateCollectionPermissions) (*types.MsgUpdateCollectionPermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		CollectionId:  msg.CollectionId,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	for _, addressMapping := range msg.AddressMappings {
		if err := k.CreateAddressMapping(ctx, addressMapping); err != nil {
			return nil, err
		}
	}
	
	err = k.ValidatePermissionsUpdate(ctx, collection.Permissions, msg.Permissions, msg.Creator)
	if err != nil {
		return nil, err
	}

	//iterate through the non-nil values
	if msg.Permissions.CanDeleteCollection != nil {
		collection.Permissions.CanDeleteCollection = msg.Permissions.CanDeleteCollection
	}

	if msg.Permissions.CanArchive != nil {
		collection.Permissions.CanArchive = msg.Permissions.CanArchive
	}

	if msg.Permissions.CanUpdateContractAddress != nil {
		collection.Permissions.CanUpdateContractAddress = msg.Permissions.CanUpdateContractAddress
	}

	if msg.Permissions.CanUpdateOffChainBalancesMetadata != nil {
		collection.Permissions.CanUpdateOffChainBalancesMetadata = msg.Permissions.CanUpdateOffChainBalancesMetadata
	}

	if msg.Permissions.CanUpdateCustomData != nil {
		collection.Permissions.CanUpdateCustomData = msg.Permissions.CanUpdateCustomData
	}

	if msg.Permissions.CanUpdateStandards != nil {
		collection.Permissions.CanUpdateStandards = msg.Permissions.CanUpdateStandards
	}

	if msg.Permissions.CanUpdateManager != nil {
		collection.Permissions.CanUpdateManager = msg.Permissions.CanUpdateManager
	}

	if msg.Permissions.CanUpdateCollectionMetadata != nil {
		collection.Permissions.CanUpdateCollectionMetadata = msg.Permissions.CanUpdateCollectionMetadata
	}

	if msg.Permissions.CanCreateMoreBadges != nil {
		collection.Permissions.CanCreateMoreBadges = msg.Permissions.CanCreateMoreBadges
	}

	if msg.Permissions.CanUpdateBadgeMetadata != nil {
		collection.Permissions.CanUpdateBadgeMetadata = msg.Permissions.CanUpdateBadgeMetadata
	}

	if msg.Permissions.CanUpdateInheritedBalances != nil {
		collection.Permissions.CanUpdateInheritedBalances = msg.Permissions.CanUpdateInheritedBalances
	}

	if msg.Permissions.CanUpdateCollectionApprovedTransfers != nil {
		collection.Permissions.CanUpdateCollectionApprovedTransfers = msg.Permissions.CanUpdateCollectionApprovedTransfers
	}

	collection.Permissions = msg.Permissions

	if err := k.SetCollectionInStore(ctx, collection); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		),
	)

	return &types.MsgUpdateCollectionPermissionsResponse{}, nil
}
