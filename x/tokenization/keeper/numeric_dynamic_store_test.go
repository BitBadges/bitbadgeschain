package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	approvalcriteria "github.com/bitbadges/bitbadgeschain/x/tokenization/approval_criteria"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestNumericDynamicStore_SetNumericValues tests setting various numeric values
func TestNumericDynamicStore_SetNumericValues(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	address := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"

	// Create store with default value 50
	msgCreate := types.NewMsgCreateDynamicStore(creator, sdkmath.NewUint(50))
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msgCreate)
	require.NoError(t, err)

	// Set value 0
	msg := types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, address, sdkmath.NewUint(0))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msg)
	require.NoError(t, err)

	val, found := suite.app.TokenizationKeeper.GetDynamicStoreValueFromStore(ctx, resp.StoreId, address)
	require.True(t, found)
	require.True(t, val.Value.IsZero())

	// Set value 1
	msg = types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, address, sdkmath.NewUint(1))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msg)
	require.NoError(t, err)

	val, found = suite.app.TokenizationKeeper.GetDynamicStoreValueFromStore(ctx, resp.StoreId, address)
	require.True(t, found)
	require.True(t, val.Value.Equal(sdkmath.NewUint(1)))

	// Set value 100
	msg = types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, address, sdkmath.NewUint(100))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msg)
	require.NoError(t, err)

	val, found = suite.app.TokenizationKeeper.GetDynamicStoreValueFromStore(ctx, resp.StoreId, address)
	require.True(t, found)
	require.True(t, val.Value.Equal(sdkmath.NewUint(100)))

	// Set max uint64 value
	maxVal := sdkmath.NewUint(^uint64(0))
	msg = types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, address, maxVal)
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msg)
	require.NoError(t, err)

	val, found = suite.app.TokenizationKeeper.GetDynamicStoreValueFromStore(ctx, resp.StoreId, address)
	require.True(t, found)
	require.True(t, val.Value.Equal(maxVal))

	// Verify default value returns for uninitialized address
	queryResp, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: resp.StoreId.String(),
		Address: "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme",
	})
	require.NoError(t, err)
	require.True(t, queryResp.Value.Value.Equal(sdkmath.NewUint(50)))
}

// TestNumericDynamicStore_ComparisonOperators tests all comparison operators
func TestNumericDynamicStore_ComparisonOperators(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"

	// Create store with default 0 and set initiator value to 50
	msgCreate := types.NewMsgCreateDynamicStore(creator, sdkmath.NewUint(0))
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msgCreate)
	require.NoError(t, err)

	msgSet := types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, initiator, sdkmath.NewUint(50))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msgSet)
	require.NoError(t, err)

	// Build a minimal collection for the checker
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
	}

	dynamicStoreService := &testDynamicStoreService{
		suite: suite,
	}

	tests := []struct {
		name     string
		operator string
		compVal  sdkmath.Uint
		pass     bool
	}{
		{"eq match", "eq", sdkmath.NewUint(50), true},
		{"eq no match", "eq", sdkmath.NewUint(49), false},
		{"ne match", "ne", sdkmath.NewUint(49), true},
		{"ne no match", "ne", sdkmath.NewUint(50), false},
		{"gt match", "gt", sdkmath.NewUint(49), true},
		{"gt boundary no match", "gt", sdkmath.NewUint(50), false},
		{"gte match equal", "gte", sdkmath.NewUint(50), true},
		{"gte match above", "gte", sdkmath.NewUint(49), true},
		{"gte no match", "gte", sdkmath.NewUint(51), false},
		{"lt match", "lt", sdkmath.NewUint(51), true},
		{"lt boundary no match", "lt", sdkmath.NewUint(50), false},
		{"lte match equal", "lte", sdkmath.NewUint(50), true},
		{"lte match below", "lte", sdkmath.NewUint(51), true},
		{"lte no match", "lte", sdkmath.NewUint(49), false},
		{"empty operator (legacy) - non-zero passes", "", sdkmath.NewUint(0), true},
		{"unknown operator fails", "invalid", sdkmath.NewUint(50), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			approval := &types.CollectionApproval{
				ApprovalCriteria: &types.ApprovalCriteria{
					DynamicStoreChallenges: []*types.DynamicStoreChallenge{
						{
							StoreId:             resp.StoreId,
							OwnershipCheckParty: "initiator",
							ComparisonOperator:  tc.operator,
							ComparisonValue:     tc.compVal,
						},
					},
				},
			}

			checker := approvalcriteria.NewDynamicStoreChallengesChecker(dynamicStoreService)
			detErr, err := checker.Check(ctx, approval, collection, "", "", initiator, "", "", nil, nil, "", false)
			if tc.pass {
				require.NoError(t, err, "expected pass for %s, got detErr: %s", tc.name, detErr)
			} else {
				require.Error(t, err, "expected fail for %s", tc.name)
			}
		})
	}
}

