package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ValidateApproval validates collection approvals (non-timeline version)
func ValidateApproval(ctx sdk.Context, approvals []*CollectionApproval, canChangeValues bool) error {
	return ValidateCollectionApprovals(ctx, approvals, canChangeValues)
}

// ValidateCollectionMetadata validates collection metadata
func ValidateCollectionMetadata(metadata *CollectionMetadata) error {
	if metadata == nil {
		return nil // nil is allowed when not updating metadata
	}
	return ValidateURI(metadata.Uri)
}

// ValidateStandards validates standards (non-timeline version)
func ValidateStandards(standards []string) error {
	// Standards is just a list of strings, no validation needed
	return nil
}

// ValidateCustomData validates custom data (non-timeline version)
func ValidateCustomData(customData string) error {
	// Custom data is just a string, no validation needed
	return nil
}

// ValidateManager validates a manager address
func ValidateManager(manager string) error {
	if len(manager) == 0 {
		return nil // empty manager is allowed
	}
	_, err := sdk.AccAddressFromBech32(manager)
	return err
}

// ValidateIsArchived validates isArchived (non-timeline version)
func ValidateIsArchived(isArchived bool) error {
	// IsArchived is just a bool, no validation needed
	return nil
}

// ValidateUserOutgoingApproval validates user outgoing approvals (non-timeline version)
func ValidateUserOutgoingApproval(ctx sdk.Context, approvals []*UserOutgoingApproval, address string, canChangeValues bool) error {
	return ValidateUserOutgoingApprovals(ctx, approvals, address, canChangeValues)
}

// ValidateUserIncomingApproval validates user incoming approvals (non-timeline version)
func ValidateUserIncomingApproval(ctx sdk.Context, approvals []*UserIncomingApproval, address string, canChangeValues bool) error {
	return ValidateUserIncomingApprovals(ctx, approvals, address, canChangeValues)
}
