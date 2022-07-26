package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewBadge(goCtx context.Context, msg *types.MsgNewBadge) (*types.MsgNewBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := ValidateBadgeID(msg.Id); err != nil {
		return nil, err
	}

	if k.HasBadge(ctx, msg.Id) {
		return nil, ErrBadgeExists
	}

	
	m := make(map[uint64]uint64)

	badge := types.BitBadge{
		Id:        							msg.Id,
		Uri:	   							msg.Uri,
		Creator:   							msg.Creator,
		Manager:   							msg.Manager,
		PermissionFlags: 					msg.Permissions,
		SubassetUriFormat: 					msg.SubassetUris,
		SubassetsTotalSupply: 				m,
		NextSubassetId: 					0,
		FrozenOrUnfrozenAddressesDigest:	msg.FreezeAddressesDigest,
	}

	k.SetBadge(ctx, badge);
	
	// TODO: Handling the message
	_ = ctx

	return &types.MsgNewBadgeResponse{}, nil
}
