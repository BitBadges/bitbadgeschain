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

	if err := ValidateURI(msg.Uri); err != nil {
		return nil, err
	}

	if err := ValidateURI(msg.SubassetUris); err != nil {
		return nil, err
	}

	if err := ValidateAddress(msg.Manager); err != nil {
		return nil, err
	}

	if err := ValidatePermissions(msg.Permissions); err != nil {
		return nil, err
	}

	//TODO: Validate freeze digest string

	if k.StoreHasBadgeID(ctx, msg.Id) {
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

	k.SetBadgeInStore(ctx, badge);
	
	_ = ctx

	return &types.MsgNewBadgeResponse{}, nil
}
