package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// CliCmd returns the generic catch-all forwarder to the Node.js
// bitbadges-cli binary. Every top-level JS CLI subcommand is reachable
// through this one command — `bitbadgeschaind cli <anything>` — without
// needing a dedicated Go file per alias.
//
// This is the canonical form going forward. The existing short aliases
// (sdk / api / builder) stay around for one release cycle for backwards
// compat but redirect through this same forwarder.
//
// Examples:
//
//	bitbadgeschaind cli sdk review tx.json
//	bitbadgeschaind cli api tokens get-collection 1
//	bitbadgeschaind cli builder templates vault --backing-coin USDC
//	bitbadgeschaind cli builder create-with-burner --msg-file col.json --manager bb1...
//	bitbadgeschaind cli config set apiKey <key>
//
// Any future top-level subcommand added to bitbadges-cli is automatically
// reachable via `bitbadgeschaind cli <subcmd>` — no Go changes required.
func CliCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cli [subcommand] [args...]",
		Short: "Forward any subcommand to the Node.js bitbadges-cli",
		Long: `Generic forwarder for the Node.js bitbadges-cli binary.

Usage:
  bitbadgeschaind cli <subcommand> [args...]

Where <subcommand> is any top-level bitbadges-cli subcommand:
  sdk       — SDK analysis, review, interpret, address tools, docs
  api       — 104+ indexer API routes from your terminal
  builder   — template builders, create-with-burner, burner wallets,
              review/verify/simulate/doctor, session tools, builder tools
  config    — manage ~/.bitbadges/config.json (base URL, API keys per network)

New top-level bitbadges-cli subcommands automatically reach here without
needing a Go alias — this wrapper is the single bridge between the chain
binary and the JS CLI going forward.

Requires: Node.js + bitbadges-cli (npm install -g bitbadges)`,
		DisableFlagParsing: true, // Pass all flags through to Node.js CLI
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "Usage: bitbadgeschaind cli <subcommand> [args...]")
				fmt.Fprintln(os.Stderr, "Run `bitbadgeschaind cli --help` for the list of subcommands.")
				return nil
			}
			return execNodeCLI(args[0], args[1:])
		},
	}
	return cmd
}
