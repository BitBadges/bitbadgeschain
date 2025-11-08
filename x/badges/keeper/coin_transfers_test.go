package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestUserLevelRoyalties() {
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.UserRoyalties = &types.UserRoyalties{
		Percentage:    sdkmath.NewUint(1000), // 10%
		PayoutAddress: charlie,
	}
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	charlieAddr, err := sdk.AccAddressFromBech32(charlie)
	suite.Require().Nil(err, "error getting charlie address")
	charlieBalance := suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge")
	suite.Require().Equal(charlieBalance.Amount, sdkmath.NewInt(100000000000))

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
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

				ApprovalId: "test",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
					},
					CoinTransfers: []*types.CoinTransfer{
						{
							To: alice,
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(100), Denom: "ubadge"},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "error updating user approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error deducting outgoing approvals")

	charlieBalance = suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge")
	suite.Require().Equal(charlieBalance.Amount, sdkmath.NewInt(100000000000+10)) //10% of 100 ubadge
}

func (suite *TestSuite) TestCannotHaveMoreThanOneUserRoyalties() {
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	wctx := sdk.WrapSDKContext(suite.ctx)
	collectionsToCreate[0].CollectionApprovals[1].OwnershipTimes = GetBottomHalfUintRanges()
	collectionsToCreate[0].CollectionApprovals[1].ApprovalCriteria.UserRoyalties = &types.UserRoyalties{
		Percentage:    sdkmath.NewUint(1000), // 10%
		PayoutAddress: charlie,
	}
	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, GetCollectionsToCreate()[0].CollectionApprovals[0])
	collectionsToCreate[0].CollectionApprovals[2].OwnershipTimes = GetTopHalfUintRanges()
	collectionsToCreate[0].CollectionApprovals[2].ApprovalCriteria.UserRoyalties = &types.UserRoyalties{
		Percentage:    sdkmath.NewUint(2000), // 20%
		PayoutAddress: charlie,
	}
	collectionsToCreate[0].CollectionApprovals[2].ApprovalId = "test2"

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	charlieAddr, err := sdk.AccAddressFromBech32(charlie)
	suite.Require().Nil(err, "error getting charlie address")
	charlieBalance := suite.app.BankKeeper.GetBalance(suite.ctx, charlieAddr, "ubadge")
	suite.Require().Equal(charlieBalance.Amount, sdkmath.NewInt(100000000000))

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
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

				ApprovalId: "test",
				ApprovalCriteria: &types.OutgoingApprovalCriteria{
					MaxNumTransfers: &types.MaxNumTransfers{
						OverallMaxNumTransfers: sdkmath.NewUint(1000),
					},
					ApprovalAmounts: &types.ApprovalAmounts{
						PerFromAddressApprovalAmount: sdkmath.NewUint(1),
					},
					CoinTransfers: []*types.CoinTransfer{
						{
							To: alice,
							Coins: []*sdk.Coin{
								{Amount: sdkmath.NewInt(100), Denom: "ubadge"},
							},
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "error updating user approvals")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.BadgesKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error deducting outgoing approvals")
}
