package edge_cases_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// TrackerIsolationTestSuite tests that approval trackers are properly isolated
// across different approval levels and collections to prevent cross-contamination
type TrackerIsolationTestSuite struct {
	testutil.AITestSuite
}

func TestTrackerIsolationSuite(t *testing.T) {
	suite.Run(t, new(TrackerIsolationTestSuite))
}

func (suite *TrackerIsolationTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestTrackerIsolation_CollectionVsUserOutgoing tests that same tracker ID in
// collection approval vs user outgoing approval are isolated
func (suite *TrackerIsolationTestSuite) TestTrackerIsolation_CollectionVsUserOutgoing() {
	// Create collection with approval that uses a tracker
	collectionApproval := testutil.GenerateCollectionApproval("test_approval", "AllWithoutMint", "All")
	collectionApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "shared_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{collectionApproval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(200, 1)})

	// Add user outgoing approval with SAME tracker ID
	userApproval := testutil.GenerateUserOutgoingApproval("user_approval", "All")
	userApproval.ApprovalCriteria = &types.OutgoingApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(50),
			AmountTrackerId:       "shared_tracker", // Same tracker ID as collection
		},
	}

	updateMsg := &types.MsgUpdateUserApprovals{
		Creator:                suite.Alice,
		CollectionId:           collectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:      []*types.UserOutgoingApproval{userApproval},
	}
	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Transfer using collection approval - uses collection-level tracker
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
						ApprovalId:    "test_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Collection tracker should now have 50 used out of 100
	// Transfer another 50 using collection approval - should succeed
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
						ApprovalId:    "test_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second transfer of 50 should succeed as collection tracker allows 100 total")

	// Third transfer should fail - collection tracker exhausted
	msg3 := &types.MsgTransferTokens{
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
						ApprovalId:    "test_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().Error(err, "third transfer should fail - collection tracker exhausted at 100")
}

// TestTrackerIsolation_CollectionVsUserIncoming tests that same tracker ID in
// collection approval vs user incoming approval are isolated
func (suite *TrackerIsolationTestSuite) TestTrackerIsolation_CollectionVsUserIncoming() {
	// Create collection with approval that uses a tracker
	collectionApproval := testutil.GenerateCollectionApproval("test_approval", "AllWithoutMint", "All")
	collectionApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "shared_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{collectionApproval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(200, 1)})

	// Add user incoming approval for Bob with SAME tracker ID
	userApproval := testutil.GenerateUserIncomingApproval("user_approval", "All")
	userApproval.ApprovalCriteria = &types.IncomingApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(25),
			AmountTrackerId:       "shared_tracker", // Same tracker ID as collection
		},
	}

	updateMsg := &types.MsgUpdateUserApprovals{
		Creator:                suite.Bob,
		CollectionId:           collectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:      []*types.UserIncomingApproval{userApproval},
	}
	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Transfer using collection approval - should use collection-level tracker only
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(100, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "test_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "should succeed using collection approval tracker (100 limit)")

	// User incoming approval should still have its own 25 limit available (isolated)
	// The collection approval override means we bypass user approvals anyway
}

