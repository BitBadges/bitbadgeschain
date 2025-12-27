package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgPurgeApprovals = "purge_approvals"

var _ sdk.Msg = &MsgPurgeApprovals{}

func NewMsgPurgeApprovals(creator string, collectionId Uint, purgeExpired bool, approverAddress string, purgeCounterpartyApprovals bool, approvalsToPurge []*ApprovalIdentifierDetails) *MsgPurgeApprovals {
	return &MsgPurgeApprovals{
		Creator:                    creator,
		CollectionId:               collectionId,
		PurgeExpired:               purgeExpired,
		ApproverAddress:            approverAddress,
		PurgeCounterpartyApprovals: purgeCounterpartyApprovals,
		ApprovalsToPurge:           approvalsToPurge,
	}
}

func (msg *MsgPurgeApprovals) Route() string {
	return RouterKey
}

func (msg *MsgPurgeApprovals) Type() string {
	return TypeMsgPurgeApprovals
}

func (msg *MsgPurgeApprovals) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgPurgeApprovals) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPurgeApprovals) ValidateBasic() error {
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

	// Determine target address (who we're purging approvals for)
	targetAddress := msg.ApproverAddress
	if targetAddress == "" {
		targetAddress = msg.Creator
	}

	// Validate approver address if provided
	if msg.ApproverAddress != "" {
		_, err := sdk.AccAddressFromBech32(msg.ApproverAddress)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddress, "invalid approver address (%s)", err)
		}
	}

	// Validate approvalsToPurge - must not be empty
	if len(msg.ApprovalsToPurge) == 0 {
		return sdkerrors.Wrapf(ErrInvalidRequest, "approvalsToPurge cannot be empty")
	}

	// Validate each approval in approvalsToPurge
	for i, approval := range msg.ApprovalsToPurge {
		if approval == nil {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval at index %d cannot be nil", i)
		}

		if approval.ApprovalId == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval ID at index %d cannot be empty", i)
		}

		if approval.ApprovalLevel == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approval level at index %d cannot be empty", i)
		}

		// Validate approval level
		switch approval.ApprovalLevel {
		case "collection", "incoming", "outgoing":
			// Valid levels
		default:
			return sdkerrors.Wrapf(ErrInvalidRequest, "invalid approval level at index %d: %s (must be 'collection', 'incoming', or 'outgoing')", i, approval.ApprovalLevel)
		}

		// Validate approver address in approval if provided
		if approval.ApproverAddress != "" {
			_, err := sdk.AccAddressFromBech32(approval.ApproverAddress)
			if err != nil {
				return sdkerrors.Wrapf(ErrInvalidAddress, "invalid approver address in approval at index %d (%s)", i, err)
			}
		}

		// For collection-level approvals, approverAddress should be empty
		if approval.ApprovalLevel == "collection" && approval.ApproverAddress != "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approver address must be empty for collection-level approvals at index %d", i)
		}

		// For user-level approvals, approverAddress should be provided
		if (approval.ApprovalLevel == "incoming" || approval.ApprovalLevel == "outgoing") && approval.ApproverAddress == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "approver address must be provided for %s-level approvals at index %d", approval.ApprovalLevel, i)
		}
	}

	// Rule 1: If creator is updating their own approvals, purgeExpired must be true and purgeCounterpartyApprovals must be false
	if targetAddress == msg.Creator {
		if !msg.PurgeExpired {
			return sdkerrors.Wrapf(ErrInvalidRequest, "when purging own approvals, purgeExpired must be true")
		}
		if msg.PurgeCounterpartyApprovals {
			return sdkerrors.Wrapf(ErrInvalidRequest, "when purging own approvals, purgeCounterpartyApprovals must be false")
		}
	}

	// Rule 2: If updating someone else's approvals, we can set either purgeExpired or purgeCounterpartyApprovals
	// (This is handled in the keeper implementation based on auto-deletion options)

	return nil
}
