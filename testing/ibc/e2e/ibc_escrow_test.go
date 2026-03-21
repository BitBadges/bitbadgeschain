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

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	gammtypes "github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
	ibctest "github.com/bitbadges/bitbadgeschain/testing/ibc"
)

// EscrowTestSuite tests that GAMM and custom-hooks IBC transfers correctly
// escrow tokens before sending packets, using the full two-chain IBC stack.
type EscrowTestSuite struct {
	IBCTestSuite
}

func TestEscrowTestSuite(t *testing.T) {
	suite.Run(t, new(EscrowTestSuite))
}

// createPoolOnChainA creates a balancer pool with the given denoms on chain A
// and returns the pool ID. Funds the creator with the pool assets first.
func (s *EscrowTestSuite) createPoolOnChainA(denomA, denomB string, amountA, amountB sdkmath.Int) uint64 {
	app := s.GetBitBadgesApp(s.ChainA)
	ctx := s.ChainA.GetContext()
	creator := s.ChainA.SenderAccount.GetAddress()

	// Fund the creator with pool assets
	poolCoins := sdk.NewCoins(
		sdk.NewCoin(denomA, amountA),
		sdk.NewCoin(denomB, amountB),
	)
	err := ibctest.FundAccount(s.ChainA, creator, poolCoins)
	s.Require().NoError(err)

	// Create the pool
	poolAssets := []balancer.PoolAsset{
		{Token: sdk.NewCoin(denomA, amountA), Weight: osmomath.NewInt(1)},
		{Token: sdk.NewCoin(denomB, amountB), Weight: osmomath.NewInt(1)},
	}

	poolParams := balancer.PoolParams{
		SwapFee: osmomath.ZeroDec(),
		ExitFee: osmomath.ZeroDec(),
	}

	msg := balancer.NewMsgCreateBalancerPool(creator, poolParams, poolAssets)
	poolID, err := app.PoolManagerKeeper.CreatePool(ctx, msg)
	s.Require().NoError(err)

	// Commit the block so pool state is finalized
	s.Coordinator.CommitBlock(s.ChainA)

	return poolID
}

// TestGAMMSwapWithIBCTransfer_EscrowE2E performs a real swap+IBC transfer
// across two chains and verifies that tokens are properly escrowed on the
// source chain and received on the destination chain.
func (s *EscrowTestSuite) TestGAMMSwapWithIBCTransfer_EscrowE2E() {
	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denomIn := "ubadge"
	denomOut := "ustake"
	poolLiquidity := sdkmath.NewInt(10_000_000)
	swapAmount := sdkmath.NewInt(100_000)

	// Create pool with ubadge/ustake on chain A
	poolID := s.createPoolOnChainA(denomIn, denomOut, poolLiquidity, poolLiquidity)
	s.T().Logf("Created pool %d with %s/%s", poolID, denomIn, denomOut)

	// Fund sender with tokens to swap
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denomIn, swapAmount.MulRaw(2))))
	s.Require().NoError(err)

	// Record balances before swap+IBC
	senderDenomOutBefore := s.GetBalance(s.ChainA, sender, denomOut)
	senderDenomInBefore := s.GetBalance(s.ChainA, sender, denomIn)
	s.T().Logf("Sender %s before: %s, %s before: %s", denomIn, senderDenomInBefore, denomOut, senderDenomOutBefore)

	// Build the SwapExactAmountInWithIBCTransfer message
	swapIBCMsg := &gammtypes.MsgSwapExactAmountInWithIBCTransfer{
		Sender: sender.String(),
		Routes: []poolmanagertypes.SwapAmountInRoute{
			{PoolId: poolID, TokenOutDenom: denomOut},
		},
		TokenIn:           sdk.NewCoin(denomIn, swapAmount),
		TokenOutMinAmount: osmomath.NewInt(1),
		IbcTransferInfo: gammtypes.IBCTransferInfo{
			SourceChannel: s.TransferPath.EndpointA.ChannelID,
			Receiver:      receiver.String(),
		},
	}

	// Send the message on chain A
	res, err := s.ChainA.SendMsgs(swapIBCMsg)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	// Parse packet from events — this proves a real IBC packet was created
	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)
	s.T().Logf("IBC packet sent with sequence: %d", packet.Sequence)

	// CRITICAL CHECK: After sending the packet but BEFORE relay,
	// verify the sender does NOT hold the swapped tokens.
	// If escrow is working, the swap output (denomOut) was moved to the
	// IBC transfer module's escrow address.
	senderDenomOutAfterSend := s.GetBalance(s.ChainA, sender, denomOut)
	denomOutDelta := senderDenomOutAfterSend.Amount.Sub(senderDenomOutBefore.Amount)
	s.T().Logf("Sender %s delta after send (before relay): %s", denomOut, denomOutDelta)
	s.Require().True(denomOutDelta.IsZero(),
		"sender should NOT hold swapped tokens — they must be escrowed. Delta: %s", denomOutDelta)

	// Verify the escrow address holds the tokens
	escrowAddr := transfertypes.GetEscrowAddress(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
	)
	escrowBalance := s.GetBalance(s.ChainA, escrowAddr, denomOut)
	s.T().Logf("Escrow address %s balance: %s", escrowAddr, escrowBalance)
	s.Require().True(escrowBalance.Amount.GT(sdkmath.ZeroInt()),
		"escrow address should hold the transferred tokens")

	// Now relay the packet to destination
	err = s.TransferPath.RelayPacket(packet)
	s.Require().NoError(err)

	// Verify receiver got the IBC-denominated tokens on chain B
	ibcDenom := s.GetIBCDenom(s.TransferPath, denomOut)
	receiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	s.T().Logf("Receiver got %s of %s", receiverBalance.Amount, ibcDenom)
	s.Require().True(receiverBalance.Amount.GT(sdkmath.ZeroInt()),
		"receiver should have received IBC tokens on chain B")

	// The escrow balance should equal what the receiver got
	s.Require().Equal(escrowBalance.Amount, receiverBalance.Amount,
		"escrow amount on chain A should equal received amount on chain B")
}

