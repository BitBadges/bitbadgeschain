package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	badgetypes "github.com/bitbadges/bitbadgeschain/x/badges/types"
)

const TypeMsgCreateProtocol = "create_protocol"

var _ sdk.Msg = &MsgCreateProtocol{}

func NewMsgCreateProtocol(creator string, name string, uri string, customData string) *MsgCreateProtocol {
	return &MsgCreateProtocol{
		Creator:    creator,
		Name:       name,
		Uri:        uri,
		CustomData: customData,
	}
}

func (msg *MsgCreateProtocol) Route() string {
	return RouterKey
}

func (msg *MsgCreateProtocol) Type() string {
	return TypeMsgCreateProtocol
}

func (msg *MsgCreateProtocol) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateProtocol) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateProtocol) ValidateBasic() error {
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
