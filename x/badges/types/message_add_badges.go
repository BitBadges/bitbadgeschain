package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgMintAndDistributeBadges = "mint_and_distribute_badge"

var _ sdk.Msg = &MsgMintAndDistributeBadges{}

func NewMsgMintAndDistributeBadges(
	creator string,
	collectionId sdkmath.Uint,
	badgesToCreate []*Balance,
	transfers []*Transfer,
	collectionMetadataTimeline []*CollectionMetadataTimeline,
	badgeMetadataTimeline []*BadgeMetadataTimeline,
	offChainBalancesMetadataTimeline []*OffChainBalancesMetadataTimeline,
	approvedTransfersTimeline []*CollectionApprovedTransferTimeline,
	inheritedBalancesTimeline []*InheritedBalancesTimeline,
	addressMappings []*AddressMapping,
) *MsgMintAndDistributeBadges {
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

	if msg.CollectionId.IsNil() || msg.CollectionId.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection id")
	}

	if msg.BadgesToCreate == nil {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid badges to create")
	}

	msg.BadgesToCreate, err = ValidateBalances(msg.BadgesToCreate)
	if err != nil {
		return err
	}

	if err := ValidateOffChainBalancesMetadataTimeline(msg.OffChainBalancesMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateBadgeMetadataTimeline(msg.BadgeMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateCollectionMetadataTimeline(msg.CollectionMetadataTimeline); err != nil {
		return err
	}


	for _, transfer := range msg.Transfers {
		err = ValidateTransfer(transfer)
		if err != nil {
			return err
		}
	}

	if err := ValidateApprovedTransferTimeline(msg.ApprovedTransfersTimeline); err != nil {
		return err
	}
	
	if err := ValidateInheritedBalancesTimeline(msg.InheritedBalancesTimeline); err != nil {
		return err
	}

	if len(msg.Transfers) > 0 || len(msg.ApprovedTransfersTimeline) > 0 {
		if msg.OffChainBalancesMetadataTimeline != nil && len(msg.OffChainBalancesMetadataTimeline) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "transfers and/or claims are set but collection has balances type = off-chain")
		}

		if msg.InheritedBalancesTimeline != nil && len(msg.InheritedBalancesTimeline) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "transfers and/or claims are set but collection has balances type = inherited")
		}
	}

	for _, mapping := range msg.AddressMappings {
		if err := ValidateAddressMapping(mapping); err != nil {
			return err
		}
	}

	return nil
}
