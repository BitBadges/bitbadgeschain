package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestTransferBadgeForceful() {
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
	suite.Require().Equal([]*types.BalanceToIds{
		{
			Ids: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), keeper.GetBalanceForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), keeper.GetBalanceForSubbadgeId(0, aliceBalanceInfo.BalanceAmounts))
}

func (suite *TestSuite) TestTransferBadgePending() {
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
	bobBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceToIds{
		{
			Ids: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].SendRequest)

	aliceBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].SendRequest)
}

func (suite *TestSuite) TestTransferBadgeForceBurn() {
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
	bobBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.BalanceToIds{
		{
			Ids: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].SendRequest)

	aliceBalanceInfo, _ := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated+1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].SendRequest)

	err = HandlePendingTransfers(suite, wctx, alice, true, 0, []*types.NumberRange{{Start: 0, End: 0}}, false)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), keeper.GetBalanceForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(5000), keeper.GetBalanceForSubbadgeId(0, aliceBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated+1, []uint64{firstAccountNumCreated}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	aliceBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated+1)
	suite.Require().Equal(uint64(0), keeper.GetBalanceForSubbadgeId(0, aliceBalanceInfo.BalanceAmounts))
}

func (suite *TestSuite) TestApprovalsApproved() {
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
	suite.Require().Equal([]*types.BalanceToIds{
		{
			Ids: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = SetApproval(suite, wctx, bob, 1000000, firstAccountNumCreated+1, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error setting approval")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	// suite.Require().Equal(uint64(1000000-5000), bobBalanceInfo.Approvals[0].Amount)

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo, _ = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(1000000-5000), bobBalanceInfo.Approvals[0].ApprovalAmounts[0].Balance)
}

func (suite *TestSuite) TestApprovalsNotEnoughApproved() {
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
	suite.Require().Equal([]*types.BalanceToIds{
		{
			Ids: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = SetApproval(suite, wctx, bob, 10, firstAccountNumCreated+1, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().Nil(err, "Error setting approval")

	err = TransferBadge(suite, wctx, charlie, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().EqualError(err, keeper.ErrInsufficientApproval.Error())
}

func (suite *TestSuite) TestApprovalsNotApproved() {
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
	suite.Require().Equal([]*types.BalanceToIds{
		{
			Ids: []*types.NumberRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), keeper.GetBalanceForSubbadgeId(0, bobBalanceInfo.BalanceAmounts))

	err = TransferBadge(suite, wctx, charlie, firstAccountNumCreated, []uint64{firstAccountNumCreated + 1}, []uint64{5000}, 0, []*types.NumberRange{{Start: 0, End: 0}}, 0)
	suite.Require().EqualError(err, keeper.ErrInsufficientApproval.Error())
}
