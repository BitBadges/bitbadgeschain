package keeper_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Helper function to create a voting challenge
func createVotingChallenge(proposalId string, quorumThreshold sdkmath.Uint, voters []*types.Voter) *types.VotingChallenge {
	return &types.VotingChallenge{
		ProposalId:      proposalId,
		QuorumThreshold: quorumThreshold,
		Voters:          voters,
		Uri:             "",
		CustomData:      "",
	}
}

// Helper function to cast a vote
func castVote(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, approvalLevel string, approverAddress string, approvalId string, proposalId string, yesWeight sdkmath.Uint) error {
	msg := &types.MsgCastVote{
		Creator:         creator,
		CollectionId:    collectionId,
		ApprovalLevel:   approvalLevel,
		ApproverAddress: approverAddress,
		ApprovalId:      approvalId,
		ProposalId:      proposalId,
		YesWeight:       yesWeight,
	}
	return msg.ValidateBasic()
}

// Helper function to cast a vote and execute it
func castVoteAndExecute(suite *TestSuite, ctx context.Context, creator string, collectionId sdkmath.Uint, approvalLevel string, approverAddress string, approvalId string, proposalId string, yesWeight sdkmath.Uint) error {
	msg := &types.MsgCastVote{
		Creator:         creator,
		CollectionId:    collectionId,
		ApprovalLevel:   approvalLevel,
		ApproverAddress: approverAddress,
		ApprovalId:      approvalId,
		ProposalId:      proposalId,
		YesWeight:       yesWeight,
	}
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	_, err := suite.msgServer.CastVote(ctx, msg)
	return err
}

// Helper function to set incoming approval
func SetIncomingApproval(suite *TestSuite, ctx context.Context, msg *types.MsgSetIncomingApproval) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	_, err := suite.msgServer.SetIncomingApproval(ctx, msg)
	return err
}

// Helper function to set outgoing approval
func SetOutgoingApproval(suite *TestSuite, ctx context.Context, msg *types.MsgSetOutgoingApproval) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	_, err := suite.msgServer.SetOutgoingApproval(ctx, msg)
	return err
}

// TestVotingChallenge_ValidVotes tests successful transfer with valid votes meeting threshold
func (suite *TestSuite) TestVotingChallenge_ValidVotes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Cast votes - alice votes 100% yes, bob votes 100% yes
	// Total weight: 200, yes weight: 200, percentage: 100% >= 50% threshold
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err, "Alice should be able to cast vote")

	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err, "Bob should be able to cast vote")

	// Verify votes are stored
	voteKey := keeper.ConstructVotingTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "proposal-1", alice)
	vote, found := suite.app.BadgesKeeper.GetVoteFromStore(suite.ctx, voteKey)
	suite.Require().True(found, "Alice's vote should be stored")
	suite.Require().Equal(sdkmath.NewUint(100), vote.YesWeight, "Alice's vote should be 100% yes")

	// Now try to transfer - should succeed because threshold is met
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with valid votes meeting threshold")
}

// TestVotingChallenge_NoVotes tests failure when no votes are cast
func (suite *TestSuite) TestVotingChallenge_NoVotes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId:  "AllWithoutMint",
		TransferTimes:      GetFullUintRanges(),
		TokenIds:           GetFullUintRanges(),
		OwnershipTimes:     GetFullUintRanges(),
		ApprovalId:         "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Try to transfer without casting votes - should fail
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Transfer should fail when no votes are cast")
}

// TestVotingChallenge_InsufficientVotes tests failure when threshold not met
func (suite *TestSuite) TestVotingChallenge_InsufficientVotes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge requiring 50% threshold
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Cast vote - only alice votes 100% yes
	// Total weight: 200, yes weight: 100, percentage: 50% (exactly threshold, should pass)
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err, "Alice should be able to cast vote")

	// Try to transfer - should succeed because 50% exactly meets threshold
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with exactly threshold percentage")

	// Now test with insufficient votes - alice votes 30% yes
	// Total weight: 200, yes weight: 30, percentage: 15% < 50% threshold
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(30))
	suite.Require().NoError(err, "Alice should be able to update vote")

	// Try to transfer - should fail because threshold not met
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Transfer should fail when threshold not met")
}

