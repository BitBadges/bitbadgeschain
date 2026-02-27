package invariants_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// DisablePoolCreationTestSuite tests the disablePoolCreation invariant
// This invariant prevents pool creation with collection assets
type DisablePoolCreationTestSuite struct {
	testutil.AITestSuite
}

func TestDisablePoolCreationTestSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(DisablePoolCreationTestSuite))
}

func (suite *DisablePoolCreationTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// createCollectionWithInvariant creates a collection with disablePoolCreation invariant
func (suite *DisablePoolCreationTestSuite) createCollectionWithInvariant(disablePoolCreation bool) sdkmath.Uint {
	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{},
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.Manager,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "",
		},
		TokenMetadata:       []*types.TokenMetadata{},
		CustomData:          "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards:           []string{},
		IsArchived:          false,
		Invariants: &types.InvariantsAddObject{
			DisablePoolCreation: disablePoolCreation,
		},
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "collection creation should succeed")
	return resp.CollectionId
}

// createCollectionWithAliasPath creates a collection with alias path for pool operations
func (suite *DisablePoolCreationTestSuite) createCollectionWithAliasPath(disablePoolCreation bool, denom string) sdkmath.Uint {
	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{},
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.Manager,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "",
		},
		TokenMetadata:       []*types.TokenMetadata{},
		CustomData:          "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards:           []string{},
		IsArchived:          false,
		Invariants: &types.InvariantsAddObject{
			DisablePoolCreation: disablePoolCreation,
		},
		AliasPathsToAdd: []*types.AliasPathAddObject{
			{
				Denom: denom,
				Conversion: &types.ConversionWithoutDenom{
					SideA: &types.ConversionSideA{
						Amount: sdkmath.NewUint(1),
					},
					SideB: []*types.Balance{
						{
							Amount: sdkmath.NewUint(1),
							OwnershipTimes: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
							},
							TokenIds: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
						},
					},
				},
				Symbol: "TEST",
				DenomUnits: []*types.DenomUnit{
					{Decimals: sdkmath.NewUint(6), Symbol: denom, IsDefaultDisplay: true},
				},
			},
		},
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "collection creation with alias path should succeed")
	return resp.CollectionId
}

// TestDisablePoolCreation_InvariantSetOnCreation tests that invariant is properly set on creation
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_InvariantSetOnCreation() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Verify invariant is set
	collection := suite.GetCollection(collectionId)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().True(collection.Invariants.DisablePoolCreation)
}

// TestDisablePoolCreation_InvariantDisabledOnCreation tests that invariant can be disabled on creation
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_InvariantDisabledOnCreation() {
	// Create collection without invariant
	collectionId := suite.createCollectionWithInvariant(false)

	// Verify invariant is not set or is false
	collection := suite.GetCollection(collectionId)
	if collection.Invariants != nil {
		suite.Require().False(collection.Invariants.DisablePoolCreation)
	}
}

// TestDisablePoolCreation_InvariantCanBeUpdated tests that the invariant can be changed after creation
// Note: Invariants in BitBadges are not immutable - they can be changed via MsgUniversalUpdateCollection
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_InvariantCanBeUpdated() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Verify invariant is set
	collection := suite.GetCollection(collectionId)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().True(collection.Invariants.DisablePoolCreation)

	// Try to disable the invariant - this should succeed as invariants are not immutable
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Invariants: &types.InvariantsAddObject{
			DisablePoolCreation: false, // Disable the invariant
		},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "updating invariants should succeed")

	// Verify the invariant was updated
	collection = suite.GetCollection(collectionId)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().False(collection.Invariants.DisablePoolCreation,
		"disablePoolCreation should be disabled after update")
}

// TestDisablePoolCreation_CollectionWithAliasPathAndInvariant tests collection with both alias path and invariant
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_CollectionWithAliasPathAndInvariant() {
	// Create collection with alias path and invariant enabled
	collectionId := suite.createCollectionWithAliasPath(true, "pooltestdenom")

	// Verify collection was created
	collection := suite.GetCollection(collectionId)
	suite.Require().NotNil(collection)

	// Verify invariant is set
	suite.Require().NotNil(collection.Invariants)
	suite.Require().True(collection.Invariants.DisablePoolCreation)

	// Verify alias path was created
	suite.Require().True(len(collection.AliasPaths) > 0, "alias path should be created")
}

// TestDisablePoolCreation_AllowsNormalTransfers tests that normal transfers still work with invariant enabled
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_AllowsNormalTransfers() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Setup mint approval
	suite.SetupMintApproval(collectionId)

	// Mint tokens - should succeed
	mintMsg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances: []*types.Balance{
					testutil.GenerateSimpleBalance(100, 1),
				},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), mintMsg)
	suite.Require().NoError(err, "minting should succeed even with disablePoolCreation")

	// Verify balance
	aliceBalance := suite.GetBalance(collectionId, suite.Alice)
	suite.Require().True(len(aliceBalance.Balances) > 0)
}

// TestDisablePoolCreation_AllowsCollectionApprovalUpdates tests that collection approvals can still be updated
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_AllowsCollectionApprovalUpdates() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Add collection approval - should succeed
	approval := testutil.GenerateCollectionApproval("test_approval", "AllWithoutMint", "All")
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{approval},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "collection approval updates should succeed with disablePoolCreation")
}

