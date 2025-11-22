package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"context"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/keeper"
	ratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

type KeeperTestSuite struct {
	suite.Suite

	keeper     keeper.Keeper
	bankKeeper *MockBankKeeper
	ctx        sdk.Context
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	// Setup store keys
	storeKey := storetypes.NewKVStoreKey(ratelimittypes.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey("mem_" + ratelimittypes.StoreKey)

	// Setup database
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(suite.T(), stateStore.LoadLatestVersion())

	// Setup codec
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	// Setup bank keeper (minimal setup for testing)
	// For testing, we'll create a mock bank keeper interface
	// The actual bank keeper setup is complex, so we'll use a simpler approach
	bankStoreKey := storetypes.NewKVStoreKey(banktypes.StoreKey)
	stateStore.MountStoreWithDB(bankStoreKey, storetypes.StoreTypeIAVL, db)

	// Create a minimal bank keeper for testing
	// Note: This is a simplified setup - in production, bank keeper needs more dependencies
	bankKeeper := &MockBankKeeper{
		supplies: make(map[string]sdk.Coin),
	}

	// Setup context
	ctx := sdk.NewContext(stateStore, cmtproto.Header{Height: 1}, false, log.NewNopLogger())

	// Setup rate limit keeper
	// Use a valid bech32 address for testing
	authority := "cosmos1w6t0l7z0yerj49ehnqwqaayxqpe3u7e23edgma"
	suite.keeper = keeper.NewKeeper(
		cdc,
		storeKey,
		bankKeeper,
		authority,
	)

	// Set default params for testing
	defaultParams := ratelimittypes.DefaultParams()
	// Add a test config: channel-0, uatom, max shift 300000 per 1000 blocks
	defaultParams.RateLimits = []ratelimittypes.RateLimitConfig{
		{
			ChannelId: "channel-0",
			Denom:     "uatom",
			SupplyShiftLimits: []ratelimittypes.TimeframeLimit{
				{
					MaxAmount:         sdkmath.NewInt(300000), // 300,000 tokens max shift
					TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
					TimeframeDuration: 1000, // 1000 blocks
				},
			},
		},
	}
	suite.keeper.SetParams(ctx, defaultParams)
	suite.bankKeeper = bankKeeper
	suite.ctx = ctx
}

// getAckError extracts the error message from an acknowledgement
func getAckError(ack ibcexported.Acknowledgement) string {
	channelAck, ok := ack.(channeltypes.Acknowledgement)
	if !ok {
		return ""
	}
	if errResp, ok := channelAck.Response.(*channeltypes.Acknowledgement_Error); ok {
		return errResp.Error
	}
	return ""
}

// MockBankKeeper is a simple mock for testing
type MockBankKeeper struct {
	supplies map[string]sdk.Coin
}

func (m *MockBankKeeper) GetSupply(ctx context.Context, denom string) sdk.Coin {
	if coin, ok := m.supplies[denom]; ok {
		return coin
	}
	return sdk.Coin{Denom: denom, Amount: sdkmath.ZeroInt()}
}

func (m *MockBankKeeper) GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins {
	return sdk.Coins{}
}

func (m *MockBankKeeper) MintCoins(ctx sdk.Context, moduleName string, coins sdk.Coins) error {
	for _, coin := range coins {
		if existing, ok := m.supplies[coin.Denom]; ok {
			m.supplies[coin.Denom] = sdk.NewCoin(coin.Denom, existing.Amount.Add(coin.Amount))
		} else {
			m.supplies[coin.Denom] = coin
		}
	}
	return nil
}

var _ ratelimittypes.BankKeeper = (*MockBankKeeper)(nil)

func (suite *KeeperTestSuite) TestGetSetChannelFlow() {
	channelID := "channel-0"
	denom := "uatom"

	// Initially, flow should not exist
	flow, found := suite.keeper.GetChannelFlow(suite.ctx, channelID, denom)
	suite.Require().False(found)
	suite.Require().True(flow.NetFlow.IsZero())

	// Set flow
	newFlow := ratelimittypes.ChannelFlow{
		NetFlow: sdkmath.NewInt(1000),
	}
	suite.keeper.SetChannelFlow(suite.ctx, channelID, denom, newFlow)

	// Get flow
	flow, found = suite.keeper.GetChannelFlow(suite.ctx, channelID, denom)
	suite.Require().True(found)
	suite.Require().Equal(sdkmath.NewInt(1000), flow.NetFlow)
}

func (suite *KeeperTestSuite) TestGetSetChannelFlowWindow() {
	channelID := "channel-0"
	denom := "uatom"

	// Initially, window should not exist
	window, found := suite.keeper.GetChannelFlowWindow(suite.ctx, channelID, denom)
	suite.Require().False(found)

	// Set window
	newWindow := ratelimittypes.ChannelFlowWindow{
		WindowStart:    100,
		WindowDuration: 1000,
	}
	suite.keeper.SetChannelFlowWindow(suite.ctx, channelID, denom, newWindow)

	// Get window
	window, found = suite.keeper.GetChannelFlowWindow(suite.ctx, channelID, denom)
	suite.Require().True(found)
	suite.Require().Equal(int64(100), window.WindowStart)
	suite.Require().Equal(int64(1000), window.WindowDuration)
}

func (suite *KeeperTestSuite) TestResetChannelFlowWindow() {
	channelID := "channel-0"
	denom := "uatom"
	windowDuration := int64(1000)

	// Set initial window at height 100
	suite.ctx = suite.ctx.WithBlockHeight(100)
	window := ratelimittypes.ChannelFlowWindow{
		WindowStart:    100,
		WindowDuration: windowDuration,
	}
	suite.keeper.SetChannelFlowWindow(suite.ctx, channelID, denom, window)
	suite.keeper.SetChannelFlow(suite.ctx, channelID, denom, ratelimittypes.ChannelFlow{NetFlow: sdkmath.NewInt(500)})

	// Move to height 1100 (window should expire)
	suite.ctx = suite.ctx.WithBlockHeight(1100)
	suite.keeper.ResetChannelFlowWindow(suite.ctx, channelID, denom, windowDuration)

	// Window should be reset
	window, found := suite.keeper.GetChannelFlowWindow(suite.ctx, channelID, denom)
	suite.Require().True(found)
	suite.Require().Equal(int64(1100), window.WindowStart)

	// Flow should be reset to zero
	flow, _ := suite.keeper.GetChannelFlow(suite.ctx, channelID, denom)
	suite.Require().True(flow.NetFlow.IsZero())
}

func (suite *KeeperTestSuite) TestUpdateChannelFlow() {
	// This test is kept for backward compatibility but UpdateChannelFlow is deprecated
	// The function now does nothing as it's been replaced by timeframe-based tracking
	channelID := "channel-0"
	denom := "uatom"

	// UpdateChannelFlow is deprecated and no longer updates flow
	// This test verifies it doesn't crash
	suite.keeper.UpdateChannelFlow(suite.ctx, channelID, denom, sdkmath.NewInt(100))

	// Flow should not exist since UpdateChannelFlow is deprecated
	_, found := suite.keeper.GetChannelFlow(suite.ctx, channelID, denom)
	suite.Require().False(found)
}

func (suite *KeeperTestSuite) TestCheckRateLimit_WithinLimit() {
	channelID := "channel-0"
	denom := "uatom"

	// Config has max shift of 300,000 tokens
	// Try to transfer 200,000 (within limit)
	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(200000), true, "")
	suite.Require().True(ack.Success())
}

