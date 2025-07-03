package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUniversalUpdateCollection = "update_collection"

var _ sdk.Msg = &MsgUniversalUpdateCollection{}

func NewMsgUniversalUpdateCollection(creator string, creatorOverride string) *MsgUniversalUpdateCollection {
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
	sorted := sdk.MustSortJSON(bz)
	return sorted
}

func (msg *MsgUniversalUpdateCollection) ValidateBasic() error {
	return msg.CheckAndCleanMsg(sdk.Context{}, false)
}

func (msg *MsgUniversalUpdateCollection) CheckAndCleanMsg(ctx sdk.Context, canChangeValues bool) error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.CollectionId.IsNil() {
		return sdkerrors.Wrapf(ErrInvalidRequest, "invalid collection id")
	}

	if msg.ValidBadgeIds != nil {
		err = ValidateRangesAreValid(msg.ValidBadgeIds, false, false)
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

	if err := ValidateCollectionApprovals(ctx, msg.CollectionApprovals, canChangeValues); err != nil {
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

	if _, err := ValidateBalances(ctx, msg.DefaultBalances.Balances, canChangeValues); err != nil {
		return err
	}

	if err := ValidateUserIncomingApprovals(ctx, msg.DefaultBalances.IncomingApprovals, msg.Creator, canChangeValues); err != nil {
		return err
	}

	if err := ValidateUserOutgoingApprovals(ctx, msg.DefaultBalances.OutgoingApprovals, msg.Creator, canChangeValues); err != nil {
		return err
	}

	if msg.DefaultBalances.UserPermissions == nil {
		msg.DefaultBalances.UserPermissions = &UserPermissions{
			CanUpdateIncomingApprovals:                         []*UserIncomingApprovalPermission{},
			CanUpdateOutgoingApprovals:                         []*UserOutgoingApprovalPermission{},
			CanUpdateAutoApproveSelfInitiatedOutgoingTransfers: []*ActionPermission{},
			CanUpdateAutoApproveSelfInitiatedIncomingTransfers: []*ActionPermission{},
			CanUpdateAutoApproveAllIncomingTransfers:           []*ActionPermission{},
		}
	}

	//IMPORTANT: Default balances should only be able to specify general rules for incoming and outgoing approvals
	//           They should not be able to specify approval criteria like coin transfers, approval amounts, etc.
	//           This prevents default approval attacks where a creator can set an approval that approves a coin transfer
	//           and then use that approval to transfer coins without their permissions
	//           We can allow stuff on a more fine-grained level in the future but for now, we just disallow this
	for _, incomingApproval := range msg.DefaultBalances.IncomingApprovals {
		approvalCriteria := CastIncomingApprovalCriteriaToCollectionApprovalCriteria(incomingApproval.ApprovalCriteria)
		if approvalCriteria != nil && !CollectionApprovalHasNoSideEffects(approvalCriteria) {
			return sdkerrors.Wrapf(ErrInvalidRequest, "incoming approval criteria must be nil for default balances")
		}
	}

	for _, outgoingApproval := range msg.DefaultBalances.OutgoingApprovals {
		approvalCriteria := CastOutgoingApprovalCriteriaToCollectionApprovalCriteria(outgoingApproval.ApprovalCriteria)
		if approvalCriteria != nil && !CollectionApprovalHasNoSideEffects(approvalCriteria) {
			return sdkerrors.Wrapf(ErrInvalidRequest, "outgoing approval criteria must be nil for default balances")
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
		if msg.BalancesType != "Standard" && msg.BalancesType != "Off-Chain - Indexed" && msg.BalancesType != "Off-Chain - Non-Indexed" && msg.BalancesType != "Non-Public" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "balances type must be Standard, Inherited, Off-Chain - Indexed, or Off-Chain - Non-Indexed")
		}

		if msg.BalancesType != "Standard" {
			if len(msg.CollectionApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balance type is off-chain or non-public but claims and/or transfers are set")
			}

			if len(msg.DefaultBalances.IncomingApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balance type is off-chain or non-public but default approvals are set")
			}

			if len(msg.DefaultBalances.OutgoingApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balance type is off-chain or non-public but default approvals are set")
			}

			if len(msg.DefaultBalances.Balances) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balance type is off-chain or non-public but default balances are set")
			}

			if len(msg.DefaultBalances.UserPermissions.CanUpdateIncomingApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balance type is off-chain or non-public but default user permissions are being set")
			}

			if len(msg.DefaultBalances.UserPermissions.CanUpdateOutgoingApprovals) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balance type is off-chain or non-public but default user permissions are being set")
			}

			if len(msg.DefaultBalances.UserPermissions.CanUpdateAutoApproveSelfInitiatedIncomingTransfers) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balance type is off-chain or non-public but default user permissions are being set")
			}

			if len(msg.DefaultBalances.UserPermissions.CanUpdateAutoApproveSelfInitiatedOutgoingTransfers) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balance type is off-chain or non-public but default user permissions are being set")
			}

			if len(msg.DefaultBalances.UserPermissions.CanUpdateAutoApproveAllIncomingTransfers) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balance type is off-chain or non-public but default user permissions are being set")
			}
		}

		if msg.BalancesType == "Standard" || msg.BalancesType == "Non-Public" {
			if len(msg.OffChainBalancesMetadataTimeline) > 0 {
				return sdkerrors.Wrapf(ErrInvalidRequest, "balances metadata denotes on-chain balances but off-chain balances are set")
			}
		}

		if msg.BalancesType == "Off-Chain - Non-Indexed" {
			for _, offChainUriObj := range msg.OffChainBalancesMetadataTimeline {
				//must contain "{address}" in uri
				if !strings.Contains(offChainUriObj.OffChainBalancesMetadata.Uri, "{address}") {
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

	for _, path := range msg.CosmosCoinWrapperPathsToAdd {
		if path.Denom == "" {
			return sdkerrors.Wrapf(ErrInvalidRequest, "denom cannot be empty")
		}

		err = ValidateRangesAreValid(path.BadgeIds, false, true)
		if err != nil {
			return err
		}

		err = ValidateRangesAreValid(path.OwnershipTimes, false, true)
		if err != nil {
			return err
		}

		// Validate that only one denom unit per path has isDefaultDisplay set to true
		defaultDisplayCount := 0
		decimalsSet := make(map[string]bool)
		for _, denomUnit := range path.DenomUnits {
			if denomUnit.IsDefaultDisplay {
				defaultDisplayCount++
			}

			// Check that decimals is not 0
			if denomUnit.Decimals.IsZero() {
				return sdkerrors.Wrapf(ErrInvalidRequest, "denom unit decimals cannot be 0")
			}

			// Check for duplicate decimals
			decimalsStr := denomUnit.Decimals.String()
			if _, ok := decimalsSet[decimalsStr]; ok {
				return sdkerrors.Wrapf(ErrInvalidRequest, "duplicate denom unit decimals: %s", decimalsStr)
			}
			decimalsSet[decimalsStr] = true
		}
		if defaultDisplayCount > 1 {
			return sdkerrors.Wrapf(ErrInvalidRequest, "only one denom unit per path can have isDefaultDisplay set to true, found %d", defaultDisplayCount)
		}
	}

	return nil
}
