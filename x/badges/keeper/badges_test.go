package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestGetBadge() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
				Permissions: 62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	badge, err := suite.app.BadgesKeeper.GetBadgeE(suite.ctx, 0)
	suite.Require().Nil(err, "Error getting badge: %s")
	suite.Require().Equal(badge.Id, uint64(0))

	badge, err = suite.app.BadgesKeeper.GetBadgeE(suite.ctx, 1)
	suite.Require().EqualError(err, keeper.ErrBadgeNotExists.Error())
}

func (suite *TestSuite) TestGetBadgeAndAssertSubbadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
				Permissions:             62,
				SubassetSupplys:         []uint64{1},
				SubassetAmountsToCreate: []uint64{1},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	_, err = suite.app.BadgesKeeper.GetBadgeAndAssertSubbadgeRangesAreValid(suite.ctx, 0, []*types.IdRange{
		{
			Start: 0,
			End:   0,
		},
	})
	suite.Require().Nil(err, "Error getting badge: %s")

	_, err = suite.app.BadgesKeeper.GetBadgeAndAssertSubbadgeRangesAreValid(suite.ctx, 0, []*types.IdRange{
		{
			Start: 20,
			End:   10,
		},
	})
	suite.Require().EqualError(err, keeper.ErrInvalidSubbadgeRange.Error())

	_, err = suite.app.BadgesKeeper.GetBadgeAndAssertSubbadgeRangesAreValid(suite.ctx, 0, []*types.IdRange{
		{
			Start: 0,
			End:   10,
		},
	})
	suite.Require().EqualError(err, keeper.ErrSubBadgeNotExists.Error())
}

func (suite *TestSuite) TestCreateSubassets() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
				Permissions: 62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	badge, err := GetBadge(suite, wctx, 0)
	suite.Require().Nil(err, "Error getting badge: %s")
	balanceInfo := types.UserBalanceInfo{}

	badge, balanceInfo, err = keeper.CreateSubassets(badge, balanceInfo, []uint64{1}, []uint64{1})
	suite.Require().Nil(err, "Error creating subassets: %s")
	suite.Require().Equal(badge.SubassetSupplys, []*types.BalanceObject(nil))
	suite.Require().Equal(balanceInfo.BalanceAmounts[0].Balance, uint64(1))
	suite.Require().Equal(balanceInfo.BalanceAmounts[0].IdRanges, []*types.IdRange{
		{
			Start: 0,
			End:   0,
		},
	})

	badge, balanceInfo, err = keeper.CreateSubassets(badge, balanceInfo, []uint64{1}, []uint64{1})
	suite.Require().Nil(err, "Error creating subassets: %s")
	suite.Require().Equal(badge.SubassetSupplys, []*types.BalanceObject(nil))
	suite.Require().Equal(balanceInfo.BalanceAmounts[0].Balance, uint64(1))
	suite.Require().Equal(balanceInfo.BalanceAmounts[0].IdRanges, []*types.IdRange{
		{
			Start: 0,
			End:   1,
		},
	})

	badge, balanceInfo, err = keeper.CreateSubassets(badge, balanceInfo, []uint64{0}, []uint64{1})
	suite.Require().Nil(err, "Error creating subassets: %s")
	suite.Require().Equal(badge.SubassetSupplys, []*types.BalanceObject(nil))
	suite.Require().Equal(balanceInfo.BalanceAmounts[0].Balance, uint64(1))
	suite.Require().Equal(balanceInfo.BalanceAmounts[0].IdRanges, []*types.IdRange{
		{
			Start: 0,
			End:   2,
		},
	})

	badge, balanceInfo, err = keeper.CreateSubassets(badge, balanceInfo, []uint64{1}, []uint64{math.MaxUint64})
	suite.Require().EqualError(err, keeper.ErrOverflow.Error())
}
