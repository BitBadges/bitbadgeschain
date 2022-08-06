package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgFreezeAddress = "freeze_address"

var _ sdk.Msg = &MsgFreezeAddress{}

func NewMsgFreezeAddress(creator string, addresses NumberRange, badgeId uint64, add bool) *MsgFreezeAddress {
	return &MsgFreezeAddress{
		Creator:    creator,
		Addresses:  &addresses,
		BadgeId:    badgeId,
		Add:        add,
	}
}

func (msg *MsgFreezeAddress) Route() string {
	return RouterKey
}

func (msg *MsgFreezeAddress) Type() string {
	return TypeMsgFreezeAddress
}

func (msg *MsgFreezeAddress) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgFreezeAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFreezeAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
