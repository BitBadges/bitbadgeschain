package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// GRPCQueryHandlersTestSuite tests all GRPC query handlers
type GRPCQueryHandlersTestSuite struct {
	TestSuite
}

func TestGRPCQueryHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCQueryHandlersTestSuite))
}

// =============================================================================
// GetCollection Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetCollection_Success() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query the collection
	response, err := suite.app.TokenizationKeeper.GetCollection(wctx, &types.QueryGetCollectionRequest{
		CollectionId: "1",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().NotNil(response.Collection)
	suite.Require().Equal(sdkmath.NewUint(1), response.Collection.CollectionId)
}

func (suite *GRPCQueryHandlersTestSuite) TestGetCollection_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetCollection(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetCollection_NotFound() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetCollection(wctx, &types.QueryGetCollectionRequest{
		CollectionId: "999",
	})
	suite.Require().Error(err)
}

// =============================================================================
// GetBalance Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetBalance_Success() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query bob's balance
	response, err := suite.app.TokenizationKeeper.GetBalance(wctx, &types.QueryGetBalanceRequest{
		CollectionId: "1",
		Address:      bob,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().NotNil(response.Balance)
}

func (suite *GRPCQueryHandlersTestSuite) TestGetBalance_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetBalance(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetBalance_CollectionNotFound() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetBalance(wctx, &types.QueryGetBalanceRequest{
		CollectionId: "999",
		Address:      bob,
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not found")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetBalance_AddressWithNoTokens() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query alice's balance (she has no tokens)
	response, err := suite.app.TokenizationKeeper.GetBalance(wctx, &types.QueryGetBalanceRequest{
		CollectionId: "1",
		Address:      alice,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	// Should return default/empty balance
}

// =============================================================================
// GetAddressList Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetAddressList_Success() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Query a built-in list
	response, err := suite.app.TokenizationKeeper.GetAddressList(wctx, &types.QueryGetAddressListRequest{
		ListId: "AllWithoutMint",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().NotNil(response.List)
	suite.Require().Equal("AllWithoutMint", response.List.ListId)
}

func (suite *GRPCQueryHandlersTestSuite) TestGetAddressList_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetAddressList(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetAddressList_NotFound() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetAddressList(wctx, &types.QueryGetAddressListRequest{
		ListId: "non-existent-list-id",
	})
	suite.Require().Error(err)
}

// =============================================================================
// GetDynamicStore Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetDynamicStore_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetDynamicStore(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetDynamicStore_NotFound() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: "999",
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not found")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetDynamicStore_Success() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a dynamic store
	_, err := suite.msgServer.CreateDynamicStore(wctx, &types.MsgCreateDynamicStore{
		Creator:      bob,
		DefaultValue: true,
	})
	suite.Require().NoError(err)

	// Query the dynamic store
	response, err := suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: "1",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().NotNil(response.Store)
	suite.Require().Equal(sdkmath.NewUint(1), response.Store.StoreId)
	suite.Require().Equal(true, response.Store.DefaultValue)
}

// =============================================================================
// GetDynamicStoreValue Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetDynamicStoreValue_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetDynamicStoreValue_StoreNotFound() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: "999",
		Address: bob,
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not found")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetDynamicStoreValue_ReturnsDefaultValue() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a dynamic store with default value true
	_, err := suite.msgServer.CreateDynamicStore(wctx, &types.MsgCreateDynamicStore{
		Creator:      bob,
		DefaultValue: true,
	})
	suite.Require().NoError(err)

	// Query a value that hasn't been set - should return default
	response, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: "1",
		Address: alice,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().NotNil(response.Value)
	suite.Require().Equal(true, response.Value.Value)
}

func (suite *GRPCQueryHandlersTestSuite) TestGetDynamicStoreValue_ReturnsSetValue() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a dynamic store with default false
	_, err := suite.msgServer.CreateDynamicStore(wctx, &types.MsgCreateDynamicStore{
		Creator:      bob,
		DefaultValue: false,
	})
	suite.Require().NoError(err)

	// Set a value for alice to true (creator sets for address)
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, &types.MsgSetDynamicStoreValue{
		Creator: bob,
		StoreId: sdkmath.NewUint(1),
		Address: alice,
		Value:   true,
	})
	suite.Require().NoError(err)

	// Query alice's value
	response, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: "1",
		Address: alice,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().NotNil(response.Value)
	suite.Require().Equal(true, response.Value.Value)
}

// =============================================================================
// IsAddressReservedProtocol Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestIsAddressReservedProtocol_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.IsAddressReservedProtocol(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestIsAddressReservedProtocol_NotReserved() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	response, err := suite.app.TokenizationKeeper.IsAddressReservedProtocol(wctx, &types.QueryIsAddressReservedProtocolRequest{
		Address: bob,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().False(response.IsReservedProtocol)
}

func (suite *GRPCQueryHandlersTestSuite) TestIsAddressReservedProtocol_MintAddress() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// "Mint" is a reserved protocol address
	response, err := suite.app.TokenizationKeeper.IsAddressReservedProtocol(wctx, &types.QueryIsAddressReservedProtocolRequest{
		Address: "Mint",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().True(response.IsReservedProtocol)
}

// =============================================================================
// GetAllReservedProtocolAddresses Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetAllReservedProtocolAddresses_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetAllReservedProtocolAddresses(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetAllReservedProtocolAddresses_Success() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	response, err := suite.app.TokenizationKeeper.GetAllReservedProtocolAddresses(wctx, &types.QueryGetAllReservedProtocolAddressesRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	// Response should be valid (may or may not have addresses initially)
}

// =============================================================================
// GetApprovalTracker Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetApprovalTracker_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetApprovalTracker(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetApprovalTracker_NotFound() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetApprovalTracker(wctx, &types.QueryGetApprovalTrackerRequest{
		CollectionId:    "1",
		ApproverAddress: bob,
		ApprovalId:      "non-existent",
		AmountTrackerId: "non-existent",
		ApprovalLevel:   "collection",
		TrackerType:     "amount",
		ApprovedAddress: alice,
	})
	suite.Require().Error(err)
}

