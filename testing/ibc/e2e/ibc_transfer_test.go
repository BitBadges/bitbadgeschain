//go:build test
// +build test

package e2e

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v10/testing"
	"github.com/stretchr/testify/suite"

	ibctest "github.com/bitbadges/bitbadgeschain/testing/ibc"
)

// TransferTestSuite tests ICS20 token transfers
type TransferTestSuite struct {
	IBCTestSuite
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}

// TestBasicTransfer tests a simple A->B transfer
func (s *TransferTestSuite) TestBasicTransfer() {
	s.T().Log("Testing basic ICS20 transfer A->B")

	// Get sender and receiver addresses
	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	// Get initial balance
	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Fund sender with ubadge tokens
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	initialSenderBalance := s.GetBalance(s.ChainA, sender, denom)

	s.T().Logf("Initial sender balance: %s", initialSenderBalance)
	s.Require().True(initialSenderBalance.Amount.GTE(amount), "sender should have sufficient balance")

	// Create transfer message
	token := sdk.NewCoin(denom, amount)
	timeoutHeight := clienttypes.NewHeight(1, 110)

	msg := transfertypes.NewMsgTransfer(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
		token,
		sender.String(),
		receiver.String(),
		timeoutHeight,
		0, // no timestamp timeout
		"", // no memo
	)

	// Send transfer
	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	// Parse packet from events
	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)
	s.T().Logf("Packet sent with sequence: %d", packet.Sequence)

	// Relay the packet
	err = s.TransferPath.RelayPacket(packet)
	s.Require().NoError(err)

	// Verify sender balance decreased
	finalSenderBalance := s.GetBalance(s.ChainA, sender, denom)
	s.T().Logf("Final sender balance: %s", finalSenderBalance)
	expectedSenderBalance := initialSenderBalance.Sub(token)
	s.Require().Equal(expectedSenderBalance.Amount, finalSenderBalance.Amount, "sender balance should decrease by transfer amount")

	// Verify receiver got the IBC tokens
	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)
	receiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	s.T().Logf("Receiver IBC balance: %s", receiverBalance)
	s.Require().Equal(amount, receiverBalance.Amount, "receiver should have received the tokens")
}

// TestTransferAndReceive tests the full send/relay/receive flow
func (s *TransferTestSuite) TestTransferAndReceive() {
	s.T().Log("Testing full transfer and receive flow")

	// Get sender and receiver addresses
	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	// Fund sender with extra tokens for this test
	denom := "ubadge"
	amount := sdkmath.NewInt(5000000)
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	// Get initial balances (after funding, before transfer)
	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)
	initialReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	initialSenderBalance := s.GetBalance(s.ChainA, sender, denom)
	s.T().Logf("Initial sender balance: %s", initialSenderBalance)
	s.T().Logf("Initial receiver IBC balance: %s", initialReceiverBalance)

	// Send transfer
	token := sdk.NewCoin(denom, amount)
	timeoutHeight := clienttypes.NewHeight(1, 110)

	msg := transfertypes.NewMsgTransfer(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
		token,
		sender.String(),
		receiver.String(),
		timeoutHeight,
		0,
		"",
	)

	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	// Parse and relay packet
	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Relay the packet (receive on B, ack on A)
	err = s.TransferPath.RelayPacket(packet)
	s.Require().NoError(err)

	// Verify the transfer was successful (check delta, not absolute)
	receiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	receiverDelta := receiverBalance.Amount.Sub(initialReceiverBalance.Amount)
	s.T().Logf("Receiver balance after first transfer: %s (delta: %s)", receiverBalance, receiverDelta)
	s.Require().Equal(amount, receiverDelta, "receiver should have received the tokens")

	// Now send the tokens back (B->A)
	s.T().Log("Sending tokens back from B to A")

	returnMsg := transfertypes.NewMsgTransfer(
		s.TransferPath.EndpointB.ChannelConfig.PortID,
		s.TransferPath.EndpointB.ChannelID,
		sdk.NewCoin(ibcDenom, amount),
		receiver.String(),
		sender.String(),
		timeoutHeight,
		0,
		"",
	)

	res, err = s.ChainB.SendMsgs(returnMsg)
	s.Require().NoError(err)

	// Parse and relay return packet
	returnPacket, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	err = s.TransferPath.RelayPacket(returnPacket)
	s.Require().NoError(err)

	// Verify tokens returned to original denom on chain A
	finalSenderBalance := s.GetBalance(s.ChainA, sender, denom)
	s.T().Logf("Final sender balance after return: %s", finalSenderBalance)

	// The sender should have received back the original tokens
	s.Require().Equal(initialSenderBalance.Amount, finalSenderBalance.Amount, "sender should have original balance back")

	// Receiver should have same IBC balance as before this test (delta = 0)
	finalReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	finalReceiverDelta := finalReceiverBalance.Amount.Sub(initialReceiverBalance.Amount)
	s.T().Logf("Final receiver balance: %s (delta from initial: %s)", finalReceiverBalance, finalReceiverDelta)
	s.Require().True(finalReceiverDelta.IsZero(), "receiver should have same IBC balance as before (tokens returned)")
}

