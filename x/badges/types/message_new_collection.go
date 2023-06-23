package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgNewCollection = "new_collection"

var _ sdk.Msg = &MsgNewCollection{}

func NewMsgNewCollection(
	creator string,
	standard sdk.Uint,
	badgesToCreate []*Balance,
	collectionMetadata *CollectionMetadata,
	badgeMetadata []*BadgeMetadata,
	permissions *CollectionPermissions,
	approvedTransfers []*CollectionApprovedTransfer,
	customData string,
	transfers []*Transfer,
	offChainBalancesMetadata *OffChainBalancesMetadata,
	contractAddress string,
	balancesType sdk.Uint,
	inheritedBalances []*InheritedBalance,
) *MsgNewCollection {
	for _, transfer := range transfers {
		for _, balance := range transfer.Balances {
			balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
		}
	}

	for _, badgeBalanceToCreate := range badgesToCreate {
		badgeBalanceToCreate.BadgeIds = SortAndMergeOverlapping(badgeBalanceToCreate.BadgeIds)
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

	for _, balance := range inheritedBalances {
		balance.BadgeIds = SortAndMergeOverlapping(balance.BadgeIds)
		balance.ParentBadgeIds = SortAndMergeOverlapping(balance.ParentBadgeIds)
	}

	//TODO: permissions sort and merge overlapping

	return &MsgNewCollection{
		Creator:            creator,
		CollectionMetadata: collectionMetadata,
		BadgeMetadata:      badgeMetadata,
		BadgesToCreate:     badgesToCreate,
		ApprovedTransfers:  approvedTransfers,
		CustomData:         customData,
		Permissions:        permissions,
		Standard:           standard,
		Transfers:          transfers,
		OffChainBalancesMetadata:   offChainBalancesMetadata,
		ContractAddress:    contractAddress,
		BalancesType:       balancesType,
		InheritedBalances:  inheritedBalances,
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
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

	if err := ValidatePermissions(msg.Permissions, false); err != nil {
		return err
	}

	if msg.BadgeMetadata != nil {
		if err := ValidateBadgeMetadata(msg.BadgeMetadata); err != nil {
			return err
		}
	}

	if err := ValidateBalances(msg.BadgesToCreate); err != nil {
		return err
	}

	if msg.InheritedBalances != nil {
		if len(msg.InheritedBalances) > 0 {
			if len(msg.Transfers) > 0 || len(msg.ApprovedTransfers) > 0 {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "balances metadata denotes off-chain balances but claims and/or transfers are set")
			}
		}
	}

	if !msg.BalancesType.IsZero() {
		//We have off-chain or inherited balances

		if len(msg.Transfers) > 0 || len(msg.ApprovedTransfers) > 0 {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "balances metadata denotes off-chain balances but claims and/or transfers are set")
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

	for _, approvedTransfer := range msg.ApprovedTransfers {
		err = ValidateCollectionApprovedTransfer(*approvedTransfer)
		if err != nil {
			return err
		}
	}

	if msg.Standard.IsNil() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "invalid standard (%s)", msg.Standard)
	}

	return nil
}
