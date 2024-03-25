package keeper_test

import (
	"math"

	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestCreateBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			BadgeIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	collectionsToCreate[0].Transfers = []*types.Transfer{
		{
			From:        "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting badge: %s")
	balance := &types.UserBalanceStore{}

	totalSupplys, err := GetUserBalance(suite, wctx, sdk.NewUint(1), "Total")
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, totalSupplys.Balances, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			BadgeIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	})

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertUintsEqual(suite, balance.Balances[0].Amount, sdkmath.NewUint(1))
	AssertUintRangesEqual(suite, balance.Balances[0].BadgeIds, []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(1),
		},
	})

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			BadgeIds:       GetTwoUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	})
	suite.Require().Nil(err, "Error creating badges: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetTwoUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge")

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	_, err = types.SubtractBalance(suite.ctx, balance.Balances, &types.Balance{
		BadgeIds: []*types.UintRange{
			GetOneUintRange()[0],
			GetTwoUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
		Amount:         sdkmath.NewUint(1),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	totalSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Total")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(suite.ctx, totalSupplys.Balances, &types.Balance{
		BadgeIds: []*types.UintRange{
			GetOneUintRange()[0],
			GetTwoUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
		Amount:         sdkmath.NewUint(1),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	unmintedSupplys, err := GetUserBalance(suite, wctx, sdk.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, unmintedSupplys.Balances, []*types.Balance{})

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			BadgeIds:       GetTwoUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	})
	suite.Require().Nil(err, "Error creating badges: %s")

	totalSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Total")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(suite.ctx, totalSupplys.Balances, &types.Balance{
		BadgeIds: []*types.UintRange{
			GetTwoUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
		Amount:         sdkmath.NewUint(2),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	unmintedSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(suite.ctx, unmintedSupplys.Balances, &types.Balance{
		BadgeIds: []*types.UintRange{
			GetTwoUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
		Amount:         sdkmath.NewUint(2),
	}, false)
	suite.Require().Error(err, "Error subtracting balances: %s")

	unmintedSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(suite.ctx, unmintedSupplys.Balances, &types.Balance{
		BadgeIds: []*types.UintRange{
			GetTwoUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
		Amount:         sdkmath.NewUint(1),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	_, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount:         types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
			BadgeIds:       GetTopHalfUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	})
	suite.Require().Error(err, "Error creating badges: %s")

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount:         types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
			BadgeIds:       GetBottomHalfUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
		{
			Amount:         types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
			BadgeIds:       GetTopHalfUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	})
	suite.Require().Nil(err, "Error creating badges: %s")
	// AssertUintsEqual(suite, collection.NextBadgeId, sdkmath.NewUint(uint64(math.MaxUint64)).Add(sdkmath.NewUint(1)))

	unmintedSupplys, err = GetUserBalance(suite, wctx, sdk.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	_, err = types.SubtractBalance(suite.ctx, unmintedSupplys.Balances, &types.Balance{
		BadgeIds: []*types.UintRange{
			GetTopHalfUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
		Amount:         types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	collection, err = suite.app.BadgesKeeper.CreateBadges(suite.ctx, collection, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			BadgeIds:       GetTwoUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
		},
	})
	suite.Require().Nil(err, "Error creating badges: %s")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        alice,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetTwoUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating badges: %s")
}

func (suite *TestSuite) TestCreateBadgesIdGreaterThanMax() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].BadgesToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			BadgeIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	collectionsToCreate[0].Transfers = []*types.Transfer{
		{
			From:        "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: []*types.UintRange{
						{
							Start: sdkmath.NewUint(1),
							End:   sdkmath.NewUint(math.MaxUint64).Add(sdkmath.NewUint(1)),
						},
					},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "Error creating badge: %s")
}

func (suite *TestSuite) TestDuplicateBadgeIDs() {
	// wctx := sdk.WrapSDKContext(suite.ctx)

	currBalances := []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(2),
					End:   sdkmath.NewUint(1000),
				},
			},
			OwnershipTimes: GetFullUintRanges(),
		},
		{
			Amount:         sdkmath.NewUint(2),
			BadgeIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	currBalances, err := types.SubtractBalance(suite.ctx, currBalances, &types.Balance{
		Amount:         sdkmath.NewUint(1),
		BadgeIds:       GetOneUintRange(),
		OwnershipTimes: GetFullUintRanges(),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	suite.Require().Equal(1, len(currBalances))
	suite.Require().Equal(sdkmath.NewUint(1), currBalances[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(1), currBalances[0].BadgeIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(1000), currBalances[0].BadgeIds[0].End)
	suite.Require().Equal(1, len(currBalances[0].BadgeIds))
}

func (suite *TestSuite) TestBadgeIdsWeirdJSThing() {
	// wctx := sdk.WrapSDKContext(suite.ctx)

	currBalances := []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			BadgeIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(10000),
				},
			},
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	currBalances, err := types.SubtractBalance(suite.ctx, currBalances, &types.Balance{
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(2),
				End:   sdkmath.NewUint(2),
			},
		},
		OwnershipTimes: GetFullUintRanges(),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	currBalances, err = types.SubtractBalance(suite.ctx, currBalances, &types.Balance{
		Amount: sdkmath.NewUint(1),
		BadgeIds: []*types.UintRange{
			{
				Start: sdkmath.NewUint(1),
				End:   sdkmath.NewUint(1),
			},
		},
		OwnershipTimes: GetFullUintRanges(),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	suite.Require().Equal(1, len(currBalances))
	suite.Require().Equal(sdkmath.NewUint(1), currBalances[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(3), currBalances[0].BadgeIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(10000), currBalances[0].BadgeIds[0].End)
	suite.Require().Equal(1, len(currBalances[0].BadgeIds))
}
