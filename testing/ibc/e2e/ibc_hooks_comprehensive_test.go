//go:build test
// +build test

package e2e

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v10/testing"
	"github.com/stretchr/testify/suite"

	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"

	ibctest "github.com/bitbadges/bitbadgeschain/testing/ibc"
	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	"github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
	ibchookstypes "github.com/bitbadges/bitbadgeschain/x/ibc-hooks/types"
)

// ComprehensiveHooksTestSuite tests the full swap_and_action hook pipeline
// including success paths, failure/recovery, and edge cases.
// Each test gets a fresh two-chain setup to avoid shared state issues.
type ComprehensiveHooksTestSuite struct {
	IBCTestSuite
}

func TestComprehensiveHooksTestSuite(t *testing.T) {
	suite.Run(t, new(ComprehensiveHooksTestSuite))
}

// SetupTest reinitializes chains and transfer path before each test
// to ensure tests are fully independent.
func (s *ComprehensiveHooksTestSuite) SetupTest() {
	s.SetupSuite()
}

// createPoolOnChainB creates a balancer pool with the given denoms on chain B
// and returns the pool ID. Funds the creator with the pool assets first.
func (s *ComprehensiveHooksTestSuite) createPoolOnChainB(denomA, denomB string, amountA, amountB sdkmath.Int) uint64 {
	app := s.GetBitBadgesApp(s.ChainB)
	ctx := s.ChainB.GetContext()
	creator := s.ChainB.SenderAccount.GetAddress()

	// Fund the creator with pool assets
	poolCoins := sdk.NewCoins(
		sdk.NewCoin(denomA, amountA),
		sdk.NewCoin(denomB, amountB),
	)
	err := ibctest.FundAccount(s.ChainB, creator, poolCoins)
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
	s.Coordinator.CommitBlock(s.ChainB)

	return poolID
}

// createPoolOnChainA creates a balancer pool with the given denoms on chain A
// and returns the pool ID.
func (s *ComprehensiveHooksTestSuite) createPoolOnChainA(denomA, denomB string, amountA, amountB sdkmath.Int) uint64 {
	app := s.GetBitBadgesApp(s.ChainA)
	ctx := s.ChainA.GetContext()
	creator := s.ChainA.SenderAccount.GetAddress()

	poolCoins := sdk.NewCoins(
		sdk.NewCoin(denomA, amountA),
		sdk.NewCoin(denomB, amountB),
	)
	err := ibctest.FundAccount(s.ChainA, creator, poolCoins)
	s.Require().NoError(err)

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

	s.Coordinator.CommitBlock(s.ChainA)
	return poolID
}

// sendIBCTransferWithMemo sends an IBC transfer from ChainA to ChainB with the given memo.
// Returns the packet and any error.
func (s *ComprehensiveHooksTestSuite) sendIBCTransferWithMemo(
	sender, receiver sdk.AccAddress,
	token sdk.Coin,
	memo string,
) (channeltypes.Packet, error) {
	timeoutHeight := clienttypes.NewHeight(1, 110)

	msg := transfertypes.NewMsgTransfer(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
		token,
		sender.String(),
		receiver.String(),
		timeoutHeight,
		0,
		memo,
	)

	res, err := s.ChainA.SendMsgs(msg)
	if err != nil {
		return channeltypes.Packet{}, err
	}

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	return packet, err
}

// recvAndParseAck receives a packet on ChainB and returns the raw ack bytes.
func (s *ComprehensiveHooksTestSuite) recvAndParseAck(packet channeltypes.Packet) []byte {
	err := s.TransferPath.EndpointB.UpdateClient()
	s.Require().NoError(err)

	res, err := s.TransferPath.EndpointB.RecvPacketWithResult(packet)
	s.Require().NoError(err)

	ack, err := ibctesting.ParseAckFromEvents(res.Events)
	s.Require().NoError(err)
	return ack
}

