package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestNewCollection() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.BadgesToCreate = []*types.Balance{}

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
	collection, _ := GetCollection(suite, wctx, sdkmath.NewUint(1))

	// Verify nextId increments correctly
	nextId := suite.app.BadgesKeeper.GetNextCollectionId(suite.ctx)
	AssertUintsEqual(suite, sdkmath.NewUint(2), nextId)

	// Verify badge details are correct
	AssertUintsEqual(suite, sdkmath.NewUint(1), collection.NextBadgeId)
}

func (suite *TestSuite) TestNewCollectionDifferentBalancesTypes() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.BadgesToCreate = []*types.Balance{}
	collectionsToCreate[0].Collection.BalancesType = sdkmath.NewUint(1)
	

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "Error creating badge: %s")

	collectionsToCreate = GetCollectionsToCreate()
	collectionsToCreate[0].Collection.BadgesToCreate = []*types.Balance{}
	collectionsToCreate[0].Collection.BalancesType = sdkmath.NewUint(1)
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline = nil
	
	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")


	collectionsToCreate = GetCollectionsToCreate()
	collectionsToCreate[0].Collection.BadgesToCreate = []*types.Balance{}
	collectionsToCreate[0].Collection.BalancesType = sdkmath.NewUint(2)
	

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "Error creating badge: %s")

	collectionsToCreate = GetCollectionsToCreate()
	collectionsToCreate[0].Collection.BadgesToCreate = []*types.Balance{}
	collectionsToCreate[0].Collection.BalancesType = sdkmath.NewUint(2)
	collectionsToCreate[0].Collection.CollectionApprovedTransfersTimeline = nil
	
	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Nil(err, "Error creating badge: %s")
}


func (suite *TestSuite) TestNewCollectionDuplicateBadgeIds() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := sdk.AccAddressFromBech32(alice)
	suite.Require().Nil(err, "Address %s failed to parse")

	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Collection.Transfers = []*types.Transfer{
		{
			From: "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					BadgeIds: []*types.UintRange{
						GetOneUintRange()[0],
						GetOneUintRange()[0],
					},
					OwnershipTimes: GetFullUintRanges(),
				},
			},
		},
	}
	

	err = CreateCollections(suite, wctx, collectionsToCreate)
	suite.Require().Error(err, "Error creating badge: %s")
}