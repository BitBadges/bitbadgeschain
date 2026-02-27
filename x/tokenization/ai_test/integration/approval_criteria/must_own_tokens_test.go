package approval_criteria_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// MustOwnTokensTestSuite tests the MustOwnTokens approval criteria field
type MustOwnTokensTestSuite struct {
	testutil.AITestSuite
}

func TestMustOwnTokensSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(MustOwnTokensTestSuite))
}

func (suite *MustOwnTokensTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestCrossCollectionOwnership_Verified tests that cross-collection ownership is verified
func (suite *MustOwnTokensTestSuite) TestCrossCollectionOwnership_Verified() {
	// First create a "gate" collection that users must own tokens from
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create second collection with MustOwnTokens requirement
	approval := testutil.GenerateCollectionApproval("gate_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1), // Must own at least 1
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipCheckParty: "initiator",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice owns gate token - transfer should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer by gate token owner should succeed")
}

// TestCrossCollectionOwnership_Denied tests denial when ownership requirement not met
func (suite *MustOwnTokensTestSuite) TestCrossCollectionOwnership_Denied() {
	// Create a gate collection - Bob does NOT have tokens
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create second collection with MustOwnTokens requirement
	approval := testutil.GenerateCollectionApproval("gate_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipCheckParty: "initiator",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Bob, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Bob does NOT own gate token - transfer should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Bob,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Bob,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer by non-gate-token owner should fail")
}

// TestAmountRange_Enforced tests that amountRange requirement is enforced
func (suite *MustOwnTokensTestSuite) TestAmountRange_Enforced() {
	// Create gate collection - Alice owns only 1 token
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create collection requiring at least 5 tokens
	approval := testutil.GenerateCollectionApproval("amount_gate", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(5), // Must own at least 5
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipCheckParty: "initiator",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice owns only 1, not 5 - should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer with insufficient token ownership should fail")
}

// TestAmountRange_MaxEnforced tests that max amount in amountRange is enforced
func (suite *MustOwnTokensTestSuite) TestAmountRange_MaxEnforced() {
	// Create gate collection - Alice owns 10 tokens
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(10, 1)})

	// Create collection requiring between 1-5 tokens (exclusive whale restriction)
	approval := testutil.GenerateCollectionApproval("max_amount_gate", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1), // Must own between 1-5
					End:   sdkmath.NewUint(5),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipCheckParty: "initiator",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice owns 10, exceeds max of 5 - should fail
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer by whale (exceeding max) should fail")
}

// TestOwnershipTimes_Respected tests that ownershipTimes requirement is respected
func (suite *MustOwnTokensTestSuite) TestOwnershipTimes_Respected() {
	// Create gate collection with specific ownership times
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{
		{
			Amount: sdkmath.NewUint(1),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(2000)}, // Specific time range
			},
		},
	})

	// Create collection requiring ownership during specific time
	approval := testutil.GenerateCollectionApproval("time_gate", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(2000)}, // Must own during this time
				},
				OwnershipCheckParty: "initiator",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice owns during correct time - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer with ownership during correct time should succeed")
}

// TestMustSatisfyForAllAssets_True tests mustSatisfyForAllAssets=true behavior
func (suite *MustOwnTokensTestSuite) TestMustSatisfyForAllAssets_True() {
	// Create gate collection - Alice owns token 1 but not token 2
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create collection requiring ALL tokens 1 and 2
	approval := testutil.GenerateCollectionApproval("all_assets_gate", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}, // Must own both 1 and 2
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				MustSatisfyForAllAssets: true, // Must satisfy for ALL tokens
				OwnershipCheckParty:     "initiator",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice only owns token 1, not token 2 - should fail with mustSatisfyForAllAssets=true
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "transfer without owning ALL required tokens should fail")
}

