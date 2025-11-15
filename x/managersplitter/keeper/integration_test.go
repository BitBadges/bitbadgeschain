package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/keeper"
	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bitbadgesapp "github.com/bitbadges/bitbadgeschain/app"

	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/testutil"
)

const (
	// Note these are alphanumerically sorted (needed for approvals test)
	alice   = "bb1e0w5t53nrq7p66fye6c8p0ynyhf6y24lke5430"
	bob     = "bb1jmjfq0tplp9tmx4v9uemw72y4d2wa5nrjmmk3q"
	charlie = "bb1xyxs3skf3f4jfqeuv89yyaqvjc6lffav9altme"
)

type TestSuite struct {
	suite.Suite

	app         *bitbadgesapp.App
	ctx         sdk.Context
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

// SetupTest initializes the test suite
func (suite *TestSuite) SetupTest() {
	app := bitbadgesapp.Setup(
		false,
	)

	ctx := app.BaseApp.NewContext(false)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.AppCodec().InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.ManagerSplitterKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.app = app
	suite.ctx = ctx
	suite.msgServer = keeper.NewMsgServerImpl(app.ManagerSplitterKeeper)
	suite.queryClient = queryClient

	bob_acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.MustAccAddressFromBech32(bob))
	alice_acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.MustAccAddressFromBech32(alice))
	charlie_acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, sdk.MustAccAddressFromBech32(charlie))

	suite.app.AccountKeeper.SetAccount(suite.ctx, bob_acc)
	suite.app.AccountKeeper.SetAccount(suite.ctx, alice_acc)
	suite.app.AccountKeeper.SetAccount(suite.ctx, charlie_acc)

	suite.ctx = suite.ctx.WithBlockTime(time.Now())

	ubadgeAmount := sdkmath.NewInt(100 * 1e9) // 100 BADGE

	banktestutil.FundAccount(suite.ctx, suite.app.BankKeeper, sdk.MustAccAddressFromBech32(bob), sdk.NewCoins(sdk.NewInt64Coin("ubadge", ubadgeAmount.Int64())))
	banktestutil.FundAccount(suite.ctx, suite.app.BankKeeper, sdk.MustAccAddressFromBech32(alice), sdk.NewCoins(sdk.NewInt64Coin("ubadge", ubadgeAmount.Int64())))
	banktestutil.FundAccount(suite.ctx, suite.app.BankKeeper, sdk.MustAccAddressFromBech32(charlie), sdk.NewCoins(sdk.NewInt64Coin("ubadge", ubadgeAmount.Int64())))
}

func TestManagerSplitterKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
