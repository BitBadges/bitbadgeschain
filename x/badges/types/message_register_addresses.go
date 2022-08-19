package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterAddresses = "register_addresses"

var _ sdk.Msg = &MsgRegisterAddresses{}

func NewMsgRegisterAddresses(creator string, addressesToRegister []string) *MsgRegisterAddresses {
	return &MsgRegisterAddresses{
		Creator:             creator,
		AddressesToRegister: addressesToRegister,
	}
}

func (msg *MsgRegisterAddresses) Route() string {
	return RouterKey
}

func (msg *MsgRegisterAddresses) Type() string {
	return TypeMsgRegisterAddresses
}

func (msg *MsgRegisterAddresses) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterAddresses) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterAddresses) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.AddressesToRegister) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "no addresses to register")
	}

	for _, address := range msg.AddressesToRegister {
		_, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
		}
	}
	return nil
}
