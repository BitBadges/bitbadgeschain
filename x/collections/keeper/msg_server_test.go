package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/trevormil/bitbadgeschain/testutil/keeper"
	"github.com/trevormil/bitbadgeschain/x/collections/keeper"
	"github.com/trevormil/bitbadgeschain/x/collections/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.CollectionsKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
