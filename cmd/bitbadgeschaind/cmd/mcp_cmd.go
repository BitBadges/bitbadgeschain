package cmd

import (
	"github.com/spf13/cobra"
)

// McpCmd returns a cobra command that delegates to the Node.js bitbadges-cli
// for builder MCP tool invocation (list/call session builder tools directly,
// no MCP protocol / subprocess required).
//
// Examples:
//
//	bitbadgeschaind mcp list
//	bitbadgeschaind mcp call get_skill_instructions --args '{"skillId":"smart-token"}'
//	bitbadgeschaind mcp call set_collection_metadata --session demo --args '{"name":"Foo",...}'
//	bitbadgeschaind mcp session show demo
func McpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp [args...]",
		Short: "BitBadges Builder MCP tools — list/call session builders (via bitbadges-cli)",
		Long: `Invoke bitbadges-builder-mcp tools directly from the chain CLI.
Tool handlers are called as library functions — no MCP protocol, no subprocess.

Supports 50+ tools spanning session builders (add_approval, set_collection_metadata,
get_transaction, ...), query tools (query_collection, simulate_transaction, ...),
audits (audit_collection, verify_standards), and utilities (convert_address,
diagnose_error, get_skill_instructions, ...).

Session state persists per-id under ~/.bitbadges/sessions/<id>.json so you can
compose a collection across multiple invocations.

Requires: Node.js + bitbadges-cli (npm install -g bitbadgesjs-sdk)`,
		DisableFlagParsing: true, // Pass all flags through to Node.js CLI
		RunE: func(cmd *cobra.Command, args []string) error {
			return execNodeCLI("mcp", args)
		},
	}
	return cmd
}
