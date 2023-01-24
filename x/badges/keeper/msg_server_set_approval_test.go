package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestSetApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 0)

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 0)
	bobbalance, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.MaxSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobbalance.Balances)[0].Balance)

	err = SetApproval(suite, wctx, bob, aliceAccountNum, 0, []*types.Balance{
		{
			Balance:  1000,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	bobbalance, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(aliceAccountNum), bobbalance.Approvals[0].Address)
	suite.Require().Equal(uint64(1000), bobbalance.Approvals[0].Balances[0].Balance)

	err = SetApproval(suite, wctx, bob, charlieAccountNum, 0, []*types.Balance{
		{
			Balance:  500,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	err = SetApproval(suite, wctx, bob, aliceAccountNum, 0, []*types.Balance{
		{
			Balance:  500,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	bobbalance, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(charlieAccountNum), bobbalance.Approvals[1].Address)
	suite.Require().Equal(uint64(500), bobbalance.Approvals[1].Balances[0].Balance)

	suite.Require().Equal(uint64(aliceAccountNum), bobbalance.Approvals[0].Address)
	suite.Require().Equal(uint64(500), bobbalance.Approvals[0].Balances[0].Balance)
}

func (suite *TestSuite) TestSetApprovalNoPrevBalanceInStore() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 0)

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 0)
	bobbalance, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.MaxSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobbalance.Balances)[0].Balance)

	err = SetApproval(suite, wctx, charlie, aliceAccountNum, 0, []*types.Balance{
		{
			Balance:  0,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")
}

func (suite *TestSuite) TestApproveSelf() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 0)

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 0)
	bobbalance, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.MaxSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobbalance.Balances)[0].Balance)

	err = SetApproval(suite, wctx, bob, bobAccountNum, 0, []*types.Balance{
		{
			Balance:  1000,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().EqualError(err, keeper.ErrAccountCanNotEqualCreator.Error())
}
