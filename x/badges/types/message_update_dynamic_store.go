package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateDynamicStore = "msg_update_dynamic_store"

var _ sdk.Msg = &MsgUpdateDynamicStore{}

func NewMsgUpdateDynamicStore(creator string, storeId sdkmath.Uint, defaultValue sdkmath.Uint) *MsgUpdateDynamicStore {
	return &MsgUpdateDynamicStore{
		Creator:      creator,
		StoreId:      storeId,
		DefaultValue: defaultValue,
	}
}

func (msg *MsgUpdateDynamicStore) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDynamicStore) Type() string {
	return TypeMsgUpdateDynamicStore
}

func (msg *MsgUpdateDynamicStore) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
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
