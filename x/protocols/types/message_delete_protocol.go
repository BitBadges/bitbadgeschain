package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeleteProtocol = "delete_protocol"

var _ sdk.Msg = &MsgDeleteProtocol{}

func NewMsgDeleteProtocol(creator string, name string) *MsgDeleteProtocol {
	return &MsgDeleteProtocol{
		Creator: creator,
		Name: name,
	}
}

func (msg *MsgDeleteProtocol) Route() string {
	return RouterKey
}

func (msg *MsgDeleteProtocol) Type() string {
	return TypeMsgDeleteProtocol
}

func (msg *MsgDeleteProtocol) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteProtocol) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteProtocol) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}
	
	return nil
}
