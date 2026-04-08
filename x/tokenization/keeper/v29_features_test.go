package keeper_test

import (
	"math"
	"time"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ==================== #175: AltTimeChecks — Months, WeeksOfMonth, WeeksOfYear ====================

func (suite *TestSuite) TestAltTimeChecks_OfflineMonths() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
					ApprovalId:        "test-offline-months",
					ApprovalCriteria: &types.ApprovalCriteria{
						AltTimeChecks: &types.AltTimeChecks{
							// Offline during Q2 (April, May, June)
							OfflineMonths: []*types.UintRange{
								{Start: sdkmath.NewUint(4), End: sdkmath.NewUint(6)},
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
			TokensToCreate: []*types.Balance{{Amount: sdkmath.NewUint(10), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection: []*types.ActionPermission{}, CanUpdateStandards: []*types.ActionPermission{},
				CanUpdateCustomData: []*types.ActionPermission{}, CanDeleteCollection: []*types.ActionPermission{},
				CanUpdateManager: []*types.ActionPermission{}, CanUpdateCollectionMetadata: []*types.ActionPermission{},
				CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{}, CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
				CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{{PermanentlyPermittedTimes: GetFullUintRanges()}},
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Test 1: Transfer during offline month (May) — should be denied
	testTime := time.Date(2024, 5, 15, 10, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)

	eventTracking := &keeper.EventTracking{
		ApprovalsUsed: &[]keeper.ApprovalsUsed{},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetFullUintRanges(), GetFullUintRanges(),
		alice, bob, alice, sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "test-offline-months", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		false, false, false, nil, eventTracking,
	)
	suite.Require().NotNil(err, "transfer should be denied during offline month (May)")

	// Test 2: Transfer during online month (January) — should be allowed
	testTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)
	eventTracking = &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetFullUintRanges(), GetFullUintRanges(),
		alice, bob, alice, sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "test-offline-months", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		false, false, false, nil, eventTracking,
	)
	suite.Require().Nil(err, "transfer should be allowed during online month (January)")
}

func (suite *TestSuite) TestAltTimeChecks_OfflineDaysOfMonth() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
					ApprovalId:        "test-offline-wom",
					ApprovalCriteria: &types.ApprovalCriteria{
						AltTimeChecks: &types.AltTimeChecks{
							// Offline during first 7 days of month (days 1-7)
							OfflineDaysOfMonth: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(7)},
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
			TokensToCreate: []*types.Balance{{Amount: sdkmath.NewUint(10), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection: []*types.ActionPermission{}, CanUpdateStandards: []*types.ActionPermission{},
				CanUpdateCustomData: []*types.ActionPermission{}, CanDeleteCollection: []*types.ActionPermission{},
				CanUpdateManager: []*types.ActionPermission{}, CanUpdateCollectionMetadata: []*types.ActionPermission{},
				CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{}, CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
				CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{{PermanentlyPermittedTimes: GetFullUintRanges()}},
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Test 1: Day 3 of month (week 1) — should be denied
	testTime := time.Date(2024, 3, 3, 10, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)
	eventTracking := &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetFullUintRanges(), GetFullUintRanges(),
		alice, bob, alice, sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "test-offline-wom", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		false, false, false, nil, eventTracking,
	)
	suite.Require().NotNil(err, "transfer should be denied during days 1-7 of month")

	// Test 2: Day 15 of month (week 3) — should be allowed
	testTime = time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)
	eventTracking = &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetFullUintRanges(), GetFullUintRanges(),
		alice, bob, alice, sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "test-offline-wom", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		false, false, false, nil, eventTracking,
	)
	suite.Require().Nil(err, "transfer should be allowed during day 15 of month")
}

func (suite *TestSuite) TestAltTimeChecks_OfflineWeeksOfYear() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
					ApprovalId:        "test-offline-woy",
					ApprovalCriteria: &types.ApprovalCriteria{
						AltTimeChecks: &types.AltTimeChecks{
							// Offline during ISO week 1 of year
							OfflineWeeksOfYear: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
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
			TokensToCreate: []*types.Balance{{Amount: sdkmath.NewUint(10), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection: []*types.ActionPermission{}, CanUpdateStandards: []*types.ActionPermission{},
				CanUpdateCustomData: []*types.ActionPermission{}, CanDeleteCollection: []*types.ActionPermission{},
				CanUpdateManager: []*types.ActionPermission{}, CanUpdateCollectionMetadata: []*types.ActionPermission{},
				CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{}, CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
				CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{{PermanentlyPermittedTimes: GetFullUintRanges()}},
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Test 1: Jan 2, 2024 is ISO week 1 — should be denied
	testTime := time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)
	eventTracking := &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetFullUintRanges(), GetFullUintRanges(),
		alice, bob, alice, sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "test-offline-woy", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		false, false, false, nil, eventTracking,
	)
	suite.Require().NotNil(err, "transfer should be denied during ISO week 1")

	// Test 2: March 15, 2024 is ISO week 11 — should be allowed
	testTime = time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)
	eventTracking = &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetFullUintRanges(), GetFullUintRanges(),
		alice, bob, alice, sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "test-offline-woy", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		false, false, false, nil, eventTracking,
	)
	suite.Require().Nil(err, "transfer should be allowed during ISO week 11")
}

func (suite *TestSuite) TestAltTimeChecks_TimezoneOffset() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
					ApprovalId:        "test-timezone",
					ApprovalCriteria: &types.ApprovalCriteria{
						AltTimeChecks: &types.AltTimeChecks{
							// Offline during hours 9-17 (business hours) in EST (UTC-5)
							OfflineHours: []*types.UintRange{
								{Start: sdkmath.NewUint(9), End: sdkmath.NewUint(17)},
							},
							TimezoneOffsetMinutes: sdkmath.NewUint(300), // 300 minutes = 5 hours
						TimezoneOffsetNegative: true,                // negative = west of UTC (EST = UTC-5)
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
			TokensToCreate: []*types.Balance{{Amount: sdkmath.NewUint(10), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			Permissions: &types.CollectionPermissions{
				CanArchiveCollection: []*types.ActionPermission{}, CanUpdateStandards: []*types.ActionPermission{},
				CanUpdateCustomData: []*types.ActionPermission{}, CanDeleteCollection: []*types.ActionPermission{},
				CanUpdateManager: []*types.ActionPermission{}, CanUpdateCollectionMetadata: []*types.ActionPermission{},
				CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{}, CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
				CanUpdateValidTokenIds: []*types.TokenIdsActionPermission{{PermanentlyPermittedTimes: GetFullUintRanges()}},
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err)

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Test 1: 15:00 UTC = 10:00 EST (within offline hours 9-17 EST) — should be denied
	testTime := time.Date(2024, 6, 15, 15, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)
	eventTracking := &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetFullUintRanges(), GetFullUintRanges(),
		alice, bob, alice, sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "test-timezone", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		false, false, false, nil, eventTracking,
	)
	suite.Require().NotNil(err, "15:00 UTC = 10:00 EST should be denied (within 9-17 EST offline)")

	// Test 2: 13:00 UTC = 08:00 EST (outside offline hours 9-17 EST) — should be allowed
	testTime = time.Date(2024, 6, 15, 13, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)
	eventTracking = &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetFullUintRanges(), GetFullUintRanges(),
		alice, bob, alice, sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "test-timezone", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		false, false, false, nil, eventTracking,
	)
	suite.Require().Nil(err, "13:00 UTC = 08:00 EST should be allowed (outside 9-17 EST offline)")

	// Test 3: 23:00 UTC = 18:00 EST (outside offline hours 9-17 EST) — should be allowed
	testTime = time.Date(2024, 6, 15, 23, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(testTime)
	eventTracking = &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetFullUintRanges(), GetFullUintRanges(),
		alice, bob, alice, sdkmath.NewUint(1),
		[]*types.MerkleProof{},
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "test-timezone", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		false, false, false, nil, eventTracking,
	)
	suite.Require().Nil(err, "23:00 UTC = 18:00 EST should be allowed (outside 9-17 EST offline)")
}

// ==================== #173: Voting Challenge Reset + Delay ====================

func (suite *TestSuite) TestVotingChallenge_ResetAfterExecution() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	votingChallenge := &types.VotingChallenge{
		ProposalId:          "vault-proposal",
		QuorumThreshold:     sdkmath.NewUint(50),
		Voters:              []*types.Voter{{Address: alice, Weight: sdkmath.NewUint(100)}, {Address: bob, Weight: sdkmath.NewUint(100)}},
		ResetAfterExecution: true,
		DelayAfterQuorum:    sdkmath.NewUint(0),
	}

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	if collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts != nil {
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(100)
	}

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId: "AllWithoutMint", FromListId: "Mint", InitiatedByListId: "AllWithoutMint",
		TransferTimes: GetFullUintRanges(), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges(),
		ApprovalId: "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1000), AmountTrackerId: "mint-tracker"},
			ApprovalAmounts: &types.ApprovalAmounts{PerFromAddressApprovalAmount: sdkmath.NewUint(1000), AmountTrackerId: "mint-tracker"},
			OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint tokens to bob
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator: bob, CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From: "Mint", ToAddresses: []string{bob},
			Balances: []*types.Balance{{Amount: sdkmath.NewUint(10), TokenIds: GetTopHalfUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{{ApprovalId: "mint-test", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		}},
	})
	suite.Require().NoError(err)

	// Cast votes for first transfer
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "vault-proposal", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "vault-proposal", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// First transfer should succeed
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator: alice, CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From: bob, ToAddresses: []string{alice},
			Balances: []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetTopHalfUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
		}},
	})
	suite.Require().NoError(err, "first transfer should succeed with votes")

	// After reset, votes should be cleared — second transfer should fail without new votes
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator: alice, CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From: bob, ToAddresses: []string{alice},
			Balances: []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetTopHalfUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
		}},
	})
	suite.Require().Error(err, "second transfer should fail — votes were reset after first execution")

	// Verify votes are actually deleted
	voteKey := keeper.ConstructVotingTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "vault-proposal", alice)
	_, found := suite.app.TokenizationKeeper.GetVoteFromStore(suite.ctx, voteKey)
	suite.Require().False(found, "Alice's vote should be deleted after reset")
}

