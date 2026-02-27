package msg_handlers_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type UpdateUserApprovalsTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestUpdateUserApprovalsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(UpdateUserApprovalsTestSuite))
}

func (suite *UpdateUserApprovalsTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestUpdateUserApprovals_OutgoingApprovals tests updating outgoing approvals
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_OutgoingApprovals() {
	// Create outgoing approval
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")

	msg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoingApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating outgoing approvals should succeed")

	// Verify approval was saved
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.OutgoingApprovals), "should have one outgoing approval")
	suite.Require().Equal("outgoing1", balance.OutgoingApprovals[0].ApprovalId)
	suite.Require().Equal("All", balance.OutgoingApprovals[0].ToListId)
}

// TestUpdateUserApprovals_IncomingApprovals tests updating incoming approvals
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_IncomingApprovals() {
	// Create incoming approval
	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")

	msg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{incomingApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating incoming approvals should succeed")

	// Verify approval was saved
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.IncomingApprovals), "should have one incoming approval")
	suite.Require().Equal("incoming1", balance.IncomingApprovals[0].ApprovalId)
	suite.Require().Equal("All", balance.IncomingApprovals[0].FromListId)
}

// TestUpdateUserApprovals_MultipleApprovals tests updating multiple approvals at once
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_MultipleApprovals() {
	// Create multiple outgoing approvals
	outgoing1 := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	outgoing2 := testutil.GenerateUserOutgoingApproval("outgoing2", suite.Bob)

	// Create multiple incoming approvals
	incoming1 := testutil.GenerateUserIncomingApproval("incoming1", "All")
	incoming2 := testutil.GenerateUserIncomingApproval("incoming2", suite.Bob)

	msg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoing1, outgoing2},
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{incoming1, incoming2},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating multiple approvals should succeed")

	// Verify all approvals were saved
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(2, len(balance.OutgoingApprovals), "should have two outgoing approvals")
	suite.Require().Equal(2, len(balance.IncomingApprovals), "should have two incoming approvals")
}

// TestUpdateUserApprovals_AutoApproveFlags tests updating auto-approve flags
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_AutoApproveFlags() {
	// Update auto-approve self-initiated incoming transfers
	msg := &types.MsgUpdateUserApprovals{
		Creator:                                      suite.Alice,
		CollectionId:                                 suite.CollectionId,
		UpdateAutoApproveSelfInitiatedIncomingTransfers: true,
		AutoApproveSelfInitiatedIncomingTransfers:       true,
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating auto-approve self-initiated incoming should succeed")

	// Verify flag was saved
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().True(balance.AutoApproveSelfInitiatedIncomingTransfers, "auto-approve self-initiated incoming should be true")
}

// TestUpdateUserApprovals_AutoApproveSelfInitiatedOutgoing tests auto-approve self-initiated outgoing
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_AutoApproveSelfInitiatedOutgoing() {
	msg := &types.MsgUpdateUserApprovals{
		Creator:                                      suite.Alice,
		CollectionId:                                 suite.CollectionId,
		UpdateAutoApproveSelfInitiatedOutgoingTransfers: true,
		AutoApproveSelfInitiatedOutgoingTransfers:       true,
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating auto-approve self-initiated outgoing should succeed")

	// Verify flag was saved
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().True(balance.AutoApproveSelfInitiatedOutgoingTransfers, "auto-approve self-initiated outgoing should be true")
}

// TestUpdateUserApprovals_AutoApproveAllIncoming tests auto-approve all incoming transfers
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_AutoApproveAllIncoming() {
	msg := &types.MsgUpdateUserApprovals{
		Creator:                               suite.Alice,
		CollectionId:                          suite.CollectionId,
		UpdateAutoApproveAllIncomingTransfers: true,
		AutoApproveAllIncomingTransfers:       true,
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating auto-approve all incoming should succeed")

	// Verify flag was saved
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().True(balance.AutoApproveAllIncomingTransfers, "auto-approve all incoming should be true")
}

// TestUpdateUserApprovals_UserPermissions tests updating user permissions
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_UserPermissions() {
	// Create user permissions that lock outgoing approval updates
	userPermissions := &types.UserPermissions{
		CanUpdateOutgoingApprovals: []*types.UserOutgoingApprovalPermission{
			{
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
				ApprovalId:                "All",
				PermanentlyPermittedTimes: []*types.UintRange{},
				PermanentlyForbiddenTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
			},
		},
	}

	msg := &types.MsgUpdateUserApprovals{
		Creator:               suite.Alice,
		CollectionId:          suite.CollectionId,
		UpdateUserPermissions: true,
		UserPermissions:       userPermissions,
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating user permissions should succeed")

	// Verify permissions were saved
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().NotNil(balance.UserPermissions, "user permissions should exist")
	suite.Require().Equal(1, len(balance.UserPermissions.CanUpdateOutgoingApprovals), "should have one outgoing approval permission")
}

// TestUpdateUserApprovals_PermissionEnforcement tests that permissions are enforced
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_PermissionEnforcement() {
	// First, set permissions that forbid updating outgoing approvals
	userPermissions := &types.UserPermissions{
		CanUpdateOutgoingApprovals: []*types.UserOutgoingApprovalPermission{
			{
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
				ApprovalId:                "All",
				PermanentlyPermittedTimes: []*types.UintRange{},
				PermanentlyForbiddenTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
			},
		},
	}

	setPermissionsMsg := &types.MsgUpdateUserApprovals{
		Creator:               suite.Alice,
		CollectionId:          suite.CollectionId,
		UpdateUserPermissions: true,
		UserPermissions:       userPermissions,
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setPermissionsMsg)
	suite.Require().NoError(err, "setting initial permissions should succeed")

	// Now try to update outgoing approvals - should fail due to permissions
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	updateApprovalsMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoingApproval},
	}

	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateApprovalsMsg)
	suite.Require().Error(err, "updating outgoing approvals should fail when permissions forbid it")
}

// TestUpdateUserApprovals_ArchivedCollectionFails tests that updates fail for archived collections
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_ArchivedCollectionFails() {
	// Archive the collection
	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err, "archiving collection should succeed")

	// Try to update user approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	updateMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoingApproval},
	}

	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().Error(err, "updating user approvals on archived collection should fail")
}

