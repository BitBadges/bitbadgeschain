package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewBadge(goCtx context.Context, msg *types.MsgNewBadge) (*types.MsgNewBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	NextBadgeId := k.GetNextAssetId(ctx)
	k.IncrementNextAssetId(ctx)

	badge := types.BitBadge{
		Id:                NextBadgeId,
		Uri:               msg.Uri,
		Manager:           CreatorAccountNum,
		PermissionFlags:   msg.Permissions,
		SubassetUriFormat: msg.SubassetUris,
		// SubassetsTotalSupply: []*types.Subasset{},
		// NextSubassetId:       0,
		// FreezeAddresses:      []uint64{},
	}

	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	return &types.MsgNewBadgeResponse{
		Id: NextBadgeId,
	}, nil
}
