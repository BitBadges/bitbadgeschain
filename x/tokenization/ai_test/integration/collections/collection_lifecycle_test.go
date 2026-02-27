package collections_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type CollectionLifecycleTestSuite struct {
	testutil.AITestSuite
}

func TestCollectionLifecycleSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(CollectionLifecycleTestSuite))
}

// TestCollectionLifecycle_CreateUpdateArchiveDelete tests the full Create -> Update -> Archive -> Delete workflow
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_CreateUpdateArchiveDelete() {
	// Step 1: Create a collection
	collectionId := suite.CreateTestCollection(suite.Manager)
	suite.Require().True(collectionId.GT(sdkmath.NewUint(0)), "collection should be created with valid ID")

	// Verify collection exists
	collection := suite.GetCollection(collectionId)
	suite.Require().Equal(suite.Manager, collection.Manager)
	suite.Require().False(collection.IsArchived, "collection should not be archived initially")

	// Step 2: Update the collection (add metadata)
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                    suite.Manager,
		CollectionId:               collectionId,
		UpdateCollectionMetadata:   true,
		CollectionMetadata:         testutil.GenerateCollectionMetadata("https://updated.example.com/metadata", "updated custom data"),
		UpdateCustomData:           true,
		CustomData:                 "new custom data",
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "updating collection should succeed")

	// Verify update
	updatedCollection := suite.GetCollection(collectionId)
	suite.Require().Equal("https://updated.example.com/metadata", updatedCollection.CollectionMetadata.Uri)
	suite.Require().Equal("new custom data", updatedCollection.CustomData)

	// Step 3: Archive the collection
	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}

	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err, "archiving collection should succeed")

	// Verify archive
	archivedCollection := suite.GetCollection(collectionId)
	suite.Require().True(archivedCollection.IsArchived, "collection should be archived")

	// Step 4: Unarchive the collection
	unarchiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       false,
	}

	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), unarchiveMsg)
	suite.Require().NoError(err, "unarchiving collection should succeed")

	// Verify unarchive
	unarchivedCollection := suite.GetCollection(collectionId)
	suite.Require().False(unarchivedCollection.IsArchived, "collection should be unarchived")

	// Step 5: Delete the collection
	deleteMsg := &types.MsgDeleteCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
	}

	_, err = suite.MsgServer.DeleteCollection(sdk.WrapSDKContext(suite.Ctx), deleteMsg)
	suite.Require().NoError(err, "deleting collection should succeed")

	// Verify deletion
	_, found := suite.Keeper.GetCollectionFromStore(suite.Ctx, collectionId)
	suite.Require().False(found, "collection should not exist after deletion")
}

// TestCollectionLifecycle_ArchiveCollectionAfterUpdates tests archiving a collection after making updates
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_ArchiveCollectionAfterUpdates() {
	// Create collection
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Setup mint approval and mint some tokens
	suite.SetupMintApproval(collectionId)
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintTokens(collectionId, suite.Alice, mintBalances)

	// Add collection approvals
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Update metadata
	metadataMsg := &types.MsgUniversalUpdateCollection{
		Creator:                  suite.Manager,
		CollectionId:             collectionId,
		UpdateCollectionMetadata: true,
		CollectionMetadata:       testutil.GenerateCollectionMetadata("https://final.example.com/metadata", ""),
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), metadataMsg)
	suite.Require().NoError(err)

	// Now archive the collection
	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}

	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err, "archiving collection after updates should succeed")

	// Verify collection is archived but data is preserved
	archivedCollection := suite.GetCollection(collectionId)
	suite.Require().True(archivedCollection.IsArchived)
	suite.Require().Equal("https://final.example.com/metadata", archivedCollection.CollectionMetadata.Uri)

	// Verify balances are still there
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	suite.Require().Greater(len(aliceBalance.Balances), 0, "balances should be preserved after archiving")
}