func (suite *KeeperTestSuite) TestCheckRateLimit_ExceedsLimit() {
	channelID := "channel-0"
	denom := "uatom"

	// Config has max shift of 300,000 tokens
	// Try to transfer 400,000 (exceeds limit)
	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(400000), true, "")
	suite.Require().False(ack.Success())
	suite.Require().Contains(getAckError(ack), "rate limit exceeded")
}

func (suite *KeeperTestSuite) TestCheckRateLimit_NetFlow() {
	channelID := "channel-0"
	denom := "uatom"

	// Config has max shift of 300,000 tokens per 1000 blocks
	// Manually set existing flow and window for testing
	currentHeight := suite.ctx.BlockHeight()
	window := ratelimittypes.ChannelFlowWindow{
		WindowStart:    currentHeight,
		WindowDuration: 1000, // 1000 blocks
	}
	suite.keeper.SetChannelFlowWindowWithTimeframe(suite.ctx, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000, window)
	flow := ratelimittypes.ChannelFlow{NetFlow: sdkmath.NewInt(200000)}
	suite.keeper.SetChannelFlowWithTimeframe(suite.ctx, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000, flow)

	// Try to add more inflow that would exceed limit
	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(150000), true, "")
	suite.Require().False(ack.Success()) // Total would be 350,000 > 300,000

	// But 100,000 more would be OK
	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(100000), true, "")
	suite.Require().True(ack.Success()) // Total would be 300,000 = limit
}

