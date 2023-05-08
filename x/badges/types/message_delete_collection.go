package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeleteCollection = "delete_collection"

var _ sdk.Msg = &MsgDeleteCollection{}

func NewMsgDeleteCollection(creator string, collectionId sdk.Uint) *MsgDeleteCollection {
	return &MsgDeleteCollection{
		Creator:      creator,
		CollectionId: collectionId,
	}
}

func (msg *MsgDeleteCollection) Route() string {
	return RouterKey
}

func (msg *MsgDeleteCollection) Type() string {
	return TypeMsgDeleteCollection
}

func (msg *MsgDeleteCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteCollection) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsZero() || msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid collection id")
	}
	
	return nil
}
