package keeper_test

import (
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func legacyChallengeTrackerKeyForTest(collectionID sdkmath.Uint, approverAddress, approvalLevel, approvalID, challengeID string, leafIndex sdkmath.Uint) string {
	return strings.Join([]string{collectionID.String(), approverAddress, approvalLevel, approvalID, challengeID, leafIndex.String()}, keeper.BalanceKeyDelimiter)
}

func legacyVotingKeyForTest(collectionID sdkmath.Uint, approverAddress, approvalLevel, approvalID, proposalID, voterAddress string) string {
	return strings.Join([]string{collectionID.String(), approverAddress, approvalLevel, approvalID, proposalID, voterAddress}, keeper.BalanceKeyDelimiter)
}

func legacyVotingChallengeKeyForTest(collectionID sdkmath.Uint, approverAddress, approvalLevel, approvalID, proposalID string) string {
	return strings.Join([]string{collectionID.String(), approverAddress, approvalLevel, approvalID, proposalID}, keeper.BalanceKeyDelimiter)
}

func TestApprovalTrackerKeysSeparateDelimitedIdentifiers(t *testing.T) {
	collectionID := sdkmath.NewUint(1)

	first := keeper.ConstructApprovalTrackerKey(collectionID, "", "mint", "daily-total", "collection", "overall", "")
	second := keeper.ConstructApprovalTrackerKey(collectionID, "", "mint-daily", "total", "collection", "overall", "")

	require.NotEqual(t, first, second)
}

func TestChallengeTrackerKeysSeparateDelimitedIdentifiers(t *testing.T) {
	collectionID := sdkmath.NewUint(1)

	first := keeper.ConstructUsedClaimChallengeKey(collectionID, "", "collection", "mint-daily", "claim", sdkmath.ZeroUint())
	second := keeper.ConstructUsedClaimChallengeKey(collectionID, "", "collection", "mint", "daily-claim", sdkmath.ZeroUint())

	require.NotEqual(t, first, second)
}

func TestVotingKeysSeparateDelimitedIdentifiers(t *testing.T) {
	collectionID := sdkmath.NewUint(1)

	firstVote := keeper.ConstructVotingTrackerKey(collectionID, "", "collection", "mint-daily", "proposal", "voter")
	secondVote := keeper.ConstructVotingTrackerKey(collectionID, "", "collection", "mint", "daily-proposal", "voter")
	require.NotEqual(t, firstVote, secondVote)

	firstChallenge := keeper.ConstructVotingChallengeTrackerKey(collectionID, "", "collection", "mint-daily", "proposal")
	secondChallenge := keeper.ConstructVotingChallengeTrackerKey(collectionID, "", "collection", "mint", "daily-proposal")
	require.NotEqual(t, firstChallenge, secondChallenge)
}

func TestAutoScannablePrioritizedApprovalRejectsStaleVersion(t *testing.T) {
	approval := &types.CollectionApproval{
		ApprovalId: "limited-transfer",
		Version:    sdkmath.NewUint(2),
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(10),
				AmountTrackerId:        "lifetime",
			},
		},
	}
	transfer := &types.Transfer{
		PrioritizedApprovals: []*types.ApprovalIdentifierDetails{{
			ApprovalId:    approval.ApprovalId,
			ApprovalLevel: "collection",
			Version:       sdkmath.NewUint(1),
		}},
		OnlyCheckPrioritizedCollectionApprovals: true,
	}

	_, err := keeper.FilterApprovalsWithPrioritizedHandling(
		[]*types.CollectionApproval{approval}, transfer, "collection", "",
	)
	require.ErrorIs(t, err, types.ErrMismatchedVersions)
}

