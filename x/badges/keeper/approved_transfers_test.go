package keeper_test

import (
	"fmt"
	"math"

	"bitbadgeschain/x/badges/keeper"
	"bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	types1 "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// For legacy conversions
func DeductUserOutgoingApprovals(suite *TestSuite, ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string,
	requester string, amount sdkmath.Uint, solutions []*types.MerkleProof, prioritizedApprovals []*types.ApprovalIdentifierDetails,
	onlyCheckPrioritizedCollectionApprovals bool,
	onlyCheckProritizedIncomingApprovals bool,
	onlyCheckPrioritizedOutgoingApprovals bool,
	zkProofSolutions []*types.ZkProofSolution) error {
	return suite.app.BadgesKeeper.DeductUserOutgoingApprovals(ctx, collection, overallTransferBalances, &types.Transfer{
		From:        from,
		ToAddresses: []string{to},
		Balances: []*types.Balance{
			{
				BadgeIds:       badgeIds,
				OwnershipTimes: times,
				Amount:         amount,
			},
		},
		MerkleProofs:                            solutions,
		PrioritizedApprovals:                    prioritizedApprovals,
		OnlyCheckPrioritizedCollectionApprovals: onlyCheckPrioritizedCollectionApprovals,
		OnlyCheckPrioritizedIncomingApprovals:   onlyCheckProritizedIncomingApprovals,
		OnlyCheckPrioritizedOutgoingApprovals:   onlyCheckPrioritizedOutgoingApprovals,
		ZkProofSolutions:                        zkProofSolutions,
	}, from, to, requester, userBalance)
}

func DeductUserIncomingApprovals(suite *TestSuite, ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.MerkleProof, prioritizedApprovals []*types.ApprovalIdentifierDetails,
	onlyCheckPrioritizedCollectionApprovals bool,
	onlyCheckProritizedIncomingApprovals bool,
	onlyCheckPrioritizedOutgoingApprovals bool,
	zkProofSolutions []*types.ZkProofSolution) error {
	return suite.app.BadgesKeeper.DeductUserIncomingApprovals(ctx, collection, overallTransferBalances, &types.Transfer{
		From:        from,
		ToAddresses: []string{to},
		Balances: []*types.Balance{
			{
				BadgeIds:       badgeIds,
				OwnershipTimes: times,
				Amount:         amount,
			},
		},
		MerkleProofs:                            solutions,
		PrioritizedApprovals:                    prioritizedApprovals,
		OnlyCheckPrioritizedCollectionApprovals: onlyCheckPrioritizedCollectionApprovals,
		OnlyCheckPrioritizedIncomingApprovals:   onlyCheckProritizedIncomingApprovals,
		OnlyCheckPrioritizedOutgoingApprovals:   onlyCheckPrioritizedOutgoingApprovals,
		ZkProofSolutions:                        zkProofSolutions,
	}, to, requester, userBalance)
}

func DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite *TestSuite, ctx sdk.Context, overallTransferBalances []*types.Balance,
	collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string,
	amount sdkmath.Uint, solutions []*types.MerkleProof, prioritizedApprovals []*types.ApprovalIdentifierDetails,
	onlyCheckPrioritizedCollectionApprovals bool,
	onlyCheckProritizedIncomingApprovals bool,
	onlyCheckPrioritizedOutgoingApprovals bool,
	zkProofSolutions []*types.ZkProofSolution) ([]*keeper.UserApprovalsToCheck, error) {
	return suite.app.BadgesKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx, collection,
		&types.Transfer{
			From:        from,
			ToAddresses: []string{to},
			Balances: []*types.Balance{
				{
					BadgeIds:       badgeIds,
					OwnershipTimes: times,
					Amount:         amount,
				},
			},
			MerkleProofs:                            solutions,
			PrioritizedApprovals:                    prioritizedApprovals,
			OnlyCheckPrioritizedCollectionApprovals: onlyCheckPrioritizedCollectionApprovals,
			OnlyCheckPrioritizedIncomingApprovals:   onlyCheckProritizedIncomingApprovals,
			OnlyCheckPrioritizedOutgoingApprovals:   onlyCheckPrioritizedOutgoingApprovals,
			ZkProofSolutions:                        zkProofSolutions,
		}, to, requester)
}

