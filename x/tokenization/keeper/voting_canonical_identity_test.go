package keeper_test

import (
	"strings"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestV35LegacyVoterCanChangeVote() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	challenge := &types.VotingChallenge{
		ProposalId: "legacy-vote", QuorumThreshold: sdkmath.NewUint(100),
		Voters:           []*types.Voter{{Address: bob, Weight: sdkmath.OneUint()}},
		DelayAfterQuorum: sdkmath.NewUint(100), ResetAfterExecution: true,
	}
	suite.Require().NoError(CreateCollections(suite, wctx, []*types.MsgNewCollection{
		v35MintCollection(&types.ApprovalCriteria{VotingChallenges: []*types.VotingChallenge{challenge}, OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true}, GetFullUintRanges()),
	}))
	collection, err := GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	approval := collection.CollectionApprovals[0]
	upper := strings.ToUpper(bob)
	approval.ApprovalCriteria.VotingChallenges[0].Voters[0].Address = upper
	suite.Require().NoError(k.SetCollectionInStore(suite.ctx, collection, false))
	key := keeper.ConstructVotingTrackerKey(sdkmath.OneUint(), "", "collection", approval.ApprovalId, "legacy-vote", upper)
	suite.Require().NoError(k.SetVoteInStore(suite.ctx, key, &types.VoteProof{ProposalId: "legacy-vote", Voter: upper, YesWeight: sdkmath.NewUint(100), VotedAt: sdkmath.OneUint()}))
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	server := keeper.NewMsgServerImpl(k)
	msg := &types.MsgCastVote{Creator: upper, CollectionId: sdkmath.OneUint(), ApprovalLevel: "collection", ApprovalId: approval.ApprovalId, ProposalId: "legacy-vote", YesWeight: sdkmath.ZeroUint()}
	_, err = server.CastVote(wctx, msg)
	suite.Require().ErrorContains(err, "canonical form")
	msg.Creator = bob
	_, err = server.CastVote(wctx, msg)
	suite.Require().NoError(err)
	vote, found := k.GetVoteFromStore(suite.ctx, key)
	suite.Require().True(found)
	suite.Require().Equal(upper, vote.Voter)
	suite.Require().True(vote.YesWeight.IsZero())
	trackerKey := keeper.ConstructVotingChallengeTrackerKey(sdkmath.OneUint(), "", "collection", approval.ApprovalId, "legacy-vote")
	tracker, found := k.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().True(tracker.QuorumReachedTimestamp.IsZero())
	msg.YesWeight = sdkmath.NewUint(100)
	_, err = server.CastVote(wctx, msg)
	suite.Require().NoError(err)
	tracker, found = k.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(uint64(suite.ctx.BlockTime().UnixMilli())), tracker.QuorumReachedTimestamp)
	msg.Creator = charlie
	_, err = server.CastVote(wctx, msg)
	suite.Require().ErrorContains(err, "Voter not found")
	suite.Require().NoError(k.DeleteAllVotesForProposal(suite.ctx, sdkmath.OneUint(), "", "collection", approval.ApprovalId, "legacy-vote", approval.ApprovalCriteria.VotingChallenges[0].Voters))
	_, found = k.GetVoteFromStore(suite.ctx, key)
	suite.Require().False(found)
}

