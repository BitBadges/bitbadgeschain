package challenges_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type VotingChallengeTestSuite struct {
	testutil.AITestSuite
}

func TestVotingChallengeTestSuite(t *testing.T) {
	suite.Run(t, new(VotingChallengeTestSuite))
}

func (suite *VotingChallengeTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// createVotingChallenge creates a VotingChallenge with the given parameters
func createVotingChallenge(proposalId string, quorumThreshold uint64, voters []*types.Voter) *types.VotingChallenge {
	return &types.VotingChallenge{
		ProposalId:      proposalId,
		QuorumThreshold: sdkmath.NewUint(quorumThreshold),
		Voters:          voters,
		Uri:             "",
		CustomData:      "",
	}
}

// TestVotingChallenge_QuorumMetPasses tests that transfer succeeds when quorum is met
func (suite *VotingChallengeTestSuite) TestVotingChallenge_QuorumMetPasses() {
	// Create voters: Alice (weight 50), Bob (weight 50)
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(50)},
		{Address: suite.Bob, Weight: sdkmath.NewUint(50)},
	}

	// 50% quorum threshold
	votingChallenge := createVotingChallenge("proposal_quorum_met", 50, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Cast vote from Alice (100% yes weight)
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_quorum_met",
		YesWeight:       sdkmath.NewUint(100), // 100% of Alice's weight (50) goes to yes
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().NoError(err, "Alice should be able to cast vote")

	// Bob also votes yes
	voteMsg2 := &types.MsgCastVote{
		Creator:         suite.Bob,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_quorum_met",
		YesWeight:       sdkmath.NewUint(100), // 100% of Bob's weight (50) goes to yes
	}

	_, err = suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg2)
	suite.Require().NoError(err, "Bob should be able to cast vote")

	// Now transfer should succeed (100% yes votes >= 50% threshold)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when quorum is met")
}

// TestVotingChallenge_QuorumNotMetFails tests that transfer fails when quorum is not met
func (suite *VotingChallengeTestSuite) TestVotingChallenge_QuorumNotMetFails() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(50)},
		{Address: suite.Bob, Weight: sdkmath.NewUint(50)},
	}

	// 75% quorum threshold
	votingChallenge := createVotingChallenge("proposal_quorum_not_met", 75, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Only Alice votes (50% of total weight)
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_quorum_not_met",
		YesWeight:       sdkmath.NewUint(100), // 100% yes
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().NoError(err, "Alice should be able to cast vote")

	// Transfer should fail (only 50% yes votes < 75% threshold)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when quorum is not met")
}

// TestVotingChallenge_VoterWeightsCalculatedCorrectly tests that voter weights are calculated correctly
func (suite *VotingChallengeTestSuite) TestVotingChallenge_VoterWeightsCalculatedCorrectly() {
	// Weighted voting: Alice (30), Bob (70)
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(30)},
		{Address: suite.Bob, Weight: sdkmath.NewUint(70)},
	}

	// 60% threshold
	votingChallenge := createVotingChallenge("proposal_weighted", 60, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Alice votes yes (30% of total)
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_weighted",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().NoError(err)

	// Transfer should fail with only Alice's vote (30% < 60%)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail with only 30% yes votes")

	// Bob votes yes (70% of total) - now total is 100%
	voteMsg2 := &types.MsgCastVote{
		Creator:         suite.Bob,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_weighted",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err = suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg2)
	suite.Require().NoError(err)

	// Now transfer should succeed (100% > 60%)
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed with 100% yes votes")
}

// TestVotingChallenge_NonVoterCannotVote tests that non-voters cannot cast votes
func (suite *VotingChallengeTestSuite) TestVotingChallenge_NonVoterCannotVote() {
	// Only Alice is a voter
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	votingChallenge := createVotingChallenge("proposal_non_voter", 50, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Bob (not a voter) tries to vote
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Bob,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_non_voter",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().Error(err, "non-voter should not be able to cast vote")
}

// TestVotingChallenge_VoteUpdateRecalculatesQuorum tests that updating a vote recalculates quorum
func (suite *VotingChallengeTestSuite) TestVotingChallenge_VoteUpdateRecalculatesQuorum() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	// 50% threshold
	votingChallenge := createVotingChallenge("proposal_update", 50, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(20, 1)})

	// Alice votes 100% yes initially
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_update",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().NoError(err)

	// Transfer should succeed (100% > 50%)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "first transfer should succeed with 100% yes votes")

	// Alice updates vote to 0% yes (100% no)
	voteMsg2 := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_update",
		YesWeight:       sdkmath.NewUint(0), // Change to 0% yes
	}

	_, err = suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg2)
	suite.Require().NoError(err)

	// Transfer should now fail (0% < 50%)
	transferMsg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg2)
	suite.Require().Error(err, "second transfer should fail after vote update to 0% yes")
}

// TestVotingChallenge_ProposalIdScoping tests that proposalId scopes vote tracking
func (suite *VotingChallengeTestSuite) TestVotingChallenge_ProposalIdScoping() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	// Create two challenges with different proposal IDs
	votingChallenge1 := createVotingChallenge("proposal_1", 50, voters)
	votingChallenge2 := createVotingChallenge("proposal_2", 50, voters)

	approval1 := testutil.GenerateCollectionApproval("voting_approval_1", "AllWithoutMint", "All")
	approval1.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge1},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	approval2 := testutil.GenerateCollectionApproval("voting_approval_2", "AllWithoutMint", "All")
	approval2.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge2},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval1, approval2})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(30, 1)})

	// Vote for proposal_1 only
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval_1",
		ProposalId:      "proposal_1",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().NoError(err)

	// Transfer using approval_1 should succeed
	transferMsg1 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("voting_approval_1"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg1)
	suite.Require().NoError(err, "transfer using approval_1 should succeed")

	// Transfer using approval_2 should fail (no votes for proposal_2)
	transferMsg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					testutil.GeneratePrioritizedApproval("voting_approval_2"),
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg2)
	suite.Require().Error(err, "transfer using approval_2 should fail - no votes for proposal_2")
}

