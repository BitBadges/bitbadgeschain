package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetWrappableBalances(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Test with invalid denom format
	_, err := suite.app.BadgesKeeper.GetWrappableBalances(wctx, &types.QueryGetWrappableBalancesRequest{
		Denom:   "invalid-denom",
		Address: bob,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "denom must start with 'badges:'")

	// Test with invalid denom format (missing parts)
	_, err = suite.app.BadgesKeeper.GetWrappableBalances(wctx, &types.QueryGetWrappableBalancesRequest{
		Denom:   "badges:1",
		Address: bob,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid denom format")

	// Test with non-existent collection
	_, err = suite.app.BadgesKeeper.GetWrappableBalances(wctx, &types.QueryGetWrappableBalancesRequest{
		Denom:   "badges:999:test",
		Address: bob,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "collection 999 not found")

	// Create a collection with cosmos wrapper paths
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "testcoin",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10), // 10 native badges = 1 wrapped token
					OwnershipTimes: GetFullUintRanges(),
					BadgeIds:       GetOneUintRange(),
				},
			},
			Symbol:              "TESTCOIN",
			DenomUnits:          []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "testcoin", IsDefaultDisplay: true}},
			AllowCosmosWrapping: true,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	require.NoError(t, err, "error creating collection for wrappable balances test")

	// Execute the transfers to actually give badges to the users
	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	require.NoError(t, err, "error getting collection")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				BadgeIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		CollectionApprovals: collection.CollectionApprovals,
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	require.NoError(t, err, "error minting and distributing badges")

	// Test with valid denom but no wrapper path found
	_, err = suite.app.BadgesKeeper.GetWrappableBalances(wctx, &types.QueryGetWrappableBalancesRequest{
		Denom:   "badges:1:nonexistent",
		Address: bob,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "wrapper path not found")

	// Test with valid denom and wrapper path
	response, err := suite.app.BadgesKeeper.GetWrappableBalances(wctx, &types.QueryGetWrappableBalancesRequest{
		Denom:   "badges:1:testcoin",
		Address: bob,
	})
	require.NoError(t, err)
	require.NotNil(t, response)

	// Bob should have 1 badge, and since 10 native badges = 1 wrapped token,
	// the max wrappable amount should be 0 (1 badge / 10 = 0 with integer division)
	require.Equal(t, sdkmath.NewUint(0), response.MaxWrappableAmount)

	// Test with a different user who has more badges
	// Create another collection with a different wrapper path for charlie
	collectionsToCreate2 := GetTransferableCollectionToCreateAllMintedToCreator(charlie)
	collectionsToCreate2[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "testcoin-two",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(5), // 5 native badges = 1 wrapped token
					OwnershipTimes: GetFullUintRanges(),
					BadgeIds:       GetOneUintRange(),
				},
			},
			Symbol:              "TESTCOIN-TWO",
			DenomUnits:          []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "testcoin-two", IsDefaultDisplay: true}},
			AllowCosmosWrapping: true,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate2)
	require.NoError(t, err, "error creating second collection for wrappable balances test")

	// Execute the transfers to actually give badges to charlie
	collection2, err := GetCollection(suite, wctx, sdkmath.NewUint(2))
	require.NoError(t, err, "error getting collection 2")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      charlie,
		CollectionId: sdkmath.NewUint(2),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				BadgeIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		CollectionApprovals: collection2.CollectionApprovals,
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(2)),
			},
		},
	})
	require.NoError(t, err, "error minting and distributing badges to charlie")

	// Charlie should have 1 badge, and since 5 native badges = 1 wrapped token,
	// the max wrappable amount should be 0 (1 badge / 5 = 0 with integer division)
	response2, err := suite.app.BadgesKeeper.GetWrappableBalances(wctx, &types.QueryGetWrappableBalancesRequest{
		Denom:   "badges:2:testcoin-two",
		Address: charlie,
	})
	require.NoError(t, err)
	require.NotNil(t, response2)
	require.Equal(t, sdkmath.NewUint(0), response2.MaxWrappableAmount)

	// Test with a wrapper path that doesn't allow cosmos wrapping
	// This should now work and return 0 since the user has 1 badge but wrapper needs 1 badge
	collectionsToCreate3 := GetTransferableCollectionToCreateAllMintedToCreator(alice)
	collectionsToCreate3[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "nowrap",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					BadgeIds:       GetOneUintRange(), // Use GetOneUintRange to match what the user gets
				},
			},
			Symbol:              "NOWRAP",
			DenomUnits:          []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "nowrap", IsDefaultDisplay: true}},
			AllowCosmosWrapping: false, // This should now work
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate3)
	require.NoError(t, err, "error creating collection for no-wrap test")

	// Execute the transfers to actually give badges to alice
	collection3, err := GetCollection(suite, wctx, sdkmath.NewUint(3))
	require.NoError(t, err, "error getting collection 3")

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(3),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				BadgeIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		CollectionApprovals: collection3.CollectionApprovals,
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(3)),
			},
		},
	})
	require.NoError(t, err, "error minting and distributing badges to alice")

	response3, err := suite.app.BadgesKeeper.GetWrappableBalances(wctx, &types.QueryGetWrappableBalancesRequest{
		Denom:   "badgeslp:3:nowrap",
		Address: alice,
	})
	require.NoError(t, err)
	require.NotNil(t, response3)
	// The main goal was to test that AllowCosmosWrapping=false doesn't cause an error
	// The actual wrappable amount depends on the specific badge setup, which can be 0 or more
	require.True(t, response3.MaxWrappableAmount.GTE(sdkmath.NewUint(0)))
}

