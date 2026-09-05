package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// TestDefaultParamsCommunityPoolSwapDenom pins the default community-pool swap
// denom so any future change to it is a conscious decision. Stage 1 of the
// two-stage plan is native ubadge; stage 2 (via governance) is canonical USDC
// once a ubadge/USDC pool exists. See DefaultParams for the rationale.
func TestDefaultParamsCommunityPoolSwapDenom(t *testing.T) {
	params := types.DefaultParams()
	require.Equal(t, "ubadge",
		params.TakerFeeParams.CommunityPoolDenomToSwapNonWhitelistedAssetsTo)
}

// TestDefaultParamsValidate ensures the shipped defaults always pass their own
// validation.
func TestDefaultParamsValidate(t *testing.T) {
	require.NoError(t, types.DefaultParams().Validate())
}

// A taker fee of exactly 1 would zero every exact-in swap and divide by zero on exact-out.
func TestDefaultTakerFeeMustBeBelowOne(t *testing.T) {
	params := types.DefaultParams()
	params.TakerFeeParams.DefaultTakerFee = types.OneDec
	require.Error(t, params.Validate())

	params.TakerFeeParams.DefaultTakerFee = types.OneDec.Sub(types.OneDec.Quo(types.OneDec.MulInt64(100)))
	require.NoError(t, params.Validate())
}
