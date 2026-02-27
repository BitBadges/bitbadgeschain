package approval_criteria_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// AutoDeletionTestSuite tests the AutoDeletionOptions approval criteria field
type AutoDeletionTestSuite struct {
	testutil.AITestSuite
}

func TestAutoDeletionSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(AutoDeletionTestSuite))
}

func (suite *AutoDeletionTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// TestAfterOneUse_True_DeletesApproval tests that afterOneUse=true deletes approval after one use
func (suite *AutoDeletionTestSuite) TestAfterOneUse_True_DeletesApproval() {
	// Create approval with afterOneUse=true
	approval := testutil.GenerateCollectionApproval("one_time_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOneUse: true, // Delete after first use
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Verify approval exists before use
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "one_time_approval" {
			found = true
			break
		}
	}
	suite.Require().True(found, "approval should exist before use")

	// First transfer should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first transfer should succeed")

	// Verify approval is deleted after use
	collectionAfter := suite.GetCollection(collectionId)
	foundAfter := false
	for _, app := range collectionAfter.CollectionApprovals {
		if app.ApprovalId == "one_time_approval" {
			foundAfter = true
			break
		}
	}
	suite.Require().False(foundAfter, "approval should be deleted after one use")

	// Second transfer should fail (no matching approval)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 2)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "second transfer should fail after approval deleted")
}

// TestAfterOneUse_False_PreservesApproval tests that afterOneUse=false preserves approval after use
func (suite *AutoDeletionTestSuite) TestAfterOneUse_False_PreservesApproval() {
	// Create approval with afterOneUse=false (default)
	approval := testutil.GenerateCollectionApproval("persistent_approval", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOneUse: false, // Do not delete after use
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Multiple transfers should all succeed
	for i := 0; i < 3; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d should succeed with persistent approval", i+1)
	}

	// Verify approval still exists
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "persistent_approval" {
			found = true
			break
		}
	}
	suite.Require().True(found, "approval should still exist")
}

// TestAfterOverallMaxNumTransfers_True_DeletesAtLimit tests that afterOverallMaxNumTransfers=true deletes at limit
func (suite *AutoDeletionTestSuite) TestAfterOverallMaxNumTransfers_True_DeletesAtLimit() {
	// Create approval with max 3 transfers and auto-delete when reached
	approval := testutil.GenerateCollectionApproval("limited_auto_delete", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(3),
			AmountTrackerId:        "delete_tracker",
		},
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOverallMaxNumTransfers: true, // Delete after reaching max transfers
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Make 3 transfers (reach the limit)
	for i := 0; i < 3; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d should succeed", i+1)
	}

	// Verify approval is deleted after reaching max
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "limited_auto_delete" {
			found = true
			break
		}
	}
	suite.Require().False(found, "approval should be deleted after reaching max transfers")
}

// TestAfterOverallMaxNumTransfers_False_PreservesAtLimit tests that approval is preserved when false
func (suite *AutoDeletionTestSuite) TestAfterOverallMaxNumTransfers_False_PreservesAtLimit() {
	// Create approval with max 3 transfers but no auto-delete
	approval := testutil.GenerateCollectionApproval("limited_no_delete", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(3),
			AmountTrackerId:        "no_delete_tracker",
		},
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOverallMaxNumTransfers: false, // Don't delete after reaching max
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Make 3 transfers (reach the limit)
	for i := 0; i < 3; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d should succeed", i+1)
	}

	// Verify approval still exists (just exhausted)
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "limited_no_delete" {
			found = true
			break
		}
	}
	suite.Require().True(found, "approval should still exist even after reaching max")

	// 4th transfer should still fail (limit reached, not deleted)
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 4)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "4th transfer should fail (limit reached)")
}

// TestApprovalRemovedFromCollection tests that approval is actually removed from collection
func (suite *AutoDeletionTestSuite) TestApprovalRemovedFromCollection() {
	// Create two approvals - one with auto-delete, one without
	approval1 := testutil.GenerateCollectionApproval("auto_delete_me", "AllWithoutMint", "All")
	approval1.ApprovalCriteria = &types.ApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOneUse: true,
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	approval2 := testutil.GenerateCollectionApproval("keep_me", "AllWithoutMint", "All")
	approval2.ApprovalCriteria = &types.ApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOneUse: false,
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval1, approval2})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Verify both exist initially
	collectionBefore := suite.GetCollection(collectionId)
	countBefore := 0
	for _, app := range collectionBefore.CollectionApprovals {
		if app.ApprovalId == "auto_delete_me" || app.ApprovalId == "keep_me" {
			countBefore++
		}
	}
	suite.Require().Equal(2, countBefore, "both approvals should exist initially")

	// Use the auto-delete approval
	msg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "transfer should succeed")

	// Verify only keep_me exists now (auto_delete_me was used first due to first-match)
	collectionAfter := suite.GetCollection(collectionId)
	hasAutoDelete := false
	hasKeep := false
	for _, app := range collectionAfter.CollectionApprovals {
		if app.ApprovalId == "auto_delete_me" {
			hasAutoDelete = true
		}
		if app.ApprovalId == "keep_me" {
			hasKeep = true
		}
	}
	suite.Require().False(hasAutoDelete, "auto_delete_me should be removed")
	suite.Require().True(hasKeep, "keep_me should still exist")
}

// TestAutoDeletionWithMint tests auto-deletion with mint approvals
func (suite *AutoDeletionTestSuite) TestAutoDeletionWithMint() {
	// Create mint approval with auto-delete after one use
	approval := testutil.GenerateCollectionApproval("one_time_mint", types.MintAddress, "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOneUse: true,
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// First mint should succeed
	msg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first mint should succeed")

	// Verify mint approval is deleted
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "one_time_mint" {
			found = true
			break
		}
	}
	suite.Require().False(found, "one-time mint approval should be deleted")

	// Second mint should fail
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "second mint should fail after approval deleted")
}

