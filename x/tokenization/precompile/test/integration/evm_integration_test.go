package tokenization_test

import (
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"

	keepertest "github.com/bitbadges/bitbadgeschain/x/tokenization/testutil/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EVMIntegrationTestSuite is a test suite for real EVM integration tests
// Note: Full EVM integration requires setting up a complete app with EVM keeper
// This is a placeholder structure that can be expanded when full app setup is available
type EVMIntegrationTestSuite struct {
	suite.Suite
	TokenizationKeeper *tokenizationkeeper.Keeper
	Ctx                sdk.Context
	Precompile         *tokenization.Precompile

	// Test addresses
	AliceEVM common.Address
	BobEVM   common.Address
	Alice    sdk.AccAddress
	Bob      sdk.AccAddress

	// Test data
	CollectionId sdkmath.Uint
}

func TestEVMIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(EVMIntegrationTestSuite))
}

// SetupTest initializes the test suite
func (suite *EVMIntegrationTestSuite) SetupTest() {
	keeper, ctx := keepertest.TokenizationKeeper(suite.T())
	suite.TokenizationKeeper = keeper
	suite.Ctx = ctx
	suite.Precompile = tokenization.NewPrecompile(keeper)

	// Create test addresses
	suite.AliceEVM = common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0")
	suite.BobEVM = common.HexToAddress("0x8ba1f109551bD432803012645Hac136c22C9e7")

	suite.Alice = sdk.AccAddress(suite.AliceEVM.Bytes())
	suite.Bob = sdk.AccAddress(suite.BobEVM.Bytes())

	// Set up test collection
	suite.CollectionId = suite.createTestCollection()
}

// createTestCollection creates a test collection with balances
func (suite *EVMIntegrationTestSuite) createTestCollection() sdkmath.Uint {
	msgServer := tokenizationkeeper.NewMsgServerImpl(suite.TokenizationKeeper)

	// Create collection
	createMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:      suite.Alice.String(),
		CollectionId: sdkmath.NewUint(0),
		ValidTokenIds: []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		UpdateValidTokenIds:   true,
		CollectionPermissions: &tokenizationtypes.CollectionPermissions{},
	}

	resp, err := msgServer.UniversalUpdateCollection(suite.Ctx, createMsg)
	suite.Require().NoError(err)
	collectionId := resp.CollectionId

	// Set up mint approval
	getFullUintRanges := func() []*tokenizationtypes.UintRange {
		return []*tokenizationtypes.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		}
	}

	mintApproval := &tokenizationtypes.CollectionApproval{
		ApprovalId:        "mint_approval",
		FromListId:        tokenizationtypes.MintAddress,
		ToListId:          "AllWithoutMint",
		InitiatedByListId: "AllWithoutMint",
		TransferTimes:     getFullUintRanges(),
		TokenIds:          getFullUintRanges(),
		OwnershipTimes:    getFullUintRanges(),
		ApprovalCriteria: &tokenizationtypes.ApprovalCriteria{
			MaxNumTransfers: &tokenizationtypes.MaxNumTransfers{
				OverallMaxNumTransfers: sdkmath.NewUint(1000),
				AmountTrackerId:        "mint-tracker",
			},
			ApprovalAmounts: &tokenizationtypes.ApprovalAmounts{
				PerFromAddressApprovalAmount: sdkmath.NewUint(1000),
				AmountTrackerId:              "mint-tracker",
			},
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		Version: sdkmath.NewUint(0),
	}

	updateApprovalsMsg := &tokenizationtypes.MsgUniversalUpdateCollection{
		Creator:                   suite.Alice.String(),
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*tokenizationtypes.CollectionApproval{mintApproval},
	}
	_, err = msgServer.UniversalUpdateCollection(suite.Ctx, updateApprovalsMsg)
	suite.Require().NoError(err)

	// Mint tokens to Alice using TransferTokens from Mint address
	// Note: This is a simplified test setup - full minting would use MsgMintAndDistributeTokens
	// For this test, we'll just verify the precompile structure is correct

	return collectionId
}

// TestPrecompileRegistration tests that the precompile is properly registered
func (suite *EVMIntegrationTestSuite) TestPrecompileRegistration() {
	// Verify precompile address
	suite.Equal(tokenization.TokenizationPrecompileAddress, suite.Precompile.ContractAddress.Hex())

	// Verify ABI is loaded
	suite.NotNil(suite.Precompile.ABI)

	// Verify key methods exist
	methods := []string{
		"transferTokens",
		"setIncomingApproval",
		"setOutgoingApproval",
		"getCollection",
		"getBalance",
		"getBalanceAmount",
		"getTotalSupply",
	}

	for _, methodName := range methods {
		method, found := suite.Precompile.ABI.Methods[methodName]
		suite.True(found, "method %s should exist", methodName)
		suite.NotNil(method, "method %s should not be nil", methodName)
	}
}

// TestPrecompileRequiredGas tests gas calculation
func (suite *EVMIntegrationTestSuite) TestPrecompileRequiredGas() {
	method := suite.Precompile.ABI.Methods["transferTokens"]
	suite.Require().NotNil(method)

	// Test with valid method ID
	// Transaction methods add a 200k buffer to base gas for Cosmos SDK operations
	const txBuffer = 200_000
	gas := suite.Precompile.RequiredGas(method.ID[:])
	suite.Equal(uint64(tokenization.GasTransferTokensBase+txBuffer), gas)

	// Test with invalid input (too short)
	gas = suite.Precompile.RequiredGas([]byte{0x12, 0x34})
	suite.Equal(uint64(0), gas)
}

// TestPrecompileExecuteDirectly tests direct execution (bypassing EVM)
// This is a simplified test - full EVM integration would require app setup
func (suite *EVMIntegrationTestSuite) TestPrecompileExecuteDirectly() {
	method := suite.Precompile.ABI.Methods["getCollection"]
	suite.Require().NotNil(method)

	// Build JSON query
	queryJson, err := helpers.BuildGetCollectionQueryJSON(suite.CollectionId.BigInt())
	suite.Require().NoError(err)

	// Pack method with JSON string
	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.Require().NoError(err)

	// Note: Full EVM integration would require:
	// 1. Setting up app with EVM keeper
	// 2. Creating EVM instance
	// 3. Creating contract context
	// 4. Calling Run() method which internally calls Execute()
	// For now, we verify the precompile structure is correct
	suite.NotNil(input)
	suite.Greater(len(input), 4) // Should have method ID + packed args
}

// NOTE: Full EVM integration tests require:
// 1. Complete app setup with EVM keeper (see app/test_helpers.go)
// 2. EVM instance creation
// 3. Contract deployment and interaction
// 4. Event verification through EVM
// These tests are placeholders that can be expanded when full app integration is available