func (suite *TestSuite) TestVotingChallenge_DelayAfterQuorum() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	delayMs := uint64(172800000) // 48 hours in ms

	votingChallenge := &types.VotingChallenge{
		ProposalId:          "delay-proposal",
		QuorumThreshold:     sdkmath.NewUint(50),
		Voters:              []*types.Voter{{Address: alice, Weight: sdkmath.NewUint(100)}, {Address: bob, Weight: sdkmath.NewUint(100)}},
		ResetAfterExecution: true,
		DelayAfterQuorum:    sdkmath.NewUint(delayMs),
	}

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	if collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts != nil {
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(100)
	}

	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId: "AllWithoutMint", FromListId: "Mint", InitiatedByListId: "AllWithoutMint",
		TransferTimes: GetFullUintRanges(), TokenIds: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges(),
		ApprovalId: "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{OverallMaxNumTransfers: sdkmath.NewUint(1000), AmountTrackerId: "mint-tracker"},
			ApprovalAmounts: &types.ApprovalAmounts{PerFromAddressApprovalAmount: sdkmath.NewUint(1000), AmountTrackerId: "mint-tracker"},
			OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	baseTime := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(baseTime)
	wctx = sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint tokens to bob
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator: bob, CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From: "Mint", ToAddresses: []string{bob},
			Balances: []*types.Balance{{Amount: sdkmath.NewUint(10), TokenIds: GetTopHalfUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{{ApprovalId: "mint-test", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)}},
		}},
	})
	suite.Require().NoError(err)

	// Cast votes to reach quorum
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "delay-proposal", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "delay-proposal", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Try transfer immediately — should fail (delay not elapsed)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator: alice, CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From: bob, ToAddresses: []string{alice},
			Balances: []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetTopHalfUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
		}},
	})
	suite.Require().Error(err, "transfer should fail — 48h delay not elapsed")

	// Advance time by 24 hours — still not enough
	suite.ctx = suite.ctx.WithBlockTime(baseTime.Add(24 * time.Hour))
	wctx = sdk.WrapSDKContext(suite.ctx)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator: alice, CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From: bob, ToAddresses: []string{alice},
			Balances: []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetTopHalfUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
		}},
	})
	suite.Require().Error(err, "transfer should fail — only 24h elapsed, need 48h")

	// Advance time to 49 hours — delay should have elapsed
	suite.ctx = suite.ctx.WithBlockTime(baseTime.Add(49 * time.Hour))
	wctx = sdk.WrapSDKContext(suite.ctx)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator: alice, CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From: bob, ToAddresses: []string{alice},
			Balances: []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetTopHalfUintRanges(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
		}},
	})
	suite.Require().NoError(err, "transfer should succeed after 48h delay")
}

