package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgAddCustomData{}

const TypeMsgAddCustomData = "add_custom_data"

func NewMsgAddCustomData(creator string, data string) *MsgAddCustomData {
	return &MsgAddCustomData{
		Creator: creator,
		Data:    data,
	}
}

func (msg *MsgAddCustomData) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddCustomData) Route() string {
	return RouterKey
}

func (msg *MsgAddCustomData) Type() string {
	return TypeMsgAddCustomData
}

func (msg *MsgAddCustomData) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddCustomData) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Data) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequest, "data cannot be empty")
	}

	return nil
}
