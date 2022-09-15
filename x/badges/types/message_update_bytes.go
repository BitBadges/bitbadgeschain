package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateBytes = "update_bytes"

var _ sdk.Msg = &MsgUpdateBytes{}

func NewMsgUpdateBytes(creator string, badgeId uint64, newBytes string) *MsgUpdateBytes {
	return &MsgUpdateBytes{
		Creator:  creator,
		BadgeId:  badgeId,
		NewBytes: newBytes,
	}
}

func (msg *MsgUpdateBytes) Route() string {
	return RouterKey
}

func (msg *MsgUpdateBytes) Type() string {
	return TypeMsgUpdateBytes
}

func (msg *MsgUpdateBytes) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateBytes) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateBytes) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateBytes(msg.NewBytes); err != nil {
		return err
	}

	return nil
}
