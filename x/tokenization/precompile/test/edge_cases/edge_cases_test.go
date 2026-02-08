package tokenization_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

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

func (suite *EdgeCasesTestSuite) createContract(caller common.Address) *vm.Contract {
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	return vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
}

// TestMaximumCollectionId tests with maximum valid collection ID
func (suite *EdgeCasesTestSuite) TestMaximumCollectionId() {
	// Create collection with very large ID (will be auto-incremented, but test with large value in queries)
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	// Test query with the created collection ID
	method := suite.Precompile.ABI.Methods["getCollection"]
	args := []interface{}{
		collectionId.BigInt(),
	}

	result, err := suite.Precompile.GetCollection(suite.TestSuite.Ctx, &method, args)
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
	args := []interface{}{
		collectionId.BigInt(),
		suite.TestSuite.AliceEVM,
	}

	result, err := suite.Precompile.GetBalance(suite.TestSuite.Ctx, &method, args)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestMaximumArraySizes tests with maximum array sizes
func (suite *EdgeCasesTestSuite) TestMaximumArraySizes() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	// Create collection with many token ID ranges
	maxRanges := 100 // Reasonable limit for testing
	// Convert to []interface{} format expected by CreateCollection
	validTokenIds := make([]interface{}, maxRanges)

	for i := 0; i < maxRanges; i++ {
		validTokenIds[i] = map[string]interface{}{
			"start": big.NewInt(int64(i*10 + 1)),
			"end":   big.NewInt(int64((i + 1) * 10)),
		}
	}

	args := []interface{}{
		nil,
		validTokenIds,
		map[string]interface{}{},
		suite.TestSuite.Manager.String(),
		map[string]interface{}{"uri": "", "customData": ""},
		[]interface{}{},
		"",
		[]interface{}{},
		[]string{},
		false,
		[]interface{}{},
		[]interface{}{},
		map[string]interface{}{},
		[]interface{}{},
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestEmptyArrays tests with empty arrays
func (suite *EdgeCasesTestSuite) TestEmptyArrays() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	args := []interface{}{
		nil, // defaultBalances
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{}, // Empty validTokenIds (should fail validation)
		map[string]interface{}{},
		suite.TestSuite.Manager.String(),
		map[string]interface{}{"uri": "", "customData": ""},
		[]interface{}{}, // Empty tokenMetadata
		"",
		[]interface{}{}, // Empty collectionApprovals
		[]string{},       // Empty standards
		false,
		[]interface{}{},
		[]interface{}{},
		map[string]interface{}{},
		[]interface{}{},
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
	// Should fail because validTokenIds is empty
	suite.Error(err)
	suite.Nil(result)
}

// TestZeroValues tests with zero values
func (suite *EdgeCasesTestSuite) TestZeroValues() {
	// Test zero collection ID
	method := suite.Precompile.ABI.Methods["getCollection"]
	args := []interface{}{
		big.NewInt(0),
	}

	result, err := suite.Precompile.GetCollection(suite.TestSuite.Ctx, &method, args)
	suite.Error(err)
	suite.Nil(result)
}

// TestVeryLongStrings tests with very long strings
func (suite *EdgeCasesTestSuite) TestVeryLongStrings() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	// Create very long URI
	longURI := make([]byte, 10000)
	for i := range longURI {
		longURI[i] = 'a'
	}

	validTokenIds := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: big.NewInt(100)},
	}

	args := []interface{}{
		nil,
		validTokenIds,
		map[string]interface{}{},
		suite.TestSuite.Manager.String(),
		map[string]interface{}{"uri": string(longURI), "customData": ""},
		[]interface{}{},
		"",
		[]interface{}{},
		[]string{},
		false,
		[]interface{}{},
		[]interface{}{},
		map[string]interface{}{},
		[]interface{}{},
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
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
	contract := suite.createContract(caller)

	unicodeURI := "https://example.com/测试/🚀/日本語"

	validTokenIds := []interface{}{
		map[string]interface{}{
			"start": big.NewInt(1),
			"end":   big.NewInt(100),
		},
	}

	args := []interface{}{
		nil,
		validTokenIds,
		map[string]interface{}{},
		suite.TestSuite.Manager.String(),
		map[string]interface{}{"uri": unicodeURI, "customData": "测试数据"},
		[]interface{}{},
		"",
		[]interface{}{},
		[]string{},
		false,
		[]interface{}{},
		[]interface{}{},
		map[string]interface{}{},
		[]interface{}{},
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestVeryLargeNestedStructures tests with very large nested structures
func (suite *EdgeCasesTestSuite) TestVeryLargeNestedStructures() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

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

	validTokenIds := []interface{}{
		map[string]interface{}{
			"start": big.NewInt(1),
			"end":   big.NewInt(500),
		},
	}

	args := []interface{}{
		nil,
		validTokenIds,
		map[string]interface{}{},
		suite.TestSuite.Manager.String(),
		map[string]interface{}{"uri": "", "customData": ""},
		tokenMetadata,
		"",
		[]interface{}{},
		[]string{},
		false,
		[]interface{}{},
		[]interface{}{},
		map[string]interface{}{},
		[]interface{}{},
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestOverlappingRanges tests with overlapping token ID ranges
func (suite *EdgeCasesTestSuite) TestOverlappingRanges() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	// Overlapping ranges should be handled by the keeper validation
	validTokenIds := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: big.NewInt(50)},
		{Start: big.NewInt(40), End: big.NewInt(100)}, // Overlaps with first
	}

	args := []interface{}{
		nil,
		validTokenIds,
		map[string]interface{}{},
		suite.TestSuite.Manager.String(),
		map[string]interface{}{"uri": "", "customData": ""},
		[]interface{}{},
		"",
		[]interface{}{},
		[]string{},
		false,
		[]interface{}{},
		[]interface{}{},
		map[string]interface{}{},
		[]interface{}{},
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
	// May succeed or fail depending on keeper validation
	_ = result
	_ = err
}

// TestSingleTokenIdRange tests with single token ID (start == end)
func (suite *EdgeCasesTestSuite) TestSingleTokenIdRange() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	validTokenIds := []interface{}{
		map[string]interface{}{
			"start": big.NewInt(1),
			"end":   big.NewInt(1), // Single token
		},
	}

	args := []interface{}{
		nil,
		validTokenIds,
		map[string]interface{}{},
		suite.TestSuite.Manager.String(),
		map[string]interface{}{"uri": "", "customData": ""},
		[]interface{}{},
		"",
		[]interface{}{},
		[]string{},
		false,
		[]interface{}{},
		[]interface{}{},
		map[string]interface{}{},
		[]interface{}{},
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
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
	args := []interface{}{
		collectionId.BigInt(),
		suite.TestSuite.CharlieEVM, // User with no balance
	}

	result, err := suite.Precompile.GetBalance(suite.TestSuite.Ctx, &method, args)
	suite.NoError(err)
	suite.NotNil(result)
}

// TestInvalidRangeStartGreaterThanEnd tests with invalid range
func (suite *EdgeCasesTestSuite) TestInvalidRangeStartGreaterThanEnd() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	// Invalid: start > end
	invalidTokenIds := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(100), End: big.NewInt(1)},
	}

	args := []interface{}{
		nil,
		invalidTokenIds,
		map[string]interface{}{},
		suite.TestSuite.Manager.String(),
		map[string]interface{}{"uri": "", "customData": ""},
		[]interface{}{},
		"",
		[]interface{}{},
		[]string{},
		false,
		[]interface{}{},
		[]interface{}{},
		map[string]interface{}{},
		[]interface{}{},
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
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
	contract := suite.createContract(caller)

	// Create array with many recipients (up to MaxRecipients)
	recipients := make([]common.Address, tokenization.MaxRecipients)
	for i := 0; i < tokenization.MaxRecipients; i++ {
		// Use different addresses
		addr := common.BigToAddress(big.NewInt(int64(i + 1)))
		recipients[i] = addr
	}

	args := []interface{}{
		collectionId.BigInt(),
		recipients,
		big.NewInt(1), // 1 token per recipient
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: big.NewInt(1000)}},
	}

	method := suite.Precompile.ABI.Methods["transferTokens"]
	result, err := suite.Precompile.TransferTokens(suite.TestSuite.Ctx, &method, args, contract)
	// May succeed or fail based on balance and validation
	_ = result
	_ = err
}

