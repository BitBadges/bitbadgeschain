package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateUserApprovedTransfers = "update_user_approved_transfers"

var _ sdk.Msg = &MsgUpdateUserApprovedTransfers{}

func NewMsgUpdateUserApprovedTransfers(creator string) *MsgUpdateUserApprovedTransfers {
	return &MsgUpdateUserApprovedTransfers{
		Creator: creator,
	}
}

func (msg *MsgUpdateUserApprovedTransfers) Route() string {
	return RouterKey
}

func (msg *MsgUpdateUserApprovedTransfers) Type() string {
	return TypeMsgUpdateUserApprovedTransfers
}

func (msg *MsgUpdateUserApprovedTransfers) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateUserApprovedTransfers) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateUserApprovedTransfers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
