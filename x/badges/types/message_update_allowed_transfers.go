package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateCollectionApprovedTransfers = "update_allowed_transfers"

var _ sdk.Msg = &MsgUpdateCollectionApprovedTransfers{}

func NewMsgUpdateCollectionApprovedTransfers(creator string, collectionId sdkmath.Uint, approvedTransfersTimeline []*CollectionApprovedTransferTimeline) *MsgUpdateCollectionApprovedTransfers {
	return &MsgUpdateCollectionApprovedTransfers{
		Creator:           creator,
		CollectionId:      collectionId,
		ApprovedTransfersTimeline: approvedTransfersTimeline,
	}
}

func (msg *MsgUpdateCollectionApprovedTransfers) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCollectionApprovedTransfers) Type() string {
	return TypeMsgUpdateCollectionApprovedTransfers
}

func (msg *MsgUpdateCollectionApprovedTransfers) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCollectionApprovedTransfers) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCollectionApprovedTransfers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err := ValidateApprovedTransferTimeline(msg.ApprovedTransfersTimeline); err != nil {
		return err
	}

	if msg.CollectionId.IsNil() || msg.CollectionId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidCollectionID, "collection id cannot be 0")
	}

	for _, mapping := range msg.AddressMappings {
		if err := ValidateAddressMapping(mapping); err != nil {
			return err
		}
	}

	return nil
}
