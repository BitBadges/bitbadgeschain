package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	approvalcriteria "github.com/bitbadges/bitbadgeschain/x/tokenization/approval_criteria"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestGlobalKillSwitch_CreateStore_DefaultsToEnabled tests that new stores default to globalEnabled = true
func TestGlobalKillSwitch_CreateStore_DefaultsToEnabled(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

	// Create a dynamic store
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: true,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Query the store to verify globalEnabled = true
	queryResp, err := suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: createResp.StoreId.String(),
	})
	require.NoError(t, err)
	require.NotNil(t, queryResp)
	require.NotNil(t, queryResp.Store)
	require.True(t, queryResp.Store.GlobalEnabled, "new stores should default to globalEnabled = true")
}

// TestGlobalKillSwitch_UpdateGlobalEnabled tests updating the globalEnabled field
func TestGlobalKillSwitch_UpdateGlobalEnabled(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

	// Create a dynamic store
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: true,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Verify it starts as enabled
	queryResp, err := suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: createResp.StoreId.String(),
	})
	require.NoError(t, err)
	require.True(t, queryResp.Store.GlobalEnabled, "store should start as enabled")

	// Disable the global kill switch
	updateMsg := &types.MsgUpdateDynamicStore{
		Creator:       creator,
		StoreId:      createResp.StoreId,
		DefaultValue: true, // Keep defaultValue unchanged
		GlobalEnabled: false,
	}

	_, err = suite.msgServer.UpdateDynamicStore(wctx, updateMsg)
	require.NoError(t, err)

	// Verify it's now disabled
	queryResp, err = suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: createResp.StoreId.String(),
	})
	require.NoError(t, err)
	require.False(t, queryResp.Store.GlobalEnabled, "store should be disabled after update")

	// Re-enable the global kill switch
	updateMsg.GlobalEnabled = true
	_, err = suite.msgServer.UpdateDynamicStore(wctx, updateMsg)
	require.NoError(t, err)

	// Verify it's enabled again
	queryResp, err = suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: createResp.StoreId.String(),
	})
	require.NoError(t, err)
	require.True(t, queryResp.Store.GlobalEnabled, "store should be enabled after re-enabling")
}

