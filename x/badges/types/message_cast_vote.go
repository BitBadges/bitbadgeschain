package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCastVote = "msg_cast_vote"

var _ sdk.Msg = &MsgCastVote{}

func NewMsgCastVote(creator string, collectionId sdkmath.Uint, approvalLevel string, approverAddress string, approvalId string, proposalId string, yesWeight sdkmath.Uint) *MsgCastVote {
	return &MsgCastVote{
		Creator:        creator,
		CollectionId:   collectionId,
		ApprovalLevel:  approvalLevel,
		ApproverAddress: approverAddress,
		ApprovalId:     approvalId,
		ProposalId:     proposalId,
		YesWeight:      yesWeight,
	}
}

func (msg *MsgCastVote) Route() string {
	return RouterKey
}

func (msg *MsgCastVote) Type() string {
	return TypeMsgCastVote
}

func (msg *MsgCastVote) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCastVote) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCastVote) ValidateBasic() error {
	if len(msg.Creator) == 0 {
		return sdkerrors.Wrapf(ErrInvalidAddress, "creator address cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "collectionId cannot be zero")
	}

	if msg.ApprovalLevel != "collection" && msg.ApprovalLevel != "incoming" && msg.ApprovalLevel != "outgoing" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "approvalLevel must be 'collection', 'incoming', or 'outgoing', got %s", msg.ApprovalLevel)
	}

	// For collection-level approvals, approverAddress should be empty
	if msg.ApprovalLevel == "collection" {
		if msg.ApproverAddress != "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approverAddress must be empty for collection-level approvals")
		}
	} else {
		// For incoming/outgoing approvals, approverAddress should be valid
		if msg.ApproverAddress == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approverAddress cannot be empty for %s-level approvals", msg.ApprovalLevel)
		}
		if err := ValidateAddress(msg.ApproverAddress, false); err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddress, "invalid approverAddress (%s)", err)
		}
	}

	if msg.ApprovalId == "" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "approvalId cannot be empty")
	}

	if msg.ProposalId == "" {
		return sdkerrors.Wrapf(ErrInvalidRequest, "proposalId cannot be empty")
	}

	if msg.YesWeight.GT(sdkmath.NewUint(100)) {
		return sdkerrors.Wrapf(ErrInvalidRequest, "yesWeight must be between 0 and 100, got %s", msg.YesWeight.String())
	}

	return nil
}

