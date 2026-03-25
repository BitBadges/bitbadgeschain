package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// GetBalanceForTokenTestSuite tests the GetBalanceForToken GRPC query handler
type GetBalanceForTokenTestSuite struct {
	TestSuite
}

func TestGetBalanceForTokenTestSuite(t *testing.T) {
	suite.Run(t, new(GetBalanceForTokenTestSuite))
}

// =============================================================================
// Suite-based tests
// =============================================================================

func (suite *GetBalanceForTokenTestSuite) TestGetBalanceForToken_Success() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query bob's balance for token ID 1 (with default time = block time)
	response, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "1",
		Address:      bob,
		TokenId:      "1",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	// Bob was minted 1 of all tokens, so balance should be "1"
	suite.Require().Equal("1", response.Balance)
}

func (suite *GetBalanceForTokenTestSuite) TestGetBalanceForToken_WithExplicitTime() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query with explicit time (current block time in ms)
	blockTimeMs := fmt.Sprintf("%d", suite.ctx.BlockTime().UnixMilli())
	response, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "1",
		Address:      bob,
		TokenId:      "1",
		Time:         blockTimeMs,
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().Equal("1", response.Balance)
}

func (suite *GetBalanceForTokenTestSuite) TestGetBalanceForToken_DefaultTimeEmptyString() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query with empty time string (should default to block time)
	response, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "1",
		Address:      bob,
		TokenId:      "1",
		Time:         "",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().Equal("1", response.Balance)
}

func (suite *GetBalanceForTokenTestSuite) TestGetBalanceForToken_DefaultTimeZero() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query with "0" time (should default to block time)
	response, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "1",
		Address:      bob,
		TokenId:      "1",
		Time:         "0",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	suite.Require().Equal("1", response.Balance)
}

func (suite *GetBalanceForTokenTestSuite) TestGetBalanceForToken_NonExistentToken() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Query alice's balance (she has no tokens) for a specific token
	response, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "1",
		Address:      alice,
		TokenId:      "1",
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(response)
	// Alice has no tokens, so balance should be "0"
	suite.Require().Equal("0", response.Balance)
}

func (suite *GetBalanceForTokenTestSuite) TestGetBalanceForToken_NilRequest() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, nil)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "invalid request")
}

func (suite *GetBalanceForTokenTestSuite) TestGetBalanceForToken_CollectionNotFound() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "999",
		Address:      bob,
		TokenId:      "1",
	})
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "not found")
}

func (suite *GetBalanceForTokenTestSuite) TestGetBalanceForToken_AfterTransfer() {
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection with tokens minted to bob
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Transfer token 1 from bob to alice
	_, err = suite.msgServer.TransferTokens(wctx, &types.MsgTransferTokens{
		Creator:      bob,
		CollectionId: sdkmath.NewUint(1),
		Transfers: []*types.Transfer{{
			From:        bob,
			ToAddresses: []string{alice},
			Balances: []*types.Balance{{
				Amount:         sdkmath.NewUint(1),
				TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
				OwnershipTimes: GetFullUintRanges(),
			}},
		}},
	})
	suite.Require().NoError(err)

	// Bob should now have 0 for token 1
	response, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "1",
		Address:      bob,
		TokenId:      "1",
	})
	suite.Require().NoError(err)
	suite.Require().Equal("0", response.Balance)

	// Alice should now have 1 for token 1
	response, err = suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "1",
		Address:      alice,
		TokenId:      "1",
	})
	suite.Require().NoError(err)
	suite.Require().Equal("1", response.Balance)
}

// =============================================================================
// Standalone tests (using keepertest pattern)
// =============================================================================

func TestGetBalanceForTokenQuery(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(suite, wctx, collectionsToCreate)
	require.NoError(t, err)

	// Query bob's balance for token 1
	response, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "1",
		Address:      bob,
		TokenId:      "1",
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, "1", response.Balance)
}

func TestGetBalanceForTokenQuery_NonExistentReturnsZero(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	wctx := sdk.WrapSDKContext(suite.ctx)

	// Create a collection
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	err := CreateCollections(suite, wctx, collectionsToCreate)
	require.NoError(t, err)

	// Query alice's balance (she has no tokens)
	response, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "1",
		Address:      alice,
		TokenId:      "1",
	})
	require.NoError(t, err)
	require.NotNil(t, response)
	require.Equal(t, "0", response.Balance)
}

func TestGetBalanceForTokenQuery_CollectionNotFound(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	wctx := sdk.WrapSDKContext(suite.ctx)

	_, err := suite.app.TokenizationKeeper.GetBalanceForToken(wctx, &types.QueryGetBalanceForTokenRequest{
		CollectionId: "999",
		Address:      bob,
		TokenId:      "1",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}
