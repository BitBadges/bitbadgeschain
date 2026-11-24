package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	ratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

// TestUniqueSenderLimitsIsolatedByDenom verifies that unique sender tracking on
// the same channel and timeframe is scoped per denom so traffic on one denom
// does not exhaust the unique sender quota of another denom.
func (suite *KeeperTestSuite) TestUniqueSenderLimitsIsolatedByDenom() {
	channelID := "channel-0"
	denomA := "uatom"
	denomB := "uusdc"

	params := ratelimittypes.DefaultParams()
	params.RateLimits = []ratelimittypes.RateLimitConfig{
		{
			ChannelId: channelID,
			Denom:     denomA,
			UniqueSenderLimits: []ratelimittypes.UniqueSenderLimit{{
				MaxUniqueSenders:  2,
				TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
				TimeframeDuration: 1,
			}},
		},
		{
			ChannelId: channelID,
			Denom:     denomB,
			UniqueSenderLimits: []ratelimittypes.UniqueSenderLimit{{
				MaxUniqueSenders:  2,
				TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
				TimeframeDuration: 1,
			}},
		},
	}
	suite.Require().NoError(suite.keeper.SetParams(suite.ctx, params))

	// Saturate unique sender quota on denomA
	for _, sender := range []string{"cosmos1senderA1", "cosmos1senderA2"} {
		ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denomA, sdkmath.NewInt(100), true, sender)
		suite.Require().True(ack.Success())
		suite.keeper.AddUniqueSender(suite.ctx, channelID, denomA, sender, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)
	}

	// 3rd unique sender on denomA is rejected
	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denomA, sdkmath.NewInt(100), true, "cosmos1senderA3")
	suite.Require().False(ack.Success())

	// DenomB on the same channel and block must remain unaffected
	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denomB, sdkmath.NewInt(100), true, "cosmos1senderB1")
	suite.Require().True(ack.Success())
	suite.keeper.AddUniqueSender(suite.ctx, channelID, denomB, "cosmos1senderB1", ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)

	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denomB, sdkmath.NewInt(100), true, "cosmos1senderB2")
	suite.Require().True(ack.Success())
	suite.keeper.AddUniqueSender(suite.ctx, channelID, denomB, "cosmos1senderB2", ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)

	// 3rd unique sender on denomB is now rejected independently
	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denomB, sdkmath.NewInt(100), true, "cosmos1senderB3")
	suite.Require().False(ack.Success())
}
