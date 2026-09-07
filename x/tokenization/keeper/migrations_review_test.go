package keeper_test

import (
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

func (suite *TestSuite) TestV35DynamicStoreCollisionPreservesRevocation() {
	k := suite.app.TokenizationKeeper
	suite.Require().NoError(k.SetDynamicStoreValueInStore(suite.ctx, sdkmath.OneUint(), alice, true))
	suite.Require().NoError(k.SetDynamicStoreValueInStore(suite.ctx, sdkmath.OneUint(), strings.ToUpper(alice), false))
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	value, found := k.GetDynamicStoreValueFromStore(suite.ctx, sdkmath.OneUint(), alice)
	suite.Require().True(found)
	suite.Require().False(value.Value)
}

func (suite *TestSuite) TestV35BalanceCollisionCorrectsHolderCount() {
	k := suite.app.TokenizationKeeper
	suite.Require().NoError(CreateCollections(suite, suite.ctx, GetTransferableCollectionToCreateAllMintedToCreator(bob)))
	store := suite.ctx.KVStore(suite.app.GetKey(types.StoreKey))
	key := func(address string) []byte {
		return append(append([]byte{}, keeper.UserBalanceKey...), []byte(keeper.ConstructBalanceKey(address, sdkmath.OneUint()))...)
	}
	store.Set(key(strings.ToUpper(bob)), store.Get(key(bob)))
	stats, _ := k.GetCollectionStatsFromStore(suite.ctx, sdkmath.OneUint())
	stats.HolderCount = sdkmath.NewUint(2)
	suite.Require().NoError(k.SetCollectionStatsInStore(suite.ctx, sdkmath.OneUint(), stats))
	suite.Require().NoError(k.MigrateV35CanonicalAddresses(suite.ctx))
	stats, _ = k.GetCollectionStatsFromStore(suite.ctx, sdkmath.OneUint())
	suite.Require().Equal(sdkmath.OneUint(), stats.HolderCount)
}
