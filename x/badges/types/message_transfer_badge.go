package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTransferBadge = "transfer_badge"

var _ sdk.Msg = &MsgTransferBadge{}

func NewMsgTransferBadge(creator string, from uint64, toAddresses []uint64, amounts []uint64, badgeId uint64, subbadgeRanges []*IdRange, expirationTime uint64) *MsgTransferBadge {
	return &MsgTransferBadge{
		Creator:        creator,
		From:           from,
		ToAddresses:    toAddresses,
		Amounts:        amounts,
		BadgeId:        badgeId,
		SubbadgeRanges: subbadgeRanges,
		ExpirationTime: expirationTime,
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

	if msg.SubbadgeRanges == nil {
		return ErrRangesIsNil
	}

	for _, subbadgeRange := range msg.SubbadgeRanges {
		if subbadgeRange == nil || subbadgeRange.Start > subbadgeRange.End {
			return ErrStartGreaterThanEnd
		}
	}

	if duplicateInArray(msg.ToAddresses) {
		return ErrDuplicateAddresses
	}

	for _, toAddress := range msg.ToAddresses {
		//Can't send to same address
		if toAddress == msg.From {
			return ErrSenderAndReceiverSame
		}
	}

	for _, amount := range msg.Amounts {
		if amount == uint64(0) {
			return ErrAmountEqualsZero
		}
	}
	return nil
}
