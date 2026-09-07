package keeper_test

import (
	"math"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func rollbackTestMintApproval() *types.CollectionApproval {
	return &types.CollectionApproval{
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
	}
}

func rollbackTestTransferApproval(id string, coinTransfers []*types.CoinTransfer) *types.CollectionApproval {
	return &types.CollectionApproval{
		ApprovalId:        id,
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		TokenIds:          GetFullUintRanges(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			MaxNumTransfers: &types.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        id + "-tracker",
			},
			ApprovalAmounts: &types.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(uint64(math.MaxUint64)),
				AmountTrackerId:              id + "-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
			CoinTransfers:                  coinTransfers,
		},
	}
}

// A priced approval executes its coin transfer before its num-transfers threshold is
// checked. When the threshold fails and a later approval serves the transfer, the
// coin transfer is rolled back and must not be charged a protocol fee or reported.
func (suite *TestSuite) TestRolledBackCoinTransferIsNotChargedOrReported() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	colA := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	priced := rollbackTestTransferApproval("priced-once", []*types.CoinTransfer{
		{To: charlie, Coins: []*sdk.Coin{{Amount: sdkmath.NewInt(1_000_000_000), Denom: "ubadge"}}},
	})
	priced.ApprovalCriteria.MaxNumTransfers.OverallMaxNumTransfers = sdkmath.NewUint(1)
	free := rollbackTestTransferApproval("free-fallback", nil)
	colA[0].CollectionApprovals = []*types.CollectionApproval{rollbackTestMintApproval(), priced, free}
	colA[0].Transfers[0].PrioritizedApprovals = []*types.ApprovalIdentifierDetails{
		{ApprovalId: "mint-test", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
	}
	suite.Require().Nil(CreateCollections(suite, wctx, colA))

	suite.Require().Nil(TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From:        "Mint",
			ToAddresses: []string{bob},
			Balances:    []*types.Balance{{Amount: sdkmath.NewUint(9), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{ApprovalId: "mint-test", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
			},
		}},
	}))

	bobAddr := sdk.MustAccAddressFromBech32(bob)
	charlieAddr := sdk.MustAccAddressFromBech32(charlie)

	doTransfer := func() *types.MsgTransferTokensResponse {
		res, err := suite.msgServer.TransferTokens(wctx, &types.MsgTransferTokens{
			Creator:      bob,
			CollectionId: sdkmath.NewUint(1),
			Transfers: []*types.Transfer{{
				From:        bob,
				ToAddresses: []string{alice},
				Balances:    []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}},
				PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
					{ApprovalId: "priced-once", ApprovalLevel: "collection", ApproverAddress: "", Version: sdkmath.NewUint(0)},
				},
			}},
		})
		suite.Require().Nil(err, "transfer failed: %v", err)
		return res
	}

	// First use: priced-once succeeds (1 of 1), charlie is paid, bob pays the amount plus the protocol fee.
	res1 := doTransfer()
	suite.Require().Len(res1.ApprovalsUsed, 1)
	suite.Require().Equal("priced-once", res1.ApprovalsUsed[0].ApprovalId)
	suite.Require().Len(res1.CoinTransfers, 2, "payment plus protocol fee")

	bobBefore := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, "ubadge").Amount
	charlieBefore := suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge").Amount

	collection, found := suite.app.TokenizationKeeper.GetCollectionFromStore(suite.ctx, sdkmath.OneUint())
	suite.Require().True(found)
	tracking := &keeper.EventTracking{ApprovalsUsed: &[]keeper.ApprovalsUsed{}, CoinTransfers: &[]keeper.CoinTransfers{}}
	_, err := DeductCollectionApprovalsAndGetUserApprovalsToCheck(suite, suite.ctx,
		[]*types.Balance{{Amount: sdkmath.OneUint(), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}},
		collection, GetOneUintRange(), GetFullUintRanges(), bob, alice, bob, sdkmath.OneUint(), nil,
		[]*types.ApprovalIdentifierDetails{{ApprovalId: "priced-once", ApprovalLevel: "collection", Version: sdkmath.NewUint(0)}}, false, false, false, nil, tracking)
	suite.Require().NoError(err)
	suite.Require().Empty(*tracking.CoinTransfers, "failed approval attempts must not leak into their caller's payment accounting")

	// Second use: priced-once runs its coin transfer, then fails the num-transfers
	// threshold and is rolled back; free-fallback serves the transfer.
	res2 := doTransfer()

	bobAfter := suite.app.BankKeeper.GetBalance(suite.ctx, bobAddr, "ubadge").Amount
	charlieAfter := suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge").Amount

	suite.Require().Len(res2.ApprovalsUsed, 1)
	suite.Require().Equal("free-fallback", res2.ApprovalsUsed[0].ApprovalId)
	suite.Require().True(charlieAfter.Equal(charlieBefore), "charlie must not receive the rolled-back coin transfer")
	suite.Require().Len(res2.CoinTransfers, 0, "rolled-back coin transfers must not be reported")
	suite.Require().True(bobAfter.Equal(bobBefore), "bob must not pay a protocol fee for a rolled-back coin transfer (delta=%s)", bobAfter.Sub(bobBefore))
}
