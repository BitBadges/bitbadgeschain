package collections_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/trevormil/bitbadgeschain/testutil/keeper"
	"github.com/trevormil/bitbadgeschain/testutil/nullify"
	"github.com/trevormil/bitbadgeschain/x/collections"
	"github.com/trevormil/bitbadgeschain/x/collections/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CollectionsKeeper(t)
	collections.InitGenesis(ctx, *k, genesisState)
	got := collections.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
