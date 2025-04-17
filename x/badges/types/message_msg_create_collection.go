package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreateCollection = "msg_create_collection"

var _ sdk.Msg = &MsgCreateCollection{}

func NewMsgCreateCollection(creator string, creatorOverride string) *MsgCreateCollection {
	return &MsgCreateCollection{
		Creator:         creator,
		CreatorOverride: creatorOverride,
	}
}

func (msg *MsgCreateCollection) Route() string {
	return RouterKey
}

func (msg *MsgCreateCollection) Type() string {
	return TypeMsgCreateCollection
}

func (msg *MsgCreateCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateCollection) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CreatorOverride != "" {
		_, err = sdk.AccAddressFromBech32(msg.CreatorOverride)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator override address (%s)", err)
		}
	}

	return nil
}
