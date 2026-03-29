package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_MsgSetDynamicStoreValue(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := suite.app.AccountKeeper.GetModuleAddress("badges").String()
	address := suite.app.AccountKeeper.GetModuleAddress("badges").String()
	if creator == "" {
		creator = "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q" // fallback to a valid address
		address = creator
	}

	// Create a dynamic store with defaultValue = 1 (was true)
	msgCreate := types.NewMsgCreateDynamicStore(creator, sdkmath.NewUint(1))
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msgCreate)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Query an uninitialized address, should return defaultValue = 1
	queryResp, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: resp.StoreId.String(),
		Address: "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme", // some other valid address
	})
	require.NoError(t, err)
	require.NotNil(t, queryResp)
	require.True(t, queryResp.Value.Value.Equal(sdkmath.NewUint(1)))

	// Set a value for the address to 0 (was false)
	msg := types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, address, sdkmath.NewUint(0))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msg)
	require.NoError(t, err)

	// Query the address, should return 0
	queryResp2, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: resp.StoreId.String(),
		Address: address,
	})
	require.NoError(t, err)
	require.NotNil(t, queryResp2)
	require.True(t, queryResp2.Value.Value.IsZero())

	// Update the defaultValue to 0 (was false)
	msgUpdate := types.NewMsgUpdateDynamicStore(creator, resp.StoreId, sdkmath.NewUint(0))
	_, err = suite.msgServer.UpdateDynamicStore(wctx, msgUpdate)
	require.NoError(t, err)

	// Query a new uninitialized address, should now return 0
	queryResp3, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: resp.StoreId.String(),
		Address: "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme",
	})
	require.NoError(t, err)
	require.NotNil(t, queryResp3)
	require.True(t, queryResp3.Value.Value.IsZero())

	// Set a numeric value > 1
	msg = types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, address, sdkmath.NewUint(100))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msg)
	require.NoError(t, err)

	queryResp4, err := suite.app.TokenizationKeeper.GetDynamicStoreValue(wctx, &types.QueryGetDynamicStoreValueRequest{
		StoreId: resp.StoreId.String(),
		Address: address,
	})
	require.NoError(t, err)
	require.True(t, queryResp4.Value.Value.Equal(sdkmath.NewUint(100)))

	// Test that only creator can set values
	wrongCreator := "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	msgWrongCreator := types.NewMsgSetDynamicStoreValue(wrongCreator, resp.StoreId, address, sdkmath.NewUint(0))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msgWrongCreator)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Only the creator can set values in the dynamic store")

	// Test invalid address
	msgInvalidAddress := types.NewMsgSetDynamicStoreValue(creator, resp.StoreId, "invalid-address", sdkmath.NewUint(0))
	_, err = suite.msgServer.SetDynamicStoreValue(wctx, msgInvalidAddress)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid address")
}
