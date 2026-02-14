package edge_cases_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type DefaultBalanceVersionIncrementTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestDefaultBalanceVersionIncrementSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(DefaultBalanceVersionIncrementTestSuite))
}

func (suite *DefaultBalanceVersionIncrementTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestDefaultBalanceVersionIncrement_FirstAccessIncrementsVersions verifies that
// the first access to a user's balance initializes approval version trackers.
// This test addresses HIGH-006: Default Balance Version Increment on First Access.
func (suite *DefaultBalanceVersionIncrementTestSuite) TestDefaultBalanceVersionIncrement_FirstAccessIncrementsVersions() {
	// Create a new collection with default approvals
	// Note: Cannot use reserved IDs like "default-incoming" or "default-outgoing"
	defaultIncomingApproval := testutil.GenerateUserIncomingApproval("test-incoming-1", "All")
	defaultOutgoingApproval := testutil.GenerateUserOutgoingApproval("test-outgoing-1", "All")

	createMsg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			IncomingApprovals: []*types.UserIncomingApproval{defaultIncomingApproval},
			OutgoingApprovals: []*types.UserOutgoingApproval{defaultOutgoingApproval},
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.Manager,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "",
		},
		TokenMetadata:       []*types.TokenMetadata{},
		CustomData:          "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards:           []string{},
		IsArchived:          false,
	}
	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), createMsg)
	suite.Require().NoError(err)
	collectionId := resp.CollectionId
	collection := suite.GetCollection(collectionId)

	// First access to balance - should initialize version trackers
	// Note: GetBalanceOrApplyDefault doesn't save the balance, but IncrementApprovalVersion
	// does save version trackers to the store
	balance1, appliedDefault1, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice)
	suite.Require().True(appliedDefault1, "default should be applied on first access")
	suite.Require().Equal(1, len(balance1.IncomingApprovals), "should have default incoming approval")
	suite.Require().Equal(1, len(balance1.OutgoingApprovals), "should have default outgoing approval")

	// Verify versions were initialized (should be 0 on first call to IncrementApprovalVersion)
	version1Incoming := balance1.IncomingApprovals[0].Version
	version1Outgoing := balance1.OutgoingApprovals[0].Version
	suite.Require().True(version1Incoming.IsZero(), "incoming approval version should be 0 on first access")
	suite.Require().True(version1Outgoing.IsZero(), "outgoing approval version should be 0 on first access")

	// Second access - version trackers should exist in store, so versions should increment
	// (GetBalanceOrApplyDefault doesn't save the balance, so it will apply defaults again,
	// but IncrementApprovalVersion will return the incremented version from the store)
	balance2, appliedDefault2, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice)
	suite.Require().True(appliedDefault2, "default should be applied again since balance wasn't saved")

	// Versions should have incremented because version trackers exist in store
	version2Incoming := balance2.IncomingApprovals[0].Version
	version2Outgoing := balance2.OutgoingApprovals[0].Version
	suite.Require().True(version2Incoming.GT(version1Incoming), "incoming approval version should increment on second access")
	suite.Require().True(version2Outgoing.GT(version1Outgoing), "outgoing approval version should increment on second access")
}

// TestDefaultBalanceVersionIncrement_VersionTrackerInitialized verifies that
// version trackers are properly initialized in the store on first access.
// This test verifies that the side effect (version initialization) occurs as documented.
func (suite *DefaultBalanceVersionIncrementTestSuite) TestDefaultBalanceVersionIncrement_VersionTrackerInitialized() {
	// Create a new collection with default approvals
	defaultIncomingApproval := testutil.GenerateUserIncomingApproval("test-incoming-2", "All")
	defaultOutgoingApproval := testutil.GenerateUserOutgoingApproval("test-outgoing-2", "All")

	createMsg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			IncomingApprovals: []*types.UserIncomingApproval{defaultIncomingApproval},
			OutgoingApprovals: []*types.UserOutgoingApproval{defaultOutgoingApproval},
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.Manager,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "",
		},
		TokenMetadata:       []*types.TokenMetadata{},
		CustomData:          "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards:           []string{},
		IsArchived:          false,
	}
	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), createMsg)
	suite.Require().NoError(err)
	collectionId := resp.CollectionId
	collection := suite.GetCollection(collectionId)

	// First access initializes version trackers in the store
	// This is the side effect documented in HIGH-006
	balance1, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice)
	suite.Require().Equal(1, len(balance1.IncomingApprovals), "should have default incoming approval")
	suite.Require().Equal(1, len(balance1.OutgoingApprovals), "should have default outgoing approval")

	// Verify versions were initialized (should be 0 on first call)
	version1Incoming := balance1.IncomingApprovals[0].Version
	version1Outgoing := balance1.OutgoingApprovals[0].Version
	suite.Require().True(version1Incoming.IsZero(), "incoming approval version should be 0 on first access")
	suite.Require().True(version1Outgoing.IsZero(), "outgoing approval version should be 0 on first access")

	// Second access - version trackers exist in store, so versions should increment
	// This demonstrates that the version tracker was initialized in the store
	balance2, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice)
	version2Incoming := balance2.IncomingApprovals[0].Version
	version2Outgoing := balance2.OutgoingApprovals[0].Version
	suite.Require().True(version2Incoming.GT(version1Incoming), "incoming approval version should increment because tracker exists in store")
	suite.Require().True(version2Outgoing.GT(version1Outgoing), "outgoing approval version should increment because tracker exists in store")
}

// TestDefaultBalanceVersionIncrement_NoDefaultApprovals verifies that
// version trackers are not initialized when there are no default approvals.
func (suite *DefaultBalanceVersionIncrementTestSuite) TestDefaultBalanceVersionIncrement_NoDefaultApprovals() {
	collection := suite.GetCollection(suite.CollectionId)
	// Collection created in SetupTest has no default approvals

	// Collection has no default approvals
	// First access to balance - should not initialize any version trackers
	balance1, appliedDefault1, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice)
	suite.Require().True(appliedDefault1, "default should be applied on first access")
	suite.Require().Equal(0, len(balance1.IncomingApprovals), "should have no incoming approvals")
	suite.Require().Equal(0, len(balance1.OutgoingApprovals), "should have no outgoing approvals")
}
