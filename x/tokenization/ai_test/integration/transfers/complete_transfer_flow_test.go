package transfers_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type CompleteTransferFlowTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestCompleteTransferFlowSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(CompleteTransferFlowTestSuite))
}

func (suite *CompleteTransferFlowTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestCompleteTransferFlow_AllThreeApprovals tests complete transfer flow with all three approval levels
func (suite *CompleteTransferFlowTestSuite) TestCompleteTransferFlow_AllThreeApprovals() {
	// Setup collection approval
	collectionApproval := testutil.GenerateCollectionApproval("collection1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{collectionApproval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set outgoing approval for Alice
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	// Set incoming approval for Bob
	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: suite.CollectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err)

	// Perform transfer - all three approval levels should be checked
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
	suite.Require().NoError(err, "transfer should succeed with all three approval levels")

	// Verify balances
	aliceBalance := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalance := suite.GetBalance(suite.CollectionId, suite.Bob)

	// Alice should have 50 tokens remaining (100 - 50)
	suite.Require().Greater(len(aliceBalance.Balances), 0, "Alice should have remaining balance")
	// Bob should have 50 tokens
	suite.Require().Greater(len(bobBalance.Balances), 0, "Bob should have received tokens")
}

// TestCompleteTransferFlow_MissingCollectionApproval tests that transfer fails without collection approval
func (suite *CompleteTransferFlowTestSuite) TestCompleteTransferFlow_MissingCollectionApproval() {
	// Mint tokens
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set only user approvals (missing collection approval)
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err := suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: suite.CollectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err)

	// Attempt transfer - should fail without collection approval
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
	suite.Require().Error(err, "transfer should fail without collection approval")
}

// TestCompleteTransferFlow_MissingOutgoingApproval tests that transfer fails without outgoing approval
func (suite *CompleteTransferFlowTestSuite) TestCompleteTransferFlow_MissingOutgoingApproval() {
	// Setup collection approval
	collectionApproval := testutil.GenerateCollectionApproval("collection1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{collectionApproval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set only incoming approval (missing outgoing approval)
	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: suite.CollectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err)

	// Attempt transfer - should fail without outgoing approval
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
	suite.Require().Error(err, "transfer should fail without outgoing approval")
}

// TestCompleteTransferFlow_MissingIncomingApproval tests that transfer fails without incoming approval
func (suite *CompleteTransferFlowTestSuite) TestCompleteTransferFlow_MissingIncomingApproval() {
	// Setup collection approval
	collectionApproval := testutil.GenerateCollectionApproval("collection1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{collectionApproval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set only outgoing approval (missing incoming approval)
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	// Attempt transfer - should fail without incoming approval
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
	suite.Require().Error(err, "transfer should fail without incoming approval")
}

