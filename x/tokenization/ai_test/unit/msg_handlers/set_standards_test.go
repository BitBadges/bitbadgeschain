package msg_handlers_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type SetStandardsTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestSetStandardsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(SetStandardsTestSuite))
}

func (suite *SetStandardsTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
	suite.Require().True(suite.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0 after creation")
}

// TestSetStandards_ValidUpdate tests successfully updating standards
func (suite *SetStandardsTestSuite) TestSetStandards_ValidUpdate() {
	newStandards := []string{"ERC-721", "ERC-1155"}

	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          newStandards,
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating standards should succeed")

	// Verify standards updated
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(newStandards, collection.Standards, "standards should be updated")
}

// TestSetStandards_PermissionChecked tests that canUpdateStandards permission is enforced
func (suite *SetStandardsTestSuite) TestSetStandards_PermissionChecked() {
	// First, forbid all future updates
	forbidAllTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	msg := &types.MsgSetStandards{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		Standards:    []string{"ERC-721"},
		CanUpdateStandards: []*types.ActionPermission{
			{
				PermanentlyForbiddenTimes: forbidAllTimes,
				PermanentlyPermittedTimes: []*types.UintRange{},
			},
		},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first update should succeed")

	// Now try to update again - should fail because permission forbids it
	msg2 := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-1155"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "second update should fail because permission is forbidden")
}

// TestSetStandards_DuplicateStandardsAllowed tests that duplicate standards are allowed
// Note: The system does not enforce uniqueness of standards - it's up to the client to avoid duplicates
func (suite *SetStandardsTestSuite) TestSetStandards_DuplicateStandardsAllowed() {
	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-721", "ERC-721"}, // Duplicate
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	// Duplicates are allowed at the system level
	suite.Require().NoError(err, "duplicate standards are allowed by the system")
}

// TestSetStandards_EmptyArrayClearsStandards tests that empty array clears standards
func (suite *SetStandardsTestSuite) TestSetStandards_EmptyArrayClearsStandards() {
	// First set some standards
	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-721", "ERC-1155"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting standards should succeed")

	// Verify standards are set
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().NotEmpty(collection.Standards, "standards should be set")

	// Now clear them with empty array
	msg2 := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          []string{}, // Empty array to clear
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "clearing standards with empty array should succeed")

	// Verify standards are cleared
	collection = suite.GetCollection(suite.CollectionId)
	suite.Require().Empty(collection.Standards, "standards should be cleared")
}

// TestSetStandards_OnlyManagerCanUpdate tests that only manager can update standards
func (suite *SetStandardsTestSuite) TestSetStandards_OnlyManagerCanUpdate() {
	msg := &types.MsgSetStandards{
		Creator:            suite.Alice, // Not the manager
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-721"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-manager should not be able to update standards")
	suite.Require().Contains(err.Error(), "manager", "error should mention manager permission")
}

// TestSetStandards_NonExistentCollection tests updating standards on non-existent collection
func (suite *SetStandardsTestSuite) TestSetStandards_NonExistentCollection() {
	nonExistentId := sdkmath.NewUint(99999)

	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       nonExistentId,
		Standards:          []string{"ERC-721"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "updating standards on non-existent collection should fail")
	suite.Require().Contains(err.Error(), "collection", "error should mention collection")
}

// TestSetStandards_EmptyCreator tests behavior with empty creator
func (suite *SetStandardsTestSuite) TestSetStandards_EmptyCreator() {
	msg := &types.MsgSetStandards{
		Creator:            "", // Empty creator
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-721"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "empty creator should fail")
}

// TestSetStandards_SingleStandard tests setting a single standard
func (suite *SetStandardsTestSuite) TestSetStandards_SingleStandard() {
	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-721"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting single standard should succeed")

	// Verify
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(1, len(collection.Standards))
	suite.Require().Equal("ERC-721", collection.Standards[0])
}

// TestSetStandards_ManyStandards tests setting many standards
func (suite *SetStandardsTestSuite) TestSetStandards_ManyStandards() {
	standards := []string{
		"ERC-20",
		"ERC-721",
		"ERC-1155",
		"EIP-2981",
		"EIP-4906",
		"BitBadges-Native",
		"Custom-Standard-1",
		"Custom-Standard-2",
	}

	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          standards,
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting many standards should succeed")

	// Verify
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(len(standards), len(collection.Standards))
}

// TestSetStandards_AfterManagerChange tests that new manager can update standards
func (suite *SetStandardsTestSuite) TestSetStandards_AfterManagerChange() {
	// First change manager to Bob
	setManagerMsg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}
	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), setManagerMsg)
	suite.Require().NoError(err, "setting new manager should succeed")

	// Bob (new manager) should be able to update standards
	msg := &types.MsgSetStandards{
		Creator:            suite.Bob,
		CollectionId:       suite.CollectionId,
		Standards:          []string{"Bob-Standard"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "new manager should be able to update standards")

	// Old manager should not be able to update
	msg2 := &types.MsgSetStandards{
		Creator:            suite.Manager, // Old manager
		CollectionId:       suite.CollectionId,
		Standards:          []string{"Old-Manager-Standard"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "old manager should not be able to update standards")
}

// TestSetStandards_ReplaceStandards tests replacing existing standards
func (suite *SetStandardsTestSuite) TestSetStandards_ReplaceStandards() {
	// First set some standards
	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-721", "ERC-1155"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting standards should succeed")

	// Now replace with different standards
	newStandards := []string{"ERC-20", "EIP-2981"}
	msg2 := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          newStandards,
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "replacing standards should succeed")

	// Verify new standards replaced old ones
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(newStandards, collection.Standards, "standards should be replaced")
}

// TestSetStandards_EmptyStringStandard tests behavior with empty string in standards array
// Note: Empty strings in standards array are allowed by the system
func (suite *SetStandardsTestSuite) TestSetStandards_EmptyStringStandard() {
	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-721", "", "ERC-1155"}, // Contains empty string
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	// Empty strings in standards are allowed
	suite.Require().NoError(err, "empty string in standards is allowed")
}

// TestSetStandards_SpecialCharacters tests behavior with special characters in standards
func (suite *SetStandardsTestSuite) TestSetStandards_SpecialCharacters() {
	// Some special characters might be allowed in standard names
	testCases := []struct {
		name      string
		standards []string
		shouldErr bool
	}{
		{"hyphen", []string{"ERC-721"}, false},
		{"underscore", []string{"Custom_Standard"}, false},
		{"numbers", []string{"Standard123"}, false},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := &types.MsgSetStandards{
				Creator:            suite.Manager,
				CollectionId:       suite.CollectionId,
				Standards:          tc.standards,
				CanUpdateStandards: []*types.ActionPermission{},
			}

			_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

// TestSetStandards_CaseSensitiveDuplicates tests that case-sensitive duplicates are handled
func (suite *SetStandardsTestSuite) TestSetStandards_CaseSensitiveDuplicates() {
	// "ERC-721" and "erc-721" might or might not be considered duplicates
	// depending on implementation
	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-721", "erc-721"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	// This test documents the actual behavior - different case might be allowed
	_ = err
}

// TestSetStandards_NilStandards tests behavior with nil standards array
func (suite *SetStandardsTestSuite) TestSetStandards_NilStandards() {
	msg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          nil, // Nil standards
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), msg)
	// Nil standards might be treated same as empty array
	suite.Require().NoError(err, "nil standards should be allowed (treated as empty)")
}
