package keeper_test

// import (
// sdkmath "cosmossdk.io/math"
// 	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
// 	"github.com/bitbadges/bitbadgeschain/x/badges/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/ethereum/go-ethereum/common/math"
// )

// func (suite *TestSuite) TestFreezeAddressesDirectlyWhenCreatingNewBadge() {
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
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(uint64(types.AddressOptions_None)),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(uint64(types.AddressOptions_None)),
// 						},
// 					},
// 				},
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)

// 	//Create badge 1 with supply > 1
// 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(1),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badge")
// 	// badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(5000),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(1),
// 							End:   sdkmath.NewUint(1),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transferring badge")

// 	err = UpdateCollectionApprovedTransfers(suite, wctx, bob, sdkmath.NewUint(1), []*types.CollectionApprovedTransfer{
// 		{
// 			From: &types.AddressMapping{
// 				IncludeOnlySpecified: false,
// 			},
// 			To: &types.AddressMapping{
// 				IncludeOnlySpecified: false,
// 				Addresses:            []string{bob},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error updating allowed transfers")

// 	err = TransferBadge(suite, wctx, alice, sdkmath.NewUint(1), alice, []*types.Transfer{
// 		{
// 			ToAddresses: []string{bob},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(5000),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(1),
// 							End:   sdkmath.NewUint(1),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
// }

// func (suite *TestSuite) TestTransferBadgeForcefulUnfrozenByDefault() {
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
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(uint64(types.AddressOptions_None)),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(uint64(types.AddressOptions_None)),
// 						},
// 					},
// 				},
// 			},
// 			Amount:  sdkmath.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

// 	//Create badge 1 with supply > 1
// 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdkmath.NewUint(10000),
// 			Amount: sdkmath.NewUint(1),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

// 	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

// 	suite.Require().Equal(sdkmath.NewUint(2), badge.NextBadgeId)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, //0 to 0 range so it will be nil
// 			Amount:   sdkmath.NewUint(10000),
// 		},
// 	}, badge.MaxSupplys)
// 	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Nil(err)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(5000),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(1),
// 							End:   sdkmath.NewUint(1),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transferring badge")

// 	bobbalance, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
// 	fetchedBalance, err = keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(5000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	alicebalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)

// 	fetchedBalance, err = keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, alicebalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(5000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	err = UpdateCollectionApprovedTransfers(suite, wctx, bob, sdkmath.NewUint(1), []*types.CollectionApprovedTransfer{
// 		{
// 			From: &types.AddressMapping{
// 				Addresses:            []string{},
// 				IncludeOnlySpecified: true,
// 				ManagerOptions:       sdkmath.NewUint(uint64(types.AddressOptions_None)),
// 			},
// 			To: &types.AddressMapping{
// 				Addresses:            []string{},
// 				IncludeOnlySpecified: true,
// 				ManagerOptions:       sdkmath.NewUint(uint64(types.AddressOptions_None)),
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error freezing address")

// 	badge, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

// 	err = TransferBadge(suite, wctx, alice, sdkmath.NewUint(1), alice, []*types.Transfer{
// 		{
// 			ToAddresses: []string{bob},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(5000),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(1),
// 							End:   sdkmath.NewUint(1),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
// }

// //TODO:
// //TODO: also test transfer approvedTransfers

// //TODO: also test manager approved transfers with transfer approvedTransfers

// // func (suite *TestSuite) TestTransferBadgeForcefulFrozenByDefault() {
// // 	wctx := sdk.WrapSDKContext(suite.ctx)

// // 	collectionsToCreate := []CollectionsToCreate{
// // 		{
// // 			Badge: types.MsgNewCollection{
// // 				CollectionMetadata: "https://example.com",
// // BadgeMetadata: "https://example.com/{id}",
// // 				Permissions: 63,
// // 			},
// // 			Amount:  sdkmath.NewUint(1),
// // 			Creator: bob,
// // 		},
// // 	}

// // 	CreateCollections(suite, wctx, collectionsToCreate)
// // 	badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

// // 	//Create badge 1 with supply > 1
// // 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// // 		{
// // 			Supply: sdkmath.NewUint(10000),
// // 			Amount: sdkmath.NewUint(1),
// // 		},
// // 	},)
// // 	suite.Require().Nil(err, "Error creating badge")
// // 	badge, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

// // 	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

// // 	suite.Require().Equal(sdkmath.NewUint(2), badge.NextBadgeId)
// // 	suite.Require().Equal([]*types.Balance{
// // 		{
// // 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, //0 to 0 range so it will be nil
// // 			Amount: sdkmath.NewUint(10000),
// // 		},
// // 	}, badge.MaxSupplys)
// // 	suite.Require().Equal(sdkmath.NewUint(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)[0].Amount)

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// // 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

// // 	err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.IdRange{{Start: bob, End: bob}})
// // 	suite.Require().Nil(err, "Error unfreezing address")

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// // 	suite.Require().Nil(err, "Error transferring after unfreeze")

