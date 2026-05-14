package cmd

import (
	"github.com/spf13/cobra"
)

// Cobra GroupID values for `bb --help`. These visually separate the
// chain-native commands from the BitBadges SDK CLI surface (registered
// via thin forwarders) so `bb --help` answers "which binary owns this?"
// at a glance. See outputs/flagship-plans/cli-v2-design.md § "Locked
// Decisions" (decision 7) in the bitbadges-autopilot repo for context.
const (
	groupChain        = "chain"
	groupSDKBuild     = "sdk-build"
	groupSDKStandards = "sdk-standards"
	groupSDKIndexer   = "sdk-indexer"
	groupSDKSwap      = "sdk-swap"
	groupSDKDev       = "sdk-dev"
	groupLocal        = "local"
)

// registerHelpGroups installs the Cobra groups on the root command. Each
// top-level command is then tagged with the appropriate GroupID via
// tagChainCommandGroups (for chain natives) and the per-forwarder
// registration in sdk_forwarders.go (for SDK CLI commands).
func registerHelpGroups(rootCmd *cobra.Command) {
	rootCmd.AddGroup(
		&cobra.Group{ID: groupChain, Title: "Chain operations:"},
		&cobra.Group{ID: groupSDKBuild, Title: "BitBadges SDK — Build & Deploy:"},
		&cobra.Group{ID: groupSDKStandards, Title: "BitBadges SDK — Standards:"},
		&cobra.Group{ID: groupSDKIndexer, Title: "BitBadges SDK — Indexer & Auth:"},
		&cobra.Group{ID: groupSDKSwap, Title: "BitBadges SDK — Swap / DEX:"},
		&cobra.Group{ID: groupSDKDev, Title: "BitBadges SDK — Dev:"},
		&cobra.Group{ID: groupLocal, Title: "Local state:"},
	)
}

// chainNativeGroups maps every chain-native top-level command name to
// its group. Names match what Cobra actually registers (verified via
// `bitbadgeschaind --help`).
//
// Anything not in this map and not in sdkForwarderGroups (see
// sdk_forwarders.go) keeps Cobra's default ungrouped placement, which
// is fine for built-ins like `help`.
var chainNativeGroups = map[string]string{
	// Cosmos SDK core
	"start":          groupChain,
	"init":           groupChain,
	"status":         groupChain,
	"version":        groupChain,
	"tx":             groupChain,
	"query":          groupChain, // alias `q`
	"keys":           groupChain,
	"sign-arbitrary": groupChain,
	"comet":          groupChain, // aliases `cometbft`, `tendermint`
	"config":         groupChain, // confix-provided; chain owns this name
	"genesis":        groupChain,
	"debug":          groupChain,
	"prune":          groupChain,
	"snapshots":      groupChain,
	// Cosmos EVM additions
	"export":       groupChain,
	"rollback":     groupChain,
	"index-eth-tx": groupChain,
	// Existing back-compat forwarder. Keep grouped under chain so it
	// doesn't show up ungrouped during the deprecation runway.
	"cli": groupChain,
	// Cobra-provided completion command. Lives under "Local state" per
	// the v2 design decision 7 — it's a user-facing local helper, not a
	// chain native.
	"completion": groupLocal,
}

// tagChainCommandGroups assigns the GroupID on each chain-native
// top-level command. Must be called after every chain-native command
// is registered (i.e. after initRootCmd + autocli enhancement).
//
// Forces Cobra's lazy-registered defaults (`help`, `completion`) to be
// installed up front so they land in the right group too — otherwise
// they get auto-added at render time, after this pass, and show under
// "Additional Commands" instead of "Local state".
func tagChainCommandGroups(rootCmd *cobra.Command) {
	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultCompletionCmd()

	for _, cmd := range rootCmd.Commands() {
		if cmd.GroupID != "" {
			continue
		}
		if groupID, ok := chainNativeGroups[cmd.Name()]; ok {
			cmd.GroupID = groupID
		}
	}
}
