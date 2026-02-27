//go:build test
// +build test

package e2e

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// PoolTestSuite tests GAMM pool operations
type PoolTestSuite struct {
	CrossModuleTestSuite
}

func TestPoolTestSuite(t *testing.T) {
	suite.Run(t, new(PoolTestSuite))
}

// TestCreateBalancerPool tests creating a Balancer pool
func (s *PoolTestSuite) TestCreateBalancerPool() {
	s.T().Log("Testing Balancer pool creation")

	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(1000000)),
	)

	s.Require().Greater(poolId, uint64(0))
	s.T().Logf("Created pool with ID: %d", poolId)

	// Verify pool exists
	pool, err := s.App.GammKeeper.GetPoolAndPoke(s.Ctx, poolId)
	s.Require().NoError(err)
	s.Require().NotNil(pool)

	// Verify pool denoms
	denoms, err := s.App.GammKeeper.GetPoolDenoms(s.Ctx, poolId)
	s.Require().NoError(err)
	s.Require().Len(denoms, 2)
	s.T().Logf("Pool denoms: %v", denoms)
}

// TestSwapExactAmountIn tests swapping tokens in a pool
func (s *PoolTestSuite) TestSwapExactAmountIn() {
	s.T().Log("Testing SwapExactAmountIn")

	// Create pool
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)

	sender := s.TestAccs[1]
	tokenIn := sdk.NewCoin("foo", sdkmath.NewInt(100000))

	// Get initial balance
	initialBarBalance := s.GetCoinBalance(sender, "bar")

	// Perform swap
	tokenOut, err := s.SwapTokens(sender, poolId, tokenIn, "bar")
	s.Require().NoError(err)
	s.Require().True(tokenOut.Amount.GT(sdkmath.NewInt(0)))

	// Verify balance changed
	finalBarBalance := s.GetCoinBalance(sender, "bar")
	s.Require().True(finalBarBalance.Amount.GT(initialBarBalance.Amount))

	s.T().Logf("Swapped %s for %s", tokenIn, tokenOut)
}

// TestSwapExactAmountOut tests swapping to receive exact output
func (s *PoolTestSuite) TestSwapExactAmountOut() {
	s.T().Log("Testing SwapExactAmountOut")

	// Create pool with liquidity
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)

	sender := s.TestAccs[1]
	tokenOutAmount := sdkmath.NewInt(50000)

	// Fund sender with enough tokens
	s.FundAcc(sender, sdk.NewCoins(sdk.NewCoin("foo", sdkmath.NewInt(1000000))))

	initialFooBalance := s.GetCoinBalance(sender, "foo")

	msg := &gammtypes.MsgSwapExactAmountOut{
		Sender: sender.String(),
		Routes: []poolmanagertypes.SwapAmountOutRoute{
			{PoolId: poolId, TokenInDenom: "foo"},
		},
		TokenOut:         sdk.NewCoin("bar", tokenOutAmount),
		TokenInMaxAmount: sdkmath.NewInt(1000000),
	}

	_, err := s.gammMsgServer.SwapExactAmountOut(s.Ctx, msg)
	s.Require().NoError(err)

	// Verify we received exactly the expected amount
	finalBarBalance := s.GetCoinBalance(sender, "bar")
	s.Require().Equal(tokenOutAmount, finalBarBalance.Amount)

	// Verify foo was spent
	finalFooBalance := s.GetCoinBalance(sender, "foo")
	s.Require().True(finalFooBalance.Amount.LT(initialFooBalance.Amount))

	s.T().Logf("Received exactly %s bar, spent %s foo", tokenOutAmount, initialFooBalance.Amount.Sub(finalFooBalance.Amount))
}

// TestJoinPool tests joining a pool with liquidity
func (s *PoolTestSuite) TestJoinPool() {
	s.T().Log("Testing JoinPool")

	// Create pool
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)

	joiner := s.TestAccs[1]

	// Fund joiner
	s.FundAcc(joiner, sdk.NewCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(1000000)),
	))

	// Get pool share denom
	pool, err := s.App.GammKeeper.GetPoolAndPoke(s.Ctx, poolId)
	s.Require().NoError(err)
	lpDenom := gammtypes.GetPoolShareDenom(poolId)

	initialLPBalance := s.GetCoinBalance(joiner, lpDenom)
	s.Require().True(initialLPBalance.IsZero())

	// Join pool
	shareOutAmount := pool.GetTotalShares().Quo(sdkmath.NewInt(100)) // 1% of pool
	msg := &gammtypes.MsgJoinPool{
		Sender:         joiner.String(),
		PoolId:         poolId,
		ShareOutAmount: shareOutAmount,
		TokenInMaxs: sdk.NewCoins(
			sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
			sdk.NewCoin("bar", sdkmath.NewInt(1000000)),
		),
	}

	_, err = s.gammMsgServer.JoinPool(s.Ctx, msg)
	s.Require().NoError(err)

	// Verify LP tokens received
	finalLPBalance := s.GetCoinBalance(joiner, lpDenom)
	s.Require().Equal(shareOutAmount, finalLPBalance.Amount)

	s.T().Logf("Joined pool, received %s LP tokens", finalLPBalance)
}

