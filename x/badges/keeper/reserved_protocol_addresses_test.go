package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
)

// TestReservedProtocolAddressCosmosCoinWrapperPath tests that cosmos coin wrapper path addresses are auto-reservedProtocol
func (suite *TestSuite) TestReservedProtocolAddressCosmosCoinWrapperPath() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with cosmos coin wrapper paths
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "testcoin",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					TokenIds:       GetOneUintRange(),
				},
			},
			AllowCosmosWrapping: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "test-approval",
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

	// Get the collection and verify the wrapper path address exists
	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	suite.Require().Equal(1, len(collection.CosmosCoinWrapperPaths), "Collection should have one cosmos coin wrapper path")

	wrapperPathAddress := collection.CosmosCoinWrapperPaths[0].Address

	// Verify the address is reservedProtocol
	isReservedProtocol := suite.app.BadgesKeeper.IsAddressReservedProtocolInStore(suite.ctx, wrapperPathAddress)
	suite.Require().True(isReservedProtocol, "Cosmos coin wrapper path address should be reservedProtocol")
}

// TestReservedProtocolAddressGammPool tests that gamm pool addresses are auto-reservedProtocol
func (suite *TestSuite) TestReservedProtocolAddressGammPool() {
	// Create a gamm pool
	poolId := suite.createTestPool()

	// Get the pool address
	poolAddress := poolmanagertypes.NewPoolAddress(poolId).String()

	// Verify the address is reservedProtocol
	isReservedProtocol := suite.app.BadgesKeeper.IsAddressReservedProtocolInStore(suite.ctx, poolAddress)
	suite.Require().True(isReservedProtocol, "Gamm pool address should be reservedProtocol")
}

// Helper function to create a test pool
func (suite *TestSuite) createTestPool() uint64 {
	// Fund bob with pool assets
	bobAcc := sdk.MustAccAddressFromBech32(bob)
	poolAssets := sdk.NewCoins(
		sdk.NewCoin("ubadge", sdkmath.NewInt(1000000)),
		sdk.NewCoin("stake", sdkmath.NewInt(1000000)),
	)

	// Fund the account
	banktestutil.FundAccount(suite.ctx, suite.app.BankKeeper, bobAcc, poolAssets)

	// Create a simple balancer pool for testing
	poolAssetsList := []balancer.PoolAsset{
		{
			Token:  sdk.NewCoin("ubadge", sdkmath.NewInt(1000000)),
			Weight: sdkmath.NewInt(1),
		},
		{
			Token:  sdk.NewCoin("stake", sdkmath.NewInt(1000000)),
			Weight: sdkmath.NewInt(1),
		},
	}

	poolParams := balancer.PoolParams{
		SwapFee: sdkmath.LegacyZeroDec(),
		ExitFee: sdkmath.LegacyZeroDec(),
	}

	msg := balancer.NewMsgCreateBalancerPool(
		bobAcc,
		poolParams,
		poolAssetsList,
	)

	// Use PoolManagerKeeper to create the pool
	poolId, err := suite.app.PoolManagerKeeper.CreatePool(suite.ctx, msg)
	suite.Require().Nil(err, "Error creating pool")
	suite.Require().NotZero(poolId, "Pool ID should not be zero")

	return poolId
}

// TestReservedProtocolAddressGovernanceSetUnset tests that only governance can set/unset reservedProtocol addresses
func (suite *TestSuite) TestReservedProtocolAddressGovernanceSetUnset() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	authority := suite.app.BadgesKeeper.GetAuthority()
	testAddress := bob

	// Test 1: Non-authority cannot set reservedProtocol address
	msg := &types.MsgSetReservedProtocolAddress{
		Authority:     bob, // Not the authority
		Address:       testAddress,
		IsReservedProtocol: true,
	}
	_, err := suite.msgServer.SetReservedProtocolAddress(wctx, msg)
	suite.Require().Error(err, "Non-authority should not be able to set reservedProtocol address")
	suite.Require().Contains(err.Error(), "invalid authority", "Error should mention invalid authority")

	// Test 2: Authority can set reservedProtocol address
	msg = &types.MsgSetReservedProtocolAddress{
		Authority:     authority,
		Address:       testAddress,
		IsReservedProtocol: true,
	}
	_, err = suite.msgServer.SetReservedProtocolAddress(wctx, msg)
	suite.Require().Nil(err, "Authority should be able to set reservedProtocol address")

	// Verify it's reservedProtocol
	isReservedProtocol := suite.app.BadgesKeeper.IsAddressReservedProtocolInStore(suite.ctx, testAddress)
	suite.Require().True(isReservedProtocol, "Address should be reservedProtocol")

	// Test 3: Authority can unset reservedProtocol address
	msg = &types.MsgSetReservedProtocolAddress{
		Authority:     authority,
		Address:       testAddress,
		IsReservedProtocol: false,
	}
	_, err = suite.msgServer.SetReservedProtocolAddress(wctx, msg)
	suite.Require().Nil(err, "Authority should be able to unset reservedProtocol address")

	// Verify it's not reservedProtocol
	isReservedProtocol = suite.app.BadgesKeeper.IsAddressReservedProtocolInStore(suite.ctx, testAddress)
	suite.Require().False(isReservedProtocol, "Address should not be reservedProtocol")
}

