package keeper_test

import (
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/badges/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_MsgCreateDynamicStore(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	msg := types.NewMsgCreateDynamicStore(creator, false)
	resp, err := suite.msgServer.CreateDynamicStore(wctx, msg)
	require.NoError(t, err)
	require.NotNil(t, resp)

	msg = types.NewMsgCreateDynamicStore("", false)
	_, err = suite.msgServer.CreateDynamicStore(wctx, msg)
	require.Error(t, err)
}