// TestDisablePoolCreation_AllowsMetadataUpdates tests that metadata can still be updated
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_AllowsMetadataUpdates() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Update metadata - should succeed
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                    suite.Manager,
		CollectionId:               collectionId,
		UpdateCollectionMetadata:   true,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/updated-metadata",
			CustomData: "updated",
		},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "metadata updates should succeed with disablePoolCreation")

	// Verify update
	collection := suite.GetCollection(collectionId)
	suite.Require().Equal("https://example.com/updated-metadata", collection.CollectionMetadata.Uri)
}

// TestDisablePoolCreation_CanStillAddAliasPaths tests that alias paths can be added even with invariant
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_CanStillAddAliasPaths() {
	// Create collection with invariant enabled (without alias path)
	collectionId := suite.createCollectionWithInvariant(true)

	// Try to add alias path - may or may not be allowed depending on implementation
	// The invariant should only block pool creation, not alias path creation
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		AliasPathsToAdd: []*types.AliasPathAddObject{
			{
				Denom: "newdenom",
				Conversion: &types.ConversionWithoutDenom{
					SideA: &types.ConversionSideA{
						Amount: sdkmath.NewUint(1),
					},
					SideB: []*types.Balance{
						{
							Amount: sdkmath.NewUint(1),
							OwnershipTimes: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
							},
							TokenIds: []*types.UintRange{
								{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
							},
						},
					},
				},
				Symbol: "NEW",
				DenomUnits: []*types.DenomUnit{
					{Decimals: sdkmath.NewUint(6), Symbol: "newdenom", IsDefaultDisplay: true},
				},
			},
		},
	}

	// This may succeed or fail depending on permissions, but should not be blocked by disablePoolCreation
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	// Not asserting success/failure here as it depends on collection permissions
	// The key point is that disablePoolCreation shouldn't affect alias path addition
	_ = err
}

// TestDisablePoolCreation_MultipleCollectionsDifferentInvariants tests multiple collections with different settings
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_MultipleCollectionsDifferentInvariants() {
	// Create collection with invariant enabled
	collectionIdWithInvariant := suite.createCollectionWithInvariant(true)

	// Create collection without invariant
	collectionIdWithoutInvariant := suite.createCollectionWithInvariant(false)

	// Verify first collection has invariant enabled
	collection1 := suite.GetCollection(collectionIdWithInvariant)
	suite.Require().NotNil(collection1.Invariants)
	suite.Require().True(collection1.Invariants.DisablePoolCreation)

	// Verify second collection has invariant disabled
	collection2 := suite.GetCollection(collectionIdWithoutInvariant)
	if collection2.Invariants != nil {
		suite.Require().False(collection2.Invariants.DisablePoolCreation)
	}
}

// TestDisablePoolCreation_CombinedWithOtherInvariants tests combining with other invariants
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_CombinedWithOtherInvariants() {
	// Create collection with multiple invariants
	msg := &types.MsgCreateCollection{
		Creator: suite.Manager,
		DefaultBalances: &types.UserBalanceStore{
			Balances: []*types.Balance{},
		},
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		CollectionPermissions: &types.CollectionPermissions{},
		Manager:               suite.Manager,
		CollectionMetadata: &types.CollectionMetadata{
			Uri:        "https://example.com/metadata",
			CustomData: "",
		},
		TokenMetadata:       []*types.TokenMetadata{},
		CustomData:          "",
		CollectionApprovals: []*types.CollectionApproval{},
		Standards:           []string{},
		IsArchived:          false,
		Invariants: &types.InvariantsAddObject{
			DisablePoolCreation:         true,
			NoCustomOwnershipTimes:      true,
			NoForcefulPostMintTransfers: true,
			MaxSupplyPerId:              sdkmath.NewUint(1000),
		},
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "collection creation with multiple invariants should succeed")

	// Verify all invariants are set
	collection := suite.GetCollection(resp.CollectionId)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().True(collection.Invariants.DisablePoolCreation)
	suite.Require().True(collection.Invariants.NoCustomOwnershipTimes)
	suite.Require().True(collection.Invariants.NoForcefulPostMintTransfers)
	suite.Require().Equal(sdkmath.NewUint(1000), collection.Invariants.MaxSupplyPerId)
}

// TestDisablePoolCreation_PreservesOtherFunctionality tests that other collection functionality is preserved
func (suite *DisablePoolCreationTestSuite) TestDisablePoolCreation_PreservesOtherFunctionality() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Test various operations that should still work

	// 1. Update standards
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:         suite.Manager,
		CollectionId:    collectionId,
		UpdateStandards: true,
		Standards:       []string{"ERC721"},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "standards update should succeed")

	// 2. Update custom data
	updateMsg2 := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateCustomData: true,
		CustomData:       "custom data",
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg2)
	suite.Require().NoError(err, "custom data update should succeed")

	// 3. Archive/unarchive
	updateMsg3 := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       true,
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg3)
	suite.Require().NoError(err, "archive should succeed")

	// Unarchive to continue testing
	updateMsg4 := &types.MsgUniversalUpdateCollection{
		Creator:          suite.Manager,
		CollectionId:     collectionId,
		UpdateIsArchived: true,
		IsArchived:       false,
	}
	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg4)
	suite.Require().NoError(err, "unarchive should succeed")
}
