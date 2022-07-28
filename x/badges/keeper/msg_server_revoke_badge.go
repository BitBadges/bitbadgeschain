package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)


func (k msgServer) RevokeBadge(goCtx context.Context,  msg *types.MsgRevokeBadge) (*types.MsgRevokeBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	permissions := GetPermissions(badge.PermissionFlags)

	if badge.Manager != msg.Creator {
		return nil, ErrSenderIsNotManager
	}

	if !permissions.can_revoke {
		return nil, ErrInvalidPermissions
	}

	manager, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	
	manager_account := k.accountKeeper.GetAccount(ctx, manager)
	if manager_account == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", manager)
	}
	manager_account_num := manager_account.GetAccountNumber()

	address, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}
	

	address_account := k.accountKeeper.GetAccount(ctx, address)
	if address_account == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", address)
	}
	address_account_num := address_account.GetAccountNumber()

	balance_id := GetFullSubassetID(address_account_num, msg.BadgeId, msg.SubbadgeId)
	manager_balance_id := GetFullSubassetID(manager_account_num, msg.BadgeId, msg.SubbadgeId)

	err = k.Keeper.RevokeBadge(ctx, balance_id, manager_balance_id, msg.Amount)
	if err != nil {
		return nil, err
	}


    // TODO: Handling the message
    _ = ctx

	return &types.MsgRevokeBadgeResponse{
		Message: "Success!",
	}, nil
}
