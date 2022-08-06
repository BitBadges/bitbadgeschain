package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRevokeBadge = "revoke_badge"

var _ sdk.Msg = &MsgRevokeBadge{}

func NewMsgRevokeBadge(creator string, addresses []uint64, amounts []uint64, badgeId uint64, subbadgeRanges []*NumberRange) *MsgRevokeBadge {
	return &MsgRevokeBadge{
		Creator:    creator,
		Addresses:  addresses,
		Amounts:    amounts,
		BadgeId:    badgeId,
		SubbadgeRanges: subbadgeRanges,
	}
}

func (msg *MsgRevokeBadge) Route() string {
	return RouterKey
}

func (msg *MsgRevokeBadge) Type() string {
	return TypeMsgRevokeBadge
}

func (msg *MsgRevokeBadge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRevokeBadge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func duplicateInArray(arr []uint64) bool {
	visited := make(map[uint64]bool, 0)
	for i:=0; i<len(arr); i++{
	   if visited[arr[i]] == true {
		  return true
	   } else {
		  visited[arr[i]] = true
	   }
	}
	return false
}
 

func (msg *MsgRevokeBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Amounts) != len(msg.Addresses) {
		return ErrInvalidAmountsAndAddressesLength
	}

	for _, amount := range msg.Amounts {
		if amount == 0 {
			return ErrAmountEqualsZero
		}
	}

	if msg.SubbadgeRanges == nil {
		return ErrRangesIsNil
	}

	for _, subbadgeRange := range msg.SubbadgeRanges {
		if subbadgeRange == nil || subbadgeRange.Start > subbadgeRange.End {
			return ErrStartGreaterThanEnd
		}
	}

	if duplicateInArray(msg.Addresses) {
		return ErrDuplicateAddresses
	}

	return nil
}
