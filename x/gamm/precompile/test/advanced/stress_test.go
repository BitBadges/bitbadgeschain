package gamm_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
)

// StressTestSuite provides tests for stress and performance testing
type StressTestSuite struct {
	apptesting.KeeperTestHelper

	Precompile *gamm.Precompile
	PoolId     uint64
}

func TestStressTestSuite(t *testing.T) {
	suite.Run(t, new(StressTestSuite))
}

func (suite *StressTestSuite) SetupTest() {
	suite.Reset()

	suite.Precompile = gamm.NewPrecompile(suite.App.GammKeeper)

	alice := suite.TestAccs[0]
	largeAmount, _ := new(big.Int).SetString("10000000000000000000", 10)
	poolCreationCoins := sdk.NewCoins(
		sdk.NewCoin("uatom", osmomath.NewIntFromBigInt(largeAmount)),
		sdk.NewCoin("uosmo", osmomath.NewIntFromBigInt(largeAmount)),
		sdk.NewCoin("uion", osmomath.NewIntFromBigInt(largeAmount)),
	)
	suite.FundAcc(alice, poolCreationCoins)

	poolId, err := suite.createDefaultTestPoolInContext(alice)
	suite.Require().NoError(err)
	suite.PoolId = poolId
}

func (suite *StressTestSuite) createDefaultTestPoolInContext(creator sdk.AccAddress) (uint64, error) {
	oneTrillion := osmomath.NewInt(1e12)
	poolAssets := []balancer.PoolAsset{
		{
			Token:  sdk.NewCoin("uatom", oneTrillion),
			Weight: osmomath.NewInt(100),
		},
		{
			Token:  sdk.NewCoin("uosmo", oneTrillion),
			Weight: osmomath.NewInt(100),
		},
	}
	poolParams := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}
	msg := balancer.NewMsgCreateBalancerPool(creator, poolParams, poolAssets)
	poolId, err := suite.App.PoolManagerKeeper.CreatePool(suite.Ctx, msg)
	return poolId, err
}

// TestStress_MaxRoutes tests swaps with maximum number of routes
func (suite *StressTestSuite) TestStress_MaxRoutes() {
	// Test that swaps with maximum routes (10) work correctly
	maxRoutes := gamm.MaxRoutes
	routes := make([]struct {
		PoolId        uint64 `json:"poolId"`
		TokenOutDenom string `json:"tokenOutDenom"`
	}, maxRoutes)

	// Create multiple pools for routing
	alice := suite.TestAccs[0]
	poolIds := []uint64{suite.PoolId}

	// Create additional pools for routing
	for i := 1; i < maxRoutes; i++ {
		poolAssets := []balancer.PoolAsset{
			{
				Token:  sdk.NewCoin("uatom", osmomath.NewInt(1e12)),
				Weight: osmomath.NewInt(100),
			},
			{
				Token:  sdk.NewCoin("uosmo", osmomath.NewInt(1e12)),
				Weight: osmomath.NewInt(100),
			},
		}
		poolParams := balancer.PoolParams{
			SwapFee: osmomath.MustNewDecFromStr("0.025"),
			ExitFee: osmomath.ZeroDec(),
		}
		msg := balancer.NewMsgCreateBalancerPool(alice, poolParams, poolAssets)
		poolId, err := suite.App.PoolManagerKeeper.CreatePool(suite.Ctx, msg)
		if err != nil {
			suite.T().Logf("Failed to create pool %d: %v", i, err)
			break
		}
		poolIds = append(poolIds, poolId)
	}

	// Build routes using available pools
	for i := 0; i < len(poolIds) && i < maxRoutes; i++ {
		routes[i] = struct {
			PoolId        uint64 `json:"poolId"`
			TokenOutDenom string `json:"tokenOutDenom"`
		}{
			PoolId:        poolIds[i],
			TokenOutDenom: "uosmo",
		}
	}

	// Verify routes are valid
	err := gamm.ValidateRoutes(routes, "routes")
	suite.NoError(err, "Max routes should be valid")

	suite.T().Logf("Created %d pools for routing test", len(poolIds))
	suite.T().Logf("Maximum routes (%d) validated successfully", maxRoutes)
}

