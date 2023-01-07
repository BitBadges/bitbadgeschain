package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewBadge = "new_badge"

var _ sdk.Msg = &MsgNewBadge{}

func NewMsgNewBadge(creator string, standard uint64, defaultSupply uint64, subassetsToCreate []*SubassetSupplyAndAmount, uriObject *UriObject, permissions uint64, freezeAddressRanges []*IdRange, bytesToStore string, whitelistedRecipients []*WhitelistMintInfo) *MsgNewBadge {
	return &MsgNewBadge{
		Creator:                 	creator,
		Uri:                     	uriObject,
		DefaultSubassetSupply:   	defaultSupply,
		SubassetSupplysAndAmounts:	subassetsToCreate,
		FreezeAddressRanges:    	freezeAddressRanges,
		ArbitraryBytes:          	bytesToStore,
		Permissions:             	permissions,
		Standard:                	standard,
		WhitelistedRecipients:   	whitelistedRecipients,
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
	bz := AminoCdc.MustMarshalJSON(msg)
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

	err = ValidateRangesAreValid(msg.FreezeAddressRanges)
	if err != nil {
		return err
	}

	return nil
}
