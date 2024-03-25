package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgGlobalArchive = "global_archive"

var _ sdk.Msg = &MsgGlobalArchive{}

func NewMsgGlobalArchive(creator string, archive bool) *MsgGlobalArchive {
	return &MsgGlobalArchive{
		Creator: creator,
		Archive: archive,
	}
}

func (msg *MsgGlobalArchive) Route() string {
	return RouterKey
}

func (msg *MsgGlobalArchive) Type() string {
	return TypeMsgGlobalArchive
}

func (msg *MsgGlobalArchive) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgGlobalArchive) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgGlobalArchive) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
