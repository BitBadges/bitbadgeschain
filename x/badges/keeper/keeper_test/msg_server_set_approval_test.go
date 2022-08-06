package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestSetApproval() {
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

	err = SetApproval(suite, wctx, bob, 1000, firstAccountNumCreated+1, 0, types.SubbadgeRange{Start: 0, End: 0})
	suite.Require().Nil(err, "Error setting approval")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(firstAccountNumCreated+1), bobBalanceInfo.Approvals[0].Address)
	suite.Require().Equal(uint64(1000), bobBalanceInfo.Approvals[0].Amounts[0])

	err = SetApproval(suite, wctx, bob, 500, firstAccountNumCreated+2, 0, types.SubbadgeRange{Start: 0, End: 0})
	suite.Require().Nil(err, "Error setting approval")

	err = SetApproval(suite, wctx, bob, 500, firstAccountNumCreated+1, 0, types.SubbadgeRange{Start: 0, End: 0})
	suite.Require().Nil(err, "Error setting approval")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(firstAccountNumCreated+2), bobBalanceInfo.Approvals[1].Address)
	suite.Require().Equal(uint64(500), bobBalanceInfo.Approvals[1].Amounts[0])

	suite.Require().Equal(uint64(firstAccountNumCreated+1), bobBalanceInfo.Approvals[0].Address)
	suite.Require().Equal(uint64(500), bobBalanceInfo.Approvals[0].Amounts[0])
}

func (suite *TestSuite) TestApproveSelf() {
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

	err = SetApproval(suite, wctx, bob, 1000, firstAccountNumCreated, 0, types.SubbadgeRange{Start: 0, End: 0})
	suite.Require().EqualError(err, keeper.ErrSenderAndReceiverSame.Error())
}
