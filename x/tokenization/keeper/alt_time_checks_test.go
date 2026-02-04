package keeper_test

import (
	"math"
	"time"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestAltTimeChecks_CollectionApproval_OfflineHours tests that collection approvals deny transfers during offline hours
func (suite *TestSuite) TestAltTimeChecks_CollectionApproval_OfflineHours() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create address list (check if it already exists)
	_, err := GetAddressList(suite, wctx, "AllWithoutMint")
	if err != nil {
		err = suite.app.TokenizationKeeper.CreateAddressList(suite.ctx, &types.AddressList{
			ListId:    "AllWithoutMint",
			Addresses: []string{alice, bob},
		})
		suite.Require().Nil(err, "error creating address list")
	}

	// Create a collection with offline hours (9am-5pm UTC, hours 9-16)
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ToListId:          "AllWithoutMint",
					FromListId:        "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TransferTimes:     GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TokenIds:          GetFullUintRanges(),
					ApprovalId:        "test-offline-hours",
					ApprovalCriteria: &types.ApprovalCriteria{
						AltTimeChecks: &types.AltTimeChecks{
							OfflineHours: []*types.UintRange{
								{
									Start: sdkmath.NewUint(9),  // 9am UTC
									End:   sdkmath.NewUint(16), // 4pm UTC (inclusive, so 5pm is hour 16)
								},
							},
						},
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
							AmountTrackerId:        "test-tracker",
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(uint64(math.MaxUint64)),
							AmountTrackerId:              "test-tracker",
						},
					},
				},
			},
			TokensToCreate: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					TokenIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection:         []*types.ActionPermission{},
				CanUpdateStandards:           []*types.ActionPermission{},
				CanUpdateCustomData:          []*types.ActionPermission{},
				CanDeleteCollection:          []*types.ActionPermission{},
				CanUpdateManager:             []*types.ActionPermission{},
				CanUpdateCollectionMetadata:  []*types.ActionPermission{},
				CanUpdateTokenMetadata:       []*types.TokenIdsActionPermission{},
				CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
				CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
					{
						PermanentlyPermittedTimes: GetFullUintRanges(),
					},
				},
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Test 1: Transfer during offline hours (10am UTC) - should be denied
	// Set block time to 10am UTC on a specific date (e.g., 2024-01-01 10:00:00 UTC)
	testTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking := &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-offline-hours",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		},
		false,
		false,
		false,
		nil,
		eventTracking,
	)

	suite.Require().NotNil(err, "transfer should be denied during offline hours")
	// The error might be wrapped, so check for either the direct message or the wrapped message
	suite.Require().True(
		contains(err.Error(), "alt time check failed") || contains(err.Error(), "transfer denied") || contains(err.Error(), "inadequate approvals"),
		"error should mention alt time check or transfer denial: %s", err.Error(),
	)

	// Test 2: Transfer outside offline hours (8am UTC) - should be allowed
	testTime = time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking = &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-offline-hours",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		},
		false,
		false,
		false,
		nil,
		eventTracking,
	)

	suite.Require().Nil(err, "transfer should be allowed outside offline hours")
}

// TestAltTimeChecks_CollectionApproval_OfflineDays tests that collection approvals deny transfers during offline days
func (suite *TestSuite) TestAltTimeChecks_CollectionApproval_OfflineDays() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create address list (check if it already exists)
	_, err := GetAddressList(suite, wctx, "AllWithoutMint")
	if err != nil {
		err = suite.app.TokenizationKeeper.CreateAddressList(suite.ctx, &types.AddressList{
			ListId:    "AllWithoutMint",
			Addresses: []string{alice, bob},
		})
		suite.Require().Nil(err, "error creating address list")
	}

	// Create a collection with offline days (Monday-Friday, days 1-5)
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ToListId:          "AllWithoutMint",
					FromListId:        "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TransferTimes:     GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TokenIds:          GetFullUintRanges(),
					ApprovalId:        "test-offline-days",
					ApprovalCriteria: &types.ApprovalCriteria{
						AltTimeChecks: &types.AltTimeChecks{
							OfflineDays: []*types.UintRange{
								{
									Start: sdkmath.NewUint(1), // Monday
									End:   sdkmath.NewUint(5), // Friday
								},
							},
						},
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
							AmountTrackerId:        "test-tracker",
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(uint64(math.MaxUint64)),
							AmountTrackerId:              "test-tracker",
						},
					},
				},
			},
			TokensToCreate: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					TokenIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection:         []*types.ActionPermission{},
				CanUpdateStandards:           []*types.ActionPermission{},
				CanUpdateCustomData:          []*types.ActionPermission{},
				CanDeleteCollection:          []*types.ActionPermission{},
				CanUpdateManager:             []*types.ActionPermission{},
				CanUpdateCollectionMetadata:  []*types.ActionPermission{},
				CanUpdateTokenMetadata:       []*types.TokenIdsActionPermission{},
				CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
				CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
					{
						PermanentlyPermittedTimes: GetFullUintRanges(),
					},
				},
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Test 1: Transfer on Monday (day 1) - should be denied
	// January 1, 2024 is a Monday
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking := &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-offline-days",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		},
		false,
		false,
		false,
		nil,
		eventTracking,
	)

	suite.Require().NotNil(err, "transfer should be denied on offline days")
	// The error might be wrapped, so check for either the direct message or the wrapped message
	suite.Require().True(
		contains(err.Error(), "alt time check failed") || contains(err.Error(), "transfer denied") || contains(err.Error(), "inadequate approvals"),
		"error should mention alt time check or transfer denial: %s", err.Error(),
	)

	// Test 2: Transfer on Saturday (day 6) - should be allowed
	// January 6, 2024 is a Saturday
	testTime = time.Date(2024, 1, 6, 12, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking = &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-offline-days",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		},
		false,
		false,
		false,
		nil,
		eventTracking,
	)

	suite.Require().Nil(err, "transfer should be allowed on non-offline days")
}

