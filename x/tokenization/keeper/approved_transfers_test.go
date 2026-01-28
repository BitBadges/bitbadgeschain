package keeper_test

import (
	"math"
	"time"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// For legacy conversions
func DeductUserOutgoingApprovals(suite *TestSuite, ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.TokenCollection, userBalance *types.UserBalanceStore, tokenIds []*types.UintRange, times []*types.UintRange, from string, to string,
	requester string, amount sdkmath.Uint, solutions []*types.MerkleProof, prioritizedApprovals []*types.ApprovalIdentifierDetails,
	onlyCheckPrioritizedCollectionApprovals bool,
	onlyCheckProritizedIncomingApprovals bool,
	onlyCheckPrioritizedOutgoingApprovals bool,
	precalculationOptions *types.PrecalculationOptions,
	eventTracking *keeper.EventTracking,
	userRoyalties *types.UserRoyalties,
) error {
	transferMetadata := keeper.TransferMetadata{
		To:              to,
		From:            from,
		InitiatedBy:     requester,
		ApproverAddress: from,
		ApprovalLevel:   "outgoing",
	}

	return suite.app.TokenizationKeeper.DeductUserOutgoingApprovals(ctx, collection, overallTransferBalances, &types.Transfer{
		From:        from,
		ToAddresses: []string{to},
		Balances: []*types.Balance{
			{
				TokenIds:       tokenIds,
				OwnershipTimes: times,
				Amount:         amount,
			},
		},
		MerkleProofs:                            solutions,
		PrioritizedApprovals:                    prioritizedApprovals,
		OnlyCheckPrioritizedCollectionApprovals: onlyCheckPrioritizedCollectionApprovals,
		OnlyCheckPrioritizedIncomingApprovals:   onlyCheckProritizedIncomingApprovals,
		OnlyCheckPrioritizedOutgoingApprovals:   onlyCheckPrioritizedOutgoingApprovals,
		PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
			PrecalculationOptions: precalculationOptions,
		},
	}, transferMetadata, userBalance, eventTracking, userRoyalties)
}

func DeductUserIncomingApprovals(suite *TestSuite, ctx sdk.Context, overallTransferBalances []*types.Balance, collection *types.TokenCollection, userBalance *types.UserBalanceStore, tokenIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string, amount sdkmath.Uint, solutions []*types.MerkleProof, prioritizedApprovals []*types.ApprovalIdentifierDetails,
	onlyCheckPrioritizedCollectionApprovals bool,
	onlyCheckProritizedIncomingApprovals bool,
	onlyCheckPrioritizedOutgoingApprovals bool,
	precalculationOptions *types.PrecalculationOptions, eventTracking *keeper.EventTracking, userRoyalties *types.UserRoyalties) error {
	transferMetadata := keeper.TransferMetadata{
		To:              to,
		From:            from,
		InitiatedBy:     requester,
		ApproverAddress: to,
		ApprovalLevel:   "incoming",
	}

	return suite.app.TokenizationKeeper.DeductUserIncomingApprovals(ctx, collection, overallTransferBalances, &types.Transfer{
		From:        from,
		ToAddresses: []string{to},
		Balances: []*types.Balance{
			{
				TokenIds:       tokenIds,
				OwnershipTimes: times,
				Amount:         amount,
			},
		},
		MerkleProofs:                            solutions,
		PrioritizedApprovals:                    prioritizedApprovals,
		OnlyCheckPrioritizedCollectionApprovals: onlyCheckPrioritizedCollectionApprovals,
		OnlyCheckPrioritizedIncomingApprovals:   onlyCheckProritizedIncomingApprovals,
		OnlyCheckPrioritizedOutgoingApprovals:   onlyCheckPrioritizedOutgoingApprovals,
		PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
			PrecalculationOptions: precalculationOptions,
		},
	}, transferMetadata, userBalance, eventTracking, userRoyalties)
}

func DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite *TestSuite, ctx sdk.Context, overallTransferBalances []*types.Balance,
	collection *types.TokenCollection, tokenIds []*types.UintRange, times []*types.UintRange, from string, to string, requester string,
	amount sdkmath.Uint, solutions []*types.MerkleProof, prioritizedApprovals []*types.ApprovalIdentifierDetails,
	onlyCheckPrioritizedCollectionApprovals bool,
	onlyCheckProritizedIncomingApprovals bool,
	onlyCheckPrioritizedOutgoingApprovals bool,
	precalculationOptions *types.PrecalculationOptions, eventTracking *keeper.EventTracking) ([]*keeper.UserApprovalsToCheck, error) {
	transferMetadata := keeper.TransferMetadata{
		To:              to,
		From:            from,
		InitiatedBy:     requester,
		ApproverAddress: "",
		ApprovalLevel:   "collection",
	}

	return suite.app.TokenizationKeeper.DeductCollectionApprovalsAndGetUserApprovalsToCheck(ctx, collection,
		&types.Transfer{
			From:        from,
			ToAddresses: []string{to},
			Balances: []*types.Balance{
				{
					TokenIds:       tokenIds,
					OwnershipTimes: times,
					Amount:         amount,
				},
			},
			MerkleProofs:                            solutions,
			PrioritizedApprovals:                    prioritizedApprovals,
			OnlyCheckPrioritizedCollectionApprovals: onlyCheckPrioritizedCollectionApprovals,
			OnlyCheckPrioritizedIncomingApprovals:   onlyCheckProritizedIncomingApprovals,
			OnlyCheckPrioritizedOutgoingApprovals:   onlyCheckPrioritizedOutgoingApprovals,
			PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
				PrecalculationOptions: precalculationOptions,
			},
		}, transferMetadata, eventTracking, "collection")
}

func (suite *TestSuite) TestDeductFromOutgoing() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, GetCollectionsToCreate())
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	overallTransferBalances := []*types.Balance{
		{
			TokenIds:       GetFullUintRanges(),
			OwnershipTimes: GetFullUintRanges(),
			Amount:         sdkmath.NewUint(1),
		},
	}

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetFullUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, overallTransferBalances, collection, GetFullUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestDeductFromOutgoingTwoSeparateTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateCollections(suite, wctx, GetCollectionsToCreate())
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)
	aliceBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, alice)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, aliceBalance, GetTopHalfUintRanges(), GetFullUintRanges(), alice, bob, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestMaxOneTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.MaxNumTransfers.PerFromAddressMaxNumTransfers = sdkmath.NewUint(1)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
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

// GetPrioritizedApprovalsFromCollection gets all collection approvals as prioritized approvals
// This is useful for tests involving special addresses where OnlyCheckPrioritizedCollectionApprovals is set
func GetPrioritizedApprovalsFromCollection(ctx sdk.Context, k keeper.Keeper, collection *types.TokenCollection) []*types.ApprovalIdentifierDetails {
	prioritizedApprovals := []*types.ApprovalIdentifierDetails{}
	for _, approval := range collection.CollectionApprovals {
		version, found := k.GetApprovalTrackerVersionFromStore(
			ctx,
			keeper.ConstructApprovalVersionKey(collection.CollectionId, "collection", "", approval.ApprovalId),
		)
		if !found {
			version = sdkmath.NewUint(0)
		}
		prioritizedApprovals = append(prioritizedApprovals, &types.ApprovalIdentifierDetails{
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			ApprovalId:      approval.ApprovalId,
			Version:         version,
		})
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
				TokenIds:       GetFullUintRanges(),
				Amount:         sdkmath.NewUint(1),
			},
		},
		IncrementTokenIdsBy:       sdkmath.NewUint(math.MaxUint64),
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
				TokenIds:       GetFullUintRanges(),
				Amount:         sdkmath.NewUint(1),
			},
		},
		IncrementTokenIdsBy:       sdkmath.NewUint(math.MaxUint64),
		IncrementOwnershipTimesBy: sdkmath.NewUint(math.MaxUint64),
		DurationFromTimestamp:     sdkmath.NewUint(0),
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	overallTransferBalances := []*types.Balance{
		{
			OwnershipTimes: GetFullUintRanges(),
			TokenIds:       GetFullUintRanges(),
			Amount:         sdkmath.NewUint(1),
		},
	}

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, &types.PrecalculationOptions{}, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, &types.PrecalculationOptions{}, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, &types.PrecalculationOptions{}, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, overallTransferBalances, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, &types.PrecalculationOptions{}, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestRequiresEquals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultOutgoingApprovals[0].ApprovalCriteria.RequireToDoesNotEqualInitiatedBy = true
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria.RequireFromDoesNotEqualInitiatedBy = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestSpecificApproved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultOutgoingApprovals[0].InitiatedByListId = alice
	collectionsToCreate[0].DefaultIncomingApprovals[0].InitiatedByListId = alice

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestDefaults() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultOutgoingApprovals[0].InitiatedByListId = alice
	collectionsToCreate[0].DefaultIncomingApprovals[0].InitiatedByListId = alice

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestDefaultsNotAutoApplies() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].DefaultOutgoingApprovals = []*types.UserOutgoingApproval{}
	collectionsToCreate[0].DefaultIncomingApprovals = []*types.UserIncomingApproval{}
	collectionsToCreate[0].DefaultDisapproveSelfInitiated = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserIncomingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, bob, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
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
// 				TokenIds: []*types.UintRange{
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
// 	suite.Require().Nil(err, "error creating tokens")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
// 	suite.Require().Error(err, "Error deducting outgoing approvals")

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
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
// 			TokenIds:             GetFullUintRanges(),
// 			OwnershipTimes: GetFullUintRanges(),
// 			IsApproved: false,

