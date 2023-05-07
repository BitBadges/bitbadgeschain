package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNewBadges() {
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
				AllowedTransfers: []*types.TransferMapping{
					{
						From: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
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
			Supply: 10,
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
			Amount:  10,
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 1, End: 1}}, bobbalance.Balances)
	suite.Require().Equal(uint64(10), fetchedBalance[0].Amount)

	//Create badge 2 with supply == 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 1,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	badge, _ = GetCollection(suite, wctx, 1)
	bobbalance, _ = GetUserBalance(suite, wctx, 1, bob)

	suite.Require().Equal(uint64(3), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 2, End: 2}}, //0 to 0 range so it will be nil
			Amount:  1,
		},
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Amount:  10,
		},
	}, badge.MaxSupplys)
	bobbalance, _ = GetUserBalance(suite, wctx, 1, bob)
	suite.Require().Equal(uint64(1), bobbalance.Balances[0].Amount)
	suite.Require().Equal(uint64(2), bobbalance.Balances[0].BadgeIds[0].Start)

	//Create badge 2 with supply == 10
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 2,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 1)
	bobbalance, _ = GetUserBalance(suite, wctx, 1, bob)

	suite.Require().Equal(uint64(5), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 2, End: 2}}, //0 to 0 range so it will be nil
			Amount:  1,
		},
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}, {Start: 3, End: 4}}, //0 to 0 range so it will be nil
			Amount:  10,
		},
	}, badge.MaxSupplys)
	suite.Require().Equal(uint64(10), bobbalance.Balances[1].Amount)
	suite.Require().Equal(uint64(1), bobbalance.Balances[1].BadgeIds[0].Start)
	suite.Require().Equal(uint64(1), bobbalance.Balances[1].BadgeIds[0].End)
	suite.Require().Equal(uint64(3), bobbalance.Balances[1].BadgeIds[1].Start)
	suite.Require().Equal(uint64(4), bobbalance.Balances[1].BadgeIds[1].End)
}

func (suite *TestSuite) TestNewBadgesDirectlyUponCreatingNewBadge() {
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
				AllowedTransfers: []*types.TransferMapping{
					{
						From: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 1)

	CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 1,
		},
	})

	badge, _ = GetCollection(suite, wctx, 1)

	bobbalance, _ := GetUserBalance(suite, wctx, 1, bob)

	suite.Require().Equal(uint64(2), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Amount:  10,
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 1, End: 1}}, bobbalance.Balances)
	suite.Require().Equal(uint64(10), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	//Create badge 2 with supply == 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 1,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	badge, _ = GetCollection(suite, wctx, 1)
	bobbalance, _ = GetUserBalance(suite, wctx, 1, bob)

	suite.Require().Equal(uint64(3), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 2, End: 2}}, //0 to 0 range so it will be nil
			Amount:  1,
		},
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Amount:  10,
		},
	}, badge.MaxSupplys)
	suite.Require().Equal(uint64(1), bobbalance.Balances[0].Amount)
	suite.Require().Equal(uint64(2), bobbalance.Balances[0].BadgeIds[0].Start)

	//Create badge 2 with supply == 10
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 2,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 1)
	bobbalance, _ = GetUserBalance(suite, wctx, 1, bob)

	suite.Require().Equal(uint64(5), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 2, End: 2}}, //0 to 0 range so it will be nil
			Amount:  1,
		},
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}, {Start: 3, End: 4}}, //0 to 0 range so it will be nil
			Amount:  10,
		},
	},
		badge.MaxSupplys)
	suite.Require().Equal(uint64(10), bobbalance.Balances[1].Amount)
	suite.Require().Equal(uint64(1), bobbalance.Balances[1].BadgeIds[0].Start)
	suite.Require().Equal(uint64(1), bobbalance.Balances[1].BadgeIds[0].End)
	suite.Require().Equal(uint64(3), bobbalance.Balances[1].BadgeIds[1].Start)
	suite.Require().Equal(uint64(4), bobbalance.Balances[1].BadgeIds[1].End)
}

func (suite *TestSuite) TestNewBadgesNotManager() {
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
				AllowedTransfers: []*types.TransferMapping{
					{
						From: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	err := CreateBadgesAndMintAllToCreator(suite, wctx, alice, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 1,
		},
	})
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}

func (suite *TestSuite) TestNewBadgeBadgeNotExists() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateBadgesAndMintAllToCreator(suite, wctx, alice, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 1,
		},
	})
	suite.Require().EqualError(err, keeper.ErrCollectionNotExists.Error())
}

func (suite *TestSuite) TestNewBadgeCreateIsLocked() {
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
				Permissions: 0,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 1,
		},
	})
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}
