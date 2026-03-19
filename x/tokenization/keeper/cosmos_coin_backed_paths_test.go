package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// createBackedCollectionAndAddApproval creates a backed collection, then adds the backed minting
// approval with the proper backing address as list ID. Returns the collection with the approval.
// The approvalId parameter is used as the approval ID for the backed minting approval.
func (suite *TestSuite) createBackedCollectionAndAddApproval(
	wctx sdk.Context,
	collectionsToCreate []*types.MsgNewCollection,
	approvalId string,
) (*types.TokenCollection, error) {
	err := CreateCollections(suite, sdk.WrapSDKContext(wctx), collectionsToCreate)
	if err != nil {
		return nil, err
	}

	collection, err := GetCollection(suite, sdk.WrapSDKContext(wctx), sdkmath.NewUint(1))
	if err != nil {
		return nil, err
	}

	backingAddr := collection.Invariants.CosmosCoinBackedPath.Address

	// Add backed minting approval with proper backing address
	approvals := collection.CollectionApprovals
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        approvalId,
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        backingAddr,
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true,
			MustPrioritize:     true,
		},
	})
	// Also add approval for backing direction (to backing address)
	approvals = append(approvals, &types.CollectionApproval{
		ApprovalId:        approvalId + "-back",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          backingAddr,
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			AllowBackedMinting: true,
			MustPrioritize:     true,
		},
	})

	err = UpdateCollectionApprovals(suite, sdk.WrapSDKContext(wctx), &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: approvals,
	})
	if err != nil {
		return nil, err
	}

	// Re-fetch to get updated versions
	collection, err = GetCollection(suite, sdk.WrapSDKContext(wctx), sdkmath.NewUint(1))
	return collection, err
}

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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "backed-transfer")
	suite.Require().Nil(err, "error creating backed collection")
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "unback-transfer")
	suite.Require().Nil(err, "error creating backed collection")
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "transfer-test")
	suite.Require().Nil(err, "error creating backed collection")
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "error-test")
	suite.Require().Nil(err, "error creating backed collection")
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "multi-denom-test")
	suite.Require().Nil(err, "error creating backed collection")
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "balance-test")
	suite.Require().Nil(err, "error creating backed collection")
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "unback-balance-test")
	suite.Require().Nil(err, "error creating backed collection")
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "conversion-test")
	suite.Require().Nil(err, "error creating backed collection")
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "ibcamount-test")
	suite.Require().Nil(err, "error creating backed collection")
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
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

// TestCosmosCoinBackedPathsUnbackOnBehalfOf tests unbacking (deposit) on behalf of another user.
// Initiator (bob) pays IBC coins, tokens go to alice.
func (suite *TestSuite) TestCosmosCoinBackedPathsUnbackOnBehalfOf() {
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
					Denom:  "ibc/onbehalf-unback",
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "backed-transfer")
	suite.Require().Nil(err, "error creating backed collection")
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Fund bob (the initiator) with IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Bob (Creator/initiator) unbacks tokens FROM backing address TO alice
	// Bob pays the IBC coins, alice receives the tokens
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Error unbacking on behalf of alice")

	// Verify alice has 1 token
	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(1), aliceBalance.Balances[0].Amount, "Alice should have 1 token")

	// Verify bob's IBC coins were debited (not alice's)
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(0), bobIbcBalance.Amount, "Bob should have 0 IBC coins after paying for unbacking")

	// Verify alice has no IBC coins (she didn't pay)
	aliceAccAddr, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err)
	aliceIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAccAddr, backedPath.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(0), aliceIbcBalance.Amount, "Alice should have 0 IBC coins (she didn't pay)")
}

