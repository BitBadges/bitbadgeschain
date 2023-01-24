package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTransferBadge = "transfer_badge"

var _ sdk.Msg = &MsgTransferBadge{}

func NewMsgTransferBadge(creator string, collectionId uint64, from uint64, transfers []*Transfers) *MsgTransferBadge {
	return &MsgTransferBadge{
		Creator:      creator,
		CollectionId: collectionId,
		From:         from,
		Transfers:    transfers,
	}
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

	if msg.Transfers == nil || len(msg.Transfers) == 0 {
		return ErrInvalidLengthBalances
	}

	for _, transfer := range msg.Transfers {
		for _, balance := range transfer.Balances {
			if balance == nil {
				return ErrInvalidLengthBalances
			}

			if balance.Balance == 0 {
				return ErrAmountEqualsZero
			}

			err = ValidateRangesAreValid(balance.BadgeIds)
			if err != nil {
				return err
			}
		}

		if duplicateInArray(transfer.ToAddresses) {
			return ErrDuplicateAddresses
		}

		err = ValidateNoElementIsX(transfer.ToAddresses, msg.From)
		if err != nil {
			return err
		}
	}

	return nil
}
