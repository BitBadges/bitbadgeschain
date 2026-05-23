package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// Protocol fee model (additive): a 0.1% fee (amount / ProtocolFeeDenominator) is
// charged to the initiator on top of each coin transfer. The recipient receives
// the full quoted amount; the fee is routed to the community pool. These tests
// exercise the fee with amounts large enough to produce a non-zero fee, and assert
// all three sides of the transfer (payer, recipient, community pool).

func (suite *TestSuite) protocolFeeCollection(coinTransfers []*types.CoinTransfer) []*types.MsgNewCollection {
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].CollectionApprovals[0].FromListId = "Mint"
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.CoinTransfers = coinTransfers
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesFromOutgoingApprovals = true
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.OverridesToIncomingApprovals = true
	return collectionsToCreate
}

func (suite *TestSuite) mintTransferTo(recipient string) error {
	return TransferTokens(suite, sdk.WrapSDKContext(suite.ctx), &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{recipient},
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

func (suite *TestSuite) ubadgeBalance(addr string) sdkmath.Int {
	return suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(addr), "ubadge").Amount
}

func (suite *TestSuite) communityPoolUbadge() sdkmath.Int {
	distrAddr := authtypes.NewModuleAddress(distrtypes.ModuleName)
	return suite.app.BankKeeper.GetBalance(suite.ctx, distrAddr, "ubadge").Amount
}

// TestProtocolFeeChargedOnTopOfTransfer verifies the payer is debited the transfer
// amount PLUS the 0.1% fee, the recipient receives the full transfer amount, and the
// fee lands in the community pool.
func (suite *TestSuite) TestProtocolFeeChargedOnTopOfTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	const transferAmount = int64(100000) // 0.1% = 100, non-zero fee
	const expectedFee = transferAmount / keeper.ProtocolFeeDenominator

	collectionsToCreate := suite.protocolFeeCollection([]*types.CoinTransfer{
		{
			To:    alice,
			Coins: []*sdk.Coin{{Amount: sdkmath.NewInt(transferAmount), Denom: "ubadge"}},
		},
	})
	suite.Require().Nil(CreateCollections(suite, wctx, collectionsToCreate), "error creating tokens")

	bobBefore := suite.ubadgeBalance(bob)
	aliceBefore := suite.ubadgeBalance(alice)
	poolBefore := suite.communityPoolUbadge()

	suite.Require().Nil(suite.mintTransferTo(alice), "error transferring tokens")

	suite.Require().Equal(bobBefore.SubRaw(transferAmount).SubRaw(expectedFee), suite.ubadgeBalance(bob),
		"payer should be charged transfer amount + protocol fee")
	suite.Require().Equal(aliceBefore.AddRaw(transferAmount), suite.ubadgeBalance(alice),
		"recipient should receive the full transfer amount (fee not skimmed)")
	suite.Require().Equal(poolBefore.AddRaw(expectedFee), suite.communityPoolUbadge(),
		"community pool should receive the protocol fee")
}

// TestProtocolFeeAggregatesAcrossTransfers verifies the fee is computed on the summed
// per-denom total across all coin transfers in the message, not per individual transfer.
func (suite *TestSuite) TestProtocolFeeAggregatesAcrossTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	const toAlice = int64(100000)
	const toCharlie = int64(50000)
	const expectedFee = (toAlice + toCharlie) / keeper.ProtocolFeeDenominator // 0.1% of 150000 = 150

	collectionsToCreate := suite.protocolFeeCollection([]*types.CoinTransfer{
		{To: alice, Coins: []*sdk.Coin{{Amount: sdkmath.NewInt(toAlice), Denom: "ubadge"}}},
		{To: charlie, Coins: []*sdk.Coin{{Amount: sdkmath.NewInt(toCharlie), Denom: "ubadge"}}},
	})
	suite.Require().Nil(CreateCollections(suite, wctx, collectionsToCreate), "error creating tokens")

	bobBefore := suite.ubadgeBalance(bob)
	aliceBefore := suite.ubadgeBalance(alice)
	charlieBefore := suite.ubadgeBalance(charlie)
	poolBefore := suite.communityPoolUbadge()

	suite.Require().Nil(suite.mintTransferTo(charlie), "error transferring tokens")

	suite.Require().Equal(bobBefore.SubRaw(toAlice).SubRaw(toCharlie).SubRaw(expectedFee), suite.ubadgeBalance(bob),
		"payer should be charged both transfers + a single aggregated fee")
	suite.Require().Equal(aliceBefore.AddRaw(toAlice), suite.ubadgeBalance(alice),
		"alice should receive her full amount")
	suite.Require().Equal(charlieBefore.AddRaw(toCharlie), suite.ubadgeBalance(charlie),
		"charlie should receive his full amount")
	suite.Require().Equal(poolBefore.AddRaw(expectedFee), suite.communityPoolUbadge(),
		"community pool should receive the aggregated fee")
}

// TestProtocolFeeRoundsToZeroForSmallTransfers verifies amounts below the fee
// denominator pay no fee (0.1% rounds down to 0), and nothing reaches the community pool.
func (suite *TestSuite) TestProtocolFeeRoundsToZeroForSmallTransfers() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	const transferAmount = int64(500) // 0.1% = 0.5, rounds down to 0

	collectionsToCreate := suite.protocolFeeCollection([]*types.CoinTransfer{
		{
			To:    alice,
			Coins: []*sdk.Coin{{Amount: sdkmath.NewInt(transferAmount), Denom: "ubadge"}},
		},
	})
	suite.Require().Nil(CreateCollections(suite, wctx, collectionsToCreate), "error creating tokens")

	bobBefore := suite.ubadgeBalance(bob)
	aliceBefore := suite.ubadgeBalance(alice)
	poolBefore := suite.communityPoolUbadge()

	suite.Require().Nil(suite.mintTransferTo(alice), "error transferring tokens")

	suite.Require().Equal(bobBefore.SubRaw(transferAmount), suite.ubadgeBalance(bob),
		"payer should be charged only the transfer amount when the fee rounds to zero")
	suite.Require().Equal(aliceBefore.AddRaw(transferAmount), suite.ubadgeBalance(alice),
		"recipient should receive the full transfer amount")
	suite.Require().Equal(poolBefore, suite.communityPoolUbadge(),
		"community pool should be unchanged when the fee rounds to zero")
}
