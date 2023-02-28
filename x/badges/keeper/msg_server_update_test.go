package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestUpdateURIs() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				CollectionUri: "https://example.com",
				Permissions:   62 + 128,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, 0, "https://example.com", []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},)
	suite.Require().Nil(err, "Error updating uris")
	badge, _ := GetCollection(suite, wctx, 0)
	suite.Require().Equal("https://example.com", badge.CollectionUri)
	// suite.Require().Equal("https://example.com/{id}", badge.BadgeUri)

	err = UpdatePermissions(suite, wctx, bob, 0, 60+128)
	suite.Require().Nil(err, "Error updating permissions")

	err = UpdateBytes(suite, wctx, bob, 0, "example.com/")
	suite.Require().Nil(err, "Error updating permissions")
}

func (suite *TestSuite) TestCantUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   0,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, 0, "https://example.com", []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},)
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())

	err = UpdatePermissions(suite, wctx, bob, 0, 123)
	suite.Require().EqualError(err, types.ErrInvalidPermissionsUpdateLocked.Error())

	err = UpdateBytes(suite, wctx, bob, 0, "example.com/")
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}

func (suite *TestSuite) TestCantUpdateNotManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris:            []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},
				Permissions:   0,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, alice, 0, "https://example.com", []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End: math.MaxUint64,
							},
						},
					},
				},)
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

	err = UpdatePermissions(suite, wctx, alice, 0, 77)
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

	err = UpdateBytes(suite, wctx, alice, 0, "example.com/")
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}
