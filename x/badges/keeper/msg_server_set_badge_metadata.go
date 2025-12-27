package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetTokenMetadata(goCtx context.Context, msg *types.MsgSetTokenMetadata) (*types.MsgSetTokenMetadataResponse, error) {
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
		UpdateTokenMetadata:         true,
		TokenMetadata:               msg.TokenMetadata,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateTokenMetadata: msg.CanUpdateTokenMetadata,
			// Copy existing permissions for other fields
			CanDeleteCollection:          collection.CollectionPermissions.CanDeleteCollection,
			CanArchiveCollection:         collection.CollectionPermissions.CanArchiveCollection,
			CanUpdateStandards:           collection.CollectionPermissions.CanUpdateStandards,
			CanUpdateCustomData:          collection.CollectionPermissions.CanUpdateCustomData,
			CanUpdateManager:             collection.CollectionPermissions.CanUpdateManager,
			CanUpdateValidTokenIds:       collection.CollectionPermissions.CanUpdateValidTokenIds,
			CanUpdateCollectionMetadata:  collection.CollectionPermissions.CanUpdateCollectionMetadata,
			CanUpdateCollectionApprovals: collection.CollectionPermissions.CanUpdateCollectionApprovals,
		},
	}

	// Call the existing UniversalUpdateCollection handler
	response, err := k.UniversalUpdateCollection(goCtx, universalMsg)
	if err != nil {
		return nil, err
	}

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "set_token_metadata"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgSetTokenMetadataResponse{
		CollectionId: response.CollectionId,
	}, nil
}
