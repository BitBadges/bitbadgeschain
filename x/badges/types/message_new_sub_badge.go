package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewSubBadge = "new_sub_badge"

var _ sdk.Msg = &MsgNewSubBadge{}

func NewMsgNewSubBadge(creator string, id uint64, supplys []uint64, amountsToCreate []uint64) *MsgNewSubBadge {
	return &MsgNewSubBadge{
		Creator:         creator,
		Id:              id,
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

	for i, _ := range msg.Supplys {
		if msg.AmountsToCreate[i] == 0 {
			return ErrAmountEqualsZero
		}
	}
	return nil
}
