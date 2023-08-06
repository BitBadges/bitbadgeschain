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
			BadgeIds:      GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			Amount:        sdkmath.NewUint(1),
		},
	}

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	_, err = suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestDeductFromOutgoingTwoSeparateTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, GetCollectionsToCreate())
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)



	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestMaxOneTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].ApprovalDetails[0].MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestClaimIncrementsExceedsBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].ApprovalDetails[0].PredeterminedBalances = &types.PredeterminedBalances{
		OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].ApprovalDetails[0].PredeterminedBalances.IncrementedBalances = &types.IncrementedBalances{
		StartBalances: []*types.Balance{
			{
				OwnershipTimes: GetFullUintRanges(),
				BadgeIds:       GetFullUintRanges(),
				Amount:         sdkmath.NewUint(1),
			},
		},
		IncrementBadgeIdsBy: sdkmath.NewUint(math.MaxUint64),
		IncrementOwnershipTimesBy: sdkmath.NewUint(math.MaxUint64),
	}
	
	collectionsToCreate[0].DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].ApprovalDetails[0].PredeterminedBalances = &types.PredeterminedBalances{
		OrderCalculationMethod:  &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}
	collectionsToCreate[0].DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].ApprovalDetails[0].PredeterminedBalances.IncrementedBalances = &types.IncrementedBalances{
		StartBalances: []*types.Balance{
			{
				OwnershipTimes: GetFullUintRanges(),
				BadgeIds:             GetFullUintRanges(),
				Amount:               sdkmath.NewUint(1),
			},
		},
		IncrementBadgeIdsBy: sdkmath.NewUint(math.MaxUint64),
		IncrementOwnershipTimesBy: sdkmath.NewUint(math.MaxUint64),
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	overallTransferBalances := []*types.Balance{
		{
			OwnershipTimes: GetFullUintRanges(),
			BadgeIds:             GetFullUintRanges(),
			Amount:               sdkmath.NewUint(1),
		},
	}

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestRequiresEquals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].ApprovalDetails[0].RequireToDoesNotEqualInitiatedBy = true
	collectionsToCreate[0].DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].ApprovalDetails[0].RequireFromDoesNotEqualInitiatedBy = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestSpecificApproved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].InitiatedByMappingId = alice
	collectionsToCreate[0].DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].InitiatedByMappingId = alice

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

}

func (suite *TestSuite) TestDefaults() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].InitiatedByMappingId = alice
	collectionsToCreate[0].DefaultApprovedIncomingTransfersTimeline[0].ApprovedIncomingTransfers[0].InitiatedByMappingId = alice

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserIncomingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestFirstMatchOnly() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0].InitiatedByMappingId = alice
	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{{
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
			{
				ToMappingId:          "AllWithoutMint",
				InitiatedByMappingId: alice,
				TransferTimes:        GetFullUintRanges(),
				BadgeIds: []*types.UintRange{
					{
						Start: sdkmath.NewUint(1),
						End:   sdkmath.NewUint(1),
					},
				},
				OwnershipTimes: GetFullUintRanges(),
				AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
					{
						IsApproved: false,
					},
				},
				ApprovalDetails: []*types.OutgoingApprovalDetails{
					{
						MerkleChallenges:                []*types.MerkleChallenge{},
						ApprovalId:                 "test-alice",
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						},
					},
				},
			},
			collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0],
		},
		TimelineTimes: GetFullUintRanges(),
	},
	}
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestFirstMatchOnlyWrongTime() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{{
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
			{
				ToMappingId:          "AllWithoutMint",
				InitiatedByMappingId: alice,
				TransferTimes:        []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				BadgeIds:             GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
				AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
					{
						IsApproved: false,
					},
				},
				ApprovalDetails: []*types.OutgoingApprovalDetails{
					{
						MerkleChallenges:                []*types.MerkleChallenge{},
						ApprovalId:                 "test-alice",
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						},
					},
				},
				
			},
			collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline[0].ApprovedOutgoingTransfers[0],
		},
		TimelineTimes: GetFullUintRanges(),
	},
	}
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestCombinations() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{{
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
			{
				ToMappingId:          "AllWithoutMint",
				InitiatedByMappingId: alice,
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
				BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
					{
						IsApproved: true,
					},
					{
						BadgeIdsOptions: &types.ValueOptions{ InvertDefault: true },
						IsApproved:      false,
					},
				},
				ApprovalDetails: []*types.OutgoingApprovalDetails{
					{
						MerkleChallenges:                []*types.MerkleChallenge{},
						ApprovalId:                 "test",
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						},
					},
				},
				
			},
		},
		TimelineTimes: GetFullUintRanges(),
	},
	}
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestCombinationsOrder() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{{
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
			{
				ToMappingId:          "AllWithoutMint",
				InitiatedByMappingId: alice,
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
				BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
					{
						IsApproved: false,
					},
					{
						InitiatedByMappingOptions: &types.ValueOptions{ InvertDefault: true },
						IsApproved:         true,
					},
				},
				ApprovalDetails: []*types.OutgoingApprovalDetails{
					{
						MerkleChallenges:                []*types.MerkleChallenge{},
						ApprovalId:                 "test",
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						},
					},
				},
				
			},
		},
		TimelineTimes: GetFullUintRanges(),
	},
	}
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestNotExplicitlyDefined() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	newOutgoingTimeline := []*types.UserApprovedOutgoingTransferTimeline{{
		ApprovedOutgoingTransfers: []*types.UserApprovedOutgoingTransfer{
			{
				ToMappingId:          "AllWithoutMint",
				InitiatedByMappingId: alice,
				TransferTimes:        GetFullUintRanges(),
				OwnershipTimes: GetFullUintRanges(),
				BadgeIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				AllowedCombinations: []*types.IsUserOutgoingTransferAllowed{
					{
						IsApproved: true,
					},
				},
				ApprovalDetails: []*types.OutgoingApprovalDetails{
					{
						MerkleChallenges:                []*types.MerkleChallenge{},
						ApprovalId:                 "test",
						MaxNumTransfers: &types.MaxNumTransfers{
							OverallMaxNumTransfers: sdkmath.NewUint(1000),
						},
						ApprovalAmounts: &types.ApprovalAmounts{
							PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						},
					},
				},
			},
		},
		TimelineTimes: GetFullUintRanges(),
	},
	}
	collectionsToCreate[0].DefaultApprovedOutgoingTransfersTimeline = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = suite.app.BadgesKeeper.DeductUserOutgoingApprovals(suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturned() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	// collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].OverridesFromApprovedOutgoingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(2, len(x), "Error deducting outgoing approvals")
	suite.Require().True(x[0].Outgoing != x[1].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesOutgoing() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].OverridesFromApprovedOutgoingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(1, len(x), "Error deducting outgoing approvals")
	suite.Require().False(x[0].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesIncoming() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].OverridesToApprovedIncomingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(1, len(x), "Error deducting outgoing approvals")
	suite.Require().True(x[0].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesBoth() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].OverridesToApprovedIncomingTransfers = true
	collectionsToCreate[0].CollectionApprovedTransfersTimeline[0].CollectionApprovedTransfers[0].ApprovalDetails[0].OverridesFromApprovedOutgoingTransfers = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(0, len(x), "Error deducting outgoing approvals")
}


//TODO: Test transfer tracker ID after update approved transfers
