package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *TestSuite) TestDeleteCollectionRefusesMintEscrowCoins() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	suite.Require().NoError(CreateCollections(suite, ctx, GetTransferableCollectionToCreateAllMintedToCreator(bob)))
	collection, err := GetCollection(suite, ctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	escrow := sdk.MustAccAddressFromBech32(collection.MintEscrowAddress)
	coins := sdk.NewCoins(sdk.NewInt64Coin("ubadge", 1))
	suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coins))
	suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", escrow, coins))
	msg := &types.MsgDeleteCollection{Creator: bob, CollectionId: collection.CollectionId}
	suite.Require().ErrorContains(DeleteCollection(suite, ctx, msg), "mint escrow")
	_, err = GetCollection(suite, ctx, collection.CollectionId)
	suite.Require().NoError(err)
	suite.Require().Equal(coins, suite.app.BankKeeper.GetAllBalances(suite.ctx, escrow))
	suite.Require().NoError(suite.app.BankKeeper.SendCoins(suite.ctx, escrow, sdk.MustAccAddressFromBech32(bob), coins))
	suite.Require().NoError(DeleteCollection(suite, ctx, msg))
}

func (suite *PoolIntegrationTestSuite) TestDeleteCollectionRefusesAliasSupply() {
	collection, denom, err := suite.createCollectionWithAliasPath(bob, "delete-test")
	suite.Require().NoError(err)
	coins := sdk.NewCoins(sdk.NewInt64Coin(denom, 1))
	suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, "tokenization", coins))
	msg := &types.MsgDeleteCollection{Creator: bob, CollectionId: collection.CollectionId}
	suite.Require().ErrorContains(DeleteCollection(&suite.TestSuite, suite.ctx, msg), "alias")
	suite.Require().NoError(suite.app.BankKeeper.BurnCoins(suite.ctx, "tokenization", coins))
	suite.Require().NoError(DeleteCollection(&suite.TestSuite, suite.ctx, msg))
}

func (suite *PoolIntegrationTestSuite) TestDeleteCollectionRefusesRegisteredAliasPoolWithNoBankSupply() {
	collection, denom, err := suite.createCollectionWithAliasPath(bob, "pooled-alias")
	suite.Require().NoError(err)
	pool, err := balancer.NewBalancerPool(1, balancer.PoolParams{SwapFee: sdkmath.LegacyZeroDec(), ExitFee: sdkmath.LegacyZeroDec()}, []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin(denom, 100), Weight: sdkmath.OneInt()},
		{Token: sdk.NewInt64Coin("ubadge", 100), Weight: sdkmath.OneInt()},
	}, suite.ctx.BlockTime())
	suite.Require().NoError(err)
	suite.Require().NoError(suite.app.GammKeeper.OverwritePoolV15MigrationUnsafe(suite.ctx, &pool))
	suite.Require().True(suite.app.BankKeeper.GetSupply(suite.ctx, denom).IsZero())
	suite.Require().ErrorContains(DeleteCollection(&suite.TestSuite, suite.ctx, &types.MsgDeleteCollection{Creator: bob, CollectionId: collection.CollectionId}), "pool")
	_, found := suite.app.TokenizationKeeper.GetCollectionFromStore(suite.ctx, collection.CollectionId)
	suite.Require().True(found)
}