// buildSwapAndActionMemo constructs a swap_and_action memo JSON with optional fields.
func buildSwapAndActionMemo(operations []map[string]interface{}, minAssetDenom, minAssetAmount string, postSwapAction map[string]interface{}, extras map[string]interface{}) string {
	swapAndAction := map[string]interface{}{
		"user_swap": map[string]interface{}{
			"swap_exact_asset_in": map[string]interface{}{
				"swap_venue_name": "bitbadges",
				"operations":      operations,
			},
		},
		"min_asset": map[string]interface{}{
			"native": map[string]interface{}{
				"denom":  minAssetDenom,
				"amount": minAssetAmount,
			},
		},
		"post_swap_action": postSwapAction,
	}
	for k, v := range extras {
		swapAndAction[k] = v
	}

	memo := map[string]interface{}{
		"swap_and_action": swapAndAction,
	}
	b, _ := json.Marshal(memo)
	return string(b)
}

// deriveIntermediateReceiver computes the intermediate sender address that
// the IBC hooks middleware derives for a given original sender and dest channel.
// For hooks to work, the IBC packet receiver MUST be this address so that
// the IBC transfer delivers tokens to the address the hook will swap from.
func (s *ComprehensiveHooksTestSuite) deriveIntermediateReceiver(originalSender string) string {
	destChannel := s.TransferPath.EndpointB.ChannelID
	bech32Prefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	intermediate, err := ibchookstypes.DeriveIntermediateSender(destChannel, originalSender, bech32Prefix)
	s.Require().NoError(err)
	return intermediate
}

// ---------------------------------------------------------------------------
// SUCCESS PATH TESTS
// ---------------------------------------------------------------------------

// TestSuccessfulSwapAndTransfer sends tokens from A to B with a swap_and_action
// memo. On ChainB a pool exists between the IBC denom and a local denom. The
// hook swaps the received IBC tokens and locally transfers the output to the
// recipient. Verifies recipient gets the swapped tokens.
func (s *ComprehensiveHooksTestSuite) TestSuccessfulSwapAndTransfer() {
	s.T().Log("=== TestSuccessfulSwapAndTransfer ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	outputDenom := "ustake"
	transferAmount := sdkmath.NewInt(500_000)
	poolLiquidity := sdkmath.NewInt(10_000_000)

	// Step 1: Compute the IBC denom that will exist on Chain B for ubadge
	ibcDenom := s.GetIBCDenom(s.TransferPath, baseDenom)
	s.T().Logf("IBC denom on Chain B: %s", ibcDenom)

	// Step 2: Create a pool on Chain B between the IBC denom and ustake
	s.T().Log("Creating pool on Chain B: ibcDenom/ustake")
	poolID := s.createPoolOnChainB(ibcDenom, outputDenom, poolLiquidity, poolLiquidity)
	s.T().Logf("Created pool %d on Chain B", poolID)

	// Step 3: Fund sender on Chain A
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(3))))
	s.Require().NoError(err)

	// Step 4: Record receiver's ustake balance before
	receiverStakeBefore := s.GetBalance(s.ChainB, receiver, outputDenom)
	s.T().Logf("Receiver %s balance before: %s", outputDenom, receiverStakeBefore)

	// Step 5: Build the swap_and_action memo
	memo := buildSwapAndActionMemo(
		[]map[string]interface{}{
			{"pool": fmt.Sprintf("%d", poolID), "denom_in": ibcDenom, "denom_out": outputDenom},
		},
		outputDenom, "1",
		map[string]interface{}{
			"transfer": map[string]interface{}{
				"to_address": receiver.String(),
			},
		},
		nil,
	)
	s.T().Logf("Memo: %s", memo)

	// Step 6: The IBC packet receiver must be the intermediate sender so the
	// IBC transfer delivers tokens to the address the hook will swap from.
	intermediateReceiver := s.deriveIntermediateReceiver(sender.String())
	intermediateAddr, err := sdk.AccAddressFromBech32(intermediateReceiver)
	s.Require().NoError(err)
	s.T().Logf("Intermediate receiver: %s", intermediateReceiver)

	// Step 7: Send IBC transfer with hook memo to the intermediate address
	s.T().Log("Sending IBC transfer A->B with swap_and_action memo")
	packet, err := s.sendIBCTransferWithMemo(sender, intermediateAddr, sdk.NewCoin(baseDenom, transferAmount), memo)
	s.Require().NoError(err)

	// Step 8: Receive on Chain B and parse ack
	ack := s.recvAndParseAck(packet)
	s.T().Logf("Acknowledgement: %s", string(ack))

	// Step 9: Assert success ack (no "error" in ack)
	s.Require().NotContains(string(ack), "error",
		"hook should succeed: swap + local transfer")

	// Step 10: Verify receiver got the swapped ustake tokens
	receiverStakeAfter := s.GetBalance(s.ChainB, receiver, outputDenom)
	stakeGained := receiverStakeAfter.Amount.Sub(receiverStakeBefore.Amount)
	s.T().Logf("Receiver %s gained: %s", outputDenom, stakeGained)
	s.Require().True(stakeGained.IsPositive(),
		"receiver should have gained ustake tokens from the swap")
}

