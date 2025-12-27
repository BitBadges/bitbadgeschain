package messages

import (
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

func TestMsgPurgeApprovals_ValidateBasic_InvalidCreator(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "invalid_address",
		CollectionId: sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator should fail")
}

// ============================================================================
// Collection ID Validation
// ============================================================================

func TestMsgPurgeApprovals_ValidateBasic_ZeroCollectionId(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(0),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "zero collection ID should fail")
}

// ============================================================================
// ApprovalsToPurge Validation
// ============================================================================

func TestMsgPurgeApprovals_ValidateBasic_EmptyApprovalsToPurge(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:          "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId:     sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty approvalsToPurge should fail")
	require.Contains(t, err.Error(), "cannot be empty", "error should mention cannot be empty")
}

func TestMsgPurgeApprovals_ValidateBasic_NilApprovalInArray(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			nil, // Invalid: nil approval
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "nil approval should fail")
	require.Contains(t, err.Error(), "cannot be nil", "error should mention cannot be nil")
}

func TestMsgPurgeApprovals_ValidateBasic_EmptyApprovalId(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "", // Invalid: empty
				ApprovalLevel:   "collection",
				ApproverAddress: "",
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty approval ID should fail")
}

func TestMsgPurgeApprovals_ValidateBasic_EmptyApprovalLevel(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "", // Invalid: empty
				ApproverAddress: "",
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty approval level should fail")
}

func TestMsgPurgeApprovals_ValidateBasic_InvalidApprovalLevel(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "invalid_level", // Invalid
				ApproverAddress: "",
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid approval level should fail")
	require.Contains(t, err.Error(), "must be", "error should mention valid levels")
}

// ============================================================================
// Approver Address Validation
// ============================================================================

func TestMsgPurgeApprovals_ValidateBasic_CollectionLevelWithApproverAddress(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "collection",
				ApproverAddress: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430", // Invalid: should be empty
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "collection level with approver address should fail")
	require.Contains(t, err.Error(), "must be empty", "error should mention must be empty")
}

func TestMsgPurgeApprovals_ValidateBasic_IncomingLevelWithoutApproverAddress(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "incoming",
				ApproverAddress: "", // Invalid: should be provided
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "incoming level without approver address should fail")
	require.Contains(t, err.Error(), "must be provided", "error should mention must be provided")
}

func TestMsgPurgeApprovals_ValidateBasic_OutgoingLevelWithoutApproverAddress(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "outgoing",
				ApproverAddress: "", // Invalid: should be provided
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "outgoing level without approver address should fail")
}

func TestMsgPurgeApprovals_ValidateBasic_InvalidApproverAddress(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "incoming",
				ApproverAddress: "invalid_address",
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid approver address should fail")
}

// ============================================================================
// Business Logic Validation
// ============================================================================

func TestMsgPurgeApprovals_ValidateBasic_PurgingOwnApprovals_PurgeExpiredFalse(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:         "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId:    sdkmath.NewUint(1),
		PurgeExpired:    false, // Invalid: should be true when purging own
		ApproverAddress: "",    // Empty means purging own
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "incoming",
				ApproverAddress: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "purgeExpired should be true when purging own approvals")
	require.Contains(t, err.Error(), "must be true", "error should mention must be true")
}

func TestMsgPurgeApprovals_ValidateBasic_PurgingOwnApprovals_PurgeCounterpartyTrue(t *testing.T) {
	msg := &types.MsgPurgeApprovals{
		Creator:                    "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId:               sdkmath.NewUint(1),
		PurgeExpired:               true,
		PurgeCounterpartyApprovals: true, // Invalid: should be false when purging own
		ApproverAddress:            "",
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "incoming",
				ApproverAddress: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "purgeCounterpartyApprovals should be false when purging own")
	require.Contains(t, err.Error(), "must be false", "error should mention must be false")
}

func TestMsgPurgeApprovals_ValidateBasic_ValidCollectionLevel(t *testing.T) {
	// For collection-level, approverAddress must be empty in the approval
	// If msg.ApproverAddress is empty, targetAddress = creator, so purgeExpired must be true
	msg := &types.MsgPurgeApprovals{
		Creator:         "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId:    sdkmath.NewUint(1),
		PurgeExpired:    true, // Required when targetAddress == creator
		ApproverAddress: "",   // Empty means targetAddress = creator
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "collection",
				ApproverAddress: "", // Must be empty for collection-level
			},
		},
	}

	err := msg.ValidateBasic()
	require.NoError(t, err, "valid collection level purge should pass")
}

func TestMsgPurgeApprovals_ValidateBasic_ValidIncomingLevel(t *testing.T) {
	// When msg.ApproverAddress is different from creator, targetAddress = msg.ApproverAddress
	// So it's purging someone else's approvals, which doesn't require purgeExpired=true
	msg := &types.MsgPurgeApprovals{
		Creator:         "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId:    sdkmath.NewUint(1),
		ApproverAddress: "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", // Different from creator
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "incoming",
				ApproverAddress: "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", // Must match msg.ApproverAddress
			},
		},
	}

	err := msg.ValidateBasic()
	require.NoError(t, err, "valid incoming level purge should pass")
}

func TestMsgPurgeApprovals_ValidateBasic_ValidOutgoingLevel(t *testing.T) {
	// When msg.ApproverAddress is different from creator, targetAddress = msg.ApproverAddress
	// So it's purging someone else's approvals, which doesn't require purgeExpired=true
	msg := &types.MsgPurgeApprovals{
		Creator:         "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId:    sdkmath.NewUint(1),
		ApproverAddress: "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", // Different from creator
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test_approval",
				ApprovalLevel:   "outgoing",
				ApproverAddress: "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q", // Must match msg.ApproverAddress
			},
		},
	}

	err := msg.ValidateBasic()
	require.NoError(t, err, "valid outgoing level purge should pass")
}
