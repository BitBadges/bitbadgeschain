package keeper_test

import (
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
				BadgeUri:      "https://example.com/{id}",
				Permissions:   127,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create badge 1 with supply > 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = TransferManager(suite, wctx, bob, 0, aliceAccountNum)
	suite.Require().Nil(err, "Error transferring manager")

	badge, _ := GetCollection(suite, wctx, 0)
	suite.Require().Equal(aliceAccountNum, badge.Manager)
}

func (suite *TestSuite) TestRequestTransferManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   127,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create badge 1 with supply > 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = RequestTransferManager(suite, wctx, alice, 0, false)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = TransferManager(suite, wctx, bob, 0, aliceAccountNum)
	suite.Require().Nil(err, "Error transferring manager")

	badge, _ := GetCollection(suite, wctx, 0)
	suite.Require().Equal(aliceAccountNum, badge.Manager)
}

func (suite *TestSuite) TestRemovedRequestTransferManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   127,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create badge 1 with supply > 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = RequestTransferManager(suite, wctx, alice, 0, false)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = TransferManager(suite, wctx, bob, 0, aliceAccountNum)
	suite.Require().EqualError(err, keeper.ErrAddressNeedsToOptInAndRequestManagerTransfer.Error())
}

func (suite *TestSuite) TestRemovedRequestTransferManagerBadPermissions() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   23,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create badge 1 with supply > 1
	err = CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}

func (suite *TestSuite) TestManagerCantBeTransferred() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   0,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge")

	err = TransferManager(suite, wctx, bob, 0, aliceAccountNum)
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}