// // 	err = TransferBadge(suite, wctx, alice, alice, []string{bob}, []string{5000}, 0, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// // 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
// // }

// // func (suite *TestSuite) TestTransferBadgeForcefulFrozenByDefaultAddAndRemove() {
// // 	wctx := sdk.WrapSDKContext(suite.ctx)

// // 	collectionsToCreate := []CollectionsToCreate{
// // 		{
// // 			Badge: types.MsgNewCollection{
// // 				CollectionMetadata: "https://example.com",
// // BadgeMetadata: "https://example.com/{id}",
// // 				Permissions: 63,
// // 			},
// // 			Amount:  sdkmath.NewUint(1),
// // 			Creator: bob,
// // 		},
// // 	}

// // 	CreateCollections(suite, wctx, collectionsToCreate)
// // 	badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

// // 	//Create badge 1 with supply > 1
// // 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// // 		{
// // 			Supply: sdkmath.NewUint(10000),
// // 			Amount: sdkmath.NewUint(1),
// // 		},
// // 	},)
// // 	suite.Require().Nil(err, "Error creating badge")
// // 	badge, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

// // 	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

// // 	suite.Require().Equal(sdkmath.NewUint(2), badge.NextBadgeId)
// // 	suite.Require().Equal([]*types.Balance{
// // 		{
// // 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, //0 to 0 range so it will be nil
// // 			Amount: sdkmath.NewUint(10000),
// // 		},
// // 	}, badge.MaxSupplys)
// // 	suite.Require().Equal(sdkmath.NewUint(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)[0].Amount)

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// // 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

// // 	err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.IdRange{{Start: bob, End: bob}})
// // 	suite.Require().Nil(err, "Error unfreezing address")

// // 	err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), false, []*types.IdRange{{Start: bob, End: sdkmath.NewUint(0)}})
// // 	suite.Require().Nil(err, "Error unfreezing address")

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// // 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

// // 	err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.IdRange{{Start: bob, End: bob}})
// // 	suite.Require().Nil(err, "Error unfreezing address")

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// // 	suite.Require().Nil(err, "Error transferring after unfreeze")

// // 	err = TransferBadge(suite, wctx, alice, alice, []string{bob}, []string{5000}, 0, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, 0, 0)
// // 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
// // }

// // func (suite *TestSuite) TestFreezeCantFreeze() {
// // 	wctx := sdk.WrapSDKContext(suite.ctx)

// // 	collectionsToCreate := []CollectionsToCreate{
// // 		{
// // 			Badge: types.MsgNewCollection{
// // 				CollectionMetadata: "https://example.com",
// // 				BadgeMetadata: "https://example.com/{id}",
// // 				Permissions: sdkmath.NewUint(0),
// // 			},
// // 			Amount:  sdkmath.NewUint(1),
// // 			Creator: bob,
// // 		},
// // 	}

// // 	CreateCollections(suite, wctx, collectionsToCreate)

// // 	err := FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.IdRange{{Start: bob, End: bob}})
// // 	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
// // }

// // func (suite *TestSuite) TestTransferBadgeForcefulUnfrozenByDefaultOmitEmptyCase() {
// // 	wctx := sdk.WrapSDKContext(suite.ctx)

// // 	collectionsToCreate := []CollectionsToCreate{
// // 		{
// // 			Badge: types.MsgNewCollection{
// // 				CollectionMetadata: "https://example.com",
// // BadgeMetadata: "https://example.com/{id}",
// // 				Permissions: sdkmath.NewUint(62),
// // 			},
// // 			Amount:  sdkmath.NewUint(1),
// // 			Creator: bob,
// // 		},
// // 	}

// // 	CreateCollections(suite, wctx, collectionsToCreate)
// // 	badge, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

// // 	//Create badge 1 with supply > 1
// // 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdkmath.NewUint(1), []*types.BadgeSupplyAndAmount{
// // 		{
// // 			Supply: sdkmath.NewUint(10000),
// // 			Amount: sdkmath.NewUint(1),
// // 		},
// // 	},)
// // 	suite.Require().Nil(err, "Error creating badge")
// // 	badge, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

// // 	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

// // 	suite.Require().Equal(sdkmath.NewUint(2), badge.NextBadgeId)
// // 	suite.Require().Equal([]*types.Balance{
// // 		{
// // 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, //0 to 0 range so it will be nil
// // 			Amount: sdkmath.NewUint(10000),
// // 		},
// // 	}, badge.MaxSupplys)
// // 	suite.Require().Equal(sdkmath.NewUint(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)[0].Amount)

// // 	err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}})
// // 	suite.Require().Nil(err, "Error freezing address")

// // 	badge, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

// // 	err = FreezeAddresses(suite, wctx, bob, sdkmath.NewUint(0), true, []*types.IdRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2),}})
// // 	suite.Require().Nil(err, "Error freezing address")

// // 	badge, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))
// // 	suite.Require().Equal(badge.FreezeRanges, []*types.IdRange{{Start: sdkmath.NewUint(0), End: sdkmath.NewUint(1)}})
// // }
