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

func TestMsgSetDynamicStoreValue_ValidateBasic_EmptyCreator(t *testing.T) {
	msg := &types.MsgSetDynamicStoreValue{
		Creator: "",
		StoreId: sdkmath.NewUint(1),
		Address: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		Value:   true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty creator should fail")
}

func TestMsgSetDynamicStoreValue_ValidateBasic_InvalidCreator(t *testing.T) {
	msg := &types.MsgSetDynamicStoreValue{
		Creator: "invalid_address",
		StoreId: sdkmath.NewUint(1),
		Address: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		Value:   true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator should fail")
}

func TestMsgSetDynamicStoreValue_ValidateBasic_ZeroStoreId(t *testing.T) {
	msg := &types.MsgSetDynamicStoreValue{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		StoreId: sdkmath.NewUint(0),
		Address: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		Value:   true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "zero store ID should fail")
	require.Contains(t, err.Error(), "cannot be zero", "error should mention cannot be zero")
}

func TestMsgSetDynamicStoreValue_ValidateBasic_EmptyAddress(t *testing.T) {
	msg := &types.MsgSetDynamicStoreValue{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		StoreId: sdkmath.NewUint(1),
		Address: "",
		Value:   true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty address should fail")
}

func TestMsgSetDynamicStoreValue_ValidateBasic_InvalidAddress(t *testing.T) {
	msg := &types.MsgSetDynamicStoreValue{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		StoreId: sdkmath.NewUint(1),
		Address: "invalid_address",
		Value:   true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid address should fail")
}

func TestMsgSetDynamicStoreValue_ValidateBasic_Valid(t *testing.T) {
	msg := &types.MsgSetDynamicStoreValue{
		Creator: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		StoreId: sdkmath.NewUint(1),
		Address: "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q",
		Value:   true,
	}

	err := msg.ValidateBasic()
	require.NoError(t, err, "valid message should pass")
}

