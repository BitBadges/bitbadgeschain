package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdatePermissions = "update_permissions"

var _ sdk.Msg = &MsgUpdatePermissions{}

func NewMsgUpdatePermissions(creator string, badgeId uint64, permissions uint64) *MsgUpdatePermissions {
	return &MsgUpdatePermissions{
		Creator:     creator,
		BadgeId:     badgeId,
		Permissions: permissions,
	}
}

func (msg *MsgUpdatePermissions) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePermissions) Type() string {
	return TypeMsgUpdatePermissions
}

func (msg *MsgUpdatePermissions) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePermissions) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePermissions) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
