package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateDynamicStore = "msg_update_dynamic_store"

var _ sdk.Msg = &MsgUpdateDynamicStore{}

func NewMsgUpdateDynamicStore(creator string, storeId sdkmath.Uint, defaultValue bool) *MsgUpdateDynamicStore {
	return &MsgUpdateDynamicStore{
		Creator:       creator,
		StoreId:       storeId,
		DefaultValue:  defaultValue,
		GlobalEnabled: true, // Default to enabled for backward compatibility
	}
}

// NewMsgUpdateDynamicStoreWithGlobalEnabled creates a new MsgUpdateDynamicStore with explicit globalEnabled
func NewMsgUpdateDynamicStoreWithGlobalEnabled(creator string, storeId sdkmath.Uint, defaultValue bool, globalEnabled bool) *MsgUpdateDynamicStore {
	return &MsgUpdateDynamicStore{
		Creator:       creator,
		StoreId:       storeId,
		DefaultValue:  defaultValue,
		GlobalEnabled: globalEnabled,
	}
}

func (msg *MsgUpdateDynamicStore) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDynamicStore) Type() string {
	return TypeMsgUpdateDynamicStore
}

func (msg *MsgUpdateDynamicStore) GetSigners() []sdk.AccAddress {
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDynamicStore) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDynamicStore) ValidateBasic() error {
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
