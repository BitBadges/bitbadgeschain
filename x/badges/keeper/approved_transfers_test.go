package keeper_test

import (
	"math"
	"time"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// For legacy conversions
func DeductUserOutgoingApprovals(suite *TestSuite, ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string,
	requester string, amount sdkmath.Uint, solutions []*types.MerkleProof, prioritizedApprovals []*types.ApprovalIdentifierDetails,
	onlyCheckPrioritizedCollectionApprovals bool,
	onlyCheckProritizedIncomingApprovals bool,
	onlyCheckPrioritizedOutgoingApprovals bool,
	overrideTimestamp sdkmath.Uint,
	approvalsUsed *[]keeper.ApprovalsUsed,
	coinTransfersUsed *[]keeper.CoinTransfers) error {
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
		OverrideTimestamp:                       overrideTimestamp,
	}, from, to, requester, userBalance, approvalsUsed, coinTransfersUsed)
}

func DeductUserIncomingApprovals(suite *TestSuite, ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.BadgeCollection, userBalance *types.UserBalanceStore, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.MerkleProof, prioritizedApprovals []*types.ApprovalIdentifierDetails,
	onlyCheckPrioritizedCollectionApprovals bool,
	onlyCheckProritizedIncomingApprovals bool,
	onlyCheckPrioritizedOutgoingApprovals bool,
	overrideTimestamp sdkmath.Uint, approvalsUsed *[]keeper.ApprovalsUsed, coinTransfersUsed *[]keeper.CoinTransfers) error {
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
		OverrideTimestamp:                       overrideTimestamp,
	}, to, requester, userBalance, approvalsUsed, coinTransfersUsed)
}

func DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite *TestSuite, ctx sdk.Context, overallTransferBalances []*types.Balance,
	collection *types.BadgeCollection, badgeIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string,
	amount sdkmath.Uint, solutions []*types.MerkleProof, prioritizedApprovals []*types.ApprovalIdentifierDetails,
	onlyCheckPrioritizedCollectionApprovals bool,
	onlyCheckProritizedIncomingApprovals bool,
	onlyCheckPrioritizedOutgoingApprovals bool,
	overrideTimestamp sdkmath.Uint, approvalsUsed *[]keeper.ApprovalsUsed, coinTransfersUsed *[]keeper.CoinTransfers) ([]*keeper.UserApprovalsToCheck, error) {
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
			OverrideTimestamp:                       overrideTimestamp,
		}, to, requester, approvalsUsed, coinTransfersUsed)
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestDeductFromOutgoingTwoSeparateTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, GetCollectionsToCreate())
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

