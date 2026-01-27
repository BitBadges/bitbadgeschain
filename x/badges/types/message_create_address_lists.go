package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreateAddressLists = "create_address_lists"

var _ sdk.Msg = &MsgCreateAddressLists{}

func NewMsgCreateAddressLists(creator string, addressLists []*AddressListInput) *MsgCreateAddressLists {
	return &MsgCreateAddressLists{
		Creator:      creator,
		AddressLists: addressLists,
	}
}

func (msg *MsgCreateAddressLists) Route() string {
	return RouterKey
}

func (msg *MsgCreateAddressLists) Type() string {
	return TypeMsgCreateAddressLists
}

func (msg *MsgCreateAddressLists) GetSigners() []sdk.AccAddress {
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateAddressLists) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAddressLists) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	for _, list := range msg.AddressLists {
		if err := ValidateAddressListInput(list); err != nil {
			return err
		}
	}

	return nil
}
