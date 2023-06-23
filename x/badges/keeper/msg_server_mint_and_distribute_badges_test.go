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
				CollectionMetadata: "https://example.com",
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(62),
				ApprovedTransfers: []*types.CollectionApprovedTransfer{
					{
						From: &types.AddressMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10),
			Amount: sdk.NewUint(1),
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))
	bobbalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

	suite.Require().Equal(sdk.NewUint(2), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(10),
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)
	suite.Require().Equal(sdk.NewUint(10), fetchedBalance[0].Amount)

	//Create badge 2 with supply == 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(1),
			Amount: sdk.NewUint(1),
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))
	bobbalance, _ = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

	suite.Require().Equal(sdk.NewUint(3), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(2), End: sdk.NewUint(2)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(1),
		},
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(10),
		},
	}, badge.MaxSupplys)
	bobbalance, _ = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)
	suite.Require().Equal(sdk.NewUint(1), bobbalance.Balances[0].Amount)
	suite.Require().Equal(sdk.NewUint(2), bobbalance.Balances[0].BadgeIds[0].Start)

	//Create badge 2 with supply == 10
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10),
			Amount: sdk.NewUint(2),
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))
	bobbalance, _ = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

	suite.Require().Equal(sdk.NewUint(5), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(2), End: sdk.NewUint(2)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(1),
		},
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}, {Start: sdk.NewUint(3), End: sdk.NewUint(4)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(10),
		},
	}, badge.MaxSupplys)
	suite.Require().Equal(sdk.NewUint(10), bobbalance.Balances[1].Amount)
	suite.Require().Equal(sdk.NewUint(1), bobbalance.Balances[1].BadgeIds[0].Start)
	suite.Require().Equal(sdk.NewUint(1), bobbalance.Balances[1].BadgeIds[0].End)
	suite.Require().Equal(sdk.NewUint(3), bobbalance.Balances[1].BadgeIds[1].Start)
	suite.Require().Equal(sdk.NewUint(4), bobbalance.Balances[1].BadgeIds[1].End)
}

func (suite *TestSuite) TestNewBadgesDirectlyUponCreatingNewBadge() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionMetadata: "https://example.com",
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(62),
				ApprovedTransfers: []*types.CollectionApprovedTransfer{
					{
						From: &types.AddressMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10),
			Amount: sdk.NewUint(1),
		},
	})

	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

	bobbalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

	suite.Require().Equal(sdk.NewUint(2), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(10),
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)
	suite.Require().Equal(sdk.NewUint(10), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	//Create badge 2 with supply == 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(1),
			Amount: sdk.NewUint(1),
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))
	bobbalance, _ = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

	suite.Require().Equal(sdk.NewUint(3), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(2), End: sdk.NewUint(2)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(1),
		},
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(10),
		},
	}, badge.MaxSupplys)
	suite.Require().Equal(sdk.NewUint(1), bobbalance.Balances[0].Amount)
	suite.Require().Equal(sdk.NewUint(2), bobbalance.Balances[0].BadgeIds[0].Start)

	//Create badge 2 with supply == 10
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10),
			Amount: sdk.NewUint(2),
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))
	bobbalance, _ = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

	suite.Require().Equal(sdk.NewUint(5), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(2), End: sdk.NewUint(2)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(1),
		},
		{
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}, {Start: sdk.NewUint(3), End: sdk.NewUint(4)}}, //0 to 0 range so it will be nil
			Amount:   sdk.NewUint(10),
		},
	},
		badge.MaxSupplys)
	suite.Require().Equal(sdk.NewUint(10), bobbalance.Balances[1].Amount)
	suite.Require().Equal(sdk.NewUint(1), bobbalance.Balances[1].BadgeIds[0].Start)
	suite.Require().Equal(sdk.NewUint(1), bobbalance.Balances[1].BadgeIds[0].End)
	suite.Require().Equal(sdk.NewUint(3), bobbalance.Balances[1].BadgeIds[1].Start)
	suite.Require().Equal(sdk.NewUint(4), bobbalance.Balances[1].BadgeIds[1].End)
}

func (suite *TestSuite) TestNewBadgesNotManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionMetadata: "https://example.com",
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(62),
				ApprovedTransfers: []*types.CollectionApprovedTransfer{
					{
						From: &types.AddressMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	err := CreateBadgesAndMintAllToCreator(suite, wctx, alice, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10),
			Amount: sdk.NewUint(1),
		},
	})
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}

func (suite *TestSuite) TestNewBadgeBadgeNotExists() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateBadgesAndMintAllToCreator(suite, wctx, alice, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10),
			Amount: sdk.NewUint(1),
		},
	})
	suite.Require().EqualError(err, keeper.ErrCollectionNotExists.Error())
}

func (suite *TestSuite) TestNewBadgeCreateIsLocked() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionMetadata: "https://example.com",
				BadgeMetadata: []*types.BadgeMetadata{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(0),
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10),
			Amount: sdk.NewUint(1),
		},
	})
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}
