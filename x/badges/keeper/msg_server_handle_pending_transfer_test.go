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
				Uri: validUri,
				Permissions: 46,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id: 0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated + 1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].SendRequest)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].SendRequest)

	err = HandlePendingTransfer(suite, wctx, alice, true, 0, 0, 0)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)

	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated + 1, firstAccountNumCreated, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated + 1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	//Just test that it can readd pending transfer after removal
	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	
	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
}


func (suite *TestSuite) TestHandleAcceptIncomingRequestWithApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: validUri,
				Permissions: 46,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id: 0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = SetApproval(suite, wctx, bob, 100000, firstAccountNumCreated+1, 0, 0)
	suite.Require().Nil(err, "Error approving badge")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated, firstAccountNumCreated + 1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].SendRequest)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].SendRequest)

	err = HandlePendingTransfer(suite, wctx, alice, true, 0, 0, 0)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)
	suite.Require().Equal(uint64(100000 - 5000), bobBalanceInfo.Approvals[0].Amount)

	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated + 1, firstAccountNumCreated, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated + 1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	//Just test that it can readd pending transfer after removal
	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	
	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
}


func (suite *TestSuite) TestHandleRejectIncomingRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: validUri,
				Permissions: 46,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id: 0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated + 1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].SendRequest)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].SendRequest)

	err = HandlePendingTransfer(suite, wctx, alice, false, 0, 0, 0)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)

	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)
}


func (suite *TestSuite) TestHandleRejectIncomingRequestWithApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: validUri,
				Permissions: 46,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id: 0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = SetApproval(suite, wctx, bob, 100000, firstAccountNumCreated + 1, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated, firstAccountNumCreated + 1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].SendRequest)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].SendRequest)

	err = HandlePendingTransfer(suite, wctx, alice, false, 0, 0, 0)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)
	suite.Require().Equal(uint64(100000), bobBalanceInfo.Approvals[0].Amount)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Approvals[0].AddressNum)


	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)
	suite.Require().Equal([]*types.Approval(nil), aliceBalanceInfo.Approvals)
}

func (suite *TestSuite) TestHandleCancelOutgoingRequestWithApproval() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: validUri,
				Permissions: 46,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id: 0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = SetApproval(suite, wctx, bob, 100000, firstAccountNumCreated + 1, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	err = TransferBadge(suite, wctx, alice, firstAccountNumCreated, firstAccountNumCreated + 1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].SendRequest)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].SendRequest)

	err = HandlePendingTransfer(suite, wctx, bob, false, 0, 0, 0)
	suite.Require().Nil(err, "Error accepting badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), bobBalanceInfo.Pending)
	suite.Require().Equal(uint64(100000), bobBalanceInfo.Approvals[0].Amount)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Approvals[0].AddressNum)

	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Balance)
	suite.Require().Equal([]*types.PendingTransfer(nil), aliceBalanceInfo.Pending)
	suite.Require().Equal([]*types.Approval(nil), aliceBalanceInfo.Approvals)
}


func (suite *TestSuite) TestHandleCancelOutgoingRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: validUri,
				Permissions: 46,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: bob,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)
	badge := GetBadge(suite, wctx, 0)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")
	badge = GetBadge(suite, wctx, 0)
	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(1), badge.NextSubassetId)
	suite.Require().Equal([]*types.Subasset{
		{
			Id: 0,
			Supply: 10000,
		},
	}, badge.SubassetsTotalSupply)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)

	err = TransferBadge(suite, wctx, bob, firstAccountNumCreated, firstAccountNumCreated + 1, 5000, 0, 0)
	suite.Require().Nil(err, "Error transferring badge")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, bobBalanceInfo.Pending[0].SendRequest)

	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, aliceBalanceInfo.Pending[0].SendRequest)

	err = HandlePendingTransfer(suite, wctx, bob, true, 0, 0, 0)
	suite.Require().EqualError(err, keeper.ErrCantAcceptOwnTransferRequest.Error())
}

