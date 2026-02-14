package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestDynamicStoreFlow(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	// Test creating a dynamic store
	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

	createMsg := &types.MsgCreateDynamicStore{
		Creator: creator,
	}

	createResp, err := suite.msgServer.CreateDynamicStore(wctx, createMsg)
	require.NoError(t, err)
	require.NotNil(t, createResp)
	require.Equal(t, sdkmath.NewUint(1), createResp.StoreId)

	// Test querying the dynamic store
	queryResp, err := suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: "1",
	})
	require.NoError(t, err)
	require.NotNil(t, queryResp)
	require.NotNil(t, queryResp.Store)
	require.Equal(t, sdkmath.NewUint(1), queryResp.Store.StoreId)
	require.Equal(t, creator, queryResp.Store.CreatedBy)

	// Test updating the dynamic store (no-op since no data/metadata)
	updateMsg := &types.MsgUpdateDynamicStore{
		Creator: creator,
		StoreId: sdkmath.NewUint(1),
	}

	updateResp, err := suite.msgServer.UpdateDynamicStore(wctx, updateMsg)
	require.NoError(t, err)
	require.NotNil(t, updateResp)

	// Test that only creator can update
	wrongCreator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	updateMsgWrongCreator := &types.MsgUpdateDynamicStore{
		Creator: wrongCreator,
		StoreId: sdkmath.NewUint(1),
	}

	_, err = suite.msgServer.UpdateDynamicStore(wctx, updateMsgWrongCreator)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Only the creator can update the dynamic store")

	// Test deleting the dynamic store
	deleteMsg := &types.MsgDeleteDynamicStore{
		Creator: creator,
		StoreId: sdkmath.NewUint(1),
	}

	deleteResp, err := suite.msgServer.DeleteDynamicStore(wctx, deleteMsg)
	require.NoError(t, err)
	require.NotNil(t, deleteResp)

	// Verify the deletion
	_, err = suite.app.TokenizationKeeper.GetDynamicStore(wctx, &types.QueryGetDynamicStoreRequest{
		StoreId: "1",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "dynamic store not found")

	// Test that only creator can delete (create a new store first)
	createMsg2 := &types.MsgCreateDynamicStore{
		Creator: creator,
	}

	createResp2, err := suite.msgServer.CreateDynamicStore(wctx, createMsg2)
	require.NoError(t, err)
	require.NotNil(t, createResp2)

	deleteMsgWrongCreator := &types.MsgDeleteDynamicStore{
		Creator: wrongCreator,
		StoreId: createResp2.StoreId,
	}

	_, err = suite.msgServer.DeleteDynamicStore(wctx, deleteMsgWrongCreator)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Only the creator can delete the dynamic store")
}
