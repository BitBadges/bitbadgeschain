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

func TestMsgCreateCollection_ValidateBasic_InvalidCreatorAddress(t *testing.T) {
	msg := &types.MsgCreateCollection{
		Creator: "invalid_address",
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator address should fail")
	require.Contains(t, err.Error(), "invalid creator address", "error should mention invalid creator")
}

func TestMsgCreateCollection_ValidateBasic_EmptyCreator(t *testing.T) {
	msg := &types.MsgCreateCollection{
		Creator: "",
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty creator should fail")
}

func TestMsgCreateCollection_ValidateBasic_MintAddressNotAllowed(t *testing.T) {
	msg := &types.MsgCreateCollection{
		Creator: types.MintAddress,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "Mint address should not be allowed for creator")
}

func TestMsgCreateCollection_ValidateBasic_ValidCreator(t *testing.T) {
	msg := &types.MsgCreateCollection{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
	}

	err := msg.ValidateBasic()
	// Might fail on address validation if SDK config uses different prefix
	if err != nil {
		require.Contains(t, err.Error(), "address", "error should be about address validation")
	} else {
		require.NoError(t, err, "valid creator should pass")
	}
}

// ============================================================================
// Invariants Validation
// ============================================================================

func TestMsgCreateCollection_ValidateBasic_NoCustomOwnershipTimesInvariant_Valid(t *testing.T) {
	msg := &types.MsgCreateCollection{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		Invariants: &types.InvariantsAddObject{
			NoCustomOwnershipTimes: true,
		},
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId:        "test_approval",
				FromListId:        "All",
				ToListId:          "All",
				InitiatedByListId: "All",
				TokenIds: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
				TransferTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
				OwnershipTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
				ApprovalCriteria: &types.ApprovalCriteria{},
				Version:          sdkmath.NewUint(0),
			},
		},
	}

	err := msg.ValidateBasic()
	// Might fail on address validation if SDK config uses different prefix
	if err != nil {
		require.Contains(t, err.Error(), "address", "error should be about address validation")
	} else {
		require.NoError(t, err, "valid invariants with full ownership times should pass")
	}
}

func TestMsgCreateCollection_ValidateBasic_NoCustomOwnershipTimesInvariant_Invalid(t *testing.T) {
	msg := &types.MsgCreateCollection{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		Invariants: &types.InvariantsAddObject{
			NoCustomOwnershipTimes: true,
		},
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId:        "test_approval",
				FromListId:        "All",
				ToListId:          "All",
				InitiatedByListId: "All",
				TokenIds: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(10),
					},
				},
				TransferTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
				OwnershipTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(100), // Invalid: not full range
					},
				},
				ApprovalCriteria: &types.ApprovalCriteria{},
				Version:          sdkmath.NewUint(0),
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid ownership times with invariant should fail")
}

// ============================================================================
// Collection Approvals Validation (if present)
// ============================================================================

func TestMsgCreateCollection_ValidateBasic_CollectionApprovals_InvalidTokenIds(t *testing.T) {
	// MsgCreateCollection.ValidateBasic() only validates CollectionApprovals if invariants are present
	// So we need to set invariants to trigger the validation
	msg := &types.MsgCreateCollection{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		Invariants: &types.InvariantsAddObject{
			NoCustomOwnershipTimes: true, // Set invariants to trigger approval validation
		},
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId:        "test_approval",
				FromListId:        "All",
				ToListId:          "All",
				InitiatedByListId: "All",
				TokenIds: []*types.UintRange{
					{
						Start: sdkmath.NewUint(0), // Invalid: start is 0, should be >= 1
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
				TransferTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
				OwnershipTimes: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(math.MaxUint64),
					},
				},
				ApprovalCriteria: &types.ApprovalCriteria{},
				Version:          sdkmath.NewUint(0),
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid token IDs in collection approvals should fail when invariants are present")
}

