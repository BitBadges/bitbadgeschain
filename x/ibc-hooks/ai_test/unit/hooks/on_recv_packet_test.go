package hooks

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/ibc-hooks/ai_test/testutil"
	ibc_hooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
)

type OnRecvPacketTestSuite struct {
	testutil.AITestSuite
}

func TestOnRecvPacketTestSuite(t *testing.T) {
	suite.Run(t, new(OnRecvPacketTestSuite))
}

func (suite *OnRecvPacketTestSuite) TestOnRecvPacket_NoHooks() {
	// Test with no hooks registered (nil hooks)
	packet := testutil.GenerateTestPacket("transfer", "channel-0", []byte("test data"))
	relayer, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	// IBC v10: OnRecvPacket uses channelID (not packetID) in ibc-hooks wrapper
	ack := suite.IBCMiddleware.OnRecvPacket(suite.Ctx, packet.GetDestChannel(), packet, relayer)
	suite.Require().NotNil(ack)
	suite.Require().True(ack.Success(), "should succeed when no hooks")
}

func (suite *OnRecvPacketTestSuite) TestOnRecvPacket_OverrideHook() {
	// Create a mock override hook
	mockHook := &MockOnRecvPacketOverrideHook{
		shouldSucceed: true,
	}

	// Set the hook in ICS4Middleware
	ics4Middleware := ibc_hooks.NewICS4Middleware(suite.MockICS4, mockHook)
	suite.ICS4Middleware = &ics4Middleware
	suite.IBCMiddleware = ibc_hooks.NewIBCMiddleware(suite.MockApp, suite.ICS4Middleware)

	packet := testutil.GenerateTestPacket("transfer", "channel-0", []byte("test data"))
	relayer, err := sdk.AccAddressFromBech32(suite.Alice)
	suite.Require().NoError(err)

	// IBC v10: OnRecvPacket uses channelID (not packetID) in ibc-hooks wrapper
	ack := suite.IBCMiddleware.OnRecvPacket(suite.Ctx, packet.GetDestChannel(), packet, relayer)
	suite.Require().NotNil(ack)
	suite.Require().True(ack.Success(), "override hook should control the result")
	suite.Require().True(mockHook.wasCalled, "override hook should be called")
}

// MockOnRecvPacketOverrideHook is a mock implementation of OnRecvPacketOverrideHooks
type MockOnRecvPacketOverrideHook struct {
	shouldSucceed bool
	wasCalled     bool
}

// IBC v10: OnRecvPacketOverride requires channelID parameter
func (m *MockOnRecvPacketOverrideHook) OnRecvPacketOverride(im ibc_hooks.IBCMiddleware, ctx sdk.Context, channelID string, packet channeltypes.Packet, relayer sdk.AccAddress) ibcexported.Acknowledgement {
	m.wasCalled = true
	if m.shouldSucceed {
		return channeltypes.NewResultAcknowledgement([]byte("success"))
	}
	return channeltypes.NewErrorAcknowledgement(fmt.Errorf("mock error"))
}

