package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestRequestTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: validUri,
				Permissions: 62,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: bob,
		},
		{
			Badge: types.MsgNewBadge{
				Uri: validUri,
				Permissions: 62,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: alice,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")

	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.PendingNonce)

	err = RequestTransferBadge(suite, wctx, alice, firstAccountNumCreated, 5000, 0, 0)
	suite.Require().Nil(err, "Error requesting transfer")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(5000), bobBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, bobBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, bobBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), bobBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(false, bobBalanceInfo.Pending[0].SendRequest)

	aliceBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(5000), aliceBalanceInfo.Pending[0].Amount)
	suite.Require().Equal(uint64(1), aliceBalanceInfo.PendingNonce)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].ApprovedBy)
	suite.Require().Equal(firstAccountNumCreated + 1, aliceBalanceInfo.Pending[0].To)
	suite.Require().Equal(firstAccountNumCreated, aliceBalanceInfo.Pending[0].From)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].ThisPendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.Pending[0].OtherPendingNonce)
	suite.Require().Equal(true, aliceBalanceInfo.Pending[0].SendRequest)
}

func (suite *TestSuite) TestRequestTransferToSelf() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: validUri,
				Permissions: 62,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: bob,
		},
		{
			Badge: types.MsgNewBadge{
				Uri: validUri,
				Permissions: 62,
				FreezeAddressesDigest: "",
				SubassetUris: validUri,
			},
			Amount: 1,
			Creator: alice,
		},
	}

	CreateBadges(suite, wctx, badgesToCreate)

	//Create subbadge 1 with supply > 1
	err := CreateSubBadge(suite, wctx, bob, 0, 10000)
	suite.Require().Nil(err, "Error creating subbadge")

	bobBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	aliceBalanceInfo := GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated + 1)
	suite.Require().Equal(uint64(10000), bobBalanceInfo.Balance)
	suite.Require().Equal(uint64(0), bobBalanceInfo.PendingNonce)
	suite.Require().Equal(uint64(0), aliceBalanceInfo.PendingNonce)

	err = RequestTransferBadge(suite, wctx, bob, firstAccountNumCreated, 5000, 0, 0)
	suite.Require().EqualError(err, keeper.ErrSenderAndReceiverSame.Error())
}