package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgTransferBadges = "transfer_badge"

var _ sdk.Msg = &MsgTransferBadges{}

func NewMsgTransferBadges(creator string, collectionId sdkmath.Uint, transfers []*Transfer) *MsgTransferBadges {
	return &MsgTransferBadges{
		Creator:      creator,
		CollectionId: collectionId,
		Transfers:    transfers,
	}
}

func (msg *MsgTransferBadges) Route() string {
	return RouterKey
}

func (msg *MsgTransferBadges) Type() string {
	return TypeMsgTransferBadges
}

func (msg *MsgTransferBadges) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferBadges) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferBadges) CheckAndCleanMsg(canChangeValues bool) error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Transfers == nil || len(msg.Transfers) == 0 {
		return sdkerrors.Wrapf(ErrInvalidTransfers, "transfers cannot be empty")
	}

	for _, transfer := range msg.Transfers {
		err = ValidateTransfer(transfer, canChangeValues)
		if err != nil {
			return err
		}
	}

	return nil
}

func (msg *MsgTransferBadges) ValidateBasic() error {
	return msg.CheckAndCleanMsg(false)
}
