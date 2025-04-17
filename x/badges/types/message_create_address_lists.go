package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreateAddressLists = "create_address_lists"

var _ sdk.Msg = &MsgCreateAddressLists{}

func NewMsgCreateAddressLists(creator string, addressLists []*AddressList, creatorOverride string) *MsgCreateAddressLists {
	return &MsgCreateAddressLists{
		Creator:         creator,
		AddressLists:    addressLists,
		CreatorOverride: creatorOverride,
	}
}

func (msg *MsgCreateAddressLists) Route() string {
	return RouterKey
}

func (msg *MsgCreateAddressLists) Type() string {
	return TypeMsgCreateAddressLists
}

func (msg *MsgCreateAddressLists) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
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

	if msg.CreatorOverride != "" {
		_, err = sdk.AccAddressFromBech32(msg.CreatorOverride)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator override address (%s)", err)
		}
	}

	for _, list := range msg.AddressLists {
		if err := ValidateAddressList(list); err != nil {
			return err
		}
	}

	return nil
}
