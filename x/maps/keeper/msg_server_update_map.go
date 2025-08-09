package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/maps/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	badgetypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdateMap(goCtx context.Context, msg *types.MsgUpdateMap) (*types.MsgUpdateMapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx

	currMap, found := k.GetMapFromStore(ctx, msg.MapId)
	if !found {
		return nil, sdkerrors.Wrap(ErrMapDoesNotExist, "Failed to get map from store")
	}

	collection := &badgetypes.TokenCollection{}
	if !currMap.InheritManagerTimelineFrom.IsNil() && !currMap.InheritManagerTimelineFrom.IsZero() {
		collectionRes, err := k.badgesKeeper.GetCollection(ctx, &badgetypes.QueryGetCollectionRequest{CollectionId: currMap.InheritManagerTimelineFrom.String()})
		if err != nil {
			return nil, sdkerrors.Wrap(ErrInvalidMapId, "Could not find collection in store")
		}

		collection = collectionRes.Collection
	}

	currManager := types.GetCurrentManagerForMap(ctx, currMap, collection)
	if currManager != msg.Creator {
		return nil, sdkerrors.Wrapf(ErrNotMapCreator, "current manager is %s but got %s", currManager, msg.Creator)
	}

	if msg.UpdateManagerTimeline {
		if err := k.badgesKeeper.ValidateManagerUpdate(ctx, types.CastManagerTimelineArray(currMap.ManagerTimeline), types.CastManagerTimelineArray(msg.ManagerTimeline), types.CastTimedUpdatePermissions(currMap.Permissions.CanUpdateManager)); err != nil {
			return nil, err
		}
		currMap.ManagerTimeline = msg.ManagerTimeline
	}

	if msg.UpdateMetadataTimeline {
		if err := k.badgesKeeper.ValidateCollectionMetadataUpdate(ctx, types.CastMetadataTimelineArray(currMap.MetadataTimeline), types.CastMetadataTimelineArray(msg.MetadataTimeline), types.CastTimedUpdatePermissions(currMap.Permissions.CanUpdateMetadata)); err != nil {
			return nil, err
		}
		currMap.MetadataTimeline = msg.MetadataTimeline
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

		if err := k.badgesKeeper.ValidateTimedUpdatePermissionUpdate(ctx, types.CastTimedUpdatePermissions(currMap.Permissions.CanUpdateManager), types.CastTimedUpdatePermissions(msg.Permissions.CanUpdateManager)); err != nil {
			return nil, err
		}

		if err := k.badgesKeeper.ValidateTimedUpdatePermissionUpdate(ctx, types.CastTimedUpdatePermissions(currMap.Permissions.CanUpdateMetadata), types.CastTimedUpdatePermissions(msg.Permissions.CanUpdateMetadata)); err != nil {
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
