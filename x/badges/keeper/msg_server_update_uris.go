package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdateUris(goCtx context.Context, msg *types.MsgUpdateUris) (*types.MsgUpdateUrisResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	if badge.Manager != CreatorAccountNum {
		return nil, ErrSenderIsNotManager
	}

	permissions := types.GetPermissions(badge.PermissionFlags)
	if !permissions.CanUpdateUris() {
		return nil, ErrInvalidPermissions
	}

	badge.Uri = msg.Uri
	badge.SubassetUriFormat = msg.SubassetUri

	if err := k.UpdateBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	return &types.MsgUpdateUrisResponse{}, nil
}
