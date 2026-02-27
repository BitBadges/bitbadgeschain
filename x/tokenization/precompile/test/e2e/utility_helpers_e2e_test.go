package tokenization_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
)

// UtilityHelpersE2ETestSuite tests the new utility helper precompile methods
type UtilityHelpersE2ETestSuite struct {
	suite.Suite
	TestSuite *helpers.TestSuite
}

func TestUtilityHelpersE2ETestSuite(t *testing.T) {
	suite.Run(t, new(UtilityHelpersE2ETestSuite))
}

func (suite *UtilityHelpersE2ETestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite()
}

// ============ rangeContains Tests ============

func (suite *UtilityHelpersE2ETestSuite) TestRangeContains_ValueInRange() {
	// Test: 10 <= 15 <= 20 should be true
	method := suite.TestSuite.Precompile.ABI.Methods["rangeContains"]
	input, err := method.Inputs.Pack(
		big.NewInt(10), // start
		big.NewInt(20), // end
		big.NewInt(15), // value
	)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err, "failed to unpack result")
	suite.Len(unpacked, 1)
	contains := unpacked[0].(bool)
	suite.True(contains, "value 15 should be in range [10, 20]")
}

func (suite *UtilityHelpersE2ETestSuite) TestRangeContains_ValueAtStart() {
	// Test: 10 <= 10 <= 20 should be true (inclusive start)
	method := suite.TestSuite.Precompile.ABI.Methods["rangeContains"]
	input, err := method.Inputs.Pack(
		big.NewInt(10), // start
		big.NewInt(20), // end
		big.NewInt(10), // value == start
	)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	contains := unpacked[0].(bool)
	suite.True(contains, "value 10 should be in range [10, 20] (inclusive start)")
}

func (suite *UtilityHelpersE2ETestSuite) TestRangeContains_ValueAtEnd() {
	// Test: 10 <= 20 <= 20 should be true (inclusive end)
	method := suite.TestSuite.Precompile.ABI.Methods["rangeContains"]
	input, err := method.Inputs.Pack(
		big.NewInt(10), // start
		big.NewInt(20), // end
		big.NewInt(20), // value == end
	)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	contains := unpacked[0].(bool)
	suite.True(contains, "value 20 should be in range [10, 20] (inclusive end)")
}

func (suite *UtilityHelpersE2ETestSuite) TestRangeContains_ValueBeforeRange() {
	// Test: 10 <= 5 should be false
	method := suite.TestSuite.Precompile.ABI.Methods["rangeContains"]
	input, err := method.Inputs.Pack(
		big.NewInt(10), // start
		big.NewInt(20), // end
		big.NewInt(5),  // value < start
	)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	contains := unpacked[0].(bool)
	suite.False(contains, "value 5 should NOT be in range [10, 20]")
}

func (suite *UtilityHelpersE2ETestSuite) TestRangeContains_ValueAfterRange() {
	// Test: 25 <= 20 should be false
	method := suite.TestSuite.Precompile.ABI.Methods["rangeContains"]
	input, err := method.Inputs.Pack(
		big.NewInt(10), // start
		big.NewInt(20), // end
		big.NewInt(25), // value > end
	)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	contains := unpacked[0].(bool)
	suite.False(contains, "value 25 should NOT be in range [10, 20]")
}

// ============ rangesOverlap Tests ============

func (suite *UtilityHelpersE2ETestSuite) TestRangesOverlap_Overlapping() {
	// Test: [10, 20] and [15, 25] should overlap
	method := suite.TestSuite.Precompile.ABI.Methods["rangesOverlap"]
	input, err := method.Inputs.Pack(
		big.NewInt(10), // start1
		big.NewInt(20), // end1
		big.NewInt(15), // start2
		big.NewInt(25), // end2
	)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	overlap := unpacked[0].(bool)
	suite.True(overlap, "[10, 20] and [15, 25] should overlap")
}

func (suite *UtilityHelpersE2ETestSuite) TestRangesOverlap_Touching() {
	// Test: [10, 20] and [20, 30] should overlap (touching at boundary)
	method := suite.TestSuite.Precompile.ABI.Methods["rangesOverlap"]
	input, err := method.Inputs.Pack(
		big.NewInt(10), // start1
		big.NewInt(20), // end1
		big.NewInt(20), // start2
		big.NewInt(30), // end2
	)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	overlap := unpacked[0].(bool)
	suite.True(overlap, "[10, 20] and [20, 30] should overlap at boundary")
}

