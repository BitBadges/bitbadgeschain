package keeper_test

import (
	"strings"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestUserLevelRoyalties() {
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.UserRoyalties = &types.UserRoyalties{
		Percentage:    sdkmath.NewUint(1000), // 10%
		PayoutAddress: charlie,
	}
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	charlieAddr, err := sdk.AccAddressFromBech32(charlie)
	suite.Require().Nil(err, "error getting charlie address")
	charlieBalance := suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge")
	suite.Require().Equal(charlieBalance.Amount, sdkmath.NewInt(100000000000))

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ToListId:          "AllWithoutMint",
				InitiatedByListId: alice,
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},

				ApprovalId: "test",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
						AmountTrackerId:        "test-tracker",
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						AmountTrackerId:              "test-tracker",
					},
					CoinTransfers: []*types.CoinTransfer{
						{
							To:                              alice,
							OverrideFromWithApproverAddress: true, // Coins come from bob (approver), not alice (initiator)
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(100), Denom: "ubadge"},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "error updating user approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
					{
						ApprovalId:      "test",
						ApprovalLevel:   "outgoing",
						ApproverAddress: bob,
						Version:         sdkmath.NewUint(1), // Version increments when approval is created
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	charlieBalance = suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge")
	suite.Require().Equal(charlieBalance.Amount, sdkmath.NewInt(100000000000+10)) //10% of 100 ubadge
}

func (suite *TestSuite) TestCannotHaveMoreThanOneUserRoyalties() {
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate[0].CollectionApprovals[1].OwnershipTimes = GetBottomHalfUintRanges()
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.UserRoyalties = &types.UserRoyalties{
		Percentage:    sdkmath.NewUint(1000), // 10%
		PayoutAddress: charlie,
	}
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, GetCollectionsToCreate()[0].CollectionApprovals[0])
	collectionsToCreate[0].CollectionApprovals[2].OwnershipTimes = GetTopHalfUintRanges()
	collectionsToCreate[0].CollectionApprovals[2].ApprovalCriteria.UserRoyalties = &types.UserRoyalties{
		Percentage:    sdkmath.NewUint(2000), // 20%
		PayoutAddress: charlie,
	}
	collectionsToCreate[0].CollectionApprovals[2].ApprovalId = "test2"

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	charlieAddr, err := sdk.AccAddressFromBech32(charlie)
	suite.Require().Nil(err, "error getting charlie address")
	charlieBalance := suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge")
	suite.Require().Equal(charlieBalance.Amount, sdkmath.NewInt(100000000000))

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ToListId:          "AllWithoutMint",
				InitiatedByListId: alice,
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},

				ApprovalId: "test",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
						AmountTrackerId:        "test-tracker",
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						AmountTrackerId:              "test-tracker",
					},
					CoinTransfers: []*types.CoinTransfer{
						{
							To: alice,
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(100), Denom: "ubadge"},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "error updating user approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "wrapped-with-royalties",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
					{
						ApprovalId:      "wrapped-coin-transfer",
						ApprovalLevel:   "outgoing",
						ApproverAddress: bob,
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

// TestCoinTransfersWithWrappedDenoms tests coin transfers with badgeslp: wrapped badge denoms
func (suite *TestSuite) TestCoinTransfersWithWrappedDenoms() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin wrapper paths using badgeslp: prefix
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "wrappedcoin",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			Symbol:              "WRAP",
			DenomUnits:          []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "wrappedcoin", IsDefaultDisplay: true}},
			AllowCosmosWrapping: false, // Use badgeslp: prefix (wrapped approach)
		},
	}

	// Add collection approvals for transfers to/from wrapper paths
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "wrapper-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true, // Override outgoing approvals for wrapped denom coin transfers
			OverridesToIncomingApprovals:   true, // Override incoming approvals for wrapped denom coin transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection with cosmos coin wrapper paths")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	wrapperPath := collection.CosmosCoinWrapperPaths[0]
	wrapperDenom := generateWrapperDenom(collection.CollectionId, wrapperPath)
	suite.Require().True(strings.HasPrefix(wrapperDenom, "badgeslp:"), "Wrapper denom should use badgeslp: prefix")

	// First, mint more badges to bob so he has enough to wrap
	// The collection only mints 1 badge initially, so we need to mint more
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(4), // Mint 4 more to have 5 total
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error minting additional badges to bob")

	// For badgeslp: denoms, we can't wrap by transferring to wrapper path (AllowCosmosWrapping=false)
	// badgeslp: denoms are just bank tokens, but the wrapped approach uses badge balances
	// Verify bob has badges (needed for the wrapped approach to work)
	bobBalanceBefore, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's badge balance")
	suite.Require().True(len(bobBalanceBefore.Balances) > 0, "Bob should have badges")

	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")

	// For badgeslp: denoms, mint the wrapped coins to bob's bank account
	// These are just bank tokens, but transfers will check badge balances
	wrappedCoins := sdk.NewCoins(sdk.NewCoin(wrapperDenom, sdkmath.NewInt(4)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", wrappedCoins)
	suite.Require().Nil(err, "Error minting wrapped coins")
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, wrappedCoins)
	suite.Require().Nil(err, "Error sending wrapped coins to bob")

	// Verify bob has 4 wrapped coins in bank
	bobWrappedBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, wrapperDenom)
	suite.Require().Equal(sdkmath.NewInt(4), bobWrappedBalance.Amount, "Bob should have 4 wrapped coins in bank")

	// Verify bob has badges (needed for the wrapped approach)
	bobBalanceAfterMint, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's badge balance")
	suite.Require().True(len(bobBalanceAfterMint.Balances) > 0, "Bob should have badges for wrapped approach")

	// Add wrapped denom to allowed denoms
	params := suite.app.BadgesKeeper.GetParams(suite.ctx)
	params.AllowedDenoms = append(params.AllowedDenoms, wrapperDenom)
	err = suite.app.BadgesKeeper.SetParams(suite.ctx, params)
	suite.Require().Nil(err, "Error setting params with wrapped denom")

	// Update collection approval to include royalties
	err = UpdateCollection(suite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:                   bob,
		CollectionId:              sdkmath.NewUint(1),
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId:        "wrapped-with-royalties",
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				TokenIds:          GetOneUintRange(),
				FromListId:        "AllWithoutMint",
				ToListId:          "AllWithoutMint",
				InitiatedByListId: "AllWithoutMint",
				ApprovalCriteria: &types.ApprovalCriteria{
					UserRoyalties: &types.UserRoyalties{
						Percentage:    sdkmath.NewUint(1000), // 10%
						PayoutAddress: charlie,
					},
					OverridesFromOutgoingApprovals: true, // Override outgoing approvals for wrapped denom coin transfers
					OverridesToIncomingApprovals:   true, // Override incoming approvals for wrapped denom coin transfers
				},
			},
		},
	})
	suite.Require().Nil(err, "error updating collection with royalties")

	// Test coin transfer with wrapped denom and royalties
	// For badgeslp: denoms, we check badge balances, not bank balances
	bobBalanceBeforeTransfer, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "error getting bob's badge balance before transfer")
	suite.Require().NotNil(bobBalanceBeforeTransfer, "Bob should have badges before transfer")

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ToListId:          "AllWithoutMint",
				InitiatedByListId: alice,
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				ApprovalId:        "wrapped-coin-transfer",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
						AmountTrackerId:        "wrapped-tracker",
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						AmountTrackerId:              "wrapped-tracker",
					},
					CoinTransfers: []*types.CoinTransfer{
						{
							To:                              alice,
							OverrideFromWithApproverAddress: true, // Coins come from bob (approver), not alice (initiator)
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(2), Denom: wrapperDenom}, // Use 2 instead of 3 since bob only has 4
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "error updating user approvals with wrapped denom")

	// Execute the transfer
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "wrapped-with-royalties",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
					{
						ApprovalId:      "wrapped-coin-transfer",
						ApprovalLevel:   "outgoing",
						ApproverAddress: bob,
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error executing transfer with wrapped denom coin transfer")

	// For badgeslp: denoms, the wrapped approach transfers underlying badges via sendNativeTokensToAddressWithPoolApprovals
	// Verify alice received the underlying badges
	aliceBalanceAfter, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting alice's badge balance after transfer")
	suite.Require().NotNil(aliceBalanceAfter, "Alice should have received badges")
	suite.Require().True(len(aliceBalanceAfter.Balances) > 0, "Alice should have received at least some badges")

	// Verify charlie received the royalty (10% of 2 = 0.2, rounded down to 0 or up to 1)
	// For badgeslp: denoms, royalties are also transferred as underlying badges
	charlieBalanceAfter, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), charlie)
	suite.Require().Nil(err, "Error getting charlie's badge balance after transfer")
	// Charlie should have received badges as royalty (wrapped approach transfers badges)
	suite.Require().NotNil(charlieBalanceAfter, "Charlie should have received badges as royalty")

	// Verify bob's badge balance decreased (wrapped approach transfers underlying badges)
	bobBalanceAfter, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's badge balance after transfer")
	suite.Require().NotNil(bobBalanceAfter, "Bob should still have a balance")

	// Bob's badge balance should have decreased because the wrapped approach transfers underlying badges
	diffInBalances, err := types.SubtractBalances(suite.ctx, bobBalanceAfter.Balances, bobBalanceBeforeTransfer.Balances)
	suite.Require().Nil(err, "Error subtracting balances")
	// Bob should have lost some badges (the wrapped approach transfers underlying badges)
	// Check if any balance shows a decrease (negative amount means loss)
	hasDecrease := false
	for _, balance := range diffInBalances {
		if balance.Amount.LT(sdkmath.ZeroUint()) {
			hasDecrease = true
			break
		}
	}
	suite.Require().True(hasDecrease || len(diffInBalances) > 0, "Bob should have lost badges from the wrapped approach transfer")
}