// TestGlobalKillSwitch_BlocksApprovalsWhenDisabled tests that disabled kill switch blocks approvals
func TestGlobalKillSwitch_BlocksApprovalsWhenDisabled(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"

	// Create a dynamic store
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: true,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Set a value for the initiator to true (so per-address check would pass)
	setValueMsg := &types.MsgSetDynamicStoreValue{
		Creator: creator,
		StoreId: createResp.StoreId,
		Address: initiator,
		Value:   true,
	}
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	// Create a collection with an approval that uses this dynamic store
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId: "1",
				ApprovalCriteria: &types.ApprovalCriteria{
					DynamicStoreChallenges: []*types.DynamicStoreChallenge{
						{
							StoreId: createResp.StoreId,
						},
					},
				},
			},
		},
	}

	// Store the collection
	err = suite.app.TokenizationKeeper.SetCollectionInStore(ctx, collection, false)
	require.NoError(t, err)

	// Test approval with globalEnabled = true (should pass per-address check)
	approval := collection.CollectionApprovals[0]
	checkers := suite.app.TokenizationKeeper.GetApprovalCriteriaCheckers(approval)
	require.NotEmpty(t, checkers)

	// Find the DynamicStoreChallengesChecker
	var dynamicStoreChecker approvalcriteria.ApprovalCriteriaChecker
	for _, checker := range checkers {
		if checker.Name() == "DynamicStoreChallenges" {
			dynamicStoreChecker = checker
			break
		}
	}
	require.NotNil(t, dynamicStoreChecker, "should have DynamicStoreChallengesChecker")

	// Check with globalEnabled = true (should pass)
	detErrMsg, err := dynamicStoreChecker.Check(ctx, approval, collection, "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme", "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme", initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.NoError(t, err, "approval should pass when globalEnabled = true and initiator has permission")

	// Disable the global kill switch
	updateMsg := &types.MsgUpdateDynamicStore{
		Creator:       creator,
		StoreId:      createResp.StoreId,
		DefaultValue: true,
		GlobalEnabled: false,
	}
	_, err = suite.msgServer.UpdateDynamicStore(wctx, updateMsg)
	require.NoError(t, err)

	// Check with globalEnabled = false (should fail immediately)
	detErrMsg, err = dynamicStoreChecker.Check(ctx, approval, collection, "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme", "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme", initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.Error(t, err, "approval should fail when globalEnabled = false")
	require.Contains(t, err.Error(), "globally disabled", "error should mention globally disabled")
	require.Contains(t, detErrMsg, "globally disabled", "deterministic error message should mention globally disabled")
}

// TestGlobalKillSwitch_AllowsApprovalsWhenEnabled tests that enabled kill switch allows per-address logic
func TestGlobalKillSwitch_AllowsApprovalsWhenEnabled(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	initiator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"

	// Create a dynamic store with defaultValue = false
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: false,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Create a collection with an approval that uses this dynamic store
	collection := &types.TokenCollection{
		CollectionId: sdkmath.NewUint(1),
		ValidTokenIds: []*types.UintRange{
			{Start: sdkmath.NewUint(1), End: sdkmath.NewUint(1)},
		},
		CollectionApprovals: []*types.CollectionApproval{
			{
				ApprovalId: "1",
				ApprovalCriteria: &types.ApprovalCriteria{
					DynamicStoreChallenges: []*types.DynamicStoreChallenge{
						{
							StoreId: createResp.StoreId,
						},
					},
				},
			},
		},
	}

	// Store the collection
	err = suite.app.TokenizationKeeper.SetCollectionInStore(ctx, collection, false)
	require.NoError(t, err)

	approval := collection.CollectionApprovals[0]
	checkers := suite.app.TokenizationKeeper.GetApprovalCriteriaCheckers(approval)
	require.NotEmpty(t, checkers)

	// Find the DynamicStoreChallengesChecker
	var dynamicStoreChecker approvalcriteria.ApprovalCriteriaChecker
	for _, checker := range checkers {
		if checker.Name() == "DynamicStoreChallenges" {
			dynamicStoreChecker = checker
			break
		}
	}
	require.NotNil(t, dynamicStoreChecker, "should have DynamicStoreChallengesChecker")

	// Test with globalEnabled = true but initiator has no value (should use defaultValue = false, fail)
	_, err = dynamicStoreChecker.Check(ctx, approval, collection, "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme", "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme", initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.Error(t, err, "should fail when defaultValue = false and no per-address value")
	require.Contains(t, err.Error(), "does not have permission", "error should mention permission")

	// Set value for initiator to true
	setValueMsg := &types.MsgSetDynamicStoreValue{
		Creator: creator,
		StoreId: createResp.StoreId,
		Address: initiator,
		Value:   true,
	}
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, setValueMsg)
	require.NoError(t, err)

	// Now should pass (globalEnabled = true, initiator has value = true)
	_, err = dynamicStoreChecker.Check(ctx, approval, collection, "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme", "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme", initiator, "collection", "", []*types.MerkleProof{}, []*types.ETHSignatureProof{}, "", false)
	require.NoError(t, err, "should pass when globalEnabled = true and initiator has permission")
}

// TestGlobalKillSwitch_OnlyCreatorCanUpdate tests that only the creator can update globalEnabled
func TestGlobalKillSwitch_OnlyCreatorCanUpdate(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	wrongCreator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"

	// Create a dynamic store
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: true,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Try to update with wrong creator
	updateMsg := &types.MsgUpdateDynamicStore{
		Creator:       wrongCreator,
		StoreId:      createResp.StoreId,
		DefaultValue: true,
		GlobalEnabled: false,
	}

	_, err = suite.msgServer.UpdateDynamicStore(wctx, updateMsg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Only the creator can update the dynamic store")
}

// TestGlobalKillSwitch_UpdatePreservesDefaultValue tests that updating globalEnabled doesn't affect defaultValue
func TestGlobalKillSwitch_UpdatePreservesDefaultValue(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

	// Create a dynamic store with defaultValue = true
	createMsg := &types.MsgCreateDynamicStore{
		Creator:      creator,
		DefaultValue: true,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	// Update globalEnabled but keep defaultValue = true
	updateMsg := &types.MsgUpdateDynamicStore{
		Creator:       creator,
		StoreId:      createResp.StoreId,
		DefaultValue: true,
		GlobalEnabled: false,
	}

	_, err = suite.msgServer.UpdateDynamicStore(wctx, updateMsg)
	require.NoError(t, err)

	// Verify defaultValue is still true
	queryResp, err := suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: createResp.StoreId.String(),
	})
	require.NoError(t, err)
	require.True(t, queryResp.Store.DefaultValue, "defaultValue should remain unchanged")
	require.False(t, queryResp.Store.GlobalEnabled, "globalEnabled should be updated")
}

