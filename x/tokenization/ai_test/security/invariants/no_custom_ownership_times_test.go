package invariants_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

// NoCustomOwnershipTimesTestSuite tests the noCustomOwnershipTimes invariant
// This invariant requires all ownership times to be full ranges [1, MaxUint64]
type NoCustomOwnershipTimesTestSuite struct {
	testutil.AITestSuite
}

func TestNoCustomOwnershipTimesTestSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(NoCustomOwnershipTimesTestSuite))
}

func (suite *NoCustomOwnershipTimesTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
}

// createCollectionWithInvariant creates a collection with noCustomOwnershipTimes invariant
func (suite *NoCustomOwnershipTimesTestSuite) createCollectionWithInvariant(noCustomOwnershipTimes bool) sdkmath.Uint {
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
			NoCustomOwnershipTimes: noCustomOwnershipTimes,
		},
	}

	resp, err := suite.MsgServer.CreateCollection(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "collection creation should succeed")
	return resp.CollectionId
}

// TestNoCustomOwnershipTimes_InvariantCanBeUpdated tests that invariant can be updated after creation
// Note: Invariants in BitBadges are not immutable - they can be changed via MsgUniversalUpdateCollection
func (suite *NoCustomOwnershipTimesTestSuite) TestNoCustomOwnershipTimes_InvariantCanBeUpdated() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Verify invariant is set
	collection := suite.GetCollection(collectionId)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().True(collection.Invariants.NoCustomOwnershipTimes)

	// Try to update and disable the invariant - this should succeed as invariants are not immutable
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Invariants: &types.InvariantsAddObject{
			NoCustomOwnershipTimes: false, // Disable the invariant
		},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err, "updating invariants should succeed")

	// Verify invariant was updated
	collection = suite.GetCollection(collectionId)
	suite.Require().NotNil(collection.Invariants)
	suite.Require().False(collection.Invariants.NoCustomOwnershipTimes,
		"noCustomOwnershipTimes invariant should be disabled after update")
}

// TestNoCustomOwnershipTimes_FullRangeTransferAllowed tests that transfers with full ownership times succeed
func (suite *NoCustomOwnershipTimesTestSuite) TestNoCustomOwnershipTimes_FullRangeTransferAllowed() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Setup mint approval
	suite.SetupMintApproval(collectionId)

	// Mint tokens with full ownership times [1, MaxUint64] - should succeed
	fullOwnershipTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	mintMsg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(100),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: fullOwnershipTimes,
					},
				},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), mintMsg)
	suite.Require().NoError(err, "transfer with full ownership times should succeed")
}

// TestNoCustomOwnershipTimes_RestrictedTimesRejected tests that transfers with restricted ownership times fail
func (suite *NoCustomOwnershipTimesTestSuite) TestNoCustomOwnershipTimes_RestrictedTimesRejected() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Setup mint approval
	suite.SetupMintApproval(collectionId)

	// Try to mint tokens with restricted ownership times - should fail
	restrictedOwnershipTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(200)}, // Not full range
	}

	mintMsg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(100),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: restrictedOwnershipTimes,
					},
				},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), mintMsg)
	suite.Require().Error(err, "transfer with restricted ownership times should fail")
	suite.Require().Contains(err.Error(), "noCustomOwnershipTimes")
}

// TestNoCustomOwnershipTimes_StartNotOneRejected tests that ownership times not starting at 1 fail
func (suite *NoCustomOwnershipTimesTestSuite) TestNoCustomOwnershipTimes_StartNotOneRejected() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Setup mint approval
	suite.SetupMintApproval(collectionId)

	// Try with start = 2 (not 1) - should fail
	invalidStartOwnershipTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(2), End: sdkmath.NewUint(math.MaxUint64)},
	}

	mintMsg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(100),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: invalidStartOwnershipTimes,
					},
				},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), mintMsg)
	suite.Require().Error(err, "transfer with ownership times not starting at 1 should fail")
}

