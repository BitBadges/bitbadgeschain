package cmd

import (
	"github.com/spf13/cobra"
)

// BuilderCmd returns a cobra command that delegates to the Node.js
// bitbadges-cli for BitBadges Builder tool invocation (list/call session
// builder tools directly, handlers called in-process by the JS CLI).
//
// Examples:
//
//	bitbadgeschaind builder list
//	bitbadgeschaind builder call get_skill_instructions --args '{"skillId":"smart-token"}'
//	bitbadgeschaind builder call set_collection_metadata --session demo --args '{"name":"Foo",...}'
//	bitbadgeschaind builder session show demo
func BuilderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "builder [args...]",
		Short: "BitBadges Builder tools — list/call session builders (via bitbadges-cli)",
		Long: `Invoke BitBadges Builder tools directly from the chain CLI.
Tool handlers are called in-process by the JS CLI — no subprocess, no MCP
protocol round-trip.

Supports 50+ tools spanning session builders (add_approval, set_collection_metadata,
get_transaction, ...), query tools (query_collection, simulate_transaction, ...),
reviews (review_collection), and utilities
(convert_address, diagnose_error, get_skill_instructions, ...).

Session state persists per-id under ~/.bitbadges/sessions/<id>.json so you can
compose a collection across multiple invocations.

Requires: Node.js + bitbadges-cli (npm install -g bitbadgesjs-sdk)`,
		DisableFlagParsing: true, // Pass all flags through to Node.js CLI
		RunE: func(cmd *cobra.Command, args []string) error {
			return execNodeCLI("builder", args)
		},
	}
	return cmd
}
