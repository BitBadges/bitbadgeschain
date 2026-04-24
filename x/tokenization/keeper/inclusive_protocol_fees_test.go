package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Inclusive protocol fee: the payer sends exactly the coin transfer amount.
// The chain skims 0.1% for the community pool out of that amount (not on top).
// See backlog #0335.

// configureSingleCoinTransferApproval wires up a mint approval whose coinTransfer
// sends `amount` of ubadge from bob (creator) to alice on behalf of a mint.
func (suite *TestSuite) configureSingleCoinTransferApproval(amount sdkmath.Int) {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.CoinTransfers = []*types.CoinTransfer{
		{
			To: alice,
			Coins: []*sdk.Coin{
				{Amount: amount, Denom: "ubadge"},
			},
		},
	}
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating collection")
}

func (suite *TestSuite) runMintTransferForCoinTransferApproval() error {
	wctx := sdk.WrapSDKContext(suite.ctx)
	return TransferTokens(suite, wctx, &types.MsgTransferTokens{
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
}

// TestInclusiveProtocolFeeDebitsPayerExactly verifies the payer is debited exactly
// the coin transfer amount (not amount + fee). The fee is carved out of the payment.
func (suite *TestSuite) TestInclusiveProtocolFeeDebitsPayerExactly() {
	const gross int64 = 10_000 // 0.1% = 10 → non-zero fee
	const expectedFee int64 = 10

	suite.configureSingleCoinTransferApproval(sdkmath.NewInt(gross))

	bobAddr := sdk.MustAccAddressFromBech32(bob)
	aliceAddr := sdk.MustAccAddressFromBech32(alice)
	poolAddr := suite.app.DistrKeeper.GetDistributionAccount(suite.ctx).GetAddress()

	bobBefore := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, "ubadge").Amount
	aliceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, "ubadge").Amount
	poolBefore := suite.app.BankKeeper.GetBalance(suite.ctx, poolAddr, "ubadge").Amount

	err := suite.runMintTransferForCoinTransferApproval()
	suite.Require().Nil(err, "transfer should succeed")

	bobAfter := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, "ubadge").Amount
	aliceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, "ubadge").Amount
	poolAfter := suite.app.BankKeeper.GetBalance(suite.ctx, poolAddr, "ubadge").Amount

	// Payer sends exactly gross — not gross + fee.
	suite.Require().Equal(
		sdkmath.NewInt(gross),
		bobBefore.Sub(bobAfter),
		"payer should be debited exactly the coin transfer amount (inclusive fee)",
	)

	// Recipient gets gross - fee.
	suite.Require().Equal(
		sdkmath.NewInt(gross-expectedFee),
		aliceAfter.Sub(aliceBefore),
		"recipient should receive gross minus the inclusive protocol fee",
	)

	// Community pool received the fee.
	suite.Require().Equal(
		sdkmath.NewInt(expectedFee),
		poolAfter.Sub(poolBefore),
		"community pool should receive the protocol fee",
	)
}

// TestInclusiveProtocolFeeRoundsToZero verifies small amounts below the fee
// denominator (1000 ubadge) pay no fee — matches the previous behavior at
// the low end, so existing integrations with small amounts are unaffected.
func (suite *TestSuite) TestInclusiveProtocolFeeRoundsToZero() {
	const gross int64 = 100 // 0.1% = 0.1 → rounds down to 0

	suite.configureSingleCoinTransferApproval(sdkmath.NewInt(gross))

	bobAddr := sdk.MustAccAddressFromBech32(bob)
	aliceAddr := sdk.MustAccAddressFromBech32(alice)
	poolAddr := suite.app.DistrKeeper.GetDistributionAccount(suite.ctx).GetAddress()

	bobBefore := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, "ubadge").Amount
	aliceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, "ubadge").Amount
	poolBefore := suite.app.BankKeeper.GetBalance(suite.ctx, poolAddr, "ubadge").Amount

	err := suite.runMintTransferForCoinTransferApproval()
	suite.Require().Nil(err, "transfer should succeed")

	bobAfter := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, "ubadge").Amount
	aliceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, "ubadge").Amount
	poolAfter := suite.app.BankKeeper.GetBalance(suite.ctx, poolAddr, "ubadge").Amount

	suite.Require().Equal(sdkmath.NewInt(gross), bobBefore.Sub(bobAfter), "payer debited gross")
	suite.Require().Equal(sdkmath.NewInt(gross), aliceAfter.Sub(aliceBefore), "recipient gets full amount when fee rounds to zero")
	suite.Require().Equal(sdkmath.ZeroInt(), poolAfter.Sub(poolBefore), "no fee routed to community pool for sub-denominator amounts")
}