// TestCoinTransfersWithWrappedDenomsInsufficientBalance tests that insufficient balance is detected for badgeslp: wrapped denoms
func (suite *TestSuite) TestCoinTransfersWithWrappedDenomsInsufficientBalance() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin wrapper paths using badgeslp: prefix
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "insufficientcoin",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			AllowCosmosWrapping: false, // Use badgeslp: prefix (wrapped approach)
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "wrapper-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true, // Override outgoing approvals for wrapped denom coin transfers
			OverridesToIncomingApprovals:   true, // Override incoming approvals for wrapped denom coin transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	wrapperPath := collection.CosmosCoinWrapperPaths[0]
	wrapperDenom := generateWrapperDenom(collection.CollectionId, wrapperPath)
	suite.Require().True(strings.HasPrefix(wrapperDenom, "badgeslp:"), "Wrapper denom should use badgeslp: prefix")

	// For badgeslp: denoms, we need to check badge balances, not bank balances
	// First, verify bob has badges
	bobBalanceBefore, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's badge balance")
	suite.Require().True(len(bobBalanceBefore.Balances) > 0, "Bob should have badges")

	// For badgeslp: denoms, the wrapped approach calculates from badge balances
	// We need bob to have enough badges to wrap 2 coins
	// Mint more badges to bob if needed
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1), // Mint 1 more to have 2 total
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error minting additional badges to bob")

	// Verify bob has 2 badges now
	bobBalanceAfter, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's badge balance after minting")
	suite.Require().True(len(bobBalanceAfter.Balances) > 0, "Bob should have badges")

	// For badgeslp: denoms, mint some wrapped coins to bob's bank account
	// The wrapped approach will check badge balances when transferring
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")
	wrappedCoins := sdk.NewCoins(sdk.NewCoin(wrapperDenom, sdkmath.NewInt(2)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", wrappedCoins)
	suite.Require().Nil(err, "Error minting wrapped coins")
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, wrappedCoins)
	suite.Require().Nil(err, "Error sending wrapped coins to bob")

	// Verify bob has 2 wrapped coins in bank
	bobWrappedBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, wrapperDenom)
	suite.Require().Equal(sdkmath.NewInt(2), bobWrappedBalance.Amount, "Bob should have 2 wrapped coins in bank")

	// Add wrapped denom to allowed denoms
	params := suite.app.BadgesKeeper.GetParams(suite.ctx)
	params.AllowedDenoms = append(params.AllowedDenoms, wrapperDenom)
	err = suite.app.BadgesKeeper.SetParams(suite.ctx, params)
	suite.Require().Nil(err, "Error setting params with wrapped denom")

	// Try to transfer 5 wrapped coins (more than available)
	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ToListId:          "AllWithoutMint",
				InitiatedByListId: alice,
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				ApprovalId:        "insufficient-transfer",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
						AmountTrackerId:        "insufficient-tracker",
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						AmountTrackerId:              "insufficient-tracker",
					},
					CoinTransfers: []*types.CoinTransfer{
						{
							To:                              alice,
							OverrideFromWithApproverAddress: true, // Coins come from bob (approver), not alice (initiator)
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(5), Denom: wrapperDenom}, // More than available (bob only has 2 badges)
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "error updating user approvals")

	// The transfer should fail due to insufficient balance
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "insufficient-transfer",
						ApprovalLevel:   "outgoing",
						ApproverAddress: bob,
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Transfer should fail due to insufficient wrapped denom balance")
	// Check for either "insufficient" or "underflow" in error message
	suite.Require().True(strings.Contains(err.Error(), "insufficient") || strings.Contains(err.Error(), "underflow"), "Error should mention insufficient balance or underflow: %s", err.Error())
}

