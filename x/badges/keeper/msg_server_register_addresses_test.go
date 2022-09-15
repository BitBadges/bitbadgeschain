package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestRegisterAddresses() {
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

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{1010}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrAccountNotRegistered.Error())

	err = RegisterAddresses(suite, wctx, bob, []string{"cosmos1f6k8dr0hzafe78uhnmcg7e9qm07ejue8mmd8xq"})
	suite.Require().Nil(err, "Error registering addresses")

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{1010}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring to rregistered address")
}
