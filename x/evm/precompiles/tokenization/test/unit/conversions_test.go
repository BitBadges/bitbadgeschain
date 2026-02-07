package tokenization_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokenization "github.com/bitbadges/bitbadgeschain/x/evm/precompiles/tokenization"
)

type ConversionsTestSuite struct {
	suite.Suite
}

func TestConversionsTestSuite(t *testing.T) {
	suite.Run(t, new(ConversionsTestSuite))
}

func (suite *ConversionsTestSuite) TestConvertBalance_ValidInput() {
	balanceMap := map[string]interface{}{
		"amount": big.NewInt(100),
		"ownershipTimes": []interface{}{
			map[string]interface{}{
				"start": big.NewInt(1),
				"end":   big.NewInt(1000),
			},
		},
		"tokenIds": []interface{}{
			map[string]interface{}{
				"start": big.NewInt(1),
				"end":   big.NewInt(10),
			},
		},
	}

	balance, err := tokenization.ConvertBalance(balanceMap)
	suite.NoError(err)
	suite.NotNil(balance)
	suite.Equal(sdkmath.NewUint(100), balance.Amount)
	suite.Len(balance.OwnershipTimes, 1)
	suite.Len(balance.TokenIds, 1)
	suite.Equal(sdkmath.NewUint(1), balance.OwnershipTimes[0].Start)
	suite.Equal(sdkmath.NewUint(1000), balance.OwnershipTimes[0].End)
	suite.Equal(sdkmath.NewUint(1), balance.TokenIds[0].Start)
	suite.Equal(sdkmath.NewUint(10), balance.TokenIds[0].End)
}

func (suite *ConversionsTestSuite) TestConvertBalance_InvalidAmount() {
	balanceMap := map[string]interface{}{
		"amount":         "invalid",
		"ownershipTimes": []interface{}{},
		"tokenIds":       []interface{}{},
	}

	balance, err := tokenization.ConvertBalance(balanceMap)
	suite.Error(err)
	suite.Nil(balance)
}

func (suite *ConversionsTestSuite) TestConvertBalance_Overflow() {
	balanceMap := map[string]interface{}{
		"amount":         new(big.Int).Lsh(big.NewInt(1), 256), // 2^256 - overflow
		"ownershipTimes": []interface{}{},
		"tokenIds":       []interface{}{},
	}

	balance, err := tokenization.ConvertBalance(balanceMap)
	suite.Error(err)
	suite.Nil(balance)
}

func (suite *ConversionsTestSuite) TestConvertBalance_EmptyArrays() {
	// Empty arrays are not allowed by validation - ranges cannot be empty
	balanceMap := map[string]interface{}{
		"amount":         big.NewInt(100),
		"ownershipTimes": []interface{}{},
		"tokenIds":       []interface{}{},
	}

	balance, err := tokenization.ConvertBalance(balanceMap)
	suite.Error(err)
	suite.Nil(balance)
	suite.Contains(err.Error(), "cannot be empty")
}

func (suite *ConversionsTestSuite) TestConvertUintRangeArray_Valid() {
	ranges := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: big.NewInt(100)},
		{Start: big.NewInt(101), End: big.NewInt(200)},
	}

	result, err := tokenization.ConvertUintRangeArray(ranges)
	suite.NoError(err)
	suite.Len(result, 2)
	suite.Equal(sdkmath.NewUint(1), result[0].Start)
	suite.Equal(sdkmath.NewUint(100), result[0].End)
	suite.Equal(sdkmath.NewUint(101), result[1].Start)
	suite.Equal(sdkmath.NewUint(200), result[1].End)
}

func (suite *ConversionsTestSuite) TestConvertUintRangeArray_InvalidRange() {
	ranges := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(100), End: big.NewInt(1)}, // start > end
	}

	result, err := tokenization.ConvertUintRangeArray(ranges)
	suite.Error(err)
	suite.Nil(result)
}

func (suite *ConversionsTestSuite) TestConvertCollectionMetadata() {
	metadata := tokenization.ConvertCollectionMetadata("https://example.com/uri", "custom data")
	suite.NotNil(metadata)
	suite.Equal("https://example.com/uri", metadata.Uri)
	suite.Equal("custom data", metadata.CustomData)
}

func (suite *ConversionsTestSuite) TestConvertTokenMetadata_Valid() {
	tokenIds := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: big.NewInt(10)},
	}

	metadata, err := tokenization.ConvertTokenMetadata("https://example.com/token", "data", tokenIds)
	suite.NoError(err)
	suite.NotNil(metadata)
	suite.Equal("https://example.com/token", metadata.Uri)
	suite.Equal("data", metadata.CustomData)
	suite.Len(metadata.TokenIds, 1)
}

func (suite *ConversionsTestSuite) TestConvertActionPermission_Valid() {
	permittedTimes := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(1), End: big.NewInt(1000)},
	}

	forbiddenTimes := []struct {
		Start *big.Int `json:"start"`
		End   *big.Int `json:"end"`
	}{
		{Start: big.NewInt(2000), End: big.NewInt(3000)},
	}

	permission, err := tokenization.ConvertActionPermission(permittedTimes, forbiddenTimes)
	suite.NoError(err)
	suite.NotNil(permission)
	suite.Len(permission.PermanentlyPermittedTimes, 1)
	suite.Len(permission.PermanentlyForbiddenTimes, 1)
}

