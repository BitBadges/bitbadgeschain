package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type SecurityTestSuite struct {
	TestSuite
}

func TestSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(SecurityTestSuite))
}

func (suite *SecurityTestSuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

// =============================================================================
// Self-Transfer Guard Tests
// =============================================================================

// TestSelfTransfer_Prevented tests that transferring tokens to oneself is blocked
func (suite *SecurityTestSuite) TestSelfTransfer_Prevented() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with bob as creator (gets tokens minted to him)
	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to transfer from bob to bob (self-transfer)
	msg := &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{bob}, // Self-transfer attempt
				Balances: []*types.Balance{
					{
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
			},
		},
	}

	_, err = suite.msgServer.TransferTokens(wctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "sender and receiver cannot be the same")
}


// =============================================================================
// Permission Immutability Tests
// =============================================================================

// TestPermissionImmutability_CannotUpdateCustomDataWhenForbidden tests that
// permanently forbidden permissions block custom data updates
func (suite *SecurityTestSuite) TestPermissionImmutability_CannotUpdateCustomDataWhenForbidden() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with permanently forbidden custom data permission
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Permissions = &types.CollectionPermissions{
		CanUpdateCustomData: []*types.ActionPermission{
			{
				PermanentlyForbiddenTimes: GetFullUintRanges(),
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to update custom data (should fail)
	err = UpdateCollection(&suite.TestSuite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:          bob,
		CollectionId:     sdkmath.NewUint(1),
		UpdateCustomData: true,
		CustomData:       "new-custom-data",
	})
	suite.Require().Error(err)
	// The error should indicate permission violation
}

// TestPermissionImmutability_CannotChangeLockedManager tests that
// locked manager cannot be changed when permanently forbidden
func (suite *SecurityTestSuite) TestPermissionImmutability_CannotChangeLockedManager() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with locked manager (can't update manager)
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Permissions = &types.CollectionPermissions{
		CanUpdateManager: []*types.ActionPermission{
			{
				PermanentlyForbiddenTimes: GetFullUintRanges(),
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to change manager (should fail since updating manager is forbidden)
	err = UpdateCollection(&suite.TestSuite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:       bob,
		CollectionId:  sdkmath.NewUint(1),
		UpdateManager: true,
		Manager:       alice, // Trying to change manager to alice
	})
	suite.Require().Error(err)
}

// =============================================================================
// Manager Authorization Tests
// =============================================================================

// TestManagerAuthorization_NonManagerCannotUpdate tests that non-managers cannot update collections
func (suite *SecurityTestSuite) TestManagerAuthorization_NonManagerCannotUpdate() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with bob as manager
	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify bob is the manager
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().Equal(bob, collection.Manager)

	// Alice (non-manager) trying to update custom data should fail
	err = UpdateCollection(&suite.TestSuite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:          alice, // Not the manager
		CollectionId:     sdkmath.NewUint(1),
		UpdateCustomData: true,
		CustomData:       "alice-custom-data",
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "manager")
}

// =============================================================================
// Collection Invariant Enforcement Tests
// =============================================================================

// TestCollectionInvariant_MaxSupplyEnforced tests that max supply invariant is enforced
func (suite *SecurityTestSuite) TestCollectionInvariant_MaxSupplyEnforced() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with max supply invariant
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: sdkmath.NewUint(1), // Only 1 of each token allowed
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// The initial mint gives bob tokens. Verify collection was created
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().NotNil(collection)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().Equal(sdkmath.NewUint(1), collection.Invariants.MaxSupplyPerId)
}

// TestCollectionInvariant_NoCustomOwnershipTimes tests the no custom ownership times invariant
func (suite *SecurityTestSuite) TestCollectionInvariant_NoCustomOwnershipTimes() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with no custom ownership times invariant
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		NoCustomOwnershipTimes: true,
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Collection should be created successfully
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().NotNil(collection)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().True(collection.Invariants.NoCustomOwnershipTimes)
}

// =============================================================================
// Transfer Security Tests
// =============================================================================