// TestDirectIBCTransfer_EscrowE2E tests a standard IBC transfer to confirm
// the baseline escrow behavior works correctly end-to-end.
func (s *EscrowTestSuite) TestDirectIBCTransfer_EscrowE2E() {
	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := sdkmath.NewInt(1_000_000)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	initialSenderBalance := s.GetBalance(s.ChainA, sender, denom)
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

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Check escrow BEFORE relay
	escrowAddr := transfertypes.GetEscrowAddress(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
	)
	escrowBalance := s.GetBalance(s.ChainA, escrowAddr, denom)
	s.Require().True(escrowBalance.Amount.GTE(amount),
		"escrow should hold tokens before relay")

	// Sender should have lost the tokens
	senderAfterSend := s.GetBalance(s.ChainA, sender, denom)
	s.Require().Equal(initialSenderBalance.Amount.Sub(amount), senderAfterSend.Amount,
		"sender balance should decrease by transfer amount")

	// Relay and verify receiver
	err = s.TransferPath.RelayPacket(packet)
	s.Require().NoError(err)

	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)
	receiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	s.Require().Equal(amount, receiverBalance.Amount,
		"receiver should have received full amount")
}

// TestGAMMSwapWithIBCTransfer_TimeoutRefund verifies that when a GAMM
// swap+IBC transfer times out, the escrowed tokens are correctly refunded
// to the sender on the source chain.
func (s *EscrowTestSuite) TestGAMMSwapWithIBCTransfer_TimeoutRefund() {
	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denomIn := "ubadge"
	denomOut := "ustake"
	poolLiquidity := sdkmath.NewInt(10_000_000)
	swapAmount := sdkmath.NewInt(100_000)

	// Create pool
	poolID := s.createPoolOnChainA(denomIn, denomOut, poolLiquidity, poolLiquidity)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denomIn, swapAmount.MulRaw(2))))
	s.Require().NoError(err)

	// Record sender's denomOut balance and escrow balance before
	senderDenomOutBefore := s.GetBalance(s.ChainA, sender, denomOut)
	escrowAddr := transfertypes.GetEscrowAddress(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
	)
	escrowBefore := s.GetBalance(s.ChainA, escrowAddr, denomOut)

	// Use a timeout height that will expire quickly
	currentHeight := s.ChainB.GetContext().BlockHeight()
	timeoutHeight := clienttypes.NewHeight(1, uint64(currentHeight+1))

	// Do a swap first, then send the output via MsgTransfer with short timeout.
	// This exercises the same escrow code path that ExecuteIBCTransfer now uses
	// (transferKeeper.Transfer), while letting us control the timeout height.
	app := s.GetBitBadgesApp(s.ChainA)
	ctx := s.ChainA.GetContext()

	// Do the swap directly via keeper
	tokenOutAmount, err := app.PoolManagerKeeper.RouteExactAmountIn(
		ctx,
		sender,
		[]poolmanagertypes.SwapAmountInRoute{{PoolId: poolID, TokenOutDenom: denomOut}},
		sdk.NewCoin(denomIn, swapAmount),
		osmomath.NewInt(1),
		nil,
	)
	s.Require().NoError(err)
	s.T().Logf("Swap produced %s %s", tokenOutAmount, denomOut)

	// Commit the swap
	s.Coordinator.CommitBlock(s.ChainA)

	// Now send the swapped tokens via IBC with short timeout
	token := sdk.NewCoin(denomOut, tokenOutAmount)

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

	// Tokens should be escrowed now
	senderDenomOutAfterSend := s.GetBalance(s.ChainA, sender, denomOut)
	s.T().Logf("Sender %s after send: %s (before: %s)", denomOut, senderDenomOutAfterSend, senderDenomOutBefore)

	// Advance chain B past timeout
	for i := 0; i < 5; i++ {
		s.ChainB.NextBlock()
	}

	// Update client and trigger timeout
	err = s.TransferPath.EndpointA.UpdateClient()
	s.Require().NoError(err)

	err = s.TransferPath.EndpointA.TimeoutPacket(packet)
	s.Require().NoError(err)

	// After timeout: sender should get the tokens back (un-escrowed)
	senderDenomOutAfterTimeout := s.GetBalance(s.ChainA, sender, denomOut)
	refundDelta := senderDenomOutAfterTimeout.Amount.Sub(senderDenomOutBefore.Amount)
	s.T().Logf("Sender %s after timeout refund: %s (delta from start: %s)", denomOut, senderDenomOutAfterTimeout, refundDelta)

	// The refund should give back exactly the swapped amount
	s.Require().Equal(tokenOutAmount, refundDelta,
		"timeout should refund the full escrowed amount back to sender")

	// Escrow address should be back to its pre-test level (other tests may have
	// legitimately escrowed tokens on this channel that remain in escrow)
	escrowAfter := s.GetBalance(s.ChainA, escrowAddr, denomOut)
	escrowDelta := escrowAfter.Amount.Sub(escrowBefore.Amount)
	s.Require().True(escrowDelta.IsZero(),
		"escrow delta should be zero after timeout refund (all escrowed tokens returned)")
}
