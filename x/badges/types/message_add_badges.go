package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgMintAndDistributeBadges = "mint_and_distribute_badge"

var _ sdk.Msg = &MsgMintAndDistributeBadges{}

func NewMsgMintAndDistributeBadges(
	creator string,
	collectionId sdk.Uint,
	badgesToCreate []*Balance,
	transfers []*Transfer,
	collectionMetadataTimeline []*CollectionMetadataTimeline,
	badgeMetadataTimeline []*BadgeMetadataTimeline,
	offChainBalancesMetadataTimeline []*OffChainBalancesMetadataTimeline,
	approvedTransfersTimeline []*CollectionApprovedTransferTimeline,
	inheritedBalancesTimeline []*InheritedBalancesTimeline,
	addressMappings []*AddressMapping,
) *MsgMintAndDistributeBadges {
	// err := *new(error)
	// for _, transfer := range transfers {
	// 	transfer.Balances, err = HandleDuplicateBadgeIds(transfer.Balances)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// for _, timelineVal := range badgeMetadataTimeline {
	// 	for _, badgeMetadata := range timelineVal.BadgeMetadata {
	// 		badgeMetadata.BadgeIds = SortAndMergeOverlapping(badgeMetadata.BadgeIds)
	// 	}
	// }

	// for _, timelineVal := range approvedTransfersTimeline {
	// 	for _, approvedTransfer := range timelineVal.ApprovedTransfers {
	// 		approvedTransfer.BadgeIds = SortAndMergeOverlapping(approvedTransfer.BadgeIds)
	// 		approvedTransfer.TransferTimes = SortAndMergeOverlapping(approvedTransfer.TransferTimes)

	// 		for _, balance := range approvedTransfer.Claim.StartAmounts {
	// 			balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
	// 		}
	// 	}
	// }

	// badgesToCreate, err = HandleDuplicateBadgeIds(badgesToCreate)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, timelineVal := range inheritedBalancesTimeline {
	// 	for _, balance := range timelineVal.InheritedBalances {
	// 		balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
	// 		balance.ParentBadgeIds = SortAndMergeOverlapping(balance.ParentBadgeIds)
	// 	}
	// }

	return &MsgMintAndDistributeBadges{
		Creator:            creator,
		CollectionId:       collectionId,
		BadgesToCreate:     badgesToCreate,
		Transfers:          transfers,
		CollectionMetadataTimeline: collectionMetadataTimeline,
		BadgeMetadataTimeline:      badgeMetadataTimeline,
		OffChainBalancesMetadataTimeline:   offChainBalancesMetadataTimeline,
		ApprovedTransfersTimeline:  approvedTransfersTimeline,
		InheritedBalancesTimeline:  inheritedBalancesTimeline,
	}
}

func (msg *MsgMintAndDistributeBadges) Route() string {
	return RouterKey
}

func (msg *MsgMintAndDistributeBadges) Type() string {
	return TypeMsgMintAndDistributeBadges
}

func (msg *MsgMintAndDistributeBadges) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintAndDistributeBadges) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintAndDistributeBadges) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsZero() || msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection id")
	}

	if msg.BadgesToCreate == nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid badges to create")
	}

	if err := ValidateBalances(msg.BadgesToCreate); err != nil {
		return err
	}

	if msg.OffChainBalancesMetadataTimeline != nil {
		for _, timelineVal := range msg.OffChainBalancesMetadataTimeline {	
			err = ValidateURI(timelineVal.OffChainBalancesMetadata.Uri)
			if err != nil {
				return err
			}
		}
	}

	if msg.BadgeMetadataTimeline != nil && len(msg.BadgeMetadataTimeline) > 0 {
		for _, timelineVal := range msg.BadgeMetadataTimeline {
			err = ValidateBadgeMetadata(timelineVal.BadgeMetadata)
			if err != nil {
				return err
			}
		}
	}

	if msg.CollectionMetadataTimeline != nil {
		for _, timelineVal := range msg.CollectionMetadataTimeline {
			err = ValidateURI(timelineVal.CollectionMetadata.Uri)
			if err != nil {
				return err
			}
		}
	}


	for _, transfer := range msg.Transfers {
		err = ValidateTransfer(transfer)
		if err != nil {
			return err
		}
	}

	if msg.ApprovedTransfersTimeline != nil {
		for _, timelineVal := range msg.ApprovedTransfersTimeline {
			for _, approvedTransfer := range timelineVal.ApprovedTransfers {
				err = ValidateCollectionApprovedTransfer(*approvedTransfer)
				if err != nil {
					return err
				}
			}
		}
	}

	if msg.InheritedBalancesTimeline != nil {
		for _, timelineVal := range msg.InheritedBalancesTimeline {
			for _, balance := range timelineVal.InheritedBalances {
				err = ValidateRangesAreValid(balance.BadgeIds, true)
				if err != nil {
					return err
				}

				err = ValidateRangesAreValid(balance.ParentBadgeIds, true)
				if err != nil {
					return err
				}

				if balance.ParentCollectionId.IsZero() || balance.ParentCollectionId.IsNil() {
					return sdkerrors.Wrapf(ErrInvalidRequest, "invalid parent collection id")
				}
	
			}
		}
	}

	if len(msg.Transfers) > 0 || len(msg.ApprovedTransfersTimeline) > 0 {
		if msg.OffChainBalancesMetadataTimeline != nil && len(msg.OffChainBalancesMetadataTimeline) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "transfers and/or claims are set but collection has balances type = off-chain")
		}

		if msg.InheritedBalancesTimeline != nil && len(msg.InheritedBalancesTimeline) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "transfers and/or claims are set but collection has balances type = inherited")
		}
	}

	return nil
}