// TestExitPool tests exiting a pool
func (s *PoolTestSuite) TestExitPool() {
	s.T().Log("Testing ExitPool")

	// Create pool
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)

	// Pool creator (TestAccs[0]) has LP tokens from pool creation
	exiter := s.TestAccs[0]
	lpDenom := gammtypes.GetPoolShareDenom(poolId)

	initialLPBalance := s.GetCoinBalance(exiter, lpDenom)
	s.Require().True(initialLPBalance.Amount.GT(sdkmath.NewInt(0)), "should have LP tokens")

	// Exit with a portion
	shareInAmount := initialLPBalance.Amount.Quo(sdkmath.NewInt(10)) // 10%

	msg := &gammtypes.MsgExitPool{
		Sender:        exiter.String(),
		PoolId:        poolId,
		ShareInAmount: shareInAmount,
		TokenOutMins:  sdk.NewCoins(),
	}

	_, err := s.gammMsgServer.ExitPool(s.Ctx, msg)
	s.Require().NoError(err)

	// Verify LP tokens decreased
	finalLPBalance := s.GetCoinBalance(exiter, lpDenom)
	s.Require().Equal(initialLPBalance.Amount.Sub(shareInAmount), finalLPBalance.Amount)

	s.T().Logf("Exited pool, LP balance went from %s to %s", initialLPBalance, finalLPBalance)
}

// TestMultiHopSwap tests swapping through multiple pools
func (s *PoolTestSuite) TestMultiHopSwap() {
	s.T().Log("Testing multi-hop swap")

	// Create two pools: foo<->bar and bar<->baz
	poolId1 := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)

	poolId2 := s.CreatePoolWithCoins(
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
		sdk.NewCoin("baz", sdkmath.NewInt(10000000)),
	)

	sender := s.TestAccs[1]
	tokenIn := sdk.NewCoin("foo", sdkmath.NewInt(100000))
	s.FundAcc(sender, sdk.NewCoins(tokenIn))

	// Swap foo -> bar -> baz
	msg := &gammtypes.MsgSwapExactAmountIn{
		Sender: sender.String(),
		Routes: []poolmanagertypes.SwapAmountInRoute{
			{PoolId: poolId1, TokenOutDenom: "bar"},
			{PoolId: poolId2, TokenOutDenom: "baz"},
		},
		TokenIn:           tokenIn,
		TokenOutMinAmount: sdkmath.NewInt(1),
	}

	res, err := s.gammMsgServer.SwapExactAmountIn(s.Ctx, msg)
	s.Require().NoError(err)

	// Verify we got baz
	finalBazBalance := s.GetCoinBalance(sender, "baz")
	s.Require().Equal(res.TokenOutAmount, finalBazBalance.Amount)

	s.T().Logf("Multi-hop swap: %s foo -> %s baz", tokenIn.Amount, res.TokenOutAmount)
}

// TestPoolSpotPrice tests spot price calculation
func (s *PoolTestSuite) TestPoolSpotPrice() {
	s.T().Log("Testing pool spot price")

	// Create imbalanced pool (more bar than foo)
	poolId := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(1000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(2000000)),
	)

	// Calculate spot price - the direction depends on AMM formula
	spotPrice, err := s.App.GammKeeper.CalculateSpotPrice(s.Ctx, poolId, "foo", "bar")
	s.Require().NoError(err)
	s.Require().True(spotPrice.GT(osmomath.ZeroBigDec()), "Spot price should be positive")
	s.T().Logf("Spot price (foo/bar): %s", spotPrice)

	reverseSpotPrice, err := s.App.GammKeeper.CalculateSpotPrice(s.Ctx, poolId, "bar", "foo")
	s.Require().NoError(err)
	s.Require().True(reverseSpotPrice.GT(osmomath.ZeroBigDec()), "Reverse spot price should be positive")
	s.T().Logf("Spot price (bar/foo): %s", reverseSpotPrice)

	// The product of spot prices should be approximately 1 (inverse relationship)
	product := spotPrice.Mul(reverseSpotPrice)
	s.Require().True(product.GT(osmomath.NewBigDec(99).Quo(osmomath.NewBigDec(100))), "Product should be close to 1")
	s.Require().True(product.LT(osmomath.NewBigDec(101).Quo(osmomath.NewBigDec(100))), "Product should be close to 1")
}