// TestCosmosCoinBackedPathsBackOnBehalfOf tests backing (withdraw) on behalf of another user.
// Alice has tokens, bob (initiator) initiates backing alice's tokens, bob receives IBC coins.
func (suite *TestSuite) TestCosmosCoinBackedPathsBackOnBehalfOf() {
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
					Denom:  "ibc/onbehalf-back",
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "backed-transfer")
	suite.Require().Nil(err, "error creating backed collection")
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// First, give alice some tokens via unbacking (alice pays her own IBC coins)
	aliceAccAddr, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", aliceAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Error giving alice tokens")

	// Verify alice has 1 token
	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(1), aliceBalance.Balances[0].Amount)

	// Alice sets outgoing approval for bob to transfer her tokens
	err = SetOutgoingApproval(suite, wctx, &types.MsgSetOutgoingApproval{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Approval: &types.UserOutgoingApproval{
			ToListId:          "AllWithoutMint",
			InitiatedByListId: bob,
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetOneUintRange(),
			ApprovalId:        "bob-can-back",
		},
	})
	suite.Require().Nil(err, "Error setting outgoing approval")

	// Fund special address with IBC coins for backing
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Bob (initiator) backs alice's tokens TO backing address
	// Alice's tokens get backed, bob receives the IBC coins
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        alice,
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
	suite.Require().Nil(err, "Error backing on behalf of alice")

	// Verify alice has 0 tokens
	aliceBalance, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err)
	suite.Require().Equal(0, len(aliceBalance.Balances), "Alice should have 0 tokens after backing")

	// Verify bob received the IBC coins (not alice)
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(1), bobIbcBalance.Amount, "Bob should have 1 IBC coin after backing")

	// Verify alice has no IBC coins (she didn't receive them)
	aliceIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAccAddr, backedPath.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(0), aliceIbcBalance.Amount, "Alice should have 0 IBC coins")
}

// TestCosmosCoinBackedPathsBackOnBehalfOfNoApproval tests that backing on behalf fails without outgoing approval.
func (suite *TestSuite) TestCosmosCoinBackedPathsBackOnBehalfOfNoApproval() {
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
	// Clear default outgoing approvals so alice must explicitly approve bob
	collectionsToCreate[0].DefaultOutgoingApprovals = []*types.UserOutgoingApproval{}
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{
			Conversion: &types.Conversion{
				SideA: &types.ConversionSideAWithDenom{
					Amount: sdkmath.NewUint(1),
					Denom:  "ibc/noapproval-back",
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "backed-transfer")
	suite.Require().Nil(err, "error creating backed collection")
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Give alice tokens
	aliceAccAddr, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", aliceAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err)

	// Fund special address
	backedPathAccAddr, err := sdk.AccAddressFromBech32(backedPath.Address)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", backedPathAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(1))})

	// Bob tries to back alice's tokens WITHOUT alice's outgoing approval - should fail
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        alice,
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
	suite.Require().Error(err, "Should fail without outgoing approval from alice")
}

// TestCosmosCoinBackedPathsUnbackOnBehalfOfConversionRate tests on-behalf unbacking with conversion rates.
func (suite *TestSuite) TestCosmosCoinBackedPathsUnbackOnBehalfOfConversionRate() {
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
					Amount: sdkmath.NewUint(10), // 10 IBC coins per token
					Denom:  "ibc/onbehalf-conversion",
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

	collection, err := suite.createBackedCollectionAndAddApproval(suite.ctx, collectionsToCreate, "backed-transfer")
	suite.Require().Nil(err, "error creating backed collection")
	backedPath := collection.Invariants.CosmosCoinBackedPath

	// Fund bob (initiator) with 10 IBC coins
	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err)
	suite.app.BankKeeper.MintCoins(suite.ctx, "mint", sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(10))})
	suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", bobAccAddr, sdk.Coins{sdk.NewCoin(backedPath.Conversion.SideA.Denom, sdkmath.NewInt(10))})

	// Bob unbacks 1 token on behalf of alice (bob pays 10 IBC coins)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
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
				PrioritizedApprovals: GetPrioritizedApprovalsFromCollection(suite.ctx, suite.app.TokenizationKeeper, collection),
			},
		},
	})
	suite.Require().Nil(err, "Error unbacking on behalf with conversion rate")

	// Verify alice has 1 token
	aliceBalance, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err)
	suite.Require().Equal(sdkmath.NewUint(1), aliceBalance.Balances[0].Amount)

	// Verify bob paid 10 IBC coins
	bobIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, backedPath.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(0), bobIbcBalance.Amount, "Bob should have 0 IBC coins after paying 10")

	// Verify alice didn't pay anything
	aliceAccAddr, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err)
	aliceIbcBalance := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAccAddr, backedPath.Conversion.SideA.Denom)
	suite.Require().Equal(sdkmath.NewInt(0), aliceIbcBalance.Amount, "Alice should have 0 IBC coins")
}
