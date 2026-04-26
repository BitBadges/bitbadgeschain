package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/x/tokenization/types"
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

// Regression: when the next-id counter is zero (uninitialized), the first
// create must assign id 1 AND advance the counter to 2 so the second create
// gets id 2 — not id 1, which would silently overwrite the first store.
func TestKeeper_MsgCreateDynamicStore_AssignsSequentialIdsFromZero(t *testing.T) {
	suite := new(TestSuite)
	suite.SetT(t)
	suite.SetupTest()
	ctx := suite.ctx
	wctx := sdk.WrapSDKContext(ctx)

	suite.app.TokenizationKeeper.SetNextDynamicStoreId(ctx, sdkmath.NewUint(0))

	creator := "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"

	first, err := suite.msgServer.CreateDynamicStore(wctx, types.NewMsgCreateDynamicStore(creator, false))
	require.NoError(t, err)
	require.Equal(t, sdkmath.NewUint(1), first.StoreId)

	second, err := suite.msgServer.CreateDynamicStore(wctx, types.NewMsgCreateDynamicStore(creator, false))
	require.NoError(t, err)
	require.Equal(t, sdkmath.NewUint(2), second.StoreId, "second create must get a new id, not clobber id 1")

	third, err := suite.msgServer.CreateDynamicStore(wctx, types.NewMsgCreateDynamicStore(creator, false))
	require.NoError(t, err)
	require.Equal(t, sdkmath.NewUint(3), third.StoreId)

	got, found := suite.app.TokenizationKeeper.GetDynamicStoreFromStore(ctx, sdkmath.NewUint(1))
	require.True(t, found)
	require.Equal(t, creator, got.CreatedBy)
}