// TestSuccessfulSwapAndIBCForward sends tokens A->B with swap_and_action
// that performs a swap on B and then forwards the output via IBC transfer
// (back to A or to another chain). Verifies the IBC packet is created.
func (s *ComprehensiveHooksTestSuite) TestSuccessfulSwapAndIBCForward() {
	s.T().Log("=== TestSuccessfulSwapAndIBCForward ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()
	finalReceiver := sender // forward back to Chain A

	baseDenom := "ubadge"
	// The output denom must be ICS20-compatible and have a channel to forward on.
	// We use a native denom on ChainB that we'll swap into, then IBC forward it.
	outputDenom := "ustake"
	transferAmount := sdkmath.NewInt(500_000)
	poolLiquidity := sdkmath.NewInt(10_000_000)

	ibcDenom := s.GetIBCDenom(s.TransferPath, baseDenom)

	// Create pool on Chain B
	poolID := s.createPoolOnChainB(ibcDenom, outputDenom, poolLiquidity, poolLiquidity)
	s.T().Logf("Created pool %d on Chain B", poolID)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(3))))
	s.Require().NoError(err)

	// The post_swap_action is ibc_transfer back to Chain A.
	// On Chain B, the channel back to A is EndpointB's channel.
	// We also need a recover_address (required by validation).
	memo := buildSwapAndActionMemo(
		[]map[string]interface{}{
			{"pool": fmt.Sprintf("%d", poolID), "denom_in": ibcDenom, "denom_out": outputDenom},
		},
		outputDenom, "1",
		map[string]interface{}{
			"ibc_transfer": map[string]interface{}{
				"ibc_info": map[string]interface{}{
					"source_channel":  s.TransferPath.EndpointB.ChannelID,
					"receiver":        finalReceiver.String(),
					"recover_address": receiver.String(),
				},
			},
		},
		nil,
	)

	// IBC packet receiver must be the intermediate sender
	intermediateReceiver := s.deriveIntermediateReceiver(sender.String())
	intermediateAddr, err := sdk.AccAddressFromBech32(intermediateReceiver)
	s.Require().NoError(err)

	s.T().Log("Sending IBC transfer A->B with swap + IBC forward memo")
	packet, err := s.sendIBCTransferWithMemo(sender, intermediateAddr, sdk.NewCoin(baseDenom, transferAmount), memo)
	s.Require().NoError(err)

	// Receive on Chain B
	ack := s.recvAndParseAck(packet)
	s.T().Logf("Acknowledgement: %s", string(ack))

	// The ack should be success if the swap and IBC forward packet creation both worked
	s.Require().NotContains(string(ack), "error",
		"hook should succeed: swap + IBC forward")

	// We verify the escrow on Chain B: ustake should be escrowed for the forward transfer
	escrowAddr := transfertypes.GetEscrowAddress(
		s.TransferPath.EndpointB.ChannelConfig.PortID,
		s.TransferPath.EndpointB.ChannelID,
	)
	escrowBalance := s.GetBalance(s.ChainB, escrowAddr, outputDenom)
	s.T().Logf("Escrow on Chain B for %s: %s", outputDenom, escrowBalance)
	s.Require().True(escrowBalance.Amount.IsPositive(),
		"escrow should hold ustake tokens for the IBC forward transfer")
}

