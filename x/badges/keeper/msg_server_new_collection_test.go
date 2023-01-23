package keeper_test

import (
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
				BadgeUri: "https://example.com/{id}",
				CollectionUri: "https://example.com",
				Permissions: 62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	badge, _ := GetCollection(suite, wctx, 0)

	// Verify nextId increments correctly
	nextId := suite.app.BadgesKeeper.GetNextCollectionId(suite.ctx)
	suite.Require().Equal(uint64(1), nextId)

	// Verify badge details are correct
	suite.Require().Equal(uint64(0), badge.NextBadgeId)
	suite.Require().Equal("https://example.com" , badge.CollectionUri)
	suite.Require().Equal("https://example.com/{id}" , badge.BadgeUri)
	suite.Require().Equal([]*types.Balance(nil), badge.MaxSupplys)
	suite.Require().Equal(bobAccountNum, badge.Manager) //7 is the first ID it creates
	suite.Require().Equal(perms, badge.Permissions)
	suite.Require().Equal([]*types.TransferMapping(nil), badge.DisallowedTransfers)
	suite.Require().Equal([]*types.TransferMapping(nil), badge.ManagerApprovedTransfers)
	suite.Require().Equal(uint64(0), badge.CollectionId)

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	// Verify nextId increments correctly
	nextId = suite.app.BadgesKeeper.GetNextCollectionId(suite.ctx)
	suite.Require().Equal(uint64(2), nextId)
	badge, _ = GetCollection(suite, wctx, 1)
	suite.Require().Equal(uint64(1), badge.CollectionId)
}


func (suite *TestSuite) TestNewBadgesWhitelistRecipients() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := uint64(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				BadgeUri: "https://example.com/{id}",
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
						ToAddresses: []uint64{aliceAccountNum, charlieAccountNum},
						Balances: []*types.Balance{
							{
								Balance: 5,
								BadgeIds: []*types.IdRange{
									{
										Start: 0,
										End: 4,
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
	suite.Require().Equal(uint64(1), nextId)

	collection, _ := GetCollection(suite, wctx, 0)
	
	unmintedBalances := types.UserBalance{
		Balances: collection.UnmintedSupplys,
	}

	suite.Require().Equal(uint64(10), unmintedBalances.Balances[0].Balance)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: 5,
			End: 9,
		},
	}, unmintedBalances.Balances[0].BadgeIds)

	aliceBalance, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5), aliceBalance.Balances[0].Balance)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: 0,
			End: 4,
		},
	}, aliceBalance.Balances[0].BadgeIds)

	charlieBalance, _ := GetUserBalance(suite, wctx, 0, charlieAccountNum)
	suite.Require().Equal(uint64(5), charlieBalance.Balances[0].Balance)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: 0,
			End: 4,
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
				BadgeUri: "https://example.com/{id}",
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
						ToAddresses: []uint64{aliceAccountNum, charlieAccountNum},
						Balances: []*types.Balance{
							{
								Balance: 6,
								BadgeIds: []*types.IdRange{
									{
										Start: 0,
										End: 4,
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