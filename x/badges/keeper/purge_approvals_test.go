package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type PurgeApprovalsTestSuite struct {
	TestSuite
}

func TestPurgeApprovalsTestSuite(t *testing.T) {
	suite.Run(t, new(PurgeApprovalsTestSuite))
}

func (suite *PurgeApprovalsTestSuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

func (suite *PurgeApprovalsTestSuite) TestPurgeOwnExpiredApprovals() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection first
	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		BalancesType:   sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{{BadgeIds: GetFullUintRanges()}},
	}})

	now := uint64(time.Now().UnixMilli())
	expiredTime := sdkmath.NewUint(now - 1000000)
	startTime := sdkmath.NewUint(now - 2000000)

	approval := &types.UserOutgoingApproval{
		ToListId:          alice,
		InitiatedByListId: bob,
		TransferTimes:     []*types.UintRange{{Start: startTime, End: expiredTime}},
		BadgeIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "expired-approval",
		ApprovalCriteria: &types.OutgoingApprovalCriteria{
			AutoDeletionOptions: &types.AutoDeletionOptions{},
		},
	}

	suite.T().Logf("Creating approval: %+v", approval)

	err := UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		OutgoingApprovals:       []*types.UserOutgoingApproval{approval},
		UpdateOutgoingApprovals: true,
	})
	suite.Require().NoError(err)

	balance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.T().Logf("Bob's balance after creation: %+v", balance)
	suite.Require().Len(balance.OutgoingApprovals, 1)

	msg := &types.MsgPurgeApprovals{
		Creator:                    bob,
		CollectionId:               sdkmath.NewUint(1),
		PurgeExpired:               true,
		ApproverAddress:            "",
		PurgeCounterpartyApprovals: false,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{
			ApprovalId:      "expired-approval",
			ApprovalLevel:   "outgoing",
			ApproverAddress: bob,
			Version:         sdkmath.NewUint(0),
		}},
	}
	_, err = suite.msgServer.PurgeApprovals(wctx, msg)
	suite.Require().NoError(err)

	balance, err = GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 0)
}

func (suite *PurgeApprovalsTestSuite) TestPurgeAnotherUserExpiredApprovalsNotAllowed() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        alice,
		BalancesType:   sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{{BadgeIds: GetFullUintRanges()}},
	}})

	now := uint64(time.Now().UnixMilli())
	expiredTime := sdkmath.NewUint(now - 1000000)
	startTime := sdkmath.NewUint(now - 2000000)

	approval := &types.UserOutgoingApproval{
		ToListId:          charlie,
		InitiatedByListId: alice,
		TransferTimes:     []*types.UintRange{{Start: startTime, End: expiredTime}},
		BadgeIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "expired-approval",
		ApprovalCriteria: &types.OutgoingApprovalCriteria{
			AutoDeletionOptions: &types.AutoDeletionOptions{},
		},
	}

	err := UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 alice,
		CollectionId:            sdkmath.NewUint(1),
		OutgoingApprovals:       []*types.UserOutgoingApproval{approval},
		UpdateOutgoingApprovals: true,
	})
	suite.Require().NoError(err)

	balance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 1)

	msg := &types.MsgPurgeApprovals{
		Creator:                    bob,
		CollectionId:               sdkmath.NewUint(1),
		PurgeExpired:               true,
		ApproverAddress:            alice,
		PurgeCounterpartyApprovals: false,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{
			ApprovalId:      "expired-approval",
			ApprovalLevel:   "outgoing",
			ApproverAddress: alice,
			Version:         sdkmath.NewUint(0),
		}},
	}
	_, err = suite.msgServer.PurgeApprovals(wctx, msg)
	suite.Require().NoError(err)

	balance, err = GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 1) // Approval should remain
}

func (suite *PurgeApprovalsTestSuite) TestPurgeAnotherUserExpiredApprovalsAllowed() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        alice,
		BalancesType:   sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{{BadgeIds: GetFullUintRanges()}},
	}})

	now := uint64(time.Now().UnixMilli())
	expiredTime := sdkmath.NewUint(now - 1000000)
	startTime := sdkmath.NewUint(now - 2000000)

	approval := &types.UserOutgoingApproval{
		ToListId:          charlie,
		InitiatedByListId: alice,
		TransferTimes:     []*types.UintRange{{Start: startTime, End: expiredTime}},
		BadgeIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "expired-approval",
		ApprovalCriteria: &types.OutgoingApprovalCriteria{
			AutoDeletionOptions: &types.AutoDeletionOptions{AllowPurgeIfExpired: true},
		},
	}

	err := UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 alice,
		CollectionId:            sdkmath.NewUint(1),
		OutgoingApprovals:       []*types.UserOutgoingApproval{approval},
		UpdateOutgoingApprovals: true,
	})
	suite.Require().NoError(err)

	balance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 1)

	msg := &types.MsgPurgeApprovals{
		Creator:                    bob,
		CollectionId:               sdkmath.NewUint(1),
		PurgeExpired:               true,
		ApproverAddress:            alice,
		PurgeCounterpartyApprovals: false,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{
			ApprovalId:      "expired-approval",
			ApprovalLevel:   "outgoing",
			ApproverAddress: alice,
			Version:         sdkmath.NewUint(0),
		}},
	}
	_, err = suite.msgServer.PurgeApprovals(wctx, msg)
	suite.Require().NoError(err)

	balance, err = GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 0)
}

