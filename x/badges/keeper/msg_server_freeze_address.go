package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) FreezeAddress(goCtx context.Context, msg *types.MsgFreezeAddress) (*types.MsgFreezeAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, badge, err := k.UniversalValidate(ctx, UniversalValidationParams{
		Creator:       msg.Creator,
		BadgeId:       msg.BadgeId,
		MustBeManager: true,
		CanFreeze:     true,
	})
	if err != nil {
		return nil, err
	}

	//For convenience, we will use the same logic as balanceObjects with setting all addresses in digest to balance = 1
	//We will set all new addresses that we want to add to balance == 1 and all addresses that we want to remove to balance == 0
	newBalanceObjects := []*types.BalanceObject{{
		IdRanges: badge.FreezeRanges,
		Balance:  1,
	}}

	for _, addressRange := range msg.AddressRanges {
		for targetAddress := addressRange.Start; targetAddress <= addressRange.End; targetAddress++ {
			newAmount := uint64(0)
			if msg.Add {
				newAmount = 1
			}

			newBalanceObjects = UpdateBalanceForId(targetAddress, newAmount, newBalanceObjects)
		}
	}

	if len(newBalanceObjects) > 0 {
		badge.FreezeRanges = newBalanceObjects[0].IdRanges
	} else {
		badge.FreezeRanges = nil
	}

	err = k.SetBadgeInStore(ctx, badge)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "badges"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "FreezeAddress"),
			sdk.NewAttribute("BadgeID", fmt.Sprint(msg.BadgeId)),
			sdk.NewAttribute("AddressRanges", fmt.Sprint(msg.AddressRanges)),
			sdk.NewAttribute("Add", fmt.Sprint(msg.Add)),
		),
	)

	return &types.MsgFreezeAddressResponse{}, nil
}
