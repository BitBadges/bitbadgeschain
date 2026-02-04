package keeper_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_MsgDeleteDynamicStore(t *testing.T) {
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

	msg := types.NewMsgDeleteDynamicStore(creator, resp.StoreId)
	_, err = suite.msgServer.DeleteDynamicStore(wctx, msg)
	require.NoError(t, err)

	msg = types.NewMsgDeleteDynamicStore("", resp.StoreId)
	_, err = suite.msgServer.DeleteDynamicStore(wctx, msg)
	require.Error(t, err)

	msg = types.NewMsgDeleteDynamicStore(creator, types.NewUintFromString("0"))
	_, err = suite.msgServer.DeleteDynamicStore(wctx, msg)
	require.Error(t, err)
}
