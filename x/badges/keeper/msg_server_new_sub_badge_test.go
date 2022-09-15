package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNewSubBadges() {
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

	CreateBadges(suite, wctx, badgesToCreate)
	badge, _ := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(10), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	//Create subbadge 2 with supply == 1
	err = CreateSubBadges(suite, wctx, bob, 0, []uint64{1}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(2), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	}, badge.SubassetSupplys)
	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(1), bobBalanceInfo.BalanceAmounts[0].Balance)
	suite.Require().Equal(uint64(1), bobBalanceInfo.BalanceAmounts[0].IdRanges[0].Start)

	//Create subbadge 2 with supply == 10
	err = CreateSubBadges(suite, wctx, bob, 0, []uint64{10}, []uint64{2})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(4), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}, {Start: 2, End: 3}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(10), bobBalanceInfo.BalanceAmounts[1].Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.BalanceAmounts[1].IdRanges[0].Start)
	suite.Require().Equal(uint64(0), bobBalanceInfo.BalanceAmounts[1].IdRanges[0].End)
	suite.Require().Equal(uint64(2), bobBalanceInfo.BalanceAmounts[1].IdRanges[1].Start)
	suite.Require().Equal(uint64(3), bobBalanceInfo.BalanceAmounts[1].IdRanges[1].End)
}

func (suite *TestSuite) TestNewSubbadgesDirectlyUponCreatingNewBadge() {
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

				SubassetSupplys:         []uint64{10},
				SubassetAmountsToCreate: []uint64{1},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)
	badge, _ := GetBadge(suite, wctx, 0)

	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(10), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	//Create subbadge 2 with supply == 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{1}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(2), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(1), bobBalanceInfo.BalanceAmounts[0].Balance)
	suite.Require().Equal(uint64(1), bobBalanceInfo.BalanceAmounts[0].IdRanges[0].Start)

	//Create subbadge 2 with supply == 10
	err = CreateSubBadges(suite, wctx, bob, 0, []uint64{10}, []uint64{2})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(4), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}, {Start: 2, End: 3}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	},
		badge.SubassetSupplys)
	suite.Require().Equal(uint64(10), bobBalanceInfo.BalanceAmounts[1].Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.BalanceAmounts[1].IdRanges[0].Start)
	suite.Require().Equal(uint64(0), bobBalanceInfo.BalanceAmounts[1].IdRanges[0].End)
	suite.Require().Equal(uint64(2), bobBalanceInfo.BalanceAmounts[1].IdRanges[1].Start)
	suite.Require().Equal(uint64(3), bobBalanceInfo.BalanceAmounts[1].IdRanges[1].End)
}

func (suite *TestSuite) TestNewSubBadgesNotManager() {
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

	CreateBadges(suite, wctx, badgesToCreate)
	err := CreateSubBadges(suite, wctx, alice, 0, []uint64{10}, []uint64{1})
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}

func (suite *TestSuite) TestNewSubBadgeBadgeNotExists() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	err := CreateSubBadges(suite, wctx, alice, 0, []uint64{10}, []uint64{1})
	suite.Require().EqualError(err, keeper.ErrBadgeNotExists.Error())
}

func (suite *TestSuite) TestNewSubBadgeCreateIsLocked() {
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

	CreateBadges(suite, wctx, badgesToCreate)
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10}, []uint64{1})
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}