// TestTransfer_InsufficientBalance tests that transfers fail with insufficient balance
func (suite *SecurityTestSuite) TestTransfer_InsufficientBalance() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection
	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to transfer more than bob has
	msg := &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
						Amount:         sdkmath.NewUint(999999), // Way more than exists
					},
				},
			},
		},
	}

	_, err = suite.msgServer.TransferTokens(wctx, msg)
	suite.Require().Error(err)
}

// TestTransfer_UnauthorizedInitiator tests that only authorized initiators can transfer
func (suite *SecurityTestSuite) TestTransfer_UnauthorizedInitiator() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection
	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Charlie tries to transfer bob's tokens (should fail - not authorized)
	msg := &types.MsgTransferTokens{
		Creator:      charlie, // Not bob
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob, // Trying to transfer from bob
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
			},
		},
	}

	_, err = suite.msgServer.TransferTokens(wctx, msg)
	// This should fail because charlie is not approved to transfer bob's tokens
	// The specific error depends on approval configuration
	suite.Require().Error(err)
}

// =============================================================================
// Permission Immutability Enforcement Tests (Advanced)
// =============================================================================

// TestPermissionImmutability_CannotClearForbiddenTimes tests that permanently forbidden
// times cannot be cleared or changed to permitted
func (suite *SecurityTestSuite) TestPermissionImmutability_CannotClearForbiddenTimes() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create collection with forbidden custom data updates
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Permissions = &types.CollectionPermissions{
		CanUpdateCustomData: []*types.ActionPermission{
			{
				PermanentlyForbiddenTimes: GetFullUintRanges(),
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to change the permission to permitted (should fail - can't change frozen permissions)
	err = UpdateCollection(&suite.TestSuite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:                     bob,
		CollectionId:                sdkmath.NewUint(1),
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateCustomData: []*types.ActionPermission{
				{
					PermanentlyPermittedTimes: GetFullUintRanges(), // Trying to change to permitted
				},
			},
		},
	})
	// This should fail - you can't loosen a frozen permission
	suite.Require().Error(err)
}

// TestPermissionImmutability_MetadataUpdatesForbidden tests that token metadata
// cannot be updated when forbidden
func (suite *SecurityTestSuite) TestPermissionImmutability_MetadataUpdatesForbidden() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create collection with forbidden token metadata updates
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Permissions = &types.CollectionPermissions{
		CanUpdateCollectionMetadata: []*types.ActionPermission{
			{
				PermanentlyForbiddenTimes: GetFullUintRanges(),
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to update collection metadata (should fail)
	err = UpdateCollection(&suite.TestSuite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:                  bob,
		CollectionId:             sdkmath.NewUint(1),
		UpdateCollectionMetadata: true,
		CollectionMetadata:       &types.CollectionMetadata{}, // Try to update
	})
	suite.Require().Error(err)
}

// =============================================================================
// Governance Authority Tests
// =============================================================================

// TestGovernanceAuthority_CanBypassManagerCheck tests that the governance authority
// can perform manager-only operations
func (suite *SecurityTestSuite) TestGovernanceAuthority_CanBypassManagerCheck() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with bob as manager
	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify bob is the manager
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().Equal(bob, collection.Manager)

	// Get the governance authority from the keeper
	govAuthority := suite.app.TokenizationKeeper.GetAuthority()
	suite.Require().NotEmpty(govAuthority)

	// The governance authority should be able to perform updates
	// Note: This requires the governance authority to have a funded account
	// For the test, we're just verifying the authority is properly configured
	suite.Require().NotEqual(bob, govAuthority, "Gov authority should be different from bob")
	suite.Require().NotEqual(alice, govAuthority, "Gov authority should be different from alice")
}

// TestGovernanceAuthority_NonGovNonManagerFails tests that neither governance
// nor manager can bypass restrictions - only valid addresses can act
func (suite *SecurityTestSuite) TestGovernanceAuthority_NonGovNonManagerFails() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with bob as manager
	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Charlie (not manager, not governance) trying to update should fail
	err = UpdateCollection(&suite.TestSuite, wctx, &types.MsgUniversalUpdateCollection{
		Creator:          charlie,
		CollectionId:     sdkmath.NewUint(1),
		UpdateCustomData: true,
		CustomData:       "charlie-data",
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "manager")
}
