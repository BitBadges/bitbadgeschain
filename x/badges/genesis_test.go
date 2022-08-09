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
		Params:      types.DefaultParams(),
		PortId:      types.PortID,
		NextBadgeId: 0,
		Badges:      []*types.BitBadge{},
		Balances:    []*types.BadgeBalanceInfo{},
		BalanceIds:  []string{},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BadgesKeeper(t)
	badges.InitGenesis(ctx, *k, genesisState)
	got := badges.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	require.Equal(t, genesisState.NextBadgeId, got.NextBadgeId)
	require.Equal(t, genesisState.Badges, got.Badges)
	require.Equal(t, genesisState.Balances, got.Balances)
	require.Equal(t, genesisState.BalanceIds, got.BalanceIds)
	// this line is used by starport scaffolding # genesis/test/assert
}
