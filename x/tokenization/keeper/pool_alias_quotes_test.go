package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	poolmanager "github.com/bitbadges/bitbadgeschain/x/poolmanager"
	pooltypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *PoolIntegrationTestSuite) setupAliasQuotePools() (string, uint64, uint64) {
	collection, alias, err := suite.createCollectionWithAliasPath(bob, "quotetest")
	suite.Require().NoError(err)
	collection.CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(5_000_000)
	suite.Require().NoError(suite.app.TokenizationKeeper.SetCollectionInStore(suite.ctx, collection, false))
	suite.Require().NoError(TransferTokens(&suite.TestSuite, sdk.WrapSDKContext(suite.ctx), &types.MsgTransferTokens{
		Creator: bob, CollectionId: collection.CollectionId,
		Transfers: []*types.Transfer{{From: "Mint", ToAddresses: []string{bob},
			Balances:             []*types.Balance{{Amount: sdkmath.NewUint(4_000_000), TokenIds: GetOneUintRange(), OwnershipTimes: GetFullUintRanges()}},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{{ApprovalId: "mint-test", ApprovalLevel: "collection", Version: sdkmath.NewUint(0)}},
		}},
	}))
	for _, addr := range []string{bob, alice} {
		coins := sdk.NewCoins(sdk.NewInt64Coin("quotea", 4_000_000), sdk.NewInt64Coin("quoteb", 4_000_000))
		suite.Require().NoError(suite.app.BankKeeper.MintCoins(suite.ctx, "mint", coins))
		suite.Require().NoError(suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, "mint", sdk.MustAccAddressFromBech32(addr), coins))
	}
	createPool := func(denom string) uint64 {
		assets := []balancer.PoolAsset{
			{Token: sdk.NewInt64Coin(alias, 1_000_000), Weight: osmomath.OneInt()},
			{Token: sdk.NewInt64Coin(denom, 1_000_000), Weight: osmomath.OneInt()},
		}
		id, err := suite.app.PoolManagerKeeper.CreatePool(suite.ctx, balancer.NewMsgCreateBalancerPool(
			sdk.MustAccAddressFromBech32(bob), balancer.NewPoolParams(osmomath.ZeroDec(), osmomath.ZeroDec()), assets))
		suite.Require().NoError(err)
		suite.app.PoolManagerKeeper.SetDenomPairTakerFee(suite.ctx, alias, denom, osmomath.MustNewDecFromStr("0.01"))
		suite.app.PoolManagerKeeper.SetDenomPairTakerFee(suite.ctx, denom, alias, osmomath.MustNewDecFromStr("0.01"))
		return id
	}
	return alias, createPool("quotea"), createPool("quoteb")
}

func (suite *PoolIntegrationTestSuite) TestAliasQuotesMatchExecutedSwaps() {
	for _, aliasInput := range []bool{true, false} {
		for _, exactIn := range []bool{true, false} {
			name := "native"
			if aliasInput {
				name = "alias"
			}
			if exactIn {
				name += "-exact-in"
			} else {
				name += "-exact-out"
			}
			suite.Run(name, func() {
				suite.SetupTest()
				alias, _, poolID := suite.setupAliasQuotePools()
				inDenom, outDenom := "quoteb", alias
				if aliasInput {
					inDenom, outDenom = alias, "quoteb"
				}
				k := suite.app.PoolManagerKeeper
				pool, err := suite.app.GammKeeper.GetPool(suite.ctx, poolID)
				suite.Require().NoError(err)
				sender := sdk.MustAccAddressFromBech32(bob)
				if exactIn {
					input := sdk.NewInt64Coin(inDenom, 10_000)
					netInput := input
					if !aliasInput {
						netInput, _ = poolmanager.CalcTakerFeeExactIn(input, osmomath.MustNewDecFromStr("0.01"))
					}
					raw, err := suite.app.GammKeeper.CalcOutAmtGivenIn(suite.ctx, pool, netInput, outDenom, osmomath.ZeroDec())
					suite.Require().NoError(err)
					route := []pooltypes.SwapAmountInRoute{{PoolId: poolID, TokenOutDenom: outDenom}}
					quote, err := k.MultihopEstimateOutGivenExactAmountIn(suite.ctx, route, input)
					suite.Require().NoError(err)
					suite.Require().Equal(raw.Amount.String(), quote.String())
					executed, err := k.RouteExactAmountIn(suite.ctx, sender, route, input, quote, nil)
					suite.Require().NoError(err)
					suite.Require().Equal(quote.String(), executed.String())
				} else {
					output := sdk.NewInt64Coin(outDenom, 10_000)
					raw, err := suite.app.GammKeeper.CalcInAmtGivenOut(suite.ctx, pool, output, inDenom, osmomath.ZeroDec())
					suite.Require().NoError(err)
					expected := raw
					if !aliasInput {
						expected, _ = poolmanager.CalcTakerFeeExactOut(raw, osmomath.MustNewDecFromStr("0.01"))
					}
					route := []pooltypes.SwapAmountOutRoute{{PoolId: poolID, TokenInDenom: inDenom}}
					quote, err := k.MultihopEstimateInGivenExactAmountOut(suite.ctx, route, output)
					suite.Require().NoError(err)
					suite.Require().Equal(expected.Amount.String(), quote.String())
					executed, err := k.RouteExactAmountOut(suite.ctx, sender, route, quote, output)
					suite.Require().NoError(err)
					suite.Require().Equal(quote.String(), executed.String())
				}
			})
		}
	}
}

