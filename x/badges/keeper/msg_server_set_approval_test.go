package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestSetApproval() {
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

	err = SetApproval(suite, wctx, bob, 1000, firstAccountNumCreated + 1, 0, 0)
	suite.Require().Nil(err, "Error setting approval")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)
	suite.Require().Equal(uint64(firstAccountNumCreated + 1), bobBalanceInfo.Approvals[0].AddressNum)
	suite.Require().Equal(uint64(1000), bobBalanceInfo.Approvals[0].Amount)

	err = SetApproval(suite, wctx, bob, 500, firstAccountNumCreated + 2, 0, 0)
	suite.Require().Nil(err, "Error setting approval")

	err = SetApproval(suite, wctx, bob, 500, firstAccountNumCreated + 1, 0, 0)
	suite.Require().Nil(err, "Error setting approval")

	bobBalanceInfo = GetBadgeBalance(suite, wctx, 0, 0, firstAccountNumCreated)

	suite.Require().Equal(uint64(firstAccountNumCreated + 2), bobBalanceInfo.Approvals[0].AddressNum)
	suite.Require().Equal(uint64(500), bobBalanceInfo.Approvals[0].Amount)

	suite.Require().Equal(uint64(firstAccountNumCreated + 1), bobBalanceInfo.Approvals[1].AddressNum)
	suite.Require().Equal(uint64(500), bobBalanceInfo.Approvals[1].Amount)
}

func (suite *TestSuite) TestApproveSelf() {
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

	err = SetApproval(suite, wctx, bob, 1000, firstAccountNumCreated, 0, 0)
	suite.Require().EqualError(err, keeper.ErrSenderAndReceiverSame.Error())
}