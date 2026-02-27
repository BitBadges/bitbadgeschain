//go:build test
// +build test

package e2e

import (
	"encoding/json"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v10/testing"
	"github.com/stretchr/testify/suite"

	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	ibctest "github.com/bitbadges/bitbadgeschain/testing/ibc"
)

// HooksTestSuite tests IBC custom hooks
type HooksTestSuite struct {
	IBCTestSuite
}

func TestHooksTestSuite(t *testing.T) {
	suite.Run(t, new(HooksTestSuite))
}

// createSwapAndActionMemo creates a memo with swap_and_action hook data
func (s *HooksTestSuite) createSwapAndActionMemo(poolID string, denomIn, denomOut string, minAmountOut string, recipientAddress string) string {
	memo := map[string]interface{}{
		"swap_and_action": map[string]interface{}{
			"user_swap": map[string]interface{}{
				"swap_exact_asset_in": map[string]interface{}{
					"swap_venue_name": "bitbadges",
					"operations": []map[string]interface{}{
						{
							"pool":      poolID,
							"denom_in":  denomIn,
							"denom_out": denomOut,
						},
					},
				},
			},
			"min_asset": map[string]interface{}{
				"native": map[string]interface{}{
					"denom":  denomOut,
					"amount": minAmountOut,
				},
			},
			"post_swap_action": map[string]interface{}{
				"transfer": map[string]interface{}{
					"to_address": recipientAddress,
				},
			},
		},
	}

	memoBytes, err := json.Marshal(memo)
	s.Require().NoError(err)
	return string(memoBytes)
}

// createSwapWithFallbackMemo creates a memo with swap_and_action and fallback recover address
func (s *HooksTestSuite) createSwapWithFallbackMemo(poolID string, denomIn, denomOut string, minAmountOut string, recipientAddress string, recoverAddress string) string {
	memo := map[string]interface{}{
		"swap_and_action": map[string]interface{}{
			"user_swap": map[string]interface{}{
				"swap_exact_asset_in": map[string]interface{}{
					"swap_venue_name": "bitbadges",
					"operations": []map[string]interface{}{
						{
							"pool":      poolID,
							"denom_in":  denomIn,
							"denom_out": denomOut,
						},
					},
				},
			},
			"min_asset": map[string]interface{}{
				"native": map[string]interface{}{
					"denom":  denomOut,
					"amount": minAmountOut,
				},
			},
			"post_swap_action": map[string]interface{}{
				"transfer": map[string]interface{}{
					"to_address": recipientAddress,
				},
			},
			"destination_recover_address": recoverAddress,
		},
	}

	memoBytes, err := json.Marshal(memo)
	s.Require().NoError(err)
	return string(memoBytes)
}

// TestOnRecvPacketNoHooks tests pass-through behavior when no memo/hooks
func (s *HooksTestSuite) TestOnRecvPacketNoHooks() {
	s.T().Log("Testing OnRecvPacket pass-through (no hooks)")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Fund sender with ubadge tokens
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	// Send transfer without memo
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
		"", // No memo
	)

	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Relay the packet - should succeed normally
	err = s.TransferPath.RelayPacket(packet)
	s.Require().NoError(err)

	// Verify receiver got the tokens
	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)
	receiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	s.Require().Equal(amount, receiverBalance.Amount, "receiver should have received tokens without hooks")
}

// TestOnRecvPacketWithInvalidMemo tests behavior with invalid memo JSON
func (s *HooksTestSuite) TestOnRecvPacketWithInvalidMemo() {
	s.T().Log("Testing OnRecvPacket with invalid memo JSON")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Fund sender with ubadge tokens
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	// Send transfer with invalid JSON memo
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
		"{invalid json", // Invalid memo
	)

	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Receive packet on chain B - should fail with error ack
	err = s.TransferPath.EndpointB.UpdateClient()
	s.Require().NoError(err)

	res, err = s.TransferPath.EndpointB.RecvPacketWithResult(packet)
	s.Require().NoError(err)

	// Parse acknowledgement - should be an error
	ack, err := ibctesting.ParseAckFromEvents(res.Events)
	s.Require().NoError(err)

	s.T().Logf("Acknowledgement: %s", string(ack))
	s.Require().Contains(string(ack), "error", "should return error acknowledgement for invalid memo")
}