// TestSuccessfulMultiHopSwap sends tokens A->B with a swap_and_action that
// chains two swap operations: pool1 (ibcDenom->intermediary) then pool2
// (intermediary->finalDenom). Verifies recipient gets the final denom.
func (s *ComprehensiveHooksTestSuite) TestSuccessfulMultiHopSwap() {
	s.T().Log("=== TestSuccessfulMultiHopSwap ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	intermediaryDenom := "umid"
	finalDenom := "ufinal"
	transferAmount := sdkmath.NewInt(500_000)
	poolLiquidity := sdkmath.NewInt(10_000_000)

	ibcDenom := s.GetIBCDenom(s.TransferPath, baseDenom)

	// Create two pools on Chain B
	s.T().Log("Creating pool 1 on Chain B: ibcDenom -> umid")
	pool1 := s.createPoolOnChainB(ibcDenom, intermediaryDenom, poolLiquidity, poolLiquidity)

	s.T().Log("Creating pool 2 on Chain B: umid -> ufinal")
	pool2 := s.createPoolOnChainB(intermediaryDenom, finalDenom, poolLiquidity, poolLiquidity)
	s.T().Logf("Created pools %d and %d on Chain B", pool1, pool2)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(3))))
	s.Require().NoError(err)

	receiverFinalBefore := s.GetBalance(s.ChainB, receiver, finalDenom)

	// Build memo with two operations
	memo := buildSwapAndActionMemo(
		[]map[string]interface{}{
			{"pool": fmt.Sprintf("%d", pool1), "denom_in": ibcDenom, "denom_out": intermediaryDenom},
			{"pool": fmt.Sprintf("%d", pool2), "denom_in": intermediaryDenom, "denom_out": finalDenom},
		},
		finalDenom, "1",
		map[string]interface{}{
			"transfer": map[string]interface{}{
				"to_address": receiver.String(),
			},
		},
		nil,
	)

	intermediateReceiver := s.deriveIntermediateReceiver(sender.String())
	intermediateAddr, err := sdk.AccAddressFromBech32(intermediateReceiver)
	s.Require().NoError(err)

	s.T().Log("Sending IBC transfer A->B with multi-hop swap memo")
	packet, err := s.sendIBCTransferWithMemo(sender, intermediateAddr, sdk.NewCoin(baseDenom, transferAmount), memo)
	s.Require().NoError(err)

	ack := s.recvAndParseAck(packet)
	s.T().Logf("Acknowledgement: %s", string(ack))

	s.Require().NotContains(string(ack), "error",
		"multi-hop swap should succeed")

	receiverFinalAfter := s.GetBalance(s.ChainB, receiver, finalDenom)
	gained := receiverFinalAfter.Amount.Sub(receiverFinalBefore.Amount)
	s.T().Logf("Receiver gained %s %s", gained, finalDenom)
	s.Require().True(gained.IsPositive(),
		"receiver should have gained ufinal tokens from multi-hop swap")
}

// ---------------------------------------------------------------------------
// FAILURE + RECOVERY TESTS
// ---------------------------------------------------------------------------

