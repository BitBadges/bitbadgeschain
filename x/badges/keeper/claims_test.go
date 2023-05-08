package keeper_test

// func (suite *TestSuite) TestSendAllToClaims() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				CollectionUri: "https://example.com",
// 				BadgeUris: []*types.BadgeUri{
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

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

// 	claimToAdd := types.Claim{
// 		UndistributedBalances: []*types.Balance{{
// 			Amount: sdk.NewUint(10),
// 			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}},
// 		}},
// 	}

// 	err := CreateBadges(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(10),
// 			Amount: sdk.NewUint(1),
// 		},
// 	},
// 		[]*types.Transfer{},
// 		[]*types.Claim{
// 			&claimToAdd,
// 		}, "https://example.com",
// 		[]*types.BadgeUri{
// 			{
// 				Uri: "https://example.com/{id}",
// 				BadgeIds: []*types.IdRange{
// 					{
// 						Start: sdk.NewUint(1),
// 						End:   sdk.NewUint(math.MaxUint64),
// 					},
// 				},
// 			},
// 		}, "")
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// 	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
// 			Amount: sdk.NewUint(10),
// 		},
// 	}, badge.MaxSupplys)

// 	claim, _ := GetClaim(suite, wctx, sdk.NewUint(1), sdk.NewUint(1))
// 	suite.Require().Nil(err, "Error getting claim")
// 	suite.Require().Equal(claimToAdd, claim)
// }
