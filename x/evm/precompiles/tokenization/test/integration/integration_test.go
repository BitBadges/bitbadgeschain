package tokenization_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tokenization "github.com/bitbadges/bitbadgeschain/x/evm/precompiles/tokenization"
	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

// TestIntegration_PrecompileSetup verifies the precompile can be instantiated
// Full integration tests with actual transfers require complete collection setup
// which should be done in the tokenization module's integration test suite
func TestIntegration_PrecompileSetup(t *testing.T) {
	tokenizationKeeper, ctx := keepertest.TokenizationKeeper(t)
	precompile := tokenization.NewPrecompile(tokenizationKeeper)

	require.NotNil(t, precompile)
	require.NotNil(t, ctx)

	// Verify the precompile has the correct address
	require.Equal(t, tokenization.TokenizationPrecompileAddress, precompile.ContractAddress.Hex())

	// Verify ABI is loaded
	require.NotNil(t, precompile.ABI)
	method, found := precompile.ABI.Methods["transferTokens"]
	require.True(t, found)
	require.NotNil(t, method)
}
