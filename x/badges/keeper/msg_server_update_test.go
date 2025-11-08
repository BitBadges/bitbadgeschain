package keeper_test

//We have some tests in update_checks_test that call MsgUpdateMetadata

// import (
// sdkmath "cosmossdk.io/math"
// 	"math"

// 	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
// 	"github.com/bitbadges/bitbadgeschain/x/badges/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// func (suite *TestSuite) TestUpdateURIs() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	_, err := sdk.AccAddressFromBech32(alice)
// 	suite.Require().Nil(err, "Address %s failed to parse")

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				TokenMetadata: []*types.TokenMetadata{
// 					{
// 						Uri: "https://example.com/{id}",
// 						TokenIds: []*types.UintRange{
// 							{
// 								Start: sdkmath.NewUint(1),
// 								End:   sdkmath.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				CollectionMetadata: "https://example.com",
// 				Permissions:        sdkmath.NewUint(62),
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err = CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating token: %s")

// 	err = UpdateURIs(suite, wctx, bob, sdkmath.NewUint(1), "https://example.com", []*types.TokenMetadata{
// 		{
// 			Uri: "https://example.com/{id}",
// 			TokenIds: []*types.UintRange{
// 				{
// 					Start: sdkmath.NewUint(1),
// 					End:   sdkmath.NewUint(math.MaxUint64),
// 				},
// 			},
// 		},
// 	}, "")
// 	suite.Require().Nil(err, "Error updating uris")
// 	badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	suite.Require().Equal("https://example.com", badge.CollectionMetadata)
// 	// suite.Require().Equal("https://example.com/{id}", badge.TokenMetadata)

// 	err = UpdateCollectionPermissions(suite, wctx, bob, sdkmath.NewUint(1), sdkmath.NewUint(60-32))
// 	suite.Require().Nil(err, "Error updating permissions")

// 	err = UpdateCustomData(suite, wctx, bob, sdkmath.NewUint(1), "example.com/")
// 	suite.Require().Nil(err, "Error updating bytes")
// }

// func (suite *TestSuite) TestCantUpdate() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	_, err := sdk.AccAddressFromBech32(alice)
// 	suite.Require().Nil(err, "Address %s failed to parse")

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				CollectionMetadata: "https://example.com",
// 				TokenMetadata: []*types.TokenMetadata{
// 					{
// 						Uri: "https://example.com/{id}",
// 						TokenIds: []*types.UintRange{
// 							{
// 								Start: sdkmath.NewUint(1),
// 								End:   sdkmath.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdkmath.NewUint(0),
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err = CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating token: %s")

// 	err = UpdateURIs(suite, wctx, bob, sdkmath.NewUint(1), "https://example.com/test2222", []*types.TokenMetadata{
// 		{
// 			Uri: "https://example.com/{id}/edited",
// 			TokenIds: []*types.UintRange{
// 				{
// 					Start: sdkmath.NewUint(1),
// 					End:   sdkmath.NewUint(math.MaxUint64),
// 				},
// 			},
// 		},
// 	}, "")
// 	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())

// 	err = UpdateCollectionPermissions(suite, wctx, bob, sdkmath.NewUint(1), sdkmath.NewUint(123-64))
// 	suite.Require().EqualError(err, types.ErrInvalidPermissionsUpdateLocked.Error())

// 	err = UpdateCustomData(suite, wctx, bob, sdkmath.NewUint(1), "example.com/")
// 	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
// }

// func (suite *TestSuite) TestCantUpdateNotManager() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	_, err := sdk.AccAddressFromBech32(alice)
// 	suite.Require().Nil(err, "Address %s failed to parse")

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				CollectionMetadata: "https://example.com",
// 				TokenMetadata: []*types.TokenMetadata{
// 					{
// 						Uri: "https://example.com/{id}",
// 						TokenIds: []*types.UintRange{
// 							{
// 								Start: sdkmath.NewUint(1),
// 								End:   sdkmath.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdkmath.NewUint(0),
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err = CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating token: %s")

// 	err = UpdateURIs(suite, wctx, alice, sdkmath.NewUint(1), "https://example.com", []*types.TokenMetadata{
// 		{
// 			Uri: "https://example.com/{id}",
// 			TokenIds: []*types.UintRange{
// 				{
// 					Start: sdkmath.NewUint(1),
// 					End:   sdkmath.NewUint(math.MaxUint64),
// 				},
// 			},
// 		},
// 	}, "")
// 	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

// 	err = UpdateCollectionPermissions(suite, wctx, alice, sdkmath.NewUint(1), sdkmath.NewUint(77))
// 	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

// 	err = UpdateCustomData(suite, wctx, alice, sdkmath.NewUint(1), "example.com/")
// 	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
// }

// func (suite *TestSuite) TestUpdateBalanceURIs() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	_, err := sdk.AccAddressFromBech32(alice)
// 	suite.Require().Nil(err, "Address %s failed to parse")

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				TokenMetadata: []*types.TokenMetadata{
// 					{
// 						Uri: "https://example.com/{id}",
// 						TokenIds: []*types.UintRange{
// 							{
// 								Start: sdkmath.NewUint(1),
// 								End:   sdkmath.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				CollectionMetadata: "https://example.com",
// 				Permissions:        sdkmath.NewUint(62 + 64),
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err = CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating token: %s")

// 	err = UpdateURIs(suite, wctx, bob, sdkmath.NewUint(1), "", []*types.TokenMetadata{}, "https://balance.com/{id}")
// 	suite.Require().Nil(err, "Error updating uris")
// 	badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	suite.Require().Equal("https://balance.com/{id}", badge.OffChainBalancesMetadata)
// 	// suite.Require().Equal("https://example.com/{id}", badge.TokenMetadata)
// }

// func (suite *TestSuite) TestCantUpdateBalanceUris() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	_, err := sdk.AccAddressFromBech32(alice)
// 	suite.Require().Nil(err, "Address %s failed to parse")

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Collection: types.MsgNewCollection{
// 				CollectionMetadata: "https://example.com",
// 				TokenMetadata: []*types.TokenMetadata{
// 					{
// 						Uri: "https://example.com/{id}",
// 						TokenIds: []*types.UintRange{
// 							{
// 								Start: sdkmath.NewUint(1),
// 								End:   sdkmath.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdkmath.NewUint(0),
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err = CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating token: %s")

// 	err = UpdateURIs(suite, wctx, bob, sdkmath.NewUint(1), "", []*types.TokenMetadata{}, "https://balance.com/{id}")
// 	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
// }
