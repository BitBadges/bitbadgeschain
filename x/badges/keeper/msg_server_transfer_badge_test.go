package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestTransferBadgesForceful() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
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
	suite.Require().Nil(err, "Error transferring badge")

	bobbalance, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	fetchedBalance, err = types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), bobbalance.Balances)
	AssertUintsEqual(suite, sdkmath.NewUint(0), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	alicebalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	fetchedBalance, err = types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), alicebalance.Balances)
	AssertUintsEqual(suite, sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)
}

func (suite *TestSuite) TestTransferBadgesHandleDuplicateIDs() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(suite.ctx, GetOneUintRange(), GetOneUintRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: []*types.UintRange{
							GetOneUintRange()[0],
							GetOneUintRange()[0],
						},
						OwnershipTimes: GetFullUintRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge")

}

func (suite *TestSuite) TestTransferBadgesNotApprovedCollectionLevel() {
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

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
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
	suite.Require().Error(err, "Error transferring badge")
}

func (suite *TestSuite) TestTransferBadgesNotApprovedIncoming() {
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

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        bob,
				ToAddresses: []string{alice},
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
	suite.Require().Error(err, "Error transferring badge")
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

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount:         sdkmath.NewUint(1),
				BadgeIds:       GetOneUintRange(),
				OwnershipTimes: GetFullUintRanges(),
			},
		},
		CollectionApprovals: collection.CollectionApprovals,
		// InheritedBalancesTimeline: 				 collection.InheritedBalancesTimeline,
		CollectionMetadataTimeline:       collection.CollectionMetadataTimeline,
		BadgeMetadataTimeline:            collection.BadgeMetadataTimeline,
		OffChainBalancesMetadataTimeline: collection.OffChainBalancesMetadataTimeline,
	})
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
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
	suite.Require().Nil(err, "Error transferring badge")

	unmintedSupplys, err = GetUserBalance(suite, wctx, sdkmath.NewUint(1), "Mint")
	suite.Require().Nil(err, "Error getting user balance: %s")
	AssertBalancesEqual(suite, []*types.Balance{}, unmintedSupplys.Balances)

	err = TransferBadges(suite, wctx, &types.MsgTransferBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:        "Mint",
				ToAddresses: []string{alice},
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
	suite.Require().Error(err, "Error transferring badge")
}
