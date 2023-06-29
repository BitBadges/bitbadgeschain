package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgTransferBadge = "transfer_badge"

var _ sdk.Msg = &MsgTransferBadge{}

func NewMsgTransferBadge(creator string, collectionId sdk.Uint, transfers []*Transfer) *MsgTransferBadge {
	for _, transfer := range transfers {
		for _, balance := range transfer.Balances {
			balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
		}
	}

	return &MsgTransferBadge{
		Creator:      creator,
		CollectionId: collectionId,
		Transfers:    transfers,
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
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Transfers == nil || len(msg.Transfers) == 0 {
		return ErrInvalidLengthBalances
	}

	for _, transfer := range msg.Transfers {
		err = ValidateTransfer(transfer)
		if err != nil {
			return err
		}
	}

	return nil
}
