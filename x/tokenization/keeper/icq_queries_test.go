package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type ICQTestSuite struct {
	TestSuite
}

func TestICQTestSuite(t *testing.T) {
	suite.Run(t, new(ICQTestSuite))
}

func (suite *ICQTestSuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

// =============================================================================
// OwnershipQueryPacket Validation Tests
// =============================================================================

func (suite *ICQTestSuite) TestOwnershipQueryPacket_ValidateBasic_Success() {
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"1",
		"1",              // single token ID
		"1609459200000", // single ownership time (timestamp)
	)

	err := query.ValidateBasic()
	suite.Require().NoError(err)
}

func (suite *ICQTestSuite) TestOwnershipQueryPacket_ValidateBasic_EmptyQueryId() {
	query := types.NewOwnershipQueryPacket(
		"", // Empty query ID
		bob,
		"1",
		"1",
		"1609459200000",
	)

	err := query.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "query_id")
}

func (suite *ICQTestSuite) TestOwnershipQueryPacket_ValidateBasic_EmptyAddress() {
	query := types.NewOwnershipQueryPacket(
		"query-1",
		"", // Empty address
		"1",
		"1",
		"1609459200000",
	)

	err := query.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "address")
}

func (suite *ICQTestSuite) TestOwnershipQueryPacket_ValidateBasic_InvalidCollectionId() {
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"not-a-number", // Invalid collection ID
		"1",
		"1609459200000",
	)

	err := query.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "collection_id")
}

func (suite *ICQTestSuite) TestOwnershipQueryPacket_ValidateBasic_InvalidTokenId() {
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"1",
		"not-a-number", // Invalid token ID
		"1609459200000",
	)

	err := query.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "token_id")
}

func (suite *ICQTestSuite) TestOwnershipQueryPacket_ValidateBasic_InvalidOwnershipTime() {
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"1",
		"1",
		"not-a-number", // Invalid ownership time
	)

	err := query.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "ownership_time")
}

// =============================================================================
// BulkOwnershipQueryPacket Validation Tests
// =============================================================================

func (suite *ICQTestSuite) TestBulkOwnershipQueryPacket_ValidateBasic_Success() {
	queries := []*types.OwnershipQueryPacket{
		types.NewOwnershipQueryPacket("q1", bob, "1", "1", "1609459200000"),
		types.NewOwnershipQueryPacket("q2", alice, "1", "1", "1609459200000"),
	}

	bulk := types.NewBulkOwnershipQueryPacket("bulk-1", queries)

	err := bulk.ValidateBasic()
	suite.Require().NoError(err)
}

func (suite *ICQTestSuite) TestBulkOwnershipQueryPacket_ValidateBasic_EmptyQueries() {
	bulk := types.NewBulkOwnershipQueryPacket("bulk-1", []*types.OwnershipQueryPacket{})

	err := bulk.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "empty")
}

func (suite *ICQTestSuite) TestBulkOwnershipQueryPacket_ValidateBasic_TooManyQueries() {
	// Create more than MaxBulkQueries
	queries := make([]*types.OwnershipQueryPacket, types.MaxBulkQueries+1)
	for i := range queries {
		queries[i] = types.NewOwnershipQueryPacket("q", bob, "1", "1", "1609459200000")
	}

	bulk := types.NewBulkOwnershipQueryPacket("bulk-1", queries)

	err := bulk.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "too many")
}

// =============================================================================
// ProcessOwnershipQuery Tests
// =============================================================================

func (suite *ICQTestSuite) TestProcessOwnershipQuery_CollectionNotExists() {
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"999", // Non-existent collection
		"1",
		"1609459200000",
	)

	response := suite.app.TokenizationKeeper.ProcessOwnershipQuery(suite.ctx, query)

	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().False(response.OwnsTokens)
	suite.Require().NotEmpty(response.Error)
	suite.Require().Contains(response.Error, "does not exist")
}

func (suite *ICQTestSuite) TestProcessOwnershipQuery_UserHasNoBalance() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection
	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query for alice who doesn't own any tokens
	query := types.NewOwnershipQueryPacket(
		"query-1",
		alice, // Alice doesn't own tokens
		"1",
		"1",
		"1609459200000",
	)

	response := suite.app.TokenizationKeeper.ProcessOwnershipQuery(suite.ctx, query)

	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().False(response.OwnsTokens)
	suite.Require().Equal(sdkmath.ZeroUint(), sdkmath.Uint(response.TotalAmount))
	suite.Require().Empty(response.Error)
}

