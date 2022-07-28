package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) NewSubBadge(goCtx context.Context, msg *types.MsgNewSubBadge) (*types.MsgNewSubBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// if err := ValidateBadgeID(msg.Id); err != nil {
	// 	return nil, err
	// }

	if msg.Supply == 0 {
		return nil, ErrSupplyEqualsZero
	}

	badge, found := k.GetBadgeFromStore(ctx, msg.Id)
	if !found {
		return nil, ErrBadgeNotExists
	}

	if badge.Manager != msg.Creator {
		return nil, ErrSenderIsNotManager
	}

	permission_flags := GetPermissions(badge.PermissionFlags)

	if !permission_flags.can_create {
		return nil, ErrInvalidPermissions
	}

	//By default, we assume non fungible (i.e. supply == 1) so we don't store if supply == 1
	subasset_id := badge.NextSubassetId
	if msg.Supply != 1 {
		badge.SubassetsTotalSupply = append(badge.SubassetsTotalSupply, &types.Subasset{
			Id:     subasset_id,
			Supply: msg.Supply,
		})
	}
	badge.NextSubassetId += 1

	manager_address, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	manager_account := k.accountKeeper.GetAccount(ctx, manager_address)
	if manager_account == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "manager account %s does not exist", msg.Creator)
	}

	manager_balance_id := GetFullSubassetID(
		manager_account.GetAccountNumber(),
		msg.Id,
		subasset_id,
	)

	if err := k.AddToBadgeBalance(ctx, manager_balance_id, msg.Supply); err != nil {
		return nil, err
	}

	//Don't have to garbage collect since we are minting so balance > 0

	k.UpdateBadgeInStore(ctx, badge)

	_ = ctx

	return &types.MsgNewSubBadgeResponse{
		SubassetId: subasset_id,
		Message:    "Subbadge created successfully",
	}, nil
}
