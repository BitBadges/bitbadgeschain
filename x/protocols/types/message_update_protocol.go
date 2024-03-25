package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	badgetypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

const TypeMsgUpdateProtocol = "update_protocol"

var _ sdk.Msg = &MsgUpdateProtocol{}

func NewMsgUpdateProtocol(creator string, name string, uri string, customData string, isFrozen bool) *MsgUpdateProtocol {
	return &MsgUpdateProtocol{
		Creator:    creator,
		Name:       name,
		Uri:        uri,
		CustomData: customData,
		IsFrozen:   isFrozen,
	}
}

func (msg *MsgUpdateProtocol) Route() string {
	return RouterKey
}

func (msg *MsgUpdateProtocol) Type() string {
	return TypeMsgUpdateProtocol
}

func (msg *MsgUpdateProtocol) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateProtocol) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateProtocol) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}

	if badgetypes.ValidateURI(msg.Uri) != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "uri cannot be invalid")
	}

	return nil
}
