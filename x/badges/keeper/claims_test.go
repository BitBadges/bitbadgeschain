package keeper_test

// func (suite *TestSuite) TestSendAllToClaims() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				CollectionMetadata: "https://example.com",
// 				BadgeMetadata: []*types.BadgeMetadata{
// 					{
// 						Uri: "https://example.com/{id}",
// 						BadgeIds: []*types.IdRange{
// 							{
// 								Start: sdkmath.NewUint(1),
// 								End:   sdkmath.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdkmath.NewUint(62),
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

// 	claimToAdd := types.Claim{
// 		Balances: []*types.Balance{{
// 			Amount: sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 		}},
// 	}

// 	err := CreateBadges(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10),
// 			Amount: sdkmath.NewUint(1),
// 		},
// 	},
// 		[]*types.Transfer{},
// 		[]*types.Claim{
// 			&claimToAdd,
// 		}, "https://example.com",
// 		[]*types.BadgeMetadata{
// 			{
// 				Uri: "https://example.com/{id}",
// 				BadgeIds: []*types.IdRange{
// 					{
// 						Start: sdkmath.NewUint(1),
// 						End:   sdkmath.NewUint(math.MaxUint64),
// 					},
// 				},
// 			},
// 		}, "")
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

// 	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, //0 to 0 range so it will be nil
// 			Amount: sdkmath.NewUint(10),
// 		},
// 	}, badge.MaxSupplys)

// 	claim, _ := GetClaim(suite, wctx, sdkmath.NewUint(1), sdkmath.NewUint(1))
// 	suite.Require().Nil(err, "Error getting claim")
// 	suite.Require().Equal(claimToAdd, claim)
// }
