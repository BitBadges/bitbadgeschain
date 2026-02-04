package attack_scenarios_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type ApprovalBypassAttackTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestApprovalBypassAttackSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(ApprovalBypassAttackTestSuite))
}

func (suite *ApprovalBypassAttackTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestApprovalBypass_WithoutCollectionApproval tests that transfers fail without collection approval
func (suite *ApprovalBypassAttackTestSuite) TestApprovalBypass_WithoutCollectionApproval() {
	// Mint tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set only outgoing approval (missing collection approval)
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err := suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
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

	// Attempt transfer - should fail without collection approval
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(5, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail without collection approval")
}

// TestApprovalBypass_WithoutOutgoingApproval tests that transfers fail without outgoing approval
func (suite *ApprovalBypassAttackTestSuite) TestApprovalBypass_WithoutOutgoingApproval() {
	// Setup collection approval
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
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
				testutil.GenerateSimpleBalance(5, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail without outgoing approval")
}

// TestApprovalBypass_WithoutIncomingApproval tests that transfers fail without incoming approval
func (suite *ApprovalBypassAttackTestSuite) TestApprovalBypass_WithoutIncomingApproval() {
	// Setup collection approval
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
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
				testutil.GenerateSimpleBalance(5, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail without incoming approval")
}

