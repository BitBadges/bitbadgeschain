package protocols_test

import (
	"testing"

	keepertest "github.com/bitbadges/bitbadgeschain/testutil/keeper"
	"github.com/bitbadges/bitbadgeschain/testutil/nullify"
	"github.com/bitbadges/bitbadgeschain/x/protocols"
	"github.com/bitbadges/bitbadgeschain/x/protocols/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ProtocolsKeeper(t)
	protocols.InitGenesis(ctx, *k, genesisState)
	got := protocols.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
