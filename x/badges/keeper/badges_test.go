package keeper_test

import (
	"math"

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
			BadgeIds: GetOneIdRange(),
			OwnershipTimes: GetFullIdRanges(),
		},
	}
	collectionsToCreate[0].Collection.Transfers = []*types.Transfer{
		{
			From: "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetOneIdRange(),
					OwnershipTimes: GetFullIdRanges(),
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


	AssertBalancesEqual(suite, collection.TotalSupplys, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetOneIdRange(),
			OwnershipTimes: GetFullIdRanges(),
		},
	})

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertUintsEqual(suite, balance.Balances[0].Amount, sdkmath.NewUint(1))
	AssertIdRangesEqual(suite, balance.Balances[0].BadgeIds, []*types.IdRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(1),
		},
	})

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetTwoIdRanges(),
			OwnershipTimes: GetFullIdRanges(),
		},
	}, []*types.Transfer{
		{
			From: "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetTwoIdRanges(),
					OwnershipTimes: GetFullIdRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error creating badges: %s")

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	_, err = types.SubtractBalance(balance.Balances, &types.Balance{ 
		BadgeIds: []*types.IdRange{
			GetOneIdRange()[0],
			GetTwoIdRanges()[0],
		}, 
		OwnershipTimes: GetFullIdRanges(), 
		Amount: sdkmath.NewUint(1),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")


	_, err = types.SubtractBalance(collection.TotalSupplys, &types.Balance{ 
		BadgeIds: []*types.IdRange{
		GetOneIdRange()[0],
		GetTwoIdRanges()[0],
		}, 
		OwnershipTimes: GetFullIdRanges(), 
		Amount: sdkmath.NewUint(1),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")

	AssertBalancesEqual(suite, collection.UnmintedSupplys, []*types.Balance{})

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetTwoIdRanges(),
			OwnershipTimes: GetFullIdRanges(),
		},
	}, []*types.Transfer{})
	suite.Require().Nil(err, "Error creating badges: %s")

	_, err = types.SubtractBalance(collection.TotalSupplys, &types.Balance{ 
		BadgeIds: []*types.IdRange{
			GetTwoIdRanges()[0],
		}, 
		OwnershipTimes: GetFullIdRanges(), 
		Amount: sdkmath.NewUint(2),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")

	_, err = types.SubtractBalance(collection.UnmintedSupplys, &types.Balance{ 
		BadgeIds: []*types.IdRange{
			GetTwoIdRanges()[0],
		}, 
		OwnershipTimes: GetFullIdRanges(), 
		Amount: sdkmath.NewUint(2),
	})
	suite.Require().Error(err, "Error subtracting balances: %s")

	_, err = types.SubtractBalance(collection.UnmintedSupplys, &types.Balance{ 
		BadgeIds: []*types.IdRange{
		GetTwoIdRanges()[0],
		}, 
		OwnershipTimes: GetFullIdRanges(), 
		Amount: sdkmath.NewUint(1),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")

	
	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
			BadgeIds: GetTopHalfIdRanges(),
			OwnershipTimes: GetFullIdRanges(),
		},
	}, []*types.Transfer{})
	suite.Require().Nil(err, "Error creating badges: %s")
	AssertUintsEqual(suite, collection.NextBadgeId, sdkmath.NewUint(uint64(math.MaxUint64)).Add(sdkmath.NewUint(1)))

	_, err = types.SubtractBalance(collection.UnmintedSupplys, &types.Balance{ 
		BadgeIds: []*types.IdRange{
			GetTopHalfIdRanges()[0],
		}, 
		OwnershipTimes: GetFullIdRanges(), 
		Amount: types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
	})
	suite.Require().Nil(err, "Error subtracting balances: %s")


	_, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetTwoIdRanges(),
			OwnershipTimes: GetFullIdRanges(),
		},
	}, []*types.Transfer{
		{
			From: alice,
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetTwoIdRanges(),
					OwnershipTimes: GetFullIdRanges(),
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating badges: %s")
}
