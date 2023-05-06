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
	bobbalance, _ := GetUserBalance(suite, wctx, 1, bob)

	suite.Require().Equal(uint64(2), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Amount:  10000,
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 1, End: 1}}, bobbalance.Balances)
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, alice, 1, []*types.Balance{
		{
			Amount:  1000,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	bobbalance, _ = GetUserBalance(suite, wctx, 1, bob)
	suite.Require().Equal(alice, bobbalance.Approvals[0].Address)
	suite.Require().Equal(uint64(1000), bobbalance.Approvals[0].Balances[0].Amount)

	err = SetApproval(suite, wctx, bob, charlie, 1, []*types.Balance{
		{
			Amount:  500,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	err = SetApproval(suite, wctx, bob, alice, 1, []*types.Balance{
		{
			Amount:  500,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	bobbalance, _ = GetUserBalance(suite, wctx, 1, bob)

	suite.Require().Equal(charlie, bobbalance.Approvals[1].Address)
	suite.Require().Equal(uint64(500), bobbalance.Approvals[1].Balances[0].Amount)

	suite.Require().Equal(alice, bobbalance.Approvals[0].Address)
	suite.Require().Equal(uint64(500), bobbalance.Approvals[0].Balances[0].Amount)
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
	bobbalance, _ := GetUserBalance(suite, wctx, 1, bob)

	suite.Require().Equal(uint64(2), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Amount:  10000,
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 1, End: 1}}, bobbalance.Balances)
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, charlie, alice, 1, []*types.Balance{
		{
			Amount:  0,
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
	bobbalance, _ := GetUserBalance(suite, wctx, 1, bob)

	suite.Require().Equal(uint64(2), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Amount:  10000,
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 1, End: 1}}, bobbalance.Balances)
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, bob, 1, []*types.Balance{
		{
			Amount:  1000,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
		},
	})
	suite.Require().EqualError(err, keeper.ErrAccountCanNotEqualCreator.Error())
}