// TestTransferTimeout tests timeout handling
func (s *TransferTestSuite) TestTransferTimeout() {
	s.T().Log("Testing transfer timeout handling")

	// Get sender address
	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Fund sender with ubadge tokens
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	// Get initial balance
	initialSenderBalance := s.GetBalance(s.ChainA, sender, denom)
	s.T().Logf("Initial sender balance: %s", initialSenderBalance)

	// Create transfer with a timeout that will expire soon
	// We need to use a timeout height that is valid when sending but will be passed on chain B
	token := sdk.NewCoin(denom, amount)

	// Get current chain B height and set timeout just slightly ahead
	currentHeight := s.ChainB.GetContext().BlockHeight()
	// Use a timeout height that's just 1 block ahead of current
	timeoutHeight := clienttypes.NewHeight(1, uint64(currentHeight+1))

	msg := transfertypes.NewMsgTransfer(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
		token,
		sender.String(),
		receiver.String(),
		timeoutHeight,
		0,
		"",
	)

	// Send transfer
	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	// Parse packet
	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Advance chain B's height well beyond the timeout
	for i := 0; i < 5; i++ {
		s.ChainB.NextBlock()
	}

	// Update client on chain A to reflect chain B's new height
	err = s.TransferPath.EndpointA.UpdateClient()
	s.Require().NoError(err)

	// Trigger timeout on chain A
	err = s.TransferPath.EndpointA.TimeoutPacket(packet)
	s.Require().NoError(err)

	// Verify sender got tokens back after timeout
	finalSenderBalance := s.GetBalance(s.ChainA, sender, denom)
	s.T().Logf("Final sender balance after timeout: %s", finalSenderBalance)
	s.Require().Equal(initialSenderBalance.Amount, finalSenderBalance.Amount, "sender should have original balance after timeout refund")
}

// TestTransferAcknowledgement tests acknowledgement verification
func (s *TransferTestSuite) TestTransferAcknowledgement() {
	s.T().Log("Testing transfer acknowledgement")

	// Get sender and receiver addresses
	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Fund sender with ubadge tokens
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	// Get initial receiver balance (to check delta later)
	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)
	initialReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)

	// Send transfer
	token := sdk.NewCoin(denom, amount)
	timeoutHeight := clienttypes.NewHeight(1, 110)

	msg := transfertypes.NewMsgTransfer(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
		token,
		sender.String(),
		receiver.String(),
		timeoutHeight,
		0,
		"",
	)

	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	// Parse packet
	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Receive packet on chain B
	err = s.TransferPath.EndpointB.UpdateClient()
	s.Require().NoError(err)

	res, err = s.TransferPath.EndpointB.RecvPacketWithResult(packet)
	s.Require().NoError(err)

	// Parse acknowledgement from receive result
	ack, err := ibctesting.ParseAckFromEvents(res.Events)
	s.Require().NoError(err)

	s.T().Logf("Acknowledgement received: %s", string(ack))

	// Verify it's a success acknowledgement
	s.Require().True(len(ack) > 0, "acknowledgement should not be empty")

	// The ack should be a success (contains "result" for success, "error" for failure)
	// Success ack in ICS20 is the hash of the packet data
	s.Require().NotContains(string(ack), "error", "acknowledgement should be successful")

	// Acknowledge the packet on chain A
	err = s.TransferPath.EndpointA.UpdateClient()
	s.Require().NoError(err)

	err = s.TransferPath.EndpointA.AcknowledgePacket(packet, ack)
	s.Require().NoError(err)

	// Verify receiver has the tokens (check delta, not absolute)
	receiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	receiverDelta := receiverBalance.Amount.Sub(initialReceiverBalance.Amount)
	s.Require().Equal(amount, receiverDelta, "receiver should have the transferred tokens")
}

// TestMultipleTransfers tests multiple sequential transfers
func (s *TransferTestSuite) TestMultipleTransfers() {
	s.T().Log("Testing multiple sequential transfers")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	singleAmount := sdkmath.NewInt(100000)
	numTransfers := 5

	// Fund sender with enough tokens
	totalAmount := singleAmount.MulRaw(int64(numTransfers * 2))
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, totalAmount)))
	s.Require().NoError(err)

	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)
	initialReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)

	// Execute multiple transfers
	for i := 0; i < numTransfers; i++ {
		token := sdk.NewCoin(denom, singleAmount)
		timeoutHeight := clienttypes.NewHeight(1, 110)

		msg := transfertypes.NewMsgTransfer(
			s.TransferPath.EndpointA.ChannelConfig.PortID,
			s.TransferPath.EndpointA.ChannelID,
			token,
			sender.String(),
			receiver.String(),
			timeoutHeight,
			0,
			"",
		)

		res, err := s.ChainA.SendMsgs(msg)
		s.Require().NoError(err)

		packet, err := ibctesting.ParsePacketFromEvents(res.Events)
		s.Require().NoError(err)

		err = s.TransferPath.RelayPacket(packet)
		s.Require().NoError(err)

		s.T().Logf("Transfer %d completed", i+1)
	}

	// Verify total received amount
	finalReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	expectedIncrease := singleAmount.MulRaw(int64(numTransfers))
	actualIncrease := finalReceiverBalance.Amount.Sub(initialReceiverBalance.Amount)
	s.Require().Equal(expectedIncrease, actualIncrease, "receiver should have received all transfers")
}
