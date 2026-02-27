package edge_cases_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// TransactionAtomicityTestSuite tests that transactions are atomic -
// either all operations succeed or all are rolled back
type TransactionAtomicityTestSuite struct {
	testutil.AITestSuite
}

func TestTransactionAtomicitySuite(t *testing.T) {
	suite.Run(t, new(TransactionAtomicityTestSuite))
}

func (suite *TransactionAtomicityTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestAtomicity_MultipleTransfersProcessedSequentially documents that when
// a message contains multiple transfers, they are processed sequentially.
// If a later transfer fails, earlier successful transfers are NOT rolled back.
// This is important behavior to understand for security considerations.
func (suite *TransactionAtomicityTestSuite) TestAtomicity_MultipleTransfersProcessedSequentially() {
	// Create collection with limited transfer approval
	approval := testutil.GenerateCollectionApproval("transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint 100 tokens to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Record Alice's balance before attempted transfer
	aliceBalanceBefore := suite.GetBalance(collectionId, suite.Alice)
	bobBalanceBefore := suite.GetBalance(collectionId, suite.Bob)

	// Calculate total balance before
	aliceTotalBefore := suite.calculateTotalBalance(aliceBalanceBefore.Balances)
	bobTotalBefore := suite.calculateTotalBalance(bobBalanceBefore.Balances)

	// Attempt a message with multiple transfers:
	// Transfer 1: Alice -> Bob: 50 tokens (will succeed)
	// Transfer 2: Alice -> Charlie: 100 tokens (will fail - insufficient balance after first transfer)
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
						ApprovalId:    "transfer_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(100, 1)}, // Would exceed remaining balance
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
	suite.Require().Error(err, "transaction should fail due to insufficient balance for second transfer")

	// IMPORTANT: With sequential processing, the first transfer succeeds before the second fails
	// This documents actual behavior - transfers are processed in order within a message
	aliceBalanceAfter := suite.GetBalance(collectionId, suite.Alice)
	bobBalanceAfter := suite.GetBalance(collectionId, suite.Bob)
	charlieBalanceAfter := suite.GetBalance(collectionId, suite.Charlie)

	aliceTotalAfter := suite.calculateTotalBalance(aliceBalanceAfter.Balances)
	bobTotalAfter := suite.calculateTotalBalance(bobBalanceAfter.Balances)
	charlieTotalAfter := suite.calculateTotalBalance(charlieBalanceAfter.Balances)

	// Document actual behavior: first transfer was executed before failure
	suite.Require().True(aliceTotalAfter.Equal(aliceTotalBefore.Sub(sdkmath.NewUint(50))),
		"Alice's balance should reflect the first successful transfer: expected=%s, got=%s",
		aliceTotalBefore.Sub(sdkmath.NewUint(50)), aliceTotalAfter)
	suite.Require().True(bobTotalAfter.Equal(bobTotalBefore.Add(sdkmath.NewUint(50))),
		"Bob should have received tokens from first transfer: expected=%s, got=%s",
		bobTotalBefore.Add(sdkmath.NewUint(50)), bobTotalAfter)
	suite.Require().True(charlieTotalAfter.Equal(sdkmath.NewUint(0)),
		"Charlie should not have received any tokens (second transfer failed)")

	// Total conservation still holds
	totalBefore := aliceTotalBefore.Add(bobTotalBefore)
	totalAfter := aliceTotalAfter.Add(bobTotalAfter).Add(charlieTotalAfter)
	suite.Require().True(totalBefore.Equal(totalAfter),
		"Total tokens should be conserved despite partial failure")
}

// TestAtomicity_FailedTransferRollsBackBalanceChanges tests that a failed
// transfer does not partially update balances
func (suite *TransactionAtomicityTestSuite) TestAtomicity_FailedTransferRollsBackBalanceChanges() {
	// Create collection with approval that has a limit
	approval := testutil.GenerateCollectionApproval("limited_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(50), // Only allow 50 tokens
			AmountTrackerId:       "limit_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Record initial state
	aliceBalanceBefore := suite.GetBalance(collectionId, suite.Alice)
	aliceTotalBefore := suite.calculateTotalBalance(aliceBalanceBefore.Balances)

	// Try to transfer more than approval allows
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(75, 1)}, // Exceeds 50 limit
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "limited_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer exceeding approval limit should fail")

	// Verify Alice's balance is unchanged
	aliceBalanceAfter := suite.GetBalance(collectionId, suite.Alice)
	aliceTotalAfter := suite.calculateTotalBalance(aliceBalanceAfter.Balances)

	suite.Require().True(aliceTotalBefore.Equal(aliceTotalAfter),
		"Alice's balance should be unchanged: before=%s, after=%s", aliceTotalBefore, aliceTotalAfter)
}

// TestAtomicity_FailedTransferConsumesApprovalTracker documents an IMPORTANT behavior:
// When a transfer fails due to exceeding the approval limit (partial success),
// the approval tracker IS STILL CONSUMED for the amount that was "partially" processed.
//
// SECURITY NOTE: This means failed transfers can drain approval limits!
// Applications should be aware that attempting a transfer that exceeds the
// remaining approval limit will consume all remaining approval, even though
// no tokens are actually transferred.
func (suite *TransactionAtomicityTestSuite) TestAtomicity_FailedTransferConsumesApprovalTracker() {
	// Create collection with tracked approval - large enough for testing
	approval := testutil.GenerateCollectionApproval("tracked_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		ApprovalAmounts: &types.ApprovalAmounts{
			OverallApprovalAmount: sdkmath.NewUint(100),
			AmountTrackerId:       "test_tracker",
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens to Alice - IMPORTANT: mint more than the approval limit
	// This ensures any failures are due to approval limit, not balance
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(200, 1)})

	// First, do a successful transfer of 40 tokens
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(40, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "tracked_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err, "first transfer of 40 should succeed")

	// Verify Bob received the tokens
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	bobTotal := suite.calculateTotalBalance(bobBalance.Balances)
	suite.Require().True(bobTotal.Equal(sdkmath.NewUint(40)),
		"Bob should have 40 tokens from first transfer")

	// Tracker now has 60 remaining (100 - 40)
	// Alice still has 160 tokens (200 - 40)
	// Try a transfer exceeding approval limit - this will fail with "partial success"
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(80, 1)}, // Exceeds 60 approval remaining
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "tracked_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer exceeding approval limit should fail")
	suite.Require().Contains(err.Error(), "partial success",
		"error should mention partial success - this indicates the tracker was partially consumed")

	// IMPORTANT BEHAVIOR: The failed transfer consumed the remaining approval!
	// Even attempting to transfer 1 token now fails
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(1, 1)}, // Even 1 token fails
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "tracked_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().Error(err, "approval tracker should be exhausted after failed partial transfer")
	suite.Require().Contains(err.Error(), "no balances found for approval",
		"error should indicate approval is exhausted")

	// Verify balances: Alice still has all tokens (no actual transfer occurred from failed attempts)
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	aliceTotal := suite.calculateTotalBalance(aliceBalance.Balances)
	suite.Require().True(aliceTotal.Equal(sdkmath.NewUint(160)),
		"Alice should have 160 remaining (200 - 40 from successful transfer only)")

	charlieBalance := suite.GetBalance(collectionId, suite.Charlie)
	charlieTotal := suite.calculateTotalBalance(charlieBalance.Balances)
	suite.Require().True(charlieTotal.Equal(sdkmath.NewUint(0)),
		"Charlie should have 0 - no tokens were transferred despite approval being consumed")
}