func (suite *UtilityHelpersE2ETestSuite) TestRangesOverlap_NoOverlap() {
	// Test: [10, 20] and [25, 35] should not overlap
	method := suite.TestSuite.Precompile.ABI.Methods["rangesOverlap"]
	input, err := method.Inputs.Pack(
		big.NewInt(10), // start1
		big.NewInt(20), // end1
		big.NewInt(25), // start2
		big.NewInt(35), // end2
	)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	overlap := unpacked[0].(bool)
	suite.False(overlap, "[10, 20] and [25, 35] should NOT overlap")
}

func (suite *UtilityHelpersE2ETestSuite) TestRangesOverlap_ContainedRange() {
	// Test: [10, 30] and [15, 25] should overlap (one contains the other)
	method := suite.TestSuite.Precompile.ABI.Methods["rangesOverlap"]
	input, err := method.Inputs.Pack(
		big.NewInt(10), // start1
		big.NewInt(30), // end1
		big.NewInt(15), // start2
		big.NewInt(25), // end2
	)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	overlap := unpacked[0].(bool)
	suite.True(overlap, "[10, 30] and [15, 25] should overlap (contained)")
}

// ============ searchInRanges Tests ============

func (suite *UtilityHelpersE2ETestSuite) TestSearchInRanges_Found() {
	// Test: value 50 in ranges [{"start":"1","end":"100"}]
	method := suite.TestSuite.Precompile.ABI.Methods["searchInRanges"]
	rangesJson := `[{"start":"1","end":"100"}]`
	input, err := method.Inputs.Pack(rangesJson, big.NewInt(50))
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	found := unpacked[0].(bool)
	suite.True(found, "value 50 should be found in range [1, 100]")
}

func (suite *UtilityHelpersE2ETestSuite) TestSearchInRanges_NotFound() {
	// Test: value 150 not in ranges [{"start":"1","end":"100"}]
	method := suite.TestSuite.Precompile.ABI.Methods["searchInRanges"]
	rangesJson := `[{"start":"1","end":"100"}]`
	input, err := method.Inputs.Pack(rangesJson, big.NewInt(150))
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	found := unpacked[0].(bool)
	suite.False(found, "value 150 should NOT be found in range [1, 100]")
}

func (suite *UtilityHelpersE2ETestSuite) TestSearchInRanges_MultipleRanges() {
	// Test: value 75 in multiple ranges
	method := suite.TestSuite.Precompile.ABI.Methods["searchInRanges"]
	rangesJson := `[{"start":"1","end":"50"},{"start":"60","end":"100"}]`
	input, err := method.Inputs.Pack(rangesJson, big.NewInt(75))
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	found := unpacked[0].(bool)
	suite.True(found, "value 75 should be found in second range [60, 100]")
}

func (suite *UtilityHelpersE2ETestSuite) TestSearchInRanges_ValueInGap() {
	// Test: value 55 NOT in gap between ranges
	method := suite.TestSuite.Precompile.ABI.Methods["searchInRanges"]
	rangesJson := `[{"start":"1","end":"50"},{"start":"60","end":"100"}]`
	input, err := method.Inputs.Pack(rangesJson, big.NewInt(55))
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	found := unpacked[0].(bool)
	suite.False(found, "value 55 should NOT be found (in gap between ranges)")
}

func (suite *UtilityHelpersE2ETestSuite) TestSearchInRanges_EmptyArray() {
	// Test: empty ranges array
	method := suite.TestSuite.Precompile.ABI.Methods["searchInRanges"]
	rangesJson := `[]`
	input, err := method.Inputs.Pack(rangesJson, big.NewInt(50))
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	found := unpacked[0].(bool)
	suite.False(found, "value should NOT be found in empty ranges array")
}

// ============ getBalanceForIdAndTime Tests ============

func (suite *UtilityHelpersE2ETestSuite) TestGetBalanceForIdAndTime_Found() {
	// Test: Find balance for token ID 5 at time 1000
	method := suite.TestSuite.Precompile.ABI.Methods["getBalanceForIdAndTime"]
	balancesJson := `[{"amount":"100","badgeIds":[{"start":"1","end":"10"}],"ownershipTimes":[{"start":"0","end":"2000"}]}]`
	input, err := method.Inputs.Pack(balancesJson, big.NewInt(5), big.NewInt(1000))
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	amount := unpacked[0].(*big.Int)
	suite.True(amount.Cmp(big.NewInt(100)) == 0, "should find balance of 100")
}

