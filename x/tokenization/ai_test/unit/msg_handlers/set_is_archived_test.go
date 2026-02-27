package msg_handlers_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type SetIsArchivedTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestSetIsArchivedSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(SetIsArchivedTestSuite))
}

func (suite *SetIsArchivedTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
	suite.Require().True(suite.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0 after creation")
}

// TestSetIsArchived_ArchiveCollection tests successfully archiving a collection
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_ArchiveCollection() {
	// Verify collection is not archived initially
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().False(collection.IsArchived, "collection should not be archived initially")

	msg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "archiving collection should succeed")

	// Verify collection is archived
	collection = suite.GetCollection(suite.CollectionId)
	suite.Require().True(collection.IsArchived, "collection should be archived")
}

// TestSetIsArchived_UnarchiveCollection tests successfully unarchiving a collection
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_UnarchiveCollection() {
	// First archive the collection
	archiveMsg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err, "archiving collection should succeed")

	// Verify archived
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().True(collection.IsArchived, "collection should be archived")

	// Now unarchive
	unarchiveMsg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           false,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), unarchiveMsg)
	suite.Require().NoError(err, "unarchiving collection should succeed")

	// Verify unarchived
	collection = suite.GetCollection(suite.CollectionId)
	suite.Require().False(collection.IsArchived, "collection should be unarchived")
}

// TestSetIsArchived_PermissionChecked tests that canArchiveCollection permission is enforced
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_PermissionChecked() {
	// First, forbid all future archive changes
	forbidAllTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	msg := &types.MsgSetIsArchived{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		IsArchived:   true,
		CanArchiveCollection: []*types.ActionPermission{
			{
				PermanentlyForbiddenTimes: forbidAllTimes,
				PermanentlyPermittedTimes: []*types.UintRange{},
			},
		},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first archive should succeed")

	// Now try to unarchive - should fail because permission forbids it
	msg2 := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           false,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "unarchiving should fail because permission is forbidden")
}

// TestSetIsArchived_ArchivedCollectionRejectsModifications tests that archived collection rejects modifications
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_ArchivedCollectionRejectsModifications() {
	// First archive the collection
	archiveMsg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err, "archiving collection should succeed")

	// Try to update manager - should fail because collection is archived
	setManagerMsg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), setManagerMsg)
	suite.Require().Error(err, "updating manager on archived collection should fail")

	// Try to update collection metadata - should fail
	setMetadataMsg := &types.MsgSetCollectionMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://new.example.com",
			CustomData: "",
		},
		CanUpdateCollectionMetadata: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCollectionMetadata(sdk.WrapSDKContext(suite.Ctx), setMetadataMsg)
	suite.Require().Error(err, "updating metadata on archived collection should fail")

	// Try to update custom data - should fail
	setCustomDataMsg := &types.MsgSetCustomData{
		Creator:             suite.Manager,
		CollectionId:        suite.CollectionId,
		CustomData:          "new data",
		CanUpdateCustomData: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetCustomData(sdk.WrapSDKContext(suite.Ctx), setCustomDataMsg)
	suite.Require().Error(err, "updating custom data on archived collection should fail")

	// Try to update standards - should fail
	setStandardsMsg := &types.MsgSetStandards{
		Creator:            suite.Manager,
		CollectionId:       suite.CollectionId,
		Standards:          []string{"ERC-721"},
		CanUpdateStandards: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetStandards(sdk.WrapSDKContext(suite.Ctx), setStandardsMsg)
	suite.Require().Error(err, "updating standards on archived collection should fail")
}

// TestSetIsArchived_OnlyManagerCanArchive tests that only manager can archive/unarchive
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_OnlyManagerCanArchive() {
	msg := &types.MsgSetIsArchived{
		Creator:              suite.Alice, // Not the manager
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-manager should not be able to archive collection")
	suite.Require().Contains(err.Error(), "manager", "error should mention manager permission")
}

// TestSetIsArchived_NonExistentCollection tests archiving non-existent collection
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_NonExistentCollection() {
	nonExistentId := sdkmath.NewUint(99999)

	msg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         nonExistentId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "archiving non-existent collection should fail")
	suite.Require().Contains(err.Error(), "collection", "error should mention collection")
}

// TestSetIsArchived_EmptyCreator tests behavior with empty creator
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_EmptyCreator() {
	msg := &types.MsgSetIsArchived{
		Creator:              "", // Empty creator
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "empty creator should fail")
}

// TestSetIsArchived_CollectionStillAccessibleWhenArchived tests that archived collection is still readable
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_CollectionStillAccessibleWhenArchived() {
	// Archive the collection
	msg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "archiving collection should succeed")

	// Collection should still be accessible (readable)
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().NotNil(collection, "archived collection should still be accessible")
	suite.Require().True(collection.IsArchived, "collection should be marked as archived")
	suite.Require().Equal(suite.Manager, collection.Manager, "collection data should be intact")
}