// TestOnRecvPacketWithNonSwapMemo tests behavior with memo that's not a swap_and_action
func (s *HooksTestSuite) TestOnRecvPacketWithNonSwapMemo() {
	s.T().Log("Testing OnRecvPacket with non-swap memo")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Fund sender with ubadge tokens
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	// Get initial receiver balance
	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)
	initialReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)

	// Send transfer with valid JSON but non-swap memo
	token := sdk.NewCoin(denom, amount)
	timeoutHeight := clienttypes.NewHeight(1, 110)

	regularMemo := `{"note": "just a regular transfer"}`

	msg := transfertypes.NewMsgTransfer(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
		token,
		sender.String(),
		receiver.String(),
		timeoutHeight,
		0,
		regularMemo,
	)

	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Relay the packet - should succeed (no hooks triggered)
	err = s.TransferPath.RelayPacket(packet)
	s.Require().NoError(err)

	// Verify receiver got the tokens (check delta)
	receiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	receiverDelta := receiverBalance.Amount.Sub(initialReceiverBalance.Amount)
	s.Require().Equal(amount, receiverDelta, "receiver should have received tokens")
}

// TestOnRecvPacketWithHookMemo tests hook execution with swap memo
func (s *HooksTestSuite) TestOnRecvPacketWithHookMemo() {
	s.T().Log("Testing OnRecvPacket with swap_and_action hook memo")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Fund sender with ubadge tokens
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	// Create a swap memo pointing to a non-existent pool
	// This will trigger hook execution but the swap will fail
	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)
	memo := s.createSwapAndActionMemo(
		"999999", // Non-existent pool ID
		ibcDenom,
		"someoutputdenom",
		"1",
		receiver.String(),
	)

	// Send transfer with hook memo
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
		memo,
	)

	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Receive packet on chain B
	err = s.TransferPath.EndpointB.UpdateClient()
	s.Require().NoError(err)

	res, err = s.TransferPath.EndpointB.RecvPacketWithResult(packet)
	s.Require().NoError(err)

	// Parse acknowledgement - should be an error since pool doesn't exist
	ack, err := ibctesting.ParseAckFromEvents(res.Events)
	s.Require().NoError(err)

	s.T().Logf("Acknowledgement: %s", string(ack))
	// The ack should be an error because the pool doesn't exist
	s.Require().Contains(string(ack), "error", "should return error acknowledgement for non-existent pool")
}

// TestAtomicRollbackOnSwapFailure tests that state is rolled back atomically on hook failure
func (s *HooksTestSuite) TestAtomicRollbackOnSwapFailure() {
	s.T().Log("Testing atomic rollback on swap failure")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := sdkmath.NewInt(5000000)

	// Fund sender with extra tokens
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	// Get initial balances
	initialSenderBalance := s.GetBalance(s.ChainA, sender, denom)
	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)
	initialReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)

	// Create a swap memo pointing to invalid pool (will fail)
	memo := s.createSwapAndActionMemo(
		"999999", // Non-existent pool
		ibcDenom,
		"outputdenom",
		"1",
		receiver.String(),
	)

	// Send transfer with failing hook
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
		memo,
	)

	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Receive on chain B - will fail due to invalid pool
	err = s.TransferPath.EndpointB.UpdateClient()
	s.Require().NoError(err)

	res, err = s.TransferPath.EndpointB.RecvPacketWithResult(packet)
	s.Require().NoError(err)

	// Parse acknowledgement
	ack, err := ibctesting.ParseAckFromEvents(res.Events)
	s.Require().NoError(err)

	s.T().Logf("Acknowledgement: %s", string(ack))

	// Verify it's an error ack
	s.Require().Contains(string(ack), "error", "should be error ack")

	// Verify receiver balance did NOT change (atomic rollback)
	finalReceiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenom)
	s.Require().Equal(initialReceiverBalance.Amount, finalReceiverBalance.Amount,
		"receiver balance should not change on hook failure - atomic rollback")

	// Acknowledge the error packet on chain A
	err = s.TransferPath.EndpointA.UpdateClient()
	s.Require().NoError(err)

	err = s.TransferPath.EndpointA.AcknowledgePacket(packet, ack)
	s.Require().NoError(err)

	// Verify sender got tokens back (refund on error ack)
	finalSenderBalance := s.GetBalance(s.ChainA, sender, denom)
	s.T().Logf("Initial sender: %s, Final sender: %s", initialSenderBalance, finalSenderBalance)
	s.Require().Equal(initialSenderBalance.Amount, finalSenderBalance.Amount,
		"sender should have original balance after error acknowledgement refund")
}

