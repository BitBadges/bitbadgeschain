package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSelfDestructBadge = "self_destruct_badge"

var _ sdk.Msg = &MsgSelfDestructBadge{}

func NewMsgSelfDestructBadge(creator string, badgeId uint64) *MsgSelfDestructBadge {
	return &MsgSelfDestructBadge{
		Creator: creator,
		BadgeId: badgeId,
	}
}

func (msg *MsgSelfDestructBadge) Route() string {
	return RouterKey
}

func (msg *MsgSelfDestructBadge) Type() string {
	return TypeMsgSelfDestructBadge
}

func (msg *MsgSelfDestructBadge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSelfDestructBadge) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSelfDestructBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
