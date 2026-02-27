package msg_handlers_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type SetCollectionApprovalsTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestSetCollectionApprovalsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(SetCollectionApprovalsTestSuite))
}

func (suite *SetCollectionApprovalsTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
	suite.Require().True(suite.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0 after creation")
}

// TestSetCollectionApprovals_ValidUpdate tests successfully updating collection approvals
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_ValidUpdate() {
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All"),
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating collection approvals should succeed")

	// Verify approvals updated
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().NotEmpty(collection.CollectionApprovals, "collection approvals should be set")
	suite.Require().Equal("approval1", collection.CollectionApprovals[0].ApprovalId)
}

// TestSetCollectionApprovals_PermissionChecked tests that canUpdateCollectionApprovals permission is enforced
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_PermissionChecked() {
	// First, set some approvals with permission that forbids future updates for this approval ID
	forbidAllTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All"),
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CollectionApprovals: approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{
			{
				ApprovalId: "approval1",
				FromListId: "AllWithoutMint",
				ToListId:   "All",
				TransferTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				OwnershipTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				InitiatedByListId:         "All",
				PermanentlyForbiddenTimes: forbidAllTimes,
				PermanentlyPermittedTimes: []*types.UintRange{},
			},
		},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first update should succeed")

	// Now try to update the same approval - should fail
	approvals2 := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "All", "All"), // Changed from/to
	}

	msg2 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals2,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err = suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "update should fail because permission is forbidden")
}

// TestSetCollectionApprovals_DuplicateApprovalIDsRejected tests that duplicate approval IDs are rejected
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_DuplicateApprovalIDsRejected() {
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All"),
		testutil.GenerateCollectionApproval("approval1", "All", "All"), // Same ID
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "duplicate approval IDs should be rejected")
}

// TestSetCollectionApprovals_ApprovalStructureValidated tests that approval structure is validated
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_ApprovalStructureValidated() {
	// Create an approval with invalid structure (invalid token ID range)
	approvals := []*types.CollectionApproval{
		{
			ApprovalId:        "invalid_approval",
			FromListId:        "AllWithoutMint",
			ToListId:          "All",
			InitiatedByListId: "All",
			TransferTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(1)}, // Invalid: start > end
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
			ApprovalCriteria: &types.ApprovalCriteria{},
			Version:          sdkmath.NewUint(0),
		},
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "invalid approval structure should be rejected")
}

// TestSetCollectionApprovals_OnlyManagerCanUpdate tests that only manager can update collection approvals
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_OnlyManagerCanUpdate() {
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All"),
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Alice, // Not the manager
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-manager should not be able to update collection approvals")
	suite.Require().Contains(err.Error(), "manager", "error should mention manager permission")
}

// TestSetCollectionApprovals_NonExistentCollection tests updating approvals on non-existent collection
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_NonExistentCollection() {
	nonExistentId := sdkmath.NewUint(99999)

	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All"),
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 nonExistentId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "updating approvals on non-existent collection should fail")
}

// TestSetCollectionApprovals_EmptyCreator tests behavior with empty creator
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_EmptyCreator() {
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All"),
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      "", // Empty creator
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "empty creator should fail")
}

// TestSetCollectionApprovals_EmptyApprovals tests clearing all approvals
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_EmptyApprovals() {
	// First set some approvals
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All"),
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting approvals should succeed")

	// Now clear them
	msg2 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          []*types.CollectionApproval{}, // Empty
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err = suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "clearing approvals should succeed")

	// Verify approvals are cleared
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Empty(collection.CollectionApprovals, "collection approvals should be cleared")
}