// TestAtomicity_StateUnchangedAfterFailure verifies complete state consistency
// after a failed transaction
func (suite *TransactionAtomicityTestSuite) TestAtomicity_StateUnchangedAfterFailure() {
	// Create collection
	approval := testutil.GenerateCollectionApproval("test_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens to multiple users
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})
	suite.MintTokens(collectionId, suite.Bob, []*types.Balance{testutil.GenerateSimpleBalance(50, 1)})

	// Capture full state before
	aliceBalanceBefore := suite.GetBalance(collectionId, suite.Alice)
	bobBalanceBefore := suite.GetBalance(collectionId, suite.Bob)
	charlieBalanceBefore := suite.GetBalance(collectionId, suite.Charlie)
	collectionBefore := suite.GetCollection(collectionId)

	// Attempt an invalid transfer (to invalid address)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{"invalid_address_format"}, // Invalid address
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer to invalid address should fail")

	// Verify all state is unchanged
	aliceBalanceAfter := suite.GetBalance(collectionId, suite.Alice)
	bobBalanceAfter := suite.GetBalance(collectionId, suite.Bob)
	charlieBalanceAfter := suite.GetBalance(collectionId, suite.Charlie)
	collectionAfter := suite.GetCollection(collectionId)

	suite.Require().Equal(
		suite.calculateTotalBalance(aliceBalanceBefore.Balances),
		suite.calculateTotalBalance(aliceBalanceAfter.Balances),
		"Alice's balance should be unchanged")
	suite.Require().Equal(
		suite.calculateTotalBalance(bobBalanceBefore.Balances),
		suite.calculateTotalBalance(bobBalanceAfter.Balances),
		"Bob's balance should be unchanged")
	suite.Require().Equal(
		suite.calculateTotalBalance(charlieBalanceBefore.Balances),
		suite.calculateTotalBalance(charlieBalanceAfter.Balances),
		"Charlie's balance should be unchanged")
	suite.Require().Equal(len(collectionBefore.CollectionApprovals), len(collectionAfter.CollectionApprovals),
		"Collection approvals should be unchanged")
}

