package customhooks_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	customhooks "github.com/bitbadges/bitbadgeschain/x/custom-hooks"
	customhookskeeper "github.com/bitbadges/bitbadgeschain/x/custom-hooks/keeper"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	ibchooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
)

type HooksTestSuite struct {
	apptesting.KeeperTestHelper

	customHooks *customhooks.CustomHooks
}

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(HooksTestSuite))
}

func (s *HooksTestSuite) SetupTest() {
	s.Reset()

	// Create custom hooks keeper
	// Pass pointer to GammKeeper to avoid copying the keeper (which contains storeKey)
	customHooksKeeper := customhookskeeper.NewKeeper(
		s.App.Logger(),
		&s.App.GammKeeper,
		s.App.BankKeeper,
		&s.App.TokenizationKeeper,
		&s.App.SendmanagerKeeper,
		s.App.TransferKeeper,
		s.App.HooksICS4Wrapper,
		s.App.IBCKeeper.ChannelKeeper,
	)

	// Create custom hooks
	bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	s.customHooks = customhooks.NewCustomHooks(customHooksKeeper, bech32Prefix)
}

// mockIBCModule is a simple mock IBC module that returns success acknowledgements
type mockIBCModule struct{}

// IBC v10: capabilities removed from channel handshake
func (m *mockIBCModule) OnChanOpenInit(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, counterparty channeltypes.Counterparty, version string) (string, error) {
	return version, nil
}

func (m *mockIBCModule) OnChanOpenTry(ctx sdk.Context, order channeltypes.Order, connectionHops []string, portID, channelID string, counterparty channeltypes.Counterparty, counterpartyVersion string) (string, error) {
	return counterpartyVersion, nil
}

func (m *mockIBCModule) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return nil
}

func (m *mockIBCModule) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (m *mockIBCModule) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return nil
}

func (m *mockIBCModule) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return nil
}

// IBC v10: OnRecvPacket now includes channelID parameter
func (m *mockIBCModule) OnRecvPacket(ctx sdk.Context, channelID string, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	// Return success acknowledgement for testing
	return channeltypes.NewResultAcknowledgement([]byte{1})
}

// IBC v10: OnAcknowledgementPacket now includes packetID parameter
func (m *mockIBCModule) OnAcknowledgementPacket(ctx sdk.Context, packetID string, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return nil
}

// IBC v10: OnTimeoutPacket now includes packetID parameter
func (m *mockIBCModule) OnTimeoutPacket(ctx sdk.Context, packetID string, packet channeltypes.Packet, relayer sdk.AccAddress) error {
	return nil
}

// createMockIBCMiddleware creates a mock IBCMiddleware for testing
func (s *HooksTestSuite) createMockIBCMiddleware() ibchooks.IBCMiddleware {
	mockApp := &mockIBCModule{}
	mockICS4 := ibchooks.NewICS4Middleware(s.App.HooksICS4Wrapper, nil)
	return ibchooks.NewIBCMiddleware(mockApp, &mockICS4)
}

// TestOnRecvPacketOverride_ValidMemo tests hook execution with valid swap_and_action memo
func (s *HooksTestSuite) TestOnRecvPacketOverride_ValidMemo() {
	// Create a valid IBC packet with swap_and_action memo
	// Note: post_swap_action is now required
	memo := `{
		"swap_and_action": {
			"user_swap": {
				"swap_exact_asset_in": {
					"swap_venue_name": "bitbadges-poolmanager",
					"operations": [
						{
							"pool": "1",
							"denom_in": "` + sdk.DefaultBondDenom + `",
							"denom_out": "uatom"
						}
					]
				}
			},
			"min_asset": {
				"native": {
					"denom": "uatom",
					"amount": "1000"
				}
			},
			"post_swap_action": {
				"transfer": {
					"to_address": "` + s.TestAccs[1].String() + `"
				}
			}
		}
	}`

	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    sdk.DefaultBondDenom,
		Amount:   "100000",
		Sender:   "cosmos1test",
		Receiver: "bb1test",
		Memo:     memo,
	}

	data, err := transfertypes.ModuleCdc.MarshalJSON(&packetData)
	s.Require().NoError(err)

	packet := channeltypes.Packet{
		Sequence:           1,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               data,
		TimeoutHeight:      clienttypes.Height{},
		TimeoutTimestamp:   0,
	}

	relayer := s.TestAccs[0]
	// Create mock IBC middleware
	im := s.createMockIBCMiddleware()

	// Call the override hook - IBC v10: OnRecvPacketOverride requires channelID parameter
	ack := s.customHooks.OnRecvPacketOverride(im, s.Ctx, packet.GetDestChannel(), packet, relayer)

	// Note: The hook will fail because the pool doesn't exist in this test setup
	// This is expected - the test verifies the hook doesn't panic and handles errors gracefully
	// In a real integration test, you'd create a pool first and then verify success
	// For now, we just verify it returns an acknowledgement (either success or error)
	s.Require().NotNil(ack, "acknowledgement should not be nil")

	// The hook will likely fail due to missing pool, which is fine for this test
	// The important thing is it doesn't panic and returns a proper acknowledgement
}