// Legacy helper
// Now that we switched to disallowing auto scan for truthy approval criteria, we need to specify the prioritized approvals manually
func GetDefaultPrioritizedApprovals(ctx sdk.Context, k keeper.Keeper, collectionId sdkmath.Uint) []*types.ApprovalIdentifierDetails {
	prioritizedApprovals := []*types.ApprovalIdentifierDetails{
		{
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      "mint-test",
			Version:         sdkmath.NewUint(0),
		},
		{
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      GetBobApproval().ApprovalId,
			Version:         sdkmath.NewUint(0),
		},
		{
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      "test",
			Version:         sdkmath.NewUint(0),
		},
		{
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      "tessdgfst",
			Version:         sdkmath.NewUint(0),
		},
		{
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      "test2",
			Version:         sdkmath.NewUint(0),
		},
	}

	otherIds := []string{
		"asadsdas", // most common
		"fasdfasdf",
		"target approval",
		"random approval",
		"asadsdasfghdsfasdfasdf",
		"asadsdasfghaaadsd",
		"asadsdasfghd",
		"testsdfgsdf",
		"testsgdfs",
	}

	for _, otherId := range otherIds {
		prioritizedApprovals = append(prioritizedApprovals, &types.ApprovalIdentifierDetails{
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      otherId,
			Version:         sdkmath.NewUint(0),
		})
	}

	addresses := []string{bob, alice, charlie}

	for _, address := range addresses {
		prioritizedApprovals = append(prioritizedApprovals, &types.ApprovalIdentifierDetails{
			ApprovalLevel:   "incoming",
			ApproverAddress: address,
			ApprovalId:      "test",
			Version:         sdkmath.NewUint(0),
		})

		prioritizedApprovals = append(prioritizedApprovals, &types.ApprovalIdentifierDetails{
			ApprovalLevel:   "outgoing",
			ApproverAddress: address,
			ApprovalId:      "test",
			Version:         sdkmath.NewUint(0),
		})
	}

	for _, approval := range prioritizedApprovals {
		version, found := k.GetApprovalTrackerVersionFromStore(ctx, keeper.ConstructApprovalVersionKey(collectionId, approval.ApprovalLevel, approval.ApproverAddress, approval.ApprovalId))
		if found {
			approval.Version = version
		}
	}

	return prioritizedApprovals
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
		DurationFromTimestamp:     sdkmath.NewUint(0),
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
		DurationFromTimestamp:     sdkmath.NewUint(0),
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
// 	suite.Require().Error(err, "Error deducting outgoing approvals")

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
// 	suite.Require().Nil(err, "Error deducting outgoing approvals")

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturned() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	// collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
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
// 		sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
// 	suite.Require().Error(err, "Error deducting outgoing approvals")
// }

const ProtocolFee = 100000000

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
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000-100-ProtocolFee), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000+100), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")

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
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

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
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")
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
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000-200-ProtocolFee), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000+100), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000+100), suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge").Amount, "Error deducting outgoing approvals")
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
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000-300-(ProtocolFee*2)), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000+200), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000+100), suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge").Amount, "Error deducting outgoing approvals")
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
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000-300-(ProtocolFee*2)), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000+200), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000+100), suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge").Amount, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestVersionControlCollectionApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	version, found := suite.app.BadgesKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "collection", "", "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(0), version, "Error getting approval tracker version")

	// Update the collection approvals
	collectionsToCreate[0].CollectionApprovals[0].BadgeIds = GetTwoUintRanges()
	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		CollectionId:        sdkmath.NewUint(1),
		Creator:             bob,
		CollectionApprovals: collectionsToCreate[0].CollectionApprovals,
	})
	suite.Require().Nil(err, "error updating collection approvals")

	version, found = suite.app.BadgesKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "collection", "", "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(1), version, "Error getting approval tracker version")

	storedApprovals := collectionsToCreate[0].CollectionApprovals

	// Should persist version even after setting empty and resetting
	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{}
	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		CollectionId:        sdkmath.NewUint(1),
		Creator:             bob,
		CollectionApprovals: []*types.CollectionApproval{},
	})
	suite.Require().Nil(err, "error updating collection approvals")

	version, found = suite.app.BadgesKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "collection", "", "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(1), version, "Error getting approval tracker version")

	// Update it back and check incremented
	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		CollectionId:        sdkmath.NewUint(1),
		Creator:             bob,
		CollectionApprovals: storedApprovals,
	})
	suite.Require().Nil(err, "error updating collection approvals")

	version, found = suite.app.BadgesKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "collection", "", "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(2), version, "Error getting approval tracker version")
}

func (suite *TestSuite) TestVersionControlUserApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	defaultApproval := &types.UserIncomingApproval{
		FromListId:        "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetFullUintRanges(),

		ApprovalId: "test",
		ApprovalCriteria: &types.IncomingApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1),
			},
		},
	}

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateIncomingApprovals: true,
		IncomingApprovals: []*types.UserIncomingApproval{
			defaultApproval,
		},
	})
	suite.Require().Nil(err, "error updating user approvals")

	version, found := suite.app.BadgesKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "incoming", bob, "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(0), version, "Error getting approval tracker version")

	defaultApproval.BadgeIds = GetOneUintRange()
	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{defaultApproval},
	})
	suite.Require().Nil(err, "error updating user approvals")

	version, found = suite.app.BadgesKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "incoming", bob, "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(1), version, "Error getting approval tracker version")

	// Should persist version even after setting empty and resetting
	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{},
	})
	suite.Require().Nil(err, "error updating user approvals")

	version, found = suite.app.BadgesKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "incoming", bob, "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(1), version, "Error getting approval tracker version")

	// Update it back and check incremented
	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{defaultApproval},
	})
	suite.Require().Nil(err, "error updating user approvals")

	version, found = suite.app.BadgesKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "incoming", bob, "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(2), version, "Error getting approval tracker version")
}