// TestNoCustomOwnershipTimes_EndNotMaxRejected tests that ownership times not ending at MaxUint64 fail
func (suite *NoCustomOwnershipTimesTestSuite) TestNoCustomOwnershipTimes_EndNotMaxRejected() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Setup mint approval
	suite.SetupMintApproval(collectionId)

	// Try with end < MaxUint64 - should fail
	invalidEndOwnershipTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64 - 1)},
	}

	mintMsg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(100),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: invalidEndOwnershipTimes,
					},
				},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), mintMsg)
	suite.Require().Error(err, "transfer with ownership times not ending at MaxUint64 should fail")
}

// TestNoCustomOwnershipTimes_AdjacentRangesMergedAndPass tests that adjacent ranges that merge to full range pass
// Note: The system merges adjacent ranges before checking invariants. If multiple adjacent ranges
// together form a full range [1, MaxUint64], they get merged and pass the invariant check.
func (suite *NoCustomOwnershipTimesTestSuite) TestNoCustomOwnershipTimes_AdjacentRangesMergedAndPass() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Setup mint approval
	suite.SetupMintApproval(collectionId)

	// Multiple adjacent ranges that together span full range
	// The system merges these to [1, MaxUint64] before invariant check
	multipleRangesOwnershipTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1000)},
		{Start: sdkmath.NewUint(1001), End: sdkmath.NewUint(math.MaxUint64)},
	}

	mintMsg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(100),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: multipleRangesOwnershipTimes,
					},
				},
			},
		},
	}

	// This should succeed because adjacent ranges are merged to a full range
	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), mintMsg)
	suite.Require().NoError(err, "transfer with adjacent ranges that merge to full range should succeed")
}

// TestNoCustomOwnershipTimes_InvariantDisabledAllowsCustomTimes tests that custom times work when invariant is disabled
func (suite *NoCustomOwnershipTimesTestSuite) TestNoCustomOwnershipTimes_InvariantDisabledAllowsCustomTimes() {
	// Create collection WITHOUT invariant enabled
	collectionId := suite.createCollectionWithInvariant(false)

	// Verify invariant is not set
	collection := suite.GetCollection(collectionId)
	if collection.Invariants != nil {
		suite.Require().False(collection.Invariants.NoCustomOwnershipTimes)
	}

	// Setup mint approval
	suite.SetupMintApproval(collectionId)

	// Transfer with restricted ownership times - should succeed
	restrictedOwnershipTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(200)},
	}

	mintMsg := &types.MsgTransferTokens{
		Creator:      suite.Manager,
		CollectionId: collectionId,
		Transfers: []*types.Transfer{
			{
				From:        types.MintAddress,
				ToAddresses: []string{suite.Alice},
				Balances: []*types.Balance{
					{
						Amount:         sdkmath.NewUint(100),
						TokenIds:       []*types.UintRange{{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}},
						OwnershipTimes: restrictedOwnershipTimes,
					},
				},
			},
		},
	}

	_, err := suite.MsgServer.TransferTokens(sdk.WrapSDKContext(suite.Ctx), mintMsg)
	suite.Require().NoError(err, "transfer with custom ownership times should succeed when invariant is disabled")
}

// TestNoCustomOwnershipTimes_CollectionApprovalOwnershipTimesValidated tests that collection approval ownership times are validated
func (suite *NoCustomOwnershipTimesTestSuite) TestNoCustomOwnershipTimes_CollectionApprovalOwnershipTimesValidated() {
	// Create collection with invariant enabled
	collectionId := suite.createCollectionWithInvariant(true)

	// Try to add a collection approval with restricted ownership times - should fail
	restrictedOwnershipTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(200)}, // Not full range
	}

	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              collectionId,
		UpdateCollectionApprovals: true,
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId:        "restricted_approval",
				FromListId:        "AllWithoutMint",
				ToListId:          "All",
				InitiatedByListId: "All",
				TransferTimes: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
				},
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
				},
				OwnershipTimes:   restrictedOwnershipTimes, // Invalid with invariant
				ApprovalCriteria: &types.ApprovalCriteria{},
				Version:          sdkmath.NewUint(0),
			},
		},
	}

	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().Error(err, "collection approval with restricted ownership times should fail")
}