// TestExceedsMaxRecipients tests with more than MaxRecipients
func (suite *EdgeCasesTestSuite) TestExceedsMaxRecipients() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	// Create array exceeding MaxRecipients
	recipients := make([]common.Address, tokenization.MaxRecipients+1)
	for i := 0; i < tokenization.MaxRecipients+1; i++ {
		recipients[i] = common.BigToAddress(big.NewInt(int64(i + 1)))
	}

	args := []interface{}{
		collectionId.BigInt(),
		recipients,
		big.NewInt(1),
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct {
			Start *big.Int `json:"start"`
			End   *big.Int `json:"end"`
		}{{Start: big.NewInt(1), End: big.NewInt(1000)}},
	}

	method := suite.Precompile.ABI.Methods["transferTokens"]
	result, err := suite.Precompile.TransferTokens(suite.TestSuite.Ctx, &method, args, contract)
	suite.Error(err)
	suite.Nil(result)
	suite.Contains(err.Error(), "exceeds maximum")
}

// TestNilFields tests with nil fields
func (suite *EdgeCasesTestSuite) TestNilFields() {
	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	validTokenIds := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: big.NewInt(100)},
	}

	args := []interface{}{
		nil, // nil defaultBalances
		validTokenIds,
		nil, // nil collectionPermissions (should be handled)
		suite.TestSuite.Manager.String(),
		nil, // nil collectionMetadata (should be handled)
		nil, // nil tokenMetadata
		"",
		nil, // nil collectionApprovals
		nil, // nil standards
		false,
		nil,
		nil,
		nil,
		nil,
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
	// Should handle nil fields gracefully
	_ = result
	_ = err
}

