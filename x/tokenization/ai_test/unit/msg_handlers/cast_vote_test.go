package msg_handlers_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type CastVoteTestSuite struct {
	testutil.AITestSuite
}

func TestCastVoteSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(CastVoteTestSuite))
}

// createCollectionWithVotingChallenge creates a collection with a voting challenge approval
func (suite *CastVoteTestSuite) createCollectionWithVotingChallenge(proposalId string, quorumThreshold uint64, voters []*types.Voter) sdkmath.Uint {
	votingChallenge := &types.VotingChallenge{
		ProposalId:      proposalId,
		QuorumThreshold: sdkmath.NewUint(quorumThreshold),
		Voters:          voters,
	}

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	return collectionId
}

// TestCastVote_Success tests casting a valid vote
func (suite *CastVoteTestSuite) TestCastVote_Success() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(50)},
		{Address: suite.Bob, Weight: sdkmath.NewUint(50)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	// Cast vote as Alice
	msg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(100), // 100% yes
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "casting vote should succeed")

	// Verify vote was recorded
	voteKey := keeper.ConstructVotingTrackerKey(
		collectionId,
		"",
		"collection",
		"voting_approval",
		"proposal1",
		suite.Alice,
	)
	vote, found := suite.Keeper.GetVoteFromStore(suite.Ctx, voteKey)
	suite.Require().True(found, "vote should be stored")
	suite.Require().Equal(sdkmath.NewUint(100), vote.YesWeight, "yes weight should be 100")
	suite.Require().Equal(suite.Alice, vote.Voter, "voter should be Alice")
}

// TestCastVote_UpdateExistingVote tests updating an existing vote
func (suite *CastVoteTestSuite) TestCastVote_UpdateExistingVote() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	// Cast initial vote with 30% yes
	initialMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(30),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), initialMsg)
	suite.Require().NoError(err, "initial vote should succeed")

	// Verify initial vote
	voteKey := keeper.ConstructVotingTrackerKey(
		collectionId,
		"",
		"collection",
		"voting_approval",
		"proposal1",
		suite.Alice,
	)
	vote, found := suite.Keeper.GetVoteFromStore(suite.Ctx, voteKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(30), vote.YesWeight, "initial yes weight should be 30")

	// Update vote to 70% yes
	updateMsg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(70),
	}

	_, err = suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "updating vote should succeed")

	// Verify updated vote
	vote, found = suite.Keeper.GetVoteFromStore(suite.Ctx, voteKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(70), vote.YesWeight, "updated yes weight should be 70")
}

// TestCastVote_NonVoterRejected tests that non-voter in voters list is rejected
func (suite *CastVoteTestSuite) TestCastVote_NonVoterRejected() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(50)},
		{Address: suite.Bob, Weight: sdkmath.NewUint(50)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	// Try to cast vote as Charlie (not in voters list)
	msg := &types.MsgCastVote{
		Creator:         suite.Charlie, // Charlie is not in the voters list
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-voter should be rejected")
	suite.Require().Contains(err.Error(), "Voter not found", "error should mention voter not found")
}

// TestCastVote_YesWeightValid tests valid yesWeight values (0-100)
func (suite *CastVoteTestSuite) TestCastVote_YesWeightValid() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	// Test valid weights: 0, 50, 100
	validWeights := []uint64{0, 50, 100}

	for _, weight := range validWeights {
		msg := &types.MsgCastVote{
			Creator:         suite.Alice,
			CollectionId:    collectionId,
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      "voting_approval",
			ProposalId:      "proposal1",
			YesWeight:       sdkmath.NewUint(weight),
		}

		_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "yesWeight %d should be valid", weight)

		// Verify vote
		voteKey := keeper.ConstructVotingTrackerKey(
			collectionId,
			"",
			"collection",
			"voting_approval",
			"proposal1",
			suite.Alice,
		)
		vote, found := suite.Keeper.GetVoteFromStore(suite.Ctx, voteKey)
		suite.Require().True(found)
		suite.Require().Equal(sdkmath.NewUint(weight), vote.YesWeight)
	}
}