// TestAltTimeChecks_IncomingApproval tests that incoming approvals respect altTimeChecks
func (suite *TestSuite) TestAltTimeChecks_IncomingApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create address list for alice and bob BEFORE creating the collection
	err := suite.app.TokenizationKeeper.CreateAddressList(suite.ctx, &types.AddressList{
		ListId:    "aliceAndBob",
		Addresses: []string{alice, bob},
		Whitelist: true,
	})
	suite.Require().Nil(err, "error creating address list")

	// Create a collection with default incoming approval that has offline hours
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			DefaultIncomingApprovals: []*types.UserIncomingApproval{
				{
					FromListId:        "aliceAndBob",
					InitiatedByListId: "aliceAndBob",
					TransferTimes:     GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TokenIds:          GetFullUintRanges(),
					ApprovalId:        "test-incoming-offline",
					ApprovalCriteria: &types.IncomingApprovalCriteria{
						AltTimeChecks: &types.AltTimeChecks{
							OfflineHours: []*types.UintRange{
								{
									Start: sdkmath.NewUint(9),
									End:   sdkmath.NewUint(17),
								},
							},
						},
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
							AmountTrackerId:        "test-tracker",
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(uint64(math.MaxUint64)),
							AmountTrackerId:              "test-tracker",
						},
					},
				},
			},
			TokensToCreate: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					TokenIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection:         []*types.ActionPermission{},
				CanUpdateStandards:           []*types.ActionPermission{},
				CanUpdateCustomData:          []*types.ActionPermission{},
				CanDeleteCollection:          []*types.ActionPermission{},
				CanUpdateManager:             []*types.ActionPermission{},
				CanUpdateCollectionMetadata:  []*types.ActionPermission{},
				CanUpdateTokenMetadata:       []*types.TokenIdsActionPermission{},
				CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
				CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
					{
						PermanentlyPermittedTimes: GetFullUintRanges(),
					},
				},
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Get user balance with default approvals applied (this will set versions correctly)
	userBalance, _, _ := suite.app.TokenizationKeeper.GetBalanceOrApplyDefault(suite.ctx, collection, bob)

	// Get the correct version for the approval
	version, found := suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(collection.CollectionId, "incoming", bob, "test-incoming-offline"))
	if !found {
		version = sdkmath.NewUint(0)
	}

	// Test 1: Transfer during offline hours (10am UTC) - should be denied
	testTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking := &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	err = DeductUserIncomingApprovals(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		userBalance,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-incoming-offline",
				ApprovalLevel:   "incoming",
				ApproverAddress: bob,
				Version:         version,
			},
		},
		false,
		true, // Only check prioritized incoming approvals
		false,
		nil,
		eventTracking,
		nil,
	)

	suite.Require().NotNil(err, "transfer should be denied during offline hours")
	// The error might be wrapped, so check for either the direct message or the wrapped message
	suite.Require().True(
		contains(err.Error(), "alt time check failed") || contains(err.Error(), "transfer denied") || contains(err.Error(), "inadequate approvals"),
		"error should mention alt time check or transfer denial: %s", err.Error(),
	)

	// Test 2: Transfer outside offline hours (8am UTC) - should be allowed
	testTime = time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	// Get the balance again to ensure we have the latest version
	userBalance, _, _ = suite.app.TokenizationKeeper.GetBalanceOrApplyDefault(suite.ctx, collection, bob)
	version, found = suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(collection.CollectionId, "incoming", bob, "test-incoming-offline"))
	if !found {
		version = sdkmath.NewUint(0)
	}

	eventTracking = &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	err = DeductUserIncomingApprovals(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		userBalance,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-incoming-offline",
				ApprovalLevel:   "incoming",
				ApproverAddress: bob,
				Version:         version,
			},
		},
		false,
		true, // Only check prioritized incoming approvals
		false,
		nil,
		eventTracking,
		nil,
	)

	suite.Require().Nil(err, "transfer should be allowed outside offline hours")
}

