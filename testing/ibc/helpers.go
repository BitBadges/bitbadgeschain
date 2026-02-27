//go:build test
// +build test

package ibc

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibctesting "github.com/cosmos/ibc-go/v10/testing"

	"github.com/bitbadges/bitbadgeschain/app"
)

// NewCoordinator creates a new ibctesting.Coordinator with the specified number of chains
// All chains use the BitBadges app with proper EVM-compatible genesis
func NewCoordinator(t *testing.T, numChains int) *ibctesting.Coordinator {
	t.Helper()

	// Create coordinator using our custom app initializer
	// Note: ibc-go's setupWithGenesisValSet will call our DefaultTestingAppInit
	// then creates its own bank genesis. We handle this by ensuring our
	// app's default genesis has the proper EVM configuration.
	coordinator := ibctesting.NewCoordinator(t, numChains)
	return coordinator
}

// NewTransferPath creates a new IBC transfer path between two chains
func NewTransferPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = transfertypes.PortID
	path.EndpointB.ChannelConfig.PortID = transfertypes.PortID
	path.EndpointA.ChannelConfig.Version = transfertypes.V1
	path.EndpointB.ChannelConfig.Version = transfertypes.V1
	return path
}

// GetBitBadgesApp extracts the BitBadges App from a TestChain
func GetBitBadgesApp(chain *ibctesting.TestChain) *app.App {
	testingApp, ok := chain.App.(*app.App)
	if !ok {
		panic("failed to get BitBadges App from TestChain")
	}
	return testingApp
}

// FundAccount funds an account on the given chain with the specified coins
func FundAccount(chain *ibctesting.TestChain, addr sdk.AccAddress, coins sdk.Coins) error {
	bitbadgesApp := GetBitBadgesApp(chain)
	ctx := chain.GetContext()

	// Mint coins to the account
	if err := bitbadgesApp.BankKeeper.MintCoins(ctx, transfertypes.ModuleName, coins); err != nil {
		return err
	}

	// Send from module account to target account
	if err := bitbadgesApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, transfertypes.ModuleName, addr, coins); err != nil {
		return err
	}

	return nil
}

// GetBalance returns the balance of an account on the given chain
func GetBalance(chain *ibctesting.TestChain, addr sdk.AccAddress, denom string) sdk.Coin {
	bitbadgesApp := GetBitBadgesApp(chain)
	ctx := chain.GetContext()
	return bitbadgesApp.BankKeeper.GetBalance(ctx, addr, denom)
}

// CreateTransferMsg creates a MsgTransfer for ICS20 token transfer
func CreateTransferMsg(
	sourcePort, sourceChannel string,
	token sdk.Coin,
	sender, receiver string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	memo string,
) *transfertypes.MsgTransfer {
	return transfertypes.NewMsgTransfer(
		sourcePort,
		sourceChannel,
		token,
		sender,
		receiver,
		timeoutHeight,
		timeoutTimestamp,
		memo,
	)
}

// GetIBCDenom returns the IBC denom for a token transferred through a specific path
func GetIBCDenom(portID, channelID, baseDenom string) string {
	return transfertypes.ParseDenomTrace(
		transfertypes.GetPrefixedDenom(portID, channelID, baseDenom),
	).IBCDenom()
}

// DefaultTransferAmount returns the default amount used in transfer tests
func DefaultTransferAmount() sdkmath.Int {
	return sdkmath.NewInt(1000000) // 1 unit with 6 decimals
}

// RelayPacket is a helper to relay a packet from source to destination
// This commits the packet on the receiving chain and sends the ack back
func RelayPacket(path *ibctesting.Path, packet channeltypes.Packet) error {
	// The ibctesting.Path has RelayPacket method
	return path.RelayPacket(packet)
}

// CommitBlock commits the current block and advances to the next
func CommitBlock(chain *ibctesting.TestChain) {
	chain.NextBlock()
}

// GetSenderAddress returns a funded sender address from the chain's validators
func GetSenderAddress(chain *ibctesting.TestChain) sdk.AccAddress {
	return chain.SenderAccount.GetAddress()
}

// SetupTransferPath creates and sets up a transfer path between two chains
// This creates clients, connections, and channels
func SetupTransferPath(coordinator *ibctesting.Coordinator, chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := NewTransferPath(chainA, chainB)
	coordinator.Setup(path)
	return path
}