// TestVotingChallenge_WeightedVotes tests 50/50, 70/30 vote splits
func (suite *TestSuite) TestVotingChallenge_WeightedVotes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	// Increase approval amount limit to allow multiple transfers
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(1000)

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Test 50/50 split - should meet 50% threshold
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(50))
	suite.Require().NoError(err, "Alice should be able to cast 50/50 vote")

	// Transfer should succeed
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with 50/50 vote")

	// Mint more badges to bob for the second test
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting more badges to bob")

	// Test 70/30 split - should meet 50% threshold
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(70))
	suite.Require().NoError(err, "Alice should be able to cast 70/30 vote")

	// Transfer should succeed
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with 70/30 vote")
}

// TestVotingChallenge_VoteUpdate tests that votes can be updated
func (suite *TestSuite) TestVotingChallenge_VoteUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Cast initial vote
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(30))
	suite.Require().NoError(err, "Alice should be able to cast initial vote")

	// Verify initial vote
	voteKey := keeper.ConstructVotingTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "proposal-1", alice)
	vote, found := suite.app.BadgesKeeper.GetVoteFromStore(suite.ctx, voteKey)
	suite.Require().True(found, "Vote should be stored")
	suite.Require().Equal(sdkmath.NewUint(30), vote.YesWeight, "Initial vote should be 30%")

	// Update vote
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(70))
	suite.Require().NoError(err, "Alice should be able to update vote")

	// Verify updated vote
	vote, found = suite.app.BadgesKeeper.GetVoteFromStore(suite.ctx, voteKey)
	suite.Require().True(found, "Vote should still be stored")
	suite.Require().Equal(sdkmath.NewUint(70), vote.YesWeight, "Updated vote should be 70%")
}

// TestVotingChallenge_MultipleVoters tests multiple voters casting votes
func (suite *TestSuite) TestVotingChallenge_MultipleVoters() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge with 3 voters
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
			{Address: charlie, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Cast votes from all three voters
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err, "Alice should be able to cast vote")

	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(50))
	suite.Require().NoError(err, "Bob should be able to cast vote")

	err = castVoteAndExecute(suite, wctx, charlie, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(0))
	suite.Require().NoError(err, "Charlie should be able to cast vote")

	// Verify all votes are stored
	aliceVoteKey := keeper.ConstructVotingTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "proposal-1", alice)
	bobVoteKey := keeper.ConstructVotingTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "proposal-1", bob)
	charlieVoteKey := keeper.ConstructVotingTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "proposal-1", charlie)

	aliceVote, found := suite.app.BadgesKeeper.GetVoteFromStore(suite.ctx, aliceVoteKey)
	suite.Require().True(found, "Alice's vote should be stored")
	suite.Require().Equal(sdkmath.NewUint(100), aliceVote.YesWeight)

	bobVote, found := suite.app.BadgesKeeper.GetVoteFromStore(suite.ctx, bobVoteKey)
	suite.Require().True(found, "Bob's vote should be stored")
	suite.Require().Equal(sdkmath.NewUint(50), bobVote.YesWeight)

	charlieVote, found := suite.app.BadgesKeeper.GetVoteFromStore(suite.ctx, charlieVoteKey)
	suite.Require().True(found, "Charlie's vote should be stored")
	suite.Require().Equal(sdkmath.NewUint(0), charlieVote.YesWeight)

	// Calculate: alice 100% of 100 = 100, bob 50% of 100 = 50, charlie 0% of 100 = 0
	// Total yes weight: 150, total possible: 300, percentage: 50% (exactly threshold)
	// Transfer should succeed
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with multiple voters meeting threshold")
}

// TestVotingChallenge_ZeroThreshold tests with 0% threshold (should always pass)
func (suite *TestSuite) TestVotingChallenge_ZeroThreshold() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge with 0% threshold
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(0), // 0% threshold - should always pass
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Transfer should succeed even with no votes (0% threshold)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with 0% threshold even with no votes")
}

