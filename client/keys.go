package client

import (
	cosmosevmclient "github.com/cosmos/evm/client"
	"github.com/spf13/cobra"
)

// KeyCommands wraps Cosmos EVM's KeyCommands to provide Ethereum-compatible key management.
// It strips empty-name separator commands (such as flags.LineBreak) that otherwise sort to
// the top of Cobra's alphabetical help listing as a blank row.
func KeyCommands(defaultNodeHome string, defaultToEthKeys bool) *cobra.Command {
	cmd := cosmosevmclient.KeyCommands(defaultNodeHome, defaultToEthKeys)
	for _, sub := range cmd.Commands() {
		if sub.Name() == "" {
			cmd.RemoveCommand(sub)
		}
	}
	return cmd
}
