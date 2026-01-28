package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestCollectionApprovalAllowBackedMintingFalse tests that collection approvals with allowBackedMinting=false
// block backing transfers
func (suite *TestSuite) TestCollectionApprovalAllowBackedMintingFalse() {
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

	// Add collection approval with allowBackedMinting=false
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "backed-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: false, // Explicitly set to false
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Fund bob with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Try to unback - should fail because allowBackedMinting=false
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Error(err, "Should fail because allowBackedMinting=false")
	suite.Require().Contains(err.Error(), "does not allow backed minting operations", "Error should mention backed minting")
}

// TestCollectionApprovalAllowBackedMintingTrue tests that collection approvals with allowBackedMinting=true
// allow backing transfers
func (suite *TestSuite) TestCollectionApprovalAllowBackedMintingTrue() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin backed paths
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

	// Add collection approval with allowBackedMinting=true
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "backed-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true, // Explicitly set to true
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Fund bob with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Try to unback - should succeed because allowBackedMinting=true
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Should succeed because allowBackedMinting=true")

	// Verify bob has 1 token
	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount, "Bob should have 1 token after unbacking")
}

// TestCollectionApprovalAllowBackedMintingNil tests that collection approvals with nil ApprovalCriteria
// block backing transfers (defaults to false)
func (suite *TestSuite) TestCollectionApprovalAllowBackedMintingNil() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin backed paths
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

	// Add collection approval with nil ApprovalCriteria
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "backed-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria:  nil, // nil should default to false
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Fund bob with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Try to unback - should fail because ApprovalCriteria is nil (defaults to false)
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Error(err, "Should fail because ApprovalCriteria is nil")
	suite.Require().Contains(err.Error(), "does not allow backed minting operations", "Error should mention backed minting")
}

// TestCollectionApprovalAllowSpecialWrappingFalse tests that collection approvals with allowSpecialWrapping=false
// block wrapping transfers
func (suite *TestSuite) TestCollectionApprovalAllowSpecialWrappingFalse() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin wrapper paths
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Transfers = []*types.Transfer{}
	filteredApprovals := []*types.CollectionApproval{}
	for _, approval := range collectionsToCreate[0].CollectionApprovals {
		if approval.FromListId != "Mint" {
			filteredApprovals = append(filteredApprovals, approval)
		}
	}
	collectionsToCreate[0].CollectionApprovals = filteredApprovals
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom:  "wrapped",
			Symbol: "WBDG",
			Conversion: &types.ConversionWithoutDenom{
				SideA: &types.ConversionSideA{
					Amount: sdkmath.NewUint(1),
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

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	suite.Require().NotNil(collection.CosmosCoinWrapperPaths)
	suite.Require().Equal(1, len(collection.CosmosCoinWrapperPaths))
	wrapperPath := collection.CosmosCoinWrapperPaths[0]

	// Add collection approval with allowSpecialWrapping=false
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "wrapper-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowSpecialWrapping: false, // Explicitly set to false
		},
	})

	// Update the collection with the new approval
	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: collectionsToCreate[0].CollectionApprovals,
	})
	suite.Require().Nil(err, "error updating collection approvals")

	// Fund bob with wrapped coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	wrappedDenom := "wrapped"
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(wrappedDenom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(wrappedDenom, sdkmath.NewInt(1))})

	// Try to unwrap - should fail because allowSpecialWrapping=false
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        wrapperPath.Address,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Error(err, "Should fail because allowSpecialWrapping=false")
	suite.Require().Contains(err.Error(), "does not allow special wrapping operations", "Error should mention special wrapping")
}

// TestCollectionApprovalAllowSpecialWrappingTrue tests that collection approvals with allowSpecialWrapping=true
// allow wrapping transfers
func (suite *TestSuite) TestCollectionApprovalAllowSpecialWrappingTrue() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin wrapper paths
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	// Keep transfers so bob gets tokens initially - don't clear them
	// Keep all approvals including mint approval for the initial transfer
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom:  "wrapped",
			Symbol: "WBDG",
			Conversion: &types.ConversionWithoutDenom{
				SideA: &types.ConversionSideA{
					Amount: sdkmath.NewUint(1),
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

	// Add collection approval with allowSpecialWrapping=true
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "wrapper-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowSpecialWrapping: true, // Explicitly set to true
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	wrapperPath := collection.CosmosCoinWrapperPaths[0]

	// First, wrap tokens (transfer FROM bob TO wrapper address) to give bob wrapped coins
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{wrapperPath.Address},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Should succeed wrapping tokens")

	// Get the wrapped coin amount
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	// The wrapped denom will be prefixed with the collection ID
	fullDenom := generateWrappedWrapperDenom(collection.CollectionId, wrapperPath)
	bobBalanceDenom := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, fullDenom)
	bobAmount := sdkmath.NewUintFromBigInt(bobBalanceDenom.Amount.BigInt())
	suite.Require().True(bobAmount.GTE(sdkmath.NewUint(1)), "Bob should have wrapped coins")

	// Now try to unwrap - should succeed because allowSpecialWrapping=true
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        wrapperPath.Address,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Should succeed because allowSpecialWrapping=true")

	// Verify bob has 1 token
	bobBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(1), bobBalance.Balances[0].Amount, "Bob should have 1 token after unwrapping")
}

// TestBidirectionalBackingTransfer tests that backing transfers work in both directions
func (suite *TestSuite) TestBidirectionalBackingTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin backed paths
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

	// Add collection approval with allowBackedMinting=true
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "backed-transfer",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true,
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Test direction 1: from special address to user (unbacking)
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Unbacking should succeed")

	// Test direction 2: from user to special address (backing)
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Backing should succeed")
}
