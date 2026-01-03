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

func TestMsgDeleteOutgoingApproval_ValidateBasic_InvalidCreator(t *testing.T) {
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      "invalid_address",
		CollectionId: sdkmath.NewUint(1),
		ApprovalId:   "test_approval",
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator should fail")
}

func TestMsgDeleteOutgoingApproval_ValidateBasic_NilCollectionId(t *testing.T) {
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.Uint{},
		ApprovalId:   "test_approval",
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "nil collection ID should fail")
}

func TestMsgDeleteOutgoingApproval_ValidateBasic_ZeroCollectionId(t *testing.T) {
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(0),
		ApprovalId:   "test_approval",
	}

	// Collection ID 0 is allowed for auto-prev resolution (checks happen later)
	err := msg.ValidateBasic()
	require.NoError(t, err, "zero collection ID is allowed for auto-prev")
}

func TestMsgDeleteOutgoingApproval_ValidateBasic_EmptyApprovalId(t *testing.T) {
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalId:   "",
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty approval ID should fail")
	require.Contains(t, err.Error(), "cannot be empty", "error should mention cannot be empty")
}

func TestMsgDeleteOutgoingApproval_ValidateBasic_Valid(t *testing.T) {
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ApprovalId:   "test_approval",
	}

	err := msg.ValidateBasic()
	require.NoError(t, err, "valid message should pass")
}

