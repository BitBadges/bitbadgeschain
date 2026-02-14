package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgDeleteIncomingApproval = "delete_incoming_approval"

var _ sdk.Msg = &MsgDeleteIncomingApproval{}

func NewMsgDeleteIncomingApproval(creator string, collectionId Uint, approvalId string) *MsgDeleteIncomingApproval {
	return &MsgDeleteIncomingApproval{
		Creator:      creator,
		CollectionId: collectionId,
		ApprovalId:   approvalId,
	}
}

func (msg *MsgDeleteIncomingApproval) Route() string {
	return RouterKey
}

func (msg *MsgDeleteIncomingApproval) Type() string {
	return TypeMsgDeleteIncomingApproval
}

func (msg *MsgDeleteIncomingApproval) GetSigners() []sdk.AccAddress {
	// MustAccAddressFromBech32 panics if address is invalid, which is expected
	// since ValidateBasic() should have already validated the address
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteIncomingApproval) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteIncomingApproval) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate collection ID
	// Allow collectionId = 0 for auto-prev resolution (used in multi-msg transactions)
	// The actual validation and resolution happens in resolveCollectionIdWithAutoPrev
	if msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidCollectionID, "collection ID cannot be nil")
	}

	// Validate approval ID
	if msg.ApprovalId == "" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "approval ID cannot be empty")
	}

	return nil
}
