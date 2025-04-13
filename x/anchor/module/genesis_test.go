package anchor_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/nullify"
	"github.com/bitbadges/bitbadgeschain/x/anchor/types"

	keepertest "github.com/bitbadges/bitbadgeschain/testutil/keeper"
	anchor "github.com/bitbadges/bitbadgeschain/x/anchor/module"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.AnchorKeeper(t)
	anchor.InitGenesis(ctx, k, genesisState)
	got := anchor.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
