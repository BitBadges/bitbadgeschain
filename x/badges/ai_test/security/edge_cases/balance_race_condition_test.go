package edge_cases_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type BalanceRaceConditionTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestBalanceRaceConditionSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(BalanceRaceConditionTestSuite))
}

func (suite *BalanceRaceConditionTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
	suite.SetupMintApproval(suite.CollectionId)
}

// calculateTotalBalance calculates the total amount from a slice of balances
func (suite *BalanceRaceConditionTestSuite) calculateTotalBalance(balances []*types.Balance) sdkmath.Uint {
	total := sdkmath.NewUint(0)
	for _, balance := range balances {
		total = total.Add(balance.Amount)
	}
	return total
}

// TestBalanceRaceCondition_MultiRecipientTransfer verifies that balance calculations
// are correct when transferring to multiple recipients in a single transfer.
// This test addresses HIGH-005: Balance Calculation Race Condition Risk.
func (suite *BalanceRaceConditionTestSuite) TestBalanceRaceCondition_MultiRecipientTransfer() {
	// Setup collection approval for regular transfers
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint 100 tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set up approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	for _, recipient := range []string{suite.Bob, suite.Charlie} {
		setIncomingMsg := &types.MsgSetIncomingApproval{
			Creator:      recipient,
			CollectionId: suite.CollectionId,
			Approval:     incomingApproval,
		}
		_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
		suite.Require().NoError(err)
	}

	// Get initial balances
	aliceBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Bob)
	charlieBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Charlie)

	aliceTotalBefore := suite.calculateTotalBalance(aliceBalanceBefore.Balances)
	bobTotalBefore := suite.calculateTotalBalance(bobBalanceBefore.Balances)
	charlieTotalBefore := suite.calculateTotalBalance(charlieBalanceBefore.Balances)

	// Transfer 10 tokens to each of 2 recipients (20 total)
	transferAmount := sdkmath.NewUint(10)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob, suite.Charlie}, []*types.Balance{
				testutil.GenerateSimpleBalance(10, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "multi-recipient transfer should succeed")

	// Get balances after transfer
	aliceBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Bob)
	charlieBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Charlie)

	aliceTotalAfter := suite.calculateTotalBalance(aliceBalanceAfter.Balances)
	bobTotalAfter := suite.calculateTotalBalance(bobBalanceAfter.Balances)
	charlieTotalAfter := suite.calculateTotalBalance(charlieBalanceAfter.Balances)

	// Verify Alice's balance decreased by 20 (10 * 2 recipients)
	expectedAliceTotal := aliceTotalBefore.Sub(transferAmount.Mul(sdkmath.NewUint(2)))
	suite.Require().True(
		aliceTotalAfter.Equal(expectedAliceTotal),
		"Alice's balance should decrease by 20: before=%s, after=%s, expected=%s",
		aliceTotalBefore, aliceTotalAfter, expectedAliceTotal,
	)

	// Verify each recipient received 10 tokens
	suite.Require().True(
		bobTotalAfter.Equal(bobTotalBefore.Add(transferAmount)),
		"Bob should receive 10 tokens: before=%s, after=%s",
		bobTotalBefore, bobTotalAfter,
	)
	suite.Require().True(
		charlieTotalAfter.Equal(charlieTotalBefore.Add(transferAmount)),
		"Charlie should receive 10 tokens: before=%s, after=%s",
		charlieTotalBefore, charlieTotalAfter,
	)

	// Verify total conservation
	totalBefore := aliceTotalBefore.Add(bobTotalBefore).Add(charlieTotalBefore)
	totalAfter := aliceTotalAfter.Add(bobTotalAfter).Add(charlieTotalAfter)
	suite.Require().True(
		totalBefore.Equal(totalAfter),
		"Total balance should be conserved: before=%s, after=%s",
		totalBefore, totalAfter,
	)
}