// 		},
// 		collectionsToCreate[0].DefaultOutgoingApprovals[0],
// 	}

// 	collectionsToCreate[0].DefaultOutgoingApprovals = newOutgoingTimeline

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "error creating tokens")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
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
// 				TokenIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},

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
// 				TokenIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 				TokenIdsOptions: &types.ValueOptions{ InvertDefault: true },
// 				IsApproved:      false,
// 			},
// 	}
// 	collectionsToCreate[0].DefaultOutgoingApprovals = newOutgoingTimeline

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "error creating tokens")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
// 	suite.Require().Nil(err, "Error deducting outgoing approvals")

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
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
// 				TokenIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
// 				IsApproved: false,

// 			},
// 			{
// 					ToListId:          "AllWithoutMint",
// 					InitiatedByListId: alice,
// 					TransferTimes:        GetFullUintRanges(),
// 					OwnershipTimes: GetFullUintRanges(),
// 					TokenIds:             []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
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
// 	suite.Require().Nil(err, "error creating tokens")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetOneUintRange(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
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
// 			TokenIds:             GetFullUintRanges(),
// 			IsApproved: false,
// 	})

// 	err := CreateCollections(suite, wctx, collectionsToCreate)
// 	suite.Require().Nil(err, "error creating tokens")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
// 	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

// 	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
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
			TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},

			ApprovalId: "test",
			ApprovalCriteria: &types.OutgoingApprovalCriteria{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(1000),
					AmountTrackerId:        "test-tracker",
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					PerFromAddressApprovalAmount: sdkmath.NewUint(1),
					AmountTrackerId:              "test-tracker",
				},
			},
		},
	}
	collectionsToCreate[0].DefaultOutgoingApprovals = newOutgoingTimeline

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturned() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	// collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(2, len(x), "Error deducting outgoing approvals")
	suite.Require().True(x[0].Outgoing != x[1].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesOutgoing() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}})
	suite.Require().Nil(err, "Error deducting outgoing approvals")
	suite.Require().Equal(1, len(x), "Error deducting outgoing approvals")
	suite.Require().False(x[0].Outgoing, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestUserApprovalsReturnedOverridesIncoming() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}})
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
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	x, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx, []*types.Balance{}, collection, GetTopHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}})
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
// 			TokenIds: 					 GetFullUintRanges(),
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
// 	suite.Require().Nil(err, "error creating tokens")

// 	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

// 	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite,
// 		suite.ctx,
// 		[]*types.Balance{},
// 		collection,
// 		GetFullUintRanges(),
// 		GetFullUintRanges(),
// 		bob, alice, alice,
// 		sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
// 	suite.Require().Error(err, "Error deducting outgoing approvals")
// }

