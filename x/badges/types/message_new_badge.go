package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewBadge = "new_badge"

var _ sdk.Msg = &MsgNewBadge{}

func NewMsgNewBadge(creator string, standard uint64, defaultSupply uint64, amountsToCreate []uint64, supplysToCreate []uint64, uriObject *UriObject, permissions uint64, freezeAddressRanges []*IdRange, bytesToStore []byte) *MsgNewBadge {
	return &MsgNewBadge{
		Creator:                 creator,
		Uri:                     uriObject,
		DefaultSubassetSupply:   defaultSupply,
		SubassetAmountsToCreate: amountsToCreate,
		SubassetSupplys:         supplysToCreate,
		FreezeAddressRanges:     freezeAddressRanges,
		ArbitraryBytes:          bytesToStore,
		Permissions:             permissions,
		Standard:                standard,
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

	err = ValidateNoElementIsX(msg.SubassetAmountsToCreate, 0)
	if err != nil {
		return err
	}

	err = ValidateNoElementIsX(msg.SubassetSupplys, 0)
	if err != nil {
		return err
	}

	err = ValidateRangesAreValid(msg.FreezeAddressRanges)
	if err != nil {
		return err
	}

	return nil
}
