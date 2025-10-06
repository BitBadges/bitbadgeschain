package keeper

import (
	"context"
	"encoding/json"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetCollectionApprovals(goCtx context.Context, msg *types.MsgSetCollectionApprovals) (*types.MsgSetCollectionApprovalsResponse, error) {
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
		UpdateCollectionApprovals:   true,
		CollectionApprovals:         msg.CollectionApprovals,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateCollectionApprovals: msg.CanUpdateCollectionApprovals,
			// Copy existing permissions for other fields
			CanDeleteCollection:         collection.CollectionPermissions.CanDeleteCollection,
			CanArchiveCollection:        collection.CollectionPermissions.CanArchiveCollection,
			CanUpdateStandards:          collection.CollectionPermissions.CanUpdateStandards,
			CanUpdateCustomData:         collection.CollectionPermissions.CanUpdateCustomData,
			CanUpdateManager:            collection.CollectionPermissions.CanUpdateManager,
			CanUpdateValidBadgeIds:      collection.CollectionPermissions.CanUpdateValidBadgeIds,
			CanUpdateCollectionMetadata: collection.CollectionPermissions.CanUpdateCollectionMetadata,
			CanUpdateBadgeMetadata:      collection.CollectionPermissions.CanUpdateBadgeMetadata,
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
		sdk.NewAttribute("msg_type", "set_collection_approvals"),
		sdk.NewAttribute("msg", string(msgBytes)),
	)

	return &types.MsgSetCollectionApprovalsResponse{
		CollectionId: response.CollectionId,
	}, nil
}
