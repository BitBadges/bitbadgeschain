package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestGetCollection() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com/",
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

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	badge, err := suite.app.BadgesKeeper.GetCollectionE(suite.ctx, 0)
	suite.Require().Nil(err, "Error getting badge: %s")
	suite.Require().Equal(badge.CollectionId, uint64(0))

	badge, err = suite.app.BadgesKeeper.GetCollectionE(suite.ctx, 1)
	suite.Require().EqualError(err, keeper.ErrCollectionNotExists.Error())
}

func (suite *TestSuite) TestGetBadgeAndAssertBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
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
				CollectionUri: "https://example.com",
				Permissions:   62,
				BadgeSupplys: []*types.BadgeSupplyAndAmount{
					{
						Supply: 1,
						Amount: 1,
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	_, err = suite.app.BadgesKeeper.GetCollectionAndAssertBadgeIdsAreValid(suite.ctx, 0, []*types.IdRange{
		{
			Start: 0,
			End:   0,
		},
	})
	suite.Require().Nil(err, "Error getting badge: %s")

	_, err = suite.app.BadgesKeeper.GetCollectionAndAssertBadgeIdsAreValid(suite.ctx, 0, []*types.IdRange{
		{
			Start: 20,
			End:   10,
		},
	})
	suite.Require().EqualError(err, keeper.ErrInvalidBadgeRange.Error())

	_, err = suite.app.BadgesKeeper.GetCollectionAndAssertBadgeIdsAreValid(suite.ctx, 0, []*types.IdRange{
		{
			Start: 0,
			End:   10,
		},
	})
	suite.Require().EqualError(err, keeper.ErrBadgeNotExists.Error())
}

func (suite *TestSuite) TestCreateBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
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
				CollectionUri: "https://example.com",
				Permissions:   62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	badge, err := GetCollection(suite, wctx, 0)
	suite.Require().Nil(err, "Error getting badge: %s")
	balance := types.UserBalance{}

	badge, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, badge, []*types.BadgeSupplyAndAmount{
		{
			Supply: 1,
			Amount: 1,
		},
	}, []*types.Transfers{
		{
			ToAddresses: []uint64{bobAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 1,
					BadgeIds: []*types.IdRange{
						{
							Start: 0,
							End:   0,
						},
					},
				},
			},
		},
	}, []*types.Claim{}, bob)
	suite.Require().Nil(err, "Error creating subassets: %s")

	suite.Require().Equal(badge.MaxSupplys, []*types.Balance{
		{
			Balance: 1,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   0,
				},
			},
		},
	})

	balance, err = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Nil(err, "Error getting user balance: %s")
	suite.Require().Equal(balance.Balances[0].Balance, uint64(1))
	suite.Require().Equal(balance.Balances[0].BadgeIds, []*types.IdRange{
		{
			Start: 0,
			End:   0,
		},
	})

	badge, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, badge, []*types.BadgeSupplyAndAmount{
		{
			Supply: 1,
			Amount: 1,
		},
	}, []*types.Transfers{
		{
			ToAddresses: []uint64{bobAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 1,
					BadgeIds: []*types.IdRange{
						{
							Start: 1,
							End:   1,
						},
					},
				},
			},
		},
	}, []*types.Claim{}, bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	balance, err = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Nil(err, "Error getting user balance: %s")
	suite.Require().Nil(err, "Error creating subassets: %s")
	suite.Require().Equal(badge.MaxSupplys, []*types.Balance{
		{
			Balance: 1,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   1,
				},
			},
		},
	})
	suite.Require().Equal(balance.Balances[0].Balance, uint64(1))
	suite.Require().Equal(balance.Balances[0].BadgeIds, []*types.IdRange{
		{
			Start: 0,
			End:   1,
		},
	})

	badge, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, badge, []*types.BadgeSupplyAndAmount{
		{
			Supply: 1,
			Amount: 1,
		},
	}, []*types.Transfers{
		{
			ToAddresses: []uint64{bobAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 1,
					BadgeIds: []*types.IdRange{
						{
							Start: 2,
							End:   2,
						},
					},
				},
			},
		},
	}, []*types.Claim{}, bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	balance, err = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Nil(err, "Error getting user balance: %s")
	suite.Require().Nil(err, "Error creating subassets: %s")
	suite.Require().Equal(badge.MaxSupplys, []*types.Balance{
		{
			Balance: 1,
			BadgeIds: []*types.IdRange{
				{
					Start: 0,
					End:   2,
				},
			},
		},
	})
	suite.Require().Equal(balance.Balances[0].Balance, uint64(1))
	suite.Require().Equal(balance.Balances[0].BadgeIds, []*types.IdRange{
		{
			Start: 0,
			End:   2,
		},
	})

	badge, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, badge, []*types.BadgeSupplyAndAmount{
		{
			Supply: 1,
			Amount: math.MaxUint64,
		},
	}, []*types.Transfers{
		{
			ToAddresses: []uint64{bobAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 1,
					BadgeIds: []*types.IdRange{
						{
							Start: 3,
							End:   3,
						},
					},
				},
			},
		},
	}, []*types.Claim{}, bob)
	suite.Require().EqualError(err, keeper.ErrOverflow.Error())
}
