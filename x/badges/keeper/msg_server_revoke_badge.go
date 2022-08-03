package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) RevokeBadge(goCtx context.Context, msg *types.MsgRevokeBadge) (*types.MsgRevokeBadgeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, _, permissions, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, msg.Addresses, msg.BadgeId, msg.SubbadgeId, true,
	)
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "RevokeBadge"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("SubbadgeId", fmt.Sprint(msg.SubbadgeId)),
			sdk.NewAttribute("Addresses", fmt.Sprint(msg.Addresses)),
			sdk.NewAttribute("Amounts", fmt.Sprint(msg.Amounts)),
		),
	)

	return &types.MsgRevokeBadgeResponse{}, nil
}
