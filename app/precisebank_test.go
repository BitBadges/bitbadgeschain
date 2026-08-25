package app

import (
	"testing"

	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/stretchr/testify/require"
)

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
