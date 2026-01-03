package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/maps/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokentypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdateMap(goCtx context.Context, msg *types.MsgUpdateMap) (*types.MsgUpdateMapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	currMap, found := k.GetMapFromStore(ctx, msg.MapId)
	if !found {
		return nil, sdkerrors.Wrap(ErrMapDoesNotExist, "Failed to get map from store")
	}

	collection := &tokentypes.TokenCollection{}
	if !currMap.InheritManagerFrom.IsNil() && !currMap.InheritManagerFrom.IsZero() {
		collectionRes, err := k.badgesKeeper.GetCollection(ctx, &tokentypes.QueryGetCollectionRequest{CollectionId: currMap.InheritManagerFrom.String()})
		if err != nil {
			return nil, sdkerrors.Wrap(ErrInvalidMapId, "Could not find collection in store")
		}

		collection = collectionRes.Collection
	}

	currManager := types.GetCurrentManagerForMap(ctx, currMap, collection)
	if currManager != msg.Creator {
		return nil, sdkerrors.Wrapf(ErrNotMapCreator, "current manager is %s but got %s", currManager, msg.Creator)
	}

	if msg.UpdateManager {
		oldManager := currMap.Manager
		newManager := msg.Manager
		if err := k.badgesKeeper.ValidateManagerUpdate(ctx, oldManager, newManager, types.CastActionPermissions(currMap.Permissions.CanUpdateManager)); err != nil {
			return nil, err
		}
		currMap.Manager = msg.Manager
	}

	if msg.UpdateMetadata {
		oldMetadata := types.CastMapMetadataToCollectionMetadata(currMap.Metadata)
		newMetadata := types.CastMapMetadataToCollectionMetadata(msg.Metadata)
		if err := k.badgesKeeper.ValidateCollectionMetadataUpdate(ctx, oldMetadata, newMetadata, types.CastActionPermissions(currMap.Permissions.CanUpdateMetadata)); err != nil {
			return nil, err
		}
		currMap.Metadata = msg.Metadata
	}

	if msg.UpdatePermissions {
		if err := types.ValidatePermissions(msg.Permissions, true); err != nil {
			return nil, err
		}

		if err := types.ValidatePermissions(currMap.Permissions, true); err != nil {
			return nil, err
		}

		if err := k.badgesKeeper.ValidateActionPermissionUpdate(ctx, types.CastActionPermissions(currMap.Permissions.CanDeleteMap), types.CastActionPermissions(msg.Permissions.CanDeleteMap)); err != nil {
			return nil, err
		}

		if err := k.badgesKeeper.ValidateActionPermissionUpdate(ctx, types.CastActionPermissions(currMap.Permissions.CanUpdateManager), types.CastActionPermissions(msg.Permissions.CanUpdateManager)); err != nil {
			return nil, err
		}

		if err := k.badgesKeeper.ValidateActionPermissionUpdate(ctx, types.CastActionPermissions(currMap.Permissions.CanUpdateMetadata), types.CastActionPermissions(msg.Permissions.CanUpdateMetadata)); err != nil {
			return nil, err
		}

		currMap.Permissions = msg.Permissions
	}

	//Add protocol to store
	err := k.SetMapInStore(ctx, currMap)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Failed to add protocol to store")
	}

	return &types.MsgUpdateMapResponse{}, nil
}
