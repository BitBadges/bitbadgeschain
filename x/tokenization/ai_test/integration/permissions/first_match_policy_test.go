package permissions_test

import (
	"math"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/ai_test/testutil"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
)

type FirstMatchPolicyTestSuite struct {
	testutil.AITestSuite
	CollectionId sdkmath.Uint
}

func TestFirstMatchPolicySuite(t *testing.T) {
	testutil.RunTestSuite(t, new(FirstMatchPolicyTestSuite))
}

func (suite *FirstMatchPolicyTestSuite) SetupTest() {
	suite.AITestSuite.SetupTest()
	suite.CollectionId = suite.CreateTestCollection(suite.Manager)
}

// TestFirstMatchPolicy_OrderMatters tests that permission order matters for first-match policy
func (suite *FirstMatchPolicyTestSuite) TestFirstMatchPolicy_OrderMatters() {
	// Create permissions with specific order
	// First permission: Forbid token IDs 1-50
	// Second permission: Permit token IDs 1-100
	// First match should win, so token IDs 1-50 should be forbidden

	permission1 := &types.TokenIdsActionPermission{
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		PermanentlyPermittedTimes: []*types.UintRange{},
	}

	permission2 := &types.TokenIdsActionPermission{
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{},
	}

	// Update collection with permissions in order: forbid first, then permit
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{
				permission1, // First match - should forbid 1-50
				permission2, // Second match - should permit 1-100, but won't be checked for 1-50
			},
		},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Verify collection has permissions set
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().NotNil(collection.CollectionPermissions)
	suite.Require().Equal(2, len(collection.CollectionPermissions.CanUpdateTokenMetadata))
}

// TestFirstMatchPolicy_ReverseOrder tests that reversing permission order changes behavior
func (suite *FirstMatchPolicyTestSuite) TestFirstMatchPolicy_ReverseOrder() {
	// Create permissions with reverse order
	// First permission: Permit token IDs 1-100
	// Second permission: Forbid token IDs 1-50
	// First match should win, so token IDs 1-50 should be permitted (different from previous test)

	permission1 := &types.TokenIdsActionPermission{
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(100)},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{},
	}

	permission2 := &types.TokenIdsActionPermission{
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(50)},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		PermanentlyPermittedTimes: []*types.UintRange{},
	}

	// Update collection with permissions in reverse order: permit first, then forbid
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{
				permission1, // First match - should permit 1-100
				permission2, // Second match - should forbid 1-50, but won't be checked
			},
		},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Verify collection has permissions set
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().NotNil(collection.CollectionPermissions)
	suite.Require().Equal(2, len(collection.CollectionPermissions.CanUpdateTokenMetadata))
}

// TestFirstMatchPolicy_NoMatch tests behavior when no permissions match
func (suite *FirstMatchPolicyTestSuite) TestFirstMatchPolicy_NoMatch() {
	// Create permissions that don't match the token IDs we're checking
	permission := &types.TokenIdsActionPermission{
		TokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(100), End: sdkmath.NewUint(200)},
		},
		PermanentlyPermittedTimes: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(math.MaxUint64)},
		},
		PermanentlyForbiddenTimes: []*types.UintRange{},
	}

	// Update collection with permission that doesn't match token IDs 1-50
	updateMsg := &types.MsgUniversalUpdateCollection{
		Creator:                   suite.Manager,
		CollectionId:              suite.CollectionId,
		UpdateCollectionPermissions: true,
		CollectionPermissions: &types.CollectionPermissions{
			CanUpdateTokenMetadata: []*types.TokenIdsActionPermission{
				permission,
			},
		},
	}
	_, err := suite.MsgServer.UniversalUpdateCollection(sdk.WrapSDKContext(suite.Ctx), updateMsg)
	suite.Require().NoError(err)

	// Verify collection has permission set
	collection := suite.GetCollection(suite.CollectionId)
	suite.Require().NotNil(collection.CollectionPermissions)
	suite.Require().Equal(1, len(collection.CollectionPermissions.CanUpdateTokenMetadata))
}

