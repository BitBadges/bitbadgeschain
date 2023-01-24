package keeper

import (
	"context"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) RegisterAddresses(goCtx context.Context, msg *types.MsgRegisterAddresses) (*types.MsgRegisterAddressesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	start := uint64(0)
	end := uint64(0)

	for i, address := range msg.AddressesToRegister {
		convertedAddress, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
		}

		newNum := k.Keeper.GetOrCreateAccountNumberForAccAddressBech32(ctx, convertedAddress) //This panics but is saved
		if i == 0 {
			start = newNum
		}
		end = newNum
	}

	return &types.MsgRegisterAddressesResponse{
		RegisteredAddressNumbers: GetIdRangeToInsert(start, end),
	}, nil
}
