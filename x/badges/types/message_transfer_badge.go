package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTransferBadge = "transfer_badge"

var _ sdk.Msg = &MsgTransferBadge{}

func NewMsgTransferBadge(creator string, from uint64, to uint64, amount uint64, badgeId uint64, subbadgeId uint64) *MsgTransferBadge {
	return &MsgTransferBadge{
		Creator:    creator,
		From:       from,
		To:         to,
		Amount:     amount,
		BadgeId:    badgeId,
		SubbadgeId: subbadgeId,
	}
}

func (msg *MsgTransferBadge) Route() string {
	return RouterKey
}

func (msg *MsgTransferBadge) Type() string {
	return TypeMsgTransferBadge
}

func (msg *MsgTransferBadge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferBadge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	//Can't send to same address
	if msg.To == msg.From {
		return ErrSenderAndReceiverSame
	}
	return nil
}