// TestSetIsArchived_UnarchiveThenModify tests that unarchived collection can be modified again
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_UnarchiveThenModify() {
	// Archive the collection
	archiveMsg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err, "archiving collection should succeed")

	// Unarchive the collection
	unarchiveMsg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           false,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), unarchiveMsg)
	suite.Require().NoError(err, "unarchiving collection should succeed")

	// Now modifications should work again
	setManagerMsg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), setManagerMsg)
	suite.Require().NoError(err, "updating manager on unarchived collection should succeed")
}

// TestSetIsArchived_CanArchiveWhileArchived tests archiving an already archived collection
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_CanArchiveWhileArchived() {
	// Archive the collection
	archiveMsg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err, "archiving collection should succeed")

	// Try to archive again (no-op)
	archiveMsg2 := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           true, // Already archived
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), archiveMsg2)
	// This should succeed (it's a no-op) or fail depending on implementation
	// The important thing is it doesn't change the archived state
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().True(collection.IsArchived, "collection should still be archived")
}

// TestSetIsArchived_AfterManagerChange tests that new manager can archive/unarchive
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_AfterManagerChange() {
	// First change manager to Bob
	setManagerMsg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}
	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), setManagerMsg)
	suite.Require().NoError(err, "setting new manager should succeed")

	// Bob (new manager) should be able to archive
	msg := &types.MsgSetIsArchived{
		Creator:              suite.Bob,
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "new manager should be able to archive collection")

	// Old manager should not be able to archive
	msg2 := &types.MsgSetIsArchived{
		Creator:              suite.Manager, // Old manager
		CollectionId:         suite.CollectionId,
		IsArchived:           false,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "old manager should not be able to unarchive collection")
}

// TestSetIsArchived_ArchivedCollectionCanStillBeUnarchived tests that archived collection can be unarchived
// even while other modifications are blocked
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_ArchivedCollectionCanStillBeUnarchived() {
	// Archive the collection
	archiveMsg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           true,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err, "archiving collection should succeed")

	// Verify other modifications are blocked
	setManagerMsg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}
	_, err = suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), setManagerMsg)
	suite.Require().Error(err, "modifications should be blocked while archived")

	// But unarchiving should still work
	unarchiveMsg := &types.MsgSetIsArchived{
		Creator:              suite.Manager,
		CollectionId:         suite.CollectionId,
		IsArchived:           false,
		CanArchiveCollection: []*types.ActionPermission{},
	}

	_, err = suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), unarchiveMsg)
	suite.Require().NoError(err, "unarchiving should still work on archived collection")
}

// TestSetIsArchived_MultipleArchiveUnarchiveCycles tests multiple archive/unarchive cycles
func (suite *SetIsArchivedTestSuite) TestSetIsArchived_MultipleArchiveUnarchiveCycles() {
	for i := 0; i < 3; i++ {
		// Archive
		archiveMsg := &types.MsgSetIsArchived{
			Creator:              suite.Manager,
			CollectionId:         suite.CollectionId,
			IsArchived:           true,
			CanArchiveCollection: []*types.ActionPermission{},
		}
		_, err := suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
		suite.Require().NoError(err, "archiving should succeed in cycle %d", i)

		// Verify archived
		collection := suite.GetCollection(suite.CollectionId)
		suite.Require().True(collection.IsArchived, "should be archived in cycle %d", i)

		// Unarchive
		unarchiveMsg := &types.MsgSetIsArchived{
			Creator:              suite.Manager,
			CollectionId:         suite.CollectionId,
			IsArchived:           false,
			CanArchiveCollection: []*types.ActionPermission{},
		}
		_, err = suite.MsgServer.SetIsArchived(sdk.WrapSDKContext(suite.Ctx), unarchiveMsg)
		suite.Require().NoError(err, "unarchiving should succeed in cycle %d", i)

		// Verify unarchived
		collection = suite.GetCollection(suite.CollectionId)
		suite.Require().False(collection.IsArchived, "should be unarchived in cycle %d", i)
	}
}