// =============================================================================
// GetChallengeTracker Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetChallengeTracker_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetChallengeTracker(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetChallengeTracker_NotFound() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Query for non-existent tracker - may return 0 or error depending on implementation
	response, err := suite.app.TokenizationKeeper.GetChallengeTracker(wctx, &types.QueryGetChallengeTrackerRequest{
		CollectionId:       "1",
		ApproverAddress:    bob,
		ApprovalLevel:      "collection",
		ApprovalId:         "non-existent",
		ChallengeTrackerId: "non-existent",
		LeafIndex:          "0",
	})
	// Either returns error or returns "0" for unused tracker
	if err == nil {
		suite.Require().Equal("0", response.NumUsed)
	}
}

// =============================================================================
// GetETHSignatureTracker Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetETHSignatureTracker_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetETHSignatureTracker(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetETHSignatureTracker_NotFound_ReturnsZero() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Non-existent signature should return 0
	response, err := suite.app.TokenizationKeeper.GetETHSignatureTracker(wctx, &types.QueryGetETHSignatureTrackerRequest{
		CollectionId:       "1",
		ApproverAddress:    bob,
		ApprovalLevel:      "collection",
		ApprovalId:         "test",
		ChallengeTrackerId: "test",
		Signature:          "0x1234567890",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().Equal("0", response.NumUsed)
}

// =============================================================================
// GetVote Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetVote_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetVote(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetVote_NotFound() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetVote(wctx, &types.QueryGetVoteRequest{
		CollectionId:    "1",
		ApproverAddress: bob,
		ApprovalLevel:   "collection",
		ApprovalId:      "test",
		ProposalId:      "proposal-1",
		VoterAddress:    alice,
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not found")
}

// =============================================================================
// GetVotes Tests
// =============================================================================

