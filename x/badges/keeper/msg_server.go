package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// resolveCollectionIdWithAutoPrev resolves a collection ID, handling the case where collectionId is 0
// (used for multi-msg transactions). If collectionId is 0, it uses the next collection ID minus 1.
// Returns the resolved collection ID and an error if underflow would occur.
func (k msgServer) resolveCollectionIdWithAutoPrev(ctx sdk.Context, collectionId sdkmath.Uint) (sdkmath.Uint, error) {
	if collectionId.Equal(sdkmath.NewUint(0)) {
		nextCollectionId := k.GetNextCollectionId(ctx)
		// Prevent underflow by checking if nextCollectionId is greater than 0
		if nextCollectionId.IsZero() {
			return sdkmath.Uint{}, sdkerrors.Wrapf(types.ErrInvalidRequest, "cannot calculate collection ID: next collection ID is zero")
		}
		return nextCollectionId.Sub(sdkmath.NewUint(1)), nil
	}
	return collectionId, nil
}
