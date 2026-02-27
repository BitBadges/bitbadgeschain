package approval_criteria_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// RequireEqualityFlagsTestSuite tests the require equality flag approval criteria fields
type RequireEqualityFlagsTestSuite struct {
	testutil.AITestSuite
}

func TestRequireEqualityFlagsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(RequireEqualityFlagsTestSuite))
}

func (suite *RequireEqualityFlagsTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestRequireToEqualsInitiatedBy_True_Enforced tests that requireToEqualsInitiatedBy=true is enforced
func (suite *RequireEqualityFlagsTestSuite) TestRequireToEqualsInitiatedBy_True_Enforced() {
	// Create approval requiring to == initiatedBy (self-claim only)
	approval := testutil.GenerateCollectionApproval("selfclaimonly", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireToEqualsInitiatedBy:     true, // Recipient must be the initiator
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Alice initiates and receives - should succeed (to == initiatedBy)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Initiator is Alice
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice}, // Recipient is also Alice
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "self-claim should succeed when to == initiatedBy")

	// Alice initiates but Bob receives - should fail (to != initiatedBy)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Initiator is Alice
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob}, // Recipient is Bob, not Alice
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "claim for another should fail when requireToEqualsInitiatedBy=true")
}

// TestRequireFromEqualsInitiatedBy_True_Enforced tests that requireFromEqualsInitiatedBy=true is enforced
func (suite *RequireEqualityFlagsTestSuite) TestRequireFromEqualsInitiatedBy_True_Enforced() {
	// Create approval requiring from == initiatedBy (only transfer your own tokens)
	approval := testutil.GenerateCollectionApproval("selftransferonly", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireFromEqualsInitiatedBy:   true, // Sender must be the initiator
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice initiates and sends from herself - should succeed (from == initiatedBy)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Initiator is Alice
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice, // Sender is also Alice
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "self-transfer should succeed when from == initiatedBy")
}

// TestRequireToDoesNotEqualInitiatedBy_True_Enforced tests that requireToDoesNotEqualInitiatedBy=true is enforced
func (suite *RequireEqualityFlagsTestSuite) TestRequireToDoesNotEqualInitiatedBy_True_Enforced() {
	// Create approval requiring to != initiatedBy (no self-receiving)
	approval := testutil.GenerateCollectionApproval("noselfreceive", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireToDoesNotEqualInitiatedBy: true, // Recipient must NOT be the initiator
		OverridesFromOutgoingApprovals:   true,
		OverridesToIncomingApprovals:     true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Alice initiates but Bob receives - should succeed (to != initiatedBy)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Initiator is Alice
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob}, // Recipient is Bob, not Alice
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer to another should succeed when requireToDoesNotEqualInitiatedBy=true")

	// Alice initiates and receives - should fail (to == initiatedBy, but we require !=)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Initiator is Alice
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice}, // Recipient is also Alice
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "self-receive should fail when requireToDoesNotEqualInitiatedBy=true")
}

// TestRequireFromDoesNotEqualInitiatedBy_True_Enforced tests that requireFromDoesNotEqualInitiatedBy=true is enforced
func (suite *RequireEqualityFlagsTestSuite) TestRequireFromDoesNotEqualInitiatedBy_True_Enforced() {
	// Create approval requiring from != initiatedBy (delegated transfer only)
	approval := testutil.GenerateCollectionApproval("delegatedonly", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireFromDoesNotEqualInitiatedBy: true, // Sender must NOT be the initiator
		OverridesFromOutgoingApprovals:     true,
		OverridesToIncomingApprovals:       true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Alice initiates and sends from herself - should fail (from == initiatedBy, but we require !=)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Initiator is Alice
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice, // Sender is also Alice
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "self-send should fail when requireFromDoesNotEqualInitiatedBy=true")
}

