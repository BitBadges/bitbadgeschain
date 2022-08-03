package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgHandlePendingTransfer = "handle_pending_transfer"

var _ sdk.Msg = &MsgHandlePendingTransfer{}

func NewMsgHandlePendingTransfer(creator string, accept bool, badgeId uint64, subbadgeId uint64, startingNonce uint64, endingNonce uint64) *MsgHandlePendingTransfer {
	return &MsgHandlePendingTransfer{
		Creator:       creator,
		Accept:        accept,
		BadgeId:       badgeId,
		SubbadgeId:    subbadgeId,
		StartingNonce: startingNonce,
		EndingNonce:   endingNonce,
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
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgHandlePendingTransfer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