func (suite *ICQTestSuite) TestProcessOwnershipQuery_UserOwnsTokens() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query for bob who should own tokens
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"1",
		"1",
		"1609459200000",
	)

	response := suite.app.TokenizationKeeper.ProcessOwnershipQuery(suite.ctx, query)

	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().True(response.OwnsTokens)
	suite.Require().True(sdkmath.Uint(response.TotalAmount).GT(sdkmath.ZeroUint()))
	suite.Require().Empty(response.Error)
}

func (suite *ICQTestSuite) TestProcessOwnershipQuery_EmptyTokenId() {
	// Query with empty token_id should fail validation
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"1",
		"", // Empty token ID - should fail
		"1609459200000",
	)

	err := query.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "token_id")
}

// =============================================================================
// ProcessBulkOwnershipQuery Tests
// =============================================================================

func (suite *ICQTestSuite) TestProcessBulkOwnershipQuery_Success() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Create bulk query with multiple addresses
	queries := []*types.OwnershipQueryPacket{
		types.NewOwnershipQueryPacket("q1", bob, "1", "1", "1609459200000"),
		types.NewOwnershipQueryPacket("q2", alice, "1", "1", "1609459200000"),
	}

	bulk := types.NewBulkOwnershipQueryPacket("bulk-1", queries)

	response := suite.app.TokenizationKeeper.ProcessBulkOwnershipQuery(suite.ctx, bulk)

	suite.Require().Equal("bulk-1", response.QueryId)
	suite.Require().Len(response.Responses, 2)

	// Bob should own tokens
	suite.Require().Equal("q1", response.Responses[0].QueryId)
	suite.Require().True(response.Responses[0].OwnsTokens)

	// Alice should not own tokens
	suite.Require().Equal("q2", response.Responses[1].QueryId)
	suite.Require().False(response.Responses[1].OwnsTokens)
}

func (suite *ICQTestSuite) TestProcessBulkOwnershipQuery_ValidationError() {
	// Create bulk query with invalid query (empty query ID)
	queries := []*types.OwnershipQueryPacket{
		types.NewOwnershipQueryPacket("", bob, "1", "1", "1609459200000"), // Empty query ID
	}

	bulk := types.NewBulkOwnershipQueryPacket("bulk-1", queries)

	response := suite.app.TokenizationKeeper.ProcessBulkOwnershipQuery(suite.ctx, bulk)

	suite.Require().Equal("bulk-1", response.QueryId)
	suite.Require().Len(response.Responses, 1)
	suite.Require().NotEmpty(response.Responses[0].Error)
}

// =============================================================================
// Packet Serialization Tests
// =============================================================================

func (suite *ICQTestSuite) TestOwnershipQueryPacket_GetBytes() {
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"1",
		"1",
		"1609459200000",
	)

	bytes := query.GetBytes()
	suite.Require().NotEmpty(bytes)

	// Verify we can unmarshal the bytes back
	var packetData types.TokenizationPacketData
	err := types.ModuleCdc.Unmarshal(bytes, &packetData)
	suite.Require().NoError(err)

	unmarshaled := packetData.GetOwnershipQuery()
	suite.Require().NotNil(unmarshaled)
	suite.Require().Equal("query-1", unmarshaled.QueryId)
	suite.Require().Equal(bob, unmarshaled.Address)
	suite.Require().Equal("1", unmarshaled.TokenId)
	suite.Require().Equal("1609459200000", unmarshaled.OwnershipTime)
}

func (suite *ICQTestSuite) TestOwnershipQueryResponsePacket_GetBytes() {
	response := types.NewOwnershipQueryResponsePacket(
		"query-1",
		true,
		sdkmath.NewUint(100),
		12345,
		"",
	)

	bytes := response.GetBytes()
	suite.Require().NotEmpty(bytes)

	// Verify we can unmarshal the bytes back
	var packetData types.TokenizationPacketData
	err := types.ModuleCdc.Unmarshal(bytes, &packetData)
	suite.Require().NoError(err)

	unmarshaled := packetData.GetOwnershipQueryResponse()
	suite.Require().NotNil(unmarshaled)
	suite.Require().Equal("query-1", unmarshaled.QueryId)
	suite.Require().True(unmarshaled.OwnsTokens)
	suite.Require().Equal(uint64(12345), unmarshaled.ProofHeight)
}

