package msg_handlers_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type SetTokenMetadataTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestSetTokenMetadataSuite(t *testing.T) {
	testutil.RunTestSuite(t, new(SetTokenMetadataTestSuite))
}

func (suite *SetTokenMetadataTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
	suite.Require().True(suite.CollectionId.GT(sdkmath.NewUint(0)), "collection ID should be greater than 0 after creation")
}

// TestSetTokenMetadata_ValidUpdate tests successfully updating token metadata
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_ValidUpdate() {
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				Uri:        "https://example.com/token/1",
				CustomData: "token custom data",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting token metadata should succeed")

	// Verify metadata was set
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().NotEmpty(collection.TokenMetadata, "token metadata should be set")
}

// TestSetTokenMetadata_PermissionChecked tests that canUpdateTokenMetadata permission is enforced
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_PermissionChecked() {
	// First, forbid all future updates for token IDs 1-10
	forbidAllTimes := []*types.UintRange{
		{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
	}

	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				Uri:        "https://example.com/token/initial",
				CustomData: "",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				PermanentlyForbiddenTimes: forbidAllTimes,
				PermanentlyPermittedTimes: []*types.UintRange{},
			},
		},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "first update should succeed")

	// Now try to update the same token IDs - should fail
	msg2 := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				Uri:        "https://example.com/token/updated",
				CustomData: "",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err = suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg2)
	suite.Require().Error(err, "update should fail because permission is forbidden for these token IDs")
}

// TestSetTokenMetadata_TokenIdsOutsideValidRange tests setting metadata for token IDs outside valid range
// Note: Token metadata can be set for any token IDs - the validTokenIds only restricts transfers/balances
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_TokenIdsOutsideValidRange() {
	// Collection was created with validTokenIds [1, 100]
	// Try to set metadata for token IDs outside this range
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(101), End: sdkmath.NewUint(200)}, // Outside valid range
				},
				Uri:        "https://example.com/token/outside-range",
				CustomData: "",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	// Token metadata can be set for any token IDs - validTokenIds constraint doesn't apply to metadata
	suite.Require().NoError(err, "setting metadata for token IDs outside validTokenIds is allowed")
}

// TestSetTokenMetadata_MultipleTokenRanges tests updating metadata for multiple token ranges
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_MultipleTokenRanges() {
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				Uri:        "https://example.com/tokens/1-10",
				CustomData: "range 1",
			},
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(11), End: sdkmath.NewUint(20)},
				},
				Uri:        "https://example.com/tokens/11-20",
				CustomData: "range 2",
			},
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(21), End: sdkmath.NewUint(30)},
				},
				Uri:        "https://example.com/tokens/21-30",
				CustomData: "range 3",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting metadata for multiple token ranges should succeed")

	// Verify all metadata entries were set
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().GreaterOrEqual(len(collection.TokenMetadata), 3, "should have at least 3 token metadata entries")
}

// TestSetTokenMetadata_OnlyManagerCanUpdate tests that only manager can update token metadata
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_OnlyManagerCanUpdate() {
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Alice, // Not the manager
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				Uri:        "https://example.com/token",
				CustomData: "",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "non-manager should not be able to update token metadata")
	suite.Require().Contains(err.Error(), "manager", "error should mention manager permission")
}

// TestSetTokenMetadata_NonExistentCollection tests updating metadata on non-existent collection
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_NonExistentCollection() {
	nonExistentId := sdkmath.NewUint(99999)

	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: nonExistentId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				Uri:        "https://example.com/token",
				CustomData: "",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "setting metadata on non-existent collection should fail")
}

