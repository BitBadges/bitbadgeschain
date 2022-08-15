package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateUris = "update_uris"

var _ sdk.Msg = &MsgUpdateUris{}

func NewMsgUpdateUris(creator string, badgeId uint64, uri *UriObject) *MsgUpdateUris {
	return &MsgUpdateUris{
		Creator:     creator,
		BadgeId:     badgeId,
		Uri:         uri,
	}
}

func (msg *MsgUpdateUris) Route() string {
	return RouterKey
}

func (msg *MsgUpdateUris) Type() string {
	return TypeMsgUpdateUris
}

func (msg *MsgUpdateUris) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateUris) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateUris) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	//Validate well-formedness of the message entries
	if err := ValidateURI(*msg.Uri); err != nil {
		return err
	}

	return nil
}
