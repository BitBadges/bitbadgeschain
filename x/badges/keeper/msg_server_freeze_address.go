package keeper

import (
	"context"
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) FreezeAddress(goCtx context.Context, msg *types.MsgFreezeAddress) (*types.MsgFreezeAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	CreatorAccountNum, badge, permissions, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, msg.Addresses, msg.BadgeId, msg.SubbadgeId, true,
	)
	ctx.GasMeter().ConsumeGas(FixedCostPerMsg, "fixed cost per transaction")
	if err != nil {
		return nil, err
	}

	if !permissions.CanFreeze() {
		return nil, ErrInvalidPermissions
	}

	ctx.GasMeter().ConsumeGas(FreezeOrUnfreezeAddress * uint64(len(msg.Addresses)), "pay per address frozen / unfrozen")
	found := false
	for _, targetAddress := range msg.Addresses {
		new_addresses := []uint64{}
		for _, address := range badge.FreezeAddresses {
			if address == targetAddress {
				found = true
			} else {
				new_addresses = append(new_addresses, address)
			}
		}

		if found && msg.Add {
			return nil, ErrAddressAlreadyFrozen
		} else if !found && msg.Add {
			badge.FreezeAddresses = append(badge.FreezeAddresses, targetAddress)
		} else if found && !msg.Add {
			badge.FreezeAddresses = new_addresses
		} else {
			return nil, ErrAddressNotFrozen
		}
	}

	//sort the addresses in order
	sort.Slice(badge.FreezeAddresses, func(i, j int) bool { return badge.FreezeAddresses[i] < badge.FreezeAddresses[j] })

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
			sdk.NewAttribute("Addresses", fmt.Sprint(msg.Addresses)),
			sdk.NewAttribute("Add", fmt.Sprint(msg.Add)),
		),
	)

	return &types.MsgFreezeAddressResponse{}, nil
}
