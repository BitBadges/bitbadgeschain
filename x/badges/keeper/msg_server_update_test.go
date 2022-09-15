package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestUpdateURIs() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: &types.UriObject{
					Uri:                    "example.com/",
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					InsertIdIdx:            10,
				},
				Permissions: 62 + 128,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, 0, &types.UriObject{
		Uri:                    "example.com/",
		Scheme:                 1,
		IdxRangeToRemove:       &types.IdRange{},
		InsertSubassetBytesIdx: 0,

		InsertIdIdx: 10,
	})
	suite.Require().Nil(err, "Error updating uris")
	badge, _ := GetBadge(suite, wctx, 0)
	suite.Require().Equal(&types.UriObject{
		Uri:                    "example.com/",
		Scheme:                 1,
		IdxRangeToRemove:       &types.IdRange{},
		InsertSubassetBytesIdx: 0,

		InsertIdIdx: 10,
	}, badge.Uri)

	err = UpdatePermissions(suite, wctx, bob, 0, 60+128)
	suite.Require().Nil(err, "Error updating permissions")

	err = UpdateBytes(suite, wctx, bob, 0, "example.com/")
	suite.Require().Nil(err, "Error updating permissions")
}

func (suite *TestSuite) TestCantUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: &types.UriObject{
					Uri:                    "example.com/",
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,

					InsertIdIdx: 10,
				},
				Permissions: 0,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, 0, &types.UriObject{
		Uri:                    "example.com/",
		Scheme:                 1,
		IdxRangeToRemove:       &types.IdRange{},
		InsertSubassetBytesIdx: 0,

		InsertIdIdx: 10,
	})
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

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: &types.UriObject{
					Uri:                    "example.com/",
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,

					InsertIdIdx: 10,
				},
				Permissions: 0,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, alice, 0, &types.UriObject{
		Uri:                    "example.com/",
		Scheme:                 1,
		IdxRangeToRemove:       &types.IdRange{},
		InsertSubassetBytesIdx: 0,

		InsertIdIdx: 10,
	})
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

	err = UpdatePermissions(suite, wctx, alice, 0, 77)
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

	err = UpdateBytes(suite, wctx, alice, 0, "example.com/")
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}
