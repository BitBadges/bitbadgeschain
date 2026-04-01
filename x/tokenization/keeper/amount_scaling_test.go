package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Helper to create a collection with amount scaling enabled on the mint approval (index 0).
// The approval uses incrementedBalances with startBalances as the 1x base and allowAmountScaling=true.
// CoinTransfers charge `coinAmount` ubadge per base unit (initiator pays, sent to alice).
// maxMultiplier caps the scaling (must be > 0 when scaling is on).
func (suite *TestSuite) createScalingCollection(baseTokenAmount uint64, coinAmount int64, maxMultiplier uint64) {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	// The mint approval is at index 0 (prepended by GetTransferableCollectionToCreateAllMintedToCreator).
	mintApproval := collectionsToCreate[0].CollectionApprovals[0]
	mintApproval.ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(baseTokenAmount),
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
			AllowAmountScaling:        true,
			MaxScalingMultiplier:      sdkmath.NewUint(maxMultiplier),
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	if coinAmount > 0 {
		mintApproval.ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{
			{
				To: alice,
				Coins: []*sdk.Coin{
					{Amount: sdkmath.NewInt(coinAmount), Denom: "ubadge"},
				},
			},
		}
	}

	// Remove the initial mint transfer since predetermined balances
	// require exact/multiple match — the default transfer mints full range at amount 1
	// which won't match our specific base tokenIds.
	collectionsToCreate[0].Transfers = []*types.Transfer{}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating scaling collection")
}

// TestAmountScaling1x verifies that a 1x multiplier works (transfer exactly equals base).
func (suite *TestSuite) TestAmountScaling1x() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.createScalingCollection(1, 100, 100)

	aliceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")

	err := TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                    GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Nil(err, "1x scaling transfer should succeed")

	aliceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(aliceBefore.Amount.Add(sdkmath.NewInt(100)), aliceAfter.Amount, "alice should receive 100 ubadge (1x)")
}

// TestAmountScaling5x verifies that a 5x multiplier scales coinTransfers by 5.
func (suite *TestSuite) TestAmountScaling5x() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.createScalingCollection(1, 100, 100)

	aliceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")

	err := TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(5),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                    GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Nil(err, "5x scaling transfer should succeed")

	aliceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(aliceBefore.Amount.Add(sdkmath.NewInt(500)), aliceAfter.Amount, "alice should receive 500 ubadge (5x)")
}

// TestAmountScalingNotEvenlyDivisible verifies rejection when transfer is not evenly divisible.
func (suite *TestSuite) TestAmountScalingNotEvenlyDivisible() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.createScalingCollection(3, 100, 100)

	// 7 is not evenly divisible by base of 3
	err := TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(7),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                    GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Error(err, "non-evenly divisible transfer should be rejected")
}

// TestAmountScalingDisabled verifies exact match is still enforced when scaling is off.
func (suite *TestSuite) TestAmountScalingDisabled() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	approval := collectionsToCreate[0].CollectionApprovals[0]
	approval.ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
			AllowAmountScaling:        false, // disabled
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	collectionsToCreate[0].Transfers = []*types.Transfer{}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")

	// Transfer of 5 should fail because scaling is off and base is 1
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(5),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                    GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Error(err, "5x transfer should fail when scaling is disabled")
}

// TestAmountScalingNoCoinTransfers verifies scaling works with no coin transfers (free any-quantity).
func (suite *TestSuite) TestAmountScalingNoCoinTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.createScalingCollection(1, 0, 100)

	err := TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(10),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                    GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Nil(err, "scaling with no coin transfers should succeed")
}

// TestAmountScalingValidateBasicNonZeroIncrements verifies rejection at ValidateBasic.
func (suite *TestSuite) TestAmountScalingValidateBasicNonZeroIncrements() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	approval := collectionsToCreate[0].CollectionApprovals[0]
	approval.ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(1), // non-zero - should be rejected
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
			AllowAmountScaling:        true,
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "allowAmountScaling with non-zero incrementTokenIdsBy should be rejected")
}

