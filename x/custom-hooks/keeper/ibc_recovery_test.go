package keeper_test

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/lfdt-ibc/ibc-go/v11/modules/apps/transfer/types"
	channeltypes "github.com/lfdt-ibc/ibc-go/v11/modules/core/04-channel/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/custom-hooks/keeper"
	customhookstypes "github.com/bitbadges/bitbadgeschain/x/custom-hooks/types"
	ibchookstypes "github.com/bitbadges/bitbadgeschain/x/ibc-hooks/types"
	tokenizationkeeper "github.com/bitbadges/bitbadgeschain/x/tokenization/keeper"
)

type recordingTransferKeeper struct {
	lastMsg  *transfertypes.MsgTransfer
	nextSeq  uint64
	sendBank func(ctx context.Context, from sdk.AccAddress, coin sdk.Coin) error
}

type openChannelKeeper struct{}

func (openChannelKeeper) GetChannel(_ sdk.Context, _, _ string) (channeltypes.Channel, bool) {
	return channeltypes.Channel{State: channeltypes.OPEN}, true
}

func (r *recordingTransferKeeper) DenomPathFromHash(_ sdk.Context, denom string) (string, error) {
	return denom, nil
}

func (r *recordingTransferKeeper) Transfer(ctx context.Context, msg *transfertypes.MsgTransfer) (*transfertypes.MsgTransferResponse, error) {
	cp := *msg
	r.lastMsg = &cp
	if r.sendBank != nil {
		fromAddr, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			return nil, err
		}
		if err := r.sendBank(ctx, fromAddr, msg.Token); err != nil {
			return nil, err
		}
	}
	r.nextSeq++
	return &transfertypes.MsgTransferResponse{Sequence: r.nextSeq}, nil
}

