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
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	// Remove mint approvals since they're not allowed when cosmosCoinBackedPath is set
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/1234567890ABCDEF",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "backed-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Required for backing transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	suite.Require().NotNil(collection.Invariants)
	suite.Require().NotNil(collection.Invariants.CosmosCoinBackedPath)
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Verify initial balance is zero
	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	suite.Require().Equal(0, len(bobBalance.Balances), "Bob should have 0 tokens initially")

	// Fund bob with IBC coins (when unbacking FROM special address, user sends IBC coins TO special address)
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Unback to get tokens (backing minting process)
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Error unbacking token")

	// Verify bob has 1 token
	bobBalance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount, "Bob should have 1 token after unbacking")

	// Verify bob sent IBC coin to special address
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(0), bobIbcBalance.Amount, "Bob should have 0 IBC coins after unbacking")

	// Fund special address with IBC coins (when backing TO special address, special address sends IBC coins to user)
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Back the token
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Error backing token")

	// Verify bob received IBC coin from special address
	bobIbcBalance = suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(1), bobIbcBalance.Amount, "Bob should have 1 IBC coin after backing")

	// Verify bob has 0 tokens
	bobBalance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	suite.Require().Equal(0, len(bobBalance.Balances), "Bob should have 0 tokens after backing")
}

// TestCosmosCoinBackedPathsUnback tests unbacking IBC coins back to tokens
func (suite *TestSuite) TestCosmosCoinBackedPathsUnback() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/unwraptest",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
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
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Required for backing transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Fund bob with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Unback to get token
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Fund special address for backing
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Back the token
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Fund bob again for second unback
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Unback again
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Verify bob has token back
	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount, "Bob should have 1 token after unback")
}

// TestCosmosCoinBackedPathsTransferToOtherUser tests transferring IBC coins between users
func (suite *TestSuite) TestCosmosCoinBackedPathsTransferToOtherUser() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/transfertest",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
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
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Required for backing transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Fund bob with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Unback to get token
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Fund special address for backing
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Back the token
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Transfer IBC coin to alice
	aliceAccAddr, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err)
	err = suite.app.BankKeeper.SendCoins(suite.ctx, backedPathAccAddr, aliceAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.Require().Nil(err)

	// Alice needs IBC coins to unback (when unbacking, user sends IBC coins to special address)
	// Alice already has 1 IBC coin from the transfer above, which is enough

	// Alice unbacks (sends IBC coin to special address)
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Verify alice has token
	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(1), aliceBalance.Balances[0].Amount, "Alice should have 1 token")
}

// TestCosmosCoinBackedPathsErrors tests various error scenarios
func (suite *TestSuite) TestCosmosCoinBackedPathsErrors() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/errortest",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
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
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Required for backing transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Test backing without tokens - should fail
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
	suite.Require().Error(err, "Should error when trying to back without tokens")
}

// TestCosmosCoinBackedPathsMultipleDenoms tests collections with cosmos coin backed paths
func (suite *TestSuite) TestCosmosCoinBackedPathsMultipleDenoms() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/coin-one",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
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
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Required for backing transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	backedPath1 := collection.Invariants.CosmosCoinBackedPath

	// Fund bob with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath1.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath1.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Unback
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        backedPath1.Address,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Fund special address for backing
	backedPath1AccAddr, err := sdk.AccAddressFromBech32(backedPath1.Address)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath1.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", backedPath1AccAddr, sdk.Coins{sdk.NewCoin(backedPath1.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Back
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Verify bob received IBC coin
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath1.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(1), bobIbcBalance.Amount, "Bob should have 1 IBC coin")
}

// TestCosmosCoinBackedPathsInadequateBalance tests inadequate IBC coin balance scenarios
func (suite *TestSuite) TestCosmosCoinBackedPathsInadequateBalance() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/balance-test",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
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
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Required for backing transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Try to back without tokens - should fail
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Error(err, "Should error when trying to back without tokens")
	suite.Require().Contains(err.Error(), "inadequate balances", "Error should mention inadequate balances")
}

// TestCosmosCoinBackedPathsUnbackInadequateBalance tests unbacking when IBC coins are not available
func (suite *TestSuite) TestCosmosCoinBackedPathsUnbackInadequateBalance() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/unback-balance-test",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
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
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Required for backing transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Try to unback without IBC coins - should fail
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Error(err, "Should error when trying to unback without IBC coins")
	suite.Require().Contains(err.Error(), "insufficient funds", "Error should mention insufficient funds")
}

// TestCosmosCoinBackedPathsConversionRate tests different conversion rates
func (suite *TestSuite) TestCosmosCoinBackedPathsConversionRate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/conversion-test",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(5), // 5 tokens = 1 IBC coin
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
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
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Required for backing transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Fund bob with IBC coins (5 tokens = 1 IBC coin)
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Unback 5 tokens
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        backedPath.Address,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(5),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Verify bob has 5 tokens
	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	fetchedBobBalances, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobBalance.Balances)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(5), fetchedBobBalances[0].Amount, "Bob should have 5 tokens")

	// Fund special address for backing
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Back 5 tokens
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Verify bob has 0 tokens
	bobBalance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	fetchedBobBalances, err = types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobBalance.Balances)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(0), fetchedBobBalances[0].Amount, "Bob should have 0 tokens after backing")
}

// TestCosmosCoinBackedPathsIbcAmount tests the ibcAmount field functionality
func (suite *TestSuite) TestCosmosCoinBackedPathsIbcAmount() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(10), // 10 IBC coins per conversion unit
					Denom:  "ibc/ibcamount-test",
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
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
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Required for backing transfers
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	backedPath := collection.Invariants.CosmosCoinBackedPath
	suite.Require().Equal(sdkmath.NewUint(10), backedPath.Conversion.SideA.Amount, "IbcAmount should be 10")

	// Fund bob with 10 IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(10))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(10))})

	// Unback 1 token (requires 10 IBC coins)
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Verify bob has 1 token
	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	fetchedBobBalances, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobBalance.Balances)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBobBalances[0].Amount, "Bob should have 1 token")

	// Fund special address with 10 IBC coins for backing
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(10))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(10))})

	// Back 1 token (should send 10 IBC coins to bob)
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.BadgesKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Verify bob received 10 IBC coins
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(10), bobIbcBalance.Amount, "Bob should have 10 IBC coins after backing")

	// Verify bob has 0 tokens
	bobBalance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	fetchedBobBalances, err = types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobBalance.Balances)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(0), fetchedBobBalances[0].Amount, "Bob should have 0 tokens after backing")
}
