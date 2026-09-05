package keeper_test

import (
	"math"
	"time"

	sdkmath "cosmossdk.io/math"

	ratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

const (
	windowTestChannel = "channel-0"
	windowTestDenom   = "uatom"
	windowTestAddr    = "cosmos1windowtest"
)

var windowTestT0 = time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)

func (suite *KeeperTestSuite) debitHourFlow(amount int64) {
	suite.keeper.ResetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, 1)
	flow, _ := suite.keeper.GetChannelFlowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, 1)
	flow.NetFlow = flow.NetFlow.Sub(sdkmath.NewInt(amount))
	suite.keeper.SetChannelFlowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, 1, flow)
}

func (suite *KeeperTestSuite) hourFlow() sdkmath.Int {
	flow, _ := suite.keeper.GetChannelFlowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, 1)
	return flow.NetFlow
}

// TestHourWindowStartsAtBlockTime checks that an HOUR window is anchored to
// the block timestamp, not the block height.
func (suite *KeeperTestSuite) TestHourWindowStartsAtBlockTime() {
	suite.ctx = suite.ctx.WithBlockHeight(100).WithBlockTime(windowTestT0)
	suite.debitHourFlow(500)

	window, found := suite.keeper.GetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, 1)
	suite.Require().True(found)
	suite.Require().Equal(windowTestT0.Unix(), window.WindowStart, "window start must be the block time in unix seconds")
	suite.Require().Equal(int64(3600), window.WindowDuration, "window duration must be in seconds")
}

// TestHourWindowExpiresByBlockTime checks that an HOUR window rolls over once
// the block time has advanced by an hour, however few blocks that took.
func (suite *KeeperTestSuite) TestHourWindowExpiresByBlockTime() {
	suite.ctx = suite.ctx.WithBlockHeight(100).WithBlockTime(windowTestT0)
	suite.debitHourFlow(500)

	// One block later but one hour of wall-clock time: window expired.
	suite.ctx = suite.ctx.WithBlockHeight(101).WithBlockTime(windowTestT0.Add(time.Hour))
	suite.keeper.ResetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, 1)
	suite.Require().True(suite.hourFlow().IsZero(), "flow must reset after an hour of block time, got %s", suite.hourFlow())
}

// TestHourWindowDoesNotExpireByHeightAlone checks that fast blocks do not
// shorten an HOUR window.
func (suite *KeeperTestSuite) TestHourWindowDoesNotExpireByHeightAlone() {
	suite.ctx = suite.ctx.WithBlockHeight(100).WithBlockTime(windowTestT0)
	suite.debitHourFlow(500)

	// 5000 blocks later but only ten minutes of block time: still the same window.
	suite.ctx = suite.ctx.WithBlockHeight(5100).WithBlockTime(windowTestT0.Add(10 * time.Minute))
	suite.keeper.ResetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, 1)
	suite.Require().True(suite.hourFlow().Equal(sdkmath.NewInt(-500)), "flow must survive until an hour of block time has passed, got %s", suite.hourFlow())
}

// TestBlockWindowUsesHeight checks that BLOCK windows keep height semantics.
func (suite *KeeperTestSuite) TestBlockWindowUsesHeight() {
	suite.ctx = suite.ctx.WithBlockHeight(100).WithBlockTime(windowTestT0)
	suite.keeper.ResetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 10)
	window, _ := suite.keeper.GetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 10)
	suite.Require().Equal(int64(100), window.WindowStart)
	suite.Require().Equal(int64(10), window.WindowDuration)
	suite.keeper.SetChannelFlowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 10, ratelimittypes.ChannelFlow{NetFlow: sdkmath.NewInt(-5)})

	// A day of block time but only 9 blocks: not expired.
	suite.ctx = suite.ctx.WithBlockHeight(109).WithBlockTime(windowTestT0.Add(24 * time.Hour))
	suite.keeper.ResetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 10)
	flow, _ := suite.keeper.GetChannelFlowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 10)
	suite.Require().True(flow.NetFlow.Equal(sdkmath.NewInt(-5)))

	// Tenth block: expired.
	suite.ctx = suite.ctx.WithBlockHeight(110)
	suite.keeper.ResetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 10)
	flow, _ = suite.keeper.GetChannelFlowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 10)
	suite.Require().True(flow.NetFlow.IsZero())
}

// TestTimeframeDurationInSeconds checks the HOUR/DAY conversion.
func (suite *KeeperTestSuite) TestTimeframeDurationInSeconds() {
	suite.Require().Equal(int64(7200), ratelimittypes.TimeframeDurationInSeconds(ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, 2))
	suite.Require().Equal(int64(172800), ratelimittypes.TimeframeDurationInSeconds(ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY, 2))
}