// TestSwapAndActionOutboundIBCUsesDerivedSenderAndRecoversOnFailure verifies that:
// 1. Outbound IBC transfers use the derived intermediate sender rather than recover_address.
// 2. Pending recovery metadata is persisted under (port, channel, sequence).
// 3. Error acknowledgements and timeouts refund tokens from the intermediate sender to recover_address.
// 4. Successful acknowledgements clear the recovery record without moving funds.
func (s *KeeperTestSuite) TestSwapAndActionOutboundIBCUsesDerivedSenderAndRecoversOnFailure() {
	poolID := s.prepareTestPool()

	derivedSender := s.TestAccs[0]
	recoverAddr := s.TestAccs[1]
	escrowSink := s.TestAccs[2]

	// Fund derivedSender with input tokens for swap and give recoverAddr a distinct initial balance.
	s.FundAcc(derivedSender, sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(1_000_000))})
	s.FundAcc(recoverAddr, sdk.Coins{sdk.NewCoin("uatom", osmomath.NewInt(500_000))})

	recorder := &recordingTransferKeeper{
		sendBank: func(ctx context.Context, from sdk.AccAddress, coin sdk.Coin) error {
			return s.App.BankKeeper.SendCoins(ctx, from, escrowSink, sdk.NewCoins(coin))
		},
	}

	k := keeper.NewKeeper(
		s.App.Logger(),
		&s.App.GammKeeper,
		s.App.BankKeeper,
		s.App.TokenizationKeeper,
		&s.App.SendmanagerKeeper,
		recorder,
		s.App.HooksICS4Wrapper,
		openChannelKeeper{},
		tokenizationkeeper.NewMsgServerImpl(s.App.TokenizationKeeper),
	).WithStoreKey(s.App.GetKey(ibchookstypes.StoreKey))

	timeoutTs := uint64(s.Ctx.BlockTime().UnixNano()) + uint64(10*60*1e9)
	swapAndAction := &customhookstypes.SwapAndAction{
		UserSwap: &customhookstypes.UserSwap{
			SwapExactAssetIn: &customhookstypes.SwapExactAssetIn{
				SwapVenueName: "bitbadges-poolmanager",
				Operations: []customhookstypes.Operation{{
					Pool:     strconv.FormatUint(poolID, 10),
					DenomIn:  sdk.DefaultBondDenom,
					DenomOut: "uatom",
				}},
			},
		},
		MinAsset: &customhookstypes.MinAsset{
			Native: &customhookstypes.NativeAsset{
				Denom:  "uatom",
				Amount: "1000",
			},
		},
		TimeoutTimestamp: &timeoutTs,
		PostSwapAction: &customhookstypes.PostSwapAction{
			IBCTransfer: &customhookstypes.IBCTransferInfo{
				IBCInfo: &customhookstypes.IBCInfo{
					SourceChannel:  "channel-0",
					Receiver:       "osmo1receiver000000000000000000000000000000",
					RecoverAddress: recoverAddr.String(),
				},
			},
		},
	}

	recoverBalBefore := s.App.BankKeeper.GetBalance(s.Ctx, recoverAddr, "uatom").Amount
	ack := k.ExecuteSwapAndAction(s.Ctx, derivedSender, swapAndAction, sdk.NewCoin(sdk.DefaultBondDenom, osmomath.NewInt(100_000)), "osmo1original")
	s.Require().True(ack.Success(), "expected swap+ibc to succeed, got error: %s", getAckError(ack))
	s.Require().NotNil(recorder.lastMsg)

	// Outbound IBC transfer must debit derivedSender, not recoverAddr.
	s.Require().Equal(derivedSender.String(), recorder.lastMsg.Sender)
	s.Require().Equal(recoverBalBefore, s.App.BankKeeper.GetBalance(s.Ctx, recoverAddr, "uatom").Amount,
		"recover_address balance must remain untouched during outbound send")

	// Pending recovery record must be persisted.
	rec, found := k.GetPendingIBCRecovery(s.Ctx, transfertypes.PortID, "channel-0", 1)
	s.Require().True(found)
	s.Require().Equal(derivedSender.String(), rec.Sender)
	s.Require().Equal(recoverAddr.String(), rec.RecoverAddress)
	s.Require().Equal("uatom", rec.Denom)
	sentAmt := recorder.lastMsg.Token.Amount

	// Simulate ibc-transfer refunding derivedSender on error ack, then invoke HandleIBCAcknowledgement.
	s.Require().NoError(s.App.BankKeeper.SendCoins(s.Ctx, escrowSink, derivedSender, sdk.NewCoins(sdk.NewCoin("uatom", sentAmt))))
	pkt := channeltypes.Packet{
		SourcePort:    transfertypes.PortID,
		SourceChannel: "channel-0",
		Sequence:      1,
	}
	errAckBz := transfertypes.ModuleCdc.MustMarshalJSON(&channeltypes.Acknowledgement{
		Response: &channeltypes.Acknowledgement_Error{Error: "remote failure"},
	})
	k.HandleIBCAcknowledgement(s.Ctx, pkt, errAckBz)

	// Tokens must now have been forwarded to recoverAddr and record cleared.
	s.Require().Equal(recoverBalBefore.Add(sentAmt), s.App.BankKeeper.GetBalance(s.Ctx, recoverAddr, "uatom").Amount)
	_, found = k.GetPendingIBCRecovery(s.Ctx, transfertypes.PortID, "channel-0", 1)
	s.Require().False(found)

	// Timeout path: store a fresh pending recovery record and simulate timeout refund to derivedSender.
	s.Require().NoError(k.SetPendingIBCRecovery(s.Ctx, transfertypes.PortID, "channel-0", 2, keeper.PendingIBCRecovery{
		Sender:         derivedSender.String(),
		RecoverAddress: recoverAddr.String(),
		Denom:          "uatom",
		Amount:         "5000",
	}))
	s.FundAcc(derivedSender, sdk.Coins{sdk.NewCoin("uatom", osmomath.NewInt(5000))})
	recoverBalBeforeTimeout := s.App.BankKeeper.GetBalance(s.Ctx, recoverAddr, "uatom").Amount
	k.HandleIBCTimeout(s.Ctx, channeltypes.Packet{SourcePort: transfertypes.PortID, SourceChannel: "channel-0", Sequence: 2})
	s.Require().Equal(recoverBalBeforeTimeout.Add(osmomath.NewInt(5000)), s.App.BankKeeper.GetBalance(s.Ctx, recoverAddr, "uatom").Amount)
	_, found = k.GetPendingIBCRecovery(s.Ctx, transfertypes.PortID, "channel-0", 2)
	s.Require().False(found)

	// Success ack path: clears recovery record without forwarding funds.
	s.Require().NoError(k.SetPendingIBCRecovery(s.Ctx, transfertypes.PortID, "channel-0", 3, keeper.PendingIBCRecovery{
		Sender:         derivedSender.String(),
		RecoverAddress: recoverAddr.String(),
		Denom:          "uatom",
		Amount:         "2500",
	}))
	s.FundAcc(derivedSender, sdk.Coins{sdk.NewCoin("uatom", osmomath.NewInt(2500))})
	recoverBalBeforeSuccess := s.App.BankKeeper.GetBalance(s.Ctx, recoverAddr, "uatom").Amount
	okAckBz := transfertypes.ModuleCdc.MustMarshalJSON(&channeltypes.Acknowledgement{
		Response: &channeltypes.Acknowledgement_Result{Result: []byte{1}},
	})
	k.HandleIBCAcknowledgement(s.Ctx, channeltypes.Packet{SourcePort: transfertypes.PortID, SourceChannel: "channel-0", Sequence: 3}, okAckBz)
	s.Require().Equal(recoverBalBeforeSuccess, s.App.BankKeeper.GetBalance(s.Ctx, recoverAddr, "uatom").Amount)
	_, found = k.GetPendingIBCRecovery(s.Ctx, transfertypes.PortID, "channel-0", 3)
	s.Require().False(found)
}
