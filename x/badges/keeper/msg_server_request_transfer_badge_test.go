package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestRequestTransfer() {
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
			Creator: alice,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)
	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.PendingNonce)

	err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 5000, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error requesting transfer")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(0), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().EqualError(err, keeper.ErrNotApproved.Error())

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().Nil(err, "Error accepting transfer")

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
}

func (suite *TestSuite) TestRequestTransferForcefulAccept() {
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
			Creator: alice,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)
	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.PendingNonce)

	err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 5000, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error requesting transfer")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(0), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().EqualError(err, keeper.ErrNotApproved.Error())

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, true, true)
	suite.Require().Nil(err, "Error accepting transfer")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
}

func (suite *TestSuite) TestRequestTransferFrozen() {
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
			Creator: alice,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)
	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.PendingNonce)

	err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 5000, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error requesting transfer")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(0), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, aliceBalanceInfo.Pending[0].Sent)

	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
	suite.Require().Nil(err, "Error freezing address")

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().Nil(err, "Error freezing address")

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().EqualError(err, keeper.ErrAddressFrozen.Error())
}

func (suite *TestSuite) TestRequestTransferFrozenThenUnrozen() {
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
			Creator: alice,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)
	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.PendingNonce)

	err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 5000, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error requesting transfer")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(0), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, aliceBalanceInfo.Pending[0].Sent)

	err = FreezeAddresses(suite, wctx, bob, 0, true, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
	suite.Require().Nil(err, "Error freezing address")

	err = FreezeAddresses(suite, wctx, bob, 0, false, []*types.IdRange{{Start: bobAccountNum, End: bobAccountNum}})
	suite.Require().Nil(err, "Error unfreezing address")

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().Nil(err, "Error accepting transfer")

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().Nil(err, "Error accepting transfer")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
}
func (suite *TestSuite) TestRequestTransferToSelf() {
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
			Creator: alice,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)
	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.PendingNonce)

	err = RequestTransferBadge(suite, wctx, bob, bobAccountNum, 5000, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrAccountCanNotEqualCreator.Error())
}

func (suite *TestSuite) TestTryToAcceptTranferRequestBeforeMarkedAsApproved() {
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
			Creator: alice,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	bobBalanceInfo, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)
	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.PendingNonce)

	err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 5000, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error requesting transfer")

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().EqualError(err, keeper.ErrNotApproved.Error())
}
