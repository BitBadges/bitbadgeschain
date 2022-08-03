package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RevokeBadge(goCtx context.Context, msg *types.MsgRevokeBadge) (*types.MsgRevokeBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, _, permissions, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, []uint64{msg.Address}, msg.BadgeId, msg.SubbadgeId, true,
	)
	if err != nil {
		return nil, err
	}

	if CreatorAccountNum == msg.Address {
		return nil, ErrSenderAndReceiverSame
	}

	if !permissions.CanRevoke() {
		return nil, ErrInvalidPermissions
	}

	AddressBalanceKey := GetBalanceKey(msg.Address, msg.BadgeId, msg.SubbadgeId)
	ManagerBalanceKey := GetBalanceKey(CreatorAccountNum, msg.BadgeId, msg.SubbadgeId)

	err = k.RemoveFromBadgeBalance(ctx, AddressBalanceKey, msg.Amount)
	if err != nil {
		return nil, err
	}

	err = k.AddToBadgeBalance(ctx, ManagerBalanceKey, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgRevokeBadgeResponse{}, nil
}
