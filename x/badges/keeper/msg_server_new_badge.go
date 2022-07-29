package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewBadge(goCtx context.Context, msg *types.MsgNewBadge) (*types.MsgNewBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Creator will already be registered, so we can do this and panic if it fails
	creator_account_num := k.Keeper.MustGetAccountNumberForAddressString(ctx, msg.Creator)

	//Validate well-formedness of the message entries
	if err := ValidateURI(msg.Uri); err != nil {
		return nil, err
	}

	if err := ValidateURI(msg.SubassetUris); err != nil {
		return nil, err
	}

	if err := ValidatePermissions(msg.Permissions); err != nil {
		return nil, err
	}

	//TODO: Validate freeze digest string


	//Get next badge ID and increment
	next_badge_id := k.GetNextAssetId(ctx)
	k.SetNextAssetId(ctx, next_badge_id + 1)

	//Create and store the badge
	badge := types.BitBadge{
		Id:                              next_badge_id,
		Uri:                             msg.Uri,
		Manager:                         creator_account_num,
		PermissionFlags:                 msg.Permissions,
		SubassetUriFormat:               msg.SubassetUris,
		SubassetsTotalSupply:            []*types.Subasset{},
		NextSubassetId:                  0,
		FrozenOrUnfrozenAddressesDigest: msg.FreezeAddressesDigest,
	}
	if err := k.SetBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}
	
	return &types.MsgNewBadgeResponse{
		Id:      next_badge_id,
	}, nil
}
