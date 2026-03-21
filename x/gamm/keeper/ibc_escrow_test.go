package keeper_test

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"

	"github.com/bitbadges/bitbadgeschain/third_party/osmomath"
	"github.com/bitbadges/bitbadgeschain/x/gamm/keeper"
	"github.com/bitbadges/bitbadgeschain/x/gamm/types"
	poolmanagertypes "github.com/bitbadges/bitbadgeschain/x/poolmanager/types"
)

// mockICS4Wrapper records SendPacket calls without actually sending IBC packets.
type mockICS4Wrapper struct {
	sendPacketCalled bool
}

func (m *mockICS4Wrapper) SendPacket(
	ctx sdk.Context,
	sourcePort string,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	data []byte,
) (uint64, error) {
	m.sendPacketCalled = true
	return 1, nil
}

// mockChannelKeeper returns a valid channel for any query.
type mockChannelKeeper struct{}

func (m *mockChannelKeeper) GetChannel(ctx sdk.Context, portID, channelID string) (channeltypes.Channel, bool) {
	return channeltypes.Channel{
		State:    channeltypes.OPEN,
		Ordering: channeltypes.UNORDERED,
		Counterparty: channeltypes.Counterparty{
			PortId:    "transfer",
			ChannelId: "channel-0",
		},
		ConnectionHops: []string{"connection-0"},
		Version:        "ics20-1",
	}, true
}

// mockTransferKeeper simulates the IBC transfer module's Transfer method.
// It escrows tokens by sending them from the sender to the "transfer" module
// account, mirroring what the real ibc-go transfer keeper does.
type mockTransferKeeper struct {
	transferCalled bool
	bankKeeper     types.BankKeeper
	lastMsg        *transfertypes.MsgTransfer
}

func (m *mockTransferKeeper) DenomPathFromHash(ctx sdk.Context, denom string) (string, error) {
	return denom, nil
}

func (m *mockTransferKeeper) Transfer(ctx context.Context, msg *transfertypes.MsgTransfer) (*transfertypes.MsgTransferResponse, error) {
	m.transferCalled = true
	m.lastMsg = msg

	// Simulate escrow: move tokens from sender to the "transfer" module account,
	// exactly as the real IBC transfer keeper does for native tokens.
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	if err := m.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, transfertypes.ModuleName, sdk.NewCoins(msg.Token)); err != nil {
		return nil, err
	}

	return &transfertypes.MsgTransferResponse{Sequence: 1}, nil
}

// TestExecuteIBCTransfer_EscrowsTokens verifies that ExecuteIBCTransfer
// delegates to the transfer keeper's Transfer method, which escrows tokens
// before sending the IBC packet.
func (s *KeeperTestSuite) TestExecuteIBCTransfer_EscrowsTokens() {
	s.SetupTest()

	sender := s.TestAccs[0]
	transferAmount := sdk.NewCoin("stake", osmomath.NewInt(500_000))

	// Fund sender
	s.FundAcc(sender, sdk.NewCoins(transferAmount))

	// Record balance before IBC transfer
	balanceBefore := s.App.BankKeeper.GetAllBalances(s.Ctx, sender)
	s.Require().True(balanceBefore.AmountOf("stake").Equal(transferAmount.Amount),
		"sender should have the funded amount before transfer")

	mockTransfer := &mockTransferKeeper{bankKeeper: s.App.BankKeeper}

	testKeeper := keeper.NewKeeper(
		s.App.AppCodec(),
		s.App.GetKey(types.StoreKey),
		s.App.GetSubspace(types.ModuleName),
		s.App.AccountKeeper,
		s.App.BankKeeper,
		s.App.DistrKeeper,
		s.App.TokenizationKeeper,
		&s.App.SendmanagerKeeper,
		mockTransfer,
		&mockICS4Wrapper{},
		&mockChannelKeeper{},
	)

	ibcInfo := &types.IBCTransferInfo{
		SourceChannel: "channel-0",
		Receiver:      "cosmos1receiver",
	}

	err := testKeeper.ExecuteIBCTransfer(s.Ctx, sender, ibcInfo, transferAmount)
	s.Require().NoError(err, "ExecuteIBCTransfer should succeed")

	// Verify Transfer was called (delegation to transfer module)
	s.Require().True(mockTransfer.transferCalled, "transferKeeper.Transfer should have been called")
	s.Require().Equal("channel-0", mockTransfer.lastMsg.SourceChannel)
	s.Require().Equal("cosmos1receiver", mockTransfer.lastMsg.Receiver)
	s.Require().Equal(transferAmount, mockTransfer.lastMsg.Token)

	// Verify sender balance decreased (tokens were escrowed)
	balanceAfter := s.App.BankKeeper.GetAllBalances(s.Ctx, sender)
	s.Require().True(balanceAfter.AmountOf("stake").IsZero(),
		"sender balance should be zero after IBC transfer — tokens must be escrowed")
}

