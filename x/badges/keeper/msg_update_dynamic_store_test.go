package keeper_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
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
	msgCreate := types.NewMsgCreateDynamicStore(creator, false)
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msgCreate)
	require.NoError(t, err)
	require.NotNil(t, resp)

	msg := types.NewMsgUpdateDynamicStore(creator, resp.StoreId, false)
	_, err = suite.msgServer.UpdateDynamicStore(wctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgUpdateDynamicStore("", resp.StoreId, false)
	_, err = suite.msgServer.UpdateDynamicStore(wctx, msg)
	require.Error(t, err)

	msg = types.NewMsgUpdateDynamicStore(creator, types.NewUintFromString("0"), false)
	_, err = suite.msgServer.UpdateDynamicStore(wctx, msg)
	require.Error(t, err)
}