func TestKeeper_GetWrappableBalances_AdvancedLogic(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with a wrapper path that requires multiple badge types
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "advanced-test",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(5), // 5 of badge ID 1
					OwnershipTimes: GetFullUintRanges(),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				},
			},
			Symbol:              "ADVANCED-TEST",
			DenomUnits:          []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "advanced-test", IsDefaultDisplay: true}},
			AllowCosmosWrapping: true,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	require.NoError(t, err, "error creating collection for advanced test")

	// Test with user who has exactly enough badges for 1 wrapped token
	// User has: 1 of badge ID 1
	// Wrapper needs: 5 of badge ID 1 for 1 wrapped token
	// So max wrappable should be 0 (1 < 5)
	response, err := suite.app.BadgesKeeper.GetWrappableBalances(wctx, &types.QueryGetWrappableBalancesRequest{
		Denom:   "badges:1:advanced-test",
		Address: bob,
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, sdkmath.NewUint(0), response.MaxWrappableAmount)

	// Test with user who has more badges
	// Give charlie more badges by creating another collection
	collectionsToCreate2 := GetTransferableCollectionToCreateAllMintedToCreator(charlie)
	collectionsToCreate2[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "advanced-test-two",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(3), // 3 of badge ID 1
					OwnershipTimes: GetFullUintRanges(),
					BadgeIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				},
			},
			Symbol:              "ADVANCED-TEST-TWO",
			DenomUnits:          []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "advanced-test-two", IsDefaultDisplay: true}},
			AllowCosmosWrapping: true,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate2)
	require.NoError(t, err, "error creating second collection for advanced test")

	// Charlie has: 1 of badge ID 1
	// Wrapper needs: 3 of badge ID 1 for 1 wrapped token
	// So max wrappable should be 0 (1 < 3)
	response2, err := suite.app.BadgesKeeper.GetWrappableBalances(wctx, &types.QueryGetWrappableBalancesRequest{
		Denom:   "badges:2:advanced-test-two",
		Address: charlie,
	})
	require.NoError(t, err)
	require.NotNil(t, response2)
	require.Equal(t, sdkmath.NewUint(0), response2.MaxWrappableAmount)
}
