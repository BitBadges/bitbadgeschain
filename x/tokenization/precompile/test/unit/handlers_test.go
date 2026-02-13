package tokenization_test

import (
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

type HandlersTestSuite struct {
	suite.Suite
	TestSuite  *helpers.TestSuite
	Precompile *tokenization.Precompile
}

func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}

func (suite *HandlersTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite()
	suite.Precompile = suite.TestSuite.Precompile
}

func (suite *HandlersTestSuite) createContract(caller common.Address) *vm.Contract {
	precompileAddr := common.HexToAddress(tokenization.TokenizationPrecompileAddress)
	valueUint256, _ := uint256.FromBig(big.NewInt(0))
	return vm.NewContract(caller, precompileAddr, valueUint256, 1000000, nil)
}

func (suite *HandlersTestSuite) TestCreateCollection_Valid() {
	caller := suite.TestSuite.AliceEVM

	// Build JSON message
	msg := map[string]interface{}{
		"defaultBalances":            nil,
		"validTokenIds":              []map[string]interface{}{{"start": "1", "end": "100"}},
		"collectionPermissions":     map[string]interface{}{},
		"manager":                    suite.TestSuite.Manager.String(),
		"collectionMetadata":         map[string]interface{}{"uri": "https://example.com", "customData": "test"},
		"tokenMetadata":              []interface{}{},
		"customData":                 "custom data",
		"collectionApprovals":        []interface{}{},
		"standards":                  []string{"ERC721"},
		"isArchived":                 false,
		"mintEscrowCoinsToTransfer":  []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                 map[string]interface{}{},
		"aliasPathsToAdd":            []interface{}{},
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

	// Unpack result
	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)

	collectionIdBig, ok := unpacked[0].(*big.Int)
	suite.True(ok)
	suite.Greater(collectionIdBig.Uint64(), uint64(0))
}

