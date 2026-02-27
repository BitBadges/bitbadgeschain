package msg_handlers_test

import (
	"math"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type PurgeApprovalsTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestPurgeApprovalsSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(PurgeApprovalsTestSuite))
}

func (suite *PurgeApprovalsTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// createExpiredApproval creates an approval with transfer times in the past
func (suite *PurgeApprovalsTestSuite) createExpiredApproval(approvalId string, toListId string) *types.UserOutgoingApproval {
	// Current time in milliseconds
	currentTimeMs := uint64(suite.Ctx.BlockTime().UnixMilli())

	return &types.UserOutgoingApproval{
		ApprovalId:        approvalId,
		ToListId:          toListId,
		InitiatedByListId: "All",
		TransferTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(currentTimeMs - 1000)}, // Ended 1 second ago
		},
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		ApprovalCriteria: &types.OutgoingApprovalCriteria{},
		Version:          sdkmath.NewUint(0),
	}
}

// createActiveApproval creates an approval with transfer times in the future
func (suite *PurgeApprovalsTestSuite) createActiveApproval(approvalId string, toListId string) *types.UserOutgoingApproval {
	return testutil.GenerateUserOutgoingApproval(approvalId, toListId)
}

// createExpiredIncomingApproval creates an incoming approval with transfer times in the past
func (suite *PurgeApprovalsTestSuite) createExpiredIncomingApproval(approvalId string, fromListId string) *types.UserIncomingApproval {
	// Current time in milliseconds
	currentTimeMs := uint64(suite.Ctx.BlockTime().UnixMilli())

	return &types.UserIncomingApproval{
		ApprovalId:        approvalId,
		FromListId:        fromListId,
		InitiatedByListId: "All",
		TransferTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(currentTimeMs - 1000)}, // Ended 1 second ago
		},
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		ApprovalCriteria: &types.IncomingApprovalCriteria{},
		Version:          sdkmath.NewUint(0),
	}
}

// TestPurgeApprovals_PurgeExpiredRemovesExpiredApprovals tests that purgeExpired removes approvals with past transferTimes
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_PurgeExpiredRemovesExpiredApprovals() {
	// Create an expired outgoing approval
	expiredApproval := suite.createExpiredApproval("expired1", "All")

	// Set the approval for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{expiredApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get the approval version after setting
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.OutgoingApprovals))
	approvalVersion := balance.OutgoingApprovals[0].Version

	// Purge expired approvals (self-purge)
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expired1",
				ApprovalLevel:   "outgoing",
				ApproverAddress: suite.Alice,
				Version:         approvalVersion,
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err, "purging expired approvals should succeed")
	suite.Require().True(resp.NumPurged.GT(sdkmath.NewUint(0)), "should have purged at least one approval")

	// Verify approval was removed
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balanceAfter.OutgoingApprovals), "expired approval should be purged")
}

// TestPurgeApprovals_ApproverAddressSpecifiesTarget tests that approverAddress specifies whose approvals to purge
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_ApproverAddressSpecifiesTarget() {
	// Create an expired approval for Alice
	expiredApproval := suite.createExpiredApproval("expired1", "All")

	// Set the approval for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{expiredApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get the approval version
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.OutgoingApprovals))
	approvalVersion := balance.OutgoingApprovals[0].Version

	// Try to purge Alice's approvals as Alice (self-purge, no approverAddress needed)
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:         suite.Alice,
		CollectionId:    suite.CollectionId,
		PurgeExpired:    true,
		ApproverAddress: "", // Empty means target is creator (Alice)
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expired1",
				ApprovalLevel:   "outgoing",
				ApproverAddress: suite.Alice,
				Version:         approvalVersion,
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	suite.Require().True(resp.NumPurged.GT(sdkmath.NewUint(0)))

	// Verify Alice's approval was purged
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balanceAfter.OutgoingApprovals))
}

