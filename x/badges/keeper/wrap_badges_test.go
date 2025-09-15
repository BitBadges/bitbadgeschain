package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestWrapBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "test-coin",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					BadgeIds:       GetOneUintRange(),
				},
			},
			AllowCosmosWrapping: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "asadsdas",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesToIncomingApprovals:   true,
			OverridesFromOutgoingApprovals: true,
		},
	})
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	bobBalanceBefore, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance")
	suite.Require().Equal(sdkmath.NewUint(1), bobBalanceBefore.Balances[0].Amount, "Error creating tokens")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	denomAddress := collection.CosmosCoinWrapperPaths[0].Address

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{denomAddress},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error wrapping tokens")

	//1. ensure tokens were burned
	bobBalanceAfter, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance")

	diffInBalances, err := types.SubtractBalances(suite.ctx, bobBalanceAfter.Balances, bobBalanceBefore.Balances)
	suite.Require().Nil(err, "Error subtracting balances")

	// len 1, amount 1, badgeIds 1, full ownership times
	suite.Require().Equal(1, len(diffInBalances), "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(1), diffInBalances[0].Amount, "Error burning tokens")
	suite.Require().Equal(1, len(diffInBalances[0].BadgeIds), "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(1), diffInBalances[0].BadgeIds[0].Start, "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(1), diffInBalances[0].BadgeIds[0].End, "Error burning tokens")
	suite.Require().Equal(GetFullUintRanges(), diffInBalances[0].OwnershipTimes, "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(1), diffInBalances[0].OwnershipTimes[0].Start, "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(18446744073709551615), diffInBalances[0].OwnershipTimes[0].End, "Error burning tokens")

	// //2. ensure tokens were wrapped
	collection, err = GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")

	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting user address")
	fullDenom := generateBadgeDenom(collection.CollectionId, collection.CosmosCoinWrapperPaths[0])

	bobBalanceDenom := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, fullDenom)
	bobAmount := sdkmath.NewUintFromBigInt(bobBalanceDenom.Amount.BigInt())
	suite.Require().Equal(sdkmath.NewUint(1), bobAmount, "Error wrapping tokens")

	// Unwrap the tokens
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        denomAddress,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         bobAmount,
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error unwrapping tokens")

	// Ensure tokens were unwrapped
	bobBalanceAfterUnwrap, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance")
	suite.Require().Equal(bobBalanceBefore.Balances, bobBalanceAfterUnwrap.Balances, "Error unwrapping tokens")

	// Ensure the denom was burned
	bobBalanceDenomAfterUnwrap := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, fullDenom)
	suite.Require().Equal(sdkmath.NewInt(0), bobBalanceDenomAfterUnwrap.Amount, "Error unwrapping tokens")
}

func (suite *TestSuite) TestWrapBadgesErrors() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "test-coin",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					BadgeIds:       GetOneUintRange(),
				},
			},
			AllowCosmosWrapping: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "asadsdas",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesToIncomingApprovals:   true,
			OverridesFromOutgoingApprovals: true,
		},
	})
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	bobBalanceBefore, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance")
	suite.Require().Equal(sdkmath.NewUint(1), bobBalanceBefore.Balances[0].Amount, "Error creating tokens")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	denomAddress := collection.CosmosCoinWrapperPaths[0].Address

	balances := []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1),
			BadgeIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	// Test more than one balance
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{denomAddress},
				Balances: append(balances, &types.Balance{
					Amount:         sdkmath.NewUint(1),
					BadgeIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				}),
			},
		},
	})
	suite.Require().Error(err, "Error wrapping tokens")

	// Test wrong token IDs
	newBalancesClone := make([]*types.Balance, len(balances))
	copy(newBalancesClone, balances)
	newBalancesClone[0].BadgeIds[0].Start = sdkmath.NewUint(2)
	newBalancesClone[0].BadgeIds[0].End = sdkmath.NewUint(2)

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{denomAddress},
				Balances:    newBalancesClone,
			},
		},
	})
	suite.Require().Error(err, "Error wrapping tokens")

	// Test wrong ownership times
	newBalancesClone[0].OwnershipTimes[0].Start = sdkmath.NewUint(2)
	newBalancesClone[0].OwnershipTimes[0].End = sdkmath.NewUint(2)

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{denomAddress},
				Balances:    newBalancesClone,
			},
		},
	})
	suite.Require().Error(err, "Error wrapping tokens")
}