// TestCollectionLifecycle_ArchivedCollectionRejectsModifications tests that archived collections reject modifications
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_ArchivedCollectionRejectsModifications() {
	// Create and archive collection
	collectionId := suite.CreateTestCollection(suite.Manager)

	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err)

	// Try to update metadata - should fail
	metadataMsg := &types.MsgUniversalUpdateCollection{
		Creator:                  suite.Manager,
		CollectionId:             collectionId,
		UpdateCollectionMetadata: true,
		CollectionMetadata:       testutil.GenerateCollectionMetadata("https://new.example.com", ""),
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), metadataMsg)
	suite.Require().Error(err, "updating metadata on archived collection should fail")

	// Try to update custom data - should fail
	customDataMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateCustomData: true,
		CustomData:       "new data",
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), customDataMsg)
	suite.Require().Error(err, "updating custom data on archived collection should fail")

	// Try to update collection approvals - should fail
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	approvalsMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), approvalsMsg)
	suite.Require().Error(err, "updating approvals on archived collection should fail")

	// Try to update valid token IDs - should fail
	validTokensMsg := &types.MsgUniversalUpdateCollection{
		Creator:             suite.Manager,
		CollectionId:        collectionId,
		UpdateValidTokenIds: true,
		ValidTokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), validTokensMsg)
	suite.Require().Error(err, "updating valid token IDs on archived collection should fail")
}

// TestCollectionLifecycle_ArchivedCollectionRejectsTransfers tests that archived collections reject transfers
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_ArchivedCollectionRejectsTransfers() {
	// Create collection and setup for transfers
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Setup mint approval and mint tokens
	suite.SetupMintApproval(collectionId)
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintTokens(collectionId, suite.Alice, mintBalances)

	// Setup transfer approvals
	collectionApproval := testutil.GenerateCollectionApproval("transfer1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{collectionApproval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Setup user approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Bob,
		CollectionId: collectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().NoError(err)

	// Now archive the collection
	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err)

	// Try to transfer - should fail
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			testutil.GenerateTransfer(suite.Alice, []string{suite.Bob}, []*types.Balance{
				testutil.GenerateSimpleBalance(10, 1),
			}),
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	suite.Require().Error(err, "transfer on archived collection should fail")
}

// TestCollectionLifecycle_UnarchiveCollection tests unarchiving a collection
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_UnarchiveCollection() {
	// Create and archive collection
	collectionId := suite.CreateTestCollection(suite.Manager)

	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err)

	// Verify archived
	archivedCollection := suite.GetCollection(collectionId)
	suite.Require().True(archivedCollection.IsArchived)

	// Unarchive
	unarchiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       false,
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), unarchiveMsg)
	suite.Require().NoError(err, "unarchiving collection should succeed")

	// Verify unarchived
	unarchivedCollection := suite.GetCollection(collectionId)
	suite.Require().False(unarchivedCollection.IsArchived)

	// Verify modifications work again
	metadataMsg := &types.MsgUniversalUpdateCollection{
		Creator:                  suite.Manager,
		CollectionId:             collectionId,
		UpdateCollectionMetadata: true,
		CollectionMetadata:       testutil.GenerateCollectionMetadata("https://after-unarchive.example.com", ""),
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), metadataMsg)
	suite.Require().NoError(err, "updating metadata after unarchive should succeed")

	// Verify update worked
	finalCollection := suite.GetCollection(collectionId)
	suite.Require().Equal("https://after-unarchive.example.com", finalCollection.CollectionMetadata.Uri)
}

// TestCollectionLifecycle_DeleteActiveCollection tests deleting an active (non-archived) collection
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_DeleteActiveCollection() {
	// Create collection
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Verify collection exists and is not archived
	collection := suite.GetCollection(collectionId)
	suite.Require().False(collection.IsArchived)

	// Delete without archiving first
	deleteMsg := &types.MsgDeleteCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
	}

	_, err := suite.MsgServer.DeleteCollection(sdk.WrapSDKContext(suite.Ctx), deleteMsg)
	suite.Require().NoError(err, "deleting active collection should succeed")

	// Verify deletion
	_, found := suite.Keeper.GetCollectionFromStore(suite.Ctx, collectionId)
	suite.Require().False(found, "collection should not exist after deletion")
}

