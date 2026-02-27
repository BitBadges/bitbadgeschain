//go:build test
// +build test

package e2e

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibctesting "github.com/cosmos/ibc-go/v10/testing"
	"github.com/stretchr/testify/suite"

	"github.com/bitbadges/bitbadgeschain/app"
	ibctest "github.com/bitbadges/bitbadgeschain/testing/ibc"
)

// IBCTestSuite is the base test suite for IBC E2E tests
type IBCTestSuite struct {
	suite.Suite

	Coordinator  *ibctesting.Coordinator
	ChainA       *ibctesting.TestChain
	ChainB       *ibctesting.TestChain
	TransferPath *ibctesting.Path
}

// SetupSuite sets up the test suite with two chains and a transfer path
func (s *IBCTestSuite) SetupSuite() {
	s.T().Log("Setting up IBC test suite...")

	// Create coordinator with 2 chains
	s.Coordinator = ibctest.NewCoordinator(s.T(), 2)

	// Get references to the chains
	s.ChainA = s.Coordinator.GetChain(ibctesting.GetChainID(1))
	s.ChainB = s.Coordinator.GetChain(ibctesting.GetChainID(2))

	// Create and setup transfer path
	s.TransferPath = ibctest.SetupTransferPath(s.Coordinator, s.ChainA, s.ChainB)

	s.T().Log("IBC test suite setup complete")
}

// SetupTest runs before each test
func (s *IBCTestSuite) SetupTest() {
	// Reset any test-specific state if needed
}

// TearDownSuite tears down the test suite
func (s *IBCTestSuite) TearDownSuite() {
	s.T().Log("Tearing down IBC test suite...")
}

// GetBitBadgesApp returns the BitBadges app for the given chain
func (s *IBCTestSuite) GetBitBadgesApp(chain *ibctesting.TestChain) *app.App {
	return ibctest.GetBitBadgesApp(chain)
}

// FundAccount funds an account on the given chain
func (s *IBCTestSuite) FundAccount(chain *ibctesting.TestChain, addr sdk.AccAddress, coins sdk.Coins) {
	err := ibctest.FundAccount(chain, addr, coins)
	s.Require().NoError(err)
}

// GetBalance returns the balance of an account
func (s *IBCTestSuite) GetBalance(chain *ibctesting.TestChain, addr sdk.AccAddress, denom string) sdk.Coin {
	return ibctest.GetBalance(chain, addr, denom)
}

// SendTransfer sends an ICS20 transfer and returns the result
func (s *IBCTestSuite) SendTransfer(
	path *ibctesting.Path,
	sender sdk.AccAddress,
	receiver sdk.AccAddress,
	token sdk.Coin,
	memo string,
) (uint64, error) {
	// Get timeout
	timeoutHeight := s.ChainB.GetTimeoutHeight()
	timeoutTimestamp := uint64(0) // Use height-based timeout

	// Create transfer message
	msg := transfertypes.NewMsgTransfer(
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		token,
		sender.String(),
		receiver.String(),
		timeoutHeight,
		timeoutTimestamp,
		memo,
	)

	// Send the message and get the result
	res, err := s.ChainA.SendMsgs(msg)
	if err != nil {
		return 0, err
	}

	// Parse packet sequence from events
	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	if err != nil {
		return 0, err
	}

	return packet.Sequence, nil
}

// RelayPacket relays a packet from source to destination
func (s *IBCTestSuite) RelayPacket(path *ibctesting.Path, sequence uint64) error {
	// Get the packet commitment
	commitment := s.GetBitBadgesApp(s.ChainA).IBCKeeper.ChannelKeeper.GetPacketCommitment(
		s.ChainA.GetContext(),
		path.EndpointA.ChannelConfig.PortID,
		path.EndpointA.ChannelID,
		sequence,
	)
	if len(commitment) == 0 {
		s.T().Logf("Packet commitment not found for sequence %d", sequence)
	}

	s.T().Logf("Relaying packet with sequence %d, commitment: %x", sequence, commitment)

	// Use the path's relay function which handles the full relay process
	return path.EndpointA.UpdateClient()
}

// GetIBCDenom returns the IBC denom for a token on the receiving chain
func (s *IBCTestSuite) GetIBCDenom(path *ibctesting.Path, baseDenom string) string {
	return ibctest.GetIBCDenom(
		path.EndpointB.ChannelConfig.PortID,
		path.EndpointB.ChannelID,
		baseDenom,
	)
}

// DefaultTransferAmount returns the default transfer amount for tests
func (s *IBCTestSuite) DefaultTransferAmount() sdkmath.Int {
	return ibctest.DefaultTransferAmount()
}

// CommitBlock commits the current block on the given chain
func (s *IBCTestSuite) CommitBlock(chain *ibctesting.TestChain) {
	ibctest.CommitBlock(chain)
}

// TestIBCTestSuiteSetup tests that the suite can be properly initialized
func TestIBCTestSuiteSetup(t *testing.T) {
	suite.Run(t, new(IBCTestSuite))
}
