package keeper_test

import (
	"math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestSendAllToClaims() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	collectionsToCreate := []CollectionsToCreate{
		{
			Collection: types.MsgNewCollection{
				CollectionUri: "https://example.com",
				BadgeUris: []*types.BadgeUri{
					{
						Uri: "https://example.com/{id}",
						BadgeIds: []*types.IdRange{
							{
								Start: 1,
								End:   math.MaxUint64,
							},
						},
					},
				},
				Permissions: 62,
			},
			Amount:  1,
			Creator: bob,
		},
	}

	CreateCollections(suite, wctx, collectionsToCreate)
	badge, _ := GetCollection(suite, wctx, 1)

	claimToAdd := types.Claim{
		Balances: []*types.Balance{{
			Amount:  10,
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}},
		}},
	}

	err := CreateBadges(suite, wctx, bob, 1, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10,
			Amount: 1,
		},
	},
		[]*types.Transfers{},
		[]*types.Claim{
			&claimToAdd,
		}, "https://example.com",
		[]*types.BadgeUri{
			{
				Uri: "https://example.com/{id}",
				BadgeIds: []*types.IdRange{
					{
						Start: 1,
						End:   math.MaxUint64,
					},
				},
			},
		}, "")
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 1)

	suite.Require().Equal([]*types.Balance(nil), badge.UnmintedSupplys)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 1, End: 1}}, //0 to 0 range so it will be nil
			Amount:  10,
		},
	}, badge.MaxSupplys)

	claim, _ := GetClaim(suite, wctx, 1, 1)
	suite.Require().Nil(err, "Error getting claim")
	suite.Require().Equal(claimToAdd, claim)
}
