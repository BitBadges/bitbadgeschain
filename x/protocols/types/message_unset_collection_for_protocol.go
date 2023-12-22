package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnsetCollectionForProtocol = "unset_collection_for_protocol"

var _ sdk.Msg = &MsgUnsetCollectionForProtocol{}

func NewMsgUnsetCollectionForProtocol(creator string, name string) *MsgUnsetCollectionForProtocol {
	return &MsgUnsetCollectionForProtocol{
		Creator: creator,
		Name:    name,
	}
}

func (msg *MsgUnsetCollectionForProtocol) Route() string {
	return RouterKey
}

func (msg *MsgUnsetCollectionForProtocol) Type() string {
	return TypeMsgUnsetCollectionForProtocol
}

func (msg *MsgUnsetCollectionForProtocol) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUnsetCollectionForProtocol) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnsetCollectionForProtocol) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Name) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}

	return nil
}
