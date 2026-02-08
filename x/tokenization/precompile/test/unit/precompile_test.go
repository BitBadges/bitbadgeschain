package tokenization_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

func TestPrecompile_RequiredGas(t *testing.T) {
	precompile := createTestPrecompile(t)

	// Test with valid method ID - get the method selector manually
	methodID := precompile.ABI.Methods["transferTokens"].ID
	gas := precompile.RequiredGas(methodID[:])
	require.Equal(t, uint64(tokenization.GasTransferTokensBase), gas)

	// Test with invalid input (too short)
	gas = precompile.RequiredGas([]byte{0x12, 0x34})
	require.Equal(t, uint64(0), gas)
}

func TestPrecompile_TransferTokens_Structure(t *testing.T) {
	tokenizationKeeper, _ := keepertest.TokenizationKeeper(t)
	precompile := tokenization.NewPrecompile(tokenizationKeeper)

	// Verify precompile is created correctly
	require.NotNil(t, precompile)
	require.NotNil(t, precompile.ABI)
	// Note: tokenizationKeeper is unexported, so we can't check it directly in tokenization_test package
	require.Equal(t, tokenization.TokenizationPrecompileAddress, precompile.ContractAddress.Hex())

	// Verify the transferTokens method exists
	method, found := precompile.ABI.Methods["transferTokens"]
	require.True(t, found)
	require.NotNil(t, method)

	// Verify method signature
	require.Equal(t, 5, len(method.Inputs))  // collectionId, toAddresses, amount, tokenIds, ownershipTimes
	require.Equal(t, 1, len(method.Outputs)) // success bool
}

func createTestPrecompile(t *testing.T) *tokenization.Precompile {
	tokenizationKeeper, _ := keepertest.TokenizationKeeper(t)
	return tokenization.NewPrecompile(tokenizationKeeper)
}
