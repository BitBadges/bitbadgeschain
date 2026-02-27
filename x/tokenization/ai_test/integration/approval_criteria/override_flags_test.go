package approval_criteria_test

import (
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// OverrideFlagsTestSuite tests the override flag approval criteria fields
type OverrideFlagsTestSuite struct {
	testutil.AITestSuite
}

func TestOverrideFlagsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(OverrideFlagsTestSuite))
}

func (suite *OverrideFlagsTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestOverridesFromOutgoingApprovals_True tests that overridesFromOutgoingApprovals=true bypasses user outgoing check
func (suite *OverrideFlagsTestSuite) TestOverridesFromOutgoingApprovals_True() {
	// Create collection approval that overrides user outgoing approvals
	approval := testutil.GenerateCollectionApproval("override_outgoing", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true, // Bypass user's outgoing approval
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice has NOT set any outgoing approval, but collection approval overrides
	// Transfer should succeed due to override flag
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
	suite.Require().NoError(err, "transfer should succeed with outgoing override flag")
}

// TestOverridesFromOutgoingApprovals_False tests that without override, user outgoing approval is checked
func (suite *OverrideFlagsTestSuite) TestOverridesFromOutgoingApprovals_False() {
	// Create collection approval WITHOUT override
	approval := testutil.GenerateCollectionApproval("no_override_outgoing", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: false, // Does not bypass user's outgoing approval
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice has NOT set any outgoing approval
	// Transfer should fail because user outgoing approval is checked and not set
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
	suite.Require().Error(err, "transfer should fail without user outgoing approval when override is false")
}

// TestOverridesToIncomingApprovals_True tests that overridesToIncomingApprovals=true bypasses user incoming check
func (suite *OverrideFlagsTestSuite) TestOverridesToIncomingApprovals_True() {
	// Create collection approval that overrides user incoming approvals
	approval := testutil.GenerateCollectionApproval("override_incoming", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true, // Bypass user's incoming approval
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Bob has NOT set any incoming approval, but collection approval overrides
	// Transfer should succeed due to override flag
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
	suite.Require().NoError(err, "transfer should succeed with incoming override flag")
}

// TestOverridesToIncomingApprovals_False tests that without override, user incoming approval is checked
func (suite *OverrideFlagsTestSuite) TestOverridesToIncomingApprovals_False() {
	// Create collection approval WITHOUT incoming override
	approval := testutil.GenerateCollectionApproval("no_override_incoming", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   false, // Does not bypass user's incoming approval
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Bob has NOT set any incoming approval
	// Transfer should fail because user incoming approval is checked and not set
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
	suite.Require().Error(err, "transfer should fail without user incoming approval when override is false")
}

// TestBothOverrides_True tests both override flags set to true
func (suite *OverrideFlagsTestSuite) TestBothOverrides_True() {
	// Create collection approval with both overrides
	approval := testutil.GenerateCollectionApproval("both_overrides", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true, // Bypass outgoing
		OverridesToIncomingApprovals:   true, // Bypass incoming
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Neither Alice nor Bob have set any user-level approvals
	// Transfer should succeed due to both override flags
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
	suite.Require().NoError(err, "transfer should succeed with both override flags")
}

// TestBothOverrides_False tests both override flags set to false
func (suite *OverrideFlagsTestSuite) TestBothOverrides_False() {
	// Create collection approval with neither override
	approval := testutil.GenerateCollectionApproval("no_overrides", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: false, // Don't bypass outgoing
		OverridesToIncomingApprovals:   false, // Don't bypass incoming
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Neither Alice nor Bob have set any user-level approvals
	// Transfer should fail due to missing user approvals
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
	suite.Require().Error(err, "transfer should fail without user approvals when no overrides")
}

// TestOverrides_WithUserApprovals tests that user approvals work when overrides are false
func (suite *OverrideFlagsTestSuite) TestOverrides_WithUserApprovals() {
	// Create collection approval WITHOUT overrides
	approval := testutil.GenerateCollectionApproval("no_overrides_user_approvals", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: false,
		OverridesToIncomingApprovals:   false,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice sets outgoing approval
	outgoingApproval := testutil.GenerateUserOutgoingApproval("alice_outgoing", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     outgoingApproval,
	}
	_, err := suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err, "Alice should set outgoing approval")

	// Bob sets incoming approval
	incomingApproval := testutil.GenerateUserIncomingApproval("bob_incoming", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: collectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err, "Bob should set incoming approval")

	// Now transfer should succeed with user approvals in place
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
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer should succeed with user approvals")
}

// TestOverridesOnlyWorkAtCollectionLevel tests that overrides only work at collection level
func (suite *OverrideFlagsTestSuite) TestOverridesOnlyWorkAtCollectionLevel() {
	// Create collection approval with overrides
	collectionApproval := testutil.GenerateCollectionApproval("collection_override", "AllWithoutMint", "All")
	collectionApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{collectionApproval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Verify the override flags are set at collection level
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "collection_override" {
			found = true
			suite.Require().NotNil(app.ApprovalCriteria, "approval criteria should exist")
			suite.Require().True(app.ApprovalCriteria.OverridesFromOutgoingApprovals,
				"outgoing override should be true")
			suite.Require().True(app.ApprovalCriteria.OverridesToIncomingApprovals,
				"incoming override should be true")
			break
		}
	}
	suite.Require().True(found, "collection approval should exist with override flags")
}

// TestMintApproval_RequiresOverrides tests that mint approvals require overrides
func (suite *OverrideFlagsTestSuite) TestMintApproval_RequiresOverrides() {
	// Create mint approval with overrides (required for Mint address)
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true, // Required for Mint
		OverridesToIncomingApprovals:   true, // Required for recipients
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{mintApproval})

	// Mint should work with proper override flags
	msg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "mint with override flags should succeed")
}

// TestPartialOutgoingOverride tests override only for outgoing
func (suite *OverrideFlagsTestSuite) TestPartialOutgoingOverride() {
	// Create collection approval that overrides only outgoing
	approval := testutil.GenerateCollectionApproval("outgoing_only", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,  // Override outgoing
		OverridesToIncomingApprovals:   false, // Do NOT override incoming
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice doesn't need outgoing approval, but Bob needs incoming approval
	// Should fail without Bob's incoming approval
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
	suite.Require().Error(err, "transfer should fail without Bob's incoming approval")

	// Now Bob sets incoming approval
	incomingApproval := testutil.GenerateUserIncomingApproval("bob_incoming", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: collectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err, "Bob should set incoming approval")

	// Now transfer should succeed
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer should succeed with Bob's incoming approval")
}

// TestPartialIncomingOverride tests override only for incoming
func (suite *OverrideFlagsTestSuite) TestPartialIncomingOverride() {
	// Create collection approval that overrides only incoming
	approval := testutil.GenerateCollectionApproval("incoming_only", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: false, // Do NOT override outgoing
		OverridesToIncomingApprovals:   true,  // Override incoming
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Bob doesn't need incoming approval, but Alice needs outgoing approval
	// Should fail without Alice's outgoing approval
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
	suite.Require().Error(err, "transfer should fail without Alice's outgoing approval")

	// Now Alice sets outgoing approval
	outgoingApproval := testutil.GenerateUserOutgoingApproval("alice_outgoing", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err, "Alice should set outgoing approval")

	// Now transfer should succeed
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer should succeed with Alice's outgoing approval")
}

// TestMultipleApprovalsWithDifferentOverrides tests multiple approvals with different override settings
func (suite *OverrideFlagsTestSuite) TestMultipleApprovalsWithDifferentOverrides() {
	// Create two collection approvals with different override settings
	approval1 := testutil.GenerateCollectionApproval("override_both", "AllWithoutMint", "All")
	approval1.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	approval2 := testutil.GenerateCollectionApproval("no_overrides", "AllWithoutMint", "All")
	approval2.ApprovalCriteria = &types.ApprovalCriteria{
		OverridesFromOutgoingApprovals: false,
		OverridesToIncomingApprovals:   false,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval1, approval2})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 10, 1, math.MaxUint64)})

	// Transfer should use first matching approval (first-match policy)
	// The first approval has overrides, so it should succeed
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
	suite.Require().NoError(err, "transfer should succeed using first approval with overrides")
}
