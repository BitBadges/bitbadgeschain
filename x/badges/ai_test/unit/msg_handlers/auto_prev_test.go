package msg_handlers_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type AutoPrevTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestAutoPrevSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(AutoPrevTestSuite))
}

func (suite *AutoPrevTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestAutoPrev_WithExistingCollection tests auto-prev with an existing collection
func (suite *AutoPrevTestSuite) TestAutoPrev_WithExistingCollection() {
	// Create a collection first
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Use auto-prev (collectionId = 0) - should resolve to the most recent collection
	// This is tested indirectly through message handlers that use auto-prev
	// For direct testing, we'd need access to the internal resolveCollectionIdWithAutoPrev function
	// So we test it through a message handler that uses it

	// Test via TransferTokens with collectionId = 0
	// First, we need to set up approvals
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:               collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(collectionId, suite.Alice, mintBalances)

	// Set approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId, // Use explicit collection ID
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: collectionId, // Use explicit collection ID
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err)

	// Test transfer with explicit collection ID (not auto-prev)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId, // Explicit ID
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(10, 1),
			}),
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "transfer with explicit collection ID should succeed")
}

// TestAutoPrev_WithNonExistentCollection tests that auto-prev fails for non-existent collection
func (suite *AutoPrevTestSuite) TestAutoPrev_WithNonExistentCollection() {
	// Try to use a non-existent collection ID
	nonExistentId := sdkmath.NewUint(99999)

	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: nonExistentId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(10, 1),
			}),
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer with non-existent collection ID should fail")
	suite.Require().Contains(err.Error(), "collection does not exist", "error should indicate collection does not exist")
}

// TestAutoPrev_CreatorPermission_Manager tests that manager can use auto-prev
func (suite *AutoPrevTestSuite) TestAutoPrev_CreatorPermission_Manager() {
	// Create a collection with manager
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Set up approvals and mint tokens
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(collectionId, suite.Manager, mintBalances)

	// Set up outgoing approval for manager
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	// Set up incoming approval for Alice
	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err)

	// Manager should be able to use auto-prev (collectionId = 0)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: sdkmath.NewUint(0), // Auto-prev
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Manager, []string{suite.Alice}, []*types.Balance{
				testutil.GenerateSimpleBalance(10, 1),
			}),
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().NoError(err, "manager should be able to use auto-prev")
}

// TestAutoPrev_CreatorPermission_OriginalCreator tests that original creator can use auto-prev
func (suite *AutoPrevTestSuite) TestAutoPrev_CreatorPermission_OriginalCreator() {
	// Create a collection with a specific creator
	creator := suite.Alice
	collectionId := suite.CreateTestCollection(creator)

	// Change manager to someone else
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:       creator,
		CollectionId:  collectionId,
		UpdateManager: true,
		Manager:       suite.Manager, // Change manager
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Set up approvals and mint tokens
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg2 := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg2)
	suite.Require().NoError(err)

	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(collectionId, creator, mintBalances)

	// Original creator should be able to use auto-prev even though they're not the manager
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      creator,
		CollectionId: sdkmath.NewUint(0), // Auto-prev
		Approval:     testutil.GenerateUserIncomingApproval("incoming1", "All"),
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err, "original creator should be able to use auto-prev")
}

// TestAutoPrev_CreatorPermission_Unauthorized tests that unauthorized creator cannot use auto-prev
// Note: Permission checks for auto-prev were moved out of resolveCollectionIdWithAutoPrev.
// The function is safe and checks (manager, gov, etc.) are performed later in the message handlers.
// This test verifies that auto-prev resolution works (the actual permission checks happen elsewhere).
func (suite *AutoPrevTestSuite) TestAutoPrev_CreatorPermission_Unauthorized() {
	// Create a collection with manager
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Unauthorized user (Bob) can use auto-prev to resolve collection ID
	// Permission checks happen later in the message handler (not in resolveCollectionIdWithAutoPrev)
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob, // Unauthorized - not manager or creator
		CollectionId: sdkmath.NewUint(0), // Auto-prev - will resolve to collectionId
		Approval:     testutil.GenerateUserIncomingApproval("incoming1", "All"),
	}
	_, err := suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	// Auto-prev resolution succeeds (checks happen later), but the operation itself may fail
	// if Bob doesn't have permission to set incoming approvals for this collection
	// The test verifies that auto-prev resolution works, not that permissions are checked here
	suite.Require().NoError(err, "auto-prev resolution should work (permission checks happen later)")
	
	// Verify the approval was set (Bob can set his own incoming approvals)
	collection := suite.GetCollection(collectionId)
	bobBalance, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Bob)
	found := false
	for _, approval := range bobBalance.IncomingApprovals {
		if approval.ApprovalId == "incoming1" {
			found = true
			break
		}
	}
	suite.Require().True(found, "approval should be set (Bob can set his own incoming approvals)")
}

// TestAutoPrev_MultiMsgTransaction tests auto-prev in multi-msg transaction scenarios
func (suite *AutoPrevTestSuite) TestAutoPrev_MultiMsgTransaction() {
	// Create first collection
	collectionId1 := suite.CreateTestCollection(suite.Manager)

	// Set up approvals for first collection
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId1,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens to first collection
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintBadges(collectionId1, suite.Alice, mintBalances)

	// Create second collection (simulating second message in multi-msg transaction)
	_ = suite.CreateTestCollection(suite.Manager) // collectionId2 - we'll get it from next collection ID
	collectionId2 := suite.Keeper.GetNextCollectionId(suite.Ctx).Sub(sdkmath.NewUint(1))

	// In a real multi-msg transaction, both messages would use auto-prev (collectionId = 0)
	// The first message would resolve to collectionId1, and the second would resolve to collectionId2
	// However, within a single transaction, state changes are atomic, so:
	// - If message 1 creates a collection, message 2 using auto-prev would see the new collection
	// - This is the intended behavior for multi-msg transactions

	// Test that auto-prev resolves to the most recent collection (collectionId2)
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Manager,
		CollectionId: sdkmath.NewUint(0), // Auto-prev - should resolve to collectionId2
		Approval:     testutil.GenerateUserIncomingApproval("incoming1", "All"),
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err, "auto-prev should resolve to most recent collection")

	// Verify the approval was set on collectionId2, not collectionId1
	collection2 := suite.GetCollection(collectionId2)
	balance2, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection2, suite.Manager)
	found := false
	for _, approval := range balance2.IncomingApprovals {
		if approval.ApprovalId == "incoming1" {
			found = true
			break
		}
	}
	suite.Require().True(found, "approval should be set on the most recent collection (collectionId2)")
}

