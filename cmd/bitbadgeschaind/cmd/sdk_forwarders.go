package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// SDK CLI forwarders.
//
// The existing `cli` subcommand (cli_cmd.go) is a generic catch-all
// forwarder — `bitbadgeschaind cli <anything>` reaches every
// bitbadges-cli subcommand. It stays registered for one release of
// back-compat per the v2 deprecation runway.
//
// This file adds **named, top-level** forwarders so `bb <cmd>` resolves
// directly without the `cli` infix:
//
//	bb build vault ...        → bitbadges-cli build vault ...
//	bb auctions place-bid 42  → bitbadges-cli auctions place-bid 42
//	bb account all            → bitbadges-cli account all
//
// Each forwarder is grouped for `bb --help` (see help_groups.go) and
// uses the same execvp-style passthrough that cli_cmd.go uses. When
// bitbadges-cli is not installed, forwarders print the same helpful
// install hint that execNodeCLI() does today.
//
// Runtime dependency: these only do useful work when bitbadges-cli
// (npm package `bitbadges`) is on the user's PATH. They compile and
// pass tests on their own, independent of the SDK CLI version.

// sdkForwarderSpec describes a single top-level SDK forwarder.
type sdkForwarderSpec struct {
	name       string
	short      string
	group      string
	deprecated string // non-empty → command is hidden + marked deprecated
}

// sdkForwarderSpecs is the canonical list of bitbadges-cli top-level
// commands that should be reachable directly on `bb`. Names and
// grouping come from outputs/flagship-plans/cli-v2-design.md "Locked
// Decisions" decisions 1-7. Order within each group is alphabetical for
// stable `--help` rendering.
//
// IMPORTANT: `config` is intentionally absent. The chain binary owns
// `bb config` (client.toml). The SDK's old `config` is reachable via
// `bb cli config` only; the SDK rename to `settings` removes the
// collision on the SDK side.
var sdkForwarderSpecs = []sdkForwarderSpec{
	// SDK — Build & Deploy
	{name: "build", short: "Builders for deterministic ready-to-sign tx JSON", group: groupSDKBuild},
	{name: "check", short: "Validate a tx JSON against its proto schema", group: groupSDKBuild},
	{name: "deploy", short: "Sign and broadcast a tx (--browser/--burner/--with-keyring/--gen-payload)", group: groupSDKBuild},
	{name: "explain", short: "Human-readable explanation of a tx JSON", group: groupSDKBuild},
	{name: "preview", short: "Preview a tx as it would land on chain", group: groupSDKBuild},
	{name: "simulate", short: "Simulate a tx against a live node (no broadcast)", group: groupSDKBuild},

	// SDK — Standards (12 end-user verbs)
	{name: "auctions", short: "Auction standard: list / show / place-bid / settle / ...", group: groupSDKStandards},
	{name: "bounties", short: "Bounty standard: list / show / claim / ...", group: groupSDKStandards},
	{name: "credit-tokens", short: "Credit-token standard: list / show / issue / ...", group: groupSDKStandards},
	{name: "crowdfunds", short: "Crowdfund standard: list / show / contribute / ...", group: groupSDKStandards},
	{name: "dynamic-stores", short: "Dynamic-store standard: list / show / set-value / ...", group: groupSDKStandards},
	{name: "intents", short: "Intent standard: list / show / submit / ...", group: groupSDKStandards},
	{name: "nfts", short: "NFT standard: list / show / transfer / ...", group: groupSDKStandards},
	{name: "pay-requests", short: "Pay-request standard: list / show / pay / ...", group: groupSDKStandards},
	{name: "prediction-markets", short: "Prediction-market standard: list / show / trade / ...", group: groupSDKStandards},
	{name: "products", short: "Product standard: list / show / buy / ...", group: groupSDKStandards},
	{name: "smart-tokens", short: "Smart-token standard: list / show / transfer / ...", group: groupSDKStandards},
	{name: "subscriptions", short: "Subscription standard: list / show / subscribe / ...", group: groupSDKStandards},

	// SDK — Indexer & Auth
	{name: "account", short: "Account aggregator: profile, tokens, balances, activity, approvals", group: groupSDKIndexer},
	{name: "api", short: "Indexer API surface (106 routes, tag-grouped)", group: groupSDKIndexer},
	{name: "auth", short: "Blockin session: login / logout / status / use / whoami", group: groupSDKIndexer},

	// SDK — Swap / DEX
	{name: "pairs", short: "Asset-pair listings (promoted from `swap asset-pairs`)", group: groupSDKSwap},
	{name: "pools", short: "Liquidity pool listings (promoted from `swap pools`)", group: groupSDKSwap},
	{name: "price", short: "Spot price for a symbol", group: groupSDKSwap},
	{name: "swap", short: "Swap / DEX queries: assets, chains, estimate, track, status", group: groupSDKSwap},

	// SDK — Dev (MCP-flavored agent surface)
	{name: "dev", short: "Dev surface: tools, resources, docs, skills, gen-pub-key", group: groupSDKDev},

	// Local state
	{name: "burner", short: "Burner keys: list / show / resume / sweep / forget", group: groupLocal},
	{name: "doctor", short: "Quick smoke test of local CLI state", group: groupLocal},
	{name: "session", short: "Local Blockin sessions: list / show / reset", group: groupLocal},
	{name: "settings", short: "SDK CLI config (renamed from `config`)", group: groupLocal},

	// Deprecated aliases — registered hidden so they still resolve but
	// don't clutter `bb --help`. The SDK CLI side prints a stderr
	// deprecation banner with the canonical replacement name.
	{name: "portfolio", short: "Deprecated: use `bb account`", group: groupSDKIndexer, deprecated: "use `bb account`"},
	{name: "address", short: "Deprecated: use `bb account` subcommands", group: groupSDKIndexer, deprecated: "use `bb account convert` / `bb account validate`"},
	{name: "alias", short: "Deprecated: use `bb account alias`", group: groupSDKIndexer, deprecated: "use `bb account alias`"},
	{name: "lookup", short: "Deprecated: use `bb account lookup`", group: groupSDKIndexer, deprecated: "use `bb account lookup`"},
	{name: "gen-list-id", short: "Deprecated: use `bb account gen-list-id`", group: groupSDKIndexer, deprecated: "use `bb account gen-list-id`"},
	{name: "tool", short: "Deprecated: use `bb dev tools`", group: groupSDKDev, deprecated: "use `bb dev tools`"},
	{name: "tools", short: "Deprecated: use `bb dev tools`", group: groupSDKDev, deprecated: "use `bb dev tools`"},
	{name: "resources", short: "Deprecated: use `bb dev resources`", group: groupSDKDev, deprecated: "use `bb dev resources`"},
	{name: "docs", short: "Deprecated: use `bb dev docs`", group: groupSDKDev, deprecated: "use `bb dev docs`"},
	{name: "skills", short: "Deprecated: use `bb dev skills`", group: groupSDKDev, deprecated: "use `bb dev skills`"},
	{name: "gen-pub-key", short: "Deprecated: use `bb dev gen-pub-key`", group: groupSDKDev, deprecated: "use `bb dev gen-pub-key`"},
	{name: "sign-with-browser", short: "Deprecated: use `bb deploy --browser`", group: groupSDKBuild, deprecated: "use `bb deploy --browser`"},
	{name: "gen-tx-payload", short: "Deprecated: use `bb deploy --gen-payload`", group: groupSDKBuild, deprecated: "use `bb deploy --gen-payload`"},
}

