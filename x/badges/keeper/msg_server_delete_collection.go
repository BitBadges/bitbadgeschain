package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) DeleteCollection(goCtx context.Context, msg *types.MsgDeleteCollection) (*types.MsgDeleteCollectionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		CollectionId:  msg.CollectionId,
		MustBeManager: true,
		CanDelete: 	   true,	
	})

	k.DeleteCollectionFromStore(ctx, badge.CollectionId)
	if err != nil {
		return nil, err
	}

	return &types.MsgDeleteCollectionResponse{}, nil
}
