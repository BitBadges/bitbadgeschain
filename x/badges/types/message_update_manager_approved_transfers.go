package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateManagerApprovedTransfers = "update_manager_approved_transfers"

var _ sdk.Msg = &MsgUpdateManagerApprovedTransfers{}

func NewMsgUpdateManagerApprovedTransfers(creator string, collectionId uint64, managerApprovedTransfers []*TransferMapping) *MsgUpdateManagerApprovedTransfers {
	return &MsgUpdateManagerApprovedTransfers{
		Creator:                  creator,
		CollectionId:             collectionId,
		ManagerApprovedTransfers: managerApprovedTransfers,
	}
}

func (msg *MsgUpdateManagerApprovedTransfers) Route() string {
	return RouterKey
}

func (msg *MsgUpdateManagerApprovedTransfers) Type() string {
	return TypeMsgUpdateManagerApprovedTransfers
}

func (msg *MsgUpdateManagerApprovedTransfers) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateManagerApprovedTransfers) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateManagerApprovedTransfers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
