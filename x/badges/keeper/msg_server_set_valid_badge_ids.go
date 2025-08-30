package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	oldtypes "github.com/bitbadges/bitbadgeschain/x/badges/types/v13"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CastOldSetValidBadgeIdsToNewType(oldMsg *oldtypes.MsgSetValidBadgeIds) (*types.MsgSetValidBadgeIds, error) {
	// Convert to JSON
	jsonBytes, err := json.Marshal(oldMsg)
	if err != nil {
		return nil, err
	}

	// Unmarshal into new type
	var newMsg types.MsgSetValidBadgeIds
	if err := json.Unmarshal(jsonBytes, &newMsg); err != nil {
		return nil, err
	}

	return &newMsg, nil
}

func (k msgServer) SetValidBadgeIdsV13(goCtx context.Context, msg *oldtypes.MsgSetValidBadgeIds) (*types.MsgSetValidBadgeIdsResponse, error) {
	newMsg, err := CastOldSetValidBadgeIdsToNewType(msg)
	if err != nil {
		return nil, err
	}
	return k.SetValidBadgeIds(goCtx, newMsg)
}

func (k msgServer) SetValidBadgeIdsV14(goCtx context.Context, msg *types.MsgSetValidBadgeIds) (*types.MsgSetValidBadgeIdsResponse, error) {
	return k.SetValidBadgeIds(goCtx, msg)
}

func (k msgServer) SetValidBadgeIds(goCtx context.Context, msg *types.MsgSetValidBadgeIds) (*types.MsgSetValidBadgeIdsResponse, error) {
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
		UpdateValidBadgeIds:         true,
		ValidBadgeIds:               msg.ValidBadgeIds,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateValidBadgeIds: msg.CanUpdateValidBadgeIds,
			// Copy existing permissions for other fields
			CanDeleteCollection:               collection.CollectionPermissions.CanDeleteCollection,
			CanArchiveCollection:              collection.CollectionPermissions.CanArchiveCollection,
			CanUpdateOffChainBalancesMetadata: collection.CollectionPermissions.CanUpdateOffChainBalancesMetadata,
			CanUpdateStandards:                collection.CollectionPermissions.CanUpdateStandards,
			CanUpdateCustomData:               collection.CollectionPermissions.CanUpdateCustomData,
			CanUpdateManager:                  collection.CollectionPermissions.CanUpdateManager,
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
				sdk.NewAttribute("msg_type", "set_valid_badge_ids"),
				sdk.NewAttribute("msg", string(msgBytes)),
			),
		)
	}

	return &types.MsgSetValidBadgeIdsResponse{
		CollectionId: response.CollectionId,
	}, nil
}
