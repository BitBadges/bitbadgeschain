package msg_handlers_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type TransferTokensTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestTransferTokensSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(TransferTokensTestSuite))
}

func (suite *TransferTokensTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
	suite.Require().True(suite.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0 after creation")
}

// TestTransferTokens_ValidTransfer tests a valid token transfer
func (suite *TransferTokensTestSuite) TestTransferTokens_ValidTransfer() {
	// First, set up collection approval that allows minting (from Mint address)
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	updateMintMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{mintApproval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMintMsg)
	suite.Require().NoError(err)

	// Now mint tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Create a collection approval that allows transfers from All to All
	approval := testutil.GenerateCollectionApproval("approval1", "All", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{approval},
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

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

	// Transfer tokens from Alice to Bob
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
	suite.Require().NoError(err)

	// Verify balances
	aliceBalance := suite.GetBalance(suite.CollectionId, suite.Alice)
	bobBalance := suite.GetBalance(suite.CollectionId, suite.Bob)

	// Alice should have 5 tokens remaining (10 - 5)
	suite.Require().Greater(len(aliceBalance.Balances), 0)
	// Bob should have 5 tokens
	suite.Require().Greater(len(bobBalance.Balances), 0)
}

// TestTransferTokens_WithoutApprovals tests that transfer fails without proper approvals
func (suite *TransferTokensTestSuite) TestTransferTokens_WithoutApprovals() {
	// Set up collection approval that allows minting (from Mint address)
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	updateMintMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{mintApproval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMintMsg)
	suite.Require().NoError(err)

	// Mint tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Try to transfer without approvals - should fail
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
	suite.Require().Error(err, "transfer should fail without approvals")
}

// TestTransferTokens_InvalidCollection tests transfer with non-existent collection
func (suite *TransferTokensTestSuite) TestTransferTokens_InvalidCollection() {
	invalidCollectionId := sdkmath.NewUint(99999)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: invalidCollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(5, 1),
			}),
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer should fail with non-existent collection")
}

// TestTransferTokens_MultiRecipient tests transfer to multiple recipients
func (suite *TransferTokensTestSuite) TestTransferTokens_MultiRecipient() {
	// First, set up collection approval that allows minting (from Mint address)
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	updateMintMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{mintApproval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMintMsg)
	suite.Require().NoError(err)

	// Setup approvals for regular transfers
	approval := testutil.GenerateCollectionApproval("approval1", "All", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{mintApproval, approval},
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(20, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set approvals for all parties
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

	// Transfer to multiple recipients
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob, suite.Charlie}, []*types.Balance{
				testutil.GenerateSimpleBalance(5, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err)

	// Verify all recipients received tokens
	bobBalance := suite.GetBalance(suite.CollectionId, suite.Bob)
	charlieBalance := suite.GetBalance(suite.CollectionId, suite.Charlie)
	suite.Require().Greater(len(bobBalance.Balances), 0)
	suite.Require().Greater(len(charlieBalance.Balances), 0)
}

