package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestHandleAcceptIncomingRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
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
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), keeper.GetBalanceForId(0, aliceBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(0), keeper.GetBalanceForId(0, aliceBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	//Just test that it can readd pending transfer after removal
	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
}

func (suite *TestSuite) TestHandleAcceptIncomingRequestWithApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
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

	err = SetApproval(suite, wctx, bob, 100000, firstAccountNumCreated+1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error approving badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)
	suite.Require().Equal(uint64(100000-5000), bobBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), keeper.GetBalanceForId(0, aliceBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(0), keeper.GetBalanceForId(0, aliceBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	//Just test that it can readd pending transfer after removal
	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
}

func (suite *TestSuite) TestHandleRejectIncomingRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
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
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().Nil(err, "Error accepting badge")

	err = HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(0), keeper.GetBalanceForId(0, aliceBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)
}

func (suite *TestSuite) TestHandleRejectIncomingRequestWithApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
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

	err = SetApproval(suite, wctx, bob, 100000, firstAccountNumCreated+1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, alice, false, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().Nil(err, "Error accepting badge")

	err = HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)
	suite.Require().Equal(uint64(100000), bobBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Approvals[0].Address)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(0), keeper.GetBalanceForId(0, aliceBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)
	suite.Require().Equal([]*types.Approval(nil), aliceBalanceInfo.Approvals)
}

func (suite *TestSuite) TestHandleCancelOutgoingRequestWithApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
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

	err = SetApproval(suite, wctx, bob, 100000, firstAccountNumCreated+1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, bob, false, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForId(0, bobBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)
	suite.Require().Equal(uint64(100000), bobBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Approvals[0].Address)

	aliceBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(0), keeper.GetBalanceForId(0, aliceBalanceInfo.BalanceAmounts))
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)
	suite.Require().Equal([]*types.Approval(nil), aliceBalanceInfo.Approvals)
}

func (suite *TestSuite) TestHandleCancelOutgoingRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
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

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, uint64(suite.ctx.BlockTime().Unix()+100000))
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().EqualError(err, keeper.ErrCantAcceptOwnTransferRequest.Error())
}

func (suite *TestSuite) TestBadgeDoesntExist() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
				SubassetUris: validUri,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadges(suite, wctx, bob, 0, []uint64{10000}, []uint64{1})
	suite.Require().Nil(err, "Error creating subbadge")

	err = HandlePendingTransfers(suite, wctx, bob, true, 1000, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().EqualError(err, keeper.ErrBadgeNotExists.Error())
}

func (suite *TestSuite) TestAcceptExpiredTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
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

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, uint64(suite.ctx.BlockTime().Unix()-1))
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].Sent)

	aliceBalanceInfo, _ := GetUserBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].Sent)

	err = HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().EqualError(err, keeper.ErrPendingTransferExpired.Error())
}

func (suite *TestSuite) TestNonexistentPendingTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
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

	err = HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 0, End: 0}}, false)
	suite.Require().EqualError(err, keeper.ErrNoPendingTransferFound.Error())
}

func (suite *TestSuite) TestPendingBinarySearch() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:          validUri,
				Permissions:  46,
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

	for i := 0; i < 100; i++ {
		err = RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 1, 0, []*types.IdRange{{Start: 0, End: 0}}, 0)
		suite.Require().Nil(err, "Error transferring badge")
	}

	err = HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 5, End: 0}}, false)
	suite.Require().Nil(err, "Error handling badge")

	err = HandlePendingTransfers(suite, wctx, bob, true, 0, []*types.IdRange{{Start: 95, End: 0}}, false)
	suite.Require().Nil(err, "Error handling badge")

	err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 5, End: 0}}, false)
	suite.Require().Nil(err, "Error handling badge")

	err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.IdRange{{Start: 95, End: 0}}, false)
	suite.Require().Nil(err, "Error handling badge")
}