// TestVotingChallenge_QuorumThreshold0Passes tests that 0% threshold always passes
func (suite *VotingChallengeTestSuite) TestVotingChallenge_QuorumThreshold0Passes() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	// 0% threshold
	votingChallenge := createVotingChallenge("proposal_zero_threshold", 0, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// No votes cast at all - transfer should still succeed with 0% threshold
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed with 0% threshold even without votes")
}

// TestVotingChallenge_QuorumThreshold100RequiresAll tests that 100% threshold requires all voters
func (suite *VotingChallengeTestSuite) TestVotingChallenge_QuorumThreshold100RequiresAll() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(50)},
		{Address: suite.Bob, Weight: sdkmath.NewUint(50)},
	}

	// 100% threshold
	votingChallenge := createVotingChallenge("proposal_100_threshold", 100, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(20, 1)})

	// Only Alice votes yes (50%)
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_100_threshold",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().NoError(err)

	// Transfer should fail (only 50% < 100%)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail with only 50% of 100% required threshold")

	// Bob also votes yes (100% total)
	voteMsg2 := &types.MsgCastVote{
		Creator:         suite.Bob,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_100_threshold",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err = suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg2)
	suite.Require().NoError(err)

	// Now transfer should succeed (100% = 100%)
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed when all voters vote yes")
}

// TestVotingChallenge_PartialYesWeight tests that partial yes weight is calculated correctly
func (suite *VotingChallengeTestSuite) TestVotingChallenge_PartialYesWeight() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	// 40% threshold
	votingChallenge := createVotingChallenge("proposal_partial", 40, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(20, 1)})

	// Alice votes with 30% yes weight (30% of total)
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_partial",
		YesWeight:       sdkmath.NewUint(30), // Only 30% yes
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().NoError(err)

	// Transfer should fail (30% < 40%)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail with 30% yes < 40% threshold")

	// Alice updates to 50% yes weight
	voteMsg2 := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_partial",
		YesWeight:       sdkmath.NewUint(50), // 50% yes
	}

	_, err = suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg2)
	suite.Require().NoError(err)

	// Now transfer should succeed (50% >= 40%)
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed with 50% yes >= 40% threshold")
}

// TestVotingChallenge_InvalidYesWeightRejected tests that yesWeight > 100 is rejected
func (suite *VotingChallengeTestSuite) TestVotingChallenge_InvalidYesWeightRejected() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	votingChallenge := createVotingChallenge("proposal_invalid", 50, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Try to vote with yesWeight > 100
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_invalid",
		YesWeight:       sdkmath.NewUint(150), // Invalid: > 100
	}

	err := voteMsg.ValidateBasic()
	suite.Require().Error(err, "yesWeight > 100 should be rejected in ValidateBasic")
}

// TestVotingChallenge_MultipleVotersPartialYes tests quorum calculation with multiple partial votes
func (suite *VotingChallengeTestSuite) TestVotingChallenge_MultipleVotersPartialYes() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(60)},
		{Address: suite.Bob, Weight: sdkmath.NewUint(40)},
	}

	// 50% threshold
	votingChallenge := createVotingChallenge("proposal_multi_partial", 50, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Alice votes 50% yes (contributes 30 to yes out of 100 total)
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_multi_partial",
		YesWeight:       sdkmath.NewUint(50), // 50% of Alice's 60 weight = 30 yes
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().NoError(err)

	// Bob votes 75% yes (contributes 30 to yes)
	voteMsg2 := &types.MsgCastVote{
		Creator:         suite.Bob,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_multi_partial",
		YesWeight:       sdkmath.NewUint(75), // 75% of Bob's 40 weight = 30 yes
	}

	_, err = suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg2)
	suite.Require().NoError(err)

	// Total yes: 30 (Alice) + 30 (Bob) = 60 out of 100 = 60% >= 50%
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer should succeed with 60% combined yes votes >= 50% threshold")
}

// TestVotingChallenge_NoVotesCastFails tests that transfer fails when no votes are cast
func (suite *VotingChallengeTestSuite) TestVotingChallenge_NoVotesCastFails() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	// 50% threshold - requires at least some votes
	votingChallenge := createVotingChallenge("proposal_no_votes", 50, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// No votes cast - try to transfer
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail when no votes are cast and threshold > 0")
}

// TestVotingChallenge_InvalidApprovalLevelRejected tests that invalid approval level is rejected
func (suite *VotingChallengeTestSuite) TestVotingChallenge_InvalidApprovalLevelRejected() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	votingChallenge := createVotingChallenge("proposal_level", 50, voters)

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Try to vote with invalid approval level
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "invalid_level",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal_level",
		YesWeight:       sdkmath.NewUint(100),
	}

	err := voteMsg.ValidateBasic()
	suite.Require().Error(err, "invalid approval level should be rejected")
}

// TestVotingChallenge_VotingChallengeNotFoundRejected tests voting for non-existent challenge
func (suite *VotingChallengeTestSuite) TestVotingChallenge_VotingChallengeNotFoundRejected() {
	// Create collection without any voting challenges
	approval := testutil.GenerateCollectionApproval("simple_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Try to vote for non-existent challenge
	voteMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "simple_approval",
		ProposalId:      "non_existent_proposal",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), voteMsg)
	suite.Require().Error(err, "voting for non-existent challenge should fail")
}
