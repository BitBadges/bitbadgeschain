package keeper

import (
	"context"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (k msgServer) FreezeAddress(goCtx context.Context, msg *types.MsgFreezeAddress) (*types.MsgFreezeAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, badge, permissions, err := k.Keeper.UniversalValidateMsgAndReturnMsgInfo(
		ctx, msg.Creator, []uint64{msg.Address}, msg.BadgeId, msg.SubbadgeId, true,
	)
	if err != nil {
		return nil, err
	}

	if !permissions.CanFreeze() {
		return nil, ErrInvalidPermissions
	}

	found := false
	new_addresses := []uint64{}
	//TODO: binary search (they are sorted)
	for _, address := range badge.FreezeAddresses {
		if address == msg.Address {
			found = true
		} else {
			new_addresses = append(new_addresses, address)
		}
	}

	if found && msg.Add {
		return nil, ErrAddressAlreadyFrozen
	} else if !found && msg.Add {
		badge.FreezeAddresses = append(badge.FreezeAddresses, msg.Address)
	} else if found && !msg.Add {
		badge.FreezeAddresses = new_addresses
	} else {
		return nil, ErrAddressNotFrozen
	}

	//sort the addresses in order
	sort.Slice(badge.FreezeAddresses, func(i, j int) bool { return badge.FreezeAddresses[i] < badge.FreezeAddresses[j] })

	err = k.UpdateBadgeInStore(ctx, badge)
	if err != nil {
		return nil, err
	}

	return &types.MsgFreezeAddressResponse{}, nil
}
