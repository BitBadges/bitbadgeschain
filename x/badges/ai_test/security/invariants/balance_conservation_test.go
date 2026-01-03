package invariants_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type BalanceConservationTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestBalanceConservationSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(BalanceConservationTestSuite))
}

func (suite *BalanceConservationTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestBalanceConservation_Transfer tests that total balances are conserved during transfers
func (suite *BalanceConservationTestSuite) TestBalanceConservation_Transfer() {
	// Setup approvals for regular transfers
	approval := testutil.GenerateCollectionApproval("approval1", "All", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint 100 tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: suite.CollectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err)

	// Get initial balances
	aliceBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Bob)

	// Calculate total before transfer
	totalBefore := suite.calculateTotalBalance(aliceBalanceBefore.Balances)
	totalBefore = totalBefore.Add(suite.calculateTotalBalance(bobBalanceBefore.Balances))

	// Transfer 50 tokens from Alice to Bob
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(50, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err)

	// Get balances after transfer
	aliceBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Bob)

	// Calculate total after transfer
	totalAfter := suite.calculateTotalBalance(aliceBalanceAfter.Balances)
	totalAfter = totalAfter.Add(suite.calculateTotalBalance(bobBalanceAfter.Balances))

	// Total should be conserved
	suite.Require().True(totalBefore.Equal(totalAfter), "total balance should be conserved: before=%s, after=%s", totalBefore, totalAfter)
}

// calculateTotalBalance calculates the total amount from a slice of balances
func (suite *BalanceConservationTestSuite) calculateTotalBalance(balances []*types.Balance) sdkmath.Uint {
	total := sdkmath.NewUint(0)
	for _, balance := range balances {
		total = total.Add(balance.Amount)
	}
	return total
}

// TestBalanceConservation_MultiTransfer tests balance conservation across multiple transfers
func (suite *BalanceConservationTestSuite) TestBalanceConservation_MultiTransfer() {
	// Setup approvals for regular transfers
	approval := testutil.GenerateCollectionApproval("approval1", "All", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens to multiple addresses
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)
	suite.MintBadges(suite.CollectionId, suite.Bob, mintBalances)

	// Set approvals for all
	for _, addr := range []string{suite.Alice, suite.Bob, suite.Charlie} {
		outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
		setOutgoingMsg := &types.MsgSetOutgoingApproval{
			Creator:      addr,
			CollectionId: suite.CollectionId,
			Approval:     outgoingApproval,
		}
		_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
		suite.Require().NoError(err)

		incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
		setIncomingMsg := &types.MsgSetIncomingApproval{
			Creator:      addr,
			CollectionId: suite.CollectionId,
			Approval:     incomingApproval,
		}
		_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
		suite.Require().NoError(err)
	}

	// Get initial total
	aliceBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Bob)
	charlieBalanceBefore := suite.GetBalance(suite.CollectionId, suite.Charlie)

	totalBefore := suite.calculateTotalBalance(aliceBalanceBefore.Balances)
	totalBefore = totalBefore.Add(suite.calculateTotalBalance(bobBalanceBefore.Balances))
	totalBefore = totalBefore.Add(suite.calculateTotalBalance(charlieBalanceBefore.Balances))

	// Perform multiple transfers
	transfers := []struct {
		from string
		to   string
		amount uint64
	}{
		{suite.Alice, suite.Bob, 30},
		{suite.Bob, suite.Charlie, 20},
		{suite.Charlie, suite.Alice, 10},
	}

	for _, transfer := range transfers {
		transferMsg := &types.MsgTransferTokens{
			Creator:      transfer.from,
			CollectionId: suite.CollectionId,
			Transfers: []*types.Transfer{
				testutil.GenerateTransfer(transfer.from, []string{transfer.to}, []*types.Balance{
					testutil.GenerateSimpleBalance(transfer.amount, 1),
				}),
			},
		}
		_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
		suite.Require().NoError(err)
	}

	// Get final balances
	aliceBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Bob)
	charlieBalanceAfter := suite.GetBalance(suite.CollectionId, suite.Charlie)

	totalAfter := suite.calculateTotalBalance(aliceBalanceAfter.Balances)
	totalAfter = totalAfter.Add(suite.calculateTotalBalance(bobBalanceAfter.Balances))
	totalAfter = totalAfter.Add(suite.calculateTotalBalance(charlieBalanceAfter.Balances))

	// Total should be conserved
	suite.Require().True(totalBefore.Equal(totalAfter), "total balance should be conserved after multiple transfers")
}