func (suite *KeeperTestSuite) TestCheckRateLimit_Outflow() {
	channelID := "channel-0"
	denom := "uatom"

	// Config has max shift of 300,000 tokens per 1000 blocks
	// Manually set existing outflow (negative flow) and window for testing
	currentHeight := suite.ctx.BlockHeight()
	window := ratelimittypes.ChannelFlowWindow{
		WindowStart:    currentHeight,
		WindowDuration: 1000, // 1000 blocks
	}
	suite.keeper.SetChannelFlowWindowWithTimeframe(suite.ctx, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000, window)
	flow := ratelimittypes.ChannelFlow{NetFlow: sdkmath.NewInt(-200000)}
	suite.keeper.SetChannelFlowWithTimeframe(suite.ctx, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1000, flow)

	// Try to add more outflow
	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(150000), false, "")
	suite.Require().False(ack.Success()) // Net flow would be -350,000, abs > 300,000
}

func (suite *KeeperTestSuite) TestCheckRateLimit_Disabled() {
	// Set disabled params on the existing keeper (empty configs = no rate limits)
	params := ratelimittypes.DefaultParams()
	suite.keeper.SetParams(suite.ctx, params)

	// Should allow any transfer when no config matches
	ack := suite.keeper.CheckRateLimit(suite.ctx, "channel-0", "uatom", sdkmath.NewInt(1000000), true, "")
	suite.Require().True(ack.Success())
}

func (suite *KeeperTestSuite) TestCheckRateLimit_NoSupply() {
	channelID := "channel-0"
	denom := "newtoken" // Token with no supply

	// Don't set supply in mock - GetSupply will return zero coin
	// With the updated keeper logic, zero supply should allow transfers
	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(1000000), true, "")
	suite.Require().True(ack.Success()) // Should allow transfer if supply is zero or not found
}

// TestCheckRateLimit_MultipleTimeframes tests multiple timeframe limits
func (suite *KeeperTestSuite) TestCheckRateLimit_MultipleTimeframes() {
	channelID := "channel-0"
	denom := "uatom"
	senderAddr := "cosmos1test123"

	// Set up config with multiple timeframe limits
	params := ratelimittypes.DefaultParams()
	params.RateLimits = []ratelimittypes.RateLimitConfig{
		{
			ChannelId: channelID,
			Denom:     denom,
			SupplyShiftLimits: []ratelimittypes.TimeframeLimit{
				{
					MaxAmount:         sdkmath.NewInt(100000), // 100k per block
					TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
					TimeframeDuration: 1,
				},
				{
					MaxAmount:         sdkmath.NewInt(500000), // 500k per hour (600 blocks)
					TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR,
					TimeframeDuration: 1,
				},
				{
					MaxAmount:         sdkmath.NewInt(2000000), // 2M per day (14400 blocks)
					TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY,
					TimeframeDuration: 1,
				},
			},
		},
	}
	suite.keeper.SetParams(suite.ctx, params)

	// Test: Transfer within all limits
	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(50000), true, senderAddr)
	suite.Require().True(ack.Success())

	// Test: Transfer exceeds block limit
	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(150000), true, senderAddr)
	suite.Require().False(ack.Success())
	suite.Require().Contains(getAckError(ack), "rate limit exceeded")
}

