package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// sdkTopLevelCommands mirrors HELP_GROUPS in bitbadgesjs
// packages/bitbadgesjs-sdk/src/cli/index.ts. Every non-deprecated verb the
// SDK CLI exposes at its top level must be reachable as `bb <verb>`.
// `tx` and `completion` are owned by the chain binary and are covered by
// their own tests below.
var sdkTopLevelCommands = []string{
	// Build & ship a transaction
	"build", "check", "explain", "simulate", "preview", "deploy",
	// Indexer access
	"api", "auth",
	// Local state
	"settings", "burner", "session",
	// Discovery
	"doctor",
	// Dev / agent surface
	"dev",
	// Account & lookup
	"account", "amount", "url",
	// Swap & DEX
	"swap", "pools", "pairs", "balances", "price", "assets",
	// Standards (end-user actions)
	"pay-requests", "bounties", "subscriptions", "intents", "credit-tokens",
	"products", "crowdfunds", "auctions", "prediction-markets", "smart-tokens",
	"nfts", "custom-2fa", "dynamic-stores",
}

func newRootWithChainNatives() *cobra.Command {
	root := &cobra.Command{Use: "bb"}
	registerHelpGroups(root)
	for name := range chainNativeGroups {
		root.AddCommand(&cobra.Command{Use: name})
	}
	return root
}

func TestEverySDKTopLevelCommandHasAVisibleForwarder(t *testing.T) {
	root := newRootWithChainNatives()
	registerSDKForwarders(root)

	for _, name := range sdkTopLevelCommands {
		found, _, err := root.Find([]string{name})
		if err != nil || found == nil || found.Name() != name {
			t.Errorf("bb %s does not resolve to a forwarder", name)
			continue
		}
		if found.Hidden || found.Deprecated != "" {
			t.Errorf("bb %s is hidden or deprecated but the SDK still ships it as a live verb", name)
		}
		if !found.DisableFlagParsing {
			t.Errorf("bb %s must pass flags through verbatim", name)
		}
		if found.GroupID == "" {
			t.Errorf("bb %s has no help group", name)
		}
	}
}

func TestForwardersNeverShadowChainNatives(t *testing.T) {
	root := newRootWithChainNatives()
	before := len(root.Commands())
	registerSDKForwarders(root)

	seen := map[string]int{}
	for _, c := range root.Commands() {
		seen[c.Name()]++
	}
	for name, n := range seen {
		if n > 1 {
			t.Errorf("%q registered %d times", name, n)
		}
	}
	if len(root.Commands()) <= before {
		t.Fatal("no forwarders were registered")
	}
}

func TestChainTxCommandForwardsSDKStatusAndWait(t *testing.T) {
	tx := txCommand()
	for _, sub := range []string{"status", "wait"} {
		found, _, err := tx.Find([]string{sub})
		if err != nil || found == nil || found.Name() != sub {
			t.Errorf("bb tx %s does not resolve", sub)
			continue
		}
		if !found.DisableFlagParsing {
			t.Errorf("bb tx %s must pass flags through to bitbadges-cli", sub)
		}
	}
}

// sdkStandardSubcommands mirrors the subcommands each standards verb exposes
// in bitbadgesjs (bitbadges npm 0.45.8, `bitbadges-cli <verb> --help`). The
// `bb --help` one-liners for these verbs advertise example subcommands; every
// one they name must exist, or agents reading `bb --help` try verbs that fail
// with "unknown command" (QA session 17: 11 of 13 lines named a missing verb).
var sdkStandardSubcommands = map[string][]string{
	"auctions":           {"list", "show", "status", "place-bid", "cancel-bid", "accept-bid", "build"},
	"bounties":           {"list", "show", "status", "accept", "deny", "claim-refund", "build"},
	"credit-tokens":      {"list", "show", "quote", "purchase", "build"},
	"crowdfunds":         {"list", "show", "status", "contribute", "withdraw", "refund", "build"},
	"custom-2fa":         {"mint"},
	"dynamic-stores":     {"create", "update", "delete", "set-value", "add", "remove", "show", "get-value", "list-values", "batch", "by-creator", "search"},
	"intents":            {"list", "show", "create", "fill", "cancel", "build"},
	"nfts":               {"bid", "list", "cancel", "buy", "sell", "orders", "history"},
	"pay-requests":       {"list", "show", "status", "pay", "deny", "build", "build-v2"},
	"prediction-markets": {"list", "show", "status", "buy-yes", "buy-no", "sell-yes", "sell-no", "cancel", "deposit", "quote", "redeem", "resolve", "build"},
	"products":           {"list", "show", "purchase", "build"},
	"smart-tokens":       {"list", "show", "status", "deposit", "withdraw", "build"},
	"subscriptions":      {"list", "status", "claim", "enable-renewal", "cancel", "subscribe", "charge-due", "build"},
}

func TestStandardsHelpOnlyNamesRealSubcommands(t *testing.T) {
	byName := map[string]string{}
	for _, spec := range sdkForwarderSpecs {
		byName[spec.name] = spec.short
	}
	for verb, subs := range sdkStandardSubcommands {
		short, ok := byName[verb]
		if !ok {
			t.Errorf("no forwarder for standards verb %q", verb)
			continue
		}
		known := map[string]bool{}
		for _, s := range subs {
			known[s] = true
		}
		rest := short
		if i := strings.Index(rest, ":"); i >= 0 {
			rest = rest[i+1:]
		}
		if !strings.Contains(rest, " / ") {
			continue // prose description, no subcommand list
		}
		for _, tok := range strings.Split(rest, "/") {
			tok = strings.TrimSpace(tok)
			if tok == "" || tok == "..." {
				continue
			}
			if !known[tok] {
				t.Errorf("bb --help for %q advertises %q, which bitbadges-cli %s does not have", verb, tok, verb)
			}
		}
	}
}