func (suite *ConversionsTestSuite) TestConvertManagerAddress_EVMAddress() {
	evmAddr := common.HexToAddress("0x1111111111111111111111111111111111111111")
	manager, err := tokenization.ConvertManagerAddress(evmAddr)
	suite.NoError(err)
	suite.NotEmpty(manager)
	// Should be a valid Cosmos address
	_, err = sdk.AccAddressFromBech32(manager)
	suite.NoError(err)
}

func (suite *ConversionsTestSuite) TestConvertManagerAddress_String() {
	manager, err := tokenization.ConvertManagerAddress("bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430")
	suite.NoError(err)
	suite.Equal("bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430", manager)
}

func (suite *ConversionsTestSuite) TestConvertManagerAddress_InvalidType() {
	manager, err := tokenization.ConvertManagerAddress(123)
	suite.Error(err)
	suite.Empty(manager)
}

func (suite *ConversionsTestSuite) TestConvertInvariantsAddObject_Empty() {
	invariants, err := tokenization.ConvertInvariantsAddObject(nil)
	suite.NoError(err)
	suite.NotNil(invariants)
	suite.False(invariants.NoCustomOwnershipTimes)
	suite.True(invariants.MaxSupplyPerId.IsZero())
}

func (suite *ConversionsTestSuite) TestConvertInvariantsAddObject_WithFields() {
	invariantsMap := map[string]interface{}{
		"noCustomOwnershipTimes":      true,
		"maxSupplyPerId":              big.NewInt(1000),
		"noForcefulPostMintTransfers": true,
		"disablePoolCreation":         true,
	}

	invariants, err := tokenization.ConvertInvariantsAddObject(invariantsMap)
	suite.NoError(err)
	suite.NotNil(invariants)
	suite.True(invariants.NoCustomOwnershipTimes)
	suite.Equal(sdkmath.NewUint(1000), invariants.MaxSupplyPerId)
	suite.True(invariants.NoForcefulPostMintTransfers)
	suite.True(invariants.DisablePoolCreation)
}

func (suite *ConversionsTestSuite) TestConvertApprovalCriteria_BooleanFields() {
	criteriaMap := map[string]interface{}{
		"requireToEqualsInitiatedBy":       true,
		"requireFromEqualsInitiatedBy":     false,
		"requireToDoesNotEqualInitiatedBy": true,
		"overridesFromOutgoingApprovals":   false,
		"mustPrioritize":                   true,
		"allowBackedMinting":               false,
		"allowSpecialWrapping":             true,
	}

	// ConvertApprovalCriteria returns (*ApprovalCriteria, error)
	criteria, err := tokenization.ConvertApprovalCriteria(criteriaMap)
	suite.NoError(err)
	suite.NotNil(criteria)
	suite.True(criteria.RequireToEqualsInitiatedBy)
	suite.False(criteria.RequireFromEqualsInitiatedBy)
	suite.True(criteria.RequireToDoesNotEqualInitiatedBy)
	suite.False(criteria.OverridesFromOutgoingApprovals)
	suite.True(criteria.MustPrioritize)
	suite.False(criteria.AllowBackedMinting)
	suite.True(criteria.AllowSpecialWrapping)
}

func (suite *ConversionsTestSuite) TestConvertUserOutgoingApproval_Valid() {
	appMap := map[string]interface{}{
		"approvalId":        "test-approval",
		"toListId":          "test-list",
		"initiatedByListId": "initiator-list",
		"uri":               "https://example.com",
		"customData":        "data",
		"transferTimes": []interface{}{
			map[string]interface{}{
				"start": big.NewInt(1),
				"end":   big.NewInt(1000),
			},
		},
		"tokenIds": []interface{}{
			map[string]interface{}{
				"start": big.NewInt(1),
				"end":   big.NewInt(10),
			},
		},
		"ownershipTimes": []interface{}{
			map[string]interface{}{
				"start": big.NewInt(1),
				"end":   big.NewInt(1000),
			},
		},
	}

	approval, err := tokenization.ConvertUserOutgoingApproval(appMap)
	suite.NoError(err)
	suite.NotNil(approval)
	suite.Equal("test-approval", approval.ApprovalId)
	suite.Equal("test-list", approval.ToListId)
	suite.NotNil(approval.ApprovalCriteria)
}

func (suite *ConversionsTestSuite) TestConvertUserIncomingApproval_Valid() {
	appMap := map[string]interface{}{
		"approvalId":        "test-approval",
		"fromListId":        "test-list",
		"initiatedByListId": "initiator-list",
		"uri":               "https://example.com",
		"customData":        "data",
		"transferTimes": []interface{}{
			map[string]interface{}{
				"start": big.NewInt(1),
				"end":   big.NewInt(1000),
			},
		},
	}

	approval, err := tokenization.ConvertUserIncomingApproval(appMap)
	suite.NoError(err)
	suite.NotNil(approval)
	suite.Equal("test-approval", approval.ApprovalId)
	suite.Equal("test-list", approval.FromListId)
	suite.NotNil(approval.ApprovalCriteria)
}