// TestAmountScalingValidateBasicDurationFromTimestamp verifies rejection when durationFromTimestamp is set.
func (suite *TestSuite) TestAmountScalingValidateBasicDurationFromTimestamp() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	approval := collectionsToCreate[0].CollectionApprovals[0]
	approval.ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(1000), // non-zero - should be rejected
			AllowAmountScaling:        true,
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "allowAmountScaling with non-zero durationFromTimestamp should be rejected")
}

// TestAmountScalingValidateBasicAllowOverrideTimestamp verifies rejection when allowOverrideTimestamp is set.
func (suite *TestSuite) TestAmountScalingValidateBasicAllowOverrideTimestamp() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	approval := collectionsToCreate[0].CollectionApprovals[0]
	approval.ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
			AllowOverrideTimestamp:    true, // should be rejected
			AllowAmountScaling:        true,
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "allowAmountScaling with allowOverrideTimestamp should be rejected")
}

// TestAmountScalingValidateBasicAllowOverrideWithAnyValidToken verifies rejection.
func (suite *TestSuite) TestAmountScalingValidateBasicAllowOverrideWithAnyValidToken() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	approval := collectionsToCreate[0].CollectionApprovals[0]
	approval.ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:            sdkmath.NewUint(0),
			IncrementOwnershipTimesBy:       sdkmath.NewUint(0),
			DurationFromTimestamp:           sdkmath.NewUint(0),
			AllowOverrideWithAnyValidToken: true, // should be rejected
			AllowAmountScaling:              true,
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "allowAmountScaling with allowOverrideWithAnyValidToken should be rejected")
}

// TestAmountScalingValidateBasicEmptyStartBalances verifies rejection when startBalances is empty.
func (suite *TestSuite) TestAmountScalingValidateBasicEmptyStartBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	approval := collectionsToCreate[0].CollectionApprovals[0]
	approval.ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances:             []*types.Balance{}, // empty - should be rejected
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
			AllowAmountScaling:        true,
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "allowAmountScaling with empty startBalances should be rejected")
}

// TestAmountScalingMultipleTransfers verifies tracker increments correctly across multiple scaled transfers.
func (suite *TestSuite) TestAmountScalingMultipleTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.createScalingCollection(1, 100, 100)

	aliceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")

	// Transfer 3x
	err := TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(3),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                    GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Nil(err, "first 3x transfer should succeed")

	// Transfer 2x
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(2),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                    GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Nil(err, "second 2x transfer should succeed")

	aliceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(alice), "ubadge")
	suite.Require().Equal(aliceBefore.Amount.Add(sdkmath.NewInt(500)), aliceAfter.Amount, "alice should receive 300+200=500 ubadge total")
}

// TestAmountScalingMaxMultiplierEnforced verifies the cap is enforced at runtime.
func (suite *TestSuite) TestAmountScalingMaxMultiplierEnforced() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	suite.createScalingCollection(1, 100, 3) // max 3x

	// 3x should succeed
	err := TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(3),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                    GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Nil(err, "3x transfer should succeed with max 3")

	// 4x should fail
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{charlie},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(4),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals:                    GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
				OnlyCheckPrioritizedCollectionApprovals: true,
			},
		},
	})
	suite.Require().Error(err, "4x transfer should fail with max 3")
}

// TestAmountScalingMaxMultiplierZeroRejected verifies ValidateBasic rejects max=0 with scaling on.
func (suite *TestSuite) TestAmountScalingMaxMultiplierZeroRejected() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	approval := collectionsToCreate[0].CollectionApprovals[0]
	approval.ApprovalCriteria.PredeterminedBalances = &types.PredeterminedBalances{
		IncrementedBalances: &types.IncrementedBalances{
			StartBalances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			IncrementTokenIdsBy:       sdkmath.NewUint(0),
			IncrementOwnershipTimesBy: sdkmath.NewUint(0),
			DurationFromTimestamp:     sdkmath.NewUint(0),
			AllowAmountScaling:        true,
			MaxScalingMultiplier:      sdkmath.NewUint(0), // zero — should be rejected
		},
		OrderCalculationMethod: &types.PredeterminedOrderCalculationMethod{
			UseOverallNumTransfers: true,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "maxScalingMultiplier == 0 with scaling on should be rejected")
}