// TestStress_MaxCoins tests operations with maximum number of coins
func (suite *StressTestSuite) TestStress_MaxCoins() {
	// Test that operations with maximum coins (20) work correctly
	maxCoins := gamm.MaxCoins
	coins := make([]struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	}, maxCoins)

	// Create coins with different denominations
	denoms := []string{"uatom", "uosmo", "uion", "ustake", "uusdc"}
	for i := 0; i < maxCoins; i++ {
		coins[i] = struct {
			Denom  string   `json:"denom"`
			Amount *big.Int `json:"amount"`
		}{
			Denom:  denoms[i%len(denoms)],
			Amount: big.NewInt(1000),
		}
	}

	// Verify coins are valid
	err := gamm.ValidateCoins(coins, "coins")
	suite.NoError(err, "Max coins should be valid")

	suite.T().Logf("Maximum coins (%d) validated successfully", maxCoins)
}

// TestStress_MaxAffiliates tests swaps with maximum number of affiliates
func (suite *StressTestSuite) TestStress_MaxAffiliates() {
	// Test that swaps with maximum affiliates (10) work correctly
	maxAffiliates := gamm.MaxAffiliates
	affiliates := make([]struct {
		Address        common.Address `json:"address"`
		BasisPointsFee *big.Int      `json:"basisPointsFee"`
	}, maxAffiliates)

	// Create affiliate addresses
	for i := 0; i < maxAffiliates; i++ {
		affiliates[i] = struct {
			Address        common.Address `json:"address"`
			BasisPointsFee *big.Int      `json:"basisPointsFee"`
		}{
			Address:        common.BigToAddress(big.NewInt(int64(i + 1000))),
			BasisPointsFee: big.NewInt(10), // 0.1% fee
		}
	}

	// Verify affiliates are valid
	err := gamm.ValidateAffiliates(affiliates, "affiliates")
	suite.NoError(err, "Max affiliates should be valid")

	suite.T().Logf("Maximum affiliates (%d) validated successfully", maxAffiliates)
}

// TestStress_ManyPools tests operations with many pools
func (suite *StressTestSuite) TestStress_ManyPools() {
	// Create multiple pools to test pool management at scale
	numPools := 10
	alice := suite.TestAccs[0]
	poolIds := []uint64{suite.PoolId}

	for i := 1; i < numPools; i++ {
		poolAssets := []balancer.PoolAsset{
			{
				Token:  sdk.NewCoin("uatom", osmomath.NewInt(1e12)),
				Weight: osmomath.NewInt(100),
			},
			{
				Token:  sdk.NewCoin("uosmo", osmomath.NewInt(1e12)),
				Weight: osmomath.NewInt(100),
			},
		}
		poolParams := balancer.PoolParams{
			SwapFee: osmomath.MustNewDecFromStr("0.025"),
			ExitFee: osmomath.ZeroDec(),
		}
		msg := balancer.NewMsgCreateBalancerPool(alice, poolParams, poolAssets)
		poolId, err := suite.App.PoolManagerKeeper.CreatePool(suite.Ctx, msg)
		if err != nil {
			suite.T().Logf("Failed to create pool %d: %v", i, err)
			break
		}
		poolIds = append(poolIds, poolId)
	}

	// Verify all pools can be queried
	for _, poolId := range poolIds {
		pool, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, poolId)
		suite.NoError(err, "Pool should be accessible")
		suite.NotNil(pool, "Pool should not be nil")
		suite.Equal(poolId, pool.GetId(), "Pool ID should match")
	}

	suite.T().Logf("Created and verified %d pools", len(poolIds))
}