func (suite *ConversionsTestSuite) TestConvertBalanceArray_Valid() {
	balancesRaw := []interface{}{
		map[string]interface{}{
			"amount": big.NewInt(100),
			"ownershipTimes": []interface{}{
				map[string]interface{}{
					"start": big.NewInt(1),
					"end":   big.NewInt(1000),
				},
			},
			"tokenIds": []interface{}{
				map[string]interface{}{
					"start": big.NewInt(1),
					"end":   big.NewInt(10),
				},
			},
		},
		map[string]interface{}{
			"amount": big.NewInt(200),
			"ownershipTimes": []interface{}{
				map[string]interface{}{
					"start": big.NewInt(1),
					"end":   big.NewInt(1000),
				},
			},
			"tokenIds": []interface{}{
				map[string]interface{}{
					"start": big.NewInt(11),
					"end":   big.NewInt(20),
				},
			},
		},
	}

	balances, err := tokenization.ConvertBalanceArray(balancesRaw)
	suite.NoError(err)
	suite.Len(balances, 2)
	suite.Equal(sdkmath.NewUint(100), balances[0].Amount)
	suite.Equal(sdkmath.NewUint(200), balances[1].Amount)
}

func (suite *ConversionsTestSuite) TestConvertBalanceArray_InvalidEntry() {
	balancesRaw := []interface{}{
		"not a map",
	}

	balances, err := tokenization.ConvertBalanceArray(balancesRaw)
	suite.Error(err)
	suite.Nil(balances)
}

func (suite *ConversionsTestSuite) TestConvertCollectionPermissions_Empty() {
	permsMap := map[string]interface{}{}

	perms, err := tokenization.ConvertCollectionPermissions(permsMap)
	suite.NoError(err)
	suite.NotNil(perms)
}

func (suite *ConversionsTestSuite) TestConvertCollectionPermissions_WithFields() {
	permsMap := map[string]interface{}{
		"canDeleteCollection": []interface{}{
			map[string]interface{}{
				"permanentlyPermittedTimes": []interface{}{},
				"permanentlyForbiddenTimes": []interface{}{},
			},
		},
	}

	perms, err := tokenization.ConvertCollectionPermissions(permsMap)
	suite.NoError(err)
	suite.NotNil(perms)
	// Should have at least one permission if conversion succeeded
}

func (suite *ConversionsTestSuite) TestConvertAddressListInput_Valid() {
	addresses := []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
	}

	input := tokenization.ConvertAddressListInput("test-list", addresses, true, "https://example.com", "data")
	suite.NotNil(input)
	suite.Equal("test-list", input.ListId)
	suite.True(input.Whitelist)
	suite.Equal("https://example.com", input.Uri)
	suite.Equal("data", input.CustomData)
	suite.Len(input.Addresses, 2)
	// Addresses should be converted to Cosmos format
	for _, addr := range input.Addresses {
		_, err := sdk.AccAddressFromBech32(addr)
		suite.NoError(err, "Address should be valid Cosmos address: %s", addr)
	}
}

func (suite *ConversionsTestSuite) TestConvertCosmosCoinWrapperPathAddObjectArray_Empty() {
	paths, err := tokenization.ConvertCosmosCoinWrapperPathAddObjectArray([]interface{}{})
	suite.NoError(err)
	suite.NotNil(paths)
	suite.Len(paths, 0)
}

func (suite *ConversionsTestSuite) TestConvertCosmosCoinWrapperPathAddObjectArray_Valid() {
	pathsRaw := []interface{}{
		map[string]interface{}{
			"denom":                          "test-denom",
			"symbol":                         "TEST",
			"allowOverrideWithAnyValidToken": true,
			"denomUnits":                     []interface{}{},
			"metadata": map[string]interface{}{
				"uri":        "https://example.com",
				"customData": "data",
			},
		},
	}

	paths, err := tokenization.ConvertCosmosCoinWrapperPathAddObjectArray(pathsRaw)
	suite.NoError(err)
	suite.Len(paths, 1)
	suite.Equal("test-denom", paths[0].Denom)
	suite.Equal("TEST", paths[0].Symbol)
	suite.True(paths[0].AllowOverrideWithAnyValidToken)
}

func (suite *ConversionsTestSuite) TestConvertAliasPathAddObjectArray_Valid() {
	pathsRaw := []interface{}{
		map[string]interface{}{
			"denom":      "test-denom",
			"symbol":     "TEST",
			"denomUnits": []interface{}{},
			"metadata": map[string]interface{}{
				"uri":        "https://example.com",
				"customData": "data",
			},
		},
	}

	paths, err := tokenization.ConvertAliasPathAddObjectArray(pathsRaw)
	suite.NoError(err)
	suite.Len(paths, 1)
	suite.Equal("test-denom", paths[0].Denom)
	suite.Equal("TEST", paths[0].Symbol)
}
