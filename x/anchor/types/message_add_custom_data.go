package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddCustomData = "add_custom_data"

var _ sdk.Msg = &MsgAddCustomData{}

func NewMsgAddCustomData(creator string, data string) *MsgAddCustomData {
	return &MsgAddCustomData{
		Creator: creator,
		Data:    data,
	}
}

func (msg *MsgAddCustomData) Route() string {
	return RouterKey
}

func (msg *MsgAddCustomData) Type() string {
	return TypeMsgAddCustomData
}

func (msg *MsgAddCustomData) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddCustomData) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddCustomData) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Data) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "data cannot be empty")
	}
	
	return nil
}