func (suite *TestSuite) TestMsgsToExecute() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// msg := &bankv1beta1.MsgSend{
	// 	FromAddress: alice,
	// 	ToAddress:   bob,
	// 	Amount: []*basev1beta1.Coin{
	// 		{
	// 			Denom:  "ubadge",
	// 			Amount: "1",
	// 		},
	// 	},
	// }

	msg := &types.MsgDeleteCollection{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
	}

	protoMsg, err := types1.NewAnyWithValue(msg)
	suite.Require().Nil(err, "error creating any")
	suite.Require().NotNil(protoMsg, "error creating any")

	// fmt.Println(protoMsg)
	fmt.Println(protoMsg.GetCachedValue())

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.MsgsToExecute = []*types.GenericMsgToExecute{
		{
			UseApproverAddressAsCreator: true,
			MsgsToExecute:               protoMsg,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	// suite.Require().Equal(2, len(x), "Error deducting outgoing approvals")
	// suite.Require().True(x[0].Outgoing != x[1].Outgoing, "Error deducting outgoing approvals")
}

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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestDeductFromOutgoingTwoSeparateTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, GetCollectionsToCreate())
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
// 	suite.Require().Error(err, "Error deducting outgoing approvals")

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
// 	suite.Require().Nil(err, "Error deducting outgoing approvals")

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

			ApprovalId: "test",
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturned() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	// collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
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

// 	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite,
// 		suite.ctx,
// 		[]*types.Balance{},
// 		collection,
// 		GetFullUintRanges(),
// 		GetFullUintRanges(),
// 		bob, alice, alice,
// 		sdkmath.NewUint(1), []*types.MerkleProof{}, nil, false, false, false, []*types.ZkProofSolution{})
// 	suite.Require().Error(err, "Error deducting outgoing approvals")
// }

func (suite *TestSuite) TestCoinTransfersWithApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{
		{
			To: alice,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100),
					Denom:  "ubadge",
				},
			},
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(1000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(0), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						BadgeIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(900), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")

}

func (suite *TestSuite) TestCoinTransfersWithApprovalsUnderflow() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{
		{
			To: alice,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100000),
					Denom:  "ubadge",
				},
			},
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(1000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(0), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						BadgeIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(1000), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(0), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestCoinTransfersWithApprovalsMultiple() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{
		{
			To: alice,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100),
					Denom:  "ubadge",
				},
			},
		},
		{
			To: charlie,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100),
					Denom:  "ubadge",
				},
			},
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(1000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(0), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						BadgeIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(800), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100), suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge").Amount, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestCoinTransfersWithOverflowIntoNextApprovals() {
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{
		{
			To: alice,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100),
					Denom:  "ubadge",
				},
			},
		},
		{
			To: charlie,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100),
					Denom:  "ubadge",
				},
			},
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.MaxNumTransfers.OverallMaxNumTransfers = sdkmath.NewUint(1)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetFullUintRanges(),
		ApprovalId:        "test2",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1),
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	})
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{
		{
			To: alice,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100),
					Denom:  "ubadge",
				},
			},
		},
	}

	wctx := sdk.WrapSDKContext(suite.ctx)
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(1000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(0), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						BadgeIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						BadgeIds:       GetTwoUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(700), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(200), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100), suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge").Amount, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestWeirdBootstrapThing() {
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{
		{
			To: alice,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100),
					Denom:  "ubadge",
				},
			},
		},
		{
			To: charlie,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100),
					Denom:  "ubadge",
				},
			},
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.MaxNumTransfers.OverallMaxNumTransfers = sdkmath.NewUint(1)
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ToListId:          "AllWithoutMint",
		FromListId:        "Mint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetFullUintRanges(),
		ApprovalId:        "test2",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1),
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
	})
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{
		{
			To: alice,
			Coins: []*sdk.Coin{
				{
					Amount: sdkmath.NewInt(100),
					Denom:  "ubadge",
				},
			},
		},
	}

	wctx := sdk.WrapSDKContext(suite.ctx)
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(1000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(0), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						BadgeIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						BadgeIds:       GetTwoUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(700), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(200), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100), suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge").Amount, "Error deducting outgoing approvals")
}
