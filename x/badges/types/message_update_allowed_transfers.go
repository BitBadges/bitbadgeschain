package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateCollectionApprovedTransfers = "update_allowed_transfers"

var _ sdk.Msg = &MsgUpdateCollectionApprovedTransfers{}

func NewMsgUpdateCollectionApprovedTransfers(creator string, collectionId sdk.Uint, approvedTransfersTimeline []*CollectionApprovedTransferTimeline) *MsgUpdateCollectionApprovedTransfers {
	// for _, approvedTransfer := range approvedTransfers {
	// 	approvedTransfer.BadgeIds = SortAndMergeOverlapping(approvedTransfer.BadgeIds)
	// 	approvedTransfer.TransferTimes = SortAndMergeOverlapping(approvedTransfer.TransferTimes)

	// 	for _, balance := range approvedTransfer.Claim.StartAmounts {
	// 		balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
	// 	}
	// }
	
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

	for _, timelineVal := range msg.ApprovedTransfersTimeline {
		for _, approvedTransfer := range timelineVal.ApprovedTransfers {
			err = ValidateCollectionApprovedTransfer(*approvedTransfer)
			if err != nil {
				return err
			}
		}
	}

	if msg.CollectionId.IsZero() || msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidBadgeID, "collection id cannot be 0")
	}

	return nil
}
