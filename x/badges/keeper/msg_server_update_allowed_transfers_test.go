package keeper_test

// import (
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
// 								Start: sdk.NewUint(1),
// 								End:   sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdk.NewUint(23),
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
// 						},
// 					},
// 				},
// 			},
// 			Amount:  sdk.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)

// 	//Create badge 1 with supply > 1
// 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(10000),
// 			Amount: sdk.NewUint(1),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badge")
// 	// badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

// 	err = TransferBadge(suite, wctx, bob, sdk.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdk.NewUint(5000),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdk.NewUint(1),
// 							End:   sdk.NewUint(1),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transferring badge")

// 	err = UpdateCollectionApprovedTransfers(suite, wctx, bob, sdk.NewUint(1), []*types.CollectionApprovedTransfer{
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

// 	err = TransferBadge(suite, wctx, alice, sdk.NewUint(1), alice, []*types.Transfer{
// 		{
// 			ToAddresses: []string{bob},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdk.NewUint(5000),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdk.NewUint(1),
// 							End:   sdk.NewUint(1),
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
// 								Start: sdk.NewUint(1),
// 								End:   sdk.NewUint(math.MaxUint64),
// 							},
// 						},
// 					},
// 				},
// 				Permissions: sdk.NewUint(23),
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
// 						},
// 					},
// 				},
// 			},
// 			Amount:  sdk.NewUint(1),
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

// 	//Create badge 1 with supply > 1
// 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: sdk.NewUint(10000),
// 			Amount: sdk.NewUint(1),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// 	bobbalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

// 	suite.Require().Equal(sdk.NewUint(2), badge.NextBadgeId)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
// 			Amount:   sdk.NewUint(10000),
// 		},
// 	}, badge.MaxSupplys)
// 	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Nil(err)
// 	suite.Require().Equal(sdk.NewUint(10000), fetchedBalance[0].Amount)

// 	err = TransferBadge(suite, wctx, bob, sdk.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdk.NewUint(5000),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdk.NewUint(1),
// 							End:   sdk.NewUint(1),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transferring badge")

// 	bobbalance, _ = GetUserBalance(suite, wctx, sdk.NewUint(1), bob)
// 	fetchedBalance, err = keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdk.NewUint(5000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	alicebalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), alice)

// 	fetchedBalance, err = keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, alicebalance.Balances)
// 	suite.Require().Equal(sdk.NewUint(5000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	err = UpdateCollectionApprovedTransfers(suite, wctx, bob, sdk.NewUint(1), []*types.CollectionApprovedTransfer{
// 		{
// 			From: &types.AddressMapping{
// 				Addresses:            []string{},
// 				IncludeOnlySpecified: true,
// 				ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
// 			},
// 			To: &types.AddressMapping{
// 				Addresses:            []string{},
// 				IncludeOnlySpecified: true,
// 				ManagerOptions:       sdk.NewUint(uint64(types.AddressOptions_None)),
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error freezing address")

// 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// 	err = TransferBadge(suite, wctx, alice, sdk.NewUint(1), alice, []*types.Transfer{
// 		{
// 			ToAddresses: []string{bob},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdk.NewUint(5000),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdk.NewUint(1),
// 							End:   sdk.NewUint(1),
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
// // 			Amount:  sdk.NewUint(1),
// // 			Creator: bob,
// // 		},
// // 	}

// // 	CreateCollections(suite, wctx, collectionsToCreate)
// // 	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

// // 	//Create badge 1 with supply > 1
// // 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
// // 		{
// // 			Supply: sdk.NewUint(10000),
// // 			Amount: sdk.NewUint(1),
// // 		},
// // 	},)
// // 	suite.Require().Nil(err, "Error creating badge")
// // 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// // 	bobbalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

// // 	suite.Require().Equal(sdk.NewUint(2), badge.NextBadgeId)
// // 	suite.Require().Equal([]*types.Balance{
// // 		{
// // 			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
// // 			Amount: sdk.NewUint(10000),
// // 		},
// // 	}, badge.MaxSupplys)
// // 	suite.Require().Equal(sdk.NewUint(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)[0].Amount)

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, 0, 0)
// // 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

// // 	err = FreezeAddresses(suite, wctx, bob, sdk.NewUint(0), true, []*types.IdRange{{Start: bob, End: bob}})
// // 	suite.Require().Nil(err, "Error unfreezing address")

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, 0, 0)
// // 	suite.Require().Nil(err, "Error transferring after unfreeze")

// // 	err = TransferBadge(suite, wctx, alice, alice, []string{bob}, []string{5000}, 0, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, 0, 0)
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
// // 			Amount:  sdk.NewUint(1),
// // 			Creator: bob,
// // 		},
// // 	}

// // 	CreateCollections(suite, wctx, collectionsToCreate)
// // 	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

