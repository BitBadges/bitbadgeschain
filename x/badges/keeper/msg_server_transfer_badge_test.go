package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestTransferBadgeForceful() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(GetOneIdRange(), GetOneIdRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = TransferBadge(suite, wctx, &types.MsgTransferBadge{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:         bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetOneIdRange(),
						Times: GetFullIdRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge")

	bobbalance, _ = GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	fetchedBalance, err = types.GetBalancesForIds(GetOneIdRange(), GetOneIdRange(), bobbalance.Balances)
	AssertUintsEqual(suite, sdkmath.NewUint(0), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	alicebalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), alice)
	fetchedBalance, err = types.GetBalancesForIds(GetOneIdRange(), GetOneIdRange(), alicebalance.Balances)
	AssertUintsEqual(suite, sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)
}


func (suite *TestSuite) TestTransferBadgeHandleDuplicateIDs() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(GetOneIdRange(), GetOneIdRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = TransferBadge(suite, wctx, &types.MsgTransferBadge{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:         bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: []*types.IdRange{
							GetOneIdRange()[0],
							GetOneIdRange()[0],
						},
						Times: GetFullIdRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge")

}

func (suite *TestSuite) TestTransferBadgeNotApprovedCollectionLevel() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(GetOneIdRange(), GetOneIdRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = UpdateCollectionApprovedTransfers(suite, wctx, &types.MsgUpdateCollectionApprovedTransfers{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		ApprovedTransfersTimeline: []*types.CollectionApprovedTransferTimeline{
			{
				ApprovedTransfers: []*types.CollectionApprovedTransfer{},
				Times: GetFullIdRanges(),
			},
		},
	})
	suite.Require().Nil(err, "Error updating approved transfers")

	err = TransferBadge(suite, wctx, &types.MsgTransferBadge{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:         bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetOneIdRange(),
						Times: GetFullIdRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge")
}


func (suite *TestSuite) TestTransferBadgeNotApprovedIncoming() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)

	fetchedBalance, err := types.GetBalancesForIds(GetOneIdRange(), GetOneIdRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	err = UpdateUserApprovedTransfers(suite, wctx, &types.MsgUpdateUserApprovedTransfers{
		Creator:      alice,
		CollectionId: sdkmath.NewUint(1),
		ApprovedIncomingTransfersTimeline: []*types.UserApprovedIncomingTransferTimeline{
			{
				// ApprovedIncomingTransfers: []*types.CollectionApprovedTransfer{},
				Times: GetFullIdRanges(),
			},
		},
	})
	suite.Require().Nil(err, "Error updating approved transfers")

	err = TransferBadge(suite, wctx, &types.MsgTransferBadge{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:         bob,
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetOneIdRange(),
						Times: GetFullIdRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge")
}

func (suite *TestSuite) TestTransferBadgeFromMintAddress() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating collections")

	bobbalance, _ := GetUserBalance(suite, wctx, sdkmath.NewUint(1), bob)
	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	fetchedBalance, err := types.GetBalancesForIds(GetOneIdRange(), GetOneIdRange(), bobbalance.Balances)
	suite.Require().Equal(sdkmath.NewUint(1), fetchedBalance[0].Amount)
	suite.Require().Nil(err)

	AssertBalancesEqual(suite, []*types.Balance{}, collection.UnmintedSupplys)

	err = MintAndDistributeBadges(suite, wctx, &types.MsgMintAndDistributeBadges{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		BadgesToCreate: []*types.Balance{
			{
				Amount: sdkmath.NewUint(1),
				BadgeIds: GetOneIdRange(),
				Times: GetFullIdRanges(),
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, &types.MsgTransferBadge{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:         "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetOneIdRange(),
						Times: GetFullIdRanges(),
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transferring badge")

	AssertBalancesEqual(suite, []*types.Balance{}, collection.UnmintedSupplys)

	err = TransferBadge(suite, wctx, &types.MsgTransferBadge{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{
			{
				From:         "Mint",
				ToAddresses: []string{alice},
				Balances: []*types.Balance{
					{
						Amount: sdkmath.NewUint(1),
						BadgeIds: GetOneIdRange(),
						Times: GetFullIdRanges(),
					},
				},
			},
		},
	})
	suite.Require().Error(err, "Error transferring badge")
}