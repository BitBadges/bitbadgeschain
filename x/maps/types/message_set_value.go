package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetValue = "set_value"

var _ sdk.Msg = &MsgSetValue{}

func NewMsgSetValue(creator string, mapId string, key string, value string, options *SetOptions) *MsgSetValue {
	return &MsgSetValue{
		Creator:                  creator,
		MapId:                    mapId,
		Key:                      key,
		Value:                    value,
		Options:                  options,
	}
}


func (msg *MsgSetValue) Route() string {
	return RouterKey
}

func (msg *MsgSetValue) Type() string {
	return TypeMsgSetValue
}

func (msg *MsgSetValue) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetValue) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetValue) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.MapId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "map ID cannot be empty")
	}

	if len(msg.Key) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "key cannot be empty")
	}

	return nil
}