// TestAllFlagsDefault tests default behavior when all flags are false
func (suite *RequireEqualityFlagsTestSuite) TestAllFlagsDefault() {
	// Create approval with all equality flags at default (false)
	approval := testutil.GenerateCollectionApproval("defaultflags", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireToEqualsInitiatedBy:         false, // No constraint
		RequireFromEqualsInitiatedBy:       false, // No constraint
		RequireToDoesNotEqualInitiatedBy:   false, // No constraint
		RequireFromDoesNotEqualInitiatedBy: false, // No constraint
		OverridesFromOutgoingApprovals:     true,
		OverridesToIncomingApprovals:       true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Any transfer should work with default flags
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
	suite.Require().NoError(err, "transfer should succeed with default flags")

	// Transfer to Charlie should also work (different recipient)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Charlie},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "transfer to different recipient should succeed with default flags")
}

// TestConflictingFlags tests behavior with conflicting flag combinations
func (suite *RequireEqualityFlagsTestSuite) TestConflictingFlags() {
	// Create approval with conflicting flags (both require and require not)
	// This should make it impossible to satisfy
	approval := testutil.GenerateCollectionApproval("conflictingflags", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireToEqualsInitiatedBy:       true, // Require to == initiatedBy
		RequireToDoesNotEqualInitiatedBy: true, // AND require to != initiatedBy (impossible!)
		OverridesFromOutgoingApprovals:   true,
		OverridesToIncomingApprovals:     true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Self-transfer should fail (violates RequireToDoesNotEqualInitiatedBy)
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().Error(err, "self-transfer should fail with conflicting flags")

	// Transfer to other should also fail (violates RequireToEqualsInitiatedBy)
	msg2 := &types.MsgTransferTokens{
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
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer to other should fail with conflicting flags")
}

// TestCombinedFromAndToFlags tests combined from and to equality flags
func (suite *RequireEqualityFlagsTestSuite) TestCombinedFromAndToFlags() {
	// Create approval requiring both from == initiatedBy AND to == initiatedBy
	// This approval requires both from and to to equal the initiator
	// However, self-transfers (from == to) are not allowed by the protocol
	// So this approval combination effectively blocks all transfers
	approval := testutil.GenerateCollectionApproval("selfonly", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireToEqualsInitiatedBy:     true, // to must be initiator
		RequireFromEqualsInitiatedBy:   true, // from must be initiator
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Transfer to other should fail (to != initiator)
	// When both RequireToEqualsInitiatedBy and RequireFromEqualsInitiatedBy are true,
	// it would require from == to == initiator, but self-transfers are not allowed
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
	suite.Require().Error(err, "transfer to other should fail when RequireToEqualsInitiatedBy is true")

	// Transfer from other initiator should also fail (from != initiator)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Bob, // Bob initiates
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice, // but Alice is sender (from != initiator)
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "transfer should fail when from != initiator")
}

// TestRequireNotEquals_Combined tests combined "not equals" flags
func (suite *RequireEqualityFlagsTestSuite) TestRequireNotEquals_Combined() {
	// Create approval requiring both from != initiatedBy AND to != initiatedBy
	// This would be for delegated transfers where neither party is the initiator
	approval := testutil.GenerateCollectionApproval("delegatedthirdparty", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireToDoesNotEqualInitiatedBy:   true, // to must NOT be initiator
		RequireFromDoesNotEqualInitiatedBy: true, // from must NOT be initiator
		OverridesFromOutgoingApprovals:     true,
		OverridesToIncomingApprovals:       true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateSimpleBalance(100, 1)})

	// Manager (third party) initiates transfer from Alice to Bob
	// This should succeed if we properly set up user approvals
	// But in this simple test, it will fail due to the override structure
	// The point is to verify the flag logic is applied
	msg := &types.MsgTransferTokens{
		Creator:      suite.Manager, // Manager initiates (not Alice, not Bob)
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice, // Alice sends (not Manager)
				ToAddresses: []string{suite.Bob}, // Bob receives (not Manager)
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	// This should succeed because:
	// - from (Alice) != initiatedBy (Manager) - satisfies RequireFromDoesNotEqualInitiatedBy
	// - to (Bob) != initiatedBy (Manager) - satisfies RequireToDoesNotEqualInitiatedBy
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "third-party delegated transfer should succeed")

	// Alice initiates own transfer - should fail
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice, // Alice initiates
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice, // Alice sends (== initiator, violates flag)
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "self-initiated transfer should fail")
}

// TestEqualityFlagsWithMint tests equality flags with mint transfers
func (suite *RequireEqualityFlagsTestSuite) TestEqualityFlagsWithMint() {
	// Create mint approval requiring self-claim only
	approval := testutil.GenerateCollectionApproval("selfmintclaim", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireToEqualsInitiatedBy:     true, // Can only mint to yourself
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Alice mints to herself - should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
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
	suite.Require().NoError(err, "self-mint should succeed")

	// Bob mints to himself - should succeed
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Bob,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "Bob's self-mint should succeed")

	// Manager trying to mint to Bob - should fail
	msg3 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg3)
	suite.Require().Error(err, "manager minting to Bob should fail")
}

// TestApprovalStructure_EqualityFlags tests that equality flags are properly stored
func (suite *RequireEqualityFlagsTestSuite) TestApprovalStructure_EqualityFlags() {
	// Create approval with all equality flags set
	approval := testutil.GenerateCollectionApproval("allflags", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		RequireToEqualsInitiatedBy:         true,
		RequireFromEqualsInitiatedBy:       true,
		RequireToDoesNotEqualInitiatedBy:   false, // Can't have both true
		RequireFromDoesNotEqualInitiatedBy: false, // Can't have both true
		OverridesFromOutgoingApprovals:     true,
		OverridesToIncomingApprovals:       true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Verify structure was saved correctly
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "allflags" {
			found = true
			criteria := app.ApprovalCriteria
			suite.Require().NotNil(criteria, "approval criteria should exist")
			suite.Require().True(criteria.RequireToEqualsInitiatedBy, "RequireToEqualsInitiatedBy should be true")
			suite.Require().True(criteria.RequireFromEqualsInitiatedBy, "RequireFromEqualsInitiatedBy should be true")
			suite.Require().False(criteria.RequireToDoesNotEqualInitiatedBy, "RequireToDoesNotEqualInitiatedBy should be false")
			suite.Require().False(criteria.RequireFromDoesNotEqualInitiatedBy, "RequireFromDoesNotEqualInitiatedBy should be false")
			break
		}
	}
	suite.Require().True(found, "approval should exist")
}
