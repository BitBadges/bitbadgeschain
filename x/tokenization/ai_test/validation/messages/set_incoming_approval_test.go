package messages

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/validation"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func init() {
	// Ensure SDK config is initialized for address validation
	validation.EnsureSDKConfig()
}

// ============================================================================
// Creator Address Validation
// ============================================================================

func TestMsgSetIncomingApproval_ValidateBasic_InvalidCreator(t *testing.T) {
	msg := &types.MsgSetIncomingApproval{
		Creator:      "invalid_address",
		CollectionId: sdkmath.NewUint(1),
		Approval: &types.UserIncomingApproval{
			ApprovalId:        "test_approval",
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
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator should fail")
}

// ============================================================================
// Collection ID Validation
// ============================================================================

func TestMsgSetIncomingApproval_ValidateBasic_NilCollectionId(t *testing.T) {
	msg := &types.MsgSetIncomingApproval{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.Uint{},
		Approval: &types.UserIncomingApproval{
			ApprovalId:        "test_approval",
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
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "nil collection ID should fail")
	require.Contains(t, err.Error(), "cannot be nil", "error should mention cannot be nil")
}

func TestMsgSetIncomingApproval_ValidateBasic_ZeroCollectionId(t *testing.T) {
	msg := &types.MsgSetIncomingApproval{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(0),
		Approval: &types.UserIncomingApproval{
			ApprovalId:        "test_approval",
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
	}

	// Collection ID 0 is allowed for auto-prev resolution (checks happen later)
	err := msg.ValidateBasic()
	require.NoError(t, err, "zero collection ID is allowed for auto-prev")
}

// ============================================================================
// Approval Validation
// ============================================================================

func TestMsgSetIncomingApproval_ValidateBasic_NilApproval(t *testing.T) {
	msg := &types.MsgSetIncomingApproval{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		Approval:     nil,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "nil approval should fail")
	require.Contains(t, err.Error(), "cannot be nil", "error should mention cannot be nil")
}

func TestMsgSetIncomingApproval_ValidateBasic_InvalidApproval(t *testing.T) {
	msg := &types.MsgSetIncomingApproval{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		Approval: &types.UserIncomingApproval{
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
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid approval should fail")
}

func TestMsgSetIncomingApproval_ValidateBasic_ValidApproval(t *testing.T) {
	msg := &types.MsgSetIncomingApproval{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		Approval: &types.UserIncomingApproval{
			ApprovalId:        "test_approval",
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
	}

	err := msg.ValidateBasic()
	require.NoError(t, err, "valid approval should pass")
}

