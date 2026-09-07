package sendmanager_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"

	sendmanager "github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile"
	"github.com/bitbadges/bitbadgeschain/x/sendmanager/precompile/test/helpers"
)

// Gas metering of the sendmanager precompile: store access is charged at the
// SDK's rates, the up-front requirement grows with the input, and the coin
// list is priced per element and capped.
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

func (suite *GasMeteringTestSuite) sendInput(coins int) []byte {
	amount := make([]map[string]interface{}, coins)
	for i := range amount {
		amount[i] = map[string]interface{}{"denom": fmt.Sprintf("coin%03d", i), "amount": "1"}
	}
	bz, err := json.Marshal(map[string]interface{}{
		"to_address": suite.TestSuite.Bob.String(),
		"amount":     amount,
	})
	suite.Require().NoError(err)
	method := suite.TestSuite.Precompile.ABI.Methods["send"]
	input, err := helpers.PackMethodCall(&method, string(bz))
	suite.Require().NoError(err)
	return input
}

func (suite *GasMeteringTestSuite) TestStoreAccessIsMetered() {
	suite.Equal(storetypes.KVGasConfig(), suite.TestSuite.Precompile.KvGasConfig)
	suite.Equal(storetypes.TransientGasConfig(), suite.TestSuite.Precompile.TransientKVGasConfig)
}

func (suite *GasMeteringTestSuite) TestRequiredGasGrowsWithInputSize() {
	small := suite.sendInput(1)
	large := suite.sendInput(sendmanager.MaxCoins)
	suite.Require().Greater(len(large), len(small))
	suite.Greater(suite.TestSuite.Precompile.RequiredGas(large), suite.TestSuite.Precompile.RequiredGas(small))
}

func (suite *GasMeteringTestSuite) TestSendChargesPerCoin() {
	gasFor := func(coins int) uint64 {
		ctx := suite.TestSuite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, suite.sendInput(coins))
		// Alice holds none of these denoms, so the keeper rejects the send;
		// the per-element charge is taken before dispatch either way.
		_, _ = suite.TestSuite.Precompile.Execute(ctx, contract, false)
		return ctx.GasMeter().GasConsumed()
	}
	one := gasFor(1)
	many := gasFor(sendmanager.MaxCoins)
	suite.GreaterOrEqual(many-one, uint64((sendmanager.MaxCoins-1)*sendmanager.GasPerCoin))
}

func (suite *GasMeteringTestSuite) TestSendRejectsTooManyCoins() {
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, suite.sendInput(sendmanager.MaxCoins+1))
	_, err := suite.TestSuite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Require().Error(err)
	suite.Contains(err.Error(), "exceeds maximum allowed size")
}