// TestPurgeApprovals_CounterpartyPurgeWorks tests that purgeCounterpartyApprovals works
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_CounterpartyPurgeWorks() {
	// Create an expired incoming approval for Alice that allows counterparty purge
	expiredApproval := suite.createExpiredIncomingApproval("expired1", suite.Bob) // Only from Bob
	expiredApproval.ApprovalCriteria = &types.IncomingApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AllowCounterpartyPurge: true,
			AllowPurgeIfExpired:    true,
		},
	}
	expiredApproval.InitiatedByListId = suite.Bob // Only Bob can initiate

	// Set the approval for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{expiredApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get the approval version
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.IncomingApprovals))
	approvalVersion := balance.IncomingApprovals[0].Version

	// Bob tries to purge Alice's approval using counterparty purge
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:                    suite.Bob,
		CollectionId:               suite.CollectionId,
		PurgeExpired:               false,
		ApproverAddress:            suite.Alice,
		PurgeCounterpartyApprovals: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expired1",
				ApprovalLevel:   "incoming",
				ApproverAddress: suite.Alice,
				Version:         approvalVersion,
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err, "counterparty purge should succeed")
	suite.Require().True(resp.NumPurged.GT(sdkmath.NewUint(0)), "should have purged approval")

	// Verify approval was purged
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balanceAfter.IncomingApprovals))
}

// TestPurgeApprovals_ApprovalsToPurgeArray tests that approvalsToPurge array specifies specific approvals
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_ApprovalsToPurgeArray() {
	// Create two expired approvals
	expired1 := suite.createExpiredApproval("expired1", "All")
	expired2 := suite.createExpiredApproval("expired2", suite.Bob)

	// Set both approvals for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{expired1, expired2},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get approval versions
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(2, len(balance.OutgoingApprovals))

	var version1, version2 sdkmath.Uint
	for _, a := range balance.OutgoingApprovals {
		if a.ApprovalId == "expired1" {
			version1 = a.Version
		} else if a.ApprovalId == "expired2" {
			version2 = a.Version
		}
	}

	// Only purge expired1, not expired2
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expired1",
				ApprovalLevel:   "outgoing",
				ApproverAddress: suite.Alice,
				Version:         version1,
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(1), resp.NumPurged, "should purge exactly 1 approval")

	// Verify only expired1 was purged, expired2 remains
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balanceAfter.OutgoingApprovals))
	suite.Require().Equal("expired2", balanceAfter.OutgoingApprovals[0].ApprovalId)

	// Now purge expired2
	purgeMsg2 := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expired2",
				ApprovalLevel:   "outgoing",
				ApproverAddress: suite.Alice,
				Version:         version2,
			},
		},
	}

	resp2, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg2)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(1), resp2.NumPurged)

	// Verify all purged
	balanceFinal := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balanceFinal.OutgoingApprovals))
}

// TestPurgeApprovals_RespectsAllowCounterpartyPurge tests that allowCounterpartyPurge setting is respected
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_RespectsAllowCounterpartyPurge() {
	// Create an expired approval WITHOUT allowCounterpartyPurge
	expiredApproval := suite.createExpiredIncomingApproval("expired1", suite.Bob)
	expiredApproval.ApprovalCriteria = &types.IncomingApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AllowCounterpartyPurge: false, // Not allowed
			AllowPurgeIfExpired:    false,
		},
	}
	expiredApproval.InitiatedByListId = suite.Bob

	// Set the approval for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{expiredApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get the approval version
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.IncomingApprovals))
	approvalVersion := balance.IncomingApprovals[0].Version

	// Bob tries to purge Alice's approval - should not be able to purge due to settings
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:                    suite.Bob,
		CollectionId:               suite.CollectionId,
		PurgeExpired:               false,
		ApproverAddress:            suite.Alice,
		PurgeCounterpartyApprovals: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expired1",
				ApprovalLevel:   "incoming",
				ApproverAddress: suite.Alice,
				Version:         approvalVersion,
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	// NumPurged should be 0 because counterparty purge is not allowed
	suite.Require().True(resp.NumPurged.Equal(sdkmath.NewUint(0)), "should not purge when allowCounterpartyPurge is false")

	// Verify approval still exists
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balanceAfter.IncomingApprovals))
}

// TestPurgeApprovals_RespectsAllowPurgeIfExpired tests that allowPurgeIfExpired setting is respected
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_RespectsAllowPurgeIfExpired() {
	// Create an expired approval WITH allowPurgeIfExpired
	expiredApproval := suite.createExpiredIncomingApproval("expired1", "All")
	expiredApproval.ApprovalCriteria = &types.IncomingApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AllowCounterpartyPurge: false,
			AllowPurgeIfExpired:    true, // Allow others to purge if expired
		},
	}

	// Set the approval for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{expiredApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get the approval version
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.IncomingApprovals))
	approvalVersion := balance.IncomingApprovals[0].Version

	// Bob tries to purge Alice's expired approval - should succeed due to AllowPurgeIfExpired
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:                    suite.Bob,
		CollectionId:               suite.CollectionId,
		PurgeExpired:               false,
		ApproverAddress:            suite.Alice,
		PurgeCounterpartyApprovals: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expired1",
				ApprovalLevel:   "incoming",
				ApproverAddress: suite.Alice,
				Version:         approvalVersion,
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	suite.Require().True(resp.NumPurged.GT(sdkmath.NewUint(0)), "should purge when allowPurgeIfExpired is true")

	// Verify approval was purged
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balanceAfter.IncomingApprovals))
}