// TestCheckRateLimit_UniqueSenders tests unique sender tracking
func (suite *KeeperTestSuite) TestCheckRateLimit_UniqueSenders() {
	channelID := "channel-0"
	denom := "uatom"

	// Set up config with unique sender limit
	params := ratelimittypes.DefaultParams()
	params.RateLimits = []ratelimittypes.RateLimitConfig{
		{
			ChannelId: channelID,
			Denom:     denom,
			UniqueSenderLimits: []ratelimittypes.UniqueSenderLimit{
				{
					MaxUniqueSenders:  3, // Max 3 unique senders per block
					TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
					TimeframeDuration: 1,
				},
			},
		},
	}
	suite.keeper.SetParams(suite.ctx, params)

	// First 3 unique senders should be allowed
	sender1 := "cosmos1sender1"
	sender2 := "cosmos1sender2"
	sender3 := "cosmos1sender3"
	sender4 := "cosmos1sender4"

	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(1000), true, sender1)
	suite.Require().True(ack.Success())
	suite.keeper.AddUniqueSender(suite.ctx, channelID, sender1, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)

	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(1000), true, sender2)
	suite.Require().True(ack.Success())
	suite.keeper.AddUniqueSender(suite.ctx, channelID, sender2, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)

	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(1000), true, sender3)
	suite.Require().True(ack.Success())
	suite.keeper.AddUniqueSender(suite.ctx, channelID, sender3, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)

	// 4th unique sender should be rejected
	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(1000), true, sender4)
	suite.Require().False(ack.Success())
	suite.Require().Contains(getAckError(ack), "rate limit exceeded")

	// But same sender (sender1) should still be allowed
	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(1000), true, sender1)
	suite.Require().True(ack.Success())
}

// TestCheckRateLimit_AddressLimits tests per-address transfer limits
func (suite *KeeperTestSuite) TestCheckRateLimit_AddressLimits() {
	channelID := "channel-0"
	denom := "uatom"
	senderAddr := "cosmos1test123"

	// Set up config with per-address limits
	params := ratelimittypes.DefaultParams()
	params.RateLimits = []ratelimittypes.RateLimitConfig{
		{
			ChannelId: channelID,
			Denom:     denom,
			AddressLimits: []ratelimittypes.AddressLimit{
				{
					MaxTransfers:      5,                     // Max 5 transfers per block
					MaxAmount:         sdkmath.NewInt(10000), // Max 10k per block
					TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
					TimeframeDuration: 1,
				},
			},
		},
	}
	suite.keeper.SetParams(suite.ctx, params)

	// First 5 transfers should be allowed
	for i := 0; i < 5; i++ {
		ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(1000), true, senderAddr)
		suite.Require().True(ack.Success())
		// Update tracking
		suite.keeper.ResetAddressTransferWindow(suite.ctx, senderAddr, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)
		data, _ := suite.keeper.GetAddressTransferData(suite.ctx, senderAddr, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)
		data.TransferCount++
		data.TotalAmount = data.TotalAmount.Add(sdkmath.NewInt(1000))
		suite.keeper.SetAddressTransferData(suite.ctx, senderAddr, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1, data)
	}

	// 6th transfer should be rejected (exceeds transfer count limit)
	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(1000), true, senderAddr)
	suite.Require().False(ack.Success())
	suite.Require().Contains(getAckError(ack), "rate limit exceeded")
}

