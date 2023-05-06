package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNewCollections() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := uint64(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

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

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	badge, _ := GetCollection(suite, wctx, 1)

	// Verify nextId increments correctly
	nextId := suite.app.BadgesKeeper.GetNextCollectionId(suite.ctx)
	suite.Require().Equal(uint64(2), nextId)

	// Verify badge details are correct
	suite.Require().Equal(uint64(1), badge.NextBadgeId)
	suite.Require().Equal("https://example.com", badge.CollectionUri)
	suite.Require().Equal([]*types.BadgeUri{
		{
			Uri: "https://example.com/{id}",
			BadgeIds: []*types.IdRange{
				{
					Start: 1,
					End:   math.MaxUint64,
				},
			},
		},
	}, badge.BadgeUris)
	suite.Require().Equal([]*types.Balance(nil), badge.MaxSupplys)
	suite.Require().Equal(bob, badge.Manager) //7 is the first ID it creates
	suite.Require().Equal(perms, badge.Permissions)
	suite.Require().Equal([]*types.TransferMapping(nil), badge.AllowedTransfers)
	suite.Require().Equal([]*types.TransferMapping(nil), badge.ManagerApprovedTransfers)
	suite.Require().Equal(uint64(1), badge.CollectionId)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	// Verify nextId increments correctly
	nextId = suite.app.BadgesKeeper.GetNextCollectionId(suite.ctx)
	suite.Require().Equal(uint64(3), nextId)
	badge, _ = GetCollection(suite, wctx, 2)
	suite.Require().Equal(uint64(2), badge.CollectionId)
}

func (suite *TestSuite) TestNewBadgesWhitelistRecipients() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := uint64(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

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
				BadgeSupplys: []*types.BadgeSupplyAndAmount{
					{
						Supply: 10,
						Amount: 10,
					},
				},
				Permissions: perms,
				Transfers: []*types.Transfers{
					{
						ToAddresses: []string{alice, charlie},
						Balances: []*types.Balance{
							{
								Amount: 5,
								BadgeIds: []*types.IdRange{
									{
										Start: 1,
										End:   5,
									},
								},
							},
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	// Verify nextId increments correctly
	nextId := suite.app.BadgesKeeper.GetNextCollectionId(suite.ctx)
	suite.Require().Equal(uint64(2), nextId)

	collection, _ := GetCollection(suite, wctx, 1)

	unmintedBalances := types.UserBalanceStore{
		Balances: collection.UnmintedSupplys,
	}

	suite.Require().Equal(uint64(10), unmintedBalances.Balances[0].Amount)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: 6,
			End:   10,
		},
	}, unmintedBalances.Balances[0].BadgeIds)

	aliceBalance, _ := GetUserBalance(suite, wctx, 1, alice)
	suite.Require().Equal(uint64(5), aliceBalance.Balances[0].Amount)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: 1,
			End:   5,
		},
	}, aliceBalance.Balances[0].BadgeIds)

	charlieBalance, _ := GetUserBalance(suite, wctx, 1, charlie)
	suite.Require().Equal(uint64(5), charlieBalance.Balances[0].Amount)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: 1,
			End:   5,
		},
	}, charlieBalance.Balances[0].BadgeIds)
}

func (suite *TestSuite) TestNewBadgesWhitelistRecipientsOverflow() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := uint64(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

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
				BadgeSupplys: []*types.BadgeSupplyAndAmount{
					{
						Supply: 10,
						Amount: 10,
					},
				},
				Permissions: perms,
				Transfers: []*types.Transfers{
					{
						ToAddresses: []string{alice, charlie},
						Balances: []*types.Balance{
							{
								Amount: 6,
								BadgeIds: []*types.IdRange{
									{
										Start: 0,
										End:   4,
									},
								},
							},
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
}
