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

	ibctest "github.com/bitbadges/bitbadgeschain/testing/ibc"
)

// DenomResolutionTestSuite verifies that custom hooks resolve IBC denoms
// correctly for both foreign tokens arriving and native tokens returning.
type DenomResolutionTestSuite struct {
	IBCTestSuite
}

func TestDenomResolutionTestSuite(t *testing.T) {
	suite.Run(t, new(DenomResolutionTestSuite))
}

// TestForeignTokenDenomResolution tests that when a foreign token arrives via IBC
// with a hook memo, the hook correctly resolves the local ibc/HASH denom.
// Chain A sends "ubadge" to Chain B — on Chain B this is a foreign token
// that becomes ibc/HASH(transfer/channelB/ubadge).
func (s *DenomResolutionTestSuite) TestForeignTokenDenomResolution() {
	s.T().Log("Testing foreign token denom resolution in custom hooks")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Fund sender
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(2))))
	s.Require().NoError(err)

	// Compute what the IBC denom should be on Chain B
	expectedIBCDenom := s.GetIBCDenom(s.TransferPath, denom)
	s.T().Logf("Expected IBC denom on Chain B: %s", expectedIBCDenom)

	// Create a hook memo that references the IBC denom as denom_in
	// The swap will fail (no pool), but we verify the hook correctly identifies
	// the denom by checking the error message references the right denom
	memo := map[string]interface{}{
		"swap_and_action": map[string]interface{}{
			"user_swap": map[string]interface{}{
				"swap_exact_asset_in": map[string]interface{}{
					"swap_venue_name": "bitbadges",
					"operations": []map[string]interface{}{
						{
							"pool":      "1",
							"denom_in":  expectedIBCDenom,
							"denom_out": "someoutput",
						},
					},
				},
			},
			"min_asset": map[string]interface{}{
				"native": map[string]interface{}{
					"denom":  "someoutput",
					"amount": "1",
				},
			},
			"post_swap_action": map[string]interface{}{
				"transfer": map[string]interface{}{
					"to_address": receiver.String(),
				},
			},
		},
	}
	memoBytes, err := json.Marshal(memo)
	s.Require().NoError(err)

	// Send the IBC transfer with hook memo
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
		string(memoBytes),
	)

	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Receive on Chain B
	err = s.TransferPath.EndpointB.UpdateClient()
	s.Require().NoError(err)

	res, err = s.TransferPath.EndpointB.RecvPacketWithResult(packet)
	s.Require().NoError(err)

	// Parse ack — the hook should have executed (and failed due to no pool),
	// but critically it should NOT have failed with a "denom mismatch" error.
	// A denom mismatch would mean extractDenomFromPacketOnRecv produced wrong denom.
	ack, err := ibctesting.ParseAckFromEvents(res.Events)
	s.Require().NoError(err)
	s.T().Logf("Ack: %s", string(ack))

	// The error should be about the pool/swap failing, NOT about denom validation
	s.Require().Contains(string(ack), "error", "should fail due to no pool")
	s.Require().NotContains(string(ack), "denom mismatch",
		"should NOT fail due to denom mismatch — extractDenomFromPacketOnRecv should resolve correctly")
}

