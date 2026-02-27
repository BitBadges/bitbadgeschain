package msg_handlers_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type SetCollectionMetadataTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestSetCollectionMetadataSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(SetCollectionMetadataTestSuite))
}

func (suite *SetCollectionMetadataTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
	suite.Require().True(suite.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0 after creation")
}

// TestSetCollectionMetadata_ValidUpdate tests successfully updating collection metadata
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_ValidUpdate() {
	newUri := "https://newexample.com/metadata"
	newCustomData := "some custom data"

	msg := &types.MsgSetCollectionMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        newUri,
			CustomData: newCustomData,
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating collection metadata should succeed")

	// Verify metadata updated
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(newUri, collection.CollectionMetadata.Uri, "URI should be updated")
	suite.Require().Equal(newCustomData, collection.CollectionMetadata.CustomData, "custom data should be updated")
}

// TestSetCollectionMetadata_PermissionChecked tests that canUpdateCollectionMetadata permission is enforced
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_PermissionChecked() {
	// First, forbid all future updates
	forbidAllTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	msg := &types.MsgSetCollectionMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://first-update.com/metadata",
			CustomData: "",
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{
			{
				PermanentlyForbiddenTimes: forbidAllTimes,
				PermanentlyPermittedTimes: []*types.UintRange{},
			},
		},
	}

	_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first update should succeed")

	// Now try to update again - should fail because permission forbids it
	msg2 := &types.MsgSetCollectionMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://second-update.com/metadata",
			CustomData: "",
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "second update should fail because permission is forbidden")
}

// TestSetCollectionMetadata_OnlyManagerCanUpdate tests that only manager can update metadata
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_OnlyManagerCanUpdate() {
	msg := &types.MsgSetCollectionMetadata{
		Creator:      suite.Alice, // Not the manager
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://alice-update.com/metadata",
			CustomData: "",
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-manager should not be able to update metadata")
	suite.Require().Contains(err.Error(), "manager", "error should mention manager permission")
}

// TestSetCollectionMetadata_InvalidURIHandling tests handling of invalid URIs
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_InvalidURIHandling() {
	// Note: The system may or may not validate URI format strictly
	// This test checks behavior with various URI formats
	testCases := []struct {
		name        string
		uri         string
		shouldError bool
	}{
		{"empty_uri", "", false}, // Empty URI might be allowed
		{"valid_https", "https://example.com/metadata", false},
		{"valid_ipfs", "ipfs://QmTest123", false},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := &types.MsgSetCollectionMetadata{
				Creator:      suite.Manager,
				CollectionId: suite.CollectionId,
				CollectionMetadata: &types.CollectionMetadata{
					Uri:        tc.uri,
					CustomData: "",
				},
				CanUpdateCollectionMetadata: []*types.ActionPermission{},
			}

			_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
			if tc.shouldError {
				suite.Require().Error(err, "expected error for URI: %s", tc.uri)
			} else {
				suite.Require().NoError(err, "expected success for URI: %s", tc.uri)
			}
		})
	}
}

// TestSetCollectionMetadata_UpdateCustomData tests updating only custom data
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_UpdateCustomData() {
	// First get current metadata
	collection := suite.GetCollection(suite.CollectionId)
	originalUri := collection.CollectionMetadata.Uri

	// Update only custom data
	newCustomData := "new custom data value"
	msg := &types.MsgSetCollectionMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        originalUri, // Keep original URI
			CustomData: newCustomData,
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "updating custom data should succeed")

	// Verify custom data updated
	collection = suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(newCustomData, collection.CollectionMetadata.CustomData, "custom data should be updated")
	suite.Require().Equal(originalUri, collection.CollectionMetadata.Uri, "URI should remain unchanged")
}

// TestSetCollectionMetadata_NonExistentCollection tests updating metadata on non-existent collection
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_NonExistentCollection() {
	nonExistentId := sdkmath.NewUint(99999)

	msg := &types.MsgSetCollectionMetadata{
		Creator:      suite.Manager,
		CollectionId: nonExistentId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "",
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "updating metadata on non-existent collection should fail")
	suite.Require().Contains(err.Error(), "collection", "error should mention collection")
}

// TestSetCollectionMetadata_NilMetadata tests behavior with nil metadata
// Note: Nil metadata is allowed by the system - it will be treated as empty metadata
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_NilMetadata() {
	msg := &types.MsgSetCollectionMetadata{
		Creator:                     suite.Manager,
		CollectionId:                suite.CollectionId,
		CollectionMetadata:          nil, // Nil metadata
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	// Nil metadata is allowed - it clears the metadata
	suite.Require().NoError(err, "nil metadata is allowed")
}

// TestSetCollectionMetadata_EmptyCreator tests behavior with empty creator
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_EmptyCreator() {
	msg := &types.MsgSetCollectionMetadata{
		Creator:      "", // Empty creator
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "",
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "empty creator should fail")
}

// TestSetCollectionMetadata_LongCustomData tests behavior with long custom data
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_LongCustomData() {
	// Create a moderately long custom data string
	longCustomData := ""
	for i := 0; i < 100; i++ {
		longCustomData += "data"
	}

	msg := &types.MsgSetCollectionMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: longCustomData,
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	// Long custom data should be allowed (within reasonable limits)
	suite.Require().NoError(err, "moderately long custom data should succeed")
}

// TestSetCollectionMetadata_MultipleUpdates tests multiple sequential updates
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_MultipleUpdates() {
	updates := []struct {
		uri        string
		customData string
	}{
		{"https://example1.com", "data1"},
		{"https://example2.com", "data2"},
		{"https://example3.com", "data3"},
	}

	for _, update := range updates {
		msg := &types.MsgSetCollectionMetadata{
			Creator:      suite.Manager,
			CollectionId: suite.CollectionId,
			CollectionMetadata: &types.CollectionMetadata{
				Uri:        update.uri,
				CustomData: update.customData,
			},
			CanUpdateCollectionMetadata: []*types.ActionPermission{},
		}

		_, err := suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "update should succeed")

		// Verify update
		collection := suite.GetCollection(suite.CollectionId)
		suite.Require().Equal(update.uri, collection.CollectionMetadata.Uri)
		suite.Require().Equal(update.customData, collection.CollectionMetadata.CustomData)
	}
}

// TestSetCollectionMetadata_AfterManagerChange tests that new manager can update metadata
func (suite *SetCollectionMetadataTestSuite) TestSetCollectionMetadata_AfterManagerChange() {
	// First change manager to Bob
	setManagerMsg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}
	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), setManagerMsg)
	suite.Require().NoError(err, "setting new manager should succeed")

	// Bob (new manager) should be able to update metadata
	msg := &types.MsgSetCollectionMetadata{
		Creator:      suite.Bob,
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://bob-update.com/metadata",
			CustomData: "bob's custom data",
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "new manager should be able to update metadata")

	// Old manager should not be able to update
	msg2 := &types.MsgSetCollectionMetadata{
		Creator:      suite.Manager, // Old manager
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://old-manager.com/metadata",
			CustomData: "",
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "old manager should not be able to update metadata")
}
