package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/math"
)

func (suite *TestSuite) TestFreezeAddressesDirectlyWhenCreatingNewBadge() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   62,

				DisallowedTransfers: []*types.TransferMapping{
					{
						From: &types.Addresses{
							AccountNums: []*types.IdRange{
								{Start: aliceAccountNum},
							},
							ManagerOptions: types.ManagerOptions_Neutral,
						},
						To: &types.Addresses{
							AccountNums: []*types.IdRange{
								{Start: 0, End: math.MaxUint64},
							},
							ManagerOptions: types.ManagerOptions_Neutral,
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	// badge, _ := GetCollection(suite, wctx, 0)

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 5000,
					BadgeIds: []*types.IdRange{
						{
							Start: 0,
							End:   0,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, alice, 0, aliceAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{bobAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 5000,
					BadgeIds: []*types.IdRange{
						{
							Start: 0,
							End:   0,
						},
					},
				},
			},
		},
	})
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

func (suite *TestSuite) TestTransferBadgeForcefulUnfrozenByDefault() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   23,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 0)

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 0)

	bobbalance, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.MaxSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobbalance.Balances)[0].Balance)

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 5000,
					BadgeIds: []*types.IdRange{
						{
							Start: 0,
							End:   0,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge")

	bobbalance, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobbalance.Balances)[0].Balance)

	alicebalance, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)

	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, alicebalance.Balances)[0].Balance)

	err = UpdateDisallowedTransfers(suite, wctx, bob, 0, []*types.TransferMapping{
		{
			From: &types.Addresses{
				AccountNums:    []*types.IdRange{{Start: aliceAccountNum, End: aliceAccountNum}},
				ManagerOptions: types.ManagerOptions_Neutral,
			},
			To: &types.Addresses{
				AccountNums:    []*types.IdRange{{Start: 0, End: math.MaxUint64}},
				ManagerOptions: types.ManagerOptions_Neutral,
			},
		},
	})
	suite.Require().Nil(err, "Error freezing address")

	badge, _ = GetCollection(suite, wctx, 0)

	err = TransferBadge(suite, wctx, alice, 0, aliceAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{bobAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 5000,
					BadgeIds: []*types.IdRange{
						{
							Start: 0,
							End:   0,
						},
					},
				},
			},
		},
	})
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

//TODO:

// func (suite *TestSuite) TestTransferBadgeForcefulFrozenByDefault() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Badge: types.MsgNewCollection{
// 				CollectionUri: "https://example.com",
// BadgeUri: "https://example.com/{id}",
// 				Permissions: 63,
// 			},
// 			Amount:  1,
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, 0)

// 	//Create badge 1 with supply > 1
// 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: 10000,
// 			Amount: 1,
// 		},
// 	},)
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, 0)

// 	bobbalance, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

// 	suite.Require().Equal(uint64(1), badge.NextBadgeId)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
// 			Balance:  10000,
// 		},
// 	}, badge.MaxSupplys)
// 	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobbalance.Balances)[0].Balance)

// 	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
// 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

// 	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
// 	suite.Require().Nil(err, "Error unfreezing address")

// 	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
// 	suite.Require().Nil(err, "Error transferring after unfreeze")

// 	err = TransferBadge(suite, wctx, alice, aliceAccountNum, []uint64{bobAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
// 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
// }

// func (suite *TestSuite) TestTransferBadgeForcefulFrozenByDefaultAddAndRemove() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Badge: types.MsgNewCollection{
// 				CollectionUri: "https://example.com",
// BadgeUri: "https://example.com/{id}",
// 				Permissions: 63,
// 			},
// 			Amount:  1,
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, 0)

// 	//Create badge 1 with supply > 1
// 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: 10000,
// 			Amount: 1,
// 		},
// 	},)
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, 0)

// 	bobbalance, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

// 	suite.Require().Equal(uint64(1), badge.NextBadgeId)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
// 			Balance:  10000,
// 		},
// 	}, badge.MaxSupplys)
// 	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobbalance.Balances)[0].Balance)

// 	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
// 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

// 	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
// 	suite.Require().Nil(err, "Error unfreezing address")

// 	err = FreezeAddresses(suite, wctx, bob, 0, false, []*types.IdRange{{Start: bobAccountNum, End: 0}})
// 	suite.Require().Nil(err, "Error unfreezing address")

// 	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
// 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

// 	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
// 	suite.Require().Nil(err, "Error unfreezing address")

// 	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
// 	suite.Require().Nil(err, "Error transferring after unfreeze")

// 	err = TransferBadge(suite, wctx, alice, aliceAccountNum, []uint64{bobAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
// 	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
// }

// func (suite *TestSuite) TestFreezeCantFreeze() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Badge: types.MsgNewCollection{
// 				CollectionUri: "https://example.com",
// 				BadgeUri: "https://example.com/{id}",
// 				Permissions: 0,
// 			},
// 			Amount:  1,
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)

// 	err := FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
// 	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
// }

// func (suite *TestSuite) TestTransferBadgeForcefulUnfrozenByDefaultOmitEmptyCase() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := []CollectionsToCreate{
// 		{
// 			Badge: types.MsgNewCollection{
// 				CollectionUri: "https://example.com",
// BadgeUri: "https://example.com/{id}",
// 				Permissions: 62,
// 			},
// 			Amount:  1,
// 			Creator: bob,
// 		},
// 	}

// 	CreateCollections(suite, wctx, collectionsToCreate)
// 	badge, _ := GetCollection(suite, wctx, 0)

// 	//Create badge 1 with supply > 1
// 	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
// 		{
// 			Supply: 10000,
// 			Amount: 1,
// 		},
// 	},)
// 	suite.Require().Nil(err, "Error creating badge")
// 	badge, _ = GetCollection(suite, wctx, 0)

// 	bobbalance, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

// 	suite.Require().Equal(uint64(1), badge.NextBadgeId)
// 	suite.Require().Equal([]*types.Balance{
// 		{
// 			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
// 			Balance:  10000,
// 		},
// 	}, badge.MaxSupplys)
// 	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobbalance.Balances)[0].Balance)

// 	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: 0, End: 0}})
// 	suite.Require().Nil(err, "Error freezing address")

// 	badge, _ = GetCollection(suite, wctx, 0)

// 	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: 1, End: 0}})
// 	suite.Require().Nil(err, "Error freezing address")

// 	badge, _ = GetCollection(suite, wctx, 0)
// 	suite.Require().Equal(badge.FreezeRanges, []*types.IdRange{{Start: 0, End: 1}})
// }