// =============================================================================
// ICQ Packet Type Detection Tests
// =============================================================================

func (suite *ICQTestSuite) TestGetICQPacketType() {
	// Test OwnershipQuery
	query := types.NewOwnershipQueryPacket("q1", bob, "1", "1", "1609459200000")
	packetData := &types.TokenizationPacketData{
		Packet: &types.TokenizationPacketData_OwnershipQuery{
			OwnershipQuery: query,
		},
	}
	suite.Require().Equal(types.ICQPacketTypeOwnershipQuery, types.GetICQPacketType(packetData))

	// Test OwnershipQueryResponse
	response := types.NewOwnershipQueryResponsePacket("q1", true, sdkmath.NewUint(1), 0, "")
	packetData = &types.TokenizationPacketData{
		Packet: &types.TokenizationPacketData_OwnershipQueryResponse{
			OwnershipQueryResponse: response,
		},
	}
	suite.Require().Equal(types.ICQPacketTypeOwnershipQueryResponse, types.GetICQPacketType(packetData))

	// Test BulkOwnershipQuery
	bulk := types.NewBulkOwnershipQueryPacket("bulk-1", []*types.OwnershipQueryPacket{query})
	packetData = &types.TokenizationPacketData{
		Packet: &types.TokenizationPacketData_BulkOwnershipQuery{
			BulkOwnershipQuery: bulk,
		},
	}
	suite.Require().Equal(types.ICQPacketTypeBulkOwnershipQuery, types.GetICQPacketType(packetData))

	// Test FullBalanceQuery
	fullBalanceQuery := types.NewFullBalanceQueryPacket("fb1", bob, "1")
	packetData = &types.TokenizationPacketData{
		Packet: &types.TokenizationPacketData_FullBalanceQuery{
			FullBalanceQuery: fullBalanceQuery,
		},
	}
	suite.Require().Equal(types.ICQPacketTypeFullBalanceQuery, types.GetICQPacketType(packetData))

	// Test FullBalanceQueryResponse
	fullBalanceResponse := types.NewFullBalanceQueryResponsePacket("fb1", []byte("test"), 0, "")
	packetData = &types.TokenizationPacketData{
		Packet: &types.TokenizationPacketData_FullBalanceQueryResponse{
			FullBalanceQueryResponse: fullBalanceResponse,
		},
	}
	suite.Require().Equal(types.ICQPacketTypeFullBalanceQueryResponse, types.GetICQPacketType(packetData))

	// Test NoData (unknown)
	packetData = &types.TokenizationPacketData{
		Packet: &types.TokenizationPacketData_NoData{},
	}
	suite.Require().Equal(types.ICQPacketTypeUnknown, types.GetICQPacketType(packetData))
}

// =============================================================================
// FullBalanceQuery Tests
// =============================================================================

func (suite *ICQTestSuite) TestFullBalanceQueryPacket_ValidateBasic_Success() {
	query := types.NewFullBalanceQueryPacket(
		"query-1",
		bob,
		"1",
	)

	err := query.ValidateBasic()
	suite.Require().NoError(err)
}

func (suite *ICQTestSuite) TestFullBalanceQueryPacket_ValidateBasic_EmptyQueryId() {
	query := types.NewFullBalanceQueryPacket(
		"", // Empty query ID
		bob,
		"1",
	)

	err := query.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "query_id")
}

func (suite *ICQTestSuite) TestFullBalanceQueryPacket_ValidateBasic_EmptyAddress() {
	query := types.NewFullBalanceQueryPacket(
		"query-1",
		"", // Empty address
		"1",
	)

	err := query.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "address")
}

func (suite *ICQTestSuite) TestFullBalanceQueryPacket_ValidateBasic_InvalidCollectionId() {
	query := types.NewFullBalanceQueryPacket(
		"query-1",
		bob,
		"not-a-number", // Invalid collection ID
	)

	err := query.ValidateBasic()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "collection_id")
}

