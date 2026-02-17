package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type DeleteApprovalTestSuite struct {
	TestSuite
}

func TestDeleteApprovalTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteApprovalTestSuite))
}

func (suite *DeleteApprovalTestSuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

// TestDeleteOutgoingApproval_Success tests successfully deleting an outgoing approval
func (suite *DeleteApprovalTestSuite) TestDeleteOutgoingApproval_Success() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection first
	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		TokensToCreate: []*types.Balance{{TokenIds: GetFullUintRanges()}},
	}})

	// Create an outgoing approval for bob
	approval := &types.UserOutgoingApproval{
		ToListId:          alice,
		InitiatedByListId: bob,
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "test-approval-to-delete",
		ApprovalCriteria:  &types.OutgoingApprovalCriteria{},
	}

	err := UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		OutgoingApprovals:       []*types.UserOutgoingApproval{approval},
		UpdateOutgoingApprovals: true,
	})
	suite.Require().NoError(err)

	// Verify approval exists
	balance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 1)
	suite.Require().Equal("test-approval-to-delete", balance.OutgoingApprovals[0].ApprovalId)

	// Delete the approval
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		ApprovalId:   "test-approval-to-delete",
	}
	_, err = suite.msgServer.DeleteOutgoingApproval(wctx, msg)
	suite.Require().NoError(err)

	// Verify approval is deleted
	balance, err = GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 0)
}

// TestDeleteIncomingApproval_Success tests successfully deleting an incoming approval
func (suite *DeleteApprovalTestSuite) TestDeleteIncomingApproval_Success() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection first
	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		TokensToCreate: []*types.Balance{{TokenIds: GetFullUintRanges()}},
	}})

	// Create an incoming approval for alice
	approval := &types.UserIncomingApproval{
		FromListId:        bob,
		InitiatedByListId: alice,
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "test-incoming-approval",
		ApprovalCriteria:  &types.IncomingApprovalCriteria{},
	}

	err := UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 alice,
		CollectionId:            sdkmath.NewUint(1),
		IncomingApprovals:       []*types.UserIncomingApproval{approval},
		UpdateIncomingApprovals: true,
	})
	suite.Require().NoError(err)

	// Verify approval exists
	balance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().NoError(err)
	suite.Require().Len(balance.IncomingApprovals, 1)

	// Delete the approval
	msg := &types.MsgDeleteIncomingApproval{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		ApprovalId:   "test-incoming-approval",
	}
	_, err = suite.msgServer.DeleteIncomingApproval(wctx, msg)
	suite.Require().NoError(err)

	// Verify approval is deleted
	balance, err = GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().NoError(err)
	suite.Require().Len(balance.IncomingApprovals, 0)
}

// TestDeleteOutgoingApproval_NonExistent tests deleting a non-existent approval
func (suite *DeleteApprovalTestSuite) TestDeleteOutgoingApproval_NonExistent() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection first
	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		TokensToCreate: []*types.Balance{{TokenIds: GetFullUintRanges()}},
	}})

	// Try to delete an approval that doesn't exist
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		ApprovalId:   "non-existent-approval",
	}
	_, err := suite.msgServer.DeleteOutgoingApproval(wctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "approval")
}

// TestDeleteIncomingApproval_NonExistent tests deleting a non-existent incoming approval
func (suite *DeleteApprovalTestSuite) TestDeleteIncomingApproval_NonExistent() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection first
	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		TokensToCreate: []*types.Balance{{TokenIds: GetFullUintRanges()}},
	}})

	// Try to delete an approval that doesn't exist
	msg := &types.MsgDeleteIncomingApproval{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		ApprovalId:   "non-existent-approval",
	}
	_, err := suite.msgServer.DeleteIncomingApproval(wctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "approval")
}

// TestDeleteOutgoingApproval_NonExistentCollection tests deleting from a non-existent collection
func (suite *DeleteApprovalTestSuite) TestDeleteOutgoingApproval_NonExistentCollection() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Try to delete from a collection that doesn't exist
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(999),
		ApprovalId:   "some-approval",
	}
	_, err := suite.msgServer.DeleteOutgoingApproval(wctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "does not exist")
}

// TestDeleteIncomingApproval_NonExistentCollection tests deleting from a non-existent collection
func (suite *DeleteApprovalTestSuite) TestDeleteIncomingApproval_NonExistentCollection() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Try to delete from a collection that doesn't exist
	msg := &types.MsgDeleteIncomingApproval{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(999),
		ApprovalId:   "some-approval",
	}
	_, err := suite.msgServer.DeleteIncomingApproval(wctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "does not exist")
}

// TestDeleteOutgoingApproval_MultipleApprovals tests deleting one approval while keeping others
func (suite *DeleteApprovalTestSuite) TestDeleteOutgoingApproval_MultipleApprovals() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection first
	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		TokensToCreate: []*types.Balance{{TokenIds: GetFullUintRanges()}},
	}})

	// Create multiple outgoing approvals
	approval1 := &types.UserOutgoingApproval{
		ToListId:          alice,
		InitiatedByListId: bob,
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "approval-1",
		ApprovalCriteria:  &types.OutgoingApprovalCriteria{},
	}
	approval2 := &types.UserOutgoingApproval{
		ToListId:          charlie,
		InitiatedByListId: bob,
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "approval-2",
		ApprovalCriteria:  &types.OutgoingApprovalCriteria{},
	}

	err := UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		OutgoingApprovals:       []*types.UserOutgoingApproval{approval1, approval2},
		UpdateOutgoingApprovals: true,
	})
	suite.Require().NoError(err)

	// Verify both approvals exist
	balance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 2)

	// Delete only the first approval
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		ApprovalId:   "approval-1",
	}
	_, err = suite.msgServer.DeleteOutgoingApproval(wctx, msg)
	suite.Require().NoError(err)

	// Verify only approval-2 remains
	balance, err = GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 1)
	suite.Require().Equal("approval-2", balance.OutgoingApprovals[0].ApprovalId)
}

// TestDeleteOutgoingApproval_DeleteTwice tests attempting to delete the same approval twice
func (suite *DeleteApprovalTestSuite) TestDeleteOutgoingApproval_DeleteTwice() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection first
	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		TokensToCreate: []*types.Balance{{TokenIds: GetFullUintRanges()}},
	}})

	// Create an outgoing approval
	approval := &types.UserOutgoingApproval{
		ToListId:          alice,
		InitiatedByListId: bob,
		TransferTimes:     GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "delete-twice-test",
		ApprovalCriteria:  &types.OutgoingApprovalCriteria{},
	}

	err := UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		OutgoingApprovals:       []*types.UserOutgoingApproval{approval},
		UpdateOutgoingApprovals: true,
	})
	suite.Require().NoError(err)

	// Delete the approval first time - should succeed
	msg := &types.MsgDeleteOutgoingApproval{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		ApprovalId:   "delete-twice-test",
	}
	_, err = suite.msgServer.DeleteOutgoingApproval(wctx, msg)
	suite.Require().NoError(err)

	// Try to delete the same approval again - should fail
	_, err = suite.msgServer.DeleteOutgoingApproval(wctx, msg)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "approval")
}