// TestAltTimeChecks_OutgoingApproval tests that outgoing approvals respect altTimeChecks
func (suite *TestSuite) TestAltTimeChecks_OutgoingApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create address list for alice and bob BEFORE creating the collection
	err := suite.app.TokenizationKeeper.CreateAddressList(suite.ctx, &types.AddressList{
		ListId:    "aliceAndBob",
		Addresses: []string{alice, bob},
		Whitelist: true,
	})
	suite.Require().Nil(err, "error creating address list")

	// Create a collection with default outgoing approval that has offline hours
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			DefaultOutgoingApprovals: []*types.UserOutgoingApproval{
				{
					ToListId:          "aliceAndBob",
					InitiatedByListId: "aliceAndBob",
					TransferTimes:     GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TokenIds:          GetFullUintRanges(),
					ApprovalId:        "test-outgoing-offline",
					ApprovalCriteria: &types.OutgoingApprovalCriteria{
						AltTimeChecks: &types.AltTimeChecks{
							OfflineHours: []*types.UintRange{
								{
									Start: sdkmath.NewUint(9),
									End:   sdkmath.NewUint(17),
								},
							},
						},
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
							AmountTrackerId:        "test-tracker",
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(uint64(math.MaxUint64)),
							AmountTrackerId:              "test-tracker",
						},
					},
				},
			},
			TokensToCreate: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					TokenIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection:         []*types.ActionPermission{},
				CanUpdateStandards:           []*types.ActionPermission{},
				CanUpdateCustomData:          []*types.ActionPermission{},
				CanDeleteCollection:          []*types.ActionPermission{},
				CanUpdateManager:             []*types.ActionPermission{},
				CanUpdateCollectionMetadata:  []*types.ActionPermission{},
				CanUpdateTokenMetadata:       []*types.TokenIdsActionPermission{},
				CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
				CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
					{
						PermanentlyPermittedTimes: GetFullUintRanges(),
					},
				},
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Get user balance with default approvals applied (this will set versions correctly)
	userBalance, _, _ := suite.app.TokenizationKeeper.GetBalanceOrApplyDefault(suite.ctx, collection, alice)

	// Get the correct version for the approval
	version, found := suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(collection.CollectionId, "outgoing", alice, "test-outgoing-offline"))
	if !found {
		version = sdkmath.NewUint(0)
	}

	// Test 1: Transfer during offline hours (10am UTC) - should be denied
	testTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking := &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	err = DeductUserOutgoingApprovals(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		userBalance,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-outgoing-offline",
				ApprovalLevel:   "outgoing",
				ApproverAddress: alice,
				Version:         version,
			},
		},
		false,
		false,
		true, // Only check prioritized outgoing approvals
		nil,
		eventTracking,
		nil,
	)

	suite.Require().NotNil(err, "transfer should be denied during offline hours")
	// The error might be wrapped, so check for either the direct message or the wrapped message
	suite.Require().True(
		contains(err.Error(), "alt time check failed") || contains(err.Error(), "transfer denied") || contains(err.Error(), "inadequate approvals"),
		"error should mention alt time check or transfer denial: %s", err.Error(),
	)

	// Test 2: Transfer outside offline hours (8am UTC) - should be allowed
	testTime = time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	// Get the balance again to ensure we have the latest version
	userBalance, _, _ = suite.app.TokenizationKeeper.GetBalanceOrApplyDefault(suite.ctx, collection, alice)
	version, found = suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(collection.CollectionId, "outgoing", alice, "test-outgoing-offline"))
	if !found {
		version = sdkmath.NewUint(0)
	}

	eventTracking = &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	err = DeductUserOutgoingApprovals(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		userBalance,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-outgoing-offline",
				ApprovalLevel:   "outgoing",
				ApproverAddress: alice,
				Version:         version,
			},
		},
		false,
		false,
		true, // Only check prioritized outgoing approvals
		nil,
		eventTracking,
		nil,
	)

	suite.Require().Nil(err, "transfer should be allowed outside offline hours")
}

