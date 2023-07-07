package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestCreateBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.BadgesToCreate = []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	collectionsToCreate[0].Collection.Transfers = []*types.Transfer{
		{
			From: "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
	}
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].FromMappingId = "Mint"
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesToApprovedIncomingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	
	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting badge: %s")
	balance := &types.UserBalanceStore{}


	totalSupplys, err := GetUserBalance(suite, wctx, sdk.NewUint(1), "Total")
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, totalSupplys.Balances, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	})

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertUintsEqual(suite, balance.Balances[0].Amount, sdkmath.NewUint(1))
	AssertUintRangesEqual(suite, balance.Balances[0].BadgeIds, []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(1),
		},
	})

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetTwoUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}, []*types.Transfer{
		{
			From: "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetTwoUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
	}, bob)
	suite.Require().Nil(err, "Error creating badges: %s")

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	_, err = types.SubtractBalance(balance.Balances, &types.Balance{ 
		BadgeIds: []*types.UintRange{
			GetOneUintRange()[0],
			GetTwoUintRanges()[0],
		}, 
		OwnershipTimes: GetFullUintRanges(), 
		Amount: sdkmath.NewUint(1),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")


	totalSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Total")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(totalSupplys.Balances, &types.Balance{ 
		BadgeIds: []*types.UintRange{
		GetOneUintRange()[0],
		GetTwoUintRanges()[0],
		}, 
		OwnershipTimes: GetFullUintRanges(), 
		Amount: sdkmath.NewUint(1),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")

	unmintedSupplys, err := GetUserBalance(suite, wctx, sdk.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, unmintedSupplys.Balances, []*types.Balance{})

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetTwoUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}, []*types.Transfer{}, bob)
	suite.Require().Nil(err, "Error creating badges: %s")

	totalSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Total")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(totalSupplys.Balances, &types.Balance{ 
		BadgeIds: []*types.UintRange{
			GetTwoUintRanges()[0],
		}, 
		OwnershipTimes: GetFullUintRanges(), 
		Amount: sdkmath.NewUint(2),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")

	unmintedSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(unmintedSupplys.Balances, &types.Balance{ 
		BadgeIds: []*types.UintRange{
			GetTwoUintRanges()[0],
		}, 
		OwnershipTimes: GetFullUintRanges(), 
		Amount: sdkmath.NewUint(2),
	})
	suite.Require().Error(err, "Error subtracting balances: %s")

	unmintedSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(unmintedSupplys.Balances, &types.Balance{ 
		BadgeIds: []*types.UintRange{
			GetTwoUintRanges()[0],
		}, 
		OwnershipTimes: GetFullUintRanges(), 
		Amount: sdkmath.NewUint(1),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")

	_, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
			BadgeIds: GetTopHalfUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}, []*types.Transfer{}, bob)
	suite.Require().Error(err, "Error creating badges: %s")
	
	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
			BadgeIds: GetBottomHalfUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
		{
			Amount: types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
			BadgeIds: GetTopHalfUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}, []*types.Transfer{}, bob)
	suite.Require().Nil(err, "Error creating badges: %s")
	// AssertUintsEqual(suite, collection.NextBadgeId, sdkmath.NewUint(uint64(math.MaxUint64)).Add(sdkmath.NewUint(1)))

	unmintedSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(unmintedSupplys.Balances, &types.Balance{ 
		BadgeIds: []*types.UintRange{
			GetTopHalfUintRanges()[0],
		}, 
		OwnershipTimes: GetFullUintRanges(), 
		Amount: types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")


	_, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetTwoUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}, []*types.Transfer{
		{
			From: alice,
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetTwoUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
	}, bob)
	suite.Require().Error(err, "Error creating badges: %s")
}