func (suite *TestSuite) TestV35DualSpellingVotersKeepSeparateWeights() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	challenge := &types.VotingChallenge{
		ProposalId: "two-slots", QuorumThreshold: sdkmath.NewUint(40), DelayAfterQuorum: sdkmath.OneUint(),
		Voters: []*types.Voter{{Address: bob, Weight: sdkmath.OneUint()}, {Address: charlie, Weight: sdkmath.NewUint(3)}},
	}
	suite.Require().NoError(CreateCollections(suite, wctx, []*types.MsgNewCollection{
		v35MintCollection(&types.ApprovalCriteria{VotingChallenges: []*types.VotingChallenge{challenge}, OverridesFromOutgoingApprovals: true, OverridesToIncomingApprovals: true}, GetFullUintRanges()),
	}))
	collection, err := GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	approval := collection.CollectionApprovals[0]
	challenge = approval.ApprovalCriteria.VotingChallenges[0]
	upper := strings.ToUpper(bob)
	challenge.Voters[1].Address = upper
	suite.Require().NoError(k.SetCollectionInStore(suite.ctx, collection, false))
	for i, voter := range challenge.Voters {
		key := keeper.ConstructVotingTrackerKey(collection.CollectionId, "", "collection", approval.ApprovalId, challenge.ProposalId, voter.Address)
		suite.Require().NoError(k.SetVoteInStore(suite.ctx, key, &types.VoteProof{ProposalId: challenge.ProposalId, Voter: voter.Address, YesWeight: sdkmath.NewUint(uint64(25 + 75*i)), VotedAt: sdkmath.OneUint()}))
	}
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	stored, err := GetCollection(suite, wctx, collection.CollectionId)
	suite.Require().NoError(err)
	suite.Require().Equal(challenge.Voters, stored.CollectionApprovals[0].ApprovalCriteria.VotingChallenges[0].Voters)
	for i, voter := range challenge.Voters {
		key := keeper.ConstructVotingTrackerKey(collection.CollectionId, "", "collection", approval.ApprovalId, challenge.ProposalId, voter.Address)
		vote, found := k.GetVoteFromStore(suite.ctx, key)
		suite.Require().True(found)
		suite.Require().Equal(sdkmath.NewUint(uint64(25+75*i)), vote.YesWeight, "migration preserves each existing vote")
	}
	server := keeper.NewMsgServerImpl(k)
	msg := &types.MsgCastVote{Creator: bob, CollectionId: collection.CollectionId, ApprovalLevel: "collection", ApprovalId: approval.ApprovalId, ProposalId: challenge.ProposalId, YesWeight: sdkmath.NewUint(50)}
	_, err = server.CastVote(wctx, msg)
	suite.Require().NoError(err)
	for _, voter := range challenge.Voters {
		key := keeper.ConstructVotingTrackerKey(collection.CollectionId, "", "collection", approval.ApprovalId, challenge.ProposalId, voter.Address)
		vote, found := k.GetVoteFromStore(suite.ctx, key)
		suite.Require().True(found)
		suite.Require().Equal(voter.Address, vote.Voter)
		suite.Require().Equal(sdkmath.NewUint(50), vote.YesWeight)
	}
	trackerKey := keeper.ConstructVotingChallengeTrackerKey(collection.CollectionId, "", "collection", approval.ApprovalId, challenge.ProposalId)
	tracker, found := k.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().True(tracker.QuorumReachedTimestamp.IsZero(), "floor(1*50%%)+floor(3*50%%)=1, not merged-weight 2")
	for _, checker := range k.GetApprovalCriteriaCheckers(approval) {
		if checker.Name() == "VotingChallenges" {
			_, err := checker.Check(suite.ctx, approval, collection, alice, "Mint", alice, "collection", "", nil, nil, "", true)
			suite.Require().ErrorContains(err, "threshold not met")
		}
	}
	msg.YesWeight = sdkmath.NewUint(100)
	_, err = server.CastVote(wctx, msg)
	suite.Require().NoError(err)
	tracker, found = k.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(uint64(suite.ctx.BlockTime().UnixMilli())), tracker.QuorumReachedTimestamp)
}

