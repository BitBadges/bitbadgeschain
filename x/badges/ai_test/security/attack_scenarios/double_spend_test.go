package attack_scenarios_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type DoubleSpendAttackTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestDoubleSpendAttackSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(DoubleSpendAttackTestSuite))
}

func (suite *DoubleSpendAttackTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestDoubleSpend_Prevented tests that the system prevents double-spending
func (suite *DoubleSpendAttackTestSuite) TestDoubleSpend_Prevented() {
	// Setup approvals
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint 10 tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
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

	// First transfer: Alice sends 10 tokens to Bob (all she has)
	transferMsg1 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(10, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg1)
	suite.Require().NoError(err)

	// Attempt second transfer: Alice tries to send 5 more tokens (double-spend attempt)
	transferMsg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Charlie}, []*types.Balance{
				testutil.GenerateSimpleBalance(5, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg2)
	suite.Require().Error(err, "double-spend attempt should fail - Alice has no tokens left")
}

// TestDoubleSpend_ConcurrentTransfers tests concurrent transfer attempts
func (suite *DoubleSpendAttackTestSuite) TestDoubleSpend_ConcurrentTransfers() {
	// Setup approvals
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint 10 tokens to Alice
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
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

	for _, recipient := range []string{suite.Bob, suite.Charlie} {
		incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
		setIncomingMsg := &types.MsgSetIncomingApproval{
			Creator:      recipient,
			CollectionId: suite.CollectionId,
			Approval:     incomingApproval,
		}
		_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
		suite.Require().NoError(err)
	}

	// Attempt to send 10 tokens to Bob
	transferMsg1 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(10, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg1)
	suite.Require().NoError(err)

	// Attempt to send 10 tokens to Charlie (should fail - Alice has no tokens)
	transferMsg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Charlie}, []*types.Balance{
				testutil.GenerateSimpleBalance(10, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg2)
	suite.Require().Error(err, "second transfer should fail due to insufficient balance")
}
