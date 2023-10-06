package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateUserApprovals = "update_user_approved_transfers"

var _ sdk.Msg = &MsgUpdateUserApprovals{}

func NewMsgUpdateUserApprovals(creator string) *MsgUpdateUserApprovals {
	return &MsgUpdateUserApprovals{
		Creator: creator,
	}
}

func (msg *MsgUpdateUserApprovals) Route() string {
	return RouterKey
}

func (msg *MsgUpdateUserApprovals) Type() string {
	return TypeMsgUpdateUserApprovals
}

func (msg *MsgUpdateUserApprovals) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateUserApprovals) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateUserApprovals) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateUserIncomingApprovals(msg.IncomingApprovals, msg.Creator); err != nil {
		return err
	}

	if err := ValidateUserOutgoingApprovals(msg.OutgoingApprovals, msg.Creator); err != nil {
		return err
	}

	if msg.UserPermissions == nil {
		msg.UserPermissions = &UserPermissions{}
	}

	err = ValidateUserPermissions(msg.UserPermissions)
	if err != nil {
		return err
	}

	return nil
}
