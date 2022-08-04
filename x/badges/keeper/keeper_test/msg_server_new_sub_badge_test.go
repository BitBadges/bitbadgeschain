package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestNewSubBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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

	CreateBadges(suite, wctx, badgesToCreate)
	badge, _ := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			StartId: 0,
			EndId:   0,
			Supply:  10,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10), keeper.GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(0, bobBalanceInfo.IdsForBalances, bobBalanceInfo.Balances))

	//Create subbadge 2 with supply == 1
	err = CreateSubBadges(suite, wctx, bob, 0, []uint64{1}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 1, firstAccountNumCreated)

	suite.Require().Equal(uint64(2), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			StartId: 0,
			EndId:   0,
			Supply:  10,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(1), keeper.GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(1, bobBalanceInfo.IdsForBalances, bobBalanceInfo.Balances))

	//Create subbadge 2 with supply == 10
	err = CreateSubBadges(suite, wctx, bob, 0, []uint64{10}, []uint64{2})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 2, firstAccountNumCreated)

	suite.Require().Equal(uint64(4), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			StartId: 0,
			EndId:   0,
			Supply:  10,
		},
		{
			StartId: 2,
			EndId:   3,
			Supply:  10,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10), keeper.GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(2, bobBalanceInfo.IdsForBalances, bobBalanceInfo.Balances))
}

func (suite *TestSuite) TestNewSubBadgesNotManager() {
	wctx := sdk.WrapSDKContext(suite.ctx)

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
				Uri:          validUri,
				Permissions:  0,
				SubassetUris: validUri,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10}, []uint64{1})
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}
