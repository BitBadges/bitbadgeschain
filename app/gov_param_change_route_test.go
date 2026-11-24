package app

import (
	"testing"

	"github.com/stretchr/testify/require"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// TestLegacyGovRouterKeepsBaseGovRoute pins the gov legacy-content router
// wiring after x/params removal in Cosmos SDK v0.55.
func TestLegacyGovRouterKeepsBaseGovRoute(t *testing.T) {
	app := Setup(false)
	router := app.GovKeeper.LegacyRouter()

	require.True(t, router.HasRoute(govtypes.RouterKey),
		"legacy gov router must keep the base gov route")
}
