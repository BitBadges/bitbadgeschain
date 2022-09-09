package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRevokeBadge = "revoke_badge"

var _ sdk.Msg = &MsgRevokeBadge{}

func NewMsgRevokeBadge(creator string, addresses []uint64, amounts []uint64, badgeId uint64, subbadgeRanges []*IdRange) *MsgRevokeBadge {
	return &MsgRevokeBadge{
		Creator:        creator,
		Addresses:      addresses,
		Amounts:        amounts,
		BadgeId:        badgeId,
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
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func duplicateInArray(arr []uint64) bool {
	visited := make(map[uint64]bool, 0)
	for i := 0; i < len(arr); i++ {
		if visited[arr[i]] {
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

	err = ValidateNoElementIsX(msg.Amounts, 0)
	if err != nil {
		return err
	}

	err = ValidateRangesAreValid(msg.SubbadgeRanges)
	if err != nil {
		return err
	}

	if duplicateInArray(msg.Addresses) {
		return ErrDuplicateAddresses
	}

	if duplicateInArray(msg.Amounts) {
		return ErrDuplicateAmounts
	}

	return nil
}