// TestCollectionLifecycle_DeleteWithOutstandingBalances tests deleting a collection that has outstanding balances
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_DeleteWithOutstandingBalances() {
	// Create collection
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Setup mint approval and mint tokens to multiple users
	suite.SetupMintApproval(collectionId)

	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
		testutil.GenerateSimpleBalance(50, 2),
	}
	suite.MintTokens(collectionId, suite.Alice, mintBalances)
	suite.MintTokens(collectionId, suite.Bob, mintBalances)
	suite.MintTokens(collectionId, suite.Charlie, mintBalances)

	// Verify balances exist
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	suite.Require().Greater(len(aliceBalance.Balances), 0, "Alice should have balances")
	bobBalance := suite.GetBalance(collectionId, suite.Bob)
	suite.Require().Greater(len(bobBalance.Balances), 0, "Bob should have balances")

	// Delete collection with outstanding balances
	deleteMsg := &types.MsgDeleteCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
	}

	_, err := suite.MsgServer.DeleteCollection(sdk.WrapSDKContext(suite.Ctx), deleteMsg)
	suite.Require().NoError(err, "deleting collection with balances should succeed")

	// Verify collection is gone
	_, found := suite.Keeper.GetCollectionFromStore(suite.Ctx, collectionId)
	suite.Require().False(found, "collection should not exist after deletion")

	// Verify balances are purged
	allBalances, _, allCollectionIds, err := suite.Keeper.GetUserBalancesFromStore(suite.Ctx)
	suite.Require().NoError(err)
	for i, id := range allCollectionIds {
		suite.Require().False(id.Equal(collectionId), "No balances should remain for deleted collection")
		_ = allBalances[i]
	}
}

// TestCollectionLifecycle_NonManagerCannotArchive tests that non-managers cannot archive
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_NonManagerCannotArchive() {
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Alice (non-manager) tries to archive
	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Alice,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().Error(err, "non-manager should not be able to archive collection")
}

// TestCollectionLifecycle_NonManagerCannotDelete tests that non-managers cannot delete
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_NonManagerCannotDelete() {
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Alice (non-manager) tries to delete
	deleteMsg := &types.MsgDeleteCollection{
		Creator:      suite.Alice,
		CollectionId: collectionId,
	}

	_, err := suite.MsgServer.DeleteCollection(sdk.WrapSDKContext(suite.Ctx), deleteMsg)
	suite.Require().Error(err, "non-manager should not be able to delete collection")
}

// TestCollectionLifecycle_ArchivePermissionLocked tests that locked archive permissions prevent archiving
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_ArchivePermissionLocked() {
	// Create collection with locked archive permissions
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Lock the archive permission (forbid archiving forever)
	lockPermissionsMsg := &types.MsgUniversalUpdateCollection{
		Creator:                     suite.Manager,
		CollectionId:                collectionId,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanArchiveCollection: []*types.ActionPermission{
				{
					PermanentlyPermittedTimes: []*types.UintRange{},
					PermanentlyForbiddenTimes: []*types.UintRange{
						{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
					},
				},
			},
		},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), lockPermissionsMsg)
	suite.Require().NoError(err)

	// Now try to archive - should fail
	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().Error(err, "archiving should fail when permission is locked")
}

// TestCollectionLifecycle_DeletePermissionLocked tests that locked delete permissions prevent deletion
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_DeletePermissionLocked() {
	// Create collection with locked delete permissions
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Lock the delete permission (forbid deleting forever)
	lockPermissionsMsg := &types.MsgUniversalUpdateCollection{
		Creator:                     suite.Manager,
		CollectionId:                collectionId,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanDeleteCollection: []*types.ActionPermission{
				{
					PermanentlyPermittedTimes: []*types.UintRange{},
					PermanentlyForbiddenTimes: []*types.UintRange{
						{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
					},
				},
			},
		},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), lockPermissionsMsg)
	suite.Require().NoError(err)

	// Now try to delete - should fail
	deleteMsg := &types.MsgDeleteCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
	}
	_, err = suite.MsgServer.DeleteCollection(sdk.WrapSDKContext(suite.Ctx), deleteMsg)
	suite.Require().Error(err, "deleting should fail when permission is locked")
}

