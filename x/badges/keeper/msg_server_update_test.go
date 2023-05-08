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
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				CollectionUri: "https://example.com",
				Permissions: sdk.NewUint(62),
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, sdk.NewUint(1), "https://example.com", []*types.BadgeUri{
		{
			Uri: "https://example.com/{id}",
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(1),
					End:   sdk.NewUint(math.MaxUint64),
				},
			},
		},
	}, "")
	suite.Require().Nil(err, "Error updating uris")
	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	suite.Require().Equal("https://example.com", badge.CollectionUri)
	// suite.Require().Equal("https://example.com/{id}", badge.BadgeUri)

	err = UpdatePermissions(suite, wctx, bob, sdk.NewUint(1), sdk.NewUint(60-32))
	suite.Require().Nil(err, "Error updating permissions")

	err = UpdateBytes(suite, wctx, bob, sdk.NewUint(1), "example.com/")
	suite.Require().Nil(err, "Error updating bytes")
}

func (suite *TestSuite) TestCantUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(0),
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, sdk.NewUint(1), "https://example.com/test2222", []*types.BadgeUri{
		{
			Uri: "https://example.com/{id}/edited",
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(1),
					End:   sdk.NewUint(math.MaxUint64),
				},
			},
		},
	}, "")
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())

	err = UpdatePermissions(suite, wctx, bob, sdk.NewUint(1), sdk.NewUint(123-64))
	suite.Require().EqualError(err, types.ErrInvalidPermissionsUpdateLocked.Error())

	err = UpdateBytes(suite, wctx, bob, sdk.NewUint(1), "example.com/")
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
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(0),
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, alice, sdk.NewUint(1), "https://example.com", []*types.BadgeUri{
		{
			Uri: "https://example.com/{id}",
			BadgeIds: []*types.IdRange{
				{
					Start: sdk.NewUint(1),
					End:   sdk.NewUint(math.MaxUint64),
				},
			},
		},
	}, "")
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

	err = UpdatePermissions(suite, wctx, alice, sdk.NewUint(1), sdk.NewUint(77))
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

	err = UpdateBytes(suite, wctx, alice, sdk.NewUint(1), "example.com/")
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}


func (suite *TestSuite) TestUpdateBalanceURIs() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				CollectionUri: "https://example.com",
				Permissions: sdk.NewUint(62 + 64),
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, sdk.NewUint(1), "", []*types.BadgeUri{	}, "https://balance.com/{id}")
	suite.Require().Nil(err, "Error updating uris")
	badge, _ := GetCollection(suite, wctx, sdk.NewUint(1))
	suite.Require().Equal("https://balance.com/{id}", badge.BalancesUri)
	// suite.Require().Equal("https://example.com/{id}", badge.BadgeUri)
}

func (suite *TestSuite) TestCantUpdateBalanceUris() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: sdk.NewUint(1),
								End:   sdk.NewUint(math.MaxUint64),
							},
						},
					},
				},
				Permissions: sdk.NewUint(0),
			},
			Amount:  sdk.NewUint(1),
			Creator: bob,
		},
	}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, sdk.NewUint(1), "", []*types.BadgeUri{}, "https://balance.com/{id}")
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}