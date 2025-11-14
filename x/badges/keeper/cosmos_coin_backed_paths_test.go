package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestCosmosCoinBackedPathsBasic tests the basic functionality of cosmos coin backed paths
// by creating a collection with backed paths and transferring tokens to the special alias address
func (suite *TestSuite) TestCosmosCoinBackedPathsBasic() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin backed paths
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinBackedPathsToAdd = []*types.CosmosCoinBackedPathAddObject{
		{
			IbcDenom: "ibc/1234567890ABCDEF",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(1), // Default: 1 IBC coin per conversion unit
		},
	}

	// Add collection approvals for transfers to/from backed paths
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "backed-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  &types.ApprovalCriteria{},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection with cosmos coin backed paths")

	// Verify the collection was created with backed paths
	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	suite.Require().Equal(1, len(collection.CosmosCoinBackedPaths), "Collection should have one cosmos coin backed path")

	backedPath := collection.CosmosCoinBackedPaths[0]
	suite.Require().Equal("ibc/1234567890ABCDEF", backedPath.IbcDenom, "Backed path ibc denom should match")

	// Get initial balance
	bobBalanceBefore, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance")
	suite.Require().Equal(sdkmath.NewUint(1), bobBalanceBefore.Balances[0].Amount, "Bob should have 1 token initially")

	// Fund bob with IBC coins first (since backed paths use existing coins)
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})

	// Transfer token to the special alias address (backed path address)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
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
	suite.Require().Nil(err, "Error transferring token to backed path address")

	// Verify the token was transferred (burned from bob)
	bobBalanceAfter, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance after transfer")

	diffInBalances, err := types.SubtractBalances(suite.ctx, bobBalanceAfter.Balances, bobBalanceBefore.Balances)
	suite.Require().Nil(err, "Error subtracting balances")
	suite.Require().Equal(1, len(diffInBalances), "Bob should have lost one badge")
	suite.Require().Equal(sdkmath.NewUint(1), diffInBalances[0].Amount, "Bob should have lost exactly 1 token")

	// Verify the IBC coin was sent from bob to the backed path address
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err, "Error getting backed path address")
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(0), bobIbcBalance.Amount, "Bob should have 0 IBC coins after transfer")

	backedPathIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, backedPathAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(1), backedPathIbcBalance.Amount, "Backed path address should have 1 IBC coin")
}

// TestCosmosCoinBackedPathsUnback tests unbacking IBC coins back to tokens
func (suite *TestSuite) TestCosmosCoinBackedPathsUnback() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin backed paths
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinBackedPathsToAdd = []*types.CosmosCoinBackedPathAddObject{
		{
			IbcDenom: "ibc/unwraptest",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(1),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "unback-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  &types.ApprovalCriteria{},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.CosmosCoinBackedPaths[0]

	// Fund bob with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})

	// Back the token first
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
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
	suite.Require().Nil(err, "Error backing token")

	// Get the backed coin amount
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(0), bobIbcBalance.Amount, "Bob should have 0 IBC coins after backing")

	// Now unback the coin back to tokens
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        backedPath.Address,
				ToAddresses: []string{bob},
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
	suite.Require().Nil(err, "Error unbacking coin")

	// Verify the token was restored
	bobBalanceAfterUnback, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance after unback")
	suite.Require().Equal(sdkmath.NewUint(1), bobBalanceAfterUnback.Balances[0].Amount, "Bob should have 1 token after unback")

	// Verify the IBC coin was sent back to bob
	bobIbcBalanceAfterUnback := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(1), bobIbcBalanceAfterUnback.Amount, "Bob should have 1 IBC coin after unback")
}

// TestCosmosCoinBackedPathsTransferToOtherUser tests transferring IBC coins between users
func (suite *TestSuite) TestCosmosCoinBackedPathsTransferToOtherUser() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin backed paths
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinBackedPathsToAdd = []*types.CosmosCoinBackedPathAddObject{
		{
			IbcDenom: "ibc/transfertest",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(1),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "transfer-test",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  &types.ApprovalCriteria{},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.CosmosCoinBackedPaths[0]

	// Fund bob with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})

	// Back the badge
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
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
	suite.Require().Nil(err, "Error backing token")

	// Transfer the IBC coin to alice
	aliceAccAddr, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Error getting alice's address")
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err, "Error getting backed path address")

	err = suite.app.BankKeeper.SendCoins(suite.ctx, backedPathAccAddr, aliceAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})
	suite.Require().Nil(err, "Error transferring IBC coin to alice")

	// Verify alice has the IBC coin
	aliceIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(1), aliceIbcBalance.Amount, "Alice should have 1 IBC coin")

	// Alice needs to send the IBC coin back to the backed path address to unback
	err = suite.app.BankKeeper.SendCoins(suite.ctx, aliceAccAddr, backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})
	suite.Require().Nil(err, "Error sending IBC coin back to backed path address")

	// Now alice should be able to unback the coin
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        backedPath.Address,
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
	suite.Require().Nil(err, "Error unbacking coin as alice")

	// Verify alice now has the badge
	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err, "Error getting alice's balance")
	suite.Require().Equal(sdkmath.NewUint(1), aliceBalance.Balances[0].Amount, "Alice should have 1 token after unback")

	// Verify alice received the IBC coin back
	aliceIbcBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(1), aliceIbcBalanceAfter.Amount, "Alice should have 1 IBC coin after unback")
}