// TestAutoDeletionOptions_AllFalse tests behavior when all auto-deletion flags are false
func (suite *AutoDeletionTestSuite) TestAutoDeletionOptions_AllFalse() {
	// Create approval with all auto-deletion flags false
	approval := testutil.GenerateCollectionApproval("no_auto_delete", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOneUse:                 false,
			AfterOverallMaxNumTransfers: false,
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Multiple transfers should all succeed and approval should persist
	for i := 0; i < 5; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d should succeed", i+1)
	}

	// Verify approval still exists
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "no_auto_delete" {
			found = true
			break
		}
	}
	suite.Require().True(found, "approval should persist with all auto-delete flags false")
}

// TestAutoDeletionOptions_Nil tests behavior when AutoDeletionOptions is nil
func (suite *AutoDeletionTestSuite) TestAutoDeletionOptions_Nil() {
	// Create approval with nil AutoDeletionOptions
	approval := testutil.GenerateCollectionApproval("nil_auto_delete", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		AutoDeletionOptions:            nil, // Nil auto-deletion options
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// Multiple transfers should succeed (nil means no auto-deletion)
	for i := 0; i < 3; i++ {
		msg := &types.MsgTransferTokens{
			Creator:      suite.Alice,
			CollectionId: collectionId,
			Transfers: []*types.Transfer{
				{
					From:        suite.Alice,
					ToAddresses: []string{suite.Bob},
					Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, uint64(i+1))},
				},
			},
		}
		_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg)
		suite.Require().NoError(err, "transfer %d should succeed with nil AutoDeletionOptions", i+1)
	}

	// Verify approval still exists
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "nil_auto_delete" {
			found = true
			break
		}
	}
	suite.Require().True(found, "approval should persist with nil AutoDeletionOptions")
}

// TestAllowCounterpartyPurge tests the allowCounterpartyPurge flag
func (suite *AutoDeletionTestSuite) TestAllowCounterpartyPurge() {
	// Create approval that allows counterparty purge
	approval := testutil.GenerateCollectionApproval("counterparty_purgeable", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AllowCounterpartyPurge: true, // Counterparty can purge
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Verify the flag was set correctly
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "counterparty_purgeable" {
			found = true
			suite.Require().NotNil(app.ApprovalCriteria.AutoDeletionOptions, "auto deletion options should exist")
			suite.Require().True(app.ApprovalCriteria.AutoDeletionOptions.AllowCounterpartyPurge,
				"allowCounterpartyPurge should be true")
			break
		}
	}
	suite.Require().True(found, "approval should exist")
}

// TestAllowPurgeIfExpired tests the allowPurgeIfExpired flag
func (suite *AutoDeletionTestSuite) TestAllowPurgeIfExpired() {
	// Create approval that allows purge if expired
	approval := testutil.GenerateCollectionApproval("expired_purgeable", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AllowPurgeIfExpired: true, // Others can purge if expired
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})

	// Verify the flag was set correctly
	collection := suite.GetCollection(collectionId)
	found := false
	for _, app := range collection.CollectionApprovals {
		if app.ApprovalId == "expired_purgeable" {
			found = true
			suite.Require().NotNil(app.ApprovalCriteria.AutoDeletionOptions, "auto deletion options should exist")
			suite.Require().True(app.ApprovalCriteria.AutoDeletionOptions.AllowPurgeIfExpired,
				"allowPurgeIfExpired should be true")
			break
		}
	}
	suite.Require().True(found, "approval should exist")
}

// TestCombinedAutoDeletionFlags tests combined auto-deletion flag scenarios
func (suite *AutoDeletionTestSuite) TestCombinedAutoDeletionFlags() {
	// Create approval with max transfers AND auto-delete when reached
	approval := testutil.GenerateCollectionApproval("combined_flags", "AllWithoutMint", "All")
	approval.ApprovalCriteria = &types.ApprovalCriteria{
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(2),
			AmountTrackerId:        "combined_tracker",
		},
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOneUse:                 false, // Don't delete after one use
			AfterOverallMaxNumTransfers: true,  // Delete after reaching max
			AllowCounterpartyPurge:      true,
			AllowPurgeIfExpired:         true,
		},
		OverridesFromOutgoingApprovals: true,
		OverridesToIncomingApprovals:   true,
	}

	collectionId := suite.CreateTestCollectionWithApprovals(suite.Manager, []*types.CollectionApproval{approval})
	suite.MintTokens(collectionId, suite.Alice, []*types.Balance{testutil.GenerateBalance(100, 1, 100, 1, math.MaxUint64)})

	// First transfer - should succeed, approval should persist
	msg1 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 1)},
			},
		},
	}
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg1)
	suite.Require().NoError(err, "first transfer should succeed")

	// Verify approval still exists
	collectionAfter1 := suite.GetCollection(collectionId)
	found1 := false
	for _, app := range collectionAfter1.CollectionApprovals {
		if app.ApprovalId == "combined_flags" {
			found1 = true
			break
		}
	}
	suite.Require().True(found1, "approval should still exist after first use")

	// Second transfer - should succeed, then approval should be deleted (max reached)
	msg2 := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances:    []*types.Balance{testutil.GenerateSimpleBalance(10, 2)},
			},
		},
	}
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().NoError(err, "second transfer should succeed")

	// Verify approval is now deleted
	collectionAfter2 := suite.GetCollection(collectionId)
	found2 := false
	for _, app := range collectionAfter2.CollectionApprovals {
		if app.ApprovalId == "combined_flags" {
			found2 = true
			break
		}
	}
	suite.Require().False(found2, "approval should be deleted after max transfers reached")
}
