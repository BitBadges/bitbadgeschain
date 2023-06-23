package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgMintAndDistributeBadges = "mint_and_distribute_badge"

var _ sdk.Msg = &MsgMintAndDistributeBadges{}

func NewMsgMintAndDistributeBadges(
	creator string,
	collectionId sdk.Uint,
	badgesToCreate []*Balance,
	transfers []*Transfer,
	collectionMetadata *CollectionMetadata,
	badgeMetadata []*BadgeMetadata,
	offChainBalancesMetadata *OffChainBalancesMetadata,
	approvedTransfers []*CollectionApprovedTransfer,
	inheritedBalances []*InheritedBalance,
	addressMappings []*AddressMapping,
) *MsgMintAndDistributeBadges {
	for _, transfer := range transfers {
		for _, balance := range transfer.Balances {
			balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
		}
	}

	for _, badgeMetadata := range badgeMetadata {
		badgeMetadata.BadgeIds = SortAndMergeOverlapping(badgeMetadata.BadgeIds)
	}

	for _, approvedTransfer := range approvedTransfers {
		approvedTransfer.BadgeIds = SortAndMergeOverlapping(approvedTransfer.BadgeIds)
		approvedTransfer.TransferTimes = SortAndMergeOverlapping(approvedTransfer.TransferTimes)

		for _, balance := range approvedTransfer.Claim.StartAmounts {
			balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
		}
	}

	for _, badgeBalanceToCreate := range badgesToCreate {
		badgeBalanceToCreate.BadgeIds = SortAndMergeOverlapping(badgeBalanceToCreate.BadgeIds)
	}

	for _, balance := range inheritedBalances {
		balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
		balance.ParentBadgeIds = SortAndMergeOverlapping(balance.ParentBadgeIds)
	}

	return &MsgMintAndDistributeBadges{
		Creator:            creator,
		CollectionId:       collectionId,
		BadgesToCreate:     badgesToCreate,
		Transfers:          transfers,
		CollectionMetadata: collectionMetadata,
		BadgeMetadata:      badgeMetadata,
		OffChainBalancesMetadata:   offChainBalancesMetadata,
		ApprovedTransfers:  approvedTransfers,
		InheritedBalances:  inheritedBalances,
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsZero() || msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid collection id")
	}

	if msg.BadgesToCreate == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid badges to create")
	}

	if err := ValidateBalances(msg.BadgesToCreate); err != nil {
		return err
	}

	if msg.OffChainBalancesMetadata != nil {
		err = ValidateURI(msg.OffChainBalancesMetadata.Uri)
		if err != nil {
			return err
		}
	}

	if msg.BadgeMetadata != nil && len(msg.BadgeMetadata) > 0 {
		err = ValidateBadgeMetadata(msg.BadgeMetadata)
		if err != nil {
			return err
		}
	}

	if msg.CollectionMetadata != nil {
		err = ValidateURI(msg.CollectionMetadata.Uri)
		if err != nil {
			return err
		}
	}

	for _, transfer := range msg.Transfers {
		err = ValidateTransfer(transfer)
		if err != nil {
			return err
		}
	}

	for _, approvedTransfer := range msg.ApprovedTransfers {
		err = ValidateCollectionApprovedTransfer(*approvedTransfer)
		if err != nil {
			return err
		}
	}

	if msg.OffChainBalancesMetadata != nil {
		if msg.OffChainBalancesMetadata.Uri != "" || msg.OffChainBalancesMetadata.CustomData != "" {
			//We have off-chain balances

			if len(msg.Transfers) > 0 || len(msg.ApprovedTransfers) > 0 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "balances metadata denotes off-chain balances but claims and/or transfers are set")
			}
		}
	}

	if msg.InheritedBalances != nil {
		if len(msg.InheritedBalances) > 0 {
			if len(msg.Transfers) > 0 || len(msg.ApprovedTransfers) > 0 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "balances metadata denotes off-chain balances but claims and/or transfers are set")
			}
		}
	}

	return nil
}
