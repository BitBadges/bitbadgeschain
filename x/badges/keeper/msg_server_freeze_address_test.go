package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestFreezeAddressesDirectlyWhenCreatingNewBadge() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  62,
				SubassetUris: validUri,
				FreezeAddressRanges: []*types.IdRange{
					{Start: firstAccountNumCreated + 1},
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

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

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

	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(5000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)

	suite.Require().Equal(uint64(5000), keeper.GetBalanceForId(0, aliceBalanceInfo.BalanceAmounts))

	err = FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: firstAccountNumCreated + 1, End: firstAccountNumCreated + 1}}, 0, 0, true)
	suite.Require().Nil(err, "Error freezing address")

	badge, _ = GetBadge(suite, wctx, 0)

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
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

	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

	err = FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: firstAccountNumCreated, End: firstAccountNumCreated}}, 0, 0, true)
	suite.Require().Nil(err, "Error unfreezing address")

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring after unfreeze")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
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

	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceObject{
		{
			IdRanges: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.SubassetSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

	err = FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: firstAccountNumCreated, End: firstAccountNumCreated}}, 0, 0, true)
	suite.Require().Nil(err, "Error unfreezing address")

	err = FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: firstAccountNumCreated, End: firstAccountNumCreated}}, 0, 0, false)
	suite.Require().Nil(err, "Error unfreezing address")

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())

	err = FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: firstAccountNumCreated, End: firstAccountNumCreated}}, 0, 0, true)
	suite.Require().Nil(err, "Error unfreezing address")

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring after unfreeze")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

func (suite *TestSuite) TestFreezeCantFreeze() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  0,
				SubassetUris: validUri,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	err := FreezeAddresses(suite, wctx, bob, []*types.IdRange{{Start: firstAccountNumCreated, End: firstAccountNumCreated}}, 0, 0, true)
	suite.Require().EqualError(err, keeper.ErrInvalidPermissions.Error())
}
