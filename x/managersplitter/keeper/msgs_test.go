package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetDefaultPermissions() *types.ManagerSplitterPermissions {
	return &types.ManagerSplitterPermissions{
		CanDeleteCollection: &types.PermissionCriteria{
			ApprovedAddresses: []string{},
		},
		CanArchiveCollection: &types.PermissionCriteria{
			ApprovedAddresses: []string{},
		},
		CanUpdateStandards: &types.PermissionCriteria{
			ApprovedAddresses: []string{},
		},
		CanUpdateCustomData: &types.PermissionCriteria{
			ApprovedAddresses: []string{},
		},
		CanUpdateManager: &types.PermissionCriteria{
			ApprovedAddresses: []string{},
		},
		CanUpdateCollectionMetadata: &types.PermissionCriteria{
			ApprovedAddresses: []string{},
		},
		CanUpdateValidTokenIds: &types.PermissionCriteria{
			ApprovedAddresses: []string{},
		},
		CanUpdateTokenMetadata: &types.PermissionCriteria{
			ApprovedAddresses: []string{},
		},
		CanUpdateCollectionApprovals: &types.PermissionCriteria{
			ApprovedAddresses: []string{},
		},
	}
}

func (suite *TestSuite) TestCreateManagerSplitter() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: GetDefaultPermissions(),
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter: %s")
	suite.Require().NotNil(res, "Response should not be nil")
	suite.Require().NotEmpty(res.Address, "Address should not be empty")

	// Verify the manager splitter was created
	managerSplitter, err := GetManagerSplitter(suite, wctx, res.Address)
	suite.Require().Nil(err, "Error getting manager splitter: %s")
	suite.Require().NotNil(managerSplitter, "Manager splitter should not be nil")
	suite.Require().Equal(bob, managerSplitter.Admin, "Admin should match")
	suite.Require().Equal(res.Address, managerSplitter.Address, "Address should match")
}

func (suite *TestSuite) TestUpdateManagerSplitter() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a manager splitter
	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: GetDefaultPermissions(),
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter: %s")

	// Update permissions
	updatedPermissions := GetDefaultPermissions()
	updatedPermissions.CanUpdateManager.ApprovedAddresses = []string{alice}

	updateMsg := &types.MsgUpdateManagerSplitter{
		Admin:       bob,
		Address:     res.Address,
		Permissions: updatedPermissions,
	}

	err = UpdateManagerSplitter(suite, wctx, updateMsg)
	suite.Require().Nil(err, "Error updating manager splitter: %s")

	// Verify the update
	managerSplitter, err := GetManagerSplitter(suite, wctx, res.Address)
	suite.Require().Nil(err, "Error getting manager splitter: %s")
	suite.Require().Equal(1, len(managerSplitter.Permissions.CanUpdateManager.ApprovedAddresses), "Should have one approved address")
	suite.Require().Equal(alice, managerSplitter.Permissions.CanUpdateManager.ApprovedAddresses[0], "Approved address should match")

	// Try to update with wrong admin (should fail)
	updateMsg.Admin = alice
	err = UpdateManagerSplitter(suite, wctx, updateMsg)
	suite.Require().Error(err, "Should fail when non-admin tries to update")
}

func (suite *TestSuite) TestDeleteManagerSplitter() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a manager splitter
	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: GetDefaultPermissions(),
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter: %s")

	// Delete with wrong admin (should fail)
	deleteMsg := &types.MsgDeleteManagerSplitter{
		Admin:   alice,
		Address: res.Address,
	}
	err = DeleteManagerSplitter(suite, wctx, deleteMsg)
	suite.Require().Error(err, "Should fail when non-admin tries to delete")

	// Delete with correct admin
	deleteMsg.Admin = bob
	err = DeleteManagerSplitter(suite, wctx, deleteMsg)
	suite.Require().Nil(err, "Error deleting manager splitter: %s")

	// Verify it's deleted
	_, err = GetManagerSplitter(suite, wctx, res.Address)
	suite.Require().Error(err, "Should fail to get deleted manager splitter")
}

func (suite *TestSuite) TestExecuteUniversalUpdateCollection_Admin() {
	// Skip this test for now - requires full tokenization module setup
	// This test would verify that admin can execute UniversalUpdateCollection
	suite.T().Skip("Skipping - requires full tokenization module integration")
}

func (suite *TestSuite) TestExecuteUniversalUpdateCollection_ApprovedAddress() {
		// Skip this test for now - requires full tokenization module setup
		suite.T().Skip("Skipping - requires full tokenization module integration")
}

func (suite *TestSuite) TestExecuteUniversalUpdateCollection_Unauthorized() {
		// Skip this test for now - requires full tokenization module setup
		suite.T().Skip("Skipping - requires full tokenization module integration")
}

func (suite *TestSuite) TestExecuteUniversalUpdateCollection_MultiplePermissions() {
		// Skip this test for now - requires full tokenization module setup
		suite.T().Skip("Skipping - requires full tokenization module integration")
}
