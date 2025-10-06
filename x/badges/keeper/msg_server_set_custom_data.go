package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetCustomData(goCtx context.Context, msg *types.MsgSetCustomData) (*types.MsgSetCustomDataResponse, error) {
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
		UpdateCustomDataTimeline:    true,
		CustomDataTimeline:          msg.CustomDataTimeline,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateCustomData: msg.CanUpdateCustomData,
			// Copy existing permissions for other fields
			CanDeleteCollection:               collection.CollectionPermissions.CanDeleteCollection,
			CanArchiveCollection:              collection.CollectionPermissions.CanArchiveCollection,
			CanUpdateOffChainBalancesMetadata: collection.CollectionPermissions.CanUpdateOffChainBalancesMetadata,
			CanUpdateStandards:                collection.CollectionPermissions.CanUpdateStandards,
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
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "set_custom_data"),
		sdk.NewAttribute("msg", string(msgBytes)),
	)

	return &types.MsgSetCustomDataResponse{
		CollectionId: response.CollectionId,
	}, nil
}
