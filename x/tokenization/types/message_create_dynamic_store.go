package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreateDynamicStore = "msg_create_dynamic_store"

var _ sdk.Msg = &MsgCreateDynamicStore{}

func NewMsgCreateDynamicStore(creator string, defaultValue bool) *MsgCreateDynamicStore {
	return &MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: defaultValue,
	}
}

func (msg *MsgCreateDynamicStore) Route() string {
	return RouterKey
}

func (msg *MsgCreateDynamicStore) Type() string {
	return TypeMsgCreateDynamicStore
}

func (msg *MsgCreateDynamicStore) GetSigners() []sdk.AccAddress {
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateDynamicStore) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateDynamicStore) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return sdkerrors.Wrapf(ErrInvalidAddress, "creator address cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
