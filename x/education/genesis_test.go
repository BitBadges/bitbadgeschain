package education_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/trevormil/bitbadgeschain/testutil/keeper"
	"github.com/trevormil/bitbadgeschain/testutil/nullify"
	"github.com/trevormil/bitbadgeschain/x/education"
	"github.com/trevormil/bitbadgeschain/x/education/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.EducationKeeper(t)
	education.InitGenesis(ctx, *k, genesisState)
	got := education.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	// this line is used by starport scaffolding # genesis/test/assert
}