func (suite *TestSuite) TestMaxOneTransferWithResetIntervals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(10)
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(0)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(0)

	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals = &types.ResetTimeIntervals{
		StartTime:      sdkmath.NewUint(1000),
		IntervalLength: sdkmath.NewUint(100),
	}

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1000))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1099))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1100))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1200))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestResetIntervalsWithFutureTime() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(2)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(10)
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(0)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(0)

	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.MaxNumTransfers.ResetTimeIntervals = &types.ResetTimeIntervals{
		StartTime:      sdkmath.NewUint(10000),
		IntervalLength: sdkmath.NewUint(100),
	}

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1000))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(20000))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)), false, false, false, sdkmath.NewUint(0), &[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestAutoDeletingOutgoingApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(10)
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(0)
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.AutoDeletionOptions = &types.AutoDeletionOptions{
		AfterOneUse: true,
	}

	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		{
			FromListId:        "Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			BadgeIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "a1",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
			},
		},
		{
			FromListId:        "!Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			BadgeIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "a2",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesToIncomingApprovals: true,
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	bobBalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Equal(1, len(bobBalance.OutgoingApprovals), "Outgoing approvals should not be deleted")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	for i := 0; i < 2; i++ {
		err = suite.app.BadgesKeeper.HandleTransfers(
			suite.ctx,
			collection,
			[]*types.Transfer{
				{
					From:        bob,
					ToAddresses: []string{alice},
					Balances: []*types.Balance{
						{
							BadgeIds:       GetFullUintRanges(),
							Amount:         sdkmath.NewUint(1),
							OwnershipTimes: GetFullUintRanges(),
						},
					},
					PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
						{
							ApprovalId:      collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalId,
							ApprovalLevel:   "outgoing",
							ApproverAddress: bob,
							Version:         sdkmath.NewUint(0),
						},
					},
					OnlyCheckPrioritizedOutgoingApprovals: true,
				},
			},
			alice,
		)
		if i == 0 {
			suite.Require().Nil(err, "Error deducting outgoing approvals")

			//Transfer it back to bob
			err = suite.app.BadgesKeeper.HandleTransfers(
				suite.ctx,
				collection,
				[]*types.Transfer{
					{
						From:        alice,
						ToAddresses: []string{bob},
						Balances: []*types.Balance{
							{
								Amount:         sdkmath.NewUint(1),
								BadgeIds:       GetFullUintRanges(),
								OwnershipTimes: GetFullUintRanges(),
							},
						},
						PrioritizedApprovals: []*types.ApprovalIdentifierDetails{},
					},
				},
				alice,
			)
			suite.Require().Nil(err, "Error deducting outgoing approvals")

		} else {
			suite.Require().Error(err, "Error deducting outgoing approvals")
		}
	}

	bobBalance, _ = GetUserBalance(suite, wctx, collection.CollectionId, bob)
	suite.Require().Equal(0, len(bobBalance.OutgoingApprovals), "Outgoing approvals should be deleted")
}

func (suite *TestSuite) TestAutoDeletingIncomingApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(10)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(0)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.AutoDeletionOptions = &types.AutoDeletionOptions{
		AfterOneUse: true,
	}

	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		{
			FromListId:        "Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			BadgeIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "a1",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
			},
		},
		{
			FromListId:        "!Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			BadgeIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "a2",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	aliceBalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Equal(1, len(aliceBalance.IncomingApprovals), "Incoming approvals should not be deleted")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	for i := 0; i < 2; i++ {

		transfers := []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						BadgeIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetBottomHalfUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalId,
						ApprovalLevel:   "incoming",
						ApproverAddress: alice,
						Version:         sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedIncomingApprovals: true,
			},
		}
		if i == 1 {
			transfers[0].Balances[0].OwnershipTimes = GetTopHalfUintRanges()
		}
		err = suite.app.BadgesKeeper.HandleTransfers(
			suite.ctx,
			collection,
			transfers,
			alice,
		)
		if i == 0 {
			suite.Require().Nil(err, "Error deducting outgoing approvals")

		} else {
			suite.Require().Error(err, "Error deducting outgoing approvals")
		}
	}

	aliceBalance, _ = GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(0, len(aliceBalance.IncomingApprovals), "Incoming approvals should be deleted")
}

func (suite *TestSuite) TestAutoDeletingCollectionApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(10)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(0)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.AutoDeletionOptions = &types.AutoDeletionOptions{
		AfterOneUse: true,
	}

	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		{
			FromListId:        "Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			BadgeIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "a1",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				AutoDeletionOptions: &types.AutoDeletionOptions{
					AfterOneUse: true,
				},
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating badges")

	// We mint to bob in the CreateCollections.transfers so it should be deleted automatically
	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Equal(0, len(collection.CollectionApprovals), "Collection approvals should be deleted")
}
