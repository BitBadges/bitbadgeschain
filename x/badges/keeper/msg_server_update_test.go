package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestUpdateURIs() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  62,
				SubassetUris: validUri,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, 0, "http://google.com", "http://google.com")
	suite.Require().Nil(err, "Error updating uris")
	badge, _ := GetBadge(suite, wctx, 0)
	suite.Require().Equal("http://google.com", badge.Uri)
	suite.Require().Equal("http://google.com", badge.SubassetUriFormat)

	err = UpdatePermissions(suite, wctx, bob, 0, 60)
	suite.Require().Nil(err, "Error updating permissions")
}

func (suite *TestSuite) TestCantUpdate() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  0,
				SubassetUris: validUri,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, bob, 0, "http://google.com", "http://google.com")
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())

	err = UpdatePermissions(suite, wctx, bob, 0, 123)
	suite.Require().EqualError(err, types.ErrInvalidPermissionsUpdateLocked.Error())
}

func (suite *TestSuite) TestCantUpdateNotManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  0,
				SubassetUris: validUri,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	err = UpdateURIs(suite, wctx, alice, 0, "http://google.com", "http://google.com")
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

	err = UpdatePermissions(suite, wctx, alice, 0, 77)
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}
