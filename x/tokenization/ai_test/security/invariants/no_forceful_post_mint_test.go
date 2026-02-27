package invariants_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// NoForcefulPostMintTestSuite tests the noForcefulPostMintTransfers invariant
// This invariant ensures that post-mint transfers cannot use approval overrides
type NoForcefulPostMintTestSuite struct {
	testutil.AITestSuite
}

func TestNoForcefulPostMintTestSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(NoForcefulPostMintTestSuite))
}

func (suite *NoForcefulPostMintTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// fullUintRanges returns full uint ranges [1, MaxUint64]
func fullUintRanges() []*types.UintRange {
	return []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}
}

// createCollectionWithInvariant creates a collection with noForcefulPostMintTransfers invariant
func (suite *NoForcefulPostMintTestSuite) createCollectionWithInvariant(noForcefulPostMint bool) sdkmath.Uint {
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
			NoForcefulPostMintTransfers: noForcefulPostMint,
		},
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "collection creation should succeed")
	return resp.CollectionId
}

// TestNoForcefulPostMint_MintOperationsCanUseOverrides tests that mint operations can use overrides
func (suite *NoForcefulPostMintTestSuite) TestNoForcefulPostMint_MintOperationsCanUseOverrides() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Add mint approval with overrides - should succeed since from=Mint
	mintApproval := &types.CollectionApproval{
		ApprovalId:        "mint_with_override",
		FromListId:        types.MintAddress, // Mint address
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     fullUintRanges(),
		TokenIds:          fullUintRanges(),
		OwnershipTimes:    fullUintRanges(),
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true, // Allowed for Mint
			OverridesToIncomingApprovals:   true, // Allowed for Mint
		},
		Version: sdkmath.NewUint(0),
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{mintApproval},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "mint approval with overrides should succeed when invariant is enabled")

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

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), mintMsg)
	suite.Require().NoError(err, "minting with overrides should succeed")
}

// TestNoForcefulPostMint_PostMintCannotUseOutgoingOverrides tests that post-mint transfers cannot use outgoing overrides
func (suite *NoForcefulPostMintTestSuite) TestNoForcefulPostMint_PostMintCannotUseOutgoingOverrides() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Try to add approval for non-Mint address with outgoing override - should fail
	badApproval := &types.CollectionApproval{
		ApprovalId:        "bad_outgoing_override",
		FromListId:        "AllWithoutMint", // Not Mint
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     fullUintRanges(),
		TokenIds:          fullUintRanges(),
		OwnershipTimes:    fullUintRanges(),
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true, // Not allowed for non-Mint
		},
		Version: sdkmath.NewUint(0),
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{badApproval},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().Error(err, "approval with outgoing override for non-Mint should fail")
	suite.Require().Contains(err.Error(), "overridesFromOutgoingApprovals")
}

// TestNoForcefulPostMint_PostMintCannotUseIncomingOverrides tests that post-mint transfers cannot use incoming overrides
func (suite *NoForcefulPostMintTestSuite) TestNoForcefulPostMint_PostMintCannotUseIncomingOverrides() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Try to add approval for non-Mint address with incoming override - should fail
	badApproval := &types.CollectionApproval{
		ApprovalId:        "bad_incoming_override",
		FromListId:        "AllWithoutMint", // Not Mint
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     fullUintRanges(),
		TokenIds:          fullUintRanges(),
		OwnershipTimes:    fullUintRanges(),
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesToIncomingApprovals: true, // Not allowed for non-Mint
		},
		Version: sdkmath.NewUint(0),
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{badApproval},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().Error(err, "approval with incoming override for non-Mint should fail")
	suite.Require().Contains(err.Error(), "overridesToIncomingApprovals")
}

// TestNoForcefulPostMint_InvariantEnforcesUserApprovalRespect tests that the invariant respects user-level approvals
func (suite *NoForcefulPostMintTestSuite) TestNoForcefulPostMint_InvariantEnforcesUserApprovalRespect() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// First, set up mint approval properly
	suite.SetupMintApproval(collectionId)

	// Add normal transfer approval (without overrides)
	normalApproval := &types.CollectionApproval{
		ApprovalId:        "normal_transfer",
		FromListId:        "AllWithoutMint",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     fullUintRanges(),
		TokenIds:          fullUintRanges(),
		OwnershipTimes:    fullUintRanges(),
		ApprovalCriteria: &types.ApprovalCriteria{
			// No overrides - respects user approvals
			OverridesFromOutgoingApprovals: false,
			OverridesToIncomingApprovals:   false,
		},
		Version: sdkmath.NewUint(0),
	}

	// Get current collection to preserve mint approval
	collection := suite.GetCollection(collectionId)
	approvals := append(collection.CollectionApprovals, normalApproval)

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       approvals,
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "normal transfer approval without overrides should succeed")

	// Mint tokens to Alice
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
	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), mintMsg)
	suite.Require().NoError(err)

	// Without user approvals, transfer should fail (because we respect user approvals)
	transferMsg := &types.MsgTransferTokens{
		Creator:      suite.Alice,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        suite.Alice,
				ToAddresses: []string{suite.Bob},
				Balances: []*types.Balance{
					testutil.GenerateSimpleBalance(50, 1),
				},
			},
		},
	}

	_, err = suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), transferMsg)
	// This should fail without proper user-level approvals
	suite.Require().Error(err, "transfer should fail without user-level approvals")
}

