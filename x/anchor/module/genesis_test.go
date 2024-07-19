package anchor_test

import (
	"testing"

	keepertest "bitbadgeschain/testutil/keeper"
	"bitbadgeschain/testutil/nullify"
	anchor "bitbadgeschain/x/anchor/module"
	"bitbadgeschain/x/anchor/types"

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
