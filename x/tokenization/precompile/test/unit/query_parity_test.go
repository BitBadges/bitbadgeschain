package tokenization_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/gogoproto/proto"

	tokenization "github.com/bitbadges/bitbadgeschain/x/tokenization/precompile"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/precompile/test/helpers"
	tokenizationtypes "github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type QueryParityTestSuite struct {
	suite.Suite
	TestSuite  *helpers.TestSuite
	Precompile *tokenization.Precompile
}

func TestQueryParityTestSuite(t *testing.T) {
	suite.Run(t, new(QueryParityTestSuite))
}

func (suite *QueryParityTestSuite) SetupTest() {
	suite.TestSuite = helpers.NewTestSuite()
	suite.TestSuite.Ctx = suite.TestSuite.Ctx.
		WithBlockHeight(4242).
		WithBlockTime(time.UnixMilli(1735689600123).UTC())
	suite.Precompile = suite.TestSuite.Precompile

	collectionId, err := suite.TestSuite.CreateTestCollection(suite.TestSuite.Alice.String())
	suite.Require().NoError(err)

	// Attach wrapper paths (exact + placeholder) and an alias path to the collection
	collection, found := suite.TestSuite.Keeper.GetCollectionFromStore(suite.TestSuite.Ctx, collectionId)
	suite.Require().True(found)
	collection.CosmosCoinWrapperPaths = []*tokenizationtypes.CosmosCoinWrapperPath{
		{
			Denom:  "exact-wrapper",
			Symbol: "EXW",
			Conversion: &tokenizationtypes.ConversionWithoutDenom{
				SideA: &tokenizationtypes.ConversionSideA{Amount: sdkmath.NewUint(1)},
				SideB: []*tokenizationtypes.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2000000000000)}},
					},
				},
			},
		},
		{
			Denom:                          "item-{id}",
			Symbol:                         "ITM",
			AllowOverrideWithAnyValidToken: true,
			Conversion: &tokenizationtypes.ConversionWithoutDenom{
				SideA: &tokenizationtypes.ConversionSideA{Amount: sdkmath.NewUint(1)},
				SideB: []*tokenizationtypes.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2000000000000)}},
					},
				},
			},
		},
	}
	collection.AliasPaths = []*tokenizationtypes.AliasPath{
		{
			Denom:  "pool-1",
			Symbol: "LP1",
			Conversion: &tokenizationtypes.ConversionWithoutDenom{
				SideA: &tokenizationtypes.ConversionSideA{Amount: sdkmath.NewUint(1)},
				SideB: []*tokenizationtypes.Balance{
					{
						Amount:         sdkmath.NewUint(1),
						TokenIds:       []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(2)}},
						OwnershipTimes: []*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(2000000000000)}},
					},
				},
			},
		},
	}
	suite.Require().NoError(suite.TestSuite.Keeper.SetCollectionInStore(suite.TestSuite.Ctx, collection, false))

	// Give Alice balance for token 1..100 active during [1000, 1735689600999]
	err = suite.TestSuite.CreateTestBalance(
		collectionId,
		suite.TestSuite.Alice.String(),
		sdkmath.NewUint(750),
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)}},
		[]*tokenizationtypes.UintRange{{Start: sdkmath.NewUint(1000), End: sdkmath.NewUint(1735689600999)}},
	)
	suite.Require().NoError(err)
}

func (suite *QueryParityTestSuite) TestNewQueryMethods_MetadataAndGas() {
	for _, name := range []string{"resolveDenom", "getBalanceForToken"} {
		method, found := suite.Precompile.ABI.Methods[name]
		suite.Require().True(found, "method %s must exist in ABI", name)
		suite.False(suite.Precompile.IsTransaction(&method), "method %s must be read-only", name)
		gas := suite.Precompile.RequiredGas(method.ID)
		suite.Greater(gas, uint64(0), "method %s must have non-zero RequiredGas", name)
	}
}

func (suite *QueryParityTestSuite) TestGetCollection_ReturnsImmutableCreatedMetadata() {
	method := suite.Precompile.ABI.Methods["getCollection"]
	queryJson, err := helpers.BuildGetCollectionQueryJSON(suite.TestSuite.CollectionId.BigInt())
	suite.Require().NoError(err)

	input, err := helpers.PackMethodWithJSON(&method, queryJson)
	suite.Require().NoError(err)

	result, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, input)
	suite.Require().NoError(err)

	unpacked, err := method.Outputs.Unpack(result)
	suite.Require().NoError(err)
	suite.Require().Len(unpacked, 1)

	respBytes, ok := unpacked[0].([]byte)
	suite.Require().True(ok)

	var resp tokenizationtypes.QueryGetCollectionResponse
	suite.Require().NoError(proto.Unmarshal(respBytes, &resp))
	suite.Require().NotNil(resp.Collection)
	suite.Equal(uint64(4242), resp.Collection.CreatedBlock)
	suite.Equal(uint64(1735689600123), resp.Collection.CreatedTimestamp)
}

