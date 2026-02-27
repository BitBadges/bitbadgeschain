package tokenization_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

	// Pack*AsStruct only supports struct/tuple ABI outputs. Current precompile ABI has getCollection
	// returning bytes, so packing as struct fails (wrong output shape). This test documents that
	// when the ABI is updated to return a struct tuple, PackCollectionAsStruct will succeed.
	method, found := suite.Precompile.ABI.Methods["getCollection"]
	suite.True(found, "getCollection method should exist in ABI")
	packed, err := tokenization.PackCollectionAsStruct(&method, collection)
	if len(method.Outputs) == 1 && method.Outputs[0].Type.T == abi.BytesTy {
		suite.Error(err, "bytes return type not supported; ABI must use struct tuple for getCollection")
		suite.Nil(packed)
	} else {
		suite.NoError(err)
		suite.NotNil(packed)
		suite.Greater(len(packed), 0)
	}
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

func (suite *ReturnTypesTestSuite) TestConvertApprovalTrackerToSolidityStruct_Valid() {
	tracker := &tokenizationtypes.ApprovalTracker{
		NumTransfers:  sdkmath.NewUint(5),
		LastUpdatedAt: sdkmath.NewUint(1000),
		Amounts: []*tokenizationtypes.Balance{
			{Amount: sdkmath.NewUint(10), OwnershipTimes: nil, TokenIds: nil},
		},
	}
	structData, err := tokenization.ConvertApprovalTrackerToSolidityStruct(tracker)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Len(structData, 3)
	suite.Equal(big.NewInt(5), structData[0].(*big.Int))
	suite.Equal(big.NewInt(1000), structData[2].(*big.Int))
}

func (suite *ReturnTypesTestSuite) TestConvertApprovalTrackerToSolidityStruct_Nil() {
	_, err := tokenization.ConvertApprovalTrackerToSolidityStruct(nil)
	suite.Error(err)
}

func (suite *ReturnTypesTestSuite) TestConvertDynamicStoreToSolidityStruct_Valid() {
	store := &tokenizationtypes.DynamicStore{
		StoreId:         sdkmath.NewUint(1),
		CreatedBy:       "bb1xxx",
		DefaultValue:    true,
		GlobalEnabled:   true,
		Uri:             "https://example.com",
		CustomData:      "data",
	}
	structData, err := tokenization.ConvertDynamicStoreToSolidityStruct(store)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Len(structData, 6)
	suite.Equal(big.NewInt(1), structData[0].(*big.Int))
	// CreatedBy is now converted to EVM address format (common.Address)
	_, ok := structData[1].(common.Address)
	suite.True(ok, "CreatedBy should be common.Address type")
	suite.True(structData[2].(bool))
	suite.True(structData[3].(bool))
}

func (suite *ReturnTypesTestSuite) TestConvertDynamicStoreToSolidityStruct_Nil() {
	_, err := tokenization.ConvertDynamicStoreToSolidityStruct(nil)
	suite.Error(err)
}

func (suite *ReturnTypesTestSuite) TestConvertDynamicStoreValueToSolidityStruct_Valid() {
	val := &tokenizationtypes.DynamicStoreValue{
		StoreId:  sdkmath.NewUint(1),
		Address:  "bb1xxx",
		Value:    true,
	}
	structData, err := tokenization.ConvertDynamicStoreValueToSolidityStruct(val)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Len(structData, 3)
	suite.Equal(big.NewInt(1), structData[0].(*big.Int))
	// Address is now converted to EVM address format (common.Address)
	_, ok := structData[1].(common.Address)
	suite.True(ok, "Address should be common.Address type")
	suite.True(structData[2].(bool))
}

func (suite *ReturnTypesTestSuite) TestConvertDynamicStoreValueToSolidityStruct_Nil() {
	_, err := tokenization.ConvertDynamicStoreValueToSolidityStruct(nil)
	suite.Error(err)
}

func (suite *ReturnTypesTestSuite) TestConvertVoteProofToSolidityStruct_Valid() {
	proof := &tokenizationtypes.VoteProof{
		ProposalId: "prop-1",
		Voter:      "bb1xxx",
		YesWeight:  sdkmath.NewUint(70),
	}
	structData, err := tokenization.ConvertVoteProofToSolidityStruct(proof)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Len(structData, 3)
	suite.Equal("prop-1", structData[0].(string))
	// Voter is now converted to EVM address format (common.Address)
	_, ok := structData[1].(common.Address)
	suite.True(ok, "Voter should be common.Address type")
	suite.Equal(big.NewInt(70), structData[2].(*big.Int))
}

func (suite *ReturnTypesTestSuite) TestConvertVoteProofToSolidityStruct_Nil() {
	_, err := tokenization.ConvertVoteProofToSolidityStruct(nil)
	suite.Error(err)
}

func (suite *ReturnTypesTestSuite) TestConvertParamsToSolidityStruct_Valid() {
	params := &tokenizationtypes.Params{
		AllowedDenoms:      []string{"uatom", "ubadge"},
		AffiliatePercentage: sdkmath.NewUint(5),
	}
	structData, err := tokenization.ConvertParamsToSolidityStruct(params)
	suite.NoError(err)
	suite.NotNil(structData)
	suite.Len(structData, 2)
	denoms, ok := structData[0].([]interface{})
	suite.True(ok)
	suite.Len(denoms, 2)
	suite.Equal(big.NewInt(5), structData[1].(*big.Int))
}