func (suite *TestSuite) TestVotingChallenge_VoteTimestamp() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	votingChallenge := &types.VotingChallenge{
		ProposalId:          "timestamp-proposal",
		QuorumThreshold:     sdkmath.NewUint(50),
		Voters:              []*types.Voter{{Address: alice, Weight: sdkmath.NewUint(100)}},
		ResetAfterExecution: false,
		DelayAfterQuorum:    sdkmath.NewUint(0),
	}

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Cast vote at a specific time
	voteTime := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(voteTime)
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "timestamp-proposal", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Verify the vote has a timestamp
	voteKey := keeper.ConstructVotingTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "timestamp-proposal", alice)
	vote, found := suite.app.TokenizationKeeper.GetVoteFromStore(suite.ctx, voteKey)
	suite.Require().True(found)
	expectedTimestamp := sdkmath.NewUint(uint64(voteTime.UnixMilli()))
	suite.Require().True(vote.VotedAt.Equal(expectedTimestamp), "vote timestamp should match block time: got %s, want %s", vote.VotedAt.String(), expectedTimestamp.String())
}

func (suite *TestSuite) TestVotingChallenge_QuorumDropResetsDelay() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	delayMs := uint64(3600000) // 1 hour

	votingChallenge := &types.VotingChallenge{
		ProposalId:          "quorum-drop-proposal",
		QuorumThreshold:     sdkmath.NewUint(100), // need both voters
		Voters:              []*types.Voter{{Address: alice, Weight: sdkmath.NewUint(100)}, {Address: bob, Weight: sdkmath.NewUint(100)}},
		ResetAfterExecution: true,
		DelayAfterQuorum:    sdkmath.NewUint(delayMs),
	}

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	baseTime := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(baseTime)
	wctx = sdk.WrapSDKContext(suite.ctx)

	// Both vote — quorum reached
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "quorum-drop-proposal", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "quorum-drop-proposal", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Verify quorum tracker was set
	trackerKey := keeper.ConstructVotingChallengeTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "quorum-drop-proposal")
	tracker, found := suite.app.TokenizationKeeper.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().False(tracker.QuorumReachedTimestamp.IsZero(), "quorum timestamp should be set")

	// Bob removes his vote (votes 0% yes) — quorum drops
	suite.ctx = suite.ctx.WithBlockTime(baseTime.Add(30 * time.Minute))
	wctx = sdk.WrapSDKContext(suite.ctx)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "quorum-drop-proposal", sdkmath.NewUint(0))
	suite.Require().NoError(err)

	// Verify quorum timestamp was cleared
	tracker, found = suite.app.TokenizationKeeper.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().True(tracker.QuorumReachedTimestamp.IsZero(), "quorum timestamp should be cleared when quorum drops")
}