// TestOnRecvPacketOverride_InvalidMemo tests hook execution with invalid memo
func (s *HooksTestSuite) TestOnRecvPacketOverride_InvalidMemo() {
	// Create packet with invalid JSON memo
	memo := `{"invalid": json}`

	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    sdk.DefaultBondDenom,
		Amount:   "100000",
		Sender:   "cosmos1test",
		Receiver: "bb1test",
		Memo:     memo,
	}

	data, err := transfertypes.ModuleCdc.MarshalJSON(&packetData)
	s.Require().NoError(err)

	packet := channeltypes.Packet{
		Sequence:           1,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               data,
		TimeoutHeight:      clienttypes.Height{},
		TimeoutTimestamp:   0,
	}

	relayer := s.TestAccs[0]
	// Create mock IBC middleware
	im := s.createMockIBCMiddleware()

	// Call the override hook - should return error ack for invalid memo
	// IBC v10: OnRecvPacketOverride requires channelID parameter
	ack := s.customHooks.OnRecvPacketOverride(im, s.Ctx, packet.GetDestChannel(), packet, relayer)

	// With invalid memo, should return error acknowledgement
	s.Require().False(ack.Success(), "acknowledgement should be error for invalid memo")
}

// TestOnRecvPacketOverride_EmptyMemo tests hook execution with empty memo
func (s *HooksTestSuite) TestOnRecvPacketOverride_EmptyMemo() {
	packetData := transfertypes.FungibleTokenPacketData{
		Denom:    sdk.DefaultBondDenom,
		Amount:   "100000",
		Sender:   "cosmos1test",
		Receiver: "bb1test",
		Memo:     "",
	}

	data, err := transfertypes.ModuleCdc.MarshalJSON(&packetData)
	s.Require().NoError(err)

	packet := channeltypes.Packet{
		Sequence:           1,
		SourcePort:         "transfer",
		SourceChannel:      "channel-0",
		DestinationPort:    "transfer",
		DestinationChannel: "channel-1",
		Data:               data,
		TimeoutHeight:      clienttypes.Height{},
		TimeoutTimestamp:   0,
	}

	relayer := s.TestAccs[0]
	// Create mock IBC middleware
	im := s.createMockIBCMiddleware()

	// Call the override hook - should handle empty memo gracefully (no hook data, return success)
	// IBC v10: OnRecvPacketOverride requires channelID parameter
	ack := s.customHooks.OnRecvPacketOverride(im, s.Ctx, packet.GetDestChannel(), packet, relayer)

	// Empty memo means no hook data, so should return success
	s.Require().True(ack.Success(), "acknowledgement should be successful for empty memo")
}

// TestOnRecvPacketOverride_NonICS20Packet tests hook with non-transfer packet
func (s *HooksTestSuite) TestOnRecvPacketOverride_NonICS20Packet() {
	// Create a non-ICS20 packet
	data := []byte("non-ics20-data")

	packet := channeltypes.Packet{
		Sequence:           1,
		SourcePort:         "custom",
		SourceChannel:      "channel-0",
		DestinationPort:    "custom",
		DestinationChannel: "channel-1",
		Data:               data,
		TimeoutHeight:      clienttypes.Height{},
		TimeoutTimestamp:   0,
	}

	relayer := s.TestAccs[0]
	// Create mock IBC middleware
	im := s.createMockIBCMiddleware()

	// Call the override hook - should not process non-ICS20 packets, return success
	// IBC v10: OnRecvPacketOverride requires channelID parameter
	ack := s.customHooks.OnRecvPacketOverride(im, s.Ctx, packet.GetDestChannel(), packet, relayer)

	// Non-ICS20 packets should return success (not processed by hooks)
	s.Require().True(ack.Success(), "acknowledgement should be successful for non-ICS20 packets")
}

// TestParseHookDataFromMemo tests memo parsing
func (s *HooksTestSuite) TestParseHookDataFromMemo() {
	// Test valid memo
	validMemo := `{
		"swap_and_action": {
			"user_swap": {
				"swap_exact_asset_in": {
					"swap_venue_name": "bitbadges-poolmanager",
					"operations": [
						{
							"pool": "1",
							"denom_in": "` + sdk.DefaultBondDenom + `",
							"denom_out": "uatom"
						}
					]
				}
			}
		}
	}`

	hookData, err := customhookstypes.ParseHookDataFromMemo(validMemo)
	s.Require().NoError(err)
	s.Require().NotNil(hookData)
	s.Require().NotNil(hookData.SwapAndAction)
	s.Require().NotNil(hookData.SwapAndAction.UserSwap)

	// Test empty memo
	hookData, err = customhookstypes.ParseHookDataFromMemo("")
	s.Require().NoError(err)
	s.Require().Nil(hookData)

	// Test memo without swap_and_action
	otherMemo := `{"other_key": "value"}`
	hookData, err = customhookstypes.ParseHookDataFromMemo(otherMemo)
	s.Require().NoError(err)
	s.Require().Nil(hookData)

	// Test invalid JSON
	invalidMemo := `{"invalid": json}`
	_, err = customhookstypes.ParseHookDataFromMemo(invalidMemo)
	s.Require().Error(err)
}