// TestSwapFailsWithFallback_TokensGoToRecoverAddress tests that when a swap
// fails and destination_recover_address is set, the original IBC tokens are
// sent to the recover address (success ack, no rollback).
func (s *ComprehensiveHooksTestSuite) TestSwapFailsWithFallback_TokensGoToRecoverAddress() {
	s.T().Log("=== TestSwapFailsWithFallback_TokensGoToRecoverAddress ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	transferAmount := sdkmath.NewInt(500_000)
	ibcDenom := s.GetIBCDenom(s.TransferPath, baseDenom)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(3))))
	s.Require().NoError(err)

	// Record receiver's IBC denom balance before (the recover address is receiver)
	receiverIBCBefore := s.GetBalance(s.ChainB, receiver, ibcDenom)

	// Build memo: non-existent pool (will fail), but with destination_recover_address
	memo := buildSwapAndActionMemo(
		[]map[string]interface{}{
			{"pool": "999999", "denom_in": ibcDenom, "denom_out": "nonexistent"},
		},
		"nonexistent", "1",
		map[string]interface{}{
			"transfer": map[string]interface{}{
				"to_address": receiver.String(),
			},
		},
		map[string]interface{}{
			"destination_recover_address": receiver.String(),
		},
	)

	// Use intermediate receiver so IBC tokens land at the hook address
	intermediateReceiver := s.deriveIntermediateReceiver(sender.String())
	intermediateAddr, err := sdk.AccAddressFromBech32(intermediateReceiver)
	s.Require().NoError(err)

	s.T().Log("Sending IBC transfer with swap that will fail + recover address set")
	packet, err := s.sendIBCTransferWithMemo(sender, intermediateAddr, sdk.NewCoin(baseDenom, transferAmount), memo)
	s.Require().NoError(err)

	ack := s.recvAndParseAck(packet)
	s.T().Logf("Acknowledgement: %s", string(ack))

	// With a recover address, the hook should return SUCCESS ack (tokens go to recover)
	// The swap fails, but fallback sends IBC tokens to the recover address
	s.Require().NotContains(string(ack), "error",
		"with destination_recover_address set, should return success ack (fallback)")

	// Verify recover address (receiver) got the original IBC tokens
	receiverIBCAfter := s.GetBalance(s.ChainB, receiver, ibcDenom)
	ibcGained := receiverIBCAfter.Amount.Sub(receiverIBCBefore.Amount)
	s.T().Logf("Recover address gained %s of %s", ibcGained, ibcDenom)
	s.Require().True(ibcGained.IsPositive(),
		"recover address should have received the original IBC tokens")
}

// TestSwapFailsWithoutFallback_AtomicRollback tests that when the swap fails
// and no recover address is set, the entire IBC receive is rolled back (error
// ack) and the sender on Chain A gets refunded.
func (s *ComprehensiveHooksTestSuite) TestSwapFailsWithoutFallback_AtomicRollback() {
	s.T().Log("=== TestSwapFailsWithoutFallback_AtomicRollback ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	transferAmount := sdkmath.NewInt(500_000)
	ibcDenom := s.GetIBCDenom(s.TransferPath, baseDenom)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(3))))
	s.Require().NoError(err)

	initialSenderBalance := s.GetBalance(s.ChainA, sender, baseDenom)
	initialReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)

	// Build memo: non-existent pool, NO recover address
	memo := buildSwapAndActionMemo(
		[]map[string]interface{}{
			{"pool": "999999", "denom_in": ibcDenom, "denom_out": "nonexistent"},
		},
		"nonexistent", "1",
		map[string]interface{}{
			"transfer": map[string]interface{}{
				"to_address": receiver.String(),
			},
		},
		nil, // no extras, no recover address
	)

	s.T().Log("Sending IBC transfer with swap that will fail, no recover address")
	packet, err := s.sendIBCTransferWithMemo(sender, receiver, sdk.NewCoin(baseDenom, transferAmount), memo)
	s.Require().NoError(err)

	ack := s.recvAndParseAck(packet)
	s.T().Logf("Acknowledgement: %s", string(ack))

	// Should be error ack
	s.Require().Contains(string(ack), "error",
		"without recover address, failed swap should return error ack")

	// Receiver on Chain B should NOT have gotten any tokens
	finalReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	s.Require().Equal(initialReceiverBalance.Amount, finalReceiverBalance.Amount,
		"receiver balance should not change on atomic rollback")

	// Acknowledge error packet on Chain A so sender gets refund
	err = s.TransferPath.EndpointA.UpdateClient()
	s.Require().NoError(err)
	err = s.TransferPath.EndpointA.AcknowledgePacket(packet, ack)
	s.Require().NoError(err)

	// Sender should have original balance restored
	finalSenderBalance := s.GetBalance(s.ChainA, sender, baseDenom)
	s.T().Logf("Sender balance: initial=%s, final=%s", initialSenderBalance, finalSenderBalance)
	s.Require().Equal(initialSenderBalance.Amount, finalSenderBalance.Amount,
		"sender should be refunded after error ack")
}