// TestMustSatisfyForAllAssets_False tests mustSatisfyForAllAssets=false behavior
func (suite *MustOwnTokensTestSuite) TestMustSatisfyForAllAssets_False() {
	// Create gate collection - Alice owns only token 1
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create collection requiring ANY of tokens 1 or 2
	approval := testutil.GenerateCollectionApproval("any_asset_gate", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2)}, // Token 1 or 2
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				MustSatisfyForAllAssets: false, // Only need to satisfy for ANY token
				OwnershipCheckParty:     "initiator",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice owns token 1 - should succeed with mustSatisfyForAllAssets=false
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer with ANY required token ownership should succeed")
}

// TestOwnershipCheckParty_Initiator tests ownershipCheckParty="initiator"
func (suite *MustOwnTokensTestSuite) TestOwnershipCheckParty_Initiator() {
	// Create gate collection - only Alice owns the token
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create collection checking initiator's ownership
	approval := testutil.GenerateCollectionApproval("initiator_check", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipCheckParty: "initiator", // Check initiator (creator of tx)
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice is initiator and owns gate token - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Alice is the initiator
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer by initiator with gate token should succeed")
}

// TestOwnershipCheckParty_Sender tests ownershipCheckParty="sender"
func (suite *MustOwnTokensTestSuite) TestOwnershipCheckParty_Sender() {
	// Create gate collection - only Alice owns the token
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create collection checking sender's ownership
	approval := testutil.GenerateCollectionApproval("sender_check", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipCheckParty: "sender", // Check sender (from address)
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice is sender and owns gate token - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice, // Alice is the sender
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer where sender has gate token should succeed")
}

// TestOwnershipCheckParty_Recipient tests ownershipCheckParty="recipient"
func (suite *MustOwnTokensTestSuite) TestOwnershipCheckParty_Recipient() {
	// Create gate collection - only Bob owns the token
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Bob, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create collection checking recipient's ownership
	approval := testutil.GenerateCollectionApproval("recipient_check", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipCheckParty: "recipient", // Check recipient (to address)
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Bob is recipient and owns gate token - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob}, // Bob is the recipient
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer to recipient with gate token should succeed")

	// Transfer to Charlie (who doesn't own gate token) should fail
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie}, // Charlie doesn't have gate token
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer to recipient without gate token should fail")
}

// TestMultipleMustOwnTokens_Requirements tests multiple MustOwnTokens requirements
func (suite *MustOwnTokensTestSuite) TestMultipleMustOwnTokens_Requirements() {
	// Create two gate collections
	gateCollection1 := suite.CreateTestCollection(suite.Manager)
	gateCollection2 := suite.CreateTestCollection(suite.Manager)

	// Alice owns both gate tokens
	suite.MintTokens(gateCollection1, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})
	suite.MintTokens(gateCollection2, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create collection requiring ownership from both collections
	approval := testutil.GenerateCollectionApproval("multi_gate", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollection1,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipCheckParty: "initiator",
			},
			{
				CollectionId: gateCollection2,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipCheckParty: "initiator",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice owns both gate tokens - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer with all gate tokens should succeed")
}

// TestOverrideWithCurrentTime tests the overrideWithCurrentTime flag
func (suite *MustOwnTokensTestSuite) TestOverrideWithCurrentTime() {
	// Create gate collection
	gateCollectionId := suite.CreateTestCollection(suite.Manager)
	suite.MintTokens(gateCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(1, 1)})

	// Create collection with overrideWithCurrentTime
	approval := testutil.GenerateCollectionApproval("current_time_gate", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MustOwnTokens: []*types.MustOwnTokens{
			{
				CollectionId: gateCollectionId,
				AmountRange: &types.UintRange{
					Start: sdkmath.NewUint(1),
					End:   sdkmath.NewUint(math.MaxUint64),
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}, // Will be overridden
				},
				OverrideWithCurrentTime: true, // Override to current time
				OwnershipCheckParty:     "initiator",
			},
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	mainCollectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(mainCollectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// With overrideWithCurrentTime, the ownership check uses current block time
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: mainCollectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer with overrideWithCurrentTime should work")
}
