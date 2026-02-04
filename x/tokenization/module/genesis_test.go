package tokenization_test

import (
	"testing"
)

func TestGenesis(t *testing.T) {
	// genesisState := types.GenesisState{
	// 	Params:           types.DefaultParams(),
	// 	PortId:           types.PortID,
	// 	NextCollectionId: sdkmath.NewUint(1),
	// 	Collections:      []*types.TokenCollection{},
	// 	Balances:         []*types.UserBalanceStore{},
	// 	BalanceStoreKeys: []string{},
	// 	// Claims:           []*types.Claim{},
	// 	// ClaimStoreKeys:   []string{},
	// 	// NextClaimId:      sdkmath.NewUint(1),
	// 	// this line is used by starport scaffolding # genesis/test/state
	// }

	// k, ctx := keepertest.TokenizationKeeper(t)
	// badges.InitGenesis(ctx, *k, genesisState)
	// got := badges.ExportGenesis(ctx, *k)
	// require.NotNil(t, got)

	// nullify.Fill(&genesisState)
	// nullify.Fill(got)

	// require.Equal(t, genesisState.PortId, got.PortId)

	// require.Equal(t, genesisState.NextCollectionId, got.NextCollectionId)
	// require.Equal(t, genesisState.Collections, got.Collections)
	// require.Equal(t, genesisState.Balances, got.Balances)
	// require.Equal(t, genesisState.BalanceStoreKeys, got.BalanceStoreKeys)
	// require.Equal(t, genesisState.NextClaimId, got.NextClaimId)
	// require.Equal(t, genesisState.Claims, got.Claims)
	// require.Equal(t, genesisState.ClaimStoreKeys, got.ClaimStoreKeys)

	// this line is used by starport scaffolding # genesis/test/assert
}
