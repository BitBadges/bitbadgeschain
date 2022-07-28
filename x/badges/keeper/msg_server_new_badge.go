package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewBadge(goCtx context.Context, msg *types.MsgNewBadge) (*types.MsgNewBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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

	badge_id := k.GetNextAssetId(ctx)

	//probably redundant
	// if err := ValidateBadgeID(badge_id); err != nil {
	// 	return nil, err
	// }

	if k.StoreHasBadgeID(ctx, badge_id) {
		return nil, ErrBadgeExists
	}

	badge := types.BitBadge{
		Id:                              badge_id,
		Uri:                             msg.Uri,
		Manager:                         msg.Manager,
		PermissionFlags:                 msg.Permissions,
		SubassetUriFormat:               msg.SubassetUris,
		SubassetsTotalSupply:            []*types.Subasset{},
		NextSubassetId:                  0,
		FrozenOrUnfrozenAddressesDigest: msg.FreezeAddressesDigest,
	}

	k.SetBadgeInStore(ctx, badge)

	_ = ctx

	k.SetNextAssetId(ctx, badge_id+1)

	return &types.MsgNewBadgeResponse{
		Id:      badge_id,
		Message: "Badge created successfully",
	}, nil
}
