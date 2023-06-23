package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgForkCollection = "fork_collection"

var _ sdk.Msg = &MsgForkCollection{}

func NewMsgForkCollection(creator string, collectionId sdk.Uint) *MsgForkCollection {
	return &MsgForkCollection{
		Creator:            creator,
		ParentCollectionId: collectionId,
	}
}

func (msg *MsgForkCollection) Route() string {
	return RouterKey
}

func (msg *MsgForkCollection) Type() string {
	return TypeMsgForkCollection
}

func (msg *MsgForkCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgForkCollection) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgForkCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
