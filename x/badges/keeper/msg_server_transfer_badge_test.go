package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestTransferBadgeForceful() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   62,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

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
}

func (suite *TestSuite) TestApprovalsApproved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   62,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, aliceAccountNum, 0, []*types.Balance{
		{
			Balance:  1000000,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	bobbalance, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	// suite.Require().Equal(uint64(1000000-5000), bobbalance.Approvals[0].Amount)

	err = TransferBadge(suite, wctx, alice, 0, bobAccountNum, []*types.Transfers{
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
	suite.Require().Equal(uint64(1000000-5000), bobbalance.Approvals[0].Balances[0].Balance)
}

func (suite *TestSuite) TestApprovalsNotEnoughApproved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   62,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, aliceAccountNum, 0, []*types.Balance{
		{
			Balance:  10,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	err = TransferBadge(suite, wctx, charlie, 0, bobAccountNum, []*types.Transfers{
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
	suite.Require().EqualError(err, keeper.ErrApprovalForAddressDoesntExist.Error())
}

func (suite *TestSuite) TestApprovalsNotApprovedAtAll() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   62,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = TransferBadge(suite, wctx, charlie, 0, bobAccountNum, []*types.Transfers{
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
	suite.Require().EqualError(err, keeper.ErrApprovalForAddressDoesntExist.Error())
}

func (suite *TestSuite) TestApprovalsNotApprovedEnough() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   62,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, charlieAccountNum, 0, []*types.Balance{
		{
			Balance:  10,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	err = TransferBadge(suite, wctx, charlie, 0, bobAccountNum, []*types.Transfers{
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
	suite.Require().EqualError(err, keeper.ErrUnderflow.Error()) //underflow
}

func (suite *TestSuite) TestApprovalsApprovedJustEnough() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   62,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, charlieAccountNum, 0, []*types.Balance{
		{
			Balance:  10,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	err = TransferBadge(suite, wctx, charlie, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
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
	suite.Require().Nil(err, "Error transferring valid approved")
}

func (suite *TestSuite) TestApprovalOverflow() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   46,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = SetApproval(suite, wctx, bob, charlieAccountNum, 0, []*types.Balance{
		{
			Balance:  math.MaxUint64,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	err = TransferBadge(suite, wctx, charlie, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
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
	suite.Require().Nil(err, "Error transferring valid approved")

	err = SetApproval(suite, wctx, bob, charlieAccountNum, 0, []*types.Balance{
		{
			Balance:  math.MaxUint64,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	})
	suite.Require().Nil(err, "Error setting approval")

	// err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, []uint64{0})
	// suite.Require().Nil(err, "Error setting approval")
}

func (suite *TestSuite) TestTransferUnderflowNotEnoughBalance() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   46,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: math.MaxUint64,
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
	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
}

func (suite *TestSuite) TestPendingTransferUnderflowNotEnoughBalance() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   62,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)
	suite.Require().Nil(err)

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: math.MaxUint64,
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
	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
}

func (suite *TestSuite) TestTransferInvalidBadgeIdRanges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   46,
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
	suite.Require().Equal(uint64(10000), fetchedBalance[0].Balance)

	err = TransferBadge(suite, wctx, charlie, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 10,
							End:   1,
						},
					},
				},
			},
		},
	})
	suite.Require().EqualError(err, keeper.ErrInvalidBadgeRange.Error())

	err = TransferBadge(suite, wctx, charlie, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 0,
							End:   math.MaxUint64,
						},
					},
				},
			},
		},
	})
	suite.Require().EqualError(err, keeper.ErrBadgeNotExists.Error())
}

func (suite *TestSuite) TestTransferBadgeNeedToMergeWithNextAndPrev() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   46,
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
			Amount: 10000,
		},
	})
	suite.Require().Nil(err, "Error creating badges")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 0,
							End:   499,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 501,
							End:   1000,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 500,
							End:   500,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")
}

func (suite *TestSuite) TestTransferBadgeNeedToMergeWithJustNext() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   46,
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
			Amount: 10000,
		},
	})
	suite.Require().Nil(err, "Error creating badges")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 501,
							End:   1000,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 500,
							End:   500,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")
}

func (suite *TestSuite) TestTransferBadgeBinarySearchInsertIdx() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   46,
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
			Amount: 10000,
		},
	})
	suite.Require().Nil(err, "Error creating badges")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 0,
							End:   100,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 200,
							End:   300,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 400,
							End:   500,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 600,
							End:   700,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 800,
							End:   900,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 1000,
							End:   1100,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 150,
							End: 150,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{aliceAccountNum},
			Balances: []*types.Balance{
				{
					Balance: 10,
					BadgeIds: []*types.IdRange{
						{
							Start: 950,
							End:  950,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")
}
