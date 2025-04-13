package maps_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/nullify"
	"github.com/bitbadges/bitbadgeschain/x/maps/types"

	keepertest "github.com/bitbadges/bitbadgeschain/testutil/keeper"
	maps "github.com/bitbadges/bitbadgeschain/x/maps/module"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MapsKeeper(t)
	maps.InitGenesis(ctx, k, genesisState)
	got := maps.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
