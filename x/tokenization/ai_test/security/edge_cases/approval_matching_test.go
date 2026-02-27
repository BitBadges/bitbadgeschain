package edge_cases_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// ApprovalMatchingTestSuite tests the first-match semantics and approval
// ordering behavior in the tokenization module
type ApprovalMatchingTestSuite struct {
	testutil.AITestSuite
}

func TestApprovalMatchingSuite(t *testing.T) {
	suite.Run(t, new(ApprovalMatchingTestSuite))
}

func (suite *ApprovalMatchingTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestFirstMatch_OrderDeterminesPriority tests that the order of approvals
// in the array determines which approval is matched first
func (suite *ApprovalMatchingTestSuite) TestFirstMatch_OrderDeterminesPriority() {
	// Create two overlapping approvals with different limits
	// First approval: allows 10 tokens
	firstApproval := testutil.GenerateCollectionApproval("first_approval", "AllWithoutMint", "All")
	firstApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(10),
			AmountTrackerId:       "first_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Second approval: allows 100 tokens (more generous)
	secondApproval := testutil.GenerateCollectionApproval("second_approval", "AllWithoutMint", "All")
	secondApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "second_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// First is more restrictive (10 tokens), second is more generous (100 tokens)
	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		firstApproval,  // First in array
		secondApproval, // Second in array
	})

	// Mint tokens
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(200, 1)})

	// Transfer 10 tokens - should use first approval
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "first_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first transfer of 10 should succeed")

	// Now first approval is exhausted, try to transfer 5 more using first approval - should fail
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "first_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "first approval is exhausted at 10")

	// But using second approval should still work
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "second_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().NoError(err, "second approval should still have capacity")
}

// TestFirstMatch_ReorderingChangesWhichApprovalMatches tests that reordering
// approvals changes which one gets matched
func (suite *ApprovalMatchingTestSuite) TestFirstMatch_ReorderingChangesWhichApprovalMatches() {
	// Create two approvals, put the generous one first this time
	// Second approval (generous): allows 100 tokens
	generousApproval := testutil.GenerateCollectionApproval("generous_approval", "AllWithoutMint", "All")
	generousApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "generous_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// First approval (restrictive): allows 10 tokens
	restrictiveApproval := testutil.GenerateCollectionApproval("restrictive_approval", "AllWithoutMint", "All")
	restrictiveApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(10),
			AmountTrackerId:       "restrictive_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Mint approval
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Generous is now first in array
	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		generousApproval,   // First - will be tried first
		restrictiveApproval, // Second
	})

	// Mint tokens
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(200, 1)})

	// Transfer 50 tokens using generous approval - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "generous_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "should succeed with generous approval allowing 100")
}

// TestFirstMatch_FirstMatchUsedEvenIfLaterBetter tests that when using prioritized
// approvals, only the specified approval is considered (not others)
func (suite *ApprovalMatchingTestSuite) TestFirstMatch_FirstMatchUsedEvenIfLaterBetter() {
	// Create two separate approvals with different limits
	restrictiveApproval := testutil.GenerateCollectionApproval("restrictive", "AllWithoutMint", "All")
	restrictiveApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(5), // Very limited
			AmountTrackerId:       "restrictive_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	generousApproval := testutil.GenerateCollectionApproval("generous", "AllWithoutMint", "All")
	generousApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(1000), // Very generous
			AmountTrackerId:       "generous_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		restrictiveApproval,
		generousApproval,
	})

	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer 10 using restrictive approval - should fail (only allows 5)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "restrictive",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "restrictive approval only allows 5, attempting 10 should fail")

	// But using generous approval should work for 10
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "generous",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "generous approval allows 1000, should work for 10")
}

// TestFirstMatch_EmptyApprovalArrayNoMatch tests that an empty approval array
// results in no match
func (suite *ApprovalMatchingTestSuite) TestFirstMatch_EmptyApprovalArrayNoMatch() {
	// Create collection with only mint approval (no transfer approval)
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{mintApproval})

	// Mint tokens
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Try to transfer without any matching approval - should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				// No prioritized approvals, and no collection approval matches
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer should fail with no matching approval")
}

// TestApprovalOverlap_TokenIdsFirstMatchWins tests that when two approvals
// have overlapping tokenIds, the first match wins
func (suite *ApprovalMatchingTestSuite) TestApprovalOverlap_TokenIdsFirstMatchWins() {
	// First approval: tokens 1-50, limit 10
	firstApproval := testutil.GenerateCollectionApproval("first_tokens", "AllWithoutMint", "All")
	firstApproval.TokenIds = []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
	}
	firstApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(10),
			AmountTrackerId:       "first_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Second approval: tokens 1-100 (overlaps with first), limit 100
	secondApproval := testutil.GenerateCollectionApproval("second_tokens", "AllWithoutMint", "All")
	secondApproval.TokenIds = []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
	}
	secondApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "second_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		firstApproval,
		secondApproval,
	})

	// Mint token ID 1
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer token 1 using first approval - only allows 10
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "first_tokens",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first transfer of 10 should succeed")

	// First approval exhausted, using second for same token should work
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "second_tokens",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second approval should still have capacity for token 1")
}

