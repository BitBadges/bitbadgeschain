package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestCreateDynamicStore_SequentialIds verifies that stores get sequential IDs
func TestCreateDynamicStore_SequentialIds(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

	// Create first store - should get ID 1
	msg1 := types.NewMsgCreateDynamicStore(creator, false)
	resp1, err := suite.msgServer.CreateDynamicStore(wctx, msg1)
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Equal(t, sdkmath.NewUint(1), resp1.StoreId, "First store should get ID 1")

	// Verify next ID was incremented
	nextID := suite.app.BadgesKeeper.GetNextDynamicStoreId(ctx)
	require.Equal(t, sdkmath.NewUint(2), nextID, "Next ID should be 2 after creating first store")

	// Create second store - should get ID 2
	msg2 := types.NewMsgCreateDynamicStore(creator, false)
	resp2, err := suite.msgServer.CreateDynamicStore(wctx, msg2)
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Equal(t, sdkmath.NewUint(2), resp2.StoreId, "Second store should get ID 2")

	// Verify next ID was incremented again
	nextID = suite.app.BadgesKeeper.GetNextDynamicStoreId(ctx)
	require.Equal(t, sdkmath.NewUint(3), nextID, "Next ID should be 3 after creating second store")

	// Create third store - should get ID 3
	msg3 := types.NewMsgCreateDynamicStore(creator, false)
	resp3, err := suite.msgServer.CreateDynamicStore(wctx, msg3)
	require.NoError(t, err)
	require.NotNil(t, resp3)
	require.Equal(t, sdkmath.NewUint(3), resp3.StoreId, "Third store should get ID 3")

	// Verify all stores exist with correct IDs
	store1, found1 := suite.app.BadgesKeeper.GetDynamicStoreFromStore(ctx, sdkmath.NewUint(1))
	require.True(t, found1, "Store 1 should exist")
	require.Equal(t, sdkmath.NewUint(1), store1.StoreId)

	store2, found2 := suite.app.BadgesKeeper.GetDynamicStoreFromStore(ctx, sdkmath.NewUint(2))
	require.True(t, found2, "Store 2 should exist")
	require.Equal(t, sdkmath.NewUint(2), store2.StoreId)

	store3, found3 := suite.app.BadgesKeeper.GetDynamicStoreFromStore(ctx, sdkmath.NewUint(3))
	require.True(t, found3, "Store 3 should exist")
	require.Equal(t, sdkmath.NewUint(3), store3.StoreId)
}

// TestCreateDynamicStore_UninitializedId verifies that the handler correctly handles
// uninitialized next ID (when it's 0). This tests the edge case where the store
// hasn't been initialized via genesis or test setup.
func TestCreateDynamicStore_UninitializedId(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Manually reset the next ID to 0 to simulate uninitialized state
	suite.app.BadgesKeeper.SetNextDynamicStoreId(ctx, sdkmath.NewUint(0))

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

	// Create first store - should get ID 1 (from the handler's check)
	msg := types.NewMsgCreateDynamicStore(creator, false)
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, sdkmath.NewUint(1), resp.StoreId, "First store should get ID 1 even when starting from 0")

	// Verify next ID was properly incremented to 2 (not 1!)
	nextID := suite.app.BadgesKeeper.GetNextDynamicStoreId(ctx)
	require.Equal(t, sdkmath.NewUint(2), nextID, "Next ID should be 2 after creating first store from uninitialized state")

	// Create second store - should get ID 2
	msg2 := types.NewMsgCreateDynamicStore(creator, false)
	resp2, err := suite.msgServer.CreateDynamicStore(wctx, msg2)
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Equal(t, sdkmath.NewUint(2), resp2.StoreId, "Second store should get ID 2")

	// Verify next ID is now 3
	nextID = suite.app.BadgesKeeper.GetNextDynamicStoreId(ctx)
	require.Equal(t, sdkmath.NewUint(3), nextID, "Next ID should be 3 after creating second store")
}

// TestCreateDynamicStore_WithPreInitializedId verifies that the handler correctly
// handles when the next ID is already initialized to a value > 1
func TestCreateDynamicStore_WithPreInitializedId(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Manually set next ID to 5 (simulating existing stores)
	suite.app.BadgesKeeper.SetNextDynamicStoreId(ctx, sdkmath.NewUint(5))

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

	// Create store - should get ID 5 (not 1!)
	msg := types.NewMsgCreateDynamicStore(creator, false)
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, sdkmath.NewUint(5), resp.StoreId, "Store should get ID 5 when next ID is pre-initialized to 5")

	// Verify next ID was incremented to 6
	nextID := suite.app.BadgesKeeper.GetNextDynamicStoreId(ctx)
	require.Equal(t, sdkmath.NewUint(6), nextID, "Next ID should be 6 after creating store with pre-initialized ID 5")
}

