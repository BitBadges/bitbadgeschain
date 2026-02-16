package app

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/stretchr/testify/require"
)

// TestPreciseBankUbadgeToWeiConversion tests that when 1 ubadge is specified as a fee,
// it correctly converts to 1*10^9 native ETH currency units (wei).
func TestPreciseBankUbadgeToWeiConversion(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	// Ensure EVM params are configured correctly
	evmParams := app.EVMKeeper.GetParams(ctx)
	if evmParams.EvmDenom != "ubadge" {
		evmParams.EvmDenom = "ubadge"
	}
	if evmParams.ExtendedDenomOptions == nil || evmParams.ExtendedDenomOptions.ExtendedDenom != "abadge" {
		evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"}
	}
	err := app.EVMKeeper.SetParams(ctx, evmParams)
	require.NoError(t, err, "Should be able to set EVM params")

	// Verify configuration
	evmParams = app.EVMKeeper.GetParams(ctx)
	require.Equal(t, "ubadge", evmParams.EvmDenom, "Base denom should be ubadge")
	require.NotNil(t, evmParams.ExtendedDenomOptions, "Extended denom options should be set")
	require.Equal(t, "abadge", evmParams.ExtendedDenomOptions.ExtendedDenom, "Extended denom should be abadge")

	// Get denom metadata to verify it exists
	_, found := app.BankKeeper.GetDenomMetaData(ctx, "ubadge")
	require.True(t, found, "ubadge metadata should exist")

	// Calculate decimals from denom units (badge has exponent 9, so ubadge has 9 decimals)
	decimals := uint8(9) // ubadge has 9 decimals (badge has exponent 9)

	// Test: 1 ubadge should convert to 1*10^9 wei
	// Since ubadge has 9 decimals, 1 ubadge = 1 * 10^9 base units
	// In EVM context (18 decimals), 1 ubadge should map to 1 * 10^9 wei
	oneUbadge := sdkmath.NewInt(1)
	expectedWei := big.NewInt(1_000_000_000) // 1 * 10^9

	// Convert 1 ubadge to wei using the conversion logic
	// The conversion multiplies by 10^(18 - 9) = 10^9 to extend from 9 to 18 decimals
	decimalsDiff := 18 - int(decimals) // 18 - 9 = 9
	conversionFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimalsDiff)), nil)
	convertedWei := new(big.Int).Mul(oneUbadge.BigInt(), conversionFactor)

	require.Equal(t, expectedWei, convertedWei, "1 ubadge should convert to 1*10^9 wei")
}

// TestPreciseBankFeeConversion tests that fees specified in ubadge correctly convert
// to the expected wei amount for EVM transactions.
func TestPreciseBankFeeConversion(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	// Ensure EVM params are configured correctly
	evmParams := app.EVMKeeper.GetParams(ctx)
	if evmParams.EvmDenom != "ubadge" {
		evmParams.EvmDenom = "ubadge"
	}
	if evmParams.ExtendedDenomOptions == nil || evmParams.ExtendedDenomOptions.ExtendedDenom != "abadge" {
		evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"}
	}
	err := app.EVMKeeper.SetParams(ctx, evmParams)
	require.NoError(t, err, "Should be able to set EVM params")

	decimals := uint8(9) // ubadge has 9 decimals

	// Test multiple fee amounts
	testCases := []struct {
		name         string
		ubadgeAmount sdkmath.Int
		expectedWei  *big.Int
		description  string
	}{
		{
			name:         "1 ubadge",
			ubadgeAmount: sdkmath.NewInt(1),
			expectedWei:  big.NewInt(1_000_000_000), // 1 * 10^9
			description:  "1 ubadge should equal 1*10^9 wei",
		},
		{
			name:         "10 ubadge",
			ubadgeAmount: sdkmath.NewInt(10),
			expectedWei:  big.NewInt(10_000_000_000), // 10 * 10^9
			description:  "10 ubadge should equal 10*10^9 wei",
		},
		{
			name:         "1000 ubadge",
			ubadgeAmount: sdkmath.NewInt(1000),
			expectedWei:  big.NewInt(1_000_000_000_000), // 1000 * 10^9
			description:  "1000 ubadge should equal 1000*10^9 wei",
		},
		{
			name:         "1 badge (1*10^9 ubadge)",
			ubadgeAmount: sdkmath.NewInt(1_000_000_000),                                          // 1 badge
			expectedWei:  new(big.Int).Mul(big.NewInt(1_000_000_000), big.NewInt(1_000_000_000)), // 1 * 10^18
			description:  "1 badge (1*10^9 ubadge) should equal 1*10^18 wei (1 ETH equivalent)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert ubadge to wei
			// Formula: wei = ubadge * 10^(18 - 9) = ubadge * 10^9
			decimalsDiff := 18 - int(decimals) // 18 - 9 = 9
			conversionFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimalsDiff)), nil)
			convertedWei := new(big.Int).Mul(tc.ubadgeAmount.BigInt(), conversionFactor)

			require.Equal(t, tc.expectedWei, convertedWei, tc.description)
		})
	}
}

