package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) TransferManager(goCtx context.Context, msg *types.MsgTransferManager) (*types.MsgTransferManagerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.Keeper.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	if err := k.AssertAccountNumbersAreRegistered(ctx, []uint64{msg.Address}); err != nil {
		return nil, ErrAccountsAreNotRegistered
	}

	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	//Transfer to new manager but first need to check privileges
	if badge.Manager != CreatorAccountNum {
		return nil, ErrSenderIsNotManager
	}

	permissions := types.GetPermissions(badge.PermissionFlags)
	if !permissions.CanManagerTransfer() {
		return nil, ErrInvalidPermissions
	}

	requested := k.HasAddressRequestedManagerTransfer(ctx, msg.BadgeId, msg.Address)
	if !requested {
		return nil, ErrAddressNeedsToOptInAndRequestManagerTransfer
	}

	//TODO: other permissions such as remove force mint, etc

	badge.Manager = msg.Address

	if err := k.UpdateBadgeInStore(ctx, badge); err != nil {
		return nil, err
	}

	return &types.MsgTransferManagerResponse{}, nil
}
