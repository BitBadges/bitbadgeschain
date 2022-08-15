package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewSubBadge = "new_sub_badge"

var _ sdk.Msg = &MsgNewSubBadge{}

func NewMsgNewSubBadge(creator string, badgeId uint64, supplys []uint64, amountsToCreate []uint64) *MsgNewSubBadge {
	return &MsgNewSubBadge{
		Creator:         creator,
		BadgeId:         badgeId,
		Supplys:         supplys,
		AmountsToCreate: amountsToCreate,
	}
}

func (msg *MsgNewSubBadge) Route() string {
	return RouterKey
}

func (msg *MsgNewSubBadge) Type() string {
	return TypeMsgNewSubBadge
}

func (msg *MsgNewSubBadge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgNewSubBadge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNewSubBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Supplys) != len(msg.AmountsToCreate) {
		return ErrInvalidSupplyAndAmounts
	}

	err = ValidateNoElementIsX(msg.Supplys, 0)
	if err != nil {
		return err
	}

	err = ValidateNoElementIsX(msg.AmountsToCreate, 0)
	if err != nil {
		return err
	}

	return nil
}
