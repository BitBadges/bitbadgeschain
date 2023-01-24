package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestRegisterAddresses() {
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

	//Create badge 1 with supply > 1
	err := CreateBadgesAndMintAllToCreator(suite, wctx, bob, 0, []*types.BadgeSupplyAndAmount{
		{
			Supply: 10000,
			Amount: 1,
		},
	})
	suite.Require().Nil(err, "Error creating badge")
	badge, _ = GetCollection(suite, wctx, 0)
	bobbalance, _ := GetUserBalance(suite, wctx, 0, bobAccountNum)

	suite.Require().Equal(uint64(1), badge.NextBadgeId)
	suite.Require().Equal([]*types.Balance{
		{
			BadgeIds: []*types.IdRange{{Start: 0, End: 0}}, //0 to 0 range so it will be nil
			Balance:  10000,
		},
	}, badge.MaxSupplys)
	suite.Require().Equal(uint64(10000), keeper.GetBalancesForIdRanges([]*types.IdRange{{Start: 0}}, bobbalance.Balances)[0].Balance)

	// err = TransferBadge(suite, wctx, bob, bobAccountNum, []uint64{1010}, []uint64{5000}, 0, []*types.IdRange{{Start: 0, End: 0}}, 0, 0)
	// suite.Require().EqualError(err, keeper.ErrAccountNotRegistered.Error())

	registered := suite.app.AccountKeeper.HasAccount(suite.ctx, sdk.MustAccAddressFromBech32("cosmos1f6k8dr0hzafe78uhnmcg7e9qm07ejue8mmd8xq"))
	suite.Require().False(registered, "Account should not be registered")

	err = RegisterAddresses(suite, wctx, bob, []string{"cosmos1f6k8dr0hzafe78uhnmcg7e9qm07ejue8mmd8xq"})
	suite.Require().Nil(err, "Error registering addresses")

	err = TransferBadge(suite, wctx, bob, 0, bobAccountNum, []*types.Transfers{
		{
			ToAddresses: []uint64{1010},
			Balances: []*types.Balance{
				{
					Balance: 5000,
					BadgeIds: []*types.IdRange{
						{
							Start: 0,
							End:   0,
						},
					},
				},
			},
		},
	})
	suite.Require().Nil(err, "Error transfering badge")
	suite.Require().Nil(err, "Error transferring to rregistered address")

	registered = suite.app.AccountKeeper.HasAccount(suite.ctx, sdk.MustAccAddressFromBech32("cosmos1f6k8dr0hzafe78uhnmcg7e9qm07ejue8mmd8xq"))
	suite.Require().True(registered, "Account should be registered")
}
