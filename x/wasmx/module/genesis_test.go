package wasmx_test

import (
	"testing"

	keepertest "bitbadgeschain/testutil/keeper"
	"bitbadgeschain/testutil/nullify"
	wasmx "bitbadgeschain/x/wasmx/module"
	"bitbadgeschain/x/wasmx/types"

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