// TestCollectionLifecycle_MultipleArchiveUnarchiveCycles tests multiple archive/unarchive cycles
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_MultipleArchiveUnarchiveCycles() {
	collectionId := suite.CreateTestCollection(suite.Manager)

	for i := 0; i < 3; i++ {
		// Archive
		archiveMsg := &types.MsgUniversalUpdateCollection{
			Creator:          suite.Manager,
			CollectionId:     collectionId,
			UpdateIsArchived: true,
			IsArchived:       true,
		}
		_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
		suite.Require().NoError(err, "archive cycle %d should succeed", i+1)

		collection := suite.GetCollection(collectionId)
		suite.Require().True(collection.IsArchived, "should be archived in cycle %d", i+1)

		// Unarchive
		unarchiveMsg := &types.MsgUniversalUpdateCollection{
			Creator:          suite.Manager,
			CollectionId:     collectionId,
			UpdateIsArchived: true,
			IsArchived:       false,
		}
		_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), unarchiveMsg)
		suite.Require().NoError(err, "unarchive cycle %d should succeed", i+1)

		collection = suite.GetCollection(collectionId)
		suite.Require().False(collection.IsArchived, "should be unarchived in cycle %d", i+1)
	}
}

// TestCollectionLifecycle_ArchivePreservesAllData tests that archiving preserves all collection data
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_ArchivePreservesAllData() {
	// Create collection with various data
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Setup mint approval and mint tokens
	suite.SetupMintApproval(collectionId)
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(100, 1),
	}
	suite.MintTokens(collectionId, suite.Alice, mintBalances)

	// Add various data
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionMetadata:  true,
		CollectionMetadata:        testutil.GenerateCollectionMetadata("https://test.example.com", "custom data"),
		UpdateCustomData:          true,
		CustomData:                "test custom data",
		UpdateStandards:           true,
		Standards:                 []string{"standard1", "standard2"},
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Archive
	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err)

	// Verify all data is preserved
	archivedCollection := suite.GetCollection(collectionId)
	suite.Require().True(archivedCollection.IsArchived)
	suite.Require().Equal("https://test.example.com", archivedCollection.CollectionMetadata.Uri)
	suite.Require().Equal("test custom data", archivedCollection.CustomData)
	suite.Require().Equal(2, len(archivedCollection.Standards))
	suite.Require().Greater(len(archivedCollection.CollectionApprovals), 0)

	// Verify balances preserved
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	suite.Require().Greater(len(aliceBalance.Balances), 0)
}

// TestCollectionLifecycle_ArchivedCollectionRejectsUserApprovalUpdates tests that user approval updates fail on archived collections
func (suite *CollectionLifecycleTestSuite) TestCollectionLifecycle_ArchivedCollectionRejectsUserApprovalUpdates() {
	collectionId := suite.CreateTestCollection(suite.Manager)

	// Archive
	archiveMsg := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), archiveMsg)
	suite.Require().NoError(err)

	// Try to set outgoing approval
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().Error(err, "setting outgoing approval on archived collection should fail")

	// Try to set incoming approval
	incomingApproval := testutil.GenerateUserIncomingApproval("incoming1", "All")
	setIncomingMsg := &types.MsgSetIncomingApproval{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Approval:     incomingApproval,
	}
	_, err = suite.MsgServer.SetIncomingApproval(sdk.WrapSDKContext(suite.Ctx), setIncomingMsg)
	suite.Require().Error(err, "setting incoming approval on archived collection should fail")

	// Try to update user approvals
	updateUserMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            collectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{outgoingApproval},
	}
	_, err = suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), updateUserMsg)
	suite.Require().Error(err, "updating user approvals on archived collection should fail")
}