// TestPoolWithSwapFee tests pool with swap fees
func (s *PoolTestSuite) TestPoolWithSwapFee() {
	s.T().Log("Testing pool with swap fee")

	// Create pool with 1% swap fee
	poolAssets := []balancer.PoolAsset{
		{Weight: sdkmath.NewInt(1), Token: sdk.NewCoin("foo", sdkmath.NewInt(10000000))},
		{Weight: sdkmath.NewInt(1), Token: sdk.NewCoin("bar", sdkmath.NewInt(10000000))},
	}

	fundCoins := sdk.NewCoins(
		sdk.NewCoin("ubadge", sdkmath.NewInt(10000000000)),
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)
	s.FundAcc(s.TestAccs[0], fundCoins)

	msg := balancer.NewMsgCreateBalancerPool(s.TestAccs[0], balancer.PoolParams{
		SwapFee: osmomath.NewDecWithPrec(1, 2), // 1%
		ExitFee: osmomath.ZeroDec(),
	}, poolAssets)

	poolIdWithFee, err := s.App.PoolManagerKeeper.CreatePool(s.Ctx, msg)
	s.Require().NoError(err)

	// Create identical pool without fee for comparison
	poolIdNoFee := s.CreatePoolWithCoins(
		sdk.NewCoin("foo", sdkmath.NewInt(10000000)),
		sdk.NewCoin("bar", sdkmath.NewInt(10000000)),
	)

	sender := s.TestAccs[1]
	tokenIn := sdk.NewCoin("foo", sdkmath.NewInt(100000))

	// Swap in pool with fee
	s.FundAcc(sender, sdk.NewCoins(tokenIn))
	res1, err := s.gammMsgServer.SwapExactAmountIn(s.Ctx, &gammtypes.MsgSwapExactAmountIn{
		Sender:            sender.String(),
		Routes:            []poolmanagertypes.SwapAmountInRoute{{PoolId: poolIdWithFee, TokenOutDenom: "bar"}},
		TokenIn:           tokenIn,
		TokenOutMinAmount: sdkmath.NewInt(1),
	})
	s.Require().NoError(err)

	// Swap in pool without fee
	s.FundAcc(sender, sdk.NewCoins(tokenIn))
	res2, err := s.gammMsgServer.SwapExactAmountIn(s.Ctx, &gammtypes.MsgSwapExactAmountIn{
		Sender:            sender.String(),
		Routes:            []poolmanagertypes.SwapAmountInRoute{{PoolId: poolIdNoFee, TokenOutDenom: "bar"}},
		TokenIn:           tokenIn,
		TokenOutMinAmount: sdkmath.NewInt(1),
	})
	s.Require().NoError(err)

	// Pool with fee should give less output
	s.Require().True(res1.TokenOutAmount.LT(res2.TokenOutAmount))
	s.T().Logf("With 1%% fee: %s bar, without fee: %s bar", res1.TokenOutAmount, res2.TokenOutAmount)
}

// TestMultiplePoolsAndSwaps tests creating multiple pools and swapping
func (s *PoolTestSuite) TestMultiplePoolsAndSwaps() {
	s.T().Log("Testing multiple pools and swaps")

	// Create several pools
	numPools := 5
	var poolIds []uint64

	for i := 0; i < numPools; i++ {
		denomA := fmt.Sprintf("denom%d", i)
		denomB := fmt.Sprintf("denom%d", i+1)

		// Fund with custom denoms
		s.FundAcc(s.TestAccs[0], sdk.NewCoins(
			sdk.NewCoin(denomA, sdkmath.NewInt(10000000)),
			sdk.NewCoin(denomB, sdkmath.NewInt(10000000)),
		))

		poolId := s.CreatePoolWithCoins(
			sdk.NewCoin(denomA, sdkmath.NewInt(5000000)),
			sdk.NewCoin(denomB, sdkmath.NewInt(5000000)),
		)
		poolIds = append(poolIds, poolId)
	}

	s.Require().Len(poolIds, numPools)
	s.T().Logf("Created %d pools: %v", numPools, poolIds)

	// Verify all pools exist and have correct denoms
	for i, poolId := range poolIds {
		pool, err := s.App.GammKeeper.GetPoolAndPoke(s.Ctx, poolId)
		s.Require().NoError(err)
		s.Require().NotNil(pool)

		denoms, err := s.App.GammKeeper.GetPoolDenoms(s.Ctx, poolId)
		s.Require().NoError(err)
		s.Require().Len(denoms, 2)
		s.T().Logf("Pool %d (ID %d): %v", i, poolId, denoms)
	}
}
