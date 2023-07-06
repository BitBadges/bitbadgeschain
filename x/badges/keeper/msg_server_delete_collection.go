package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteCollection(goCtx context.Context, msg *types.MsgDeleteCollection) (*types.MsgDeleteCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	collection, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:             msg.Creator,
		CollectionId:        msg.CollectionId,
		MustBeManager:       true,
	})
	if err != nil {
		return nil, err
	}

	//Check deleted permission is valid for current time
	err = k.CheckActionPermission(ctx, collection.Permissions.CanDeleteCollection)
	if err != nil {
		return nil, err
	}

	k.DeleteCollectionFromStore(ctx, collection.CollectionId)
	if err != nil {
		return nil, err
	}

	

	return &types.MsgDeleteCollectionResponse{}, nil
}
