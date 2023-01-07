package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgHandlePendingTransfer = "handle_pending_transfer"

var _ sdk.Msg = &MsgHandlePendingTransfer{}

func NewMsgHandlePendingTransfer(creator string, badgeId uint64, nonceRanges []*IdRange, actions []uint64) *MsgHandlePendingTransfer {
	return &MsgHandlePendingTransfer{
		Creator:        creator,
		BadgeId:        badgeId,
		NonceRanges:    nonceRanges,
		Actions:        actions,
	}
}

func (msg *MsgHandlePendingTransfer) Route() string {
	return RouterKey
}

func (msg *MsgHandlePendingTransfer) Type() string {
	return TypeMsgHandlePendingTransfer
}

func (msg *MsgHandlePendingTransfer) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgHandlePendingTransfer) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgHandlePendingTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	err = ValidateRangesAreValid(msg.NonceRanges)
	if err != nil {
		return err
	}

	err = ValidateActionsAreValid(msg.Actions, len(msg.NonceRanges))
	if err != nil {
		return err
	}
	return nil
}