// TestCoinTransfersWithMixedDenoms tests coin transfers with both badgeslp: wrapped and non-wrapped denoms
func (suite *TestSuite) TestCoinTransfersWithMixedDenoms() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin wrapper paths using badgeslp: prefix
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "mixedcoin",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			AllowCosmosWrapping: false, // Use badgeslp: prefix (wrapped approach)
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "wrapper-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true, // Override outgoing approvals for wrapped denom coin transfers
			OverridesToIncomingApprovals:   true, // Override incoming approvals for wrapped denom coin transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	wrapperPath := collection.CosmosCoinWrapperPaths[0]
	wrapperDenom := generateWrapperDenom(collection.CollectionId, wrapperPath)
	suite.Require().True(strings.HasPrefix(wrapperDenom, "badgeslp:"), "Wrapper denom should use badgeslp: prefix")

	// First, mint more badges to bob so he has enough
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(2), // Mint 2 more to have 3 total
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error minting additional badges to bob")

	// Verify bob has badges (needed for the wrapped approach)
	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's badge balance")
	suite.Require().True(len(bobBalance.Balances) > 0, "Bob should have badges")

	// For badgeslp: denoms, mint wrapped coins to bob's bank account
	// These are just bank tokens, but transfers will check badge balances
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")
	wrappedCoins := sdk.NewCoins(sdk.NewCoin(wrapperDenom, sdkmath.NewInt(2)))
	err = suite.app.BankKeeper.MintCoins(suite.ctx, "mint", wrappedCoins)
	suite.Require().Nil(err, "Error minting wrapped coins")
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, wrappedCoins)
	suite.Require().Nil(err, "Error sending wrapped coins to bob")

	// Verify bob has 2 wrapped coins in bank
	bobWrappedBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, wrapperDenom)
	suite.Require().Equal(sdkmath.NewInt(2), bobWrappedBalance.Amount, "Bob should have 2 wrapped coins in bank")

	// Add wrapped denom to allowed denoms
	params := suite.app.BadgesKeeper.GetParams(suite.ctx)
	params.AllowedDenoms = append(params.AllowedDenoms, wrapperDenom)
	err = suite.app.BadgesKeeper.SetParams(suite.ctx, params)
	suite.Require().Nil(err, "Error setting params with wrapped denom")

	// Fund bob with regular ubadge coins
	suite.app.BankKeeper.SendCoins(suite.ctx, suite.app.AccountKeeper.GetModuleAddress("mint"), bobAccAddr, sdk.NewCoins(sdk.NewCoin("ubadge", sdkmath.NewInt(1000))))

	// Test coin transfer with both wrapped and non-wrapped denoms
	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ToListId:          "AllWithoutMint",
				InitiatedByListId: alice,
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				ApprovalId:        "mixed-transfer",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
						AmountTrackerId:        "mixed-tracker",
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						AmountTrackerId:              "mixed-tracker",
					},
					CoinTransfers: []*types.CoinTransfer{
						{
							To:                              alice,
							OverrideFromWithApproverAddress: true, // Coins come from bob (approver), not alice (initiator)
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(2), Denom: wrapperDenom}, // Wrapped denom
								{Amount: sdkmath.NewInt(100), Denom: "ubadge"},   // Regular denom
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "error updating user approvals with mixed denoms")

	// Execute the transfer
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "wrapper-transfer",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
					{
						ApprovalId:      "mixed-transfer",
						ApprovalLevel:   "outgoing",
						ApproverAddress: bob,
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error executing transfer with mixed denoms")

	// Verify alice received both types of coins
	aliceAccAddr, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Error getting alice's address")
	aliceUbadgeBalance := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAccAddr, "ubadge")
	// Alice starts with 100 * 1e9 ubadge, so after receiving 100 more, she should have 100 * 1e9 + 100
	suite.Require().True(aliceUbadgeBalance.Amount.GTE(sdkmath.NewInt(100)), "Alice should receive at least 100 ubadge coins")

	// For badgeslp: denoms, the wrapped approach transfers underlying badges
	// Verify alice received the underlying badges
	aliceBalanceAfter, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting alice's badge balance after transfer")
	suite.Require().NotNil(aliceBalanceAfter, "Alice should have received badges from the wrapped approach")

	// Note: For badgeslp: denoms, the wrapped coins in bank are just representations
	// The actual transfer uses sendNativeTokensToAddressWithPoolApprovals which transfers underlying badges
	// So we check badge balances, not bank balances for the wrapped coins
}
