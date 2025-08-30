package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldSetStandardsToNewType(oldMsg *oldtypes.MsgSetStandards) (*types.MsgSetStandards, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newMsg types.MsgSetStandards
	if err := json.Unmarshal(jsonBytes, &newMsg); err != nil {
		return nil, err
	}

	return &newMsg, nil
}

func (k msgServer) SetStandardsV13(goCtx context.Context, msg *oldtypes.MsgSetStandards) (*types.MsgSetStandardsResponse, error) {
	newMsg, err := CastOldSetStandardsToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.SetStandards(goCtx, newMsg)
}

func (k msgServer) SetStandardsV14(goCtx context.Context, msg *types.MsgSetStandards) (*types.MsgSetStandardsResponse, error) {
	return k.SetStandards(goCtx, msg)
}

func (k msgServer) SetStandards(goCtx context.Context, msg *types.MsgSetStandards) (*types.MsgSetStandardsResponse, error) {
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
		Creator:                     msg.Creator,
		CollectionId:                msg.CollectionId,
		UpdateStandardsTimeline:     true,
		StandardsTimeline:           msg.StandardsTimeline,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateStandards: msg.CanUpdateStandards,
			// Copy existing permissions for other fields
			CanDeleteCollection:               collection.CollectionPermissions.CanDeleteCollection,
			CanArchiveCollection:              collection.CollectionPermissions.CanArchiveCollection,
			CanUpdateOffChainBalancesMetadata: collection.CollectionPermissions.CanUpdateOffChainBalancesMetadata,
			CanUpdateCustomData:               collection.CollectionPermissions.CanUpdateCustomData,
			CanUpdateManager:                  collection.CollectionPermissions.CanUpdateManager,
			CanUpdateValidBadgeIds:            collection.CollectionPermissions.CanUpdateValidBadgeIds,
			CanUpdateCollectionMetadata:       collection.CollectionPermissions.CanUpdateCollectionMetadata,
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
				sdk.NewAttribute("msg_type", "set_standards"),
				sdk.NewAttribute("msg", string(msgBytes)),
			),
		)
	}

	return &types.MsgSetStandardsResponse{
		CollectionId: response.CollectionId,
	}, nil
}
