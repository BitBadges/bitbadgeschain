package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldSetCollectionMetadataToNewType(oldMsg *oldtypes.MsgSetCollectionMetadata) (*types.MsgSetCollectionMetadata, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newMsg types.MsgSetCollectionMetadata
	if err := json.Unmarshal(jsonBytes, &newMsg); err != nil {
		return nil, err
	}

	return &newMsg, nil
}

func (k msgServer) SetCollectionMetadataV13(goCtx context.Context, msg *oldtypes.MsgSetCollectionMetadata) (*types.MsgSetCollectionMetadataResponse, error) {
	newMsg, err := CastOldSetCollectionMetadataToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.SetCollectionMetadata(goCtx, newMsg)
}

func (k msgServer) SetCollectionMetadataV14(goCtx context.Context, msg *types.MsgSetCollectionMetadata) (*types.MsgSetCollectionMetadataResponse, error) {
	return k.SetCollectionMetadata(goCtx, msg)
}

func (k msgServer) SetCollectionMetadata(goCtx context.Context, msg *types.MsgSetCollectionMetadata) (*types.MsgSetCollectionMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Get existing collection to fetch current permissions
	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	// Construct the full UniversalUpdateCollection message
	universalMsg := &types.MsgUniversalUpdateCollection{
		Creator:                          msg.Creator,
		CollectionId:                     msg.CollectionId,
		UpdateCollectionMetadataTimeline: true,
		CollectionMetadataTimeline:       msg.CollectionMetadataTimeline,
		UpdateCollectionPermissions:      true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateCollectionMetadata: msg.CanUpdateCollectionMetadata,
			// Copy existing permissions for other fields
			CanDeleteCollection:               collection.CollectionPermissions.CanDeleteCollection,
			CanArchiveCollection:              collection.CollectionPermissions.CanArchiveCollection,
			CanUpdateOffChainBalancesMetadata: collection.CollectionPermissions.CanUpdateOffChainBalancesMetadata,
			CanUpdateStandards:                collection.CollectionPermissions.CanUpdateStandards,
			CanUpdateCustomData:               collection.CollectionPermissions.CanUpdateCustomData,
			CanUpdateManager:                  collection.CollectionPermissions.CanUpdateManager,
			CanUpdateValidBadgeIds:            collection.CollectionPermissions.CanUpdateValidBadgeIds,
			CanUpdateBadgeMetadata:            collection.CollectionPermissions.CanUpdateBadgeMetadata,
			CanUpdateCollectionApprovals:      collection.CollectionPermissions.CanUpdateCollectionApprovals,
		},
	}

	// Call the existing UniversalUpdateCollection handler
	response, err := k.UniversalUpdateCollection(goCtx, universalMsg)
	if err != nil {
		return nil, err
	}

	msgBytes, err := json.Marshal(msg)
	if err == nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
				sdk.NewAttribute("msg_type", "set_collection_metadata"),
				sdk.NewAttribute("msg", string(msgBytes)),
			),
		)
	}

	return &types.MsgSetCollectionMetadataResponse{
		CollectionId: response.CollectionId,
	}, nil
}