// TestSwapFailsSlippageExceeded creates a real pool but sets min_asset very
// high so the swap output does not meet the slippage requirement.
func (s *ComprehensiveHooksTestSuite) TestSwapFailsSlippageExceeded() {
	s.T().Log("=== TestSwapFailsSlippageExceeded ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	outputDenom := "ustake"
	transferAmount := sdkmath.NewInt(100_000) // relatively small swap
	poolLiquidity := sdkmath.NewInt(10_000_000)

	ibcDenom := s.GetIBCDenom(s.TransferPath, baseDenom)

	// Create pool on Chain B
	poolID := s.createPoolOnChainB(ibcDenom, outputDenom, poolLiquidity, poolLiquidity)
	s.T().Logf("Created pool %d on Chain B", poolID)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(3))))
	s.Require().NoError(err)

	// Build memo with min_asset set absurdly high -- much more than pool can output
	memo := buildSwapAndActionMemo(
		[]map[string]interface{}{
			{"pool": fmt.Sprintf("%d", poolID), "denom_in": ibcDenom, "denom_out": outputDenom},
		},
		outputDenom, "999999999999999", // impossibly high min output
		map[string]interface{}{
			"transfer": map[string]interface{}{
				"to_address": receiver.String(),
			},
		},
		nil,
	)

	intermediateReceiver := s.deriveIntermediateReceiver(sender.String())
	intermediateAddr, err := sdk.AccAddressFromBech32(intermediateReceiver)
	s.Require().NoError(err)

	s.T().Log("Sending IBC transfer with swap that will fail slippage check")
	packet, err := s.sendIBCTransferWithMemo(sender, intermediateAddr, sdk.NewCoin(baseDenom, transferAmount), memo)
	s.Require().NoError(err)

	ack := s.recvAndParseAck(packet)
	s.T().Logf("Acknowledgement: %s", string(ack))

	// Should be error ack because slippage exceeded
	s.Require().Contains(string(ack), "error",
		"swap should fail because min_asset exceeds what pool can output")
}

// ---------------------------------------------------------------------------
// EDGE CASE TESTS
// ---------------------------------------------------------------------------

// TestZeroAmountTransfer sends a minimal (1) amount with a hook memo.
func (s *ComprehensiveHooksTestSuite) TestZeroAmountTransfer() {
	s.T().Log("=== TestZeroAmountTransfer ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	outputDenom := "ustake"
	poolLiquidity := sdkmath.NewInt(10_000_000)

	ibcDenom := s.GetIBCDenom(s.TransferPath, baseDenom)

	// Create pool on Chain B
	poolID := s.createPoolOnChainB(ibcDenom, outputDenom, poolLiquidity, poolLiquidity)

	// Fund sender with a small amount
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, sdkmath.NewInt(100))))
	s.Require().NoError(err)

	memo := buildSwapAndActionMemo(
		[]map[string]interface{}{
			{"pool": fmt.Sprintf("%d", poolID), "denom_in": ibcDenom, "denom_out": outputDenom},
		},
		outputDenom, "1",
		map[string]interface{}{
			"transfer": map[string]interface{}{
				"to_address": receiver.String(),
			},
		},
		nil,
	)

	// Send with amount = 1 (minimal) to intermediate receiver
	intermediateReceiver := s.deriveIntermediateReceiver(sender.String())
	intermediateAddr, err := sdk.AccAddressFromBech32(intermediateReceiver)
	s.Require().NoError(err)

	s.T().Log("Sending IBC transfer with amount=1")
	packet, err := s.sendIBCTransferWithMemo(sender, intermediateAddr, sdk.NewCoin(baseDenom, sdkmath.NewInt(1)), memo)
	s.Require().NoError(err)

	ack := s.recvAndParseAck(packet)
	s.T().Logf("Acknowledgement for 1-token transfer: %s", string(ack))

	// The swap may succeed or fail depending on pool math for amount=1.
	// We just verify the chain does not panic and produces a valid ack.
	s.Require().True(len(ack) > 0, "should produce a valid acknowledgement")
}