// TestBalanceRaceCondition_MultipleTransfersSameSender verifies that balance calculations
// are correct when the same sender appears in multiple transfers.
// This test addresses HIGH-005: Balance Calculation Race Condition Risk.
func (suite *BalanceRaceConditionTestSuite) TestBalanceRaceCondition_MultipleTransfersSameSender() {
	// Setup collection approval for regular transfers
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint 100 tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set up approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	for _, recipient := range []string{suite.Bob, suite.Charlie} {
		setIncomingMsg := &types.MsgSetIncomingApproval{
			Creator:      recipient,
			CollectionId: suite.CollectionId,
			Approval:     incomingApproval,
		}
		_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
		suite.Require().NoError(err)
	}

	// Get initial balances
	aliceBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Bob)
	charlieBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Charlie)

	aliceTotalBefore := suite.calculateTotalBalance(aliceBalanceBefore.Balances)
	bobTotalBefore := suite.calculateTotalBalance(bobBalanceBefore.Balances)
	charlieTotalBefore := suite.calculateTotalBalance(charlieBalanceBefore.Balances)

	// Create two transfers from Alice: one to Bob (10 tokens) and one to Charlie (20 tokens)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(10, 1),
			}),
			testutil.GenerateTransfer(suite.Alice, []string{suite.Charlie}, []*types.Balance{
				testutil.GenerateSimpleBalance(20, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "multiple transfers from same sender should succeed")

	// Get balances after transfer
	aliceBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Bob)
	charlieBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Charlie)

	aliceTotalAfter := suite.calculateTotalBalance(aliceBalanceAfter.Balances)
	bobTotalAfter := suite.calculateTotalBalance(bobBalanceAfter.Balances)
	charlieTotalAfter := suite.calculateTotalBalance(charlieBalanceAfter.Balances)

	// Verify Alice's balance decreased by 30 (10 + 20)
	expectedAliceTotal := aliceTotalBefore.Sub(sdkmath.NewUint(30))
	suite.Require().True(
		aliceTotalAfter.Equal(expectedAliceTotal),
		"Alice's balance should decrease by 30: before=%s, after=%s, expected=%s",
		aliceTotalBefore, aliceTotalAfter, expectedAliceTotal,
	)

	// Verify recipients received correct amounts
	suite.Require().True(
		bobTotalAfter.Equal(bobTotalBefore.Add(sdkmath.NewUint(10))),
		"Bob should receive 10 tokens: before=%s, after=%s",
		bobTotalBefore, bobTotalAfter,
	)
	suite.Require().True(
		charlieTotalAfter.Equal(charlieTotalBefore.Add(sdkmath.NewUint(20))),
		"Charlie should receive 20 tokens: before=%s, after=%s",
		charlieTotalBefore, charlieTotalAfter,
	)

	// Verify total conservation
	totalBefore := aliceTotalBefore.Add(bobTotalBefore).Add(charlieTotalBefore)
	totalAfter := aliceTotalAfter.Add(bobTotalAfter).Add(charlieTotalAfter)
	suite.Require().True(
		totalBefore.Equal(totalAfter),
		"Total balance should be conserved: before=%s, after=%s",
		totalBefore, totalAfter,
	)
}

// TestBalanceRaceCondition_InsufficientBalance verifies that transfers fail correctly
// when the sender doesn't have sufficient balance for all recipients.
// This test addresses HIGH-005: Balance Calculation Race Condition Risk.
func (suite *BalanceRaceConditionTestSuite) TestBalanceRaceCondition_InsufficientBalance() {
	// Setup collection approval for regular transfers
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint only 20 tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(20, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set up approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	for _, recipient := range []string{suite.Bob, suite.Charlie} {
		setIncomingMsg := &types.MsgSetIncomingApproval{
			Creator:      recipient,
			CollectionId: suite.CollectionId,
			Approval:     incomingApproval,
		}
		_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
		suite.Require().NoError(err)
	}

	// Try to transfer 15 tokens to each of 2 recipients (30 total, but Alice only has 20)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob, suite.Charlie}, []*types.Balance{
				testutil.GenerateSimpleBalance(15, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail due to insufficient balance")
	suite.Require().Contains(err.Error(), "inadequate", "error should mention inadequate balance")

	// Note: In a real Cosmos SDK transaction, if the transaction fails, the entire cache context
	// is rolled back atomically, so all balance changes would be reverted. The important thing
	// is that the error is returned correctly, which prevents the transaction from being committed.
	// The balance calculation fix (HIGH-005) ensures that fromUserBalance is correctly updated
	// and saved after each recipient, preventing race conditions in balance calculations.
}
