package badges_test

import (
	"testing"

	keepertest "github.com/bitbadges/bitbadgeschain/testutil/keeper"
	"github.com/bitbadges/bitbadgeschain/testutil/nullify"
	"github.com/bitbadges/bitbadgeschain/x/badges"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params:      types.DefaultParams(),
		PortId:      types.PortID,
		NextCollectionId: 0,
		Collections:      []*types.BadgeCollection{},
		Balances:    []*types.UserBalance{},
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

	require.Equal(t, genesisState.NextCollectionId, got.NextCollectionId)
	require.Equal(t, genesisState.Collections, got.Collections)
	require.Equal(t, genesisState.Balances, got.Balances)
	require.Equal(t, genesisState.BalanceIds, got.BalanceIds)
	// this line is used by starport scaffolding # genesis/test/assert
}