func (suite *GRPCQueryHandlersTestSuite) TestGetVotes_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetVotes(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GRPCQueryHandlersTestSuite) TestGetVotes_EmptyResults() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	response, err := suite.app.TokenizationKeeper.GetVotes(wctx, &types.QueryGetVotesRequest{
		CollectionId:    "1",
		ApproverAddress: bob,
		ApprovalLevel:   "collection",
		ApprovalId:      "test",
		ProposalId:      "non-existent-proposal",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().Empty(response.Votes)
}

// =============================================================================
// Standalone Tests (using keepertest pattern)
// =============================================================================

func TestGetCollectionQuery(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(suite, wctx, collectionsToCreate)
	require.NoError(t, err)

	// Query the collection
	response, err := suite.app.TokenizationKeeper.GetCollection(wctx, &types.QueryGetCollectionRequest{
		CollectionId: "1",
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, sdkmath.NewUint(1), response.Collection.CollectionId)
}

func TestGetBalanceQuery(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(suite, wctx, collectionsToCreate)
	require.NoError(t, err)

	// Query bob's balance
	response, err := suite.app.TokenizationKeeper.GetBalance(wctx, &types.QueryGetBalanceRequest{
		CollectionId: "1",
		Address:      bob,
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Balance)
}

func TestGetAddressListQuery(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Query built-in list
	response, err := suite.app.TokenizationKeeper.GetAddressList(wctx, &types.QueryGetAddressListRequest{
		ListId: "AllWithoutMint",
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, "AllWithoutMint", response.List.ListId)
}

func TestIsAddressReservedProtocolQuery(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Test regular address is not reserved
	response, err := suite.app.TokenizationKeeper.IsAddressReservedProtocol(wctx, &types.QueryIsAddressReservedProtocolRequest{
		Address: bob,
	})
	require.NoError(t, err)
	require.False(t, response.IsReservedProtocol)

	// Test another regular address is not reserved
	response, err = suite.app.TokenizationKeeper.IsAddressReservedProtocol(wctx, &types.QueryIsAddressReservedProtocolRequest{
		Address: alice,
	})
	require.NoError(t, err)
	require.False(t, response.IsReservedProtocol)
}

func TestGetAllReservedProtocolAddressesQuery(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	response, err := suite.app.TokenizationKeeper.GetAllReservedProtocolAddresses(wctx, &types.QueryGetAllReservedProtocolAddressesRequest{})
	require.NoError(t, err)
	require.NotNil(t, response)
	// Response is valid (may or may not contain addresses initially)
}

func TestDynamicStoreQueries(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)
	msgServer := suite.msgServer

	// Create a dynamic store with default true
	_, err := msgServer.CreateDynamicStore(wctx, &types.MsgCreateDynamicStore{
		Creator:      bob,
		DefaultValue: true,
	})
	require.NoError(t, err)

	// Test GetDynamicStore
	storeResponse, err := suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: "1",
	})
	require.NoError(t, err)
	require.Equal(t, true, storeResponse.Store.DefaultValue)

	// Test GetDynamicStoreValue returns default
	valueResponse, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: "1",
		Address: alice,
	})
	require.NoError(t, err)
	require.Equal(t, true, valueResponse.Value.Value)

	// Set a value to false (creator sets for address)
	_, err = msgServer.SetDynamicStoreValue(wctx, &types.MsgSetDynamicStoreValue{
		Creator: bob,
		StoreId: sdkmath.NewUint(1),
		Address: alice,
		Value:   false,
	})
	require.NoError(t, err)

	// Test GetDynamicStoreValue returns custom value
	valueResponse, err = suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: "1",
		Address: alice,
	})
	require.NoError(t, err)
	require.Equal(t, false, valueResponse.Value.Value)
}

func TestGetETHSignatureTrackerQuery(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Query non-existent signature - should return 0
	response, err := suite.app.TokenizationKeeper.GetETHSignatureTracker(wctx, &types.QueryGetETHSignatureTrackerRequest{
		CollectionId:       "1",
		ApproverAddress:    bob,
		ApprovalLevel:      "collection",
		ApprovalId:         "test",
		ChallengeTrackerId: "test",
		Signature:          "0xabcdef",
	})
	require.NoError(t, err)
	require.Equal(t, "0", response.NumUsed)
}

func TestGetVotesQuery(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Query votes for non-existent proposal - should return empty
	response, err := suite.app.TokenizationKeeper.GetVotes(wctx, &types.QueryGetVotesRequest{
		CollectionId:    "1",
		ApproverAddress: bob,
		ApprovalLevel:   "collection",
		ApprovalId:      "test",
		ProposalId:      "non-existent",
	})
	require.NoError(t, err)
	require.Empty(t, response.Votes)
}