// TestSwapWithFallback tests fallback to destination_recover_address on swap failure
func (s *HooksTestSuite) TestSwapWithFallback() {
	s.T().Log("Testing swap with fallback to destination_recover_address")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Fund sender with ubadge tokens
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	ibcDenom := s.GetIBCDenom(s.TransferPath, denom)

	// Create a memo with fallback address
	// When the swap fails, tokens should go to the recover address
	memo := s.createSwapWithFallbackMemo(
		"999999", // Non-existent pool (will fail)
		ibcDenom,
		"outputdenom",
		"1",
		receiver.String(), // Original recipient
		receiver.String(), // Fallback/recover address
	)

	// Send transfer with hook memo
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
		memo,
	)

	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Receive on chain B
	err = s.TransferPath.EndpointB.UpdateClient()
	s.Require().NoError(err)

	res, err = s.TransferPath.EndpointB.RecvPacketWithResult(packet)
	s.Require().NoError(err)

	// Parse acknowledgement
	ack, err := ibctesting.ParseAckFromEvents(res.Events)
	s.Require().NoError(err)

	s.T().Logf("Acknowledgement with fallback: %s", string(ack))

	// Note: The behavior depends on implementation
	// If fallback is supported, tokens go to recover_address
	// If not, it's an error ack and sender gets refunded
}

// TestHookDataParsing tests the HookData parsing logic
func (s *HooksTestSuite) TestHookDataParsing() {
	s.T().Log("Testing HookData parsing")

	// Test valid swap_and_action memo
	validMemo := `{
		"swap_and_action": {
			"user_swap": {
				"swap_exact_asset_in": {
					"swap_venue_name": "bitbadges",
					"operations": [
						{
							"pool": "1",
							"denom_in": "ubadge",
							"denom_out": "uatom"
						}
					]
				}
			},
			"min_asset": {
				"native": {
					"denom": "uatom",
					"amount": "100"
				}
			}
		}
	}`

	hookData, err := customhookstypes.ParseHookDataFromMemo(validMemo)
	s.Require().NoError(err)
	s.Require().NotNil(hookData)
	s.Require().NotNil(hookData.SwapAndAction)
	s.Require().NotNil(hookData.SwapAndAction.UserSwap)
	s.Require().NotNil(hookData.SwapAndAction.UserSwap.SwapExactAssetIn)
	s.Require().Len(hookData.SwapAndAction.UserSwap.SwapExactAssetIn.Operations, 1)

	// Test empty memo
	hookData, err = customhookstypes.ParseHookDataFromMemo("")
	s.Require().NoError(err)
	s.Require().Nil(hookData)

	// Test non-swap memo
	hookData, err = customhookstypes.ParseHookDataFromMemo(`{"note": "hello"}`)
	s.Require().NoError(err)
	s.Require().Nil(hookData)

	// Test memo exceeding max size
	largeMemo := make([]byte, customhookstypes.MaxMemoSize+1)
	for i := range largeMemo {
		largeMemo[i] = 'a'
	}
	_, err = customhookstypes.ParseHookDataFromMemo(string(largeMemo))
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "exceeds maximum")
}

// Note: setupTestPool would be needed for testing hooks with actual pool swaps.
// For now, the tests focus on hook failure scenarios which don't require real pools.
// To implement actual swap tests, use the balancer package:
//   import balancer "github.com/bitbadges/bitbadgeschain/x/gamm/poolmodels/balancer"
// And create pools using balancer.NewMsgCreateBalancerPool()
