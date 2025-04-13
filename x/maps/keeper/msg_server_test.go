package keeper_test

import (
	"context"
	"testing"

	"github.com/bitbadges/bitbadgeschain/x/maps/keeper"
	"github.com/bitbadges/bitbadgeschain/x/maps/types"

	keepertest "github.com/bitbadges/bitbadgeschain/testutil/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.MapsKeeper(t)
	return keeper.NewMsgServerImpl(k), sdk.WrapSDKContext(ctx)
}

func TestMsgServer(t *testing.T) {
	ms, ctx := setupMsgServer(t)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}
