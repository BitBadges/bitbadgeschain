package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeleteMap = "delete_map"

var _ sdk.Msg = &MsgDeleteMap{}

func NewMsgDeleteMap(creator string, mapId string) *MsgDeleteMap {
	return &MsgDeleteMap{
		Creator:                  creator,
		MapId:                    mapId,
	}
}


func (msg *MsgDeleteMap) Route() string {
	return RouterKey
}

func (msg *MsgDeleteMap) Type() string {
	return TypeMsgDeleteMap
}

func (msg *MsgDeleteMap) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteMap) GetSignBytes() []byte {
	bz := Amino.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteMap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.MapId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "map ID cannot be empty")
	}
	
	return nil
}
