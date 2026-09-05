package gamm_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	"github.com/ethereum/go-ethereum/core/vm"

	gamm "github.com/bitbadges/bitbadgeschain/x/gamm/precompile"
	"github.com/bitbadges/bitbadgeschain/x/gamm/precompile/test/helpers"
)

// Gas metering of the gamm precompile: store access is charged at the SDK's
// rates, the up-front requirement grows with the input, and the per-element
// prices and size caps declared in gas.go / security.go apply to the parsed
// message.
type GasMeteringTestSuite struct {
	suite.Suite
	TestSuite *helpers.TestSuite
}

func TestGasMeteringTestSuite(t *testing.T) {
	suite.Run(t, new(GasMeteringTestSuite))
}

func (suite *GasMeteringTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite(suite.T())
}

func (suite *GasMeteringTestSuite) swapInput(routes int, memo string) []byte {
	routeList := make([]map[string]interface{}, routes)
	for i := range routeList {
		routeList[i] = map[string]interface{}{"pool_id": uint64(i + 1), "token_out_denom": "uion"}
	}
	msg := map[string]interface{}{
		"routes":               routeList,
		"token_in":             map[string]interface{}{"denom": "uosmo", "amount": "1000"},
		"token_out_min_amount": "1",
	}
	methodName := "swapExactAmountIn"
	if memo != "" {
		methodName = "swapExactAmountInWithIBCTransfer"
		msg["ibc_transfer_info"] = map[string]interface{}{
			"source_channel":    "channel-0",
			"receiver":          "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu",
			"memo":              memo,
			"timeout_timestamp": uint64(1) << 62,
		}
	}
	bz, err := json.Marshal(msg)
	suite.Require().NoError(err)
	method := suite.TestSuite.Precompile.ABI.Methods[methodName]
	input, err := helpers.PackMethodWithJSON(&method, string(bz))
	suite.Require().NoError(err)
	return input
}

func (suite *GasMeteringTestSuite) contract(input []byte) *vm.Contract {
	c := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	c.Input = input
	return c
}

func (suite *GasMeteringTestSuite) execute(input []byte) error {
	_, err := suite.TestSuite.Precompile.Execute(suite.TestSuite.Ctx, suite.contract(input), false)
	return err
}

func (suite *GasMeteringTestSuite) TestStoreAccessIsMetered() {
	suite.Equal(storetypes.KVGasConfig(), suite.TestSuite.Precompile.KvGasConfig)
	suite.Equal(storetypes.TransientGasConfig(), suite.TestSuite.Precompile.TransientKVGasConfig)
}

func (suite *GasMeteringTestSuite) TestRequiredGasGrowsWithInputSize() {
	small := suite.swapInput(1, "")
	large := suite.swapInput(gamm.MaxRoutes, "")
	suite.Require().Greater(len(large), len(small))
	suite.Greater(suite.TestSuite.Precompile.RequiredGas(large), suite.TestSuite.Precompile.RequiredGas(small))
}

func (suite *GasMeteringTestSuite) TestSwapChargesPerRoute() {
	gasFor := func(routes int) uint64 {
		ctx := suite.TestSuite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		contract := suite.contract(suite.swapInput(routes, ""))
		// The pools do not exist, so the keeper rejects the swap; the
		// per-element charge is taken before dispatch either way.
		_, _ = suite.TestSuite.Precompile.Execute(ctx, contract, false)
		return ctx.GasMeter().GasConsumed()
	}
	one := gasFor(1)
	ten := gasFor(gamm.MaxRoutes)
	suite.GreaterOrEqual(ten-one, uint64((gamm.MaxRoutes-1)*gamm.GasPerRoute))
}

func (suite *GasMeteringTestSuite) TestSwapRejectsTooManyRoutes() {
	err := suite.execute(suite.swapInput(gamm.MaxRoutes+1, ""))
	suite.Require().Error(err)
	suite.Contains(err.Error(), "exceeds maximum allowed size")
}

func (suite *GasMeteringTestSuite) TestSwapWithIBCTransferRejectsLongMemo() {
	err := suite.execute(suite.swapInput(1, strings.Repeat("m", gamm.MaxMemoLength+1)))
	suite.Require().Error(err)
	suite.Contains(err.Error(), "exceeds maximum allowed length")
}
