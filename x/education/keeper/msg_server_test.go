package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/trevormil/bitbadgeschain/testutil/keeper"
	"github.com/trevormil/bitbadgeschain/x/education/keeper"
	"github.com/trevormil/bitbadgeschain/x/education/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.EducationKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
