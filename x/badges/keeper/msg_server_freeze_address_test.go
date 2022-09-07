package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestFreezeAddressesDirectlyWhenCreatingNewBadge() {
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

				FreezeAddressRanges: []*types.IdRange{
					{Start: aliceAccountNum},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")
	// badge, _ := GetBadge(suite, wctx, 0)

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, alice, aliceAccountNum, []uint64{bobAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

func (suite *TestSuite) TestTransferBadgeForcefulUnfrozenByDefault() {
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

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)

	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)

	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: aliceAccountNum, End: aliceAccountNum}})
	suite.Require().Nil(err, "Error freezing address")

	badge, _ = GetBadge(suite, wctx, 0)

	err = TransferBadge(suite, wctx, alice, aliceAccountNum, []uint64{bobAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

func (suite *TestSuite) TestTransferBadgeForcefulFrozenByDefault() {
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
				Permissions: 63,
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

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
	suite.Require().Nil(err, "Error unfreezing address")

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring after unfreeze")

	err = TransferBadge(suite, wctx, alice, aliceAccountNum, []uint64{bobAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

func (suite *TestSuite) TestTransferBadgeForcefulFrozenByDefaultAddAndRemove() {
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
				Permissions: 63,
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

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
	suite.Require().Nil(err, "Error unfreezing address")

	err = FreezeAddresses(suite, wctx, bob, 0, false, []*types.IdRange{{Start: bobAccountNum, End: 0}})
	suite.Require().Nil(err, "Error unfreezing address")

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
	suite.Require().Nil(err, "Error unfreezing address")

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring after unfreeze")

	err = TransferBadge(suite, wctx, alice, aliceAccountNum, []uint64{bobAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

func (suite *TestSuite) TestFreezeCantFreeze() {
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
				Permissions: 0,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	err := FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}

func (suite *TestSuite) TestTransferBadgeForcefulUnfrozenByDefaultOmitEmptyCase() {
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

	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: 0, End: 0}})
	suite.Require().Nil(err, "Error freezing address")

	badge, _ = GetBadge(suite, wctx, 0)

	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: 1, End: 0}})
	suite.Require().Nil(err, "Error freezing address")

	badge, _ = GetBadge(suite, wctx, 0)
	suite.Require().Equal(badge.FreezeRanges, []*types.IdRange{{Start: 0, End: 1}})
}
