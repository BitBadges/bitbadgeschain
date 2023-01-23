package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRequestTransferManager = "request_transfer_manager"

var _ sdk.Msg = &MsgRequestTransferManager{}

func NewMsgRequestTransferManager(creator string, collectionId uint64, addRequest bool) *MsgRequestTransferManager {
	return &MsgRequestTransferManager{
		Creator: creator,
		CollectionId: collectionId,
		AddRequest: addRequest,
	}
}

func (msg *MsgRequestTransferManager) Route() string {
	return RouterKey
}

func (msg *MsgRequestTransferManager) Type() string {
	return TypeMsgRequestTransferManager
}

func (msg *MsgRequestTransferManager) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRequestTransferManager) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRequestTransferManager) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
