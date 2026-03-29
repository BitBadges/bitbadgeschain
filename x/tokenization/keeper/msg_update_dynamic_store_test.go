package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_MsgUpdateDynamicStore(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	msgCreate := types.NewMsgCreateDynamicStore(creator, sdkmath.NewUint(0))
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msgCreate)
	require.NoError(t, err)
	require.NotNil(t, resp)

	msg := types.NewMsgUpdateDynamicStore(creator, resp.StoreId, sdkmath.NewUint(0))
	_, err = suite.msgServer.UpdateDynamicStore(wctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgUpdateDynamicStore("", resp.StoreId, sdkmath.NewUint(0))
	_, err = suite.msgServer.UpdateDynamicStore(wctx, msg)
	require.Error(t, err)

	msg = types.NewMsgUpdateDynamicStore(creator, types.NewUintFromString("0"), sdkmath.NewUint(0))
	_, err = suite.msgServer.UpdateDynamicStore(wctx, msg)
	require.Error(t, err)

	// Test updating to a numeric value
	msg = types.NewMsgUpdateDynamicStore(creator, resp.StoreId, sdkmath.NewUint(42))
	_, err = suite.msgServer.UpdateDynamicStore(wctx, msg)
	require.NoError(t, err)

	store, found := suite.app.TokenizationKeeper.GetDynamicStoreFromStore(ctx, resp.StoreId)
	require.True(t, found)
	require.True(t, store.DefaultValue.Equal(sdkmath.NewUint(42)))
}
