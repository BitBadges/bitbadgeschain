package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestRevokeBadge() {
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

	err = RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().Nil(err, "Error revoking badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), keeper.GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(0, bobBalanceInfo.IdsForBalances, bobBalanceInfo.Balances))

	aliceBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(0), keeper.GetBadgeBalanceFromIDsAndBalancesForSubbadgeId(0, aliceBalanceInfo.IdsForBalances, aliceBalanceInfo.Balances))
}

func (suite *TestSuite) TestRevokeBadgeTooMuch() {
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

	err = RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated + 1}, []uint64{7000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrBadgeBalanceTooLow.Error())
}

func (suite *TestSuite) TestRevokeBadgeFromSelf() {
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

	accs := suite.app.AccountKeeper.GetAllAccounts(suite.ctx)
	_ = accs

	err = RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrSenderAndReceiverSame.Error())
}

func (suite *TestSuite) TestNewSubBadgeRevokeIsLocked() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  58,
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

	accs := suite.app.AccountKeeper.GetAllAccounts(suite.ctx)
	_ = accs

	err = RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}

func (suite *TestSuite) TestNewSubBadgeNotManager() {
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

	accs := suite.app.AccountKeeper.GetAllAccounts(suite.ctx)
	_ = accs

	err = RevokeBadges(suite, wctx, alice, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, types.NumberRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}
