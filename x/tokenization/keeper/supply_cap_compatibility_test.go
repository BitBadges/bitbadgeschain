package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func (suite *TestSuite) TestSupplyCapCompatibilityStore() {
	collection := suite.createBackedCollection()
	k := suite.app.TokenizationKeeper
	collection.Invariants.MaxSupplyPerId = sdkmath.OneUint()
	err := k.SetCollectionInStore(suite.ctx, collection, false)
	suite.Require().ErrorContains(err, "maxSupplyPerId cannot be used with cosmosCoinBackedPath")
	stored, found := k.GetCollectionFromStore(suite.ctx, collection.CollectionId)
	suite.Require().True(found)
	suite.Require().True(stored.Invariants.MaxSupplyPerId.IsNil() || stored.Invariants.MaxSupplyPerId.IsZero())

	collection.Invariants.MaxSupplyPerId = sdkmath.ZeroUint()
	suite.Require().NoError(k.SetCollectionInStore(suite.ctx, collection, false))
	collection.Invariants.CosmosCoinBackedPath = nil
	collection.Invariants.MaxSupplyPerId = sdkmath.OneUint()
	suite.Require().NoError(k.SetCollectionInStore(suite.ctx, collection, false))
}

func (suite *TestSuite) TestSupplyCapCompatibilityHandler() {
	collection := suite.createBackedCollection()
	k := suite.app.TokenizationKeeper
	next := k.GetNextCollectionId(suite.ctx)
	for _, id := range []sdkmath.Uint{sdkmath.ZeroUint(), collection.CollectionId} {
		_, err := keeper.NewMsgServerImpl(k).UniversalUpdateCollection(suite.ctx, &types.MsgUniversalUpdateCollection{
			Creator: bob, CollectionId: id,
			Invariants: &types.InvariantsAddObject{
				MaxSupplyPerId:       sdkmath.OneUint(),
				CosmosCoinBackedPath: &types.CosmosCoinBackedPathAddObject{Conversion: collection.Invariants.CosmosCoinBackedPath.Conversion},
			},
		})
		suite.Require().ErrorContains(err, "maxSupplyPerId cannot be used with cosmosCoinBackedPath")
		suite.Require().Equal(next, k.GetNextCollectionId(suite.ctx))
	}
}