// TestSetCollectionApprovals_MultipleApprovals tests setting multiple approvals
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_MultipleApprovals() {
	// Note: "All" list cannot be used with Mint approvals in same set because
	// "All" includes Mint address. Use "AllWithoutMint" for non-mint approvals.
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "AllWithoutMint")
	mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals = true
	mintApproval.ApprovalCriteria.OverridesToIncomingApprovals = true

	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "AllWithoutMint"),
		testutil.GenerateCollectionApproval("approval2", "AllWithoutMint", "AllWithoutMint"),
		mintApproval,
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting multiple approvals should succeed")

	// Verify all approvals are set
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(3, len(collection.CollectionApprovals), "should have 3 approvals")
}

// TestSetCollectionApprovals_AfterManagerChange tests that new manager can update approvals
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_AfterManagerChange() {
	// First change manager to Bob
	setManagerMsg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}
	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), setManagerMsg)
	suite.Require().NoError(err, "setting new manager should succeed")

	// Bob (new manager) should be able to update collection approvals
	approvals := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("bob_approval", "AllWithoutMint", "All"),
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Bob,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err = suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "new manager should be able to update collection approvals")

	// Old manager should not be able to update
	msg2 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager, // Old manager
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err = suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "old manager should not be able to update collection approvals")
}

// TestSetCollectionApprovals_MintApproval tests setting mint approval with required flags
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_MintApproval() {
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals = true
	mintApproval.ApprovalCriteria.OverridesToIncomingApprovals = true

	approvals := []*types.CollectionApproval{mintApproval}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting mint approval should succeed")

	// Verify mint approval is set correctly
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().NotEmpty(collection.CollectionApprovals)
	suite.Require().Equal(types.MintAddress, collection.CollectionApprovals[0].FromListId)
}

// TestSetCollectionApprovals_EmptyApprovalId tests behavior with empty approval ID
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_EmptyApprovalId() {
	approvals := []*types.CollectionApproval{
		{
			ApprovalId:        "", // Empty approval ID
			FromListId:        "AllWithoutMint",
			ToListId:          "All",
			InitiatedByListId: "All",
			TransferTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
			ApprovalCriteria: &types.ApprovalCriteria{},
			Version:          sdkmath.NewUint(0),
		},
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "empty approval ID should be rejected")
}

// TestSetCollectionApprovals_InvalidListId tests behavior with invalid list IDs
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_InvalidListId() {
	// Using a non-existent list ID that isn't a reserved ID
	approvals := []*types.CollectionApproval{
		{
			ApprovalId:        "approval1",
			FromListId:        "NonExistentList", // Not a valid reserved ID or created list
			ToListId:          "All",
			InitiatedByListId: "All",
			TransferTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
			OwnershipTimes: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
			},
			ApprovalCriteria: &types.ApprovalCriteria{},
			Version:          sdkmath.NewUint(0),
		},
	}

	msg := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "invalid list ID should be rejected")
}

// TestSetCollectionApprovals_ReplaceApprovals tests replacing existing approvals
func (suite *SetCollectionApprovalsTestSuite) TestSetCollectionApprovals_ReplaceApprovals() {
	// First set some approvals - use AllWithoutMint to avoid Mint address conflicts
	approvals1 := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "AllWithoutMint"),
		testutil.GenerateCollectionApproval("approval2", "AllWithoutMint", "AllWithoutMint"),
	}

	msg1 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals1,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err := suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err, "first update should succeed")

	// Now replace with different approvals
	approvals2 := []*types.CollectionApproval{
		testutil.GenerateCollectionApproval("new_approval", "AllWithoutMint", "All"),
	}

	msg2 := &types.MsgSetCollectionApprovals{
		Creator:                      suite.Manager,
		CollectionId:                 suite.CollectionId,
		CollectionApprovals:          approvals2,
		CanUpdateCollectionApprovals: []*types.CollectionApprovalPermission{},
	}

	_, err = suite.MsgServer.SetCollectionApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "replacing approvals should succeed")

	// Verify new approvals replaced old ones
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(1, len(collection.CollectionApprovals), "should have 1 approval after replacement")
	suite.Require().Equal("new_approval", collection.CollectionApprovals[0].ApprovalId)
}
