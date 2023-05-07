package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestTransferManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
				Permissions: 127,
				AllowedTransfers: []*types.TransferMapping{
					{
						From: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create badge 1 with supply > 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	err = RequestTransferManager(suite, wctx, alice, 1, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = TransferManager(suite, wctx, bob, 1, alice)
	suite.Require().Nil(err, "Error transferring manager")

	badge, _ := GetCollection(suite, wctx, 1)
	suite.Require().Equal(alice, badge.Manager)
}

func (suite *TestSuite) TestRequestTransferManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
				Permissions: 127,
				AllowedTransfers: []*types.TransferMapping{
					{
						From: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create badge 1 with supply > 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	err = RequestTransferManager(suite, wctx, alice, 1, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = RequestTransferManager(suite, wctx, alice, 1, false)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = RequestTransferManager(suite, wctx, alice, 1, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = TransferManager(suite, wctx, bob, 1, alice)
	suite.Require().Nil(err, "Error transferring manager")

	badge, _ := GetCollection(suite, wctx, 1)
	suite.Require().Equal(alice, badge.Manager)
}

func (suite *TestSuite) TestRemovedRequestTransferManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
				Permissions: 127,
				AllowedTransfers: []*types.TransferMapping{
					{
						From: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create badge 1 with supply > 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	err = RequestTransferManager(suite, wctx, alice, 1, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = RequestTransferManager(suite, wctx, alice, 1, false)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = TransferManager(suite, wctx, bob, 1, alice)
	suite.Require().EqualError(err, keeper.ErrAddressNeedsToOptInAndRequestManagerTransfer.Error())
}

func (suite *TestSuite) TestRemovedRequestTransferManagerBadPermissions() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
				Permissions: 23,
				AllowedTransfers: []*types.TransferMapping{
					{
						From: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
						To: &types.AddressesMapping{
							IncludeOnlySpecified: false,
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create badge 1 with supply > 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	err = RequestTransferManager(suite, wctx, alice, 1, true)
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}

func (suite *TestSuite) TestManagerCantBeTransferred() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
				Permissions: 0,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	err = TransferManager(suite, wctx, bob, 1, alice)
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}
