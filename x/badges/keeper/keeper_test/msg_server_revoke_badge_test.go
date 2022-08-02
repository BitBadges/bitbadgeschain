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
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id:     0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated+1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Balance)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Balance)

	err = RevokeBadge(suite, wctx, bob, firstAccountNumCreated+1, 5000, 0, 0)
	suite.Require().Nil(err, "Error revoking badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Balance)
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
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id:     0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated+1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Balance)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Balance)

	err = RevokeBadge(suite, wctx, bob, firstAccountNumCreated+1, 7000, 0, 0)
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
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id:     0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated+1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Balance)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Balance)

	accs := suite.app.AccountKeeper.GetAllAccounts(suite.ctx)
	_ = accs

	err = RevokeBadge(suite, wctx, bob, firstAccountNumCreated, 5000, 0, 0)
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
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id:     0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated+1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Balance)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Balance)

	accs := suite.app.AccountKeeper.GetAllAccounts(suite.ctx)
	_ = accs

	err = RevokeBadge(suite, wctx, bob, firstAccountNumCreated+1, 5000, 0, 0)
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
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id:     0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated+1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Balance)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Balance)

	accs := suite.app.AccountKeeper.GetAllAccounts(suite.ctx)
	_ = accs

	err = RevokeBadge(suite, wctx, alice, firstAccountNumCreated, 5000, 0, 0)
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}
