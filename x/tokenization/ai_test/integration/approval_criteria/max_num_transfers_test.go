package approval_criteria_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// MaxNumTransfersTestSuite tests the MaxNumTransfers approval criteria field
type MaxNumTransfersTestSuite struct {
	testutil.AITestSuite
}

func TestMaxNumTransfersSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(MaxNumTransfersTestSuite))
}

func (suite *MaxNumTransfersTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestOverallMaxNumTransfers_EnforcesLimit tests that overall max num transfers limit is enforced
func (suite *MaxNumTransfersTestSuite) TestOverallMaxNumTransfers_EnforcesLimit() {
	// Create approval with overallMaxNumTransfers = 3
	approval := testutil.GenerateCollectionApproval("max_transfers_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(3),
			AmountTrackerId:        "transfer_tracker1",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 10, 1, math.MaxUint64)})

	// First three transfers should succeed
	for i := 0; i < 3; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d within limit should succeed", i+1)
	}

	// Fourth transfer should fail (exceeded max transfers)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 4)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer exceeding max count should fail")
}

// TestOverallMaxNumTransfers_ZeroMeansUnlimited tests that 0 means unlimited transfers
func (suite *MaxNumTransfersTestSuite) TestOverallMaxNumTransfers_ZeroMeansUnlimited() {
	// Create approval with overallMaxNumTransfers = 0 (unlimited)
	approval := testutil.GenerateCollectionApproval("unlimited_transfers", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(0), // 0 = unlimited
			AmountTrackerId:        "unlimited_transfer_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Many transfers should all succeed
	for i := 0; i < 10; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer with unlimited count should succeed")
	}
}

// TestPerToAddressMaxNumTransfers_EnforcesLimit tests per-to-address max transfers limit
func (suite *MaxNumTransfersTestSuite) TestPerToAddressMaxNumTransfers_EnforcesLimit() {
	// Create approval with perToAddressMaxNumTransfers = 2
	approval := testutil.GenerateCollectionApproval("per_to_max", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			PerToAddressMaxNumTransfers: sdkmath.NewUint(2),
			AmountTrackerId:             "per_to_transfer_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// First two transfers to Bob should succeed
	for i := 0; i < 2; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d to Bob should succeed", i+1)
	}

	// Third transfer to Bob should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 3)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "third transfer to Bob should fail")

	// Transfer to Charlie should succeed (different recipient)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 3)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "transfer to different recipient should succeed")
}

// TestPerFromAddressMaxNumTransfers_EnforcesLimit tests per-from-address max transfers limit
func (suite *MaxNumTransfersTestSuite) TestPerFromAddressMaxNumTransfers_EnforcesLimit() {
	// Create approval with perFromAddressMaxNumTransfers = 2
	approval := testutil.GenerateCollectionApproval("per_from_max", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			PerFromAddressMaxNumTransfers: sdkmath.NewUint(2),
			AmountTrackerId:               "per_from_transfer_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})
	suite.MintTokens(collectionId, suite.Bob, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// First two transfers from Alice should succeed
	for i := 0; i < 2; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Charlie},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d from Alice should succeed", i+1)
	}

	// Third transfer from Alice should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 3)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "third transfer from Alice should fail")

	// Transfer from Bob should succeed (different sender)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Bob,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Bob,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 3)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "transfer from different sender should succeed")
}

// TestPerInitiatedByAddressMaxNumTransfers_EnforcesLimit tests per-initiator max transfers limit
func (suite *MaxNumTransfersTestSuite) TestPerInitiatedByAddressMaxNumTransfers_EnforcesLimit() {
	// Create approval with perInitiatedByAddressMaxNumTransfers = 2
	approval := testutil.GenerateCollectionApproval("per_initiator_max", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			PerInitiatedByAddressMaxNumTransfers: sdkmath.NewUint(2),
			AmountTrackerId:                      "per_initiator_transfer_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// First two transfers initiated by Alice should succeed
	for i := 0; i < 2; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d initiated by Alice should succeed", i+1)
	}

	// Third transfer initiated by Alice should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 3)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "third transfer initiated by Alice should fail")
}

// TestTransfersBeyondMax_Rejected tests that transfers beyond max are rejected
func (suite *MaxNumTransfersTestSuite) TestTransfersBeyondMax_Rejected() {
	// Create approval with max of 1 transfer
	approval := testutil.GenerateCollectionApproval("single_transfer", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(1),
			AmountTrackerId:        "single_transfer_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// First transfer should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first transfer should succeed")

	// Second transfer should fail (max is 1)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "second transfer should fail when max is 1")
}

// TestTrackerIncrements_Correctly tests that the tracker increments correctly
func (suite *MaxNumTransfersTestSuite) TestTrackerIncrements_Correctly() {
	// Create approval with max of 5 transfers
	approval := testutil.GenerateCollectionApproval("tracked_transfers", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(5),
			AmountTrackerId:        "tracked_transfer_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Make 5 successful transfers
	for i := 0; i < 5; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d should succeed", i+1)
	}

	// 6th transfer should fail - tracker should be at 5
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 6)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "6th transfer should fail as tracker incremented correctly")
}

// TestCombinedMaxNumTransfersLimits tests multiple max transfer limit types combined
func (suite *MaxNumTransfersTestSuite) TestCombinedMaxNumTransfersLimits() {
	// Create approval with both overall and per-to limits
	approval := testutil.GenerateCollectionApproval("combined_max", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers:      sdkmath.NewUint(4), // Overall limit of 4
			PerToAddressMaxNumTransfers: sdkmath.NewUint(2), // Per-to limit of 2
			AmountTrackerId:             "combined_max_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Two transfers to Bob - succeeds
	for i := 0; i < 2; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d to Bob should succeed", i+1)
	}

	// Third transfer to Bob should fail (per-to limit)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 3)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "third transfer to Bob should fail due to per-to limit")

	// Two transfers to Charlie - succeeds (per-to resets for Charlie)
	for i := 0; i < 2; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Charlie},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, uint64(i+3))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d to Charlie should succeed", i+1)
	}

	// Any more transfers should fail (overall limit of 4 reached)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Manager},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 5)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "5th transfer should fail due to overall limit")
}

// TestMaxNumTransfers_WithDifferentTrackerIds tests that different tracker IDs are independent
func (suite *MaxNumTransfersTestSuite) TestMaxNumTransfers_WithDifferentTrackerIds() {
	// Create two approvals with different tracker IDs
	approval1 := testutil.GenerateCollectionApproval("tracker_a_approval", "AllWithoutMint", "All")
	approval1.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(2),
			AmountTrackerId:        "tracker_a",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	approval2 := testutil.GenerateCollectionApproval("tracker_b_approval", "AllWithoutMint", "All")
	approval2.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(2),
			AmountTrackerId:        "tracker_b",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval1, approval2})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Use up tracker_a's limit (2 transfers)
	for i := 0; i < 2; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d using first approval should succeed", i+1)
	}

	// Due to first-match policy, third transfer will try to use first approval (which is exhausted)
	// Then fall back to second approval with tracker_b
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 3)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "third transfer should use second approval with different tracker")
}
