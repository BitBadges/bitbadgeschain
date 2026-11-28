package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

// annotationSDKForwarder marks Cobra commands that delegate verbatim to the
// Node.js bitbadges-cli binary. rootCmd.PersistentPreRunE checks this
// annotation to bypass Cosmos client.toml / server config loading so forwarded
// commands work cleanly with an empty, missing, or read-only --home directory.
const annotationSDKForwarder = "bitbadges.sdk_forwarder"

// isSDKForwarderCmd reports whether cmd is an SDK CLI passthrough forwarder.
func isSDKForwarderCmd(cmd *cobra.Command) bool {
	if cmd == nil || cmd.Annotations == nil {
		return false
	}
	return cmd.Annotations[annotationSDKForwarder] == "true"
}

// cliArgsProvider returns the raw CLI arguments (excluding argv[0]) that Cobra
// dispatched from. Overridden only in unit tests if needed; production uses os.Args[1:].
var cliArgsProvider = func() []string {
	if len(os.Args) <= 1 {
		return nil
	}
	return os.Args[1:]
}

// lookupCmdFlag looks up a flag (or shorthand) defined on cmd or inherited from
// any ancestor's PersistentFlags.
func lookupCmdFlag(c *cobra.Command, name string, shorthand bool) *pflag.Flag {
	for cur := c; cur != nil; cur = cur.Parent() {
		if cur == c {
			if shorthand {
				if f := cur.Flags().ShorthandLookup(name); f != nil {
					return f
				}
			} else {
				if f := cur.Flags().Lookup(name); f != nil {
					return f
				}
			}
		}
		if shorthand {
			if f := cur.PersistentFlags().ShorthandLookup(name); f != nil {
				return f
			}
		} else {
			if f := cur.PersistentFlags().Lookup(name); f != nil {
				return f
			}
		}
	}
	return nil
}

// extractForwarderArgs locates cmd along its command path in rawArgs, confirms
// that the leading root/parent flags plus trailing subcommand tokens reconstruct
// cobraArgs verbatim, and returns ONLY the tokens appearing after cmd.
//
// Because SDK forwarders set DisableFlagParsing: true, Cobra leaves leading
// persistent root flags (such as --home, --log_level, --log_format, --trace)
// in args alongside the subcommand's own flags. Inspecting rawArgs by position
// relative to the matched command token strips leading Cosmos root flags
// (`bb --log_level info --home /tmp/empty dev --help` -> `["--help"]`) while
// preserving any flags that appear AFTER the subcommand (`bb dev --home /foo`
// -> `["--home", "/foo"]`, `bb cli --home /foo` -> `["--home", "/foo"]`).
func extractForwarderArgs(cmd *cobra.Command, rawArgs, cobraArgs []string) ([]string, bool) {
	if cmd == nil {
		return nil, false
	}

	var path []*cobra.Command
	for cur := cmd; cur != nil && cur.HasParent(); cur = cur.Parent() {
		path = append(path, cur)
	}
	if len(path) == 0 {
		return nil, false
	}
	slices.Reverse(path)

	pos := 0
	prefixTokens := make([]string, 0)
	parent := cmd.Root()

	for _, target := range path {
		matched := false
		for pos < len(rawArgs) {
			tok := rawArgs[pos]
			switch {
			case tok == "--":
				return nil, false
			case strings.HasPrefix(tok, "--") && len(tok) > 2:
				prefixTokens = append(prefixTokens, tok)
				if !strings.Contains(tok, "=") {
					f := lookupCmdFlag(parent, tok[2:], false)
					if f == nil || f.NoOptDefVal == "" {
						if pos+1 >= len(rawArgs) {
							return nil, false
						}
						pos++
						prefixTokens = append(prefixTokens, rawArgs[pos])
					}
				}
				pos++
			case strings.HasPrefix(tok, "-") && len(tok) == 2 && tok != "--":
				prefixTokens = append(prefixTokens, tok)
				f := lookupCmdFlag(parent, tok[1:], true)
				if f == nil || f.NoOptDefVal == "" {
					if pos+1 >= len(rawArgs) {
						return nil, false
					}
					pos++
					prefixTokens = append(prefixTokens, rawArgs[pos])
				}
				pos++
			case strings.HasPrefix(tok, "-") && len(tok) > 2:
				prefixTokens = append(prefixTokens, tok)
				pos++
			default:
				if tok == target.CalledAs() || tok == target.Name() || target.HasAlias(tok) {
					pos++
					matched = true
				} else {
					return nil, false
				}
			}
			if matched {
				break
			}
		}
		if !matched {
			return nil, false
		}
		parent = target
	}

	suffix := rawArgs[pos:]
	reconstructed := make([]string, 0, len(prefixTokens)+len(suffix))
	reconstructed = append(reconstructed, prefixTokens...)
	reconstructed = append(reconstructed, suffix...)
	if !slices.Equal(reconstructed, cobraArgs) {
		return nil, false
	}

	out := make([]string, len(suffix))
	copy(out, suffix)
	return out, true
}