// TestCheckRateLimit_AddressAmountLimit tests per-address amount limits
func (suite *KeeperTestSuite) TestCheckRateLimit_AddressAmountLimit() {
	channelID := "channel-0"
	denom := "uatom"
	senderAddr := "cosmos1test123"

	// Set up config with per-address amount limit
	params := ratelimittypes.DefaultParams()
	params.RateLimits = []ratelimittypes.RateLimitConfig{
		{
			ChannelId: channelID,
			Denom:     denom,
			AddressLimits: []ratelimittypes.AddressLimit{
				{
					MaxTransfers:      0,                     // No transfer count limit
					MaxAmount:         sdkmath.NewInt(10000), // Max 10k per block
					TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
					TimeframeDuration: 1,
				},
			},
		},
	}
	suite.keeper.SetParams(suite.ctx, params)

	// Transfer within limit
	ack := suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(5000), true, senderAddr)
	suite.Require().True(ack.Success())

	// Update tracking
	suite.keeper.ResetAddressTransferWindow(suite.ctx, senderAddr, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)
	data, _ := suite.keeper.GetAddressTransferData(suite.ctx, senderAddr, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1)
	data.TransferCount++
	data.TotalAmount = data.TotalAmount.Add(sdkmath.NewInt(5000))
	suite.keeper.SetAddressTransferData(suite.ctx, senderAddr, channelID, denom, ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 1, data)

	// Another transfer that would exceed total amount limit
	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(6000), true, senderAddr)
	suite.Require().False(ack.Success())
	suite.Require().Contains(getAckError(ack), "rate limit exceeded")

	// But a smaller transfer within limit should work
	ack = suite.keeper.CheckRateLimit(suite.ctx, channelID, denom, sdkmath.NewInt(4000), true, senderAddr)
	suite.Require().True(ack.Success())
}

// TestTimeframeDurationInBlocks tests timeframe conversion
func (suite *KeeperTestSuite) TestTimeframeDurationInBlocks() {
	// Test with default block time (3 seconds for BitBadges chain)
	blockTimeSeconds := int64(3)

	// Test block timeframe (should return as-is)
	blocks := ratelimittypes.TimeframeDurationInBlocks(ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK, 100, blockTimeSeconds)
	suite.Require().Equal(int64(100), blocks)

	// Test hour timeframe (1 hour = 3600 seconds = 1200 blocks at 3s/block)
	blocks = ratelimittypes.TimeframeDurationInBlocks(ratelimittypes.TimeframeType_TIMEFRAME_TYPE_HOUR, 1, blockTimeSeconds)
	suite.Require().Equal(int64(1200), blocks)

	// Test day timeframe (1 day = 86400 seconds = 28800 blocks at 3s/block)
	blocks = ratelimittypes.TimeframeDurationInBlocks(ratelimittypes.TimeframeType_TIMEFRAME_TYPE_DAY, 1, blockTimeSeconds)
	suite.Require().Equal(int64(28800), blocks)
}

// TestValidateParams_EmptyDenom tests that empty denoms are rejected
func (suite *KeeperTestSuite) TestValidateParams_EmptyDenom() {
	params := ratelimittypes.DefaultParams()
	params.RateLimits = []ratelimittypes.RateLimitConfig{
		{
			ChannelId: "channel-0",
			Denom:     "", // Empty denom should be rejected
			SupplyShiftLimits: []ratelimittypes.TimeframeLimit{
				{
					MaxAmount:         sdkmath.NewInt(300000),
					TimeframeType:     ratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
					TimeframeDuration: 1000,
				},
			},
		},
	}

	// Validation should fail
	err := params.Validate()
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "denom must be specified")

	// SetParams should also fail
	err = suite.keeper.SetParams(suite.ctx, params)
	suite.Require().Error(err)
	suite.Require().Contains(err.Error(), "denom must be specified")
}
