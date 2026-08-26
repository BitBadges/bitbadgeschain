package cmd

import (
	"testing"

	cmtcfg "github.com/cometbft/cometbft/config"
	cosmosevmserverconfig "github.com/cosmos/evm/server/config"
)

// TestInitCometBFTConfigAgreesWithTheEVMMempoolDefault pins that the two
// config files this binary generates can actually run together.
//
// `bitbadgeschaind init` writes both config.toml (from initCometBFTConfig) and
// app.toml (from initAppConfig). cosmos/evm v0.7 cross-validates them at
// startup: if the EVM mempool is enabled it demands config.toml's
// mempool.type = "app", and refuses to start otherwise.
//
// v34 enables the EVM mempool by default — cosmos/evm's own start command
// registers mempool.max-txs with a default of 0 and then force-sets it, so
// app.toml cannot turn it off (see server/start.go:188-190). But
// initCometBFTConfig returned CometBFT's stock default, whose mempool type is
// "flood". The result was that a freshly initialised v34 node could not start
// at all:
//
//	EVM mempool enabled, but comet-bft has invalid config.toml:mempool.type
//	(want 'app', got 'flood'): error in app.toml
//
// The upgrade path hid this, because there the pre-upgrade hook rewrites an
// existing config.toml. What it did NOT cover is every path that creates a
// *new* config: a validator adding a sentry, an operator rebuilding a node for
// state sync, anyone joining the network after v34, or a local dev chain.
// Those all run `init` and would have hit a node that refuses to boot with no
// documented remedy.
//
// This asserts the invariant rather than the literal string, by running the
// same cross-check the node runs at startup.
func TestInitCometBFTConfigAgreesWithTheEVMMempoolDefault(t *testing.T) {
	cmtConfig := initCometBFTConfig()

	_, rawAppConfig := initAppConfig()
	appConfig, ok := rawAppConfig.(EVMAppConfig)
	if !ok {
		t.Fatalf("initAppConfig returned %T, want EVMAppConfig", rawAppConfig)
	}

	evmConfig := &cosmosevmserverconfig.Config{
		Config:  appConfig.Config,
		EVM:     appConfig.EVM,
		JSONRPC: appConfig.JSONRPC,
		TLS:     appConfig.TLS,
	}

	// The pair is validated exactly as `init` writes it, with no runtime
	// modelling. cosmos/evm's start command force-sets mempool.max-txs to 0
	// (server/start.go:188-190), so the node would run with the EVM mempool
	// enabled even if app.toml disagreed - but that force-set is upstream's own
	// workaround ("explicitly override the app.toml default value"), and a
	// config pair that only works because of it is one upstream release away
	// from not working. The files must agree on their own.
	if err := cosmosevmserverconfig.ValidateCrossConfig(cmtConfig, evmConfig); err != nil {
		t.Fatalf("the config.toml and app.toml that `init` generates cannot run together: %v\n"+
			"config.toml mempool.type = %q. A node created with `bitbadgeschaind init` will "+
			"refuse to start, and unlike the upgrade path there is no pre-upgrade hook to fix it.",
			err, cmtConfig.Mempool.Type)
	}
}

// TestInitCometBFTConfigUsesTheAppMempool is the narrow, readable companion to
// the cross-check above: it names the value rather than deriving it, so a
// reader who breaks the test can see immediately what is expected.
func TestInitCometBFTConfigUsesTheAppMempool(t *testing.T) {
	if got := initCometBFTConfig().Mempool.Type; got != cmtcfg.MempoolTypeApp {
		t.Fatalf("initCometBFTConfig().Mempool.Type = %q, want %q", got, cmtcfg.MempoolTypeApp)
	}
}

// TestInitAppConfigEnablesTheEVMMempoolOnDisk pins app.toml's mempool.max-txs
// to 0 - the same "EVM mempool enabled" value the start command force-sets at
// runtime (cosmos/evm server/start.go:188-190).
//
// The SDK default this replaced was -1, which app.toml documents as "disable
// the mempool". That produced a freshly initialised node whose two config
// files contradicted each other as written: config.toml said mempool.type =
// "app" while app.toml said the app-side mempool is off, a pair
// ValidateCrossConfig rejects in its symmetric direction. The node still
// booted, but only because the start command's flag force-set shadows the
// file - app.toml's max-txs is dead configuration on the start path, and an
// operator reading or editing it was reading a lie. If upstream ever drops
// its force-set workaround, an init'ed node with -1 on disk stops booting.
func TestInitAppConfigEnablesTheEVMMempoolOnDisk(t *testing.T) {
	_, rawAppConfig := initAppConfig()
	appConfig, ok := rawAppConfig.(EVMAppConfig)
	if !ok {
		t.Fatalf("initAppConfig returned %T, want EVMAppConfig", rawAppConfig)
	}

	if got := appConfig.Mempool.MaxTxs; got != 0 {
		t.Fatalf("initAppConfig() writes mempool.max-txs = %d, want 0.\n"+
			"The runtime force-sets 0 whatever this file says, so any other value here is "+
			"a lie to the operator - and -1 makes the generated config.toml/app.toml pair "+
			"mutually inconsistent as written.", got)
	}
}