// TestReservedProtocolAddressGetAll tests that we can query all reservedProtocol addresses
func (suite *TestSuite) TestReservedProtocolAddressGetAll() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	authority := suite.app.BadgesKeeper.GetAuthority()

	// Set a few addresses as reservedProtocol
	testAddresses := []string{bob, alice, charlie}
	for _, addr := range testAddresses {
		msg := &types.MsgSetReservedProtocolAddress{
			Authority:     authority,
			Address:       addr,
			IsReservedProtocol: true,
		}
		_, err := suite.msgServer.SetReservedProtocolAddress(wctx, msg)
		suite.Require().Nil(err, "Error setting reservedProtocol address")
	}

	// Get all reservedProtocol addresses directly from keeper
	allReservedProtocol := suite.app.BadgesKeeper.GetAllReservedProtocolAddressesFromStore(suite.ctx)
	suite.Require().NotNil(allReservedProtocol, "Response should not be nil")

	// Verify all test addresses are in the list
	addressMap := make(map[string]bool)
	for _, addr := range allReservedProtocol {
		addressMap[addr] = true
	}

	for _, addr := range testAddresses {
		suite.Require().True(addressMap[addr], "Address %s should be in reservedProtocol addresses", addr)
	}

	// Clean up - unset the test addresses
	for _, addr := range testAddresses {
		msg := &types.MsgSetReservedProtocolAddress{
			Authority:     authority,
			Address:       addr,
			IsReservedProtocol: false,
		}
		_, err := suite.msgServer.SetReservedProtocolAddress(wctx, msg)
		suite.Require().Nil(err, "Error unsetting reservedProtocol address")
	}
}

// TestReservedProtocolAddressForcefulTransferError tests that forceful transfers from reservedProtocol addresses fail
func (suite *TestSuite) TestReservedProtocolAddressForcefulTransferError() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "test-approval-forceful",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "All", // Allow from Mint too
		ToListId:          "All",
		InitiatedByListId: "All",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesToIncomingApprovals:   true, // Allow transfers to any address
			OverridesFromOutgoingApprovals: true, // This triggers the reserved protocol check
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	// Manually blacklist an address
	authority := suite.app.BadgesKeeper.GetAuthority()
	reservedProtocolAddress := alice

	msg := &types.MsgSetReservedProtocolAddress{
		Authority:     authority,
		Address:       reservedProtocolAddress,
		IsReservedProtocol: true,
	}
	_, err = suite.msgServer.SetReservedProtocolAddress(wctx, msg)
	suite.Require().Nil(err, "Error setting reservedProtocol address")

	// Give alice some tokens first so she can attempt a transfer
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
	suite.Require().Nil(err, "Error transferring tokens to alice")

	// Try to forcefully transfer FROM the reservedProtocol address - should fail
	// Bob initiates the transfer from alice's address (forceful transfer)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob, // Bob initiates
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        reservedProtocolAddress, // From reservedProtocol address (alice)
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
	suite.Require().Error(err, "Forceful transfer from reservedProtocol address should fail")

	// Clean up
	msg = &types.MsgSetReservedProtocolAddress{
		Authority:     authority,
		Address:       reservedProtocolAddress,
		IsReservedProtocol: false,
	}
	_, err = suite.msgServer.SetReservedProtocolAddress(wctx, msg)
	suite.Require().Nil(err, "Error unsetting reservedProtocol address")
}

// TestReservedProtocolAddressNonForcefulTransferSucceeds tests that non-forceful transfers from reservedProtocol addresses succeed
func (suite *TestSuite) TestReservedProtocolAddressNonForcefulTransferSucceeds() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "test-approval-non-forceful",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetOneUintRange(),
		FromListId:        "All", // Allow from Mint too
		ToListId:          "All",
		InitiatedByListId: "All",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesToIncomingApprovals: true, // Allow transfers to any address
			// No OverridesFromOutgoingApprovals - should not trigger reserved protocol check
		},
	})

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	// Manually blacklist an address
	authority := suite.app.BadgesKeeper.GetAuthority()
	reservedProtocolAddress := alice

	msg := &types.MsgSetReservedProtocolAddress{
		Authority:     authority,
		Address:       reservedProtocolAddress,
		IsReservedProtocol: true,
	}
	_, err = suite.msgServer.SetReservedProtocolAddress(wctx, msg)
	suite.Require().Nil(err, "Error setting reservedProtocol address")

	// Give alice some tokens first
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
	suite.Require().Nil(err, "Error transferring tokens to alice")

	// Try to transfer FROM the reservedProtocol address WITHOUT OverridesFromOutgoingApprovals - should succeed
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        reservedProtocolAddress, // From reservedProtocol address
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
	suite.Require().Nil(err, "Non-forceful transfer from reservedProtocol address should succeed")

	// Clean up
	msg = &types.MsgSetReservedProtocolAddress{
		Authority:     authority,
		Address:       reservedProtocolAddress,
		IsReservedProtocol: false,
	}
	_, err = suite.msgServer.SetReservedProtocolAddress(wctx, msg)
	suite.Require().Nil(err, "Error unsetting reservedProtocol address")
}