func (suite *TestSuite) TestLegacyCollidingTrackerStateIsPreservedAndSeparated() {
	collectionID := sdkmath.NewUint(1)
	legacyKey := strings.Join([]string{"1", "", "mint", "daily-total", "collection", "overall", ""}, keeper.BalanceKeyDelimiter)
	legacyTracker := types.ApprovalTracker{NumTransfers: sdkmath.NewUint(2)}
	suite.Require().NoError(suite.app.TokenizationKeeper.SetApprovalTrackerInStoreViaKey(suite.ctx, legacyKey, legacyTracker))

	first, found := suite.app.TokenizationKeeper.GetApprovalTrackerFromStore(
		suite.ctx, collectionID, "", "mint", "daily-total", "collection", "overall", "",
	)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(2), first.NumTransfers)

	second, found := suite.app.TokenizationKeeper.GetApprovalTrackerFromStore(
		suite.ctx, collectionID, "", "mint-daily", "total", "collection", "overall", "",
	)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(2), second.NumTransfers)

	first.NumTransfers = sdkmath.NewUint(3)
	suite.Require().NoError(suite.app.TokenizationKeeper.SetApprovalTrackerInStore(
		suite.ctx, collectionID, "", "mint", "daily-total", first, "collection", "overall", "",
	))

	first, found = suite.app.TokenizationKeeper.GetApprovalTrackerFromStore(
		suite.ctx, collectionID, "", "mint", "daily-total", "collection", "overall", "",
	)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(3), first.NumTransfers)

	second, found = suite.app.TokenizationKeeper.GetApprovalTrackerFromStore(
		suite.ctx, collectionID, "", "mint-daily", "total", "collection", "overall", "",
	)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(2), second.NumTransfers)
}

func (suite *TestSuite) TestLegacyChallengeStateIsPreservedAndSeparated() {
	collectionID := sdkmath.NewUint(1)
	legacyKey := strings.Join([]string{"1", "", "collection", "mint-daily", "claim", "0"}, keeper.BalanceKeyDelimiter)
	suite.Require().NoError(suite.app.TokenizationKeeper.SetChallengeTrackerInStore(suite.ctx, legacyKey, sdkmath.NewUint(2)))

	first, err := suite.app.TokenizationKeeper.GetChallengeTrackerFromStore(
		suite.ctx, collectionID, "", "collection", "mint-daily", "claim", sdkmath.ZeroUint(),
	)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(2), first)

	second, err := suite.app.TokenizationKeeper.GetChallengeTrackerFromStore(
		suite.ctx, collectionID, "", "collection", "mint", "daily-claim", sdkmath.ZeroUint(),
	)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(2), second)

	first, err = suite.app.TokenizationKeeper.IncrementChallengeTrackerInStore(
		suite.ctx, collectionID, "", "collection", "mint-daily", "claim", sdkmath.ZeroUint(),
	)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(3), first)

	second, err = suite.app.TokenizationKeeper.GetChallengeTrackerFromStore(
		suite.ctx, collectionID, "", "collection", "mint", "daily-claim", sdkmath.ZeroUint(),
	)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(2), second)
}

func (suite *TestSuite) TestLegacyVotingStateIsPreservedAndSeparated() {
	collectionID := sdkmath.NewUint(1)
	legacyVoteKey := strings.Join([]string{"1", "", "collection", "mint-daily", "proposal", "voter"}, keeper.BalanceKeyDelimiter)
	legacyVote := &types.VoteProof{ProposalId: "proposal", Voter: "voter", YesWeight: sdkmath.NewUint(25)}
	suite.Require().NoError(suite.app.TokenizationKeeper.SetVoteInStore(suite.ctx, legacyVoteKey, legacyVote))

	firstKey := keeper.ConstructVotingTrackerKey(collectionID, "", "collection", "mint-daily", "proposal", "voter")
	secondKey := keeper.ConstructVotingTrackerKey(collectionID, "", "collection", "mint", "daily-proposal", "voter")
	first, found := suite.app.TokenizationKeeper.GetVoteFromStore(suite.ctx, firstKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(25), first.YesWeight)
	second, found := suite.app.TokenizationKeeper.GetVoteFromStore(suite.ctx, secondKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(25), second.YesWeight)

	first.YesWeight = sdkmath.NewUint(75)
	suite.Require().NoError(suite.app.TokenizationKeeper.SetVoteInStore(suite.ctx, firstKey, first))
	first, found = suite.app.TokenizationKeeper.GetVoteFromStore(suite.ctx, firstKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(75), first.YesWeight)
	second, found = suite.app.TokenizationKeeper.GetVoteFromStore(suite.ctx, secondKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(25), second.YesWeight)
}

