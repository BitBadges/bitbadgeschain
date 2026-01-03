package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetValidTokenIds(goCtx context.Context, msg *types.MsgSetValidTokenIds) (*types.MsgSetValidTokenIdsResponse, error) {
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
		UpdateValidTokenIds:         true,
		ValidTokenIds:               msg.ValidTokenIds,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateValidTokenIds: msg.CanUpdateValidTokenIds,
			// Copy existing permissions for other fields
			CanDeleteCollection:                collection.CollectionPermissions.CanDeleteCollection,
			CanArchiveCollection:                collection.CollectionPermissions.CanArchiveCollection,
			CanUpdateStandards:                  collection.CollectionPermissions.CanUpdateStandards,
			CanUpdateCustomData:                 collection.CollectionPermissions.CanUpdateCustomData,
			CanUpdateManager:                    collection.CollectionPermissions.CanUpdateManager,
			CanUpdateCollectionMetadata:         collection.CollectionPermissions.CanUpdateCollectionMetadata,
			CanUpdateTokenMetadata:              collection.CollectionPermissions.CanUpdateTokenMetadata,
			CanUpdateCollectionApprovals:        collection.CollectionPermissions.CanUpdateCollectionApprovals,
			CanAddMoreAliasPaths:                collection.CollectionPermissions.CanAddMoreAliasPaths,
			CanAddMoreCosmosCoinWrapperPaths:    collection.CollectionPermissions.CanAddMoreCosmosCoinWrapperPaths,
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
		sdk.NewAttribute("msg_type", "set_valid_token_ids"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgSetValidTokenIdsResponse{
		CollectionId: response.CollectionId,
	}, nil
}
