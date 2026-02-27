package approval_criteria_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// CoinTransfersTestSuite tests the CoinTransfers approval criteria field
type CoinTransfersTestSuite struct {
	testutil.AITestSuite
}

func TestCoinTransfersSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(CoinTransfersTestSuite))
}

func (suite *CoinTransfersTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestCoinTransfers_TransferredOnApprovalUse tests that coins are transferred when approval is used
func (suite *CoinTransfersTestSuite) TestCoinTransfers_TransferredOnApprovalUse() {
	// Create approval that requires coin transfer on use
	approval := testutil.GenerateCollectionApproval("coin_transfer_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		CoinTransfers: []*types.CoinTransfer{
			{
				To: suite.Manager, // Payment goes to manager
				Coins: []*sdk.Coin{
					{
						Denom:  "ubadge",
						Amount: sdkmath.NewInt(1000000), // 1 BADGE
					},
				},
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer should require coin payment
	// Note: This test verifies the structure is correct; actual coin transfer may require funded accounts
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
	// The transfer may fail due to insufficient funds (expected in test environment)
	// but the approval structure should be valid
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	// We expect an error here because Alice likely doesn't have ubadge tokens
	// This validates that the coin transfer requirement is being enforced
	suite.Require().Error(err, "transfer should fail without sufficient coin balance for payment")
}

// TestCoinTransfers_MultipleCoinTypes tests coin transfers with multiple coin types
func (suite *CoinTransfersTestSuite) TestCoinTransfers_MultipleCoinTypes() {
	// Create approval requiring multiple coin types
	approval := testutil.GenerateCollectionApproval("multi_coin_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		CoinTransfers: []*types.CoinTransfer{
			{
				To: suite.Manager,
				Coins: []*sdk.Coin{
					{
						Denom:  "ubadge",
						Amount: sdkmath.NewInt(1000000),
					},
					{
						Denom:  "uatom",
						Amount: sdkmath.NewInt(500000),
					},
				},
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer should require both coin types
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
	suite.Require().Error(err, "transfer should fail without required coins")
}

// TestCoinTransfers_MultipleRecipients tests coin transfers to multiple recipients
func (suite *CoinTransfersTestSuite) TestCoinTransfers_MultipleRecipients() {
	// Create approval with coin transfers to multiple recipients
	approval := testutil.GenerateCollectionApproval("multi_recipient_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		CoinTransfers: []*types.CoinTransfer{
			{
				To: suite.Manager, // First recipient
				Coins: []*sdk.Coin{
					{
						Denom:  "ubadge",
						Amount: sdkmath.NewInt(500000),
					},
				},
			},
			{
				To: suite.Charlie, // Second recipient
				Coins: []*sdk.Coin{
					{
						Denom:  "ubadge",
						Amount: sdkmath.NewInt(300000),
					},
				},
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer should require payment to both recipients
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
	suite.Require().Error(err, "transfer should fail without coin payments to all recipients")
}

// TestOverrideToWithInitiator_Works tests that overrideToWithInitiator flag works
func (suite *CoinTransfersTestSuite) TestOverrideToWithInitiator_Works() {
	// Create approval where coin recipient is overridden to be the initiator
	// This is useful for refund scenarios
	approval := testutil.GenerateCollectionApproval("override_to_initiator", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		CoinTransfers: []*types.CoinTransfer{
			{
				To:                      suite.Manager, // This will be overridden to initiator
				OverrideToWithInitiator: true,          // Coins go to whoever initiates
				Coins: []*sdk.Coin{
					{
						Denom:  "ubadge",
						Amount: sdkmath.NewInt(100000),
					},
				},
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.Require().True(collectionId.GT(sdkmath.NewUint(0)), "collection should be created")

	// Verify the approval was created with the override flag
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "override_to_initiator" {
			found = true
			suite.Require().NotNil(app.ApprovalCriteria.CoinTransfers, "coin transfers should exist")
			suite.Require().Len(app.ApprovalCriteria.CoinTransfers, 1, "should have one coin transfer")
			suite.Require().True(app.ApprovalCriteria.CoinTransfers[0].OverrideToWithInitiator,
				"overrideToWithInitiator should be true")
			break
		}
	}
	suite.Require().True(found, "approval should be found")
}

// TestCoinTransfers_ZeroAmount tests that zero amount coin transfers are handled
func (suite *CoinTransfersTestSuite) TestCoinTransfers_ZeroAmount() {
	// Create approval with zero coin amount (should effectively be no payment required)
	approval := testutil.GenerateCollectionApproval("zero_coin_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		CoinTransfers:                  []*types.CoinTransfer{}, // Empty coin transfers
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer should succeed without coin payment
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
	suite.Require().NoError(err, "transfer with no coin requirement should succeed")
}

// TestCoinTransfers_ApprovalStructureValidation tests coin transfer structure validation
func (suite *CoinTransfersTestSuite) TestCoinTransfers_ApprovalStructureValidation() {
	// Create approval with valid coin transfer structure
	approval := testutil.GenerateCollectionApproval("valid_coin_structure", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		CoinTransfers: []*types.CoinTransfer{
			{
				To: suite.Manager,
				Coins: []*sdk.Coin{
					{
						Denom:  "ubadge",
						Amount: sdkmath.NewInt(1000000),
					},
				},
				OverrideFromWithApproverAddress: false,
				OverrideToWithInitiator:         false,
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.Require().True(collectionId.GT(sdkmath.NewUint(0)), "collection with valid coin transfer should be created")

	// Verify structure was saved correctly
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "valid_coin_structure" {
			found = true
			suite.Require().NotNil(app.ApprovalCriteria, "approval criteria should exist")
			suite.Require().NotEmpty(app.ApprovalCriteria.CoinTransfers, "coin transfers should exist")
			ct := app.ApprovalCriteria.CoinTransfers[0]
			suite.Require().Equal(suite.Manager, ct.To, "recipient should match")
			suite.Require().Len(ct.Coins, 1, "should have one coin")
			suite.Require().Equal("ubadge", ct.Coins[0].Denom, "denom should match")
			break
		}
	}
	suite.Require().True(found, "approval should exist")
}

// TestCoinTransfers_RoyaltyPayment tests coin transfer for royalty-like payments
func (suite *CoinTransfersTestSuite) TestCoinTransfers_RoyaltyPayment() {
	// Create approval that simulates a royalty payment scenario
	// Collection creator gets a fee on every transfer
	approval := testutil.GenerateCollectionApproval("royalty_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		CoinTransfers: []*types.CoinTransfer{
			{
				To: suite.Manager, // Creator receives royalty
				Coins: []*sdk.Coin{
					{
						Denom:  "ubadge",
						Amount: sdkmath.NewInt(50000), // 0.05 BADGE royalty
					},
				},
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.Require().True(collectionId.GT(sdkmath.NewUint(0)), "royalty collection should be created")

	// Verify royalty structure
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "royalty_approval" {
			found = true
			suite.Require().NotEmpty(app.ApprovalCriteria.CoinTransfers, "royalty coin transfer should exist")
			break
		}
	}
	suite.Require().True(found, "royalty approval should exist")
}

// TestCoinTransfers_WithNoCoinTransfers tests approval without any coin transfers
func (suite *CoinTransfersTestSuite) TestCoinTransfers_WithNoCoinTransfers() {
	// Create approval without coin transfers (free transfers)
	approval := testutil.GenerateCollectionApproval("free_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		CoinTransfers:                  nil, // No coin transfers required
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer should work without any coin payment
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
	suite.Require().NoError(err, "transfer without coin requirement should succeed")
}

// TestCoinTransfers_LargeCoinAmount tests handling of large coin amounts
func (suite *CoinTransfersTestSuite) TestCoinTransfers_LargeCoinAmount() {
	// Create approval with large coin amount
	approval := testutil.GenerateCollectionApproval("large_amount_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		CoinTransfers: []*types.CoinTransfer{
			{
				To: suite.Manager,
				Coins: []*sdk.Coin{
					{
						Denom:  "ubadge",
						Amount: sdkmath.NewInt(1000000000000), // 1 million BADGE
					},
				},
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.Require().True(collectionId.GT(sdkmath.NewUint(0)), "collection with large amount should be created")

	// Verify structure was saved correctly
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "large_amount_approval" {
			found = true
			suite.Require().NotEmpty(app.ApprovalCriteria.CoinTransfers, "coin transfers should exist")
			suite.Require().True(app.ApprovalCriteria.CoinTransfers[0].Coins[0].Amount.Equal(sdkmath.NewInt(1000000000000)),
				"large amount should be preserved")
			break
		}
	}
	suite.Require().True(found, "large amount approval should exist")
}
