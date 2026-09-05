//go:build test
// +build test

package e2e

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v11/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v11/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v11/testing"

	ibctest "github.com/bitbadges/bitbadgeschain/testing/ibc"
	ibcratelimittypes "github.com/bitbadges/bitbadgeschain/x/ibc-rate-limit/types"
)

const sendWindowBlocks = int64(100)

// setSendRateLimitParams configures an outbound supply-shift limit on chain A
// for its side of the transfer channel.
func (s *RateLimitTestSuite) setSendRateLimitParams(denom string, maxAmount sdkmath.Int) {
	bitbadgesApp := s.GetBitBadgesApp(s.ChainA)
	ctx := s.ChainA.GetContext()

	params := ibcratelimittypes.Params{
		RateLimits: []ibcratelimittypes.RateLimitConfig{
			{
				ChannelId: s.TransferPath.EndpointA.ChannelID,
				Denom:     denom,
				SupplyShiftLimits: []ibcratelimittypes.TimeframeLimit{
					{
						MaxAmount:         maxAmount,
						TimeframeType:     ibcratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
						TimeframeDuration: sendWindowBlocks,
					},
				},
			},
		},
	}

	s.Require().NoError(bitbadgesApp.IBCRateLimitKeeper.SetParams(ctx, params))
	s.Coordinator.CommitBlock(s.ChainA)
}

func (s *RateLimitTestSuite) clearSendRateLimitParams() {
	bitbadgesApp := s.GetBitBadgesApp(s.ChainA)
	ctx := s.ChainA.GetContext()
	s.Require().NoError(bitbadgesApp.IBCRateLimitKeeper.SetParams(ctx, ibcratelimittypes.DefaultParams()))
	s.Coordinator.CommitBlock(s.ChainA)
}

// chainANetFlow reads chain A's tracked net flow for its transfer channel.
func (s *RateLimitTestSuite) chainANetFlow(denom string) sdkmath.Int {
	bitbadgesApp := s.GetBitBadgesApp(s.ChainA)
	flow, _ := bitbadgesApp.IBCRateLimitKeeper.GetChannelFlowWithTimeframe(
		s.ChainA.GetContext(),
		s.TransferPath.EndpointA.ChannelID,
		denom,
		ibcratelimittypes.TimeframeType_TIMEFRAME_TYPE_BLOCK,
		sendWindowBlocks,
	)
	return flow.NetFlow
}

func (s *RateLimitTestSuite) newOutboundTransfer(sender sdk.AccAddress, receiver string, token sdk.Coin) *transfertypes.MsgTransfer {
	return transfertypes.NewMsgTransfer(
		s.TransferPath.EndpointA.ChannelConfig.PortID,
		s.TransferPath.EndpointA.ChannelID,
		token,
		sender.String(),
		receiver,
		clienttypes.NewHeight(1, 1000),
		0,
		"",
	)
}

// TestSendPacketRespectsRateLimit sends real MsgTransfers through the full
// transfer stack and checks that the outbound supply-shift limit is enforced
// and that an in-limit send debits the tracked flow.
func (s *RateLimitTestSuite) TestSendPacketRespectsRateLimit() {
	sender := s.ChainA.SenderAccount.GetAddress()
	receiver := s.ChainB.SenderAccount.GetAddress()
	denom := "ubadge"

	maxAmount := sdkmath.NewInt(1_000_000)
	s.setSendRateLimitParams(denom, maxAmount)
	defer s.clearSendRateLimitParams()

	s.Require().NoError(ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, maxAmount.MulRaw(10)))))
	s.Require().True(s.chainANetFlow(denom).IsZero(), "flow must start at zero")

	// Above the limit: the transfer tx must fail and nothing may be tracked.
	over := s.newOutboundTransfer(sender, receiver.String(), sdk.NewCoin(denom, maxAmount.AddRaw(1)))
	_, err := s.ChainA.SendMsgs(over)
	s.Require().Error(err, "transfer above the outbound limit must be rejected")
	s.Require().Contains(err.Error(), "rate limit exceeded")
	s.Require().True(s.chainANetFlow(denom).IsZero(), "rejected send must not be tracked")

	// Below the limit: the transfer succeeds and the flow is debited.
	amount := sdkmath.NewInt(400_000)
	under := s.newOutboundTransfer(sender, receiver.String(), sdk.NewCoin(denom, amount))
	res, err := s.ChainA.SendMsgs(under)
	s.Require().NoError(err)
	s.Require().True(s.chainANetFlow(denom).Equal(amount.Neg()), "outbound send must debit the flow, got %s", s.chainANetFlow(denom))

	// A successful ack leaves the debit in place.
	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)
	s.Require().NoError(s.TransferPath.RelayPacket(packet))
	s.Require().True(s.chainANetFlow(denom).Equal(amount.Neg()), "success ack must not change the flow, got %s", s.chainANetFlow(denom))

	// The remaining quota is enforced cumulatively within the window.
	over2 := s.newOutboundTransfer(sender, receiver.String(), sdk.NewCoin(denom, maxAmount.Sub(amount).AddRaw(1)))
	_, err = s.ChainA.SendMsgs(over2)
	s.Require().Error(err, "cumulative outbound sends above the limit must be rejected")
	s.Require().Contains(err.Error(), "rate limit exceeded")
}

// TestSendPacketErrorAckRestoresFlow sends to a receiver the counterparty
// rejects and checks that the resulting error ack refunds exactly the amount
// that was debited at send time.
func (s *RateLimitTestSuite) TestSendPacketErrorAckRestoresFlow() {
	sender := s.ChainA.SenderAccount.GetAddress()
	denom := "ubadge"

	maxAmount := sdkmath.NewInt(1_000_000)
	s.setSendRateLimitParams(denom, maxAmount)
	defer s.clearSendRateLimitParams()

	s.Require().NoError(ibctest.FundAccount(s.ChainA, sender, sdk.NewCoins(sdk.NewCoin(denom, maxAmount.MulRaw(10)))))

	amount := sdkmath.NewInt(300_000)
	msg := s.newOutboundTransfer(sender, "not-a-valid-bech32-receiver", sdk.NewCoin(denom, amount))
	res, err := s.ChainA.SendMsgs(msg)
	s.Require().NoError(err)
	s.Require().True(s.chainANetFlow(denom).Equal(amount.Neg()), "outbound send must debit the flow, got %s", s.chainANetFlow(denom))

	packet, err := ibctesting.ParsePacketFromEvents(res.Events)
	s.Require().NoError(err)
	s.Require().NoError(s.TransferPath.RelayPacket(packet))

	s.Require().True(s.chainANetFlow(denom).IsZero(), "error ack must refund the debited amount, got %s", s.chainANetFlow(denom))
}
