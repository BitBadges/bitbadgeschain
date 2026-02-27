package msg_handlers_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type SetCustomDataTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestSetCustomDataSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(SetCustomDataTestSuite))
}

func (suite *SetCustomDataTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
	suite.Require().True(suite.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0 after creation")
}

// TestSetCustomData_ValidUpdate tests successfully updating custom data
func (suite *SetCustomDataTestSuite) TestSetCustomData_ValidUpdate() {
	newCustomData := "new custom data value"

	msg := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CustomData:          newCustomData,
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating custom data should succeed")

	// Verify custom data updated
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(newCustomData, collection.CustomData, "custom data should be updated")
}

// TestSetCustomData_PermissionChecked tests that canUpdateCustomData permission is enforced
func (suite *SetCustomDataTestSuite) TestSetCustomData_PermissionChecked() {
	// First, forbid all future updates
	forbidAllTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	msg := &types.MsgSetCustomData{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		CustomData:   "first update",
		CanUpdateCustomData: []*types.ActionPermission{
			{
				PermanentlyForbiddenTimes: forbidAllTimes,
				PermanentlyPermittedTimes: []*types.UintRange{},
			},
		},
	}

	_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first update should succeed")

	// Now try to update again - should fail because permission forbids it
	msg2 := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CustomData:          "second update",
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "second update should fail because permission is forbidden")
}

// TestSetCustomData_EmptyStringAllowed tests that empty string is allowed for custom data
func (suite *SetCustomDataTestSuite) TestSetCustomData_EmptyStringAllowed() {
	// First set some custom data
	msg := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CustomData:          "some data",
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting custom data should succeed")

	// Now clear it with empty string
	msg2 := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CustomData:          "", // Empty string to clear
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "clearing custom data with empty string should succeed")

	// Verify custom data is cleared
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Empty(collection.CustomData, "custom data should be cleared")
}

// TestSetCustomData_OnlyManagerCanUpdate tests that only manager can update custom data
func (suite *SetCustomDataTestSuite) TestSetCustomData_OnlyManagerCanUpdate() {
	msg := &types.MsgSetCustomData{
		Creator:             suite.Alice, // Not the manager
		CollectionId:        suite.CollectionId,
		CustomData:          "alice's data",
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-manager should not be able to update custom data")
	suite.Require().Contains(err.Error(), "manager", "error should mention manager permission")
}

// TestSetCustomData_NonExistentCollection tests updating custom data on non-existent collection
func (suite *SetCustomDataTestSuite) TestSetCustomData_NonExistentCollection() {
	nonExistentId := sdkmath.NewUint(99999)

	msg := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        nonExistentId,
		CustomData:          "some data",
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "updating custom data on non-existent collection should fail")
	suite.Require().Contains(err.Error(), "collection", "error should mention collection")
}

// TestSetCustomData_EmptyCreator tests behavior with empty creator
func (suite *SetCustomDataTestSuite) TestSetCustomData_EmptyCreator() {
	msg := &types.MsgSetCustomData{
		Creator:             "", // Empty creator
		CollectionId:        suite.CollectionId,
		CustomData:          "some data",
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "empty creator should fail")
}

// TestSetCustomData_LongCustomData tests behavior with long custom data
func (suite *SetCustomDataTestSuite) TestSetCustomData_LongCustomData() {
	// Create a moderately long custom data string
	longCustomData := ""
	for i := 0; i < 100; i++ {
		longCustomData += "data"
	}

	msg := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CustomData:          longCustomData,
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	// Long custom data should be allowed (within reasonable limits)
	suite.Require().NoError(err, "moderately long custom data should succeed")
}

// TestSetCustomData_MultipleUpdates tests multiple sequential custom data updates
func (suite *SetCustomDataTestSuite) TestSetCustomData_MultipleUpdates() {
	updates := []string{"update1", "update2", "update3", "update4", "update5"}

	for _, data := range updates {
		msg := &types.MsgSetCustomData{
			Creator:             suite.Manager,
			CollectionId:        suite.CollectionId,
			CustomData:          data,
			CanUpdateCustomData: []*types.ActionPermission{},
		}

		_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "update should succeed")

		// Verify update
		collection := suite.GetCollection(suite.CollectionId)
		suite.Require().Equal(data, collection.CustomData)
	}
}

// TestSetCustomData_SpecialCharacters tests behavior with special characters in custom data
func (suite *SetCustomDataTestSuite) TestSetCustomData_SpecialCharacters() {
	testCases := []string{
		`{"key": "value"}`,                  // JSON
		`<xml>data</xml>`,                   // XML-like
		"data with\nnewlines",               // Newlines
		"data with\ttabs",                   // Tabs
		"unicode: \u4e2d\u6587",             // Unicode
		"special: !@#$%^&*()_+-=[]{}|;':\"", // Special chars
	}

	for _, data := range testCases {
		msg := &types.MsgSetCustomData{
			Creator:             suite.Manager,
			CollectionId:        suite.CollectionId,
			CustomData:          data,
			CanUpdateCustomData: []*types.ActionPermission{},
		}

		_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "special characters should be allowed: %s", data)
	}
}

// TestSetCustomData_AfterManagerChange tests that new manager can update custom data
func (suite *SetCustomDataTestSuite) TestSetCustomData_AfterManagerChange() {
	// First change manager to Bob
	setManagerMsg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}
	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), setManagerMsg)
	suite.Require().NoError(err, "setting new manager should succeed")

	// Bob (new manager) should be able to update custom data
	msg := &types.MsgSetCustomData{
		Creator:             suite.Bob,
		CollectionId:        suite.CollectionId,
		CustomData:          "bob's custom data",
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "new manager should be able to update custom data")

	// Old manager should not be able to update
	msg2 := &types.MsgSetCustomData{
		Creator:             suite.Manager, // Old manager
		CollectionId:        suite.CollectionId,
		CustomData:          "old manager's data",
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "old manager should not be able to update custom data")
}

// TestSetCustomData_SameValueUpdate tests updating custom data to the same value (no-op)
func (suite *SetCustomDataTestSuite) TestSetCustomData_SameValueUpdate() {
	customData := "same value"

	// First update
	msg := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CustomData:          customData,
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first update should succeed")

	// Same value update (no-op)
	msg2 := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CustomData:          customData, // Same value
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "same value update should succeed")
}

// TestSetCustomData_WhitespaceOnlyData tests behavior with whitespace-only custom data
func (suite *SetCustomDataTestSuite) TestSetCustomData_WhitespaceOnlyData() {
	msg := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CustomData:          "   ",
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), msg)
	// Whitespace-only data should be allowed
	suite.Require().NoError(err, "whitespace-only data should be allowed")
}
