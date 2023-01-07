package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestSetApproval() {
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
	err := CreateSubBadges(suite, wctx, bob, 0, []*types.SubassetSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	err = SetApproval(suite, wctx, bob, 1000, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
	suite.Require().Nil(err, "Error setting approval")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(aliceAccountNum), bobBalanceInfo.Approvals[0].Address)
	suite.Require().Equal(uint64(1000), bobBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance)

	err = SetApproval(suite, wctx, bob, 500, charlieAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
	suite.Require().Nil(err, "Error setting approval")

	err = SetApproval(suite, wctx, bob, 500, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
	suite.Require().Nil(err, "Error setting approval")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(charlieAccountNum), bobBalanceInfo.Approvals[1].Address)
	suite.Require().Equal(uint64(500), bobBalanceInfo.Approvals[1].ApprovalAmounts[0].Balance)

	suite.Require().Equal(uint64(aliceAccountNum), bobBalanceInfo.Approvals[0].Address)
	suite.Require().Equal(uint64(500), bobBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance)
}

func (suite *TestSuite) TestSetApprovalNoPrevBalanceInStore() {
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
	err := CreateSubBadges(suite, wctx, bob, 0, []*types.SubassetSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	err = SetApproval(suite, wctx, charlie, 1000, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
	suite.Require().Nil(err, "Error setting approval")
}

func (suite *TestSuite) TestApproveSelf() {
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
	err := CreateSubBadges(suite, wctx, bob, 0, []*types.SubassetSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating subbadge")
	badge, _ = GetBadge(suite, wctx, 0)
	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	err = SetApproval(suite, wctx, bob, 1000, bobAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
	suite.Require().EqualError(err, keeper.ErrAccountCanNotEqualCreator.Error())
}