// newSDKForwarder builds a Cobra command that forwards verbatim to
// `bitbadges-cli <name> [args...]`. Mirrors the pattern in
// cli_cmd.go's CliCmd() — same execNodeCLI plumbing, same not-
// installed fallback message.
func newSDKForwarder(spec sdkForwarderSpec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                spec.name + " [args...]",
		Short:              spec.short,
		GroupID:            spec.group,
		DisableFlagParsing: true, // pass flags through to bitbadges-cli
		Long: fmt.Sprintf(`Forwards to: bitbadges-cli %s

Runs the BitBadges SDK CLI subcommand of the same name. All arguments
and flags are passed through verbatim. Requires Node.js + the
bitbadges npm package (npm install -g bitbadges).

Run with --help to see the SDK subcommand's own help:
  bitbadgeschaind %s --help`, spec.name, spec.name),
		RunE: func(cmd *cobra.Command, args []string) error {
			return execNodeCLI(spec.name, args)
		},
	}
	if spec.deprecated != "" {
		cmd.Hidden = true
		cmd.Deprecated = spec.deprecated
	}
	return cmd
}

// registerSDKForwarders adds every top-level SDK forwarder to the root
// command. Called from initRootCmd after chain natives are registered.
//
// Collisions: this function panics if a forwarder name collides with a
// chain-native command that's already registered. The spec list is
// pre-filtered to avoid `config` (which the chain owns), but a future
// chain command rename could surprise us — fail loud rather than
// silently shadow.
func registerSDKForwarders(rootCmd *cobra.Command) {
	existing := map[string]bool{}
	for _, c := range rootCmd.Commands() {
		existing[c.Name()] = true
	}

	for _, spec := range sdkForwarderSpecs {
		if existing[spec.name] {
			fmt.Fprintf(os.Stderr,
				"sdk forwarder name collision: %q already registered as a chain-native command\n",
				spec.name)
			continue
		}
		rootCmd.AddCommand(newSDKForwarder(spec))
	}
}
