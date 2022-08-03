package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RevokeBadge(goCtx context.Context, msg *types.MsgRevokeBadge) (*types.MsgRevokeBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, _, permissions, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, msg.Addresses, msg.BadgeId, msg.SubbadgeId, true,
	)
	if err != nil {
		return nil, err
	}

	if !permissions.CanRevoke() {
		return nil, ErrInvalidPermissions
	}

	for i, revokeAddress := range msg.Addresses {
		if revokeAddress == CreatorAccountNum {
			return nil, ErrSenderAndReceiverSame
		}

		AddressBalanceKey := GetBalanceKey(revokeAddress, msg.BadgeId, msg.SubbadgeId)
		ManagerBalanceKey := GetBalanceKey(CreatorAccountNum, msg.BadgeId, msg.SubbadgeId)

		revokeAmount := msg.Amounts[i]
		err = k.RemoveFromBadgeBalance(ctx, AddressBalanceKey, revokeAmount)
		if err != nil {
			return nil, err
		}

		err = k.AddToBadgeBalance(ctx, ManagerBalanceKey, revokeAmount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgRevokeBadgeResponse{}, nil
}
