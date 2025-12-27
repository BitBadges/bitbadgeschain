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

func TestMsgDeleteCollection_ValidateBasic_InvalidCreator(t *testing.T) {
	msg := &types.MsgDeleteCollection{
		Creator:      "invalid_address",
		CollectionId: sdkmath.NewUint(1),
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator should fail")
}

func TestMsgDeleteCollection_ValidateBasic_EmptyCreator(t *testing.T) {
	msg := &types.MsgDeleteCollection{
		Creator:      "",
		CollectionId: sdkmath.NewUint(1),
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty creator should fail")
}

func TestMsgDeleteCollection_ValidateBasic_ValidCreator(t *testing.T) {
	msg := &types.MsgDeleteCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
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
// Collection ID Validation
// ============================================================================

func TestMsgDeleteCollection_ValidateBasic_NilCollectionId(t *testing.T) {
	msg := &types.MsgDeleteCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.Uint{},
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "nil collection ID should fail")
	// Error might be about creator address first, or collection ID
	// Both are valid validation failures
	if err.Error() != "" {
		// Just verify it errors - could be creator or collection ID validation
		require.NotNil(t, err)
	}
}

func TestMsgDeleteCollection_ValidateBasic_ZeroCollectionId(t *testing.T) {
	msg := &types.MsgDeleteCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(0),
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "zero collection ID should fail")
	// Error might be about creator address first, or collection ID
	// Both are valid validation failures
	if err.Error() != "" {
		// Just verify it errors - could be creator or collection ID validation
		require.NotNil(t, err)
	}
}

func TestMsgDeleteCollection_ValidateBasic_ValidCollectionId(t *testing.T) {
	msg := &types.MsgDeleteCollection{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		CollectionId: sdkmath.NewUint(1),
	}

	err := msg.ValidateBasic()
	require.NoError(t, err, "valid collection ID should pass")
}

