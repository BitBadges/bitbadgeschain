package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewBadge = "new_badge"

var _ sdk.Msg = &MsgNewBadge{}

func NewMsgNewBadge(creator string, uri string, permissions uint64, freezeAddressesDigest string, subassetUris string) *MsgNewBadge {
	return &MsgNewBadge{
		Creator:               creator,
		Uri:                   uri,
		Permissions:           permissions,
		FreezeAddressesDigest: freezeAddressesDigest,
		SubassetUris:          subassetUris,
	}
}

func (msg *MsgNewBadge) Route() string {
	return RouterKey
}

func (msg *MsgNewBadge) Type() string {
	return TypeMsgNewBadge
}

func (msg *MsgNewBadge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgNewBadge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNewBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	//Validate well-formedness of the message entries
	if err := ValidateURI(msg.Uri); err != nil {
		return err
	}

	if err := ValidateURI(msg.SubassetUris); err != nil {
		return err
	}

	if err := ValidatePermissions(msg.Permissions); err != nil {
		return err
	}

	//TODO: Validate freeze digest string

	return nil
}