// TestPurgeApprovals_ReturnsNumPurgedCorrectly tests that numPurged count is returned correctly
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_ReturnsNumPurgedCorrectly() {
	// Create multiple expired approvals
	expired1 := suite.createExpiredApproval("expired1", "All")
	expired2 := suite.createExpiredApproval("expired2", suite.Bob)
	expired3 := suite.createExpiredApproval("expired3", suite.Charlie)

	// Set all approvals for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{expired1, expired2, expired3},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get approval versions
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(3, len(balance.OutgoingApprovals))

	var version1, version2, version3 sdkmath.Uint
	for _, a := range balance.OutgoingApprovals {
		switch a.ApprovalId {
		case "expired1":
			version1 = a.Version
		case "expired2":
			version2 = a.Version
		case "expired3":
			version3 = a.Version
		}
	}

	// Purge all three approvals
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{ApprovalId: "expired1", ApprovalLevel: "outgoing", ApproverAddress: suite.Alice, Version: version1},
			{ApprovalId: "expired2", ApprovalLevel: "outgoing", ApproverAddress: suite.Alice, Version: version2},
			{ApprovalId: "expired3", ApprovalLevel: "outgoing", ApproverAddress: suite.Alice, Version: version3},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(3), resp.NumPurged, "numPurged should be 3")

	// Verify all purged
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balanceAfter.OutgoingApprovals))
}

// TestPurgeApprovals_ActiveApprovalsNotPurged tests that active (non-expired) approvals are not purged
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_ActiveApprovalsNotPurged() {
	// Create an active (non-expired) approval
	activeApproval := suite.createActiveApproval("active1", "All")

	// Set the approval for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{activeApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get the approval version
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.OutgoingApprovals))
	approvalVersion := balance.OutgoingApprovals[0].Version

	// Try to purge the active approval
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "active1",
				ApprovalLevel:   "outgoing",
				ApproverAddress: suite.Alice,
				Version:         approvalVersion,
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	// Active approvals should not be purged
	suite.Require().Equal(sdkmath.NewUint(0), resp.NumPurged, "active approval should not be purged")

	// Verify approval still exists
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balanceAfter.OutgoingApprovals))
}

// TestPurgeApprovals_SelfPurgeAlwaysAllowed tests that self-purge of expired approvals is always allowed
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_SelfPurgeAlwaysAllowed() {
	// Create an expired approval WITHOUT any auto-deletion options
	expiredApproval := suite.createExpiredApproval("expired1", "All")
	expiredApproval.ApprovalCriteria = &types.OutgoingApprovalCriteria{
		// No AutoDeletionOptions set
	}

	// Set the approval for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{expiredApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get the approval version
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.OutgoingApprovals))
	approvalVersion := balance.OutgoingApprovals[0].Version

	// Self-purge should always work for expired approvals
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expired1",
				ApprovalLevel:   "outgoing",
				ApproverAddress: suite.Alice,
				Version:         approvalVersion,
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	suite.Require().True(resp.NumPurged.GT(sdkmath.NewUint(0)), "self-purge of expired approvals should always work")

	// Verify approval was purged
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balanceAfter.OutgoingApprovals))
}

// TestPurgeApprovals_VersionMismatchDoesNotPurge tests that wrong version prevents purge
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_VersionMismatchDoesNotPurge() {
	// Create an expired approval
	expiredApproval := suite.createExpiredApproval("expired1", "All")

	// Set the approval for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{expiredApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Try to purge with wrong version
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expired1",
				ApprovalLevel:   "outgoing",
				ApproverAddress: suite.Alice,
				Version:         sdkmath.NewUint(999), // Wrong version
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	// Version mismatch should prevent purge
	suite.Require().Equal(sdkmath.NewUint(0), resp.NumPurged, "version mismatch should prevent purge")

	// Verify approval still exists
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balanceAfter.OutgoingApprovals))
}

