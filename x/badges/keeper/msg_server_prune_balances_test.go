package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestPruneBalances() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: &types.UriObject{
					Uri:                    []byte("example.com/"),
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
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
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

	err = PruneBalances(suite, wctx, bob, []uint64{bobAccountNum}, []uint64{0})
	suite.Require().EqualError(err, keeper.ErrCantPruneBalanceYet.Error())

	err = PruneBalances(suite, wctx, bob, []uint64{bobAccountNum}, []uint64{10})
	suite.Require().EqualError(err, keeper.ErrCantPruneBalanceYet.Error())

	err = SelfDestructBadge(suite, wctx, bob, 0)
	suite.Require().Nil(err, "Error self destructing badge")

	badge, err = GetBadge(suite, wctx, 0)
	suite.Require().NotNil(err, "We should get a not exists error here now")

	_, err = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Nil(err, "Error getting balance info")

	err = PruneBalances(suite, wctx, bob, []uint64{bobAccountNum}, []uint64{0})
	suite.Require().Nil(err, "Error pruning balances")

	bobBalanceInfo, err = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Nil(err, "Error pruning balances")
	suite.Require().Equal(0, len(bobBalanceInfo.BalanceAmounts))
	suite.Require().Equal(0, len(bobBalanceInfo.Approvals))
	suite.Require().Equal(uint64(0), (bobBalanceInfo.PendingNonce))
	suite.Require().Equal(0, len(bobBalanceInfo.Pending))

}
