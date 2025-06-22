package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateUserApprovals = "update_user_approved_transfers"

var _ sdk.Msg = &MsgUpdateUserApprovals{}

func NewMsgUpdateUserApprovals(creator string, creatorOverride string) *MsgUpdateUserApprovals {
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
	return msg.CheckAndCleanMsg(sdk.Context{}, false)
}

func (msg *MsgUpdateUserApprovals) CheckAndCleanMsg(ctx sdk.Context, canChangeValues bool) error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateUserIncomingApprovals(ctx, msg.IncomingApprovals, msg.Creator, canChangeValues); err != nil {
		return err
	}

	if err := ValidateUserOutgoingApprovals(ctx, msg.OutgoingApprovals, msg.Creator, canChangeValues); err != nil {
		return err
	}

	if msg.UserPermissions == nil {
		msg.UserPermissions = &UserPermissions{}
	}

	err = ValidateUserPermissions(msg.UserPermissions, canChangeValues)
	if err != nil {
		return err
	}

	return nil
}
