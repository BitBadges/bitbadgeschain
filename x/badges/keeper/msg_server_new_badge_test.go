package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNewBadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := uint64(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: &types.UriObject{
					Uri:                    "example.com/",
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					InsertIdIdx:            10,
				},
				Permissions: 62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	badge, _ := GetBadge(suite, wctx, 0)

	// Verify nextId increments correctly
	nextId := suite.app.BadgesKeeper.GetNextBadgeId(suite.ctx)
	suite.Require().Equal(uint64(1), nextId)

	// Verify badge details are correct
	suite.Require().Equal(uint64(0), badge.NextSubassetId)
	suite.Require().Equal(&types.UriObject{
		Uri:                    "example.com/",
		Scheme:                 1,
		IdxRangeToRemove:       &types.IdRange{},
		InsertSubassetBytesIdx: 0,
		InsertIdIdx:            10,
	}, badge.Uri)
	suite.Require().Equal([]*types.BalanceObject(nil), badge.SubassetSupplys)
	suite.Require().Equal(bobAccountNum, badge.Manager) //7 is the first ID it creates
	suite.Require().Equal(perms, badge.Permissions)
	suite.Require().Equal([]*types.IdRange(nil), badge.FreezeRanges)
	suite.Require().Equal(uint64(0), badge.Id)

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	// Verify nextId increments correctly
	nextId = suite.app.BadgesKeeper.GetNextBadgeId(suite.ctx)
	suite.Require().Equal(uint64(2), nextId)
	badge, _ = GetBadge(suite, wctx, 1)
	suite.Require().Equal(uint64(1), badge.Id)
}


func (suite *TestSuite) TestNewBadgesWhitelistRecipients() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := uint64(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: &types.UriObject{
					Uri:                    "example.com/",
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					InsertIdIdx:            10,
				},
				SubassetSupplysAndAmounts: []*types.SubassetSupplyAndAmount{
					{
						Supply: 10,
						Amount: 10,
					},
				},
				Permissions: perms,
				WhitelistedRecipients: []*types.WhitelistMintInfo{
					{
						Addresses: []uint64{aliceAccountNum, charlieAccountNum},
						BalanceAmounts: []*types.BalanceObject{
							{
								Balance: 5,
								IdRanges: []*types.IdRange{
									{
										Start: 0,
										End: 4,
									},
								},
							},
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	// badge, _ := GetBadge(suite, wctx, 0)

	// Verify nextId increments correctly
	nextId := suite.app.BadgesKeeper.GetNextBadgeId(suite.ctx)
	suite.Require().Equal(uint64(1), nextId)

	bobBalance, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)
	suite.Require().Equal(uint64(10), bobBalance.BalanceAmounts[0].Balance)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: 5,
			End: 9,
		},
	}, bobBalance.BalanceAmounts[0].IdRanges)

	aliceBalance, _ := GetUserBalance(suite, wctx, 0, aliceAccountNum)
	suite.Require().Equal(uint64(5), aliceBalance.BalanceAmounts[0].Balance)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: 0,
			End: 4,
		},
	}, aliceBalance.BalanceAmounts[0].IdRanges)

	charlieBalance, _ := GetUserBalance(suite, wctx, 0, charlieAccountNum)
	suite.Require().Equal(uint64(5), charlieBalance.BalanceAmounts[0].Balance)
	suite.Require().Equal([]*types.IdRange{
		{
			Start: 0,
			End: 4,
		},
	}, charlieBalance.BalanceAmounts[0].IdRanges)
}


func (suite *TestSuite) TestNewBadgesWhitelistRecipientsOverflow() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	perms := uint64(62)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri: &types.UriObject{
					Uri:                    "example.com/",
					Scheme:                 1,
					IdxRangeToRemove:       &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					InsertIdIdx:            10,
				},
				SubassetSupplysAndAmounts: []*types.SubassetSupplyAndAmount{
					{
						Supply: 10,
						Amount: 10,
					},
				},
				Permissions: perms,
				WhitelistedRecipients: []*types.WhitelistMintInfo{
					{
						Addresses: []uint64{aliceAccountNum, charlieAccountNum},
						BalanceAmounts: []*types.BalanceObject{
							{
								Balance: 6,
								IdRanges: []*types.IdRange{
									{
										Start: 0,
										End: 4,
									},
								},
							},
						},
					},
				},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err = CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().EqualError(err, keeper.ErrUnderflow.Error())
}