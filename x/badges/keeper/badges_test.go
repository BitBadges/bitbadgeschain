package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/keeper"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
)

func (suite *TestSuite) TestGetBadge() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:         &types.UriObject{
					Uri: 	[]byte("example.com/"),
					Scheme: 1,
					IdxRangeToRemove: &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					InsertIdIdx: 10,
				},
				Permissions:  62,
				
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	badge, err := suite.app.BadgesKeeper.GetBadgeE(suite.ctx, 0)
	suite.Require().Nil(err, "Error getting badge: %s")
	suite.Require().Equal(badge.Id, uint64(0))

	badge, err = suite.app.BadgesKeeper.GetBadgeE(suite.ctx, 1)
	suite.Require().EqualError(err, keeper.ErrBadgeNotExists.Error())
}

func (suite *TestSuite) TestGetBadgeAndAssertSubbadges() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	
	badgesToCreate := []BadgesToCreate{
		{
			Badge: types.MsgNewBadge{
				Uri:         &types.UriObject{
					Uri: 	[]byte("example.com/"),
					Scheme: 1,
					IdxRangeToRemove: &types.IdRange{},
					InsertSubassetBytesIdx: 0,
					InsertIdIdx: 10,
				},
				Permissions:  62,
				SubassetSupplys: []uint64{1},
				SubassetAmountsToCreate: []uint64{1},
			},
			Amount:  1,
			Creator: bob,
		},
	}

	err := CreateBadges(suite, wctx, badgesToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")

	_, err = suite.app.BadgesKeeper.GetBadgeAndAssertSubbadgeRangesAreValid(suite.ctx, 0, []*types.IdRange{
		{
			Start: 0,
			End:   0,
		},
	})
	suite.Require().Nil(err, "Error getting badge: %s")

	_, err = suite.app.BadgesKeeper.GetBadgeAndAssertSubbadgeRangesAreValid(suite.ctx, 0, []*types.IdRange{
		{
			Start: 20,
			End:   10,
		},
	})
	suite.Require().EqualError(err, keeper.ErrInvalidSubbadgeRange.Error())

	_, err = suite.app.BadgesKeeper.GetBadgeAndAssertSubbadgeRangesAreValid(suite.ctx, 0, []*types.IdRange{
		{
			Start: 0,
			End:   10,
		},
	})
	suite.Require().EqualError(err, keeper.ErrSubBadgeNotExists.Error())
}