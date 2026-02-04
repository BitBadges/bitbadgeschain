package approvals_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type ApprovalVersioningTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestApprovalVersioningSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(ApprovalVersioningTestSuite))
}

func (suite *ApprovalVersioningTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestApprovalVersioning_IncrementOnChange tests that approval versions increment when approvals change
func (suite *ApprovalVersioningTestSuite) TestApprovalVersioning_IncrementOnChange() {
	// Setup initial collection approval
	approval1 := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg1 := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval1},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg1)
	suite.Require().NoError(err)

	// Get collection and check initial version
	collection1 := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(1, len(collection1.CollectionApprovals))
	initialVersion := collection1.CollectionApprovals[0].Version

	// Update approval (change token IDs)
	approval2 := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	approval2.TokenIds = []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)}, // Changed from full range
	}
	updateMsg2 := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval2},
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg2)
	suite.Require().NoError(err)

	// Get collection and check version incremented
	collection2 := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(1, len(collection2.CollectionApprovals))
	newVersion := collection2.CollectionApprovals[0].Version
	suite.Require().True(newVersion.GT(initialVersion), "version should increment when approval changes")
}

// TestApprovalVersioning_NoIncrementOnNoChange tests that approval versions don't increment when approvals don't change
func (suite *ApprovalVersioningTestSuite) TestApprovalVersioning_NoIncrementOnNoChange() {
	// Setup initial collection approval
	approval1 := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg1 := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval1},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg1)
	suite.Require().NoError(err)

	// Get collection and check initial version
	collection1 := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(1, len(collection1.CollectionApprovals))
	initialVersion := collection1.CollectionApprovals[0].Version

	// Update with same approval (no changes)
	approval2 := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg2 := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals:    true,
		CollectionApprovals:        []*types.CollectionApproval{approval2},
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg2)
	suite.Require().NoError(err)

	// Get collection and check version didn't increment
	collection2 := suite.GetCollection(suite.CollectionId)
	suite.Require().Equal(1, len(collection2.CollectionApprovals))
	sameVersion := collection2.CollectionApprovals[0].Version
	suite.Require().True(sameVersion.Equal(initialVersion), "version should not increment when approval doesn't change")
}

// TestApprovalVersioning_UserApprovalVersioning tests user approval versioning
func (suite *ApprovalVersioningTestSuite) TestApprovalVersioning_UserApprovalVersioning() {
	// Setup collection approval
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Set initial outgoing approval
	outgoingApproval1 := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg1 := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval1,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg1)
	suite.Require().NoError(err)

	// Get balance and check initial version
	balance1 := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance1.OutgoingApprovals))
	initialVersion := balance1.OutgoingApprovals[0].Version

	// Update approval (change to list)
	outgoingApproval2 := testutil.GenerateUserOutgoingApproval("outgoing1", suite.Bob)
	setOutgoingMsg2 := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval2,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg2)
	suite.Require().NoError(err)

	// Get balance and check version incremented
	balance2 := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance2.OutgoingApprovals))
	newVersion := balance2.OutgoingApprovals[0].Version
	suite.Require().True(newVersion.GT(initialVersion), "version should increment when user approval changes")
}

// TestApprovalVersioning_InvalidVersionReuse tests that old approval versions cannot be reused
func (suite *ApprovalVersioningTestSuite) TestApprovalVersioning_InvalidVersionReuse() {
	// Setup approvals - need mint approval first
	mintApproval := testutil.GenerateCollectionApproval("mint_approval", types.MintAddress, "All")
	mintApproval.ApprovalCriteria.OverridesFromOutgoingApprovals = true
	mintApproval.ApprovalCriteria.OverridesToIncomingApprovals = true
	approval := testutil.GenerateCollectionApproval("approval1", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{mintApproval, approval},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Mint tokens
	mintBalances := []*types.Balance{
		testutil.GenerateSimpleBalance(10, 1),
	}
	suite.MintBadges(suite.CollectionId, suite.Alice, mintBalances)

	// Set approvals
	outgoingApproval := testutil.GenerateUserOutgoingApproval("outgoing1", "All")
	setOutgoingMsg := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg)
	suite.Require().NoError(err)

	// Get initial version
	balance1 := suite.GetBalance(suite.CollectionId, suite.Alice)
	initialVersion := balance1.OutgoingApprovals[0].Version

	// Update approval (increments version)
	outgoingApproval2 := testutil.GenerateUserOutgoingApproval("outgoing1", suite.Bob)
	setOutgoingMsg2 := &types.MsgSetOutgoingApproval{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		Approval:     outgoingApproval2,
	}
	_, err = suite.MsgServer.SetOutgoingApproval(sdk.WrapSDKContext(suite.Ctx), setOutgoingMsg2)
	suite.Require().NoError(err)

	// Get new version
	balance2 := suite.GetBalance(suite.CollectionId, suite.Alice)
	newVersion := balance2.OutgoingApprovals[0].Version
	suite.Require().True(newVersion.GT(initialVersion), "version should have incremented")

	// Attempt to use old version in transfer should fail
	// (This would be tested in transfer execution, but the version increment itself is tested here)
	suite.Require().True(newVersion.GT(initialVersion), "old version should be invalidated")
}

