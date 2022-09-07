package keeper

import (
	"context"
	"fmt"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	//For convenience, we will use the same logic as balanceObjects with setting all addresses currently in digest to balance = 1
	//We will set all new addresses that we want to add to balance == 1 and all addresses that we want to remove to balance == 0
	newBalanceObjects := []*types.BalanceObject{}
	if len(badge.FreezeRanges) > 0 {
		newBalanceObjects = append(newBalanceObjects, &types.BalanceObject{
			IdRanges: badge.FreezeRanges, //This needs to be non empty to avoid omit empty case
			Balance:  1,
		})
	}

	newAmount := uint64(0)
	if msg.Add {
		newAmount = 1
	}
	newBalanceObjects = UpdateBalancesForIdRanges(msg.AddressRanges, newAmount, newBalanceObjects)

	if len(newBalanceObjects) > 0 {
		badge.FreezeRanges = newBalanceObjects[0].IdRanges //[0] because we will only ever have one since balances == 0 are not stored
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