// TestAltTimeChecks_CombinedHoursAndDays tests that both offline hours and days work together
func (suite *TestSuite) TestAltTimeChecks_CombinedHoursAndDays() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create address list (check if it already exists)
	_, err := GetAddressList(suite, wctx, "AllWithoutMint")
	if err != nil {
		err = suite.app.TokenizationKeeper.CreateAddressList(suite.ctx, &types.AddressList{
			ListId:    "AllWithoutMint",
			Addresses: []string{alice, bob},
		})
		suite.Require().Nil(err, "error creating address list")
	}

	// Create a collection with both offline hours (9am-5pm) and offline days (Monday-Friday)
	collectionsToCreate := []*types.MsgNewCollection{
		{
			Creator: bob,
			CollectionApprovals: []*types.CollectionApproval{
				{
					ToListId:          "AllWithoutMint",
					FromListId:        "AllWithoutMint",
					InitiatedByListId: "AllWithoutMint",
					TransferTimes:     GetFullUintRanges(),
					OwnershipTimes:    GetFullUintRanges(),
					TokenIds:          GetFullUintRanges(),
					ApprovalId:        "test-combined",
					ApprovalCriteria: &types.ApprovalCriteria{
						AltTimeChecks: &types.AltTimeChecks{
							OfflineHours: []*types.UintRange{
								{
									Start: sdkmath.NewUint(9),
									End:   sdkmath.NewUint(17),
								},
							},
							OfflineDays: []*types.UintRange{
								{
									Start: sdkmath.NewUint(1), // Monday
									End:   sdkmath.NewUint(5), // Friday
								},
							},
						},
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
							AmountTrackerId:        "test-tracker",
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(uint64(math.MaxUint64)),
							AmountTrackerId:              "test-tracker",
						},
					},
				},
			},
			TokensToCreate: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10),
					TokenIds:       GetFullUintRanges(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection:         []*types.ActionPermission{},
				CanUpdateStandards:           []*types.ActionPermission{},
				CanUpdateCustomData:          []*types.ActionPermission{},
				CanDeleteCollection:          []*types.ActionPermission{},
				CanUpdateManager:             []*types.ActionPermission{},
				CanUpdateCollectionMetadata:  []*types.ActionPermission{},
				CanUpdateTokenMetadata:       []*types.TokenIdsActionPermission{},
				CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
				CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{
					{
						PermanentlyPermittedTimes: GetFullUintRanges(),
					},
				},
			},
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Test 1: Monday at 10am (both offline day and hour) - should be denied
	testTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC) // Monday
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking := &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-combined",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		},
		false,
		false,
		false,
		nil,
		eventTracking,
	)

	suite.Require().NotNil(err, "transfer should be denied on Monday at 10am")
	// The error might be wrapped, so check for either the direct message or the wrapped message
	suite.Require().True(
		contains(err.Error(), "alt time check failed") || contains(err.Error(), "transfer denied") || contains(err.Error(), "inadequate approvals"),
		"error should mention alt time check or transfer denial: %s", err.Error(),
	)

	// Test 2: Saturday at 10am (offline hour but not offline day) - should be denied (hour check takes precedence)
	testTime = time.Date(2024, 1, 6, 10, 0, 0, 0, time.UTC) // Saturday
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking = &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-combined",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		},
		false,
		false,
		false,
		nil,
		eventTracking,
	)

	suite.Require().NotNil(err, "transfer should be denied on Saturday at 10am (offline hour)")

	// Test 3: Monday at 8am (offline day but not offline hour) - should be denied (day check takes precedence)
	testTime = time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC) // Monday
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking = &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-combined",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		},
		false,
		false,
		false,
		nil,
		eventTracking,
	)

	suite.Require().NotNil(err, "transfer should be denied on Monday at 8am (offline day)")

	// Test 4: Saturday at 8am (neither offline day nor hour) - should be allowed
	testTime = time.Date(2024, 1, 6, 8, 0, 0, 0, time.UTC) // Saturday
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking = &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite,
		suite.ctx,
		[]*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		collection,
		GetFullUintRanges(),
		GetFullUintRanges(),
		alice,
		bob,
		alice,
		sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test-combined",
				ApprovalLevel:   "collection",
				ApproverAddress: "",
				Version:         sdkmath.NewUint(0),
			},
		},
		false,
		false,
		false,
		nil,
		eventTracking,
	)

	suite.Require().Nil(err, "transfer should be allowed on Saturday at 8am")
}