func (suite *PoolIntegrationTestSuite) TestAliasIntermediateExactOutUsesOnlyRequiredInput() {
	alias, firstID, secondID := suite.setupAliasQuotePools()
	k := suite.app.PoolManagerKeeper
	first, err := suite.app.GammKeeper.GetPool(suite.ctx, firstID)
	suite.Require().NoError(err)
	second, err := suite.app.GammKeeper.GetPool(suite.ctx, secondID)
	suite.Require().NoError(err)
	output := sdk.NewInt64Coin("quoteb", 10_000)
	intermediate, err := suite.app.GammKeeper.CalcInAmtGivenOut(suite.ctx, second, output, alias, osmomath.ZeroDec())
	suite.Require().NoError(err)
	rawInput, err := suite.app.GammKeeper.CalcInAmtGivenOut(suite.ctx, first, intermediate, "quotea", osmomath.ZeroDec())
	suite.Require().NoError(err)
	expected, fee := poolmanager.CalcTakerFeeExactOut(rawInput, osmomath.MustNewDecFromStr("0.01"))
	suite.Require().True(fee.IsPositive(), "native first-hop fees remain charged")
	route := []pooltypes.SwapAmountOutRoute{{PoolId: firstID, TokenInDenom: "quotea"}, {PoolId: secondID, TokenInDenom: alias}}
	quote, err := k.MultihopEstimateInGivenExactAmountOut(suite.ctx, route, output)
	suite.Require().NoError(err)
	suite.Require().Equal(expected.Amount.String(), quote.String())
	sender := sdk.MustAccAddressFromBech32(alice)
	before := suite.app.BankKeeper.GetBalance(suite.ctx, sender, "quotea")
	beforeOut := suite.app.BankKeeper.GetBalance(suite.ctx, sender, "quoteb")
	failedCtx, _ := suite.ctx.CacheContext()
	_, err = k.RouteExactAmountOut(failedCtx, sender, route, expected.Amount.SubRaw(1), output)
	suite.Require().Error(err, "one below the total input including fees must fail")
	executed, err := k.RouteExactAmountOut(suite.ctx, sender, route, expected.Amount, output)
	suite.Require().NoError(err, "the exact required maximum must succeed")
	suite.Require().Equal(expected.Amount.String(), executed.String())
	suite.Require().Equal(expected.Amount.String(), before.Amount.Sub(suite.app.BankKeeper.GetBalance(suite.ctx, sender, "quotea").Amount).String())
	suite.Require().Equal(output.Amount.String(), suite.app.BankKeeper.GetBalance(suite.ctx, sender, "quoteb").Amount.Sub(beforeOut.Amount).String())
	remaining, err := suite.app.TokenizationKeeper.GetSpendableCoinAmountTokenizationLPOnly(suite.ctx, sender, alias)
	suite.Require().NoError(err)
	suite.Require().True(remaining.IsZero(), "no unnecessary intermediate alias tokens may remain")
	suite.Require().True(suite.app.BankKeeper.GetBalance(suite.ctx, sender, alias).IsZero())
}