// TestPurgeApprovals_InvalidCollectionFails tests that purge fails for non-existent collections
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_InvalidCollectionFails() {
	invalidCollectionId := sdkmath.NewUint(99999)

	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: invalidCollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "test",
				ApprovalLevel:   "outgoing",
				ApproverAddress: suite.Alice,
				Version:         sdkmath.NewUint(0),
			},
		},
	}

	_, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().Error(err, "purge should fail for non-existent collection")
}

// TestPurgeApprovals_IncomingApprovals tests purging incoming approvals
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_IncomingApprovals() {
	// Create an expired incoming approval
	expiredApproval := suite.createExpiredIncomingApproval("expiredIncoming1", "All")

	// Set the approval for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{expiredApproval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get the approval version
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(1, len(balance.IncomingApprovals))
	approvalVersion := balance.IncomingApprovals[0].Version

	// Purge the incoming approval
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{
				ApprovalId:      "expiredIncoming1",
				ApprovalLevel:   "incoming",
				ApproverAddress: suite.Alice,
				Version:         approvalVersion,
			},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(1), resp.NumPurged)

	// Verify incoming approval was purged
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balanceAfter.IncomingApprovals))
}

// TestPurgeApprovals_MixedIncomingAndOutgoing tests purging both incoming and outgoing approvals
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_MixedIncomingAndOutgoing() {
	// Create expired incoming and outgoing approvals
	expiredIncoming := suite.createExpiredIncomingApproval("expiredIncoming1", "All")
	expiredOutgoing := suite.createExpiredApproval("expiredOutgoing1", "All")

	// Set both approvals for Alice
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{expiredIncoming},
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{expiredOutgoing},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get approval versions
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	incomingVersion := balance.IncomingApprovals[0].Version
	outgoingVersion := balance.OutgoingApprovals[0].Version

	// Purge both
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{ApprovalId: "expiredIncoming1", ApprovalLevel: "incoming", ApproverAddress: suite.Alice, Version: incomingVersion},
			{ApprovalId: "expiredOutgoing1", ApprovalLevel: "outgoing", ApproverAddress: suite.Alice, Version: outgoingVersion},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(2), resp.NumPurged)

	// Verify both were purged
	balanceAfter := suite.GetBalance(suite.CollectionId, suite.Alice)
	suite.Require().Equal(0, len(balanceAfter.IncomingApprovals))
	suite.Require().Equal(0, len(balanceAfter.OutgoingApprovals))
}

// TestPurgeApprovals_TimeAdvanceExpiration tests that approvals expire correctly with time advancement
func (suite *PurgeApprovalsTestSuite) TestPurgeApprovals_TimeAdvanceExpiration() {
	// Create an approval that expires 1 hour from now
	currentTimeMs := uint64(suite.Ctx.BlockTime().UnixMilli())
	oneHourMs := uint64(time.Hour.Milliseconds())

	approval := &types.UserOutgoingApproval{
		ApprovalId:        "soonExpired",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(currentTimeMs + oneHourMs)},
		},
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		OwnershipTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		ApprovalCriteria: &types.OutgoingApprovalCriteria{},
		Version:          sdkmath.NewUint(0),
	}

	// Set the approval
	setMsg := &types.MsgUpdateUserApprovals{
		Creator:                 suite.Alice,
		CollectionId:            suite.CollectionId,
		UpdateOutgoingApprovals: true,
		OutgoingApprovals:       []*types.UserOutgoingApproval{approval},
	}

	_, err := suite.MsgServer.UpdateUserApprovals(sdk.WrapSDKContext(suite.Ctx), setMsg)
	suite.Require().NoError(err)

	// Get version
	balance := suite.GetBalance(suite.CollectionId, suite.Alice)
	approvalVersion := balance.OutgoingApprovals[0].Version

	// Try to purge now - should fail because not expired yet
	purgeMsg := &types.MsgPurgeApprovals{
		Creator:      suite.Alice,
		CollectionId: suite.CollectionId,
		PurgeExpired: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{
			{ApprovalId: "soonExpired", ApprovalLevel: "outgoing", ApproverAddress: suite.Alice, Version: approvalVersion},
		},
	}

	resp, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(0), resp.NumPurged, "should not purge active approval")

	// Advance time by 2 hours
	suite.AdvanceBlockTime(2 * time.Hour)

	// Now try to purge - should succeed because expired
	resp2, err := suite.MsgServer.PurgeApprovals(sdk.WrapSDKContext(suite.Ctx), purgeMsg)
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(1), resp2.NumPurged, "should purge expired approval after time advance")
}
