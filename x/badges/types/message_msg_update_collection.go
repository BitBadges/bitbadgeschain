package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateCollection = "msg_update_collection"

var _ sdk.Msg = &MsgUpdateCollection{}

func NewMsgUpdateCollection(creator string, creatorOverride string) *MsgUpdateCollection {
	return &MsgUpdateCollection{
		Creator:         creator,
		CreatorOverride: creatorOverride,
	}
}

func (msg *MsgUpdateCollection) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCollection) Type() string {
	return TypeMsgUpdateCollection
}

func (msg *MsgUpdateCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCollection) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCollection) ValidateBasic() error {
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
