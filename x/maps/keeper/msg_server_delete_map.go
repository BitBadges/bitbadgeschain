package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/maps/types"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokentypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func (k msgServer) DeleteMap(goCtx context.Context, msg *types.MsgDeleteMap) (*types.MsgDeleteMapResponse, error) {
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

	//Check deleted permission is valid for current time
	err := k.badgesKeeper.CheckIfActionPermissionPermits(ctx, types.CastActionPermissions(currMap.Permissions.CanDeleteMap), "can delete map")
	if err != nil {
		return nil, err
	}

	k.DeleteMapFromStore(ctx, msg.MapId)

	return &types.MsgDeleteMapResponse{}, nil
}
