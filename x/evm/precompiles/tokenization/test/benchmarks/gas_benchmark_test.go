package tokenization

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	tokenization "github.com/bitbadges/bitbadgeschain/x/evm/precompiles/tokenization"
	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"
)

// BenchmarkTransferGas benchmarks gas calculation for transfers
func BenchmarkTransferGas(b *testing.B) {
	_, _ = keepertest.TokenizationKeeper(b) // Setup keeper for context

	testCases := []struct {
		name        string
		recipients  int
		tokenRanges int
		timeRanges  int
	}{
		{"single_recipient", 1, 1, 1},
		{"multiple_recipients", 10, 1, 1},
		{"many_ranges", 1, 10, 10},
		{"complex", 10, 10, 10},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			toAddresses := make([]common.Address, tc.recipients)
			for i := range toAddresses {
				toAddresses[i] = common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
			}

			tokenIdsRanges := make([]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}, tc.tokenRanges)
			for i := range tokenIdsRanges {
				tokenIdsRanges[i] = struct {
					Start *big.Int `json:"start"`
					End   *big.Int `json:"end"`
				}{
					Start: big.NewInt(1),
					End:   big.NewInt(100),
				}
			}

			ownershipTimesRanges := make([]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}, tc.timeRanges)
			for i := range ownershipTimesRanges {
				ownershipTimesRanges[i] = struct {
					Start *big.Int `json:"start"`
					End   *big.Int `json:"end"`
				}{
					Start: big.NewInt(1),
					End:   big.NewInt(1000000),
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				gas := tokenization.CalculateTransferGas(toAddresses, tokenIdsRanges, ownershipTimesRanges)
				require.Greater(b, gas, uint64(0))
			}
		})
	}
}

// TestGasCalculationAccuracy tests that gas calculations are reasonable
func TestGasCalculationAccuracy(t *testing.T) {
	tests := []struct {
		name           string
		recipients     int
		tokenRanges    int
		timeRanges     int
		expectedMinGas uint64
		expectedMaxGas uint64
	}{
		{
			name:           "minimal_transfer",
			recipients:     1,
			tokenRanges:    1,
			timeRanges:     1,
			expectedMinGas: tokenization.GasTransferTokensBase,
			expectedMaxGas: tokenization.GasTransferTokensBase + 10_000, // Allow some overhead
		},
		{
			name:           "moderate_transfer",
			recipients:     5,
			tokenRanges:    5,
			timeRanges:     5,
			expectedMinGas: tokenization.GasTransferTokensBase + 25_000, // 5*5000 + 5*1000 + 5*1000
			expectedMaxGas: tokenization.GasTransferTokensBase + 50_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toAddresses := make([]common.Address, tt.recipients)
			for i := range toAddresses {
				toAddresses[i] = common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
			}

			tokenIdsRanges := make([]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}, tt.tokenRanges)
			for i := range tokenIdsRanges {
				tokenIdsRanges[i] = struct {
					Start *big.Int `json:"start"`
					End   *big.Int `json:"end"`
				}{
					Start: big.NewInt(1),
					End:   big.NewInt(100),
				}
			}

			ownershipTimesRanges := make([]struct {
				Start *big.Int `json:"start"`
				End   *big.Int `json:"end"`
			}, tt.timeRanges)
			for i := range ownershipTimesRanges {
				ownershipTimesRanges[i] = struct {
					Start *big.Int `json:"start"`
					End   *big.Int `json:"end"`
				}{
					Start: big.NewInt(1),
					End:   big.NewInt(1000000),
				}
			}

			gas := tokenization.CalculateTransferGas(toAddresses, tokenIdsRanges, ownershipTimesRanges)
			require.GreaterOrEqual(t, gas, tt.expectedMinGas, "Gas should be at least the minimum")
			require.LessOrEqual(t, gas, tt.expectedMaxGas, "Gas should not exceed the maximum")
		})
	}
}

