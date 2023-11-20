package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateAddressMappings = "create_address_mappings"

var _ sdk.Msg = &MsgCreateAddressMappings{}

func NewMsgCreateAddressMappings(creator string, addressMappings []*AddressMapping) *MsgCreateAddressMappings {
	return &MsgCreateAddressMappings{
		Creator:         creator,
		AddressMappings: addressMappings,
	}
}

func (msg *MsgCreateAddressMappings) Route() string {
	return RouterKey
}

func (msg *MsgCreateAddressMappings) Type() string {
	return TypeMsgCreateAddressMappings
}

func (msg *MsgCreateAddressMappings) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateAddressMappings) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAddressMappings) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	for _, mapping := range msg.AddressMappings {
		if err := ValidateAddressMapping(mapping); err != nil {
			return err
		}
	}

	return nil
}