// TestAtomicity_MultiRecipientTransferDocumentation documents behavior when a
// transfer has multiple recipients. If the total exceeds balance, the behavior
// depends on when the check happens - recipients may be processed sequentially.
func (suite *TransactionAtomicityTestSuite) TestAtomicity_MultiRecipientTransferDocumentation() {
	// Create collection with approval
	approval := testutil.GenerateCollectionApproval("multi_recipient_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer to multiple recipients - total exactly equals balance
	// This documents that multi-recipient transfers work correctly
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob, suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)}, // 50 * 2 = 100
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "multi_recipient_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer to multiple recipients exactly matching balance should succeed")

	// Verify balances
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	charlieBalance := suite.GetBalance(collectionId, suite.Charlie)

	aliceTotal := suite.calculateTotalBalance(aliceBalance.Balances)
	bobTotal := suite.calculateTotalBalance(bobBalance.Balances)
	charlieTotal := suite.calculateTotalBalance(charlieBalance.Balances)

	suite.Require().True(aliceTotal.Equal(sdkmath.NewUint(0)),
		"Alice should have 0 after transferring all")
	suite.Require().True(bobTotal.Equal(sdkmath.NewUint(50)),
		"Bob should have 50")
	suite.Require().True(charlieTotal.Equal(sdkmath.NewUint(50)),
		"Charlie should have 50")
}

// TestAtomicity_ApprovalCheckFailureRollback tests that approval check failures
// don't partially modify state
func (suite *TransactionAtomicityTestSuite) TestAtomicity_ApprovalCheckFailureRollback() {
	// Create collection with strict approval
	approval := testutil.GenerateCollectionApproval("strict_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireToEqualsInitiatedBy: true, // Require to == initiatedBy
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Mint tokens to Alice
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Record state before
	aliceBalanceBefore := suite.GetBalance(collectionId, suite.Alice)
	aliceTotalBefore := suite.calculateTotalBalance(aliceBalanceBefore.Balances)

	// Try transfer where to != initiatedBy (Alice is initiator, sending to Bob)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Initiator is Alice
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob}, // But recipient is Bob
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(50, 1)},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalLevel: "collection",
						ApprovalId:    "strict_approval",
						Version:       sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer violating RequireToEqualsInitiatedBy should fail")

	// Verify state unchanged
	aliceBalanceAfter := suite.GetBalance(collectionId, suite.Alice)
	aliceTotalAfter := suite.calculateTotalBalance(aliceBalanceAfter.Balances)

	suite.Require().True(aliceTotalBefore.Equal(aliceTotalAfter),
		"Alice's balance should be unchanged after approval failure")
}

// Helper function to calculate total balance
func (suite *TransactionAtomicityTestSuite) calculateTotalBalance(balances []*types.Balance) sdkmath.Uint {
	total := sdkmath.NewUint(0)
	for _, balance := range balances {
		total = total.Add(balance.Amount)
	}
	return total
}
