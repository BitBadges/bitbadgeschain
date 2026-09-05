package keeper_test

import (
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Sending badge tokens into a pool sets the auto-approval flags on the pool address only;
// the sender's own settings are left untouched.
func (suite *PoolIntegrationTestSuite) TestSendCoinsToPoolSetsPoolFlagsNotSenderFlags() {
	wctx := sdk.WrapSDKContext(suite.ctx)
	collection, aliasDenom, err := suite.createCollectionWithAliasPath(bob, "flagtest")
	suite.Require().Nil(err)

	// bob opts out of incoming transfers
	suite.Require().Nil(UpdateUserApprovals(&suite.TestSuite, wctx, &types.MsgUpdateUserApprovals{
		Creator:                               bob,
		CollectionId:                          collection.CollectionId,
		UpdateAutoApproveAllIncomingTransfers: true,
		AutoApproveAllIncomingTransfers:       false,
	}))
	before, err := GetUserBalance(&suite.TestSuite, wctx, collection.CollectionId, bob)
	suite.Require().Nil(err)
	suite.Require().False(before.AutoApproveAllIncomingTransfers)

	bobAcc := sdk.MustAccAddressFromBech32(bob)
	poolAssets := []balancer.PoolAsset{
		{Token: sdk.NewInt64Coin(aliasDenom, 1), Weight: osmomath.NewInt(1)},
		{Token: sdk.NewInt64Coin("ubadge", 100), Weight: osmomath.NewInt(1)},
	}
	poolId, err := suite.app.PoolManagerKeeper.CreatePool(suite.ctx,
		balancer.NewMsgCreateBalancerPool(bobAcc, balancer.NewPoolParams(osmomath.ZeroDec(), osmomath.ZeroDec()), poolAssets))
	suite.Require().Nil(err)

	after, err := GetUserBalance(&suite.TestSuite, wctx, collection.CollectionId, bob)
	suite.Require().Nil(err)
	suite.Require().False(after.AutoApproveAllIncomingTransfers, "sender opt-out must survive a pool deposit")

	pool, err := suite.app.GammKeeper.GetPool(suite.ctx, poolId)
	suite.Require().Nil(err)
	poolBalance, err := GetUserBalance(&suite.TestSuite, wctx, collection.CollectionId, pool.GetAddress().String())
	suite.Require().Nil(err)
	suite.Require().True(poolBalance.AutoApproveAllIncomingTransfers)
	suite.Require().True(poolBalance.AutoApproveSelfInitiatedOutgoingTransfers)
	suite.Require().True(poolBalance.AutoApproveSelfInitiatedIncomingTransfers)
	AssertBalancesEqual(&suite.TestSuite, []*types.Balance{{Amount: sdkmath.NewUint(1), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}}, poolBalance.Balances)
}