// TestVotingChallenge_HundredPercentThreshold tests with 100% threshold (all must vote yes)
func (suite *TestSuite) TestVotingChallenge_HundredPercentThreshold() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge with 100% threshold
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(100), // 100% threshold - all must vote yes
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Only alice votes 100% yes - should fail (need both)
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err, "Alice should be able to cast vote")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Transfer should fail with only one voter when 100% threshold required")

	// Both vote 100% yes - should succeed
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err, "Bob should be able to cast vote")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed when all voters vote 100% yes")
}

// TestMsgCastVote_InvalidVoter tests vote casting by non-voter
func (suite *TestSuite) TestMsgCastVote_InvalidVoter() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50),
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			// Note: bob is NOT in the voters list
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to cast vote as bob (not in voters list) - should fail
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().Error(err, "Non-voter should not be able to cast vote")
}

// TestMsgCastVote_InvalidYesWeight tests vote casting with invalid yesWeight
func (suite *TestSuite) TestMsgCastVote_InvalidYesWeight() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50),
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to cast vote with yesWeight > 100 - should fail
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(101))
	suite.Require().Error(err, "yesWeight > 100 should be rejected")
}

// TestVotingChallenge_NonVoterVote tests rejection of vote from non-voter
func (suite *TestSuite) TestVotingChallenge_NonVoterVote() {
	// This is already tested in TestMsgCastVote_InvalidVoter
	suite.T().Skip("Covered by TestMsgCastVote_InvalidVoter")
}

// TestVotingChallenge_WrongProposalId tests rejection of vote with wrong proposal ID
func (suite *TestSuite) TestVotingChallenge_WrongProposalId() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50),
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to cast vote with wrong proposal ID - should fail
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "wrong-proposal", sdkmath.NewUint(100))
	suite.Require().Error(err, "Vote with wrong proposal ID should be rejected")
}

// TestQueryVote tests querying a specific vote
func (suite *TestSuite) TestQueryVote() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50),
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Cast vote
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(70))
	suite.Require().NoError(err)

	// Query the vote
	res, err := suite.app.BadgesKeeper.GetVote(wctx, &types.QueryGetVoteRequest{
		CollectionId:    "1",
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "test",
		ProposalId:      "proposal-1",
		VoterAddress:    alice,
	})
	suite.Require().NoError(err, "Query should succeed")
	suite.Require().NotNil(res.Vote, "Vote should be returned")
	suite.Require().Equal(sdkmath.NewUint(70), res.Vote.YesWeight, "Vote yesWeight should be 70")
	suite.Require().Equal(alice, res.Vote.Voter, "Vote voter should be alice")
	suite.Require().Equal("proposal-1", res.Vote.ProposalId, "Vote proposalId should match")
}

// TestQueryVotes tests querying all votes for a proposal
func (suite *TestSuite) TestQueryVotes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50),
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Cast votes from both voters
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(50))
	suite.Require().NoError(err)

	// Query all votes for the proposal
	res, err := suite.app.BadgesKeeper.GetVotes(wctx, &types.QueryGetVotesRequest{
		CollectionId:    "1",
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "test",
		ProposalId:      "proposal-1",
	})
	suite.Require().NoError(err, "Query should succeed")
	suite.Require().Equal(2, len(res.Votes), "Should return 2 votes")
}

// TestVotingChallenge_MultipleChallenges tests multiple voting challenges in same approval (all must pass)
func (suite *TestSuite) TestVotingChallenge_MultipleChallenges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with two voting challenges
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge1 := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	votingChallenge2 := createVotingChallenge(
		"proposal-2",
		sdkmath.NewUint(60), // 60% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
			{Address: charlie, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge1, votingChallenge2}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Cast votes for proposal-1: alice 100%, bob 100% = 100% >= 50% ✓
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Cast votes for proposal-2: alice 100%, bob 100%, charlie 0% = 66.67% >= 60% ✓
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-2", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-2", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, charlie, sdkmath.NewUint(1), "collection", "", "test", "proposal-2", sdkmath.NewUint(0))
	suite.Require().NoError(err)

	// Transfer should succeed - both challenges met
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed when all challenges are met")

	// Now test failure: proposal-2 doesn't meet threshold
	// Update charlie's vote to make proposal-2 fail: alice 100%, bob 50%, charlie 0% = 50% < 60%
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-2", sdkmath.NewUint(50))
	suite.Require().NoError(err)

	// Mint more badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting more badges to bob")

	// Transfer should fail - proposal-2 threshold not met
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Transfer should fail when one challenge threshold not met")
}

