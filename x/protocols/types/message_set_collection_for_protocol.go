package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetCollectionForProtocol = "set_collection_for_protocol"

var _ sdk.Msg = &MsgSetCollectionForProtocol{}

func NewMsgSetCollectionForProtocol(creator string, name string, collectionId sdk.Uint) *MsgSetCollectionForProtocol {
	return &MsgSetCollectionForProtocol{
		Creator: creator,
		Name: name,
		CollectionId: collectionId,
	}
}

func (msg *MsgSetCollectionForProtocol) Route() string {
	return RouterKey
}

func (msg *MsgSetCollectionForProtocol) Type() string {
	return TypeMsgSetCollectionForProtocol
}

func (msg *MsgSetCollectionForProtocol) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetCollectionForProtocol) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetCollectionForProtocol) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}

	if msg.CollectionId.IsNil() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "collectionId cannot be empty")
	}

	return nil
}