func (suite *PurgeApprovalsTestSuite) TestPurgeCounterpartyApprovalsAllowed() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create address list with alice BEFORE setting the approval
	err := CreateAddressLists(&suite.TestSuite, wctx, &types.MsgCreateAddressLists{
		Creator: bob,
		AddressLists: []*types.AddressList{{
			ListId:    "bobonly",
			Addresses: []string{alice},
			Whitelist: true,
		}},
	})
	suite.Require().NoError(err)

	// Fetch the address list to ensure it exists
	list, err := GetAddressList(&suite.TestSuite, wctx, "bobonly")
	suite.Require().NoError(err, "Address list 'bobonly' should exist after creation")
	suite.Require().Equal("bobonly", list.ListId)

	// Increment block time to ensure address list is committed
	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(1 * time.Second))
	wctx = sdk.WrapSDKContext(suite.ctx)

	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		BalancesType:   sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{{BadgeIds: GetFullUintRanges()}},
	}})

	now := uint64(time.Now().UnixMilli())
	expiredTime := sdkmath.NewUint(now - 1000000)
	startTime := sdkmath.NewUint(now - 2000000)

	approval := &types.UserOutgoingApproval{
		ToListId:          alice,
		InitiatedByListId: "bobonly",
		TransferTimes:     []*types.UintRange{{Start: startTime, End: expiredTime}},
		BadgeIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "counterparty-approval",
		ApprovalCriteria: &types.OutgoingApprovalCriteria{
			AutoDeletionOptions: &types.AutoDeletionOptions{AllowCounterpartyPurge: true},
		},
	}

	err = UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		OutgoingApprovals:       []*types.UserOutgoingApproval{approval},
		UpdateOutgoingApprovals: true,
	})
	suite.Require().NoError(err)

	balance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 1)

	msg := &types.MsgPurgeApprovals{
		Creator:                    alice,
		CollectionId:               sdkmath.NewUint(1),
		PurgeExpired:               false,
		ApproverAddress:            bob,
		PurgeCounterpartyApprovals: true,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{
			ApprovalId:      "counterparty-approval",
			ApprovalLevel:   "outgoing",
			ApproverAddress: bob,
			Version:         sdkmath.NewUint(0),
		}},
	}
	_, err = suite.msgServer.PurgeApprovals(wctx, msg)
	suite.Require().NoError(err)

	balance, err = GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 0)
}

func (suite *PurgeApprovalsTestSuite) TestPurgeNonExpiredApprovals() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	CreateCollections(&suite.TestSuite, wctx, []*types.MsgNewCollection{{
		Creator:        bob,
		BalancesType:   sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{{BadgeIds: GetFullUintRanges()}},
	}})

	now := uint64(time.Now().UnixMilli())
	futureTime := sdkmath.NewUint(now + 1000000)
	startTime := sdkmath.NewUint(now - 1000000)

	approval := &types.UserOutgoingApproval{
		ToListId:          alice,
		InitiatedByListId: bob,
		TransferTimes:     []*types.UintRange{{Start: startTime, End: futureTime}},
		BadgeIds:          GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		ApprovalId:        "future-approval",
		ApprovalCriteria: &types.OutgoingApprovalCriteria{
			AutoDeletionOptions: &types.AutoDeletionOptions{AllowPurgeIfExpired: true},
		},
	}

	err := UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		OutgoingApprovals:       []*types.UserOutgoingApproval{approval},
		UpdateOutgoingApprovals: true,
	})
	suite.Require().NoError(err)

	balance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 1)

	msg := &types.MsgPurgeApprovals{
		Creator:                    bob,
		CollectionId:               sdkmath.NewUint(1),
		PurgeExpired:               true,
		ApproverAddress:            "",
		PurgeCounterpartyApprovals: false,
		ApprovalsToPurge: []*types.ApprovalIdentifierDetails{{
			ApprovalId:      "future-approval",
			ApprovalLevel:   "outgoing",
			ApproverAddress: bob,
			Version:         sdkmath.NewUint(0),
		}},
	}
	_, err = suite.msgServer.PurgeApprovals(wctx, msg)
	suite.Require().NoError(err)

	balance, err = GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().Len(balance.OutgoingApprovals, 1)
}