func (suite *TestSuite) TestV35VotingApproverMigrationCollisionsAndDelay() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.Require().NoError(CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{OverridesFromOutgoingApprovals: true}, GetFullUintRanges())}))
	collection, err := GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	upper := strings.ToUpper(alice)
	challenge := &types.VotingChallenge{ProposalId: "legacy-proposal", QuorumThreshold: sdkmath.NewUint(100), DelayAfterQuorum: sdkmath.NewUint(100),
		Voters: []*types.Voter{{Address: bob, Weight: sdkmath.OneUint()}, {Address: charlie, Weight: sdkmath.OneUint()}},
	}
	approvalID := "user-approval"
	balance, _, err := k.GetBalanceOrApplyDefault(suite.ctx, collection, alice)
	suite.Require().NoError(err)
	balance.OutgoingApprovals = []*types.UserOutgoingApproval{{ApprovalId: approvalID, ToListId: "All", InitiatedByListId: "All", TransferTimes: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges(), TokenIds: GetFullUintRanges(), Version: sdkmath.ZeroUint(), ApprovalCriteria: &types.OutgoingApprovalCriteria{VotingChallenges: []*types.VotingChallenge{challenge}}}}
	suite.Require().NoError(k.SetUserBalanceInStore(suite.ctx, keeper.ConstructBalanceKey(alice, collection.CollectionId), balance, false))
	setVote := func(approver, proposal, voter string, yes, votedAt uint64) *types.VoteProof {
		vote := &types.VoteProof{ProposalId: proposal, Voter: voter, YesWeight: sdkmath.NewUint(yes), VotedAt: sdkmath.NewUint(votedAt)}
		key := keeper.ConstructVotingTrackerKey(collection.CollectionId, approver, "outgoing", approvalID, proposal, voter)
		suite.Require().NoError(k.SetVoteInStore(suite.ctx, key, vote))
		return vote
	}
	setTracker := func(approver, proposal string, timestamp uint64) {
		key := keeper.ConstructVotingChallengeTrackerKey(collection.CollectionId, approver, "outgoing", approvalID, proposal)
		suite.Require().NoError(k.SetVotingChallengeTrackerInStore(suite.ctx, key, &types.VotingChallengeTracker{QuorumReachedTimestamp: sdkmath.NewUint(timestamp)}))
	}
	canonicalVote := setVote(alice, challenge.ProposalId, bob, 100, 10)
	setVote(upper, challenge.ProposalId, bob, 0, 20)
	movedVote := setVote(upper, challenge.ProposalId, charlie, 100, 30)
	setTracker(alice, challenge.ProposalId, 10)
	setTracker(upper, challenge.ProposalId, 20)
	isolatedVote := setVote(upper, "isolated-proposal", bob, 75, 40)
	setTracker(upper, "isolated-proposal", 40)
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	for _, vote := range []*types.VoteProof{canonicalVote, movedVote, isolatedVote} {
		key := keeper.ConstructVotingTrackerKey(collection.CollectionId, alice, "outgoing", approvalID, vote.ProposalId, vote.Voter)
		stored, found := k.GetVoteFromStore(suite.ctx, key)
		suite.Require().True(found)
		suite.Require().Equal(vote, stored)
		oldKey := keeper.ConstructVotingTrackerKey(collection.CollectionId, upper, "outgoing", approvalID, vote.ProposalId, vote.Voter)
		_, found = k.GetVoteFromStore(suite.ctx, oldKey)
		suite.Require().False(found)
	}
	trackerKey := keeper.ConstructVotingChallengeTrackerKey(collection.CollectionId, alice, "outgoing", approvalID, challenge.ProposalId)
	tracker, found := k.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().True(tracker.QuorumReachedTimestamp.IsZero())
	isolatedKey := keeper.ConstructVotingChallengeTrackerKey(collection.CollectionId, alice, "outgoing", approvalID, "isolated-proposal")
	isolatedTracker, found := k.GetVotingChallengeTrackerFromStore(suite.ctx, isolatedKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(40), isolatedTracker.QuorumReachedTimestamp)
	approval := &types.CollectionApproval{ApprovalId: approvalID, ApprovalCriteria: &types.ApprovalCriteria{VotingChallenges: []*types.VotingChallenge{challenge}}}
	checkDelay := func(ctx sdk.Context, expectedError string) {
		checked := false
		for _, checker := range k.GetApprovalCriteriaCheckers(approval) {
			if checker.Name() != "VotingChallenges" {
				continue
			}
			checked = true
			_, err := checker.Check(ctx, approval, collection, bob, alice, alice, "outgoing", alice, nil, nil, "", true)
			if expectedError == "" {
				suite.Require().NoError(err)
			} else {
				suite.Require().ErrorContains(err, expectedError)
			}
		}
		suite.Require().True(checked)
	}
	checkDelay(suite.ctx, "vote again to initialize")
	server := keeper.NewMsgServerImpl(k)
	_, err = server.CastVote(wctx, &types.MsgCastVote{Creator: bob, CollectionId: collection.CollectionId, ApproverAddress: alice, ApprovalLevel: "outgoing", ApprovalId: approvalID, ProposalId: challenge.ProposalId, YesWeight: sdkmath.NewUint(100)})
	suite.Require().NoError(err)
	checkDelay(suite.ctx, "delay not elapsed")
	now := sdkmath.NewUint(uint64(suite.ctx.BlockTime().UnixMilli()))
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Second))))
	tracker, found = k.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().Equal(now, tracker.QuorumReachedTimestamp, "rerunning migration does not reset a re-established delay")
	isolatedTracker, found = k.GetVotingChallengeTrackerFromStore(suite.ctx, isolatedKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(40), isolatedTracker.QuorumReachedTimestamp)
	checkDelay(suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(100*time.Millisecond)), "")
}

