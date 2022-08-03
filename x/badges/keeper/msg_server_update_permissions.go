package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) UpdatePermissions(goCtx context.Context, msg *types.MsgUpdatePermissions) (*types.MsgUpdatePermissionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	if badge.Manager != CreatorAccountNum {
		return nil, ErrSenderIsNotManager
	}

	err := types.ValidatePermissionsUpdate(badge.PermissionFlags, msg.Permissions)
	if err != nil {
		return nil, err
	}

	badge.PermissionFlags = msg.Permissions

	if err := k.UpdateBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	return &types.MsgUpdatePermissionsResponse{}, nil
}
