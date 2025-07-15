package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSetOutgoingApproval = "set_outgoing_approval"

var _ sdk.Msg = &MsgSetOutgoingApproval{}

func NewMsgSetOutgoingApproval(creator string, collectionId Uint, approval *UserOutgoingApproval) *MsgSetOutgoingApproval {
	return &MsgSetOutgoingApproval{
		Creator:      creator,
		CollectionId: collectionId,
		Approval:     approval,
	}
}

func (msg *MsgSetOutgoingApproval) Route() string {
	return RouterKey
}

func (msg *MsgSetOutgoingApproval) Type() string {
	return TypeMsgSetOutgoingApproval
}

func (msg *MsgSetOutgoingApproval) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetOutgoingApproval) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetOutgoingApproval) ValidateBasic() error {
	return msg.CheckAndCleanMsg(sdk.Context{}, false)
}

func (msg *MsgSetOutgoingApproval) CheckAndCleanMsg(ctx sdk.Context, canChangeValues bool) error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate collection ID
	if msg.CollectionId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidCollectionID, "collection ID cannot be zero")
	}

	// Validate approval
	if msg.Approval == nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "approval cannot be nil")
	}

	if err := ValidateUserOutgoingApprovals(ctx, []*UserOutgoingApproval{msg.Approval}, msg.Creator, canChangeValues); err != nil {
		return err
	}

	return nil
}
