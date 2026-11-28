package gamm_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/precompile/test/helpers"
)

type QueryParityTestSuite struct {
	suite.Suite
	TestSuite  *helpers.TestSuite
	Precompile *gamm.Precompile
	PoolId     uint64
}

func TestQueryParityTestSuite(t *testing.T) {
	suite.Run(t, new(QueryParityTestSuite))
}

func (suite *QueryParityTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
	suite.Precompile = suite.TestSuite.Precompile

	// Fund Alice with uatom & uosmo for asymmetric pool creation and swaps
	suite.TestSuite.FundAccount(suite.TestSuite.Alice, sdk.NewCoins(
		sdk.NewCoin("uatom", osmomath.NewInt(10_000_000_000_000)),
		sdk.NewCoin("uosmo", osmomath.NewInt(10_000_000_000_000)),
	))

	// Create asymmetric pool: 1_000_000_000 uatom, 2_000_000_000 uosmo, equal weights -> spot price 2.0 uosmo/uatom
	poolId, err := suite.TestSuite.CreateTestBalancerPool(suite.TestSuite.Alice, []balancer.PoolAsset{
		{
			Token:  sdk.NewCoin("uatom", osmomath.NewInt(1_000_000_000)),
			Weight: osmomath.NewInt(100),
		},
		{
			Token:  sdk.NewCoin("uosmo", osmomath.NewInt(2_000_000_000)),
			Weight: osmomath.NewInt(100),
		},
	})
	suite.Require().NoError(err)
	suite.PoolId = poolId
}

func (suite *QueryParityTestSuite) callReadOnly(methodName string, payload interface{}) ([]interface{}, error) {
	method, found := suite.Precompile.ABI.Methods[methodName]
	suite.Require().True(found, "method %s must exist in ABI", methodName)

	jsonBytes, err := json.Marshal(payload)
	suite.Require().NoError(err)

	input, err := helpers.PackMethodWithJSON(&method, string(jsonBytes))
	suite.Require().NoError(err)

	precompileAddr := common.HexToAddress(gamm.GammPrecompileAddress)
	zeroVal, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(suite.TestSuite.AliceEVM, precompileAddr, zeroVal, 5_000_000, nil)
	contract.Input = input

	raw, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, true)
	if err != nil {
		return nil, err
	}
	return method.Outputs.Unpack(raw)
}

func (suite *QueryParityTestSuite) TestNewQueryMethods_MetadataAndGas() {
	methods := []string{
		gamm.SpotPriceMethod,
		gamm.GetTotalVolumeForPoolMethod,
		gamm.GetTakerFeesTrackerMethod,
		gamm.EstimateSwapExactAmountInMethod,
		gamm.EstimateSwapExactAmountOutMethod,
		gamm.GetTotalPoolLiquidityMethod,
		gamm.GetTradingPairTakerFeeMethod,
		gamm.GetNumPoolsMethod,
	}
	for _, name := range methods {
		method, found := suite.Precompile.ABI.Methods[name]
		suite.Require().True(found, "method %s must exist in ABI", name)
		suite.False(suite.Precompile.IsTransaction(&method), "method %s must not be a transaction", name)
		suite.Greater(suite.Precompile.RequiredGas(method.ID), uint64(0), "method %s must have RequiredGas", name)
	}
}

func (suite *QueryParityTestSuite) TestGetNumPools() {
	out, err := suite.callReadOnly(gamm.GetNumPoolsMethod, map[string]interface{}{})
	suite.Require().NoError(err)
	suite.Require().Len(out, 1)
	numPools, ok := out[0].(uint64)
	suite.Require().True(ok)
	suite.Equal(uint64(1), numPools)
}

func (suite *QueryParityTestSuite) TestSpotPrice() {
	// Matches x/gamm Querier.SpotPrice semantics (CalculateSpotPrice(poolId, base, quote) -> base/quote).
	out, err := suite.callReadOnly(gamm.SpotPriceMethod, map[string]interface{}{
		"pool_id":           suite.PoolId,
		"base_asset_denom":  "uatom",
		"quote_asset_denom": "uosmo",
	})
	suite.Require().NoError(err)
	suite.Require().Len(out, 1)
	spotPriceStr, ok := out[0].(string)
	suite.Require().True(ok)
	dec, err := osmomath.NewDecFromStr(spotPriceStr)
	suite.Require().NoError(err)
	suite.True(dec.Equal(osmomath.MustNewDecFromStr("0.5")), "expected spot price 0.5, got %s", spotPriceStr)

	// Reversed pair should return 2.0
	outRev, err := suite.callReadOnly(gamm.SpotPriceMethod, map[string]interface{}{
		"pool_id":           suite.PoolId,
		"base_asset_denom":  "uosmo",
		"quote_asset_denom": "uatom",
	})
	suite.Require().NoError(err)
	decRev, err := osmomath.NewDecFromStr(outRev[0].(string))
	suite.Require().NoError(err)
	suite.True(decRev.Equal(osmomath.MustNewDecFromStr("2.0")), "expected reverse spot price 2.0, got %s", outRev[0].(string))
}

func (suite *QueryParityTestSuite) TestGetTotalPoolLiquidity() {
	out, err := suite.callReadOnly(gamm.GetTotalPoolLiquidityMethod, map[string]interface{}{
		"pool_id": suite.PoolId,
	})
	suite.Require().NoError(err)
	suite.Require().Len(out, 1)

	coins, ok := out[0].([]struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	})
	suite.Require().True(ok)
	suite.Require().Len(coins, 2)
	coinMap := map[string]*big.Int{}
	for _, c := range coins {
		coinMap[c.Denom] = c.Amount
	}
	suite.Equal(big.NewInt(1_000_000_000), coinMap["uatom"])
	suite.Equal(big.NewInt(2_000_000_000), coinMap["uosmo"])
}

