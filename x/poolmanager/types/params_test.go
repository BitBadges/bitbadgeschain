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