// TestNativeTokenReturningDenomResolution tests that when a native token
// returns home via IBC with a hook memo, the hook resolves to the base denom.
// Chain A sends "ubadge" to Chain B, then Chain B sends ibc/HASH back to Chain A.
// On Chain A, this should resolve back to "ubadge".
func (s *DenomResolutionTestSuite) TestNativeTokenReturningDenomResolution() {
	s.T().Log("Testing native token returning home denom resolution")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := s.DefaultTransferAmount()

	// Step 1: Send ubadge from A to B (creates IBC voucher on B)
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(3))))
	s.Require().NoError(err)

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
	err = s.TransferPath.RelayPacket(packet)
	s.Require().NoError(err)

	// Verify receiver on B has the IBC tokens
	ibcDenomOnB := s.GetIBCDenom(s.TransferPath, denom)
	receiverBalance := s.GetBalance(s.ChainB, receiver, ibcDenomOnB)
	s.Require().Equal(amount, receiverBalance.Amount, "chain B should have IBC tokens")

	// Step 2: Send the IBC tokens back from B to A with a hook memo
	// On chain A, these should resolve to the native "ubadge" denom
	// The hook references "ubadge" as denom_in (the native denom on chain A)
	memo := map[string]interface{}{
		"swap_and_action": map[string]interface{}{
			"user_swap": map[string]interface{}{
				"swap_exact_asset_in": map[string]interface{}{
					"swap_venue_name": "bitbadges",
					"operations": []map[string]interface{}{
						{
							"pool":      "1",
							"denom_in":  denom, // native denom on chain A
							"denom_out": "someoutput",
						},
					},
				},
			},
			"min_asset": map[string]interface{}{
				"native": map[string]interface{}{
					"denom":  "someoutput",
					"amount": "1",
				},
			},
			"post_swap_action": map[string]interface{}{
				"transfer": map[string]interface{}{
					"to_address": sender.String(),
				},
			},
		},
	}
	memoBytes, err := json.Marshal(memo)
	s.Require().NoError(err)

	// Send from B back to A
	returnMsg := transfertypes.NewMsgTransfer(
		s.TransferPath.EndpointB.ChannelConfig.PortID,
		s.TransferPath.EndpointB.ChannelID,
		sdk.NewCoin(ibcDenomOnB, amount),
		receiver.String(),
		sender.String(),
		timeoutHeight,
		0,
		string(memoBytes),
	)
	res, err = s.ChainB.SendMsgs(returnMsg)
	s.Require().NoError(err)

	returnPacket, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)

	// Receive on Chain A
	err = s.TransferPath.EndpointA.UpdateClient()
	s.Require().NoError(err)

	res, err = s.TransferPath.EndpointA.RecvPacketWithResult(returnPacket)
	s.Require().NoError(err)

	// Parse ack
	ack, err := ibctesting.ParseAckFromEvents(res.Events)
	s.Require().NoError(err)
	s.T().Logf("Return ack: %s", string(ack))

	// Should fail due to no pool, NOT due to denom mismatch
	s.Require().Contains(string(ack), "error", "should fail due to no pool")
	s.Require().NotContains(string(ack), "denom mismatch",
		"should NOT fail due to denom mismatch — native token should resolve to base denom")
}

// TestTransferDenomConsistency verifies that the denom the IBC transfer module
// creates matches what our hook would use, by doing a plain transfer first and
// checking the balance denom, then verifying the hook references the same denom.
func (s *DenomResolutionTestSuite) TestTransferDenomConsistency() {
	s.T().Log("Testing denom consistency between IBC transfer module and hooks")

	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()

	denom := "ubadge"
	amount := sdkmath.NewInt(500000)

	// Fund and send a plain transfer (no hooks) to establish the IBC denom
	err := ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, amount.MulRaw(4))))
	s.Require().NoError(err)

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
	err = s.TransferPath.RelayPacket(packet)
	s.Require().NoError(err)

	// Get the actual IBC denom that the transfer module created on Chain B
	expectedDenom := s.GetIBCDenom(s.TransferPath, denom)
	balance := s.GetBalance(s.ChainB, receiver, expectedDenom)
	s.Require().True(balance.Amount.IsPositive(),
		"receiver should have tokens at the expected IBC denom: %s", expectedDenom)

	// Now verify the helper computes the same denom as what we expect
	// This uses the same transfertypes functions as our extractDenomFromPacketOnRecv
	d := transfertypes.ExtractDenomFromPath(denom)
	d.Trace = append(
		[]transfertypes.Hop{transfertypes.NewHop(
			s.TransferPath.EndpointB.ChannelConfig.PortID,
			s.TransferPath.EndpointB.ChannelID,
		)},
		d.Trace...,
	)
	computedDenom := d.IBCDenom()

	s.Require().Equal(expectedDenom, computedDenom,
		"ExtractDenomFromPath + prepend should produce same denom as GetIBCDenom helper")

	s.T().Logf("Verified: IBC transfer denom = computed denom = %s", expectedDenom)
}