// TestUpdateUserApprovals_InvalidCollectionFails tests that updates fail for non-existent collections
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_InvalidCollectionFails() {
	invalidCollectionId := sdkmath.NewUint(99999)

	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	msg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            invalidCollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoingApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "updating user approvals on non-existent collection should fail")
}

// TestUpdateUserApprovals_ReplaceExistingApprovals tests that updating approvals replaces existing ones
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_ReplaceExistingApprovals() {
	// Set initial outgoing approval
	outgoing1 := testutil.GenerateUserOutgoingApproval("outgoing1", suite.Bob)
	msg1 := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoing1},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Verify initial state
	balance1 := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance1.OutgoingApprovals))
	suite.Require().Equal("outgoing1", balance1.OutgoingApprovals[0].ApprovalId)

	// Replace with different approval
	outgoing2 := testutil.GenerateUserOutgoingApproval("outgoing2", suite.Charlie)
	msg2 := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoing2},
	}

	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Verify replacement
	balance2 := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance2.OutgoingApprovals))
	suite.Require().Equal("outgoing2", balance2.OutgoingApprovals[0].ApprovalId)
}

// TestUpdateUserApprovals_ClearApprovals tests clearing all approvals
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_ClearApprovals() {
	// Set initial approvals
	outgoing := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	incoming := testutil.GenerateUserIncomingApproval("incoming1", "All")
	msg1 := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoing},
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{incoming},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Clear all approvals by setting empty arrays
	msg2 := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{},
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{},
	}

	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Verify approvals are cleared
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balance.OutgoingApprovals), "outgoing approvals should be empty")
	suite.Require().Equal(0, len(balance.IncomingApprovals), "incoming approvals should be empty")
}

// TestUpdateUserApprovals_ApprovalVersioning tests that approval versions increment correctly
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_ApprovalVersioning() {
	// Create initial approval
	outgoing := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	msg1 := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoing},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Get initial version
	balance1 := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance1.OutgoingApprovals))
	initialVersion := balance1.OutgoingApprovals[0].Version

	// Update the same approval with different settings
	outgoing2 := testutil.GenerateUserOutgoingApproval("outgoing1", suite.Bob) // Same ID, different recipient
	msg2 := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoing2},
	}

	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Verify version incremented
	balance2 := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance2.OutgoingApprovals))
	newVersion := balance2.OutgoingApprovals[0].Version

	// Version should have incremented since the approval changed
	suite.Require().True(newVersion.GT(initialVersion), "version should increment when approval changes")
}

// TestUpdateUserApprovals_PartialUpdate tests that partial updates work correctly
func (suite *UpdateUserApprovalsTestSuite) TestUpdateUserApprovals_PartialUpdate() {
	// Set initial outgoing and incoming approvals
	outgoing := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	incoming := testutil.GenerateUserIncomingApproval("incoming1", "All")
	msg1 := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoing},
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{incoming},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err)

	// Now only update outgoing approvals (incoming should remain unchanged)
	newOutgoing := testutil.GenerateUserOutgoingApproval("outgoing2", suite.Bob)
	msg2 := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{newOutgoing},
		UpdateIncomingApprovals: false, // Not updating incoming
	}

	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err)

	// Verify outgoing was updated but incoming remains
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.OutgoingApprovals))
	suite.Require().Equal("outgoing2", balance.OutgoingApprovals[0].ApprovalId, "outgoing should be updated")
	suite.Require().Equal(1, len(balance.IncomingApprovals))
	suite.Require().Equal("incoming1", balance.IncomingApprovals[0].ApprovalId, "incoming should remain unchanged")
}