func (suite *HandlersTestSuite) TestCreateCollection_InvalidTokenIds() {
	caller := suite.TestSuite.AliceEVM

	// Invalid: start > end
	msg := map[string]interface{}{
		"defaultBalances":            nil,
		"validTokenIds":              []map[string]interface{}{{"start": "100", "end": "1"}}, // Invalid: start > end
		"collectionPermissions":     map[string]interface{}{},
		"manager":                    suite.TestSuite.Manager.String(),
		"collectionMetadata":         map[string]interface{}{"uri": "", "customData": ""},
		"tokenMetadata":              []interface{}{},
		"customData":                 "",
		"collectionApprovals":        []interface{}{},
		"standards":                  []string{},
		"isArchived":                 false,
		"mintEscrowCoinsToTransfer":  []interface{}{},
		"cosmosCoinWrapperPathsToAdd": []interface{}{},
		"invariants":                 map[string]interface{}{},
		"aliasPathsToAdd":            []interface{}{},
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

func (suite *HandlersTestSuite) TestDeleteCollection_Valid() {
	// Create collection first
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM

	// Build JSON message
	jsonMsg, err := helpers.BuildDeleteCollectionJSON(suite.TestSuite.Alice.String(), collectionId.BigInt())
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["deleteCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify collection is deleted
	_, found := suite.TestSuite.Keeper.GetCollectionFromStore(suite.TestSuite.Ctx, collectionId)
	suite.False(found, "Collection should be deleted")
}

func (suite *HandlersTestSuite) TestDeleteCollection_Unauthorized() {
	// Create collection as Alice
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	// Try to delete as Bob (should fail)
	caller := suite.TestSuite.BobEVM

	// Build JSON message (Bob trying to delete Alice's collection)
	jsonMsg, err := helpers.BuildDeleteCollectionJSON(suite.TestSuite.Bob.String(), collectionId.BigInt())
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["deleteCollection"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.Error(err)
	suite.Nil(result)

	// Collection should still exist
	_, found := suite.TestSuite.Keeper.GetCollectionFromStore(suite.TestSuite.Ctx, collectionId)
	suite.True(found, "Collection should still exist")
}

func (suite *HandlersTestSuite) TestTransferTokens_Valid() {
	// Create collection and balance
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

	// Also create balance for Bob (needed for incoming approvals)
	err = suite.TestSuite.CreateTestBalance(
		collectionId,
		suite.TestSuite.Bob.String(),
		sdkmath.NewUint(0), // Bob starts with 0, will receive tokens
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
	)
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM

	// Convert EVM addresses to Cosmos addresses
	toCosmos := suite.TestSuite.Bob.String()

	// Build JSON message
	jsonMsg, err := helpers.BuildTransferTokensJSON(
		collectionId.BigInt(),
		suite.TestSuite.Alice.String(),
		[]string{toCosmos},
		big.NewInt(100),
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(10)}},
		[]struct{ Start, End *big.Int }{{Start: big.NewInt(1), End: big.NewInt(1000)}},
	)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["transferTokens"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify transfer
	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)
	success, ok := unpacked[0].(bool)
	suite.True(ok)
	suite.True(success)
}

func (suite *HandlersTestSuite) TestTransferTokens_InsufficientBalance() {
	// Create collection with small balance
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	err = suite.TestSuite.CreateTestBalance(
		collectionId,
		suite.TestSuite.Alice.String(),
		sdkmath.NewUint(50), // Only 50 tokens
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)}},
	)
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM

	// Convert EVM addresses to Cosmos addresses
	toCosmos := suite.TestSuite.Bob.String()

	// Build JSON message (trying to transfer more than available)
	jsonMsg, err := helpers.BuildTransferTokensJSON(
		collectionId.BigInt(),
		suite.TestSuite.Alice.String(),
		[]string{toCosmos},
		big.NewInt(100), // More than balance
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

func (suite *HandlersTestSuite) TestSetManager_Valid() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM

	// Build JSON message
	jsonMsg, err := helpers.BuildSetManagerJSON(suite.TestSuite.Alice.String(), collectionId.BigInt(), suite.TestSuite.Bob.String())
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["setManager"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify manager was updated
	req := &tokenizationtypes.QueryGetCollectionRequest{
		CollectionId: collectionId.String(),
	}
	resp, err := suite.TestSuite.Keeper.GetCollection(suite.TestSuite.Ctx, req)
	suite.NoError(err)
	suite.Equal(suite.TestSuite.Bob.String(), resp.Collection.Manager)
}

func (suite *HandlersTestSuite) TestSetCollectionMetadata_Valid() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM

	// Build JSON message
	metadata := map[string]interface{}{
		"uri":        "https://new-uri.com",
		"customData": "new custom data",
	}
	jsonMsg, err := helpers.BuildSetCollectionMetadataJSON(suite.TestSuite.Alice.String(), collectionId.BigInt(), metadata)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["setCollectionMetadata"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify metadata was updated
	req := &tokenizationtypes.QueryGetCollectionRequest{
		CollectionId: collectionId.String(),
	}
	resp, err := suite.TestSuite.Keeper.GetCollection(suite.TestSuite.Ctx, req)
	suite.NoError(err)
	suite.Equal("https://new-uri.com", resp.Collection.CollectionMetadata.Uri)
	suite.Equal("new custom data", resp.Collection.CollectionMetadata.CustomData)
}

func (suite *HandlersTestSuite) TestSetValidTokenIds_Valid() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM

	// Build JSON message
	validTokenIds := []map[string]interface{}{
		{"start": "1", "end": "200"},
	}
	jsonMsg, err := helpers.BuildSetValidTokenIdsJSON(suite.TestSuite.Alice.String(), collectionId.BigInt(), validTokenIds)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["setValidTokenIds"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify token IDs were updated
	req := &tokenizationtypes.QueryGetCollectionRequest{
		CollectionId: collectionId.String(),
	}
	resp, err := suite.TestSuite.Keeper.GetCollection(suite.TestSuite.Ctx, req)
	suite.NoError(err)
	suite.Len(resp.Collection.ValidTokenIds, 1)
	suite.Equal(sdkmath.NewUint(1), resp.Collection.ValidTokenIds[0].Start)
	suite.Equal(sdkmath.NewUint(200), resp.Collection.ValidTokenIds[0].End)
}

func (suite *HandlersTestSuite) TestCreateDynamicStore_Valid() {
	caller := suite.TestSuite.AliceEVM

	// Build JSON message
	jsonMsg, err := helpers.BuildCreateDynamicStoreJSON(
		suite.TestSuite.Alice.String(),
		false,                 // defaultValue
		"https://example.com", // uri
		"custom data",         // customData
	)
	suite.NoError(err)

	method := suite.Precompile.ABI.Methods["createDynamicStore"]
	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Unpack result
	unpacked, err := method.Outputs.Unpack(result)
	suite.NoError(err)
	suite.Len(unpacked, 1)

	storeIdBig, ok := unpacked[0].(*big.Int)
	suite.True(ok)
	suite.Greater(storeIdBig.Uint64(), uint64(0))
}

func (suite *HandlersTestSuite) TestCreateAddressLists_Valid() {
	caller := suite.TestSuite.AliceEVM

	addressListInput := map[string]interface{}{
		"listId":     "testlist", // Use alphanumeric only (no hyphens)
		"addresses":  []string{suite.TestSuite.Bob.String(), suite.TestSuite.Charlie.String()},
		"whitelist":  true,
		"uri":        "https://example.com",
		"customData": "data",
	}

	// Build JSON message
	jsonMsg, err := helpers.BuildCreateAddressListsJSON(
		suite.TestSuite.Alice.String(),
		[]map[string]interface{}{addressListInput},
	)
	suite.NoError(err)

	method, found := suite.Precompile.ABI.Methods["createAddressLists"]
	if !found {
		suite.T().Skip("createAddressLists method not found in ABI")
		return
	}

	input, err := helpers.PackMethodWithJSON(&method, jsonMsg)
	suite.NoError(err)

	contract := suite.TestSuite.CreateMockContract(caller, input)
	result, err := suite.Precompile.Execute(suite.TestSuite.Ctx, contract, false)
	suite.NoError(err)
	suite.NotNil(result)

	// Verify address list was created
	req := &tokenizationtypes.QueryGetAddressListRequest{
		ListId: "testlist",
	}
	resp, err := suite.TestSuite.Keeper.GetAddressList(suite.TestSuite.Ctx, req)
	suite.NoError(err)
	suite.NotNil(resp.List)
	suite.Equal("testlist", resp.List.ListId)
	suite.True(resp.List.Whitelist)
}
