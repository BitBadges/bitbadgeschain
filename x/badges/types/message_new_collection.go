package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgNewCollection = "new_collection"

var _ sdk.Msg = &MsgNewCollection{}

func NewMsgNewCollection(
	creator string,
	standardTimeline []*StandardTimeline,
	badgesToCreate []*Balance,
	collectionMetadataTimeline []*CollectionMetadataTimeline,
	badgeMetadataTimeline []*BadgeMetadataTimeline,
	permissions *CollectionPermissions,
	approvedTransfersTimeline []*CollectionApprovedTransferTimeline,
	customDataTimeline []*CustomDataTimeline,
	transfers []*Transfer,
	offChainBalancesMetadataTimeline []*OffChainBalancesMetadataTimeline,
	contractAddressTimeline []*ContractAddressTimeline,
	balancesType sdk.Uint,
	inheritedBalancesTimeline []*InheritedBalancesTimeline,
	defaultApprovedOutgoingTransfersTimeline []*UserApprovedOutgoingTransferTimeline,
	defaultApprovedIncomingTransfersTimeline []*UserApprovedIncomingTransferTimeline,
) *MsgNewCollection {
	// for _, transfer := range transfers {
	// 	for _, balance := range transfer.Balances {
	// 		balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
	// 	}
	// }

	// for _, badgeBalanceToCreate := range badgesToCreate {
	// 	badgeBalanceToCreate.BadgeIds = SortAndMergeOverlapping(badgeBalanceToCreate.BadgeIds)
	// }

	// for _, badgeMetadata := range badgeMetadata {
	// 	badgeMetadata.BadgeIds = SortAndMergeOverlapping(badgeMetadata.BadgeIds)
	// }

	// for _, approvedTransfer := range approvedTransfers {
	// 	approvedTransfer.BadgeIds = SortAndMergeOverlapping(approvedTransfer.BadgeIds)
	// 	approvedTransfer.TransferTimes = SortAndMergeOverlapping(approvedTransfer.TransferTimes)

	// 	for _, balance := range approvedTransfer.Claim.StartAmounts {
	// 		balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
	// 	}
	// }

	// for _, balance := range inheritedBalances {
	// 	balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
	// 	balance.ParentBadgeIds = SortAndMergeOverlapping(balance.ParentBadgeIds)
	// }

	//TODO: permissions sort and merge overlapping



	return &MsgNewCollection{
		Creator:            creator,
		StandardsTimeline:   standardTimeline,
		BadgesToCreate:     badgesToCreate,
		CollectionMetadataTimeline: collectionMetadataTimeline,
		BadgeMetadataTimeline: badgeMetadataTimeline,
		Permissions:        permissions,

		ApprovedTransfersTimeline: approvedTransfersTimeline,
		CustomDataTimeline: customDataTimeline,
		Transfers:          transfers,
		OffChainBalancesMetadataTimeline: offChainBalancesMetadataTimeline,
		ContractAddressTimeline: contractAddressTimeline,
		BalancesType: balancesType,
		InheritedBalancesTimeline: inheritedBalancesTimeline,

		DefaultApprovedOutgoingTransfersTimeline: defaultApprovedOutgoingTransfersTimeline,
		DefaultApprovedIncomingTransfersTimeline: defaultApprovedIncomingTransfersTimeline,
	}
}

func (msg *MsgNewCollection) Route() string {
	return RouterKey
}

func (msg *MsgNewCollection) Type() string {
	return TypeMsgNewCollection
}

func (msg *MsgNewCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgNewCollection) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgNewCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
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

	if err := ValidatePermissions(msg.Permissions, false); err != nil {
		return err
	}

	if err := ValidateBalances(msg.BadgesToCreate); err != nil {
		return err
	}


	if !msg.BalancesType.IsZero() {
		//We have off-chain or inherited balances

		if len(msg.Transfers) > 0 || len(msg.ApprovedTransfersTimeline) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but claims and/or transfers are set")
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

	//TODO: Enforce irrelevant permissions to be permanently disallowed according to balances type
	//Can't update approved transfers if not on-chain
	//Can't update off-chain balances if not off-chain
	//Can't update inherited balances if not inherited

	for _, transfer := range msg.Transfers {
		err = ValidateTransfer(transfer)
		if err != nil {
			return err
		}
	}

	for _, timelineVal := range msg.ApprovedTransfersTimeline {
		for _, approvedTransfer := range timelineVal.ApprovedTransfers {
			err = ValidateCollectionApprovedTransfer(*approvedTransfer)
			if err != nil {
				return err
			}
		}
	}


	return nil
}
