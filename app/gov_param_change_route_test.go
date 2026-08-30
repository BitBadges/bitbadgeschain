package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

// TestLegacyGovRouterRoutesParamChangeProposals pins the gov legacy-content
// router wiring. Subspace-backed modules (poolmanager, gamm, ...) have no
// MsgUpdateParams, so a legacy ParameterChangeProposal (wrapped in
// MsgExecLegacyContent) is the ONLY governance path that can change their
// params. Without the params route, gov rejects such proposals at submission
// with "no handler exists for proposal type".
func TestLegacyGovRouterRoutesParamChangeProposals(t *testing.T) {
	app := Setup(false)
	router := app.GovKeeper.LegacyRouter()

	require.True(t, router.HasRoute(paramproposal.RouterKey),
		"legacy gov router must route %q (ParameterChangeProposal)", paramproposal.RouterKey)
	require.True(t, router.HasRoute(govtypes.RouterKey),
		"legacy gov router must keep the base gov route")
}
