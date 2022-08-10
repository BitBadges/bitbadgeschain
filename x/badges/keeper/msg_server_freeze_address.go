package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) FreezeAddress(goCtx context.Context, msg *types.MsgFreezeAddress) (*types.MsgFreezeAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	CreatorAccountNum, badge, err := k.UniversalValidate(ctx, 
		UniversalValidationParams{
			Creator: msg.Creator,
			BadgeId: msg.BadgeId,
			MustBeManager: true,
			CanFreeze: true,
		},
	)
	if err != nil {
		return nil, err
	}

	new_amounts := []*types.BalanceToIds{
		{
			Ids: badge.FreezeAddressRanges,
			Balance: 1,
		},
	}

	for _, addressRange := range msg.AddressRanges {
		for targetAddress := addressRange.Start; targetAddress <= addressRange.End; targetAddress++ {
			newAmount := uint64(0)
			if msg.Add {
				newAmount = 1
			}
			new_amounts = UpdateBalanceForSubbadgeId(targetAddress, newAmount, new_amounts)
		}
		if len(new_amounts) > 0 {
			badge.FreezeAddressRanges = new_amounts[0].Ids
		} else {
			badge.FreezeAddressRanges = []*types.NumberRange{}
		}
	}

	err = k.SetBadgeInStore(ctx, badge)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "FreezeAddress"),
			sdk.NewAttribute("Creator", fmt.Sprint(CreatorAccountNum)),
			sdk.NewAttribute("BadgeID", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("AddressRanges", fmt.Sprint(msg.AddressRanges)),
			sdk.NewAttribute("Add", fmt.Sprint(msg.Add)),
		),
	)

	return &types.MsgFreezeAddressResponse{}, nil
}
