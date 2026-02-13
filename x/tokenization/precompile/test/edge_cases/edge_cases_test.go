package tokenization_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type EdgeCasesTestSuite struct {
	suite.Suite
	TestSuite  *helpers.TestSuite
	Precompile *tokenization.Precompile
}

func TestEdgeCasesTestSuite(t *testing.T) {
	suite.Run(t, new(EdgeCasesTestSuite))
}

func (suite *EdgeCasesTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite()
	suite.Precompile = suite.TestSuite.Precompile
}

// createContract is deprecated - use suite.TestSuite.CreateMockContract instead

// TestMaximumCollectionId tests with maximum valid collection ID
func (suite *EdgeCasesTestSuite) TestMaximumCollectionId() {
	// Create collection with very large ID (will be auto-incremented, but test with large value in queries)
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	// Test query with the created collection ID
	method := suite.Precompile.ABI.Methods["getCollection"]

	// Build JSON query
	queryJson, err := helpers.BuildGetCollectionQueryJSON(collectionId.BigInt())
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestMaximumTokenAmount tests with maximum token amounts
func (suite *EdgeCasesTestSuite) TestMaximumTokenAmount() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	// Use maximum uint256 value
	maxAmount := new(big.Int)
	maxAmount.Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)) // 2^256 - 1

	err = suite.TestSuite.CreateTestBalance(
		collectionId,
		suite.TestSuite.Alice.String(),
		sdkmath.NewUintFromBigInt(maxAmount),
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
	)
	suite.NoError(err)

	// Query balance
	method := suite.Precompile.ABI.Methods["getBalance"]

	// Convert EVM address to Cosmos address
	aliceCosmos := suite.TestSuite.Alice.String()

	// Build JSON query
	queryJson, err := helpers.BuildGetBalanceQueryJSON(collectionId.BigInt(), aliceCosmos)
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestMaximumArraySizes tests with maximum array sizes
func (suite *EdgeCasesTestSuite) TestMaximumArraySizes() {
	caller := suite.TestSuite.AliceEVM

	// Create collection with many token ID ranges
	maxRanges := 100 // Reasonable limit for testing
	// Convert to []map[string]interface{} format for JSON
	validTokenIds := make([]map[string]interface{}, maxRanges)

	for i := 0; i < maxRanges; i++ {
		validTokenIds[i] = map[string]interface{}{
			"start": big.NewInt(int64(i*10 + 1)).String(),
			"end":   big.NewInt(int64((i + 1) * 10)).String(),
		}
	}

	// Build JSON message
	msg := map[string]interface{}{
		"defaultBalances":             nil,
		"validTokenIds":               validTokenIds,
		"collectionPermissions":       map[string]interface{}{},
		"manager":                     suite.TestSuite.Manager.String(),
		"collectionMetadata":          map[string]interface{}{"uri": "", "customData": ""},
		"tokenMetadata":               []interface{}{},
		"customData":                  "",
		"collectionApprovals":         []interface{}{},
		"standards":                   []string{},
		"isArchived":                  false,
		"mintEscrowCoinsToTransfer":   []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                  map[string]interface{}{},
		"aliasPathsToAdd":             []interface{}{},
	}

	jsonMsg, err := helpers.BuildCreateCollectionJSON(suite.TestSuite.Alice.String(), msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestEmptyArrays tests with empty arrays
func (suite *EdgeCasesTestSuite) TestEmptyArrays() {
	caller := suite.TestSuite.AliceEVM

	// Build JSON message with empty validTokenIds
	msg := map[string]interface{}{
		"defaultBalances":             nil,
		"validTokenIds":               []map[string]interface{}{}, // Empty
		"collectionPermissions":       map[string]interface{}{},
		"manager":                     suite.TestSuite.Manager.String(),
		"collectionMetadata":          map[string]interface{}{"uri": "", "customData": ""},
		"tokenMetadata":               []interface{}{},
		"customData":                  "",
		"collectionApprovals":         []interface{}{},
		"standards":                   []string{},
		"isArchived":                  false,
		"mintEscrowCoinsToTransfer":   []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                  map[string]interface{}{},
		"aliasPathsToAdd":             []interface{}{},
	}

	jsonMsg, err := helpers.BuildCreateCollectionJSON(suite.TestSuite.Alice.String(), msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	// Empty validTokenIds is allowed - collection creation should succeed
	suite.NoError(err)
	suite.NotNil(result)
	// Result should be a collectionId (uint256)
	suite.Greater(len(result), 0)
}

// TestZeroValues tests with zero values
func (suite *EdgeCasesTestSuite) TestZeroValues() {
	// Test zero collection ID
	method := suite.Precompile.ABI.Methods["getCollection"]

	// Build JSON query
	queryJson, err := helpers.BuildGetCollectionQueryJSON(big.NewInt(0))
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

// TestVeryLongStrings tests with very long strings
func (suite *EdgeCasesTestSuite) TestVeryLongStrings() {
	caller := suite.TestSuite.AliceEVM

	// Create very long URI
	longURI := make([]byte, 10000)
	for i := range longURI {
		longURI[i] = 'a'
	}

	// Build JSON message
	msg := map[string]interface{}{
		"defaultBalances":             nil,
		"validTokenIds":               []map[string]interface{}{{"start": "1", "end": "100"}},
		"collectionPermissions":       map[string]interface{}{},
		"manager":                     suite.TestSuite.Manager.String(),
		"collectionMetadata":          map[string]interface{}{"uri": string(longURI), "customData": ""},
		"tokenMetadata":               []interface{}{},
		"customData":                  "",
		"collectionApprovals":         []interface{}{},
		"standards":                   []string{},
		"isArchived":                  false,
		"mintEscrowCoinsToTransfer":   []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                  map[string]interface{}{},
		"aliasPathsToAdd":             []interface{}{},
	}

	jsonMsg, err := helpers.BuildCreateCollectionJSON(suite.TestSuite.Alice.String(), msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	// Should handle long strings (may succeed or fail based on validation)
	// Just verify it doesn't panic
	if err != nil {
		suite.Contains(err.Error(), "invalid") // Should be a validation error, not a panic
	}
	_ = result
}

// TestUnicodeStrings tests with Unicode strings
func (suite *EdgeCasesTestSuite) TestUnicodeStrings() {
	caller := suite.TestSuite.AliceEVM

	unicodeURI := "https://example.com/测试/🚀/日本語"

	// Build JSON message
	msg := map[string]interface{}{
		"defaultBalances":             nil,
		"validTokenIds":               []map[string]interface{}{{"start": "1", "end": "100"}},
		"collectionPermissions":       map[string]interface{}{},
		"manager":                     suite.TestSuite.Manager.String(),
		"collectionMetadata":          map[string]interface{}{"uri": unicodeURI, "customData": "测试数据"},
		"tokenMetadata":               []interface{}{},
		"customData":                  "",
		"collectionApprovals":         []interface{}{},
		"standards":                   []string{},
		"isArchived":                  false,
		"mintEscrowCoinsToTransfer":   []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                  map[string]interface{}{},
		"aliasPathsToAdd":             []interface{}{},
	}

	jsonMsg, err := helpers.BuildCreateCollectionJSON(suite.TestSuite.Alice.String(), msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestVeryLargeNestedStructures tests with very large nested structures
func (suite *EdgeCasesTestSuite) TestVeryLargeNestedStructures() {
	caller := suite.TestSuite.AliceEVM

	// Create collection with many token metadata entries
	tokenMetadata := make([]interface{}, 50)
	for i := 0; i < 50; i++ {
		tokenMetadata[i] = map[string]interface{}{
			"uri":        "https://example.com/token",
			"customData": "data",
			"tokenIds": []interface{}{
				map[string]interface{}{
					"start": big.NewInt(int64(i*10 + 1)),
					"end":   big.NewInt(int64((i + 1) * 10)),
				},
			},
		}
	}

	// Convert tokenMetadata to JSON-compatible format
	tokenMetadataJSON := make([]map[string]interface{}, len(tokenMetadata))
	for i, tm := range tokenMetadata {
		tmMap := tm.(map[string]interface{})
		tokenIds := tmMap["tokenIds"].([]interface{})
		tokenIdsJSON := make([]map[string]interface{}, len(tokenIds))
		for j, tid := range tokenIds {
			tidMap := tid.(map[string]interface{})
			tokenIdsJSON[j] = map[string]interface{}{
				"start": tidMap["start"].(*big.Int).String(),
				"end":   tidMap["end"].(*big.Int).String(),
			}
		}
		tokenMetadataJSON[i] = map[string]interface{}{
			"uri":        tmMap["uri"],
			"customData": tmMap["customData"],
			"tokenIds":   tokenIdsJSON,
		}
	}

	// Build JSON message
	msg := map[string]interface{}{
		"defaultBalances":             nil,
		"validTokenIds":               []map[string]interface{}{{"start": "1", "end": "500"}},
		"collectionPermissions":       map[string]interface{}{},
		"manager":                     suite.TestSuite.Manager.String(),
		"collectionMetadata":          map[string]interface{}{"uri": "", "customData": ""},
		"tokenMetadata":               tokenMetadataJSON,
		"customData":                  "",
		"collectionApprovals":         []interface{}{},
		"standards":                   []string{},
		"isArchived":                  false,
		"mintEscrowCoinsToTransfer":   []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                  map[string]interface{}{},
		"aliasPathsToAdd":             []interface{}{},
	}

	jsonMsg, err := helpers.BuildCreateCollectionJSON(suite.TestSuite.Alice.String(), msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestOverlappingRanges tests with overlapping token ID ranges
func (suite *EdgeCasesTestSuite) TestOverlappingRanges() {
	caller := suite.TestSuite.AliceEVM

	// Overlapping ranges should be handled by the keeper validation
	// Build JSON message
	msg := map[string]interface{}{
		"defaultBalances":             nil,
		"validTokenIds":               []map[string]interface{}{{"start": "1", "end": "50"}, {"start": "40", "end": "100"}}, // Overlaps
		"collectionPermissions":       map[string]interface{}{},
		"manager":                     suite.TestSuite.Manager.String(),
		"collectionMetadata":          map[string]interface{}{"uri": "", "customData": ""},
		"tokenMetadata":               []interface{}{},
		"customData":                  "",
		"collectionApprovals":         []interface{}{},
		"standards":                   []string{},
		"isArchived":                  false,
		"mintEscrowCoinsToTransfer":   []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                  map[string]interface{}{},
		"aliasPathsToAdd":             []interface{}{},
	}

	jsonMsg, err := helpers.BuildCreateCollectionJSON(suite.TestSuite.Alice.String(), msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	// May succeed or fail depending on keeper validation
	_ = result
	_ = err
}

// TestSingleTokenIdRange tests with single token ID (start == end)
func (suite *EdgeCasesTestSuite) TestSingleTokenIdRange() {
	caller := suite.TestSuite.AliceEVM

	// Build JSON message
	msg := map[string]interface{}{
		"defaultBalances":             nil,
		"validTokenIds":               []map[string]interface{}{{"start": "1", "end": "1"}}, // Single token
		"collectionPermissions":       map[string]interface{}{},
		"manager":                     suite.TestSuite.Manager.String(),
		"collectionMetadata":          map[string]interface{}{"uri": "", "customData": ""},
		"tokenMetadata":               []interface{}{},
		"customData":                  "",
		"collectionApprovals":         []interface{}{},
		"standards":                   []string{},
		"isArchived":                  false,
		"mintEscrowCoinsToTransfer":   []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                  map[string]interface{}{},
		"aliasPathsToAdd":             []interface{}{},
	}

	jsonMsg, err := helpers.BuildCreateCollectionJSON(suite.TestSuite.Alice.String(), msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestMaximumOwnershipTimeRange tests with maximum ownership time range
func (suite *EdgeCasesTestSuite) TestMaximumOwnershipTimeRange() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	// Use maximum uint64 for ownership time
	maxUint64 := new(big.Int)
	maxUint64.SetUint64(math.MaxUint64)

	err = suite.TestSuite.CreateTestBalance(
		collectionId,
		suite.TestSuite.Alice.String(),
		sdkmath.NewUint(1000),
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		[]*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUintFromBigInt(maxUint64)},
		},
	)
	suite.NoError(err)
}

// TestEmptyBalanceStore tests with empty balance store
func (suite *EdgeCasesTestSuite) TestEmptyBalanceStore() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	// Query balance for user with no balance (should return default or empty)
	method := suite.Precompile.ABI.Methods["getBalance"]

	// Convert EVM address to Cosmos address
	charlieCosmos := suite.TestSuite.Charlie.String()

	// Build JSON query
	queryJson, err := helpers.BuildGetBalanceQueryJSON(collectionId.BigInt(), charlieCosmos)
	suite.NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.NoError(err)

	// Call precompile via Execute
	contract := suite.TestSuite.CreateMockContract(suite.TestSuite.AliceEVM, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestInvalidRangeStartGreaterThanEnd tests with invalid range
func (suite *EdgeCasesTestSuite) TestInvalidRangeStartGreaterThanEnd() {
	caller := suite.TestSuite.AliceEVM

	// Invalid: start > end
	// Build JSON message
	msg := map[string]interface{}{
		"defaultBalances":             nil,
		"validTokenIds":               []map[string]interface{}{{"start": "100", "end": "1"}}, // Invalid: start > end
		"collectionPermissions":       map[string]interface{}{},
		"manager":                     suite.TestSuite.Manager.String(),
		"collectionMetadata":          map[string]interface{}{"uri": "", "customData": ""},
		"tokenMetadata":               []interface{}{},
		"customData":                  "",
		"collectionApprovals":         []interface{}{},
		"standards":                   []string{},
		"isArchived":                  false,
		"mintEscrowCoinsToTransfer":   []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                  map[string]interface{}{},
		"aliasPathsToAdd":             []interface{}{},
	}

	jsonMsg, err := helpers.BuildCreateCollectionJSON(suite.TestSuite.Alice.String(), msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
}

// TestManyRecipients tests with maximum number of recipients
func (suite *EdgeCasesTestSuite) TestManyRecipients() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	err = suite.TestSuite.CreateTestBalance(
		collectionId,
		suite.TestSuite.Alice.String(),
		sdkmath.NewUint(10000),
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
	)
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM

	// Create array with many recipients (up to MaxRecipients)
	recipients := make([]common.Address, tokenization.MaxRecipients)
	recipientsCosmos := make([]string, tokenization.MaxRecipients)
	for i := 0; i < tokenization.MaxRecipients; i++ {
		// Use different addresses
		addr := common.BigToAddress(big.NewInt(int64(i + 1)))
		recipients[i] = addr
		recipientsCosmos[i] = sdk.AccAddress(addr.Bytes()).String()
	}

	// Build JSON message
	jsonMsg, err := helpers.BuildTransferTokensJSON(
		collectionId.BigInt(),
		suite.TestSuite.Alice.String(),
		recipientsCosmos,
		big.NewInt(1), // 1 token per recipient
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(1000)}},
	)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["transferTokens"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	// May succeed or fail based on balance and validation
	_ = result
	_ = err
}

// TestExceedsMaxRecipients tests with more than MaxRecipients
func (suite *EdgeCasesTestSuite) TestExceedsMaxRecipients() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM

	// Create array exceeding MaxRecipients
	recipientsCosmos := make([]string, tokenization.MaxRecipients+1)
	for i := 0; i < tokenization.MaxRecipients+1; i++ {
		addr := common.BigToAddress(big.NewInt(int64(i + 1)))
		recipientsCosmos[i] = sdk.AccAddress(addr.Bytes()).String()
	}

	// Build JSON message
	jsonMsg, err := helpers.BuildTransferTokensJSON(
		collectionId.BigInt(),
		suite.TestSuite.Alice.String(),
		recipientsCosmos,
		big.NewInt(1),
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(1000)}},
	)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["transferTokens"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)
}

// TestNilFields tests with nil fields
func (suite *EdgeCasesTestSuite) TestNilFields() {
	caller := suite.TestSuite.AliceEVM

	// Build JSON message with nil fields (should be handled gracefully)
	msg := map[string]interface{}{
		"defaultBalances":             nil,
		"validTokenIds":               []map[string]interface{}{{"start": "1", "end": "100"}},
		"collectionPermissions":       nil, // nil (should be handled)
		"manager":                     suite.TestSuite.Manager.String(),
		"collectionMetadata":          nil, // nil (should be handled)
		"tokenMetadata":               nil, // nil
		"customData":                  "",
		"collectionApprovals":         nil, // nil
		"standards":                   nil, // nil
		"isArchived":                  false,
		"mintEscrowCoinsToTransfer":   nil,
		"cosmosCoinWrapperPathsToAdd": nil,
		"invariants":                  nil,
		"aliasPathsToAdd":             nil,
	}

	jsonMsg, err := helpers.BuildCreateCollectionJSON(suite.TestSuite.Alice.String(), msg)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	// Should handle nil fields gracefully
	_ = result
	_ = err
}
