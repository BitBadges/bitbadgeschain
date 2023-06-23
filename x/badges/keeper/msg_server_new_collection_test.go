package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNewCollections() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := sdk.NewUint(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
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
				CollectionMetadata: "https://example.com",
				Permissions:        sdk.NewUint(62),
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	// Verify nextId increments correctly
	nextId := suite.app.BadgesKeeper.GetNextCollectionId(suite.ctx)
	suite.Require().Equal(sdk.NewUint(2), nextId)

	// Verify badge details are correct
	suite.Require().Equal(sdk.NewUint(1), badge.NextBadgeId)
	suite.Require().Equal("https://example.com", badge.CollectionMetadata)
	suite.Require().Equal([]*types.BadgeMetadata{
		{
			Uri: "https://example.com/{id}",
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(1),
					End:   sdk.NewUint(math.MaxUint64),
				},
			},
		},
	}, badge.BadgeMetadata)
	suite.Require().Equal([]*types.Balance(nil), badge.MaxSupplys)
	suite.Require().Equal(bob, badge.Manager) //7 is the first ID it creates
	suite.Require().Equal(perms, badge.Permissions)
	suite.Require().Equal([]*types.CollectionApprovedTransfer(nil), badge.ApprovedTransfers)
	suite.Require().Equal([]*types.CollectionApprovedTransfer(nil), badge.ManagerApprovedTransfers)
	suite.Require().Equal(sdk.NewUint(1), badge.CollectionId)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	// Verify nextId increments correctly
	nextId = suite.app.BadgesKeeper.GetNextCollectionId(suite.ctx)
	suite.Require().Equal(sdk.NewUint(3), nextId)
	badge, _ = GetCollection(suite, wctx, sdk.NewUint(2))
	suite.Require().Equal(sdk.NewUint(2), badge.CollectionId)
}

func (suite *TestSuite) TestNewBadgesWhitelistRecipients() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := sdk.NewUint(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
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
				CollectionMetadata: "https://example.com",
				BadgesToCreate: []*types.BadgeSupplyAndAmount{
					{
						Supply: sdk.NewUint(10),
						Amount: sdk.NewUint(10),
					},
				},
				Permissions: perms,
				ApprovedTransfers: []*types.CollectionApprovedTransfer{
					{
						From: &types.AddressMapping{
							IncludeOnlySpecified: false,
							ManagerOptions:       sdk.NewUint(0),
						},
						To: &types.AddressMapping{
							IncludeOnlySpecified: false,
							ManagerOptions:       sdk.NewUint(0),
						},
					},
				},
				Transfers: []*types.Transfer{
					{
						ToAddresses: []string{alice, charlie},
						Balances: []*types.Balance{
							{
								Amount: sdk.NewUint(5),
								BadgeIds: []*types.IdRange{
									{
										Start: sdk.NewUint(1),
										End:   sdk.NewUint(5),
									},
								},
							},
						},
					},
				},
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	// Verify nextId increments correctly
	nextId := suite.app.BadgesKeeper.GetNextCollectionId(suite.ctx)
	suite.Require().Equal(sdk.NewUint(2), nextId)

	collection, _ := GetCollection(suite, wctx, sdk.NewUint(1))

	unmintedBalances := types.UserBalanceStore{
		Balances: collection.UnmintedSupplys,
	}

	suite.Require().Equal(sdk.NewUint(10), unmintedBalances.Balances[0].Amount)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: sdk.NewUint(6),
			End:   sdk.NewUint(10),
		},
	}, unmintedBalances.Balances[0].BadgeIds)

	aliceBalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), alice)
	suite.Require().Equal(sdk.NewUint(5), aliceBalance.Balances[0].Amount)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(5),
		},
	}, aliceBalance.Balances[0].BadgeIds)

	charlieBalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), charlie)
	suite.Require().Equal(sdk.NewUint(5), charlieBalance.Balances[0].Amount)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: sdk.NewUint(1),
			End:   sdk.NewUint(5),
		},
	}, charlieBalance.Balances[0].BadgeIds)
}

func (suite *TestSuite) TestNewBadgesWhitelistRecipientsOverflow() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := sdk.NewUint(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
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
				CollectionMetadata: "https://example.com",
				BadgesToCreate: []*types.BadgeSupplyAndAmount{
					{
						Supply: sdk.NewUint(10),
						Amount: sdk.NewUint(10),
					},
				},
				Permissions: perms,
				ApprovedTransfers: []*types.CollectionApprovedTransfer{
					{
						From: &types.AddressMapping{
							IncludeOnlySpecified: false,
							ManagerOptions:       sdk.NewUint(0),
						},
						To: &types.AddressMapping{
							IncludeOnlySpecified: false,
							ManagerOptions:       sdk.NewUint(0),
						},
					},
				},
				Transfers: []*types.Transfer{
					{
						ToAddresses: []string{alice, charlie},
						Balances: []*types.Balance{
							{
								Amount: sdk.NewUint(6),
								BadgeIds: []*types.IdRange{
									{
										Start: sdk.NewUint(0),
										End:   sdk.NewUint(4),
									},
								},
							},
						},
					},
				},
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
}
