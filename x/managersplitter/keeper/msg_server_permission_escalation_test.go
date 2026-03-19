package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestPermissionEscalation_NonAdminCannotUpdatePermissions verifies that a non-admin
// executor CANNOT set UpdateCollectionPermissions = true, even if they have
// canUpdateCollectionApprovals. Only admin can update collection permissions.
func (suite *TestSuite) TestPermissionEscalation_NonAdminCannotUpdatePermissions() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a manager splitter where alice can update approvals
	perms := GetDefaultPermissions()
	perms.CanUpdateCollectionApprovals.ApprovedAddresses = []string{alice}

	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: perms,
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter")

	// Alice tries to update collection permissions.
	// Even though she has canUpdateCollectionApprovals, only admin can update permissions.
	execMsg := &types.MsgExecuteUniversalUpdateCollection{
		Executor:               alice,
		ManagerSplitterAddress: res.Address,
		UniversalUpdateCollectionMsg: &tokenizationtypes.MsgUniversalUpdateCollection{
			Creator:                    alice,
			CollectionId:               sdkmath.NewUint(1),
			UpdateCollectionPermissions: true,
			CollectionPermissions:      &tokenizationtypes.CollectionPermissions{},
		},
	}

	_, err = ExecuteUniversalUpdateCollection(suite, wctx, execMsg)
	suite.Require().Error(err, "Should deny: only admin can update collection permissions")
	suite.Require().Contains(err.Error(), "only admin",
		"Error should indicate admin-only restriction")
}

// TestPermissionEscalation_AdminCanUpdatePermissions verifies that admin can
// update collection permissions (bypasses all checks).
func (suite *TestSuite) TestPermissionEscalation_AdminCanUpdatePermissions() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	perms := GetDefaultPermissions()

	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: perms,
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter")

	// Admin (bob) tries to update collection permissions. Should pass permission check.
	execMsg := &types.MsgExecuteUniversalUpdateCollection{
		Executor:               bob,
		ManagerSplitterAddress: res.Address,
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
		suite.Require().NotContains(err.Error(), "only admin",
			"Admin should bypass permission check")
	}
}

// TestPermissionEscalation_ApprovalsStillWorkForApprovals verifies that
// canUpdateCollectionApprovals still correctly controls UpdateCollectionApprovals.
func (suite *TestSuite) TestPermissionEscalation_ApprovalsStillWorkForApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	perms := GetDefaultPermissions()

	createMsg := &types.MsgCreateManagerSplitter{
		Admin:       bob,
		Permissions: perms,
	}

	res, err := CreateManagerSplitter(suite, wctx, createMsg)
	suite.Require().Nil(err, "Error creating manager splitter")

	// Charlie tries to update collection approvals without permission
	execMsg := &types.MsgExecuteUniversalUpdateCollection{
		Executor:               charlie,
		ManagerSplitterAddress: res.Address,
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