// TestHookWithEmptyOperations sends a swap_and_action with an empty operations array.
func (s *ComprehensiveHooksTestSuite) TestHookWithEmptyOperations() {
	s.T().Log("=== TestHookWithEmptyOperations ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	transferAmount := sdkmath.NewInt(500_000)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(3))))
	s.Require().NoError(err)

	// Build memo with empty operations
	memo := buildSwapAndActionMemo(
		[]map[string]interface{}{}, // empty operations
		"ustake", "1",
		map[string]interface{}{
			"transfer": map[string]interface{}{
				"to_address": receiver.String(),
			},
		},
		nil,
	)

	s.T().Log("Sending IBC transfer with empty operations array")
	packet, err := s.sendIBCTransferWithMemo(sender, receiver, sdk.NewCoin(baseDenom, transferAmount), memo)
	s.Require().NoError(err)

	ack := s.recvAndParseAck(packet)
	s.T().Logf("Acknowledgement: %s", string(ack))

	// Should be error ack because the keeper checks "no operations provided for swap"
	s.Require().Contains(string(ack), "error",
		"empty operations should produce error ack")
}

// TestHookMemoTooLarge sends a memo exceeding MaxMemoSize (64KB).
func (s *ComprehensiveHooksTestSuite) TestHookMemoTooLarge() {
	s.T().Log("=== TestHookMemoTooLarge ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	transferAmount := sdkmath.NewInt(500_000)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(3))))
	s.Require().NoError(err)

	// Build a memo that exceeds MaxMemoSize (64KB)
	// We create a valid JSON structure but pad it with a very long string value
	padSize := customhookstypes.MaxMemoSize + 1000
	padding := strings.Repeat("x", padSize)
	largeMemo := fmt.Sprintf(`{"swap_and_action":{"padding":"%s"}}`, padding)
	s.T().Logf("Large memo size: %d bytes (max: %d)", len(largeMemo), customhookstypes.MaxMemoSize)

	s.T().Log("Sending IBC transfer with oversized memo")
	_, err = s.sendIBCTransferWithMemo(sender, receiver, sdk.NewCoin(baseDenom, transferAmount), largeMemo)
	// The IBC transfer module rejects the memo at the send side (ValidateBasic)
	s.Require().Error(err, "oversized memo should be rejected at send time")
}

// TestMultipleTransfersWithHooksSameBlock sends multiple IBC transfers with
// hook memos and relays them all. Verifies all produce valid acks.
func (s *ComprehensiveHooksTestSuite) TestMultipleTransfersWithHooksSameBlock() {
	s.T().Log("=== TestMultipleTransfersWithHooksSameBlock ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	outputDenom := "ustake"
	transferAmount := sdkmath.NewInt(100_000)
	poolLiquidity := sdkmath.NewInt(50_000_000) // large pool for multiple swaps

	ibcDenom := s.GetIBCDenom(s.TransferPath, baseDenom)

	// Create pool on Chain B
	poolID := s.createPoolOnChainB(ibcDenom, outputDenom, poolLiquidity, poolLiquidity)
	s.T().Logf("Created pool %d on Chain B", poolID)

	// Fund sender with enough for multiple transfers
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(20))))
	s.Require().NoError(err)

	memo := buildSwapAndActionMemo(
		[]map[string]interface{}{
			{"pool": fmt.Sprintf("%d", poolID), "denom_in": ibcDenom, "denom_out": outputDenom},
		},
		outputDenom, "1",
		map[string]interface{}{
			"transfer": map[string]interface{}{
				"to_address": receiver.String(),
			},
		},
		nil,
	)

	// Use intermediate receiver for hook execution
	intermediateReceiver := s.deriveIntermediateReceiver(sender.String())

	numTransfers := 3
	packets := make([]channeltypes.Packet, numTransfers)

	s.T().Logf("Sending %d IBC transfers with hook memos", numTransfers)
	for i := 0; i < numTransfers; i++ {
		timeoutHeight := clienttypes.NewHeight(1, 110)
		msg := transfertypes.NewMsgTransfer(
			s.TransferPath.EndpointA.ChannelConfig.PortID,
			s.TransferPath.EndpointA.ChannelID,
			sdk.NewCoin(baseDenom, transferAmount),
			sender.String(),
			intermediateReceiver,
			timeoutHeight,
			0,
			memo,
		)
		res, sendErr := s.ChainA.SendMsgs(msg)
		s.Require().NoError(sendErr, "send %d should succeed", i)

		pkt, parseErr := ibctesting.ParsePacketFromEvents(res.Events)
		s.Require().NoError(parseErr, "parse packet %d", i)
		packets[i] = pkt
	}

	// Receive all packets on Chain B
	receiverStakeBefore := s.GetBalance(s.ChainB, receiver, outputDenom)
	for i, pkt := range packets {
		ack := s.recvAndParseAck(pkt)
		s.T().Logf("Packet %d ack: %s", i, string(ack))
		s.Require().NotContains(string(ack), "error",
			"packet %d should succeed", i)
	}

	receiverStakeAfter := s.GetBalance(s.ChainB, receiver, outputDenom)
	totalGained := receiverStakeAfter.Amount.Sub(receiverStakeBefore.Amount)
	s.T().Logf("Total %s gained from %d transfers: %s", outputDenom, numTransfers, totalGained)
	s.Require().True(totalGained.IsPositive(),
		"receiver should have gained tokens from multiple hook transfers")
}

