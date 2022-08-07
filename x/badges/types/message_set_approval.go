package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetApproval = "set_approval"

var _ sdk.Msg = &MsgSetApproval{}

func NewMsgSetApproval(creator string, amount uint64, address uint64, badgeId uint64, subbadgeRanges []*NumberRange, expirationTime uint64) *MsgSetApproval {
	return &MsgSetApproval{
		Creator:        creator,
		Amount:         amount,
		Address:        address,
		BadgeId:        badgeId,
		SubbadgeRanges: subbadgeRanges,
		ExpirationTime: expirationTime,
	}
}

func (msg *MsgSetApproval) Route() string {
	return RouterKey
}

func (msg *MsgSetApproval) Type() string {
	return TypeMsgSetApproval
}

func (msg *MsgSetApproval) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetApproval) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetApproval) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.SubbadgeRanges == nil {
		return ErrRangesIsNil
	}

	for _, subbadgeRange := range msg.SubbadgeRanges {
		if subbadgeRange == nil || subbadgeRange.Start > subbadgeRange.End {
			return ErrStartGreaterThanEnd
		}
	}
	return nil
}
