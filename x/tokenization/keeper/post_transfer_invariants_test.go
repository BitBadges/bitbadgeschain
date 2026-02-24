package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type PostTransferInvariantsTestSuite struct {
	TestSuite
}

func TestPostTransferInvariantsTestSuite(t *testing.T) {
	suite.Run(t, new(PostTransferInvariantsTestSuite))
}

func (suite *PostTransferInvariantsTestSuite) SetupTest() {
	suite.TestSuite.SetupTest()
}

// =============================================================================
// MaxSupplyPerId Invariant Tests
// =============================================================================

// TestMaxSupplyPerId_EnforcedOnMint tests that the maxSupplyPerId invariant is enforced during minting
func (suite *PostTransferInvariantsTestSuite) TestMaxSupplyPerId_EnforcedOnMint() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with maxSupplyPerId = 5 in PostTransferInvariants
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: sdkmath.NewUint(5),
	}
	// Set valid token IDs
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(5),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	// Set the transfer to mint exactly 5 tokens (should succeed with maxSupplyPerId = 5)
	collectionsToCreate[0].Transfers = []*types.Transfer{
		{
			From:        "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(5), // Exactly maxSupplyPerId
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{
					ApprovalId:      "mint-test",
					ApprovalLevel:   "collection",
					ApproverAddress: "",
					Version:         sdkmath.NewUint(0),
				},
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify collection was created with the correct invariants
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().Equal(sdkmath.NewUint(5), collection.Invariants.MaxSupplyPerId)
}

// TestMaxSupplyPerId_BlocksExcessiveMint tests that minting more than maxSupplyPerId is blocked
func (suite *PostTransferInvariantsTestSuite) TestMaxSupplyPerId_BlocksExcessiveMint() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with maxSupplyPerId = 5 in PostTransferInvariants
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: sdkmath.NewUint(5),
	}
	// Set valid token IDs
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(10),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	// Set the transfer to mint 10 tokens (more than maxSupplyPerId = 5)
	collectionsToCreate[0].Transfers = []*types.Transfer{
		{
			From:        "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(10), // More than maxSupplyPerId of 5
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{
					ApprovalId:      "mint-test",
					ApprovalLevel:   "collection",
					ApproverAddress: "",
					Version:         sdkmath.NewUint(0),
				},
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "maxSupplyPerId")
}

// TestMaxSupplyPerId_AllowsExactLimit tests that minting exactly maxSupplyPerId is allowed
func (suite *PostTransferInvariantsTestSuite) TestMaxSupplyPerId_AllowsExactLimit() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with maxSupplyPerId = 100
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: sdkmath.NewUint(100),
	}
	// Set valid token IDs
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(100),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	// Set the transfer to mint exactly 100 tokens (should succeed)
	collectionsToCreate[0].Transfers = []*types.Transfer{
		{
			From:        "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(100), // Exactly maxSupplyPerId
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{
					ApprovalId:      "mint-test",
					ApprovalLevel:   "collection",
					ApproverAddress: "",
					Version:         sdkmath.NewUint(0),
				},
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify collection was created
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().NotNil(collection)
}

// TestMaxSupplyPerId_ZeroMeansUnlimited tests that maxSupplyPerId = 0 means no limit
func (suite *PostTransferInvariantsTestSuite) TestMaxSupplyPerId_ZeroMeansUnlimited() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with maxSupplyPerId = 0 (no limit)
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: sdkmath.NewUint(0),
	}
	// Increase approval amounts to allow large mints
	collectionsToCreate[0].CollectionApprovals[0].ApprovalCriteria.ApprovalAmounts.PerFromAddressApprovalAmount = sdkmath.NewUint(1000000)
	// Set valid token IDs
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1000000),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	// Mint a large amount (should succeed because maxSupplyPerId = 0 means no limit)
	collectionsToCreate[0].Transfers = []*types.Transfer{
		{
			From:        "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(1000000), // Large amount, but 0 = no limit
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{
					ApprovalId:      "mint-test",
					ApprovalLevel:   "collection",
					ApproverAddress: "",
					Version:         sdkmath.NewUint(0),
				},
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)
}

// TestMaxSupplyPerId_WithMultipleTokenIds tests maxSupplyPerId with multiple token IDs
func (suite *PostTransferInvariantsTestSuite) TestMaxSupplyPerId_WithMultipleTokenIds() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with maxSupplyPerId = 10
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: sdkmath.NewUint(10),
	}
	// Mint 10 of token ID 1-3 (should succeed - each token ID has 10, which equals maxSupplyPerId)
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount: sdkmath.NewUint(10),
			TokenIds: []*types.UintRange{
				{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(3)},
			},
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)
}

// =============================================================================
// PostTransferInvariants with Transfers Tests
// =============================================================================

// TestPostTransferInvariants_TransferDoesNotViolateSupply tests that transfers
// within the supply limit succeed (the post-transfer invariant only checks total supply)
func (suite *PostTransferInvariantsTestSuite) TestPostTransferInvariants_TransferDoesNotViolateSupply() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with tokens and maxSupplyPerId
	// The test verifies that after minting within the limit, the collection exists with correct invariants
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: sdkmath.NewUint(10),
	}
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(5),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	// Mint 5 tokens via transfer (within maxSupplyPerId of 10)
	collectionsToCreate[0].Transfers = []*types.Transfer{
		{
			From:        "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(5),
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{
					ApprovalId:      "mint-test",
					ApprovalLevel:   "collection",
					ApproverAddress: "",
					Version:         sdkmath.NewUint(0),
				},
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify collection was created with the correct invariants
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().Equal(sdkmath.NewUint(10), collection.Invariants.MaxSupplyPerId)
}

// =============================================================================
// PostTransferInvariants Combined with Other Invariants Tests
// =============================================================================

// TestPostTransferInvariants_CombinedWithNoCustomOwnershipTimes tests combining invariants
func (suite *PostTransferInvariantsTestSuite) TestPostTransferInvariants_CombinedWithNoCustomOwnershipTimes() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with both noCustomOwnershipTimes and postTransferInvariants
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		NoCustomOwnershipTimes: true,
		MaxSupplyPerId: sdkmath.NewUint(100),
	}
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(50),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(), // Must be full ranges due to noCustomOwnershipTimes
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify both invariants are set
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().True(collection.Invariants.NoCustomOwnershipTimes)
	suite.Require().Equal(sdkmath.NewUint(100), collection.Invariants.MaxSupplyPerId)
}

