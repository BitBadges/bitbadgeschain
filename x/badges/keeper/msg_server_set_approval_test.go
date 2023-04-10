package keeper_test

import (
	"math"

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
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
				Permissions: 62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 1)

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 1)
	bobbalance, _ := GetUserBalance(suite, wctx, 1, bobAccountNum)

	suite.Require().Equal(uint64(2), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 1, End: 1}}, bobbalance.Balances)
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, aliceAccountNum, 1, []*types.Balance{
		{
			Balance:  1000,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	bobbalance, _ = GetUserBalance(suite, wctx, 1, bobAccountNum)
	suite.Require().Equal(uint64(aliceAccountNum), bobbalance.Approvals[0].Address)
	suite.Require().Equal(uint64(1000), bobbalance.Approvals[0].Balances[0].Balance)

	err = SetApproval(suite, wctx, bob, charlieAccountNum, 1, []*types.Balance{
		{
			Balance:  500,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	err = SetApproval(suite, wctx, bob, aliceAccountNum, 1, []*types.Balance{
		{
			Balance:  500,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	bobbalance, _ = GetUserBalance(suite, wctx, 1, bobAccountNum)

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
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
				Permissions: 62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 1)

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 1)
	bobbalance, _ := GetUserBalance(suite, wctx, 1, bobAccountNum)

	suite.Require().Equal(uint64(2), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 1, End: 1}}, bobbalance.Balances)
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, charlie, aliceAccountNum, 1, []*types.Balance{
		{
			Balance:  0,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
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
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
				Permissions: 62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 1)

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 1)
	bobbalance, _ := GetUserBalance(suite, wctx, 1, bobAccountNum)

	suite.Require().Equal(uint64(2), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 1, End: 1}}, bobbalance.Balances)
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, bobAccountNum, 0, []*types.Balance{
		{
			Balance:  1000,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
		},
	})
	suite.Require().EqualError(err, keeper.ErrAccountCanNotEqualCreator.Error())
}
