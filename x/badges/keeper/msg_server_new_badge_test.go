package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/trevormil/bitbadgeschain/x/badges/types"
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
					Uri:                    []byte("example.com/"),
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
		Uri:                    []byte("example.com/"),
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
