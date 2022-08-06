package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) FreezeAddress(goCtx context.Context, msg *types.MsgFreezeAddress) (*types.MsgFreezeAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum := k.MustGetAccountNumberForBech32AddressString(ctx, msg.Creator)

	badge, found := k.GetBadgeFromStore(ctx, msg.BadgeId)
	if !found {
		return nil, ErrBadgeNotExists
	}

	if badge.Manager != CreatorAccountNum {
		return nil, ErrSenderIsNotManager
	}

	permissions := types.GetPermissions(badge.PermissionFlags)
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")


	if !permissions.CanFreeze() {
		return nil, ErrInvalidPermissions
	}

	// ctx.GasMeter().ConsumeGas(FreezeOrUnfreezeAddress * uint64(len(msg.Addresses)), "pay per address frozen / unfrozen")
	
	found = false

	new_amounts := []*types.RangesToAmounts{
		{
			Ranges: badge.FreezeAddressRanges,
			Amount: 1,
		},
	}

	for targetAddress := msg.Addresses.Start; targetAddress <= msg.Addresses.End; targetAddress++ {
		newAmount := uint64(0)
		if msg.Add {
			newAmount = 1
		}
		new_amounts = UpdateBadgeBalanceBySubbadgeId(targetAddress, newAmount, new_amounts)
	}
	if len(new_amounts) > 0 {
		badge.FreezeAddressRanges = new_amounts[0].Ranges
	} else {
		badge.FreezeAddressRanges = []*types.NumberRange{}
	}
	


	err := k.SetBadgeInStore(ctx, badge)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "FreezeAddress"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeID", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("Addresses", fmt.Sprint(msg.Addresses)),
			sdk.NewAttribute("Add", fmt.Sprint(msg.Add)),
		),
	)

	return &types.MsgFreezeAddressResponse{}, nil
}