// TestCastVote_YesWeightOver100Rejected tests that yesWeight > 100 is rejected
func (suite *CastVoteTestSuite) TestCastVote_YesWeightOver100Rejected() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	// Test invalid weights > 100
	invalidWeights := []uint64{101, 150, 200, 1000}

	for _, weight := range invalidWeights {
		msg := &types.MsgCastVote{
			Creator:         suite.Alice,
			CollectionId:    collectionId,
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      "voting_approval",
			ProposalId:      "proposal1",
			YesWeight:       sdkmath.NewUint(weight),
		}

		_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().Error(err, "yesWeight %d should be rejected", weight)
		suite.Require().Contains(err.Error(), "yesWeight must be between 0 and 100", "error should mention valid range")
	}
}

// TestCastVote_ProposalIdScoping tests that votes are scoped by proposalId
func (suite *CastVoteTestSuite) TestCastVote_ProposalIdScoping() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	// Create collection with voting challenges for two different proposals
	votingChallenge1 := &types.VotingChallenge{
		ProposalId:      "proposal1",
		QuorumThreshold: sdkmath.NewUint(50),
		Voters:          voters,
	}
	votingChallenge2 := &types.VotingChallenge{
		ProposalId:      "proposal2",
		QuorumThreshold: sdkmath.NewUint(50),
		Voters:          voters,
	}

	approval := testutil.GenerateCollectionApproval("voting_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		VotingChallenges:               []*types.VotingChallenge{votingChallenge1, votingChallenge2},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Cast vote for proposal1 with 100% yes
	msg1 := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(100),
	}
	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Cast vote for proposal2 with 25% yes
	msg2 := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal2",
		YesWeight:       sdkmath.NewUint(25),
	}
	_, err = suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Verify proposal1 vote
	voteKey1 := keeper.ConstructVotingTrackerKey(
		collectionId,
		"",
		"collection",
		"voting_approval",
		"proposal1",
		suite.Alice,
	)
	vote1, found := suite.Keeper.GetVoteFromStore(suite.Ctx, voteKey1)
	suite.Require().True(found, "proposal1 vote should be stored")
	suite.Require().Equal(sdkmath.NewUint(100), vote1.YesWeight, "proposal1 yes weight should be 100")

	// Verify proposal2 vote is independent
	voteKey2 := keeper.ConstructVotingTrackerKey(
		collectionId,
		"",
		"collection",
		"voting_approval",
		"proposal2",
		suite.Alice,
	)
	vote2, found := suite.Keeper.GetVoteFromStore(suite.Ctx, voteKey2)
	suite.Require().True(found, "proposal2 vote should be stored")
	suite.Require().Equal(sdkmath.NewUint(25), vote2.YesWeight, "proposal2 yes weight should be 25")
}

// TestCastVote_MultipleVoters tests multiple voters casting votes
func (suite *CastVoteTestSuite) TestCastVote_MultipleVoters() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(40)},
		{Address: suite.Bob, Weight: sdkmath.NewUint(30)},
		{Address: suite.Charlie, Weight: sdkmath.NewUint(30)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	// All three voters cast votes
	voteData := []struct {
		voter     string
		yesWeight uint64
	}{
		{suite.Alice, 100},
		{suite.Bob, 50},
		{suite.Charlie, 0},
	}

	for _, v := range voteData {
		msg := &types.MsgCastVote{
			Creator:         v.voter,
			CollectionId:    collectionId,
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      "voting_approval",
			ProposalId:      "proposal1",
			YesWeight:       sdkmath.NewUint(v.yesWeight),
		}
		_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "vote from %s should succeed", v.voter)
	}

	// Verify all votes
	for _, v := range voteData {
		voteKey := keeper.ConstructVotingTrackerKey(
			collectionId,
			"",
			"collection",
			"voting_approval",
			"proposal1",
			v.voter,
		)
		vote, found := suite.Keeper.GetVoteFromStore(suite.Ctx, voteKey)
		suite.Require().True(found, "vote from %s should be stored", v.voter)
		suite.Require().Equal(sdkmath.NewUint(v.yesWeight), vote.YesWeight, "yes weight for %s should be %d", v.voter, v.yesWeight)
	}
}

