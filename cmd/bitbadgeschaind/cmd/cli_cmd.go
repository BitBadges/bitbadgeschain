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
// Examples:
//
//	bitbadgeschaind cli check tx.json
//	bitbadgeschaind cli api tokens get-collection 1
//	bitbadgeschaind cli build vault --backing-coin USDC
//	bitbadgeschaind cli deploy --burner --msg-file col.json --manager bb1...
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

The bitbadges-cli surface is flat — every verb is top-level. Common
subcommands (run ` + "`bitbadgeschaind cli --help`" + ` for the full grouped list):

  Build & ship a transaction:
    build, tools, tool, check, explain, simulate, preview, deploy

  Indexer access:
    api, auth

  Local state:
    config, burner, session

  Discovery:
    docs, skills, resources, doctor

  Address & lookup:
    address, alias, lookup, gen-list-id

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
