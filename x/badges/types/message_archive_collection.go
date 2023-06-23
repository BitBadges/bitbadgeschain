package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgArchiveCollection = "archive_collection"

var _ sdk.Msg = &MsgArchiveCollection{}

func NewMsgArchiveCollection(creator string, collectionId sdk.Uint) *MsgArchiveCollection {
	return &MsgArchiveCollection{
		Creator:      creator,
		CollectionId: collectionId,
	}
}

func (msg *MsgArchiveCollection) Route() string {
	return RouterKey
}

func (msg *MsgArchiveCollection) Type() string {
	return TypeMsgArchiveCollection
}

func (msg *MsgArchiveCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgArchiveCollection) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgArchiveCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
