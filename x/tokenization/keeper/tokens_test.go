package keeper_test

import (
	"math"

	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HACK: Kinda forced the legacy code. Should clean up

func (suite *TestSuite) TestCreateTokens() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       GetOneUintRange(),
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
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating token: %s")

	// collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	// suite.Require().Nil(err, "Error getting token: %s")
	balance := &types.UserBalanceStore{}

	// totalSupplys, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Total")
	// suite.Require().Nil(err, "Error getting user balance: %s")
	// AssertBalancesEqual(suite, totalSupplys.Balances, []*types.Balance{
	// 	{
	// 		Amount:         sdkmath.NewUint(1),
	// 		TokenIds:       GetOneUintRange(),
	// 		OwnershipTimes: GetFullUintRanges(),
	// 	},
	// })

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertUintsEqual(suite, balance.Balances[0].Amount, sdkmath.NewUint(1))
	AssertUintRangesEqual(suite, balance.Balances[0].TokenIds, []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(1),
		},
	})

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTwoUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token")

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")

	_, err = types.SubtractBalance(suite.ctx, balance.Balances, &types.Balance{
		TokenIds: []*types.UintRange{
			GetOneUintRange()[0],
			GetTwoUintRanges()[0],
		},
		OwnershipTimes: GetFullUintRanges(),
		Amount:         sdkmath.NewUint(1),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	// totalSupplys, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Total")
	// suite.Require().Nil(err, "Error getting user balance: %s")
	// _, err = types.SubtractBalance(suite.ctx, totalSupplys.Balances, &types.Balance{
	// 	TokenIds: []*types.UintRange{
	// 		GetOneUintRange()[0],
	// 		GetTwoUintRanges()[0],
	// 	},
	// 	OwnershipTimes: GetFullUintRanges(),
	// 	Amount:         sdkmath.NewUint(1),
	// }, false)
	// suite.Require().Nil(err, "Error subtracting balances: %s")

	// unmintedSupplys, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Mint")
	// suite.Require().Nil(err, "Error getting user balance: %s")
	// AssertBalancesEqual(suite, unmintedSupplys.Balances, []*types.Balance{})

	// totalSupplys, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Total")
	// suite.Require().Nil(err, "Error getting user balance: %s")

	// _, err = types.SubtractBalance(suite.ctx, totalSupplys.Balances, &types.Balance{
	// 	TokenIds: []*types.UintRange{
	// 		GetTwoUintRanges()[0],
	// 	},
	// 	OwnershipTimes: GetFullUintRanges(),
	// 	Amount:         sdkmath.NewUint(2),
	// }, false)
	// suite.Require().Nil(err, "Error subtracting balances: %s")

	// unmintedSupplys, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Mint")
	// suite.Require().Nil(err, "Error getting user balance: %s")
	// _, err = types.SubtractBalance(suite.ctx, unmintedSupplys.Balances, &types.Balance{
	// 	TokenIds: []*types.UintRange{
	// 		GetTwoUintRanges()[0],
	// 	},
	// 	OwnershipTimes: GetFullUintRanges(),
	// 	Amount:         sdkmath.NewUint(2),
	// }, false)
	// suite.Require().Error(err, "Error subtracting balances: %s")

	// unmintedSupplys, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Mint")
	// suite.Require().Nil(err, "Error getting user balance: %s")
	// _, err = types.SubtractBalance(suite.ctx, unmintedSupplys.Balances, &types.Balance{
	// 	TokenIds: []*types.UintRange{
	// 		GetTwoUintRanges()[0],
	// 	},
	// 	OwnershipTimes: GetFullUintRanges(),
	// 	Amount:         sdkmath.NewUint(1),
	// }, false)
	// suite.Require().Nil(err, "Error subtracting balances: %s")

	// AssertUintsEqual(suite, collection.NextTokenId, sdkmath.NewUint(uint64(math.MaxUint64)).Add(sdkmath.NewUint(1)))

	// unmintedSupplys, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Mint")
	// suite.Require().Nil(err, "Error getting user balance: %s")
	// _, err = types.SubtractBalance(suite.ctx, unmintedSupplys.Balances, &types.Balance{
	// 	TokenIds: []*types.UintRange{
	// 		GetTopHalfUintRanges()[0],
	// 	},
	// 	OwnershipTimes: GetFullUintRanges(),
	// 	Amount:         types.NewUintFromString("1000000000000000000000000000000000000000000000000000000000000"),
	// }, false)
	// suite.Require().Nil(err, "Error subtracting balances: %s")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        alice,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTwoUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error creating tokens: %s")
}

func (suite *TestSuite) TestCreateTokensIdGreaterThanMax() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       GetOneUintRange(),
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
					TokenIds: []*types.UintRange{
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
	suite.Require().Error(err, "Error creating token: %s")
}

func (suite *TestSuite) TestDuplicateTokenIDs() {
	// wctx := sdk.WrapSDKContext(suite.ctx)

	currBalances := []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			TokenIds: []*types.UintRange{
				{
					Start: sdkmath.NewUint(2),
					End:   sdkmath.NewUint(1000),
				},
			},
			OwnershipTimes: GetFullUintRanges(),
		},
		{
			Amount:         sdkmath.NewUint(2),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	currBalances, err := types.SubtractBalance(suite.ctx, currBalances, &types.Balance{
		Amount:         sdkmath.NewUint(1),
		TokenIds:       GetOneUintRange(),
		OwnershipTimes: GetFullUintRanges(),
	}, false)
	suite.Require().Nil(err, "Error subtracting balances: %s")

	suite.Require().Equal(1, len(currBalances))
	suite.Require().Equal(sdkmath.NewUint(1), currBalances[0].Amount)
	suite.Require().Equal(sdkmath.NewUint(1), currBalances[0].TokenIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(1000), currBalances[0].TokenIds[0].End)
	suite.Require().Equal(1, len(currBalances[0].TokenIds))
}

func (suite *TestSuite) TestTokenIdsWeirdJSThing() {
	// wctx := sdk.WrapSDKContext(suite.ctx)

	currBalances := []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			TokenIds: []*types.UintRange{
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
		TokenIds: []*types.UintRange{
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
		TokenIds: []*types.UintRange{
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
	suite.Require().Equal(sdkmath.NewUint(3), currBalances[0].TokenIds[0].Start)
	suite.Require().Equal(sdkmath.NewUint(10000), currBalances[0].TokenIds[0].End)
	suite.Require().Equal(1, len(currBalances[0].TokenIds))
}

func (suite *TestSuite) TestDefaultsCannotBeDoubleUsedAfterSpent() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].Transfers = []*types.Transfer{}
	collectionsToCreate[0].DefaultBalances = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	collectionsToCreate[0].CollectionApprovals[0].FromListId = "AllWithoutMint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria = &types.ApprovalCriteria{
		OverridesToIncomingApprovals:   true,
		OverridesFromOutgoingApprovals: true,
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating token: %s")

	balance := &types.UserBalanceStore{}
	// totalSupplys, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Total")
	// suite.Require().Nil(err, "Error getting user balance: %s")
	// AssertBalancesEqual(suite, totalSupplys.Balances, []*types.Balance{
	// 	{
	// 		Amount:         sdkmath.NewUint(1),
	// 		TokenIds:       GetOneUintRange(),
	// 		OwnershipTimes: GetFullUintRanges(),
	// 	},
	// })

	balance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertUintsEqual(suite, balance.Balances[0].Amount, sdkmath.NewUint(1))
	AssertUintRangesEqual(suite, balance.Balances[0].TokenIds, []*types.UintRange{
		{
			Start: sdkmath.NewUint(1),
			End:   sdkmath.NewUint(1),
		},
	})

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token")

	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, bobBalance.Balances, []*types.Balance{})

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring token")
}

func (suite *TestSuite) TestValidUpdateTokenIdsWithPermission() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating token: %s")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting token: %s")
	AssertUintRangesEqual(suite, collection.ValidTokenIds, GetOneUintRange())

	//Set permission
	err = UpdateCollection(suite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:                     bob,
		CollectionId:                sdkmath.NewUint(1),
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
				{
					TokenIds:                  []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
					PermanentlyPermittedTimes: GetFullUintRanges(),
				},
				{
					TokenIds:                  GetFullUintRanges(),
					PermanentlyForbiddenTimes: GetFullUintRanges(),
				},
			},
		},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	//Update valid token IDs
	err = UpdateCollection(suite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		UpdateValidTokenIds: true,
		ValidTokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}},
	})
	suite.Require().Nil(err, "Error updating collection permissions")

	collection, err = GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting token: %s")
	AssertUintRangesEqual(suite, collection.ValidTokenIds, []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}})

	//Update valid token IDs - invalid > 2
	err = UpdateCollection(suite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		UpdateValidTokenIds: true,
		ValidTokenIds:       GetFullUintRanges(),
	})
	suite.Require().Error(err, "Error updating collection permissions")
}
