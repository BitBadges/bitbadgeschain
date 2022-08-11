package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewBadge = "new_badge"

var _ sdk.Msg = &MsgNewBadge{}

func NewMsgNewBadge(creator string, uriObject UriObject, permissions uint64,  bytesToStore []byte, defaultSupply uint64, amountsToCreate []uint64, supplysToCreate []uint64, freezeAddressRanges []*IdRange) *MsgNewBadge {
	return &MsgNewBadge{
		Creator:                 creator,
		Uri:                     &uriObject,
		DefaultSubassetSupply:   defaultSupply,
		SubassetAmountsToCreate: amountsToCreate,
		SubassetSupplys:         supplysToCreate,
		FreezeAddressRanges:     freezeAddressRanges,
		ArbitraryBytes:          bytesToStore,
		Permissions: 		   	 permissions,
	}
}

func (msg *MsgNewBadge) Route() string {
	return RouterKey
}

func (msg *MsgNewBadge) Type() string {
	return TypeMsgNewBadge
}

func (msg *MsgNewBadge) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgNewBadge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNewBadge) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	//Validate well-formedness of the message entries
	if err := ValidateURI(*msg.Uri); err != nil {
		return err
	}


	if err := ValidatePermissions(msg.Permissions); err != nil {
		return err
	}

	if err := ValidateBytes(msg.ArbitraryBytes); err != nil {
		return err
	}

	if len(msg.SubassetAmountsToCreate) != len(msg.SubassetSupplys) {
		return ErrInvalidSupplyAndAmounts
	}

	for i, _ := range msg.SubassetSupplys {
		if msg.SubassetAmountsToCreate[i] == 0 {
			return ErrAmountEqualsZero
		}
	}

	if msg.FreezeAddressRanges == nil {
		return ErrRangesIsNil
	}

	for _, subbadgeRange := range msg.FreezeAddressRanges {
		if subbadgeRange == nil || subbadgeRange.Start > subbadgeRange.End {
			return ErrStartGreaterThanEnd
		}
	}

	return nil
}
