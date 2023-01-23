package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTransferManager = "transfer_manager"

var _ sdk.Msg = &MsgTransferManager{}

func NewMsgTransferManager(creator string, collectionId uint64, address uint64) *MsgTransferManager {
	return &MsgTransferManager{
		Creator: creator,
		CollectionId: collectionId,
		Address: address,
	}
}

func (msg *MsgTransferManager) Route() string {
	return RouterKey
}

func (msg *MsgTransferManager) Type() string {
	return TypeMsgTransferManager
}

func (msg *MsgTransferManager) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferManager) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferManager) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
