package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RequestTransferManager(goCtx context.Context, msg *types.MsgRequestTransferManager) (*types.MsgRequestTransferManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	//Redundant because this is locked so we shouldn't store anything
	if msg.Add {
		permissions := types.GetPermissions(badge.PermissionFlags)
		if !permissions.CanManagerTransfer() {
			return nil, ErrInvalidPermissions
		}

		if err := k.CreateTransferManagerRequest(ctx, msg.BadgeId, CreatorAccountNum); err != nil {
			return nil, err
		}
	} else {
		if err := k.RemoveTransferManagerRequest(ctx, msg.BadgeId, CreatorAccountNum); err != nil {
			return nil, err
		}
	}

	return &types.MsgRequestTransferManagerResponse{}, nil
}
