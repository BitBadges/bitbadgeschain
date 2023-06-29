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
			Times: GetFullIdRanges(),
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
					Times: GetFullIdRanges(),
				},
			},
		},
	}
	collectionsToCreate[0].Collection.ApprovedTransfersTimeline[0].ApprovedTransfers[0].FromMappingId = "Mint"
	collectionsToCreate[0].Collection.ApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesFromApprovedOutgoingTransfers = true
	collectionsToCreate[0].Collection.ApprovedTransfersTimeline[0].ApprovedTransfers[0].OverridesToApprovedIncomingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	
	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting badge: %s")
	balance := &types.UserBalanceStore{}


	suite.Require().Equal(collection.TotalSupplys, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetOneIdRange(),
			Times: GetFullIdRanges(),
		},
	})

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	suite.Require().Equal(balance.Balances[0].Amount, sdkmath.NewUint(1))
	suite.Require().Equal(balance.Balances[0].BadgeIds, []*types.IdRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(1),
		},
	})

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetTwoIdRanges(),
			Times: GetFullIdRanges(),
		},
	}, []*types.Transfer{
		{
			From: "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetTwoIdRanges(),
					Times: GetFullIdRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error creating badges: %s")

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	_, err = types.SubtractBalancesForIdRanges(balance.Balances, []*types.IdRange{
		GetOneIdRange()[0],
		GetTwoIdRanges()[0],
	}, GetFullIdRanges(), sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error subtracting balances: %s")


	_, err = types.SubtractBalancesForIdRanges(collection.TotalSupplys, []*types.IdRange{
		GetOneIdRange()[0],
		GetTwoIdRanges()[0],
	}, GetFullIdRanges(), sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error subtracting balances: %s")

	suite.Require().Equal(collection.UnmintedSupplys, []*types.Balance{})

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetTwoIdRanges(),
			Times: GetFullIdRanges(),
		},
	}, []*types.Transfer{})
	suite.Require().Nil(err, "Error creating badges: %s")

	_, err = types.SubtractBalancesForIdRanges(collection.TotalSupplys, []*types.IdRange{
		GetTwoIdRanges()[0],
	}, GetFullIdRanges(), sdkmath.NewUint(2))
	suite.Require().Nil(err, "Error subtracting balances: %s")

	_, err = types.SubtractBalancesForIdRanges(collection.UnmintedSupplys, []*types.IdRange{
		GetTwoIdRanges()[0],
	}, GetFullIdRanges(), sdkmath.NewUint(2))
	suite.Require().Error(err, "Error subtracting balances: %s")

	_, err = types.SubtractBalancesForIdRanges(collection.UnmintedSupplys, []*types.IdRange{
		GetTwoIdRanges()[0],
	}, GetFullIdRanges(), sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error subtracting balances: %s")

	
	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
			BadgeIds: GetTopHalfIdRanges(),
			Times: GetFullIdRanges(),
		},
	}, []*types.Transfer{})
	suite.Require().Nil(err, "Error creating badges: %s")
	suite.Require().True(collection.NextBadgeId.Equal(sdkmath.NewUint(uint64(math.MaxUint64)).Add(sdkmath.NewUint(1))))
	
	_, err = types.SubtractBalancesForIdRanges(collection.UnmintedSupplys, []*types.IdRange{
		GetTopHalfIdRanges()[0],
	}, GetFullIdRanges(), types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"))
	suite.Require().Nil(err, "Error subtracting balances: %s")


	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: GetTwoIdRanges(),
			Times: GetFullIdRanges(),
		},
	}, []*types.Transfer{
		{
			From: alice,
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: GetTwoIdRanges(),
					Times: GetFullIdRanges(),
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating badges: %s")
}