func (suite *TestSuite) TestV35VotingBalanceNamespaceCollisionResetsDelay() {
	k := suite.app.TokenizationKeeper
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.Require().NoError(CreateCollections(suite, wctx, []*types.MsgNewCollection{v35MintCollection(&types.ApprovalCriteria{OverridesFromOutgoingApprovals: true}, GetFullUintRanges())}))
	collection, err := GetCollection(suite, wctx, sdkmath.OneUint())
	suite.Require().NoError(err)
	challenge := &types.VotingChallenge{ProposalId: "balance-collision", QuorumThreshold: sdkmath.NewUint(100), DelayAfterQuorum: sdkmath.NewUint(100), Voters: []*types.Voter{{Address: bob, Weight: sdkmath.OneUint()}}}
	approvalID := "approval"
	balance, _, err := k.GetBalanceOrApplyDefault(suite.ctx, collection, alice)
	suite.Require().NoError(err)
	balance.OutgoingApprovals = []*types.UserOutgoingApproval{{ApprovalId: approvalID, ToListId: "All", InitiatedByListId: "All", TransferTimes: GetFullUintRanges(), OwnershipTimes: GetFullUintRanges(), TokenIds: GetFullUintRanges(), Version: sdkmath.ZeroUint(), ApprovalCriteria: &types.OutgoingApprovalCriteria{VotingChallenges: []*types.VotingChallenge{challenge}}}}
	upper := strings.ToUpper(alice)
	suite.Require().NoError(k.SetUserBalanceInStore(suite.ctx, keeper.ConstructBalanceKey(alice, collection.CollectionId), balance, false))
	challenge.DelayAfterQuorum = sdkmath.OneUint()
	rawStore := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	upperBalanceKey := append(append([]byte{}, keeper.UserBalanceKey...), []byte(keeper.ConstructBalanceKey(upper, collection.CollectionId))...)
	rawStore.Set(upperBalanceKey, suite.app.AppCodec().MustMarshal(balance))
	oldVoteKey := keeper.ConstructVotingTrackerKey(collection.CollectionId, upper, "outgoing", approvalID, challenge.ProposalId, bob)
	oldTrackerKey := keeper.ConstructVotingChallengeTrackerKey(collection.CollectionId, upper, "outgoing", approvalID, challenge.ProposalId)
	vote := &types.VoteProof{ProposalId: challenge.ProposalId, Voter: bob, YesWeight: sdkmath.NewUint(100), VotedAt: sdkmath.OneUint()}
	suite.Require().NoError(k.SetVoteInStore(suite.ctx, oldVoteKey, vote))
	suite.Require().NoError(k.SetVotingChallengeTrackerInStore(suite.ctx, oldTrackerKey, &types.VotingChallengeTracker{QuorumReachedTimestamp: sdkmath.OneUint()}))
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	trackerKey := keeper.ConstructVotingChallengeTrackerKey(collection.CollectionId, alice, "outgoing", approvalID, challenge.ProposalId)
	tracker, found := k.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().True(tracker.QuorumReachedTimestamp.IsZero(), "discarded uppercase approval policy must not supply the surviving policy's delay")
	storedVote, found := k.GetVoteFromStore(suite.ctx, keeper.ConstructVotingTrackerKey(collection.CollectionId, alice, "outgoing", approvalID, challenge.ProposalId, bob))
	suite.Require().True(found)
	suite.Require().Equal(vote, storedVote)
	storedBalance, _, err := k.GetBalanceOrApplyDefault(suite.ctx, collection, alice)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(100), storedBalance.OutgoingApprovals[0].ApprovalCriteria.VotingChallenges[0].DelayAfterQuorum)
	server := keeper.NewMsgServerImpl(k)
	_, err = server.CastVote(wctx, &types.MsgCastVote{Creator: bob, CollectionId: collection.CollectionId, ApproverAddress: alice, ApprovalLevel: "outgoing", ApprovalId: approvalID, ProposalId: challenge.ProposalId, YesWeight: sdkmath.NewUint(100)})
	suite.Require().NoError(err)
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Second))))
	tracker, found = k.GetVotingChallengeTrackerFromStore(suite.ctx, trackerKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(uint64(suite.ctx.BlockTime().UnixMilli())), tracker.QuorumReachedTimestamp)
}
