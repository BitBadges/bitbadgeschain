package gamm

import (
	"testing"

	"github.com/stretchr/testify/require"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
)

// BenchmarkCalculateDynamicGas benchmarks dynamic gas calculation with varying inputs
func BenchmarkCalculateDynamicGas(b *testing.B) {
	testCases := []struct {
		name          string
		baseGas       uint64
		numRoutes     int
		numCoins      int
		numAffiliates int
	}{
		{"minimal", gamm.GasJoinPoolBase, 0, 2, 0},
		{"with_routes", gamm.GasSwapExactAmountInBase, 5, 0, 0},
		{"with_coins", gamm.GasJoinPoolBase, 0, 10, 0},
		{"with_affiliates", gamm.GasSwapExactAmountInBase, 0, 0, 5},
		{"complex", gamm.GasSwapExactAmountInBase, 10, 20, 10},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				gas := gamm.CalculateDynamicGas(tc.baseGas, tc.numRoutes, tc.numCoins, tc.numAffiliates)
				require.Greater(b, gas, uint64(0))
			}
		})
	}
}

// BenchmarkJoinPoolGas benchmarks gas calculation for joinPool operations
func BenchmarkJoinPoolGas(b *testing.B) {
	testCases := []struct {
		name     string
		numCoins int
	}{
		{"base_case", 2},
		{"moderate", 10},
		{"maximum", 20},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			baseGas := uint64(gamm.GasJoinPoolBase)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				gas := gamm.CalculateDynamicGas(baseGas, 0, tc.numCoins, 0)
				require.Greater(b, gas, uint64(0))
			}
		})
	}
}

// BenchmarkSwapGas benchmarks gas calculation for swap operations
func BenchmarkSwapGas(b *testing.B) {
	testCases := []struct {
		name          string
		numRoutes     int
		numAffiliates int
	}{
		{"single_route", 1, 0},
		{"multiple_routes", 5, 0},
		{"max_routes", 10, 0},
		{"with_affiliates", 5, 5},
		{"max_with_affiliates", 10, 10},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			baseGas := uint64(gamm.GasSwapExactAmountInBase)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				gas := gamm.CalculateDynamicGas(baseGas, tc.numRoutes, 0, tc.numAffiliates)
				require.Greater(b, gas, uint64(0))
			}
		})
	}
}

// BenchmarkExitPoolGas benchmarks gas calculation for exitPool operations
func BenchmarkExitPoolGas(b *testing.B) {
	testCases := []struct {
		name     string
		numCoins int
	}{
		{"base_case", 2},
		{"moderate", 10},
		{"maximum", 20},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			baseGas := uint64(gamm.GasExitPoolBase)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				gas := gamm.CalculateDynamicGas(baseGas, 0, tc.numCoins, 0)
				require.Greater(b, gas, uint64(0))
			}
		})
	}
}

// TestGasCalculationAccuracy tests that gas calculations are within expected ranges
func TestGasCalculationAccuracy(t *testing.T) {
	tests := []struct {
		name           string
		baseGas        uint64
		numRoutes      int
		numCoins       int
		numAffiliates  int
		expectedMinGas uint64
		expectedMaxGas uint64
	}{
		{
			name:           "minimal_join",
			baseGas:        gamm.GasJoinPoolBase,
			numRoutes:      0,
			numCoins:       2,
			numAffiliates:  0,
			expectedMinGas: gamm.GasJoinPoolBase,
			expectedMaxGas: gamm.GasJoinPoolBase + 5_000, // Allow some overhead
		},
		{
			name:           "moderate_swap",
			baseGas:        gamm.GasSwapExactAmountInBase,
			numRoutes:      5,
			numCoins:       0,
			numAffiliates:  0,
			expectedMinGas: gamm.GasSwapExactAmountInBase + 20_000, // 5 * GasPerRoute
			expectedMaxGas: gamm.GasSwapExactAmountInBase + 30_000, // Allow overhead
		},
		{
			name:           "complex_swap",
			baseGas:        gamm.GasSwapExactAmountInBase,
			numRoutes:      10,
			numCoins:       0,
			numAffiliates:  10,
			expectedMinGas: gamm.GasSwapExactAmountInBase + 50_000 + 30_000, // 10*GasPerRoute + 10*GasPerAffiliate
			expectedMaxGas: gamm.GasSwapExactAmountInBase + 100_000,        // Allow overhead
		},
		{
			name:           "max_coins_join",
			baseGas:        gamm.GasJoinPoolBase,
			numRoutes:      0,
			numCoins:       20,
			numAffiliates:  0,
			expectedMinGas: gamm.GasJoinPoolBase + 35_000, // 20 * GasPerCoin
			expectedMaxGas: gamm.GasJoinPoolBase + 50_000, // Allow overhead
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gas := gamm.CalculateDynamicGas(tt.baseGas, tt.numRoutes, tt.numCoins, tt.numAffiliates)
			require.GreaterOrEqual(t, gas, tt.expectedMinGas, "Gas should be at least the minimum")
			require.LessOrEqual(t, gas, tt.expectedMaxGas, "Gas should not exceed the maximum")
		})
	}
}

// TestGasConstants tests that gas constants are properly defined
func TestGasConstants(t *testing.T) {
	// Verify transaction gas constants
	require.Greater(t, uint64(gamm.GasJoinPoolBase), uint64(0), "GasJoinPoolBase should be defined")
	require.Greater(t, uint64(gamm.GasExitPoolBase), uint64(0), "GasExitPoolBase should be defined")
	require.Greater(t, uint64(gamm.GasSwapExactAmountInBase), uint64(0), "GasSwapExactAmountInBase should be defined")
	require.Greater(t, uint64(gamm.GasSwapExactAmountInWithIBCTransferBase), uint64(0), "GasSwapExactAmountInWithIBCTransferBase should be defined")

	// Verify dynamic gas constants
	require.Greater(t, uint64(gamm.GasPerRoute), uint64(0), "GasPerRoute should be defined")
	require.Greater(t, uint64(gamm.GasPerCoin), uint64(0), "GasPerCoin should be defined")
	require.Greater(t, uint64(gamm.GasPerAffiliate), uint64(0), "GasPerAffiliate should be defined")

	// Verify query gas constants
	require.GreaterOrEqual(t, uint64(gamm.GasGetPoolBase), uint64(0), "GasGetPoolBase should be defined")
	require.GreaterOrEqual(t, uint64(gamm.GasGetPoolsBase), uint64(0), "GasGetPoolsBase should be defined")
	require.GreaterOrEqual(t, uint64(gamm.GasGetPoolTypeBase), uint64(0), "GasGetPoolTypeBase should be defined")
}

// BenchmarkGasCalculationWithLargeInputs benchmarks gas calculation with maximum allowed inputs
func BenchmarkGasCalculationWithLargeInputs(b *testing.B) {
	testCases := []struct {
		name          string
		baseGas       uint64
		numRoutes     int
		numCoins      int
		numAffiliates int
	}{
		{"max_routes", gamm.GasSwapExactAmountInBase, 10, 0, 0},
		{"max_coins", gamm.GasJoinPoolBase, 0, 20, 0},
		{"max_affiliates", gamm.GasSwapExactAmountInBase, 0, 0, 10},
		{"all_max", gamm.GasSwapExactAmountInBase, 10, 20, 10},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				gas := gamm.CalculateDynamicGas(tc.baseGas, tc.numRoutes, tc.numCoins, tc.numAffiliates)
				require.Greater(b, gas, uint64(0))
			}
		})
	}
}