// TestSwapExactAmountInWithIBCTransfer_EscrowsTokens is an end-to-end test
// that performs a real swap followed by an IBC transfer, verifying that the
// swapped tokens are escrowed via the transfer keeper before being sent.
func (s *KeeperTestSuite) TestSwapExactAmountInWithIBCTransfer_EscrowsTokens() {
	s.SetupTest()
	ctx := s.Ctx

	// Zero taker fee so swap math is simpler
	poolManagerParams := s.App.PoolManagerKeeper.GetParams(ctx)
	poolManagerParams.TakerFeeParams.DefaultTakerFee = osmomath.ZeroDec()
	s.App.PoolManagerKeeper.SetParams(ctx, poolManagerParams)

	// Create a pool with foo/bar
	s.PrepareBalancerPool() // pool 1: foo, bar, baz, stake

	sender := s.TestAccs[0]

	// Record bar balance before the swap+IBC
	barBefore := s.App.BankKeeper.GetAllBalances(ctx, sender).AmountOf("bar")

	mockTransfer := &mockTransferKeeper{bankKeeper: s.App.BankKeeper}

	testKeeper := keeper.NewKeeper(
		s.App.AppCodec(),
		s.App.GetKey(types.StoreKey),
		s.App.GetSubspace(types.ModuleName),
		s.App.AccountKeeper,
		s.App.BankKeeper,
		s.App.DistrKeeper,
		s.App.TokenizationKeeper,
		&s.App.SendmanagerKeeper,
		mockTransfer,
		&mockICS4Wrapper{},
		&mockChannelKeeper{},
	)
	testKeeper.SetPoolManager(&s.App.PoolManagerKeeper)

	msgServer := keeper.NewMsgServerImpl(&testKeeper)

	resp, err := msgServer.SwapExactAmountInWithIBCTransfer(ctx,
		&types.MsgSwapExactAmountInWithIBCTransfer{
			Sender: sender.String(),
			Routes: []poolmanagertypes.SwapAmountInRoute{
				{PoolId: 1, TokenOutDenom: "bar"},
			},
			TokenIn:           sdk.NewCoin("foo", osmomath.NewInt(10000)),
			TokenOutMinAmount: osmomath.NewInt(1),
			IbcTransferInfo: types.IBCTransferInfo{
				SourceChannel: "channel-0",
				Receiver:      "cosmos1receiver",
			},
		})
	s.Require().NoError(err, "SwapExactAmountInWithIBCTransfer should succeed")
	s.Require().True(mockTransfer.transferCalled, "transferKeeper.Transfer should have been called")

	tokenOutAmount := resp.TokenOutAmount
	s.Require().True(tokenOutAmount.GT(osmomath.ZeroInt()), "swap should produce tokens")

	// After swap+IBC: sender got tokenOut from the swap, and that amount
	// should have been escrowed for the IBC transfer. So the net change
	// in bar balance should be zero (swap gives bar, IBC escrows bar).
	barAfter := s.App.BankKeeper.GetAllBalances(ctx, sender).AmountOf("bar")
	barDelta := barAfter.Sub(barBefore)

	s.Require().True(barDelta.IsZero(),
		"swap output should be fully escrowed — net bar delta should be zero, got %s", barDelta)
}
