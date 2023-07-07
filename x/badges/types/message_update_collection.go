package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateCollection = "update_collection"

var _ sdk.Msg = &MsgUpdateCollection{}

func NewMsgUpdateCollection(creator string) *MsgUpdateCollection {
	return &MsgUpdateCollection{
		Creator: creator,
	}
}

func (msg *MsgUpdateCollection) Route() string {
	return RouterKey
}

func (msg *MsgUpdateCollection) Type() string {
	return TypeMsgUpdateCollection
}

func (msg *MsgUpdateCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateCollection) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateCollection) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection id")
	}

	if msg.BadgesToCreate != nil {
		msg.BadgesToCreate, err = ValidateBalances(msg.BadgesToCreate)
		if err != nil {
			return err
		}
	}	

	if err := ValidateIsArchivedTimeline(msg.IsArchivedTimeline); err != nil {
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

	if err := ValidateApprovedTransferTimeline(msg.CollectionApprovedTransfersTimeline); err != nil {
		return err
	}

	if err := ValidateInheritedBalancesTimeline(msg.InheritedBalancesTimeline); err != nil {
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

	if err := ValidateContractAddressTimeline(msg.ContractAddressTimeline); err != nil {
		return err
	}

	if err := ValidateCustomDataTimeline(msg.CustomDataTimeline); err != nil {
		return err
	}

	if err := ValidateStandardsTimeline(msg.StandardsTimeline); err != nil {
		return err
	}

	if msg.CollectionPermissions == nil {
		msg.CollectionPermissions = &CollectionPermissions{}
	}

	if err := ValidatePermissions(msg.CollectionPermissions); err != nil {
		return err
	}

	if err := ValidateUserApprovedIncomingTransferTimeline(msg.DefaultApprovedIncomingTransfersTimeline, msg.Creator); err != nil {
		return err
	}

	if err := ValidateUserApprovedOutgoingTransferTimeline(msg.DefaultApprovedOutgoingTransfersTimeline, msg.Creator); err != nil {
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

	if err := ValidateContractAddressTimeline(msg.ContractAddressTimeline); err != nil {
		return err
	}

	if err := ValidateCustomDataTimeline(msg.CustomDataTimeline); err != nil {
		return err
	}

	if err := ValidateStandardsTimeline(msg.StandardsTimeline); err != nil {
		return err
	}

	if msg.CollectionId.IsZero() {
		if msg.BalancesType != "Standard" && msg.BalancesType != "Inherited" && msg.BalancesType != "Off-Chain" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "balances type must be Standard, Inherited, or Off-Chain")
		}

		if msg.BalancesType != "Standard" {
			if len(msg.CollectionApprovedTransfersTimeline) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but claims and/or transfers are set")
			}
		}
	}

	if err := ValidateInheritedBalancesTimeline(msg.InheritedBalancesTimeline); err != nil {
		return err
	}

	if len(msg.CollectionApprovedTransfersTimeline) > 0 {
		if msg.OffChainBalancesMetadataTimeline != nil && len(msg.OffChainBalancesMetadataTimeline) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "transfers and/or claims are set but collection has balances type = off-chain")
		}

		if msg.InheritedBalancesTimeline != nil && len(msg.InheritedBalancesTimeline) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "transfers and/or claims are set but collection has balances type = inherited")
		}
	}

	for _, timelineVal := range msg.ManagerTimeline {
		_, err = sdk.AccAddressFromBech32(timelineVal.Manager)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidAddress, "invalid provided address (%s)", err)
		}
	}

	if err := ValidateManagerTimeline(msg.ManagerTimeline); err != nil {
		return err
	}

	

	return nil
}
