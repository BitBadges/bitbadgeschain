package cmd

import (
	"github.com/spf13/cobra"
)

// SdkCmd returns a cobra command that delegates to the Node.js bitbadges-cli
// for SDK analysis and utility operations (review, interpret, address conversion, etc.).
//
// Examples:
//
//	bitbadgeschaind sdk review tx.json
//	bitbadgeschaind sdk address convert 0x... --to bb1
func SdkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "sdk [args...]",
		Short:              "SDK analysis and utility commands (delegates to bitbadges-cli)",
		Long:               "Delegates to the Node.js bitbadges-cli for SDK commands like review, interpret, address conversion, etc.",
		DisableFlagParsing: true, // Pass all flags through to Node.js CLI
		RunE: func(cmd *cobra.Command, args []string) error {
			return execNodeCLI("sdk", args)
		},
	}
	return cmd
}
