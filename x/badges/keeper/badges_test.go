package keeper_test

// import (
// 	"math"

// 	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
// 	"github.com/bitbadges/bitbadgeschain/x/badges/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// func (suite *TestSuite) TestGetCollection() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				CollectionMetadata: "https://example.com/",
// 				BadgeMetadata: []*types.BadgeMetadata{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdk.NewUint(1),
// 								End:   sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdk.NewUint(62),
// 			},
// 			Amount:  sdk.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	badge, err := suite.app.BadgesKeeper.GetCollectionE(suite.ctx, sdk.NewUint(1))
// 	suite.Require().Nil(err, "Error getting badge: %s")
// 	suite.Require().Equal(badge.CollectionId, sdk.NewUint(1))

// 	badge, err = suite.app.BadgesKeeper.GetCollectionE(suite.ctx, sdk.NewUint(2))
// 	suite.Require().EqualError(err, keeper.ErrCollectionNotExists.Error())
// }

// func (suite *TestSuite) TestGetBadgeAndAssertBadges() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				BadgeMetadata: []*types.BadgeMetadata{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdk.NewUint(1),
// 								End:   sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				CollectionMetadata: "https://example.com",
// 				Permissions:        sdk.NewUint(62),
// 				BadgesToCreate: []*types.BadgeSupplyAndAmount{
// 					{
// 						Supply: sdk.NewUint(1),
// 						Amount: sdk.NewUint(1),
// 					},
// 				},
// 			},
// 			Amount:  sdk.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")

// 	_, err = suite.app.BadgesKeeper.GetCollectionAndAssertBadgeIdsAreValid(suite.ctx, sdk.NewUint(1), []*types.IdRange{
// 		{
// 			Start: sdk.NewUint(1),
// 			End:   sdk.NewUint(1),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error getting badge: %s")

// 	_, err = suite.app.BadgesKeeper.GetCollectionAndAssertBadgeIdsAreValid(suite.ctx, sdk.NewUint(1), []*types.IdRange{
// 		{
// 			Start: sdk.NewUint(20),
// 			End:   sdk.NewUint(10),
// 		},
// 	})
// 	suite.Require().EqualError(err, keeper.ErrInvalidBadgeRange.Error())

// 	_, err = suite.app.BadgesKeeper.GetCollectionAndAssertBadgeIdsAreValid(suite.ctx, sdk.NewUint(1), []*types.IdRange{
// 		{
// 			Start: sdk.NewUint(1),
// 			End:   sdk.NewUint(10),
// 		},
// 	})
// 	suite.Require().EqualError(err, keeper.ErrBadgeNotExists.Error())
// }

// func (suite *TestSuite) TestCreateBadges() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				BadgeMetadata: []*types.BadgeMetadata{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdk.NewUint(1),
// 								End:   sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				CollectionMetadata: "https://example.com",
// 				Permissions:        sdk.NewUint(62),
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							IncludeOnlySpecified: false,
// 						},
// 						To: &types.AddressMapping{
// 							IncludeOnlySpecified: false,
// 						},
// 					},
// 				},
// 			},
// 			Amount:  sdk.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge: %s")
// 	badge, err := GetCollection(suite, wctx, sdk.NewUint(1))
// 	suite.Require().Nil(err, "Error getting badge: %s")
// 	balance := types.UserBalanceStore{}

// 	badge, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, badge, []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(1),
// 			Amount: sdk.NewUint(1),
// 		},
// 	}, []*types.Transfer{
// 		{
// 			ToAddresses: []string{bob},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdk.NewUint(1),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdk.NewUint(1),
// 							End:   sdk.NewUint(1),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}, []*types.Claim{}, bob)
// 	suite.Require().Nil(err, "Error creating subassets: %s")

// 	suite.Require().Equal(badge.MaxSupplys, []*types.Balance{
// 		{
// 			Amount: sdk.NewUint(1),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdk.NewUint(1),
// 					End:   sdk.NewUint(1),
// 				},
// 			},
// 		},
// 	})

// 	balance, err = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)
// 	suite.Require().Nil(err, "Error getting user balance: %s")
// 	suite.Require().Equal(balance.Balances[0].Amount, sdk.NewUint(1))
// 	suite.Require().Equal(balance.Balances[0].BadgeIds, []*types.IdRange{
// 		{
// 			Start: sdk.NewUint(1),
// 			End:   sdk.NewUint(1),
// 		},
// 	})

// 	badge, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, badge, []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(1),
// 			Amount: sdk.NewUint(1),
// 		},
// 	}, []*types.Transfer{
// 		{
// 			ToAddresses: []string{bob},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdk.NewUint(1),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdk.NewUint(2),
// 							End:   sdk.NewUint(2),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}, []*types.Claim{}, bob)
// 	suite.Require().Nil(err, "Error getting user balance: %s")

// 	balance, err = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)
// 	suite.Require().Nil(err, "Error getting user balance: %s")
// 	suite.Require().Nil(err, "Error creating subassets: %s")
// 	suite.Require().Equal(badge.MaxSupplys, []*types.Balance{
// 		{
// 			Amount: sdk.NewUint(1),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdk.NewUint(1),
// 					End:   sdk.NewUint(2),
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Equal(balance.Balances[0].Amount, sdk.NewUint(1))
// 	suite.Require().Equal(balance.Balances[0].BadgeIds, []*types.IdRange{
// 		{
// 			Start: sdk.NewUint(1),
// 			End:   sdk.NewUint(2),
// 		},
// 	})

// 	badge, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, badge, []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(1),
// 			Amount: sdk.NewUint(1),
// 		},
// 	}, []*types.Transfer{
// 		{
// 			ToAddresses: []string{bob},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdk.NewUint(1),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdk.NewUint(3),
// 							End:   sdk.NewUint(3),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}, []*types.Claim{}, bob)
// 	suite.Require().Nil(err, "Error getting user balance: %s")

// 	balance, err = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)
// 	suite.Require().Nil(err, "Error getting user balance: %s")
// 	suite.Require().Nil(err, "Error creating subassets: %s")
// 	suite.Require().Equal(badge.MaxSupplys, []*types.Balance{
// 		{
// 			Amount: sdk.NewUint(1),
// 			BadgeIds: []*types.IdRange{
// 				{
// 					Start: sdk.NewUint(1),
// 					End:   sdk.NewUint(3),
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Equal(balance.Balances[0].Amount, sdk.NewUint(1))
// 	suite.Require().Equal(balance.Balances[0].BadgeIds, []*types.IdRange{
// 		{
// 			Start: sdk.NewUint(1),
// 			End:   sdk.NewUint(3),
// 		},
// 	})

// 	badge, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, badge, []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(1),
// 			Amount: types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
// 		},
// 	}, []*types.Transfer{
// 		{
// 			ToAddresses: []string{bob},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdk.NewUint(1),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdk.NewUint(4),
// 							End:   sdk.NewUint(4),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}, []*types.Claim{}, bob)
// 	// suite.Require().EqualError(err, keeper.ErrOverflow.Error())
// }