func (suite *QueryParityTestSuite) TestResolveDenom_ExactPlaceholderAndAlias() {
	method := suite.Precompile.ABI.Methods["resolveDenom"]
	cidStr := suite.TestSuite.CollectionId.String()

	callResolve := func(denom string) (*tokenizationtypes.QueryResolveDenomResponse, error) {
		queryJson, err := helpers.BuildQueryJSON(map[string]interface{}{"denom": denom})
		suite.Require().NoError(err)
		input, err := helpers.PackMethodWithJSON(&method, queryJson)
		suite.Require().NoError(err)
		raw, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, input)
		if err != nil {
			return nil, err
		}
		unpacked, err := method.Outputs.Unpack(raw)
		suite.Require().NoError(err)
		respBytes, ok := unpacked[0].([]byte)
		suite.Require().True(ok)
		var resp tokenizationtypes.QueryResolveDenomResponse
		suite.Require().NoError(proto.Unmarshal(respBytes, &resp))
		return &resp, nil
	}

	// 1. Exact wrapper match
	exactResp, err := callResolve("badges:" + cidStr + ":exact-wrapper")
	suite.Require().NoError(err)
	suite.Equal(suite.TestSuite.CollectionId, exactResp.CollectionId)
	suite.Equal("wrapper", exactResp.PathType)
	suite.Equal("exact-wrapper", exactResp.BaseDenom)
	suite.Empty(exactResp.MatchedTokenId)
	suite.Require().NotNil(exactResp.WrapperPath)
	suite.Equal("EXW", exactResp.WrapperPath.Symbol)

	// 2. Placeholder {id} match
	placeholderResp, err := callResolve("badges:" + cidStr + ":item-42")
	suite.Require().NoError(err)
	suite.Equal(suite.TestSuite.CollectionId, placeholderResp.CollectionId)
	suite.Equal("wrapper", placeholderResp.PathType)
	suite.Equal("item-42", placeholderResp.BaseDenom)
	suite.Equal("42", placeholderResp.MatchedTokenId)
	suite.Require().NotNil(placeholderResp.WrapperPath)
	suite.Equal("item-{id}", placeholderResp.WrapperPath.Denom)

	// 3. Alias match
	aliasResp, err := callResolve("badgeslp:" + cidStr + ":pool-1")
	suite.Require().NoError(err)
	suite.Equal(suite.TestSuite.CollectionId, aliasResp.CollectionId)
	suite.Equal("alias", aliasResp.PathType)
	suite.Equal("pool-1", aliasResp.BaseDenom)
	suite.Require().NotNil(aliasResp.AliasPath)
	suite.Equal("LP1", aliasResp.AliasPath.Symbol)

	// 4. Unknown denom errors
	_, err = callResolve("badges:" + cidStr + ":does-not-exist")
	suite.Error(err)

	// 5. Empty denom errors
	_, err = callResolve("")
	suite.Error(err)
}

func (suite *QueryParityTestSuite) TestGetBalanceForToken_Variants() {
	method := suite.Precompile.ABI.Methods["getBalanceForToken"]

	callBalance := func(payload map[string]interface{}) (*big.Int, error) {
		queryJson, err := helpers.BuildQueryJSON(payload)
		suite.Require().NoError(err)
		input, err := helpers.PackMethodWithJSON(&method, queryJson)
		suite.Require().NoError(err)
		raw, err := suite.TestSuite.CallPrecompileReadOnly(suite.TestSuite.AliceEVM, input)
		if err != nil {
			return nil, err
		}
		unpacked, err := method.Outputs.Unpack(raw)
		suite.Require().NoError(err)
		bal, ok := unpacked[0].(*big.Int)
		suite.Require().True(ok)
		return bal, nil
	}

	// Default timestamp (uses current block time 1735689600123 which is within [1000, 1735689600999])
	bal, err := callBalance(map[string]interface{}{
		"collectionId": suite.TestSuite.CollectionId.String(),
		"address":      suite.TestSuite.Alice.String(),
		"tokenId":      "1",
	})
	suite.Require().NoError(err)
	suite.Equal(big.NewInt(750), bal)

	// EVM 0x address + userAddress alias + explicit time
	bal, err = callBalance(map[string]interface{}{
		"collectionId": suite.TestSuite.CollectionId.String(),
		"userAddress":  suite.TestSuite.AliceEVM.Hex(),
		"tokenId":      "50",
		"time":         "5000",
	})
	suite.Require().NoError(err)
	suite.Equal(big.NewInt(750), bal)

	// Time outside ownership range -> 0
	bal, err = callBalance(map[string]interface{}{
		"collectionId": suite.TestSuite.CollectionId.String(),
		"address":      suite.TestSuite.Alice.String(),
		"tokenId":      "1",
		"time":         "500",
	})
	suite.Require().NoError(err)
	suite.Equal(0, bal.Cmp(big.NewInt(0)), fmt.Sprintf("expected 0, got %s", bal.String()))

	// Non-existent collection errors
	_, err = callBalance(map[string]interface{}{
		"collectionId": "999999",
		"address":      suite.TestSuite.Alice.String(),
		"tokenId":      "1",
	})
	suite.Error(err)
}
