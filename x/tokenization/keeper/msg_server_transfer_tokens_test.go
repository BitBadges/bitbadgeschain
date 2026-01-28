package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestTransferTokensForceful() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
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
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token")

	bobbalance, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	fetchedBalance, err = types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), bobbalance.Balances)
	AssertUintsEqual(suite, sdkmath.NewUint(0), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	alicebalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	fetchedBalance, err = types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), alicebalance.Balances)
	AssertUintsEqual(suite, sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)
}

func (suite *TestSuite) TestTransferTokensHandleDuplicateIDs() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						TokenIds: []*types.UintRange{
							GetOneUintRange()[0],
							GetOneUintRange()[0],
						},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token")

}

func (suite *TestSuite) TestTransferTokensNotApprovedCollectionLevel() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = UpdateCollectionApprovals(suite, wctx, &types.MsgUniversalUpdateCollectionApprovals{
		Creator:             bob,
		CollectionId:        sdkmath.NewUint(1),
		CollectionApprovals: []*types.CollectionApproval{},
	})
	suite.Require().Nil(err, "Error updating approved transfers")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
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
			},
		},
	})
	suite.Require().Error(err, "Error transferring token")
}

func (suite *TestSuite) TestTransferTokensNotApprovedIncoming() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = UpdateUserApprovals(suite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                 alice,
		CollectionId:            sdkmath.NewUint(1),
		UpdateIncomingApprovals: true,
		IncomingApprovals:       []*types.UserIncomingApproval{},
	})
	suite.Require().Nil(err, "Error updating approved transfers")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
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
			},
		},
	})
	suite.Require().Error(err, "Error transferring token")
}

func (suite *TestSuite) TestIncrementsWithAttemptToTransferAll() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	fetchedBalance, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	unmintedSupplys, err := GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{}, unmintedSupplys.Balances)

	err = MintAndDistributeTokens(suite, wctx, &types.MsgMintAndDistributeTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		TokensToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		CollectionApprovals: collection.CollectionApprovals,
		// InheritedBalancesTimeline: 				 collection.InheritedBalancesTimeline,
		// Note: MsgMintAndDistributeTokens is a legacy type that still uses timeline fields
		// We'll just not set these fields for now since they're optional
	})
	suite.Require().Nil(err, "Error transferring token")

	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       GetOneUintRange(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring token")

	unmintedSupplys, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{}, unmintedSupplys.Balances)

	// Mint has unlimited balances now
	err = TransferTokens(suite, wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(1000),
						TokenIds:       GetFullUintRanges(),
						OwnershipTimes: GetFullUintRanges(),
					},
				},
				PrioritizedApprovals: GetDefaultPrioritizedApprovals(suite.ctx, suite.app.TokenizationKeeper, sdkmath.NewUint(1)),
			},
		},
	})
	suite.Require().Error(err, "Error transferring token")
}