// TestVotingChallenge_IncomingApprovalLevel tests voting challenges with incoming approval level
func (suite *TestSuite) TestVotingChallenge_IncomingApprovalLevel() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with mint approval
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Set incoming approval with voting challenge for alice
	votingChallenge := createVotingChallenge(
		"incoming-proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	err = SetIncomingApproval(suite, wctx, &types.MsgSetIncomingApproval{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Approval: &types.UserIncomingApproval{
			ApprovalId:      "incoming-test",
			FromListId:      "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:   GetFullUintRanges(),
			TokenIds:        GetFullUintRanges(),
			OwnershipTimes:  GetFullUintRanges(),
			ApprovalCriteria: &types.IncomingApprovalCriteria{
				VotingChallenges: []*types.VotingChallenge{votingChallenge},
				SenderChecks:     &types.AddressChecks{},
				InitiatorChecks:  &types.AddressChecks{},
			},
		},
	})
	suite.Require().NoError(err)

	// Cast votes for incoming approval
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "incoming", alice, "incoming-test", "incoming-proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "incoming", alice, "incoming-test", "incoming-proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Transfer should succeed
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
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "incoming-test",
						ApprovalLevel:   "incoming",
						ApproverAddress: alice,
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with incoming approval votes")
}

// TestVotingChallenge_OutgoingApprovalLevel tests voting challenges with outgoing approval level
func (suite *TestSuite) TestVotingChallenge_OutgoingApprovalLevel() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with mint approval
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Set outgoing approval with voting challenge for bob
	votingChallenge := createVotingChallenge(
		"outgoing-proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	err = SetOutgoingApproval(suite, wctx, &types.MsgSetOutgoingApproval{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Approval: &types.UserOutgoingApproval{
			ApprovalId:      "outgoing-test",
			ToListId:        "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:   GetFullUintRanges(),
			TokenIds:        GetFullUintRanges(),
			OwnershipTimes:  GetFullUintRanges(),
			ApprovalCriteria: &types.OutgoingApprovalCriteria{
				VotingChallenges: []*types.VotingChallenge{votingChallenge},
				RecipientChecks:  &types.AddressChecks{},
				InitiatorChecks: &types.AddressChecks{},
			},
		},
	})
	suite.Require().NoError(err)

	// Cast votes for outgoing approval
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "outgoing", bob, "outgoing-test", "outgoing-proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "outgoing", bob, "outgoing-test", "outgoing-proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Transfer should succeed
	// Note: For outgoing approvals, we need to ensure collection approvals allow the transfer
	// The outgoing approval will be checked, but collection approval must also match
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
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with outgoing approval votes")
}

// TestVotingChallenge_UnequalWeights tests voting with voters having different weights
func (suite *TestSuite) TestVotingChallenge_UnequalWeights() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge with unequal weights
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},  // 50% of total
			{Address: bob, Weight: sdkmath.NewUint(50)},     // 25% of total
			{Address: charlie, Weight: sdkmath.NewUint(50)}, // 25% of total
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	if collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts == nil {
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts = &types.ApprovalAmounts{}
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(1000)

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Test: alice (100 weight, 100% yes) = 100 yes weight
	// Total possible: 200, yes: 100 = 50% (exactly threshold)
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Transfer should succeed
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with alice's vote meeting threshold")

	// Mint more badges
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting more badges to bob")

	// Test: alice 50%, bob 100%, charlie 0%
	// alice: 100 * 50% = 50, bob: 50 * 100% = 50, charlie: 50 * 0% = 0
	// Total yes: 100, total possible: 200, percentage: 50% (exactly threshold)
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(50))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, charlie, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(0))
	suite.Require().NoError(err)

	// Transfer should succeed
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed with unequal weights meeting threshold")
}

