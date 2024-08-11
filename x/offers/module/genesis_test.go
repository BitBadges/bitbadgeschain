package offers_test

// func TestGenesis(t *testing.T) {
// 	genesisState := types.GenesisState{
// 		Params: types.DefaultParams(),
// 		PortId: types.PortID,
// 		// this line is used by starport scaffolding # genesis/test/state
// 	}

// 	k, ctx := keepertest.OffersKeeper(t)
// 	offers.InitGenesis(ctx, k, genesisState)
// 	got := offers.ExportGenesis(ctx, k)
// 	require.NotNil(t, got)

// 	nullify.Fill(&genesisState)
// 	nullify.Fill(got)

// 	require.Equal(t, genesisState.PortId, got.PortId)

// 	// this line is used by starport scaffolding # genesis/test/assert
// }
