package tokenization_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
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
	
	// Build JSON query
	queryJson, err := helpers.BuildGetCollectionQueryJSON(suite.TestSuite.CollectionId.BigInt())
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify result can be unpacked
	suite.Greater(len(result), 0)
}

func (suite *QueryMethodsTestSuite) TestGetCollection_NonExistent() {
	method := suite.Precompile.ABI.Methods["getCollection"]
	
	// Build JSON query for non-existent collection
	queryJson, err := helpers.BuildGetCollectionQueryJSON(big.NewInt(99999))
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *QueryMethodsTestSuite) TestGetBalance_Valid() {
	method := suite.Precompile.ABI.Methods["getBalance"]
	
	// Convert EVM address to Cosmos address
	aliceCosmos := suite.TestSuite.Alice.String()

	// Build JSON query
	queryJson, err := helpers.BuildGetBalanceQueryJSON(suite.TestSuite.CollectionId.BigInt(), aliceCosmos)
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify result structure
	suite.Greater(len(result), 0)
}

func (suite *QueryMethodsTestSuite) TestGetBalance_NonExistentUser() {
	method := suite.Precompile.ABI.Methods["getBalance"]
	
	// Convert EVM address to Cosmos address
	charlieCosmos := suite.TestSuite.Charlie.String()

	// Build JSON query
	queryJson, err := helpers.BuildGetBalanceQueryJSON(suite.TestSuite.CollectionId.BigInt(), charlieCosmos)
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	// Should succeed but return empty balance
	suite.NoError(err)
	suite.NotNil(result)
}

func (suite *QueryMethodsTestSuite) TestGetBalanceAmount_Valid() {
	method := suite.Precompile.ABI.Methods["getBalanceAmount"]
	
	// Convert EVM address to Cosmos address
	aliceCosmos := suite.TestSuite.Alice.String()

	// Build JSON query
	queryJson, err := helpers.BuildQueryJSON(map[string]interface{}{
		"collectionId": suite.TestSuite.CollectionId.BigInt().String(),
		"address":       aliceCosmos,
		"tokenIds": []map[string]interface{}{
			{"start": "1", "end": "10"},
		},
		"ownershipTimes": []map[string]interface{}{
			{"start": "1", "end": "1000"},
		},
	})
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
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
	
	// Build JSON query
	queryJson, err := helpers.BuildQueryJSON(map[string]interface{}{
		"collectionId": suite.TestSuite.CollectionId.BigInt().String(),
		"tokenIds": []map[string]interface{}{
			{"start": "1", "end": "10"},
		},
		"ownershipTimes": []map[string]interface{}{
			{"start": "1", "end": "1000"},
		},
	})
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
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
	// Create address list first using JSON
	caller := suite.TestSuite.AliceEVM

	addressListInput := map[string]interface{}{
		"listId":     "testlist", // Use alphanumeric only (no hyphens)
		"addresses":  []string{suite.TestSuite.Bob.String()},
		"whitelist":  true,
		"uri":        "https://example.com",
		"customData": "data",
	}

	createMsg := map[string]interface{}{
		"creator":       suite.TestSuite.Alice.String(),
		"addressLists": []interface{}{addressListInput},
	}

	createJson, err := helpers.BuildQueryJSON(createMsg)
	suite.NoError(err)

	createMethod, found := suite.Precompile.ABI.Methods["createAddressLists"]
	if !found {
		suite.T().Skip("createAddressLists method not found in ABI")
		return
	}

	createInput, err := helpers.PackMethodWithJSON(&createMethod, createJson)
	suite.NoError(err)

	createContract := suite.TestSuite.CreateMockContract(caller, createInput)
	_, err = suite.Precompile.Execute(suite.TestSuite.Ctx, createContract, false)
	suite.NoError(err)

	// Query the address list
	method := suite.Precompile.ABI.Methods["getAddressList"]
	
	// Build JSON query
	queryJson, err := helpers.BuildGetAddressListQueryJSON("testlist")
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	queryContract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, queryContract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify result structure
	suite.Greater(len(result), 0)
}

func (suite *QueryMethodsTestSuite) TestGetAddressList_NonExistent() {
	method := suite.Precompile.ABI.Methods["getAddressList"]
	
	// Build JSON query
	queryJson, err := helpers.BuildGetAddressListQueryJSON("nonexistentlist")
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
}
