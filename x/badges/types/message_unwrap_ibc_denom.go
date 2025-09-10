package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUnwrapIBCDenom = "msg_unwrap_ibc_denom"

var _ sdk.Msg = &MsgUnwrapIBCDenom{}

func NewMsgUnwrapIBCDenom(creator string, collectionId Uint, amount *sdk.Coin, overrideTokenId string) *MsgUnwrapIBCDenom {
	return &MsgUnwrapIBCDenom{
		Creator:         creator,
		CollectionId:    collectionId,
		Amount:          amount,
		OverrideTokenId: overrideTokenId,
	}
}

func (msg *MsgUnwrapIBCDenom) Route() string {
	return RouterKey
}

func (msg *MsgUnwrapIBCDenom) Type() string {
	return TypeMsgUnwrapIBCDenom
}

func (msg *MsgUnwrapIBCDenom) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUnwrapIBCDenom) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnwrapIBCDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "collection ID cannot be nil")
	}

	if msg.CollectionId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "collection ID cannot be zero")
	}

	if msg.Amount.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "amount cannot be nil")
	}

	if msg.Amount.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "amount cannot be zero")
	}

	if msg.Amount.Denom == "" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "denom cannot be empty")
	}

	// Validate override token ID if provided
	if msg.OverrideTokenId != "" {
		// Basic validation - could be enhanced with more specific rules
		if len(msg.OverrideTokenId) == 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "override token ID cannot be empty string")
		}
	}

	return nil
}
