package middleware

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/x/ibc-hooks/ai_test/testutil"
	ibc_hooks "github.com/bitbadges/bitbadgeschain/x/ibc-hooks"
)

type SendPacketTestSuite struct {
	testutil.AITestSuite
}

func TestSendPacketTestSuite(t *testing.T) {
	suite.Run(t, new(SendPacketTestSuite))
}

func (suite *SendPacketTestSuite) TestSendPacket_NoHooks() {
	// Test with no hooks registered (nil hooks)
	chanCap := testutil.GenerateTestCapability()
	timeoutHeight := clienttypes.Height{}
	timeoutTimestamp := uint64(0)
	data := []byte("test data")

	seq, err := suite.ICS4Middleware.SendPacket(
		suite.Ctx,
		chanCap,
		"transfer",
		"channel-0",
		timeoutHeight,
		timeoutTimestamp,
		data,
	)
	suite.Require().NoError(err)
	suite.Require().Equal(uint64(1), seq)
	suite.Require().True(suite.MockICS4.SendPacketCalled, "mock ICS4 should be called")
}

func (suite *SendPacketTestSuite) TestSendPacket_OverrideHook() {
	// Create a mock override hook
	mockHook := &MockSendPacketOverrideHook{
		shouldSucceed: true,
	}

	// Set the hook in ICS4Middleware
	ics4Middleware := ibc_hooks.NewICS4Middleware(suite.MockICS4, mockHook)
	suite.ICS4Middleware = &ics4Middleware

	chanCap := testutil.GenerateTestCapability()
	timeoutHeight := clienttypes.Height{}
	timeoutTimestamp := uint64(0)
	data := []byte("test data")

	seq, err := suite.ICS4Middleware.SendPacket(
		suite.Ctx,
		chanCap,
		"transfer",
		"channel-0",
		timeoutHeight,
		timeoutTimestamp,
		data,
	)
	suite.Require().NoError(err)
	suite.Require().True(mockHook.wasCalled, "override hook should be called")
	_ = seq
}

// MockSendPacketOverrideHook is a mock implementation of SendPacketOverrideHooks
type MockSendPacketOverrideHook struct {
	shouldSucceed bool
	wasCalled     bool
}

func (m *MockSendPacketOverrideHook) SendPacketOverride(i ibc_hooks.ICS4Middleware, ctx sdk.Context, chanCap *capabilitytypes.Capability, sourcePort string, sourceChannel string, timeoutHeight clienttypes.Height, timeoutTimestamp uint64, data []byte) (uint64, error) {
	m.wasCalled = true
	if m.shouldSucceed {
		return 1, nil
	}
	return 0, fmt.Errorf("mock error")
}

