package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgDeleteOutgoingApproval = "delete_outgoing_approval"

var _ sdk.Msg = &MsgDeleteOutgoingApproval{}

func NewMsgDeleteOutgoingApproval(creator string, collectionId Uint, approvalId string) *MsgDeleteOutgoingApproval {
	return &MsgDeleteOutgoingApproval{
		Creator:      creator,
		CollectionId: collectionId,
		ApprovalId:   approvalId,
	}
}

func (msg *MsgDeleteOutgoingApproval) Route() string {
	return RouterKey
}

func (msg *MsgDeleteOutgoingApproval) Type() string {
	return TypeMsgDeleteOutgoingApproval
}

func (msg *MsgDeleteOutgoingApproval) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteOutgoingApproval) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteOutgoingApproval) ValidateBasic() error {
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