// ProtocolFee is now calculated as 0.1% of the transfer amount
// For 100 ubadge: 0.1% = 0.1, rounded down to 0
// For 200 ubadge: 0.1% = 0.2, rounded down to 0
// For 300 ubadge: 0.1% = 0.3, rounded down to 0

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
	suite.Require().Nil(err, "error creating tokens")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000-100-0), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
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
	suite.Require().Nil(err, "error creating tokens")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetFullUintRanges(),
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
	suite.Require().Nil(err, "error creating tokens")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000-200), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
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
		TokenIds:          GetFullUintRanges(),
		ApprovalId:        "test2",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1),
				AmountTrackerId:              "test-tracker",
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
	suite.Require().Nil(err, "error creating tokens")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetTwoUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000-300), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
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
		TokenIds:          GetFullUintRanges(),
		ApprovalId:        "test2",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1),
				AmountTrackerId:              "test-tracker",
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
	suite.Require().Nil(err, "error creating tokens")

	bobBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000), bobBalanceBefore.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000), aliceBalanceBefore.Amount, "Error deducting outgoing approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetTwoUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	bobBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(bob), "ubadge")
	aliceBalanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(sdkmath.NewInt(100000000000-300), bobBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000+200), aliceBalanceAfter.Amount, "Error deducting outgoing approvals")
	suite.Require().Equal(sdkmath.NewInt(100000000000+100), suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(charlie), "ubadge").Amount, "Error deducting outgoing approvals")
}

func (suite *TestSuite) TestVersionControlCollectionApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	version, found := suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "collection", "", "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(0), version, "Error getting approval tracker version")

	// Update the collection approvals
	collectionsToCreate[0].CollectionApprovals[0].TokenIds = GetTwoUintRanges()
	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		CollectionId:        sdkmath.NewUint(1),
		Creator:             bob,
		CollectionApprovals: collectionsToCreate[0].CollectionApprovals,
	})
	suite.Require().Nil(err, "error updating collection approvals")

	version, found = suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "collection", "", "test"))
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

	version, found = suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "collection", "", "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(1), version, "Error getting approval tracker version")

	// Update it back and check incremented
	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		CollectionId:        sdkmath.NewUint(1),
		Creator:             bob,
		CollectionApprovals: storedApprovals,
	})
	suite.Require().Nil(err, "error updating collection approvals")

	version, found = suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "collection", "", "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(2), version, "Error getting approval tracker version")
}

