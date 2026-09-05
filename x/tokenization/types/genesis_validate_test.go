package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func validGenesisForTest() *types.GenesisState {
	gs := types.DefaultGenesis()
	gs.Collections = []*types.TokenCollection{{CollectionId: sdkmath.NewUint(1)}, {CollectionId: sdkmath.NewUint(2)}}
	gs.NextCollectionId = sdkmath.NewUint(3)
	gs.DynamicStores = []*types.DynamicStore{{StoreId: sdkmath.NewUint(1)}}
	gs.NextDynamicStoreId = sdkmath.NewUint(2)
	gs.AddressLists = []*types.AddressList{{ListId: "listA"}}
	gs.Balances = []*types.UserBalanceStore{{}}
	gs.BalanceStoreKeys = []string{"1-addr"}
	return gs
}

func TestGenesisStateValidateStructure(t *testing.T) {
	require.NoError(t, validGenesisForTest().Validate())

	cases := map[string]func(gs *types.GenesisState){
		"balances and keys length mismatch":            func(gs *types.GenesisState) { gs.BalanceStoreKeys = nil },
		"challenge trackers length mismatch":           func(gs *types.GenesisState) { gs.ChallengeTrackers = []sdkmath.Uint{sdkmath.NewUint(1)} },
		"approval trackers length mismatch":            func(gs *types.GenesisState) { gs.ApprovalTrackers = []*types.ApprovalTracker{{}} },
		"approval tracker versions length mismatch":    func(gs *types.GenesisState) { gs.ApprovalTrackerVersionsStoreKeys = []string{"k"} },
		"eth signature trackers length mismatch":       func(gs *types.GenesisState) { gs.EthSignatureTrackerStoreKeys = []string{"k"} },
		"voting trackers length mismatch":              func(gs *types.GenesisState) { gs.VotingTrackers = []*types.VoteProof{{}} },
		"collection stats length mismatch":             func(gs *types.GenesisState) { gs.CollectionStatsIds = []sdkmath.Uint{sdkmath.NewUint(1)} },
		"voting challenge trackers length mismatch":    func(gs *types.GenesisState) { gs.VotingChallengeTrackers = []*types.VotingChallengeTracker{{}} },
		"duplicate collection id":                      func(gs *types.GenesisState) { gs.Collections[1].CollectionId = sdkmath.NewUint(1) },
		"nil collection":                               func(gs *types.GenesisState) { gs.Collections = append(gs.Collections, nil) },
		"next collection id not above max":             func(gs *types.GenesisState) { gs.NextCollectionId = sdkmath.NewUint(2) },
		"duplicate dynamic store id":                   func(gs *types.GenesisState) { gs.DynamicStores = append(gs.DynamicStores, &types.DynamicStore{StoreId: sdkmath.NewUint(1)}) },
		"next dynamic store id not above max":          func(gs *types.GenesisState) { gs.NextDynamicStoreId = sdkmath.NewUint(1) },
		"duplicate address list id":                    func(gs *types.GenesisState) { gs.AddressLists = append(gs.AddressLists, &types.AddressList{ListId: "listA"}) },
		"duplicate balance store key":                  func(gs *types.GenesisState) { gs.Balances = append(gs.Balances, &types.UserBalanceStore{}); gs.BalanceStoreKeys = append(gs.BalanceStoreKeys, "1-addr") },
		"nil next collection id":                       func(gs *types.GenesisState) { gs.NextCollectionId = sdkmath.Uint{} },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			gs := validGenesisForTest()
			mutate(gs)
			require.NotPanics(t, func() { require.Error(t, gs.Validate()) })
		})
	}
}
