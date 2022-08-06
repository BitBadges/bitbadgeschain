package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestSelfDestruct() {
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
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	
	
	bobBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	
	

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			StartId: 0,
			EndId:   0,
			Supply:  10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(0, bobBalanceInfo.IdsForBalances, bobBalanceInfo.Balances))

	err = SelfDestructBadge(suite, wctx, bob, 0)
	suite.Require().Nil(err, "Error self destructing badge")

	badge, err = GetBadge(suite, wctx, 0)
	suite.Require().NotNil(err, "We should get a not exists error here now")

	CreateBadges(suite, wctx, badgesToCreate)
	badge, _ = GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err = CreateSubBadges(suite, wctx, bob, 1, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")
	
	err = UpdatePermissions(suite, wctx, bob, 1, 58)
	suite.Require().Nil(err, "Error updating permissions")

	err = SelfDestructBadge(suite, wctx, bob, 1)
	suite.Require().EqualError(err, keeper.ErrBadgeCanNotBeSelfDestructed.Error())
}

func (suite *TestSuite) TestSelfDestructNotManager() {
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
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	
	
	bobBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	
	

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			StartId: 0,
			EndId:   0,
			Supply:  10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(0, bobBalanceInfo.IdsForBalances, bobBalanceInfo.Balances))

	err = SelfDestructBadge(suite, wctx, alice, 0)
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())

}