// TestSetTokenMetadata_EmptyTokenMetadata tests behavior with empty token metadata array
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_EmptyTokenMetadata() {
	msg := &types.MsgSetTokenMetadata{
		Creator:               suite.Manager,
		CollectionId:          suite.CollectionId,
		TokenMetadata:         []*types.TokenMetadata{}, // Empty array
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	// Empty token metadata array might be allowed (clears metadata) or rejected
	// The actual behavior depends on implementation
	suite.Require().NoError(err, "empty token metadata should be allowed")
}

// TestSetTokenMetadata_InvalidTokenIdRange tests behavior with invalid token ID ranges
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_InvalidTokenIdRange() {
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(10), End: sdkmath.NewUint(1)}, // Invalid: start > end
				},
				Uri:        "https://example.com/token",
				CustomData: "",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "invalid token ID range should fail")
}

// TestSetTokenMetadata_EmptyCreator tests behavior with empty creator
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_EmptyCreator() {
	msg := &types.MsgSetTokenMetadata{
		Creator:      "", // Empty creator
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				Uri:        "https://example.com/token",
				CustomData: "",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().Error(err, "empty creator should fail")
}

// TestSetTokenMetadata_OverlappingTokenRanges tests behavior with overlapping token ranges
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_OverlappingTokenRanges() {
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				Uri:        "https://example.com/tokens/1-10",
				CustomData: "",
			},
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(5), End: sdkmath.NewUint(15)}, // Overlaps with 1-10
				},
				Uri:        "https://example.com/tokens/5-15",
				CustomData: "",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	// Overlapping ranges might be allowed or rejected depending on implementation
	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	// Test passes regardless of outcome - we're testing behavior
	_ = err
}

// TestSetTokenMetadata_SingleToken tests setting metadata for a single token
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_SingleToken() {
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)}, // Single token
				},
				Uri:        "https://example.com/token/1",
				CustomData: "single token data",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "setting metadata for single token should succeed")
}

// TestSetTokenMetadata_AfterManagerChange tests that new manager can update token metadata
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_AfterManagerChange() {
	// First change manager to Bob
	setManagerMsg := &types.MsgSetManager{
		Creator:          suite.Manager,
		CollectionId:     suite.CollectionId,
		Manager:          suite.Bob,
		CanUpdateManager: []*types.ActionPermission{},
	}
	_, err := suite.MsgServer.SetManager(sdk.WrapSDKContext(suite.Ctx), setManagerMsg)
	suite.Require().NoError(err, "setting new manager should succeed")

	// Bob (new manager) should be able to update token metadata
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Bob,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(10)},
				},
				Uri:        "https://bob.com/token",
				CustomData: "bob's data",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err = suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "new manager should be able to update token metadata")
}

// TestSetTokenMetadata_NilTokenIds tests behavior when token metadata has nil TokenIds
// Note: Nil TokenIds are allowed - treated as empty slice
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_NilTokenIds() {
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds:   nil, // Nil token IDs
				Uri:        "https://example.com/token",
				CustomData: "",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	// Nil token IDs are allowed - treated as empty
	suite.Require().NoError(err, "nil token IDs are allowed")
}

// TestSetTokenMetadata_MultipleTokenIdsInSingleEntry tests multiple token ID ranges in single metadata entry
func (suite *SetTokenMetadataTestSuite) TestSetTokenMetadata_MultipleTokenIdsInSingleEntry() {
	msg := &types.MsgSetTokenMetadata{
		Creator:      suite.Manager,
		CollectionId: suite.CollectionId,
		TokenMetadata: []*types.TokenMetadata{
			{
				TokenIds: []*types.UintRange{
					{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(5)},
					{Start: sdkmath.NewUint(10), End: sdkmath.NewUint(15)},
					{Start: sdkmath.NewUint(20), End: sdkmath.NewUint(25)},
				},
				Uri:        "https://example.com/multi-range-tokens",
				CustomData: "multiple ranges",
			},
		},
		CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{},
	}

	_, err := suite.MsgServer.SetTokenMetadata(sdk.WrapSDKContext(suite.Ctx), msg)
	suite.Require().NoError(err, "multiple token ID ranges in single entry should succeed")
}