// // 	//Create badge 1 with supply > 1
// // 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
// // 		{
// // 			Supply: sdk.NewUint(10000),
// // 			Amount: sdk.NewUint(1),
// // 		},
// // 	},)
// // 	suite.Require().Nil(err, "Error creating badge")
// // 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// // 	bobbalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

// // 	suite.Require().Equal(sdk.NewUint(2), badge.NextBadgeId)
// // 	suite.Require().Equal([]*types.Balance{
// // 		{
// // 			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
// // 			Amount: sdk.NewUint(10000),
// // 		},
// // 	}, badge.MaxSupplys)
// // 	suite.Require().Equal(sdk.NewUint(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)[0].Amount)

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, 0, 0)
// // 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

// // 	err = FreezeAddresses(suite, wctx, bob, sdk.NewUint(0), true, []*types.IdRange{{Start: bob, End: bob}})
// // 	suite.Require().Nil(err, "Error unfreezing address")

// // 	err = FreezeAddresses(suite, wctx, bob, sdk.NewUint(0), false, []*types.IdRange{{Start: bob, End: sdk.NewUint(0)}})
// // 	suite.Require().Nil(err, "Error unfreezing address")

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, 0, 0)
// // 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

// // 	err = FreezeAddresses(suite, wctx, bob, sdk.NewUint(0), true, []*types.IdRange{{Start: bob, End: bob}})
// // 	suite.Require().Nil(err, "Error unfreezing address")

// // 	err = TransferBadge(suite, wctx, bob, bob, []string{alice}, []string{5000}, 0, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, 0, 0)
// // 	suite.Require().Nil(err, "Error transferring after unfreeze")

// // 	err = TransferBadge(suite, wctx, alice, alice, []string{bob}, []string{5000}, 0, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, 0, 0)
// // 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
// // }

// // func (suite *TestSuite) TestFreezeCantFreeze() {
// // 	wctx := sdk.WrapSDKContext(suite.ctx)

// // 	collectionsToCreate := []CollectionsToCreate{
// // 		{
// // 			Badge: types.MsgNewCollection{
// // 				CollectionMetadata: "https://example.com",
// // 				BadgeMetadata: "https://example.com/{id}",
// // 				Permissions: sdk.NewUint(0),
// // 			},
// // 			Amount:  sdk.NewUint(1),
// // 			Creator: bob,
// // 		},
// // 	}

// // 	CreateCollections(suite, wctx, collectionsToCreate)

// // 	err := FreezeAddresses(suite, wctx, bob, sdk.NewUint(0), true, []*types.IdRange{{Start: bob, End: bob}})
// // 	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
// // }

// // func (suite *TestSuite) TestTransferBadgeForcefulUnfrozenByDefaultOmitEmptyCase() {
// // 	wctx := sdk.WrapSDKContext(suite.ctx)

// // 	collectionsToCreate := []CollectionsToCreate{
// // 		{
// // 			Badge: types.MsgNewCollection{
// // 				CollectionMetadata: "https://example.com",
// // BadgeMetadata: "https://example.com/{id}",
// // 				Permissions: sdk.NewUint(62),
// // 			},
// // 			Amount:  sdk.NewUint(1),
// // 			Creator: bob,
// // 		},
// // 	}

// // 	CreateCollections(suite, wctx, collectionsToCreate)
// // 	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))

// // 	//Create badge 1 with supply > 1
// // 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, sdk.NewUint(1), []*types.BadgeSupplyAndAmount{
// // 		{
// // 			Supply: sdk.NewUint(10000),
// // 			Amount: sdk.NewUint(1),
// // 		},
// // 	},)
// // 	suite.Require().Nil(err, "Error creating badge")
// // 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// // 	bobbalance, _ := GetUserBalance(suite, wctx, sdk.NewUint(1), bob)

// // 	suite.Require().Equal(sdk.NewUint(2), badge.NextBadgeId)
// // 	suite.Require().Equal([]*types.Balance{
// // 		{
// // 			BadgeIds: []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, //0 to 0 range so it will be nil
// // 			Amount: sdk.NewUint(10000),
// // 		},
// // 	}, badge.MaxSupplys)
// // 	suite.Require().Equal(sdk.NewUint(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}}, bobbalance.Balances)[0].Amount)

// // 	err = FreezeAddresses(suite, wctx, bob, sdk.NewUint(0), true, []*types.IdRange{{Start: sdk.NewUint(1), End: sdk.NewUint(1)}})
// // 	suite.Require().Nil(err, "Error freezing address")

// // 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))

// // 	err = FreezeAddresses(suite, wctx, bob, sdk.NewUint(0), true, []*types.IdRange{{Start: sdk.NewUint(2), End: sdk.NewUint(2),}})
// // 	suite.Require().Nil(err, "Error freezing address")

// // 	badge, _ = GetCollection(suite, wctx, sdk.NewUint(1))
// // 	suite.Require().Equal(badge.FreezeRanges, []*types.IdRange{{Start: sdk.NewUint(0), End: sdk.NewUint(1)}})
// // }