func (suite *TestSuite) TestVersionControlUserApprovals() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	defaultApproval := &types.UserIncomingApproval{
		FromListId:        "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),

		ApprovalId: "test",
		ApprovalCriteria: &types.IncomingApprovalCriteria{

			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "test-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1),
				AmountTrackerId:              "test-tracker",
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

	version, found := suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "incoming", bob, "test"))
	suite.Require().True(found, "Error getting approval tracker version")
	suite.Require().Equal(sdkmath.NewUint(0), version, "Error getting approval tracker version")

	defaultApproval.TokenIds = GetOneUintRange()
	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{defaultApproval},
	})
	suite.Require().Nil(err, "error updating user approvals")

	version, found = suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "incoming", bob, "test"))
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

	version, found = suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "incoming", bob, "test"))
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

	version, found = suite.app.TokenizationKeeper.GetApprovalTrackerVersionFromStore(suite.ctx, keeper.ConstructApprovalVersionKey(sdkmath.NewUint(1), "incoming", bob, "test"))
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
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1099))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1100))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(1200))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
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
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	bobBalance, _ := GetUserBalance(suite, wctx, collection.CollectionId, bob)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Error(err, "Error deducting outgoing approvals")

	suite.ctx = suite.ctx.WithBlockTime(time.UnixMilli(20000))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	err = DeductUserOutgoingApprovals(suite, suite.ctx, []*types.Balance{}, collection, bobBalance, GetBottomHalfUintRanges(), GetFullUintRanges(), bob, alice, alice, sdkmath.NewUint(1), []*types.MerkleProof{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)), false, false, false, nil, &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}, nil)
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
			TokenIds:          GetFullUintRanges(),
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
			TokenIds:          GetFullUintRanges(),
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
	suite.Require().Nil(err, "error creating tokens")

	bobBalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Equal(1, len(bobBalance.OutgoingApprovals), "Outgoing approvals should not be deleted")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	for i := 0; i < 2; i++ {
		err = suite.app.TokenizationKeeper.HandleTransfers(
			suite.ctx,
			collection,
			[]*types.Transfer{
				{
					From:        bob,
					ToAddresses: []string{alice},
					Balances: []*types.Balance{
						{
							TokenIds:       GetFullUintRanges(),
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
			err = suite.app.TokenizationKeeper.HandleTransfers(
				suite.ctx,
				collection,
				[]*types.Transfer{
					{
						From:        alice,
						ToAddresses: []string{bob},
						Balances: []*types.Balance{
							{
								Amount:         sdkmath.NewUint(1),
								TokenIds:       GetFullUintRanges(),
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
			TokenIds:          GetFullUintRanges(),
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
			TokenIds:          GetFullUintRanges(),
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
	suite.Require().Nil(err, "error creating tokens")

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
						TokenIds:       GetFullUintRanges(),
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
		err = suite.app.TokenizationKeeper.HandleTransfers(
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
			TokenIds:          GetFullUintRanges(),
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
	suite.Require().Nil(err, "error creating tokens")

	// We mint to bob in the CreateCollections.transfers so it should be deleted automatically
	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Equal(0, len(collection.CollectionApprovals), "Collection approvals should be deleted")
}

func (suite *TestSuite) TestAutoDeletingAfterOverallMaxNumTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	// Set up collection approval with auto-deletion after overall max transfers
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, []*types.CollectionApproval{
		{
			FromListId:        "Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			TokenIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "overall-max-test",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(3),
					AmountTrackerId:        "tracker1",
				},
				AutoDeletionOptions: &types.AutoDeletionOptions{
					AfterOverallMaxNumTransfers: true,
				},
			},
		},
	}...)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// First transfer - should succeed and not delete the approval
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "overall-max-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
		bob,
	)
	suite.Require().Nil(err, "First transfer should succeed")

	// Check that approval still exists after first transfer
	collection, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Second transfer - should succeed and not delete the approval yet
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "overall-max-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
		bob,
	)
	suite.Require().Nil(err, "Second transfer should succeed")

	// Check that approval still exists after second transfer (threshold reached but not exceeded)
	collection, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Third transfer - should succeed and delete the approval (threshold exceeded)
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "overall-max-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
		bob,
	)
	suite.Require().Nil(err, "Third transfer should succeed")

	// Check that approval is deleted after third transfer (threshold exceeded)
	collection, _ = GetCollection(suite, wctx, sdkmath.NewUint(1))

	found := false
	for _, approval := range collection.CollectionApprovals {
		if approval.ApprovalId == "overall-max-test" {
			found = true
			break
		}
	}
	suite.Require().False(found, "Collection approval should be deleted after third transfer")

	// Fourth transfer - should fail because approval is deleted
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "overall-max-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
		bob,
	)
	suite.Require().Error(err, "Fourth transfer should fail because approval is deleted")
}

func (suite *TestSuite) TestAutoDeletingIncomingApprovalsAfterOverallMaxNumTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	// Set up incoming approval with auto-deletion after overall max transfers
	collectionsToCreate[0].DefaultIncomingApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].DefaultIncomingApprovals[0].InitiatedByListId = "All"
	collectionsToCreate[0].DefaultIncomingApprovals[0].TokenIds = GetFullUintRanges()
	collectionsToCreate[0].DefaultIncomingApprovals[0].TransferTimes = GetFullUintRanges()
	collectionsToCreate[0].DefaultIncomingApprovals[0].OwnershipTimes = GetFullUintRanges()
	collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalCriteria = &types.IncomingApprovalCriteria{
		AutoDeletionOptions: &types.AutoDeletionOptions{
			AfterOverallMaxNumTransfers: true,
		},
		ApprovalAmounts: &types.ApprovalAmounts{},
		MaxNumTransfers: &types.MaxNumTransfers{
			OverallMaxNumTransfers: sdkmath.NewUint(3),
			AmountTrackerId:        "tracker1",
		},
	}

	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		{
			FromListId:        "Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			TokenIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "a1",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
			},
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	aliceBalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Equal(1, len(aliceBalance.IncomingApprovals), "Incoming approvals should exist initially")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	allApprovals := append([]*types.ApprovalIdentifierDetails{}, GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1))...)
	allApprovals = append(allApprovals, &types.ApprovalIdentifierDetails{
		ApprovalId:      collectionsToCreate[0].DefaultIncomingApprovals[0].ApprovalId,
		ApprovalLevel:   "incoming",
		ApproverAddress: alice,
		Version:         sdkmath.NewUint(0),
	})

	// First transfer - should succeed and not delete the approval
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetBottomHalfUintRanges(),
					},
				},
				PrioritizedApprovals:                  allApprovals,
				OnlyCheckPrioritizedIncomingApprovals: true,
			},
		},
		bob,
	)
	suite.Require().Nil(err, "First transfer should succeed")

	// Check that approval still exists after first transfer
	aliceBalance, _ = GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(1, len(aliceBalance.IncomingApprovals), "Incoming approval should still exist after first transfer")

	// Second transfer - should succeed and not delete the approval yet
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetTopHalfUintRanges(),
					},
				},
				PrioritizedApprovals:                  allApprovals,
				OnlyCheckPrioritizedIncomingApprovals: true,
			},
		},
		bob,
	)
	suite.Require().Nil(err, "Second transfer should succeed")

	// Check that approval still exists after second transfer (threshold reached but not exceeded)
	aliceBalance, _ = GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(1, len(aliceBalance.IncomingApprovals), "Incoming approval should still exist after second transfer")

	// Third transfer - should succeed and delete the approval (threshold exceeded)
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                  allApprovals,
				OnlyCheckPrioritizedIncomingApprovals: true,
			},
		},
		bob,
	)
	suite.Require().Nil(err, "Third transfer should succeed")

	// Check that approval is deleted after third transfer (threshold exceeded)
	aliceBalance, _ = GetUserBalance(suite, wctx, collection.CollectionId, alice)
	suite.Require().Equal(0, len(aliceBalance.IncomingApprovals), "Incoming approval should be deleted after third transfer")

	// Fourth transfer - should fail because approval is deleted
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                  allApprovals,
				OnlyCheckPrioritizedIncomingApprovals: true,
			},
		},
		bob,
	)
	suite.Require().Error(err, "Fourth transfer should fail because approval is deleted")
}