func (suite *TestSuite) TestLegacyVotingChallengeStateIsPreservedAndSeparated() {
	collectionID := sdkmath.NewUint(1)
	legacyKey := strings.Join([]string{"1", "", "collection", "mint-daily", "proposal"}, keeper.BalanceKeyDelimiter)
	legacyTracker := &types.VotingChallengeTracker{QuorumReachedTimestamp: sdkmath.NewUint(10)}
	suite.Require().NoError(suite.app.TokenizationKeeper.SetVotingChallengeTrackerInStore(suite.ctx, legacyKey, legacyTracker))

	firstKey := keeper.ConstructVotingChallengeTrackerKey(collectionID, "", "collection", "mint-daily", "proposal")
	secondKey := keeper.ConstructVotingChallengeTrackerKey(collectionID, "", "collection", "mint", "daily-proposal")
	first, found := suite.app.TokenizationKeeper.GetVotingChallengeTrackerFromStore(suite.ctx, firstKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(10), first.QuorumReachedTimestamp)
	second, found := suite.app.TokenizationKeeper.GetVotingChallengeTrackerFromStore(suite.ctx, secondKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(10), second.QuorumReachedTimestamp)

	first.QuorumReachedTimestamp = sdkmath.NewUint(20)
	suite.Require().NoError(suite.app.TokenizationKeeper.SetVotingChallengeTrackerInStore(suite.ctx, firstKey, first))
	first, found = suite.app.TokenizationKeeper.GetVotingChallengeTrackerFromStore(suite.ctx, firstKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(20), first.QuorumReachedTimestamp)
	second, found = suite.app.TokenizationKeeper.GetVotingChallengeTrackerFromStore(suite.ctx, secondKey)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewUint(10), second.QuorumReachedTimestamp)
}

func (suite *TestSuite) TestGetVotesPrefersV2AndKeepsDelimitedScopesSeparate() {
	collectionID := sdkmath.NewUint(1)
	legacyKey := legacyVotingKeyForTest(collectionID, "", "collection", "mint-daily", "proposal", "voter")
	suite.Require().NoError(suite.app.TokenizationKeeper.SetVoteInStore(suite.ctx, legacyKey, &types.VoteProof{
		ProposalId: "proposal", Voter: "voter", YesWeight: sdkmath.NewUint(25),
	}))
	v2Key := keeper.ConstructVotingTrackerKey(collectionID, "", "collection", "mint-daily", "proposal", "voter")
	suite.Require().NoError(suite.app.TokenizationKeeper.SetVoteInStore(suite.ctx, v2Key, &types.VoteProof{
		ProposalId: "proposal", Voter: "voter", YesWeight: sdkmath.NewUint(75),
	}))
	otherKey := keeper.ConstructVotingTrackerKey(collectionID, "", "collection", "mint", "daily-proposal", "voter")
	suite.Require().NoError(suite.app.TokenizationKeeper.SetVoteInStore(suite.ctx, otherKey, &types.VoteProof{
		ProposalId: "daily-proposal", Voter: "voter", YesWeight: sdkmath.NewUint(100),
	}))

	response, err := suite.app.TokenizationKeeper.GetVotes(sdk.WrapSDKContext(suite.ctx), &types.QueryGetVotesRequest{
		CollectionId: collectionID.String(), ApprovalLevel: "collection", ApprovalId: "mint-daily", ProposalId: "proposal",
	})
	suite.Require().NoError(err)
	suite.Require().Len(response.Votes, 1)
	suite.Require().Equal(sdkmath.NewUint(75), response.Votes[0].YesWeight)
}