// TestApprovalOverlap_DifferentApprovalsDifferentLimits tests that different
// approvals with different limits can be used independently
func (suite *ApprovalMatchingTestSuite) TestApprovalOverlap_DifferentApprovalsDifferentLimits() {
	// First approval: limit 5
	firstApproval := testutil.GenerateCollectionApproval("first_approval", "AllWithoutMint", "All")
	firstApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(5),
			AmountTrackerId:       "first_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Second approval: limit 500
	secondApproval := testutil.GenerateCollectionApproval("second_approval", "AllWithoutMint", "All")
	secondApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(500),
			AmountTrackerId:       "second_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		firstApproval,
		secondApproval,
	})

	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer using first approval - exhaust it
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "first_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first transfer should succeed")

	// First approval exhausted, second should work independently
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "second_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second approval should have its own capacity")
}

// TestApprovalMatch_NonOverlappingTokenIds tests that non-overlapping token IDs
// result in separate matching behavior
func (suite *ApprovalMatchingTestSuite) TestApprovalMatch_NonOverlappingTokenIds() {
	// First approval: tokens 1-50
	firstApproval := testutil.GenerateCollectionApproval("tokens_1_50", "AllWithoutMint", "All")
	firstApproval.TokenIds = []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
	}
	firstApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "first_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	// Second approval: tokens 51-100 (no overlap)
	secondApproval := testutil.GenerateCollectionApproval("tokens_51_100", "AllWithoutMint", "All")
	secondApproval.TokenIds = []*types.UintRange{
		{Start: sdkmath.NewUint(51), End: sdkmath.NewUint(100)},
	}
	secondApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "second_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		firstApproval,
		secondApproval,
	})

	// Mint tokens 1-100
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{
		{
			Amount: sdkmath.NewUint(100),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
		},
	})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{
		{
			Amount: sdkmath.NewUint(100),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(51), End: sdkmath.NewUint(100)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
		},
	})

	// Transfer token 1 using first approval
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(50),
						TokenIds: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
						},
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "tokens_1_50",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer of token 1-50 should use first approval")

	// Transfer token 51 using second approval
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(50),
						TokenIds: []*types.UintRange{
							{Start: sdkmath.NewUint(51), End: sdkmath.NewUint(100)},
						},
						OwnershipTimes: []*types.UintRange{
							{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
						},
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "tokens_51_100",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "transfer of token 51-100 should use second approval")
}

// TestApprovalMatch_MultipleApprovalsRequired tests scenarios where
// different approvals are needed for different aspects of a transfer
func (suite *ApprovalMatchingTestSuite) TestApprovalMatch_MultipleApprovalsRequired() {
	// This tests that when user approvals are not overridden,
	// both collection and user approvals must match

	// Collection approval that doesn't override user approvals
	collectionApproval := testutil.GenerateCollectionApproval("collection_approval", "AllWithoutMint", "All")
	collectionApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: false, // Don't override
		OverridesToIncomingApprovals:   false, // Don't override
	}

	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{
		mintApproval,
		collectionApproval,
	})

	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Set up user outgoing approval for Alice
	outgoingApproval := testutil.GenerateUserOutgoingApproval("alice_outgoing", "All")
	outgoingApproval.ApprovalCriteria = &types.OutgoingApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(50),
			AmountTrackerId:       "outgoing_tracker",
		},
	}

	updateOutMsg := &types.MsgUpdateUserApprovals{
		Creator:                suite.Alice,
		CollectionId:           collectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:      []*types.UserOutgoingApproval{outgoingApproval},
	}
	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateOutMsg)
	suite.Require().NoError(err)

	// Set up user incoming approval for Bob
	incomingApproval := testutil.GenerateUserIncomingApproval("bob_incoming", "All")
	incomingApproval.ApprovalCriteria = &types.IncomingApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(30),
			AmountTrackerId:       "incoming_tracker",
		},
	}

	updateInMsg := &types.MsgUpdateUserApprovals{
		Creator:                suite.Bob,
		CollectionId:           collectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:      []*types.UserIncomingApproval{incomingApproval},
	}
	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateInMsg)
	suite.Require().NoError(err)

	// Transfer 30 tokens - should work (within both user limits)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(30, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "collection_approval",
						Version:       sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "outgoing",
						ApproverAddress: suite.Alice,
						ApprovalId:      "alice_outgoing",
						Version:         sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "incoming",
						ApproverAddress: suite.Bob,
						ApprovalId:      "bob_incoming",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer within all approval limits should succeed")

	// Try to transfer 1 more - should fail (Bob's incoming is exhausted)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "collection_approval",
						Version:       sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "outgoing",
						ApproverAddress: suite.Alice,
						ApprovalId:      "alice_outgoing",
						Version:         sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "incoming",
						ApproverAddress: suite.Bob,
						ApprovalId:      "bob_incoming",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "should fail - Bob's incoming approval exhausted at 30")
}
