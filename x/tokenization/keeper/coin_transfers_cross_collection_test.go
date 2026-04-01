package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TestCrossCollectionCoinTransfersViaBadgeslpAlias is an E2E test that:
//  1. Creates collection 1 with alias paths (badgeslp:1:wrappedcoin), mints tokens to bob
//  2. Creates collection 2 whose transfer approval includes a coinTransfer using badgeslp:1:wrappedcoin
//  3. Bob transfers tokens on collection 2 to alice, triggering the coinTransfer
//  4. Verifies collection 1 token balances: bob lost 3, alice gained 3
func (suite *TestSuite) TestCrossCollectionCoinTransfersViaBadgeslpAlias() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// ── Step 1: Create collection 1 with alias paths ──
	// Token ID 1 minted to bob, with alias "badgeslp:1:wrappedcoin" mapping 1:1 to token ID 1.
	col1 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	col1[0].AliasPathsToAdd = []*types.AliasPathAddObject{
		{
			Denom: "wrappedcoin",
			Conversion: &types.ConversionWithoutDenom{
				SideA: &types.ConversionSideA{
					Amount: sdkmath.NewUint(1),
				},
				SideB: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						OwnershipTimes: GetFullUintRanges(),
						TokenIds:       GetOneUintRange(),
					},
				},
			},
			Symbol:     "WRAP",
			DenomUnits: []*types.DenomUnit{{Decimals: sdkmath.NewUint(6), Symbol: "wrappedcoin", IsDefaultDisplay: true}},
		},
	}

	// Override outgoing/incoming approvals at collection level so alias transfers aren't blocked
	// by user-level per-address limits (default is 1 per from address).
	col1[0].CollectionApprovals[1].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	col1[0].CollectionApprovals[1].ApprovalCriteria.OverridesToIncomingApprovals = true

	err := CreateCollections(suite, wctx, col1)
	suite.Require().Nil(err, "error creating collection 1 with alias paths")

	// Verify alias path was created
	collection1, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err)
	suite.Require().Len(collection1.AliasPaths, 1)
	aliasDenom := generateAliasWrapperDenom(collection1.CollectionId, collection1.AliasPaths[0])
	suite.Require().Equal(keeper.AliasDenomPrefix+"1:wrappedcoin", aliasDenom)

	// Mint 9 more tokens on collection 1 so bob has 10 total of token ID 1
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(9),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "mint-test",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "error minting additional tokens on collection 1")

	// Verify bob has 10 of token ID 1 on collection 1
	bobCol1Before, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	bobToken1Before, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobCol1Before.Balances)
	suite.Require().Nil(err)
	AssertBalancesEqual(suite, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(10),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}, bobToken1Before)

	// Verify alice has 0 of token ID 1 on collection 1
	aliceCol1Before, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err)
	aliceToken1Before, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), aliceCol1Before.Balances)
	suite.Require().Nil(err)
	AssertBalancesEqual(suite, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(0),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}, aliceToken1Before)

	// ── Step 2: Create collection 2 with coinTransfers using badgeslp:1:wrappedcoin ──
	// When bob transfers on collection 2, the coinTransfer sends 3 alias coins from bob (initiator) to alice.
	col2 := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	col2[0].CollectionApprovals = []*types.CollectionApproval{
		// Mint approval
		{
			ToListId:          "AllWithoutMint",
			FromListId:        "Mint",
			InitiatedByListId: "AllWithoutMint",
			TransferTimes:     GetFullUintRanges(),
			TokenIds:          GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			ApprovalId:        "mint-test",
			ApprovalCriteria: &types.ApprovalCriteria{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(1000),
					AmountTrackerId:        "mint-test-tracker",
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
					AmountTrackerId:              "mint-test-tracker",
				},
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
			},
		},
		// Transfer approval with alias coinTransfer
		{
			ApprovalId:        "transfer-with-alias-coin",
			TransferTimes:     GetFullUintRanges(),
			OwnershipTimes:    GetFullUintRanges(),
			TokenIds:          GetFullUintRanges(),
			FromListId:        "AllWithoutMint",
			ToListId:          "AllWithoutMint",
			InitiatedByListId: "AllWithoutMint",
			ApprovalCriteria: &types.ApprovalCriteria{
				MaxNumTransfers: &types.MaxNumTransfers{
					OverallMaxNumTransfers: sdkmath.NewUint(1000),
					AmountTrackerId:        "transfer-tracker",
				},
				ApprovalAmounts: &types.ApprovalAmounts{
					PerFromAddressApprovalAmount: sdkmath.NewUint(uint64(math.MaxUint64)),
					AmountTrackerId:              "transfer-tracker",
				},
				OverridesFromOutgoingApprovals: true,
				OverridesToIncomingApprovals:   true,
				CoinTransfers: []*types.CoinTransfer{
					{
						To: alice,
						Coins: []*sdk.Coin{
							{Amount: sdkmath.NewInt(3), Denom: aliasDenom}, // 3 badgeslp:1:wrappedcoin
						},
					},
				},
			},
		},
	}

	err = CreateCollections(suite, wctx, col2)
	suite.Require().Nil(err, "error creating collection 2 with alias coinTransfers")

	// ── Step 3: Bob transfers on collection 2 (bob → alice) ──
	// Bob is Creator (initiatedBy) so coinTransfer sends from bob.
	// sendmanager routes badgeslp:1:wrappedcoin → tokenization keeper → moves collection 1 tokens.
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(2),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{
						ApprovalId:      "transfer-with-alias-coin",
						ApprovalLevel:   "collection",
						ApproverAddress: "",
						Version:         sdkmath.NewUint(0),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "transfer on collection 2 should succeed and route alias coinTransfer through sendmanager")

	// ── Step 4: Verify collection 1 balances changed correctly ──
	// Bob should have 10 - 3 = 7 of token ID 1, alice should have 3 of token ID 1 on collection 1.
	bobCol1After, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err)
	bobToken1After, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), bobCol1After.Balances)
	suite.Require().Nil(err)
	AssertBalancesEqual(suite, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(7),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}, bobToken1After)

	aliceCol1After, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	suite.Require().Nil(err)
	aliceToken1After, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetFullUintRanges(), aliceCol1After.Balances)
	suite.Require().Nil(err)
	AssertBalancesEqual(suite, []*types.Balance{
		{
			Amount:         sdkmath.NewUint(3),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}, aliceToken1After)
}
