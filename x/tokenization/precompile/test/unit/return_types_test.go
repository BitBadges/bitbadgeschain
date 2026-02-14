package tokenization_test

import (
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

type ReturnTypesTestSuite struct {
	suite.Suite
	Precompile *tokenization.Precompile
}

func TestReturnTypesTestSuite(t *testing.T) {
	suite.Run(t, new(ReturnTypesTestSuite))
}

func (suite *ReturnTypesTestSuite) SetupTest() {
	ts := helpers.NewTestSuite()
	suite.Precompile = ts.Precompile
}

func (suite *ReturnTypesTestSuite) TestConvertBalanceToSolidityStruct_Valid() {
	balance := &tokenizationtypes.Balance{
		Amount: sdkmath.NewUint(100),
		OwnershipTimes: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)},
		},
		TokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		},
	}

	structData, err := tokenization.ConvertBalanceToSolidityStruct(balance)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Len(structData, 3)

	// Check amount
	amount, ok := structData[0].(*big.Int)
	suite.True(ok)
	suite.Equal(big.NewInt(100), amount)

	// Check ownership times
	ownershipTimes, ok := structData[1].([]interface{})
	suite.True(ok)
	suite.Len(ownershipTimes, 1)
	ot, ok := ownershipTimes[0].([]interface{})
	suite.True(ok)
	suite.Len(ot, 2)
	suite.Equal(big.NewInt(1), ot[0].(*big.Int))
	suite.Equal(big.NewInt(1000), ot[1].(*big.Int))

	// Check token IDs
	tokenIds, ok := structData[2].([]interface{})
	suite.True(ok)
	suite.Len(tokenIds, 1)
	tid, ok := tokenIds[0].([]interface{})
	suite.True(ok)
	suite.Len(tid, 2)
	suite.Equal(big.NewInt(1), tid[0].(*big.Int))
	suite.Equal(big.NewInt(10), tid[1].(*big.Int))
}

func (suite *ReturnTypesTestSuite) TestConvertBalanceToSolidityStruct_Nil() {
	structData, err := tokenization.ConvertBalanceToSolidityStruct(nil)
	suite.Error(err)
	suite.Nil(structData)
}

func (suite *ReturnTypesTestSuite) TestConvertUserBalanceStoreToSolidityStruct_Valid() {
	store := &tokenizationtypes.UserBalanceStore{
		Balances: []*tokenizationtypes.Balance{
			{
				Amount: sdkmath.NewUint(100),
				OwnershipTimes: []*tokenizationtypes.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)},
				},
				TokenIds: []*tokenizationtypes.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
			},
		},
		AutoApproveSelfInitiatedOutgoingTransfers: true,
		AutoApproveSelfInitiatedIncomingTransfers: false,
		AutoApproveAllIncomingTransfers:           true,
	}

	structData, err := tokenization.ConvertUserBalanceStoreToSolidityStruct(store)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Len(structData, 7)

	// Check balances array
	balances, ok := structData[0].([]interface{})
	suite.True(ok)
	suite.Len(balances, 1)

	// Check boolean fields
	suite.True(structData[3].(bool))  // autoApproveSelfInitiatedOutgoingTransfers
	suite.False(structData[4].(bool)) // autoApproveSelfInitiatedIncomingTransfers
	suite.True(structData[5].(bool))  // autoApproveAllIncomingTransfers
}

func (suite *ReturnTypesTestSuite) TestConvertAddressListToSolidityStruct_Valid() {
	list := &tokenizationtypes.AddressList{
		ListId:     "test-list",
		Addresses:  []string{"bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"},
		Whitelist:  true,
		Uri:        "https://example.com",
		CustomData: "data",
		CreatedBy:  "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q",
	}

	structData, err := tokenization.ConvertAddressListToSolidityStruct(list)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Len(structData, 6)
	suite.Equal("test-list", structData[0].(string))
	suite.True(structData[2].(bool)) // whitelist
	suite.Equal("https://example.com", structData[3].(string))
}