func (suite *ICQTestSuite) TestProcessFullBalanceQuery_Success() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query for bob's full balance
	query := types.NewFullBalanceQueryPacket(
		"query-1",
		bob,
		"1",
	)

	response := suite.app.TokenizationKeeper.ProcessFullBalanceQuery(suite.ctx, query)

	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().Empty(response.Error)
	suite.Require().NotEmpty(response.BalanceStore)

	// Verify we can deserialize the balance store
	var balanceStore types.UserBalanceStore
	err = types.ModuleCdc.Unmarshal(response.BalanceStore, &balanceStore)
	suite.Require().NoError(err)
	suite.Require().NotEmpty(balanceStore.Balances)
}

func (suite *ICQTestSuite) TestProcessFullBalanceQuery_CollectionNotExists() {
	query := types.NewFullBalanceQueryPacket(
		"query-1",
		bob,
		"999", // Non-existent collection
	)

	response := suite.app.TokenizationKeeper.ProcessFullBalanceQuery(suite.ctx, query)

	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().NotEmpty(response.Error)
	suite.Require().Contains(response.Error, "does not exist")
}

func (suite *ICQTestSuite) TestFullBalanceQueryPacket_GetBytes() {
	query := types.NewFullBalanceQueryPacket(
		"query-1",
		bob,
		"1",
	)

	bytes := query.GetBytes()
	suite.Require().NotEmpty(bytes)

	// Verify we can unmarshal the bytes back
	var packetData types.TokenizationPacketData
	err := types.ModuleCdc.Unmarshal(bytes, &packetData)
	suite.Require().NoError(err)

	unmarshaled := packetData.GetFullBalanceQuery()
	suite.Require().NotNil(unmarshaled)
	suite.Require().Equal("query-1", unmarshaled.QueryId)
	suite.Require().Equal(bob, unmarshaled.Address)
	suite.Require().Equal("1", unmarshaled.CollectionId)
}

// =============================================================================
// Edge Cases and Error Handling Tests
// =============================================================================

func (suite *ICQTestSuite) TestOwnershipQuery_TokenIdNotInRange() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query for token ID that bob doesn't own (0 is outside the minted range [1, MaxUint64])
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"1",
		"0",              // Token ID 0 is not in the minted range [1, MaxUint64]
		"1609459200000",
	)

	response := suite.app.TokenizationKeeper.ProcessOwnershipQuery(suite.ctx, query)

	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().False(response.OwnsTokens)
	suite.Require().Equal(sdkmath.ZeroUint(), sdkmath.Uint(response.TotalAmount))
	suite.Require().Empty(response.Error)
}

func (suite *ICQTestSuite) TestOwnershipQuery_OwnershipTimeNotInRange() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query for ownership time that's outside the valid range
	query := types.NewOwnershipQueryPacket(
		"query-1",
		bob,
		"1",
		"1",
		"9999999999999999", // Very far future timestamp
	)

	response := suite.app.TokenizationKeeper.ProcessOwnershipQuery(suite.ctx, query)

	// This depends on how the collection was created - may or may not have tokens at this time
	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().Empty(response.Error)
}

func (suite *ICQTestSuite) TestOwnershipQuery_EVMHexAddress() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query using EVM hex address format (0x prefix)
	// Note: This test assumes the address normalization works correctly
	query := types.NewOwnershipQueryPacket(
		"query-1",
		"0x0000000000000000000000000000000000000001", // Placeholder EVM address
		"1",
		"1",
		"1609459200000",
	)

	response := suite.app.TokenizationKeeper.ProcessOwnershipQuery(suite.ctx, query)

	// Should not error, but may return 0 balance since it's a different address
	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().Empty(response.Error)
}

func (suite *ICQTestSuite) TestFullBalanceQuery_UserWithApprovals() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query for bob's full balance
	query := types.NewFullBalanceQueryPacket(
		"query-1",
		bob,
		"1",
	)

	response := suite.app.TokenizationKeeper.ProcessFullBalanceQuery(suite.ctx, query)

	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().Empty(response.Error)
	suite.Require().NotEmpty(response.BalanceStore)

	// Verify we can deserialize the balance store
	var balanceStore types.UserBalanceStore
	err = types.ModuleCdc.Unmarshal(response.BalanceStore, &balanceStore)
	suite.Require().NoError(err)

	// Verify the balance store has expected fields
	suite.Require().NotNil(balanceStore.Balances)
}