// TestCastVote_InvalidCollectionId tests casting vote with invalid collection ID
func (suite *CastVoteTestSuite) TestCastVote_InvalidCollectionId() {
	invalidCollectionId := sdkmath.NewUint(99999)

	msg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    invalidCollectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "should fail with invalid collection ID")
	suite.Require().Contains(err.Error(), "not found", "error should mention collection not found")
}

// TestCastVote_InvalidApprovalId tests casting vote with invalid approval ID
func (suite *CastVoteTestSuite) TestCastVote_InvalidApprovalId() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	msg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "nonexistent_approval", // Invalid approval ID
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "should fail with invalid approval ID")
	suite.Require().Contains(err.Error(), "not found", "error should mention voting challenge not found")
}

// TestCastVote_InvalidProposalId tests casting vote with invalid proposal ID
func (suite *CastVoteTestSuite) TestCastVote_InvalidProposalId() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	msg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "nonexistent_proposal", // Invalid proposal ID
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "should fail with invalid proposal ID")
	suite.Require().Contains(err.Error(), "not found", "error should mention voting challenge not found")
}

// TestCastVote_InvalidApprovalLevel tests casting vote with invalid approval level
func (suite *CastVoteTestSuite) TestCastVote_InvalidApprovalLevel() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	msg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "invalid_level", // Invalid approval level
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(100),
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "should fail with invalid approval level")
}

// TestCastVote_ZeroYesWeight tests casting vote with 0% yes weight (100% no)
func (suite *CastVoteTestSuite) TestCastVote_ZeroYesWeight() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	// Cast vote with 0% yes (100% no)
	msg := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(0), // 100% no vote
	}

	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "casting 0% yes vote should succeed")

	// Verify vote
	voteKey := keeper.ConstructVotingTrackerKey(
		collectionId,
		"",
		"collection",
		"voting_approval",
		"proposal1",
		suite.Alice,
	)
	vote, found := suite.Keeper.GetVoteFromStore(suite.Ctx, voteKey)
	suite.Require().True(found, "vote should be stored")
	suite.Require().Equal(sdkmath.NewUint(0), vote.YesWeight, "yes weight should be 0")
}

// TestCastVote_ExactBoundaryWeights tests exact boundary values (0 and 100)
func (suite *CastVoteTestSuite) TestCastVote_ExactBoundaryWeights() {
	voters := []*types.Voter{
		{Address: suite.Alice, Weight: sdkmath.NewUint(100)},
	}

	collectionId := suite.createCollectionWithVotingChallenge("proposal1", 50, voters)

	// Test exact boundary: 0
	msg0 := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(0),
	}
	_, err := suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg0)
	suite.Require().NoError(err, "yesWeight 0 should be valid")

	// Test exact boundary: 100
	msg100 := &types.MsgCastVote{
		Creator:         suite.Alice,
		CollectionId:    collectionId,
		ApprovalLevel:   "collection",
		ApproverAddress: "",
		ApprovalId:      "voting_approval",
		ProposalId:      "proposal1",
		YesWeight:       sdkmath.NewUint(100),
	}
	_, err = suite.MsgServer.CastVote(sdk.WrapSDKContext(suite.Ctx), msg100)
	suite.Require().NoError(err, "yesWeight 100 should be valid")

	// Verify final vote is 100
	voteKey := keeper.ConstructVotingTrackerKey(
		collectionId,
		"",
		"collection",
		"voting_approval",
		"proposal1",
		suite.Alice,
	)
	vote, found := suite.Keeper.GetVoteFromStore(suite.Ctx, voteKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(100), vote.YesWeight, "final yes weight should be 100")
}