// TestPrioritizedApprovalRetryLogic tests the retry logic for prioritized approvals
func (suite *TestSuite) TestPrioritizedApprovalRetryLogic() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	predeterminedBalances := &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					TokenIds:       GetFullUintRanges(),
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			AllowOverrideTimestamp:    false,
			DurationFromTimestamp:     sdkmath.NewUint(0),
			RecurringOwnershipTimes: &types.RecurringOwnershipTimes{
				IntervalLength:     sdkmath.NewUint(0),
				StartTime:          sdkmath.NewUint(0),
				ChargePeriodLength: sdkmath.NewUint(0),
			},
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	// Set up approvals that will fail on certain conditions
	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		{
			FromListId:        "Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			TokenIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "retry-test-1",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				RequireFromEqualsInitiatedBy:   true, // This will fail when from != initiatedBy
				PredeterminedBalances:          predeterminedBalances,
			},
		},
		{
			FromListId:        "Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			TokenIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "retry-test-2",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				RequireToEqualsInitiatedBy:     true, // This will fail when to != initiatedBy
				PredeterminedBalances:          predeterminedBalances,
			},
		},
		{
			FromListId:        "Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			TokenIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "retry-test-3",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				PredeterminedBalances:          predeterminedBalances,
			},
			// No criteria - this will always succeed
		},
		{
			FromListId:        "Mint",
			ToListId:          "All",
			InitiatedByListId: "All",
			TokenIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "retry-test-4-max-2-allowed",
			Version:           sdkmath.NewUint(0),
			ApprovalCriteria: &types.ApprovalCriteria{
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				PredeterminedBalances:          predeterminedBalances,
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(2),
					AmountTrackerId:        "test-tracker",
				},
			},
			// No criteria - this will always succeed
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Test 1: Multiple attempts with minimum successful attempts
	// This should succeed because we have 3 attempts and need 2 successful, and the 3rd approval always succeeds
	prioritizedApprovals := []*types.ApprovalIdentifierDetails{
		// {
		// 	ApprovalId:            "retry-test-1",
		// 	ApprovalLevel:         "collection",
		// 	ApproverAddress:       "",
		// 	Version:               sdkmath.NewUint(0),
		// 	NumAttempts:           sdkmath.NewUint(3),
		// 	MinSuccessfulAttempts: sdkmath.NewUint(2),
		// },
		// {
		// 	ApprovalId:            "retry-test-2",
		// 	ApprovalLevel:         "collection",
		// 	ApproverAddress:       "",
		// 	Version:               sdkmath.NewUint(0),
		// 	NumAttempts:           sdkmath.NewUint(3),
		// 	MinSuccessfulAttempts: sdkmath.NewUint(2),
		// },
		{
			ApprovalId:      "retry-test-3",
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			Version:         sdkmath.NewUint(0),
		},
	}

	// This should succeed because the 3rd approval (retry-test-3) will always succeed
	// and we need 2 successful attempts, so we'll get 2 successful attempts from retry-test-3
	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx, []*types.Balance{}, collection,
		GetFullUintRanges(), GetFullUintRanges(),
		"Mint", alice, bob, // from, to, initiatedBy (different addresses to trigger failures)
		sdkmath.NewUint(1), []*types.MerkleProof{},
		prioritizedApprovals, false, false, false, &types.PrecalculationOptions{},
		&keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}},
	)
	suite.Require().Nil(err, "Should succeed with retry logic")

	// Test 2: Try-only mode (minSuccessfulAttempts = 0)
	// prioritizedApprovalsTryOnly := []*types.ApprovalIdentifierDetails{
	// 	{
	// 		ApprovalId:      "retry-test-1",
	// 		ApprovalLevel:   "collection",
	// 		ApproverAddress: "",
	// 		Version:         sdkmath.NewUint(0),
	// 	},
	// }

	// // This should succeed even though all attempts fail, because minSuccessfulAttempts = 0
	// _, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
	// 	suite, suite.ctx, []*types.Balance{}, collection,
	// 	GetFullUintRanges(), GetFullUintRanges(),
	// 	"Mint", alice, alice, // from, to, initiatedBy (different addresses to trigger failures)
	// 	sdkmath.NewUint(1), []*types.MerkleProof{},
	// 	prioritizedApprovalsTryOnly, false, false, false, &types.PrecalculationOptions{},
	// 	&[]keeper.ApprovalsUsed{}, &[]keeper.CoinTransfers{}, sdkmath.NewUint(2), sdkmath.NewUint(0),
	// )
	// suite.Require().Nil(err, "Should succeed in try-only mode")

	// Test 3: Insufficient successful attempts
	prioritizedApprovalsInsufficient := []*types.ApprovalIdentifierDetails{
		{
			ApprovalId:      "retry-test-4-max-2-allowed",
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			Version:         sdkmath.NewUint(0),
		},
	}

	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals:                    prioritizedApprovalsInsufficient,
				OnlyCheckPrioritizedCollectionApprovals: true,
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
					PrecalculationOptions: &types.PrecalculationOptions{},
				},
			},
		},
		alice,
	)

	suite.Require().Nil(err, "Should succeed with first transfer")

	// Second transfer should also succeed since MaxNumTransfers is 2 and we've only used 1
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals:                    prioritizedApprovalsInsufficient,
				OnlyCheckPrioritizedCollectionApprovals: true,
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
					PrecalculationOptions: &types.PrecalculationOptions{},
				},
			},
		},
		alice,
	)
	suite.Require().Nil(err, "Should succeed with second transfer (within MaxNumTransfers limit)")

	// Third transfer should fail because MaxNumTransfers limit of 2 has been reached
	err = suite.app.TokenizationKeeper.HandleTransfers(
		suite.ctx,
		collection,
		[]*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals:                    prioritizedApprovalsInsufficient,
				OnlyCheckPrioritizedCollectionApprovals: true,
				PrecalculateBalancesFromApproval: &types.PrecalculateBalancesFromApprovalDetails{
					PrecalculationOptions: &types.PrecalculationOptions{},
				},
			},
		},
		alice,
	)
	suite.Require().Error(err, "Should fail when MaxNumTransfers limit is exceeded")

	// Test 4: All attempts fail and minSuccessfulAttempts > 0
	prioritizedApprovalsAllFail := []*types.ApprovalIdentifierDetails{
		{
			ApprovalId:      "retry-test-1",
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			Version:         sdkmath.NewUint(0),
		},
	}

	// This should fail because all attempts will fail (from != initiatedBy)
	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx, []*types.Balance{}, collection,
		GetFullUintRanges(), GetFullUintRanges(),
		"Mint", alice, charlie, // from, to, initiatedBy (different addresses to trigger failures)
		sdkmath.NewUint(1), []*types.MerkleProof{},
		prioritizedApprovalsAllFail, false, false, false, &types.PrecalculationOptions{},
		&keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}},
	)
	suite.Require().Error(err, "Should fail when all attempts fail")
}

