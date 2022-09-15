package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestTransferManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
				Permissions: 127,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create subbadge 1 with supply > 1
	err = CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = TransferManager(suite, wctx, bob, 0, aliceAccountNum)
	suite.Require().Nil(err, "Error transferring manager")

	badge, _ := GetBadge(suite, wctx, 0)
	suite.Require().Equal(aliceAccountNum, badge.Manager)
}

func (suite *TestSuite) TestRequestTransferManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
				Permissions: 127,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create subbadge 1 with supply > 1
	err = CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = RequestTransferManager(suite, wctx, alice, 0, false)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = TransferManager(suite, wctx, bob, 0, aliceAccountNum)
	suite.Require().Nil(err, "Error transferring manager")

	badge, _ := GetBadge(suite, wctx, 0)
	suite.Require().Equal(aliceAccountNum, badge.Manager)
}

func (suite *TestSuite) TestRemovedRequestTransferManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
				Permissions: 127,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create subbadge 1 with supply > 1
	err = CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = RequestTransferManager(suite, wctx, alice, 0, false)
	suite.Require().Nil(err, "Error requesting manager transfer")

	err = TransferManager(suite, wctx, bob, 0, aliceAccountNum)
	suite.Require().EqualError(err, keeper.ErrAddressNeedsToOptInAndRequestManagerTransfer.Error())
}

func (suite *TestSuite) TestRemovedRequestTransferManagerBadPermissions() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
				Permissions: 62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge")

	//Create subbadge 1 with supply > 1
	err = CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	err = RequestTransferManager(suite, wctx, alice, 0, true)
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}

func (suite *TestSuite) TestManagerCantBeTransferred() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge")

	err = TransferManager(suite, wctx, bob, 0, aliceAccountNum)
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}