func (suite *QueryParityTestSuite) TestEstimateSwapExactAmountInAndOut() {
	// Estimate exact in: 1_000_000 uatom -> uosmo
	inOut, err := suite.callReadOnly(gamm.EstimateSwapExactAmountInMethod, map[string]interface{}{
		"pool_id":  suite.PoolId,
		"token_in": "1000000uatom",
		"routes": []map[string]interface{}{
			{"pool_id": suite.PoolId, "token_out_denom": "uosmo"},
		},
	})
	suite.Require().NoError(err)
	suite.Require().Len(inOut, 1)
	amountOut, ok := inOut[0].(*big.Int)
	suite.Require().True(ok)
	suite.True(amountOut.Sign() > 0, "expected positive estimated output")

	// Estimate exact out: 1_000_000 uosmo <- uatom
	outOut, err := suite.callReadOnly(gamm.EstimateSwapExactAmountOutMethod, map[string]interface{}{
		"pool_id":   suite.PoolId,
		"token_out": "1000000uosmo",
		"routes": []map[string]interface{}{
			{"pool_id": suite.PoolId, "token_in_denom": "uatom"},
		},
	})
	suite.Require().NoError(err)
	suite.Require().Len(outOut, 1)
	amountIn, ok := outOut[0].(*big.Int)
	suite.Require().True(ok)
	suite.True(amountIn.Sign() > 0, "expected positive estimated input")

	// Empty routes must error
	_, err = suite.callReadOnly(gamm.EstimateSwapExactAmountInMethod, map[string]interface{}{
		"pool_id":  suite.PoolId,
		"token_in": "1000000uatom",
		"routes":   []map[string]interface{}{},
	})
	suite.Error(err)
}

func (suite *QueryParityTestSuite) TestGetTotalVolumeForPool() {
	// Initially empty volume
	out, err := suite.callReadOnly(gamm.GetTotalVolumeForPoolMethod, map[string]interface{}{
		"pool_id": suite.PoolId,
	})
	suite.Require().NoError(err)
	coinsBefore, ok := out[0].([]struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	})
	suite.Require().True(ok)
	suite.Len(coinsBefore, 0)

	// Populate tracked volume in poolmanager store and verify precompile reads it back
	suite.TestSuite.App.PoolManagerKeeper.SetVolume(
		suite.TestSuite.Ctx,
		suite.PoolId,
		sdk.NewCoins(sdk.NewCoin("uatom", osmomath.NewInt(5_000_000))),
	)

	out, err = suite.callReadOnly(gamm.GetTotalVolumeForPoolMethod, map[string]interface{}{
		"pool_id": suite.PoolId,
	})
	suite.Require().NoError(err)
	coinsAfter, ok := out[0].([]struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	})
	suite.Require().True(ok)
	suite.Require().Len(coinsAfter, 1)
	suite.Equal("uatom", coinsAfter[0].Denom)
	suite.Equal(big.NewInt(5_000_000), coinsAfter[0].Amount)

	// Non-existent pool must error
	_, err = suite.callReadOnly(gamm.GetTotalVolumeForPoolMethod, map[string]interface{}{
		"pool_id": uint64(99999),
	})
	suite.Error(err)
}

func (suite *QueryParityTestSuite) TestGetTakerFeesTracker() {
	out, err := suite.callReadOnly(gamm.GetTakerFeesTrackerMethod, map[string]interface{}{})
	suite.Require().NoError(err)
	suite.Require().Len(out, 3)

	_, ok0 := out[0].([]struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	})
	_, ok1 := out[1].([]struct {
		Denom  string   `json:"denom"`
		Amount *big.Int `json:"amount"`
	})
	_, ok2 := out[2].(int64)
	suite.True(ok0)
	suite.True(ok1)
	suite.True(ok2)
}

func (suite *QueryParityTestSuite) TestGetTradingPairTakerFee() {
	// Default trading pair taker fee
	out, err := suite.callReadOnly(gamm.GetTradingPairTakerFeeMethod, map[string]interface{}{
		"denom_0": "uatom",
		"denom_1": "uosmo",
	})
	suite.Require().NoError(err)
	suite.Require().Len(out, 1)
	feeStr, ok := out[0].(string)
	suite.Require().True(ok)
	_, err = osmomath.NewDecFromStr(feeStr)
	suite.Require().NoError(err)

	// Override trading pair taker fee and verify
	override := osmomath.MustNewDecFromStr("0.005")
	suite.TestSuite.App.PoolManagerKeeper.SetDenomPairTakerFee(suite.TestSuite.Ctx, "uatom", "uosmo", override)
	out, err = suite.callReadOnly(gamm.GetTradingPairTakerFeeMethod, map[string]interface{}{
		"denom_0": "uatom",
		"denom_1": "uosmo",
	})
	suite.Require().NoError(err)
	feeStr, ok = out[0].(string)
	suite.Require().True(ok)
	gotDec, err := osmomath.NewDecFromStr(feeStr)
	suite.Require().NoError(err)
	suite.True(gotDec.Equal(override), "expected override taker fee 0.005, got %s", feeStr)

	// Empty denoms must error
	_, err = suite.callReadOnly(gamm.GetTradingPairTakerFeeMethod, map[string]interface{}{
		"denom_0": "",
		"denom_1": "uosmo",
	})
	suite.Error(err)
}
