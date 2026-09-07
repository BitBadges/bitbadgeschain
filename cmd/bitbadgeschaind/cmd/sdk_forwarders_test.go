package cmd

import (
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
