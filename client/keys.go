package client

import (
	cosmosevmclient "github.com/cosmos/evm/client"
	"github.com/spf13/cobra"
)

// KeyCommands wraps Cosmos EVM's KeyCommands to provide Ethereum-compatible key management.
// This is a simple wrapper that delegates to the Cosmos EVM implementation.
func KeyCommands(defaultNodeHome string, defaultToEthKeys bool) *cobra.Command {
	return cosmosevmclient.KeyCommands(defaultNodeHome, defaultToEthKeys)
}
