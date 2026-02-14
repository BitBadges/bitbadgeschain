package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteCollection(goCtx context.Context, msg *types.MsgDeleteCollection) (*types.MsgDeleteCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, found := k.GetCollectionFromStore(ctx, msg.CollectionId)
	if !found {
		return nil, ErrCollectionNotExists
	}

	err := k.UniversalValidate(ctx, collection, UniversalValidationParams{
		Creator:       msg.Creator,
		MustBeManager: true,
	})
	if err != nil {
		return nil, err
	}

	// Check deleted permission is valid for current time
	err = k.CheckIfActionPermissionPermits(ctx, collection.CollectionPermissions.CanDeleteCollection, "can delete collection")
	if err != nil {
		return nil, err
	}

	// Purge all collection-related state before deleting the collection
	// This includes balances, approval trackers, challenge trackers, approval versions, and ETH signature trackers
	if err := k.PurgeCollectionState(ctx, collection.CollectionId); err != nil {
		return nil, err
	}

	k.DeleteCollectionFromStore(ctx, collection.CollectionId)

	msgStr, err := MarshalMessageForEvent(msg)
	if err != nil {
		return nil, err
	}

	EmitMessageAndIndexerEvents(ctx,
		sdk.NewAttribute(sdk.AttributeKeyModule, "tokenization"),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator),
		sdk.NewAttribute("msg_type", "delete_collection"),
		sdk.NewAttribute("msg", msgStr),
	)
	return &types.MsgDeleteCollectionResponse{}, nil
}
