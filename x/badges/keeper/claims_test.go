package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestSendAllToClaims() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUri:      "https://example.com/{id}",
				Permissions:   62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 0)

	claimToAdd := types.Claim{
		Balance: &types.Balance{
			Balance:  10,
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}},
		},
	}

	err := CreateBadges(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 1,
		},
	},
		[]*types.Transfers{},
		[]*types.Claim{
			&claimToAdd,
		})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 0)

	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10,
		},
	}, badge.MaxSupplys)

	claim := badge.Claims[0]
	suite.Require().Nil(err, "Error getting claim")
	suite.Require().Equal(&claimToAdd, claim)
}
