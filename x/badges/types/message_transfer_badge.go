package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTransferBadge = "transfer_badge"

var _ sdk.Msg = &MsgTransferBadge{}

func NewMsgTransferBadge(creator string, from uint64, toAddresses []uint64, amounts []uint64, badgeId uint64, subbadgeRanges []*IdRange, expirationTime uint64, cantCancelBeforeTime uint64) *MsgTransferBadge {
	return &MsgTransferBadge{
		Creator:              creator,
		From:                 from,
		ToAddresses:          toAddresses,
		Amounts:              amounts,
		BadgeId:              badgeId,
		SubbadgeRanges:       subbadgeRanges,
		ExpirationTime:       expirationTime,
		CantCancelBeforeTime: cantCancelBeforeTime,
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
	bz := AminoCdc.MustMarshalJSON(msg)
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

	if msg.ExpirationTime != 0 && msg.CantCancelBeforeTime > msg.ExpirationTime {
		return ErrCancelTimeIsGreaterThanExpirationTime
	}

	err = ValidateRangesAreValid(msg.SubbadgeRanges)
	if err != nil {
		return err
	}

	if duplicateInArray(msg.ToAddresses) {
		return ErrDuplicateAddresses
	}

	if duplicateInArray(msg.Amounts) {
		return ErrDuplicateAmounts
	}

	err = ValidateNoElementIsX(msg.ToAddresses, msg.From)
	if err != nil {
		return err
	}

	err = ValidateNoElementIsX(msg.Amounts, 0)
	if err != nil {
		return err
	}

	return nil
}