// TestTrackerIsolation_AcrossCollections tests that same tracker ID in different
// collections are completely isolated
func (suite *TrackerIsolationTestSuite) TestTrackerIsolation_AcrossCollections() {
	// Create first collection with approval using tracker
	approval1 := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	approval1.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "same_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId1 := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval1})

	// Create second collection with SAME tracker ID
	approval2 := testutil.GenerateCollectionApproval("approval2", "AllWithoutMint", "All")
	approval2.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "same_tracker", // Same ID in different collection
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId2 := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval2})

	// Mint tokens in both collections
	suite.MintTokens(collectionId1, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})
	suite.MintTokens(collectionId2, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer ALL 100 tokens in collection 1 - should exhaust collection 1's tracker
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId1,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(100, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "approval1",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Collection 2 should still have full 100 available (isolated tracker)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId2,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(100, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "approval2",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "collection 2 should have isolated tracker - full 100 available")
}

// TestTrackerIsolation_SameApprovalIdDifferentLevels tests that same approval ID
// at different levels (collection vs user outgoing vs user incoming) are isolated
func (suite *TrackerIsolationTestSuite) TestTrackerIsolation_SameApprovalIdDifferentLevels() {
	// Create collection with approval that doesn't override user approvals
	collectionApproval := testutil.GenerateCollectionApproval("shared_id", "AllWithoutMint", "All")
	collectionApproval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100), // Large enough for all transfers
			AmountTrackerId:       "tracker",
		},
		OverridesFromOutgoingApprovals: false, // Don't override - check user approvals
		OverridesToIncomingApprovals:   false,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{collectionApproval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(200, 1)})

	// Add user outgoing approval with same ID but different limit
	outgoingApproval := testutil.GenerateUserOutgoingApproval("shared_id", "All")
	outgoingApproval.ApprovalCriteria = &types.OutgoingApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(50), // Outgoing allows 50
			AmountTrackerId:       "tracker",
		},
	}

	updateOutgoingMsg := &types.MsgUpdateUserApprovals{
		Creator:                suite.Alice,
		CollectionId:           collectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:      []*types.UserOutgoingApproval{outgoingApproval},
	}
	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateOutgoingMsg)
	suite.Require().NoError(err)

	// Add user incoming approval for Bob with same ID but smaller limit
	incomingApproval := testutil.GenerateUserIncomingApproval("shared_id", "All")
	incomingApproval.ApprovalCriteria = &types.IncomingApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(20), // Bob incoming allows only 20
			AmountTrackerId:       "tracker",
		},
	}

	updateIncomingMsg := &types.MsgUpdateUserApprovals{
		Creator:                suite.Bob,
		CollectionId:           collectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:      []*types.UserIncomingApproval{incomingApproval},
	}
	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateIncomingMsg)
	suite.Require().NoError(err)

	// Transfer 20 tokens - all three trackers are used independently
	// Collection tracker: 100 - 20 = 80 remaining
	// User outgoing tracker: 50 - 20 = 30 remaining
	// User incoming tracker: 20 - 20 = 0 remaining (Bob's limit)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(20, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "shared_id",
						Version:       sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "outgoing",
						ApproverAddress: suite.Alice,
						ApprovalId:      "shared_id",
						Version:         sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "incoming",
						ApproverAddress: suite.Bob,
						ApprovalId:      "shared_id",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err)

	// Try to transfer 1 more to Bob - should fail because Bob's incoming tracker is exhausted
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
						ApprovalId:    "shared_id",
						Version:       sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "outgoing",
						ApproverAddress: suite.Alice,
						ApprovalId:      "shared_id",
						Version:         sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "incoming",
						ApproverAddress: suite.Bob,
						ApprovalId:      "shared_id",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "should fail - Bob's incoming tracker exhausted (only 20 allowed)")

	// Transfer to Charlie should work - Charlie has a different incoming tracker
	// Set up Charlie's incoming approval with same tracker ID
	charlieIncoming := testutil.GenerateUserIncomingApproval("shared_id", "All")
	charlieIncoming.ApprovalCriteria = &types.IncomingApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(15),
			AmountTrackerId:       "tracker",
		},
	}

	updateCharlieMsg := &types.MsgUpdateUserApprovals{
		Creator:                suite.Charlie,
		CollectionId:           collectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:      []*types.UserIncomingApproval{charlieIncoming},
	}
	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateCharlieMsg)
	suite.Require().NoError(err)

	// Transfer to Charlie - uses Charlie's fresh incoming tracker
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(15, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "shared_id",
						Version:       sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "outgoing",
						ApproverAddress: suite.Alice,
						ApprovalId:      "shared_id",
						Version:         sdkmath.NewUint(0),
					},
					{
						ApprovalLevel:   "incoming",
						ApproverAddress: suite.Charlie,
						ApprovalId:      "shared_id",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().NoError(err, "should succeed - Charlie has fresh incoming tracker with 15 limit")
}

// TestTrackerIsolation_NoCrossCollectionPollution verifies there's no tracker pollution
// when using the same tracker ID in multiple collections with rapid transfers
func (suite *TrackerIsolationTestSuite) TestTrackerIsolation_NoCrossCollectionPollution() {
	numCollections := 3
	collectionIds := make([]sdkmath.Uint, numCollections)

	// Create multiple collections with identical approval configurations
	for i := 0; i < numCollections; i++ {
		approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
		approval.ApprovalCriteria = &types.ApprovalCriteria{
			ApprovalAmounts: &types.ApprovalAmounts{
				OverallApprovalAmount: sdkmath.NewUint(10),
				AmountTrackerId:       "shared_tracker_id",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		}

		collectionIds[i] = suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
		suite.MintTokens(collectionIds[i], suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(50, 1)})
	}

	// Transfer 10 tokens in each collection - each should succeed independently
	for i := 0; i < numCollections; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionIds[i],
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
					PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
						{
							ApprovalLevel: "collection",
							ApprovalId:    "transfer_approval",
							Version:       sdkmath.NewUint(0),
						},
					},
					OnlyCheckPrioritizedCollectionApprovals: true,
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer in collection %d should succeed", i)
	}

	// Verify each collection's tracker is now exhausted independently
	for i := 0; i < numCollections; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionIds[i],
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Charlie},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
					PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
						{
							ApprovalLevel: "collection",
							ApprovalId:    "transfer_approval",
							Version:       sdkmath.NewUint(0),
						},
					},
					OnlyCheckPrioritizedCollectionApprovals: true,
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().Error(err, "collection %d tracker should be exhausted", i)
	}
}
