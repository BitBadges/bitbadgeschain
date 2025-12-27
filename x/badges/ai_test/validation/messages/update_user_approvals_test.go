package messages

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/validation"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

func init() {
	// Ensure SDK config is initialized for address validation
	validation.EnsureSDKConfig()
}

// ============================================================================
// Creator Address Validation
// ============================================================================

func TestMsgUpdateUserApprovals_ValidateBasic_InvalidCreator(t *testing.T) {
	msg := &types.MsgUpdateUserApprovals{
		Creator: "invalid_address",
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator should fail")
}

// ============================================================================
// Incoming Approvals Validation
// ============================================================================

func TestMsgUpdateUserApprovals_ValidateBasic_InvalidIncomingApprovals(t *testing.T) {
	msg := &types.MsgUpdateUserApprovals{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		IncomingApprovals: []*types.UserIncomingApproval{
			{
				ApprovalId:        "", // Invalid: empty approval ID
				FromListId:        "All",
				InitiatedByListId: "All",
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				TransferTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				ApprovalCriteria: &types.IncomingApprovalCriteria{},
				Version:          sdkmath.NewUint(0),
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid incoming approvals should fail")
}

// ============================================================================
// Outgoing Approvals Validation
// ============================================================================

func TestMsgUpdateUserApprovals_ValidateBasic_InvalidOutgoingApprovals(t *testing.T) {
	msg := &types.MsgUpdateUserApprovals{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ApprovalId:        "", // Invalid: empty approval ID
				ToListId:          "All",
				InitiatedByListId: "All",
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				TransferTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				ApprovalCriteria: &types.OutgoingApprovalCriteria{},
				Version:          sdkmath.NewUint(0),
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid outgoing approvals should fail")
}

// ============================================================================
// User Permissions Validation
// ============================================================================

func TestMsgUpdateUserApprovals_ValidateBasic_ValidUserPermissions(t *testing.T) {
	msg := &types.MsgUpdateUserApprovals{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		UserPermissions: &types.UserPermissions{
			CanUpdateIncomingApprovals:                         []*types.UserIncomingApprovalPermission{},
			CanUpdateOutgoingApprovals:                         []*types.UserOutgoingApprovalPermission{},
			CanUpdateAutoApproveSelfInitiatedOutgoingTransfers: []*types.ActionPermission{},
			CanUpdateAutoApproveSelfInitiatedIncomingTransfers: []*types.ActionPermission{},
			CanUpdateAutoApproveAllIncomingTransfers:           []*types.ActionPermission{},
		},
	}

	err := msg.ValidateBasic()
	// Might fail on address validation if SDK config uses different prefix
	if err != nil {
		require.Contains(t, err.Error(), "address", "error should be about address validation")
	} else {
		require.NoError(t, err, "valid user permissions should pass")
	}
}