// TestNoForcefulPostMint_InvariantDisabledAllowsOverrides tests that overrides work when invariant is disabled
func (suite *NoForcefulPostMintTestSuite) TestNoForcefulPostMint_InvariantDisabledAllowsOverrides() {
	// Create collection WITHOUT invariant enabled
	collectionId := suite.createCollectionWithInvariant(false)

	// Verify invariant is not set or is false
	collection := suite.GetCollection(collectionId)
	if collection.Invariants != nil {
		suite.Require().False(collection.Invariants.NoForcefulPostMintTransfers)
	}

	// Add approval with overrides for non-Mint address - should succeed
	overrideApproval := &types.CollectionApproval{
		ApprovalId:        "override_allowed",
		FromListId:        "AllWithoutMint",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     fullUintRanges(),
		TokenIds:          fullUintRanges(),
		OwnershipTimes:    fullUintRanges(),
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true,
			OverridesToIncomingApprovals:   true,
		},
		Version: sdkmath.NewUint(0),
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{overrideApproval},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "approval with overrides should succeed when invariant is disabled")
}

// TestNoForcefulPostMint_MixedAddressListWithMintRejected tests that a list containing Mint AND other addresses still fails
func (suite *NoForcefulPostMintTestSuite) TestNoForcefulPostMint_MixedAddressListWithMintRejected() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Create an address list that includes both Mint and other addresses
	// Note: Address list IDs must be alphanumeric only (no underscores/hyphens)
	err := suite.Keeper.CreateAddressList(suite.Ctx, &types.AddressList{
		ListId:    "mintandalice",
		Addresses: []string{types.MintAddress, suite.Alice},
		Whitelist: true,
		CreatedBy: suite.Manager,
	})
	suite.Require().NoError(err, "should be able to create address list")

	// Try to add approval using this mixed list with overrides - should fail
	badApproval := &types.CollectionApproval{
		ApprovalId:        "mixedlistoverride",
		FromListId:        "mintandalice", // Contains Mint but also others
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     fullUintRanges(),
		TokenIds:          fullUintRanges(),
		OwnershipTimes:    fullUintRanges(),
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true, // Not allowed because list includes non-Mint
		},
		Version: sdkmath.NewUint(0),
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{badApproval},
	}

	_, err = suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().Error(err, "approval with mixed address list containing Mint should fail")
}

// TestNoForcefulPostMint_BothOverridesFail tests that both overrides are rejected for non-Mint
func (suite *NoForcefulPostMintTestSuite) TestNoForcefulPostMint_BothOverridesFail() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Try to add approval with both overrides - should fail
	badApproval := &types.CollectionApproval{
		ApprovalId:        "both_overrides",
		FromListId:        suite.Alice, // Not Mint
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     fullUintRanges(),
		TokenIds:          fullUintRanges(),
		OwnershipTimes:    fullUintRanges(),
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: true, // Not allowed
			OverridesToIncomingApprovals:   true, // Not allowed
		},
		Version: sdkmath.NewUint(0),
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{badApproval},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().Error(err, "approval with both overrides for non-Mint should fail")
}

// TestNoForcefulPostMint_NoOverridesAllowed tests that approvals without overrides are allowed
func (suite *NoForcefulPostMintTestSuite) TestNoForcefulPostMint_NoOverridesAllowed() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Add approval without any overrides - should succeed
	goodApproval := &types.CollectionApproval{
		ApprovalId:        "no_overrides",
		FromListId:        "AllWithoutMint",
		ToListId:          "All",
		InitiatedByListId: "All",
		TransferTimes:     fullUintRanges(),
		TokenIds:          fullUintRanges(),
		OwnershipTimes:    fullUintRanges(),
		ApprovalCriteria: &types.ApprovalCriteria{
			OverridesFromOutgoingApprovals: false, // Respects user approvals
			OverridesToIncomingApprovals:   false, // Respects user approvals
		},
		Version: sdkmath.NewUint(0),
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals:       []*types.CollectionApproval{goodApproval},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "approval without overrides should succeed")
}

// TestNoForcefulPostMint_InvariantCanBeUpdated tests that the invariant can be updated after creation
// Note: Invariants in BitBadges are not immutable - they can be changed via MsgUniversalUpdateCollection
func (suite *NoForcefulPostMintTestSuite) TestNoForcefulPostMint_InvariantCanBeUpdated() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Verify invariant is set
	collection := suite.GetCollection(collectionId)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().True(collection.Invariants.NoForcefulPostMintTransfers)

	// Try to disable the invariant - this should succeed as invariants are not immutable
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Invariants: &types.InvariantsAddObject{
			NoForcefulPostMintTransfers: false, // Disable the invariant
		},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "updating invariants should succeed")

	// Verify invariant was updated
	collection = suite.GetCollection(collectionId)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().False(collection.Invariants.NoForcefulPostMintTransfers,
		"noForcefulPostMintTransfers should be disabled after update")
}
