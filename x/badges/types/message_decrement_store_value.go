package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgDecrementStoreValue = "msg_decrement_store_value"

var _ sdk.Msg = &MsgDecrementStoreValue{}

func NewMsgDecrementStoreValue(creator string, storeId sdkmath.Uint, address string, amount sdkmath.Uint, setToZeroOnUnderflow bool) *MsgDecrementStoreValue {
	return &MsgDecrementStoreValue{
		Creator:              creator,
		StoreId:              storeId,
		Address:              address,
		Amount:               amount,
		SetToZeroOnUnderflow: setToZeroOnUnderflow,
	}
}

func (msg *MsgDecrementStoreValue) Route() string {
	return RouterKey
}

func (msg *MsgDecrementStoreValue) Type() string {
	return TypeMsgDecrementStoreValue
}

func (msg *MsgDecrementStoreValue) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDecrementStoreValue) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDecrementStoreValue) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return sdkerrors.Wrapf(ErrInvalidAddress, "creator address cannot be empty")
	}
	if msg.StoreId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "storeId cannot be zero")
	}
	if len(msg.Address) == 0 {
		return sdkerrors.Wrapf(ErrInvalidAddress, "address cannot be empty")
	}
	if msg.Amount.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "amount cannot be zero")
	}
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if err := ValidateAddress(msg.Address, false); err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid address (%s)", err)
	}
	return nil
}
