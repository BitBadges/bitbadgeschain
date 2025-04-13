package wasmx_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/testutil/nullify"
	"github.com/bitbadges/bitbadgeschain/x/wasmx/types"

	keepertest "github.com/bitbadges/bitbadgeschain/testutil/keeper"
	wasmx "github.com/bitbadges/bitbadgeschain/x/wasmx/module"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.WasmxKeeper(t)
	wasmx.InitGenesis(ctx, k, genesisState)
	got := wasmx.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
