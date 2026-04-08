package cmd

import (
	"github.com/spf13/cobra"
)

// ApiCmd returns a cobra command that delegates to the Node.js bitbadges-cli
// for API-related operations (fetching collections, balances, etc.).
//
// Examples:
//
//	bitbadgeschaind api collection get 1
//	bitbadgeschaind api balance query bb1... 1
func ApiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "api [args...]",
		Short:              "API query commands (delegates to bitbadges-cli)",
		Long:               "Delegates to the Node.js bitbadges-cli for API commands like fetching collections, balances, and other chain data.",
		DisableFlagParsing: true, // Pass all flags through to Node.js CLI
		RunE: func(cmd *cobra.Command, args []string) error {
			return execNodeCLI("api", args)
		},
	}
	return cmd
}
