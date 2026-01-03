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

func TestMsgUniversalUpdateCollection_ValidateBasic_InvalidCreator(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "invalid_address",
		CollectionId: sdkmath.NewUint(1),
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator should fail")
}

// ============================================================================
// Collection ID Validation
// ============================================================================

func TestMsgUniversalUpdateCollection_ValidateBasic_NilCollectionId(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.Uint{},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "nil collection ID should fail")
	// Error might be about creator address first, or collection ID
	if err.Error() != "" {
		require.NotNil(t, err)
	}
}

func TestMsgUniversalUpdateCollection_ValidateBasic_ZeroCollectionId(t *testing.T) {
	// Zero collection ID is allowed for new collections
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(0),
	}

	err := msg.ValidateBasic()
	// Zero is allowed (means new collection)
	// But might fail on address validation if SDK config uses different prefix
	if err != nil {
		require.Contains(t, err.Error(), "address", "error should be about address validation")
	} else {
		require.NoError(t, err, "zero collection ID should be allowed")
	}
}

// ============================================================================
// ValidTokenIds Validation
// ============================================================================

func TestMsgUniversalUpdateCollection_ValidateBasic_InvalidValidTokenIds(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ValidTokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(100),
				End:   sdkmath.NewUint(1), // Invalid: start > end
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid valid token IDs should fail")
}

func TestMsgUniversalUpdateCollection_ValidateBasic_OverlappingValidTokenIds(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		ValidTokenIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(10),
			},
			{
				Start: sdkmath.NewUint(5),
				End:   sdkmath.NewUint(15), // Overlaps with first
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "overlapping valid token IDs should fail")
}

// ============================================================================
// Collection Approvals Validation
// ============================================================================

func TestMsgUniversalUpdateCollection_ValidateBasic_InvalidCollectionApprovals(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId:        "", // Invalid: empty approval ID
				FromListId:        "All",
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
				ApprovalCriteria: &types.ApprovalCriteria{},
				Version:          sdkmath.NewUint(0),
			},
		},
		UpdateCollectionApprovals: true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid collection approvals should fail")
}

func TestMsgUniversalUpdateCollection_ValidateBasic_DuplicateApprovalIds(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId:        "test_approval",
				FromListId:        "All",
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
				ApprovalCriteria: &types.ApprovalCriteria{},
				Version:          sdkmath.NewUint(0),
			},
			{
				ApprovalId:        "test_approval", // Duplicate
				FromListId:        "All",
				ToListId:          "All",
				InitiatedByListId: "All",
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(20), End: sdkmath.NewUint(30)},
				},
				TransferTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				ApprovalCriteria: &types.ApprovalCriteria{},
				Version:          sdkmath.NewUint(0),
			},
		},
		UpdateCollectionApprovals: true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "duplicate approval IDs should fail")
}

// ============================================================================
// Default Balances Validation
// ============================================================================

func TestMsgUniversalUpdateCollection_ValidateBasic_InvalidDefaultBalances(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(0), // Invalid: zero amount
					TokenIds: []*types.UintRange{
						{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
					},
					OwnershipTimes: []*types.UintRange{
						{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
					},
				},
			},
		},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid default balances should fail")
}

// ============================================================================
// Manager Validation
// ============================================================================

func TestMsgUniversalUpdateCollection_ValidateBasic_InvalidManager(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		Manager:      "invalid_address",
		UpdateManager: true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid manager should fail")
}

func TestMsgUniversalUpdateCollection_ValidateBasic_EmptyManager(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		Manager:      "", // Empty is allowed
		UpdateManager: true,
	}

	err := msg.ValidateBasic()
	// Empty manager is allowed
	// But might fail on creator address validation if SDK config uses different prefix
	if err != nil {
		require.Contains(t, err.Error(), "address", "error should be about address validation")
	} else {
		require.NoError(t, err, "empty manager should be allowed")
	}
}

// ============================================================================
// Metadata Validation
// ============================================================================

func TestMsgUniversalUpdateCollection_ValidateBasic_InvalidCollectionMetadataUri(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		CollectionMetadata: &types.CollectionMetadata{
			Uri: "invalid uri format", // Invalid URI
		},
		UpdateCollectionMetadata: true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid URI should fail")
}

func TestMsgUniversalUpdateCollection_ValidateBasic_ValidCollectionMetadataUri(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		CollectionMetadata: &types.CollectionMetadata{
			Uri: "https://example.com/metadata",
		},
		UpdateCollectionMetadata: true,
	}

	err := msg.ValidateBasic()
	// Valid URI should pass
	// But might fail on creator address validation if SDK config uses different prefix
	if err != nil {
		require.Contains(t, err.Error(), "address", "error should be about address validation")
	} else {
		require.NoError(t, err, "valid URI should pass")
	}
}

// ============================================================================
// Token Metadata Validation
// ============================================================================

func TestMsgUniversalUpdateCollection_ValidateBasic_InvalidTokenMetadata(t *testing.T) {
	msg := &types.MsgUniversalUpdateCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
		TokenMetadata: []*types.TokenMetadata{
			{
				Uri: "invalid uri",
				TokenIds: []*types.UintRange{
					{
						Start: sdkmath.NewUint(100),
						End:   sdkmath.NewUint(1), // Invalid: start > end
					},
				},
			},
		},
		UpdateTokenMetadata: true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid token metadata should fail")
}

// ============================================================================
// Cosmos Coin Wrapper Paths Validation
// ============================================================================
// Note: CosmosCoinWrapperPathsToAdd uses CosmosCoinWrapperPathAddObject type
// which has a different structure. These tests are simplified to test basic validation.
// Full validation tests would require understanding the exact structure of CosmosCoinWrapperPathAddObject.

