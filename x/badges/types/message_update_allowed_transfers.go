package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateAllowedTransfers = "update_allowed_transfers"

var _ sdk.Msg = &MsgUpdateAllowedTransfers{}

func NewMsgUpdateAllowedTransfers(creator string, collectionId sdk.Uint, allowedTransfers []*TransferMapping) *MsgUpdateAllowedTransfers {
	return &MsgUpdateAllowedTransfers{
		Creator:             creator,
		CollectionId:        collectionId,
		AllowedTransfers: 	 allowedTransfers,
	}
}

func (msg *MsgUpdateAllowedTransfers) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAllowedTransfers) Type() string {
	return TypeMsgUpdateAllowedTransfers
}

func (msg *MsgUpdateAllowedTransfers) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAllowedTransfers) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAllowedTransfers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	for _, transferMapping := range msg.AllowedTransfers {
		if err := ValidateTransferMapping(*transferMapping); err != nil {
			return err
		}
	}

	if msg.CollectionId.IsZero() || msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidBadgeID, "collection id cannot be 0")
	}

	return nil
}
