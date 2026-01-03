package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetIncomingApproval = "set_incoming_approval"

var _ sdk.Msg = &MsgSetIncomingApproval{}

func NewMsgSetIncomingApproval(creator string, collectionId Uint, approval *UserIncomingApproval) *MsgSetIncomingApproval {
	return &MsgSetIncomingApproval{
		Creator:      creator,
		CollectionId: collectionId,
		Approval:     approval,
	}
}

func (msg *MsgSetIncomingApproval) Route() string {
	return RouterKey
}

func (msg *MsgSetIncomingApproval) Type() string {
	return TypeMsgSetIncomingApproval
}

func (msg *MsgSetIncomingApproval) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetIncomingApproval) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetIncomingApproval) ValidateBasic() error {
	return msg.CheckAndCleanMsg(sdk.Context{}, false)
}

func (msg *MsgSetIncomingApproval) CheckAndCleanMsg(ctx sdk.Context, canChangeValues bool) error {
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

	// Validate approval
	if msg.Approval == nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "approval cannot be nil")
	}

	if err := ValidateUserIncomingApprovals(ctx, []*UserIncomingApproval{msg.Approval}, msg.Creator, canChangeValues); err != nil {
		return err
	}

	return nil
}