// TestTimeframeDurationOverflowRejected checks that durations whose second
// count does not fit in int64 fail validation for every limit kind.
func (suite *KeeperTestSuite) TestTimeframeDurationOverflowRejected() {
	tooManyHours := math.MaxInt64/3600 + 1
	tooManyDays := math.MaxInt64/86400 + 1

	suite.Require().Error(ratelimittypes.TimeframeLimit{MaxAmount: sdkmath.NewInt(1), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, TimeframeDuration: int64(tooManyHours)}.Validate())
	suite.Require().Error(ratelimittypes.TimeframeLimit{MaxAmount: sdkmath.NewInt(1), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY, TimeframeDuration: int64(tooManyDays)}.Validate())
	suite.Require().Error(ratelimittypes.UniqueSenderLimit{MaxUniqueSenders: 1, TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY, TimeframeDuration: int64(tooManyDays)}.Validate())
	suite.Require().Error(ratelimittypes.AddressLimit{MaxTransfers: 1, MaxAmount: sdkmath.NewInt(1), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, TimeframeDuration: int64(tooManyHours)}.Validate())

	// The largest representable values are still accepted.
	suite.Require().NoError(ratelimittypes.TimeframeLimit{MaxAmount: sdkmath.NewInt(1), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, TimeframeDuration: int64(tooManyHours - 1)}.Validate())
	suite.Require().NoError(ratelimittypes.TimeframeLimit{MaxAmount: sdkmath.NewInt(1), TimeframeType: ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, TimeframeDuration: math.MaxInt64}.Validate())
}

// TestMigrateV35WindowsToBlockTime checks that height-based HOUR/DAY windows
// written before v35 are converted to block-time windows with the same
// remaining lifetime (at the 3 s/block assumption they were created under),
// and that BLOCK windows are left alone.
func (suite *KeeperTestSuite) TestMigrateV35WindowsToBlockTime() {
	hour := ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR
	day := ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY
	block := ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK

	// Legacy windows: HOUR/1 started at height 100 (1200 blocks), DAY/1 at
	// height 40 (28800 blocks), BLOCK/50 at height 120.
	suite.keeper.SetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, hour, 1, ratelimittypes.ChannelFlowWindow{WindowStart: 100, WindowDuration: 1200})
	suite.keeper.SetUniqueSendersWindow(suite.ctx, windowTestChannel, day, 1, ratelimittypes.ChannelFlowWindow{WindowStart: 40, WindowDuration: 28800})
	suite.keeper.SetAddressTransferWindow(suite.ctx, windowTestAddr, windowTestChannel, windowTestDenom, hour, 2, ratelimittypes.ChannelFlowWindow{WindowStart: 100, WindowDuration: 2400})
	suite.keeper.SetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, block, 50, ratelimittypes.ChannelFlowWindow{WindowStart: 120, WindowDuration: 50})

	// Upgrade at height 400, block time T0.
	suite.ctx = suite.ctx.WithBlockHeight(400).WithBlockTime(windowTestT0)
	suite.Require().NoError(suite.keeper.MigrateV35WindowsToBlockTime(suite.ctx))

	w, _ := suite.keeper.GetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, hour, 1)
	suite.Require().Equal(windowTestT0.Unix()-300*3, w.WindowStart, "HOUR window start must be shifted back by elapsed blocks * 3 s")
	suite.Require().Equal(int64(3600), w.WindowDuration)

	w, _ = suite.keeper.GetUniqueSendersWindow(suite.ctx, windowTestChannel, day, 1)
	suite.Require().Equal(windowTestT0.Unix()-360*3, w.WindowStart)
	suite.Require().Equal(int64(86400), w.WindowDuration)

	w, _ = suite.keeper.GetAddressTransferWindow(suite.ctx, windowTestAddr, windowTestChannel, windowTestDenom, hour, 2)
	suite.Require().Equal(windowTestT0.Unix()-300*3, w.WindowStart)
	suite.Require().Equal(int64(7200), w.WindowDuration)

	w, _ = suite.keeper.GetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, block, 50)
	suite.Require().Equal(ratelimittypes.ChannelFlowWindow{WindowStart: 120, WindowDuration: 50}, w, "BLOCK windows must not change")

	// Running the migration twice must not shift windows again.
	suite.Require().NoError(suite.keeper.MigrateV35WindowsToBlockTime(suite.ctx))
	w, _ = suite.keeper.GetChannelFlowWindowWithTimeframe(suite.ctx, windowTestChannel, windowTestDenom, hour, 1)
	suite.Require().Equal(windowTestT0.Unix()-300*3, w.WindowStart, "migration must be idempotent")
}