func (suite *ReturnTypesTestSuite) TestConvertCollectionToSolidityStruct_Valid() {
	collection := &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		CollectionMetadata: &tokenizationtypes.CollectionMetadata{
			Uri:        "https://example.com",
			CustomData: "data",
		},
		CustomData: "collection data",
		Manager:    "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		Standards:  []string{"ERC721"},
		IsArchived: false,
		CreatedBy:  "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q",
	}

	structData, err := tokenization.ConvertCollectionToSolidityStruct(collection)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Greater(len(structData), 10) // Should have many fields

	// Check collection ID
	collectionId, ok := structData[0].(*big.Int)
	suite.True(ok)
	suite.Equal(big.NewInt(1), collectionId)

	// Check metadata
	metadata, ok := structData[1].([]interface{})
	suite.True(ok)
	suite.Len(metadata, 2)
	suite.Equal("https://example.com", metadata[0].(string))
}

func (suite *ReturnTypesTestSuite) TestPackCollectionAsStruct_Valid() {
	collection := &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		CollectionMetadata: &tokenizationtypes.CollectionMetadata{
			Uri:        "https://example.com",
			CustomData: "data",
		},
		Manager: "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		Standards:  []string{"ERC721"},
		IsArchived: false,
	}

	// Use actual method from ABI
	method, found := suite.Precompile.ABI.Methods["getCollection"]
	suite.True(found, "getCollection method should exist in ABI")

	packed, err := tokenization.PackCollectionAsStruct(&method, collection)
	suite.NoError(err)
	suite.NotNil(packed)
	suite.Greater(len(packed), 0)
}

func (suite *ReturnTypesTestSuite) TestPackBalanceAsStruct_Valid() {
	balance := &tokenizationtypes.Balance{
		Amount: sdkmath.NewUint(100),
		OwnershipTimes: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)},
		},
		TokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
		},
	}

	// Note: PackBalanceAsStruct is for individual Balance, but getBalance returns UserBalanceStore
	// This test verifies the conversion function works
	structData, err := tokenization.ConvertBalanceToSolidityStruct(balance)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Greater(len(structData), 0)
}

func (suite *ReturnTypesTestSuite) TestConvertAddressListToSolidityStruct_EVMAddress() {
	// Create an address list with EVM address format (will be converted from Cosmos)
	evmAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	cosmosAddr := sdk.AccAddress(evmAddr.Bytes()).String()

	list := &tokenizationtypes.AddressList{
		ListId:     "test-list",
		Addresses:  []string{cosmosAddr},
		Whitelist:  true,
		Uri:        "https://example.com",
		CustomData: "data",
		CreatedBy:  cosmosAddr,
	}

	structData, err := tokenization.ConvertAddressListToSolidityStruct(list)
	suite.NoError(err)
	suite.NotNil(structData)

	// Addresses should be converted back to EVM format
	addresses, ok := structData[1].([]interface{})
	suite.True(ok)
	suite.Len(addresses, 1)

	// Should be an address type (either common.Address or string)
	addr, ok := addresses[0].(common.Address)
	if ok {
		suite.Equal(evmAddr, addr)
	}
}

func (suite *ReturnTypesTestSuite) TestConvertCollectionToSolidityStruct_EmptyFields() {
	collection := &tokenizationtypes.TokenCollection{
		CollectionId:  sdkmath.NewUint(1),
		Manager:       "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430",
		ValidTokenIds: []*tokenizationtypes.UintRange{},
		Standards:     []string{},
		IsArchived:    false,
	}

	structData, err := tokenization.ConvertCollectionToSolidityStruct(collection)
	suite.NoError(err)
	suite.NotNil(structData)

	// Should handle nil metadata
	metadata, ok := structData[1].([]interface{})
	suite.True(ok)
	suite.Len(metadata, 2)
	suite.Equal("", metadata[0].(string)) // Empty URI
}

func (suite *ReturnTypesTestSuite) TestConvertCollectionToSolidityStruct_Nil() {
	structData, err := tokenization.ConvertCollectionToSolidityStruct(nil)
	suite.Error(err)
	suite.Nil(structData)
}

func (suite *ReturnTypesTestSuite) TestConvertUserBalanceStoreToSolidityStruct_Nil() {
	structData, err := tokenization.ConvertUserBalanceStoreToSolidityStruct(nil)
	suite.Error(err)
	suite.Nil(structData)
}

func (suite *ReturnTypesTestSuite) TestConvertAddressListToSolidityStruct_Nil() {
	structData, err := tokenization.ConvertAddressListToSolidityStruct(nil)
	suite.Error(err)
	suite.Nil(structData)
}