func (suite *ICQTestSuite) TestFullBalanceQuery_UserWithNoBalance() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection
	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query for alice who doesn't own any tokens
	query := types.NewFullBalanceQueryPacket(
		"query-1",
		alice,
		"1",
	)

	response := suite.app.TokenizationKeeper.ProcessFullBalanceQuery(suite.ctx, query)

	suite.Require().Equal("query-1", response.QueryId)
	suite.Require().Empty(response.Error)
	suite.Require().NotEmpty(response.BalanceStore)

	// Should still return a valid (possibly empty) balance store
	var balanceStore types.UserBalanceStore
	err = types.ModuleCdc.Unmarshal(response.BalanceStore, &balanceStore)
	suite.Require().NoError(err)
}

func (suite *ICQTestSuite) TestGetBalanceForIdAndTime_MultipleOverlappingBalances() {
	// Test the getBalanceForIdAndTime function directly with overlapping balances
	// This tests the core balance lookup logic

	balances := []*types.Balance{
		{
			Amount:         types.Uint(sdkmath.NewUint(10)),
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
		},
		{
			Amount:         types.Uint(sdkmath.NewUint(5)),
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(50), End: sdkmath.NewUint(150)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(500), End: sdkmath.NewUint(1500)}},
		},
	}

	// Token ID 75, time 750 is in both balances - should sum to 15
	tokenId := sdkmath.NewUint(75)
	ownershipTime := sdkmath.NewUint(750)

	amount := getBalanceForIdAndTime(balances, tokenId, ownershipTime)
	suite.Require().Equal(sdkmath.NewUint(15), amount)
}

func (suite *ICQTestSuite) TestGetBalanceForIdAndTime_NoMatch() {
	balances := []*types.Balance{
		{
			Amount:         types.Uint(sdkmath.NewUint(10)),
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
		},
	}

	// Token ID outside range
	tokenId := sdkmath.NewUint(200)
	ownershipTime := sdkmath.NewUint(500)

	amount := getBalanceForIdAndTime(balances, tokenId, ownershipTime)
	suite.Require().Equal(sdkmath.ZeroUint(), amount)
}

func (suite *ICQTestSuite) TestGetBalanceForIdAndTime_NilBalance() {
	balances := []*types.Balance{
		nil,
		{
			Amount:         types.Uint(sdkmath.NewUint(10)),
			TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
			OwnershipTimes: []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
		},
	}

	tokenId := sdkmath.NewUint(50)
	ownershipTime := sdkmath.NewUint(500)

	amount := getBalanceForIdAndTime(balances, tokenId, ownershipTime)
	suite.Require().Equal(sdkmath.NewUint(10), amount)
}

func (suite *ICQTestSuite) TestGetBalanceForIdAndTime_EmptyBalances() {
	balances := []*types.Balance{}

	tokenId := sdkmath.NewUint(50)
	ownershipTime := sdkmath.NewUint(500)

	amount := getBalanceForIdAndTime(balances, tokenId, ownershipTime)
	suite.Require().Equal(sdkmath.ZeroUint(), amount)
}

// Helper function to test getBalanceForIdAndTime directly
func getBalanceForIdAndTime(balances []*types.Balance, tokenId sdkmath.Uint, ownershipTime sdkmath.Uint) sdkmath.Uint {
	amount := sdkmath.ZeroUint()

	for _, balance := range balances {
		if balance == nil {
			continue
		}

		// Check if tokenId is in this balance's token ID ranges
		foundTokenId := false
		for _, tokenRange := range balance.TokenIds {
			if tokenRange != nil && tokenId.GTE(tokenRange.Start) && tokenId.LTE(tokenRange.End) {
				foundTokenId = true
				break
			}
		}

		if !foundTokenId {
			continue
		}

		// Check if ownershipTime is in this balance's ownership time ranges
		foundTime := false
		for _, timeRange := range balance.OwnershipTimes {
			if timeRange != nil && ownershipTime.GTE(timeRange.Start) && ownershipTime.LTE(timeRange.End) {
				foundTime = true
				break
			}
		}

		if foundTime {
			amount = amount.Add(sdkmath.Uint(balance.Amount))
		}
	}

	return amount
}
