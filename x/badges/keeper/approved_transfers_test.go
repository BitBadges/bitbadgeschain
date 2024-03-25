package keeper_test

import (
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestDeductFromOutgoing() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, GetCollectionsToCreate())
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	overallTransferBalances := []*types.Balance{
		{
			BadgeIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			Amount:         sdkmath.NewUint(1),
		},
	}

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestDeductFromOutgoingTwoSeparateTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, GetCollectionsToCreate())
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestMaxOneTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestClaimIncrementsExceedsBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.PredeterminedBalances.IncrementedBalances = &types.IncrementedBalances{

		StartBalances: []*types.Balance{
			{
				OwnershipTimes: GetFullUintRanges(),
				BadgeIds:       GetFullUintRanges(),
				Amount:         sdkmath.NewUint(1),
			},
		},
		IncrementBadgeIdsBy:       sdkmath.NewUint(math.MaxUint64),
		IncrementOwnershipTimesBy: sdkmath.NewUint(math.MaxUint64),
	}

	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.PredeterminedBalances.IncrementedBalances = &types.IncrementedBalances{
		StartBalances: []*types.Balance{
			{
				OwnershipTimes: GetFullUintRanges(),
				BadgeIds:       GetFullUintRanges(),
				Amount:         sdkmath.NewUint(1),
			},
		},
		IncrementBadgeIdsBy:       sdkmath.NewUint(math.MaxUint64),
		IncrementOwnershipTimesBy: sdkmath.NewUint(math.MaxUint64),
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	overallTransferBalances := []*types.Balance{
		{
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds:       GetFullUintRanges(),
			Amount:         sdkmath.NewUint(1),
		},
	}

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestRequiresEquals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.RequireToDoesNotEqualInitiatedBy = true
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.RequireFromDoesNotEqualInitiatedBy = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestSpecificApproved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultOutgoingApprovals[0].InitiatedByListId = alice
	collectionsToCreate[0].DefaultIncomingApprovals[0].InitiatedByListId = alice

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

}

func (suite *TestSuite) TestDefaults() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultOutgoingApprovals[0].InitiatedByListId = alice
	collectionsToCreate[0].DefaultIncomingApprovals[0].InitiatedByListId = alice

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestDefaultsNotAutoApplies() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultOutgoingApprovals = []*types.UserOutgoingApproval{}
	collectionsToCreate[0].DefaultIncomingApprovals = []*types.UserIncomingApproval{}
	collectionsToCreate[0].DefaultDisapproveSelfInitiated = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

// func (suite *TestSuite) TestFirstMatchOnly() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	collectionsToCreate[0].DefaultOutgoingApprovals[0].InitiatedByListId = alice
// 	newOutgoingTimeline := []*types.UserOutgoingApproval{
// 			{
// 				ToListId:          "AllWithoutMint",
// 				InitiatedByListId: alice,
// 				TransferTimes:        GetFullUintRanges(),
// 				BadgeIds: []*types.UintRange{
// 					{
// 						Start: sdkmath.NewUint(1),
// 						End:   sdkmath.NewUint(1),
// 					},
// 				},
// 				OwnershipTimes: GetFullUintRanges(),
// 				IsApproved: false,
// 			},
// 			collectionsToCreate[0].DefaultOutgoingApprovals[0],
// 	}
// 	collectionsToCreate[0].DefaultOutgoingApprovals = newOutgoingTimeline

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "error creating badges")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
// 	suite.Require().Error(err, "Error deducting outgoing approvals")

// 	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
// 	suite.Require().Nil(err, "Error deducting outgoing approvals")
// }

// func (suite *TestSuite) TestFirstMatchOnlyWrongTime() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()

// 	newOutgoingTimeline := []*types.UserOutgoingApproval{
// 		{
// 			ToListId:          "AllWithoutMint",
// 			InitiatedByListId: alice,
// 			TransferTimes:        []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 			BadgeIds:             GetFullUintRanges(),
// 			OwnershipTimes: GetFullUintRanges(),
// 			IsApproved: false,

// 		},
// 		collectionsToCreate[0].DefaultOutgoingApprovals[0],
// 	}

// 	collectionsToCreate[0].DefaultOutgoingApprovals = newOutgoingTimeline

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "error creating badges")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
// 	suite.Require().Nil(err, "Error deducting outgoing approvals")
// }

// func (suite *TestSuite) TestCombinations() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()

// 	newOutgoingTimeline :=  []*types.UserOutgoingApproval{
// 			{
// 				ToListId:          "AllWithoutMint",
// 				InitiatedByListId: alice,
// 				TransferTimes:        GetFullUintRanges(),
// 				OwnershipTimes: 			GetFullUintRanges(),
// 				BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},

// 				ApprovalId: "test",
// 				AmountTrackerId:                 "test",
// 				ApprovalCriteria: &types.OutgoingApprovalCriteria{

// 					MaxNumTransfers: &types.MaxNumTransfers{
// 						OverallMaxNumTransfers: sdkmath.NewUint(1000),
// 					},
// 					ApprovalAmounts: &types.ApprovalAmounts{
// 						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
// 					},
// 				},
// 			},
// 			{
// 				ToListId:          "AllWithoutMint",
// 				InitiatedByListId: alice,
// 				TransferTimes:        GetFullUintRanges(),
// 				OwnershipTimes: 			GetFullUintRanges(),
// 				BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 				BadgeIdsOptions: &types.ValueOptions{ InvertDefault: true },
// 				IsApproved:      false,
// 			},
// 	}
// 	collectionsToCreate[0].DefaultOutgoingApprovals = newOutgoingTimeline

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "error creating badges")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
// 	suite.Require().Nil(err, "Error deducting outgoing approvals")