// TestInclusiveProtocolFeeWithRoyalty verifies royalty + protocol fee are both
// carved out of the gross payment. The payer still sends exactly the quoted
// amount, and the recipient gets gross - royalty - fee.
func (suite *TestSuite) TestInclusiveProtocolFeeWithRoyalty() {
	const gross int64 = 10_000
	const royaltyBps = 1000 // 10%
	const expectedFee int64 = 10
	const expectedRoyalty int64 = 1000
	const expectedRecipient int64 = gross - expectedFee - expectedRoyalty

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.UserApprovalSettings = &types.UserApprovalSettings{
		UserRoyalties: &types.UserRoyalties{
			Percentage:    sdkmath.NewUint(royaltyBps),
			PayoutAddress: charlie,
		},
	}
	suite.Require().Nil(CreateCollections(suite, wctx, collectionsToCreate), "error creating collection")

	suite.Require().Nil(UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ToListId:          "AllWithoutMint",
				InitiatedByListId: alice,
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				ApprovalId:        "test",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
						AmountTrackerId:        "test-tracker",
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						AmountTrackerId:              "test-tracker",
					},
					CoinTransfers: []*types.CoinTransfer{
						{
							To:                              alice,
							OverrideFromWithApproverAddress: true, // bob (approver) pays
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(gross), Denom: "ubadge"},
							},
						},
					},
				},
			},
		},
	}), "error updating user approvals")

	bobAddr := sdk.MustAccAddressFromBech32(bob)
	aliceAddr := sdk.MustAccAddressFromBech32(alice)
	charlieAddr := sdk.MustAccAddressFromBech32(charlie)
	poolAddr := suite.app.DistrKeeper.GetDistributionAccount(suite.ctx).GetAddress()

	bobBefore := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, "ubadge").Amount
	aliceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, "ubadge").Amount
	charlieBefore := suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge").Amount
	poolBefore := suite.app.BankKeeper.GetBalance(suite.ctx, poolAddr, "ubadge").Amount

	err := TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{ApprovalId: "test", ApprovalLevel: "collection", Version: sdkmath.NewUint(0)},
					{ApprovalId: "test", ApprovalLevel: "outgoing", ApproverAddress: bob, Version: sdkmath.NewUint(1)},
				},
			},
		},
	})
	suite.Require().Nil(err, "transfer should succeed")

	bobAfter := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, "ubadge").Amount
	aliceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, aliceAddr, "ubadge").Amount
	charlieAfter := suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge").Amount
	poolAfter := suite.app.BankKeeper.GetBalance(suite.ctx, poolAddr, "ubadge").Amount

	suite.Require().Equal(sdkmath.NewInt(gross), bobBefore.Sub(bobAfter), "payer debited exactly gross")
	suite.Require().Equal(sdkmath.NewInt(expectedRecipient), aliceAfter.Sub(aliceBefore), "recipient gets gross - fee - royalty")
	suite.Require().Equal(sdkmath.NewInt(expectedRoyalty), charlieAfter.Sub(charlieBefore), "royalty payout receives full royalty % of gross")
	suite.Require().Equal(sdkmath.NewInt(expectedFee), poolAfter.Sub(poolBefore), "community pool receives fee % of gross")
}

// TestInclusiveProtocolFeeRejectsRoyaltyPlusFeeOverflow verifies that a royalty
// percentage high enough that royalty + fee > gross is rejected rather than
// silently shortchanging the recipient by making them negative.
func (suite *TestSuite) TestInclusiveProtocolFeeRejectsRoyaltyPlusFeeOverflow() {
	const gross int64 = 10_000

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	wctx := sdk.WrapSDKContext(suite.ctx)

	// 100% royalty leaves zero for the recipient and 0.1% for the fee — overflow.
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.UserApprovalSettings = &types.UserApprovalSettings{
		UserRoyalties: &types.UserRoyalties{
			Percentage:    sdkmath.NewUint(10000), // 100%
			PayoutAddress: charlie,
		},
	}
	suite.Require().Nil(CreateCollections(suite, wctx, collectionsToCreate), "error creating collection")

	suite.Require().Nil(UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 bob,
		CollectionId:            sdkmath.NewUint(1),
		UpdateOutgoingApprovals: true,
		OutgoingApprovals: []*types.UserOutgoingApproval{
			{
				ToListId:          "AllWithoutMint",
				InitiatedByListId: alice,
				TransferTimes:     GetFullUintRanges(),
				OwnershipTimes:    GetFullUintRanges(),
				TokenIds:          []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				ApprovalId:        "test",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
						AmountTrackerId:        "test-tracker",
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
						AmountTrackerId:              "test-tracker",
					},
					CoinTransfers: []*types.CoinTransfer{
						{
							To:                              alice,
							OverrideFromWithApproverAddress: true,
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(gross), Denom: "ubadge"},
							},
						},
					},
				},
			},
		},
	}), "error updating user approvals")

	err := TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
						Amount:         sdkmath.NewUint(1),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{ApprovalId: "test", ApprovalLevel: "collection", Version: sdkmath.NewUint(0)},
					{ApprovalId: "test", ApprovalLevel: "outgoing", ApproverAddress: bob, Version: sdkmath.NewUint(1)},
				},
			},
		},
	})
	suite.Require().Error(err, "royalty + protocol fee exceeding gross should be rejected")
}