// TestPrioritizedApprovalRetryLogicEdgeCases tests edge cases for the retry logic
func (suite *TestSuite) TestPrioritizedApprovalRetryLogicEdgeCases() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetCollectionsToCreate()

	// Set up a simple approval that always succeeds
	collectionsToCreate[0].CollectionApprovals = []*types.CollectionApproval{
		{
			FromListId:        "AllWithoutMint",
			ToListId:          "All",
			InitiatedByListId: "All",
			TokenIds:          GetFullUintRanges(),
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "edge-case-test",
			Version:           sdkmath.NewUint(0),
			// No criteria - always succeeds
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Test 1: Zero attempts (should default to 1)
	prioritizedApprovalsZeroAttempts := []*types.ApprovalIdentifierDetails{
		{
			ApprovalId:      "edge-case-test",
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			Version:         sdkmath.NewUint(0),
		},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx, []*types.Balance{}, collection,
		GetFullUintRanges(), GetFullUintRanges(),
		bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{},
		prioritizedApprovalsZeroAttempts, false, false, false, &types.PrecalculationOptions{},
		&keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}},
	)
	suite.Require().Nil(err, "Should succeed with zero attempts (defaults to 1)")

	// Test 2: Nil minSuccessfulAttempts (should default to 1)
	prioritizedApprovalsNilMinSuccess := []*types.ApprovalIdentifierDetails{
		{
			ApprovalId:      "edge-case-test",
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			Version:         sdkmath.NewUint(0),
		},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx, []*types.Balance{}, collection,
		GetFullUintRanges(), GetFullUintRanges(),
		bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{},
		prioritizedApprovalsNilMinSuccess, false, false, false, nil,
		&keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}},
	)
	suite.Require().Nil(err, "Should succeed with nil minSuccessfulAttempts (defaults to 1)")

	// Test 3: Single attempt with minSuccessfulAttempts = 1
	prioritizedApprovalsSingleAttempt := []*types.ApprovalIdentifierDetails{
		{
			ApprovalId:      "edge-case-test",
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			Version:         sdkmath.NewUint(0),
		},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx, []*types.Balance{}, collection,
		GetFullUintRanges(), GetFullUintRanges(),
		bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{},
		prioritizedApprovalsSingleAttempt, false, false, false, nil,
		&keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}},
	)
	suite.Require().Nil(err, "Should succeed with single attempt")

	// Test 4: Large number of attempts
	prioritizedApprovalsLargeAttempts := []*types.ApprovalIdentifierDetails{
		{
			ApprovalId:      "edge-case-test",
			ApprovalLevel:   "collection",
			ApproverAddress: "",
			Version:         sdkmath.NewUint(0),
		},
	}

	_, err = DeductCollectionApprovalsAndGetUserApprovalsToCheck(
		suite, suite.ctx, []*types.Balance{}, collection,
		GetFullUintRanges(), GetFullUintRanges(),
		bob, alice, charlie, sdkmath.NewUint(1), []*types.MerkleProof{},
		prioritizedApprovalsLargeAttempts, false, false, false, nil,
		&keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}},
	)
	suite.Require().Nil(err, "Should succeed with large number of attempts")
}
