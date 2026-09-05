package tokenization_test

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// Gas metering of the tokenization precompile: store access is charged at the
// SDK's rates, the up-front requirement grows with the input, and the
// per-element prices and size caps declared in gas.go / security.go apply to
// the parsed message.
type GasMeteringTestSuite struct {
	suite.Suite
	TestSuite    *helpers.TestSuite
	CollectionId sdkmath.Uint
}

func TestGasMeteringTestSuite(t *testing.T) {
	suite.Run(t, new(GasMeteringTestSuite))
}

func (suite *GasMeteringTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite()
	suite.CollectionId = suite.createOpenCollection()
}

// createOpenCollection is the helper collection plus default balances that
// auto-approve incoming transfers, so fresh recipients need no setup.
func (suite *GasMeteringTestSuite) createOpenCollection() sdkmath.Uint {
	ts := suite.TestSuite
	collectionId, err := ts.CreateTestCollection(ts.Alice.String())
	suite.Require().NoError(err)

	collection, found := ts.Keeper.GetCollectionFromStore(ts.Ctx, collectionId)
	suite.Require().True(found)
	collection.DefaultBalances = &tokenizationtypes.UserBalanceStore{
		AutoApproveAllIncomingTransfers:           true,
		AutoApproveSelfInitiatedIncomingTransfers: true,
		AutoApproveSelfInitiatedOutgoingTransfers: true,
	}
	suite.Require().NoError(ts.Keeper.SetCollectionInStore(ts.Ctx, collection, false))

	suite.Require().NoError(ts.CreateTestBalance(
		collectionId,
		ts.Alice.String(),
		sdkmath.NewUint(100_000),
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
	))
	return collectionId
}

func freshRecipients(n int) []string {
	out := make([]string, n)
	for i := range out {
		addr := common.BigToAddress(new(big.Int).Add(big.NewInt(0x10000), big.NewInt(int64(i))))
		out[i] = sdk.AccAddress(addr.Bytes()).String()
	}
	return out
}

func (suite *GasMeteringTestSuite) transferInput(recipients int) []byte {
	jsonMsg, err := helpers.BuildTransferTokensJSON(
		suite.CollectionId.BigInt(),
		suite.TestSuite.Alice.String(),
		freshRecipients(recipients),
		big.NewInt(1),
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(1000)}},
	)
	suite.Require().NoError(err)
	method := suite.TestSuite.Precompile.ABI.Methods["transferTokens"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.Require().NoError(err)
	return input
}

func (suite *GasMeteringTestSuite) TestStoreAccessIsMetered() {
	suite.Equal(storetypes.KVGasConfig(), suite.TestSuite.Precompile.KvGasConfig)
	suite.Equal(storetypes.TransientGasConfig(), suite.TestSuite.Precompile.TransientKVGasConfig)
}

func (suite *GasMeteringTestSuite) TestRequiredGasGrowsWithInputSize() {
	small := suite.transferInput(1)
	large := suite.transferInput(50)
	suite.Require().Greater(len(large), len(small))

	smallGas := suite.TestSuite.Precompile.RequiredGas(small)
	largeGas := suite.TestSuite.Precompile.RequiredGas(large)
	suite.Greater(largeGas, smallGas, "a larger transferTokens input must require more gas up front")

	wantExtra := uint64(len(large)/32-len(small)/32) * tokenization.GasPerInputChunk
	suite.Equal(smallGas+wantExtra, largeGas, "the difference is the per-chunk input price")

	// Queries are priced by input size too.
	queryMethod := suite.TestSuite.Precompile.ABI.Methods["getCollection"]
	smallQuery, err := helpers.PackMethodWithJSON(&queryMethod, `{"collectionId":"1"}`)
	suite.Require().NoError(err)
	largeQuery, err := helpers.PackMethodWithJSON(&queryMethod, fmt.Sprintf(`{"collectionId":"1%01000d"}`, 0))
	suite.Require().NoError(err)
	suite.Greater(suite.TestSuite.Precompile.RequiredGas(largeQuery), suite.TestSuite.Precompile.RequiredGas(smallQuery))
}

func (suite *GasMeteringTestSuite) TestTransferTokensChargesPerRecipient() {
	gasFor := func(recipients int) uint64 {
		ctx := suite.TestSuite.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
		contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, suite.transferInput(recipients))
		_, err := suite.TestSuite.Precompile.Execute(ctx, contract, false)
		suite.Require().NoError(err)
		return ctx.GasMeter().GasConsumed()
	}

	one := gasFor(1)
	five := gasFor(5)
	suite.GreaterOrEqual(five-one, uint64(4*tokenization.GasPerRecipient),
		"each additional recipient must cost at least GasPerRecipient on the SDK meter")
}

func (suite *GasMeteringTestSuite) TestTransferTokensRejectsOversizeRecipientList() {
	_, err := suite.TestSuite.CallPrecompile(suite.TestSuite.AliceEVM, suite.transferInput(tokenization.MaxRecipients+1))
	suite.Require().Error(err)
	suite.Contains(err.Error(), "exceeds maximum allowed size")

	_, err = suite.TestSuite.CallPrecompile(suite.TestSuite.AliceEVM, suite.transferInput(tokenization.MaxRecipients))
	suite.Require().NoError(err, "exactly MaxRecipients is allowed")
}

func (suite *GasMeteringTestSuite) TestCreateAddressListsRejectsOversizeList() {
	method := suite.TestSuite.Precompile.ABI.Methods["createAddressLists"]
	build := func(n int) []byte {
		jsonMsg, err := helpers.BuildCreateAddressListsJSON(suite.TestSuite.Alice.String(), []map[string]interface{}{{
			"listId":     fmt.Sprintf("list%d", n),
			"addresses":  freshRecipients(n),
			"whitelist":  true,
			"uri":        "",
			"customData": "",
		}})
		suite.Require().NoError(err)
		input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
		suite.Require().NoError(err)
		return input
	}

	_, err := suite.TestSuite.CallPrecompile(suite.TestSuite.AliceEVM, build(tokenization.MaxAddressListEntries+1))
	suite.Require().Error(err)
	suite.Contains(err.Error(), "exceeds maximum allowed size")

	_, err = suite.TestSuite.CallPrecompile(suite.TestSuite.AliceEVM, build(tokenization.MaxAddressListEntries))
	suite.Require().NoError(err, "exactly MaxAddressListEntries is allowed")
}

func (suite *GasMeteringTestSuite) TestSetCollectionMetadataRejectsOversizeURI() {
	method := suite.TestSuite.Precompile.ABI.Methods["setCollectionMetadata"]
	build := func(uriLen int) []byte {
		jsonMsg, err := helpers.BuildSetCollectionMetadataJSON(suite.TestSuite.Alice.String(), suite.CollectionId.BigInt(), map[string]interface{}{
			"uri":        "https://x.io/" + strings.Repeat("a", uriLen),
			"customData": "",
		})
		suite.Require().NoError(err)
		input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
		suite.Require().NoError(err)
		return input
	}

	_, err := suite.TestSuite.CallPrecompile(suite.TestSuite.AliceEVM, build(tokenization.MaxMetadataLength+1))
	suite.Require().Error(err)
	suite.Contains(err.Error(), "exceeds maximum allowed length")
}
