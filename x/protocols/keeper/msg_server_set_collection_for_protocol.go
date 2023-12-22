package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetCollectionForProtocol(goCtx context.Context, msg *types.MsgSetCollectionForProtocol) (*types.MsgSetCollectionForProtocolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	protocolName := msg.Name
	address := msg.Creator

	//Check if protocol exists
	if !k.StoreHasProtocolID(ctx, protocolName) {
		return nil, sdkerrors.Wrap(ErrProtocolDoesNotExist, "Protocol does not exist")
	}

	collectionId := msg.CollectionId
	if collectionId.Equal(sdk.NewUint(0)) {
		//Get next collection id - 1 from badges keeper
		nextCollectionId := k.badgesKeeper.GetNextCollectionId(ctx)

		collectionId = nextCollectionId.Sub(sdk.NewUint(1))
		err := k.SetProtocolCollectionInStore(ctx, protocolName, address, collectionId)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "Failed to set protocol collection in store")
		}
	} else {
		//Update protocol collection in store
		err := k.SetProtocolCollectionInStore(ctx, protocolName, address, collectionId)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "Failed to update protocol collection in store")
		}
	}

	return &types.MsgSetCollectionForProtocolResponse{}, nil
}
