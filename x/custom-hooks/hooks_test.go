package customhooks_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/third_party/apptesting"
	customhooks "github.com/bitbadges/bitbadgeschain/x/custom-hooks"
	customhookskeeper "github.com/bitbadges/bitbadgeschain/x/custom-hooks/keeper"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
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
	customHooksKeeper := customhookskeeper.NewKeeper(
		s.App.Logger(),
		s.App.PoolManagerKeeper,
		s.App.BankKeeper,
		s.App.HooksICS4Wrapper,
		s.App.IBCKeeper.ChannelKeeper,
		s.App.ScopedIBCTransferKeeper,
	)

	// Create custom hooks
	bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	s.customHooks = customhooks.NewCustomHooks(customHooksKeeper, bech32Prefix)
}

// TestOnRecvPacketAfterHook_ValidMemo tests hook execution with valid swap_and_action memo
func (s *HooksTestSuite) TestOnRecvPacketAfterHook_ValidMemo() {
	// Create a valid IBC packet with swap_and_action memo
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
	// Create a success acknowledgement
	ack := channeltypes.NewResultAcknowledgement([]byte{1})

	// Call the hook
	s.customHooks.OnRecvPacketAfterHook(s.Ctx, packet, relayer, ack)

	// Note: In a real test, you'd verify the swap was executed
	// This test mainly verifies the hook doesn't panic with valid memo
}

// TestOnRecvPacketAfterHook_InvalidMemo tests hook execution with invalid memo
func (s *HooksTestSuite) TestOnRecvPacketAfterHook_InvalidMemo() {
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
	// Create a success acknowledgement
	ack := channeltypes.NewResultAcknowledgement([]byte{1})

	// Call the hook - should not panic even with invalid memo
	s.customHooks.OnRecvPacketAfterHook(s.Ctx, packet, relayer, ack)
}

// TestOnRecvPacketAfterHook_EmptyMemo tests hook execution with empty memo
func (s *HooksTestSuite) TestOnRecvPacketAfterHook_EmptyMemo() {
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
	// Create a success acknowledgement
	ack := channeltypes.NewResultAcknowledgement([]byte{1})

	// Call the hook - should handle empty memo gracefully
	s.customHooks.OnRecvPacketAfterHook(s.Ctx, packet, relayer, ack)
}

// TestOnRecvPacketAfterHook_NonICS20Packet tests hook with non-transfer packet
func (s *HooksTestSuite) TestOnRecvPacketAfterHook_NonICS20Packet() {
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
	// Create a success acknowledgement
	ack := channeltypes.NewResultAcknowledgement([]byte{1})

	// Call the hook - should not process non-ICS20 packets
	s.customHooks.OnRecvPacketAfterHook(s.Ctx, packet, relayer, ack)
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
