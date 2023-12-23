package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUniversalUpdateCollection = "update_collection"

var _ sdk.Msg = &MsgUniversalUpdateCollection{}

func NewMsgUniversalUpdateCollection(creator string) *MsgUniversalUpdateCollection {
	return &MsgUniversalUpdateCollection{
		Creator: creator,
	}
}

func (msg *MsgUniversalUpdateCollection) Route() string {
	return RouterKey
}

func (msg *MsgUniversalUpdateCollection) Type() string {
	return TypeMsgUniversalUpdateCollection
}

func (msg *MsgUniversalUpdateCollection) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUniversalUpdateCollection) GetSignBytes() []byte {
	bz := AminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUniversalUpdateCollection) ValidateBasic() error {
	return msg.CheckAndCleanMsg(false)
}

func (msg *MsgUniversalUpdateCollection) CheckAndCleanMsg(canChangeValues bool) error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection id")
	}

	if msg.BadgesToCreate != nil {
		msg.BadgesToCreate, err = ValidateBalances(msg.BadgesToCreate, canChangeValues)
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

	if err := ValidateBadgeMetadataTimeline(msg.BadgeMetadataTimeline, canChangeValues); err != nil {
		return err
	}

	if err := ValidateCollectionMetadataTimeline(msg.CollectionMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateCollectionApprovals(msg.CollectionApprovals, canChangeValues); err != nil {
		return err
	}

	if err := ValidateOffChainBalancesMetadataTimeline(msg.OffChainBalancesMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateBadgeMetadataTimeline(msg.BadgeMetadataTimeline, canChangeValues); err != nil {
		return err
	}

	if err := ValidateCollectionMetadataTimeline(msg.CollectionMetadataTimeline); err != nil {
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

	if err := ValidatePermissions(msg.CollectionPermissions, canChangeValues); err != nil {
		return err
	}

	if msg.DefaultBalances == nil {
		msg.DefaultBalances = &UserBalanceStore{}
	}


	if _, err := ValidateBalances(msg.DefaultBalances.Balances, canChangeValues); err != nil {
		return err
	}

	if err := ValidateUserIncomingApprovals(msg.DefaultBalances.IncomingApprovals, msg.Creator, canChangeValues); err != nil {
		return err
	}

	if err := ValidateUserOutgoingApprovals(msg.DefaultBalances.OutgoingApprovals, msg.Creator, canChangeValues); err != nil {
		return err
	}

	if msg.DefaultBalances.UserPermissions == nil {
		msg.DefaultBalances.UserPermissions = &UserPermissions{
			CanUpdateIncomingApprovals:                         []*UserIncomingApprovalPermission{},
			CanUpdateOutgoingApprovals:                         []*UserOutgoingApprovalPermission{},
			CanUpdateAutoApproveSelfInitiatedOutgoingTransfers: []*ActionPermission{},
			CanUpdateAutoApproveSelfInitiatedIncomingTransfers: []*ActionPermission{},
		}
	}

	if err := ValidateUserPermissions(msg.DefaultBalances.UserPermissions, canChangeValues); err != nil {
		return err
	}

	if err := ValidateOffChainBalancesMetadataTimeline(msg.OffChainBalancesMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateBadgeMetadataTimeline(msg.BadgeMetadataTimeline, canChangeValues); err != nil {
		return err
	}

	if err := ValidateCollectionMetadataTimeline(msg.CollectionMetadataTimeline); err != nil {
		return err
	}

	if err := ValidateCustomDataTimeline(msg.CustomDataTimeline); err != nil {
		return err
	}

	if err := ValidateStandardsTimeline(msg.StandardsTimeline); err != nil {
		return err
	}

	if msg.CollectionId.IsZero() {
		if msg.BalancesType != "Standard" && msg.BalancesType != "Inherited" && msg.BalancesType != "Off-Chain - Indexed" && msg.BalancesType != "Off-Chain - Non-Indexed" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "balances type must be Standard, Inherited, Off-Chain - Indexed, or Off-Chain - Non-Indexed")
		}

		if msg.BalancesType != "Standard" {
			if len(msg.CollectionApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but claims and/or transfers are set")
			}

			if len(msg.DefaultBalances.IncomingApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but default approvals are set")
			}

			if len(msg.DefaultBalances.OutgoingApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but default approvals are set")
			}

			if (len(msg.DefaultBalances.Balances) > 0) {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but default balances are set")
			}

			if len(msg.DefaultBalances.UserPermissions.CanUpdateIncomingApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but default user permissions are being set")
			}

			if len(msg.DefaultBalances.UserPermissions.CanUpdateOutgoingApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but default user permissions are being set")
			}

			if len(msg.DefaultBalances.UserPermissions.CanUpdateAutoApproveSelfInitiatedIncomingTransfers) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but default user permissions are being set")
			}

			if len(msg.DefaultBalances.UserPermissions.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes off-chain balances but default user permissions are being set")
			}
		}

		if msg.BalancesType != "Off-Chain - Indexed" && msg.BalancesType != "Off-Chain - Non-Indexed" {
			if len(msg.OffChainBalancesMetadataTimeline) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes on-chain balances but off-chain balances are set")
			}
		}

		if msg.BalancesType == "Inherited" {
			// if msg.InheritedCollectionId.IsNil() || msg.InheritedCollectionId.IsZero() {
			// 	return sdkerrors.Wrapf(ErrInvalidRequest, "inherited collection id must be set for inherited balances")
			// }

			if msg.BadgesToCreate != nil && len(msg.BadgesToCreate) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "badges are inherited from parent so you should not specify to create any badges")
			}
		}

		if msg.BalancesType == "Off-Chain - Non-Indexed" {
			for _, offChainUriObj := range msg.OffChainBalancesMetadataTimeline {
				//must contain "{address}" in uri
				if strings.Contains(offChainUriObj.OffChainBalancesMetadata.Uri, "{address}") == false {
					return sdkerrors.Wrapf(ErrInvalidRequest, "balances type is non-indexed but uri does not contain {address} in uri")
				}
			}
		}
	}

	if len(msg.CollectionApprovals) > 0 {
		if msg.OffChainBalancesMetadataTimeline != nil && len(msg.OffChainBalancesMetadataTimeline) > 0 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "transfers and/or claims are set but collection has balances type = off-chain")
		}

		if msg.BalancesType == "Inherited" {
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
