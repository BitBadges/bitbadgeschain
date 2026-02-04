package tokenization

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

func TestPrecompile_RequiredGas(t *testing.T) {
	precompile := createTestPrecompile(t)

	// Test with valid method ID - get the method selector manually
	methodID := precompile.ABI.Methods["transferTokens"].ID
	gas := precompile.RequiredGas(methodID[:])
	require.Equal(t, uint64(GasTransferTokens), gas)

	// Test with invalid input (too short)
	gas = precompile.RequiredGas([]byte{0x12, 0x34})
	require.Equal(t, uint64(0), gas)
}

func TestPrecompile_TransferTokens_Structure(t *testing.T) {
	tokenizationKeeper, _ := keepertest.TokenizationKeeper(t)
	precompile := NewPrecompile(tokenizationKeeper)

	// Verify precompile is created correctly
	require.NotNil(t, precompile)
	require.NotNil(t, precompile.ABI)
	require.NotNil(t, precompile.tokenizationKeeper)
	require.Equal(t, TokenizationPrecompileAddress, precompile.ContractAddress.Hex())

	// Verify the transferTokens method exists
	method, found := precompile.ABI.Methods["transferTokens"]
	require.True(t, found)
	require.NotNil(t, method)

	// Verify method signature
	require.Equal(t, 5, len(method.Inputs))  // collectionId, toAddresses, amount, tokenIds, ownershipTimes
	require.Equal(t, 1, len(method.Outputs)) // success bool
}

func createTestPrecompile(t *testing.T) *Precompile {
	tokenizationKeeper, _ := keepertest.TokenizationKeeper(t)
	return NewPrecompile(tokenizationKeeper)
}
