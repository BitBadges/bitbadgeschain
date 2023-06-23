package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/bitbadges/bitbadgeschain/x/badges/types"
    "github.com/bitbadges/bitbadgeschain/x/badges/keeper"
    keepertest "github.com/bitbadges/bitbadgeschain/testutil/keeper"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.BadgesKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
