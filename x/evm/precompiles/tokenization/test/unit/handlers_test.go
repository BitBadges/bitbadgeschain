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
	contract := suite.createContract(caller)

	validTokenIds := []interface{}{
		map[string]interface{}{
			"start": big.NewInt(1),
			"end":   big.NewInt(100),
		},
	}

	args := []interface{}{
		nil, // defaultBalances
		validTokenIds,
		map[string]interface{}{},         // collectionPermissions
		suite.TestSuite.Manager.String(), // manager
		map[string]interface{}{"uri": "https://example.com", "customData": "test"}, // collectionMetadata
		[]interface{}{},          // tokenMetadata
		"custom data",            // customData
		[]interface{}{},          // collectionApprovals
		[]string{"ERC721"},       // standards
		false,                    // isArchived
		[]interface{}{},          // mintEscrowCoinsToTransfer
		[]interface{}{},          // cosmosCoinWrapperPathsToAdd
		map[string]interface{}{}, // invariants
		[]interface{}{},          // aliasPathsToAdd
	}

	method := suite.Precompile.ABI.Methods["createCollection"]
	result, err := suite.Precompile.CreateCollection(suite.TestSuite.Ctx, &method, args, contract)
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

func (suite *HandlersTestSuite) TestDeleteCollection_Valid() {
	// Create collection first
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	args := []interface{}{
		collectionId.BigInt(),
	}

	method := suite.Precompile.ABI.Methods["deleteCollection"]
	result, err := suite.Precompile.DeleteCollection(suite.TestSuite.Ctx, &method, args, contract)
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
	contract := suite.createContract(caller)

	args := []interface{}{
		collectionId.BigInt(),
	}

	method := suite.Precompile.ABI.Methods["deleteCollection"]
	result, err := suite.Precompile.DeleteCollection(suite.TestSuite.Ctx, &method, args, contract)
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
	contract := suite.createContract(caller)

	args := []interface{}{
		collectionId.BigInt(),
		[]common.Address{suite.TestSuite.BobEVM},
		big.NewInt(100),
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
	contract := suite.createContract(caller)

	// Try to transfer 100 tokens (more than available)
	args := []interface{}{
		collectionId.BigInt(),
		[]common.Address{suite.TestSuite.BobEVM},
		big.NewInt(100), // More than balance
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
}

func (suite *HandlersTestSuite) TestSetManager_Valid() {
	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.NoError(err)

	caller := suite.TestSuite.AliceEVM
	contract := suite.createContract(caller)

	args := []interface{}{
		collectionId.BigInt(),
		suite.TestSuite.Bob.String(),
	}

	method := suite.Precompile.ABI.Methods["setManager"]
	result, err := suite.Precompile.SetManager(suite.TestSuite.Ctx, &method, args, contract)
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
	contract := suite.createContract(caller)

	args := []interface{}{
		collectionId.BigInt(),
		"https://new-uri.com",
		"new custom data",
	}

	method := suite.Precompile.ABI.Methods["setCollectionMetadata"]
	result, err := suite.Precompile.SetCollectionMetadata(suite.TestSuite.Ctx, &method, args, contract)
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
	contract := suite.createContract(caller)

	newTokenIds := []interface{}{
		map[string]interface{}{
			"start": big.NewInt(1),
			"end":   big.NewInt(200),
		},
	}

	args := []interface{}{
		collectionId.BigInt(),
		newTokenIds,
		[]interface{}{}, // canUpdateValidTokenIds
	}

	method := suite.Precompile.ABI.Methods["setValidTokenIds"]
	result, err := suite.Precompile.SetValidTokenIds(suite.TestSuite.Ctx, &method, args, contract)
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
	contract := suite.createContract(caller)

	args := []interface{}{
		false,                 // defaultValue
		"https://example.com", // uri
		"custom data",         // customData
	}

	method := suite.Precompile.ABI.Methods["createDynamicStore"]
	result, err := suite.Precompile.CreateDynamicStore(suite.TestSuite.Ctx, &method, args, contract)
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
	contract := suite.createContract(caller)

	addressListInput := map[string]interface{}{
		"listId":     "testlist", // Use alphanumeric only (no hyphens)
		"addresses":  []string{suite.TestSuite.Bob.String(), suite.TestSuite.Charlie.String()},
		"whitelist":  true,
		"uri":        "https://example.com",
		"customData": "data",
	}

	args := []interface{}{
		[]interface{}{addressListInput},
	}

	method, found := suite.Precompile.ABI.Methods["createAddressLists"]
	if !found {
		// Create mock method if not in ABI (workaround for missing ABI entries)
		method = helpers.CreateMockMethod("createAddressLists", nil, helpers.CreateMockBoolOutput())
	}
	result, err := suite.Precompile.CreateAddressLists(suite.TestSuite.Ctx, &method, args, contract)
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