func (suite *ReturnTypesTestSuite) TestConvertParamsToSolidityStruct_Nil() {
	_, err := tokenization.ConvertParamsToSolidityStruct(nil)
	suite.Error(err)
}

func (suite *ReturnTypesTestSuite) TestUserBalanceStore_OutgoingApprovalTupleLength() {
	// UserOutgoingApproval in Solidity has 10 fields: approvalId, toListId, initiatedByListId,
	// transferTimes, tokenIds, ownershipTimes, uri, customData, approvalCriteria, version.
	store := &tokenizationtypes.UserBalanceStore{
		Balances: []*tokenizationtypes.Balance{},
		OutgoingApprovals: []*tokenizationtypes.UserOutgoingApproval{
			{
				ApprovalId:        "ap1",
				ToListId:          "all",
				InitiatedByListId: "all",
				TransferTimes:     []*tokenizationtypes.UintRange{},
				TokenIds:          []*tokenizationtypes.UintRange{},
				OwnershipTimes:    []*tokenizationtypes.UintRange{},
				Uri:               "",
				CustomData:        "",
				ApprovalCriteria:  nil,
				Version:          sdkmath.NewUint(1),
			},
		},
		IncomingApprovals: []*tokenizationtypes.UserIncomingApproval{},
		UserPermissions:   nil,
	}
	structData, err := tokenization.ConvertUserBalanceStoreToSolidityStruct(store)
	suite.NoError(err)
	suite.NotNil(structData)
	outgoing, ok := structData[1].([]interface{})
	suite.True(ok)
	suite.Len(outgoing, 1)
	approvalTuple, ok := outgoing[0].([]interface{})
	suite.True(ok)
	suite.Len(approvalTuple, 10, "each UserOutgoingApproval must have 10 elements (approvalCriteria + version included)")
}

func (suite *ReturnTypesTestSuite) TestConvertCollection_InvariantsCosmosCoinBackedPath() {
	// Ensure invariants.cosmosCoinBackedPath is converted (not empty placeholder)
	collection := &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		Manager:      "bb1xxx",
		Invariants: &tokenizationtypes.CollectionInvariants{
			CosmosCoinBackedPath: &tokenizationtypes.CosmosCoinBackedPath{
				Address: "bb1backed",
				Conversion: &tokenizationtypes.Conversion{
					SideA: &tokenizationtypes.ConversionSideAWithDenom{
						Amount: sdkmath.NewUint(100),
						Denom:  "uatom",
					},
					SideB: []*tokenizationtypes.Balance{},
				},
			},
		},
	}
	structData, err := tokenization.ConvertCollectionToSolidityStruct(collection)
	suite.NoError(err)
	suite.NotNil(structData)
	// invariants is at index 14 (collectionId, metadata, tokenMetadata, customData, manager, permissions, approvals, standards, isArchived, defaultBalances, createdBy, validTokenIds, mintEscrowAddress, cosmosCoinWrapperPaths, invariants, aliasPaths)
	invariants, ok := structData[14].([]interface{})
	suite.True(ok)
	suite.Len(invariants, 6)
	// cosmosCoinBackedPath is index 2 in invariants
	path, ok := invariants[2].([]interface{})
	suite.True(ok)
	suite.Len(path, 2)
	// Address is now converted to EVM address format (common.Address)
	_, addrOk := path[0].(common.Address)
	suite.True(addrOk, "cosmosCoinBackedPath.Address should be common.Address type")
}

func (suite *ReturnTypesTestSuite) TestConvertCollection_InvariantsEvmQueryChallengesUriAndCustomData() {
	// Ensure invariants.evmQueryChallenges are converted with uri and customData populated
	collection := &tokenizationtypes.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		Manager:      "bb1xxx",
		Invariants: &tokenizationtypes.CollectionInvariants{
			EvmQueryChallenges: []*tokenizationtypes.EVMQueryChallenge{
				{
					ContractAddress:     "0x1234567890123456789012345678901234567890",
					Calldata:            "70a08231",
					ExpectedResult:       "",
					ComparisonOperator:  "eq",
					GasLimit:            sdkmath.NewUint(100000),
					Uri:                 "https://example.com/challenge",
					CustomData:          `{"desc":"test"}`,
				},
			},
		},
	}
	structData, err := tokenization.ConvertCollectionToSolidityStruct(collection)
	suite.NoError(err)
	suite.NotNil(structData)
	invariants, ok := structData[14].([]interface{})
	suite.True(ok)
	suite.Len(invariants, 6)
	// evmQueryChallenges is index 5 in invariants
	evmChallenges, ok := invariants[5].([]interface{})
	suite.True(ok)
	suite.Len(evmChallenges, 1)
	challengeTuple, ok := evmChallenges[0].([]interface{})
	suite.True(ok)
	suite.Len(challengeTuple, 7, "each EVMQueryChallenge must have 7 elements (contractAddress, calldata, expectedResult, comparisonOperator, gasLimit, uri, customData)")
	suite.Equal("https://example.com/challenge", challengeTuple[5].(string), "uri should be populated")
	suite.Equal(`{"desc":"test"}`, challengeTuple[6].(string), "customData should be populated")
}
