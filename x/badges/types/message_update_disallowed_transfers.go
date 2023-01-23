package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateDisallowedTransfers = "update_disallowed_transfers"

var _ sdk.Msg = &MsgUpdateDisallowedTransfers{}

func NewMsgUpdateDisallowedTransfers(creator string, collectionId uint64, disallowedTransfers []*TransferMapping) *MsgUpdateDisallowedTransfers {
	return &MsgUpdateDisallowedTransfers{
		Creator:       creator,
		CollectionId: collectionId,
		DisallowedTransfers: disallowedTransfers,
	}
}

func (msg *MsgUpdateDisallowedTransfers) Route() string {
	return RouterKey
}

func (msg *MsgUpdateDisallowedTransfers) Type() string {
	return TypeMsgUpdateDisallowedTransfers
}

func (msg *MsgUpdateDisallowedTransfers) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateDisallowedTransfers) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateDisallowedTransfers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
