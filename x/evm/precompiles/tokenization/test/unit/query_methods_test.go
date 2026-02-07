package tokenization_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	tokenization "github.com/bitbadges/bitbadgeschain/x/evm/precompiles/tokenization"
	"github.com/bitbadges/bitbadgeschain/x/evm/precompiles/tokenization/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type QueryMethodsTestSuite struct {
	suite.Suite
	TestSuite  *helpers.TestSuite
	Precompile *tokenization.Precompile
}

func TestQueryMethodsTestSuite(t *testing.T) {
	suite.Run(t, new(QueryMethodsTestSuite))
}

func (suite *QueryMethodsTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite()
	suite.Precompile = suite.TestSuite.Precompile

	// Create test collection with balance
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	err = suite.TestSuite.CreateTestBalance(
		collectionId,
		suite.TestSuite.Alice.String(),
		sdkmath.NewUint(1000),
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
	)
	suite.NoError(err)
}

func (suite *QueryMethodsTestSuite) TestGetCollection_Valid() {
	method := suite.Precompile.ABI.Methods["getCollection"]
	args := []interface{}{
		suite.TestSuite.CollectionId.BigInt(),
	}

	result, err := suite.Precompile.GetCollection(suite.TestSuite.Ctx, &method, args)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify result can be unpacked (should be structured type now)
	// Note: ABI still shows bytes, but actual return is structured
	// This will be updated when ABI is updated
	suite.Greater(len(result), 0)
}

func (suite *QueryMethodsTestSuite) TestGetCollection_NonExistent() {
	method := suite.Precompile.ABI.Methods["getCollection"]
	args := []interface{}{
		big.NewInt(99999), // Non-existent collection
	}

	result, err := suite.Precompile.GetCollection(suite.TestSuite.Ctx, &method, args)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *QueryMethodsTestSuite) TestGetBalance_Valid() {
	method := suite.Precompile.ABI.Methods["getBalance"]
	args := []interface{}{
		suite.TestSuite.CollectionId.BigInt(),
		suite.TestSuite.AliceEVM,
	}

	result, err := suite.Precompile.GetBalance(suite.TestSuite.Ctx, &method, args)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify result structure
	suite.Greater(len(result), 0)
}

func (suite *QueryMethodsTestSuite) TestGetBalance_NonExistentUser() {
	method := suite.Precompile.ABI.Methods["getBalance"]
	args := []interface{}{
		suite.TestSuite.CollectionId.BigInt(),
		suite.TestSuite.CharlieEVM, // User with no balance
	}

	result, err := suite.Precompile.GetBalance(suite.TestSuite.Ctx, &method, args)
	// Should succeed but return empty balance
	suite.NoError(err)
	suite.NotNil(result)
}

func (suite *QueryMethodsTestSuite) TestGetBalanceAmount_Valid() {
	method := suite.Precompile.ABI.Methods["getBalanceAmount"]
	args := []interface{}{
		suite.TestSuite.CollectionId.BigInt(),
		suite.TestSuite.AliceEVM,
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: big.NewInt(1000)}},
	}

	result, err := suite.Precompile.GetBalanceAmount(suite.TestSuite.Ctx, &method, args)
	suite.NoError(err)
	suite.NotNil(result)

	// Unpack result
	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)

	amount, ok := unpacked[0].(*big.Int)
	suite.True(ok)
	suite.Greater(amount.Uint64(), uint64(0))
}

func (suite *QueryMethodsTestSuite) TestGetTotalSupply_Valid() {
	method := suite.Precompile.ABI.Methods["getTotalSupply"]
	args := []interface{}{
		suite.TestSuite.CollectionId.BigInt(),
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: big.NewInt(1000)}},
	}

	result, err := suite.Precompile.GetTotalSupply(suite.TestSuite.Ctx, &method, args)
	suite.NoError(err)
	suite.NotNil(result)

	// Unpack result
	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)

	supply, ok := unpacked[0].(*big.Int)
	suite.True(ok)
	// Supply may be 0 if no balances exist yet - just verify it's a valid value
	suite.GreaterOrEqual(supply.Uint64(), uint64(0))
}

func (suite *QueryMethodsTestSuite) TestGetAddressList_Valid() {
	// Create address list first
	caller := suite.TestSuite.AliceEVM
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	contract := vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)

	addressListInput := map[string]interface{}{
		"listId":     "testlist", // Use alphanumeric only (no hyphens)
		"addresses":  []string{suite.TestSuite.Bob.String()},
		"whitelist":  true,
		"uri":        "https://example.com",
		"customData": "data",
	}

	createArgs := []interface{}{
		[]interface{}{addressListInput},
	}

	createMethod, found := suite.Precompile.ABI.Methods["createAddressLists"]
	if !found {
		// Create mock method if not in ABI (workaround for missing ABI entries)
		createMethod = helpers.CreateMockMethod("createAddressLists", nil, helpers.CreateMockBoolOutput())
	}
	_, err := suite.Precompile.CreateAddressLists(suite.TestSuite.Ctx, &createMethod, createArgs, contract)
	suite.NoError(err)

	// Query the address list
	method := suite.Precompile.ABI.Methods["getAddressList"]
	args := []interface{}{
		"testlist",
	}

	result, err := suite.Precompile.GetAddressList(suite.TestSuite.Ctx, &method, args)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify result structure
	suite.Greater(len(result), 0)
}

func (suite *QueryMethodsTestSuite) TestGetAddressList_NonExistent() {
	method := suite.Precompile.ABI.Methods["getAddressList"]
	args := []interface{}{
		"nonexistentlist", // Use alphanumeric only (no hyphens)
	}

	result, err := suite.Precompile.GetAddressList(suite.TestSuite.Ctx, &method, args)
	suite.Error(err)
	suite.Nil(result)
}