// TestCosmosCoinBackedPathsErrors tests various error scenarios
func (suite *TestSuite) TestCosmosCoinBackedPathsErrors() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin backed paths
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinBackedPathsToAdd = []*types.CosmosCoinBackedPathAddObject{
		{
			IbcDenom: "ibc/errortest",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(1),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "error-test",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  &types.ApprovalCriteria{},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.CosmosCoinBackedPaths[0]

	// Test transferring more tokens than available
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(2), // More than available
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Should error when transferring more tokens than available")

	// Test transferring with wrong token IDs
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTwoUintRanges(), // Wrong token ID
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Should error when transferring with wrong token IDs")

	// Test transferring with wrong ownership times
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetTwoUintRanges(), // Wrong ownership times
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Should error when transferring with wrong ownership times")
}

// TestCosmosCoinBackedPathsMultipleDenoms tests collections with multiple cosmos coin backed paths
func (suite *TestSuite) TestCosmosCoinBackedPathsMultipleDenoms() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with multiple cosmos coin backed paths
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinBackedPathsToAdd = []*types.CosmosCoinBackedPathAddObject{
		{
			IbcDenom: "ibc/coin-one",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(1),
		},
		{
			IbcDenom: "ibc/coin-two",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(1),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "multi-denom-test",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  &types.ApprovalCriteria{},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection with multiple cosmos coin backed paths")

	// Verify the collection was created with multiple backed paths
	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	suite.Require().Equal(2, len(collection.CosmosCoinBackedPaths), "Collection should have two cosmos coin backed paths")

	// Test backing with first denom
	backedPath1 := collection.CosmosCoinBackedPaths[0]
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")

	// Fund bob with first IBC coin
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath1.IbcDenom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath1.IbcDenom, sdkmath.NewInt(1))})

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath1.Address},
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
	suite.Require().Nil(err, "Error backing token with first denom")

	// Verify first denom was used
	backedPath1AccAddr, err := sdk.AccAddressFromBech32(backedPath1.Address)
	suite.Require().Nil(err, "Error getting backed path 1 address")
	backedPath1IbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, backedPath1AccAddr, backedPath1.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(1), backedPath1IbcBalance.Amount, "Backed path 1 should have 1 of first IBC coin")

	// Test that we can't back the same token again (it's already backed)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath1.Address},
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
	suite.Require().Error(err, "Should error when trying to back already backed token")
}

// TestCosmosCoinBackedPathsInadequateBalance tests inadequate IBC coin balance scenarios
func (suite *TestSuite) TestCosmosCoinBackedPathsInadequateBalance() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinBackedPathsToAdd = []*types.CosmosCoinBackedPathAddObject{
		{
			IbcDenom: "ibc/balance-test",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(1),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "balance-test",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  &types.ApprovalCriteria{},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.CosmosCoinBackedPaths[0]

	// Try to back without having IBC coins - should fail
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
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
	suite.Require().Error(err, "Should error when trying to back without IBC coins")
	suite.Require().Contains(err.Error(), "insufficient funds", "Error should mention insufficient funds")
}

// TestCosmosCoinBackedPathsUnbackInadequateBalance tests unbacking when IBC coins are not available
func (suite *TestSuite) TestCosmosCoinBackedPathsUnbackInadequateBalance() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinBackedPathsToAdd = []*types.CosmosCoinBackedPathAddObject{
		{
			IbcDenom: "ibc/unback-balance-test",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(1),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "unback-balance-test",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  &types.ApprovalCriteria{},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.CosmosCoinBackedPaths[0]

	// Fund bob with IBC coins and back the token
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
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
	suite.Require().Nil(err, "Error backing token")

	// Transfer the IBC coin to alice (so backed path address doesn't have it)
	aliceAccAddr, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Error getting alice's address")
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err, "Error getting backed path address")

	err = suite.app.BankKeeper.SendCoins(suite.ctx, backedPathAccAddr, aliceAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})
	suite.Require().Nil(err, "Error transferring IBC coin to alice")

	// Try to unback - should fail because backed path address doesn't have IBC coins
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        backedPath.Address,
				ToAddresses: []string{bob},
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
	suite.Require().Error(err, "Should error when trying to unback without IBC coins")
	suite.Require().Contains(err.Error(), "insufficient funds", "Error should mention insufficient funds")

	// Alice should send the IBC coins back to the backed path address to unback
	err = suite.app.BankKeeper.SendCoins(suite.ctx, aliceAccAddr, backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})
	suite.Require().Nil(err, "Error sending IBC coin back to backed path address")

	// Now alice should be able to unback since the backed path address has the IBC coins
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        backedPath.Address,
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
	suite.Require().Nil(err, "Alice should be able to unback since the backed path address has the IBC coins")
}

