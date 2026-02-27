package msg_handlers_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type SetManagerTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestSetManagerSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(SetManagerTestSuite))
}

func (suite *SetManagerTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
	suite.Require().True(suite.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0 after creation")
}

// TestSetManager_Success tests successfully setting a new manager
func (suite *SetManagerTestSuite) TestSetManager_Success() {
	// Verify initial manager
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(suite.Manager, collection.Manager, "initial manager should be set")

	// Set new manager
	msg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{}, // Allow future updates
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting new manager should succeed")

	// Verify manager changed
	collection = suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(suite.Bob, collection.Manager, "manager should be updated to Bob")
}

// TestSetManager_OnlyCurrentManagerCanSet tests that only the current manager can change manager
func (suite *SetManagerTestSuite) TestSetManager_OnlyCurrentManagerCanSet() {
	// Try to set manager as non-manager (Alice is not the manager)
	msg := &types.MsgSetManager{
		Creator:          suite.Alice, // Not the manager
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-manager should not be able to set manager")
	suite.Require().Contains(err.Error(), "manager", "error should mention manager permission")
}

// TestSetManager_PermissionCanUpdateManagerChecked tests that canUpdateManager permission is enforced
func (suite *SetManagerTestSuite) TestSetManager_PermissionCanUpdateManagerChecked() {
	// First, set the permission to forbid all future updates
	forbidAllTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	msg := &types.MsgSetManager{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		Manager:      suite.Bob,
		CanUpdateManager: []*types.ActionPermission{
			{
				PermanentlyForbiddenTimes: forbidAllTimes,
				PermanentlyPermittedTimes: []*types.UintRange{},
			},
		},
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first update with forbidden permission should succeed")

	// Now try to update manager again - should fail because permission forbids it
	msg2 := &types.MsgSetManager{
		Creator:          suite.Bob, // Bob is now the manager
		CollectionId:     suite.CollectionId,
		Manager:          suite.Alice,
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "update should fail because permission is now forbidden")
}

// TestSetManager_EmptyStringClearsManager tests that setting empty string clears the manager
func (suite *SetManagerTestSuite) TestSetManager_EmptyStringClearsManager() {
	// Verify initial manager
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().NotEmpty(collection.Manager, "initial manager should be set")

	// Set manager to empty string
	msg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          "", // Clear manager
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "clearing manager should succeed")

	// Verify manager is cleared
	collection = suite.GetCollection(suite.CollectionId)
	suite.Require().Empty(collection.Manager, "manager should be cleared")
}

// TestSetManager_InvalidAddressRejected tests that invalid addresses are rejected
func (suite *SetManagerTestSuite) TestSetManager_InvalidAddressRejected() {
	// Try to set manager to an invalid address
	msg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          "invalid_address_format",
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "setting invalid address should fail")
}

// TestSetManager_NonExistentCollection tests setting manager on non-existent collection
func (suite *SetManagerTestSuite) TestSetManager_NonExistentCollection() {
	nonExistentId := sdkmath.NewUint(99999)

	msg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     nonExistentId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "setting manager on non-existent collection should fail")
	suite.Require().Contains(err.Error(), "collection", "error should mention collection")
}

// TestSetManager_NewManagerCanManage tests that new manager can perform manager actions
func (suite *SetManagerTestSuite) TestSetManager_NewManagerCanManage() {
	// Set Bob as the new manager
	msg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{}, // Allow future updates
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting new manager should succeed")

	// Now Bob should be able to set a new manager
	msg2 := &types.MsgSetManager{
		Creator:          suite.Bob, // Bob is now the manager
		CollectionId:     suite.CollectionId,
		Manager:          suite.Charlie,
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "new manager should be able to update manager")

	// Verify manager is now Charlie
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(suite.Charlie, collection.Manager, "manager should be updated to Charlie")
}

// TestSetManager_OldManagerCannotManageAfterTransfer tests that old manager loses privileges
func (suite *SetManagerTestSuite) TestSetManager_OldManagerCannotManageAfterTransfer() {
	// Set Bob as the new manager
	msg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting new manager should succeed")

	// Old manager (suite.Manager) should not be able to update manager anymore
	msg2 := &types.MsgSetManager{
		Creator:          suite.Manager, // Old manager
		CollectionId:     suite.CollectionId,
		Manager:          suite.Alice,
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "old manager should not be able to update manager")
}

// TestSetManager_InvalidCreatorAddress tests that invalid creator address is rejected
func (suite *SetManagerTestSuite) TestSetManager_InvalidCreatorAddress() {
	msg := &types.MsgSetManager{
		Creator:          "", // Invalid empty creator
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "empty creator should fail")
}

// TestSetManager_SelfAssignment tests that manager can set themselves as manager (no-op)
func (suite *SetManagerTestSuite) TestSetManager_SelfAssignment() {
	msg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Manager, // Same as current manager
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "self-assignment should succeed")

	// Verify manager is still the same
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(suite.Manager, collection.Manager, "manager should remain unchanged")
}
