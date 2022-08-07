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
	suite.Require().Equal([]*types.RangesToAmounts{
		{
			Ranges: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Amount: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, aliceBalanceInfo.BalanceAmounts))

	err = RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}})
	suite.Require().Nil(err, "Error revoking badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(0), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, aliceBalanceInfo.BalanceAmounts))
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
	suite.Require().Equal([]*types.RangesToAmounts{
		{
			Ranges: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Amount: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, aliceBalanceInfo.BalanceAmounts))

	err = RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated + 1}, []uint64{7000}, 0, []*types.NumberRange{{Start: 0, End: 0}})
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
	suite.Require().Equal([]*types.RangesToAmounts{
		{
			Ranges: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Amount: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, aliceBalanceInfo.BalanceAmounts))

	accs := suite.app.AccountKeeper.GetAllAccounts(suite.ctx)
	_ = accs

	err = RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}})
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
	suite.Require().Equal([]*types.RangesToAmounts{
		{
			Ranges: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Amount: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, aliceBalanceInfo.BalanceAmounts))

	accs := suite.app.AccountKeeper.GetAllAccounts(suite.ctx)
	_ = accs

	err = RevokeBadges(suite, wctx, bob, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}})
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
	suite.Require().Equal([]*types.RangesToAmounts{
		{
			Ranges: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Amount: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), keeper.GetBadgeBalanceFromBalanceAmountsForSubbadgeId(0, aliceBalanceInfo.BalanceAmounts))

	accs := suite.app.AccountKeeper.GetAllAccounts(suite.ctx)
	_ = accs

	err = RevokeBadges(suite, wctx, alice, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}})
	suite.Require().EqualError(err, keeper.ErrSenderIsNotManager.Error())
}