// TestVotingChallenge_ExactlyAtThreshold tests edge case where percentage exactly equals threshold
func (suite *TestSuite) TestVotingChallenge_ExactlyAtThreshold() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge requiring exactly 50%
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Exactly 50%: alice 100% yes = 100/200 = 50%
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Transfer should succeed (exactly at threshold)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed when exactly at threshold")

	// Mint more badges
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting more badges to bob")

	// Just below 50%: alice 99% yes = 99/200 = 49.5% < 50%
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(99))
	suite.Require().NoError(err)

	// Transfer should fail (just below threshold)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Transfer should fail when just below threshold")
}

// TestVotingChallenge_ZeroYesWeight tests edge case where voter votes 0% yes (all no)
func (suite *TestSuite) TestVotingChallenge_ZeroYesWeight() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// alice votes 0% yes (all no), bob votes 100% yes
	// Total yes: 0 + 100 = 100, total possible: 200, percentage: 50% (exactly threshold)
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(0))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Transfer should succeed (bob's vote alone meets threshold)
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed even with one voter voting 0% yes")
}

// TestVotingChallenge_WrongApproverAddress tests vote casting with wrong approver address
func (suite *TestSuite) TestVotingChallenge_WrongApproverAddress() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with mint approval
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Set incoming approval for alice
	votingChallenge := createVotingChallenge(
		"incoming-proposal-1",
		sdkmath.NewUint(50),
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
		},
	)
	err = SetIncomingApproval(suite, wctx, &types.MsgSetIncomingApproval{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Approval: &types.UserIncomingApproval{
			ApprovalId:      "incoming-test",
			FromListId:      "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:   GetFullUintRanges(),
			TokenIds:        GetFullUintRanges(),
			OwnershipTimes:  GetFullUintRanges(),
			ApprovalCriteria: &types.IncomingApprovalCriteria{
				VotingChallenges: []*types.VotingChallenge{votingChallenge},
				SenderChecks:     &types.AddressChecks{},
				InitiatorChecks:  &types.AddressChecks{},
			},
		},
	})
	suite.Require().NoError(err)

	// Try to cast vote with wrong approver address (bob instead of alice) - should fail
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "incoming", bob, "incoming-test", "incoming-proposal-1", sdkmath.NewUint(100))
	suite.Require().Error(err, "Vote with wrong approver address should be rejected")
}

// TestVotingChallenge_VotePersistence tests that votes persist across multiple transfers
func (suite *TestSuite) TestVotingChallenge_VotePersistence() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge
	collectionsToCreate := GetCollectionsToCreate()
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50), // 50% threshold
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(1000)

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Cast votes once
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Perform multiple transfers - votes should persist
	for i := 0; i < 3; i++ {
		// Mint badges to bob
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
							TokenIds:       GetTopHalfUintRanges(),
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
		suite.Require().NoError(err, "Error minting badges to bob")

		// Transfer should succeed using same votes
		err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
			Creator:      alice,
			CollectionId: sdkmath.NewUint(1),
			Transfers: []*types.Transfer{
				{
					From:        bob,
					ToAddresses: []string{alice},
					Balances: []*types.Balance{
						{
							Amount:         sdkmath.NewUint(1),
							TokenIds:       GetTopHalfUintRanges(),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
				},
			},
		})
		suite.Require().NoError(err, "Transfer %d should succeed with persistent votes", i+1)
	}

	// Verify votes are still stored
	aliceVoteKey := keeper.ConstructVotingTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "proposal-1", alice)
	bobVoteKey := keeper.ConstructVotingTrackerKey(sdkmath.NewUint(1), "", "collection", "test", "proposal-1", bob)

	aliceVote, found := suite.app.BadgesKeeper.GetVoteFromStore(suite.ctx, aliceVoteKey)
	suite.Require().True(found, "Alice's vote should still be stored")
	suite.Require().Equal(sdkmath.NewUint(100), aliceVote.YesWeight, "Alice's vote should be unchanged")

	bobVote, found := suite.app.BadgesKeeper.GetVoteFromStore(suite.ctx, bobVoteKey)
	suite.Require().True(found, "Bob's vote should still be stored")
	suite.Require().Equal(sdkmath.NewUint(100), bobVote.YesWeight, "Bob's vote should be unchanged")
}

