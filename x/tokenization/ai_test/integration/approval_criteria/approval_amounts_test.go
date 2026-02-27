package approval_criteria_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// ApprovalAmountsTestSuite tests the ApprovalAmounts approval criteria field
type ApprovalAmountsTestSuite struct {
	testutil.AITestSuite
}

func TestApprovalAmountsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(ApprovalAmountsTestSuite))
}

func (suite *ApprovalAmountsTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestOverallApprovalAmount_EnforcesLimit tests that overall approval amount limit is enforced
func (suite *ApprovalAmountsTestSuite) TestOverallApprovalAmount_EnforcesLimit() {
	// Create approval with overallApprovalAmount = 10
	approval := testutil.GenerateCollectionApproval("limitedapproval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(10),
			AmountTrackerId:       "tracker1",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// First transfer of 10 should succeed
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
	suite.Require().NoError(err, "first transfer within limit should succeed")

	// Second transfer of 1 should fail (exceeded limit)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer exceeding overall limit should fail")
}

// TestOverallApprovalAmount_ZeroMeansUnlimited tests that 0 means unlimited
func (suite *ApprovalAmountsTestSuite) TestOverallApprovalAmount_ZeroMeansUnlimited() {
	// Create approval with overallApprovalAmount = 0 (unlimited)
	approval := testutil.GenerateCollectionApproval("unlimitedapproval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(0), // 0 = unlimited
			AmountTrackerId:       "trackerunlimited",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1000, 1)})

	// Multiple large transfers should all succeed
	for i := 0; i < 5; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(100, 1)},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer with unlimited approval should succeed")
	}
}

// TestPerToAddressApprovalAmount_EnforcesLimit tests per-to-address limit enforcement
func (suite *ApprovalAmountsTestSuite) TestPerToAddressApprovalAmount_EnforcesLimit() {
	// Create approval with perToAddressApprovalAmount = 5
	approval := testutil.GenerateCollectionApproval("pertolimited", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			PerToAddressApprovalAmount: sdkmath.NewUint(5),
			AmountTrackerId:            "pertotracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// First transfer of 5 to Bob should succeed
	msg := &types.MsgTransferTokens{
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
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first transfer to Bob within limit should succeed")

	// Second transfer of 1 to Bob should fail (exceeded per-to limit)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer exceeding per-to limit should fail")

	// Transfer to Charlie should succeed (different "to" address)
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().NoError(err, "transfer to different address should succeed")
}

// TestPerFromAddressApprovalAmount_EnforcesLimit tests per-from-address limit enforcement
func (suite *ApprovalAmountsTestSuite) TestPerFromAddressApprovalAmount_EnforcesLimit() {
	// Create approval with perFromAddressApprovalAmount = 5
	approval := testutil.GenerateCollectionApproval("perfromlimited", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			PerFromAddressApprovalAmount: sdkmath.NewUint(5),
			AmountTrackerId:              "perfromtracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})
	suite.MintTokens(collectionId, suite.Bob, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// First transfer of 5 from Alice should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first transfer from Alice within limit should succeed")

	// Second transfer from Alice should fail (exceeded per-from limit)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer exceeding per-from limit should fail")

	// Transfer from Bob should succeed (different "from" address)
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Bob,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Bob,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().NoError(err, "transfer from different address should succeed")
}

// TestPerInitiatedByAddressApprovalAmount_EnforcesLimit tests per-initiated-by limit enforcement
func (suite *ApprovalAmountsTestSuite) TestPerInitiatedByAddressApprovalAmount_EnforcesLimit() {
	// Create approval with perInitiatedByAddressApprovalAmount = 5
	approval := testutil.GenerateCollectionApproval("perinitiatorlimited", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			PerInitiatedByAddressApprovalAmount: sdkmath.NewUint(5),
			AmountTrackerId:                     "perinitiatortracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// First transfer initiated by Alice (5) should succeed
	msg := &types.MsgTransferTokens{
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
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first transfer initiated by Alice should succeed")

	// Second transfer initiated by Alice should fail (exceeded per-initiator limit)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer exceeding per-initiator limit should fail")
}

// TestAmountTrackerId_ScopingBetweenApprovals tests that amountTrackerId properly scopes between approvals
func (suite *ApprovalAmountsTestSuite) TestAmountTrackerId_ScopingBetweenApprovals() {
	// Create two approvals with different tracker IDs
	approval1 := testutil.GenerateCollectionApproval("approvaltracker1", "AllWithoutMint", "All")
	approval1.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(5),
			AmountTrackerId:       "trackera", // First tracker
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	approval2 := testutil.GenerateCollectionApproval("approvaltracker2", "AllWithoutMint", "All")
	approval2.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(5),
			AmountTrackerId:       "trackerb", // Different tracker
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval1, approval2})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Use up trackera's limit
	msg := &types.MsgTransferTokens{
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
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first approval use should succeed")

	// Since there are two approvals with different trackers, first-match applies
	// The second transfer will use the second approval with trackerb
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(5, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second approval with different tracker should allow transfer")
}

// TestTransferExceedingLimit_Rejected tests that transfers exceeding limits are rejected
func (suite *ApprovalAmountsTestSuite) TestTransferExceedingLimit_Rejected() {
	// Create approval with limit of 10
	approval := testutil.GenerateCollectionApproval("strictlimit", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(10),
			AmountTrackerId:       "stricttracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Try to transfer 11 (exceeding limit of 10) - should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(11, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "single transfer exceeding limit should fail")
}

// TestCombinedApprovalAmountLimits tests multiple limit types combined
func (suite *ApprovalAmountsTestSuite) TestCombinedApprovalAmountLimits() {
	// Create approval with both overall and per-to limits
	approval := testutil.GenerateCollectionApproval("combinedlimits", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount:      sdkmath.NewUint(20), // Overall limit of 20
			PerToAddressApprovalAmount: sdkmath.NewUint(10), // Per-to limit of 10
			AmountTrackerId:            "combinedtracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer 10 to Bob - succeeds (within both limits)
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
	suite.Require().NoError(err, "first transfer within both limits should succeed")

	// Transfer 10 to Charlie - succeeds (per-to resets, overall = 20)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second transfer to different address should succeed")

	// Transfer 1 more to anyone - fails (overall limit reached)
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Manager},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().Error(err, "transfer exceeding overall limit should fail")
}

// TestApprovalAmounts_IncrementalTransfers tests incremental transfers up to the limit
func (suite *ApprovalAmountsTestSuite) TestApprovalAmounts_IncrementalTransfers() {
	// Create approval with limit of 10
	// Use a stricter approval with explicit amount limits
	approval := testutil.GenerateCollectionApproval("incrementallimit", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(10),
			AmountTrackerId:       "incrementaltracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	// Mint enough tokens for 11 transfers (need 11 to test the failure)
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Make 10 transfers of 1 each - all should succeed
	// Must use PrioritizedApprovals when AmountTrackerId is set
	for i := 0; i < 10; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
					PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
						{
							ApprovalId:    "incrementallimit",
							ApprovalLevel: "collection",
							Version:       sdkmath.NewUint(0),
						},
					},
					OnlyCheckPrioritizedCollectionApprovals: true,
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "incremental transfer %d should succeed", i+1)
	}

	// 11th transfer should fail - we've exhausted the limit of 10
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:    "incrementallimit",
						ApprovalLevel: "collection",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer after limit reached should fail")
}
