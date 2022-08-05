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
		ctx, msg.Creator, msg.Addresses, msg.BadgeId, msg.SubbadgeRange.End, true,
	)
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")
	if err != nil {
		return nil, err
	}

	if !permissions.CanRevoke() {
		return nil, ErrInvalidPermissions
	}

	ManagerBalanceKey := GetBalanceKey(CreatorAccountNum, msg.BadgeId)
	managerBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, ManagerBalanceKey)
	if !found {
		return nil, ErrBadgeBalanceNotExists
	}

	for i, revokeAddress := range msg.Addresses {
		if revokeAddress == CreatorAccountNum {
			return nil, ErrSenderAndReceiverSame
		}

		// Note that we check for duplicates in ValidateBasic, so these addresses will be unique every time
		AddressBalanceKey := GetBalanceKey(revokeAddress, msg.BadgeId)
		addressBalanceInfo, found := k.Keeper.GetBadgeBalanceFromStore(ctx, AddressBalanceKey)
		if !found {
			return nil, ErrBadgeBalanceNotExists
		}
		

		revokeAmount := msg.Amounts[i]

		for i := msg.SubbadgeRange.Start; i <= msg.SubbadgeRange.End; i++ {
			addressBalanceInfo, err = k.RemoveFromBadgeBalance(ctx, addressBalanceInfo, i, revokeAmount)
			if err != nil {
				return nil, err
			}

			managerBalanceInfo, err = k.AddToBadgeBalance(ctx, managerBalanceInfo, i, revokeAmount)
			if err != nil {
				return nil, err
			}

		}
		err = k.SetBadgeBalanceInStore(ctx, AddressBalanceKey, addressBalanceInfo)
		if err != nil {
			return nil, err
		}
	}

	err = k.SetBadgeBalanceInStore(ctx, ManagerBalanceKey, managerBalanceInfo)
	if err != nil {
		return nil, err
	}

	

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "RevokeBadge"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeId", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("Subbadge Range", fmt.Sprint(msg.SubbadgeRange)),
			sdk.NewAttribute("Addresses", fmt.Sprint(msg.Addresses)),
			sdk.NewAttribute("Amounts", fmt.Sprint(msg.Amounts)),
		),
	)

	return &types.MsgRevokeBadgeResponse{}, nil
}