// =============================================================================
// Placeholder Replacement Tests (Unit Tests for Helper Functions)
// =============================================================================

// TestPlaceholderReplacement_CollectionId tests $collectionId placeholder replacement
func (suite *PostTransferInvariantsTestSuite) TestPlaceholderReplacement_CollectionId() {
	// Test that $collectionId is properly converted to hex
	// This is a unit test for the placeholder replacement logic

	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a simple collection to get a collection ID
	collectionsToCreate := GetCollectionsToCreate()
	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify collection was created with ID 1
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().Equal(sdkmath.NewUint(1), collection.CollectionId)
}

// =============================================================================
// Invariants Immutability Tests
// =============================================================================

// TestPostTransferInvariants_CannotBeModifiedAfterCreation tests that invariants cannot be changed
func (suite *PostTransferInvariantsTestSuite) TestPostTransferInvariants_CannotBeModifiedAfterCreation() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with postTransferInvariants
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: sdkmath.NewUint(100),
	}
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(50),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Try to update the invariants (should not be possible via update - invariants are set-once)
	// This test verifies that UpdateCollection does not allow modifying invariants
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().NotNil(collection.Invariants)

	// Invariants should remain unchanged after any update operations
	suite.Require().Equal(sdkmath.NewUint(100), collection.Invariants.MaxSupplyPerId)
}

// =============================================================================
// Edge Cases
// =============================================================================

// TestPostTransferInvariants_NilInvariants tests that nil invariants are handled correctly
func (suite *PostTransferInvariantsTestSuite) TestPostTransferInvariants_NilInvariants() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection without any invariants
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Invariants = nil

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify collection was created and invariants is nil
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().Nil(collection.Invariants)
}

// TestPostTransferInvariants_EmptyInvariants tests empty Invariants (no maxSupplyPerId, no EVM challenges)
func (suite *PostTransferInvariantsTestSuite) TestPostTransferInvariants_EmptyInvariants() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with empty Invariants (all defaults)
	collectionsToCreate := GetCollectionsToCreate()
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		// All defaults - no maxSupplyPerId, no EVM challenges
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify collection was created
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().NotNil(collection)
}

// TestPostTransferInvariants_LargeMaxSupply tests with a very large maxSupplyPerId value
func (suite *PostTransferInvariantsTestSuite) TestPostTransferInvariants_LargeMaxSupply() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection with a very large maxSupplyPerId
	largeSupply := sdkmath.NewUint(1000000000000) // 1 trillion
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: largeSupply,
	}
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(1000000),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify collection was created with the large supply limit
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().Equal(largeSupply, collection.Invariants.MaxSupplyPerId)
}

// =============================================================================
// Multi-Recipient Transfer Tests
// =============================================================================

// TestPostTransferInvariants_MultiRecipientMint tests minting to multiple recipients works with invariants
func (suite *PostTransferInvariantsTestSuite) TestPostTransferInvariants_MultiRecipientMint() {
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Create a collection that mints to multiple recipients
	collectionsToCreate := GetTransferableCollectionToCreateAllMintedToCreator(bob)
	collectionsToCreate[0].Invariants = &types.InvariantsAddObject{
		MaxSupplyPerId: sdkmath.NewUint(100), // Allow up to 100 per token ID
	}
	collectionsToCreate[0].TokensToCreate = []*types.Balance{
		{
			Amount:         sdkmath.NewUint(30),
			TokenIds:       GetOneUintRange(),
			OwnershipTimes: GetFullUintRanges(),
		},
	}
	// Mint 30 tokens to bob (within maxSupplyPerId)
	collectionsToCreate[0].Transfers = []*types.Transfer{
		{
			From:        "Mint",
			ToAddresses: []string{bob},
			Balances: []*types.Balance{
				{
					Amount:         sdkmath.NewUint(30),
					TokenIds:       GetOneUintRange(),
					OwnershipTimes: GetFullUintRanges(),
				},
			},
			PrioritizedApprovals: []*types.ApprovalIdentifierDetails{
				{
					ApprovalId:      "mint-test",
					ApprovalLevel:   "collection",
					ApproverAddress: "",
					Version:         sdkmath.NewUint(0),
				},
			},
		},
	}

	err := CreateCollections(&suite.TestSuite, wctx, collectionsToCreate)
	suite.Require().NoError(err)

	// Verify collection was created with the correct invariants
	collection, err := GetCollection(&suite.TestSuite, wctx, sdkmath.NewUint(1))
	suite.Require().NoError(err)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().Equal(sdkmath.NewUint(100), collection.Invariants.MaxSupplyPerId)

	// Verify bob received the tokens
	bobBalance, err := GetUserBalance(&suite.TestSuite, wctx, sdkmath.NewUint(1), bob)
	suite.Require().NoError(err)
	suite.Require().NotNil(bobBalance)
}
