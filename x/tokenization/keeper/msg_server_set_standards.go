package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
		UpdateStandards:             true,
		Standards:                   msg.Standards,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateStandards: msg.CanUpdateStandards,
			// Copy existing permissions for other fields
			CanDeleteCollection:              collection.GetCollectionPermissions().GetCanDeleteCollection(),
			CanArchiveCollection:             collection.GetCollectionPermissions().GetCanArchiveCollection(),
			CanUpdateCustomData:              collection.GetCollectionPermissions().GetCanUpdateCustomData(),
			CanUpdateManager:                 collection.GetCollectionPermissions().GetCanUpdateManager(),
			CanUpdateValidTokenIds:           collection.GetCollectionPermissions().GetCanUpdateValidTokenIds(),
			CanUpdateCollectionMetadata:      collection.GetCollectionPermissions().GetCanUpdateCollectionMetadata(),
			CanUpdateTokenMetadata:           collection.GetCollectionPermissions().GetCanUpdateTokenMetadata(),
			CanUpdateCollectionApprovals:     collection.GetCollectionPermissions().GetCanUpdateCollectionApprovals(),
			CanAddMoreAliasPaths:             collection.GetCollectionPermissions().GetCanAddMoreAliasPaths(),
			CanAddMoreCosmosCoinWrapperPaths: collection.GetCollectionPermissions().GetCanAddMoreCosmosCoinWrapperPaths(),
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
		sdk.NewAttribute(sdk.AttributeKeyModule, "tokenization"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "set_standards"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgSetStandardsResponse{
		CollectionId: response.CollectionId,
	}, nil
}
