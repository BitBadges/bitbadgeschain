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
		Short:              "BitBadges JS SDK — review, interpret, address tools, docs (via bitbadges-cli)",
		Long: `BitBadges JavaScript SDK utilities. Delegates to the bitbadges-cli Node.js tool.

Includes: collection/transaction review, interpret, address conversion,
alias generation, token lookup, builder skill docs, and more.

Requires: Node.js + bitbadges-cli (npm install -g bitbadges)`,
		DisableFlagParsing: true, // Pass all flags through to Node.js CLI
		RunE: func(cmd *cobra.Command, args []string) error {
			return execNodeCLI("sdk", args)
		},
	}
	return cmd
}
