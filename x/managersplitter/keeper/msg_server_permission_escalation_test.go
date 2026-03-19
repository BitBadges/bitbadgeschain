package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestPermissionEscalation_ApprovalsCannotUpdatePermissions verifies that an executor
// with ONLY canUpdateCollectionApprovals CANNOT set UpdateCollectionPermissions = true.
// This is the core fix for the permission escalation vulnerability (backlog #27).
func (suite *TestSuite) TestPermissionEscalation_ApprovalsCannotUpdatePermissions() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a manager splitter where alice can update approvals but NOT permissions
	perms := GetDefaultPermissions()
	perms.CanUpdateCollectionApprovals.ApprovedAddresses = []string{alice}
	// canUpdateCollectionPermissions has NO approved addresses (empty list)

	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: perms,
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter")

	// Alice tries to update collection permissions via the manager splitter.
	// She has canUpdateCollectionApprovals but NOT canUpdateCollectionPermissions.
	// This MUST be denied.
	execMsg := &types.MsgExecuteUniversalUpdateCollection{
		Executor:                alice,
		ManagerSplitterAddress:  res.Address,
		UniversalUpdateCollectionMsg: &tokenizationtypes.MsgUniversalUpdateCollection{
			Creator:                    alice,
			CollectionId:               sdkmath.NewUint(1),
			UpdateCollectionPermissions: true,
			CollectionPermissions:      &tokenizationtypes.CollectionPermissions{},
		},
	}

	_, err = ExecuteUniversalUpdateCollection(suite, wctx, execMsg)
	suite.Require().Error(err, "Should deny: alice has canUpdateCollectionApprovals but not canUpdateCollectionPermissions")
	suite.Require().Contains(err.Error(), "canUpdateCollectionPermissions",
		"Error should reference the correct permission name")
}

// TestPermissionEscalation_CanUpdatePermissionsWithCorrectPerm verifies that an executor
// with canUpdateCollectionPermissions CAN set UpdateCollectionPermissions = true
// (though it may fail later at the tokenization module level, the permission check itself passes).
func (suite *TestSuite) TestPermissionEscalation_CanUpdatePermissionsWithCorrectPerm() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a manager splitter where alice can update collection permissions
	perms := GetDefaultPermissions()
	perms.CanUpdateCollectionPermissions.ApprovedAddresses = []string{alice}

	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: perms,
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter")

	// Alice tries to update collection permissions. She has the correct permission.
	// The permission check should pass (but it may fail later in tokenization module).
	execMsg := &types.MsgExecuteUniversalUpdateCollection{
		Executor:                alice,
		ManagerSplitterAddress:  res.Address,
		UniversalUpdateCollectionMsg: &tokenizationtypes.MsgUniversalUpdateCollection{
			Creator:                    alice,
			CollectionId:               sdkmath.NewUint(1),
			UpdateCollectionPermissions: true,
			CollectionPermissions:      &tokenizationtypes.CollectionPermissions{},
		},
	}

	_, err = ExecuteUniversalUpdateCollection(suite, wctx, execMsg)
	// If it errors, it should NOT be a permission denied error for canUpdateCollectionPermissions.
	// It may error for other reasons (e.g., collection not found in tokenization module).
	if err != nil {
		suite.Require().NotContains(err.Error(), "not approved for canUpdateCollectionPermissions",
			"Permission check should pass for executor with canUpdateCollectionPermissions")
	}
}

// TestPermissionEscalation_AdminBypassesCheck verifies that admin always bypasses
// all permission checks including canUpdateCollectionPermissions.
func (suite *TestSuite) TestPermissionEscalation_AdminBypassesCheck() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a manager splitter with no approved addresses for any permission
	perms := GetDefaultPermissions()
	// All permissions have empty approved addresses

	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: perms,
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter")

	// Admin (bob) tries to update collection permissions.
	// Admin should always bypass permission checks.
	execMsg := &types.MsgExecuteUniversalUpdateCollection{
		Executor:                bob,
		ManagerSplitterAddress:  res.Address,
		UniversalUpdateCollectionMsg: &tokenizationtypes.MsgUniversalUpdateCollection{
			Creator:                    bob,
			CollectionId:               sdkmath.NewUint(1),
			UpdateCollectionPermissions: true,
			CollectionPermissions:      &tokenizationtypes.CollectionPermissions{},
		},
	}

	_, err = ExecuteUniversalUpdateCollection(suite, wctx, execMsg)
	// If it errors, it should NOT be a permission denied error.
	if err != nil {
		suite.Require().NotContains(err.Error(), "permission denied",
			"Admin should bypass all permission checks")
	}
}

// TestPermissionEscalation_ApprovalsStillWorkForApprovals verifies that
// canUpdateCollectionApprovals still correctly controls UpdateCollectionApprovals
// (i.e., we didn't break existing functionality).
func (suite *TestSuite) TestPermissionEscalation_ApprovalsStillWorkForApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a manager splitter where charlie has NO approval permissions
	perms := GetDefaultPermissions()
	// All permissions have empty approved addresses

	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: perms,
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter")

	// Charlie tries to update collection approvals without permission
	execMsg := &types.MsgExecuteUniversalUpdateCollection{
		Executor:                charlie,
		ManagerSplitterAddress:  res.Address,
		UniversalUpdateCollectionMsg: &tokenizationtypes.MsgUniversalUpdateCollection{
			Creator:                  charlie,
			CollectionId:             sdkmath.NewUint(1),
			UpdateCollectionApprovals: true,
		},
	}

	_, err = ExecuteUniversalUpdateCollection(suite, wctx, execMsg)
	suite.Require().Error(err, "Should deny: charlie has no canUpdateCollectionApprovals permission")
	suite.Require().Contains(err.Error(), "canUpdateCollectionApprovals",
		"Error should reference canUpdateCollectionApprovals")
}