// TestCosmosCoinBackedPathsConversionRate tests different conversion rates
func (suite *TestSuite) TestCosmosCoinBackedPathsConversionRate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with a conversion rate of 5 tokens = 1 IBC coin
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinBackedPathsToAdd = []*types.CosmosCoinBackedPathAddObject{
		{
			IbcDenom: "ibc/conversion-test",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(5), // 5 tokens = 1 IBC coin
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(1),
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "conversion-test",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  &types.ApprovalCriteria{},
	})

	// Mint 5 tokens to bob
	collectionsToCreate[0].Transfers[0].Balances[0].Amount = sdkmath.NewUint(5)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.CosmosCoinBackedPaths[0]

	// Fund bob with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(1))})

	// Verify bob has 5 tokens before backing
	bobBalanceBefore, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's balance before backing")
	fetchedBobBalancesBefore, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobBalanceBefore.Balances)
	suite.Require().Nil(err, "Error fetching bob balances for IDs before backing")
	suite.Require().Equal(sdkmath.NewUint(5), fetchedBobBalancesBefore[0].Amount, "Bob should have 5 tokens before backing")

	// Back 5 tokens (should require 1 IBC coin)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(5),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error backing 5 tokens")

	// Verify bob has 0 IBC coins
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(0), bobIbcBalance.Amount, "Bob should have 0 IBC coins after backing")

	// Verify backed path address has 1 IBC coin
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err, "Error getting backed path address")
	backedPathIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, backedPathAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(1), backedPathIbcBalance.Amount, "Backed path should have 1 IBC coin")

	// Verify backed path address has the 5 tokens
	backedPathBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), backedPath.Address)
	suite.Require().Nil(err, "Error getting backed path balance")
	fetchedBackedPathBalances, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), backedPathBalance.Balances)
	suite.Require().Nil(err, "Error fetching backed path balances for IDs")
	suite.Require().Equal(sdkmath.NewUint(5), fetchedBackedPathBalances[0].Amount, "Backed path should have 5 tokens")

	// Verify bob has 0 tokens
	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's balance")
	fetchedBobBalances, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobBalance.Balances)
	suite.Require().Nil(err, "Error fetching bob balances for IDs")
	suite.Require().Equal(sdkmath.NewUint(0), fetchedBobBalances[0].Amount, "Bob should have 0 tokens after backing")
}

// TestCosmosCoinBackedPathsIbcAmount tests the ibcAmount field functionality
func (suite *TestSuite) TestCosmosCoinBackedPathsIbcAmount() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin backed paths with ibcAmount = 10
	// This means: 1 token = 10 IBC coins (instead of the default 1:1)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinBackedPathsToAdd = []*types.CosmosCoinBackedPathAddObject{
		{
			IbcDenom: "ibc/ibcamount-test",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			IbcAmount: sdkmath.NewUint(10), // 10 IBC coins per conversion unit
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "ibcamount-test",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  &types.ApprovalCriteria{},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.CosmosCoinBackedPaths[0]
	suite.Require().Equal(sdkmath.NewUint(10), backedPath.IbcAmount, "IbcAmount should be 10")

	// Fund bob with 10 IBC coins (since ibcAmount = 10, we need 10 coins for 1 token)
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting bob's address")
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(10))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.IbcDenom, sdkmath.NewInt(10))})

	// Back 1 token (should require 10 IBC coins due to ibcAmount = 10)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{backedPath.Address},
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
	suite.Require().Nil(err, "Error backing token with ibcAmount")

	// Verify bob has 0 IBC coins (all 10 were sent)
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(0), bobIbcBalance.Amount, "Bob should have 0 IBC coins after backing")

	// Verify backed path address has 10 IBC coins
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err, "Error getting backed path address")
	backedPathIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, backedPathAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(10), backedPathIbcBalance.Amount, "Backed path should have 10 IBC coins")

	// Verify bob has 0 tokens
	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's balance")
	fetchedBobBalances, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobBalance.Balances)
	suite.Require().Nil(err, "Error fetching bob balances for IDs")
	suite.Require().Equal(sdkmath.NewUint(0), fetchedBobBalances[0].Amount, "Bob should have 0 tokens after backing")

	// Now unback - should send 10 IBC coins back to bob
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        backedPath.Address,
				ToAddresses: []string{bob},
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
	suite.Require().Nil(err, "Error unbacking token")

	// Verify bob received 10 IBC coins back
	bobIbcBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.IbcDenom)
	suite.Require().Equal(sdkmath.NewInt(10), bobIbcBalanceAfter.Amount, "Bob should have 10 IBC coins after unback")

	// Verify bob has 1 token back
	bobBalanceAfter, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting bob's balance after unback")
	fetchedBobBalancesAfter, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobBalanceAfter.Balances)
	suite.Require().Nil(err, "Error fetching bob balances for IDs after unback")
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBobBalancesAfter[0].Amount, "Bob should have 1 token after unback")
}
