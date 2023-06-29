package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgForkCollection = "fork_collection"

var _ sdk.Msg = &MsgForkCollection{}

func NewMsgForkCollection(creator string, collectionId sdkmath.Uint) *MsgForkCollection {
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
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
