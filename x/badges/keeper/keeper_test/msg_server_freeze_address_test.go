package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestTransferBadgeForcefulUnfrozenByDefault() {
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

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	
	
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(0, bobBalanceInfo.IdsForBalances, bobBalanceInfo.Balances))

	aliceBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	
	
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(0, aliceBalanceInfo.IdsForBalances, aliceBalanceInfo.Balances))

	err = FreezeAddresses(suite, wctx, bob, types.NumberRange{ Start: firstAccountNumCreated + 1, End: firstAccountNumCreated + 1}, 0, 0, true)
	suite.Require().Nil(err, "Error freezing address")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

func (suite *TestSuite) TestTransferBadgeForcefulFrozenByDefault() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  63,
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

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

	err = FreezeAddresses(suite, wctx, bob, types.NumberRange{ Start: firstAccountNumCreated, End: firstAccountNumCreated}, 0, 0, true)
	suite.Require().Nil(err, "Error unfreezing address")

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().Nil(err, "Error transferring after unfreeze")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

func (suite *TestSuite) TestTransferBadgeForcefulFrozenByDefaultAddAndRemove() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  63,
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

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

	err = FreezeAddresses(suite, wctx, bob, types.NumberRange{ Start: firstAccountNumCreated, End: firstAccountNumCreated}, 0, 0, true)
	suite.Require().Nil(err, "Error unfreezing address")

	err = FreezeAddresses(suite, wctx, bob, types.NumberRange{ Start: firstAccountNumCreated, End: firstAccountNumCreated}, 0, 0, false)
	suite.Require().Nil(err, "Error unfreezing address")

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

	err = FreezeAddresses(suite, wctx, bob, types.NumberRange{ Start: firstAccountNumCreated, End: firstAccountNumCreated}, 0, 0, true)
	suite.Require().Nil(err, "Error unfreezing address")

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated+1}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().Nil(err, "Error transferring after unfreeze")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}