// 	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
// 	suite.Require().Error(err, "Error deducting outgoing approvals")
// }

// func (suite *TestSuite) TestCombinationsOrder() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()

// 	newOutgoingTimeline := []*types.UserOutgoingApproval{
// 			{
// 				ToListId:          "AllWithoutMint",
// 				InitiatedByListId: alice,
// 				TransferTimes:        GetFullUintRanges(),
// 				OwnershipTimes: GetFullUintRanges(),
// 				BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 				IsApproved: false,

// 			},
// 			{
// 					ToListId:          "AllWithoutMint",
// 					InitiatedByListId: alice,
// 					TransferTimes:        GetFullUintRanges(),
// 					OwnershipTimes: GetFullUintRanges(),
// 					BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 					InitiatedByListOptions: &types.ValueOptions{ InvertDefault: true },
// 					IsApproved:         true,
// 					ApprovalId: "test",
// 					AmountTrackerId:                 "test",
// 					ApprovalCriteria: &types.OutgoingApprovalCriteria{
// 							MaxNumTransfers: &types.MaxNumTransfers{
// 								OverallMaxNumTransfers: sdkmath.NewUint(1000),
// 							},
// 							ApprovalAmounts: &types.ApprovalAmounts{
// 								PerFromAddressApprovalAmount: sdkmath.NewUint(1),
// 							},
// 						},
// 					},

// 		}

// 	collectionsToCreate[0].DefaultOutgoingApprovals = newOutgoingTimeline

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "error creating badges")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
// 	suite.Require().Error(err, "Error deducting outgoing approvals")
// }

// func (suite *TestSuite) TestExplicitlyDisallowed() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	collectionsToCreate[0].DefaultOutgoingApprovals[0].InitiatedByListId = alice
// 	newOutgoingTimeline := []*types.UserOutgoingApproval{
// 			collectionsToCreate[0].DefaultOutgoingApprovals[0],
// 	}
// 	collectionsToCreate[0].DefaultOutgoingApprovals = newOutgoingTimeline
// 	collectionsToCreate[0].DefaultOutgoingApprovals = append(collectionsToCreate[0].DefaultOutgoingApprovals, &types.UserOutgoingApproval{

// 			ToListId:          "AllWithoutMint",
// 			InitiatedByListId: "AllWithoutMint",
// 			TransferTimes:        GetFullUintRanges(),
// 			OwnershipTimes: GetFullUintRanges(),
// 			BadgeIds:             GetFullUintRanges(),
// 			IsApproved: false,
// 	})

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "error creating badges")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
// 	suite.Require().Error(err, "Error deducting outgoing approvals")
// }

func (suite *TestSuite) TestNotExplicitlyDefined() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	newOutgoingTimeline := []*types.UserOutgoingApproval{
		{
			ToListId:          "AllWithoutMint",
			InitiatedByListId: alice,
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			BadgeIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},

			ApprovalId:         "test",
			AmountTrackerId:    "test",
			ChallengeTrackerId: "test",
			ApprovalCriteria: &types.OutgoingApprovalCriteria{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(1000),
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					PerFromAddressApprovalAmount: sdkmath.NewUint(1),
				},
			},
		},
	}
	collectionsToCreate[0].DefaultOutgoingApprovals = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturned() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	// collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(2, len(x), "Error deducting outgoing approvals")
	suite.Require().True(x[0].Outgoing != x[1].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesOutgoing() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(1, len(x), "Error deducting outgoing approvals")
	suite.Require().False(x[0].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesIncoming() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(1, len(x), "Error deducting outgoing approvals")
	suite.Require().True(x[0].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesBoth() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(0, len(x), "Error deducting outgoing approvals")
}

//TODO: Test transfer tracker ID after update approved transfers

// func (suite *TestSuite) TestAllowingButManipulatingTrackerIDs() {
// 	wctx := sdk.WrapSDKContext(suite.ctx)

// 	collectionsToCreate := GetCollectionsToCreate()
// 	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
// 		{
// 			FromListId: "AllWithMint",
// 			ToListId:   "AllWithMint",
// 			InitiatedByListId: "AllWithMint",
// 			BadgeIds: 					 GetFullUintRanges(),
// 			TransferTimes:        GetFullUintRanges(),
// 			OwnershipTimes: GetFullUintRanges(),
//
// 			ApprovalId: "test",
// 			AmountTrackerId: 							 "test",
// 			ChallengeTrackerId: "test",
// 			AmountTrackerIdOptions: &types.ValueOptions{ InvertDefault: true },
// 			ApprovalCriteria: &types.ApprovalCriteria{
// 				MaxNumTransfers: &types.MaxNumTransfers{
// 					OverallMaxNumTransfers: sdkmath.NewUint(1000),
// 				},
// 				ApprovalAmounts: &types.ApprovalAmounts{
// 					PerFromAddressApprovalAmount: sdkmath.NewUint(1),
// 				},
// 			},
// 		},
// 	}

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "error creating badges")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

// 	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(
// 		suite.ctx,
// 		[]*types.Balance{},
// 		collection,
// 		GetFullUintRanges(),
// 		GetFullUintRanges(),
// 		bob, alice, alice,
// 		sdkmath.NewUint(1), []*types.MerkleProof{}, &[]string{}, &[]string{}, nil, false)
// 	suite.Require().Error(err, "Error deducting outgoing approvals")
// }
