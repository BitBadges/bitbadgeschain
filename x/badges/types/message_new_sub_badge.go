package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewSubBadge = "new_sub_badge"

var _ sdk.Msg = &MsgNewSubBadge{}

func NewMsgNewSubBadge(creator string, badgeId uint64, supplysAndAmounts []*SubassetSupplyAndAmount) *MsgNewSubBadge {
	return &MsgNewSubBadge{
		Creator:         creator,
		BadgeId:         badgeId,
		SubassetSupplysAndAmounts: supplysAndAmounts,
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
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNewSubBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	amounts := make([]uint64, len(msg.SubassetSupplysAndAmounts))
	supplys := make([]uint64, len(msg.SubassetSupplysAndAmounts))
	for i, subasset := range msg.SubassetSupplysAndAmounts {
		amounts[i] = subasset.Amount
		supplys[i] = subasset.Supply
	}

	err = ValidateNoElementIsX(amounts, 0)
	if err != nil {
		return err
	}

	err = ValidateNoElementIsX(supplys, 0)
	if err != nil {
		return err
	}

	return nil
}
