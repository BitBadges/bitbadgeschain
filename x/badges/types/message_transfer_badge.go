package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgTransferTokens = "transfer_badge"

var _ sdk.Msg = &MsgTransferTokens{}

func NewMsgTransferTokens(creator string, collectionId sdkmath.Uint, transfers []*Transfer, creatorOverride string) *MsgTransferTokens {
	return &MsgTransferTokens{
		Creator:      creator,
		CollectionId: collectionId,
		Transfers:    transfers,
	}
}

func (msg *MsgTransferTokens) Route() string {
	return RouterKey
}

func (msg *MsgTransferTokens) Type() string {
	return TypeMsgTransferTokens
}

func (msg *MsgTransferTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTransferTokens) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTransferTokens) CheckAndCleanMsg(ctx sdk.Context, canChangeValues bool) error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Transfers == nil || len(msg.Transfers) == 0 {
		return sdkerrors.Wrapf(ErrInvalidTransfers, "transfers cannot be empty")
	}

	for _, transfer := range msg.Transfers {
		err = ValidateTransfer(ctx, transfer, canChangeValues)
		if err != nil {
			return err
		}
	}

	return nil
}

func (msg *MsgTransferTokens) ValidateBasic() error {
	return msg.CheckAndCleanMsg(sdk.Context{}, false)
}