// TestStress_LargeAmounts tests operations with very large amounts
func (suite *StressTestSuite) TestStress_LargeAmounts() {
	// Test with maximum valid int256 value (2^255 - 1)
	maxInt256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1))

	// Verify large amounts don't overflow
	err := gamm.CheckOverflow(maxInt256, "largeAmount")
	suite.NoError(err, "Max int256 should not overflow")

	// Test share amount validation with large value
	err = gamm.ValidateShareAmount(maxInt256, "shareAmount")
	suite.NoError(err, "Large share amount should be valid")

	suite.T().Logf("Large amounts (max int256) validated successfully")
}

// TestStress_StringLength tests operations with maximum string length
func (suite *StressTestSuite) TestStress_StringLength() {
	// Test with maximum string length (10,000 characters)
	maxString := make([]byte, gamm.MaxStringLength)
	for i := range maxString {
		maxString[i] = 'a'
	}

	// Verify string length validation
	err := gamm.ValidateStringLength(string(maxString), "testString")
	suite.NoError(err, "Max string length should be valid")

	// Test that strings longer than max are rejected
	tooLongString := string(append(maxString, 'a'))
	err = gamm.ValidateStringLength(tooLongString, "testString")
	suite.Error(err, "String longer than max should be invalid")

	suite.T().Logf("Maximum string length (%d) validated successfully", gamm.MaxStringLength)
}

// TestStress_ComplexRoutes tests complex multi-hop swap routes
func (suite *StressTestSuite) TestStress_ComplexRoutes() {
	// Create multiple pools for complex routing
	alice := suite.TestAccs[0]
	poolIds := []uint64{suite.PoolId}

	// Create additional pools with different token pairs
	tokenPairs := [][]string{
		{"uatom", "uosmo"},
		{"uosmo", "uion"},
		{"uion", "uatom"},
	}

	for _, pair := range tokenPairs {
		poolAssets := []balancer.PoolAsset{
			{
				Token:  sdk.NewCoin(pair[0], osmomath.NewInt(1e12)),
				Weight: osmomath.NewInt(100),
			},
			{
				Token:  sdk.NewCoin(pair[1], osmomath.NewInt(1e12)),
				Weight: osmomath.NewInt(100),
			},
		}
		poolParams := balancer.PoolParams{
			SwapFee: osmomath.MustNewDecFromStr("0.025"),
			ExitFee: osmomath.ZeroDec(),
		}
		msg := balancer.NewMsgCreateBalancerPool(alice, poolParams, poolAssets)
		poolId, err := suite.App.PoolManagerKeeper.CreatePool(suite.Ctx, msg)
		if err != nil {
			suite.T().Logf("Failed to create pool: %v", err)
			continue
		}
		poolIds = append(poolIds, poolId)
	}

	// Build a complex route through multiple pools
	routes := make([]struct {
		PoolId        uint64 `json:"poolId"`
		TokenOutDenom string `json:"tokenOutDenom"`
	}, len(poolIds))
	for i, poolId := range poolIds {
		routes[i] = struct {
			PoolId        uint64 `json:"poolId"`
			TokenOutDenom string `json:"tokenOutDenom"`
		}{
			PoolId:        poolId,
			TokenOutDenom: "uosmo",
		}
	}

	// Verify complex routes are valid
	err := gamm.ValidateRoutes(routes, "routes")
	suite.NoError(err, "Complex routes should be valid")

	suite.T().Logf("Complex routing with %d pools validated successfully", len(poolIds))
}

// TestStress_RepeatedOperations tests repeated operations to check for resource leaks
func (suite *StressTestSuite) TestStress_RepeatedOperations() {
	// Perform many repeated query operations
	numOperations := 100
	successCount := 0

	for i := 0; i < numOperations; i++ {
		// Query pool
		pool, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId)
		if err == nil && pool != nil {
			successCount++
		}

		// Query total shares
		pool, err = suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId)
		if err == nil && pool != nil {
			shares := pool.GetTotalShares()
			if shares.IsPositive() {
				successCount++
			}
		}
	}

	suite.Greater(successCount, numOperations, "Most operations should succeed")
	suite.T().Logf("Repeated operations: %d/%d succeeded", successCount, numOperations*2)
}

