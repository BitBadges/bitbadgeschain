package gamm_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
)

// EdgeCasesTestSuite provides tests for boundary conditions and edge cases
type EdgeCasesTestSuite struct {
	apptesting.KeeperTestHelper

	Precompile *gamm.Precompile
	PoolId     uint64
}

func TestEdgeCasesTestSuite(t *testing.T) {
	suite.Run(t, new(EdgeCasesTestSuite))
}

func (suite *EdgeCasesTestSuite) SetupTest() {
	suite.Reset()
	suite.Precompile = gamm.NewPrecompile(suite.App.GammKeeper)

	// Create a test pool for edge case testing
	alice := suite.TestAccs[0]
	poolCreationCoins := sdk.NewCoins(
		sdk.NewCoin("uatom", osmomath.NewInt(2_000_000_000_000_000_000)),
		sdk.NewCoin("uosmo", osmomath.NewInt(2_000_000_000_000_000_000)),
	)
	suite.FundAcc(alice, poolCreationCoins)

	oneTrillion := osmomath.NewInt(1e12)
	poolAssets := []balancer.PoolAsset{
		{Token: sdk.NewCoin("uatom", oneTrillion), Weight: osmomath.NewInt(100)},
		{Token: sdk.NewCoin("uosmo", oneTrillion), Weight: osmomath.NewInt(100)},
	}

	poolParams := balancer.PoolParams{
		SwapFee: osmomath.MustNewDecFromStr("0.025"),
		ExitFee: osmomath.ZeroDec(),
	}

	msg := balancer.NewMsgCreateBalancerPool(alice, poolParams, poolAssets)
	poolId, err := suite.App.PoolManagerKeeper.CreatePool(suite.Ctx, msg)
	suite.Require().NoError(err)
	suite.PoolId = poolId
}

// TestEdgeCases_MinimumAmounts tests with minimum valid amounts (1 wei)
func (suite *EdgeCasesTestSuite) TestEdgeCases_MinimumAmounts() {
	// Test with minimum share amount
	minShareAmount := big.NewInt(1)
	err := gamm.ValidateShareAmount(minShareAmount, "shareAmount")
	suite.NoError(err, "Minimum share amount (1) should be valid")

	// Test with minimum coin amount
	minCoin := struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	}{
		Denom:  "uatom",
		Amount: big.NewInt(1),
	}
	err = gamm.ValidateCoin(minCoin, "coin")
	suite.NoError(err, "Minimum coin amount (1) should be valid")
}

// TestEdgeCases_MaximumAmounts tests with maximum valid amounts
func (suite *EdgeCasesTestSuite) TestEdgeCases_MaximumAmounts() {
	// Test with maximum int256 value (2^255-1)
	maxInt256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1))
	err := gamm.ValidateShareAmount(maxInt256, "shareAmount")
	suite.NoError(err, "Maximum int256 value should be valid")

	// Test overflow protection
	overflowValue := new(big.Int).Lsh(big.NewInt(1), 256) // 2^256, exceeds int256
	err = gamm.CheckOverflow(overflowValue, "amount")
	suite.Error(err, "Value exceeding int256 max should cause overflow error")
}

// TestEdgeCases_EmptyPools tests operations on newly created pools
func (suite *EdgeCasesTestSuite) TestEdgeCases_EmptyPools() {
	// Verify pool exists and is accessible
	pool, err := suite.App.GammKeeper.GetPoolAndPoke(suite.Ctx, suite.PoolId)
	suite.NoError(err)
	suite.NotNil(pool)

	// Verify pool has initial liquidity
	liquidity := pool.GetTotalPoolLiquidity(suite.Ctx)
	suite.Greater(len(liquidity), 0, "Pool should have initial liquidity")
}

// TestEdgeCases_PoolNotFound tests all methods with non-existent pool IDs
func (suite *EdgeCasesTestSuite) TestEdgeCases_PoolNotFound() {
	nonExistentPoolId := uint64(99999)

	// Test ValidatePoolId with non-existent pool
	err := gamm.ValidatePoolId(nonExistentPoolId)
	suite.NoError(err, "ValidatePoolId should not check pool existence")

	// Test that GetPool would fail with non-existent pool
	// (This would be tested in handler tests, not validation tests)
}

// TestEdgeCases_ZeroShares tests exit pool with zero shares
func (suite *EdgeCasesTestSuite) TestEdgeCases_ZeroShares() {
	// Zero shares should be valid (allows zero for some operations)
	// Note: ValidateShareAmount requires > 0, but CheckOverflow allows zero
	zeroShares := big.NewInt(0)
	err := gamm.CheckOverflow(zeroShares, "shareAmount")
	suite.NoError(err, "Zero share amount should not overflow")
}

// TestEdgeCases_ExactAmounts tests swaps with exact minimum amounts
func (suite *EdgeCasesTestSuite) TestEdgeCases_ExactAmounts() {
	// Test with exact minimum token amount
	minAmount := big.NewInt(1)
	err := gamm.ValidateShareAmount(minAmount, "tokenOutMinAmount")
	suite.NoError(err, "Exact minimum amount should be valid")
}

