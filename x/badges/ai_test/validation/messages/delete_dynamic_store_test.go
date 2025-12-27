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

func TestMsgDeleteDynamicStore_ValidateBasic_EmptyCreator(t *testing.T) {
	msg := &types.MsgDeleteDynamicStore{
		Creator: "",
		StoreId: sdkmath.NewUint(1),
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty creator should fail")
}

func TestMsgDeleteDynamicStore_ValidateBasic_InvalidCreator(t *testing.T) {
	msg := &types.MsgDeleteDynamicStore{
		Creator: "invalid_address",
		StoreId: sdkmath.NewUint(1),
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator should fail")
}

func TestMsgDeleteDynamicStore_ValidateBasic_ZeroStoreId(t *testing.T) {
	msg := &types.MsgDeleteDynamicStore{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		StoreId: sdkmath.NewUint(0),
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "zero store ID should fail")
	require.Contains(t, err.Error(), "cannot be zero", "error should mention cannot be zero")
}

func TestMsgDeleteDynamicStore_ValidateBasic_Valid(t *testing.T) {
	msg := &types.MsgDeleteDynamicStore{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		StoreId: sdkmath.NewUint(1),
	}

	err := msg.ValidateBasic()
	require.NoError(t, err, "valid message should pass")
}

