package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgNewCollection = "new_collection"

var _ sdk.Msg = &MsgNewCollection{}

func NewMsgNewCollection(
	creator string,
	standardsTimeline []*StandardsTimeline,
	badgesToCreate []*Balance,
	collectionMetadataTimeline []*CollectionMetadataTimeline,
	badgeMetadataTimeline []*BadgeMetadataTimeline,
	permissions *CollectionPermissions,
	approvedTransfersTimeline []*CollectionApprovedTransferTimeline,
	customDataTimeline []*CustomDataTimeline,
	transfers []*Transfer,
	offChainBalancesMetadataTimeline []*OffChainBalancesMetadataTimeline,
	contractAddressTimeline []*ContractAddressTimeline,
	balancesType sdkmath.Uint,
	inheritedBalancesTimeline []*InheritedBalancesTimeline,
	defaultApprovedOutgoingTransfersTimeline []*UserApprovedOutgoingTransferTimeline,
	defaultApprovedIncomingTransfersTimeline []*UserApprovedIncomingTransferTimeline,
) *MsgNewCollection {
	return &MsgNewCollection{
		Creator:            creator,
		StandardsTimeline:   standardsTimeline,
		BadgesToCreate:     badgesToCreate,
		CollectionMetadataTimeline: collectionMetadataTimeline,
		BadgeMetadataTimeline: badgeMetadataTimeline,
		Permissions:        permissions,

		CollectionApprovedTransfersTimeline: approvedTransfersTimeline,
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

	if err := ValidateOffChainBalancesMetadataTimeline(msg.OffChainBalancesMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateBadgeMetadataTimeline(msg.BadgeMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateCollectionMetadataTimeline(msg.CollectionMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateContractAddressTimeline(msg.ContractAddressTimeline); err != nil {
		return err
	}

	if err := ValidateCustomDataTimeline(msg.CustomDataTimeline); err != nil {
		return err
	}

	if err := ValidateStandardsTimeline(msg.StandardsTimeline); err != nil {
		return err
	}

	if err := ValidatePermissions(msg.Permissions); err != nil {
		return err
	}

	if err := ValidateUserApprovedIncomingTransferTimeline(msg.DefaultApprovedIncomingTransfersTimeline, msg.Creator); err != nil {
		return err
	}

	if err := ValidateUserApprovedOutgoingTransferTimeline(msg.DefaultApprovedOutgoingTransfersTimeline, msg.Creator); err != nil {
		return err
	} 

	msg.BadgesToCreate, err = ValidateBalances(msg.BadgesToCreate)
	if err != nil {
		return err
	}

	if msg.BalancesType.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "balances type cannot be nil")
	}

	if !msg.BalancesType.IsZero() {
		//We have off-chain or inherited balances

		if len(msg.Transfers) > 0 || len(msg.CollectionApprovedTransfersTimeline) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but claims and/or transfers are set")
		}
	}

	if err := ValidateInheritedBalancesTimeline(msg.InheritedBalancesTimeline); err != nil {
		return err
	}

	for _, transfer := range msg.Transfers {
		err = ValidateTransfer(transfer)
		if err != nil {
			return err
		}
	}

	if err := ValidateApprovedTransferTimeline(msg.CollectionApprovedTransfersTimeline); err != nil {
		return err
	}

	for _, mapping := range msg.AddressMappings {
		if err := ValidateAddressMapping(mapping); err != nil {
			return err
		}
	}

	return nil
}