// forwarderArgs returns the exact arguments to pass through to bitbadges-cli
// for an SDK forwarder command, stripping only leading Cosmos root/parent flags
// that preceded the matched command token in the process argv.
func forwarderArgs(cmd *cobra.Command, args []string) []string {
	if raw := cliArgsProvider(); len(raw) > 0 {
		if subArgs, ok := extractForwarderArgs(cmd, raw, args); ok {
			return subArgs
		}
	}
	return args
}

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

	// SDK — Standards (13 end-user verbs)
	{name: "auctions", short: "Auction standard: list / show / place-bid / accept-bid / ...", group: groupSDKStandards},
	{name: "bounties", short: "Bounty standard: list / show / accept / deny / ...", group: groupSDKStandards},
	{name: "credit-tokens", short: "Credit-token standard: list / show / quote / purchase / ...", group: groupSDKStandards},
	{name: "crowdfunds", short: "Crowdfund standard: list / show / contribute / ...", group: groupSDKStandards},
	{name: "custom-2fa", short: "Custom 2FA standard: mint short-lived 2FA tokens", group: groupSDKStandards},
	{name: "dynamic-stores", short: "Dynamic-store standard: show / get-value / set-value / ...", group: groupSDKStandards},
	{name: "intents", short: "Intent standard: list / show / create / fill / ...", group: groupSDKStandards},
	{name: "nfts", short: "NFT orderbook: bid / list / buy / sell / orders / ...", group: groupSDKStandards},
	{name: "pay-requests", short: "Pay-request standard: list / show / pay / ...", group: groupSDKStandards},
	{name: "prediction-markets", short: "Prediction-market standard: list / show / buy-yes / buy-no / ...", group: groupSDKStandards},
	{name: "products", short: "Product standard: list / show / purchase / ...", group: groupSDKStandards},
	{name: "smart-tokens", short: "Smart-token standard: list / show / deposit / withdraw / ...", group: groupSDKStandards},
	{name: "subscriptions", short: "Subscription standard: list / status / subscribe / ...", group: groupSDKStandards},

	// SDK — Indexer & Auth
	{name: "account", short: "Account aggregator: profile, tokens, balances, activity, approvals", group: groupSDKIndexer},
	{name: "amount", short: "Amount math: to-raw / to-display / slippage / wrap-preview / ...", group: groupSDKIndexer},
	{name: "api", short: "Indexer API surface (106 routes, tag-grouped)", group: groupSDKIndexer},
	{name: "auth", short: "Blockin session: login / logout / status / use / whoami", group: groupSDKIndexer},
	{name: "url", short: "bitbadges.io links for an address, collection, token, or tx", group: groupSDKIndexer},

	// SDK — Swap / DEX
	{name: "assets", short: "Tradable asset listings: list / show / browse / price", group: groupSDKSwap},
	{name: "balances", short: "Balances across ics20 / bitbadges / assets", group: groupSDKSwap},
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
		Annotations: map[string]string{
			annotationSDKForwarder: "true",
		},
		Long: fmt.Sprintf(`Forwards to: bitbadges-cli %s

Runs the BitBadges SDK CLI subcommand of the same name. All arguments
and flags are passed through verbatim. Requires Node.js + the
bitbadges npm package (npm install -g bitbadges).

Run with --help to see the SDK subcommand's own help:
  bitbadgeschaind %s --help`, spec.name, spec.name),
		RunE: func(cmd *cobra.Command, args []string) error {
			return execNodeCLI(spec.name, forwarderArgs(cmd, args))
		},
	}
	if spec.deprecated != "" {
		cmd.Hidden = true
		cmd.Deprecated = spec.deprecated
	}
	return cmd
}

// sdkTxSubforwarders are SDK verbs that live under the chain-owned `tx`
// noun. `bb tx` is Cosmos SDK's transaction command, so the SDK's
// `tx status` and `tx wait` cannot be registered as top-level
// forwarders; instead they are attached as subcommands of the chain's
// `tx` and forward to `bitbadges-cli tx <verb> [args...]`.
var sdkTxSubforwarders = []sdkForwarderSpec{
	{name: "status", short: "Look up a broadcast tx by hash on the indexer (SDK)"},
	{name: "wait", short: "Wait for a broadcast tx to be indexed (SDK)"},
}

// registerSDKTxSubforwarders adds the SDK's `tx status` / `tx wait` to
// the chain's `tx` command. Skips a name the chain already owns.
func registerSDKTxSubforwarders(txCmd *cobra.Command) {
	existing := map[string]bool{}
	for _, c := range txCmd.Commands() {
		existing[c.Name()] = true
	}
	for _, spec := range sdkTxSubforwarders {
		if existing[spec.name] {
			fmt.Fprintf(os.Stderr,
				"sdk tx forwarder name collision: %q already registered under tx\n", spec.name)
			continue
		}
		spec := spec
		txCmd.AddCommand(&cobra.Command{
			Use:                spec.name + " [args...]",
			Short:              spec.short,
			DisableFlagParsing: true,
			Annotations: map[string]string{
				annotationSDKForwarder: "true",
			},
			Long: fmt.Sprintf(`Forwards to: bitbadges-cli tx %s

All arguments and flags are passed through verbatim. Requires the
bitbadges npm package (npm install -g bitbadges).`, spec.name),
			RunE: func(cmd *cobra.Command, args []string) error {
				return execNodeCLI("tx", append([]string{spec.name}, forwarderArgs(cmd, args)...))
			},
		})
	}
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