// TestVotingChallenge_WithOtherCriteria tests integration with other approval criteria
func (suite *TestSuite) TestVotingChallenge_WithOtherCriteria() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create collection with voting challenge AND merkle challenge
	collectionsToCreate := GetCollectionsToCreate()

	// Create Merkle tree
	aliceLeaf := "-" + alice + "-1-0-0"
	leafs := [][]byte{[]byte(aliceLeaf)}
	leafHashes := make([][]byte, len(leafs))
	for i, leaf := range leafs {
		hash := sha256.Sum256(leaf)
		leafHashes[i] = hash[:]
	}
	rootHash := hex.EncodeToString(leafHashes[0])

	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50),
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.MerkleChallenges = []*types.MerkleChallenge{
		{
			Root:                rootHash,
			ExpectedProofLength: sdkmath.NewUint(0),
			MaxUsesPerLeaf:      sdkmath.NewUint(1), // Must be 1 for non-whitelist tree
			ChallengeTrackerId:  "merkle-challenge-1",
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.VotingChallenges = []*types.VotingChallenge{votingChallenge}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	if collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts == nil {
		collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts = &types.ApprovalAmounts{}
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(1000)

	// Add mint approval
	collectionsToCreate[0].CollectionApprovals = append([]*types.CollectionApproval{{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "mint-test",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-test-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	}}, collectionsToCreate[0].CollectionApprovals...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Mint badges to bob
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
						TokenIds:       GetTopHalfUintRanges(),
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
	suite.Require().NoError(err, "Error minting badges to bob")

	// Cast votes
	err = castVoteAndExecute(suite, wctx, alice, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)
	err = castVoteAndExecute(suite, wctx, bob, sdkmath.NewUint(1), "collection", "", "test", "proposal-1", sdkmath.NewUint(100))
	suite.Require().NoError(err)

	// Transfer should succeed - both voting and merkle challenges must be satisfied
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetTopHalfUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
				MerkleProofs: []*types.MerkleProof{
					{
						Leaf:  aliceLeaf,
						Aunts: []*types.MerklePathItem{},
					},
				},
			},
		},
	})
	suite.Require().NoError(err, "Transfer should succeed when both voting and merkle challenges are satisfied")
}

// TestVotingChallenge_DuplicateVoters tests that duplicate voters are rejected in ValidateBasic
func (suite *TestSuite) TestVotingChallenge_DuplicateVoters() {
	// Test duplicate voters in a single voting challenge
	votingChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50),
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: alice, Weight: sdkmath.NewUint(50)}, // Duplicate!
			{Address: bob, Weight: sdkmath.NewUint(100)},
		},
	)

	err := votingChallenge.ValidateBasic()
	suite.Require().Error(err, "VotingChallenge with duplicate voters should be rejected")
	suite.Require().Contains(err.Error(), "duplicate voter address", "Error should mention duplicate voter address")
	suite.Require().Contains(err.Error(), alice, "Error should mention the duplicate address")

	// Test that valid voting challenge with no duplicates passes
	validChallenge := createVotingChallenge(
		"proposal-1",
		sdkmath.NewUint(50),
		[]*types.Voter{
			{Address: alice, Weight: sdkmath.NewUint(100)},
			{Address: bob, Weight: sdkmath.NewUint(100)},
			{Address: charlie, Weight: sdkmath.NewUint(100)},
		},
	)

	err = validChallenge.ValidateBasic()
	suite.Require().NoError(err, "VotingChallenge with no duplicate voters should pass validation")
}