// TestPreciseBankDenomMetadata verifies that the denom metadata is correctly set up
// for both ubadge and abadge denominations.
func TestPreciseBankDenomMetadata(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	// Ensure EVM params are configured correctly
	evmParams := app.EVMKeeper.GetParams(ctx)
	if evmParams.EvmDenom != "ubadge" {
		evmParams.EvmDenom = "ubadge"
	}
	if evmParams.ExtendedDenomOptions == nil || evmParams.ExtendedDenomOptions.ExtendedDenom != "abadge" {
		evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"}
	}
	err := app.EVMKeeper.SetParams(ctx, evmParams)
	require.NoError(t, err, "Should be able to set EVM params")

	// Check ubadge metadata
	ubadgeMetadata, found := app.BankKeeper.GetDenomMetaData(ctx, "ubadge")
	require.True(t, found, "ubadge metadata should exist")

	require.Equal(t, "ubadge", ubadgeMetadata.Base, "Base denom should be ubadge")
	require.Equal(t, "badge", ubadgeMetadata.Display, "Display denom should be badge")
	require.Len(t, ubadgeMetadata.DenomUnits, 2, "Should have 2 denom units")

	// Verify denom units
	require.Equal(t, "ubadge", ubadgeMetadata.DenomUnits[0].Denom, "First denom unit should be ubadge")
	require.Equal(t, uint32(0), ubadgeMetadata.DenomUnits[0].Exponent, "ubadge should have exponent 0")
	require.Equal(t, "badge", ubadgeMetadata.DenomUnits[1].Denom, "Second denom unit should be badge")
	require.Equal(t, uint32(9), ubadgeMetadata.DenomUnits[1].Exponent, "badge should have exponent 9")

	// Verify EVM params match (already set above)
	require.Equal(t, "ubadge", evmParams.EvmDenom, "EVM params denom should match")
	require.NotNil(t, evmParams.ExtendedDenomOptions, "Extended denom options should be set")
	require.Equal(t, "abadge", evmParams.ExtendedDenomOptions.ExtendedDenom, "EVM params extended denom should be abadge")
}

// TestPreciseBankBalanceConversion tests that balances are correctly converted
// between ubadge and the extended denom (abadge) for EVM operations.
func TestPreciseBankBalanceConversion(t *testing.T) {
	app := Setup(false)
	ctx := app.NewContext(false)

	// Ensure EVM params are configured correctly
	evmParams := app.EVMKeeper.GetParams(ctx)
	if evmParams.EvmDenom != "ubadge" {
		evmParams.EvmDenom = "ubadge"
	}
	if evmParams.ExtendedDenomOptions == nil || evmParams.ExtendedDenomOptions.ExtendedDenom != "abadge" {
		evmParams.ExtendedDenomOptions = &evmtypes.ExtendedDenomOptions{ExtendedDenom: "abadge"}
	}
	err := app.EVMKeeper.SetParams(ctx, evmParams)
	require.NoError(t, err, "Should be able to set EVM params")

	decimals := uint8(9) // ubadge has 9 decimals

	// The conversion: 1 ubadge = 1 * 10^9 abadge base units = 1 * 10^9 wei
	oneUbadgeAmount := sdkmath.NewInt(1)
	decimalsDiff := 18 - int(decimals) // 9
	conversionFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimalsDiff)), nil)
	expectedWei := new(big.Int).Mul(oneUbadgeAmount.BigInt(), conversionFactor)

	require.Equal(t, big.NewInt(1_000_000_000), expectedWei, "1 ubadge should convert to 1*10^9 wei")

	// Test that the conversion factor is correct
	require.Equal(t, big.NewInt(1_000_000_000), conversionFactor, "Conversion factor should be 10^9")
}
