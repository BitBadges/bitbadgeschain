package keeper_test

import (
	"time"

	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestHandleAcceptIncomingRequest() {
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
				Permissions: 46,
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

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{0}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().EqualError(err, keeper.ErrBalanceIsZero.Error())

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 10, 12)
	suite.Require().EqualError(err, keeper.ErrCancelTimeIsGreaterThanExpirationTime.Error())

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)

	err = TransferBadge(suite, wctx, alice, aliceAccountNum, []uint64{bobAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(0), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	//Just test that it can readd pending transfer after removal
	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
}

func (suite *TestSuite) TestHandleAcceptIncomingRequestWithApproval() {
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
				Permissions: 46,
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

	err = SetApproval(suite, wctx, bob, 100000, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
	suite.Require().Nil(err, "Error approving badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)

	err = TransferBadge(suite, wctx, alice, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)
	suite.Require().Equal(uint64(100000-5000), bobBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)

	err = TransferBadge(suite, wctx, alice, aliceAccountNum, []uint64{bobAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(0), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	//Just test that it can readd pending transfer after removal
	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
}

func (suite *TestSuite) TestHandleRejectIncomingRequest() {
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
				Permissions: 46,
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
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, false, false)
	suite.Require().Nil(err, "Error accepting badge")

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, false, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(0), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)
}

func (suite *TestSuite) TestHandleRejectIncomingRequestWithApproval() {
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
				Permissions: 46,
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

	err = SetApproval(suite, wctx, bob, 100000, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, alice, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, false, false)
	suite.Require().Nil(err, "Error accepting badge")

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, false, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)
	suite.Require().Equal(uint64(100000), bobBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Approvals[0].Address)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(0), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)
	suite.Require().Equal([]*types.Approval(nil), aliceBalanceInfo.Approvals)
}

func (suite *TestSuite) TestHandleCancelOutgoingRequestWithApproval() {
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
				Permissions: 46,
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

	err = SetApproval(suite, wctx, bob, 100000, aliceAccountNum, 0, []*types.IdRange{{Start: 0, End: 0}})
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, alice, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, false, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)
	suite.Require().Equal(uint64(100000), bobBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Approvals[0].Address)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(0), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, aliceBalanceInfo.BalanceAmounts)[0].Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)
	suite.Require().Equal([]*types.Approval(nil), aliceBalanceInfo.Approvals)
}

func (suite *TestSuite) TestHandleCancelOutgoingRequest() {
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
				Permissions: 46,
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

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, uint64(suite.ctx.BlockTime().Unix()+100000), 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().EqualError(err, keeper.ErrCantAcceptOwnTransferRequest.Error())
}

func (suite *TestSuite) TestBadgeDoesntExist() {
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
				Permissions: 46,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	err = HandlePendingTransfers(suite, wctx, bob, 1000, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().EqualError(err, keeper.ErrBadgeNotExists.Error())
}

func (suite *TestSuite) TestAcceptExpiredTransfer() {
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
				Permissions: 46,
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

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, uint64(suite.ctx.BlockTime().Unix()+1), 0)
	suite.Require().Nil(err, "Error transferring badge")

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Second * 20000))
	wctx = sdk.WrapSDKContext(suite.ctx)

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().EqualError(err, keeper.ErrPendingTransferExpired.Error())
}

func (suite *TestSuite) TestNonexistentPendingTransfer() {
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
				Permissions: 46,
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

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().Nil(err, "Error handling transfer")
}

func (suite *TestSuite) TestPendingBinarySearch() {
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
				Permissions: 46,
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

	for i := 0; i < 100; i++ {
		err = RequestTransferBadge(suite, wctx, alice, bobAccountNum, 1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
		suite.Require().Nil(err, "Error transferring badge")
	}

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 5, End: 0}}, true, false)
	suite.Require().Nil(err, "Error handling badge")

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 95, End: 0}}, true, false)
	suite.Require().Nil(err, "Error handling badge")

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 5, End: 0}}, true, false)
	suite.Require().Nil(err, "Error handling badge")

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 95, End: 0}}, true, false)
	suite.Require().Nil(err, "Error handling badge")
}

func (suite *TestSuite) TestPruneExpiredTransfer() {
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
				Permissions: 46,
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

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, uint64(suite.ctx.BlockTime().Unix()-1), 0)
	suite.Require().Nil(err, "Error transferring badge")

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Second * 20000))
	wctx = sdk.WrapSDKContext(suite.ctx)

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(0, len(aliceBalanceInfo.Pending))

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().EqualError(err, keeper.ErrPendingTransferExpired.Error())

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, false, false)
	suite.Require().Nil(err, "Error reverting transfer")
}

func (suite *TestSuite) TestCancelBeforeTimesForTransfer() {
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
				Permissions: 46,
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

	err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{aliceAccountNum}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, uint64(suite.ctx.BlockTime().Unix()+1000))
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, false, false)
	suite.Require().EqualError(err, keeper.ErrCantCancelYet.Error())

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(time.Second * 20000))
	wctx = sdk.WrapSDKContext(suite.ctx)

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, false, false)
	suite.Require().Nil(err, "error cancelling transfer")
}

func (suite *TestSuite) TestAcceptForcefullyAfterApproved() {
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
				Permissions: 46,
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
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(aliceAccountNum, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(bobAccountNum, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, 0, []*types.IdRange{{Start: 0, End: 0}}, true, false)
	suite.Require().Nil(err, "error marking transfer as approved")

	err = HandlePendingTransfers(suite, wctx, bob, 0, []*types.IdRange{{Start: 0, End: 0}}, false, false)
	suite.Require().Nil(err, "error cancelling transfer")
}