func (suite *UtilityHelpersE2ETestSuite) TestGetBalanceForIdAndTime_NotFound() {
	// Test: Token ID 15 not in range [1, 10]
	method := suite.TestSuite.Precompile.ABI.Methods["getBalanceForIdAndTime"]
	balancesJson := `[{"amount":"100","badgeIds":[{"start":"1","end":"10"}],"ownershipTimes":[{"start":"0","end":"2000"}]}]`
	input, err := method.Inputs.Pack(balancesJson, big.NewInt(15), big.NewInt(1000))
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	amount := unpacked[0].(*big.Int)
	suite.True(amount.Cmp(big.NewInt(0)) == 0, "should return 0 for non-matching token ID")
}

func (suite *UtilityHelpersE2ETestSuite) TestGetBalanceForIdAndTime_TimeNotInRange() {
	// Test: Time 3000 not in range [0, 2000]
	method := suite.TestSuite.Precompile.ABI.Methods["getBalanceForIdAndTime"]
	balancesJson := `[{"amount":"100","badgeIds":[{"start":"1","end":"10"}],"ownershipTimes":[{"start":"0","end":"2000"}]}]`
	input, err := method.Inputs.Pack(balancesJson, big.NewInt(5), big.NewInt(3000))
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	amount := unpacked[0].(*big.Int)
	suite.True(amount.Cmp(big.NewInt(0)) == 0, "should return 0 for non-matching time")
}

// ============ getReservedListId Tests ============

func (suite *UtilityHelpersE2ETestSuite) TestGetReservedListId_ValidAddress() {
	// Test: Get reserved list ID for Alice's address
	method := suite.TestSuite.Precompile.ABI.Methods["getReservedListId"]
	input, err := method.Inputs.Pack(suite.TestSuite.AliceEVM)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	listId := unpacked[0].(string)
	suite.NotEmpty(listId, "should return a non-empty list ID")
	// The list ID should be the bech32 address
	suite.Contains(listId, "bb1", "list ID should be a bech32 address starting with bb1")
}

func (suite *UtilityHelpersE2ETestSuite) TestGetReservedListId_ZeroAddress() {
	// Test: Get reserved list ID for zero address
	method := suite.TestSuite.Precompile.ABI.Methods["getReservedListId"]
	zeroAddr := common.Address{}
	input, err := method.Inputs.Pack(zeroAddr)
	suite.NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(method.ID, input...))
	suite.NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	listId := unpacked[0].(string)
	// Zero address should still produce a valid bech32 address
	suite.Contains(listId, "bb1", "zero address should still produce a bb1 address")
}

// ============ Address Conversion Round-Trip Tests ============

func (suite *UtilityHelpersE2ETestSuite) TestAddressConversion_RoundTrip() {
	// Test that EVM -> Bech32 -> EVM round-trip works correctly
	evmToBech32 := suite.TestSuite.Precompile.ABI.Methods["convertEvmAddressToBech32"]
	bech32ToEvm := suite.TestSuite.Precompile.ABI.Methods["convertBech32ToEvmAddress"]

	// Convert EVM to Bech32
	input1, err := evmToBech32.Inputs.Pack(suite.TestSuite.AliceEVM)
	suite.NoError(err)

	result1, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(evmToBech32.ID, input1...))
	suite.NoError(err)

	unpacked1, err := evmToBech32.Outputs.Unpack(result1)
	suite.NoError(err)
	suite.Len(unpacked1, 1)
	bech32Addr := unpacked1[0].(string)
	suite.Contains(bech32Addr, "bb1", "should produce bb1 address")

	// Convert back to EVM
	input2, err := bech32ToEvm.Inputs.Pack(bech32Addr)
	suite.NoError(err)

	result2, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, append(bech32ToEvm.ID, input2...))
	suite.NoError(err)

	unpacked2, err := bech32ToEvm.Outputs.Unpack(result2)
	suite.NoError(err)
	suite.Len(unpacked2, 1)
	evmAddrBack := unpacked2[0].(common.Address)
	suite.Equal(suite.TestSuite.AliceEVM, evmAddrBack, "round-trip should return original address")
}
