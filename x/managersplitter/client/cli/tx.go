package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bitbadges/bitbadgeschain/x/managersplitter/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdCreateManagerSplitter())
	cmd.AddCommand(CmdUpdateManagerSplitter())
	cmd.AddCommand(CmdDeleteManagerSplitter())
	cmd.AddCommand(CmdExecuteUniversalUpdateCollection())

	return cmd
}
