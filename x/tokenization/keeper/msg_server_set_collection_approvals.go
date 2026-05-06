package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
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
			CanDeleteCollection:              collection.GetCollectionPermissions().GetCanDeleteCollection(),
			CanArchiveCollection:             collection.GetCollectionPermissions().GetCanArchiveCollection(),
			CanUpdateStandards:               collection.GetCollectionPermissions().GetCanUpdateStandards(),
			CanUpdateCustomData:              collection.GetCollectionPermissions().GetCanUpdateCustomData(),
			CanUpdateManager:                 collection.GetCollectionPermissions().GetCanUpdateManager(),
			CanUpdateValidTokenIds:           collection.GetCollectionPermissions().GetCanUpdateValidTokenIds(),
			CanUpdateCollectionMetadata:      collection.GetCollectionPermissions().GetCanUpdateCollectionMetadata(),
			CanUpdateTokenMetadata:           collection.GetCollectionPermissions().GetCanUpdateTokenMetadata(),
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
		sdk.NewAttribute("msg_type", "set_collection_approvals"),
		sdk.NewAttribute("msg", msgStr),
	)

	return &types.MsgSetCollectionApprovalsResponse{
		CollectionId:    response.CollectionId,
		ApprovalChanges: response.ApprovalChanges,
		ReviewItems:     response.ReviewItems,
	}, nil
}
