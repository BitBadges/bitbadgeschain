package edge_cases

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/keeper"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type PathAddressGenerationTestSuite struct {
	testutil.AITestSuite
}

func TestPathAddressGenerationTestSuite(t *testing.T) {
	suite.Run(t, new(PathAddressGenerationTestSuite))
}

// TestPathAddressGeneration_EmptyDenom tests that empty denom is rejected
func (suite *PathAddressGenerationTestSuite) TestPathAddressGeneration_EmptyDenom() {
	_, err := keeper.GenerateAliasPathAddress("")
	suite.Require().Error(err, "empty denom should be rejected")
	suite.Require().Contains(err.Error(), "path string cannot be empty", "error should indicate empty path string")
}

// TestPathAddressGeneration_LongDenom tests that very long denoms are rejected
func (suite *PathAddressGenerationTestSuite) TestPathAddressGeneration_LongDenom() {
	// Create a denom that exceeds the 1024 byte limit
	longDenom := make([]byte, 1025)
	for i := range longDenom {
		longDenom[i] = 'a'
	}
	
	_, err := keeper.GenerateAliasPathAddress(string(longDenom))
	suite.Require().Error(err, "denom exceeding 1024 bytes should be rejected")
	suite.Require().Contains(err.Error(), "exceeds maximum length", "error should indicate length limit")
}

// TestPathAddressGeneration_ValidDenom tests that valid denoms generate addresses correctly
func (suite *PathAddressGenerationTestSuite) TestPathAddressGeneration_ValidDenom() {
	denom := "uatom"
	addr, err := keeper.GenerateAliasPathAddress(denom)
	suite.Require().NoError(err, "valid denom should generate address")
	suite.Require().NotEmpty(addr, "generated address should not be empty")
	
	// Verify address is valid Bech32
	_, err = sdk.AccAddressFromBech32(addr)
	suite.Require().NoError(err, "generated address should be valid Bech32")
}

// TestPathAddressGeneration_Deterministic tests that same denom generates same address
func (suite *PathAddressGenerationTestSuite) TestPathAddressGeneration_Deterministic() {
	denom := "uatom"
	addr1, err1 := keeper.GenerateAliasPathAddress(denom)
	addr2, err2 := keeper.GenerateAliasPathAddress(denom)
	
	suite.Require().NoError(err1)
	suite.Require().NoError(err2)
	suite.Require().Equal(addr1, addr2, "same denom should generate same address")
}

// TestPathAddressGeneration_ReservedProtocolAddress tests that generated addresses are marked as reserved
func (suite *PathAddressGenerationTestSuite) TestPathAddressGeneration_ReservedProtocolAddress() {
	// Create a collection with a wrapper path
	collectionId := suite.CreateTestCollection(suite.Manager)
	
	// Add a wrapper path - this should automatically mark the address as reserved
	wrapperPath := &types.CosmosCoinWrapperPathAddObject{
		Denom: "uatom",
		Conversion: &types.ConversionWithoutDenom{
			SideA: &types.ConversionSideA{
				Amount: sdkmath.NewUint(1),
			},
			SideB: []*types.Balance{
				{
					Amount: sdkmath.NewUint(1),
					TokenIds: []*types.UintRange{
						{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
					},
					OwnershipTimes: []*types.UintRange{
						{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
					},
				},
			},
		},
		Symbol: "ATOM",
		DenomUnits: []*types.DenomUnit{
			{
				Symbol:          "uatom",
				Decimals:        sdkmath.NewUint(6),
				IsDefaultDisplay: true,
			},
		},
	}
	
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                     suite.Manager,
		CollectionId:                 collectionId,
		CosmosCoinWrapperPathsToAdd: []*types.CosmosCoinWrapperPathAddObject{wrapperPath},
	}
	
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "should be able to add wrapper path")
	
	// Verify the generated address is marked as reserved
	generatedAddr, err := keeper.GenerateAliasPathAddress("uatom")
	suite.Require().NoError(err)
	
	isReserved := suite.Keeper.IsAddressReservedProtocolInStore(suite.Ctx, generatedAddr)
	suite.Require().True(isReserved, "generated path address should be marked as reserved protocol address")
}

// TestPathAddressGeneration_DifferentDenomsDifferentAddresses tests that different denoms generate different addresses
func (suite *PathAddressGenerationTestSuite) TestPathAddressGeneration_DifferentDenomsDifferentAddresses() {
	addr1, err1 := keeper.GenerateAliasPathAddress("uatom")
	addr2, err2 := keeper.GenerateAliasPathAddress("uosmo")
	
	suite.Require().NoError(err1)
	suite.Require().NoError(err2)
	suite.Require().NotEqual(addr1, addr2, "different denoms should generate different addresses")
}

