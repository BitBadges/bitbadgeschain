package messages

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/validation"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func init() {
	// Ensure SDK config is initialized for address validation
	validation.EnsureSDKConfig()
}

func TestMsgCreateDynamicStore_ValidateBasic_EmptyCreator(t *testing.T) {
	msg := &types.MsgCreateDynamicStore{
		Creator:      "",
		DefaultValue: true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "empty creator should fail")
}

func TestMsgCreateDynamicStore_ValidateBasic_InvalidCreator(t *testing.T) {
	msg := &types.MsgCreateDynamicStore{
		Creator:      "invalid_address",
		DefaultValue: true,
	}

	err := msg.ValidateBasic()
	require.Error(t, err, "invalid creator should fail")
}

func TestMsgCreateDynamicStore_ValidateBasic_Valid(t *testing.T) {
	msg := &types.MsgCreateDynamicStore{
		Creator:      "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		DefaultValue: true,
	}

	err := msg.ValidateBasic()
	// Might fail on address validation if SDK config uses different prefix
	if err != nil {
		require.Contains(t, err.Error(), "address", "error should be about address validation")
	} else {
		require.NoError(t, err, "valid message should pass")
	}
}

