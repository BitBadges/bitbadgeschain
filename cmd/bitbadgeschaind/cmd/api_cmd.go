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
		Short:              "BitBadges off-chain indexer API — collections, balances, claims, DEX (via bitbadges-cli)",
		Long: `BitBadges off-chain indexer API client. Delegates to the bitbadges-cli Node.js tool.

Queries the BitBadges indexer (not the on-chain node). Includes: collections,
balances, claims, plugins, DEX pools, asset pairs, dynamic stores, and 100+ routes.

Requires: Node.js + bitbadges-cli (npm install -g bitbadgesjs-sdk)
Configure: BITBADGES_API_KEY env var or bitbadges-cli config set apiKey <key>`,
		DisableFlagParsing: true, // Pass all flags through to Node.js CLI
		RunE: func(cmd *cobra.Command, args []string) error {
			return execNodeCLI("api", args)
		},
	}
	return cmd
}
