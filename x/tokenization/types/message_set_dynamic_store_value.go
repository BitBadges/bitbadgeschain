package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetDynamicStoreValue = "msg_set_dynamic_store_value"

var _ sdk.Msg = &MsgSetDynamicStoreValue{}

func NewMsgSetDynamicStoreValue(creator string, storeId sdkmath.Uint, address string, value bool) *MsgSetDynamicStoreValue {
	return &MsgSetDynamicStoreValue{
		Creator: creator,
		StoreId: storeId,
		Address: address,
		Value:   value,
	}
}

func (msg *MsgSetDynamicStoreValue) Route() string {
	return RouterKey
}

func (msg *MsgSetDynamicStoreValue) Type() string {
	return TypeMsgSetDynamicStoreValue
}

func (msg *MsgSetDynamicStoreValue) GetSigners() []sdk.AccAddress {
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetDynamicStoreValue) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetDynamicStoreValue) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return sdkerrors.Wrapf(ErrInvalidAddress, "creator address cannot be empty")
	}
	if msg.StoreId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "storeId cannot be zero")
	}
	if len(msg.Address) == 0 {
		return sdkerrors.Wrapf(ErrInvalidAddress, "address cannot be empty")
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
