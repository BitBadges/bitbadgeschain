package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgDeleteCollection = "delete_collection"

var _ sdk.Msg = &MsgDeleteCollection{}

func NewMsgDeleteCollection(creator string, collectionId sdkmath.Uint, creatorOverride string) *MsgDeleteCollection {
	return &MsgDeleteCollection{
		Creator:         creator,
		CollectionId:    collectionId,
		CreatorOverride: creatorOverride,
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
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsNil() || msg.CollectionId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection id")
	}

	if msg.CreatorOverride != "" {
		_, err = sdk.AccAddressFromBech32(msg.CreatorOverride)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator override address (%s)", err)
		}
	}
	return nil
}