// TestNumericDynamicStore_LegacyBackwardCompat tests backward compatibility with boolean semantics
func TestNumericDynamicStore_LegacyBackwardCompat(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"

	dynamicStoreService := &testDynamicStoreService{
		suite: suite,
	}

	collection := &types.TokenCollection{CollectionId: sdkmath.NewUint(1)}

	// Create store with default 0 (was false)
	msgCreate := types.NewMsgCreateDynamicStore(creator, sdkmath.NewUint(0))
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msgCreate)
	require.NoError(t, err)

	// With empty operator (legacy), value=0 should fail (like old false)
	approval := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			DynamicStoreChallenges: []*types.DynamicStoreChallenge{
				{StoreId: resp.StoreId, OwnershipCheckParty: "initiator"},
			},
		},
	}

	checker := approvalcriteria.NewDynamicStoreChallengesChecker(dynamicStoreService)
	_, err = checker.Check(ctx, approval, collection, "", "", initiator, "", "", nil, nil, "", false)
	require.Error(t, err, "value 0 (old false) with legacy operator should fail")

	// Set value to 1 (was true) - should now pass with empty operator
	msgSet := types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, initiator, sdkmath.NewUint(1))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msgSet)
	require.NoError(t, err)

	_, err = checker.Check(ctx, approval, collection, "", "", initiator, "", "", nil, nil, "", false)
	require.NoError(t, err, "value 1 (old true) with legacy operator should pass")

	// Any non-zero value should pass with legacy operator
	msgSet = types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, initiator, sdkmath.NewUint(999))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msgSet)
	require.NoError(t, err)

	_, err = checker.Check(ctx, approval, collection, "", "", initiator, "", "", nil, nil, "", false)
	require.NoError(t, err, "any non-zero value with legacy operator should pass")
}

// TestNumericDynamicStore_ComparisonBoundary tests edge cases around comparison boundaries
func TestNumericDynamicStore_ComparisonBoundary(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"

	dynamicStoreService := &testDynamicStoreService{
		suite: suite,
	}

	collection := &types.TokenCollection{CollectionId: sdkmath.NewUint(1)}

	// Create store
	msgCreate := types.NewMsgCreateDynamicStore(creator, sdkmath.NewUint(0))
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msgCreate)
	require.NoError(t, err)

	// Test zero value with eq 0
	approval := &types.CollectionApproval{
		ApprovalCriteria: &types.ApprovalCriteria{
			DynamicStoreChallenges: []*types.DynamicStoreChallenge{
				{
					StoreId:             resp.StoreId,
					OwnershipCheckParty: "initiator",
					ComparisonOperator:  "eq",
					ComparisonValue:     sdkmath.NewUint(0),
				},
			},
		},
	}

	checker := approvalcriteria.NewDynamicStoreChallengesChecker(dynamicStoreService)
	// Default is 0, initiator has no value set, so uses default 0. eq 0 should pass.
	_, err = checker.Check(ctx, approval, collection, "", "", initiator, "", "", nil, nil, "", false)
	require.NoError(t, err, "zero value eq 0 should pass")

	// Set max uint64 and test gt (maxUint64 - 1)
	maxVal := sdkmath.NewUint(^uint64(0))
	msgSet := types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, initiator, maxVal)
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msgSet)
	require.NoError(t, err)

	approval.ApprovalCriteria.DynamicStoreChallenges[0].ComparisonOperator = "gt"
	approval.ApprovalCriteria.DynamicStoreChallenges[0].ComparisonValue = sdkmath.NewUint(^uint64(0) - 1)
	_, err = checker.Check(ctx, approval, collection, "", "", initiator, "", "", nil, nil, "", false)
	require.NoError(t, err, "max uint64 gt (max-1) should pass")
}

// testDynamicStoreService is a test adapter for the DynamicStoreService interface
type testDynamicStoreService struct {
	suite *TestSuite
}

func (s *testDynamicStoreService) GetDynamicStore(ctx sdk.Context, storeId sdkmath.Uint) (*types.DynamicStore, bool) {
	store, found := s.suite.app.TokenizationKeeper.GetDynamicStoreFromStore(ctx, storeId)
	if !found {
		return nil, false
	}
	return &store, true
}

func (s *testDynamicStoreService) GetDynamicStoreValue(ctx sdk.Context, storeId sdkmath.Uint, address string) (*types.DynamicStoreValue, bool) {
	value, found := s.suite.app.TokenizationKeeper.GetDynamicStoreValueFromStore(ctx, storeId, address)
	if !found {
		return nil, false
	}
	return &value, true
}
