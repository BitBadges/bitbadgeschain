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
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(62),
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
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10000),
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
			Amount: sdk.NewUint(10000),
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)
	suite.Require().Equal(sdk.NewUint(10000), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, alice, sdk.NewUint(1), []*types.Balance{
		{
			Amount: sdk.NewUint(1000),
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	bobbalance, _ = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)
	suite.Require().Equal(alice, bobbalance.Approvals[0].Address)
	suite.Require().Equal(sdk.NewUint(1000), bobbalance.Approvals[0].Balances[0].Amount)

	err = SetApproval(suite, wctx, bob, charlie, sdk.NewUint(1), []*types.Balance{
		{
			Amount: sdk.NewUint(500),
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	err = SetApproval(suite, wctx, bob, alice, sdk.NewUint(1), []*types.Balance{
		{
			Amount: sdk.NewUint(500),
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	bobbalance, _ = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

	suite.Require().Equal(charlie, bobbalance.Approvals[1].Address)
	suite.Require().Equal(sdk.NewUint(500), bobbalance.Approvals[1].Balances[0].Amount)

	suite.Require().Equal(alice, bobbalance.Approvals[0].Address)
	suite.Require().Equal(sdk.NewUint(500), bobbalance.Approvals[0].Balances[0].Amount)
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
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(62),
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
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10000),
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
			Amount: sdk.NewUint(10000),
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)
	suite.Require().Equal(sdk.NewUint(10000), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, charlie, alice, sdk.NewUint(1), []*types.Balance{
		{
			Amount: sdk.NewUint(0),
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
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
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(62),
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
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
		{
			Supply: sdk.NewUint(10000),
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
			Amount: sdk.NewUint(10000),
		},
	}, badge.MaxSupplys)
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)
	suite.Require().Equal(sdk.NewUint(10000), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, bob, sdk.NewUint(1), []*types.Balance{
		{
			Amount: sdk.NewUint(1000),
			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
		},
	})
	suite.Require().EqualError(err, keeper.ErrAccountCanNotEqualCreator.Error())
}
