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
								{Start: aliceAccountNum, End: aliceAccountNum},
							},
							Options: types.AddressOptions_None,
						},
						To: &types.Addresses{
							AccountNums: []*types.IdRange{
								{Start: 0, End: math.MaxUint64},
							},
							Options: types.AddressOptions_None,
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
	fetchedBalance, err := keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0, End: 0}}, bobbalance.Balances)
	suite.Require().Nil(err)
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)

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
	fetchedBalance, err = keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0, End: 0}}, bobbalance.Balances)
	suite.Require().Equal(uint64(5000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	alicebalance, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)

	fetchedBalance, err = keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0, End: 0}}, alicebalance.Balances)
	suite.Require().Equal(uint64(5000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = UpdateDisallowedTransfers(suite, wctx, bob, 0, []*types.TransferMapping{
		{
			From: &types.Addresses{
				AccountNums:    []*types.IdRange{{Start: aliceAccountNum, End: aliceAccountNum}},
				Options: types.AddressOptions_None,
			},
			To: &types.Addresses{
				AccountNums:    []*types.IdRange{{Start: 0, End: math.MaxUint64}},
				Options: types.AddressOptions_None,
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
//TODO: also test transfer mappings

//TODO: also test manager approved transfers with transfer mappings

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
// 	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0, End: 0}}, bobbalance.Balances)[0].Balance)

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
// 	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0, End: 0}}, bobbalance.Balances)[0].Balance)

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
// 	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0, End: 0}}, bobbalance.Balances)[0].Balance)

// 	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: 0, End: 0}})
// 	suite.Require().Nil(err, "Error freezing address")

// 	badge, _ = GetCollection(suite, wctx, 0)

// 	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: 1, End: 1}})
// 	suite.Require().Nil(err, "Error freezing address")

// 	badge, _ = GetCollection(suite, wctx, 0)
// 	suite.Require().Equal(badge.FreezeRanges, []*types.IdRange{{Start: 0, End: 1}})
// }