func (suite *TestSuite) TestWrapBadgesInadequateBalanceOnTheUnwrap() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].CosmosCoinWrapperPathsToAdd = []*types.CosmosCoinWrapperPathAddObject{
		{
			Denom: "test-coin",
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1),
					OwnershipTimes: GetFullUintRanges(),
					BadgeIds:       GetOneUintRange(),
				},
			},
			AllowCosmosWrapping: true,
		},
	}

	collectionsToCreate[0].CollectionApprovals = append(collectionsToCreate[0].CollectionApprovals, &types.CollectionApproval{
		ApprovalId:        "asadsdas",
		TransferTimes:     GetFullUintRanges(),
		OwnershipTimes:    GetFullUintRanges(),
		BadgeIds:          GetOneUintRange(),
		FromListId:        "AllWithoutMint",
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesToIncomingApprovals:   true,
			OverridesFromOutgoingApprovals: true,
		},
	})
	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "error creating tokens")

	bobBalanceBefore, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance")
	suite.Require().Equal(sdkmath.NewUint(1), bobBalanceBefore.Balances[0].Amount, "Error creating tokens")

	collection, err := GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")
	denomAddress := collection.CosmosCoinWrapperPaths[0].Address

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{denomAddress},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error wrapping tokens")

	//1. ensure tokens were burned
	bobBalanceAfter, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().Nil(err, "Error getting user balance")

	diffInBalances, err := types.SubtractBalances(suite.ctx, bobBalanceAfter.Balances, bobBalanceBefore.Balances)
	suite.Require().Nil(err, "Error subtracting balances")

	// len 1, amount 1, badgeIds 1, full ownership times
	suite.Require().Equal(1, len(diffInBalances), "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(1), diffInBalances[0].Amount, "Error burning tokens")
	suite.Require().Equal(1, len(diffInBalances[0].BadgeIds), "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(1), diffInBalances[0].BadgeIds[0].Start, "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(1), diffInBalances[0].BadgeIds[0].End, "Error burning tokens")
	suite.Require().Equal(GetFullUintRanges(), diffInBalances[0].OwnershipTimes, "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(1), diffInBalances[0].OwnershipTimes[0].Start, "Error burning tokens")
	suite.Require().Equal(sdkmath.NewUint(18446744073709551615), diffInBalances[0].OwnershipTimes[0].End, "Error burning tokens")

	// //2. ensure tokens were wrapped
	collection, err = GetCollection(suite, wctx, sdkmath.NewUint(1))
	suite.Require().Nil(err, "Error getting collection")

	bobAccAddr, err := sdk.AccAddressFromBech32(bob)
	suite.Require().Nil(err, "Error getting user address")
	fullDenom := generateBadgeDenom(collection.CollectionId, collection.CosmosCoinWrapperPaths[0])

	bobBalanceDenom := suite.app.BankKeeper.GetBalance(suite.ctx, bobAccAddr, fullDenom)
	bobAmount := sdkmath.NewUintFromBigInt(bobBalanceDenom.Amount.BigInt())
	suite.Require().Equal(sdkmath.NewUint(1), bobAmount, "Error wrapping tokens")

	// Transfer some of the balance to alice
	aliceAccAddr, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Error getting user address")
	err = suite.app.BankKeeper.SendCoins(suite.ctx, bobAccAddr, aliceAccAddr, sdk.Coins{sdk.NewCoin(fullDenom, sdkmath.NewIntFromBigInt(bobAmount.BigInt()))})
	suite.Require().Nil(err, "Error sending coins")

	// Unwrap the tokens - bob should fail
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        denomAddress,
				ToAddresses: []string{bob},
				Balances: []*types.Balance{
					{
						Amount:         bobAmount,
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error unwrapping tokens")

	// Unwrap the tokens - alice should succeed
	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        denomAddress,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         bobAmount,
						BadgeIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error unwrapping tokens")
}
