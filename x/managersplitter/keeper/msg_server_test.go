package keeper_test

import (
	"context"
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/keeper"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	bitbadgesapp "github.com/bitbadges/bitbadgeschain/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	app := bitbadgesapp.Setup(false)
	ctx := app.BaseApp.NewContext(false)
	return keeper.NewMsgServerImpl(app.ManagerSplitterKeeper), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
