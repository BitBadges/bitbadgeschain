package badges_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/trevormil/bitbadgeschain/testutil/keeper"
	"github.com/trevormil/bitbadgeschain/testutil/nullify"
	"github.com/trevormil/bitbadgeschain/x/badges"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		NextAssetId: 0,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BadgesKeeper(t)
	badges.InitGenesis(ctx, *k, genesisState)
	got := badges.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	require.Equal(t, genesisState.NextAssetId, got.NextAssetId)
require.Equal(t, genesisState.NextAssetId, got.NextAssetId)
// this line is used by starport scaffolding # genesis/test/assert
}
