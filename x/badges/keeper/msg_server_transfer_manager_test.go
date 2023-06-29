package keeper_test

// import (
// sdkmath "cosmossdk.io/math"
// 	"math"

// 	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
// 	"github.com/bitbadges/bitbadgeschain/x/badges/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// func (suite *TestSuite) TestUpdateManager() {
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
// 				Permissions: sdkmath.NewUint(127),
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
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge")

// 	//Create badge 1 with supply > 1
// 	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(1),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badge")

// 	err = RequestUpdateManager(suite, wctx, alice, sdkmath.NewUint(1), true)
// 	suite.Require().Nil(err, "Error requesting manager transfer")

// 	err = UpdateManager(suite, wctx, bob, sdkmath.NewUint(1), alice)
// 	suite.Require().Nil(err, "Error transferring manager")

// 	badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	suite.Require().Equal(alice, badge.Manager)
// }

// func (suite *TestSuite) TestRequestUpdateManager() {
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
// 				Permissions: sdkmath.NewUint(127),
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
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge")

// 	//Create badge 1 with supply > 1
// 	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(1),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badge")

// 	err = RequestUpdateManager(suite, wctx, alice, sdkmath.NewUint(1), true)
// 	suite.Require().Nil(err, "Error requesting manager transfer")

// 	err = RequestUpdateManager(suite, wctx, alice, sdkmath.NewUint(1), false)
// 	suite.Require().Nil(err, "Error requesting manager transfer")

// 	err = RequestUpdateManager(suite, wctx, alice, sdkmath.NewUint(1), true)
// 	suite.Require().Nil(err, "Error requesting manager transfer")

// 	err = UpdateManager(suite, wctx, bob, sdkmath.NewUint(1), alice)
// 	suite.Require().Nil(err, "Error transferring manager")

// 	badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	suite.Require().Equal(alice, badge.Manager)
// }

// func (suite *TestSuite) TestRemovedRequestUpdateManager() {
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
// 				Permissions: sdkmath.NewUint(127),
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
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge")

// 	//Create badge 1 with supply > 1
// 	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(1),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badge")

// 	err = RequestUpdateManager(suite, wctx, alice, sdkmath.NewUint(1), true)
// 	suite.Require().Nil(err, "Error requesting manager transfer")

// 	err = RequestUpdateManager(suite, wctx, alice, sdkmath.NewUint(1), false)
// 	suite.Require().Nil(err, "Error requesting manager transfer")

// 	err = UpdateManager(suite, wctx, bob, sdkmath.NewUint(1), alice)
// 	suite.Require().EqualError(err, keeper.ErrAddressNeedsToOptInAndRequestManagerTransfer.Error())
// }

// func (suite *TestSuite) TestRemovedRequestUpdateManagerBadPermissions() {
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
// 				Permissions: sdkmath.NewUint(23),
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
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge")

// 	//Create badge 1 with supply > 1
// 	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(1),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badge")

// 	err = RequestUpdateManager(suite, wctx, alice, sdkmath.NewUint(1), true)
// 	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
// }

// func (suite *TestSuite) TestManagerCantBeTransferred() {
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
// 				Permissions: sdkmath.NewUint(0),
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "Error creating badge")

// 	err = UpdateManager(suite, wctx, bob, sdkmath.NewUint(1), alice)
// 	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
// }
