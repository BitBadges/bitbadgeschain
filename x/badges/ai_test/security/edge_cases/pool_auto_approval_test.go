package edge_cases_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/badges/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
)

type PoolAutoApprovalTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestPoolAutoApprovalSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(PoolAutoApprovalTestSuite))
}

func (suite *PoolAutoApprovalTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestPoolAutoApproval_OnlySetsIfNotAlreadySet verifies that auto-approve flags
// are only set if they're not already set, preventing unintended overrides.
// This test addresses HIGH-009: Pool Integration Auto-Approval Security.
// Note: The function no longer validates that the address is a pool/path address,
// so it can be called on any address.
func (suite *PoolAutoApprovalTestSuite) TestPoolAutoApproval_OnlySetsIfNotAlreadySet() {
	collection := suite.GetCollection(suite.CollectionId)

	// Use a regular address for testing - validation is no longer performed
	poolAddress := suite.Alice

	// First, manually set one flag to true
	balance, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, poolAddress)
	balance.AutoApproveAllIncomingTransfers = true
	err := suite.Keeper.SetBalanceForAddress(suite.Ctx, collection, poolAddress, balance)
	suite.Require().NoError(err)

	// Get balance again to verify flag is set
	balanceBefore, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, poolAddress)
	suite.Require().True(balanceBefore.AutoApproveAllIncomingTransfers, "flag should be set before calling function")

	// Call the function - it should succeed and only set flags that aren't already set
	err = suite.Keeper.SetAllAutoApprovalFlagsForPoolAddress(suite.Ctx, collection, poolAddress)
	suite.Require().NoError(err, "should succeed for any address")

	// Verify that the already-set flag remains true
	balanceAfter, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, poolAddress)
	suite.Require().True(balanceAfter.AutoApproveAllIncomingTransfers, "flag that was already set should remain true")

	// Verify that other flags were set
	suite.Require().True(balanceAfter.AutoApproveSelfInitiatedOutgoingTransfers, "flag should be set")
	suite.Require().True(balanceAfter.AutoApproveSelfInitiatedIncomingTransfers, "flag should be set")
}

// TestPoolAutoApproval_IndividualFlagCheck verifies that each flag is checked
// individually and only set if not already set.
func (suite *PoolAutoApprovalTestSuite) TestPoolAutoApproval_IndividualFlagCheck() {
	collection := suite.GetCollection(suite.CollectionId)

	// Create a collection with a path address (wrapper path)
	// This will create a path address that can be validated
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
				Symbol:           "uatom",
				Decimals:         sdkmath.NewUint(6),
				IsDefaultDisplay: true,
			},
		},
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                     suite.Manager,
		CollectionId:                suite.CollectionId,
		CosmosCoinWrapperPathsToAdd: []*types.CosmosCoinWrapperPathAddObject{wrapperPath},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Get the collection again to find the path address
	collection = suite.GetCollection(suite.CollectionId)
	suite.Require().Greater(len(collection.CosmosCoinWrapperPaths), 0, "should have wrapper path")
	pathAddress := collection.CosmosCoinWrapperPaths[0].Address

	// Get initial balance
	balanceBefore, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, pathAddress)

	// Set one flag manually to true
	balanceBefore.AutoApproveSelfInitiatedOutgoingTransfers = true
	err = suite.Keeper.SetBalanceForAddress(suite.Ctx, collection, pathAddress, balanceBefore)
	suite.Require().NoError(err)

	// Call the function - it should only set the flags that aren't already set
	err = suite.Keeper.SetAllAutoApprovalFlagsForPoolAddress(suite.Ctx, collection, pathAddress)
	suite.Require().NoError(err, "should succeed for path address")

	// Verify that the already-set flag remains true
	balanceAfter, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, pathAddress)
	suite.Require().True(balanceAfter.AutoApproveSelfInitiatedOutgoingTransfers, "flag that was already set should remain true")

	// Verify that other flags were set
	suite.Require().True(balanceAfter.AutoApproveAllIncomingTransfers, "flag should be set")
	suite.Require().True(balanceAfter.AutoApproveSelfInitiatedIncomingTransfers, "flag should be set")
}

// TestPoolAutoApproval_WorksForRegularAddress verifies that the function
// works for any address, not just pool/path addresses.
// Note: Validation was removed, so the function can be called on any address.
func (suite *PoolAutoApprovalTestSuite) TestPoolAutoApproval_WorksForRegularAddress() {
	collection := suite.GetCollection(suite.CollectionId)

	// Call with a regular user address - should succeed
	err := suite.Keeper.SetAllAutoApprovalFlagsForPoolAddress(suite.Ctx, collection, suite.Alice)
	suite.Require().NoError(err, "should succeed for regular address")

	// Verify that flags were set
	balanceAfter, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, suite.Alice)
	suite.Require().True(balanceAfter.AutoApproveAllIncomingTransfers, "flag should be set")
	suite.Require().True(balanceAfter.AutoApproveSelfInitiatedOutgoingTransfers, "flag should be set")
	suite.Require().True(balanceAfter.AutoApproveSelfInitiatedIncomingTransfers, "flag should be set")
}

// TestPoolAutoApproval_NoChangeIfAllFlagsSet verifies that the function doesn't
// make unnecessary writes if all flags are already set.
func (suite *PoolAutoApprovalTestSuite) TestPoolAutoApproval_NoChangeIfAllFlagsSet() {
	collection := suite.GetCollection(suite.CollectionId)

	// Create a collection with a path address
	wrapperPath := &types.CosmosCoinWrapperPathAddObject{
		Denom: "uosmo",
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
		Symbol: "OSMO",
		DenomUnits: []*types.DenomUnit{
			{
				Symbol:           "uosmo",
				Decimals:         sdkmath.NewUint(6),
				IsDefaultDisplay: true,
			},
		},
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                     suite.Manager,
		CollectionId:                suite.CollectionId,
		CosmosCoinWrapperPathsToAdd: []*types.CosmosCoinWrapperPathAddObject{wrapperPath},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Get the path address
	collection = suite.GetCollection(suite.CollectionId)
	pathAddress := collection.CosmosCoinWrapperPaths[0].Address

	// Set all flags manually first
	balance, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, pathAddress)
	balance.AutoApproveAllIncomingTransfers = true
	balance.AutoApproveSelfInitiatedOutgoingTransfers = true
	balance.AutoApproveSelfInitiatedIncomingTransfers = true
	err = suite.Keeper.SetBalanceForAddress(suite.Ctx, collection, pathAddress, balance)
	suite.Require().NoError(err)

	// Call the function - should succeed but not change anything
	err = suite.Keeper.SetAllAutoApprovalFlagsForPoolAddress(suite.Ctx, collection, pathAddress)
	suite.Require().NoError(err, "should succeed even if all flags already set")

	// Verify flags are still set
	balanceAfter, _, _ := suite.Keeper.GetBalanceOrApplyDefault(suite.Ctx, collection, pathAddress)
	suite.Require().True(balanceAfter.AutoApproveAllIncomingTransfers, "flag should remain set")
	suite.Require().True(balanceAfter.AutoApproveSelfInitiatedOutgoingTransfers, "flag should remain set")
	suite.Require().True(balanceAfter.AutoApproveSelfInitiatedIncomingTransfers, "flag should remain set")
}
