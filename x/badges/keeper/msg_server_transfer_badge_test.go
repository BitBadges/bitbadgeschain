package keeper_test

// import (
// sdkmath "cosmossdk.io/math"
// 	"math"

// 	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
// 	"github.com/bitbadges/bitbadgeschain/x/badges/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// func (suite *TestSuite) TestTransferBadgeForceful() {
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
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

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
// 	fetchedBalance, err = keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(5000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	alicebalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
// 	fetchedBalance, err = keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, alicebalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(5000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)
// }

// func (suite *TestSuite) TestApprovalsApproved() {
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
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	err = SetApproval(suite, wctx, bob, alice, sdkmath.NewUint(1), []*types.Balance{
// 		{
// 			Amount:   sdkmath.NewUint(1000000),
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error setting approval")

// 	bobbalance, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
// 	// suite.Require().Equal(sdkmath.NewUint(1000000-5000), bobbalance.Approvals[0].Amount)

// 	err = TransferBadge(suite, wctx, alice, sdkmath.NewUint(1), bob, []*types.Transfer{
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
// 	suite.Require().Equal(sdkmath.NewUint(1000000-5000), bobbalance.Approvals[0].Balances[0].Amount)
// }

// func (suite *TestSuite) TestApprovalsNotEnoughApproved() {
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
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	err = SetApproval(suite, wctx, bob, alice, sdkmath.NewUint(1), []*types.Balance{
// 		{
// 			Amount:   sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error setting approval")

// 	err = TransferBadge(suite, wctx, charlie, sdkmath.NewUint(1), bob, []*types.Transfer{
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
// 	suite.Require().EqualError(err, keeper.ErrApprovalForAddressDoesntExist.Error())
// }

// func (suite *TestSuite) TestApprovalsNotApprovedAtAll() {
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
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	err = TransferBadge(suite, wctx, charlie, sdkmath.NewUint(1), bob, []*types.Transfer{
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
// 	suite.Require().EqualError(err, keeper.ErrApprovalForAddressDoesntExist.Error())
// }

// func (suite *TestSuite) TestApprovalsNotApprovedEnough() {
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
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	err = SetApproval(suite, wctx, bob, charlie, sdkmath.NewUint(1), []*types.Balance{
// 		{
// 			Amount:   sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error setting approval")

// 	err = TransferBadge(suite, wctx, charlie, sdkmath.NewUint(1), bob, []*types.Transfer{
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
// 	suite.Require().EqualError(err, keeper.ErrUnderflow.Error()) //underflow
// }

// func (suite *TestSuite) TestApprovalsApprovedJustEnough() {
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
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	err = SetApproval(suite, wctx, bob, charlie, sdkmath.NewUint(1), []*types.Balance{
// 		{
// 			Amount:   sdkmath.NewUint(10),
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error setting approval")

// 	err = TransferBadge(suite, wctx, charlie, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
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
// 	suite.Require().Nil(err, "Error transferring valid approved")
// }

// func (suite *TestSuite) TestApprovalOverflow() {
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
// 				Permissions: sdkmath.NewUint(46),
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	err = SetApproval(suite, wctx, bob, charlie, sdkmath.NewUint(1), []*types.Balance{
// 		{
// 			Amount:   sdkmath.NewUint(math.MaxUint64),
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error setting approval")

// 	err = TransferBadge(suite, wctx, charlie, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
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
// 	suite.Require().Nil(err, "Error transferring valid approved")

// 	err = SetApproval(suite, wctx, bob, charlie, sdkmath.NewUint(1), []*types.Balance{
// 		{
// 			Amount:   sdkmath.NewUint(math.MaxUint64),
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error setting approval")

// 	// err = HandlePendingTransfers(suite, wctx, bob, sdkmath.NewUint(0), []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, []uint64{0})
// 	// suite.Require().Nil(err, "Error setting approval")
// }

// func (suite *TestSuite) TestTransferUnderflowNotEnoughBalance() {
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
// 				Permissions: sdkmath.NewUint(46),
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 			BadgeIds: []*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 			Amount:   sdkmath.NewUint(10000),
// 		},
// 	}, badge.MaxSupplys)
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(math.MaxUint64),
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
// 	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
// }

// func (suite *TestSuite) TestPendingTransferUnderflowNotEnoughBalance() {
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
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)
// 	suite.Require().Nil(err)

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(math.MaxUint64),
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
// 	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
// }

// func (suite *TestSuite) TestTransferInvalidBadgeIdRanges() {
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
// 				Permissions: sdkmath.NewUint(46),
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 	fetchedBalance, err := keeper.GetBalancesForIds([]*types.IdRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}}, bobbalance.Balances)
// 	suite.Require().Equal(sdkmath.NewUint(10000), fetchedBalance[0].Amount)

// 	err = TransferBadge(suite, wctx, charlie, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(10),
// 							End:   sdkmath.NewUint(1),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().EqualError(err, keeper.ErrInvalidBadgeRange.Error())

// 	err = TransferBadge(suite, wctx, charlie, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(0),
// 							End:   sdkmath.NewUint(math.MaxUint64),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().EqualError(err, keeper.ErrBadgeNotExists.Error())
// }

// func (suite *TestSuite) TestTransferBadgeNeedToMergeWithNextAndPrev() {
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
// 				Permissions: sdkmath.NewUint(46),
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 			Amount: sdkmath.NewUint(10000),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badges")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(1),
// 							End:   sdkmath.NewUint(500),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(501),
// 							End:   sdkmath.NewUint(1000),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(500),
// 							End:   sdkmath.NewUint(500),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")
// }

// func (suite *TestSuite) TestTransferBadgeNeedToMergeWithJustNext() {
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
// 				Permissions: sdkmath.NewUint(46),
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 			Amount: sdkmath.NewUint(10000),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badges")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(501),
// 							End:   sdkmath.NewUint(1000),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(500),
// 							End:   sdkmath.NewUint(500),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")
// }

// func (suite *TestSuite) TestTransferBadgeBinarySearchInsertIdx() {
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
// 				Permissions: sdkmath.NewUint(46),
// 				ApprovedTransfers: []*types.CollectionApprovedTransfer{
// 					{
// 						From: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
// 							ManagerOptions:       sdkmath.NewUint(0),
// 						},
// 						To: &types.AddressMapping{
// 							Addresses:            []string{},
// 							IncludeOnlySpecified: false,
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
// 			Amount: sdkmath.NewUint(10000),
// 		},
// 	})
// 	suite.Require().Nil(err, "Error creating badges")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(1),
// 							End:   sdkmath.NewUint(100),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(200),
// 							End:   sdkmath.NewUint(300),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(400),
// 							End:   sdkmath.NewUint(500),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(600),
// 							End:   sdkmath.NewUint(700),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(800),
// 							End:   sdkmath.NewUint(8900),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(1000),
// 							End:   sdkmath.NewUint(1100),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(150),
// 							End:   sdkmath.NewUint(150),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")

// 	err = TransferBadge(suite, wctx, bob, sdkmath.NewUint(1), bob, []*types.Transfer{
// 		{
// 			ToAddresses: []string{alice},
// 			Balances: []*types.Balance{
// 				{
// 					Amount: sdkmath.NewUint(10),
// 					BadgeIds: []*types.IdRange{
// 						{
// 							Start: sdkmath.NewUint(950),
// 							End:   sdkmath.NewUint(950),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	})
// 	suite.Require().Nil(err, "Error transfering badge")
// }
