package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgDeleteDynamicStore = "msg_delete_dynamic_store"

var _ sdk.Msg = &MsgDeleteDynamicStore{}

func NewMsgDeleteDynamicStore(creator string, storeId sdkmath.Uint) *MsgDeleteDynamicStore {
	return &MsgDeleteDynamicStore{
		Creator: creator,
		StoreId: storeId,
	}
}

func (msg *MsgDeleteDynamicStore) Route() string {
	return RouterKey
}

func (msg *MsgDeleteDynamicStore) Type() string {
	return TypeMsgDeleteDynamicStore
}

func (msg *MsgDeleteDynamicStore) GetSigners() []sdk.AccAddress {
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteDynamicStore) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteDynamicStore) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return sdkerrors.Wrapf(ErrInvalidAddress, "creator address cannot be empty")
	}
	if msg.StoreId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "storeId cannot be zero")
	}
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