// TestHookWithTransferTokensAction tests the transfer_tokens action type,
// which triggers a MsgTransferTokens on receive instead of a swap.
// Since this requires a valid tokenization collection, we test the validation
// error path (collection not found) to verify the hook is triggered.
func (s *ComprehensiveHooksTestSuite) TestHookWithTransferTokensAction() {
	s.T().Log("=== TestHookWithTransferTokensAction ===")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	baseDenom := "ubadge"
	transferAmount := sdkmath.NewInt(500_000)

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(baseDenom, transferAmount.MulRaw(3))))
	s.Require().NoError(err)

	// Build a transfer_tokens memo (not swap_and_action)
	// Use a non-existent collection to verify the hook path is triggered
	transferTokensMemo := map[string]interface{}{
		"transfer_tokens": map[string]interface{}{
			"collection_id":  "999999",
			"fail_on_error":  true,
			"transfers": []map[string]interface{}{
				{
					"from":          "Mint",
					"toAddresses":   []string{receiver.String()},
					"balances":      []interface{}{},
				},
			},
		},
	}

	memoBytes, err := json.Marshal(transferTokensMemo)
	s.Require().NoError(err)

	s.T().Log("Sending IBC transfer with transfer_tokens hook memo")
	packet, err := s.sendIBCTransferWithMemo(sender, receiver, sdk.NewCoin(baseDenom, transferAmount), string(memoBytes))
	s.Require().NoError(err)

	ack := s.recvAndParseAck(packet)
	s.T().Logf("Acknowledgement: %s", string(ack))

	// Should fail because collection 999999 doesn't exist, but this proves
	// the transfer_tokens hook path was triggered (not swap_and_action)
	s.Require().Contains(string(ack), "error",
		"should fail because collection does not exist")
	s.Require().Contains(string(ack), "transfer_tokens",
		"error should reference transfer_tokens indicating that hook path was invoked")
}
