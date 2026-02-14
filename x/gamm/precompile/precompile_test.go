package gamm

import (
	"testing"

	"github.com/stretchr/testify/require"

	gammkeeper "github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
)

func TestNewPrecompile(t *testing.T) {
	// Create a mock keeper (this would need to be properly initialized in real tests)
	// For now, this is a placeholder test structure
	var keeper gammkeeper.Keeper

	precompile := NewPrecompile(keeper)
	require.NotNil(t, precompile)
	require.Equal(t, GammPrecompileAddress, precompile.ContractAddress.Hex())
}

func TestGetABILoadError(t *testing.T) {
	// Test that ABI loading error can be retrieved
	err := GetABILoadError()
	// ABI should load successfully, so error should be nil
	require.NoError(t, err)
}

func TestPrecompileAddress(t *testing.T) {
	// Verify the precompile address is correct
	require.Equal(t, "0x0000000000000000000000000000000000001002", GammPrecompileAddress)
}