// TestEdgeCases_MultipleRoutes tests swaps with maximum routes (10)
func (suite *EdgeCasesTestSuite) TestEdgeCases_MultipleRoutes() {
	// Test with maximum number of routes
	maxRoutes := make([]struct {
		PoolId        uint64 `json:"poolId"`
		TokenOutDenom string `json:"tokenOutDenom"`
	}, gamm.MaxRoutes)
		for i := range maxRoutes {
		maxRoutes[i] = struct {
			PoolId        uint64 `json:"poolId"`
			TokenOutDenom string `json:"tokenOutDenom"`
		}{
			PoolId:        suite.PoolId,
			TokenOutDenom: "uosmo",
		}
	}

	err := gamm.ValidateRoutes(maxRoutes, "routes")
	suite.NoError(err, "Maximum number of routes should be valid")

	// Test with routes exceeding maximum
	tooManyRoutes := make([]struct {
		PoolId        uint64 `json:"poolId"`
		TokenOutDenom string `json:"tokenOutDenom"`
	}, gamm.MaxRoutes+1)
	err = gamm.ValidateRoutes(tooManyRoutes, "routes")
	suite.Error(err, "Routes exceeding maximum should be invalid")
}

// TestEdgeCases_MultipleCoins tests with maximum coins (20)
func (suite *EdgeCasesTestSuite) TestEdgeCases_MultipleCoins() {
	// Test with maximum number of coins
	maxCoins := make([]struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	}, gamm.MaxCoins)
		for i := range maxCoins {
		maxCoins[i] = struct {
			Denom  string   `json:"denom"`
			Amount *big.Int `json:"amount"`
		}{
			Denom:  "uatom",
			Amount: big.NewInt(1000),
		}
	}

	err := gamm.ValidateCoins(maxCoins, "coins")
	suite.NoError(err, "Maximum number of coins should be valid")

	// Test with coins exceeding maximum
	tooManyCoins := make([]struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	}, gamm.MaxCoins+1)
	err = gamm.ValidateCoins(tooManyCoins, "coins")
	suite.Error(err, "Coins exceeding maximum should be invalid")
}

// TestEdgeCases_InvalidDenoms tests with invalid denomination strings
func (suite *EdgeCasesTestSuite) TestEdgeCases_InvalidDenoms() {
	testCases := []struct {
		name  string
		denom string
		valid bool
	}{
		{"empty denom", "", false},
		{"too long denom", string(make([]byte, gamm.MaxStringLength+1)), false},
		{"valid denom", "uatom", true},
		{"valid denom with numbers", "uatom123", true},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			coin := struct {
				Denom  string   `json:"denom"`
				Amount *big.Int `json:"amount"`
			}{
				Denom:  tc.denom,
				Amount: big.NewInt(1000),
			}
			err := gamm.ValidateCoin(coin, "coin")
			if tc.valid {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}
		})
	}
}

// TestEdgeCases_OverflowProtection tests overflow protection for large numbers
func (suite *EdgeCasesTestSuite) TestEdgeCases_OverflowProtection() {
	testCases := []struct {
		name  string
		value *big.Int
		valid bool
	}{
		{"max int256", new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(1)), true},
		{"exceeds int256", new(big.Int).Lsh(big.NewInt(1), 256), false},
		{"negative value", big.NewInt(-1), false},
		{"zero", big.NewInt(0), true},
		{"small positive", big.NewInt(1), true},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := gamm.CheckOverflow(tc.value, "amount")
			if tc.valid {
				suite.NoError(err)
			} else {
				suite.Error(err)
			}
		})
	}
}

// TestEdgeCases_ZeroPoolId tests validation with zero pool ID
func (suite *EdgeCasesTestSuite) TestEdgeCases_ZeroPoolId() {
	err := gamm.ValidatePoolId(0)
	suite.Error(err, "Zero pool ID should be invalid")
	suite.Contains(err.Error(), "cannot be zero")
}

// TestEdgeCases_StringLength tests string length validation
func (suite *EdgeCasesTestSuite) TestEdgeCases_StringLength() {
	// Test with maximum allowed string length
	maxString := string(make([]byte, gamm.MaxStringLength))
	err := gamm.ValidateStringLength(maxString, "field")
	suite.NoError(err, "Maximum string length should be valid")

	// Test with string exceeding maximum
	tooLongString := string(make([]byte, gamm.MaxStringLength+1))
	err = gamm.ValidateStringLength(tooLongString, "field")
	suite.Error(err, "String exceeding maximum length should be invalid")
}

// TestEdgeCases_EmptyArrays tests with empty arrays (where allowed)
func (suite *EdgeCasesTestSuite) TestEdgeCases_EmptyArrays() {
	// Empty coins array should be rejected by ValidateCoins (required)
	emptyCoins := []struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	}{}
	err := gamm.ValidateCoins(emptyCoins, "coins")
	suite.Error(err, "Empty coins array should be rejected by ValidateCoins")
	suite.Contains(err.Error(), "cannot be empty", "Error should mention empty array")

	// Empty routes array should be invalid (required)
	emptyRoutes := []struct {
		PoolId        uint64 `json:"poolId"`
		TokenOutDenom string `json:"tokenOutDenom"`
	}{}
	err = gamm.ValidateRoutes(emptyRoutes, "routes")
	suite.Error(err, "Empty routes array should be invalid")
	suite.Contains(err.Error(), "cannot be empty", "Error should mention empty array")

	// Empty coins array should be allowed by ValidateCoinsAllowZero (optional minimums)
	err = gamm.ValidateCoinsAllowZero(emptyCoins, "tokenOutMins")
	suite.NoError(err, "Empty coins array should be allowed for tokenOutMins (no minimums)")
}

